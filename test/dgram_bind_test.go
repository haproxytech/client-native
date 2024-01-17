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
	"github.com/haproxytech/client-native/v5/models"
	"github.com/stretchr/testify/require"
)

func dgramBindExpectation() map[string]models.DgramBinds {
	initStructuredExpected()
	res := StructuredToDgramBindMap()
	// Add individual entries
	for k, vs := range res {
		for _, v := range vs {
			key := fmt.Sprintf("%s/%s", k, v.Name)
			res[key] = models.DgramBinds{v}
		}
	}
	return res
}

func TestGetDgramBinds(t *testing.T) {
	m := make(map[string]models.DgramBinds)

	v, dBinds, err := clientTest.GetDgramBinds("sylog-loadb", "")
	if err != nil {
		t.Error(err.Error())
	}

	if len(dBinds) != 1 {
		t.Errorf("%v dgram-binds returned, expected 1", len(dBinds))
	}
	for _, e := range dBinds {
		m[fmt.Sprintf("log_forward/sylog-loadb/%s", e.Name)] = models.DgramBinds{e}
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}
	checkDgramBinds(t, m)

	// _, dBinds, err = clientTest.GetDgramBinds("test_2", "")
	// if err != nil {
	// 	t.Error(err.Error())
	// }
	// if len(dBinds) > 0 {
	// 	t.Errorf("%v dgram-binds returned, expected 0", len(dBinds))
	// }
}

func TestGetDgramBind(t *testing.T) {
	m := make(map[string]models.DgramBinds)
	v, l, err := clientTest.GetDgramBind("webserv", "sylog-loadb", "")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}
	m["log_forward/sylog-loadb/webserv"] = models.DgramBinds{l}

	checkDgramBinds(t, m)

	_, _, err = clientTest.GetDgramBind("webserv", "test_2", "")
	if err == nil {
		t.Error("Should throw error, non existent bind")
	}
}

func checkDgramBinds(t *testing.T, got map[string]models.DgramBinds) {
	exp := dgramBindExpectation()
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

func TestCreateEditDeleteDgramBind(t *testing.T) {
	// TestCreateBind
	port := int64(4300)
	l := &models.DgramBind{
		Address:   "192.168.2.1",
		Port:      &port,
		Name:      "created",
		Interface: "eth0",
	}

	err := clientTest.CreateDgramBind("sylog-loadb", l, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, bind, err := clientTest.GetDgramBind("created", "sylog-loadb", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(bind, l) {
		fmt.Printf("Created bind: %v\n", bind)
		fmt.Printf("Given bind: %v\n", l)
		t.Error("Created bind not equal to given bind")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	err = clientTest.CreateDgramBind("sylog-loadb", l, "", version)
	if err == nil {
		t.Error("Should throw error bind already exists")
		version++
	}

	// TestEditBind
	port = int64(5300)
	l = &models.DgramBind{
		Address:     "192.168.3.1",
		Port:        &port,
		Name:        "created",
		Transparent: true,
		Interface:   "eth1",
	}

	err = clientTest.EditDgramBind("created", "sylog-loadb", l, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, bind, err = clientTest.GetDgramBind("created", "sylog-loadb", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(bind, l) {
		fmt.Printf("Edited bind: %v\n", bind)
		fmt.Printf("Given bind: %v\n", l)
		t.Error("Edited bind not equal to given bind")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	// TestDeleteBind
	err = clientTest.DeleteDgramBind("created", "sylog-loadb", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ := clientTest.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	_, _, err = clientTest.GetDgramBind("created", "sylog-loadb", "")
	if err == nil {
		t.Error("DeleteDgramBind failed, bind test still exists")
	}

	err = clientTest.DeleteDgramBind("created", "test2", "", version)
	if err == nil {
		t.Error("Should throw error, non existent bind")
		version++
	}
}
