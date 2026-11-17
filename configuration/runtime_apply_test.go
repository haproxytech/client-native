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
	"fmt"
	"maps"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	parser "github.com/haproxytech/client-native/v6/config-parser"
	parser_options "github.com/haproxytech/client-native/v6/config-parser/options"
	"github.com/haproxytech/client-native/v6/config-parser/types"
	"github.com/haproxytech/client-native/v6/runtime"
)

const runtimeApplyBase = `global
  maxconn 100

defaults d
  mode http

frontend fe from d
  maxconn 50
  bind :80
  default_backend be

backend be from d
  balance roundrobin
  # comment
  server s1 10.0.0.1:80 weight 10 check
  server s2 10.0.0.2:80 weight 10
`

func newRuntimeApplyParser(t testing.TB, cfg string) parser.Parser {
	t.Helper()
	p, err := parser.New(parser_options.String(cfg))
	require.NoError(t, err)
	return p
}

// testServerItems turns the server lines of a section into items: the
// address first, then one option per server parameter.
func testServerItems(p parser.Parser, t parser.Section, section string) ([]item, error) {
	data, err := p.Get(t, section, "server")
	if err != nil {
		return nil, err
	}
	servers, _ := data.([]types.Server)
	items := make([]item, 0, len(servers))
	for _, srv := range servers {
		it := item{name: srv.Name, options: []option{{"address", srv.Address}}}
		for _, param := range srv.Params {
			name, args, _ := strings.Cut(param.String(), " ")
			it.options = append(it.options, option{name, args})
		}
		items = append(items, it)
	}
	return items, nil
}

// stepLog records the labels of the steps that were run.
type stepLog []string

func (l *stepLog) step(label string) step {
	return func(runtime.Runtime) error {
		*l = append(*l, label)
		return nil
	}
}

func (l *stepLog) itemField(name string) itemApplier {
	return func(_ string, _, cur item) ([]step, bool) {
		return []step{l.step(name + " " + cur.name)}, true
	}
}

func (l *stepLog) itemOp(op string) func(string, item) ([]step, bool) {
	return func(_ string, it item) ([]step, bool) {
		return []step{l.step(op + " " + it.name)}, true
	}
}

func testRuntimeRegistry(t *testing.T, l *stepLog) runtimeRegistry {
	t.Helper()
	return runtimeRegistry{
		parser.Frontends: {
			"maxconn": {apply: func(c change) ([]step, bool) {
				if c.old == nil || c.new == nil {
					t.Errorf("maxconn change has nil data: %+v", c)
				}
				return []step{l.step("maxconn " + c.section)}, true
			}},
		},
		parser.Backends: {
			"server": {object: &objectSpec{
				items:  testServerItems,
				fields: map[string]itemApplier{"weight": l.itemField("weight"), "check": l.itemField("check")},
				add:    l.itemOp("add"),
				del:    l.itemOp("del"),
			}},
		},
	}
}

func runSteps(t *testing.T, steps []step) {
	t.Helper()
	for _, s := range steps {
		require.NoError(t, s(nil))
	}
}

func TestRuntimeStepsPackageRegistry(t *testing.T) {
	r := strings.Replace
	for name, tc := range map[string]struct {
		cfg   string
		ok    bool
		steps int
	}{
		"identical":      {runtimeApplyBase, true, 0},
		"comment only":   {r(runtimeApplyBase, "# comment", "# other comment", 1), true, 0},
		"weight changed": {r(runtimeApplyBase, "weight 10 check", "weight 20 check", 1), true, 1},
		"no entry":       {r(runtimeApplyBase, "balance roundrobin", "balance leastconn", 1), false, 0},
	} {
		t.Run(name, func(t *testing.T) {
			steps, ok := runtimeSteps(registry, newRuntimeApplyParser(t, runtimeApplyBase), newRuntimeApplyParser(t, tc.cfg), nil)
			require.Equal(t, tc.ok, ok)
			require.Len(t, steps, tc.steps)
		})
	}
}

