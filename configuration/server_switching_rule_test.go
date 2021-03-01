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

func TestGetServerSwitchingRules(t *testing.T) {
	v, srvRules, err := client.GetServerSwitchingRules("test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if len(srvRules) != 2 {
		t.Errorf("%v server switching rules returned, expected 2", len(srvRules))
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	for _, sr := range srvRules {
		switch *sr.Index {
		case 0:
			if sr.TargetServer != "webserv" {
				t.Errorf("%v: TargetServer not webserv: %v", *sr.Index, sr.TargetServer)
			}
			if sr.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *sr.Index, sr.Cond)
			}
			if sr.CondTest != "TRUE" {
				t.Errorf("%v: CondTest not TRUE: %v", *sr.Index, sr.CondTest)
			}
		case 1:
			if sr.TargetServer != "webserv2" {
				t.Errorf("%v: TargetServer not webserv2: %v", *sr.Index, sr.TargetServer)
			}
			if sr.Cond != "unless" {
				t.Errorf("%v: Cond not if: %v", *sr.Index, sr.Cond)
			}
			if sr.CondTest != "TRUE" {
				t.Errorf("%v: CondTest not TRUE: %v", *sr.Index, sr.CondTest)
			}
		default:
			t.Errorf("Expext only server switching rule 0 or 1, %v found", *sr.Index)
		}
	}

	_, srvRules, err = client.GetServerSwitchingRules("test_2", "")
	if err != nil {
		t.Error(err.Error())
	}
	if len(srvRules) > 0 {
		t.Errorf("%v server switching rules returned, expected 0", len(srvRules))
	}
}

func TestGetServerSwitchingRule(t *testing.T) {
	v, sr, err := client.GetServerSwitchingRule(0, "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	if sr.TargetServer != "webserv" {
		t.Errorf("%v: TargetServer not webserv: %v", sr.Index, sr.TargetServer)
	}
	if sr.Cond != "if" {
		t.Errorf("%v: Cond not if: %v", sr.Index, sr.Cond)
	}
	if sr.CondTest != "TRUE" {
		t.Errorf("%v: CondTest not TRUE: %v", sr.Index, sr.CondTest)
	}

	_, err = sr.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	_, _, err = client.GetServerSwitchingRule(3, "test2", "")
	if err == nil {
		t.Error("Should throw error, non existant server switching rule")
	}
}

func TestCreateEditDeleteServerSwitchingRule(t *testing.T) {
	id := int64(2)
	// TestCreateServerSwitchingRule
	sr := &models.ServerSwitchingRule{
		Index:        &id,
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

	v, srvRule, err := client.GetServerSwitchingRule(2, "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(srvRule, sr) {
		fmt.Printf("Created server switching rule: %v\n", srvRule)
		fmt.Printf("Given server switching rule: %v\n", sr)
		t.Error("Created server switching rule not equal to given server switching rule")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	// TestServerSwitchingRule
	sr = &models.ServerSwitchingRule{
		Index:        &id,
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

	v, srvRule, err = client.GetServerSwitchingRule(2, "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(srvRule, sr) {
		fmt.Printf("Edited server switching rule: %v\n", srvRule)
		fmt.Printf("Given server switching rule: %v\n", sr)
		t.Error("Edited server switching rule not equal to given server switching rule")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	// TestDeleteServerSwitchingRule
	err = client.DeleteServerSwitchingRule(2, "test", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ := client.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	_, _, err = client.GetServerSwitchingRule(2, "test", "")
	if err == nil {
		t.Error("DeleteServerSwitchingRule failed, server switching rule 3 still exists")
	}

	err = client.DeleteServerSwitchingRule(2, "test_2", "", version)
	if err == nil {
		t.Error("Should throw error, non existant server switching rule")
		version++
	}
}
