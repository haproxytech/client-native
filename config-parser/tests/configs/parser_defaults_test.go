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
	"errors"
	"testing"

	parser "github.com/haproxytech/client-native/v5/config-parser"
	parserErrors "github.com/haproxytech/client-native/v5/config-parser/errors"
	"github.com/haproxytech/client-native/v5/config-parser/options"
)

func TestDefaultsConfigs(t *testing.T) {
	tests := []struct {
		Name, Config string
	}{
		{"configDefaults", configDefaults},
		{"configDefaultsTwo", configDefaultsTwo},
		{"configDefaultsThree", configDefaultsThree},
		{"configDefaultsTwoSecond", configDefaultsTwoSecond},
		{"configDefaultsWithFrom", configDefaultsWithFrom},
	}
	for _, config := range tests {
		t.Run(config.Name, func(t *testing.T) {
			var buffer bytes.Buffer
			buffer.WriteString(config.Config)
			p, err := parser.New(options.Reader(&buffer))
			if err != nil {
				t.Fatalf(err.Error())
			}
			result := p.String()
			if result != config.Config {
				compare(t, config.Config, result)
				t.Fatalf("configurations does not match")
			}
		})
	}
}

func TestDefaultsConfigsSetDef(t *testing.T) {
	tests := []struct {
		Name, Config string
	}{
		{"configDefaultsTwo", configDefaultsTwo},
		{"configDefaultsThree", configDefaultsThree},
	}
	for _, config := range tests {
		t.Run(config.Name, func(t *testing.T) {
			var buffer bytes.Buffer
			buffer.WriteString(config.Config)
			p, err := parser.New(options.Reader(&buffer))
			if err != nil {
				t.Fatalf(err.Error())
			}
			err = p.SectionsDefaultsFromSet(parser.Defaults, "???", "nonexisting")
			if !errors.Is(err, parserErrors.ErrSectionMissing) {
				t.Fatalf("expected (%v) got (%v)", parserErrors.ErrSectionMissing, err)
			}
			err = p.SectionsDefaultsFromSet(parser.Defaults, "A", "nonexisting")
			if !errors.Is(err, parserErrors.ErrFromDefaultsSectionMissing) {
				t.Fatalf("expected (%v) got (%v)", parserErrors.ErrFromDefaultsSectionMissing, err)
			}
			err = p.SectionsDefaultsFromSet(parser.Defaults, "A", "withName")
			if err != nil {
				t.Fatalf("expected (%v) got (%v)", parserErrors.ErrFromDefaultsSectionMissing, err)
			}
		})
	}
}

func TestDefaultsConfigsSetCircular(t *testing.T) {
	tests := []struct {
		Name, Config string
	}{
		{"configDefaultsTwo", configDefaultsTwo},
		{"configDefaultsThree", configDefaultsThree},
	}
	for _, config := range tests {
		t.Run(config.Name, func(t *testing.T) {
			var buffer bytes.Buffer
			buffer.WriteString(config.Config)
			p, err := parser.New(options.Reader(&buffer))
			if err != nil {
				t.Fatalf(err.Error())
			}
			err = p.SectionsDefaultsFromSet(parser.Defaults, "A", "withName")
			if err != nil {
				t.Fatalf("expected (%v) got (%v)", parserErrors.ErrFromDefaultsSectionMissing, err)
			}
			err = p.SectionsDefaultsFromSet(parser.Defaults, "withName", "A")
			if !errors.Is(err, parserErrors.ErrCircularDependency) {
				t.Fatalf("expected (%v) got (%v)", parserErrors.ErrCircularDependency, err)
			}
		})
	}
}
