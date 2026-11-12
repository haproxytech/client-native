/*
Copyright 2019 HAProxy Technologies

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
	"bytes"
	"os"
	"strings"
	"testing"

	parser "github.com/haproxytech/client-native/v6/config-parser"
	"github.com/haproxytech/client-native/v6/config-parser/options"
	"github.com/haproxytech/client-native/v6/config-parser/parsers/filters"
	"github.com/haproxytech/client-native/v6/config-parser/types"
	"github.com/haproxytech/client-native/v6/models"
)

func TestWholeConfigs(t *testing.T) {
	tests := []struct {
		Name, Config string
	}{
		{"configBasic1", configBasic1},
		{"configBasic2", configBasic2},
		{"configFull", configFull},
		{"configSnippet", configSnippet},
	}
	for _, config := range tests {
		t.Run(config.Name, func(t *testing.T) {
			var buffer bytes.Buffer
			buffer.WriteString(config.Config)
			p, err := parser.New(options.Reader(&buffer))
			if err != nil {
				t.Fatal(err.Error())
			}
			result := p.String()
			if result != config.Config {
				compare(t, config.Config, result)
				t.Fatalf("configurations does not match")
			}
		})
	}
}

func TestWholeConfigsFail(t *testing.T) {
	tests := []struct {
		Name, Config string
	}{
		{"configFail1", configFail1},
		{"configFail2", configFail2},
		{"configFail3", configFail3},
		{"configFail4", configFail4},
	}
	for _, config := range tests {
		t.Run(config.Name, func(t *testing.T) {
			var buffer bytes.Buffer
			buffer.WriteString(config.Config)
			p, err := parser.New(options.Reader(&buffer))
			if err != nil {
				t.Fatal(err.Error())
			}
			result := p.String()
			if result == config.Config {
				compare(t, config.Config, result)
				t.Fatalf("configurations does not match")
			}
		})
	}
}

func compare(t *testing.T, configOriginal, configResult string) { //nolint:thelper
	original := strings.Split(configOriginal, "\n")
	result := strings.Split(configResult, "\n")
	if len(original) != len(result) {
		t.Logf("not the same size: original: %d, result: %d", len(original), len(result))
		return
	}
	for index, line := range original {
		if line != result[index] {
			t.Logf("line %d: '%s' != '%s'", index+3, line, result[index])
		}
	}
}

func TestGeneratedConfig(t *testing.T) {
	var buffer bytes.Buffer
	buffer.WriteString(generatedConfig)
	p, err := parser.New(options.DisableUnProcessed, options.Reader(&buffer))
	if err != nil {
		t.Fatal(err.Error())
	}
	result := p.String()
	for _, configLine := range configTests {
		count := strings.Count(result, configLine.Line)
		if count != configLine.Count {
			_ = os.WriteFile("/tmp/HAGEN.cfg", []byte(result), 0o644)
			t.Fatalf("line '%s' found %d times, expected %d times", configLine.Line, count, configLine.Count)
		}
	}
}

func TestHashConfig(t *testing.T) {
	var buffer bytes.Buffer
	buffer.WriteString(configBasicHash)
	p, err := parser.New(options.UseMd5Hash, options.Reader(&buffer))
	if err != nil {
		t.Fatal(err.Error())
	}
	result, err := p.StringWithHash()
	if err != nil {
		t.Fatal(err.Error())
	}
	if result != configBasicHash {
		compare(t, configBasicHash, result)
		t.Fatalf("configurations does not match")
	}
}

func TestConfigUseV2HTTPCheck(t *testing.T) {
	var buffer bytes.Buffer
	buffer.WriteString(configBasicUseV2HTTPCheck)
	p, err := parser.New(options.UseV2HTTPCheck, options.Reader(&buffer))
	if err != nil {
		t.Fatal(err.Error())
	}
	result := p.String() //nolint:ifshort
	if result != configBasicUseV2HTTPCheck {
		compare(t, configBasicUseV2HTTPCheck, result)
		t.Fatalf("configurations does not match")
	}
}

func TestListenSectionParsers(t *testing.T) {
	var buffer bytes.Buffer
	buffer.WriteString(configFull)
	p, err := parser.New(options.UseListenSectionParsers, options.Reader(&buffer))
	if err != nil {
		t.Fatal(err.Error())
	}

	result := p.String() //nolint:ifshort
	if result != configFull {
		compare(t, configFull, result)
		t.Fatalf("configurations does not match")
	}
}

func TestDefaultSectionsSkipOnWriteParsers(t *testing.T) {
	tests := []struct {
		Name                  string
		Config                string
		Result                string
		DefaultsSectionToSkip []string
	}{
		{
			"with option DefaultSectionsSkipOnWrite and section present",
			configDefaultSectionsSkipOnWrite_WithSectionName, configDefaultSectionsSkipOnWrite_WithoutSectionName,
			[]string{"to_not_serialize"},
		},
		{
			"with option DefaultSectionsSkipOnWrite and section present - multiple options",
			configDefaultSectionsSkipOnWrite_WithSectionName, configDefaultSectionsSkipOnWrite_WithoutSectionName,
			[]string{"to_not_serialize", "to_not_serialize2"},
		},
		{
			"no option DefaultSectionsSkipOnWrite and section absent",
			configDefaultSectionsSkipOnWrite_WithoutSectionName, configDefaultSectionsSkipOnWrite_WithoutSectionName,
			[]string{"to_not_serialize"},
		},
		{
			"with option DefaultSectionsSkipOnWrite and section present different option name",
			configDefaultSectionsSkipOnWrite_WithSectionName, configDefaultSectionsSkipOnWrite_WithSectionName,
			[]string{"other"},
		},
	}

	for _, tt := range tests {
		var buffer bytes.Buffer
		buffer.WriteString(tt.Config)
		p, err := parser.New(options.DefaultSectionsSkipOnWrite(tt.DefaultsSectionToSkip), options.Reader(&buffer))
		if err != nil {
			t.Fatal(err.Error())
		}

		result := p.String() //nolint:ifshort
		if result != tt.Result {
			compare(t, tt.Config, tt.Result)
			t.Fatalf("configurations does not match")
		}
	}
}

// TestFilterSequenceDirective asserts that the documented filter-sequence
// grammar is parsed into a typed parser, while the obsolete comma-list syntax
// stays unprocessed and ordinary compression filters remain typed.
func TestFilterSequenceDirective(t *testing.T) { //nolint:gocognit
	const obsolete = "filter-sequence request lua.my-filter,comp-req"
	tests := []struct {
		name    string
		section parser.Section
	}{
		{"frontend", parser.Frontends},
		{"backend", parser.Backends},
		{"listen", parser.Listen},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := tt.name + ` test
  filter comp-req
  filter comp-res
  filter-sequence request trace/A:before(trace/B)
  filter-sequence request lua.my-filter,comp-req
`
			p, err := parser.New(options.String(config), options.UseListenSectionParsers)
			if err != nil {
				t.Fatal(err)
			}

			data, err := p.Get(tt.section, "test", "filter-sequence", false)
			if err != nil {
				t.Fatalf("filter-sequence fetch: %v", err)
			}
			sequences, ok := data.([]models.FilterSequence)
			if !ok {
				t.Fatalf("unexpected filter-sequence type %T", data)
			}
			if len(sequences) != 1 {
				t.Fatalf("filter-sequence lines = %d, want 1", len(sequences))
			}
			sequence := sequences[0]
			if sequence.Side != "request" || len(sequence.Rules) != 1 {
				t.Fatalf("unexpected parsed sequence %+v", sequence)
			}
			rule := sequence.Rules[0]
			if rule.Instance.Class != "trace" || rule.Instance.ID != "A" ||
				rule.Direction != "before" ||
				rule.Target.Class != "trace" || rule.Target.ID != "B" {
				t.Fatalf("unexpected parsed rule %+v", rule)
			}

			// The obsolete comma-list syntax must remain unprocessed.
			data, err = p.Get(tt.section, "test", "", false)
			if err != nil {
				t.Fatalf("unprocessed fetch: %v", err)
			}
			unprocessed, ok := data.([]types.UnProcessed)
			if !ok {
				t.Fatalf("unexpected unprocessed type %T", data)
			}
			got := []string{}
			for _, entry := range unprocessed {
				if strings.HasPrefix(entry.Value, "filter-sequence ") {
					got = append(got, entry.Value)
				}
			}
			if len(got) != 1 || got[0] != obsolete {
				t.Fatalf("unprocessed filter-sequence lines = %v, want [%s]", got, obsolete)
			}

			// Compression filters must remain typed.
			fdata, err := p.Get(tt.section, "test", "filter", false)
			if err != nil {
				t.Fatalf("filter fetch: %v", err)
			}
			tsFilters, ok := fdata.([]types.Filter)
			if !ok {
				t.Fatalf("unexpected filter type %T", fdata)
			}
			typed := map[string]bool{}
			for _, filter := range tsFilters {
				switch filter.(type) {
				case *filters.CompReq:
					typed["comp-req"] = true
				case *filters.CompRes:
					typed["comp-res"] = true
				}
			}
			if !typed["comp-req"] || !typed["comp-res"] {
				t.Fatalf("expected typed comp-req and comp-res filters, got %v", typed)
			}

			result := p.String()
			for _, want := range []string{
				"filter comp-req",
				"filter comp-res",
				"filter-sequence request trace/A:before(trace/B)",
				"filter-sequence request lua.my-filter,comp-req",
			} {
				if !strings.Contains(result, want) {
					t.Fatalf("serialized output missing %q:\n%s", want, result)
				}
			}
		})
	}
}
