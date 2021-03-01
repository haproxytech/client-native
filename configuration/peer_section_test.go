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
	"fmt"
	"reflect"
	"testing"

	"github.com/haproxytech/client-native/v2/models"
)

func TestGetPeerSections(t *testing.T) {
	v, peerSections, err := client.GetPeerSections("")
	if err != nil {
		t.Error(err.Error())
	}

	if len(peerSections) != 1 {
		t.Errorf("%v peerSections returned, expected 2", len(peerSections))
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	if peerSections[0].Name != "mycluster" {
		t.Errorf("Expected only mycluster, %v found", peerSections[0].Name)
	}
}

func TestGetPeerSection(t *testing.T) {
	v, l, err := client.GetPeerSection("mycluster", "")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	if l.Name != "mycluster" {
		t.Errorf("Expected mycluster peerSection, %v found", l.Name)
	}

	_, err = l.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	_, _, err = client.GetPeerSection("doesnotexist", "")
	if err == nil {
		t.Error("Should throw error, non existant peer section")
	}
}

func TestCreateEditDeletePeerSection(t *testing.T) {
	f := &models.PeerSection{
		Name: "testcluster",
	}
	err := client.CreatePeerSection(f, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, peerSection, err := client.GetPeerSection("testcluster", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(peerSection, f) {
		fmt.Printf("Created peerSection: %v\n", peerSection)
		fmt.Printf("Given peerSection: %v\n", f)
		t.Error("Created peerSection not equal to given peerSection")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	err = client.CreatePeerSection(f, "", version)
	if err == nil {
		t.Error("Should throw error peerSection already exists")
		version++
	}

	err = client.DeletePeerSection("testcluster", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ := client.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	err = client.DeletePeerSection("testcluster", "", 999999)
	if err != nil {
		if confErr, ok := err.(*ConfError); ok {
			if confErr.Code() != ErrVersionMismatch {
				t.Error("Should throw ErrVersionMismatch error")
			}
		} else {
			t.Error("Should throw ErrVersionMismatch error")
		}
	}
	_, _, err = client.GetPeerSection("testcluster", "")
	if err == nil {
		t.Error("DeletePeerSection failed, peerSection testcluster still exists")
	}

	err = client.DeletePeerSection("doesnotexist", "", version)
	if err == nil {
		t.Error("Should throw error, non existant peerSection")
		version++
	}
}
