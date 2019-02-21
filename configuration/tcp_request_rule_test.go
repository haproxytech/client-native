package configuration

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/haproxytech/models"
)

func TestGetTCPRequestRules(t *testing.T) {
	tRules, err := client.GetTCPRequestRules("frontend", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if len(tRules.Data) != 4 {
		t.Errorf("%v tcp request rules returned, expected 4", len(tRules.Data))
	}

	if tRules.Version != version {
		t.Errorf("Version %v returned, expected %v", tRules.Version, version)
	}

	for _, r := range tRules.Data {
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

	rJSON, err := tRules.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	tRules, err = client.GetTCPRequestRules("backend", "test_2", "")
	if err != nil {
		t.Error(err.Error())
	}
	if len(tRules.Data) > 0 {
		t.Errorf("%v TCP Request Ruless returned, expected 0", len(tRules.Data))
	}

	if !t.Failed() {
		fmt.Println("GetTCPRequestRules succesful\nResponse: \n" + string(rJSON) + "\n")
	}
}

func TestGetTCPRequestRule(t *testing.T) {
	tRule, err := client.GetTCPRequestRule(0, "frontend", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	r := tRule.Data

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

	rJSON, err := r.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	_, err = client.GetTCPRequestRule(3, "backend", "test_2", "")
	if err == nil {
		t.Error("Should throw error, non existant TCP Request Rule")
	}

	if !t.Failed() {
		fmt.Println("GetTCPRequestRule succesful\nResponse: \n" + string(rJSON) + "\n")
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

	ondiskR, err := client.GetTCPRequestRule(4, "frontend", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(ondiskR.Data, r) {
		fmt.Printf("Created TCP request rule: %v\n", ondiskR.Data)
		fmt.Printf("Given TCP request rule: %v\n", r)
		t.Error("Created TCP request rule not equal to given TCP request rule")
	}

	if ondiskR.Version != version {
		t.Errorf("Version %v returned, expected %v", ondiskR.Version, version)
	}

	if !t.Failed() {
		fmt.Println("CreateTCPRequestRule successful")
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

	ondiskR, err = client.GetTCPRequestRule(4, "frontend", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(ondiskR.Data, r) {
		fmt.Printf("Edited TCP request rule: %v\n", ondiskR.Data)
		fmt.Printf("Given TCP request rule: %v\n", r)
		t.Error("Edited TCP request rule not equal to given TCP request rule")
	}

	if ondiskR.Version != version {
		t.Errorf("Version %v returned, expected %v", ondiskR.Version, version)
	}

	if !t.Failed() {
		fmt.Println("EditTCPRequestRule successful")
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

	_, err = client.GetTCPRequestRule(4, "frontend", "test", "")
	if err == nil {
		t.Error("DeleteTCPRequestRule failed, TCP Request Rule 3 still exists")
	}

	err = client.DeleteTCPRequestRule(2, "backend", "test_2", "", version)
	if err == nil {
		t.Error("Should throw error, non existant TCP Request Rule")
		version++
	}

	if !t.Failed() {
		fmt.Println("DeleteTCPRequestRule successful")
	}
}
