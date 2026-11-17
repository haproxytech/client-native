// Copyright 2026 HAProxy Technologies
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package configuration

import (
	"cmp"
	"fmt"
	"os"
	"slices"
	"strings"

	parser "github.com/haproxytech/client-native/v6/config-parser"
	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/configuration/options"
	"github.com/haproxytech/client-native/v6/runtime"
)

// step is one runtime call needed to make the running HAProxy match the
// committed configuration.
type step func(rt runtime.Runtime) error

// change is one differing plain attribute of a section.
type change struct {
	section  string            // section name
	old, new common.ParserData // nil when the attribute is not set
}

// fieldApplier turns a changed attribute into runtime steps; false means reload.
type fieldApplier func(c change) ([]step, bool)

// option is one keyword of an item line with its arguments, e.g. {"weight", "10"}.
type option struct{ name, args string }

// item is one named line of an object list (a server) with its effective options.
type item struct {
	name    string
	options []option
}

// get returns the arguments of the last option called name.
func (it item) get(name string) (string, bool) {
	for _, o := range slices.Backward(it.options) {
		if o.name == name {
			return o.args, true
		}
	}
	return "", false
}

// itemApplier turns one changed option of a kept item into runtime steps;
// false means reload.
type itemApplier func(section string, old, cur item) ([]step, bool)

// objectSpec describes a list of named lines inside a section.
type objectSpec struct {
	items  func(p parser.Parser, t parser.Section, section string) ([]item, error)
	fields map[string]itemApplier                       // option name -> applier; missing = reload
	add    func(section string, it item) ([]step, bool) // nil = reload
	del    func(section string, it item) ([]step, bool) // nil = reload
}

// entry is what the registry holds for one attribute: either apply or object.
type entry struct {
	apply  fieldApplier
	object *objectSpec
}

// runtimeRegistry maps section type and attribute (the configuration keyword)
// to the way a change of it is applied at runtime.
type runtimeRegistry map[parser.Section]map[string]entry

// scanOrder is the order section types are compared in: cheap and
// reload-prone types first, so a reload is detected early. Unlisted types go
// before all listed ones.
var scanOrder = []parser.Section{ //nolint:gochecknoglobals
	parser.Global, parser.HTTPHardening, parser.Defaults, parser.UserList, parser.Peers, parser.Mailers,
	parser.Resolvers, parser.Cache, parser.Ring, parser.Traces, parser.LogForward,
	parser.LogProfile, parser.HTTPErrors, parser.HealthChecks, parser.Acme,
	parser.CrtStore, parser.FCGIApp, parser.Frontends, parser.Listen, parser.Backends,
}

// prefilterMinSections: section types with fewer sections are compared
// directly; the prefilter costs a pass over both files.
const prefilterMinSections = 16

// planRuntimeApply decides how the commit of transactionID reaches the running
// HAProxy. It must run while the previous configuration file and parser are
// still in place. The returned function applies the plan and reports the
// outcome; it is nil when runtime apply is not enabled.
func (c *client) planRuntimeApply(transactionID, transactionFile string, version int64) func() {
	if c.RuntimeApply == nil && c.OnCommit == nil {
		return nil
	}
	res := options.CommitResult{TransactionID: transactionID, Version: version, Reload: true}
	var rt runtime.Runtime
	var steps []step
	func() {
		defer recoverInto(&res, "planning")
		if c.RuntimeApply != nil {
			rt = c.RuntimeApply()
		}
		if rt == nil {
			return
		}
		newP, err := c.GetParser(transactionID)
		if err != nil {
			res.Err = err
			return
		}
		same := func() map[string]struct{} { return sameSectionsInFiles(c.ConfigurationFile, transactionFile) }
		var ok bool
		steps, ok = runtimeSteps(registry, c.parser, newP, same)
		res.Reload = !ok
	}()
	return func() {
		if !res.Reload {
			func() {
				defer recoverInto(&res, "step")
				for _, s := range steps {
					if err := s(rt); err != nil {
						res.Reload, res.Err = true, err
						return
					}
				}
			}()
		}
		if c.OnCommit != nil {
			c.OnCommit(res)
		}
	}
}

// recoverInto turns a panic into a reload with the panic as the reason. It
// must be called directly by defer.
func recoverInto(res *options.CommitResult, stage string) {
	if r := recover(); r != nil {
		res.Reload, res.Err = true, fmt.Errorf("runtime apply %s panicked: %v", stage, r)
	}
}

