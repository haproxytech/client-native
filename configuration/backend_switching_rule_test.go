package configuration

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/haproxytech/models"
)

func TestGetBackendSwitchingRules(t *testing.T) {
	bckRules, err := client.GetBackendSwitchingRules("test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if len(bckRules.Data) != 1 {
		t.Errorf("%v backend switching rules returned, expected 1", len(bckRules.Data))
	}

	if bckRules.Version != version {
		t.Errorf("Version %v returned, expected %v", bckRules.Version, version)
	}

	for _, br := range bckRules.Data {
		if br.ID != 1 {
			t.Errorf("ID only backend switching rule 1, %v found", br.ID)
		}
		if br.TargetFarm != "test_2" {
			t.Errorf("%v: TargetFarm not test_2: %v", br.ID, br.TargetFarm)
		}
		if br.Cond != "if" {
			t.Errorf("%v: Cond not if: %v", br.ID, br.Cond)
		}
		if br.CondTest != "TRUE" {
			t.Errorf("%v: CondTest not TRUE: %v", br.ID, br.CondTest)
		}
	}

	brJSON, err := bckRules.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	bckRules, err = client.GetBackendSwitchingRules("test2", "")
	if err != nil {
		t.Error(err.Error())
	}
	if len(bckRules.Data) > 0 {
		t.Errorf("%v backend switching rules returned, expected 0", len(bckRules.Data))
	}

	if !t.Failed() {
		fmt.Println("GetBackendSwitchingRules succesful\nResponse: \n" + string(brJSON) + "\n")
	}
}

func TestGetBackendSwitchingRule(t *testing.T) {
	bckRule, err := client.GetBackendSwitchingRule(1, "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	br := bckRule.Data

	if br.ID != 1 {
		t.Errorf("ID only backend switching rule 1, %v found", br.ID)
	}
	if br.TargetFarm != "test_2" {
		t.Errorf("%v: TargetFarm not test_2: %v", br.ID, br.TargetFarm)
	}
	if br.Cond != "if" {
		t.Errorf("%v: Cond not if: %v", br.ID, br.Cond)
	}
	if br.CondTest != "TRUE" {
		t.Errorf("%v: CondTest not TRUE: %v", br.ID, br.CondTest)
	}

	brJSON, err := br.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	_, err = client.GetBackendSwitchingRule(3, "test2", "")
	if err == nil {
		t.Error("Should throw error, non existant backend switching rule")
	}

	if !t.Failed() {
		fmt.Println("GetBackendSwitchingRule succesful\nResponse: \n" + string(brJSON) + "\n")
	}
}

func TestCreateEditDeleteBackendSwitchingRule(t *testing.T) {
	// TestCreateBackendSwitchingRule
	br := &models.BackendSwitchingRule{
		ID:         1,
		TargetFarm: "test",
		Cond:       "unless",
		CondTest:   "TRUE",
	}

	err := client.CreateBackendSwitchingRule("test", br, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	bckRule, err := client.GetBackendSwitchingRule(1, "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(bckRule.Data, br) {
		fmt.Printf("Created backend switching rule: %v\n", bckRule.Data)
		fmt.Printf("Given backend switching rule: %v\n", br)
		t.Error("Created backend switching rule not equal to given backend switching rule")
	}

	if bckRule.Version != version {
		t.Errorf("Version %v returned, expected %v", bckRule.Version, version)
	}

	if !t.Failed() {
		fmt.Println("CreateBackendSwitchingRule successful")
	}

	// TestBackendSwitchingRule
	br = &models.BackendSwitchingRule{
		ID:         1,
		TargetFarm: "test",
		Cond:       "if",
		CondTest:   "TRUE",
	}

	err = client.EditBackendSwitchingRule(1, "test", br, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	bckRule, err = client.GetBackendSwitchingRule(1, "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(bckRule.Data, br) {
		fmt.Printf("Edited backend switching rule: %v\n", bckRule.Data)
		fmt.Printf("Given backend switching rule: %v\n", br)
		t.Error("Edited backend switching rule not equal to given backend switching rule")
	}

	if bckRule.Version != version {
		t.Errorf("Version %v returned, expected %v", bckRule.Version, version)
	}

	if !t.Failed() {
		fmt.Println("EditBackendSwitchingRule successful")
	}

	// TestBackendSwitchingRule
	err = client.DeleteBackendSwitchingRule(2, "test", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ := client.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	_, err = client.GetBackendSwitchingRule(2, "test", "")
	if err == nil {
		t.Error("DeleteBackendSwitchingRule failed, backend switching rule 2 still exists")
	}

	err = client.DeleteBackendSwitchingRule(2, "test2", "", version)
	if err == nil {
		t.Error("Should throw error, non existant backend switching rule")
		version++
	}

	if !t.Failed() {
		fmt.Println("DeleteBackendSwitchingRule successful")
	}
}
