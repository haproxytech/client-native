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

package parser_test

import (
	"strings"
	"sync"
	"testing"

	parser "github.com/haproxytech/client-native/v5/config-parser"
	"github.com/haproxytech/client-native/v5/config-parser/options"
)

const defaultsSelectorAlphaConfig = `
global
    maxconn 100

defaults alpha
    mode http
`

const defaultsSelectorBetaConfig = `
global
    maxconn 100

defaults beta
    mode tcp
`

const defaultsSelectorNamelessConfig = `
global
    maxconn 100

defaults
    mode http
`

const defaultsSelectorNoDefaultsConfig = `
global
    maxconn 100
`

const defaultsSelectorMultipleDefaultsConfig = `
global
    maxconn 100

defaults alpha
    mode http

defaults beta
    mode tcp
`

func newDefaultsSelectorParser(t *testing.T, config string) parser.Parser {
	t.Helper()
	p, err := parser.New(options.Reader(strings.NewReader(config)))
	if err != nil {
		t.Fatalf("cannot create parser: %s", err)
	}
	return p
}

// TestDefaultsSelectorIsolation verifies that each parser reports the defaults
// section selected in its own configuration instead of a shared global that the
// most recently parsed configuration overwrites.
func TestDefaultsSelectorIsolation(t *testing.T) {
	pAlpha := newDefaultsSelectorParser(t, defaultsSelectorAlphaConfig)
	pBeta := newDefaultsSelectorParser(t, defaultsSelectorBetaConfig)

	if got := parser.DefaultsSectionName(pAlpha); got != "alpha" {
		t.Errorf("alpha parser selector %q, expected %q", got, "alpha")
	}
	if got := parser.DefaultsSectionName(pBeta); got != "beta" {
		t.Errorf("beta parser selector %q, expected %q", got, "beta")
	}
}

// TestDefaultsSelectorNameless verifies that an unnamed defaults section keeps
// the generated section name as its selector.
func TestDefaultsSelectorNameless(t *testing.T) {
	p := newDefaultsSelectorParser(t, defaultsSelectorNamelessConfig)
	if got := parser.DefaultsSectionName(p); got != "unnamed_defaults_1" {
		t.Errorf("named selector %q, expected %q", got, "unnamed_defaults_1")
	}
}

// TestDefaultsSelectorReloadResets verifies that reprocessing a configuration
// without a defaults section resets the selector.
func TestDefaultsSelectorReloadResets(t *testing.T) {
	p := newDefaultsSelectorParser(t, defaultsSelectorAlphaConfig)
	if got := parser.DefaultsSectionName(p); got != "alpha" {
		t.Fatalf("initial selector %q, expected %q", got, "alpha")
	}

	if err := p.Process(strings.NewReader(defaultsSelectorNoDefaultsConfig)); err != nil {
		t.Fatalf("reload parser: %s", err)
	}
	if got := parser.DefaultsSectionName(p); got != "" {
		t.Errorf("selector after reload %q, expected empty", got)
	}
}

// TestDefaultsSelectorSectionsCreate verifies that creating a defaults section
// through the parser updates the selector through the normal processing path.
func TestDefaultsSelectorSectionsCreate(t *testing.T) {
	p := newDefaultsSelectorParser(t, defaultsSelectorNoDefaultsConfig)
	if err := p.SectionsCreate(parser.Defaults, "gamma"); err != nil {
		t.Fatalf("create defaults section: %s", err)
	}
	if got := parser.DefaultsSectionName(p); got != "gamma" {
		t.Errorf("selector after SectionsCreate %q, expected %q", got, "gamma")
	}
}

// TestDefaultsSelectorMultipleDefaultsLastSelected verifies that the last
// defaults section in a configuration is the selected one.
func TestDefaultsSelectorMultipleDefaultsLastSelected(t *testing.T) {
	p := newDefaultsSelectorParser(t, defaultsSelectorMultipleDefaultsConfig)
	if got := parser.DefaultsSectionName(p); got != "beta" {
		t.Errorf("selector %q, expected %q", got, "beta")
	}
}

