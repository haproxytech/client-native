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

package test

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/haproxytech/client-native/v6/configuration"
	"github.com/haproxytech/client-native/v6/misc"
	"github.com/haproxytech/client-native/v6/models"
	"github.com/stretchr/testify/require"
)

func serverExpectation() map[string]models.Servers {
	initStructuredExpected()
	res := StructuredToServerMap()
	// Add individual entries
	for k, vs := range res {
		for _, v := range vs {
			key := fmt.Sprintf("%s/%s", k, v.Name)
			res[key] = models.Servers{v}
		}
	}
	return res
}

func TestGetServers(t *testing.T) { //nolint:gocognit,gocyclo
	ms := make(map[string]models.Servers)
	v, servers, err := clientTest.GetServers(configuration.BackendParentName, "test", "")
	if err != nil {
		t.Error(err.Error())
	}
	ms["backend/test"] = servers

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	_, servers, err = clientTest.GetServers(configuration.BackendParentName, "test_2", "")
	if err != nil {
		t.Error(err.Error())
	}
	ms["backend/test_2"] = servers

	checkServers(t, ms)
}

func checkServers(t *testing.T, got map[string]models.Servers) {
	exp := serverExpectation()
	for k, v := range got {
		want, ok := exp[k]
		require.True(t, ok, "k=%s", k)
		require.Equal(t, len(want), len(v), "k=%s", k)
		for _, g := range v {
			for _, w := range want {
				if g.Name == w.Name {
					require.True(t, g.Equal(*w), "k=%s - diff %v", k, cmp.Diff(*g, *w))
					break
				}
			}
		}
	}
}

func TestGetServer(t *testing.T) {
	m := make(map[string]models.Servers)

	v, s, err := clientTest.GetServer("webserv", configuration.BackendParentName, "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}
	m["backend/test/webserv"] = models.Servers{s}

	_, err = s.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	_, _, err = clientTest.GetServer("webserv", configuration.BackendParentName, "test_2", "")
	if err == nil {
		t.Error("Should throw error, non existent server")
	}

	checkServers(t, m)
}

func TestGetRingServer(t *testing.T) {
	m := make(map[string]models.Servers)
	v, s, err := clientTest.GetServer("mysyslogsrv", configuration.RingParentName, "myring", "")
	if err != nil {
		t.Error(err.Error())
	}
	m["ring/myring/mysyslogsrv"] = models.Servers{s}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	v, s, err = clientTest.GetServer("s1", configuration.RingParentName, "myring", "")
	if err != nil {
		t.Error(err.Error())
	}
	m["ring/myring/s1"] = models.Servers{s}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	checkServers(t, m)

	_, _, err = clientTest.GetServer("non-existent", configuration.RingParentName, "myring", "")
	if err == nil {
		t.Error("Should throw error, non existent server")
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
		Metadata: map[string]interface{}{
			"type": "good type",
			"id":   "my-id-12",
		},
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
			Curves:         "brainpoolP384r1",
			Sigalgs:        "RSA+SHA256",
			ClientSigalgs:  "ECDSA+SHA256",
			LogBufsize:     misc.Int64P(11),
			SetProxyV2TlvFmt: &models.ServerParamsSetProxyV2TlvFmt{
				ID:    misc.StringP("0x50"),
				Value: misc.StringP("%[fc_pp_tlv(0x20)]"),
			},
			IdlePing:          misc.Int64P(10000),
			CheckReusePool:    "enabled",
			CheckPoolConnName: "bar",
			TCPMd5sig:         "secretpass",
		},
	}

	err := clientTest.CreateServer(configuration.BackendParentName, "test", s, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, server, err := clientTest.GetServer("created", configuration.BackendParentName, "test", "")
	if err != nil {
		t.Error(err.Error())
	}
	require.True(t, server.Equal(*s), "diff %v", cmp.Diff(server, s))

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	err = clientTest.CreateServer(configuration.BackendParentName, "test", s, "", version)
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
			StrictMaxconn:  true,
		},
	}

	err = clientTest.EditServer("created", configuration.BackendParentName, "test", s, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, server, err = clientTest.GetServer("created", configuration.BackendParentName, "test", "")
	if err != nil {
		t.Error(err.Error())
	}
	require.True(t, server.Equal(*s), "diff %v", cmp.Diff(server, s))

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	// TestDeleteServer
	err = clientTest.DeleteServer("created", configuration.BackendParentName, "test", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ := clientTest.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	_, _, err = clientTest.GetServer("created", configuration.BackendParentName, "test", "")
	if err == nil {
		t.Error("DeleteServer failed, server test still exists")
	}

	err = clientTest.DeleteServer("created", configuration.BackendParentName, "test_2", "", version)
	if err == nil {
		t.Error("Should throw error, non existent server")
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
			gotIpOrAddress, gotPort := configuration.ParseAddress(tt.args.address)
			if gotIpOrAddress != tt.wantIpOrAddress {
				t.Errorf("ParseAddress() gotIpOrAddress = %v, want %v", gotIpOrAddress, tt.wantIpOrAddress)
			}
			if gotPort != nil && tt.wantPort != nil && *gotPort != *tt.wantPort {
				t.Errorf("ParseAddress() gotPort = %v, want %v", gotPort, tt.wantPort)
			}
		})
	}
}
