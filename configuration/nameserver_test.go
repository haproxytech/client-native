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

func TestGetNameservers(t *testing.T) {
	v, nameservers, err := client.GetNameservers("test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if len(nameservers) != 1 {
		t.Errorf("%v nameservers returned, expected 1", len(nameservers))
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	for _, l := range nameservers {
		if l.Name != "dns1" {
			t.Errorf("Expected only dns1, %v found", l.Name)
		}
		if *l.Address != "10.0.0.1" {
			t.Errorf("%v: Address not %s: %v", "10.0.0.1", l.Name, *l.Address)
		}
		if *l.Port != 53 {
			t.Errorf("%v: Port not %s: %v", "10.0.0.1", l.Name, *l.Port)
		}
	}
}

func TestGetNameserver(t *testing.T) {
	v, l, err := client.GetNameserver("dns1", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	if l.Name != "dns1" {
		t.Errorf("Expected only happe nameserver, %v found", l.Name)
	}
	if *l.Address != "10.0.0.1" {
		t.Errorf("%v: Address not 10.0.0.1: %v", l.Name, l.Address)
	}
	if *l.Port != 53 {
		t.Errorf("%v: Port not 53: %v", l.Name, *l.Port)
	}

	_, err = l.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	_, _, err = client.GetNameserver("community", "test", "")
	if err == nil {
		t.Error("Should throw error, non existant nameserver")
	}
}

func TestCreateEditDeleteNameserver(t *testing.T) {
	_, _, err := client.GetNameserver("dns1", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	address := "192.168.1.2"
	port := int64(1024)
	e := &models.Nameserver{
		Address: &address,
		Port:    &port,
		Name:    "hapcommunity",
	}

	err = client.CreateNameserver("test", e, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, l, err := client.GetNameserver("hapcommunity", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(e, l) {
		fmt.Printf("Created nameserver: %v\n", e)
		fmt.Printf("Given nameserver: %v\n", l)
		t.Error("Created nameserver not equal to given nameserver")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	err = client.CreateNameserver("test", e, "", version)
	if err == nil {
		t.Error("Should throw error nameserver already exists")
		version++
	}

	editPort := int64(1025)
	e = &models.Nameserver{
		Address: &address,
		Port:    &editPort,
		Name:    "hapcommunity",
	}

	err = client.EditNameserver("hapcommunity", "test", e, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, l, err = client.GetNameserver("hapcommunity", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(e, l) {
		fmt.Printf("Edited nameserver: %v\n", e)
		fmt.Printf("Given lsitener: %v\n", l)
		t.Error("Edited nameserver not equal to given nameserver")
	}

	// Delete
	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	err = client.DeleteNameserver("hapcommunity", "test", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ = client.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	_, _, err = client.GetNameserver("hapcommunity", "test", "")
	if err == nil {
		t.Error("DeleteNameserver failed, nameserver still exists")
	}

	err = client.DeleteNameserver("hapcommunity", "test", "", version)
	if err == nil {
		t.Error("Should throw error, non existant nameserver")
		version++
	}
}
