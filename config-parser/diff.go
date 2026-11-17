/*
Copyright 2026 HAProxy Technologies

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package parser

import "github.com/haproxytech/client-native/v6/config-parser/common"

// DiffSection reports which attributes of section (t, name) render differently
// in a and b. fn is called once per differing attribute, in ParserSequence
// order, and stops the walk by returning false. A different `from` defaults
// section is reported first, as attribute "from". Comments are ignored.
// It returns false when a or b was not created by this package, the section is
// missing in either, or the two sections do not have the same parser sequence;
// callers must then assume the sections differ.
// fn runs while both parsers are locked and must not call any Parser method.
// DiffSection locks a then b, so concurrent calls with swapped arguments can
// deadlock.
func DiffSection(a, b Parser, t Section, name string, fn func(attribute string) bool) bool {
	pa, okA := a.(*configParser)
	pb, okB := b.(*configParser)
	if !okA || !okB {
		return false
	}
	pa.lock()
	defer pa.unLock()
	if pb != pa {
		pb.lock()
		defer pb.unLock()
	}

	sa, okA := pa.Parsers[t][name]
	sb, okB := pb.Parsers[t][name]
	if !okA || !okB {
		return false
	}
	if sa.SectionName != sb.SectionName && !fn("from") {
		return true
	}
	if len(sa.ParserSequence) != len(sb.ParserSequence) {
		return false
	}
	for i, attr := range sa.ParserSequence {
		if sb.ParserSequence[i] != attr {
			return false
		}
		la, lb := renderedLines(sa, string(attr)), renderedLines(sb, string(attr))
		if linesEqual(la, lb) {
			continue
		}
		if !fn(string(attr)) {
			return true
		}
	}
	return true
}

// SectionTypes returns the section types p holds at least one section of.
// It returns nil when p was not created by this package.
func SectionTypes(p Parser) []Section {
	cp, ok := p.(*configParser)
	if !ok {
		return nil
	}
	cp.lock()
	defer cp.unLock()
	types := make([]Section, 0, len(cp.Parsers))
	for t, sections := range cp.Parsers {
		if len(sections) > 0 {
			types = append(types, t)
		}
	}
	return types
}

func renderedLines(s *Parsers, attribute string) []common.ReturnResultLine {
	p, ok := s.Parsers[attribute]
	if !ok {
		return nil
	}
	lines, _, err := p.ResultAll()
	if err != nil {
		return nil
	}
	return lines
}

func linesEqual(a, b []common.ReturnResultLine) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].Data != b[i].Data {
			return false
		}
	}
	return true
}
