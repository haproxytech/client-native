// Copyright 2025 HAProxy Technologies
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

package runtime

import (
	"net"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/haproxytech/client-native/v6/runtime/options"
)

// captureSocket is a minimal Unix-socket server that records every byte
// a client writes after immediately sending a fixed reply. It exists so
// these tests can assert the exact bytes client-native sends to HAProxy's
// runtime API for heredoc commands.
type captureSocket struct {
	listener net.Listener
	reply    string
	mu       sync.Mutex
	received []byte
	addr     string
}

func newCaptureSocket(t *testing.T, reply string) *captureSocket {
	t.Helper()
	f, err := os.CreateTemp("", "client-native-cap-*")
	if err != nil {
		t.Fatalf("CreateTemp: %v", err)
	}
	addr := f.Name()
	_ = f.Close()
	_ = os.Remove(addr)
	l, err := net.Listen("unix", addr)
	if err != nil {
		t.Fatalf("Listen: %v", err)
	}
	s := &captureSocket{listener: l, reply: reply, addr: addr}
	go s.serve()
	t.Cleanup(func() { _ = l.Close(); _ = os.Remove(addr) })
	return s
}

func (s *captureSocket) serve() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return
		}
		go func(c net.Conn) {
			defer c.Close()
			// Send reply right away so the client's read loop sees data
			// and unblocks; afterwards the client gets EOF when we close
			// at the deferred Close, which lets readFromSocket return.
			_, _ = c.Write([]byte(s.reply))
			// Drain everything the client sends for up to 500ms.
			_ = c.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
			buf := make([]byte, 4096)
			var got []byte
			for {
				n, err := c.Read(buf)
				if n > 0 {
					got = append(got, buf[:n]...)
				}
				if err != nil {
					break
				}
			}
			s.mu.Lock()
			s.received = append(s.received, got...)
			s.mu.Unlock()
		}(conn)
	}
}

func (s *captureSocket) Received() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return string(s.received)
}

// assertHeredocTerminated checks that the bytes sent by the runtime client
// contain a `<<\n` heredoc start AND that the payload after it is followed
// by a blank line (\n\n) BEFORE any `;` command separator (which would
// otherwise be interpreted as part of the heredoc body and stall HAProxy).
func assertHeredocTerminated(t *testing.T, sent string) {
	t.Helper()
	idx := strings.Index(sent, "<<\n")
	if idx < 0 {
		t.Fatalf("no `<<\\n` heredoc start in sent bytes:\n%s", quoted(sent))
	}
	tail := sent[idx+len("<<\n"):]
	blankLine := strings.Index(tail, "\n\n")
	if blankLine < 0 {
		t.Fatalf("heredoc payload is never followed by a blank line; HAProxy will hang.\n  sent: %s", quoted(sent))
	}
	if semi := strings.Index(tail, ";"); semi >= 0 && semi < blankLine {
		t.Fatalf("`;` appears before the blank-line terminator at offset %d (blank line at %d).\n  Once the framing in runtime_single_client.go appends `;quit\\n`, HAProxy will interpret `;quit` as part of the heredoc body and stall.\n  sent: %s",
			semi, blankLine, quoted(sent))
	}
}

func quoted(s string) string { return `"` + strings.ReplaceAll(s, "\n", `\n`) + `"` }

// TestHeredocPayloadAlwaysTerminated drives each runtime command that
// builds a `<<` heredoc with a payload that explicitly does NOT end in
// "\n", in both master-worker and single-process framing modes, and
// asserts the bytes actually written to the socket are properly
// terminated. Catches the class of bug where a missing payload `\n`
// combined with the master-worker `;quit\n` framing leaves the heredoc
// without a blank-line terminator and HAProxy stalls until the socket
// deadline fires.
func TestHeredocPayloadAlwaysTerminated(t *testing.T) {
	const noTrailingNewline = "-----BEGIN CERTIFICATE-----\nMIIB...\n-----END CERTIFICATE-----"

	commands := []struct {
		name string
		call func(s *SingleRuntime) error
	}{
		{
			name: "AddMapPayload",
			call: func(s *SingleRuntime) error { return s.AddMapPayload("/etc/map", "k1 v1\nk2 v2") },
		},
		{
			name: "AddMapPayloadVersioned",
			call: func(s *SingleRuntime) error { return s.AddMapPayloadVersioned("1", "/etc/map", "k1 v1\nk2 v2") },
		},
	}

	for _, mode := range []struct {
		name             string
		masterWorkerMode bool
	}{
		{"master-worker", true},
		{"single-process", false},
	} {
		for _, cmd := range commands {
			t.Run(mode.name+"/"+cmd.name, func(t *testing.T) {
				mock := newCaptureSocket(t, " Transaction created for certificate /etc/cert.pem!\n\n")

				s := &SingleRuntime{}
				if err := s.Init(mock.addr, mode.masterWorkerMode, options.RuntimeOptions{DoNotCheckRuntimeOnInit: true}); err != nil {
					t.Fatalf("Init: %v", err)
				}
				// Return value of the call is not checked: the mock doesn't
				// emit a real success response for every command shape, so
				// some calls legitimately error. What we care about is the
				// bytes that hit the wire before that.
				_ = cmd.call(s)

				// Give the serve goroutine a moment to drain.
				time.Sleep(50 * time.Millisecond)

				assertHeredocTerminated(t, mock.Received())
			})
		}
	}
}
