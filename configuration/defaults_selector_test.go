// Copyright 2019 HAProxy Technologies
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
//

package configuration

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/haproxytech/client-native/v5/configuration/options"
	"github.com/haproxytech/client-native/v5/misc"
	"github.com/haproxytech/client-native/v5/models"
)

// defaultsSelectorConfigAlpha and defaultsSelectorConfigBeta carry distinct
// directives on their defaults sections so that selecting the wrong defaults
// section is observable through real production reads and writes.
const defaultsSelectorConfigAlpha = `
# _version=1
global
    maxconn 100

defaults alpha
    mode http
    maxconn 100
    log 127.0.0.1:514 local0
    http-check send-state
    http-check disable-on-404
    http-error status 503 content-type application/json file /alpha/503
    tcp-check connect port 443
`

const defaultsSelectorConfigBeta = `
# _version=1
global
    maxconn 100

defaults beta
    mode tcp
    maxconn 200
    log 127.0.0.1:515 local0
    http-check disable-on-404
    http-error status 429 content-type application/json file /beta/429
    tcp-check connect port 80
`

func writeDefaultsSelectorConfig(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "haproxy.cfg")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("cannot write config: %s", err)
	}
	return path
}

func newDefaultsSelectorClient(t *testing.T, path string) Configuration {
	t.Helper()
	c, err := New(context.Background(),
		options.ConfigurationFile(path),
		options.HAProxyBin("echo"),
		options.TransactionsDir(t.TempDir()),
	)
	if err != nil {
		t.Fatalf("cannot create client: %s", err)
	}
	return c
}

func newDefaultsSelectorPair(t *testing.T) (cAlpha, cBeta Configuration) {
	t.Helper()
	cAlpha = newDefaultsSelectorClient(t, writeDefaultsSelectorConfig(t, defaultsSelectorConfigAlpha))
	// Parsing the second client is what used to overwrite the shared global.
	cBeta = newDefaultsSelectorClient(t, writeDefaultsSelectorConfig(t, defaultsSelectorConfigBeta))
	return cAlpha, cBeta
}

func startDefaultsSelectorTransaction(t *testing.T, c Configuration) (version int64, transactionID string) {
	t.Helper()
	v, _, err := c.GetDefaultsConfiguration("")
	if err != nil {
		t.Fatalf("defaults before transaction: %s", err)
	}
	tx, err := c.StartTransaction(v)
	if err != nil {
		t.Fatalf("start transaction: %s", err)
	}
	t.Cleanup(func() { _ = c.DeleteTransaction(tx.ID) })
	return v, tx.ID
}

