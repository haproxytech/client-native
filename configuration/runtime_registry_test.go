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
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	parser "github.com/haproxytech/client-native/v6/config-parser"
	"github.com/haproxytech/client-native/v6/configuration/options"
	"github.com/haproxytech/client-native/v6/models"
	"github.com/haproxytech/client-native/v6/runtime"
)

var errRecorded = errors.New("recorded runtime failure")

// recordingRuntime records one line per runtime call; a call whose line
// starts with failOn returns an error.
type recordingRuntime struct {
	runtime.Runtime
	calls  []string
	failOn string
}

func (r *recordingRuntime) record(format string, args ...any) error {
	line := fmt.Sprintf(format, args...)
	r.calls = append(r.calls, line)
	if r.failOn != "" && strings.HasPrefix(line, r.failOn) {
		return fmt.Errorf("%w: %s", errRecorded, line)
	}
	return nil
}

func (r *recordingRuntime) SetFrontendMaxConn(frontend string, maxconn int) error {
	return r.record("SetFrontendMaxConn %s %d", frontend, maxconn)
}

func (r *recordingRuntime) AddServer(backend, name, attributes string) error {
	return r.record("AddServer %s %s %s", backend, name, attributes)
}

func (r *recordingRuntime) DeleteServer(backend, name string) error {
	return r.record("DeleteServer %s %s", backend, name)
}

func (r *recordingRuntime) SetServerAddr(backend, server, ip string, port int) error {
	return r.record("SetServerAddr %s %s %s %d", backend, server, ip, port)
}

func (r *recordingRuntime) SetServerState(backend, server, state string) error {
	return r.record("SetServerState %s %s %s", backend, server, state)
}

func (r *recordingRuntime) SetServerWeight(backend, server, weight string) error {
	return r.record("SetServerWeight %s %s %s", backend, server, weight)
}

func (r *recordingRuntime) SetServerCheckPort(backend, server string, port int) error {
	return r.record("SetServerCheckPort %s %s %d", backend, server, port)
}

func (r *recordingRuntime) SetServerAgentAddr(backend, server, addr string) error {
	return r.record("SetServerAgentAddr %s %s %s", backend, server, addr)
}

func (r *recordingRuntime) SetServerAgentSend(backend, server, send string) error {
	return r.record("SetServerAgentSend %s %s %s", backend, server, send)
}

func (r *recordingRuntime) EnableAgentCheck(backend, server string) error {
	return r.record("EnableAgentCheck %s %s", backend, server)
}

func (r *recordingRuntime) DisableAgentCheck(backend, server string) error {
	return r.record("DisableAgentCheck %s %s", backend, server)
}

func (r *recordingRuntime) EnableServerHealth(backend, server string) error {
	return r.record("EnableServerHealth %s %s", backend, server)
}

func (r *recordingRuntime) ShutdownSessionsServer(backend, server string) error {
	return r.record("ShutdownSessionsServer %s %s", backend, server)
}

func (r *recordingRuntime) WaitServerRemovable(backend, server string, timeout time.Duration) error {
	return r.record("WaitServerRemovable %s %s %v", backend, server, timeout)
}

// registryCalls returns the runtime calls the package registry makes to go
// from oldCfg to newCfg; nil means reload.
func registryCalls(t *testing.T, oldCfg, newCfg string) []string {
	t.Helper()
	steps, ok := runtimeSteps(registry, newRuntimeApplyParser(t, oldCfg), newRuntimeApplyParser(t, newCfg), nil)
	if !ok {
		return nil
	}
	rt := &recordingRuntime{}
	for _, s := range steps {
		require.NoError(t, s(rt))
	}
	if rt.calls == nil {
		return []string{}
	}
	return rt.calls
}

const registryBase = `defaults d
  mode http

frontend fe from d
  maxconn 50
  bind :80
  default_backend be

backend be from d
  server s1 10.0.0.1:80 weight 10
  server s2 10.0.0.2:80
`

var reloads []string // nil: the change needs a reload

