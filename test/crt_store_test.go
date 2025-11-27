// Copyright 2024 HAProxy Technologies
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
	"github.com/haproxytech/client-native/v6/configuration"
	"github.com/haproxytech/client-native/v6/models"
	"github.com/stretchr/testify/require"
)

const testCrtStoreName = "cert-bunker1"

func crtStoreExpectations() map[string]models.CrtStores {
	initStructuredExpected()
	res := StructuredToCrtStoreMap()
	// Add individual entries
	for _, vs := range res {
		for _, v := range vs {
			key := v.Name
			res[key] = models.CrtStores{v}
		}
	}
	return res
}

func checkCrtStores(t *testing.T, got map[string]models.CrtStores) {
	exp := crtStoreExpectations()
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

func TestGetCrtStores(t *testing.T) {
	m := make(map[string]models.CrtStores)
	v, stores, err := clientTest.GetCrtStores("")
	if err != nil {
		t.Error(err)
	}

	if v != version {
		t.Errorf("found version %d, expected %d", v, version)
	}

	if len(stores) != 1 {
		t.Errorf("crt_stores: found %d, expected 1", len(stores))
	}
	m[stores[0].Name] = models.CrtStores{stores[0]}

	checkCrtStores(t, m)
}

func TestGetCrtStore(t *testing.T) {
	m := make(map[string]models.CrtStores)

	v, store, err := clientTest.GetCrtStore(testCrtStoreName, "")
	if err != nil {
		t.Error(err)
	}
	m[testCrtStoreName] = models.CrtStores{store}

	if v != version {
		t.Errorf("found version %d, expected %d", v, version)
	}
	checkCrtStores(t, m)

	_, _, err = clientTest.GetCrtStore("doesnotexist", "")
	if err == nil {
		t.Error("Should throw error, non existent crt-store")
	}
}

func TestCreateEditDeleteCrtStore(t *testing.T) {
	const name = "test-store-42"

	store := &models.CrtStore{
		CrtStoreBase: models.CrtStoreBase{
			Name:    name,
			CrtBase: "/secure/certs",
			KeyBase: "/secure/keys",
		},
	}

	err := clientTest.CreateCrtStore(store, "", version)
	if err != nil {
		t.Fatal(err)
	} else {
		version++
	}

	v, created, err := clientTest.GetCrtStore(name, "")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(created, store) {
		fmt.Printf("Created CrtStore: %v\n", created)
		fmt.Printf("Given CrtStore: %+v\n", store)
		t.Fatal("Created CrtStore not equal to given CrtStore")
	}

	if v != version {
		t.Errorf("found version %d, expected %d", v, version)
	}

	err = clientTest.CreateCrtStore(store, "", version)
	if err == nil {
		t.Fatal("Should throw error CrtStore already exists")
		version++
	}

	// Modify the section.
	store.KeyBase = "/secure/pkeys"
	err = clientTest.EditCrtStore(name, store, "", version)
	if err != nil {
		t.Errorf("EditCrtStore: %v", err)
	} else {
		version++
	}

	// Check if the modification was effective.
	v, created, err = clientTest.GetCrtStore(name, "")
	if err != nil {
		t.Fatal(err)
	}
	if created == nil {
		t.Fatal("got a nil CrtStore")
	}
	if v != version {
		t.Errorf("found version %d, expected %d", v, version)
	}
	if created.KeyBase != store.KeyBase {
		t.Errorf("CrtStore KeyBase was not modified: got %s, expected %s", created.KeyBase, store.KeyBase)
	}

	// Delete the section.
	err = clientTest.DeleteCrtStore(name, "", version)
	if err != nil {
		t.Error(err)
	} else {
		version++
	}

	if v, _ := clientTest.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	err = clientTest.DeleteCrtStore(name, "", 999999)
	if err != nil {
		if confErr, ok := err.(*configuration.ConfError); ok {
			if !confErr.Is(configuration.ErrVersionMismatch) {
				t.Error("Should throw configuration.ErrVersionMismatch error")
			}
		} else {
			t.Error("Should throw configuration.ErrVersionMismatch error")
		}
	}

	_, _, err = clientTest.GetCrtStore(name, "")
	if err == nil {
		t.Errorf("DeleteCrtStore failed: %s still exists", name)
	}

	err = clientTest.DeleteCrtStore("doesnotexist", "", version)
	if err == nil {
		t.Error("Should throw error, non existent CrtStore")
		version++
	}
}

func TestCreateEditDeleteCrtLoads(t *testing.T) {
	const name = "test-store-66"

	store := &models.CrtStore{
		CrtStoreBase: models.CrtStoreBase{
			Name:    name,
			CrtBase: "/secure/certs",
			KeyBase: "/secure/keys",
		},
	}

	err := clientTest.CreateCrtStore(store, "", version)
	if err != nil {
		t.Fatal(err)
	} else {
		version++
	}

	v, created, err := clientTest.GetCrtStore(name, "")
	if err != nil {
		t.Fatal(err)
	}
	if v != version {
		t.Fatalf("found version %d, expected %d", v, version)
	}
	if !reflect.DeepEqual(created, store) {
		t.Logf("Created CrtStore: %#v\n", created)
		t.Logf("Given   CrtStore: %#v\n", store)
		t.Fatal("Created CrtStore not equal to given CrtStore")
	}

	// Add a "load crt <filename>" entry
	entry1 := &models.CrtLoad{Certificate: "c1.pem", Key: "k1.pem"}
	err = clientTest.CreateCrtLoad(name, entry1, "", v)
	if err != nil {
		t.Fatal(err)
	} else {
		version++
	}

	// Get the entry
	v, e1, err := clientTest.GetCrtLoad("c1.pem", name, "")
	if err != nil {
		t.Fatal(err)
	}
	if v != version {
		t.Fatalf("found version %d, expected %d", v, version)
	}
	if !reflect.DeepEqual(e1, entry1) {
		t.Logf("Created CrtLoad: %#v\n", e1)
		t.Logf("Given   CrtLoad: %#v\n", entry1)
		t.Fatal("Created CrtLoad not equal to given CrtLoad")
	}

	// Add another one
	entry2 := &models.CrtLoad{Certificate: "c2.pem", Ocsp: "ocsp.der", OcspUpdate: "enabled"}
	err = clientTest.CreateCrtLoad(name, entry2, "", v)
	if err != nil {
		t.Fatal(err)
	} else {
		version++
	}
	v++

	// Modify entry2
	entry2.OcspUpdate = "disabled"
	err = clientTest.EditCrtLoad("c2.pem", name, entry2, "", v)
	if err != nil {
		t.Fatal(err)
	} else {
		version++
	}

	// Check if entry2 was modified
	v, e2, err := clientTest.GetCrtLoad("c2.pem", name, "")
	if err != nil {
		t.Fatal(err)
	}
	if v != version {
		t.Fatalf("found version %d, expected %d", v, version)
	}
	if !reflect.DeepEqual(e2, entry2) {
		t.Logf("Created CrtLoad: %#v\n", e2)
		t.Logf("Given   CrtLoad: %#v\n", entry2)
		t.Fatal("Created CrtLoad not equal to given CrtLoad")
	}

	// Delete entry1
	err = clientTest.DeleteCrtLoad("c1.pem", name, "", v)
	if err != nil {
		t.Fatal(err)
	} else {
		version++
	}

	// Get all the loads
	v, loads, err := clientTest.GetCrtLoads(name, "")
	if err != nil {
		t.Fatal(err)
	}
	// Only c2.pem should be here
	if len(loads) != 1 || loads[0].Certificate != "c2.pem" {
		t.Fatal("DeleteCrtLoad() did not work correctly")
	}
}
