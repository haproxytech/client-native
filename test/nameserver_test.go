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

func nameServerExpectation() map[string]models.Nameservers {
	initStructuredExpected()
	res := StructuredToNameserverMap()
	// Add individual entries
	for k, vs := range res {
		for _, v := range vs {
			key := fmt.Sprintf("%s/%s", k, v.Name)
			res[key] = models.Nameservers{v}
		}
	}
	return res
}

func TestGetNameservers(t *testing.T) {
	m := make(map[string]models.Nameservers)
	v, nameservers, err := clientTest.GetNameservers("test", "")
	if err != nil {
		t.Error(err.Error())
	}
	m["resolvers/test"] = nameservers

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}
	checkNameservers(t, m)
}

func TestGetNameserver(t *testing.T) {
	m := make(map[string]models.Nameservers)

	v, l, err := clientTest.GetNameserver("dns1", "test", "")
	if err != nil {
		t.Error(err.Error())
	}
	m["resolvers/test/dns1"] = models.Nameservers{l}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}
	checkNameservers(t, m)

	_, _, err = clientTest.GetNameserver("community", "test", "")
	if err == nil {
		t.Error("Should throw error, non existent nameserver")
	}
}

func checkNameservers(t *testing.T, got map[string]models.Nameservers) {
	exp := nameServerExpectation()
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

func TestCreateEditDeleteNameserver(t *testing.T) {
	_, _, err := clientTest.GetNameserver("dns1", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	address := "192.168.1.2"
	port := int64(1024)
	e := &models.Nameserver{
		Address: &address,
		Port:    &port,
		Name:    "hapcommunity",
	}

	err = clientTest.CreateNameserver("test", e, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, l, err := clientTest.GetNameserver("hapcommunity", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(e, l) {
		fmt.Printf("Created nameserver: %v\n", e)
		fmt.Printf("Given nameserver: %v\n", l)
		t.Error("Created nameserver not equal to given nameserver")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	err = clientTest.CreateNameserver("test", e, "", version)
	if err == nil {
		t.Error("Should throw error nameserver already exists")
		version++
	}

	editPort := int64(1025)
	e = &models.Nameserver{
		Address: &address,
		Port:    &editPort,
		Name:    "hapcommunity",
	}

	err = clientTest.EditNameserver("hapcommunity", "test", e, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, l, err = clientTest.GetNameserver("hapcommunity", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(e, l) {
		fmt.Printf("Edited nameserver: %v\n", e)
		fmt.Printf("Given nameserver: %v\n", l)
		t.Error("Edited nameserver not equal to given nameserver")
	}

	// Delete
	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	err = clientTest.DeleteNameserver("hapcommunity", "test", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ = clientTest.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	_, _, err = clientTest.GetNameserver("hapcommunity", "test", "")
	if err == nil {
		t.Error("DeleteNameserver failed, nameserver still exists")
	}

	err = clientTest.DeleteNameserver("hapcommunity", "test", "", version)
	if err == nil {
		t.Error("Should throw error, non existent nameserver")
		version++
	}
}
