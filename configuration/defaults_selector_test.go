// Copyright 2026 HAProxy Technologies
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package configuration

import (
	"fmt"
	"slices"
	"sync"
	"testing"

	parser "github.com/haproxytech/client-native/v6/config-parser"
	parser_options "github.com/haproxytech/client-native/v6/config-parser/options"
	"github.com/haproxytech/client-native/v6/models"
)

const defaultsSelectorAlphaConfig = `# _version=1
global
  maxconn 100
defaults alpha
  mode http
  http-request set-header X-Alpha a
  tcp-check connect
  http-error status 400
  log stdout format raw local0
frontend fe
  default_backend be
backend be
  server s 127.0.0.1:80
`

const defaultsSelectorBetaConfig = `# _version=1
global
  maxconn 100
defaults beta
  mode tcp
  http-request set-header X-Beta b
  tcp-check connect
  http-error status 401
  log stdout format raw local0
frontend fe
  default_backend be
backend be
  server s 127.0.0.1:80
`

func newDefaultsSelectorParser(t *testing.T, config string) parser.Parser {
	t.Helper()
	p, err := parser.New(parser_options.String(config))
	if err != nil {
		t.Fatalf("parser.New: %v", err)
	}
	return p
}

func defaultsSectionExists(p parser.Parser, name string) bool {
	sections, err := p.SectionsGet(parser.Defaults)
	if err != nil {
		return false
	}
	return slices.Contains(sections, name)
}

func newDefaultsSelectorClient(t *testing.T, config string) *client {
	t.Helper()
	return &client{parser: newDefaultsSelectorParser(t, config)}
}

// Legacy getters must resolve the defaults name from the client's own parser,
// even after another independent parser selected a different defaults section.
// The v6.0 defaults getter does not return the section name, so this checks
// selection through the section contents instead.
func TestDefaultsSelectorLegacyGettersUseOwnParser(t *testing.T) {
	a := newDefaultsSelectorClient(t, defaultsSelectorAlphaConfig)
	_ = newDefaultsSelectorClient(t, defaultsSelectorBetaConfig)

	_, d, err := a.GetDefaultsConfiguration("")
	if err != nil {
		t.Fatalf("GetDefaultsConfiguration: %v", err)
	}
	if d.Mode != "http" {
		t.Fatalf("GetDefaultsConfiguration mode = %q, want http", d.Mode)
	}

	_, sd, err := a.GetStructuredDefaultsConfiguration("")
	if err != nil {
		t.Fatalf("GetStructuredDefaultsConfiguration: %v", err)
	}
	if sd.Name != "alpha" || sd.Mode != "http" {
		t.Fatalf("GetStructuredDefaultsConfiguration = name %q mode %q, want alpha/http", sd.Name, sd.Mode)
	}

	if parser.DefaultSectionName != "" {
		t.Fatalf("package global DefaultSectionName = %q, want empty", parser.DefaultSectionName)
	}
}

// Rule and check accessors with an implicit defaults parent must read from the
// client's own parser, not from whichever parser selected a defaults section last.
func TestDefaultsSelectorRulesUseOwnParser(t *testing.T) {
	a := newDefaultsSelectorClient(t, defaultsSelectorAlphaConfig)
	_ = newDefaultsSelectorClient(t, defaultsSelectorBetaConfig)

	t.Run("http-request", func(t *testing.T) {
		_, rules, err := a.GetHTTPRequestRules(DefaultsParentName, "", "")
		if err != nil {
			t.Fatalf("GetHTTPRequestRules: %v", err)
		}
		if len(rules) != 1 || rules[0].HdrName != "X-Alpha" {
			t.Fatalf("rules = %#v, want single X-Alpha header", rules)
		}
	})

	t.Run("tcp-check", func(t *testing.T) {
		_, checks, err := a.GetTCPChecks(DefaultsParentName, "", "")
		if err != nil {
			t.Fatalf("GetTCPChecks: %v", err)
		}
		if len(checks) != 1 {
			t.Fatalf("checks = %#v, want single check", checks)
		}
	})

	t.Run("http-error", func(t *testing.T) {
		_, rules, err := a.GetHTTPErrorRules(DefaultsParentName, "", "")
		if err != nil {
			t.Fatalf("GetHTTPErrorRules: %v", err)
		}
		if len(rules) != 1 || rules[0].Status != 400 {
			t.Fatalf("rules = %#v, want single status 400", rules)
		}
	})

	t.Run("log", func(t *testing.T) {
		_, targets, err := a.GetLogTargets(DefaultsParentName, "", "")
		if err != nil {
			t.Fatalf("GetLogTargets: %v", err)
		}
		if len(targets) != 1 || targets[0].Address != "stdout" {
			t.Fatalf("targets = %#v, want single stdout target", targets)
		}
	})
}

