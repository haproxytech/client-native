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
	"github.com/haproxytech/client-native/v6/misc"
	"github.com/haproxytech/client-native/v6/models"
	"github.com/stretchr/testify/require"
)

func tcpRequestRuleExpectation() map[string]models.TCPRequestRules {
	initStructuredExpected()
	res := StructuredToTCPRequestRuleMap()
	// Add individual entries
	for k, vs := range res {
		for _, v := range vs {
			key := fmt.Sprintf("%s/%d", k, *v.Index)
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

func TestCreateEditDeleteTCPRequestRule(t *testing.T) {
	id := int64(4)
	tOut := int64(1000)
	// TestCreateTCPRequestRule
	r := &models.TCPRequestRule{
		Index:   &id,
		Type:    "inspect-delay",
		Timeout: &tOut,
	}

	err := clientTest.CreateTCPRequestRule(configuration.FrontendParentName, "test", r, "", version)
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
		Index:    &id,
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
	N := int64(41) // number of tcp-request rules in frontend "test"
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
				Type:       models.TCPRequestRuleTypeConnection,
				Action:     models.TCPRequestRuleActionTrackDashSc0,
				Cond:       "if",
				CondTest:   "TRUE",
				TrackKey:   "src",
				TrackTable: "tr0",
			},
			expectedResult: "connection track-sc0 src table tr0 if TRUE",
		},
		{
			input: models.TCPRequestRule{
				Type:       models.TCPRequestRuleTypeConnection,
				Action:     models.TCPRequestRuleActionTrackDashSc1,
				Cond:       "if",
				CondTest:   "TRUE",
				TrackKey:   "src",
				TrackTable: "tr1",
			},
			expectedResult: "connection track-sc1 src table tr1 if TRUE",
		},
		{
			input: models.TCPRequestRule{
				Type:       models.TCPRequestRuleTypeConnection,
				Action:     models.TCPRequestRuleActionTrackDashSc2,
				Cond:       "if",
				CondTest:   "TRUE",
				TrackKey:   "src",
				TrackTable: "tr2",
			},
			expectedResult: "connection track-sc2 src table tr2 if TRUE",
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
				Type:       models.TCPRequestRuleTypeContent,
				Action:     models.TCPRequestRuleActionTrackDashSc0,
				Cond:       "if",
				CondTest:   "TRUE",
				TrackKey:   "src",
				TrackTable: "tr0",
			},
			expectedResult: "content track-sc0 src table tr0 if TRUE",
		},
		{
			input: models.TCPRequestRule{
				Type:       models.TCPRequestRuleTypeContent,
				Action:     models.TCPRequestRuleActionTrackDashSc1,
				Cond:       "if",
				CondTest:   "TRUE",
				TrackKey:   "src",
				TrackTable: "tr1",
			},
			expectedResult: "content track-sc1 src table tr1 if TRUE",
		},
		{
			input: models.TCPRequestRule{
				Type:       models.TCPRequestRuleTypeContent,
				Action:     models.TCPRequestRuleActionTrackDashSc2,
				Cond:       "if",
				CondTest:   "TRUE",
				TrackKey:   "src",
				TrackTable: "tr2",
			},
			expectedResult: "content track-sc2 src table tr2 if TRUE",
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
				Action:     models.TCPRequestRuleActionTrackDashSc0,
				Cond:       "if",
				CondTest:   "TRUE",
				TrackKey:   "src",
				TrackTable: "tr0",
			},
			expectedResult: "session track-sc0 src table tr0 if TRUE",
		},
		{
			input: models.TCPRequestRule{
				Type:       models.TCPRequestRuleTypeSession,
				Action:     models.TCPRequestRuleActionTrackDashSc1,
				Cond:       "if",
				CondTest:   "TRUE",
				TrackKey:   "src",
				TrackTable: "tr1",
			},
			expectedResult: "session track-sc1 src table tr1 if TRUE",
		},
		{
			input: models.TCPRequestRule{
				Type:       models.TCPRequestRuleTypeSession,
				Action:     models.TCPRequestRuleActionTrackDashSc2,
				Cond:       "if",
				CondTest:   "TRUE",
				TrackKey:   "src",
				TrackTable: "tr2",
			},
			expectedResult: "session track-sc2 src table tr2 if TRUE",
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
			tcpType, err := configuration.SerializeTCPRequestRule(testCase.input)
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
