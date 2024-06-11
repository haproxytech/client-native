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

func logTargetExpectation() map[string]models.LogTargets {
	initStructuredExpected()
	res := StructuredToLogTargetMap()
	// Add individual entries
	for k, vs := range res {
		for i, v := range vs {
			key := fmt.Sprintf("%s/%d", k, i)
			res[key] = models.LogTargets{v}
		}
	}
	return res
}

func TestGetLogTargets(t *testing.T) {
	mlt := make(map[string]models.LogTargets)
	v, lTargets, err := clientTest.GetLogTargets(configuration.FrontendParentName, "test", "")
	if err != nil {
		t.Error(err.Error())
	}
	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}
	mlt["frontend/test"] = lTargets

	_, lTargets, err = clientTest.GetLogTargets(configuration.BackendParentName, "test_2", "")
	if err != nil {
		t.Error(err.Error())
	}
	mlt["backend/test_2"] = lTargets

	_, lTargets, err = clientTest.GetLogTargets(configuration.GlobalParentName, "", "")
	mlt[configuration.GlobalParentName] = lTargets

	checkLogTargets(t, mlt)
}

func checkLogTargets(t *testing.T, got map[string]models.LogTargets) {
	exp := logTargetExpectation()
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

func TestGetLogTarget(t *testing.T) {
	m := make(map[string]models.LogTargets)
	v, l, err := clientTest.GetLogTarget(2, configuration.FrontendParentName, "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}
	m["frontend/test/2"] = models.LogTargets{l}

	_, err = l.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	v, l, err = clientTest.GetLogTarget(2, configuration.LogForwardParentName, "sylog-loadb", "")
	if err != nil {
		t.Error(err.Error())
	}
	m["log_forward/sylog-loadb/2"] = models.LogTargets{l}

	_, err = l.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	_, _, err = clientTest.GetLogTarget(3, configuration.BackendParentName, "test_2", "")
	if err == nil {
		t.Error("Should throw error, non existent Log Target")
	}

	checkLogTargets(t, m)
}

func TestCreateEditDeleteLogTarget(t *testing.T) {
	id := int64(3)

	// TestCreateLogTarget
	r := &models.LogTarget{
		Address:  "stdout",
		Format:   "raw",
		Facility: "daemon",
		Level:    "notice",
	}

	err := clientTest.CreateLogTarget(id, configuration.FrontendParentName, "test", r, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, ondiskR, err := clientTest.GetLogTarget(id, configuration.FrontendParentName, "test", "")
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
		Address:  "stdout",
		Format:   "rfc3164",
		Facility: "local1",
	}

	err = clientTest.EditLogTarget(3, configuration.FrontendParentName, "test", r, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, ondiskR, err = clientTest.GetLogTarget(3, configuration.FrontendParentName, "test", "")
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
	err = clientTest.DeleteLogTarget(3, configuration.FrontendParentName, "test", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ := clientTest.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	_, _, err = clientTest.GetLogTarget(3, configuration.FrontendParentName, "test", "")
	if err == nil {
		t.Error("DeleteLogTarget failed, Log Target 3 still exists")
	}

	err = clientTest.DeleteLogTarget(2, configuration.BackendParentName, "test_2", "", version)
	if err == nil {
		t.Error("Should throw error, non existent Log Target")
		version++
	}
}
