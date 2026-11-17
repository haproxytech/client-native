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

	"github.com/stretchr/testify/assert"
)

func newStrictRuntime(t *testing.T, response string) *SingleRuntime {
	t.Helper()
	socket := newMapCommandSocket(t, func(string) string { return response })
	return newMapClient(t, socket, "3.4.0").runtime
}

func TestExecuteStrict(t *testing.T) {
	tests := []struct {
		name       string
		response   string
		okWarnings []string
		wantErr    string
	}{
		{name: "empty response", response: ""},
		{name: "plain output", response: "done"},
		{name: "allowed warning", response: "[4]: nothing changed", okWarnings: []string{"nothing changed"}},
		{name: "other warning", response: "[4]: health checks are not enabled on this server", okWarnings: []string{"nothing changed"}, wantErr: "[4]  health checks are not enabled on this server [cmd]"},
		{name: "warning without allow list", response: "[4]: nothing changed", wantErr: "nothing changed"},
		{name: "error", response: "[3]: No such server.", okWarnings: []string{"No such server"}, wantErr: "[3]  No such server. [cmd]"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := newStrictRuntime(t, tt.response).executeStrict("cmd", tt.okWarnings...)
			if tt.wantErr == "" {
				assert.NoError(t, err)
				return
			}
			assert.ErrorContains(t, err, tt.wantErr)
		})
	}
	assert.ErrorIs(t, (&SingleRuntime{}).executeStrict("a;b"), ErrRuntimeInvalidChar)
}

func TestSetServerAddr(t *testing.T) {
	tests := []struct {
		name     string
		ip       string
		port     int
		response string
		wantErr  string
	}{
		{name: "changed", ip: "127.0.0.5", port: 9011, response: "[4]: IP changed from '127.0.0.1' to '127.0.0.5', port changed from '9001' to '9011' by 'stats socket command'"},
		{name: "nothing changed", ip: "127.0.0.5", port: 9011, response: "[4]: nothing changed"},
		{name: "no need to change (2.x)", ip: "127.0.0.5", port: 9011, response: "[4]: no need to change the addr"},
		{name: "empty", ip: "127.0.0.5", response: ""},
		{name: "invalid addr", ip: "not-an-ip", port: 80, response: "[4]: Invalid addr 'not-an-ip'", wantErr: "Invalid addr 'not-an-ip'"},
		{name: "invalid port", ip: "127.0.0.6", response: "[4]: problem converting port 'abc' to an int", wantErr: "problem converting port 'abc' to an int"},
		{name: "error", ip: "127.0.0.6", response: "[3]: No such server.", wantErr: "No such server."},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := newStrictRuntime(t, tt.response).SetServerAddr("be", "s1", tt.ip, tt.port)
			if tt.wantErr == "" {
				assert.NoError(t, err)
				return
			}
			assert.ErrorContains(t, err, tt.wantErr)
		})
	}
}

func TestSetServerCheckPort(t *testing.T) {
	tests := []struct {
		name     string
		response string
		wantErr  string
	}{
		{name: "ok", response: ""},
		{name: "checks disabled", response: "[4]: health checks are not enabled on this server", wantErr: "health checks are not enabled on this server"},
		{name: "not an integer", response: "[4]: provided port is not an integer", wantErr: "provided port is not an integer"},
		{name: "error", response: "[3]: No such server.", wantErr: "No such server."},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := newStrictRuntime(t, tt.response).SetServerCheckPort("be", "s1", 9012)
			if tt.wantErr == "" {
				assert.NoError(t, err)
				return
			}
			assert.ErrorContains(t, err, tt.wantErr)
		})
	}
}

func TestSetServerAgentAddr(t *testing.T) {
	tests := []struct {
		name     string
		response string
		wantErr  string
	}{
		{name: "ok", response: ""},
		{name: "agent disabled", response: "[4]: agent checks are not enabled on this server", wantErr: "agent checks are not enabled on this server"},
		{name: "error", response: "[3]: No such server.", wantErr: "No such server."},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := newStrictRuntime(t, tt.response).SetServerAgentAddr("be", "s1", "127.0.0.7")
			if tt.wantErr == "" {
				assert.NoError(t, err)
				return
			}
			assert.ErrorContains(t, err, tt.wantErr)
		})
	}
}

func TestExecuteStrictLeavesExecuteUnchanged(t *testing.T) {
	assert.NoError(t, newStrictRuntime(t, "[4]: some warning").SetServerWeight("be", "s1", "10"))
}
