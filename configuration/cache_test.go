// Copyright 2022 HAProxy Technologies
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

	"github.com/haproxytech/client-native/v4/misc"
	"github.com/haproxytech/client-native/v4/models"
)

func TestGetCaches(t *testing.T) {
	v, caches, err := clientTest.GetCaches("")
	if err != nil {
		t.Error(err.Error())
	}

	if len(caches) != 1 {
		t.Errorf("%v caches returned, expected 1", len(caches))
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	if *caches[0].Name != "test" {
		t.Errorf("Expected only test, %v found", caches[0].Name)
	}
	if caches[0].TotalMaxSize != 1024 {
		t.Errorf("Expected only test, %v found", caches[0].Name)
	}
	if caches[0].MaxObjectSize != 8 {
		t.Errorf("Expected only test, %v found", caches[0].Name)
	}
	if caches[0].MaxAge != 60 {
		t.Errorf("Expected only test, %v found", caches[0].Name)
	}
	if *caches[0].ProcessVary != true {
		t.Errorf("Expected only test, %v found", caches[0].Name)
	}
	if caches[0].MaxSecondaryEntries != 10 {
		t.Errorf("Expected only test, %v found", caches[0].Name)
	}
}

func TestGetCache(t *testing.T) {
	v, c, err := clientTest.GetCache("test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	if *c.Name != "test" {
		t.Errorf("Expected only test, %v found", c.Name)
	}
	if c.TotalMaxSize != 1024 {
		t.Errorf("Expected only test, %v found", c.Name)
	}
	if c.MaxObjectSize != 8 {
		t.Errorf("Expected only test, %v found", c.Name)
	}
	if c.MaxAge != 60 {
		t.Errorf("Expected only test, %v found", c.Name)
	}
	if *c.ProcessVary != true {
		t.Errorf("Expected only test, %v found", c.Name)
	}
	if c.MaxSecondaryEntries != 10 {
		t.Errorf("Expected only test, %v found", c.Name)
	}

	_, err = c.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	_, _, err = clientTest.GetCache("doesnotexist", "")
	if err == nil {
		t.Error("Should throw error, non existant caches section")
	}
}

func TestCreateEditDeleteCache(t *testing.T) {
	// test create
	f := &models.Cache{
		MaxAge:              60,
		MaxObjectSize:       2,
		MaxSecondaryEntries: 5,
		Name:                misc.StringP("created_cache"),
		ProcessVary:         misc.BoolP(false),
		TotalMaxSize:        2048,
	}
	err := clientTest.CreateCache(f, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, cache, err := clientTest.GetCache("created_cache", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(cache, f) {
		fmt.Printf("Created cache: %v\n", cache)
		fmt.Printf("Given cache: %v\n", f)
		t.Error("Created cache not equal to given cache")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	err = clientTest.CreateCache(f, "", version)
	if err == nil {
		t.Error("Should throw error cache already exists")
		version++
	}

	// edit cache
	f = &models.Cache{
		Name:         misc.StringP("created_cache"),
		TotalMaxSize: 1024,
	}
	err = clientTest.EditCache("created_cache", f, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, cache, err = clientTest.GetCache("created_cache", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(cache, f) {
		fmt.Printf("Created cache: %v\n", cache)
		fmt.Printf("Given cache: %v\n", f)
		t.Error("Created cache not equal to given cache")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	// delete cache
	err = clientTest.DeleteCache("created_cache", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ := clientTest.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	err = clientTest.DeleteCache("created_cache", "", 999999)
	if err != nil {
		if confErr, ok := err.(*ConfError); ok {
			if confErr.Code() != ErrVersionMismatch {
				t.Error("Should throw ErrVersionMismatch error")
			}
		} else {
			t.Error("Should throw ErrVersionMismatch error")
		}
	}
	_, _, err = clientTest.GetCache("created_cache", "")
	if err == nil {
		t.Error("DeleteCache failed, cache created_cache still exists")
	}

	err = clientTest.DeleteCache("doesnotexist", "", version)
	if err == nil {
		t.Error("Should throw error, non existant cache")
		version++
	}
}
