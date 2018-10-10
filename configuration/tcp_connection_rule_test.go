package configuration

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/haproxytech/models"
)

func TestGetTCPConnectionRules(t *testing.T) {
	tRules, err := client.GetTCPConnectionRules("test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if len(tRules.Data) != 2 {
		t.Errorf("%v TCP connection rules returned, expected 2", len(tRules.Data))
	}

	if tRules.Version != version {
		t.Errorf("Version %v returned, expected %v", tRules.Version, version)
	}

	for _, tcR := range tRules.Data {
		if tcR.ID == 1 {
			if tcR.Type != "accept" {
				t.Errorf("%v: Type not accept: %v", tcR.ID, tcR.Type)
			}
			if tcR.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", tcR.ID, tcR.Cond)
			}
			if tcR.CondTest != "TRUE" {
				t.Errorf("%v: CondTest not TRUE: %v", tcR.ID, tcR.CondTest)
			}
		} else if tcR.ID == 2 {
			if tcR.Type != "reject" {
				t.Errorf("%v: Type not accept: %v", tcR.ID, tcR.Type)
			}
			if tcR.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", tcR.ID, tcR.Cond)
			}
			if tcR.CondTest != "FALSE" {
				t.Errorf("%v: CondTest not TRUE: %v", tcR.ID, tcR.CondTest)
			}

		} else {
			t.Errorf("Expected ID 1 or 2, %v found", tcR.ID)
		}
	}

	tJSON, err := tRules.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	tRules, err = client.GetTCPConnectionRules("test2", "")
	if err != nil {
		t.Error(err.Error())
	}
	if len(tRules.Data) > 0 {
		t.Errorf("%v TCP connection rules returned, expected 0", len(tRules.Data))
	}

	if !t.Failed() {
		fmt.Println("GetTCPConnectionRules succesful\nResponse: \n" + string(tJSON) + "\n")
	}
}

func TestGetTCPConnectionRule(t *testing.T) {
	tRule, err := client.GetTCPConnectionRule(1, "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	tcR := tRule.Data

	if tcR.ID != 1 {
		t.Errorf("Expected ID 1, %v found", tcR.ID)
	}
	if tcR.Type != "accept" {
		t.Errorf("%v: Type not accept: %v", tcR.ID, tcR.Type)
	}
	if tcR.Cond != "if" {
		t.Errorf("%v: Cond not if: %v", tcR.ID, tcR.Cond)
	}
	if tcR.CondTest != "TRUE" {
		t.Errorf("%v: CondTest not TRUE: %v", tcR.ID, tcR.CondTest)
	}

	tJSON, err := tcR.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	_, err = client.GetTCPConnectionRule(3, "test2", "")
	if err == nil {
		t.Error("Should throw error, non existant TCP connection rule")
	}

	if !t.Failed() {
		fmt.Println("GetTCPConnectionRule succesful\nResponse: \n" + string(tJSON) + "\n")
	}
}

func TestCreateEditDeleteTCPConnectionRule(t *testing.T) {
	// TestCreateTCPConnectionRule
	tcR := &models.TCPRule{
		ID:       1,
		Type:     "accept",
		Cond:     "unless",
		CondTest: "TRUE",
	}

	err := client.CreateTCPConnectionRule("test", tcR, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	tRule, err := client.GetTCPConnectionRule(1, "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(tRule.Data, tcR) {
		fmt.Printf("Created TCP connection rule: %v\n", tRule.Data)
		fmt.Printf("Given TCP connection rule: %v\n", tcR)
		t.Error("Created TCP connection rule not equal to given TCP connection rule")
	}

	if tRule.Version != version {
		t.Errorf("Version %v returned, expected %v", tRule.Version, version)
	}

	if !t.Failed() {
		fmt.Println("CreateTCPConnectionRule successful")
	}

	// TestEditTCPConnectionRule
	tcR = &models.TCPRule{
		ID:       1,
		Type:     "accept",
		Cond:     "if",
		CondTest: "TRUE",
	}

	err = client.EditTCPConnectionRule(1, "test", tcR, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	tRule, err = client.GetTCPConnectionRule(1, "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(tRule.Data, tcR) {
		fmt.Printf("Edited TCP connection rule: %v\n", tRule.Data)
		fmt.Printf("Given TCP connection rule: %v\n", tcR)
		t.Error("Edited TCP connection rule not equal to given TCP connection rule")
	}

	if tRule.Version != version {
		t.Errorf("Version %v returned, expected %v", tRule.Version, version)
	}

	if !t.Failed() {
		fmt.Println("EditTCPConnectionRule successful")
	}

	// TestDeleteTCPConnectionRule
	err = client.DeleteTCPConnectionRule(3, "test", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ := client.GetVersion(); v != version {
		t.Error("Version not incremented")
	}

	_, err = client.GetTCPConnectionRule(3, "test", "")
	if err == nil {
		t.Error("DeleteTCPConnectionRule failed, TCP connection rule 2 still exists")
	}

	err = client.DeleteTCPConnectionRule(3, "test2", "", version)
	if err == nil {
		t.Error("Should throw error, non existant TCP connection rule")
		version++
	}

	if !t.Failed() {
		fmt.Println("DeleteTCPConnectionRule successful")
	}
}
