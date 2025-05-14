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
	"github.com/haproxytech/client-native/v6/configuration"
	"github.com/haproxytech/client-native/v6/models"
	"github.com/stretchr/testify/require"
)

func tcpResponseRuleExpectation() map[string]models.TCPResponseRules {
	initStructuredExpected()
	res := StructuredToTCPResponseRuleMap()
	// Add individual entries
	for k, vs := range res {
		for i, v := range vs {
			key := fmt.Sprintf("%s/%d", k, i)
			res[key] = models.TCPResponseRules{v}
		}
	}
	return res
}

func TestGetTCPResponseRules(t *testing.T) { //nolint:gocognit
	mrules := make(map[string]models.TCPResponseRules)
	v, tRules, err := clientTest.GetTCPResponseRules(configuration.BackendParentName, "test", "")
	if err != nil {
		t.Error(err.Error())
	}
	mrules["backend/test"] = tRules

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	_, tRules, err = clientTest.GetTCPResponseRules(configuration.BackendParentName, "test_2", "")
	if err != nil {
		t.Error(err.Error())
	}
	mrules["backend/test_2"] = tRules

	checkTCPResponseRules(t, mrules)
}

func checkTCPResponseRules(t *testing.T, got map[string]models.TCPResponseRules) {
	exp := tcpResponseRuleExpectation()
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

func TestGetTCPResponseRule(t *testing.T) {
	m := make(map[string]models.TCPResponseRules)

	v, r, err := clientTest.GetTCPResponseRule(0, configuration.BackendParentName, "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}
	m["backend/test/0"] = models.TCPResponseRules{r}

	checkTCPResponseRules(t, m)

	_, err = r.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	_, _, err = clientTest.GetTCPResponseRule(3, configuration.BackendParentName, "test_2", "")
	if err == nil {
		t.Error("Should throw error, non existent TCP Response Rule")
	}
}

func TestCreateEditDeleteTCPResponseRule(t *testing.T) {
	id := int64(2)
	tOut := int64(1000)
	// TestCreateTCPResponseRule
	r := &models.TCPResponseRule{
		Type:    "inspect-delay",
		Timeout: &tOut,
	}

	err := clientTest.CreateTCPResponseRule(id, configuration.BackendParentName, "test", r, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, ondiskR, err := clientTest.GetTCPResponseRule(2, configuration.BackendParentName, "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(ondiskR, r) {
		fmt.Printf("Created TCP response rule: %v\n", ondiskR)
		fmt.Printf("Given TCP response rule: %v\n", r)
		t.Error("Created TCP response rule not equal to given TCP response rule")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	// TestEditTCPResponseRule
	r = &models.TCPResponseRule{
		Type:     "content",
		Action:   "accept",
		Cond:     "if",
		CondTest: "FALSE",
	}

	err = clientTest.EditTCPResponseRule(2, configuration.BackendParentName, "test", r, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, ondiskR, err = clientTest.GetTCPResponseRule(2, configuration.BackendParentName, "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(ondiskR, r) {
		fmt.Printf("Edited TCP response rule: %v\n", ondiskR)
		fmt.Printf("Given TCP response rule: %v\n", r)
		t.Error("Edited TCP response rule not equal to given TCP response rule")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	// TestDeleteTCPResponse
	_, rules, err := clientTest.GetTCPResponseRules(configuration.BackendParentName, "test", "")
	if err != nil {
		t.Error(err.Error())
	}
	N := int64(len(rules)) - 1
	err = clientTest.DeleteTCPResponseRule(N, configuration.BackendParentName, "test", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ := clientTest.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	_, _, err = clientTest.GetTCPResponseRule(N, configuration.BackendParentName, "test", "")
	if err == nil {
		t.Errorf("DeleteTCPResponseRule failed, TCP Response Rule %d still exists", N)
	}

	err = clientTest.DeleteTCPResponseRule(18, configuration.BackendParentName, "test_2", "", version)
	if err == nil {
		t.Error("Should throw error, non existent TCP Response Rule")
		version++
	}
}
