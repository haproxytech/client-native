package configuration

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/haproxytech/models"
)

func TestGetStickRules(t *testing.T) {
	sRules, err := client.GetStickRules("test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if len(sRules.Data) != 6 {
		t.Errorf("%v stick rules returned, expected 6", len(sRules.Data))
	}

	if sRules.Version != version {
		t.Errorf("Version %v returned, expected %v", sRules.Version, version)
	}

	for _, sr := range sRules.Data {
		if *sr.ID == 0 {
			if sr.Type != "store-request" {
				t.Errorf("%v: Type not store-request: %v", *sr.ID, sr.Type)
			}
			if *sr.Pattern != "src" {
				t.Errorf("%v: Pattern not src: %v", *sr.ID, *sr.Pattern)
			}
			if sr.Table != "test" {
				t.Errorf("%v: Table not test: %v", *sr.ID, sr.Table)
			}
		} else if *sr.ID == 1 {
			if sr.Type != "match" {
				t.Errorf("%v: Type not match: %v", *sr.ID, sr.Type)
			}
			if *sr.Pattern != "src" {
				t.Errorf("%v: Pattern not src: %v", *sr.ID, *sr.Pattern)
			}
			if sr.Table != "test" {
				t.Errorf("%v: Table not test: %v", *sr.ID, sr.Table)
			}
		} else if *sr.ID == 2 {
			if sr.Type != "on" {
				t.Errorf("%v: Type not on: %v", *sr.ID, sr.Type)
			}
			if *sr.Pattern != "src" {
				t.Errorf("%v: Pattern not src: %v", *sr.ID, *sr.Pattern)
			}
			if sr.Table != "test" {
				t.Errorf("%v: Table not test: %v", *sr.ID, sr.Table)
			}
		} else if *sr.ID == 3 {
			if sr.Type != "store-response" {
				t.Errorf("%v: Type not matchandstore: %v", *sr.ID, sr.Type)
			}
			if *sr.Pattern != "src" {
				t.Errorf("%v: Pattern not src: %v", *sr.ID, *sr.Pattern)
			}
		} else if *sr.ID == 4 {
			if sr.Type != "store-response" {
				t.Errorf("%v: Type not matchandstore: %v", *sr.ID, sr.Type)
			}
			if *sr.Pattern != "src_port" {
				t.Errorf("%v: Pattern not src: %v", *sr.ID, *sr.Pattern)
			}
			if sr.Table != "test_port" {
				t.Errorf("%v: Table not test: %v", *sr.ID, sr.Table)
			}
		} else if *sr.ID == 5 {
			if sr.Type != "store-response" {
				t.Errorf("%v: Type not matchandstore: %v", *sr.ID, sr.Type)
			}
			if *sr.Pattern != "src" {
				t.Errorf("%v: Pattern not src: %v", *sr.ID, *sr.Pattern)
			}
			if sr.Table != "test" {
				t.Errorf("%v: Table not test: %v", *sr.ID, sr.Table)
			}
			if sr.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *sr.ID, sr.Cond)
			}
			if sr.CondTest != "TRUE" {
				t.Errorf("%v: Cond not if: %v", *sr.ID, sr.CondTest)
			}
		} else {
			t.Errorf("Expext only stick rule < 5, %v found", *sr.ID)
		}
	}

	srJSON, err := sRules.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	sRules, err = client.GetStickRules("test_2", "")
	if err != nil {
		t.Error(err.Error())
	}
	if len(sRules.Data) > 0 {
		t.Errorf("%v stick rules returned, expected 0", len(sRules.Data))
	}

	if !t.Failed() {
		fmt.Println("GetStickRules succesful\nResponse: \n" + string(srJSON) + "\n")
	}
}

func TestGetStickRule(t *testing.T) {
	sRule, err := client.GetStickRule(0, "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	sr := sRule.Data

	if sr.Type != "store-request" {
		t.Errorf("%v: Type not store-request: %v", *sr.ID, sr.Type)
	}
	if *sr.Pattern != "src" {
		t.Errorf("%v: Pattern not src: %v", *sr.ID, *sr.Pattern)
	}
	if sr.Table != "test" {
		t.Errorf("%v: Table not test: %v", *sr.ID, sr.Table)
	}

	srJSON, err := sr.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	_, err = client.GetStickRule(5, "test_2", "")
	if err == nil {
		t.Error("Should throw error, non existant stick rule")
	}

	if !t.Failed() {
		fmt.Println("GetStickRule succesful\nResponse: \n" + string(srJSON) + "\n")
	}
}

func TestCreateEditDeleteStickRule(t *testing.T) {
	id := int64(1)
	p := "src"
	// TestCreateStickRule
	sr := &models.StickRule{
		ID:       &id,
		Type:     "match",
		Pattern:  &p,
		Cond:     "if",
		CondTest: "TRUE",
	}

	err := client.CreateStickRule("test", sr, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	sRule, err := client.GetStickRule(1, "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(sRule.Data, sr) {
		fmt.Printf("Created stick rule: %v\n", sRule.Data)
		fmt.Printf("Given stick rule: %v\n", sr)
		t.Error("Created stick rule not equal to given stick rule")
	}

	if sRule.Version != version {
		t.Errorf("Version %v returned, expected %v", sRule.Version, version)
	}

	if !t.Failed() {
		fmt.Println("CreateStickRule successful")
	}

	// TestEditStickRule
	sr = &models.StickRule{
		ID:       &id,
		Type:     "store-request",
		Pattern:  &p,
		Table:    "test2",
		Cond:     "if",
		CondTest: "FALSE",
	}

	err = client.EditStickRule(1, "test", sr, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	sRule, err = client.GetStickRule(1, "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(sRule.Data, sr) {
		fmt.Printf("Edited stick rule: %v\n", sRule.Data)
		fmt.Printf("Given stick rule: %v\n", sr)
		t.Error("Edited stick rule not equal to given stick rule")
	}

	if sRule.Version != version {
		t.Errorf("Version %v returned, expected %v", sRule.Version, version)
	}

	if !t.Failed() {
		fmt.Println("EditStickRule successful")
	}

	// TestDeleteStickRule
	err = client.DeleteStickRule(6, "test", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ := client.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	_, err = client.GetStickRule(6, "test", "")
	if err == nil {
		t.Error("DeleteStickRule failed, stick rule 3 still exists")
	}

	err = client.DeleteStickRule(6, "test_2", "", version)
	if err == nil {
		t.Error("Should throw error, non existant stick rule")
		version++
	}

	if !t.Failed() {
		fmt.Println("DeleteStickRule successful")
	}
}
