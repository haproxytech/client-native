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

	"github.com/haproxytech/models"
)

func TestGetTCPRequestRules(t *testing.T) {
	v, tRules, err := client.GetTCPRequestRules("frontend", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if len(tRules) != 4 {
		t.Errorf("%v tcp request rules returned, expected 4", len(tRules))
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	for _, r := range tRules {
		if *r.ID == 0 {
			if r.Type != "connection" {
				t.Errorf("%v: Type not connection: %v", *r.ID, r.Type)
			}
			if r.Action != "accept" {
				t.Errorf("%v: Action not accept: %v", *r.ID, r.Action)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.ID, r.Cond)
			}
			if r.CondTest != "TRUE" {
				t.Errorf("%v: CondTest not src TRUE: %v", *r.ID, r.CondTest)
			}
		} else if *r.ID == 1 {
			if r.Type != "connection" {
				t.Errorf("%v: Type not connection: %v", *r.ID, r.Type)
			}
			if r.Action != "reject" {
				t.Errorf("%v: Action not reject: %v", *r.ID, r.Action)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.ID, r.Cond)
			}
			if r.CondTest != "FALSE" {
				t.Errorf("%v: CondTest not src FALSE: %v", *r.ID, r.CondTest)
			}
		} else if *r.ID == 2 {
			if r.Type != "content" {
				t.Errorf("%v: Type not content: %v", *r.ID, r.Type)
			}
			if r.Action != "accept" {
				t.Errorf("%v: Action not accept: %v", *r.ID, r.Action)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.ID, r.Cond)
			}
			if r.CondTest != "TRUE" {
				t.Errorf("%v: CondTest not src TRUE: %v", *r.ID, r.CondTest)
			}
		} else if *r.ID == 3 {
			if r.Type != "content" {
				t.Errorf("%v: Type not content: %v", *r.ID, r.Type)
			}
			if r.Action != "reject" {
				t.Errorf("%v: Action not reject: %v", *r.ID, r.Action)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.ID, r.Cond)
			}
			if r.CondTest != "FALSE" {
				t.Errorf("%v: CondTest not src FALSE: %v", *r.ID, r.CondTest)
			}
		} else {
			t.Errorf("Expext only tcp-request 0, 1, 2, or 3, %v found", *r.ID)
		}
	}

	_, tRules, err = client.GetTCPRequestRules("backend", "test_2", "")
	if err != nil {
		t.Error(err.Error())
	}
	if len(tRules) > 0 {
		t.Errorf("%v TCP Request Ruless returned, expected 0", len(tRules))
	}
}

func TestGetTCPRequestRule(t *testing.T) {
	v, r, err := client.GetTCPRequestRule(0, "frontend", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	if r.Type != "connection" {
		t.Errorf("%v: Type not connection: %v", *r.ID, r.Type)
	}
	if r.Action != "accept" {
		t.Errorf("%v: Action not accept: %v", *r.ID, r.Action)
	}
	if r.Cond != "if" {
		t.Errorf("%v: Cond not if: %v", *r.ID, r.Cond)
	}
	if r.CondTest != "TRUE" {
		t.Errorf("%v: CondTest not src TRUE: %v", *r.ID, r.CondTest)
	}

	_, err = r.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	_, _, err = client.GetTCPRequestRule(3, "backend", "test_2", "")
	if err == nil {
		t.Error("Should throw error, non existant TCP Request Rule")
	}
}

func TestCreateEditDeleteTCPRequestRule(t *testing.T) {
	id := int64(4)
	tOut := int64(1000)
	// TestCreateTCPRequestRule
	r := &models.TCPRequestRule{
		ID:      &id,
		Type:    "inspect-delay",
		Timeout: &tOut,
	}

	err := client.CreateTCPRequestRule("frontend", "test", r, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, ondiskR, err := client.GetTCPRequestRule(4, "frontend", "test", "")
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
		ID:       &id,
		Type:     "connection",
		Action:   "accept",
		Cond:     "if",
		CondTest: "FALSE",
	}

	err = client.EditTCPRequestRule(4, "frontend", "test", r, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, ondiskR, err = client.GetTCPRequestRule(4, "frontend", "test", "")
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
	err = client.DeleteTCPRequestRule(4, "frontend", "test", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ := client.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	_, _, err = client.GetTCPRequestRule(4, "frontend", "test", "")
	if err == nil {
		t.Error("DeleteTCPRequestRule failed, TCP Request Rule 3 still exists")
	}

	err = client.DeleteTCPRequestRule(2, "backend", "test_2", "", version)
	if err == nil {
		t.Error("Should throw error, non existant TCP Request Rule")
		version++
	}
}
