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

	"github.com/haproxytech/client-native/v4/models"
)

func TestGetBinds(t *testing.T) {
	v, binds, err := clientTest.GetBinds("frontend", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if len(binds) != 4 {
		t.Errorf("%v binds returned, expected 4", len(binds))
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	for _, l := range binds {
		if l.Name != "webserv" && l.Name != "webserv2" && l.Name != "webserv3" && l.Name != "ipv6" {
			t.Errorf("Expected only webserv,webserv2, or ipv6 binds, %v found", l.Name)
		}
		if l.Address != "192.168.1.1" && l.Address != "192.168.1.2" && l.Address != "2a01:c9c0:a3:8::3" {
			t.Errorf("%v: Address not 192.168.1.1 or 192.168.1.2 or 2a01:c9c0:a3:8::3: %v", l.Name, l.Address)
		}
		if *l.Port != 80 && *l.Port != 8080 {
			t.Errorf("%v: Port not 80 or 8080: %v", l.Name, *l.Port)
		}
		if l.Thread != "all" && l.Thread != "1/all" && l.Thread != "1/1" && l.Thread != "1/1-1" {
			t.Errorf("%v: Thread not all or 1/all or 1/1-1: %v", l.Name, l.Thread)
		}
	}

	_, binds, err = clientTest.GetBinds("frontend", "test_2", "")
	if err != nil {
		t.Error(err.Error())
	}
	if len(binds) > 0 {
		t.Errorf("%v binds returned, expected 0", len(binds))
	}
}

func TestGetBind(t *testing.T) {
	v, l, err := clientTest.GetBind("webserv", "frontend", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	if l.Name != "webserv" {
		t.Errorf("Expected only webserv or webserv2 binds, %v found", l.Name)
	}
	if l.Address != "192.168.1.1" {
		t.Errorf("%v: Address not 192.168.1.1: %v", l.Name, l.Address)
	}
	if *l.Port != 80 {
		t.Errorf("%v: Port not 80 or 8080: %v", l.Name, *l.Port)
	}
	if l.OcspUpdate != "on" {
		t.Errorf("OcspUpdate %v returned, expected on", l.OcspUpdate)
	}
	if l.Sigalgs != "RSA+SHA256" {
		t.Errorf("Sigalgs %v returned, expected on", l.Sigalgs)
	}
	if l.ClientSigalgs != "ECDSA+SHA256:RSA+SHA256" {
		t.Errorf("ClientSigalgs %v returned, expected on", l.ClientSigalgs)
	}

	_, err = l.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	_, _, err = clientTest.GetBind("webserv", "frontend", "test_2", "")
	if err == nil {
		t.Error("Should throw error, non existent bind")
	}
}

func TestCreateEditDeleteBind(t *testing.T) {
	// TestCreateBind
	port := int64(4300)
	l := &models.Bind{
		Address: "192.168.2.1",
		Port:    &port,
		BindParams: models.BindParams{
			Name:           "created",
			Ssl:            true,
			NoAlpn:         true,
			SslCertificate: "dummy.crt",
			Interface:      "eth0",
			Verify:         "optional",
			SslMinVer:      "TLSv1.3",
			SslMaxVer:      "TLSv1.3",
			Ciphers:        "ECDH+AESGCM:ECDH+CHACHA20",
			Ciphersuites:   "TLS_AES_128_GCM_SHA256:TLS_AES_256_GCM_SHA384",
			CrlFile:        "dummy.crl",
			Thread:         "1/all",
			OcspUpdate:     "on",
			Sigalgs:        "ECDSA+SHA256:RSA+SHA256",
			ClientSigalgs:  "ECDSA+SHA256",
		},
	}

	err := clientTest.CreateBind("frontend", "test", l, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, bind, err := clientTest.GetBind("created", "frontend", "test", "")
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

	err = clientTest.CreateBind("frontend", "test", l, "", version)
	if err == nil {
		t.Error("Should throw error bind already exists")
		version++
	}

	// TestEditBind
	port = int64(5300)
	tOut := int64(5)
	l = &models.Bind{
		Address: "192.168.3.1",
		Port:    &port,
		BindParams: models.BindParams{
			Name:           "created",
			NoAlpn:         false,
			Transparent:    true,
			TCPUserTimeout: &tOut,
			SslMinVer:      "TLSv1.2",
			SslMaxVer:      "TLSv1.3",
			Interface:      "eth1",
			Thread:         "odd",
			OcspUpdate:     "off",
			Sigalgs:        "ECDSA+SHA256",
			ClientSigalgs:  "ECDSA+SHA256:RSA+SHA256",
		},
	}

	err = clientTest.EditBind("created", "frontend", "test", l, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, bind, err = clientTest.GetBind("created", "frontend", "test", "")
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
	err = clientTest.DeleteBind("created", "frontend", "test", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ := clientTest.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	_, _, err = clientTest.GetBind("created", "frontend", "test", "")
	if err == nil {
		t.Error("DeleteBind failed, bind test still exists")
	}

	err = clientTest.DeleteBind("created", "frontend", "test2", "", version)
	if err == nil {
		t.Error("Should throw error, non existent bind")
		version++
	}
}
