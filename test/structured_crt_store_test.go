// Copyright 2026 HAProxy Technologies
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
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/haproxytech/client-native/v6/models"
	"github.com/stretchr/testify/require"
)

func checkStructuredCrtStores(t *testing.T, got map[string]models.CrtStores) {
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

func TestGetStructuredCrtStores(t *testing.T) {
	clientTest, filename, err := getTestClient()
	require.NoError(t, err)
	defer os.Remove(filename)
	version := int64(1)

	v, stores, err := clientTest.GetStructuredCrtStores("")
	require.NoError(t, err)
	require.Equal(t, version, v, "Version %v returned, expected %v", v, version)
	require.Len(t, stores, 1)
	checkStructuredCrtStores(t, map[string]models.CrtStores{"": stores})
}

func TestGetStructuredCrtStore(t *testing.T) {
	clientTest, filename, err := getTestClient()
	require.NoError(t, err)
	defer os.Remove(filename)
	version := int64(1)

	v, store, err := clientTest.GetStructuredCrtStore(testCrtStoreName, "")
	require.NoError(t, err)
	require.Equal(t, version, v, "Version %v returned, expected %v", v, version)
	require.Len(t, store.CrtLoads, 2)
	checkStructuredCrtStores(t, map[string]models.CrtStores{testCrtStoreName: {store}})

	_, _, err = clientTest.GetStructuredCrtStore("doesnotexist", "")
	require.Error(t, err, "should throw error, non existent crt_store")
}

func TestCreateEditStructuredCrtStore(t *testing.T) {
	clientTest, filename, err := getTestClient()
	require.NoError(t, err)
	defer os.Remove(filename)
	version := int64(1)

	store := &models.CrtStore{
		CrtStoreBase: models.CrtStoreBase{Name: "test-structured", CrtBase: "/certs", KeyBase: "/keys"},
		CrtLoads: map[string]models.CrtLoad{
			"c1.pem": {Certificate: "c1.pem", Key: "k1.pem", Alias: "c1"},
			"c2.pem": {Certificate: "c2.pem", Ocsp: "ocsp.der", OcspUpdate: models.CrtLoadOcspUpdateEnabled},
		},
	}
	require.NoError(t, clientTest.CreateStructuredCrtStore(store, "", version))
	version++

	v, got, err := clientTest.GetStructuredCrtStore("test-structured", "")
	require.NoError(t, err)
	require.Equal(t, version, v, "Version %v returned, expected %v", v, version)
	require.True(t, got.Equal(*store), "diff %v", cmp.Diff(*got, *store))

	require.Error(t, clientTest.CreateStructuredCrtStore(store, "", version), "should throw error, crt_store already exists")

	// edit: keep a single, different load entry
	store.CrtLoads = map[string]models.CrtLoad{"c3.pem": {Certificate: "c3.pem"}}
	require.NoError(t, clientTest.EditStructuredCrtStore("test-structured", store, "", version))
	version++

	v, got, err = clientTest.GetStructuredCrtStore("test-structured", "")
	require.NoError(t, err)
	require.Equal(t, version, v, "Version %v returned, expected %v", v, version)
	require.True(t, got.Equal(*store), "diff %v", cmp.Diff(*got, *store))
}
