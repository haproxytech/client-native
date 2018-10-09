package configuration

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/haproxytech/models"
)

func TestGetServerSwitchingRules(t *testing.T) {
	srvRules, err := client.GetServerSwitchingRules("test")
	if err != nil {
		t.Error(err.Error())
	}

	if len(srvRules.Data) != 2 {
		t.Errorf("%v server switching rules returned, expected 2", len(srvRules.Data))
	}

	if srvRules.Version != version {
		t.Errorf("Version %v returned, expected %v", srvRules.Version, version)
	}

	for _, sr := range srvRules.Data {
		if sr.ID == 1 {
			if sr.TargetServer != "webserv" {
				t.Errorf("%v: TargetServer not webserv: %v", sr.ID, sr.TargetServer)
			}
			if sr.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", sr.ID, sr.Cond)
			}
			if sr.CondTest != "TRUE" {
				t.Errorf("%v: CondTest not TRUE: %v", sr.ID, sr.CondTest)
			}
		} else if sr.ID == 2 {
			if sr.TargetServer != "webserv2" {
				t.Errorf("%v: TargetServer not webserv2: %v", sr.ID, sr.TargetServer)
			}
			if sr.Cond != "unless" {
				t.Errorf("%v: Cond not if: %v", sr.ID, sr.Cond)
			}
			if sr.CondTest != "TRUE" {
				t.Errorf("%v: CondTest not TRUE: %v", sr.ID, sr.CondTest)
			}
		} else {
			t.Errorf("Expext only server switching rule 1 or 2, %v found", sr.ID)
		}
	}

	srJSON, err := srvRules.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	srvRules, err = client.GetServerSwitchingRules("test2")
	if err != nil {
		t.Error(err.Error())
	}
	if len(srvRules.Data) > 0 {
		t.Errorf("%v server switching rules returned, expected 0", len(srvRules.Data))
	}

	if !t.Failed() {
		fmt.Println("GetServerSwitchingRules succesful\nResponse: \n" + string(srJSON) + "\n")
	}
}

func TestGetServerSwitchingRule(t *testing.T) {
	srvRule, err := client.GetServerSwitchingRule(1, "test")
	if err != nil {
		t.Error(err.Error())
	}

	sr := srvRule.Data

	if sr.TargetServer != "webserv" {
		t.Errorf("%v: TargetServer not webserv: %v", sr.ID, sr.TargetServer)
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

	_, err = client.GetServerSwitchingRule(3, "test2")
	if err == nil {
		t.Error("Should throw error, non existant server switching rule")
	}

	if !t.Failed() {
		fmt.Println("GetServerSwitchingRule succesful\nResponse: \n" + string(srJSON) + "\n")
	}
}

func TestCreateEditDeleteServerSwitchingRule(t *testing.T) {
	// TestCreateServerSwitchingRule
	sr := &models.ServerSwitchingRule{
		ID:           2,
		TargetServer: "webserv2",
		Cond:         "unless",
		CondTest:     "TRUE",
	}

	err := client.CreateServerSwitchingRule("test", sr, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	srvRule, err := client.GetServerSwitchingRule(2, "test")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(srvRule.Data, sr) {
		fmt.Printf("Created server switching rule: %v\n", srvRule.Data)
		fmt.Printf("Given server switching rule: %v\n", sr)
		t.Error("Created server switching rule not equal to given server switching rule")
	}

	if srvRule.Version != version {
		t.Errorf("Version %v returned, expected %v", srvRule.Version, version)
	}

	if !t.Failed() {
		fmt.Println("CreateServerSwitchingRule successful")
	}

	// TestServerSwitchingRule
	sr = &models.ServerSwitchingRule{
		ID:           2,
		TargetServer: "webserv2",
		Cond:         "if",
		CondTest:     "TRUE",
	}

	err = client.EditServerSwitchingRule(2, "test", sr, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	srvRule, err = client.GetServerSwitchingRule(2, "test")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(srvRule.Data, sr) {
		fmt.Printf("Edited server switching rule: %v\n", srvRule.Data)
		fmt.Printf("Given server switching rule: %v\n", sr)
		t.Error("Edited server switching rule not equal to given server switching rule")
	}

	if srvRule.Version != version {
		t.Errorf("Version %v returned, expected %v", srvRule.Version, version)
	}

	if !t.Failed() {
		fmt.Println("EditServerSwitchingRule successful")
	}

	// TestDeleteServerSwitchingRule
	err = client.DeleteServerSwitchingRule(3, "test", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ := client.GetVersion(); v != version {
		t.Error("Version not incremented")
	}

	_, err = client.GetServerSwitchingRule(3, "test")
	if err == nil {
		t.Error("DeleteServerSwitchingRule failed, server switching rule 3 still exists")
	}

	err = client.DeleteServerSwitchingRule(3, "test2", "", version)
	if err == nil {
		t.Error("Should throw error, non existant server switching rule")
		version++
	}

	if !t.Failed() {
		fmt.Println("DeleteServerSwitchingRule successful")
	}
}
