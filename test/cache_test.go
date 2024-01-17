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

package test

import (
	_ "embed"
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/haproxytech/client-native/v5/configuration"
	"github.com/haproxytech/client-native/v5/misc"
	"github.com/haproxytech/client-native/v5/models"
	"github.com/stretchr/testify/require"
)

func cacheExpectation() map[string]models.Caches {
	initStructuredExpected()
	res := StructuredToCacheMap()
	// Add individual entries
	for _, vs := range res {
		for _, v := range vs {
			key := *v.Name
			res[key] = models.Caches{v}
		}
	}
	return res
}

func TestGetCaches(t *testing.T) {
	m := make(map[string]models.Caches)
	v, caches, err := clientTest.GetCaches("")
	if err != nil {
		t.Error(err.Error())
	}
	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}
	m[""] = caches
	checkCaches(t, m)
}

func checkCaches(t *testing.T, got map[string]models.Caches) {
	exp := cacheExpectation()
	for k, v := range got {
		want, ok := exp[k]
		require.True(t, ok, "k=%s", k)
		require.Equal(t, len(want), len(v), "k=%s", k)
		for _, g := range v {
			for _, w := range want {
				if g.Name == w.Name {
					require.True(t, g.Equal(*w), "k=%s - diff %v", k, cmp.Diff(*g, *w))
					break
				}
			}
		}
	}
}

func TestGetCache(t *testing.T) {
	m := make(map[string]models.Caches)

	v, c, err := clientTest.GetCache("test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}
	m["test"] = models.Caches{c}

	checkCaches(t, m)
	_, err = c.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	_, _, err = clientTest.GetCache("doesnotexist", "")
	if err == nil {
		t.Error("Should throw error, non existent caches section")
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
		if confErr, ok := err.(*configuration.ConfError); ok {
			if !confErr.Is(configuration.ErrVersionMismatch) {
				t.Error("Should throw configuration.ErrVersionMismatch error")
			}
		} else {
			t.Error("Should throw configuration.ErrVersionMismatch error")
		}
	}
	_, _, err = clientTest.GetCache("created_cache", "")
	if err == nil {
		t.Error("DeleteCache failed, cache created_cache still exists")
	}

	err = clientTest.DeleteCache("doesnotexist", "", version)
	if err == nil {
		t.Error("Should throw error, non existent cache")
		version++
	}
}
