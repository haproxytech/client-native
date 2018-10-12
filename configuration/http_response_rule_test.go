package configuration

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/haproxytech/models"
)

func TestGetHTTPResponseRules(t *testing.T) {
	hRules, err := client.GetHTTPResponseRules("frontend", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if len(hRules.Data) != 3 {
		t.Errorf("%v http response rules returned, expected 3", len(hRules.Data))
	}

	if hRules.Version != version {
		t.Errorf("Version %v returned, expected %v", hRules.Version, version)
	}

	for _, r := range hRules.Data {
		if r.ID == 1 {
			if r.Type != "allow" {
				t.Errorf("%v: Type not allow: %v", r.ID, r.Type)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", r.ID, r.Cond)
			}
			if r.CondTest != "src 192.168.0.0/16" {
				t.Errorf("%v: CondTest not src 192.168.0.0/16: %v", r.ID, r.CondTest)
			}
		} else if r.ID == 2 {
			if r.Type != "set-header" {
				t.Errorf("%v: Type not set-header: %v", r.ID, r.Type)
			}
			if r.HdrName != "X-SSL" {
				t.Errorf("%v: HdrName not X-SSL: %v", r.ID, r.HdrName)
			}
			if r.HdrValue != "%[ssl_fc]" {
				t.Errorf("%v: HdrValue not [ssl_fc]: %v", r.ID, r.HdrValue)
			}
		} else if r.ID == 3 {
			if r.Type != "set-var" {
				t.Errorf("%v: Type not set-var: %v", r.ID, r.Type)
			}
			if r.VarName != "req.my_var" {
				t.Errorf("%v: VarName not req.my_var: %v", r.ID, r.VarName)
			}
			if r.VarPattern != "req.fhdr(user-agent),lower" {
				t.Errorf("%v: VarPattern not req.fhdr(user-agent),lower: %v", r.ID, r.VarPattern)
			}
		} else {
			t.Errorf("Expext only filter 1, 2 or 3, %v found", r.ID)
		}
	}

	rJSON, err := hRules.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	hRules, err = client.GetHTTPResponseRules("backend", "test2", "")
	if err != nil {
		t.Error(err.Error())
	}
	if len(hRules.Data) > 0 {
		t.Errorf("%v filters returned, expected 0", len(hRules.Data))
	}

	if !t.Failed() {
		fmt.Println("GetHTTPResponseRules succesful\nResponse: \n" + string(rJSON) + "\n")
	}
}

func TestGetHTTPResponseRule(t *testing.T) {
	hRule, err := client.GetHTTPResponseRule(1, "frontend", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	r := hRule.Data

	if r.ID != 1 {
		t.Errorf("HTTPResponse Rule ID not 1, %v found", r.ID)
	}
	if r.Type != "allow" {
		t.Errorf("%v: Type not allow: %v", r.ID, r.Type)
	}
	if r.Cond != "if" {
		t.Errorf("%v: Cond not if: %v", r.ID, r.Cond)
	}
	if r.CondTest != "src 192.168.0.0/16" {
		t.Errorf("%v: CondTest not src 192.168.0.0/16: %v", r.ID, r.CondTest)
	}

	rJSON, err := r.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	_, err = client.GetHTTPRequestRule(3, "backend", "test2", "")
	if err == nil {
		t.Error("Should throw error, non existant HTTPResponse Rule")
	}

	if !t.Failed() {
		fmt.Println("GetHTTPResponseRule succesful\nResponse: \n" + string(rJSON) + "\n")
	}
}

func TestCreateEditDeleteHTTPResponseRule(t *testing.T) {
	// TestCreateHTTPResponseRule
	r := &models.HTTPResponseRule{
		ID:       1,
		Type:     "set-log-level",
		LogLevel: "alert",
	}

	err := client.CreateHTTPResponseRule("frontend", "test", r, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	ondiskR, err := client.GetHTTPResponseRule(1, "frontend", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(ondiskR.Data, r) {
		fmt.Printf("Created HTTP response rule: %v\n", ondiskR.Data)
		fmt.Printf("Given HTTP response rule: %v\n", r)
		t.Error("Created HTTP response rule not equal to given HTTP response rule")
	}

	if ondiskR.Version != version {
		t.Errorf("Version %v returned, expected %v", ondiskR.Version, version)
	}

	if !t.Failed() {
		fmt.Println("CreateHTTPResponseRule successful")
	}

	// TestEditHTTPResponseRule
	r = &models.HTTPResponseRule{
		ID:       1,
		Type:     "set-log-level",
		LogLevel: "warning",
	}

	err = client.EditHTTPResponseRule(1, "frontend", "test", r, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	ondiskR, err = client.GetHTTPResponseRule(1, "frontend", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(ondiskR.Data, r) {
		fmt.Printf("Edited HTTP response rule: %v\n", ondiskR.Data)
		fmt.Printf("Given HTTP response rule: %v\n", r)
		t.Error("Edited HTTP response rule not equal to given HTTP response rule")
	}

	if ondiskR.Version != version {
		t.Errorf("Version %v returned, expected %v", ondiskR.Version, version)
	}

	if !t.Failed() {
		fmt.Println("EditHTTPResponseRule successful")
	}

	// TestDeleteFilter
	err = client.DeleteHTTPResponseRule(4, "frontend", "test", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ := client.GetVersion(); v != version {
		t.Error("Version not incremented")
	}

	_, err = client.GetHTTPResponseRule(4, "frontend", "test", "")
	if err == nil {
		t.Error("DeleteHTTPResponseRule failed, HTTPResponse Rule 4 still exists")
	}

	err = client.DeleteHTTPResponseRule(2, "backend", "test2", "", version)
	if err == nil {
		t.Error("Should throw error, non existant HTTPResponse Rule")
		version++
	}

	if !t.Failed() {
		fmt.Println("DeleteHTTPResponseRule successful")
	}
}
