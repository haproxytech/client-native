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
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/haproxytech/client-native/v6/configuration"
	"github.com/haproxytech/client-native/v6/models"
	"github.com/stretchr/testify/require"
)

func ringExpectation() map[string]models.Rings {
	initStructuredExpected()
	res := StructuredToRingMap()
	// Add individual entries
	for _, vs := range res {
		for _, v := range vs {
			key := v.Name
			res[key] = models.Rings{v}
		}
	}
	return res
}

func TestGetRings(t *testing.T) {
	m := make(map[string]models.Rings)
	v, rings, err := clientTest.GetRings("")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("version %v returned, expected %v", v, version)
	}
	m[""] = rings
	checkRings(t, m)
}

func TestGetRing(t *testing.T) {
	m := make(map[string]models.Rings)

	v, r, err := clientTest.GetRing("myring", "")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("version %v returned, expected %v", v, version)
	}
	m["myring"] = models.Rings{r}
	checkRings(t, m)

	_, _, err = clientTest.GetRing("doesnotexist", "")
	if err == nil {
		t.Error("should throw error, non existent rings section")
	}
}

func checkRings(t *testing.T, got map[string]models.Rings) {
	exp := ringExpectation()

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

func TestCreateEditDeleteRing(t *testing.T) {
	maxlen := int64(1300)
	size := int64(32765)
	timeoutConnect := int64(5)
	timeoutServer := int64(10)

	r := &models.Ring{
		Name:           "created_ring",
		Description:    "My local buffer 2",
		Format:         "rfc3164",
		Maxlen:         &maxlen,
		Size:           &size,
		TimeoutConnect: &timeoutConnect,
		TimeoutServer:  &timeoutServer,
	}
	err := clientTest.CreateRing(r, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, ring, err := clientTest.GetRing("created_ring", "")
	if err != nil {
		t.Error(err.Error())
	}

	if ring.Description != r.Description {
		t.Errorf("description expected %s got %s", ring.Description, r.Description)
	}

	if ring.Format != r.Format {
		t.Errorf("format expected %s got %s", ring.Format, r.Format)
	}

	if *ring.Maxlen != *r.Maxlen {
		t.Errorf("maxlen expected %v got %v", *ring.Maxlen, *r.Maxlen)
	}

	if *ring.Size != *r.Size {
		t.Errorf("size expected %v got %v", *ring.Size, *r.Size)
	}

	if *ring.TimeoutConnect != *r.TimeoutConnect {
		t.Errorf("timeout connect expected %v got %v", *ring.TimeoutConnect, *r.TimeoutConnect)
	}

	if *ring.TimeoutServer != *r.TimeoutServer {
		t.Errorf("timeout server expected %v got %v", *ring.TimeoutServer, *r.TimeoutServer)
	}

	if v != version {
		t.Errorf("version expected %v got %v", v, version)
	}

	err = clientTest.CreateRing(r, "", version)
	if err == nil {
		t.Error("should throw error ring already exists")
		version++
	}

	err = clientTest.DeleteRing("created_ring", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ := clientTest.GetVersion(""); v != version {
		t.Error("version not incremented")
	}

	err = clientTest.DeleteRing("created_ring", "", 999999)
	if err != nil {
		if confErr, ok := err.(*configuration.ConfError); ok {
			if !confErr.Is(configuration.ErrVersionMismatch) {
				t.Error("should throw configuration.ErrVersionMismatch error")
			}
		} else {
			t.Error("should throw configuration.ErrVersionMismatch error")
		}
	}
	_, _, err = clientTest.GetRing("created_ring", "")
	if err == nil {
		t.Error("deleteRing failed, ring created_ring still exists")
	}

	err = clientTest.DeleteRing("doesnotexist", "", version)
	if err == nil {
		t.Error("should throw error, non existent ring")
		version++
	}
}
