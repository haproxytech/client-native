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

package tests

import (
	"testing"

	"github.com/haproxytech/client-native/v6/config-parser/parsers"
	"github.com/haproxytech/client-native/v6/models"
)

func TestFilterRuleGrammar(t *testing.T) {
	valid := []struct {
		name    string
		line    string
		data    string
		comment string
	}{
		{
			name: "enable single class",
			line: "filter-enable trace",
			data: "filter-enable trace",
		},
		{
			name: "disable single class",
			line: "filter-disable trace",
			data: "filter-disable trace",
		},
		{
			name: "enable class and instance",
			line: "filter-enable trace/A",
			data: "filter-enable trace/A",
		},
		{
			name: "several selectors keep their order",
			line: "filter-enable trace/B trace/A http-comp-req",
			data: "filter-enable trace/B trace/A http-comp-req",
		},
		{
			name: "repeated selector",
			line: "filter-disable trace/A trace/A",
			data: "filter-disable trace/A trace/A",
		},
		{
			name:    "metadata",
			line:    "filter-enable trace/A # {\"k\":\"v\"}",
			data:    "filter-enable trace/A",
			comment: `{"k":"v"}`,
		},
	}
	for _, tt := range valid {
		t.Run(tt.name, func(t *testing.T) {
			p := &parsers.FilterRule{}
			p.Init()
			if err := ProcessLine(tt.line, p); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			result, err := p.Result()
			if err != nil {
				t.Fatalf("result: %v", err)
			}
			if len(result) != 1 {
				t.Fatalf("result lines = %d, want 1", len(result))
			}
			if result[0].Data != tt.data {
				t.Fatalf("result = %q, want %q", result[0].Data, tt.data)
			}
			if result[0].Comment != tt.comment {
				t.Fatalf("comment = %q, want %q", result[0].Comment, tt.comment)
			}
		})
	}

	invalid := []string{
		"filter-enable",
		"filter-disable",
		"filter-rule trace/A",
		"filter-enable trace/",
		"filter-enable /A",
		"filter-enable trace/A/B",
		"filter-enable trace/A!",
	}
	for _, line := range invalid {
		t.Run(line, func(t *testing.T) {
			p := &parsers.FilterRule{}
			p.Init()
			if err := ProcessLine(line, p); err == nil {
				t.Fatalf("expected error for %q", line)
			}
			if _, err := p.Result(); err == nil {
				t.Fatalf("expected empty result for %q", line)
			}
		})
	}
}

func TestFilterRuleRetainsEnableDisableOrder(t *testing.T) {
	lines := []string{
		"filter-enable trace/A",
		"filter-disable trace",
		"filter-enable trace/B trace/A",
	}
	p := &parsers.FilterRule{}
	p.Init()
	for _, line := range lines {
		if err := ProcessLine(line, p); err != nil {
			t.Fatalf("parse %q: %v", line, err)
		}
	}

	data, err := p.Get(false)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	rules := data.([]models.FilterRule)
	if len(rules) != len(lines) {
		t.Fatalf("rules = %d, want %d", len(rules), len(lines))
	}
	wantTypes := []string{"enable", "disable", "enable"}
	for i, rule := range rules {
		if rule.Type != wantTypes[i] {
			t.Fatalf("rule %d type = %q, want %q", i, rule.Type, wantTypes[i])
		}
	}

	result, err := p.Result()
	if err != nil {
		t.Fatalf("result: %v", err)
	}
	for i, line := range lines {
		if result[i].Data != line {
			t.Fatalf("result %d = %q, want %q", i, result[i].Data, line)
		}
	}
}
