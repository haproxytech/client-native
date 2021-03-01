// Copyright 2019 HAProxy Technologies
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package configuration

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/haproxytech/client-native/v2/models"
)

func TestGetACLs(t *testing.T) { //nolint:gocognit
	v, acls, err := client.GetACLs("frontend", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if len(acls) != 3 {
		t.Errorf("%v ACL rules returned, expected 3", len(acls))
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	for _, r := range acls {
		switch *r.Index {
		case 0:
			if r.ACLName != "invalid_src" {
				t.Errorf("%v: ACLName not invalid_src: %v", *r.Index, r.ACLName)
			}
			if r.Value != "0.0.0.0/7 224.0.0.0/3" {
				t.Errorf("%v: Value not 0.0.0.0/7 224.0.0.0/3: %v", *r.Index, r.Value)
			}
			if r.Criterion != "src" {
				t.Errorf("%v: Criterion not src: %v", *r.Index, r.Criterion)
			}
		case 1:
			if r.ACLName != "invalid_src" {
				t.Errorf("%v: ACLName not invalid_src: %v", *r.Index, r.ACLName)
			}
			if r.Value != "0:1023" {
				t.Errorf("%v: Value not 0:1023: %v", *r.Index, r.Value)
			}
			if r.Criterion != "src_port" {
				t.Errorf("%v: Criterion not src_port: %v", *r.Index, r.Criterion)
			}
		case 2:
			if r.ACLName != "local_dst" {
				t.Errorf("%v: ACLName not invalid_src: %v", *r.Index, r.ACLName)
			}
			if r.Value != "-i localhost" {
				t.Errorf("%v: Value not -i localhost: %v", *r.Index, r.Value)
			}
			if r.Criterion != "hdr(host)" {
				t.Errorf("%v: Criterion not hdr(host): %v", *r.Index, r.Criterion)
			}
		default:
			t.Errorf("Expext only acl 1, 2 or 3, %v found", *r.Index)
		}
	}

	_, acls, err = client.GetACLs("backend", "test_2", "")
	if err != nil {
		t.Error(err.Error())
	}
	if len(acls) > 0 {
		t.Errorf("%v ACLs returned, expected 0", len(acls))
	}
}

func TestGetACL(t *testing.T) {
	v, acl, err := client.GetACL(0, "frontend", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	if *acl.Index != 0 {
		t.Errorf("ACL ID not 0, %v found", *acl.Index)
	}
	if acl.ACLName != "invalid_src" {
		t.Errorf("%v: ACLName not invalid_src: %v", *acl.Index, acl.ACLName)
	}
	if acl.Value != "0.0.0.0/7 224.0.0.0/3" {
		t.Errorf("%v: Value not 0.0.0.0/7 224.0.0.0/3: %v", *acl.Index, acl.Value)
	}
	if acl.Criterion != "src" {
		t.Errorf("%v: Criterion not src: %v", *acl.Index, acl.Criterion)
	}

	_, err = acl.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	_, _, err = client.GetACL(3, "backend", "test_2", "")
	if err == nil {
		t.Error("Should throw error, non existent ACL")
	}
}

func TestCreateEditDeleteACL(t *testing.T) {
	id := int64(1)

	// TestCreateACL
	r := &models.ACL{
		Index:     &id,
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

	v, ondiskR, err := client.GetACL(1, "frontend", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(ondiskR, r) {
		fmt.Printf("Created ACL rule: %v\n", ondiskR)
		fmt.Printf("Given ACL rule: %v\n", r)
		t.Error("Created ACL rule not equal to given ACL rule")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	// TestEditACL
	r = &models.ACL{
		Index:     &id,
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

	v, ondiskR, err = client.GetACL(1, "frontend", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(ondiskR, r) {
		fmt.Printf("Edited ACL rule: %v\n", ondiskR)
		fmt.Printf("Given ACL rule: %v\n", r)
		t.Error("Edited ACL rule not equal to given ACL rule")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
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

	_, _, err = client.GetACL(3, "frontend", "test", "")
	if err == nil {
		t.Error("DeleteACL failed, ACL Rule 3 still exists")
	}

	err = client.DeleteACL(2, "backend", "test_2", "", version)
	if err == nil {
		t.Error("Should throw error, non existent ACL Rule")
		version++
	}
}
