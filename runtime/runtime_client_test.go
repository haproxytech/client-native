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

package runtime

import (
	"bufio"
	"errors"
	"net"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/Masterminds/semver/v3"
	native_errors "github.com/haproxytech/client-native/v6/errors"
	"github.com/haproxytech/client-native/v6/models"
	"github.com/haproxytech/client-native/v6/runtime/options"
	"github.com/stretchr/testify/assert"
)

type mapCommandSocket struct {
	net.Listener

	reply    func(string) string
	mu       sync.Mutex
	commands []string
}

func newMapCommandSocket(t *testing.T, reply func(string) string) *mapCommandSocket {
	t.Helper()
	f, err := os.CreateTemp("", "client-native-map-*")
	if err != nil {
		t.Fatal(err)
	}
	addr := f.Name()
	_ = f.Close()
	_ = os.Remove(addr)
	listener, err := net.Listen("unix", addr)
	if err != nil {
		t.Fatal(err)
	}
	socket := &mapCommandSocket{Listener: listener, reply: reply}
	go socket.serve()
	t.Cleanup(func() {
		_ = listener.Close()
		_ = os.Remove(addr)
	})
	return socket
}

func (s *mapCommandSocket) serve() {
	for {
		conn, err := s.Accept()
		if err != nil {
			return
		}
		go s.handleConnection(conn)
	}
}

func (s *mapCommandSocket) handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	command, err := reader.ReadString('\n')
	if err != nil {
		return
	}
	if strings.Contains(command, "<<") {
		for {
			line, readErr := reader.ReadString('\n')
			command += line
			if line == "\n" || readErr != nil {
				break
			}
		}
	}
	command = strings.TrimPrefix(command, "set severity-output number;")
	command = strings.TrimSuffix(command, "\n")

	s.mu.Lock()
	s.commands = append(s.commands, command)
	s.mu.Unlock()

	response := ""
	if s.reply != nil {
		response = s.reply(command)
	}
	if response == "" && strings.HasPrefix(command, "prepare map ") {
		response = "map transaction: version: 1"
	}
	_, _ = conn.Write([]byte(response))
}

func (s *mapCommandSocket) Commands() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]string(nil), s.commands...)
}

func newMapClient(t *testing.T, socket *mapCommandSocket, version string) *client {
	t.Helper()
	runtime := &SingleRuntime{}
	if err := runtime.Init(socket.Addr().String(), false, options.RuntimeOptions{DoNotCheckRuntimeOnInit: true}); err != nil {
		t.Fatal(err)
	}
	mapsDir := "/etc/haproxy/maps"
	previousVersion := haproxyVersion
	haproxyVersion = &HAProxyVersion{Version: semver.MustParse(version)}
	t.Cleanup(func() {
		haproxyVersion = previousVersion
	})
	return &client{
		options: options.RuntimeOptions{MapsDir: &mapsDir},
		runtime: runtime,
	}
}

func Test_parseMapPayload(t *testing.T) {
	type args struct {
		entries    models.MapEntries
		maxBufSize int
	}
	tests := []struct {
		name             string
		args             args
		wantExceededSize bool
		wantPayload      []string
	}{
		{
			name: "Return false and non empty slice when size < maxBufSize",
			args: args{entries: models.MapEntries{
				&models.MapEntry{Key: "k", Value: "v"},
			}, maxBufSize: 10},
			wantExceededSize: false,
			wantPayload:      []string{"k v\n"},
		},
		{
			name: "Return false and non empty slice when size == maxBufSize",
			args: args{entries: models.MapEntries{
				&models.MapEntry{Key: "key1", Value: "val1"},
				&models.MapEntry{Key: "key2", Value: "val2"},
			}, maxBufSize: 10},
			wantExceededSize: false,
			wantPayload: []string{
				"key1 val1\n",
				"key2 val2\n",
			},
		},
		{
			name: "Return true and non empty slice when size > maxBufSize",
			args: args{entries: models.MapEntries{
				&models.MapEntry{Key: "key1", Value: "val1"},
				&models.MapEntry{Key: "key2", Value: "val2"},
			}, maxBufSize: 5},
			wantExceededSize: true,
			wantPayload: []string{
				"key1 val1\n",
				"key2 val2\n",
			},
		},
		{
			name: "Keep final entry after crossing maxBufSize",
			args: args{entries: models.MapEntries{
				&models.MapEntry{Key: "a", Value: "bbbb"},
				&models.MapEntry{Key: "b", Value: "bbbb"},
			}, maxBufSize: 10},
			wantExceededSize: true,
			wantPayload: []string{
				"a bbbb\n",
				"b bbbb\n",
			},
		},
		{
			name: "Keep entries across repeated maxBufSize crossings",
			args: args{entries: models.MapEntries{
				&models.MapEntry{Key: "a", Value: "bbbb"},
				&models.MapEntry{Key: "b", Value: "bbbb"},
				&models.MapEntry{Key: "c", Value: "bbbb"},
			}, maxBufSize: 10},
			wantExceededSize: true,
			wantPayload: []string{
				"a bbbb\n",
				"b bbbb\n",
				"c bbbb\n",
			},
		},
		{
			name: "Do not repeat data after reaching maxBufSize exactly",
			args: args{entries: models.MapEntries{
				&models.MapEntry{Key: "a", Value: "b"},
				&models.MapEntry{Key: "c", Value: "ddd"},
			}, maxBufSize: 10},
			wantExceededSize: false,
			wantPayload: []string{
				"a b\nc ddd\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotExceededSize, gotPayload := parseMapPayload(tt.args.entries, tt.args.maxBufSize)
			if gotExceededSize != tt.wantExceededSize {
				t.Errorf("parseMapPayload() gotExceededSize = %v, want %v", gotExceededSize, tt.wantExceededSize)
			}
			if !assert.EqualValues(t, gotPayload, tt.wantPayload) {
				t.Errorf("parseMapPayload() gotPayload gotPayload = %v, want %v", gotPayload, tt.wantPayload)
			}
		})
	}
}

