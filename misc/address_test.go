// Copyright 2020 HAProxy Technologies
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
package misc

import (
	"testing"
)

func TestParseAddress(t *testing.T) {
	tests := []struct {
		name        string
		address     string
		requirePort bool
		want        string
	}{
		{
			name:        "IPv4 valid address without a port",
			address:     "192.168.0.1",
			requirePort: false,
			want:        "192.168.0.1",
		},
		{
			name:        "IPv4 valid address without a port that is required",
			address:     "192.168.0.2",
			requirePort: true,
			want:        "",
		},
		{
			name:        "IPv4 valid address with a port",
			address:     "192.168.0.3:80",
			requirePort: false,
			want:        "192.168.0.3:80",
		},
		{
			name:        "IPv4 valid address with a port range",
			address:     "192.168.0.3:80-81",
			requirePort: false,
			want:        "192.168.0.3:80-81",
		},
		{
			name:        "IPv4 valid address with a port that is required",
			address:     "192.168.0.4:80",
			requirePort: true,
			want:        "192.168.0.4:80",
		},
		{
			name:        "IPv4 invalid address without a port",
			address:     "192.168.0.1024",
			requirePort: false,
			want:        "",
		},
		{
			name:        "IPv4 invalid address without a port that is required",
			address:     "192.168.0.1024",
			requirePort: true,
			want:        "",
		},
		{
			name:        "Stateless Auto IP (IPv6)",
			address:     "0:0:0:0:0:0:0:0",
			requirePort: false,
			want:        "0:0:0:0:0:0:0:0",
		},
		{
			name:        "Basic (IPv6)",
			address:     "::FFFF:C0A8:1",
			requirePort: false,
			want:        "::FFFF:C0A8:1",
		},
		{
			name:        "Leading zeros (IPv6)",
			address:     "::FFFF:C0A8:0001",
			requirePort: false,
			want:        "::FFFF:C0A8:0001",
		},
		{
			name:        "Basic (IPv6)",
			address:     "0000:0000:0000:0000:0000:FFFF:C0A8:1",
			requirePort: false,
			want:        "0000:0000:0000:0000:0000:FFFF:C0A8:1",
		},
		{
			name:        "IPv4 literal (IPv6)",
			address:     "::FFFF:192.168.0.1",
			requirePort: false,
			want:        "::FFFF:192.168.0.1",
		},
		{
			name:        "With port info (IPv6)",
			address:     "[::FFFF:C0A8:1]:80",
			requirePort: false,
			want:        "[::FFFF:C0A8:1]:80",
		},
		{
			name:        "With port range (IPv6)",
			address:     "[::FFFF:C0A8:1]:80-90",
			requirePort: false,
			want:        "[::FFFF:C0A8:1]:80-90",
		},
		{
			name:        "IPv4 literal (IPv6)",
			address:     "[::FFFF:C0A6:1%1]:80",
			requirePort: false,
			want:        "[::FFFF:C0A6:1%1]:80",
		},
		{
			name:        "[::] (IPv6)",
			address:     "[::]",
			requirePort: false,
			want:        "[::]",
		},
		{
			name:        "*:80 (IPv4)",
			address:     "*:80",
			requirePort: false,
			want:        "*:80",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, _ := ParseAddress(test.address, test.requirePort); got != test.want {
				t.Errorf("ParseAddress() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestIsPrefixed(t *testing.T) {
	tests := []struct {
		name    string
		address string
		want    bool
	}{
		{
			name:    "ipv4@ prefix",
			address: "ipv4@192.168.0.1",
			want:    true,
		},
		{
			name:    "ipv6@ prefix",
			address: "ipv6@192.168.0.1",
			want:    true,
		},
		{
			name:    "no prefix",
			address: "192.168.0.2",
			want:    false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := IsPrefixed(test.address); got != test.want {
				t.Errorf("IsPrefixed() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestIsSocket(t *testing.T) {
	tests := []struct {
		name    string
		address string
		want    bool
	}{
		{
			name:    "socket",
			address: "/tmp/proxy.socket",
			want:    true,
		},
		{
			name:    "not a socket",
			address: "192.168.0.2",
			want:    false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := IsSocket(test.address); got != test.want {
				t.Errorf("IsSocket() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestHasNoPath(t *testing.T) {
	tests := []struct {
		name    string
		address string
		want    bool
	}{
		{
			name:    "has no path",
			address: "*",
			want:    true,
		},
		{
			name:    "has a path",
			address: "192.168.0.2",
			want:    false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := HasNoPath(test.address); got != test.want {
				t.Errorf("HasNoPath() = %v, want %v", got, test.want)
			}
		})
	}
}