// A transaction parser must be selected against its own defaults section.
func TestDefaultsSelectorTransactionParserIsolation(t *testing.T) {
	master := newDefaultsSelectorParser(t, defaultsSelectorAlphaConfig)
	tx := newDefaultsSelectorParser(t, defaultsSelectorBetaConfig)
	c := &client{
		parser: master,
		parsers: map[string]parser.Parser{
			"tx1": tx,
		},
	}

	_, d, err := c.GetDefaultsConfiguration("tx1")
	if err != nil {
		t.Fatalf("GetDefaultsConfiguration(tx1): %v", err)
	}
	if d.Mode != "tcp" {
		t.Fatalf("transaction defaults mode = %q, want tcp", d.Mode)
	}

	_, m, err := c.GetDefaultsConfiguration("")
	if err != nil {
		t.Fatalf("GetDefaultsConfiguration(master): %v", err)
	}
	if m.Mode != "http" {
		t.Fatalf("master defaults mode = %q, want http", m.Mode)
	}
}

func TestDefaultsSelectorConcurrentClients(t *testing.T) {
	const goroutines = 4
	const iterations = 50

	var wg sync.WaitGroup
	failures := make(chan string, goroutines)

	for g := range goroutines {
		wg.Add(1)
		go func(g int) {
			defer wg.Done()
			name := fmt.Sprintf("defaults_%d", g)
			config := fmt.Sprintf("# _version=1\nglobal\n  maxconn 1\ndefaults %s\n  mode http\n  http-request set-header X-Name %s\n", name, name)
			c := &client{parser: newDefaultsSelectorParser(t, config)}
			for range iterations {
				_, d, err := c.GetDefaultsConfiguration("")
				if err != nil || d.Mode != "http" {
					failures <- fmt.Sprintf("GetDefaultsConfiguration = %#v err %v, want %q/http", d, err, name)
					return
				}
				_, rules, err := c.GetHTTPRequestRules(DefaultsParentName, "", "")
				if err != nil || len(rules) != 1 || rules[0].HdrName != "X-Name" {
					failures <- fmt.Sprintf("GetHTTPRequestRules = %#v err %v", rules, err)
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

const defaultsSelectorGammaConfig = `# _version=1
global
  maxconn 100
defaults gamma
  mode tcp
frontend fe
  default_backend be
backend be
  server s 127.0.0.1:80
`

func newDefaultsSelectorTransactionClient(t *testing.T) (*client, parser.Parser, parser.Parser) {
	t.Helper()
	master := newDefaultsSelectorParser(t, defaultsSelectorAlphaConfig)
	tx := newDefaultsSelectorParser(t, defaultsSelectorGammaConfig)
	c := &client{
		parser:  master,
		parsers: map[string]parser.Parser{"tx1": tx},
	}
	c.TransactionClient = c
	return c, master, tx
}

// Serializing an already-existing current section must tolerate the
// ErrSectionAlreadyExists guard for the owning parser's selector, even after an
// independent parser selected a different defaults section.
func TestDefaultsSelectorSerializeExistingCurrentSection(t *testing.T) {
	a := newDefaultsSelectorClient(t, defaultsSelectorAlphaConfig)
	_ = newDefaultsSelectorParser(t, defaultsSelectorBetaConfig)

	_, alphaDefaults, err := a.GetStructuredDefaultsConfiguration("")
	if err != nil {
		t.Fatalf("GetStructuredDefaultsConfiguration: %v", err)
	}

	p := a.Parser()
	if err := serializeDefaultsSection(StructuredToParserArgs{Parser: &p}, alphaDefaults); err != nil {
		t.Fatalf("serializeDefaultsSection: %v", err)
	}

	if !defaultsSectionExists(a.parser, "alpha") {
		t.Fatal("alpha defaults section missing after serialization")
	}
	_, d, err := a.GetDefaultsConfiguration("")
	if err != nil {
		t.Fatalf("GetDefaultsConfiguration: %v", err)
	}
	if d.Mode != "http" {
		t.Fatalf("defaults mode = %q, want http", d.Mode)
	}
}

// A legacy defaults write on a transaction must target the transaction's own
// defaults selector, not the master parser's and not a package global.
func TestDefaultsSelectorTransactionLegacyDefaultsWrite(t *testing.T) {
	c, master, tx := newDefaultsSelectorTransactionClient(t)

	data := &models.Defaults{DefaultsBase: models.DefaultsBase{Name: "gamma", Mode: "http"}}
	if err := c.PushDefaultsConfiguration(data, "tx1", 0); err != nil {
		t.Fatalf("PushDefaultsConfiguration(tx1): %v", err)
	}

	if !defaultsSectionExists(tx, "gamma") {
		t.Fatal("transaction gamma defaults missing after push")
	}
	if defaultsSectionExists(tx, "alpha") {
		t.Fatal("transaction push wrote the master alpha section")
	}
	if !defaultsSectionExists(master, "alpha") || defaultsSectionExists(master, "gamma") {
		t.Fatal("master parser was modified by the transaction push")
	}
	_, d, err := c.GetDefaultsConfiguration("tx1")
	if err != nil || d.Mode != "http" {
		t.Fatalf("transaction defaults mode = %q err %v, want http", d.Mode, err)
	}
}

// A structured defaults replacement on a transaction must resolve the name from
// the transaction parser only and leave the master parser untouched.
func TestDefaultsSelectorTransactionStructuredDefaultsWrite(t *testing.T) {
	c, master, tx := newDefaultsSelectorTransactionClient(t)

	data := &models.Defaults{DefaultsBase: models.DefaultsBase{Name: "gamma", Mode: "http"}}
	if err := c.PushStructuredDefaultsConfiguration(data, "tx1", 0); err != nil {
		t.Fatalf("PushStructuredDefaultsConfiguration(tx1): %v", err)
	}

	if data.Name != "gamma" {
		t.Fatalf("data.Name = %q, want gamma", data.Name)
	}
	if !defaultsSectionExists(tx, "gamma") || defaultsSectionExists(tx, "alpha") {
		t.Fatal("transaction structured push wrote the wrong section")
	}
	if !defaultsSectionExists(master, "alpha") || defaultsSectionExists(master, "gamma") {
		t.Fatal("master parser was modified by the structured transaction push")
	}
}

// A rule write with an implicit defaults parent must resolve to the transaction's
// own selector; explicit parent names remain untouched.
func TestDefaultsSelectorTransactionEmptyParentRuleWrite(t *testing.T) {
	c, master, _ := newDefaultsSelectorTransactionClient(t)

	rule := &models.HTTPRequestRule{Type: "set-header", HdrName: "X-Tx", HdrFormat: "v"}
	if err := c.CreateHTTPRequestRule(0, DefaultsParentName, "", rule, "tx1", 0); err != nil {
		t.Fatalf("CreateHTTPRequestRule(tx1): %v", err)
	}

	_, txRules, err := c.GetHTTPRequestRules(DefaultsParentName, "", "tx1")
	if err != nil || len(txRules) != 1 || txRules[0].HdrName != "X-Tx" {
		t.Fatalf("transaction rules = %#v err %v, want single X-Tx", txRules, err)
	}

	_, masterRules, err := c.GetHTTPRequestRules(DefaultsParentName, "", "")
	if err != nil || len(masterRules) != 1 || masterRules[0].HdrName != "X-Alpha" {
		t.Fatalf("master rules = %#v err %v, want unchanged X-Alpha", masterRules, err)
	}
	if !defaultsSectionExists(master, "alpha") {
		t.Fatal("master alpha defaults section missing")
	}

	_, explicit, err := c.GetHTTPRequestRules(DefaultsParentName, "gamma", "tx1")
	if err != nil || len(explicit) != 1 {
		t.Fatalf("explicit gamma rules = %#v err %v, want single rule", explicit, err)
	}
}