// assertAlphaRuleReads checks that rule reads resolve to cAlpha's alpha
// defaults section, never to another section. An empty parentName exercises
// empty-parent resolution against the parser's selected defaults section; an
// explicit parentName pins the read to that named section.
func assertAlphaRuleReads(t *testing.T, c Configuration, parentName string, transactionID string) {
	t.Helper()

	_, targets, err := c.GetLogTargets("defaults", parentName, transactionID)
	if err != nil {
		t.Fatalf("alpha log targets: %s", err)
	}
	if len(targets) != 1 || targets[0].Address != "127.0.0.1:514" {
		t.Errorf("alpha log targets %+v, expected a single 127.0.0.1:514", targets)
	}

	_, target, err := c.GetLogTarget(0, "defaults", parentName, transactionID)
	if err != nil {
		t.Fatalf("alpha log target: %s", err)
	}
	if target.Address != "127.0.0.1:514" {
		t.Errorf("alpha log target %q, expected %q", target.Address, "127.0.0.1:514")
	}

	_, checks, err := c.GetHTTPChecks("defaults", parentName, transactionID)
	if err != nil {
		t.Fatalf("alpha http checks: %s", err)
	}
	if len(checks) != 2 {
		t.Errorf("alpha http checks %d, expected 2", len(checks))
	}

	_, check, err := c.GetHTTPCheck(0, "defaults", parentName, transactionID)
	if err != nil {
		t.Fatalf("alpha http check: %s", err)
	}
	if check.Type != "send-state" {
		t.Errorf("alpha http check type %q, expected %q", check.Type, "send-state")
	}

	_, rules, err := c.GetHTTPErrorRules("defaults", parentName, transactionID)
	if err != nil {
		t.Fatalf("alpha http error rules: %s", err)
	}
	if len(rules) != 1 || rules[0].Status != 503 {
		t.Errorf("alpha http error rules %+v, expected a single status 503", rules)
	}

	_, rule, err := c.GetHTTPErrorRule(0, "defaults", parentName, transactionID)
	if err != nil {
		t.Fatalf("alpha http error rule: %s", err)
	}
	if rule.Status != 503 || rule.ReturnContent != "/alpha/503" {
		t.Errorf("alpha http error rule status %d content %q, expected 503 /alpha/503", rule.Status, rule.ReturnContent)
	}

	_, tcpChecks, err := c.GetTCPChecks("defaults", parentName, transactionID)
	if err != nil {
		t.Fatalf("alpha tcp checks: %s", err)
	}
	if len(tcpChecks) != 1 {
		t.Errorf("alpha tcp checks %d, expected 1", len(tcpChecks))
	}

	_, tcpCheck, err := c.GetTCPCheck(0, "defaults", parentName, transactionID)
	if err != nil {
		t.Fatalf("alpha tcp check: %s", err)
	}
	if tcpCheck.PortString != "443" {
		t.Errorf("alpha tcp check port %q, expected %q", tcpCheck.PortString, "443")
	}
}

// assertBetaRuleReads checks that cBeta's own beta defaults section is intact
// and never resolves to alpha's.
func assertBetaRuleReads(t *testing.T, c Configuration) {
	t.Helper()

	_, targets, err := c.GetLogTargets("defaults", "", "")
	if err != nil {
		t.Fatalf("beta log targets: %s", err)
	}
	if len(targets) != 1 || targets[0].Address != "127.0.0.1:515" {
		t.Errorf("beta log targets %+v, expected a single 127.0.0.1:515", targets)
	}

	_, rule, err := c.GetHTTPErrorRule(0, "defaults", "", "")
	if err != nil {
		t.Fatalf("beta http error rule: %s", err)
	}
	if rule.Status != 429 || rule.ReturnContent != "/beta/429" {
		t.Errorf("beta http error rule status %d content %q, expected 429 /beta/429", rule.Status, rule.ReturnContent)
	}

	_, tcpCheck, err := c.GetTCPCheck(0, "defaults", "", "")
	if err != nil {
		t.Fatalf("beta tcp check: %s", err)
	}
	if tcpCheck.PortString != "80" {
		t.Errorf("beta tcp check port %q, expected %q", tcpCheck.PortString, "80")
	}
}

// TestDefaultsSelectorIsolationAcrossClients verifies that two independent
// clients keep their own defaults section selection instead of sharing a
// package global that the most recently parsed configuration overwrites.
func TestDefaultsSelectorIsolationAcrossClients(t *testing.T) {
	cAlpha, cBeta := newDefaultsSelectorPair(t)

	_, alpha, err := cAlpha.GetDefaultsConfiguration("")
	if err != nil {
		t.Fatalf("alpha defaults: %s", err)
	}
	if alpha.Mode != "http" {
		t.Errorf("alpha client selected mode %q, expected %q", alpha.Mode, "http")
	}

	_, beta, err := cBeta.GetDefaultsConfiguration("")
	if err != nil {
		t.Fatalf("beta defaults: %s", err)
	}
	if beta.Mode != "tcp" {
		t.Errorf("beta client selected mode %q, expected %q", beta.Mode, "tcp")
	}
}

