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

package runtime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type commandCase struct {
	name     string
	version  string
	call     func(c *client) error
	response string
	wantCmd  []string
	wantErr  string
}

func runCommandCases(t *testing.T, tests []commandCase) {
	t.Helper()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			socket := newMapCommandSocket(t, func(string) string { return tt.response })
			version := tt.version
			if version == "" {
				version = "3.4.0"
			}
			err := tt.call(newMapClient(t, socket, version))
			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.wantErr)
			}
			if tt.wantCmd == nil {
				assert.Empty(t, socket.Commands())
			} else {
				assert.Equal(t, tt.wantCmd, socket.Commands())
			}
		})
	}
}

func TestSetServerMaxconn(t *testing.T) {
	cmd := []string{"set maxconn server be/s1 100"}
	call := func(c *client) error { return c.SetServerMaxconn("be", "s1", 100) }
	runCommandCases(t, []commandCase{
		{name: "ok", call: call, wantCmd: cmd},
		{name: "zero", call: func(c *client) error { return c.SetServerMaxconn("be", "s1", 0) }, wantCmd: []string{"set maxconn server be/s1 0"}},
		{name: "no such server", call: call, response: "[3]: No such server.", wantCmd: cmd, wantErr: "No such server."},
		{name: "negative", call: func(c *client) error { return c.SetServerMaxconn("be", "s1", -1) }, wantErr: "bad request"},
	})
}

func TestSetServerAgentPort(t *testing.T) {
	cmd := []string{"set server be/s1 agent-port 5000"}
	call := func(c *client) error { return c.SetServerAgentPort("be", "s1", 5000) }
	runCommandCases(t, []commandCase{
		{name: "ok", call: call, wantCmd: cmd},
		{name: "agent disabled", call: call, response: "[4]: agent checks are not enabled on this server", wantCmd: cmd, wantErr: "agent checks are not enabled on this server"},
		{name: "no such server", call: call, response: "[3]: No such server.", wantCmd: cmd, wantErr: "No such server."},
		{name: "port zero", call: func(c *client) error { return c.SetServerAgentPort("be", "s1", 0) }, wantErr: "bad request"},
		{name: "port too big", call: func(c *client) error { return c.SetServerAgentPort("be", "s1", 65536) }, wantErr: "bad request"},
	})
}

func TestSetServerCheckAddr(t *testing.T) {
	cmd := []string{"set server be/s1 check-addr 127.0.0.2"}
	call := func(c *client) error { return c.SetServerCheckAddr("be", "s1", "127.0.0.2", 0) }
	runCommandCases(t, []commandCase{
		{name: "ok without port", call: call, wantCmd: cmd},
		{name: "ok with port", call: func(c *client) error { return c.SetServerCheckAddr("be", "s1", "127.0.0.2", 8080) }, wantCmd: []string{"set server be/s1 check-addr 127.0.0.2 port 8080"}},
		{name: "checks disabled", call: call, response: "[4]: health checks are not enabled on this server", wantCmd: cmd, wantErr: "health checks are not enabled on this server"},
		{name: "invalid addr", call: call, response: "[4]: invalid addr '127.0.0.2'", wantCmd: cmd, wantErr: "invalid addr '127.0.0.2'"},
		{name: "no such server", call: call, response: "[3]: No such server.", wantCmd: cmd, wantErr: "No such server."},
	})
}

func TestSetServerFqdn(t *testing.T) {
	cmd := []string{"set server be/s1 fqdn b.example.com"}
	call := func(c *client) error { return c.SetServerFqdn("be", "s1", "b.example.com") }
	runCommandCases(t, []commandCase{
		{name: "changed", call: call, response: "[4]: be/s1 changed its FQDN from a.example.com to b.example.com", wantCmd: cmd},
		{name: "no need to change", call: call, response: "[4]: no need to change the FDQN", wantCmd: cmd},
		{name: "invalid fqdn", call: call, response: "[4]: invalid fqdn 'b.example.com'", wantCmd: cmd, wantErr: "invalid fqdn 'b.example.com'"},
		{name: "could not update", call: call, response: "[4]: be/s1 changed its FQDN from a.example.com to b.example.com\ncould not update be/s1 FQDN", wantCmd: cmd, wantErr: "could not update be/s1 FQDN"},
		{name: "no resolution", call: call, response: "[3]: set server be/s1 fqdn failed because no resolution is configured.", wantCmd: cmd, wantErr: "no resolution is configured"},
		{name: "empty", call: call, wantCmd: cmd},
	})
}

func TestSetGlobalMaxConn(t *testing.T) {
	cmd := []string{"set maxconn global 1000"}
	call := func(c *client) error { return c.SetGlobalMaxConn(1000) }
	runCommandCases(t, []commandCase{
		{name: "ok", call: call, wantCmd: cmd},
		{name: "out of range", call: call, response: "[3]: Value out of range.", wantCmd: cmd, wantErr: "Value out of range."},
		{name: "negative", call: func(c *client) error { return c.SetGlobalMaxConn(-1) }, wantErr: "bad request"},
	})
}

