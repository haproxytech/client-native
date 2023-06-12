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
	"fmt"
	"reflect"
	"testing"

	"github.com/haproxytech/client-native/v4/misc"
	"github.com/haproxytech/client-native/v4/models"
)

func TestGetServers(t *testing.T) { //nolint:gocognit,gocyclo
	v, servers, err := clientTest.GetServers("backend", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if len(servers) != 2 {
		t.Errorf("%v servers returned, expected 2", len(servers))
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	for _, s := range servers {
		if s.Name != "webserv" && s.Name != "webserv2" {
			t.Errorf("Expected only webserv or webserv2 servers, %v found", s.Name)
		}
		if s.Address != "192.168.1.1" {
			t.Errorf("%v: Address not 192.168.1.1: %v", s.Name, s.Address)
		}
		if *s.Port != 9300 && *s.Port != 9200 {
			t.Errorf("%v: Port not 9300 or 9200: %v", s.Name, *s.Port)
		}
		if s.Ssl != "enabled" {
			t.Errorf("%v: Ssl not enabled: %v", s.Name, s.Ssl)
		}
		if s.Cookie != "BLAH" {
			t.Errorf("%v: Cookie not BLAH: %v", s.Name, s.Cookie)
		}
		if *s.Maxconn != 1000 {
			t.Errorf("%v: Maxconn not 1000: %v", s.Name, *s.Maxconn)
		}
		if *s.Weight != 10 {
			t.Errorf("%v: Weight not 10: %v", s.Name, *s.Weight)
		}
		if *s.Inter != 2000 {
			t.Errorf("%v: Inter not 2000: %v", s.Name, *s.Inter)
		}
		if s.Ws != "h1" {
			t.Errorf("%v: Ws not h1: %v", s.Name, s.Ws)
		}
		if *s.PoolLowConn != 128 {
			t.Errorf("%v: PoolLowConn not 128: %v", s.Name, *s.PoolLowConn)
		}
		if len(s.ProxyV2Options) != 2 {
			t.Errorf("%v: ProxyV2Options < 2: %d", s.Name, len(s.ProxyV2Options))
		} else {
			if s.ProxyV2Options[0] != "authority" {
				t.Errorf("%v: ProxyV2Options[0] not authority: %s", s.Name, s.ProxyV2Options[0])
			}
			if s.ProxyV2Options[1] != "crc32c" {
				t.Errorf("%v: ProxyV2Options[0] not crc32c: %s", s.Name, s.ProxyV2Options[1])
			}
		}
	}

	_, servers, err = clientTest.GetServers("backend", "test_2", "")
	if err != nil {
		t.Error(err.Error())
	}
	if len(servers) > 0 {
		t.Errorf("%v servers returned, expected 0", len(servers))
	}
}

func TestGetServer(t *testing.T) {
	v, s, err := clientTest.GetServer("webserv", "backend", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	if s.Name != "webserv" {
		t.Errorf("Expected only webserv, %v found", s.Name)
	}
	if s.Address != "192.168.1.1" {
		t.Errorf("%v: Address not 192.168.1.1: %v", s.Name, s.Address)
	}
	if *s.Port != 9200 {
		t.Errorf("%v: Port not 9200: %v", s.Name, *s.Port)
	}
	if s.Ssl != "enabled" {
		t.Errorf("%v: Ssl not enabled: %v", s.Name, s.Ssl)
	}
	if s.Cookie != "BLAH" {
		t.Errorf("%v: HTTPCookieID not BLAH: %v", s.Name, s.Cookie)
	}
	if *s.Maxconn != 1000 {
		t.Errorf("%v: MaxConnections not 1000: %v", s.Name, *s.Maxconn)
	}
	if *s.Weight != 10 {
		t.Errorf("%v: Weight not 10: %v", s.Name, *s.Weight)
	}
	if *s.Inter != 2000 {
		t.Errorf("%v: Inter not 2000: %v", s.Name, *s.Inter)
	}
	if s.Ws != "h1" {
		t.Errorf("%v: Ws not h1: %v", s.Name, s.Ws)
	}
	if *s.PoolLowConn != 128 {
		t.Errorf("%v: PoolLowConn not 128: %v", s.Name, *s.PoolLowConn)
	}
	if *s.ID != 1234 {
		t.Errorf("%v: ID not 1234: %v", s.Name, *s.ID)
	}
	if len(s.ProxyV2Options) != 2 {
		t.Errorf("%v: ProxyV2Options < 2: %d", s.Name, len(s.ProxyV2Options))
	} else {
		if s.ProxyV2Options[0] != "authority" {
			t.Errorf("%v: ProxyV2Options[0] not authority: %s", s.Name, s.ProxyV2Options[0])
		}
		if s.ProxyV2Options[1] != "crc32c" {
			t.Errorf("%v: ProxyV2Options[0] not crc32c: %s", s.Name, s.ProxyV2Options[1])
		}
	}
	if *s.PoolPurgeDelay != 10000 {
		t.Errorf("%v: PoolPurgeDelay not 10000: %v", s.Name, *s.PoolPurgeDelay)
	}

	_, err = s.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	_, _, err = clientTest.GetServer("webserv", "backend", "test_2", "")
	if err == nil {
		t.Error("Should throw error, non existant server")
	}
}

func TestGetRingServer(t *testing.T) {
	v, s, err := clientTest.GetServer("mysyslogsrv", "ring", "myring", "")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	if s.Name != "mysyslogsrv" {
		t.Errorf("Expected only mysyslogsrv, %v found", s.Name)
	}
	if s.Address != "127.0.0.1" {
		t.Errorf("%v: Address not 127.0.0.1: %v", s.Name, s.Address)
	}
	if *s.Port != 6514 {
		t.Errorf("%v: Port not 6514: %v", s.Name, *s.Port)
	}

	if s.LogProto != "octet-count" {
		t.Errorf("%v: log-proto not octet-count: %v", s.Name, s.LogProto)
	}

	v, s, err = clientTest.GetServer("s1", "ring", "myring", "")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	if s.Name != "s1" {
		t.Errorf("Expected only s1, %v found", s.Name)
	}
	if s.Address != "192.168.1.1" {
		t.Errorf("%v: Address not 192.168.1.1: %v", s.Name, s.Address)
	}
	if *s.Port != 80 {
		t.Errorf("%v: Port not 6514: %v", s.Name, *s.Port)
	}

	if s.ResolveOpts != "allow-dup-ip,ignore-weight" {
		t.Errorf("%v: resolve_opts not allow-dup-ip,ignore-weight: %v", s.Name, s.ResolveOpts)
	}

	if s.ResolveNet != "10.0.0.0/8,10.200.200.0/12" {
		t.Errorf("%v: resolve-net not 10.0.0.0/8,10.200.200.0/12: %v", s.Name, s.ResolveNet)
	}

	_, err = s.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	_, _, err = clientTest.GetServer("non-existant", "ring", "myring", "")
	if err == nil {
		t.Error("Should throw error, non existant server")
	}
}

func TestCreateEditDeleteServer(t *testing.T) {
	// TestCreateServer
	port := int64(4300)
	inter := int64(5000)
	slowStart := int64(6000)
	s := &models.Server{
		Name:    "created",
		Address: "192.168.2.1",
		Port:    &port,
		ServerParams: models.ServerParams{
			Backup:         "enabled",
			Check:          "enabled",
			Maintenance:    "enabled",
			Ssl:            "enabled",
			AgentCheck:     "enabled",
			SslCertificate: "dummy.crt",
			TLSTickets:     "enabled",
			Verify:         "none",
			Inter:          &inter,
			OnMarkedDown:   "shutdown-sessions",
			OnError:        "mark-down",
			OnMarkedUp:     "shutdown-backup-sessions",
			Slowstart:      &slowStart,
			ProxyV2Options: []string{"ssl", "unique-id"},
		},
	}

	err := clientTest.CreateServer("backend", "test", s, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, server, err := clientTest.GetServer("created", "backend", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(server, s) {
		fmt.Printf("Created server: %v\n", server)
		fmt.Printf("Given server: %v\n", s)
		t.Error("Created server not equal to given server")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	err = clientTest.CreateServer("backend", "test", s, "", version)
	if err == nil {
		t.Error("Should throw error server already exists")
		version++
	}

	// TestEditServer
	port = int64(5300)
	slowStart = int64(3000)
	s = &models.Server{
		Name:    "created",
		Address: "192.168.3.1",
		Port:    &port,
		ServerParams: models.ServerParams{
			AgentCheck:     "disabled",
			Ssl:            "enabled",
			SslCertificate: "dummy.crt",
			SslCafile:      "dummy.ca",
			TLSTickets:     "disabled",
			Verify:         "required",
			Slowstart:      &slowStart,
		},
	}

	err = clientTest.EditServer("created", "backend", "test", s, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, server, err = clientTest.GetServer("created", "backend", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(server, s) {
		fmt.Printf("Edited server: %v\n", server)
		fmt.Printf("Given server: %v\n", s)
		t.Error("Edited server not equal to given server")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	// TestDeleteServer
	err = clientTest.DeleteServer("created", "backend", "test", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ := clientTest.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	_, _, err = clientTest.GetServer("created", "backend", "test", "")
	if err == nil {
		t.Error("DeleteServer failed, server test still exists")
	}

	err = clientTest.DeleteServer("created", "backend", "test_2", "", version)
	if err == nil {
		t.Error("Should throw error, non existant server")
		version++
	}
}

func Test_parseAddress(t *testing.T) {
	type args struct {
		address string
	}
	tests := []struct {
		name            string
		args            args
		wantIpOrAddress string
		wantPort        *int64
	}{
		{
			name:            "IPv6 with brackets",
			args:            args{address: "[fd00:6:48:c85:deb:f:62:4]:80"},
			wantIpOrAddress: "fd00:6:48:c85:deb:f:62:4",
			wantPort:        misc.Int64P(80),
		},
		{
			name:            "IPv6 without brackets",
			args:            args{address: "fd00:6:48:c85:deb:f:62:4:443"},
			wantIpOrAddress: "fd00:6:48:c85:deb:f:62:4",
			wantPort:        misc.Int64P(443),
		},
		{
			name:            "IPv6 without brackets, without port",
			args:            args{address: "fd00:6:48:c85:deb:f:62:a123"},
			wantIpOrAddress: "fd00:6:48:c85:deb:f:62:a123",
			wantPort:        nil,
		},
		{
			name:            "IPv4 with port",
			args:            args{address: "10.1.1.2:8080"},
			wantIpOrAddress: "10.1.1.2",
			wantPort:        misc.Int64P(8080),
		},
		{
			name:            "IPv4 without port",
			args:            args{address: "10.1.1.2"},
			wantIpOrAddress: "10.1.1.2",
			wantPort:        nil,
		},
		{
			name:            "Socket address",
			args:            args{address: "/var/run/test_socket"},
			wantIpOrAddress: "/var/run/test_socket",
			wantPort:        nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIpOrAddress, gotPort := parseAddress(tt.args.address)
			if gotIpOrAddress != tt.wantIpOrAddress {
				t.Errorf("parseAddress() gotIpOrAddress = %v, want %v", gotIpOrAddress, tt.wantIpOrAddress)
			}
			if gotPort != nil && tt.wantPort != nil && *gotPort != *tt.wantPort {
				t.Errorf("parseAddress() gotPort = %v, want %v", gotPort, tt.wantPort)
			}
		})
	}
}
