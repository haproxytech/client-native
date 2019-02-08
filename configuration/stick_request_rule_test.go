package configuration

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/haproxytech/models"
)

func TestGetStickRequestRules(t *testing.T) {
	sRules, err := client.GetStickRequestRules("test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if len(sRules.Data) != 3 {
		t.Errorf("%v stick request rules returned, expected 3", len(sRules.Data))
	}

	if sRules.Version != version {
		t.Errorf("Version %v returned, expected %v", sRules.Version, version)
	}

	for _, sr := range sRules.Data {
		if sr.ID == 1 {
			if sr.Type != "storeonly" {
				t.Errorf("%v: Type not storeonly: %v", sr.ID, sr.Type)
			}
		} else if sr.ID == 2 {
			if sr.Type != "matchonly" {
				t.Errorf("%v: Type not matchonly: %v", sr.ID, sr.Type)
			}
		} else if sr.ID == 3 {
			if sr.Type != "matchandstore" {
				t.Errorf("%v: Type not matchandstore: %v", sr.ID, sr.Type)
			}
		} else {
			t.Errorf("Expext only stick request rule 1 or 2, %v found", sr.ID)
		}
		if sr.Pattern != "src" {
			t.Errorf("%v: Pattern not src: %v", sr.ID, sr.Pattern)
		}
		if sr.Table != "test" {
			t.Errorf("%v: Table not test: %v", sr.ID, sr.Table)
		}
	}

	srJSON, err := sRules.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	sRules, err = client.GetStickRequestRules("test2", "")
	if err != nil {
		t.Error(err.Error())
	}
	if len(sRules.Data) > 0 {
		t.Errorf("%v stick request rules returned, expected 0", len(sRules.Data))
	}

	if !t.Failed() {
		fmt.Println("GetStickRequestRules succesful\nResponse: \n" + string(srJSON) + "\n")
	}
}

func TestGetStickRequestRule(t *testing.T) {
	sRule, err := client.GetStickRequestRule(1, "test", "")
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

	srJSON, err := sr.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	_, err = client.GetStickRequestRule(5, "test2", "")
	if err == nil {
		t.Error("Should throw error, non existant stick request rule")
	}

	if !t.Failed() {
		fmt.Println("GetStickRequestRule succesful\nResponse: \n" + string(srJSON) + "\n")
	}
}

func TestCreateEditDeleteStickRequestRule(t *testing.T) {
	// TestCreateStickRequestRule
	sr := &models.StickRequestRule{
		ID:       1,
		Type:     "matchonly",
		Pattern:  "src",
		Cond:     "if",
		CondTest: "TRUE",
	}

	err := client.CreateStickRequestRule("test", sr, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	sRule, err := client.GetStickRequestRule(1, "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(sRule.Data, sr) {
		fmt.Printf("Created stick request rule: %v\n", sRule.Data)
		fmt.Printf("Given stick request rule: %v\n", sr)
		t.Error("Created stick request rule not equal to given stick request rule")
	}

	if sRule.Version != version {
		t.Errorf("Version %v returned, expected %v", sRule.Version, version)
	}

	if !t.Failed() {
		fmt.Println("CreateStickRequestRule successful")
	}

	// TestEditStickRequestRule
	sr = &models.StickRequestRule{
		ID:       1,
		Type:     "storeonly",
		Pattern:  "src",
		Table:    "test2",
		Cond:     "if",
		CondTest: "FALSE",
	}

	err = client.EditStickRequestRule(1, "test", sr, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	sRule, err = client.GetStickRequestRule(1, "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(sRule.Data, sr) {
		fmt.Printf("Edited stick request rule: %v\n", sRule.Data)
		fmt.Printf("Given stick request rule: %v\n", sr)
		t.Error("Edited stick request rule not equal to given stick request rule")
	}

	if sRule.Version != version {
		t.Errorf("Version %v returned, expected %v", sRule.Version, version)
	}

	if !t.Failed() {
		fmt.Println("EditStickRequestRule successful")
	}

	// TestDeleteStickRequestRule
	err = client.DeleteStickRequestRule(4, "test", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ := client.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	_, err = client.GetStickRequestRule(4, "test", "")
	if err == nil {
		t.Error("DeleteStickRequestRule failed, stick request rule 3 still exists")
	}

	err = client.DeleteStickRequestRule(4, "test2", "", version)
	if err == nil {
		t.Error("Should throw error, non existant stick request rule")
		version++
	}

	if !t.Failed() {
		fmt.Println("DeleteStickRequestRule successful")
	}
}
