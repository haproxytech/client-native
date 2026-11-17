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
//

package configuration

import (
	"context"
	"errors"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	parser "github.com/haproxytech/client-native/v6/config-parser"
	"github.com/haproxytech/client-native/v6/configuration/options"
	"github.com/haproxytech/client-native/v6/runtime"
)

// fakeRuntime is a non-nil runtime; calling any of its methods panics.
type fakeRuntime struct{ runtime.Runtime }

// setRegistry replaces the package registry for the duration of the test.
func setRegistry(t *testing.T, reg runtimeRegistry) {
	t.Helper()
	saved := registry
	registry = reg
	t.Cleanup(func() { registry = saved })
}

// commitRecorder collects the results passed to onCommit.
type commitRecorder []options.CommitResult

func (r *commitRecorder) onCommit(res options.CommitResult) { *r = append(*r, res) }

func (r *commitRecorder) one(t *testing.T) options.CommitResult {
	t.Helper()
	require.Len(t, *r, 1)
	return (*r)[0]
}

// newCommitClient builds a client on a configuration file holding cfg.
func newCommitClient(t testing.TB, cfg string, opts ...options.ConfigurationOption) *client {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "haproxy.cfg")
	require.NoError(t, os.WriteFile(path, []byte(cfg), 0o600))
	opts = append([]options.ConfigurationOption{
		options.ConfigurationFile(path), options.HAProxyBin("echo"), options.TransactionsDir(filepath.Join(dir, "transactions")), options.UsePersistentTransactions,
	}, opts...)
	conf, err := New(context.Background(), opts...)
	require.NoError(t, err)
	c, ok := conf.(*client)
	require.True(t, ok)
	return c
}

// newApplyClient builds a client with runtime apply enabled; rt is what the
// option's rt func returns.
func newApplyClient(t *testing.T, cfg string, rt runtime.Runtime, rec *commitRecorder, opts ...options.ConfigurationOption) *client {
	t.Helper()
	return newCommitClient(t, cfg, append(opts, options.RuntimeApply(func() runtime.Runtime { return rt }, rec.onCommit))...)
}

func currentVersion(t *testing.T, c *client) int64 {
	t.Helper()
	v, err := c.GetVersion("")
	require.NoError(t, err)
	return v
}

// postRaw pushes cfg with the current version.
func postRaw(t *testing.T, c *client, cfg string) {
	t.Helper()
	require.NoError(t, c.PostRawConfiguration(&cfg, currentVersion(t, c), false))
}

// setMaxconn sets the maxconn of frontend fe in transactionID, or with an
// implicit commit when transactionID is empty.
func setMaxconn(t *testing.T, c *client, transactionID string, maxconn int64) {
	t.Helper()
	_, fe, err := c.GetFrontend("fe", transactionID)
	require.NoError(t, err)
	fe.Maxconn = &maxconn
	version := int64(0)
	if transactionID == "" {
		version = currentVersion(t, c)
	}
	require.NoError(t, c.EditFrontend("fe", fe, transactionID, version))
}

// maxconnRegistry is testRuntimeRegistry with the frontend maxconn entry
// replaced by apply.
func maxconnRegistry(t *testing.T, l *stepLog, apply fieldApplier) runtimeRegistry {
	t.Helper()
	reg := testRuntimeRegistry(t, l)
	reg[parser.Frontends]["maxconn"] = entry{apply: apply}
	return reg
}

var maxconnChanged = strings.Replace(runtimeApplyBase, "maxconn 50", "maxconn 60", 1) //nolint:gochecknoglobals

func TestPlanRuntimeApplyDisabled(t *testing.T) {
	var l stepLog
	setRegistry(t, testRuntimeRegistry(t, &l))
	c := newCommitClient(t, runtimeApplyBase)
	postRaw(t, c, maxconnChanged)
	require.Nil(t, c.planRuntimeApply("t", c.ConfigurationFile, 1))
	require.Empty(t, l)
}

func TestRuntimeCommitNilRuntime(t *testing.T) {
	setRegistry(t, maxconnRegistry(t, nil, func(change) ([]step, bool) {
		t.Error("nothing may be compared when rt returns nil")
		return nil, false
	}))
	var rec commitRecorder
	c := newApplyClient(t, runtimeApplyBase, nil, &rec)
	version := currentVersion(t, c)
	tr, err := c.StartTransaction(version)
	require.NoError(t, err)
	setMaxconn(t, c, tr.ID, 60)
	_, err = c.CommitTransaction(tr.ID)
	require.NoError(t, err)
	res := rec.one(t)
	require.True(t, res.Reload)
	require.NoError(t, res.Err)
	require.Equal(t, tr.ID, res.TransactionID)
	require.Equal(t, version, res.Version)
}

