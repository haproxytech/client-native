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

	"github.com/haproxytech/client-native/v5/models"
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
