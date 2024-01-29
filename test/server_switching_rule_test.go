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

func serverSwitchingRuleExpectation() map[string]models.ServerSwitchingRules {
	initStructuredExpected()
	res := StructuredToServerSwitchingRuleMap()
	// Add individual entries
	for k, vs := range res {
		for _, v := range vs {
			key := fmt.Sprintf("%s/%d", k, *v.Index)
			res[key] = models.ServerSwitchingRules{v}
		}
	}
	return res
}

func TestGetServerSwitchingRules(t *testing.T) {
	mssr := make(map[string]models.ServerSwitchingRules)
	v, srvRules, err := clientTest.GetServerSwitchingRules("test", "")
	if err != nil {
		t.Error(err.Error())
	}
	mssr["backend/test"] = srvRules

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	_, srvRules, err = clientTest.GetServerSwitchingRules("test_2", "")
	if err != nil {
		t.Error(err.Error())
	}
	mssr["backend/test_2"] = srvRules

	checkServerSwitchingRules(t, mssr)
}

func checkServerSwitchingRules(t *testing.T, got map[string]models.ServerSwitchingRules) {
	exp := serverSwitchingRuleExpectation()
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

func TestGetServerSwitchingRule(t *testing.T) {
	m := make(map[string]models.ServerSwitchingRules)
	v, sr, err := clientTest.GetServerSwitchingRule(0, "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}
	m["backend/test/0"] = models.ServerSwitchingRules{sr}

	checkServerSwitchingRules(t, m)

	_, err = sr.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	_, _, err = clientTest.GetServerSwitchingRule(3, "test2", "")
	if err == nil {
		t.Error("Should throw error, non existent server switching rule")
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

	err := clientTest.CreateServerSwitchingRule("test", sr, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, srvRule, err := clientTest.GetServerSwitchingRule(2, "test", "")
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

	err = clientTest.EditServerSwitchingRule(2, "test", sr, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, srvRule, err = clientTest.GetServerSwitchingRule(2, "test", "")
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
	err = clientTest.DeleteServerSwitchingRule(2, "test", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ := clientTest.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	_, _, err = clientTest.GetServerSwitchingRule(2, "test", "")
	if err == nil {
		t.Error("DeleteServerSwitchingRule failed, server switching rule 3 still exists")
	}

	err = clientTest.DeleteServerSwitchingRule(2, "test_2", "", version)
	if err == nil {
		t.Error("Should throw error, non existent server switching rule")
		version++
	}
}
