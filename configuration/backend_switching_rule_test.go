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

	"github.com/haproxytech/models"
)

func TestGetBackendSwitchingRules(t *testing.T) {
	v, bckRules, err := client.GetBackendSwitchingRules("test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if len(bckRules) != 2 {
		t.Errorf("%v backend switching rules returned, expected 2", len(bckRules))
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	for _, br := range bckRules {
		if *br.ID == 0 {
			if br.Name != "test_2" {
				t.Errorf("%v: Name not test_2: %v", *br.ID, br.Name)
			}
			if br.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *br.ID, br.Cond)
			}
			if br.CondTest != "TRUE" {
				t.Errorf("%v: CondTest not TRUE: %v", *br.ID, br.CondTest)
			}
		} else if *br.ID == 1 {
			if br.Name != "%[req.cookie(foo)]" {
				t.Errorf("%v: Name not %%[req.cookie(foo)]: %v", *br.ID, br.Name)
			}
		}

	}

	_, bckRules, err = client.GetBackendSwitchingRules("test_2", "")
	if err != nil {
		t.Error(err.Error())
	}
	if len(bckRules) > 0 {
		t.Errorf("%v backend switching rules returned, expected 0", len(bckRules))
	}
}

func TestGetBackendSwitchingRule(t *testing.T) {
	v, br, err := client.GetBackendSwitchingRule(0, "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	if *br.ID != 0 {
		t.Errorf("ID only backend switching rule 0, %v found", *br.ID)
	}
	if br.Name != "test_2" {
		t.Errorf("%v: Name not test_2: %v", *br.ID, br.Name)
	}
	if br.Cond != "if" {
		t.Errorf("%v: Cond not if: %v", *br.ID, br.Cond)
	}
	if br.CondTest != "TRUE" {
		t.Errorf("%v: CondTest not TRUE: %v", *br.ID, br.CondTest)
	}

	_, err = br.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	_, _, err = client.GetBackendSwitchingRule(3, "test2", "")
	if err == nil {
		t.Error("Should throw error, non existant backend switching rule")
	}
}

func TestCreateEditDeleteBackendSwitchingRule(t *testing.T) {
	// TestCreateBackendSwitchingRule
	id := int64(2)
	br := &models.BackendSwitchingRule{
		ID:       &id,
		Name:     "test",
		Cond:     "unless",
		CondTest: "TRUE",
	}

	err := client.CreateBackendSwitchingRule("test", br, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, bckRule, err := client.GetBackendSwitchingRule(2, "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(bckRule, br) {
		fmt.Printf("Created backend switching rule: %v\n", bckRule)
		fmt.Printf("Given backend switching rule: %v\n", br)
		t.Error("Created backend switching rule not equal to given backend switching rule")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	// TestBackendSwitchingRule
	br = &models.BackendSwitchingRule{
		ID:   &id,
		Name: "%[req.cookie(foo)]",
	}

	err = client.EditBackendSwitchingRule(2, "test", br, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, bckRule, err = client.GetBackendSwitchingRule(2, "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(bckRule, br) {
		fmt.Printf("Edited backend switching rule: %v\n", bckRule)
		fmt.Printf("Given backend switching rule: %v\n", br)
		t.Error("Edited backend switching rule not equal to given backend switching rule")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
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

	_, _, err = client.GetBackendSwitchingRule(2, "test", "")
	if err == nil {
		t.Error("DeleteBackendSwitchingRule failed, backend switching rule 2 still exists")
	}

	err = client.DeleteBackendSwitchingRule(1, "test2", "", version)
	if err == nil {
		t.Error("Should throw error, non existant backend switching rule")
		version++
	}
}
