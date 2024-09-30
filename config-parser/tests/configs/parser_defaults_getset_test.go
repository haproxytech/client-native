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
	"testing"

	parser "github.com/haproxytech/client-native/v5/config-parser"
	"github.com/haproxytech/client-native/v5/config-parser/options"
)

func TestDefaultsConfigsGet(t *testing.T) { //nolint:gocognit
	tests := []struct {
		Name, Config string
	}{
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
			name, err := p.SectionsDefaultsFromGet(parser.Defaults, "A")
			if err != nil {
				t.Fatalf("error: %s", err.Error())
			}
			if name != "" {
				t.Fatalf("name != '': %s", name)
			}

			name, err = p.SectionsDefaultsFromGet(parser.Defaults, "B")
			if err != nil {
				t.Fatalf("error: %s", err.Error())
			}
			if name != "A" {
				t.Fatalf("name != 'A': %s", name)
			}

			name, err = p.SectionsDefaultsFromGet(parser.Frontends, "http1")
			if err != nil {
				t.Fatalf("error: %s", err.Error())
			}
			if name != "A" {
				t.Fatalf("name != 'A': %s", name)
			}

			name, err = p.SectionsDefaultsFromGet(parser.Frontends, "http2")
			if err != nil {
				t.Fatalf("error: %s", err.Error())
			}
			if name != "B" {
				t.Fatalf("name != 'B': %s", name)
			}

			name, err = p.SectionsDefaultsFromGet(parser.Backends, "default_backend1")
			if err != nil {
				t.Fatalf("error: %s", err.Error())
			}
			if name != "B" {
				t.Fatalf("name != 'B': %s", name)
			}

			name, err = p.SectionsDefaultsFromGet(parser.Backends, "default_backend2")
			if err != nil {
				t.Fatalf("error: %s", err.Error())
			}
			if name != "A" {
				t.Fatalf("name != 'A': %s", name)
			}

			result := p.String()
			if result != config.Config {
				compare(t, config.Config, result)
				t.Fatalf("configurations does not match")
			}
		})
	}
}

func TestDefaultsConfigsSet(t *testing.T) { //nolint:gocognit
	tests := []struct {
		Name, Config, Result string
	}{
		{"configDefaultsWithFrom", configDefaultsWithFrom, configDefaultsWithFromResult2},
	}
	for _, config := range tests {
		t.Run(config.Name, func(t *testing.T) {
			var buffer bytes.Buffer
			buffer.WriteString(config.Config)
			p, err := parser.New(options.Reader(&buffer))
			if err != nil {
				t.Fatalf(err.Error())
			}

			err = p.SectionsDefaultsFromSet(parser.Frontends, "http1", "B")
			if err != nil {
				t.Fatalf("error: %s", err.Error())
			}
			name, err := p.SectionsDefaultsFromGet(parser.Frontends, "http1")
			if err != nil {
				t.Fatalf("error: %s", err.Error())
			}
			if name != "B" {
				t.Fatalf("name != 'B': %s", name)
			}

			err = p.SectionsDefaultsFromSet(parser.Backends, "default_backend2", "B")
			if err != nil {
				t.Fatalf("error: %s", err.Error())
			}
			name, err = p.SectionsDefaultsFromGet(parser.Backends, "default_backend2")
			if err != nil {
				t.Fatalf("error: %s", err.Error())
			}
			if name != "B" {
				t.Fatalf("name != 'B': %s", name)
			}

			err = p.SectionsDefaultsFromSet(parser.Global, parser.GlobalSectionName, "B")
			if err == nil {
				t.Fatalf("no error!")
			}

			result := p.String()
			if result != config.Result {
				compare(t, config.Result, result)
				t.Fatalf("configurations does not match")
			}
		})
	}
}