// TestDefaultsSelectorRuleReadsIsolationAcrossClients exercises the migrated
// empty-parent rule reads (list and single) on the defaults parent after another
// client selected a different defaults section.
func TestDefaultsSelectorRuleReadsIsolationAcrossClients(t *testing.T) {
	cAlpha, cBeta := newDefaultsSelectorPair(t)

	assertAlphaRuleReads(t, cAlpha, "", "")
	assertBetaRuleReads(t, cBeta)
}

// TestDefaultsSelectorRuleMutationsIsolationAcrossClients performs create, edit
// and delete through the empty-parent defaults selection on cAlpha's transaction
// and proves the writes resolve to the transaction's own staged defaults section
// (gamma), not to cAlpha's master parser (which still selects alpha) nor to
// cBeta. gamma is seeded with the same rule families and counts as alpha so the
// edits and deletes address real entries, with distinguishable contents.
func TestDefaultsSelectorRuleMutationsIsolationAcrossClients(t *testing.T) {
	cAlpha := newDefaultsSelectorClient(t, writeDefaultsSelectorConfig(t, defaultsSelectorConfigAlpha))
	_, txID := startDefaultsSelectorTransaction(t, cAlpha)

	// Stage gamma in the transaction so the transaction selects gamma while its
	// master still selects alpha, and seed it with one rule entry per family
	// (matching alpha's counts) through explicit parent names.
	if err := cAlpha.CreateDefaultsSection(&models.Defaults{
		Name: "gamma", Mode: "tcp", Maxconn: misc.Int64P(300),
	}, txID, 0); err != nil {
		t.Fatalf("create gamma defaults in transaction: %s", err)
	}
	g0 := int64(0)
	g1 := int64(1)
	if err := cAlpha.CreateLogTarget("defaults", "gamma", &models.LogTarget{
		Index: &g0, Address: "127.0.0.1:616", Facility: "local0", Level: "notice",
	}, txID, 0); err != nil {
		t.Fatalf("seed gamma log target: %s", err)
	}
	if err := cAlpha.CreateHTTPCheck("defaults", "gamma", &models.HTTPCheck{
		Index: &g0, Type: "send", Method: "GET", Version: "HTTP/1.1", URI: "/gamma", CheckHeaders: []*models.ReturnHeader{},
	}, txID, 0); err != nil {
		t.Fatalf("seed gamma http check 0: %s", err)
	}
	if err := cAlpha.CreateHTTPCheck("defaults", "gamma", &models.HTTPCheck{
		Index: &g1, Type: "send-state",
	}, txID, 0); err != nil {
		t.Fatalf("seed gamma http check 1: %s", err)
	}
	if err := cAlpha.CreateHTTPErrorRule("defaults", "gamma", &models.HTTPErrorRule{
		Index: &g0, Type: "status", Status: 503, ReturnContentType: misc.StringP("application/json"),
		ReturnContentFormat: "file", ReturnContent: "/gamma/503",
	}, txID, 0); err != nil {
		t.Fatalf("seed gamma http error rule: %s", err)
	}
	if err := cAlpha.CreateTCPCheck("defaults", "gamma", &models.TCPCheck{
		Index: &g0, Action: "connect", PortString: "8443",
	}, txID, 0); err != nil {
		t.Fatalf("seed gamma tcp check: %s", err)
	}

	// Parse cBeta after the transaction staged gamma, so beta is the most
	// recently parsed configuration and any empty-parent leak to another parser
	// would select beta rather than the transaction's gamma.
	cBeta := newDefaultsSelectorClient(t, writeDefaultsSelectorConfig(t, defaultsSelectorConfigBeta))

	// Create through empty parent; the new entries must land in gamma.
	logID := int64(1)
	if err := cAlpha.CreateLogTarget("defaults", "", &models.LogTarget{
		Index: &logID, Address: "stdout", Facility: "daemon", Level: "notice",
	}, txID, 0); err != nil {
		t.Fatalf("create log target: %s", err)
	}
	if err := cAlpha.CreateHTTPCheck("defaults", "", &models.HTTPCheck{
		Index: &logID, Type: "send", Method: "GET", Version: "HTTP/1.1", URI: "/", CheckHeaders: []*models.ReturnHeader{},
	}, txID, 0); err != nil {
		t.Fatalf("create http check: %s", err)
	}
	if err := cAlpha.CreateHTTPErrorRule("defaults", "", &models.HTTPErrorRule{
		Index: &logID, Type: "status", Status: 429, ReturnContentType: misc.StringP("application/json"),
		ReturnContentFormat: "file", ReturnContent: "/gamma/429",
	}, txID, 0); err != nil {
		t.Fatalf("create http error rule: %s", err)
	}
	if err := cAlpha.CreateTCPCheck("defaults", "", &models.TCPCheck{
		Index: &logID, Action: "connect", PortString: "8444",
	}, txID, 0); err != nil {
		t.Fatalf("create tcp check: %s", err)
	}

	// Assert the writes against gamma explicitly, so a wrong empty-parent read
	// helper cannot mask a write that resolved to another section.
	if _, targets, err := cAlpha.GetLogTargets("defaults", "gamma", txID); err != nil || len(targets) != 2 {
		t.Fatalf("gamma log targets %d (%v), expected 2", len(targets), err)
	}
	if _, checks, err := cAlpha.GetHTTPChecks("defaults", "gamma", txID); err != nil || len(checks) != 3 {
		t.Fatalf("gamma http checks %d (%v), expected 3", len(checks), err)
	}
	if _, rules, err := cAlpha.GetHTTPErrorRules("defaults", "gamma", txID); err != nil || len(rules) != 2 {
		t.Fatalf("gamma http error rules %d (%v), expected 2", len(rules), err)
	}
	if _, tcpChecks, err := cAlpha.GetTCPChecks("defaults", "gamma", txID); err != nil || len(tcpChecks) != 2 {
		t.Fatalf("gamma tcp checks %d (%v), expected 2", len(tcpChecks), err)
	}

	// Edit the seeded entry 0 in each rule family through empty parent.
	if err := cAlpha.EditLogTarget(0, "defaults", "", &models.LogTarget{
		Index: misc.Int64P(0), Address: "127.0.0.1:9999", Facility: "local0", Level: "notice",
	}, txID, 0); err != nil {
		t.Fatalf("edit log target: %s", err)
	}
	if err := cAlpha.EditHTTPCheck(0, "defaults", "", &models.HTTPCheck{
		Index: misc.Int64P(0), Type: "send-state",
	}, txID, 0); err != nil {
		t.Fatalf("edit http check: %s", err)
	}
	if err := cAlpha.EditHTTPErrorRule(0, "defaults", "", &models.HTTPErrorRule{
		Index: misc.Int64P(0), Type: "status", Status: 503, ReturnContentType: misc.StringP("application/json"),
		ReturnContentFormat: "file", ReturnContent: "/gamma/edited",
	}, txID, 0); err != nil {
		t.Fatalf("edit http error rule: %s", err)
	}
	if err := cAlpha.EditTCPCheck(0, "defaults", "", &models.TCPCheck{
		Index: misc.Int64P(0), Action: "connect", PortString: "9443",
	}, txID, 0); err != nil {
		t.Fatalf("edit tcp check: %s", err)
	}

	if _, target, err := cAlpha.GetLogTarget(0, "defaults", "gamma", txID); err != nil || target.Address != "127.0.0.1:9999" {
		t.Fatalf("gamma log target %v (%v), expected edited address", target, err)
	}
	if _, rule, err := cAlpha.GetHTTPErrorRule(0, "defaults", "gamma", txID); err != nil || rule.ReturnContent != "/gamma/edited" {
		t.Fatalf("gamma http error rule %v (%v), expected edited content", rule, err)
	}
	if _, tcpCheck, err := cAlpha.GetTCPCheck(0, "defaults", "gamma", txID); err != nil || tcpCheck.PortString != "9443" {
		t.Fatalf("gamma tcp check %v (%v), expected edited port", tcpCheck, err)
	}

	// Delete the seeded entry 0 in each rule family through empty parent.
	if err := cAlpha.DeleteLogTarget(0, "defaults", "", txID, 0); err != nil {
		t.Fatalf("delete log target: %s", err)
	}
	if err := cAlpha.DeleteHTTPCheck(0, "defaults", "", txID, 0); err != nil {
		t.Fatalf("delete http check: %s", err)
	}
	if err := cAlpha.DeleteHTTPErrorRule(0, "defaults", "", txID, 0); err != nil {
		t.Fatalf("delete http error rule: %s", err)
	}
	if err := cAlpha.DeleteTCPCheck(0, "defaults", "", txID, 0); err != nil {
		t.Fatalf("delete tcp check: %s", err)
	}

	// gamma now holds only the created entries.
	if _, targets, err := cAlpha.GetLogTargets("defaults", "gamma", txID); err != nil || len(targets) != 1 || targets[0].Address != "stdout" {
		t.Fatalf("gamma log targets %+v (%v), expected only created stdout", targets, err)
	}
	if _, checks, err := cAlpha.GetHTTPChecks("defaults", "gamma", txID); err != nil || len(checks) != 2 {
		t.Fatalf("gamma http checks %d (%v), expected 2", len(checks), err)
	}
	if _, rules, err := cAlpha.GetHTTPErrorRules("defaults", "gamma", txID); err != nil || len(rules) != 1 || rules[0].Status != 429 {
		t.Fatalf("gamma http error rules %+v (%v), expected only created 429", rules, err)
	}
	if _, tcpChecks, err := cAlpha.GetTCPChecks("defaults", "gamma", txID); err != nil || len(tcpChecks) != 1 {
		t.Fatalf("gamma tcp checks %d (%v), expected 1", len(tcpChecks), err)
	}

	// The transaction copy of alpha, cAlpha's master parser and cBeta must all
	// be untouched by the empty-parent writes.
	assertAlphaRuleReads(t, cAlpha, "alpha", txID)
	assertAlphaRuleReads(t, cAlpha, "", "")
	assertBetaRuleReads(t, cBeta)
}

