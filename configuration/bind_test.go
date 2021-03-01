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

func TestGetBinds(t *testing.T) {
	v, binds, err := client.GetBinds("test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if len(binds) != 2 {
		t.Errorf("%v binds returned, expected 2", len(binds))
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	for _, l := range binds {
		if l.Name != "webserv" && l.Name != "webserv2" {
			t.Errorf("Expected only webserv or webserv2 binds, %v found", l.Name)
		}
		if l.Address != "192.168.1.1" {
			t.Errorf("%v: Address not 192.168.1.1: %v", l.Name, l.Address)
		}
		if *l.Port != 80 && *l.Port != 8080 {
			t.Errorf("%v: Port not 80 or 8080: %v", l.Name, *l.Port)
		}
	}

	_, binds, err = client.GetBinds("test_2", "")
	if err != nil {
		t.Error(err.Error())
	}
	if len(binds) > 0 {
		t.Errorf("%v binds returned, expected 0", len(binds))
	}
}

func TestGetBind(t *testing.T) {
	v, l, err := client.GetBind("webserv", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	if l.Name != "webserv" {
		t.Errorf("Expected only webserv or webserv2 binds, %v found", l.Name)
	}
	if l.Address != "192.168.1.1" {
		t.Errorf("%v: Address not 192.168.1.1: %v", l.Name, l.Address)
	}
	if *l.Port != 80 {
		t.Errorf("%v: Port not 80 or 8080: %v", l.Name, *l.Port)
	}

	_, err = l.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	_, _, err = client.GetBind("webserv", "test_2", "")
	if err == nil {
		t.Error("Should throw error, non existent bind")
	}
}

func TestCreateEditDeleteBind(t *testing.T) {
	// TestCreateBind
	port := int64(4300)
	l := &models.Bind{
		Name:           "created",
		Address:        "192.168.2.1",
		Port:           &port,
		Ssl:            true,
		SslCertificate: "dummy.crt",
		Interface:      "eth0",
		Verify:         "optional",
		SslMinVer:      "TLSv1.3",
		SslMaxVer:      "TLSv1.3",
	}

	err := client.CreateBind("test", l, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, bind, err := client.GetBind("created", "test", "")
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

	err = client.CreateBind("test", l, "", version)
	if err == nil {
		t.Error("Should throw error bind already exists")
		version++
	}

	// TestEditBind
	port = int64(5300)
	tOut := int64(5)
	l = &models.Bind{
		Name:           "created",
		Address:        "192.168.3.1",
		Port:           &port,
		Transparent:    true,
		TCPUserTimeout: &tOut,
		SslMinVer:      "TLSv1.2",
		SslMaxVer:      "TLSv1.3",
		Interface:      "eth1",
	}

	err = client.EditBind("created", "test", l, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, bind, err = client.GetBind("created", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(bind, l) {
		fmt.Printf("Edited bind: %v\n", bind)
		fmt.Printf("Given lsitener: %v\n", l)
		t.Error("Edited bind not equal to given bind")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	// TestDeleteBind
	err = client.DeleteBind("created", "test", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ := client.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	_, _, err = client.GetBind("created", "test", "")
	if err == nil {
		t.Error("DeleteBind failed, bind test still exists")
	}

	err = client.DeleteBind("created", "test2", "", version)
	if err == nil {
		t.Error("Should throw error, non existent bind")
		version++
	}
}