func TestRuntimeSteps(t *testing.T) {
	r := strings.Replace
	tests := []struct {
		name   string
		cfg    string
		mutate func(*objectSpec)
		wantOK bool
		want   []string
	}{
		{name: "identical", cfg: runtimeApplyBase, wantOK: true},
		{name: "no registry entry", cfg: r(runtimeApplyBase, "balance roundrobin", "balance leastconn", 1)},
		{name: "backend added", cfg: runtimeApplyBase + "\nbackend be2\n  balance roundrobin\n"},
		{name: "backend removed", cfg: r(runtimeApplyBase, "backend be from d", "backend other from d", 1)},
		{name: "frontend maxconn", cfg: r(runtimeApplyBase, "maxconn 50", "maxconn 60", 1), wantOK: true, want: []string{"maxconn fe"}},
		{name: "server added", cfg: runtimeApplyBase + "  server s3 10.0.0.3:80\n", wantOK: true, want: []string{"add s3"}},
		{name: "server removed", cfg: r(runtimeApplyBase, "  server s2 10.0.0.2:80 weight 10\n", "", 1), wantOK: true, want: []string{"del s2"}},
		{name: "weight changed", cfg: r(runtimeApplyBase, "weight 10 check", "weight 20 check", 1), wantOK: true, want: []string{"weight s1"}},
		{name: "address changed", cfg: r(runtimeApplyBase, "10.0.0.1:80", "10.0.0.9:80", 1)},
		{name: "two options changed", cfg: r(runtimeApplyBase, "weight 10 check", "weight 20", 1), wantOK: true, want: []string{"weight s1", "check s1"}},
		{name: "servers reordered", cfg: r(runtimeApplyBase, "  server s1 10.0.0.1:80 weight 10 check\n  server s2 10.0.0.2:80 weight 10\n", "  server s2 10.0.0.2:80 weight 10\n  server s1 10.0.0.1:80 weight 10 check\n", 1)},
		{name: "server inserted before kept", cfg: r(runtimeApplyBase, "  server s1", "  server s0 10.0.0.5:80\n  server s1", 1)},
		{name: "duplicate server names", cfg: runtimeApplyBase + "  server s2 10.0.0.2:80 weight 10\n"},
		{name: "nil add", cfg: runtimeApplyBase + "  server s3 10.0.0.3:80\n", mutate: func(s *objectSpec) { s.add = nil }},
		{name: "nil del", cfg: r(runtimeApplyBase, "  server s2 10.0.0.2:80 weight 10\n", "", 1), mutate: func(s *objectSpec) { s.del = nil }},
		{name: "options reordered", cfg: r(runtimeApplyBase, "weight 10 check", "check weight 10", 1)},
		{
			name:   "mixed change order",
			cfg:    r(r(runtimeApplyBase, "  server s2 10.0.0.2:80 weight 10\n", "  server s3 10.0.0.3:80\n", 1), "weight 10 check", "weight 20 check", 1),
			wantOK: true, want: []string{"del s2", "weight s1", "add s3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var l stepLog
			reg := testRuntimeRegistry(t, &l)
			if tt.mutate != nil {
				tt.mutate(reg[parser.Backends]["server"].object)
			}
			steps, ok := runtimeSteps(reg, newRuntimeApplyParser(t, runtimeApplyBase), newRuntimeApplyParser(t, tt.cfg), nil)
			require.Equal(t, tt.wantOK, ok)
			if !ok {
				require.Nil(t, steps)
				return
			}
			runSteps(t, steps)
			require.Equal(t, tt.want, []string(l))
		})
	}
}

func TestRuntimeStepsGlobalBeforeBackends(t *testing.T) {
	reg := runtimeRegistry{
		parser.Backends: {"server": {object: &objectSpec{
			items: func(parser.Parser, parser.Section, string) ([]item, error) {
				t.Fatal("backend items must not be read")
				return nil, nil
			},
		}}},
	}
	cfg := strings.NewReplacer("maxconn 100", "maxconn 200", "weight 10 check", "weight 20 check").Replace(runtimeApplyBase)
	_, ok := runtimeSteps(reg, newRuntimeApplyParser(t, runtimeApplyBase), newRuntimeApplyParser(t, cfg), nil)
	require.False(t, ok)
}

func TestDiffObjectItemGet(t *testing.T) {
	it := item{name: "s1", options: []option{{"weight", "1"}, {"check", ""}, {"weight", "2"}}}
	args, ok := it.get("weight")
	require.True(t, ok)
	require.Equal(t, "2", args)
	_, ok = it.get("maxconn")
	require.False(t, ok)
}

const sectionBlocksCfg = `# top comment
global
  maxconn 100

http-hardening hh
  all on

defaults d1
  mode http

defaults d2 from d1
  mode tcp

frontend fe from d1
  bind :80

backend be1 from d1
  # c1
  server s1 10.0.0.1:80

# note above
backend be2 # note
  ###_config-snippet_### BEGIN
  tune.ssl.default-dh-param 2048
  tune.bufsize 32768
  ###_config-snippet_### END
  server s1 10.0.0.2:80
`

var sectionBlocksKeys = []string{"global", "http-hardening hh", "defaults d1", "defaults d2", "frontend fe", "backend be1", "backend be2"} //nolint:gochecknoglobals

