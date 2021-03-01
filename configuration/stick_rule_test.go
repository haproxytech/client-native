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

func TestGetStickRules(t *testing.T) { //nolint:gocognit,gocyclo
	v, sRules, err := client.GetStickRules("test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if len(sRules) != 6 {
		t.Errorf("%v stick rules returned, expected 6", len(sRules))
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	for _, sr := range sRules {
		switch *sr.Index {
		case 0:
			if sr.Type != "store-request" {
				t.Errorf("%v: Type not store-request: %v", *sr.Index, sr.Type)
			}
			if sr.Pattern != "src" {
				t.Errorf("%v: Pattern not src: %v", *sr.Index, sr.Pattern)
			}
			if sr.Table != "test" {
				t.Errorf("%v: Table not test: %v", *sr.Index, sr.Table)
			}
		case 1:
			if sr.Type != "match" {
				t.Errorf("%v: Type not match: %v", *sr.Index, sr.Type)
			}
			if sr.Pattern != "src" {
				t.Errorf("%v: Pattern not src: %v", *sr.Index, sr.Pattern)
			}
			if sr.Table != "test" {
				t.Errorf("%v: Table not test: %v", *sr.Index, sr.Table)
			}
		case 2:
			if sr.Type != "on" {
				t.Errorf("%v: Type not on: %v", *sr.Index, sr.Type)
			}
			if sr.Pattern != "src" {
				t.Errorf("%v: Pattern not src: %v", *sr.Index, sr.Pattern)
			}
			if sr.Table != "test" {
				t.Errorf("%v: Table not test: %v", *sr.Index, sr.Table)
			}
		case 3:
			if sr.Type != "store-response" {
				t.Errorf("%v: Type not matchandstore: %v", *sr.Index, sr.Type)
			}
			if sr.Pattern != "src" {
				t.Errorf("%v: Pattern not src: %v", *sr.Index, sr.Pattern)
			}
		case 4:
			if sr.Type != "store-response" {
				t.Errorf("%v: Type not matchandstore: %v", *sr.Index, sr.Type)
			}
			if sr.Pattern != "src_port" {
				t.Errorf("%v: Pattern not src: %v", *sr.Index, sr.Pattern)
			}
			if sr.Table != "test_port" {
				t.Errorf("%v: Table not test: %v", *sr.Index, sr.Table)
			}
		case 5:
			if sr.Type != "store-response" {
				t.Errorf("%v: Type not matchandstore: %v", *sr.Index, sr.Type)
			}
			if sr.Pattern != "src" {
				t.Errorf("%v: Pattern not src: %v", *sr.Index, sr.Pattern)
			}
			if sr.Table != "test" {
				t.Errorf("%v: Table not test: %v", *sr.Index, sr.Table)
			}
			if sr.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *sr.Index, sr.Cond)
			}
			if sr.CondTest != "TRUE" {
				t.Errorf("%v: Cond not if: %v", *sr.Index, sr.CondTest)
			}
		default:
			t.Errorf("Expext only stick rule < 5, %v found", *sr.Index)
		}
	}

	_, sRules, err = client.GetStickRules("test_2", "")
	if err != nil {
		t.Error(err.Error())
	}
	if len(sRules) > 0 {
		t.Errorf("%v stick rules returned, expected 0", len(sRules))
	}
}

func TestGetStickRule(t *testing.T) {
	v, sr, err := client.GetStickRule(0, "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	if sr.Type != "store-request" {
		t.Errorf("%v: Type not store-request: %v", *sr.Index, sr.Type)
	}
	if sr.Pattern != "src" {
		t.Errorf("%v: Pattern not src: %v", *sr.Index, sr.Pattern)
	}
	if sr.Table != "test" {
		t.Errorf("%v: Table not test: %v", *sr.Index, sr.Table)
	}

	_, err = sr.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	_, _, err = client.GetStickRule(5, "test_2", "")
	if err == nil {
		t.Error("Should throw error, non existant stick rule")
	}
}

func TestCreateEditDeleteStickRule(t *testing.T) {
	id := int64(1)
	// TestCreateStickRule
	sr := &models.StickRule{
		Index:    &id,
		Type:     "match",
		Pattern:  "src",
		Cond:     "if",
		CondTest: "TRUE",
	}

	err := client.CreateStickRule("test", sr, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, sRule, err := client.GetStickRule(1, "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(sRule, sr) {
		fmt.Printf("Created stick rule: %v\n", sRule)
		fmt.Printf("Given stick rule: %v\n", sr)
		t.Error("Created stick rule not equal to given stick rule")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	// TestEditStickRule
	sr = &models.StickRule{
		Index:    &id,
		Type:     "store-request",
		Pattern:  "src",
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

	v, sRule, err = client.GetStickRule(1, "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(sRule, sr) {
		fmt.Printf("Edited stick rule: %v\n", sRule)
		fmt.Printf("Given stick rule: %v\n", sr)
		t.Error("Edited stick rule not equal to given stick rule")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
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

	_, _, err = client.GetStickRule(6, "test", "")
	if err == nil {
		t.Error("DeleteStickRule failed, stick rule 3 still exists")
	}

	err = client.DeleteStickRule(6, "test_2", "", version)
	if err == nil {
		t.Error("Should throw error, non existant stick rule")
		version++
	}
}
