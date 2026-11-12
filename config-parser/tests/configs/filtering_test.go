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
package configs //nolint:testpackage

import (
	"strings"
	"testing"

	parser "github.com/haproxytech/client-native/v6/config-parser"
	"github.com/haproxytech/client-native/v6/config-parser/options"
	"github.com/haproxytech/client-native/v6/config-parser/types"
	"github.com/haproxytech/client-native/v6/models"
)

const filterDirectiveConfig = `
global
  filter-api default

defaults named
  mode http
  filter-config trace id D
  filter-enable trace/D
  filter-disable trace
  filter-sequence request trace/D:before(trace/E)

frontend fe from named
  mode http
  filter-config http-comp-req enabled
  filter-enable trace/A # {"k":"v"}
  filter-disable trace
  filter-enable trace/B trace/A
  filter-sequence request trace/A:before(trace/B)

backend be from named
  mode http
  filter-config cache-store id C enabled
  filter-disable cache-store/C

listen li from named
  mode http
  filter-config fcgi id F
  filter-enable fcgi/F
`

func TestFilterDirectiveRoundTrip(t *testing.T) {
	p, err := parser.New(
		options.String(filterDirectiveConfig),
		options.UseListenSectionParsers,
	)
	if err != nil {
		t.Fatalf("parser.New: %v", err)
	}

	if got := p.String(); got != filterDirectiveConfig {
		t.Fatalf("round trip mismatch:\n--- got ---\n%s\n--- want ---\n%s", got, filterDirectiveConfig)
	}

	assertFilterRules(t, p, parser.Defaults, "named", []string{"enable", "disable"}, []string{"trace/D", "trace"})
	assertFilterRules(t, p, parser.Frontends, "fe", []string{"enable", "disable", "enable"}, []string{"trace/A", "trace", "trace/B trace/A"})
	assertFilterRules(t, p, parser.Backends, "be", []string{"disable"}, []string{"cache-store/C"})
	assertFilterRules(t, p, parser.Listen, "li", []string{"enable"}, []string{"fcgi/F"})

	filterAPI, err := p.Get(parser.Global, parser.GlobalSectionName, "filter-api")
	if err != nil {
		t.Fatalf("filter-api: %v", err)
	}
	value := filterAPI.(*types.StringC)
	if value.Value != "default" {
		t.Fatalf("filter-api = %q, want default", value.Value)
	}

	reparsed, err := parser.New(
		options.String(p.String()),
		options.UseListenSectionParsers,
	)
	if err != nil {
		t.Fatalf("reparse: %v", err)
	}
	if got := reparsed.String(); got != filterDirectiveConfig {
		t.Fatalf("second round trip mismatch:\n--- got ---\n%s", got)
	}
	assertFilterRules(t, reparsed, parser.Frontends, "fe", []string{"enable", "disable", "enable"}, []string{"trace/A", "trace", "trace/B trace/A"})

	data, err := reparsed.Get(parser.Frontends, "fe", "filter-rule")
	if err != nil {
		t.Fatalf("filter-rule: %v", err)
	}
	rules := data.([]models.FilterRule)
	if got := rules[0].Metadata["k"]; got != "v" {
		t.Fatalf("metadata k = %v, want v", got)
	}
}

func filterRuleSelectorText(rule models.FilterRule) string {
	parts := make([]string, 0, len(rule.Filters))
	for _, selector := range rule.Filters {
		value := selector.Class
		if selector.ID != "" {
			value += "/" + selector.ID
		}
		parts = append(parts, value)
	}
	return strings.Join(parts, " ")
}

func assertFilterRules(t *testing.T, p parser.Parser, section parser.Section, name string, wantTypes, wantSelectors []string) {
	t.Helper()
	data, err := p.Get(section, name, "filter-rule")
	if err != nil {
		t.Fatalf("%s/%s filter-rule: %v", section, name, err)
	}
	rules := data.([]models.FilterRule)
	if len(rules) != len(wantTypes) {
		t.Fatalf("%s/%s rules = %d, want %d", section, name, len(rules), len(wantTypes))
	}
	for i, rule := range rules {
		if rule.Type != wantTypes[i] {
			t.Fatalf("%s/%s rule %d type = %q, want %q", section, name, i, rule.Type, wantTypes[i])
		}
		if got := filterRuleSelectorText(rule); got != wantSelectors[i] {
			t.Fatalf("%s/%s rule %d filters = %q, want %q", section, name, i, got, wantSelectors[i])
		}
	}
}
