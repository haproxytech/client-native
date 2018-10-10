package configuration

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/haproxytech/models"
)

func TestGetStickResponseRules(t *testing.T) {
	sRules, err := client.GetStickResponseRules("test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if len(sRules.Data) != 3 {
		t.Errorf("%v stick response rules returned, expected 3", len(sRules.Data))
	}

	if sRules.Version != version {
		t.Errorf("Version %v returned, expected %v", sRules.Version, version)
	}

	for _, sr := range sRules.Data {
		if sr.Type != "storeonly" {
			t.Errorf("%v: Type not storeonly: %v", sr.ID, sr.Type)
		}
		if sr.ID == 1 {
			if sr.Pattern != "src" {
				t.Errorf("%v: Pattern not src: %v", sr.ID, sr.Pattern)
			}
		} else if sr.ID == 2 {
			if sr.Pattern != "src_port" {
				t.Errorf("%v: Pattern not src: %v", sr.ID, sr.Pattern)
			}
			if sr.Table != "test_port" {
				t.Errorf("%v: Table not test: %v", sr.ID, sr.Table)
			}
		} else if sr.ID == 3 {
			if sr.Pattern != "src" {
				t.Errorf("%v: Pattern not src: %v", sr.ID, sr.Pattern)
			}
			if sr.Table != "test" {
				t.Errorf("%v: Table not test: %v", sr.ID, sr.Table)
			}
			if sr.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", sr.ID, sr.Cond)
			}
			if sr.CondTest != "TRUE" {
				t.Errorf("%v: CondTest not TRUE: %v", sr.ID, sr.CondTest)
			}
		} else {
			t.Errorf("Expext only stick response rule 1, 2, 3 %v found", sr.ID)
		}
	}

	srJSON, err := sRules.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	sRules, err = client.GetStickResponseRules("test2", "")
	if err != nil {
		t.Error(err.Error())
	}
	if len(sRules.Data) > 0 {
		t.Errorf("%v stick response rules returned, expected 0", len(sRules.Data))
	}

	if !t.Failed() {
		fmt.Println("GetStickResponseRules succesful\nResponse: \n" + string(srJSON) + "\n")
	}
}

func TestGetStickResponseRule(t *testing.T) {
	sRule, err := client.GetStickResponseRule(3, "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	sr := sRule.Data

	if sr.Type != "storeonly" {
		t.Errorf("%v: Type not storeonly: %v", sr.ID, sr.Type)
	}
	if sr.Pattern != "src" {
		t.Errorf("%v: Pattern not src: %v", sr.ID, sr.Pattern)
	}
	if sr.Table != "test" {
		t.Errorf("%v: Table not test: %v", sr.ID, sr.Table)
	}
	if sr.Cond != "if" {
		t.Errorf("%v: Cond not if: %v", sr.ID, sr.Cond)
	}
	if sr.CondTest != "TRUE" {
		t.Errorf("%v: CondTest not TRUE: %v", sr.ID, sr.CondTest)
	}

	srJSON, err := sr.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	_, err = client.GetStickResponseRule(5, "test2", "")
	if err == nil {
		t.Error("Should throw error, non existant stick response rule")
	}

	if !t.Failed() {
		fmt.Println("GetStickResponseRule succesful\nResponse: \n" + string(srJSON) + "\n")
	}
}

func TestCreateEditDeleteStickResponseRule(t *testing.T) {
	// TestCreateStickResponseRule
	sr := &models.StickResponseRule{
		ID:       1,
		Type:     "storeonly",
		Pattern:  "src_port",
		Cond:     "if",
		CondTest: "TRUE",
	}

	err := client.CreateStickResponseRule("test", sr, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	sRule, err := client.GetStickResponseRule(1, "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(sRule.Data, sr) {
		fmt.Printf("Created stick response rule: %v\n", sRule.Data)
		fmt.Printf("Given stick response rule: %v\n", sr)
		t.Error("Created stick response rule not equal to given stick response rule")
	}

	if sRule.Version != version {
		t.Errorf("Version %v returned, expected %v", sRule.Version, version)
	}

	if !t.Failed() {
		fmt.Println("CreateStickResponseRule successful")
	}

	// TestEditStickResponseRule
	sr = &models.StickResponseRule{
		ID:       1,
		Type:     "storeonly",
		Pattern:  "src",
		Table:    "test2",
		Cond:     "if",
		CondTest: "FALSE",
	}

	err = client.EditStickResponseRule(1, "test", sr, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	sRule, err = client.GetStickResponseRule(1, "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(sRule.Data, sr) {
		fmt.Printf("Edited stick response rule: %v\n", sRule.Data)
		fmt.Printf("Given stick response rule: %v\n", sr)
		t.Error("Edited stick response rule not equal to given stick response rule")
	}

	if sRule.Version != version {
		t.Errorf("Version %v returned, expected %v", sRule.Version, version)
	}

	if !t.Failed() {
		fmt.Println("EditStickResponseRule successful")
	}

	// TestDeleteStickResponseRule
	err = client.DeleteStickResponseRule(4, "test", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ := client.GetVersion(); v != version {
		t.Error("Version not incremented")
	}

	_, err = client.GetStickResponseRule(4, "test", "")
	if err == nil {
		t.Error("DeleteStickResponseRule failed, stick response rule 3 still exists")
	}

	err = client.DeleteStickResponseRule(4, "test2", "", version)
	if err == nil {
		t.Error("Should throw error, non existant stick response rule")
		version++
	}

	if !t.Failed() {
		fmt.Println("DeleteStickResponseRule successful")
	}
}
