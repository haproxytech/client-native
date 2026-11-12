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
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

func TestFilterAPIValues(t *testing.T) {
	t.Run("absent", func(t *testing.T) {
		p := &parsers.FilterAPI{}
		p.Init()
		if _, err := p.Get(false); err == nil {
			t.Fatal("expected fetch error for absent filter-api")
		}
	})

	values := []string{"default", "downgrade-to-legacy"}
	for _, value := range values {
		t.Run(value, func(t *testing.T) {
			p := &parsers.FilterAPI{}
			p.Init()
			if err := ProcessLine("filter-api "+value, p); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			result, err := p.Result()
			if err != nil {
				t.Fatalf("result: %v", err)
			}
			if len(result) != 1 || result[0].Data != "filter-api "+value {
				t.Fatalf("result = %v, want filter-api %s", result, value)
			}
			data, err := p.Get(false)
			if err != nil {
				t.Fatalf("get: %v", err)
			}
			stringOption, ok := data.(*types.StringC)
			if !ok {
				t.Fatalf("typed value = %T, want *types.StringC", data)
			}
			if stringOption.Value != value {
				t.Fatalf("value = %q, want %q", stringOption.Value, value)
			}
		})
	}

	invalid := []string{
		"filter-api",
		"filter-api bogus",
		"filter-api default extra",
	}
	for _, line := range invalid {
		t.Run(line, func(t *testing.T) {
			p := &parsers.FilterAPI{}
			p.Init()
			if err := ProcessLine(line, p); err == nil {
				t.Fatalf("expected error for %q", line)
			}
		})
	}

	t.Run("replace and remove", func(t *testing.T) {
		p := &parsers.FilterAPI{}
		p.Init()
		if err := p.Set(&types.StringC{Value: "default"}, -1); err != nil {
			t.Fatalf("set: %v", err)
		}
		if err := p.Set(&types.StringC{Value: "downgrade-to-legacy"}, -1); err != nil {
			t.Fatalf("replace: %v", err)
		}
		result, err := p.Result()
		if err != nil {
			t.Fatalf("result: %v", err)
		}
		if len(result) != 1 || result[0].Data != "filter-api downgrade-to-legacy" {
			t.Fatalf("result = %v, want downgrade-to-legacy", result)
		}
		if err := p.Set(nil, -1); err != nil {
			t.Fatalf("clear: %v", err)
		}
		if _, err := p.Result(); err == nil {
			t.Fatal("expected fetch error after clearing filter-api")
		}
	})
}
