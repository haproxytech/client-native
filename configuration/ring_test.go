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

package configuration

import (
	"testing"

	"github.com/haproxytech/client-native/v5/misc"
	"github.com/haproxytech/client-native/v5/models"
)

func TestGetRings(t *testing.T) {
	v, rings, err := clientTest.GetRings("")
	if err != nil {
		t.Error(err.Error())
	}

	if len(rings) != 1 {
		t.Errorf("%v rings returned, expected 1", len(rings))
	}

	if v != version {
		t.Errorf("version %v returned, expected %v", v, version)
	}

	if rings[0].Name != "myring" {
		t.Errorf("expected only test, %v found", rings[0].Name)
	}
}

func TestGetRing(t *testing.T) {
	v, r, err := clientTest.GetRing("myring", "")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("version %v returned, expected %v", v, version)
	}

	if r.Name != "myring" {
		t.Errorf("expected myring ring, %v found", r.Name)
	}

	if *r.Maxlen != *misc.Int64P(1200) {
		t.Errorf("expected maxlen 1200, %v found", r.Maxlen)
	}

	_, err = r.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	_, _, err = clientTest.GetRing("doesnotexist", "")
	if err == nil {
		t.Error("should throw error, non existent rings section")
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
		if confErr, ok := err.(*ConfError); ok {
			if !confErr.Is(ErrVersionMismatch) {
				t.Error("should throw ErrVersionMismatch error")
			}
		} else {
			t.Error("should throw ErrVersionMismatch error")
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
