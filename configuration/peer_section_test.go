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

	"github.com/haproxytech/client-native/v5/models"
)

func TestGetPeerSections(t *testing.T) {
	v, peerSections, err := clientTest.GetPeerSections("")
	if err != nil {
		t.Error(err.Error())
	}

	if len(peerSections) != 1 {
		t.Errorf("%v peerSections returned, expected 1", len(peerSections))
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	if peerSections[0].Name != "mycluster" {
		t.Errorf("Expected only mycluster, %v found", peerSections[0].Name)
	}
}

func TestGetPeerSection(t *testing.T) {
	v, l, err := clientTest.GetPeerSection("mycluster", "")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("%v: Version not %v: %v", l.Name, v, version)
	}

	if l.Name != "mycluster" {
		t.Errorf("%v: Name not: %v", l.Name, l.Name)
	}

	if !l.Enabled {
		t.Errorf("%v: Enabled not true", l.Name)
	}

	if l.DefaultServer.Fall == nil {
		t.Errorf("%v: DefaultServer.Fall is nil", l.Name)
	} else if *l.DefaultServer.Fall != 2000 {
		t.Errorf("%v: DefaultServer.Fall not 2000: %v", l.Name, *l.DefaultServer.Fall)
	}
	if l.DefaultServer.Rise == nil {
		t.Errorf("%v: DefaultServer.Rise is nil", l.Name)
	} else if *l.DefaultServer.Rise != 4000 {
		t.Errorf("%v: DefaultServer.Rise not 4000: %v", l.Name, *l.DefaultServer.Rise)
	}
	if l.DefaultServer.Inter == nil {
		t.Errorf("%v: DefaultServer.Inter is nil", l.Name)
	} else if *l.DefaultServer.Inter != 5000 {
		t.Errorf("%v: DefaultServer.Inter not 5000: %v", l.Name, *l.DefaultServer.Inter)
	}
	if l.DefaultServer.HealthCheckPort == nil {
		t.Errorf("%v: DefaultServer.HealthCheckPort is nil", l.Name)
	} else if *l.DefaultServer.HealthCheckPort != 8888 {
		t.Errorf("%v: DefaultServer.HealthCheckPort not 8888: %v", l.Name, *l.DefaultServer.HealthCheckPort)
	}
	if l.DefaultServer.Slowstart == nil {
		t.Errorf("%v: DefaultServer.Slowstart is nil", l.Name)
	} else if *l.DefaultServer.Slowstart != 6000 {
		t.Errorf("%v: DefaultServer.Slowstart not 6000: %v", l.Name, *l.DefaultServer.Slowstart)
	}

	if !l.DefaultBind.V4v6 {
		t.Errorf("%v: DefaultBind.V4v6 not true", l.Name)
	}
	if !l.DefaultBind.Ssl {
		t.Errorf("%v: DefaultBind.Ssl not true", l.Name)
	}
	if l.DefaultBind.SslCertificate != "/etc/haproxy/site.pem" {
		t.Errorf("%v: DefaultBind.SslCertificate not etc/haproxy/site.pem: %v", l.Name, l.DefaultBind.SslCertificate)
	}
	if l.DefaultBind.Alpn != "h2,http/1.1" {
		t.Errorf("%v: DefaultBind.Alpn not h2,http/1.1: %v", l.Name, l.DefaultBind.Alpn)
	}

	if l.StickTable.Type != "ip" {
		t.Errorf("%v: StickTable.Type not ip: %v", l.Name, l.StickTable.Type)
	}
	if l.StickTable.Size == nil {
		t.Errorf("%v: StickTable.Size is nil", l.Name)
	} else if *l.StickTable.Size != 204800 {
		t.Errorf("%v: StickTable.Size not 204800: %v", l.Name, *l.StickTable.Size)
	}
	if l.StickTable.Expire == nil {
		t.Errorf("%v: StickTable.Expire is nil", l.Name)
	} else if *l.StickTable.Expire != 3600000 {
		t.Errorf("%v: StickTable.Expire not 3600000: %v", l.Name, *l.StickTable.Expire)
	}
	if l.StickTable.Store != "http_req_rate(10s)" {
		t.Errorf("%v: StickTable.Store not http_req_rate(10s): %v", l.Name, l.StickTable.Store)
	}

	_, err = l.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	_, _, err = clientTest.GetPeerSection("doesnotexist", "")
	if err == nil {
		t.Error("Should throw error, non existant peer section")
	}
}

func TestCreateEditDeletePeerSection(t *testing.T) {
	tOut := int64(5)
	tSize := int64(20000)
	f := &models.PeerSection{
		Name:     "testcluster",
		Disabled: true,
		DefaultServer: &models.DefaultServer{
			ServerParams: models.ServerParams{
				Fall:  &tOut,
				Inter: &tOut,
			},
		},
		DefaultBind: &models.DefaultBind{
			BindParams: models.BindParams{
				Alpn:           "h2,http/1.1",
				Ssl:            true,
				SslCertificate: "/etc/haproxy/cluster.pem",
			},
		},
		StickTable: &models.ConfigStickTable{
			Size:   &tSize,
			Expire: &tSize,
			Type:   "integer",
		},
	}
	err := clientTest.CreatePeerSection(f, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, peerSection, err := clientTest.GetPeerSection("testcluster", "")
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

	err = clientTest.CreatePeerSection(f, "", version)
	if err == nil {
		t.Error("Should throw error peerSection already exists")
		version++
	}

	err = clientTest.DeletePeerSection("testcluster", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ := clientTest.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	err = clientTest.DeletePeerSection("testcluster", "", 999999)
	if err != nil {
		if confErr, ok := err.(*ConfError); ok {
			if confErr.Code() != ErrVersionMismatch {
				t.Error("Should throw ErrVersionMismatch error")
			}
		} else {
			t.Error("Should throw ErrVersionMismatch error")
		}
	}
	_, _, err = clientTest.GetPeerSection("testcluster", "")
	if err == nil {
		t.Error("DeletePeerSection failed, peerSection testcluster still exists")
	}

	err = clientTest.DeletePeerSection("doesnotexist", "", version)
	if err == nil {
		t.Error("Should throw error, non existant peerSection")
		version++
	}
}
