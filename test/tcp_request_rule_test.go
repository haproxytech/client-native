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
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/haproxytech/client-native/v6/configuration"
	"github.com/haproxytech/client-native/v6/configuration/options"
	"github.com/haproxytech/client-native/v6/misc"
	"github.com/haproxytech/client-native/v6/models"
	"github.com/stretchr/testify/require"
)

func tcpRequestRuleExpectation() map[string]models.TCPRequestRules {
	initStructuredExpected()
	res := StructuredToTCPRequestRuleMap()
	// Add individual entries
	for k, vs := range res {
		for i, v := range vs {
			key := fmt.Sprintf("%s/%d", k, i)
			res[key] = models.TCPRequestRules{v}
		}
	}
	return res
}

func TestGetTCPRequestRules(t *testing.T) { //nolint:gocognit,gocyclo
	mrules := make(map[string]models.TCPRequestRules)
	v, tRules, err := clientTest.GetTCPRequestRules(configuration.FrontendParentName, "test", "")
	if err != nil {
		t.Error(err.Error())
	}
	mrules["frontend/test"] = tRules

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	_, tRules, err = clientTest.GetTCPRequestRules(configuration.BackendParentName, "test_2", "")
	if err != nil {
		t.Error(err.Error())
	}
	mrules["backend/test_2"] = tRules

	_, tRules, err = clientTest.GetTCPRequestRules(configuration.DefaultsParentName, "test_defaults", "")
	if err != nil {
		t.Error(err.Error())
	}
	mrules["defaults/test_defaults"] = tRules

	checkTCPRequestRules(t, mrules)
}

