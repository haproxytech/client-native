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
	"github.com/haproxytech/client-native/v6/configuration"
	"github.com/haproxytech/client-native/v6/misc"
	"github.com/haproxytech/client-native/v6/models"
	"github.com/stretchr/testify/require"
)

func hTTPResponseRuleExpectation() map[string]models.HTTPResponseRules {
	initStructuredExpected()
	res := StructuredToHTTPResponseRuleMap()
	// Add individual entries
	for k, vs := range res {
		for i, v := range vs {
			key := fmt.Sprintf("%s/%d", k, i)
			res[key] = models.HTTPResponseRules{v}
		}
	}
	return res
}

func TestGetHTTPResponseRules(t *testing.T) { //nolint:gocognit,gocyclo
	mrules := make(map[string]models.HTTPResponseRules)
	v, hRules, err := clientTest.GetHTTPResponseRules(configuration.FrontendParentName, "test", "")
	if err != nil {
		t.Error(err.Error())
	}
	mrules["frontend/test"] = hRules

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	_, hRules, err = clientTest.GetHTTPResponseRules(configuration.BackendParentName, "test_2", "")
	if err != nil {
		t.Error(err.Error())
	}
	mrules["backend/test_2"] = hRules

	_, hRules, err = clientTest.GetHTTPResponseRules(configuration.DefaultsParentName, "test_defaults", "")
	if err != nil {
		t.Error(err.Error())
	}
	mrules["defaults/test_defaults"] = hRules

	checkHTTPResponseRules(t, mrules)
}

func checkHTTPResponseRules(t *testing.T, got map[string]models.HTTPResponseRules) {
	exp := hTTPResponseRuleExpectation()
	for k, v := range got {
		want, ok := exp[k]
		require.True(t, ok, "k=%s", k)
		require.Equal(t, len(want), len(v), "k=%s", k)
		i := 0
		for _, w := range want {
			require.True(t, v[i].Equal(*w), "k=%s - diff %v", k, cmp.Diff(*v[i], *w))
			i++
		}
	}
}

func TestGetHTTPResponseRule(t *testing.T) {
	m := make(map[string]models.HTTPResponseRules)
	v, r, err := clientTest.GetHTTPResponseRule(0, configuration.FrontendParentName, "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}
	m["frontend/test/0"] = models.HTTPResponseRules{r}

	_, err = r.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	_, _, err = clientTest.GetHTTPResponseRule(3, configuration.BackendParentName, "test2", "")
	if err == nil {
		t.Error("Should throw error, non existent HTTPResponse Rule")
	}

	_, r, err = clientTest.GetHTTPResponseRule(0, configuration.FrontendParentName, "test_2", "")
	if err != nil {
		t.Error(err.Error())
	}
	m["frontend/test_2/0"] = models.HTTPResponseRules{r}

	checkHTTPResponseRules(t, m)
}

func TestCreateEditDeleteHTTPResponseRule(t *testing.T) {
	id := int64(1)
	// TestCreateHTTPResponseRule
	r := &models.HTTPResponseRule{
		Type:     "set-log-level",
		LogLevel: "alert",
	}

	err := clientTest.CreateHTTPResponseRule(id, configuration.FrontendParentName, "test", r, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, ondiskR, err := clientTest.GetHTTPResponseRule(1, configuration.FrontendParentName, "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(ondiskR, r) {
		fmt.Printf("Created HTTP response rule: %v\n", ondiskR)
		fmt.Printf("Given HTTP response rule: %v\n", r)
		t.Error("Created HTTP response rule not equal to given HTTP response rule")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	// TestEditHTTPResponseRule
	r = &models.HTTPResponseRule{
		Type:     "set-log-level",
		LogLevel: "warning",
	}

	err = clientTest.EditHTTPResponseRule(1, configuration.FrontendParentName, "test", r, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, ondiskR, err = clientTest.GetHTTPResponseRule(1, configuration.FrontendParentName, "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(ondiskR, r) {
		fmt.Printf("Edited HTTP response rule: %v\n", ondiskR)
		fmt.Printf("Given HTTP response rule: %v\n", r)
		t.Error("Edited HTTP response rule not equal to given HTTP response rule")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	// TestDeleteHTTPResponse
	_, rules, err := clientTest.GetHTTPResponseRules(configuration.FrontendParentName, "test", "")
	if err != nil {
		t.Error(err.Error())
	}
	N := int64(len(rules)) - 1
	err = clientTest.DeleteHTTPResponseRule(N, configuration.FrontendParentName, "test", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ := clientTest.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	_, _, err = clientTest.GetHTTPResponseRule(N, configuration.FrontendParentName, "test", "")
	if err == nil {
		t.Errorf("DeleteHTTPResponseRule failed, HTTPResponse Rule %d still exists", N)
	}

	err = clientTest.DeleteHTTPResponseRule(2, configuration.BackendParentName, "test_2", "", version)
	if err == nil {
		t.Error("Should throw error, non existent HTTPResponse Rule")
		version++
	}
}

func TestSerializeHTTPResponseRule(t *testing.T) {
	testCases := []struct {
		input          models.HTTPResponseRule
		expectedResult string
	}{
		{
			input: models.HTTPResponseRule{
				Type:                models.HTTPResponseRuleTypeTrackDashSc,
				Cond:                "if",
				CondTest:            "TRUE",
				TrackScKey:          "src",
				TrackScTable:        "tr0",
				TrackScStickCounter: misc.Int64P(3),
			},
			expectedResult: "track-sc3 src table tr0 if TRUE",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.expectedResult, func(t *testing.T) {
			tcpType, err := configuration.SerializeHTTPResponseRule(testCase.input)
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
