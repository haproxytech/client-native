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
	"slices"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/haproxytech/client-native/v6/models"
	"github.com/stretchr/testify/require"
)

func httpHardeningExpectation() map[string]models.HTTPHardenings {
	initStructuredExpected()
	res := StructuredToHTTPHardeningMap()
	for _, vs := range res {
		for _, v := range vs {
			key := v.Name
			res[key] = models.HTTPHardenings{v}
		}
	}
	return res
}

func TestGetHTTPHardenings(t *testing.T) {
	m := make(map[string]models.HTTPHardenings)
	v, profiles, err := clientTest.GetHTTPHardenings("")
	if err != nil {
		t.Error(err.Error())
	}
	if v != version {
		t.Errorf("version %v returned, expected %v", v, version)
	}
	m[""] = profiles
	checkHTTPHardenings(t, m)

	names := make([]string, len(profiles))
	for i, p := range profiles {
		names[i] = p.Name
	}
	if !slices.Contains(names, "hardening1") || !slices.Contains(names, "hardening2") {
		t.Fatalf("missing profiles: %v", names)
	}
}

func TestGetHTTPHardening(t *testing.T) {
	m := make(map[string]models.HTTPHardenings)
	v, profile, err := clientTest.GetHTTPHardening("hardening2", "")
	require.NoError(t, err)
	require.Equal(t, version, v)
	m["hardening2"] = models.HTTPHardenings{profile}
	checkHTTPHardenings(t, m)

	_, _, err = clientTest.GetHTTPHardening("doesnotexist", "")
	require.Error(t, err)
}

func checkHTTPHardenings(t *testing.T, got map[string]models.HTTPHardenings) {
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

func TestHTTPHardeningReservedName(t *testing.T) {
	require.Error(t, clientTest.CreateHTTPHardening(&models.HTTPHardening{Name: "auto"}, "", version))
	require.Error(t, clientTest.EditHTTPHardening("auto", &models.HTTPHardening{Name: "auto"}, "", version))
}

func TestCreateEditDeleteHTTPHardening(t *testing.T) {
	require := require.New(t)

	profile := &models.HTTPHardening{
		Name:                  "hardening3",
		From:                  "hardening1",
		Metadata:              map[string]any{"comment": "managed"},
		All:                   "off",
		RejectLowercaseMethod: "on",
	}
	require.NoError(clientTest.CreateHTTPHardening(profile, "", version))
	version++

	v, got, err := clientTest.GetHTTPHardening(profile.Name, "")
	require.NoError(err)
	require.Equal(version, v)
	require.True(profile.Equal(*got), cmp.Diff(*profile, *got))

	_, profiles, err := clientTest.GetHTTPHardenings("")
	require.NoError(err)
	names := make([]string, len(profiles))
	for i, p := range profiles {
		names[i] = p.Name
	}
	require.True(slices.Contains(names, "hardening1") && slices.Contains(names, "hardening2") && slices.Contains(names, "hardening3"), names)

	require.Error(clientTest.CreateHTTPHardening(profile, "", version))

	profile.From = "hardening2"
	profile.All = ""
	profile.RejectLowercaseMethod = ""
	profile.AllDelete = "on"
	require.NoError(clientTest.EditHTTPHardening(profile.Name, profile, "", version))
	version++

	_, profiles, err = clientTest.GetHTTPHardenings("")
	require.NoError(err)
	names = make([]string, len(profiles))
	for i, p := range profiles {
		names[i] = p.Name
	}
	require.True(slices.Contains(names, "hardening1") && slices.Contains(names, "hardening2") && slices.Contains(names, "hardening3"), names)

	v, got, err = clientTest.GetHTTPHardening(profile.Name, "")
	require.NoError(err)
	require.Equal(version, v)
	require.True(profile.Equal(*got), cmp.Diff(*profile, *got))

	// Path/body mismatch is a validation error, not a rename.
	require.Error(clientTest.EditHTTPHardening("other", profile, "", version))

	// Failed updates leave the profile untouched.
	failed := *profile
	failed.From = "missing"
	require.Error(clientTest.EditHTTPHardening(failed.Name, &failed, "", version))
	_, got, err = clientTest.GetHTTPHardening(profile.Name, "")
	require.NoError(err)
	require.True(profile.Equal(*got), cmp.Diff(*profile, *got))

	require.NoError(clientTest.DeleteHTTPHardening(profile.Name, "", version))
	version++

	_, _, err = clientTest.GetHTTPHardening(profile.Name, "")
	require.Error(err)
	require.Error(clientTest.DeleteHTTPHardening(profile.Name, "", version))
}
