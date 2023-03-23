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

package configuration

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/haproxytech/client-native/v4/models"
)

func TestGetTCPRequestRules(t *testing.T) { //nolint:gocognit,gocyclo
	v, tRules, err := clientTest.GetTCPRequestRules("frontend", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if len(tRules) != 24 {
		t.Errorf("%v tcp request rules returned, expected 24", len(tRules))
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	for _, r := range tRules {
		switch *r.Index {
		case 0:
			if r.Type != "connection" {
				t.Errorf("%v: Type not connection: %v", *r.Index, r.Type)
			}
			if r.Action != "accept" {
				t.Errorf("%v: Action not accept: %v", *r.Index, r.Action)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "TRUE" {
				t.Errorf("%v: CondTest not src TRUE: %v", *r.Index, r.CondTest)
			}
		case 1:
			if r.Type != "connection" {
				t.Errorf("%v: Type not connection: %v", *r.Index, r.Type)
			}
			if r.Action != "reject" {
				t.Errorf("%v: Action not reject: %v", *r.Index, r.Action)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "FALSE" {
				t.Errorf("%v: CondTest not src FALSE: %v", *r.Index, r.CondTest)
			}
		case 2:
			if r.Type != "content" {
				t.Errorf("%v: Type not content: %v", *r.Index, r.Type)
			}
			if r.Action != "accept" {
				t.Errorf("%v: Action not accept: %v", *r.Index, r.Action)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "TRUE" {
				t.Errorf("%v: CondTest not src TRUE: %v", *r.Index, r.CondTest)
			}
		case 3:
			if r.Type != "content" {
				t.Errorf("%v: Type not content: %v", *r.Index, r.Type)
			}
			if r.Action != "reject" {
				t.Errorf("%v: Action not reject: %v", *r.Index, r.Action)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "FALSE" {
				t.Errorf("%v: CondTest not src FALSE: %v", *r.Index, r.CondTest)
			}
		case 4:
			if r.Type != "connection" {
				t.Errorf("%v: Type not connection: %v", *r.Index, r.Type)
			}
			if r.Action != "silent-drop" {
				t.Errorf("%v: Action not silent-drop: %v", *r.Index, r.Action)
			}
			if r.Cond != "" {
				t.Errorf("%v: Cond not blank: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "" {
				t.Errorf("%v: CondTest not blank: %v", *r.Index, r.CondTest)
			}
		case 5:
			if r.Type != "connection" {
				t.Errorf("%v: Type not connection: %v", *r.Index, r.Type)
			}
			if r.Action != "silent-drop" {
				t.Errorf("%v: Action not silent-drop: %v", *r.Index, r.Action)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "TRUE" {
				t.Errorf("%v: CondTest not src TRUE: %v", *r.Index, r.CondTest)
			}
		case 6:
			if r.Type != "connection" {
				t.Errorf("%v: Type not connection: %v", *r.Index, r.Type)
			}
			if r.Action != "lua" {
				t.Errorf("%v: Action not lua: %v", *r.Index, r.Action)
			}
			if r.LuaAction != "foo" {
				t.Errorf("%v: LuaAction not foo: %v", *r.Index, r.LuaAction)
			}
			if r.LuaParams != "param1 param2" {
				t.Errorf("%v: LuaParams not param1 param2: %v", *r.Index, r.LuaParams)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "FALSE" {
				t.Errorf("%v: CondTest not src FALSE: %v", *r.Index, r.CondTest)
			}
		case 7:
			if r.Type != "connection" {
				t.Errorf("%v: Type not connection: %v", *r.Index, r.Type)
			}
			if r.Action != "sc-add-gpc" {
				t.Errorf("%v: Type not sc-add-gpc: %v", *r.Index, r.Type)
			}
			if r.ScIdx != "0" {
				t.Errorf("%v: ScIdx not 0: %v", *r.Index, r.ScIdx)
			}
			if r.ScIncID != "1" {
				t.Errorf("%v: ScID not 1: %v", *r.Index, r.ScIncID)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "FALSE" {
				t.Errorf("%v: CondTest not FALSE: %v", *r.Index, r.CondTest)
			}
		case 8:
			if r.Type != "content" {
				t.Errorf("%v: Type not content: %v", *r.Index, r.Type)
			}
			if r.Action != "lua" {
				t.Errorf("%v: Action not lua: %v", *r.Index, r.Action)
			}
			if r.LuaAction != "foo" {
				t.Errorf("%v: LuaAction not foo: %v", *r.Index, r.LuaAction)
			}
			if r.LuaParams != "param1 param2" {
				t.Errorf("%v: LuaParams not param1 param2: %v", *r.Index, r.LuaParams)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "FALSE" {
				t.Errorf("%v: CondTest not src FALSE: %v", *r.Index, r.CondTest)
			}
		case 9:
			if r.Type != "content" {
				t.Errorf("%v: Type not content: %v", *r.Index, r.Type)
			}
			if r.Action != "sc-add-gpc" {
				t.Errorf("%v: Type not sc-add-gpc: %v", *r.Index, r.Type)
			}
			if r.ScIdx != "0" {
				t.Errorf("%v: ScIdx not 0: %v", *r.Index, r.ScIdx)
			}
			if r.ScIncID != "1" {
				t.Errorf("%v: ScID not 1: %v", *r.Index, r.ScIncID)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "FALSE" {
				t.Errorf("%v: CondTest not FALSE: %v", *r.Index, r.CondTest)
			}
		case 10:
			if r.Type != "content" {
				t.Errorf("%v: Type not content: %v", *r.Index, r.Type)
			}
			if r.Action != "set-bandwidth-limit" {
				t.Errorf("%v: Action not set-bandwidth-limit: %v", *r.Index, r.Action)
			}
			if r.BandwidthLimitName != "my-limit" {
				t.Errorf("%v: BandwidthLimitName not my-limit: %v", *r.Index, r.BandwidthLimitName)
			}
			if r.BandwidthLimitLimit != "1m" {
				t.Errorf("%v: BandwidthLimitLimit not 1m: %v", *r.Index, r.BandwidthLimitLimit)
			}
			if r.BandwidthLimitPeriod != "10s" {
				t.Errorf("%v: BandwidthLimitPeriod not 10s: %v", *r.Index, r.BandwidthLimitPeriod)
			}
		case 11:
			if r.Type != "content" {
				t.Errorf("%v: Type not content: %v", *r.Index, r.Type)
			}
			if r.Action != "set-bandwidth-limit" {
				t.Errorf("%v: Action not set-bandwidth-limit: %v", *r.Index, r.Action)
			}
			if r.BandwidthLimitName != "my-limit-reverse" {
				t.Errorf("%v: BandwidthLimitName not my-limit-reverse: %v", *r.Index, r.BandwidthLimitName)
			}
			if r.BandwidthLimitLimit != "2m" {
				t.Errorf("%v: BandwidthLimitLimit not 2m: %v", *r.Index, r.BandwidthLimitLimit)
			}
			if r.BandwidthLimitPeriod != "20s" {
				t.Errorf("%v: BandwidthLimitPeriod no 20s: %v", *r.Index, r.BandwidthLimitPeriod)
			}
		case 12:
			if r.Type != "content" {
				t.Errorf("%v: Type not content: %v", *r.Index, r.Type)
			}
			if r.Action != "set-bandwidth-limit" {
				t.Errorf("%v: Action not set-bandwidth-limit: %v", *r.Index, r.Action)
			}
			if r.BandwidthLimitName != "my-limit-cond" {
				t.Errorf("%v: BandwidthLimitName not my-limit-cond: %v", *r.Index, r.BandwidthLimitName)
			}
			if r.BandwidthLimitLimit != "3m" {
				t.Errorf("%v: BandwidthLimitLimit not 3m: %v", *r.Index, r.BandwidthLimitLimit)
			}
			if r.BandwidthLimitPeriod != "" {
				t.Errorf("%v: BandwidthLimitPeriod not empty", *r.Index)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "FALSE" {
				t.Errorf("%v: CondTest not FALSE: %v", *r.Index, r.CondTest)
			}
		case 13:
			if r.Type != "connection" {
				t.Errorf("%v: Type not connection: %v", *r.Index, r.Type)
			}
			if r.Action != "set-mark" {
				t.Errorf("%v: Action not set-mark: %v", *r.Index, r.Action)
			}
			if r.MarkValue != "0x1Ab" {
				t.Errorf("%v: MarkValue not 0x1Ab: %v", *r.Index, r.MarkValue)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "FALSE" {
				t.Errorf("%v: CondTest not FALSE: %v", *r.Index, r.CondTest)
			}
		case 14:
			if r.Type != "connection" {
				t.Errorf("%v: Type not connection: %v", *r.Index, r.Type)
			}
			if r.Action != "set-src-port" {
				t.Errorf("%v: Action not set-src-port: %v", *r.Index, r.Action)
			}
			if r.Expr != "hdr(port)" {
				t.Errorf("%v: Expr not hdr(port): %v", *r.Index, r.Expr)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "FALSE" {
				t.Errorf("%v: CondTest not FALSE: %v", *r.Index, r.CondTest)
			}
		case 15:
			if r.Type != "connection" {
				t.Errorf("%v: Type not connection: %v", *r.Index, r.Type)
			}
			if r.Action != "set-tos" {
				t.Errorf("%v: Action not set-tos: %v", *r.Index, r.Action)
			}
			if r.TosValue != "1" {
				t.Errorf("%v: TosValue not 1: %v", *r.Index, r.TosValue)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "FALSE" {
				t.Errorf("%v: CondTest not FALSE: %v", *r.Index, r.CondTest)
			}
		case 16:
			if r.Type != "content" {
				t.Errorf("%v: Type not content: %v", *r.Index, r.Type)
			}
			if r.Action != "set-log-level" {
				t.Errorf("%v: Action not set-log-level: %v", *r.Index, r.Action)
			}
			if r.LogLevel != "silent" {
				t.Errorf("%v: LogLevel not silent %v", *r.Index, r.LogLevel)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "FALSE" {
				t.Errorf("%v: CondTest not FALSE: %v", *r.Index, r.CondTest)
			}
		case 17:
			if r.Type != "content" {
				t.Errorf("%v: Type not content: %v", *r.Index, r.Type)
			}
			if r.Action != "set-mark" {
				t.Errorf("%v: Action not set-mark: %v", *r.Index, r.Action)
			}
			if r.MarkValue != "0x1Ac" {
				t.Errorf("%v: MarkValue not 0x1Ac: %v", *r.Index, r.MarkValue)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "FALSE" {
				t.Errorf("%v: CondTest not FALSE: %v", *r.Index, r.CondTest)
			}
		case 18:
			if r.Type != "content" {
				t.Errorf("%v: Type not content: %v", *r.Index, r.Type)
			}
			if r.Action != "set-nice" {
				t.Errorf("%v: Action not set-nice: %v", *r.Index, r.Action)
			}
			if r.NiceValue != 2 {
				t.Errorf("%v: NiceValue not 2 %v", *r.Index, r.NiceValue)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "FALSE" {
				t.Errorf("%v: CondTest not FALSE: %v", *r.Index, r.CondTest)
			}
		case 19:
			if r.Type != "content" {
				t.Errorf("%v: Type not content: %v", *r.Index, r.Type)
			}
			if r.Action != "set-src-port" {
				t.Errorf("%v: Action not set-src-port: %v", *r.Index, r.Action)
			}
			if r.Expr != "hdr(port)" {
				t.Errorf("%v: Expr not hdr(port): %v", *r.Index, r.Expr)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "FALSE" {
				t.Errorf("%v: CondTest not FALSE: %v", *r.Index, r.CondTest)
			}
		case 20:
			if r.Type != "content" {
				t.Errorf("%v: Type not content: %v", *r.Index, r.Type)
			}
			if r.Action != "set-tos" {
				t.Errorf("%v: Action not set-tos: %v", *r.Index, r.Action)
			}
			if r.TosValue != "3" {
				t.Errorf("%v: TosValue not 3: %v", *r.Index, r.TosValue)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "FALSE" {
				t.Errorf("%v: CondTest not FALSE: %v", *r.Index, r.CondTest)
			}
		case 21:
			if r.Type != "content" {
				t.Errorf("%v: Type not content: %v", *r.Index, r.Type)
			}
			if r.Action != "set-var-fmt" {
				t.Errorf("%v: Action not set-var-fmt: %v", *r.Index, r.Action)
			}
			if r.VarName != "tn" {
				t.Errorf("%v: VarName not tn: %v", *r.Index, r.VarName)
			}
			if r.VarScope != "req" {
				t.Errorf("%v: VarScope not req: %v", *r.Index, r.VarScope)
			}
			if r.VarFormat != "ssl_c_s_tn" {
				t.Errorf("%v: VarFormat not ssl_c_s_tn: %v", *r.Index, r.VarFormat)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "FALSE" {
				t.Errorf("%v: CondTest not FALSE: %v", *r.Index, r.CondTest)
			}
		case 22:
			if r.Type != "content" {
				t.Errorf("%v: Type not content: %v", *r.Index, r.Type)
			}
			if r.Action != "switch-mode" {
				t.Errorf("%v: Action not switch-mode: %v", *r.Index, r.Action)
			}
			if r.SwitchModeProto != "my-proto" {
				t.Errorf("%v: SwitchModeProto not my-proto: %v", *r.Index, r.SwitchModeProto)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "FALSE" {
				t.Errorf("%v: CondTest not FALSE: %v", *r.Index, r.CondTest)
			}
		case 23:
			if r.Type != "session" {
				t.Errorf("%v: Type not session: %v", *r.Index, r.Type)
			}
			if r.Action != "sc-add-gpc" {
				t.Errorf("%v: Type not sc-add-gpc: %v", *r.Index, r.Type)
			}
			if r.ScIdx != "0" {
				t.Errorf("%v: ScIdx not 0: %v", *r.Index, r.ScIdx)
			}
			if r.ScIncID != "1" {
				t.Errorf("%v: ScID not 1: %v", *r.Index, r.ScIncID)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "FALSE" {
				t.Errorf("%v: CondTest not FALSE: %v", *r.Index, r.CondTest)
			}
		default:
			t.Errorf("Expect tcp-request 0-23, %v found", *r.Index)
		}
	}

	_, tRules, err = clientTest.GetTCPRequestRules("backend", "test_2", "")
	if err != nil {
		t.Error(err.Error())
	}
	if len(tRules) > 0 {
		t.Errorf("%v TCP Request Rules returned, expected 0", len(tRules))
	}
}

func TestGetTCPRequestRule(t *testing.T) {
	v, r, err := clientTest.GetTCPRequestRule(0, "frontend", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	if r.Type != "connection" {
		t.Errorf("%v: Type not connection: %v", *r.Index, r.Type)
	}
	if r.Action != "accept" {
		t.Errorf("%v: Action not accept: %v", *r.Index, r.Action)
	}
	if r.Cond != "if" {
		t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
	}
	if r.CondTest != "TRUE" {
		t.Errorf("%v: CondTest not src TRUE: %v", *r.Index, r.CondTest)
	}

	_, err = r.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	_, _, err = clientTest.GetTCPRequestRule(3, "backend", "test_2", "")
	if err == nil {
		t.Error("Should throw error, non existant TCP Request Rule")
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

	err := clientTest.CreateTCPRequestRule("frontend", "test", r, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, ondiskR, err := clientTest.GetTCPRequestRule(4, "frontend", "test", "")
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

	err = clientTest.EditTCPRequestRule(4, "frontend", "test", r, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, ondiskR, err = clientTest.GetTCPRequestRule(4, "frontend", "test", "")
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
	err = clientTest.DeleteTCPRequestRule(24, "frontend", "test", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ := clientTest.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	_, _, err = clientTest.GetTCPRequestRule(24, "frontend", "test", "")
	if err == nil {
		t.Error("DeleteTCPRequestRule failed, TCP Request Rule 22 still exists")
	}

	err = clientTest.DeleteTCPRequestRule(24, "backend", "test_2", "", version)
	if err == nil {
		t.Error("Should throw error, non existant TCP Request Rule")
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
	}

	for _, testCase := range testCases {
		t.Run(testCase.expectedResult, func(t *testing.T) {
			tcpType, err := SerializeTCPRequestRule(testCase.input)
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
