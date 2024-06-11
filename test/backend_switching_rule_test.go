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
	"github.com/haproxytech/client-native/v6/models"
	"github.com/stretchr/testify/require"
)

func bsrExpectation() map[string]models.BackendSwitchingRules {
	initStructuredExpected()
	res := StructuredToBackendSwitchingRuleMap()
	// Add individual entries
	for k, vs := range res {
		for i, v := range vs {
			key := fmt.Sprintf("%s/%d", k, i)
			res[key] = models.BackendSwitchingRules{v}
		}
	}
	return res
}

func TestGetBackendSwitchingRules(t *testing.T) {
	mbsr := make(map[string]models.BackendSwitchingRules)

	v, bckRules, err := clientTest.GetBackendSwitchingRules("test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}
	mbsr["frontend/test"] = bckRules

	_, bckRules, err = clientTest.GetBackendSwitchingRules("test_2", "")
	if err != nil {
		t.Error(err.Error())
	}
	mbsr["frontend/test_2"] = bckRules

	checkBackendSwitchingRules(t, mbsr)
}

func checkBackendSwitchingRules(t *testing.T, got map[string]models.BackendSwitchingRules) {
	exp := bsrExpectation()
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

func TestGetBackendSwitchingRule(t *testing.T) {
	mbsr := make(map[string]models.BackendSwitchingRules)

	v, br, err := clientTest.GetBackendSwitchingRule(0, "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}
	mbsr["frontend/test/0"] = models.BackendSwitchingRules{br}

	checkBackendSwitchingRules(t, mbsr)

	_, err = br.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	_, _, err = clientTest.GetBackendSwitchingRule(3, "test2", "")
	if err == nil {
		t.Error("Should throw error, non existent backend switching rule")
	}
}

func TestCreateEditDeleteBackendSwitchingRule(t *testing.T) {
	// TestCreateBackendSwitchingRule
	id := int64(2)
	br := &models.BackendSwitchingRule{
		Name:     "test",
		Cond:     "unless",
		CondTest: "TRUE",
	}

	err := clientTest.CreateBackendSwitchingRule(id, "test", br, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, bckRule, err := clientTest.GetBackendSwitchingRule(2, "test", "")
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
		Name: "%[req.cookie(foo)]",
	}

	err = clientTest.EditBackendSwitchingRule(2, "test", br, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, bckRule, err = clientTest.GetBackendSwitchingRule(2, "test", "")
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
	err = clientTest.DeleteBackendSwitchingRule(2, "test", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ := clientTest.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	_, _, err = clientTest.GetBackendSwitchingRule(2, "test", "")
	if err == nil {
		t.Error("DeleteBackendSwitchingRule failed, backend switching rule 2 still exists")
	}

	err = clientTest.DeleteBackendSwitchingRule(1, "test2", "", version)
	if err == nil {
		t.Error("Should throw error, non existent backend switching rule")
		version++
	}
}
