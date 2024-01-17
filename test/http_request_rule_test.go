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
	_ "embed"
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/haproxytech/client-native/v5/configuration"
	"github.com/haproxytech/client-native/v5/misc"
	"github.com/haproxytech/client-native/v5/models"
	"github.com/stretchr/testify/require"
)

func hTTPRequestRuleExpectation() map[string]models.HTTPRequestRules {
	initStructuredExpected()
	res := StructuredToHTTPRequestRuleMap()
	// Add individual entries
	for k, vs := range res {
		for _, v := range vs {
			key := fmt.Sprintf("%s/%d", k, *v.Index)
			res[key] = models.HTTPRequestRules{v}
		}
	}
	return res
}

func TestGetHTTPRequestRules(t *testing.T) { //nolint:gocognit,gocyclo
	mrules := make(map[string]models.HTTPRequestRules)
	v, hRules, err := clientTest.GetHTTPRequestRules(configuration.FrontendParentName, "test", "")
	if err != nil {
		t.Error(err.Error())
	}
	mrules["frontend/test"] = hRules

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	_, hRules, err = clientTest.GetHTTPRequestRules(configuration.BackendParentName, "test", "")
	if err != nil {
		t.Error(err.Error())
	}
	mrules["backend/test"] = hRules

	_, hRules, err = clientTest.GetHTTPRequestRules(configuration.BackendParentName, "test_2", "")
	if err != nil {
		t.Error(err.Error())
	}
	mrules["backend/test_2"] = hRules

	checkHTTPRequestRules(t, mrules)
}

func checkHTTPRequestRules(t *testing.T, got map[string]models.HTTPRequestRules) {
	exp := hTTPRequestRuleExpectation()
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

func TestGetHTTPRequestRule(t *testing.T) {
	m := make(map[string]models.HTTPRequestRules)
	v, r, err := clientTest.GetHTTPRequestRule(0, configuration.FrontendParentName, "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}
	m["frontend/test/0"] = models.HTTPRequestRules{r}

	_, err = r.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	_, _, err = clientTest.GetHTTPRequestRule(3, configuration.BackendParentName, "test_2", "")
	if err == nil {
		t.Error("Should throw error, non existent HTTP Request Rule")
	}

	_, r, err = clientTest.GetHTTPRequestRule(0, configuration.FrontendParentName, "test_2", "")
	if err != nil {
		t.Error("Should throw error, non existent HTTP Request Rule")
	}
	m["frontend/test_2/0"] = models.HTTPRequestRules{r}

	_, r, err = clientTest.GetHTTPRequestRule(1, configuration.FrontendParentName, "test_2", "")
	if err != nil {
		t.Error("Should throw error, non existent HTTP Request Rule")
	}
	m["frontend/test_2/1"] = models.HTTPRequestRules{r}

	checkHTTPRequestRules(t, m)
}

func TestCreateEditDeleteHTTPRequestRule(t *testing.T) {
	id := int64(1)

	// TestCreateHTTPRequestRule
	var redirCode int64 = 301
	r := &models.HTTPRequestRule{
		Index:      &id,
		Type:       "redirect",
		RedirCode:  &redirCode,
		RedirValue: "http://www.%[hdr(host)]%[capture.req.uri]",
		RedirType:  "location",
	}

	err := clientTest.CreateHTTPRequestRule(configuration.FrontendParentName, "test", r, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, ondiskR, err := clientTest.GetHTTPRequestRule(1, configuration.FrontendParentName, "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(ondiskR, r) {
		fmt.Printf("Created HTTP request rule: %v\n", ondiskR)
		fmt.Printf("Given HTTP request rule: %v\n", r)
		t.Error("Created HTTP request rule not equal to given HTTP request rule")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	// TestEditHTTPRequestRule
	r = &models.HTTPRequestRule{
		Index:      &id,
		Type:       "redirect",
		RedirCode:  &redirCode,
		RedirValue: "http://www1.%[hdr(host)]%[capture.req.uri]",
		RedirType:  "scheme",
	}

	err = clientTest.EditHTTPRequestRule(1, configuration.FrontendParentName, "test", r, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, ondiskR, err = clientTest.GetHTTPRequestRule(1, configuration.FrontendParentName, "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(ondiskR, r) {
		fmt.Printf("Edited HTTP request rule: %v\n", ondiskR)
		fmt.Printf("Given HTTP request rule: %v\n", r)
		t.Error("Edited HTTP request rule not equal to given HTTP request rule")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	// TestDeleteHTTPRequest
	err = clientTest.DeleteHTTPRequestRule(47, configuration.FrontendParentName, "test", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ := clientTest.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	_, _, err = clientTest.GetHTTPRequestRule(47, "frontend", "test", "")
	if err == nil {
		t.Error("DeleteHTTPRequestRule failed, HTTP Request Rule 47 still exists")
	}

	err = clientTest.DeleteHTTPRequestRule(2, configuration.BackendParentName, "test_2", "", version)
	if err == nil {
		t.Error("Should throw error, non existent HTTP Request Rule")
		version++
	}
}

func TestSerializeHTTPRequestRule(t *testing.T) {
	testCases := []struct {
		input          models.HTTPRequestRule
		expectedResult string
	}{
		{
			input: models.HTTPRequestRule{
				Type:                models.HTTPRequestRuleTypeTrackDashSc,
				Cond:                "if",
				CondTest:            "TRUE",
				TrackScKey:          "src",
				TrackScTable:        "tr0",
				TrackScStickCounter: misc.Int64P(3),
			},
			expectedResult: "track-sc3 src table tr0 if TRUE",
		},
		{
			input: models.HTTPRequestRule{
				Type:          models.HTTPRequestRuleTypeTrackDashSc0,
				Cond:          "if",
				CondTest:      "TRUE",
				TrackSc0Key:   "src",
				TrackSc0Table: "tr0",
			},
			expectedResult: "track-sc0 src table tr0 if TRUE",
		},
		{
			input: models.HTTPRequestRule{
				Type:          models.HTTPRequestRuleTypeTrackDashSc1,
				Cond:          "if",
				CondTest:      "TRUE",
				TrackSc1Key:   "src",
				TrackSc1Table: "tr1",
			},
			expectedResult: "track-sc1 src table tr1 if TRUE",
		},
		{
			input: models.HTTPRequestRule{
				Type:          models.HTTPRequestRuleTypeTrackDashSc2,
				Cond:          "if",
				CondTest:      "TRUE",
				TrackSc2Key:   "src",
				TrackSc2Table: "tr2",
			},
			expectedResult: "track-sc2 src table tr2 if TRUE",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.expectedResult, func(t *testing.T) {
			tcpType, err := configuration.SerializeHTTPRequestRule(testCase.input)
			if err != nil {
				t.Error(err.Error())
			}

			actual := tcpType.String()
			if actual != testCase.expectedResult {
				t.Errorf("Expected %q, got: %q", testCase.expectedResult, actual)
			}
		})
	}
}
