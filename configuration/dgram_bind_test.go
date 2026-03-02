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

	"github.com/haproxytech/client-native/v5/config-parser/params"
	"github.com/haproxytech/client-native/v5/config-parser/types"
	"github.com/haproxytech/client-native/v5/models"

	"github.com/stretchr/testify/require"
)

func TestGetDgramBinds(t *testing.T) {
	v, dBinds, err := clientTest.GetDgramBinds("sylog-loadb", "")
	if err != nil {
		t.Error(err.Error())
	}

	if len(dBinds) != 1 {
		t.Errorf("%v dgram-binds returned, expected 1", len(dBinds))
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	for _, l := range dBinds {
		if l.Name != "webserv" {
			t.Errorf("Expected webserv dgram-binds, %v found", l.Name)
		}
		if l.Address != "127.0.0.1" {
			t.Errorf("%v: Address not 127.0.0.1: %v", l.Name, l.Address)
		}
		if *l.Port != 1514 {
			t.Errorf("%v: Port not 1514 : %v", l.Name, *l.Port)
		}
	}

	// _, dBinds, err = clientTest.GetDgramBinds("test_2", "")
	// if err != nil {
	// 	t.Error(err.Error())
	// }
	// if len(dBinds) > 0 {
	// 	t.Errorf("%v dgram-binds returned, expected 0", len(dBinds))
	// }
}

func TestGetDgramBind(t *testing.T) {
	v, l, err := clientTest.GetDgramBind("webserv", "sylog-loadb", "")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	if l.Name != "webserv" {
		t.Errorf("Expected only webserv dgram-bind, %v found", l.Name)
	}
	if l.Address != "127.0.0.1" {
		t.Errorf("%v: Address not 127.0.0.1: %v", l.Name, l.Address)
	}
	if *l.Port != 1514 {
		t.Errorf("%v: Port not 1514: %v", l.Name, *l.Port)
	}

	_, err = l.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	_, _, err = clientTest.GetDgramBind("webserv", "test_2", "")
	if err == nil {
		t.Error("Should throw error, non existent bind")
	}
}

func TestCreateEditDeleteDgramBind(t *testing.T) {
	// TestCreateBind
	port := int64(4300)
	l := &models.DgramBind{
		Address:   "192.168.2.1",
		Port:      &port,
		Name:      "created",
		Interface: "eth0",
	}

	err := clientTest.CreateDgramBind("sylog-loadb", l, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, bind, err := clientTest.GetDgramBind("created", "sylog-loadb", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(bind, l) {
		fmt.Printf("Created bind: %v\n", bind)
		fmt.Printf("Given bind: %v\n", l)
		t.Error("Created bind not equal to given bind")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	err = clientTest.CreateDgramBind("sylog-loadb", l, "", version)
	if err == nil {
		t.Error("Should throw error bind already exists")
		version++
	}

	// TestEditBind
	port = int64(5300)
	l = &models.DgramBind{
		Address:     "192.168.3.1",
		Port:        &port,
		Name:        "created",
		Transparent: true,
		Interface:   "eth1",
	}

	err = clientTest.EditDgramBind("created", "sylog-loadb", l, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, bind, err = clientTest.GetDgramBind("created", "sylog-loadb", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(bind, l) {
		fmt.Printf("Edited bind: %v\n", bind)
		fmt.Printf("Given bind: %v\n", l)
		t.Error("Edited bind not equal to given bind")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	// TestDeleteBind
	err = clientTest.DeleteDgramBind("created", "sylog-loadb", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ := clientTest.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	_, _, err = clientTest.GetDgramBind("created", "sylog-loadb", "")
	if err == nil {
		t.Error("DeleteDgramBind failed, bind test still exists")
	}

	err = clientTest.DeleteDgramBind("created", "test2", "", version)
	if err == nil {
		t.Error("Should throw error, non existent bind")
		version++
	}
}

func TestParseDgramBindIPv6(t *testing.T) {
	tests := []struct {
		name            string
		path            string
		params          []params.DgramBindOption
		expectedAddress string
		expectedPort    *int64
		expectedName    string
	}{
		{
			name:            "bracketed IPv6 with port",
			path:            "[fd66:c3ec:c7fc::7c]:123",
			params:          []params.DgramBindOption{&params.BindOptionValue{Name: "name", Value: "ntp6"}},
			expectedAddress: "fd66:c3ec:c7fc::7c",
			expectedPort:    int64P(123),
			expectedName:    "ntp6",
		},
		{
			name:            "bracketed IPv6 with port and no name param",
			path:            "[2a04:4b07:21dc:218::1]:53",
			expectedAddress: "2a04:4b07:21dc:218::1",
			expectedPort:    int64P(53),
			expectedName:    "[2a04:4b07:21dc:218::1]:53",
		},
		{
			name:            "double colon with port",
			path:            ":::443",
			expectedAddress: "::",
			expectedPort:    int64P(443),
			expectedName:    ":::443",
		},
		{
			name:            "IPv4 with port",
			path:            "10.0.0.1:8080",
			params:          []params.DgramBindOption{&params.BindOptionValue{Name: "name", Value: "test"}},
			expectedAddress: "10.0.0.1",
			expectedPort:    int64P(8080),
			expectedName:    "test",
		},
		{
			name:            "wildcard with port",
			path:            ":443",
			expectedAddress: "",
			expectedPort:    int64P(443),
			expectedName:    ":443",
		},
		{
			name:            "unix socket path",
			path:            "/var/run/haproxy.sock",
			expectedAddress: "/var/run/haproxy.sock",
			expectedPort:    nil,
			expectedName:    "/var/run/haproxy.sock",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ondisk := types.DgramBind{
				Path:   tc.path,
				Params: tc.params,
			}
			result := ParseDgramBind(ondisk)
			require.NotNil(t, result, "ParseDgramBind returned nil for path %q", tc.path)
			require.Equal(t, tc.expectedAddress, result.Address, "address mismatch")
			if tc.expectedPort != nil {
				require.NotNil(t, result.Port, "expected port to be set")
				require.Equal(t, *tc.expectedPort, *result.Port, "port mismatch")
			} else {
				require.Nil(t, result.Port, "expected port to be nil")
			}
			require.Equal(t, tc.expectedName, result.Name, "name mismatch")
		})
	}
}

func TestSerializeDgramBindIPv6(t *testing.T) {
	tests := []struct {
		name         string
		dgramBind    models.DgramBind
		expectedPath string
	}{
		{
			name: "IPv6 address with port",
			dgramBind: models.DgramBind{
				Address: "fd66:c3ec:c7fc::7c",
				Port:    int64P(123),
				Name:    "ntp6",
			},
			expectedPath: "[fd66:c3ec:c7fc::7c]:123",
		},
		{
			name: "IPv6 double colon with port",
			dgramBind: models.DgramBind{
				Address: "::",
				Port:    int64P(443),
				Name:    "test",
			},
			expectedPath: "[::]:443",
		},
		{
			name: "IPv4 address with port",
			dgramBind: models.DgramBind{
				Address: "10.0.0.1",
				Port:    int64P(8080),
				Name:    "test",
			},
			expectedPath: "10.0.0.1:8080",
		},
		{
			name: "address only, no port",
			dgramBind: models.DgramBind{
				Address: "/var/run/haproxy.sock",
				Name:    "test",
			},
			expectedPath: "/var/run/haproxy.sock",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := SerializeDgramBind(tc.dgramBind)
			require.Equal(t, tc.expectedPath, result.Path, "serialized path mismatch")
		})
	}
}

func TestParseDgramBindRoundTrip(t *testing.T) {
	tests := []struct {
		name string
		path string
	}{
		{"IPv6 bracketed", "[fd66:c3ec:c7fc::7c]:123"},
		{"IPv6 full", "[2a04:4b07:21dc:218::1]:53"},
		{"IPv6 double colon", "[::]:443"},
		{"IPv4", "10.0.0.1:8080"},
		{"wildcard", ":443"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ondisk := types.DgramBind{
				Path:   tc.path,
				Params: []params.DgramBindOption{&params.BindOptionValue{Name: "name", Value: "roundtrip"}},
			}
			parsed := ParseDgramBind(ondisk)
			require.NotNil(t, parsed, "ParseDgramBind returned nil for path %q", tc.path)

			serialized := SerializeDgramBind(*parsed)
			require.Equal(t, tc.path, serialized.Path, "round-trip path mismatch")
		})
	}
}

func int64P(v int64) *int64 {
	return &v
}
