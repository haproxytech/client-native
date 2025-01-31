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
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/haproxytech/client-native/v6/configuration"
	"github.com/haproxytech/client-native/v6/models"
	"github.com/stretchr/testify/require"
)

func logProfileExpectation() map[string]models.LogProfiles {
	initStructuredExpected()
	res := StructuredToLogProfileMap()
	// Add individual entries
	for _, vs := range res {
		for _, v := range vs {
			key := v.Name
			res[key] = models.LogProfiles{v}
		}
	}
	return res
}

func TestGetLogProfiles(t *testing.T) {
	m := make(map[string]models.LogProfiles)
	v, rings, err := clientTest.GetLogProfiles("")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("version %v returned, expected %v", v, version)
	}
	m[""] = rings
	checkLogProfiles(t, m)
}

func TestGetLogProfile(t *testing.T) {
	m := make(map[string]models.LogProfiles)

	v, r, err := clientTest.GetLogProfile("logp1", "")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("version %v returned, expected %v", v, version)
	}
	m["logp1"] = models.LogProfiles{r}
	checkLogProfiles(t, m)

	_, _, err = clientTest.GetLogProfile("doesnotexist", "")
	if err == nil {
		t.Error("should throw error, non existent rings section")
	}
}

func checkLogProfiles(t *testing.T, got map[string]models.LogProfiles) {
	exp := logProfileExpectation()

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

func TestCreateEditDeleteLogProfile(t *testing.T) {
	require := require.New(t)

	r := &models.LogProfile{
		Name:   "lp1",
		LogTag: "tagx",
		Steps: models.LogProfileSteps{
			{
				Step:   "connect",
				Drop:   "disabled",
				Format: "conn: %ci",
			},
			{
				Step:   "error",
				Drop:   "disabled",
				Format: "ERR: %lol",
				Sd:     "oops",
			},
			{
				Step: "close",
				Drop: "enabled",
			},
		},
	}
	err := clientTest.CreateLogProfile(r, "", version)
	require.NoError(err)
	version++

	v, profile, err := clientTest.GetLogProfile(r.Name, "")
	require.NoError(err)
	require.Equal(v, version)

	require.Equal(profile.Name, r.Name)
	require.Equal(profile.LogTag, r.LogTag)
	require.Equal(len(profile.Steps), len(r.Steps))
	require.Equal(profile, r)

	err = clientTest.CreateLogProfile(r, "", version)
	if err == nil {
		t.Error("should throw error log-profile already exists")
		version++
	}

	// Edit
	r.Steps[1].Sd = "new sd"
	r.Steps = append(r.Steps, &models.LogProfileStep{
		Step:   "http-req",
		Drop:   "disabled",
		Format: "http req detected",
	})
	err = clientTest.EditLogProfile(r.Name, r, "", version)
	require.NoError(err)
	version++

	v, profile, err = clientTest.GetLogProfile(r.Name, "")
	require.NoError(err)
	require.Equal(v, version)
	require.Equal(profile, r)

	// Delete
	err = clientTest.DeleteLogProfile(r.Name, "", version)
	require.NoError(err)
	version++

	if v, _ := clientTest.GetVersion(""); v != version {
		t.Error("version not incremented")
	}

	err = clientTest.DeleteLogProfile(r.Name, "", 999999)
	if err != nil {
		if confErr, ok := err.(*configuration.ConfError); ok {
			if !confErr.Is(configuration.ErrVersionMismatch) {
				t.Error("should throw configuration.ErrVersionMismatch error")
			}
		} else {
			t.Error("should throw configuration.ErrVersionMismatch error")
		}
	}
	_, _, err = clientTest.GetLogProfile(r.Name, "")
	if err == nil {
		t.Errorf("deleteLogProfile failed, '%s' still exists", r.Name)
	}

	err = clientTest.DeleteLogProfile("doesnotexist", "", version)
	if err == nil {
		t.Error("should throw error, non existent log-profile")
		version++
	}
}