func TestRuntimeRegistry(t *testing.T) {
	s1 := "server s1 10.0.0.1:80 weight 10"
	s2 := "server s2 10.0.0.2:80"
	rep := func(cfg, old, cur string) string { return strings.Replace(cfg, old, cur, 1) }
	r := func(old, cur string) string { return rep(registryBase, old, cur) }
	both := func(s1New, s2New string) string { return rep(r(s1, s1New), s2, s2New) }
	withS2 := func(s2Line string) string { return r(s2, s2Line) }
	del := []string{
		"SetServerState be s2 maint", "ShutdownSessionsServer be s2",
		"WaitServerRemovable be s2 500ms", "DeleteServer be s2",
	}
	tests := []struct {
		name, old, cur string
		want           []string
	}{
		{"no change", registryBase, registryBase, []string{}},
		{"weight changed", registryBase, r(s1, "server s1 10.0.0.1:80 weight 20"), []string{"SetServerWeight be s1 20"}},
		{
			"weight removed inherits default-server",
			r("backend be from d\n", "backend be from d\n  default-server weight 5\n"),
			rep(r("backend be from d\n", "backend be from d\n  default-server weight 5\n"), s1, "server s1 10.0.0.1:80"),
			[]string{"SetServerWeight be s1 5"},
		},
		{"weight removed nothing inherited", registryBase, r(s1, "server s1 10.0.0.1:80"), reloads},
		{"address changed", registryBase, r(s1, "server s1 10.0.0.9:80 weight 10"), []string{"SetServerAddr be s1 10.0.0.9 80"}},
		{"port changed", registryBase, r(s1, "server s1 10.0.0.1:81 weight 10"), []string{"SetServerAddr be s1 10.0.0.1 81"}},
		{"host name", registryBase, r(s1, "server s1 example.com:80 weight 10"), reloads},
		{"host name to address", r(s1, "server s1 example.com:80 weight 10"), registryBase, reloads},
		{"port removed", registryBase, r(s1, "server s1 10.0.0.1 weight 10"), reloads},
		{"ipv6", withS2("server s2 [::1]:80"), withS2("server s2 [::2]:8080"), []string{"SetServerAddr be s2 ::2 8080"}},
		{"disabled added", registryBase, withS2(s2 + " disabled"), []string{"SetServerState be s2 maint"}},
		{"disabled removed", withS2(s2 + " disabled"), registryBase, []string{"SetServerState be s2 ready"}},
		{"disabled to enabled", withS2(s2 + " disabled"), withS2(s2 + " enabled"), []string{"SetServerState be s2 ready"}},
		{"enabled added", registryBase, withS2(s2 + " enabled"), []string{"SetServerState be s2 ready"}},
		{"check port", withS2(s2 + " check port 9000"), withS2(s2 + " check port 9001"), []string{"SetServerCheckPort be s2 9001"}},
		{"agent-check added", withS2(s2 + " agent-port 5000"), withS2(s2 + " agent-check agent-port 5000"), []string{"EnableAgentCheck be s2"}},
		{"agent-check removed", withS2(s2 + " agent-check agent-port 5000"), withS2(s2 + " agent-port 5000"), []string{"DisableAgentCheck be s2"}},
		{"agent-addr", withS2(s2 + " agent-addr 10.0.0.5"), withS2(s2 + " agent-addr 10.0.0.6"), []string{"SetServerAgentAddr be s2 10.0.0.6"}},
		{"agent-send plain", withS2(s2 + " agent-send hello"), withS2(s2 + " agent-send bye"), []string{"SetServerAgentSend be s2 bye"}},
		{"agent-send quoted", withS2(s2 + " agent-send hello"), withS2(s2 + ` agent-send "a b"`), reloads},
		{"server maxconn", registryBase, withS2(s2 + " maxconn 10"), reloads},
		{"server ssl", registryBase, withS2(s2 + " ssl"), reloads},
		{"add plain", registryBase, registryBase + "  server s3 10.0.0.3:80\n", []string{"AddServer be s3 10.0.0.3:80", "SetServerState be s3 ready"}},
		{
			"add with check", registryBase, registryBase + "  server s3 10.0.0.3:80 check\n",
			[]string{"AddServer be s3 10.0.0.3:80 check", "EnableServerHealth be s3", "SetServerState be s3 ready"},
		},
		{
			"add with agent-check", registryBase, registryBase + "  server s3 10.0.0.3:80 agent-check agent-port 5000\n",
			[]string{"AddServer be s3 10.0.0.3:80 agent-check agent-port 5000", "EnableAgentCheck be s3", "SetServerState be s3 ready"},
		},
		{"add disabled", registryBase, registryBase + "  server s3 10.0.0.3:80 disabled\n", []string{"AddServer be s3 10.0.0.3:80 disabled"}},
		{"add host name", registryBase, registryBase + "  server s3 example.com:80\n", reloads},
		{"add quoted", registryBase, registryBase + "  server s3 10.0.0.3:80 agent-send \"a b\"\n", reloads},
		{"add before existing", registryBase, r(s2, "server s3 10.0.0.3:80\n  "+s2), reloads},
		{"delete", registryBase, r(s2+"\n", ""), del},
		{
			"delete with agent-check", withS2(s2 + " agent-check agent-port 5000"), r(s2+"\n", ""),
			append([]string{"DisableAgentCheck be s2"}, del...),
		},
		{
			"replace server", registryBase, r(s2, "server s3 10.0.0.3:80"),
			append(append([]string{}, del...), "AddServer be s3 10.0.0.3:80", "SetServerState be s3 ready"),
		},
		{
			"mixed", registryBase, rep(r(s2+"\n", ""), s1, "server s1 10.0.0.1:80 weight 20\n  server s3 10.0.0.3:80"),
			append(append([]string{}, del...), "SetServerWeight be s1 20", "AddServer be s3 10.0.0.3:80", "SetServerState be s3 ready"),
		},
		{"both servers", registryBase, both("server s1 10.0.0.1:80 weight 20", s2+" disabled"), []string{"SetServerWeight be s1 20", "SetServerState be s2 maint"}},
		{"frontend maxconn", registryBase, r("maxconn 50", "maxconn 100"), []string{"SetFrontendMaxConn fe 100"}},
		{"frontend maxconn zero", registryBase, r("maxconn 50", "maxconn 0"), reloads},
		{"frontend maxconn removed", registryBase, r("  maxconn 50\n", ""), reloads},
		{
			"default-server changed",
			r("backend be from d\n", "backend be from d\n  default-server weight 5\n"),
			r("backend be from d\n", "backend be from d\n  default-server weight 6\n"),
			reloads,
		},
		{
			"server-template changed",
			registryBase + "  server-template srv 1-3 10.0.1.1:80\n",
			registryBase + "  server-template srv 1-4 10.0.1.1:80\n",
			reloads,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.want, registryCalls(t, tc.old, tc.cur))
		})
	}
}

