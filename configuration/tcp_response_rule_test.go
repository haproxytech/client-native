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

func TestGetTCPResponseRules(t *testing.T) { //nolint:gocognit
	v, tRules, err := clientTest.GetTCPResponseRules("test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if len(tRules) != 18 {
		t.Errorf("%v tcp response rules returned, expected 17", len(tRules))
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	for _, r := range tRules {
		switch *r.Index {
		case 0:
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
		case 1:
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
		case 2:
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
		case 3:
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
		case 4:
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
		case 5:
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
		case 6:
			if r.Type != "content" {
				t.Errorf("%v: Type not content: %v", *r.Index, r.Type)
			}
			if r.Action != "close" {
				t.Errorf("%v: Action not close: %v", *r.Index, r.Action)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "FALSE" {
				t.Errorf("%v: CondTest not FALSE: %v", *r.Index, r.CondTest)
			}
		case 7:
			if r.Type != "content" {
				t.Errorf("%v: Type not content: %v", *r.Index, r.Type)
			}
			if r.Action != "sc-add-gpc" {
				t.Errorf("%v: Type not sc-add-gpc: %v", *r.Index, r.Type)
			}
			if r.ScIdx != 0 {
				t.Errorf("%v: ScIdx not 0: %v", *r.Index, r.ScIdx)
			}
			if r.ScID != 1 {
				t.Errorf("%v: ScID not 1: %v", *r.Index, r.ScID)
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
			if r.Action != "sc-inc-gpc0" {
				t.Errorf("%v: Action not sc-inc-gpc0: %v", *r.Index, r.Action)
			}
			if r.ScID != 1 {
				t.Errorf("%v: ScID not 1 %v", *r.Index, r.ScID)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "FALSE" {
				t.Errorf("%v: CondTest not FALSE: %v", *r.Index, r.CondTest)
			}
		case 9:
			if r.Type != "content" {
				t.Errorf("%v: Type not content: %v", *r.Index, r.Type)
			}
			if r.Action != "sc-inc-gpc1" {
				t.Errorf("%v: Action not sc-inc-gpc1: %v", *r.Index, r.Action)
			}
			if r.ScID != 2 {
				t.Errorf("%v: ScID not 2 %v", *r.Index, r.ScID)
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
			if r.Action != "sc-set-gpt0" {
				t.Errorf("%v: Action not sc-set-gpt0: %v", *r.Index, r.Action)
			}
			if r.ScID != 3 {
				t.Errorf("%v: ScID not 3 %v", *r.Index, r.ScID)
			}
			if r.Expr != "hdr(Host),lower" {
				t.Errorf("%v: Expr not hdr(Host),lower: %v", *r.Index, r.Expr)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "FALSE" {
				t.Errorf("%v: CondTest not FALSE: %v", *r.Index, r.CondTest)
			}
		case 11:
			if r.Type != "content" {
				t.Errorf("%v: Type not content: %v", *r.Index, r.Type)
			}
			if r.Action != "send-spoe-group" {
				t.Errorf("%v: Action not send-spoe-group: %v", *r.Index, r.Action)
			}
			if r.SpoeEngine != "engine" {
				t.Errorf("%v: SpoeEngine not engine %v", *r.Index, r.SpoeEngine)
			}
			if r.SpoeGroup != "group" {
				t.Errorf("%v: SpoeGroup not group %v", *r.Index, r.SpoeGroup)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "FALSE" {
				t.Errorf("%v: CondTest not FALSE: %v", *r.Index, r.CondTest)
			}
		case 12:
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
		case 13:
			if r.Type != "content" {
				t.Errorf("%v: Type not content: %v", *r.Index, r.Type)
			}
			if r.Action != "set-mark" {
				t.Errorf("%v: Action not set-mark: %v", *r.Index, r.Action)
			}
			if r.MarkValue != "0x1Ab" {
				t.Errorf("%v: MarkValue not 0x1Ab %v", *r.Index, r.MarkValue)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "FALSE" {
				t.Errorf("%v: CondTest not FALSE: %v", *r.Index, r.CondTest)
			}
		case 14:
			if r.Type != "content" {
				t.Errorf("%v: Type not content: %v", *r.Index, r.Type)
			}
			if r.Action != "set-nice" {
				t.Errorf("%v: Action not set-nice: %v", *r.Index, r.Action)
			}
			if r.NiceValue != 1 {
				t.Errorf("%v: NiceValue not 1 %v", *r.Index, r.NiceValue)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "FALSE" {
				t.Errorf("%v: CondTest not FALSE: %v", *r.Index, r.CondTest)
			}
		case 15:
			if r.Type != "content" {
				t.Errorf("%v: Type not content: %v", *r.Index, r.Type)
			}
			if r.Action != "set-tos" {
				t.Errorf("%v: Action not set-tos: %v", *r.Index, r.Action)
			}
			if r.TosValue != "2" {
				t.Errorf("%v: TosValue not 2 %v", *r.Index, r.TosValue)
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
			if r.Action != "silent-drop" {
				t.Errorf("%v: Action not silent-drop: %v", *r.Index, r.Action)
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
			if r.Action != "unset-var" {
				t.Errorf("%v: Action not unset-var: %v", *r.Index, r.Action)
			}
			if r.VarName != "my_var" {
				t.Errorf("%v: VarName not my_var: %v", *r.Index, r.VarName)
			}
			if r.VarScope != "req" {
				t.Errorf("%v: VarScope not req: %v", *r.Index, r.VarScope)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "FALSE" {
				t.Errorf("%v: CondTest not FALSE: %v", *r.Index, r.CondTest)
			}
		default:
			t.Errorf("Expect only tcp-response 0 to 17 %v found", *r.Index)
		}
	}

	_, tRules, err = clientTest.GetTCPResponseRules("test_2", "")
	if err != nil {
		t.Error(err.Error())
	}
	if len(tRules) > 0 {
		t.Errorf("%v TCP Response Rules returned, expected 0", len(tRules))
	}
}

func TestGetTCPResponseRule(t *testing.T) {
	v, r, err := clientTest.GetTCPResponseRule(0, "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

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

	_, err = r.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	_, _, err = clientTest.GetTCPResponseRule(3, "test_2", "")
	if err == nil {
		t.Error("Should throw error, non existant TCP Response Rule")
	}
}

func TestCreateEditDeleteTCPResponseRule(t *testing.T) {
	id := int64(2)
	tOut := int64(1000)
	// TestCreateTCPResponseRule
	r := &models.TCPResponseRule{
		Index:   &id,
		Type:    "inspect-delay",
		Timeout: &tOut,
	}

	err := clientTest.CreateTCPResponseRule("test", r, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, ondiskR, err := clientTest.GetTCPResponseRule(2, "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(ondiskR, r) {
		fmt.Printf("Created TCP response rule: %v\n", ondiskR)
		fmt.Printf("Given TCP response rule: %v\n", r)
		t.Error("Created TCP response rule not equal to given TCP response rule")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	// TestEditTCPResponseRule
	r = &models.TCPResponseRule{
		Index:    &id,
		Type:     "content",
		Action:   "accept",
		Cond:     "if",
		CondTest: "FALSE",
	}

	err = clientTest.EditTCPResponseRule(2, "test", r, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, ondiskR, err = clientTest.GetTCPResponseRule(2, "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(ondiskR, r) {
		fmt.Printf("Edited TCP response rule: %v\n", ondiskR)
		fmt.Printf("Given TCP response rule: %v\n", r)
		t.Error("Edited TCP response rule not equal to given TCP response rule")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	// TestDeleteTCPResponse
	err = clientTest.DeleteTCPResponseRule(18, "test", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ := clientTest.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	_, _, err = clientTest.GetTCPResponseRule(18, "test", "")
	if err == nil {
		t.Error("DeleteTCPResponseRule failed, TCP Response Rule 17 still exists")
	}

	err = clientTest.DeleteTCPResponseRule(18, "test_2", "", version)
	if err == nil {
		t.Error("Should throw error, non existant TCP Response Rule")
		version++
	}
}
