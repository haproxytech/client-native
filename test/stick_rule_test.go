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
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/haproxytech/client-native/v6/models"
	"github.com/stretchr/testify/require"
)

func stickRuleExpectation() map[string]models.StickRules {
	initStructuredExpected()
	res := StructuredToStickRuleMap()
	// Add individual entries
	for k, vs := range res {
		for i, v := range vs {
			key := fmt.Sprintf("%s/%d", k, i)
			res[key] = models.StickRules{v}
		}
	}
	return res
}

func TestGetStickRules(t *testing.T) { //nolint:gocognit,gocyclo
	mrules := make(map[string]models.StickRules)
	v, sRules, err := clientTest.GetStickRules("test", "")
	if err != nil {
		t.Error(err.Error())
	}
	mrules["backend/test"] = sRules

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	_, sRules, err = clientTest.GetStickRules("test_2", "")
	if err != nil {
		t.Error(err.Error())
	}
	mrules["backend/test_2"] = sRules

	checkStickRules(t, mrules)
}

func checkStickRules(t *testing.T, got map[string]models.StickRules) {
	exp := stickRuleExpectation()
	for k, v := range got {
		want, ok := exp[k]
		require.True(t, ok, "k=%s", k)
		require.Equal(t, len(want), len(v), "k=%s", k)
		i := 0
		for _, w := range want {
			require.True(t, v[i].Equal(*w), "k=%s - diff %v", k, cmp.Diff(*v[i], *w))
			i++
		}
	}
}

func TestGetStickRule(t *testing.T) {
	m := make(map[string]models.StickRules)

	v, sr, err := clientTest.GetStickRule(0, "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}
	m["backend/test/0"] = models.StickRules{sr}
	checkStickRules(t, m)

	_, err = sr.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	_, _, err = clientTest.GetStickRule(5, "test_2", "")
	if err == nil {
		t.Error("Should throw error, non existent stick rule")
	}
}

func TestCreateEditDeleteStickRule(t *testing.T) {
	id := int64(1)
	// TestCreateStickRule
	sr := &models.StickRule{
		Type:     "match",
		Pattern:  "src",
		Cond:     "if",
		CondTest: "TRUE",
	}

	err := clientTest.CreateStickRule(id, "test", sr, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, sRule, err := clientTest.GetStickRule(1, "test", "")
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
		Type:     "store-request",
		Pattern:  "src",
		Table:    "test2",
		Cond:     "if",
		CondTest: "FALSE",
	}

	err = clientTest.EditStickRule(1, "test", sr, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, sRule, err = clientTest.GetStickRule(1, "test", "")
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
	err = clientTest.DeleteStickRule(6, "test", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ := clientTest.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	_, _, err = clientTest.GetStickRule(6, "test", "")
	if err == nil {
		t.Error("DeleteStickRule failed, stick rule 3 still exists")
	}

	err = clientTest.DeleteStickRule(6, "test_2", "", version)
	if err == nil {
		t.Error("Should throw error, non existent stick rule")
		version++
	}
}