// TestDefaultsSelectorDeleteSelectedDangling verifies that deleting the
// selected defaults section leaves the selector unchanged (dangling), matching
// the existing behavior instead of inventing a fallback.
func TestDefaultsSelectorDeleteSelectedDangling(t *testing.T) {
	p := newDefaultsSelectorParser(t, defaultsSelectorMultipleDefaultsConfig)
	if got := parser.DefaultsSectionName(p); got != "beta" {
		t.Fatalf("selector %q, expected %q", got, "beta")
	}

	if err := p.SectionsDelete(parser.Defaults, "beta"); err != nil {
		t.Fatalf("delete defaults section: %s", err)
	}
	if got := parser.DefaultsSectionName(p); got != "beta" {
		t.Errorf("selector after deleting selected defaults %q, expected dangling %q", got, "beta")
	}
}

// TestDefaultsSelectorAbsentDefaults verifies that a configuration without a
// defaults section reports an empty selector.
func TestDefaultsSelectorAbsentDefaults(t *testing.T) {
	p := newDefaultsSelectorParser(t, defaultsSelectorNoDefaultsConfig)
	if got := parser.DefaultsSectionName(p); got != "" {
		t.Errorf("selector %q, expected empty", got)
	}
}

// TestDefaultsSelectorNoNamedDefaultsFrom verifies that NoNamedDefaultsFrom
// does not suppress the defaults selector.
func TestDefaultsSelectorNoNamedDefaultsFrom(t *testing.T) {
	p, err := parser.New(
		options.Reader(strings.NewReader(defaultsSelectorAlphaConfig)),
		options.NoNamedDefaultsFrom,
	)
	if err != nil {
		t.Fatalf("cannot create parser: %s", err)
	}
	if got := parser.DefaultsSectionName(p); got != "alpha" {
		t.Errorf("selector %q, expected %q", got, "alpha")
	}
}

// foreignParser wraps a Parser and does not expose the optional
// DefaultSectionNamer interface, as a foreign implementation or wrapper would.
type foreignParser struct {
	parser.Parser
}

// TestDefaultsSelectorForeignParser verifies that a Parser implementation
// without the optional accessor still compiles and yields an empty selector,
// without falling back to a process-wide value.
func TestDefaultsSelectorForeignParser(t *testing.T) {
	inner := newDefaultsSelectorParser(t, defaultsSelectorAlphaConfig)
	wrapper := foreignParser{Parser: inner}

	if got := parser.DefaultsSectionName(wrapper); got != "" {
		t.Errorf("foreign parser selector %q, expected empty", got)
	}
	if got := parser.DefaultsSectionName(inner); got != "alpha" {
		t.Errorf("concrete parser selector %q, expected %q", got, "alpha")
	}
}

// TestDefaultsSelectorConcurrentCreation verifies that independent parsers
// created and processed concurrently inside goroutines each report their own
// selector, without a shared mutable global. This exercises the original
// concurrent Process/creation write hazard, not just concurrent reads.
func TestDefaultsSelectorConcurrentCreation(t *testing.T) {
	cases := []struct {
		config string
		want   string
	}{
		{config: defaultsSelectorAlphaConfig, want: "alpha"},
		{config: defaultsSelectorBetaConfig, want: "beta"},
		{config: defaultsSelectorNamelessConfig, want: "unnamed_defaults_1"},
	}

	const iterations = 50
	var wg sync.WaitGroup
	for _, tc := range cases {
		tc := tc
		for i := 0; i < 4; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < iterations; j++ {
					p, err := parser.New(options.Reader(strings.NewReader(tc.config)))
					if err != nil {
						t.Errorf("cannot create parser: %s", err)
						return
					}
					if got := parser.DefaultsSectionName(p); got != tc.want {
						t.Errorf("selector %q, expected %q", got, tc.want)
						return
					}
				}
			}()
		}
	}
	wg.Wait()
}
