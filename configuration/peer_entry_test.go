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

const (
	testPort = 1023
	testIP   = "192.168.1.1"
)

func TestGetPeerEntries(t *testing.T) {
	v, peerEntries, err := client.GetPeerEntries("mycluster", "")
	if err != nil {
		t.Error(err.Error())
	}

	if len(peerEntries) != 2 {
		t.Errorf("%v peerEntries returned, expected 2", len(peerEntries))
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	for _, l := range peerEntries {
		if l.Name != "hapee" && l.Name != "aggregator" {
			t.Errorf("Expected only hapee and aggregator, %v found", l.Name)
		}
		if *l.Address != testIP && *l.Address != "HARDCODEDCLUSTERIP" {
			t.Errorf("%v: Address not %s: %v", testIP, l.Name, *l.Address)
		}
		if *l.Port != testPort && *l.Port != 10023 {
			t.Errorf("%v: Port not %s: %v", testIP, l.Name, *l.Port)
		}
	}
}

func TestGetPeerEntry(t *testing.T) {
	v, l, err := client.GetPeerEntry("hapee", "mycluster", "")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	if l.Name != "hapee" {
		t.Errorf("Expected only happe peerEntry, %v found", l.Name)
	}
	if *l.Address != "192.168.1.1" {
		t.Errorf("%v: Address not 192.168.1.1: %v", l.Name, l.Address)
	}
	if *l.Port != 1023 {
		t.Errorf("%v: Port not 1023: %v", l.Name, *l.Port)
	}

	_, err = l.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	_, _, err = client.GetPeerEntry("community", "mycluster", "")
	if err == nil {
		t.Error("Should throw error, non existant peer entry")
	}
}

func TestCreateEditDeletePeerEntry(t *testing.T) {
	_, _, err := client.GetPeerEntry("hapee", "mycluster", "")
	if err != nil {
		t.Error(err.Error())
	}

	address := "192.168.1.2"
	port := int64(1024)
	e := &models.PeerEntry{
		Address: &address,
		Port:    &port,
		Name:    "hapcommunity",
	}

	err = client.CreatePeerEntry("mycluster", e, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, l, err := client.GetPeerEntry("hapcommunity", "mycluster", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(e, l) {
		fmt.Printf("Created peerEntry: %v\n", e)
		fmt.Printf("Given peerEntry: %v\n", l)
		t.Error("Created peerEntry not equal to given peerEntry")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	err = client.CreatePeerEntry("mycluster", e, "", version)
	if err == nil {
		t.Error("Should throw error peerEntry already exists")
		version++
	}

	editPort := int64(1025)
	e = &models.PeerEntry{
		Address: &address,
		Port:    &editPort,
		Name:    "hapcommunity",
	}

	err = client.EditPeerEntry("hapcommunity", "mycluster", e, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, l, err = client.GetPeerEntry("hapcommunity", "mycluster", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(e, l) {
		fmt.Printf("Edited peerEntry: %v\n", e)
		fmt.Printf("Given lsitener: %v\n", l)
		t.Error("Edited peerEntry not equal to given peerEntry")
	}

	// Delete
	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	err = client.DeletePeerEntry("hapcommunity", "mycluster", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ = client.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	_, _, err = client.GetPeerEntry("hapcommunity", "mycluster", "")
	if err == nil {
		t.Error("DeletePeerEntry failed, peer entry still exists")
	}

	err = client.DeletePeerEntry("hapcommunity", "mycluster", "", version)
	if err == nil {
		t.Error("Should throw error, non existant peer entry")
		version++
	}
}
