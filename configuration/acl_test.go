package configuration

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/haproxytech/models"
)

func TestGetACLs(t *testing.T) {
	acls, err := client.GetACLs("frontend", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if len(acls.Data) != 3 {
		t.Errorf("%v ACL rules returned, expected 3", len(acls.Data))
	}

	if acls.Version != version {
		t.Errorf("Version %v returned, expected %v", acls.Version, version)
	}

	for _, r := range acls.Data {
		if *r.ID == 0 {
			if r.ACLName != "invalid_src" {
				t.Errorf("%v: ACLName not invalid_src: %v", *r.ID, r.ACLName)
			}
			if r.Value != "0.0.0.0/7 224.0.0.0/3" {
				t.Errorf("%v: Value not 0.0.0.0/7 224.0.0.0/3: %v", *r.ID, r.Value)
			}
			if r.Criterion != "src" {
				t.Errorf("%v: Criterion not src: %v", *r.ID, r.Criterion)
			}
		} else if *r.ID == 1 {
			if r.ACLName != "invalid_src" {
				t.Errorf("%v: ACLName not invalid_src: %v", *r.ID, r.ACLName)
			}
			if r.Value != "0:1023" {
				t.Errorf("%v: Value not 0:1023: %v", *r.ID, r.Value)
			}
			if r.Criterion != "src_port" {
				t.Errorf("%v: Criterion not src_port: %v", *r.ID, r.Criterion)
			}
		} else if *r.ID == 2 {
			if r.ACLName != "local_dst" {
				t.Errorf("%v: ACLName not invalid_src: %v", *r.ID, r.ACLName)
			}
			if r.Value != "-i localhost" {
				t.Errorf("%v: Value not -i localhost: %v", *r.ID, r.Value)
			}
			if r.Criterion != "hdr(host)" {
				t.Errorf("%v: Criterion not hdr(host): %v", *r.ID, r.Criterion)
			}
		} else {
			t.Errorf("Expext only acl 1, 2 or 3, %v found", *r.ID)
		}
	}

	rJSON, err := acls.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	acls, err = client.GetACLs("backend", "test_2", "")
	if err != nil {
		t.Error(err.Error())
	}
	if len(acls.Data) > 0 {
		t.Errorf("%v ACLs returned, expected 0", len(acls.Data))
	}

	if !t.Failed() {
		fmt.Println("GetACLs succesful\nResponse: \n" + string(rJSON) + "\n")
	}
}

func TestGetACL(t *testing.T) {
	acl, err := client.GetACL(0, "frontend", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	r := acl.Data

	if *r.ID != 0 {
		t.Errorf("ACL ID not 0, %v found", *r.ID)
	}
	if r.ACLName != "invalid_src" {
		t.Errorf("%v: ACLName not invalid_src: %v", *r.ID, r.ACLName)
	}
	if r.Value != "0.0.0.0/7 224.0.0.0/3" {
		t.Errorf("%v: Value not 0.0.0.0/7 224.0.0.0/3: %v", *r.ID, r.Value)
	}
	if r.Criterion != "src" {
		t.Errorf("%v: Criterion not src: %v", *r.ID, r.Criterion)
	}

	rJSON, err := r.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	_, err = client.GetACL(3, "backend", "test_2", "")
	if err == nil {
		t.Error("Should throw error, non existant ACL")
	}

	if !t.Failed() {
		fmt.Println("GetACL succesful\nResponse: \n" + string(rJSON) + "\n")
	}
}

func TestCreateEditDeleteACL(t *testing.T) {
	id := int64(1)

	// TestCreateACL
	r := &models.ACL{
		ID:        &id,
		ACLName:   "site_dead",
		Criterion: "nbsrv(dynamic)",
		Value:     "lt 2",
	}

	err := client.CreateACL("frontend", "test", r, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	ondiskR, err := client.GetACL(1, "frontend", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(ondiskR.Data, r) {
		fmt.Printf("Created ACL rule: %v\n", ondiskR.Data)
		fmt.Printf("Given ACL rule: %v\n", r)
		t.Error("Created ACL rule not equal to given ACL rule")
	}

	if ondiskR.Version != version {
		t.Errorf("Version %v returned, expected %v", ondiskR.Version, version)
	}

	if !t.Failed() {
		fmt.Println("CreateACL successful")
	}

	// TestEditACL
	r = &models.ACL{
		ID:        &id,
		ACLName:   "site_dead",
		Criterion: "nbsrv(static)",
		Value:     "lt 4",
	}

	err = client.EditACL(1, "frontend", "test", r, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	ondiskR, err = client.GetACL(1, "frontend", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(ondiskR.Data, r) {
		fmt.Printf("Edited ACL rule: %v\n", ondiskR.Data)
		fmt.Printf("Given ACL rule: %v\n", r)
		t.Error("Edited ACL rule not equal to given ACL rule")
	}

	if ondiskR.Version != version {
		t.Errorf("Version %v returned, expected %v", ondiskR.Version, version)
	}

	if !t.Failed() {
		fmt.Println("EditACL successful")
	}

	// TestDeleteACL
	err = client.DeleteACL(3, "frontend", "test", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ := client.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	_, err = client.GetACL(3, "frontend", "test", "")
	if err == nil {
		t.Error("DeleteACL failed, ACL Rule 3 still exists")
	}

	err = client.DeleteACL(2, "backend", "test_2", "", version)
	if err == nil {
		t.Error("Should throw error, non existant ACL Rule")
		version++
	}

	if !t.Failed() {
		fmt.Println("DeleteACL successful")
	}
}