func TestSectionBlocks(t *testing.T) {
	text := newRuntimeApplyParser(t, sectionBlocksCfg).String()
	blocks, ok := sectionBlocks(text)
	require.True(t, ok)
	require.ElementsMatch(t, sectionBlocksKeys, slices.Collect(maps.Keys(blocks)))
	headers := map[string]string{
		"global": "global\n", "http-hardening hh": "http-hardening hh\n", "defaults d1": "defaults d1\n", "defaults d2": "defaults d2 from d1\n",
		"frontend fe": "frontend fe from d1\n", "backend be1": "backend be1 from d1\n", "backend be2": "backend be2 from d2 # note\n",
	}
	for key, block := range blocks {
		require.True(t, strings.HasPrefix(block, headers[key]), "block %q: %q", key, block)
		require.Contains(t, text, block)
	}
	require.Contains(t, blocks["backend be2"], "  ###_config-snippet_### BEGIN\n  tune.ssl.default-dh-param 2048\n  tune.bufsize 32768\n  ###_config-snippet_### END\n")
	require.Contains(t, blocks["backend be1"], "# note above\n")
}

func TestSectionBlocksInvalid(t *testing.T) {
	for name, text := range map[string]string{
		"unknown token":    "global\n  maxconn 1\n\nfoo bar\n\nbackend b\n  server s1 10.0.0.1:80\n",
		"duplicate header": "global\n  maxconn 1\n\nbackend b\n  server s1 10.0.0.1:80\n\nbackend b # again\n  server s1 10.0.0.1:80\n",
	} {
		t.Run(name, func(t *testing.T) {
			blocks, ok := sectionBlocks(text)
			require.False(t, ok)
			require.Nil(t, blocks)
		})
	}
}

func TestSameSections(t *testing.T) {
	without := func(key string) []string {
		return slices.DeleteFunc(slices.Clone(sectionBlocksKeys), func(k string) bool { return k == key })
	}
	tests := []struct {
		name string
		cfg  string
		want []string
	}{
		{name: "identical", cfg: sectionBlocksCfg, want: sectionBlocksKeys},
		{name: "server changed", cfg: strings.Replace(sectionBlocksCfg, "10.0.0.1:80", "10.0.0.9:80", 1), want: without("backend be1")},
		{name: "from changed", cfg: strings.Replace(sectionBlocksCfg, "backend be1 from d1", "backend be1 from d2", 1), want: without("backend be1")},
		{name: "comment changed", cfg: strings.Replace(sectionBlocksCfg, "# c1", "# c2", 1), want: without("backend be1")},
	}
	oldText := newRuntimeApplyParser(t, sectionBlocksCfg).String()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			same := sameSections(oldText, newRuntimeApplyParser(t, tt.cfg).String())
			require.ElementsMatch(t, tt.want, slices.Collect(maps.Keys(same)))
		})
	}

	bad := strings.Replace(oldText, "\nfrontend fe", "\nfoo bar\nfrontend fe", 1)
	require.Nil(t, sameSections(bad, oldText))
	require.Nil(t, sameSections(oldText, bad))
}

// prefilterBase has more than prefilterMinSections backends, two defaults, a
// frontend with maxconn and a comment in backend be1.
func prefilterBase(backends int) string {
	return strings.NewReplacer(
		"frontend fe from d\n", "defaults d2 from d\n  mode http\n\nfrontend fe from d\n  maxconn 50\n",
		"backend be1 from d\n  balance roundrobin\n", "backend be1 from d\n  balance roundrobin\n  # comment\n",
	).Replace(genConfig(backends, 2))
}

// failSame is a prefilter that fails the test when it is called.
func failSame(t *testing.T) func() map[string]struct{} {
	t.Helper()
	return func() map[string]struct{} {
		t.Error("prefilter must not be called")
		return nil
	}
}

