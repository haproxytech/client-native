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

func TestFilterSequenceGrammar(t *testing.T) {
	valid := []struct {
		name    string
		line    string
		data    string
		comment string
	}{
		{
			name: "global",
			line: "filter-sequence global trace/A:before(trace/B)",
			data: "filter-sequence global trace/A:before(trace/B)",
		},
		{
			name: "request with several rules",
			line: "filter-sequence request trace/A:before(trace/B) comp-req:after(trace/C)",
			data: "filter-sequence request trace/A:before(trace/B) comp-req:after(trace/C)",
		},
		{
			name: "response",
			line: "filter-sequence response trace:after(trace/A)",
			data: "filter-sequence response trace:after(trace/A)",
		},
		{
			name:    "metadata",
			line:    "filter-sequence request trace/A:before(trace/B) # {\"k\":\"v\"}",
			data:    "filter-sequence request trace/A:before(trace/B)",
			comment: `{"k":"v"}`,
		},
	}
	for _, tt := range valid {
		t.Run(tt.name, func(t *testing.T) {
			p := &parsers.FilterSequence{}
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
		"filter-sequence",
		"filter-sequence request",
		"filter-sequence request trace/A",
		"filter-sequence request trace/A:before",
		"filter-sequence request trace/A:before(trace/B",
		"filter-sequence request trace/A:beside(trace/B)",
		"filter-sequence request trace/A before trace/B",
		"filter-sequence request trace/A:before(trace/B) trace/C",
		"filter-sequence bogus trace/A:before(trace/B)",
		"filter-sequence request trace/A:before(trace/A)",
	}
	for _, line := range invalid {
		t.Run(line, func(t *testing.T) {
			p := &parsers.FilterSequence{}
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

func sequenceSelector(class, id string) *models.FilterSelector {
	return &models.FilterSelector{Class: class, ID: id}
}

func parseSequences(t *testing.T, lines ...string) *parsers.FilterSequence {
	t.Helper()
	p := &parsers.FilterSequence{}
	p.Init()
	for _, line := range lines {
		if err := ProcessLine(line, p); err != nil {
			t.Fatalf("parse %q: %v", line, err)
		}
	}
	return p
}

func TestFilterSequenceCycles(t *testing.T) {
	t.Run("self cycle", func(t *testing.T) {
		p := &parsers.FilterSequence{}
		p.Init()
		if err := ProcessLine("filter-sequence request trace/A:before(trace/A)", p); err == nil {
			t.Fatal("expected self cycle error")
		}
	})

	t.Run("two node cycle", func(t *testing.T) {
		p := &parsers.FilterSequence{}
		p.Init()
		line := "filter-sequence request trace/A:before(trace/B) trace/B:before(trace/A)"
		if err := ProcessLine(line, p); err == nil {
			t.Fatal("expected two node cycle error")
		}
	})

	t.Run("three node cycle", func(t *testing.T) {
		p := &parsers.FilterSequence{}
		p.Init()
		line := "filter-sequence request trace/A:before(trace/B) trace/B:before(trace/C) trace/C:before(trace/A)"
		if err := ProcessLine(line, p); err == nil {
			t.Fatal("expected three node cycle error")
		}
	})

	t.Run("cross-line cycles are left to HAProxy", func(t *testing.T) {
		p := parseSequences(t,
			"filter-sequence request trace/A:before(trace/B)",
			"filter-sequence request trace/B:before(trace/A)",
		)
		data, err := p.Get(false)
		if err != nil {
			t.Fatalf("get: %v", err)
		}
		sequences := data.([]models.FilterSequence)
		if len(sequences) != 2 {
			t.Fatalf("sequences = %d, want both lines preserved", len(sequences))
		}
	})

	t.Run("after direction is normalized", func(t *testing.T) {
		p := &parsers.FilterSequence{}
		p.Init()
		line := "filter-sequence request trace/A:after(trace/B) trace/B:after(trace/A)"
		if err := ProcessLine(line, p); err == nil {
			t.Fatal("expected cycle error after normalization")
		}
	})

	t.Run("duplicate edges are harmless", func(t *testing.T) {
		parseSequences(t,
			"filter-sequence request trace/A:before(trace/B)",
			"filter-sequence request trace/A:before(trace/B)",
		)
	})

	t.Run("sides are independent", func(t *testing.T) {
		parseSequences(t,
			"filter-sequence request trace/A:before(trace/B)",
			"filter-sequence response trace/B:before(trace/A)",
		)
	})

	t.Run("cycle is detected per side", func(t *testing.T) {
		p := parseSequences(t, "filter-sequence request trace/A:before(trace/B)")
		line := "filter-sequence response trace/A:before(trace/B) trace/B:before(trace/A)"
		if err := ProcessLine(line, p); err == nil {
			t.Fatal("expected cycle error on the response side")
		}
	})

	t.Run("dangling references are allowed", func(t *testing.T) {
		parseSequences(t,
			"filter-sequence request trace/A:before(trace/B)",
			"filter-sequence request trace/C:before(trace/D)",
		)
	})

	t.Run("failed set keeps the previous collection", func(t *testing.T) {
		p := parseSequences(t, "filter-sequence request trace/A:before(trace/B)")
		cyclic := []models.FilterSequence{
			{Side: "request", Rules: []*models.FilterSequenceRule{
				{Instance: sequenceSelector("trace", "A"), Direction: "before", Target: sequenceSelector("trace", "B")},
				{Instance: sequenceSelector("trace", "B"), Direction: "before", Target: sequenceSelector("trace", "A")},
			}},
		}
		if err := p.Set(cyclic, -1); err == nil {
			t.Fatal("expected set cycle error")
		}
		result, err := p.Result()
		if err != nil {
			t.Fatalf("result: %v", err)
		}
		if len(result) != 1 || result[0].Data != "filter-sequence request trace/A:before(trace/B)" {
			t.Fatalf("previous collection changed: %v", result)
		}
	})

	t.Run("failed indexed set keeps the previous collection", func(t *testing.T) {
		p := parseSequences(t, "filter-sequence request trace/A:before(trace/B)")
		cyclic := models.FilterSequence{
			Side: "request",
			Rules: []*models.FilterSequenceRule{
				{Instance: sequenceSelector("trace", "A"), Direction: "before", Target: sequenceSelector("trace", "B")},
				{Instance: sequenceSelector("trace", "B"), Direction: "before", Target: sequenceSelector("trace", "A")},
			},
		}
		if err := p.Set(cyclic, 0); err == nil {
			t.Fatal("expected indexed set cycle error")
		}
		result, err := p.Result()
		if err != nil {
			t.Fatalf("result: %v", err)
		}
		if len(result) != 1 || result[0].Data != "filter-sequence request trace/A:before(trace/B)" {
			t.Fatalf("previous collection changed: %v", result)
		}
	})

	t.Run("failed insert keeps the previous collection", func(t *testing.T) {
		p := parseSequences(t, "filter-sequence request trace/A:before(trace/B)")
		cyclic := models.FilterSequence{
			Side: "request",
			Rules: []*models.FilterSequenceRule{
				{Instance: sequenceSelector("trace", "A"), Direction: "before", Target: sequenceSelector("trace", "B")},
				{Instance: sequenceSelector("trace", "B"), Direction: "before", Target: sequenceSelector("trace", "A")},
			},
		}
		if err := p.Insert(cyclic, 1); err == nil {
			t.Fatal("expected insert cycle error")
		}
		result, err := p.Result()
		if err != nil {
			t.Fatalf("result: %v", err)
		}
		if len(result) != 1 || result[0].Data != "filter-sequence request trace/A:before(trace/B)" {
			t.Fatalf("previous collection changed: %v", result)
		}
	})

	t.Run("delete and clear", func(t *testing.T) {
		p := parseSequences(t,
			"filter-sequence request trace/A:before(trace/B)",
			"filter-sequence response trace/C:before(trace/D)",
		)
		if err := p.Delete(0); err != nil {
			t.Fatalf("delete: %v", err)
		}
		result, err := p.Result()
		if err != nil {
			t.Fatalf("result: %v", err)
		}
		if len(result) != 1 || result[0].Data != "filter-sequence response trace/C:before(trace/D)" {
			t.Fatalf("unexpected result after delete: %v", result)
		}
		if err := p.Set(nil, -1); err != nil {
			t.Fatalf("clear: %v", err)
		}
		if _, err := p.Result(); err == nil {
			t.Fatal("expected fetch error after clear")
		}
	})
}
