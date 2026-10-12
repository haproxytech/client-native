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
	"fmt"
	"strings"
	"sync"
	"testing"

	configparser "github.com/haproxytech/client-native/v6/config-parser"
	"github.com/haproxytech/client-native/v6/config-parser/options"
)

func namedDefaultsConfig(name string) string {
	return strings.Join([]string{
		"global",
		"  maxconn 100",
		"defaults " + name,
		"  mode http",
		"backend be",
		"  server s 127.0.0.1:80",
		"",
	}, "\n")
}

func parseSelectorParser(t *testing.T, opt ...options.ParserOption) configparser.Parser {
	t.Helper()
	p, err := configparser.New(opt...)
	if err != nil {
		t.Fatalf("parser.New: %v", err)
	}
	return p
}

// An independent parser must not change the defaults selector observed on
// another parser instance.
func TestDefaultsSelectorIsParserLocal(t *testing.T) {
	a := parseSelectorParser(t, options.String(namedDefaultsConfig("alpha")))
	b := parseSelectorParser(t, options.String(namedDefaultsConfig("beta")))

	if got := configparser.DefaultsSectionName(a); got != "alpha" {
		t.Fatalf("parser a selector = %q, want alpha", got)
	}
	if got := configparser.DefaultsSectionName(b); got != "beta" {
		t.Fatalf("parser b selector = %q, want beta", got)
	}
}

// The deprecated package global must stay untouched by the package itself.
func TestDefaultsSelectorDoesNotWritePackageGlobal(t *testing.T) {
	parseSelectorParser(t, options.String(namedDefaultsConfig("alpha")))

	if got := configparser.DefaultSectionName; got != "" {
		t.Fatalf("package global DefaultSectionName = %q, want empty", got)
	}
}

// A parser that does not implement the optional interface reports an empty
// selector rather than falling back to package state.
func TestDefaultsSelectorOptionalInterfaceCompatibility(t *testing.T) {
	p := parseSelectorParser(t, options.String(namedDefaultsConfig("alpha")))

	var wrapper struct{ configparser.Parser }
	wrapper.Parser = p

	if got := configparser.DefaultsSectionName(wrapper); got != "" {
		t.Fatalf("wrapper selector = %q, want empty", got)
	}
	if got := configparser.DefaultsSectionName(p); got != "alpha" {
		t.Fatalf("concrete selector = %q, want alpha", got)
	}
}

func TestDefaultsSelectorLiteralSemantics(t *testing.T) {
	t.Run("no defaults", func(t *testing.T) {
		p := parseSelectorParser(t, options.String("global\n  maxconn 1\n"))
		if got := configparser.DefaultsSectionName(p); got != "" {
			t.Fatalf("selector = %q, want empty", got)
		}
	})

	t.Run("named", func(t *testing.T) {
		p := parseSelectorParser(t, options.String(namedDefaultsConfig("alpha")))
		if got := configparser.DefaultsSectionName(p); got != "alpha" {
			t.Fatalf("selector = %q, want alpha", got)
		}
	})

	t.Run("unnamed", func(t *testing.T) {
		p := parseSelectorParser(t, options.String("defaults\n  mode http\n"))
		if got := configparser.DefaultsSectionName(p); got != "unnamed_defaults_1" {
			t.Fatalf("selector = %q, want unnamed_defaults_1", got)
		}
	})

	t.Run("multiple picks last", func(t *testing.T) {
		config := "defaults first\n  mode http\ndefaults second\n  mode tcp\n"
		p := parseSelectorParser(t, options.String(config))
		if got := configparser.DefaultsSectionName(p); got != "second" {
			t.Fatalf("selector = %q, want second", got)
		}
	})

	t.Run("NoNamedDefaultsFrom still tracks selector", func(t *testing.T) {
		p := parseSelectorParser(t, options.String(namedDefaultsConfig("alpha")), options.NoNamedDefaultsFrom)
		if got := configparser.DefaultsSectionName(p); got != "alpha" {
			t.Fatalf("selector = %q, want alpha", got)
		}
	})

	t.Run("SectionsCreate updates selector", func(t *testing.T) {
		p := parseSelectorParser(t, options.String("global\n  maxconn 1\n"))
		if err := p.SectionsCreate(configparser.Defaults, "created"); err != nil {
			t.Fatalf("SectionsCreate: %v", err)
		}
		if got := configparser.DefaultsSectionName(p); got != "created" {
			t.Fatalf("selector = %q, want created", got)
		}
	})

	t.Run("reload resets selector", func(t *testing.T) {
		p := parseSelectorParser(t, options.String(namedDefaultsConfig("alpha")))
		if err := p.Process(strings.NewReader("global\n  maxconn 1\n")); err != nil {
			t.Fatalf("Process: %v", err)
		}
		if got := configparser.DefaultsSectionName(p); got != "" {
			t.Fatalf("selector = %q, want empty after reload", got)
		}
	})

	t.Run("deleted current selection remains dangling", func(t *testing.T) {
		p := parseSelectorParser(t, options.String(namedDefaultsConfig("alpha")))
		if err := p.SectionsDelete(configparser.Defaults, "alpha"); err != nil {
			t.Fatalf("SectionsDelete: %v", err)
		}
		if got := configparser.DefaultsSectionName(p); got != "alpha" {
			t.Fatalf("selector = %q, want dangling alpha", got)
		}
	})
}

func TestDefaultsSelectorConcurrentParsers(t *testing.T) {
	const goroutines = 4
	const iterations = 50

	var wg sync.WaitGroup
	failures := make(chan string, goroutines)

	for g := range goroutines {
		wg.Add(1)
		go func(g int) {
			defer wg.Done()
			name := fmt.Sprintf("defaults_%d", g)
			p := parseSelectorParser(t, options.String(namedDefaultsConfig(name)))
			for range iterations {
				if got := configparser.DefaultsSectionName(p); got != name {
					failures <- fmt.Sprintf("selector = %q, want %q", got, name)
					return
				}
			}
		}(g)
	}

	wg.Wait()
	close(failures)
	for msg := range failures {
		t.Error(msg)
	}
}