// sameSectionsInFiles is sameSections on the contents of two files; nil when
// either cannot be read.
func sameSectionsInFiles(oldPath, newPath string) map[string]struct{} {
	oldText, err := os.ReadFile(oldPath)
	if err != nil {
		return nil
	}
	newText, err := os.ReadFile(newPath)
	if err != nil {
		return nil
	}
	return sameSections(string(oldText), string(newText))
}

// runtimeSteps returns the runtime calls that turn a HAProxy running oldP into
// one running newP. It returns false at the first difference reg has no entry
// for; the caller must then reload. same, when not nil, is called at most once,
// on reaching the first section type with more than prefilterMinSections
// sections; the sections whose "type name" key it returns are known to be
// textually identical and are not compared from then on.
func runtimeSteps(reg runtimeRegistry, oldP, newP parser.Parser, same func() map[string]struct{}) ([]step, bool) {
	types, ok := sectionTypes(oldP, newP)
	if !ok {
		return nil, false
	}
	var steps []step
	var skip map[string]struct{}
	for _, t := range types {
		names := sectionNames(oldP, t)
		if !slices.Equal(names, sectionNames(newP, t)) {
			return nil, false
		}
		if same != nil && len(names) > prefilterMinSections {
			skip, same = same(), nil
		}
		for _, name := range names {
			if len(skip) > 0 {
				if _, found := skip[string(t)+" "+name]; found {
					continue
				}
			}
			more, ok := sectionSteps(reg, oldP, newP, t, name)
			if !ok {
				return nil, false
			}
			steps = append(steps, more...)
		}
	}
	return steps, true
}

// sectionSteps returns the runtime calls for the differences of one section;
// false means reload.
func sectionSteps(reg runtimeRegistry, oldP, newP parser.Parser, t parser.Section, name string) ([]step, bool) {
	// fn runs under both parser locks: it must not call the parsers.
	var attrs []string
	missing := false
	ok := parser.DiffSection(oldP, newP, t, name, func(attr string) bool {
		if _, found := reg[t][attr]; !found {
			missing = true
			return false
		}
		attrs = append(attrs, attr)
		return true
	})
	if !ok || missing {
		return nil, false
	}
	var steps []step
	for _, attr := range attrs {
		var more []step
		e := reg[t][attr]
		switch {
		case e.object != nil:
			more, ok = diffObject(e.object, oldP, newP, t, name)
		case e.apply != nil:
			more, ok = e.apply(change{section: name, old: getData(oldP, t, name, attr), new: getData(newP, t, name, attr)})
		default:
			ok = false
		}
		if !ok {
			return nil, false
		}
		steps = append(steps, more...)
	}
	return steps, true
}

// sameSections returns the keys ("backend be1") of the sections whose text is
// byte-identical in two configurations rendered by the config-parser writer.
// It returns nil, meaning nothing may be skipped, when either text does not
// have the writer's layout.
func sameSections(oldText, newText string) map[string]struct{} {
	a, okA := sectionBlocks(oldText)
	b, okB := sectionBlocks(newText)
	if !okA || !okB {
		return nil
	}
	same := make(map[string]struct{}, len(a))
	for key, block := range a {
		if other, found := b[key]; found && other == block {
			same[key] = struct{}{}
		}
	}
	return same
}

// sectionBlocks splits a rendered configuration into its sections, keyed by
// "type name" ("global" when the header has no name). It returns false when
// the text does not have the layout the writer produces.
func sectionBlocks(text string) (map[string]string, bool) {
	blocks := map[string]string{}
	key, start := "", -1
	closeBlock := func(end int) bool {
		if start < 0 {
			return true
		}
		if _, dup := blocks[key]; dup {
			return false
		}
		blocks[key] = text[start:end]
		return true
	}
	for pos := 0; pos < len(text); {
		end := len(text)
		if i := strings.IndexByte(text[pos:], '\n'); i >= 0 {
			end = pos + i + 1
		}
		if c := text[pos]; c != '\n' && c != ' ' && c != '\t' && c != '#' {
			k, ok := headerKey(text[pos:end])
			if !ok || !closeBlock(pos) {
				return nil, false
			}
			key, start = k, pos
		}
		pos = end
	}
	if !closeBlock(len(text)) {
		return nil, false
	}
	return blocks, true
}

