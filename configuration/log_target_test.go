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

func TestGetLogTargets(t *testing.T) {
	v, lTargets, err := client.GetLogTargets("frontend", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if len(lTargets) != 3 {
		t.Errorf("%v log targets returned, expected 3", len(lTargets))
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	for _, l := range lTargets {
		switch *l.Index {
		case 0:
			if l.Global != true {
				t.Errorf("%v: Global not true: %v", *l.Index, l.Global)
			}
		case 1:
			if l.Nolog != true {
				t.Errorf("%v: Nolog not true: %v", *l.Index, l.Nolog)
			}
		case 2:
			if l.Address != "127.0.0.1:514" {
				t.Errorf("%v: Address not 127.0.0.1:514: %v", *l.Index, l.Address)
			}
			if l.Facility != "local0" {
				t.Errorf("%v: Facility not local0: %v", *l.Index, l.Facility)
			}
			if l.Level != "notice" {
				t.Errorf("%v: Level not notice: %v", *l.Index, l.Level)
			}
			if l.Minlevel != "notice" {
				t.Errorf("%v: Minlevel not notice: %v", *l.Index, l.Minlevel)
			}
		default:
			t.Errorf("Expext only log 0, 1, or 2, %v found", *l.Index)
		}
	}

	_, lTargets, err = client.GetLogTargets("backend", "test_2", "")
	if err != nil {
		t.Error(err.Error())
	}
	if len(lTargets) > 0 {
		t.Errorf("%v log targets returned, expected 0", len(lTargets))
	}
}

func TestGetLogTarget(t *testing.T) {
	v, l, err := client.GetLogTarget(2, "frontend", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	if *l.Index != 2 {
		t.Errorf("Log Target ID not 2, %v found", *l.Index)
	}
	if l.Address != "127.0.0.1:514" {
		t.Errorf("%v: Address not 127.0.0.1:514: %v", *l.Index, l.Address)
	}
	if l.Facility != "local0" {
		t.Errorf("%v: Facility not local0: %v", *l.Index, l.Facility)
	}
	if l.Level != "notice" {
		t.Errorf("%v: Level not notice: %v", *l.Index, l.Level)
	}
	if l.Minlevel != "notice" {
		t.Errorf("%v: Minlevel not notice: %v", *l.Index, l.Minlevel)
	}

	_, err = l.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	_, _, err = client.GetLogTarget(3, "backend", "test_2", "")
	if err == nil {
		t.Error("Should throw error, non existant Log Target")
	}
}

func TestCreateEditDeleteLogTarget(t *testing.T) {
	id := int64(3)

	// TestCreateLogTarget
	r := &models.LogTarget{
		Index:    &id,
		Address:  "stdout",
		Format:   "raw",
		Facility: "daemon",
		Level:    "notice",
	}

	err := client.CreateLogTarget("frontend", "test", r, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, ondiskR, err := client.GetLogTarget(3, "frontend", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(ondiskR, r) {
		fmt.Printf("Created Log Target: %v\n", ondiskR)
		fmt.Printf("Given Log Target: %v\n", r)
		t.Error("Created Log Target not equal to given Log Target")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	// TestEditLogTarget
	r = &models.LogTarget{
		Index:    &id,
		Address:  "stdout",
		Format:   "rfc3164",
		Facility: "local1",
	}

	err = client.EditLogTarget(3, "frontend", "test", r, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, ondiskR, err = client.GetLogTarget(3, "frontend", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(ondiskR, r) {
		fmt.Printf("Edited Log Target: %v\n", ondiskR)
		fmt.Printf("Given Log Target: %v\n", r)
		t.Error("Edited Log Target not equal to given Log Target")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	// TestDeleteFilter
	err = client.DeleteLogTarget(3, "frontend", "test", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ := client.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	_, _, err = client.GetLogTarget(3, "frontend", "test", "")
	if err == nil {
		t.Error("DeleteLogTarget failed, Log Target 3 still exists")
	}

	err = client.DeleteLogTarget(2, "backend", "test_2", "", version)
	if err == nil {
		t.Error("Should throw error, non existant Log Target")
		version++
	}
}