func TestRuntimeCommitNoChange(t *testing.T) {
	var l stepLog
	setRegistry(t, testRuntimeRegistry(t, &l))
	var rec commitRecorder
	c := newApplyClient(t, runtimeApplyBase, &fakeRuntime{}, &rec)
	tr, err := c.StartTransaction(currentVersion(t, c))
	require.NoError(t, err)
	_, err = c.CommitTransaction(tr.ID)
	require.NoError(t, err)
	res := rec.one(t)
	require.False(t, res.Reload)
	require.NoError(t, res.Err)
	require.Empty(t, l)
}

func TestRuntimeCommitCommentOnly(t *testing.T) {
	var l stepLog
	setRegistry(t, testRuntimeRegistry(t, &l))
	var rec commitRecorder
	c := newApplyClient(t, runtimeApplyBase, &fakeRuntime{}, &rec)
	postRaw(t, c, strings.Replace(runtimeApplyBase, "# comment", "# other comment", 1))
	require.False(t, rec.one(t).Reload)
	require.Empty(t, l)
}

func TestRuntimeCommitApplied(t *testing.T) {
	tests := map[string]func(t *testing.T, c *client){
		"explicit transaction": func(t *testing.T, c *client) {
			t.Helper()
			tr, err := c.StartTransaction(currentVersion(t, c))
			require.NoError(t, err)
			setMaxconn(t, c, tr.ID, 60)
			_, err = c.CommitTransaction(tr.ID)
			require.NoError(t, err)
		},
		"implicit commit": func(t *testing.T, c *client) {
			t.Helper()
			setMaxconn(t, c, "", 60)
		},
		"raw configuration": func(t *testing.T, c *client) {
			t.Helper()
			postRaw(t, c, maxconnChanged)
		},
	}
	for name, commit := range tests {
		t.Run(name, func(t *testing.T) {
			var l stepLog
			var got []runtime.Runtime
			setRegistry(t, maxconnRegistry(t, &l, func(c change) ([]step, bool) {
				return []step{func(rt runtime.Runtime) error {
					l = append(l, "maxconn "+c.section)
					got = append(got, rt)
					return nil
				}}, true
			}))
			var rec commitRecorder
			rt := &fakeRuntime{}
			c := newApplyClient(t, runtimeApplyBase, rt, &rec)
			commit(t, c)
			res := rec.one(t)
			require.False(t, res.Reload)
			require.NoError(t, res.Err)
			require.Equal(t, stepLog{"maxconn fe"}, l)
			require.Len(t, got, 1)
			require.Same(t, rt, got[0])
		})
	}
}

func TestRuntimeCommitNoRegistryEntry(t *testing.T) {
	var l stepLog
	setRegistry(t, testRuntimeRegistry(t, &l))
	var rec commitRecorder
	c := newApplyClient(t, runtimeApplyBase, &fakeRuntime{}, &rec)
	postRaw(t, c, strings.Replace(runtimeApplyBase, "balance roundrobin", "balance leastconn", 1))
	res := rec.one(t)
	require.True(t, res.Reload)
	require.NoError(t, res.Err)
	require.Empty(t, l)
}

func TestRuntimeCommitStepError(t *testing.T) {
	errStep := errors.New("step failed")
	var l stepLog
	setRegistry(t, maxconnRegistry(t, &l, func(change) ([]step, bool) {
		return []step{func(runtime.Runtime) error { return errStep }, l.step("later")}, true
	}))
	var rec commitRecorder
	c := newApplyClient(t, runtimeApplyBase, &fakeRuntime{}, &rec)
	postRaw(t, c, maxconnChanged)
	res := rec.one(t)
	require.True(t, res.Reload)
	require.ErrorIs(t, res.Err, errStep)
	require.Empty(t, l)
	data, err := os.ReadFile(c.ConfigurationFile)
	require.NoError(t, err)
	require.Contains(t, string(data), "maxconn 60")
}

func TestRuntimeCommitStepPanic(t *testing.T) {
	var l stepLog
	setRegistry(t, maxconnRegistry(t, &l, func(change) ([]step, bool) {
		return []step{func(runtime.Runtime) error { panic("boom") }, l.step("later")}, true
	}))
	var rec commitRecorder
	c := newApplyClient(t, runtimeApplyBase, &fakeRuntime{}, &rec)
	postRaw(t, c, maxconnChanged)
	res := rec.one(t)
	require.True(t, res.Reload)
	require.Error(t, res.Err)
	require.Empty(t, l)
}

func TestRuntimeCommitPlanningPanic(t *testing.T) {
	setRegistry(t, maxconnRegistry(t, nil, func(change) ([]step, bool) { panic("boom") }))
	var rec commitRecorder
	c := newApplyClient(t, runtimeApplyBase, &fakeRuntime{}, &rec)
	postRaw(t, c, maxconnChanged)
	res := rec.one(t)
	require.True(t, res.Reload)
	require.Error(t, res.Err)
}