// headerKey returns the "type name" key of a section header line.
func headerKey(line string) (string, bool) {
	typ, rest, _ := strings.Cut(strings.TrimRight(line, "\n"), " ")
	if !slices.Contains(scanOrder, parser.Section(typ)) {
		return "", false
	}
	name, _, _ := strings.Cut(strings.TrimLeft(rest, " "), " ")
	if name == "" || name == "from" || strings.HasPrefix(name, "#") {
		return typ, true
	}
	return typ + " " + name, true
}

// sectionTypes returns the section types of both parsers in scan order.
func sectionTypes(oldP, newP parser.Parser) ([]parser.Section, bool) {
	a, b := parser.SectionTypes(oldP), parser.SectionTypes(newP)
	if a == nil || b == nil {
		return nil, false
	}
	types := slices.DeleteFunc(slices.Concat(a, b), func(t parser.Section) bool { return t == parser.Comments })
	slices.SortFunc(types, func(x, y parser.Section) int {
		return cmp.Or(cmp.Compare(slices.Index(scanOrder, x), slices.Index(scanOrder, y)), cmp.Compare(x, y))
	})
	return slices.Compact(types), true
}

func sectionNames(p parser.Parser, t parser.Section) []string {
	names, err := p.SectionsGet(t)
	if err != nil {
		return nil
	}
	slices.Sort(names)
	return names
}

func getData(p parser.Parser, t parser.Section, section, attr string) common.ParserData {
	data, err := p.Get(t, section, attr)
	if err != nil {
		return nil
	}
	return data
}

// diffObject returns the runtime steps that turn the object list of section
// in oldP into the one in newP: removals, then changes, then additions.
func diffObject(spec *objectSpec, oldP, newP parser.Parser, t parser.Section, section string) ([]step, bool) {
	olds, err := spec.items(oldP, t, section)
	if err != nil {
		return nil, false
	}
	curs, err := spec.items(newP, t, section)
	if err != nil {
		return nil, false
	}
	oldIdx, okOld := indexItems(olds)
	curIdx, okCur := indexItems(curs)
	if !okOld || !okCur {
		return nil, false
	}

	if !keptInOrder(curs, oldIdx) {
		return nil, false
	}

	var steps []step
	appendSteps := func(more []step, ok bool) bool {
		steps = append(steps, more...)
		return ok
	}
	for _, old := range olds {
		if _, kept := curIdx[old.name]; !kept && (spec.del == nil || !appendSteps(spec.del(section, old))) {
			return nil, false
		}
	}
	for _, cur := range curs {
		i, kept := oldIdx[cur.name]
		if !kept || slices.Equal(olds[i].options, cur.options) {
			continue
		}
		changed := changedOptions(olds[i], cur)
		if len(changed) == 0 {
			return nil, false
		}
		for _, name := range changed {
			apply, found := spec.fields[name]
			if !found || !appendSteps(apply(section, olds[i], cur)) {
				return nil, false
			}
		}
	}
	for _, cur := range curs {
		if _, kept := oldIdx[cur.name]; !kept && (spec.add == nil || !appendSteps(spec.add(section, cur))) {
			return nil, false
		}
	}
	return steps, true
}

// indexItems maps item names to their position; false on duplicate names.
func indexItems(items []item) (map[string]int, bool) {
	idx := make(map[string]int, len(items))
	for i, it := range items {
		if _, dup := idx[it.name]; dup {
			return nil, false
		}
		idx[it.name] = i
	}
	return idx, true
}

// keptInOrder reports whether the items of curs found in oldIdx keep their
// relative order and come before every added item.
func keptInOrder(curs []item, oldIdx map[string]int) bool {
	last, added := -1, false
	for _, cur := range curs {
		i, kept := oldIdx[cur.name]
		if !kept {
			added = true
			continue
		}
		if added || i < last {
			return false
		}
		last = i
	}
	return true
}

// changedOptions returns the option names whose effective value differs
// between old and cur, in order of first appearance in cur then old.
func changedOptions(old, cur item) []string {
	var names []string
	for _, o := range slices.Concat(cur.options, old.options) {
		if slices.Contains(names, o.name) {
			continue
		}
		a, okA := old.get(o.name)
		b, okB := cur.get(o.name)
		if okA != okB || a != b {
			names = append(names, o.name)
		}
	}
	return names
}