func TestGetTCPRequestRule(t *testing.T) {
	m := make(map[string]models.TCPRequestRules)

	v, r, err := clientTest.GetTCPRequestRule(0, configuration.FrontendParentName, "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}
	m["frontend/test/0"] = models.TCPRequestRules{r}

	checkTCPRequestRules(t, m)

	_, err = r.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	_, _, err = clientTest.GetTCPRequestRule(3, configuration.BackendParentName, "test_2", "")
	if err == nil {
		t.Error("Should throw error, non existent TCP Request Rule")
	}
}

func checkTCPRequestRules(t *testing.T, got map[string]models.TCPRequestRules) {
	exp := tcpRequestRuleExpectation()
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

func TestCreateEditDeleteTCPRequestRule(t *testing.T) {
	id := int64(4)
	tOut := int64(1000)
	// TestCreateTCPRequestRule
	r := &models.TCPRequestRule{
		Type:    "inspect-delay",
		Timeout: &tOut,
	}

	err := clientTest.CreateTCPRequestRule(id, configuration.FrontendParentName, "test", r, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, ondiskR, err := clientTest.GetTCPRequestRule(4, configuration.FrontendParentName, "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(ondiskR, r) {
		fmt.Printf("Created TCP request rule: %v\n", ondiskR)
		fmt.Printf("Given TCP request rule: %v\n", r)
		t.Error("Created TCP request rule not equal to given TCP request rule")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	// TestEditTCPRequestRule
	r = &models.TCPRequestRule{
		Type:     "connection",
		Action:   "accept",
		Cond:     "if",
		CondTest: "FALSE",
	}

	err = clientTest.EditTCPRequestRule(4, configuration.FrontendParentName, "test", r, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, ondiskR, err = clientTest.GetTCPRequestRule(4, configuration.FrontendParentName, "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(ondiskR, r) {
		fmt.Printf("Edited TCP request rule: %v\n", ondiskR)
		fmt.Printf("Given TCP request rule: %v\n", r)
		t.Error("Edited TCP request rule not equal to given TCP request rule")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	// TestDeleteTCPRequest
	_, rules, err := clientTest.GetTCPRequestRules(configuration.FrontendParentName, "test", "")
	if err != nil {
		t.Error(err.Error())
	}
	N := int64(len(rules)) - 1
	err = clientTest.DeleteTCPRequestRule(N, configuration.FrontendParentName, "test", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ := clientTest.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	_, _, err = clientTest.GetTCPRequestRule(N, configuration.FrontendParentName, "test", "")
	if err == nil {
		t.Errorf("DeleteTCPRequestRule failed, TCP Request Rule %d still exists", N)
	}

	err = clientTest.DeleteTCPRequestRule(27, configuration.BackendParentName, "test_2", "", version)
	if err == nil {
		t.Error("Should throw error, non existent TCP Request Rule")
		version++
	}
}

func TestSerializeTCPRequestRule(t *testing.T) {
	testCases := []struct {
		input          models.TCPRequestRule
		expectedResult string
	}{
		{
			input: models.TCPRequestRule{
				Type:     models.TCPRequestRuleTypeConnection,
				Action:   models.TCPRequestRuleActionSilentDashDrop,
				Cond:     "if",
				CondTest: "FALSE",
			},
			expectedResult: "connection silent-drop if FALSE",
		},
		{
			input: models.TCPRequestRule{
				Type:     models.TCPRequestRuleTypeContent,
				Action:   models.TCPRequestRuleActionSilentDashDrop,
				Cond:     "if",
				CondTest: "FALSE",
			},
			expectedResult: "content silent-drop if FALSE",
		},
		{
			input: models.TCPRequestRule{
				Type:     models.TCPRequestRuleTypeSession,
				Action:   models.TCPRequestRuleActionSilentDashDrop,
				Cond:     "if",
				CondTest: "FALSE",
			},
			expectedResult: "session silent-drop if FALSE",
		},
		{
			input: models.TCPRequestRule{
				Type:     models.TCPRequestRuleTypeConnection,
				Action:   models.TCPRequestRuleActionExpectDashProxy,
				Cond:     "if",
				CondTest: "{ src 1.2.3.4 5.6.7.8 }",
			},
			expectedResult: "connection expect-proxy layer4 if { src 1.2.3.4 5.6.7.8 }",
		},
		{
			input: models.TCPRequestRule{
				Type:              models.TCPRequestRuleTypeConnection,
				Action:            models.TCPRequestRuleActionTrackDashSc,
				Cond:              "if",
				CondTest:          "TRUE",
				TrackKey:          "src",
				TrackTable:        "tr0",
				TrackStickCounter: misc.Int64P(3),
			},
			expectedResult: "connection track-sc3 src table tr0 if TRUE",
		},
		{
			input: models.TCPRequestRule{
				Type:              models.TCPRequestRuleTypeContent,
				Action:            models.TCPRequestRuleActionTrackDashSc,
				Cond:              "if",
				CondTest:          "TRUE",
				TrackKey:          "src",
				TrackTable:        "tr0",
				TrackStickCounter: misc.Int64P(3),
			},
			expectedResult: "content track-sc3 src table tr0 if TRUE",
		},
		{
			input: models.TCPRequestRule{
				Type:              models.TCPRequestRuleTypeSession,
				Action:            models.TCPRequestRuleActionTrackDashSc,
				Cond:              "if",
				CondTest:          "TRUE",
				TrackKey:          "src",
				TrackTable:        "tr0",
				TrackStickCounter: misc.Int64P(3),
			},
			expectedResult: "session track-sc3 src table tr0 if TRUE",
		},
		{
			input: models.TCPRequestRule{
				Type:       models.TCPRequestRuleTypeSession,
				Action:     models.TCPRequestRuleActionAttachDashSrv,
				ServerName: "srv8",
				Expr:       "haproxy.org",
			},
			expectedResult: "session attach-srv srv8 name haproxy.org",
		},
		{
			input: models.TCPRequestRule{
				Type:       models.TCPRequestRuleTypeSession,
				Action:     models.TCPRequestRuleActionAttachDashSrv,
				ServerName: "srv8",
			},
			expectedResult: "session attach-srv srv8",
		},
		{
			input: models.TCPRequestRule{
				Type:       models.TCPRequestRuleTypeSession,
				Action:     models.TCPRequestRuleActionAttachDashSrv,
				ServerName: "srv8",
				Cond:       "unless",
				CondTest:   "limit_exceeded",
			},
			expectedResult: "session attach-srv srv8 unless limit_exceeded",
		},
		{
			input: models.TCPRequestRule{
				Type:     models.TCPRequestRuleTypeContent,
				Action:   models.TCPRequestRuleActionSetDashBcDashMark,
				Expr:     "0xffff",
				Cond:     "if",
				CondTest: "TRUE",
			},
			expectedResult: "content set-bc-mark 0xffff if TRUE",
		},
		{
			input: models.TCPRequestRule{
				Type:      models.TCPRequestRuleTypeContent,
				Action:    models.TCPRequestRuleActionSetDashBcDashMark,
				MarkValue: "123",
			},
			expectedResult: "content set-bc-mark 123",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.expectedResult, func(t *testing.T) {
			opt := new(options.ConfigurationOptions)
			opt.PreferredTimeSuffix = "d"
			tcpType, err := configuration.SerializeTCPRequestRule(testCase.input, opt)
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
