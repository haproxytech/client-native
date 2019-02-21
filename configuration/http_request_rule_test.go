package configuration

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/haproxytech/models"
)

func TestGetHTTPRequestRules(t *testing.T) {
	hRules, err := client.GetHTTPRequestRules("frontend", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if len(hRules.Data) != 3 {
		t.Errorf("%v http request rules returned, expected 3", len(hRules.Data))
	}

	if hRules.Version != version {
		t.Errorf("Version %v returned, expected %v", hRules.Version, version)
	}

	for _, r := range hRules.Data {
		if *r.ID == 0 {
			if r.Type != "allow" {
				t.Errorf("%v: Type not allow: %v", *r.ID, r.Type)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.ID, r.Cond)
			}
			if r.CondTest != "src 192.168.0.0/16" {
				t.Errorf("%v: CondTest not src 192.168.0.0/16: %v", *r.ID, r.CondTest)
			}
		} else if *r.ID == 1 {
			if r.Type != "set-header" {
				t.Errorf("%v: Type not set-header: %v", *r.ID, r.Type)
			}
			if r.HdrName != "X-SSL" {
				t.Errorf("%v: HdrName not X-SSL: %v", *r.ID, r.HdrName)
			}
			if r.HdrFormat != "%[ssl_fc]" {
				t.Errorf("%v: HdrFormat not [ssl_fc]: %v", *r.ID, r.HdrFormat)
			}
		} else if *r.ID == 2 {
			if r.Type != "set-var" {
				t.Errorf("%v: Type not set-var: %v", *r.ID, r.Type)
			}
			if r.VarName != "my_var" {
				t.Errorf("%v: VarName not my_var: %v", *r.ID, r.VarName)
			}
			if r.VarScope != "req" {
				t.Errorf("%v: VarName not req: %v", *r.ID, r.VarScope)
			}
			if r.VarExpr != "req.fhdr(user-agent),lower" {
				t.Errorf("%v: VarPattern not req.fhdr(user-agent),lower: %v", *r.ID, r.VarExpr)
			}
		} else {
			t.Errorf("Expext only http-request 1, 2 or 3, %v found", *r.ID)
		}
	}

	rJSON, err := hRules.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	hRules, err = client.GetHTTPRequestRules("backend", "test_2", "")
	if err != nil {
		t.Error(err.Error())
	}
	if len(hRules.Data) > 0 {
		t.Errorf("%v HTTP Request Ruless returned, expected 0", len(hRules.Data))
	}

	if !t.Failed() {
		fmt.Println("GetHTTPRequestRules succesful\nResponse: \n" + string(rJSON) + "\n")
	}
}

func TestGetHTTPRequestRule(t *testing.T) {
	hRule, err := client.GetHTTPRequestRule(0, "frontend", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	r := hRule.Data

	if *r.ID != 0 {
		t.Errorf("HTTP Request Rule ID not 0, %v found", *r.ID)
	}
	if r.Type != "allow" {
		t.Errorf("%v: Type not allow: %v", *r.ID, r.Type)
	}
	if r.Cond != "if" {
		t.Errorf("%v: Cond not if: %v", *r.ID, r.Cond)
	}
	if r.CondTest != "src 192.168.0.0/16" {
		t.Errorf("%v: CondTest not src 192.168.0.0/16: %v", *r.ID, r.CondTest)
	}

	rJSON, err := r.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	_, err = client.GetHTTPRequestRule(3, "backend", "test_2", "")
	if err == nil {
		t.Error("Should throw error, non existant HTTP Request Rule")
	}

	if !t.Failed() {
		fmt.Println("GetHTTPRequestRule succesful\nResponse: \n" + string(rJSON) + "\n")
	}
}

func TestCreateEditDeleteHTTPRequestRule(t *testing.T) {
	id := int64(1)

	// TestCreateHTTPRequestRule
	r := &models.HTTPRequestRule{
		ID:         &id,
		Type:       "redirect",
		RedirCode:  301,
		RedirValue: "http://www.%[hdr(host)]%[capture.req.uri]",
		RedirType:  "location",
	}

	err := client.CreateHTTPRequestRule("frontend", "test", r, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	ondiskR, err := client.GetHTTPRequestRule(1, "frontend", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(ondiskR.Data, r) {
		fmt.Printf("Created HTTP request rule: %v\n", ondiskR.Data)
		fmt.Printf("Given HTTP request rule: %v\n", r)
		t.Error("Created HTTP request rule not equal to given HTTP request rule")
	}

	if ondiskR.Version != version {
		t.Errorf("Version %v returned, expected %v", ondiskR.Version, version)
	}

	if !t.Failed() {
		fmt.Println("CreateHTTPRequestRule successful")
	}

	// TestEditHTTPRequestRule
	r = &models.HTTPRequestRule{
		ID:         &id,
		Type:       "redirect",
		RedirCode:  302,
		RedirValue: "http://www1.%[hdr(host)]%[capture.req.uri]",
		RedirType:  "scheme",
	}

	err = client.EditHTTPRequestRule(1, "frontend", "test", r, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	ondiskR, err = client.GetHTTPRequestRule(1, "frontend", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(ondiskR.Data, r) {
		fmt.Printf("Edited HTTP request rule: %v\n", ondiskR.Data)
		fmt.Printf("Given HTTP request rule: %v\n", r)
		t.Error("Edited HTTP request rule not equal to given HTTP request rule")
	}

	if ondiskR.Version != version {
		t.Errorf("Version %v returned, expected %v", ondiskR.Version, version)
	}

	if !t.Failed() {
		fmt.Println("EditHTTPRequestRule successful")
	}

	// TestDeleteHTTPRequest
	err = client.DeleteHTTPRequestRule(3, "frontend", "test", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ := client.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	_, err = client.GetHTTPRequestRule(3, "frontend", "test", "")
	if err == nil {
		t.Error("DeleteHTTPRequestRule failed, HTTP Request Rule 3 still exists")
	}

	err = client.DeleteHTTPRequestRule(2, "backend", "test_2", "", version)
	if err == nil {
		t.Error("Should throw error, non existant HTTP Request Rule")
		version++
	}

	if !t.Failed() {
		fmt.Println("DeleteHTTPRequestRule successful")
	}
}