// TestDefaultsSelectorPushDefaultsIsolationAcrossClients verifies that
// PushDefaultsConfiguration edits the pushing client's own defaults section
// after another client selected a different one.
func TestDefaultsSelectorPushDefaultsIsolationAcrossClients(t *testing.T) {
	cAlpha, cBeta := newDefaultsSelectorPair(t)

	v, _, err := cAlpha.GetDefaultsConfiguration("")
	if err != nil {
		t.Fatalf("alpha defaults: %s", err)
	}
	if err := cAlpha.PushDefaultsConfiguration(&models.Defaults{
		Mode: "tcp", Maxconn: misc.Int64P(777),
	}, "", v); err != nil {
		t.Fatalf("push alpha defaults: %s", err)
	}

	_, alpha, err := cAlpha.GetDefaultsConfiguration("")
	if err != nil {
		t.Fatalf("alpha defaults after push: %s", err)
	}
	if alpha.Mode != "tcp" || alpha.Maxconn == nil || *alpha.Maxconn != 777 {
		t.Errorf("alpha defaults after push mode %q maxconn %v, expected tcp 777", alpha.Mode, alpha.Maxconn)
	}

	_, beta, err := cBeta.GetDefaultsConfiguration("")
	if err != nil {
		t.Fatalf("beta defaults: %s", err)
	}
	if beta.Mode != "tcp" || beta.Maxconn == nil || *beta.Maxconn != 200 {
		t.Errorf("beta defaults mode %q maxconn %v, expected tcp 200 (unchanged)", beta.Mode, beta.Maxconn)
	}
	assertBetaRuleReads(t, cBeta)
}

