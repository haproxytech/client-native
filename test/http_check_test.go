// Copyright 2022 HAProxy Technologies
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
	_ "embed"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/haproxytech/client-native/v6/configuration"
	"github.com/haproxytech/client-native/v6/misc"
	"github.com/haproxytech/client-native/v6/models"
	"github.com/stretchr/testify/require"
)

func hTTPCheckExpectation() map[string]models.HTTPChecks {
	initStructuredExpected()
	res := StructuredToHTTPCheckMap()
	// Add individual entries
	for k, vs := range res {
		for _, v := range vs {
			key := fmt.Sprintf("%s/%d", k, *v.Index)
			res[key] = models.HTTPChecks{v}
		}
	}
	return res
}

func TestGetHTTPChecks(t *testing.T) { //nolint:gocognit,gocyclo
	mchecks := make(map[string]models.HTTPChecks)

	v, checks, err := clientTest.GetHTTPChecks(configuration.BackendParentName, "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if len(checks) != 14 {
		t.Errorf("%v http checks returned, expected 14", len(checks))
	}
	mchecks["backend/test"] = checks

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	_, checks, err = clientTest.GetHTTPChecks(configuration.DefaultsParentName, "", "")
	if err != nil {
		t.Error(err.Error())
	}
	mchecks[configuration.DefaultsParentName+"/unnamed_defaults_1"] = checks

	_, checks, err = clientTest.GetHTTPChecks(configuration.BackendParentName, "test_2", "")
	if err != nil {
		t.Error(err.Error())
	}
	mchecks["backend/test_2"] = checks

	checkHttpChecks(t, mchecks)
}

func checkHttpChecks(t *testing.T, got map[string]models.HTTPChecks) {
	exp := hTTPCheckExpectation()
	for k, v := range got {
		want, ok := exp[k]
		require.True(t, ok, "k=%s", k)
		require.Equal(t, len(want), len(v), "k=%s", k)
		for _, g := range v {
			for _, w := range want {
				if *g.Index == *w.Index {
					require.True(t, g.Equal(*w), "k=%s - diff %v", k, cmp.Diff(*g, *w))
					break
				}
			}
		}
	}
}

func TestGetHTTPCheck(t *testing.T) {
	m := make(map[string]models.HTTPChecks)
	v, check, err := clientTest.GetHTTPCheck(0, configuration.BackendParentName, "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}
	m["backend/test/0"] = models.HTTPChecks{check}

	_, err = check.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	_, _, err = clientTest.GetHTTPCheck(3, configuration.BackendParentName, "test_2", "")
	if err == nil {
		t.Error("Should throw error, non existent HTTP Request Rule")
	}

	_, check, err = clientTest.GetHTTPCheck(0, configuration.DefaultsParentName, "", "")
	if err != nil {
		t.Error(err.Error())
	}
	m[configuration.DefaultsParentName+"/unnamed_defaults_1/0"] = models.HTTPChecks{check}
	checkHttpChecks(t, m)
}

func TestCreateEditDeleteHTTPCheck(t *testing.T) {
	id := int64(1)

	// TestCreateHTTPCheck
	r := &models.HTTPCheck{
		Index:        &id,
		Type:         "send",
		Method:       "GET",
		Version:      "HTTP/1.1",
		URI:          "/",
		CheckHeaders: []*models.ReturnHeader{},
	}

	err := clientTest.CreateHTTPCheck(configuration.BackendParentName, "test", r, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, ondiskR, err := clientTest.GetHTTPCheck(1, configuration.BackendParentName, "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	var givenJSONB []byte
	givenJSONB, err = r.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	var ondiskJSONB []byte
	ondiskJSONB, err = ondiskR.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	if string(givenJSONB) != string(ondiskJSONB) {
		fmt.Printf("Created HTTP check: %v\n", string(ondiskJSONB))
		fmt.Printf("Given HTTP check: %v\n", string(givenJSONB))
		t.Error("Created HTTP check not equal to given HTTP check")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	// TestEditHTTPRequestRule
	r = &models.HTTPCheck{
		Index:   &id,
		Type:    "send",
		Method:  "GET",
		Version: "HTTP/1.1",
		URI:     "/",
		CheckHeaders: []*models.ReturnHeader{
			{
				Name: misc.StringP("Host"),
				Fmt:  misc.StringP("google.com"),
			},
		},
	}

	err = clientTest.EditHTTPCheck(1, configuration.BackendParentName, "test", r, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, ondiskR, err = clientTest.GetHTTPCheck(1, configuration.BackendParentName, "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	givenJSONB, err = r.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	ondiskJSONB, err = ondiskR.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	if string(givenJSONB) != string(ondiskJSONB) {
		fmt.Printf("Created HTTP check: %v\n", string(ondiskJSONB))
		fmt.Printf("Given HTTP check: %v\n", string(givenJSONB))
		t.Error("Created HTTP check not equal to given HTTP check")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	// TestDeleteHTTPRequest
	err = clientTest.DeleteHTTPCheck(14, configuration.BackendParentName, "test", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ := clientTest.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	_, _, err = clientTest.GetHTTPCheck(14, configuration.BackendParentName, "test", "")
	if err == nil {
		t.Error("DeleteHTTPCheck failed, HTTP check 13 still exists")
	}

	err = clientTest.DeleteHTTPCheck(5, configuration.BackendParentName, "test_2", "", version)
	if err == nil {
		t.Error("Should throw error, non existent HTTP Check")
		version++
	}
}