func TestRuntimeRegistryAddServerInherited(t *testing.T) {
	tests := []struct {
		name, cfg, server, want string
		after                   []string // calls after AddServer; nil = only the ready step
	}{
		{
			name:   "backend default-server",
			cfg:    "defaults d\n  mode http\n\nbackend be from d\n  default-server inter 3s weight 5\n  server s1 10.0.0.1:80\n",
			server: "server s3 10.0.0.3:80",
			want:   "10.0.0.3:80 inter 3s weight 5",
		},
		{
			name:   "defaults default-server",
			cfg:    "defaults d\n  mode http\n  default-server inter 3s\n\nbackend be from d\n  server s1 10.0.0.1:80\n",
			server: "server s3 10.0.0.3:80",
			want:   "10.0.0.3:80 inter 3s",
		},
		{
			name: "two-level chain",
			cfg: "defaults a\n  default-server inter 1s\n\ndefaults b from a\n  default-server rise 3\n\n" +
				"backend be from b\n  default-server fall 2\n  server s1 10.0.0.1:80\n",
			server: "server s3 10.0.0.3:80 weight 7",
			want:   "10.0.0.3:80 inter 1s rise 3 fall 2 weight 7",
		},
		{
			name: "nearest wins",
			cfg: "defaults a\n  default-server weight 1 inter 1s\n\ndefaults b from a\n  default-server weight 2\n\n" +
				"backend be from b\n  default-server weight 3\n  server s1 10.0.0.1:80\n",
			server: "server s3 10.0.0.3:80",
			want:   "10.0.0.3:80 inter 1s weight 3",
		},
		{
			name: "server wins",
			cfg: "defaults a\n  default-server weight 1 inter 1s\n\ndefaults b from a\n  default-server weight 2\n\n" +
				"backend be from b\n  default-server weight 3\n  server s1 10.0.0.1:80\n",
			server: "server s3 10.0.0.3:80 weight 9",
			want:   "10.0.0.3:80 inter 1s weight 9",
		},
		{
			name:   "inherited check, no-check on server",
			cfg:    "defaults d\n  mode http\n\nbackend be from d\n  default-server check\n  server s1 10.0.0.1:80\n",
			server: "server s3 10.0.0.3:80 no-check",
			want:   "10.0.0.3:80",
		},
		{
			name:   "inherited no-check, check on server",
			cfg:    "defaults d\n  mode http\n  default-server no-check\n\nbackend be from d\n  server s1 10.0.0.1:80\n",
			server: "server s3 10.0.0.3:80 check",
			want:   "10.0.0.3:80 check",
			after:  []string{"EnableServerHealth be s3", "SetServerState be s3 ready"},
		},
		{
			name:   "enabled on server",
			cfg:    "defaults d\n  mode http\n\nbackend be from d\n  server s1 10.0.0.1:80\n",
			server: "server s3 10.0.0.3:80 weight 5 enabled",
			want:   "10.0.0.3:80 weight 5",
		},
		{
			name:   "inherited disabled, enabled on server",
			cfg:    "defaults d\n  mode http\n\nbackend be from d\n  default-server disabled\n  server s1 10.0.0.1:80\n",
			server: "server s3 10.0.0.3:80 enabled",
			want:   "10.0.0.3:80",
		},
		{
			name:   "inherited enabled, disabled on server",
			cfg:    "defaults d\n  mode http\n\nbackend be from d\n  default-server enabled\n  server s1 10.0.0.1:80\n",
			server: "server s3 10.0.0.3:80 disabled",
			want:   "10.0.0.3:80 disabled",
			after:  []string{},
		},
		{
			name:   "inherited agent-check, no-agent-check on server",
			cfg:    "defaults d\n  mode http\n\nbackend be from d\n  default-server agent-check agent-port 5000\n  server s1 10.0.0.1:80\n",
			server: "server s3 10.0.0.3:80 no-agent-check",
			want:   "10.0.0.3:80 agent-port 5000",
		},
		{
			name:   "other no- keyword passed through",
			cfg:    "defaults d\n  mode http\n\nbackend be from d\n  server s1 10.0.0.1:80\n",
			server: "server s3 10.0.0.3:80 no-ssl",
			want:   "10.0.0.3:80 no-ssl",
		},
		{
			name:   "inherited check",
			cfg:    "defaults d\n  mode http\n\nbackend be from d\n  default-server check\n  server s1 10.0.0.1:80\n",
			server: "server s3 10.0.0.3:80",
			want:   "10.0.0.3:80 check",
			after:  []string{"EnableServerHealth be s3", "SetServerState be s3 ready"},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			after := tc.after
			if after == nil {
				after = []string{"SetServerState be s3 ready"}
			}
			want := append([]string{"AddServer be s3 " + tc.want}, after...)
			require.Equal(t, want, registryCalls(t, tc.cfg, tc.cfg+"  "+tc.server+"\n"))
		})
	}
}

