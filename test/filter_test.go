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

package test

import (
	_ "embed"
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/haproxytech/client-native/v6/configuration"
	"github.com/haproxytech/client-native/v6/models"
	"github.com/stretchr/testify/require"
)

func filterExpectation() map[string]models.Filters {
	initStructuredExpected()
	res := StructuredToFilterMap()
	// Add individual entries
	for k, vs := range res {
		for _, v := range vs {
			key := fmt.Sprintf("%s/%d", k, *v.Index)
			res[key] = models.Filters{v}
		}
	}
	return res
}

func TestGetFilters(t *testing.T) { //nolint:gocognit
	mfilters := make(map[string]models.Filters)
	v, filters, err := clientTest.GetFilters(configuration.FrontendParentName, "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}
	mfilters["frontend/test"] = filters

	_, filters, err = clientTest.GetFilters(configuration.BackendParentName, "test_2", "")
	if err != nil {
		t.Error(err.Error())
	}
	mfilters["backend/test_2"] = filters

	checkFilters(t, mfilters)
}

func checkFilters(t *testing.T, got map[string]models.Filters) {
	exp := filterExpectation()
	for k, v := range got {
		want, ok := exp[k]
		require.True(t, ok, "k=%s", k)
		require.Equal(t, len(want), len(v), "k=%s", k)
		for _, g := range v {
			for _, w := range want {
				if *g.Index == *w.Index {
					require.True(t, g.Equal(*w), "k=%s - diff %v", k, cmp.Diff(*g, *w))
					break
				}
			}
		}
	}
}

func TestGetFilter(t *testing.T) {
	m := make(map[string]models.Filters)
	v, f, err := clientTest.GetFilter(0, configuration.FrontendParentName, "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}
	m["frontend/test/0"] = models.Filters{f}

	_, err = f.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}
	_, _, err = clientTest.GetFilter(3, configuration.BackendParentName, "test2", "")
	if err == nil {
		t.Error("Should throw error, non existent filter")
	}

	checkFilters(t, m)
}

func TestCreateEditDeleteFilter(t *testing.T) {
	// TestCreateFilter
	id := int64(1)
	f := &models.Filter{
		Index:      &id,
		Type:       "spoe",
		SpoeEngine: "test",
		SpoeConfig: "test.cfg",
	}

	err := clientTest.CreateFilter(configuration.FrontendParentName, "test", f, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, ondiskF, err := clientTest.GetFilter(1, configuration.FrontendParentName, "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(ondiskF, f) {
		fmt.Printf("Created filter: %v\n", ondiskF)
		fmt.Printf("Given filter: %v\n", f)
		t.Error("Created filter not equal to given filter")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	// TestEditFilter
	f = &models.Filter{
		Index:      &id,
		Type:       "spoe",
		SpoeConfig: "bla.cfg",
		SpoeEngine: "bla",
	}

	err = clientTest.EditFilter(1, configuration.FrontendParentName, "test", f, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, ondiskF, err = clientTest.GetFilter(1, configuration.FrontendParentName, "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(ondiskF, f) {
		fmt.Printf("Edited filter: %v\n", ondiskF)
		fmt.Printf("Given filter: %v\n", f)
		t.Error("Edited filter not equal to given filter")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	// TestDeleteFilter
	err = clientTest.DeleteFilter(3, "frontend", "test", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ := clientTest.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	_, filters, _ := clientTest.GetFilters(configuration.FrontendParentName, "test", "")
	_ = filters

	_, _, err = clientTest.GetFilter(6, "frontend", "test", "")
	if err == nil {
		t.Error("DeleteFilter failed, filter 5 still exists")
	}

	err = clientTest.DeleteFilter(1, configuration.BackendParentName, "test_2", "", version)
	if err == nil {
		t.Error("Should throw error, non existent filter")
		version++
	}
}
