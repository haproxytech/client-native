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

func TestFilterConfigGrammar(t *testing.T) {
	valid := []struct {
		name    string
		line    string
		data    string
		comment string
	}{
		{
			name: "trace minimal",
			line: "filter-config trace id A",
			data: "filter-config trace id A",
		},
		{
			name:    "trace complete",
			line:    "filter-config trace id A enabled quiet random-parsing random-forwarding max-fwd 0 hexdump # {\"k\":\"v\"}",
			data:    "filter-config trace id A enabled quiet random-parsing random-forwarding max-fwd 0 hexdump",
			comment: `{"k":"v"}`,
		},
		{
			name: "http compression request without enabled",
			line: "filter-config http-comp-req",
			data: "filter-config http-comp-req",
		},
		{
			name: "decompression request",
			line: "filter-config decomp-req",
			data: "filter-config decomp-req",
		},
		{
			name: "decompression response enabled",
			line: "filter-config decomp-res enabled",
			data: "filter-config decomp-res enabled",
		},
		{
			name: "http compression response enabled",
			line: "filter-config http-comp-res enabled",
			data: "filter-config http-comp-res enabled",
		},
		{
			name: "spoe",
			line: "filter-config spoe id S enabled engine e config f",
			data: "filter-config spoe id S enabled engine e config f",
		},
		{
			name: "cache store uses id as cache name",
			line: "filter-config cache-store id C enabled",
			data: "filter-config cache-store id C enabled",
		},
		{
			name: "fcgi uses id as application name",
			line: "filter-config fcgi id F",
			data: "filter-config fcgi id F",
		},
		{
			name: "bandwidth per stream",
			line: "filter-config bwlim-in id B enabled default-limit 10k default-period 1s min-size 2k",
			data: "filter-config bwlim-in id B enabled default-limit 10k default-period 1s min-size 2k",
		},
		{
			name: "bandwidth shared",
			line: "filter-config bwlim-out id C limit 1m key src table t min-size 1k",
			data: "filter-config bwlim-out id C limit 1m key src table t min-size 1k",
		},
		{
			name: "lua without arguments",
			line: "filter-config lua id my-filter",
			data: "filter-config lua id my-filter",
		},
		{
			name: "lua with arguments",
			line: "filter-config lua id my-filter enabled arg1 arg2 arg3",
			data: "filter-config lua id my-filter enabled arg1 arg2 arg3",
		},
		{
			name: "lua with arguments and no enabled token",
			line: "filter-config lua id my-filter arg1",
			data: "filter-config lua id my-filter arg1",
		},
	}
	for _, tt := range valid {
		t.Run(tt.name, func(t *testing.T) {
			p := &parsers.FilterConfig{}
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
		"filter-config",
		"filter-config bogus id A",
		"filter-config trace",
		"filter-config trace A",
		"filter-config trace A enabled",
		"filter-config trace id",
		"filter-config trace id bad!id",
		"filter-config trace id A extra",
		"filter-config trace id A max-fwd",
		"filter-config http-comp-req A",
		"filter-config http-comp-req id A",
		"filter-config decomp-req A",
		"filter-config decomp-req id A",
		"filter-config spoe id S",
		"filter-config spoe id S engine e",
		"filter-config spoe S engine e config f",
		"filter-config cache-store C enabled",
		"filter-config lua",
		"filter-config lua my-filter",
		"filter-config lua id",
		"filter-config lua id bad!id",
		"filter-config bwlim-in id B",
		"filter-config bwlim-in id B limit 10k",
		"filter-config bwlim-in id B default-limit 10k",
		"filter-config bwlim-in id B limit 10k key src table",
		"filter-config bwlim-in id B default-limit nope default-period 1s",
	}
	for _, line := range invalid {
		t.Run(line, func(t *testing.T) {
			p := &parsers.FilterConfig{}
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

func TestFilterConfigDisabledByDefault(t *testing.T) {
	p := &parsers.FilterConfig{}
	p.Init()
	if err := ProcessLine("filter-config http-comp-req", p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, err := p.Get(false)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	configs := data.([]models.FilterConfig)
	if len(configs) != 1 {
		t.Fatalf("configs = %d, want 1", len(configs))
	}
	if configs[0].Enabled {
		t.Fatal("enabled must stay false when the directive omits it")
	}
	if configs[0].TraceMaxFwd != nil {
		t.Fatal("trace max-fwd must stay nil when absent")
	}
}