func TestRuntimeStepsPrefilterEquivalence(t *testing.T) {
	base := prefilterBase(20)
	r := strings.Replace
	be1s1 := "  server s1 10.0.1.1:8080 check weight 10\n"
	tests := []struct {
		name, cfg string
		called    int
	}{
		{name: "no change", cfg: base, called: 1},
		{name: "comment only", cfg: r(base, "# comment", "# other comment", 1), called: 1},
		{name: "server changed", cfg: r(base, "10.0.1.0:8080 check weight 10", "10.0.1.0:8080 check weight 20", 1), called: 1},
		{name: "server no entry", cfg: r(base, "10.0.1.0:8080", "10.0.1.9:8080", 1), called: 1},
		{name: "server added", cfg: r(base, be1s1, be1s1+"  server s2 10.0.1.2:8080\n", 1), called: 1},
		{name: "server removed", cfg: r(base, be1s1, "", 1), called: 1},
		{name: "backend added", cfg: base + "\nbackend zz\n  balance roundrobin\n"},
		{name: "backend removed", cfg: r(base, "backend be19 from d\n", "backend zz from d\n", 1)},
		{name: "frontend maxconn", cfg: r(base, "maxconn 50", "maxconn 60", 1), called: 1},
		{name: "global changed", cfg: r(base, "maxconn 1000", "maxconn 2000", 1)},
		{name: "from changed", cfg: r(base, "backend be2 from d\n", "backend be2 from d2\n", 1), called: 1},
		{
			name:   "two backends change",
			cfg:    strings.NewReplacer("10.0.1.0:8080 check weight 10", "10.0.1.0:8080 check weight 20", "10.0.5.1:8080 check weight 10", "10.0.5.1:8080 check weight 30").Replace(base),
			called: 1,
		},
	}
	oldP := newRuntimeApplyParser(t, base)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newP := newRuntimeApplyParser(t, tt.cfg)
			run := func(same func() map[string]struct{}) (stepLog, bool) {
				var l stepLog
				steps, ok := runtimeSteps(testRuntimeRegistry(t, &l), oldP, newP, same)
				runSteps(t, steps)
				return l, ok
			}
			wantLog, wantOK := run(nil)
			called := 0
			gotLog, gotOK := run(func() map[string]struct{} {
				called++
				same := sameSections(oldP.String(), newP.String())
				require.NotNil(t, same)
				return same
			})
			require.Equal(t, tt.called, called)
			require.Equal(t, wantOK, gotOK)
			require.Equal(t, wantLog, gotLog)
		})
	}
}

func TestRuntimeStepsPrefilterNotCalledOnEarlyExit(t *testing.T) {
	base := prefilterBase(20)
	cfg := strings.Replace(base, "maxconn 1000", "maxconn 2000", 1)
	var l stepLog
	_, ok := runtimeSteps(testRuntimeRegistry(t, &l), newRuntimeApplyParser(t, base), newRuntimeApplyParser(t, cfg), failSame(t))
	require.False(t, ok)
}

func TestRuntimeStepsPrefilterNotCalledForSmallTypes(t *testing.T) {
	base := prefilterBase(prefilterMinSections)
	cfg := strings.Replace(base, "10.0.1.0:8080 check weight 10", "10.0.1.0:8080 check weight 20", 1)
	run := func(same func() map[string]struct{}) (stepLog, bool) {
		var l stepLog
		steps, ok := runtimeSteps(testRuntimeRegistry(t, &l), newRuntimeApplyParser(t, base), newRuntimeApplyParser(t, cfg), same)
		runSteps(t, steps)
		return l, ok
	}
	wantLog, wantOK := run(nil)
	gotLog, gotOK := run(failSame(t))
	require.True(t, wantOK)
	require.Equal(t, stepLog{"weight s0"}, wantLog)
	require.Equal(t, wantOK, gotOK)
	require.Equal(t, wantLog, gotLog)
}

func TestRuntimeStepsPrefilterCalledOnce(t *testing.T) {
	var sb strings.Builder
	sb.WriteString(prefilterBase(20))
	for i := range 19 {
		fmt.Fprintf(&sb, "\nfrontend fe%d from d\n  bind :%d\n", i, 8000+i)
	}
	base := sb.String()
	cfg := strings.Replace(base, "10.0.19.1:8080 check weight 10", "10.0.19.1:8080 check weight 20", 1)
	oldP, newP := newRuntimeApplyParser(t, base), newRuntimeApplyParser(t, cfg)
	called := 0
	var l stepLog
	steps, ok := runtimeSteps(testRuntimeRegistry(t, &l), oldP, newP, func() map[string]struct{} {
		called++
		return sameSections(oldP.String(), newP.String())
	})
	require.True(t, ok)
	runSteps(t, steps)
	require.Equal(t, 1, called)
	require.Equal(t, stepLog{"weight s1"}, l)
}

func TestRuntimeStepsSameIsTrusted(t *testing.T) {
	reg := runtimeRegistry{
		parser.Backends: {"server": {object: &objectSpec{
			items: func(parser.Parser, parser.Section, string) ([]item, error) {
				t.Fatal("backend items must not be read")
				return nil, nil
			},
		}}},
	}
	base := genConfig(prefilterMinSections+1, 2)
	cfg := strings.Replace(base, "10.0.1.0:8080 check weight 10", "10.0.1.0:8080 check weight 20", 1)
	same := func() map[string]struct{} { return map[string]struct{}{"backend be1": {}} }
	steps, ok := runtimeSteps(reg, newRuntimeApplyParser(t, base), newRuntimeApplyParser(t, cfg), same)
	require.True(t, ok)
	require.Empty(t, steps)
}
