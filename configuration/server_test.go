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

	"github.com/haproxytech/client-native/v2/models"
)

func TestGetServers(t *testing.T) { //nolint:gocognit,gocyclo
	v, servers, err := client.GetServers("test", "")
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

	_, servers, err = client.GetServers("test_2", "")
	if err != nil {
		t.Error(err.Error())
	}
	if len(servers) > 0 {
		t.Errorf("%v servers returned, expected 0", len(servers))
	}
}

func TestGetServer(t *testing.T) {
	v, s, err := client.GetServer("webserv", "test", "")
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

	_, err = s.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	_, _, err = client.GetServer("webserv", "test_2", "")
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
		Name:           "created",
		Address:        "192.168.2.1",
		Port:           &port,
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
	}

	err := client.CreateServer("test", s, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, server, err := client.GetServer("created", "test", "")
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

	err = client.CreateServer("test", s, "", version)
	if err == nil {
		t.Error("Should throw error server already exists")
		version++
	}

	// TestEditServer
	port = int64(5300)
	slowStart = int64(3000)
	s = &models.Server{
		Name:           "created",
		Address:        "192.168.3.1",
		Port:           &port,
		AgentCheck:     "disabled",
		Ssl:            "enabled",
		SslCertificate: "dummy.crt",
		SslCafile:      "dummy.ca",
		TLSTickets:     "disabled",
		Verify:         "required",
		Slowstart:      &slowStart,
	}

	err = client.EditServer("created", "test", s, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, server, err = client.GetServer("created", "test", "")
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
	err = client.DeleteServer("created", "test", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ := client.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	_, _, err = client.GetServer("created", "test", "")
	if err == nil {
		t.Error("DeleteServer failed, server test still exists")
	}

	err = client.DeleteServer("created", "test_2", "", version)
	if err == nil {
		t.Error("Should throw error, non existant server")
		version++
	}
}
