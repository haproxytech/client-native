// Copyright 2025 HAProxy Technologies
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
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/haproxytech/client-native/v6/configuration"
	"github.com/haproxytech/client-native/v6/misc"
	"github.com/haproxytech/client-native/v6/models"
	"github.com/stretchr/testify/require"
)

func acmeExpectation() map[string]models.AcmeProviders {
	initStructuredExpected()
	res := StructuredToAcmeMap()
	// Add individual entries
	for _, vs := range res {
		for _, v := range vs {
			key := v.Name
			res[key] = models.AcmeProviders{v}
		}
	}
	return res
}

func TestGetAcmeProviders(t *testing.T) {
	m := make(map[string]models.AcmeProviders)
	v, acmes, err := clientTest.GetAcmeProviders("")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("version %v returned, expected %v", v, version)
	}
	m[""] = acmes
	checkAcmeProviders(t, m)
}

func TestGetAcmeProvider(t *testing.T) {
	m := make(map[string]models.AcmeProviders)

	v, r, err := clientTest.GetAcmeProvider("test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("version %v returned, expected %v", v, version)
	}
	m[""] = models.AcmeProviders{r}
	checkAcmeProviders(t, m)

	_, _, err = clientTest.GetAcmeProvider("doesnotexist", "")
	if err == nil {
		t.Error("should throw error, non existent rings section")
	}
}

func checkAcmeProviders(t *testing.T, got map[string]models.AcmeProviders) {
	exp := acmeExpectation()

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

func TestCreateEditDeleteAcmeProvider(t *testing.T) {
	require := require.New(t)

	a := &models.AcmeProvider{
		Name:         "ninja",
		AccountKey:   "ninja-acme.key",
		AcmeProvider: "godaddy",
		AcmeVars: map[string]string{
			"ApiKey":   "foobar",
			"WeirdKey": "\"__, +=\"",
		},
		Bits:      misc.Int64P(2048),
		Challenge: "http-01",
		Contact:   "me@example.com",
		Curves:    "dem curves",
		Directory: "https://acme.ninja.com/directory",
		Keytype:   "ECDSA",
		Map:       "acme@virt",
	}

	err := clientTest.CreateAcmeProvider(a, "", version)
	require.NoError(err)
	version++

	v, acme, err := clientTest.GetAcmeProvider(a.Name, "")
	require.NoError(err)
	require.Equal(v, version)

	require.Equal(acme.Name, a.Name)
	require.Equal(acme.Contact, a.Contact)
	require.Equal(acme.Bits, a.Bits)
	require.Equal(acme, a)

	err = clientTest.CreateAcmeProvider(a, "", version)
	if err == nil {
		t.Error("should throw error: acme section already exists")
		version++
	}

	// Edit
	a.Contact = "new@example.com"
	a.Bits = misc.Int64P(4096)
	err = clientTest.EditAcmeProvider(a.Name, a, "", version)
	require.NoError(err)
	version++

	v, acme, err = clientTest.GetAcmeProvider(a.Name, "")
	require.NoError(err)
	require.Equal(v, version)
	require.Equal(acme, a)

	// Delete
	err = clientTest.DeleteAcmeProvider(a.Name, "", version)
	require.NoError(err)
	version++

	if v, _ := clientTest.GetVersion(""); v != version {
		t.Error("version not incremented")
	}

	err = clientTest.DeleteAcmeProvider(a.Name, "", 999999)
	if err != nil {
		if confErr, ok := err.(*configuration.ConfError); ok {
			if !confErr.Is(configuration.ErrVersionMismatch) {
				t.Error("should throw configuration.ErrVersionMismatch error")
			}
		} else {
			t.Error("should throw configuration.ErrVersionMismatch error")
		}
	}
	_, _, err = clientTest.GetAcmeProvider(a.Name, "")
	if err == nil {
		t.Errorf("deleteAcmeProvider failed, '%s' still exists", a.Name)
	}

	err = clientTest.DeleteAcmeProvider("doesnotexist", "", version)
	if err == nil {
		t.Error("should throw error, non existent section")
		version++
	}
}