func TestServerItems(t *testing.T) {
	cfg := `defaults a
  default-server check weight 1

defaults b from a
  default-server inter 2s

backend be from b
  default-server no-check
  default-server inter 3s
  server s1 10.0.0.1:80 check
  server s2 10.0.0.2:80 weight 4

backend empty from b
  mode http
`
	p := newRuntimeApplyParser(t, cfg)
	items, err := serverItems(p, parser.Backends, "be")
	require.NoError(t, err)
	require.Equal(t, []item{
		{name: "s1", options: []option{{"weight", "1"}, {"no-check", ""}, {"inter", "3s"}, {"address", "10.0.0.1:80"}, {"check", ""}}},
		{name: "s2", options: []option{{"check", ""}, {"no-check", ""}, {"inter", "3s"}, {"address", "10.0.0.2:80"}, {"weight", "4"}}},
	}, items)

	items, err = serverItems(p, parser.Backends, "empty")
	require.NoError(t, err)
	require.Empty(t, items)
}

func TestServerItemsDefaultsCycle(t *testing.T) {
	cfg := "defaults a from b\n  default-server weight 1\n\ndefaults b from a\n  default-server inter 1s\n\n" +
		"backend be from b\n  server s1 10.0.0.1:80\n"
	items, err := serverItems(newRuntimeApplyParser(t, cfg), parser.Backends, "be")
	require.NoError(t, err)
	require.Equal(t, []item{{name: "s1", options: []option{{"weight", "1"}, {"inter", "1s"}, {"address", "10.0.0.1:80"}}}}, items)
}

