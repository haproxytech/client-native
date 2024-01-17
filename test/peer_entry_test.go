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
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/haproxytech/client-native/v5/models"
	"github.com/stretchr/testify/require"
)

const (
	testPort = 1023
	testIP   = "192.168.1.1"
)

func peerEntryExpectation() map[string]models.PeerEntries {
	initStructuredExpected()
	res := StructuredToPeerEntryMap()
	for k, vs := range res {
		for _, v := range vs {
			key := fmt.Sprintf("%s/%s", k, v.Name)
			res[key] = models.PeerEntries{v}
		}
	}
	return res
}

func TestGetPeerEntries(t *testing.T) {
	m := make(map[string]models.PeerEntries)
	v, peerEntries, err := clientTest.GetPeerEntries("mycluster", "")
	if err != nil {
		t.Error(err.Error())
	}

	m["peers/mycluster"] = peerEntries

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}
	checkPeerEntries(t, m)
}

func TestGetPeerEntry(t *testing.T) {
	m := make(map[string]models.PeerEntries)

	v, l, err := clientTest.GetPeerEntry("hapee", "mycluster", "")
	if err != nil {
		t.Error(err.Error())
	}
	m["peers/mycluster/hapee"] = models.PeerEntries{l}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	checkPeerEntries(t, m)

	_, _, err = clientTest.GetPeerEntry("community", "mycluster", "")
	if err == nil {
		t.Error("Should throw error, non existent peer entry")
	}
}

func checkPeerEntries(t *testing.T, got map[string]models.PeerEntries) {
	exp := peerEntryExpectation()
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

func TestCreateEditDeletePeerEntry(t *testing.T) {
	_, _, err := clientTest.GetPeerEntry("hapee", "mycluster", "")
	if err != nil {
		t.Error(err.Error())
	}

	address := "192.168.1.2"
	port := int64(1024)
	e := &models.PeerEntry{
		Address: &address,
		Port:    &port,
		Name:    "hapcommunity",
		Shard:   2,
	}

	err = clientTest.CreatePeerEntry("mycluster", e, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, l, err := clientTest.GetPeerEntry("hapcommunity", "mycluster", "")
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

	err = clientTest.CreatePeerEntry("mycluster", e, "", version)
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

	err = clientTest.EditPeerEntry("hapcommunity", "mycluster", e, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, l, err = clientTest.GetPeerEntry("hapcommunity", "mycluster", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(e, l) {
		fmt.Printf("Edited peerEntry: %v\n", e)
		fmt.Printf("Given peerEntry: %v\n", l)
		t.Error("Edited peerEntry not equal to given peerEntry")
	}

	// Delete
	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	err = clientTest.DeletePeerEntry("hapcommunity", "mycluster", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ = clientTest.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	_, _, err = clientTest.GetPeerEntry("hapcommunity", "mycluster", "")
	if err == nil {
		t.Error("DeletePeerEntry failed, peer entry still exists")
	}

	err = clientTest.DeletePeerEntry("hapcommunity", "mycluster", "", version)
	if err == nil {
		t.Error("Should throw error, non existent peer entry")
		version++
	}
}