// TestDefaultsSelectorAbsentDefaultsEmptySuccess verifies the existing
// empty-success behavior for a configuration without a defaults section: no
// error and an empty Defaults value.
func TestDefaultsSelectorAbsentDefaultsEmptySuccess(t *testing.T) {
	c := newDefaultsSelectorClient(t, writeDefaultsSelectorConfig(t, `
# _version=1
global
    maxconn 100

frontend fe
    bind :8080
`))

	v, d, err := c.GetDefaultsConfiguration("")
	if err != nil {
		t.Fatalf("absent defaults returned error: %s", err)
	}
	if v != 1 {
		t.Errorf("version %d, expected 1", v)
	}
	if d == nil || d.Mode != "" {
		t.Errorf("absent defaults mode %q, expected empty defaults value", d.Mode)
	}
}

// TestDefaultsSelectorTransactionDifferentDefaultsFromMaster stages a gamma
// defaults section on a transaction so the transaction selector differs from its
// master parser (master alpha, staged gamma). It proves that reads and pushes
// resolve through the transaction parser, not the master, and that only the
// intended section changes.
func TestDefaultsSelectorTransactionDifferentDefaultsFromMaster(t *testing.T) {
	cAlpha, _ := newDefaultsSelectorPair(t)
	_, txID := startDefaultsSelectorTransaction(t, cAlpha)

	if err := cAlpha.CreateDefaultsSection(&models.Defaults{
		Name: "gamma", Mode: "tcp", Maxconn: misc.Int64P(300),
	}, txID, 0); err != nil {
		t.Fatalf("create gamma defaults in transaction: %s", err)
	}

	// The transaction now selects gamma; the master still selects alpha.
	_, staged, err := cAlpha.GetDefaultsConfiguration(txID)
	if err != nil {
		t.Fatalf("transaction defaults: %s", err)
	}
	if staged.Mode != "tcp" || staged.Maxconn == nil || *staged.Maxconn != 300 {
		t.Fatalf("transaction defaults mode %q maxconn %v, expected staged gamma tcp 300", staged.Mode, staged.Maxconn)
	}

	_, master, err := cAlpha.GetDefaultsConfiguration("")
	if err != nil {
		t.Fatalf("master defaults: %s", err)
	}
	if master.Mode != "http" || master.Maxconn == nil || *master.Maxconn != 100 {
		t.Fatalf("master defaults mode %q maxconn %v, expected alpha http 100", master.Mode, master.Maxconn)
	}

	// Empty-parent rule reads on the transaction must target staged gamma (no
	// rules), not master alpha (which has a log target).
	if _, targets, err := cAlpha.GetLogTargets("defaults", "", txID); err != nil {
		t.Fatalf("transaction log targets: %s", err)
	} else if len(targets) != 0 {
		t.Fatalf("transaction log targets %+v, expected none from staged gamma", targets)
	}
	assertAlphaRuleReads(t, cAlpha, "", "")

	// Push on the transaction must edit staged gamma only.
	if err := cAlpha.PushDefaultsConfiguration(&models.Defaults{
		Name: "gamma", Mode: "tcp", Maxconn: misc.Int64P(777),
	}, txID, 0); err != nil {
		t.Fatalf("push transaction defaults: %s", err)
	}

	_, staged, err = cAlpha.GetDefaultsConfiguration(txID)
	if err != nil {
		t.Fatalf("transaction defaults after push: %s", err)
	}
	if staged.Maxconn == nil || *staged.Maxconn != 777 {
		t.Errorf("transaction defaults maxconn %v, expected 777", staged.Maxconn)
	}

	_, master, err = cAlpha.GetDefaultsConfiguration("")
	if err != nil {
		t.Fatalf("master defaults after push: %s", err)
	}
	if master.Mode != "http" || master.Maxconn == nil || *master.Maxconn != 100 {
		t.Errorf("master defaults mode %q maxconn %v, expected unchanged alpha http 100", master.Mode, master.Maxconn)
	}
	assertAlphaRuleReads(t, cAlpha, "", "")
}