func TestReplaceMapPayload(t *testing.T) {
	socket := newMapCommandSocket(t, nil)
	client := newMapClient(t, socket, "2.4.0")
	entries := models.MapEntries{
		&models.MapEntry{Key: "first", Value: "one"},
		&models.MapEntry{Key: "second", Value: "two"},
	}

	if err := client.ReplaceMapPayload("example", entries); err != nil {
		t.Fatalf("ReplaceMapPayload() error = %v", err)
	}

	path := "/etc/haproxy/maps/example.map"
	assert.Equal(t, []string{
		"prepare map " + path,
		"clear map @1 " + path,
		"add map @1 " + path + " <<\nfirst one\nsecond two\n",
		"commit map @1 " + path,
	}, socket.Commands())
}

func TestReplaceMapPayloadMultipleChunksPreservesOrder(t *testing.T) {
	socket := newMapCommandSocket(t, nil)
	client := newMapClient(t, socket, "2.4.0")
	firstValue := strings.Repeat("a", maxBufSize-len("first")-2)
	secondValue := strings.Repeat("b", maxBufSize-len("second")-2)
	entries := models.MapEntries{
		&models.MapEntry{Key: "first", Value: firstValue},
		&models.MapEntry{Key: "second", Value: secondValue},
	}

	if err := client.ReplaceMapPayload("example", entries); err != nil {
		t.Fatalf("ReplaceMapPayload() error = %v", err)
	}

	path := "/etc/haproxy/maps/example.map"
	assert.Equal(t, []string{
		"prepare map " + path,
		"clear map @1 " + path,
		"add map @1 " + path + " <<\nfirst " + firstValue + "\n",
		"add map @1 " + path + " <<\nsecond " + secondValue + "\n",
		"commit map @1 " + path,
	}, socket.Commands())
}

func TestReplaceMapPayloadEmpty(t *testing.T) {
	socket := newMapCommandSocket(t, nil)
	client := newMapClient(t, socket, "2.4.0")

	if err := client.ReplaceMapPayload("example", nil); err != nil {
		t.Fatalf("ReplaceMapPayload() error = %v", err)
	}

	path := "/etc/haproxy/maps/example.map"
	assert.Equal(t, []string{
		"prepare map " + path,
		"clear map @1 " + path,
		"commit map @1 " + path,
	}, socket.Commands())
}

func TestReplaceMapPayloadUnsupported(t *testing.T) {
	socket := newMapCommandSocket(t, nil)
	client := newMapClient(t, socket, "2.3.0")

	err := client.ReplaceMapPayload("example", nil)
	if !errors.Is(err, native_errors.ErrGeneral) {
		t.Fatalf("ReplaceMapPayload() error = %v, want %v", err, native_errors.ErrGeneral)
	}
	assert.Empty(t, socket.Commands())
}

func TestReplaceMapPayloadDoesNotCommitAfterAddFailure(t *testing.T) {
	socket := newMapCommandSocket(t, func(command string) string {
		if strings.HasPrefix(command, "add map ") {
			return "[3]: add failed"
		}
		return ""
	})
	client := newMapClient(t, socket, "2.4.0")

	err := client.ReplaceMapPayload("example", models.MapEntries{
		&models.MapEntry{Key: "first", Value: "one"},
	})
	if !errors.Is(err, native_errors.ErrGeneral) {
		t.Fatalf("ReplaceMapPayload() error = %v, want %v", err, native_errors.ErrGeneral)
	}

	path := "/etc/haproxy/maps/example.map"
	assert.Equal(t, []string{
		"prepare map " + path,
		"clear map @1 " + path,
		"add map @1 " + path + " <<\nfirst one\n",
	}, socket.Commands())
}