func TestSetRateLimitGlobal(t *testing.T) {
	runCommandCases(t, []commandCase{
		{name: "connections", call: func(c *client) error { return c.SetRateLimitConnectionsGlobal(100) }, wantCmd: []string{"set rate-limit connections global 100"}},
		{name: "sessions", call: func(c *client) error { return c.SetRateLimitSessionsGlobal(200) }, wantCmd: []string{"set rate-limit sessions global 200"}},
		{name: "http-compression", call: func(c *client) error { return c.SetRateLimitHTTPCompressionGlobal(300) }, wantCmd: []string{"set rate-limit http-compression global 300"}},
		{name: "error", call: func(c *client) error { return c.SetRateLimitSessionsGlobal(200) }, response: "[3]: Value out of range.", wantCmd: []string{"set rate-limit sessions global 200"}, wantErr: "Value out of range."},
	})
}

func TestAddBackend(t *testing.T) {
	cmd := []string{"add backend dyn1 from def"}
	call := func(c *client) error { return c.AddBackend("dyn1", "def", "", "") }
	runCommandCases(t, []commandCase{
		{name: "ok", call: call, response: "[6]: New backend registered.", wantCmd: cmd},
		{name: "ok with mode and guid", call: func(c *client) error { return c.AddBackend("dyn1", "def", "http", "g1") }, response: "[6]: New backend registered.", wantCmd: []string{"add backend dyn1 from def mode http guid g1"}},
		{name: "ok with mode", call: func(c *client) error { return c.AddBackend("dyn1", "def", "tcp", "") }, wantCmd: []string{"add backend dyn1 from def mode tcp"}},
		{name: "ok with guid", call: func(c *client) error { return c.AddBackend("dyn1", "def", "", "g1") }, wantCmd: []string{"add backend dyn1 from def guid g1"}},
		{name: "name used", call: call, response: "[3]: Error during backend creation : name is already used by other proxy 'backend dyn1'.", wantCmd: cmd, wantErr: "name is already used by other proxy 'backend dyn1'."},
		{name: "no defaults", call: call, response: "[3]: Cannot find default proxy 'nope'.", wantCmd: cmd, wantErr: "Cannot find default proxy 'nope'."},
		{name: "empty name", call: func(c *client) error { return c.AddBackend("", "def", "", "") }, wantErr: "bad request"},
		{name: "empty defaults", call: func(c *client) error { return c.AddBackend("dyn1", "", "", "") }, wantErr: "bad request"},
		{name: "3.3 rejected", version: "3.3.0", call: call, wantErr: "requires HAProxy 3.4"},
	})
}

func TestDeleteBackend(t *testing.T) {
	cmd := []string{"del backend dyn1"}
	call := func(c *client) error { return c.DeleteBackend("dyn1") }
	runCommandCases(t, []commandCase{
		{name: "ok", call: call, response: "[6]: Backend deleted.", wantCmd: cmd},
		{name: "has servers", call: call, response: "[3]: Only a backend without server can be deleted.", wantCmd: cmd, wantErr: "Only a backend without server can be deleted."},
		{name: "empty name", call: func(c *client) error { return c.DeleteBackend("") }, wantErr: "bad request"},
		{name: "3.3 rejected", version: "3.3.0", call: call, wantErr: "requires HAProxy 3.4"},
	})
}

func TestPublishBackend(t *testing.T) {
	cmd := []string{"publish backend dyn1"}
	call := func(c *client) error { return c.PublishBackend("dyn1") }
	runCommandCases(t, []commandCase{
		{name: "ok", call: call, response: "[6]: Backend published.", wantCmd: cmd},
		{name: "error", call: call, response: "[3]: No such backend.", wantCmd: cmd, wantErr: "No such backend."},
		{name: "empty name", call: func(c *client) error { return c.PublishBackend("") }, wantErr: "bad request"},
		{name: "3.3 rejected", version: "3.3.0", call: call, wantErr: "requires HAProxy 3.4"},
	})
}

func TestUnpublishBackend(t *testing.T) {
	cmd := []string{"unpublish backend dyn1"}
	call := func(c *client) error { return c.UnpublishBackend("dyn1") }
	runCommandCases(t, []commandCase{
		{name: "ok", call: call, response: "[6]: Backend unpublished.", wantCmd: cmd},
		{name: "error", call: call, response: "[3]: No such backend.", wantCmd: cmd, wantErr: "No such backend."},
		{name: "empty name", call: func(c *client) error { return c.UnpublishBackend("") }, wantErr: "bad request"},
		{name: "3.3 rejected", version: "3.3.0", call: call, wantErr: "requires HAProxy 3.4"},
	})
}

func TestWaitBackendRemovable(t *testing.T) {
	cmd := []string{"wait 2000ms be-removable dyn1"}
	call := func(c *client) error { return c.WaitBackendRemovable("dyn1", 2*time.Second) }
	runCommandCases(t, []commandCase{
		{name: "ok", call: call, response: "[6]: Done.", wantCmd: cmd},
		{name: "has servers", call: call, response: "[3]: Failed. Only a backend without server can be deleted.", wantCmd: cmd, wantErr: "Failed. Only a backend without server can be deleted."},
		{name: "sub millisecond", call: func(c *client) error { return c.WaitBackendRemovable("dyn1", time.Microsecond) }, wantErr: "bad request"},
		{name: "empty name", call: func(c *client) error { return c.WaitBackendRemovable("", time.Second) }, wantErr: "bad request"},
		{name: "3.3 rejected", version: "3.3.0", call: call, wantErr: "requires HAProxy 3.4"},
	})
}
