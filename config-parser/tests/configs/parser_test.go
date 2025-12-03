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