func TestHostPort(t *testing.T) {
	tests := []struct {
		addr string
		ip   string
		port int
		ok   bool
	}{
		{"10.0.0.1", "10.0.0.1", 0, true},
		{"10.0.0.1:80", "10.0.0.1", 80, true},
		{"10.0.0.1:0", "10.0.0.1", 0, true},
		{"[::1]:80", "::1", 80, true},
		{"[::1]", "", 0, false},
		{"::1", "", 0, false},
		{"example.com", "", 0, false},
		{"example.com:80", "", 0, false},
		{"unix@/run/s.sock", "", 0, false},
		{"/run/s.sock", "", 0, false},
		{"ipv4@10.0.0.1:80", "", 0, false},
		{"10.0.0.1:+80", "", 0, false},
		{"10.0.0.1:-80", "", 0, false},
		{"10.0.0.1:80-90", "", 0, false},
		{"10.0.0.1:", "", 0, false},
		{"10.0.0.1:70000", "", 0, false},
		{"", "", 0, false},
	}
	for _, tc := range tests {
		t.Run(tc.addr, func(t *testing.T) {
			ip, port, ok := hostPort(tc.addr)
			require.Equal(t, tc.ok, ok)
			require.Equal(t, tc.ip, ip)
			require.Equal(t, tc.port, port)
		})
	}
}

func TestRuntimeRegistryCommitError(t *testing.T) {
	var rec commitRecorder
	rt := &recordingRuntime{failOn: "SetServerState"}
	c := newApplyClient(t, registryBase, rt, &rec)
	postRaw(t, c, strings.Replace(registryBase, "  server s2 10.0.0.2:80\n", "", 1))
	res := rec.one(t)
	require.True(t, res.Reload)
	require.ErrorIs(t, res.Err, errRecorded)
	require.Equal(t, []string{"SetServerState be s2 maint"}, rt.calls)
}

func TestRuntimeRegistryClient(t *testing.T) {
	tests := map[string]struct {
		do   func(t *testing.T, c *client)
		want []string
	}{
		"edit weight": {
			do: func(t *testing.T, c *client) {
				t.Helper()
				_, srv, err := c.GetServer("s1", BackendParentName, "be", "")
				require.NoError(t, err)
				weight := int64(20)
				srv.Weight = &weight
				require.NoError(t, c.EditServer("s1", BackendParentName, "be", srv, "", currentVersion(t, c)))
			},
			want: []string{"SetServerWeight be s1 20"},
		},
		"create": {
			do: func(t *testing.T, c *client) {
				t.Helper()
				port := int64(80)
				srv := &models.Server{Name: "s3", Address: "10.0.0.3", Port: &port, ServerParams: models.ServerParams{Check: "enabled"}}
				require.NoError(t, c.CreateServer(BackendParentName, "be", srv, "", currentVersion(t, c)))
			},
			want: []string{"AddServer be s3 10.0.0.3:80 check", "EnableServerHealth be s3", "SetServerState be s3 ready"},
		},
		"create maintenance disabled": {
			do: func(t *testing.T, c *client) {
				t.Helper()
				port := int64(80)
				srv := &models.Server{Name: "s3", Address: "10.0.0.3", Port: &port, ServerParams: models.ServerParams{Maintenance: "disabled"}}
				require.NoError(t, c.CreateServer(BackendParentName, "be", srv, "", currentVersion(t, c)))
				data, err := os.ReadFile(c.ConfigurationFile)
				require.NoError(t, err)
				require.Contains(t, string(data), "server s3 10.0.0.3:80 enabled")
			},
			want: []string{"AddServer be s3 10.0.0.3:80", "SetServerState be s3 ready"},
		},
		"delete": {
			do: func(t *testing.T, c *client) {
				t.Helper()
				require.NoError(t, c.DeleteServer("s2", BackendParentName, "be", "", currentVersion(t, c)))
			},
			want: []string{
				"SetServerState be s2 maint", "ShutdownSessionsServer be s2",
				"WaitServerRemovable be s2 500ms", "DeleteServer be s2",
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var rec commitRecorder
			rt := &recordingRuntime{}
			c := newCommitClient(t, registryBase, options.RuntimeApply(func() runtime.Runtime { return rt }, rec.onCommit))
			tc.do(t, c)
			res := rec.one(t)
			require.False(t, res.Reload)
			require.NoError(t, res.Err)
			require.Equal(t, tc.want, rt.calls)
		})
	}
}
