package configuration

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/haproxytech/models"
)

func TestGetTCPContentRules(t *testing.T) {
	tRules, err := client.GetTCPContentRules("frontend", "test", "request", "")
	if err != nil {
		t.Error(err.Error())
	}

	if len(tRules.Data) != 2 {
		t.Errorf("%v tcp content rules returned, expected 2", len(tRules.Data))
	}

	if tRules.Version != version {
		t.Errorf("Version %v returned, expected %v", tRules.Version, version)
	}

	for _, r := range tRules.Data {
		if r.ID == 1 {
			if r.Type != "accept" {
				t.Errorf("%v: Type not accept: %v", r.ID, r.Type)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", r.ID, r.Cond)
			}
			if r.CondTest != "TRUE" {
				t.Errorf("%v: CondTest not TRUE: %v", r.ID, r.CondTest)
			}
		} else if r.ID == 2 {
			if r.Type != "reject" {
				t.Errorf("%v: Type not reject: %v", r.ID, r.Type)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", r.ID, r.Cond)
			}
			if r.CondTest != "FALSE" {
				t.Errorf("%v: CondTest not FALSE: %v", r.ID, r.CondTest)
			}
		} else {
			t.Errorf("Expext only filter 1, 2 or 3, %v found", r.ID)
		}
	}

	tJSON, err := tRules.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	tRules, err = client.GetTCPContentRules("backend", "test2", "response", "")
	if err != nil {
		t.Error(err.Error())
	}
	if len(tRules.Data) > 0 {
		t.Errorf("%v filters returned, expected 0", len(tRules.Data))
	}

	if !t.Failed() {
		fmt.Println("GetTCPContentRules succesful\nResponse: \n" + string(tJSON) + "\n")
	}
}

func TestGetTCPContentRule(t *testing.T) {
	tRule, err := client.GetTCPContentRule(1, "backend", "test", "response", "")
	if err != nil {
		t.Error(err.Error())
	}

	r := tRule.Data

	if r.ID != 1 {
		t.Errorf("TCP Content Rule ID not 1, %v found", r.ID)
	}
	if r.Type != "accept" {
		t.Errorf("%v: Type not accept: %v", r.ID, r.Type)
	}
	if r.Cond != "if" {
		t.Errorf("%v: Cond not if: %v", r.ID, r.Cond)
	}
	if r.CondTest != "TRUE" {
		t.Errorf("%v: CondTest not TRUE: %v", r.ID, r.CondTest)
	}

	rJSON, err := r.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	_, err = client.GetTCPContentRule(3, "backend", "test2", "response", "")
	if err == nil {
		t.Error("Should throw error, non existant TCP Content Rule")
	}

	if !t.Failed() {
		fmt.Println("GetTCPContentRule succesful\nResponse: \n" + string(rJSON) + "\n")
	}
}

func TestCreateEditDeleteTCPContentRule(t *testing.T) {
	// TestCreateTCPContentRule
	r := &models.TCPRule{
		ID:       1,
		Type:     "reject",
		Cond:     "unless",
		CondTest: "FALSE",
	}

	err := client.CreateTCPContentRule("frontend", "test", "response", r, "", version)
	if err == nil {
		t.Error(fmt.Errorf("Should throw error, no response type allowed for frontend"))
	}

	err = client.CreateTCPContentRule("frontend", "test", "request", r, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	ondiskR, err := client.GetTCPContentRule(1, "frontend", "test", "request", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(ondiskR.Data, r) {
		fmt.Printf("Created TCP content rule: %v\n", ondiskR.Data)
		fmt.Printf("Given TCP content rule: %v\n", r)
		t.Error("Created TCP content rule not equal to given TCP content rule")
	}

	if ondiskR.Version != version {
		t.Errorf("Version %v returned, expected %v", ondiskR.Version, version)
	}

	if !t.Failed() {
		fmt.Println("CreateTCPContentRule successful")
	}

	// TestEditTCPContentRule
	r = &models.TCPRule{
		ID:       1,
		Type:     "accept",
		Cond:     "if",
		CondTest: "FALSE",
	}

	err = client.EditTCPContentRule(1, "frontend", "test", "request", r, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	ondiskR, err = client.GetTCPContentRule(1, "frontend", "test", "request", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(ondiskR.Data, r) {
		fmt.Printf("Edited TCP content rule: %v\n", ondiskR.Data)
		fmt.Printf("Given TCP content rule: %v\n", r)
		t.Error("Edited TCP content rule not equal to given TCP content rule")
	}

	if ondiskR.Version != version {
		t.Errorf("Version %v returned, expected %v", ondiskR.Version, version)
	}

	if !t.Failed() {
		fmt.Println("EditTCPContentRule successful")
	}

	// TestDeleteFilter
	err = client.DeleteTCPContentRule(3, "frontend", "test", "request", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ := client.GetVersion(); v != version {
		t.Error("Version not incremented")
	}

	_, err = client.GetTCPContentRule(3, "frontend", "test", "request", "")
	if err == nil {
		t.Error("DeleteTCPContentRule failed, TCP Content Rule 4 still exists")
	}

	err = client.DeleteTCPContentRule(2, "backend", "test2", "response", "", version)
	if err == nil {
		t.Error("Should throw error, non existant TCP Content Rule")
		version++
	}

	if !t.Failed() {
		fmt.Println("DeleteTCPContentRule successful")
	}
}