func TestRuntimeCommitFailed(t *testing.T) {
	var l stepLog
	setRegistry(t, testRuntimeRegistry(t, &l))
	var rec commitRecorder
	c := newApplyClient(t, runtimeApplyBase, &fakeRuntime{}, &rec)
	tr, err := c.StartTransaction(currentVersion(t, c))
	require.NoError(t, err)
	postRaw(t, c, maxconnChanged)
	require.Len(t, rec, 1)
	_, err = c.CommitTransaction(tr.ID)
	require.Error(t, err)
	require.Len(t, rec, 1)
}

func TestRuntimeCommitStepSeesNewConfiguration(t *testing.T) {
	var c *client
	var before int64
	ran := false
	setRegistry(t, maxconnRegistry(t, nil, func(change) ([]step, bool) {
		return []step{func(runtime.Runtime) error {
			ran = true
			data, err := os.ReadFile(c.ConfigurationFile)
			require.NoError(t, err)
			require.Contains(t, string(data), "maxconn 60")
			require.Equal(t, before+1, currentVersion(t, c))
			return nil
		}}, true
	}))
	var rec commitRecorder
	c = newApplyClient(t, runtimeApplyBase, &fakeRuntime{}, &rec)
	before = currentVersion(t, c)
	postRaw(t, c, maxconnChanged)
	require.True(t, ran)
	require.False(t, rec.one(t).Reload)
}

func TestRuntimeCommitFileHasSectionBlocks(t *testing.T) {
	for name, opts := range map[string][]options.ConfigurationOption{
		"plain": nil,
		"md5":   {options.UseMd5Hash},
	} {
		t.Run(name, func(t *testing.T) {
			var l stepLog
			setRegistry(t, testRuntimeRegistry(t, &l))
			var rec commitRecorder
			c := newApplyClient(t, runtimeApplyBase, &fakeRuntime{}, &rec, opts...)
			postRaw(t, c, maxconnChanged)
			require.False(t, rec.one(t).Reload)
			data, err := os.ReadFile(c.ConfigurationFile)
			require.NoError(t, err)
			blocks, ok := sectionBlocks(string(data))
			require.True(t, ok, "configuration file not recognised:\n%s", data)
			require.ElementsMatch(t, []string{"global", "defaults d", "frontend fe", "backend be"}, slices.Collect(maps.Keys(blocks)))
		})
	}
}

func TestSameSectionsInFiles(t *testing.T) {
	dir := t.TempDir()
	oldPath, newPath := filepath.Join(dir, "old.cfg"), filepath.Join(dir, "new.cfg")
	require.NoError(t, newRuntimeApplyParser(t, runtimeApplyBase).Save(oldPath))
	require.NoError(t, newRuntimeApplyParser(t, maxconnChanged).Save(newPath))
	same := sameSectionsInFiles(oldPath, newPath)
	require.ElementsMatch(t, []string{"global", "defaults d", "backend be"}, slices.Collect(maps.Keys(same)))
	require.Nil(t, sameSectionsInFiles(filepath.Join(dir, "missing.cfg"), newPath))
	require.Nil(t, sameSectionsInFiles(oldPath, filepath.Join(dir, "missing.cfg")))
}

func TestRuntimeCommitNonCanonicalFile(t *testing.T) {
	var l stepLog
	setRegistry(t, testRuntimeRegistry(t, &l))
	lines := strings.Split(genConfig(20, 2), "\n")
	for i, line := range lines {
		lines[i] = strings.TrimLeft(line, " ")
	}
	base := strings.Join(lines, "\n")
	_, ok := sectionBlocks(base)
	require.False(t, ok, "base must not have the writer layout")

	cfg := strings.Replace(base, "10.0.5.0:8080 check weight 10", "10.0.5.0:8080 check weight 20", 1)
	var want stepLog
	steps, ok := runtimeSteps(testRuntimeRegistry(t, &want), newRuntimeApplyParser(t, base), newRuntimeApplyParser(t, cfg), nil)
	require.True(t, ok)
	runSteps(t, steps)
	require.Equal(t, stepLog{"weight s0"}, want)

	var rec commitRecorder
	c := newApplyClient(t, base, &fakeRuntime{}, &rec)
	postRaw(t, c, cfg)
	require.False(t, rec.one(t).Reload)
	require.Equal(t, want, l)

	data, err := os.ReadFile(c.ConfigurationFile)
	require.NoError(t, err)
	_, ok = sectionBlocks(string(data))
	require.True(t, ok, "committed file must have the writer layout")
	l, rec = nil, nil
	postRaw(t, c, strings.Replace(string(data), "10.0.7.1:8080 check weight 10", "10.0.7.1:8080 check weight 30", 1))
	require.False(t, rec.one(t).Reload)
	require.Equal(t, stepLog{"weight s1"}, l)
}
