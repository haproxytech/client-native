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
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/haproxytech/client-native/v6/models"
	"github.com/stretchr/testify/require"
)

func TestGetStructuredHTTPHardenings(t *testing.T) {
	clientTest, filename, err := getTestClient()
	require.NoError(t, err)
	defer os.Remove(filename)
	version := int64(1)

	m := make(map[string]models.HTTPHardenings)
	v, profiles, err := clientTest.GetStructuredHTTPHardenings("")
	require.NoError(t, err)
	require.Equal(t, version, v)
	m[""] = profiles
	checkStructuredHTTPHardenings(t, m)
}

func TestGetStructuredHTTPHardening(t *testing.T) {
	clientTest, filename, err := getTestClient()
	require.NoError(t, err)
	defer os.Remove(filename)
	version := int64(1)

	m := make(map[string]models.HTTPHardenings)
	v, profile, err := clientTest.GetStructuredHTTPHardening("hardening2", "")
	require.NoError(t, err)
	require.Equal(t, version, v)
	m["hardening2"] = models.HTTPHardenings{profile}
	checkStructuredHTTPHardenings(t, m)

	_, _, err = clientTest.GetStructuredHTTPHardening("doesnotexist", "")
	require.Error(t, err)
}

func TestCreateEditStructuredHTTPHardening(t *testing.T) {
	clientTest, filename, err := getTestClient()
	require.NoError(t, err)
	defer os.Remove(filename)
	version := int64(1)

	profile := &models.HTTPHardening{
		Name: "hardening3",
		From: "hardening1",
		All:  "off",
	}
	require.NoError(t, clientTest.CreateStructuredHTTPHardening(profile, "", version))
	version++

	v, got, err := clientTest.GetStructuredHTTPHardening(profile.Name, "")
	require.NoError(t, err)
	require.Equal(t, version, v)
	require.True(t, profile.Equal(*got), cmp.Diff(*profile, *got))

	profile.From = "hardening2"
	require.NoError(t, clientTest.EditStructuredHTTPHardening(profile.Name, profile, "", version))
	version++

	v, got, err = clientTest.GetStructuredHTTPHardening(profile.Name, "")
	require.NoError(t, err)
	require.Equal(t, version, v)
	require.True(t, profile.Equal(*got), cmp.Diff(*profile, *got))
}

func checkStructuredHTTPHardenings(t *testing.T, got map[string]models.HTTPHardenings) {
	t.Helper()
	exp := httpHardeningExpectation()
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
