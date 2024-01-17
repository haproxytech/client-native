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
	"github.com/haproxytech/client-native/v5/configuration"
	"github.com/haproxytech/client-native/v5/misc"
	"github.com/haproxytech/client-native/v5/models"
	"github.com/stretchr/testify/require"
)

func hTTPErrorRulesExpectation() map[string]models.HTTPErrorRules {
	initStructuredExpected()
	res := StructuredToHTTPErrorRuleMap()
	return res
}

func TestGetHTTPErrorRules(t *testing.T) { //nolint:gocognit,gocyclo
	mr := make(map[string]models.HTTPErrorRules)

	v, checks, err := clientTest.GetHTTPErrorRules(configuration.FrontendParentName, "test", "")
	if err != nil {
		t.Error(err.Error())
	}
	if len(checks) != 1 {
		t.Errorf("%v http-error rules returned, expected 1", len(checks))
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}
	mr["frontend/test"] = checks

	_, checks, err = clientTest.GetHTTPErrorRules(configuration.DefaultsParentName, "", "")
	if err != nil {
		t.Error(err.Error())
	}
	mr[configuration.DefaultsParentName] = checks

	_, checks, err = clientTest.GetHTTPErrorRules(configuration.BackendParentName, "test_2", "")
	if err != nil {
		t.Error(err.Error())
	}
	mr["backend/test_2"] = checks

	checkErrorRules(t, mr)
}

func checkErrorRules(t *testing.T, got map[string]models.HTTPErrorRules) {
	exp := hTTPErrorRulesExpectation()
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

func TestGetHTTPErrorRule(t *testing.T) {
	v, check, err := clientTest.GetHTTPErrorRule(0, configuration.BackendParentName, "test_2", "")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	if *check.Index != 0 {
		t.Errorf("http-error rule index not 0: %v", *check.Index)
	}
	if check.Type != "status" {
		t.Errorf("%v: Action not status: %v", *check.Index, check.Type)
	}
	if check.Status != 200 {
		t.Errorf("%v: Status not 200: %v", *check.Index, check.Type)
	}
	_, err = check.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	_, _, err = clientTest.GetHTTPErrorRule(3, configuration.BackendParentName, "test", "")
	if err == nil {
		t.Error("no http-error rules in backend section named test - expected an error")
	}

	_, check, err = clientTest.GetHTTPErrorRule(1, configuration.DefaultsParentName, "", "")
	if err != nil {
		t.Error(err.Error())
	}
	if *check.Index != 1 {
		t.Errorf("http-error rule index not 1: %v", *check.Index)
	}
	if check.Type != "status" {
		t.Errorf("%v: Action not status: %v", *check.Index, check.Type)
	}
	if check.Status != 429 {
		t.Errorf("%v: Status not 429: %v", *check.Index, check.Type)
	}
}

func TestCreateEditDeleteHTTPErrorRule(t *testing.T) {
	id := int64(1)
	r := &models.HTTPErrorRule{
		Index:               &id,
		Type:                "status",
		Status:              429,
		ReturnContentType:   misc.StringP("application/json"),
		ReturnContentFormat: "file",
		ReturnContent:       "/test/429",
		ReturnHeaders: []*models.ReturnHeader{
			{
				Name: misc.StringP("Some-Header"),
				Fmt:  misc.StringP("value"),
			},
		},
	}

	// TestCreateHTTPErrorRule
	err := clientTest.CreateHTTPErrorRule(configuration.BackendParentName, "test_2", r, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, ondiskR, err := clientTest.GetHTTPErrorRule(1, configuration.BackendParentName, "test_2", "")
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
		fmt.Printf("Created HTTP error rule: %v\n", string(ondiskJSONB))
		fmt.Printf("Given HTTP error rule: %v\n", string(givenJSONB))
		t.Error("Created HTTP error rule not equal to given HTTP error rule")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	// TestEditHTTPErrorRule
	err = clientTest.EditHTTPErrorRule(1, configuration.BackendParentName, "test_2", r, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, ondiskR, err = clientTest.GetHTTPErrorRule(1, configuration.BackendParentName, "test_2", "")
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
		fmt.Printf("Created HTTP error rule: %v\n", string(ondiskJSONB))
		fmt.Printf("Given HTTP error rule: %v\n", string(givenJSONB))
		t.Error("Created HTTP error rule not equal to given HTTP error rule")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	// TestDeleteHTTPErrorRule
	err = clientTest.DeleteHTTPErrorRule(0, configuration.FrontendParentName, "test", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ := clientTest.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	_, _, err = clientTest.GetHTTPErrorRule(0, configuration.FrontendParentName, "test", "")
	if err == nil {
		t.Error("deleting http-error rule failed - http-error rule 0 still exists")
	}

	err = clientTest.DeleteHTTPErrorRule(1, configuration.DefaultsParentName, "", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ := clientTest.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	_, _, err = clientTest.GetHTTPErrorRule(1, configuration.DefaultsParentName, "", "")
	if err == nil {
		t.Error("deleting http-error rule failed - http-error rule 1 still exists")
	}

	err = clientTest.DeleteHTTPErrorRule(3, configuration.BackendParentName, "test_2", "", version)
	if err == nil {
		t.Error("deleting http-error rule that does not exist - expected an error")
		version++
	}
}
