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
	"github.com/haproxytech/client-native/v6/configuration"
	"github.com/haproxytech/client-native/v6/misc"
	"github.com/haproxytech/client-native/v6/models"
	"github.com/stretchr/testify/require"
)

func bindExpectation() map[string]models.Binds {
	initStructuredExpected()
	res := StructuredToBindMap()
	// Add individual entries
	for k, vs := range res {
		for _, v := range vs {
			key := fmt.Sprintf("%s/%s", k, v.Name)
			res[key] = models.Binds{v}
		}
	}
	return res
}

func TestGetBinds(t *testing.T) {
	mbinds := map[string]models.Binds{}

	v, binds, err := clientTest.GetBinds(configuration.FrontendParentName, "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if len(binds) != 10 {
		t.Errorf("%v binds returned, expected 10", len(binds))
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}
	mbinds["frontend/test"] = binds

	_, binds, err = clientTest.GetBinds(configuration.FrontendParentName, "test_2", "")
	if err != nil {
		t.Error(err.Error())
	}

	mbinds["frontend/test_2"] = binds

	checkBinds(t, mbinds)
}

func checkBinds(t *testing.T, got map[string]models.Binds) {
	exp := bindExpectation()
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

func TestGetBind(t *testing.T) {
	m := make(map[string]models.Binds)
	v, l, err := clientTest.GetBind("webserv", configuration.FrontendParentName, "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}
	m["frontend/test/webserv"] = models.Binds{l}
	checkBinds(t, m)

	_, _, err = clientTest.GetBind("webserv", configuration.FrontendParentName, "test_2", "")
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
		Name:    "created",
		BindParams: models.BindParams{
			Ssl:            true,
			NoAlpn:         true,
			SslCertificate: "dummy.crt:dummy2.crt",
			Interface:      "eth0",
			Verify:         "optional",
			SslMinVer:      "TLSv1.3",
			SslMaxVer:      "TLSv1.3",
			Ciphers:        "ECDH+AESGCM:ECDH+CHACHA20",
			Ciphersuites:   "TLS_AES_128_GCM_SHA256:TLS_AES_256_GCM_SHA384",
			CrlFile:        "dummy.crl",
			Thread:         "1/all",
			Sigalgs:        "ECDSA+SHA256:RSA+SHA256",
			ClientSigalgs:  "ECDSA+SHA256",
			CaVerifyFile:   "ca.pem",
			Nice:           123,
			QuicSocket:     "listener",
			Nbconn:         12,
			DefaultCrtList: []string{"foobar2.pem.rsa", "foobar2.pem.ecdsa"},
			TCPMd5sig:      "secretpass",
			Ktls:           "off",
			TCPSs:          2,
		},
	}

	err := clientTest.CreateBind(configuration.FrontendParentName, "test", l, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, bind, err := clientTest.GetBind("created", configuration.FrontendParentName, "test", "")
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

	err = clientTest.CreateBind(configuration.FrontendParentName, "test", l, "", version)
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
		Name:    "created",
		BindParams: models.BindParams{
			NoAlpn:         false,
			Transparent:    true,
			TCPUserTimeout: &tOut,
			SslMinVer:      "TLSv1.2",
			SslMaxVer:      "TLSv1.3",
			Interface:      "eth1",
			Thread:         "odd",
			Sigalgs:        "ECDSA+SHA256",
			ClientSigalgs:  "ECDSA+SHA256:RSA+SHA256",
			IdlePing:       misc.Int64P(10000),
			SslCertificate: "dummy.crt:dummy2.crt:dummy3.crt",
		},
	}

	err = clientTest.EditBind("created", configuration.FrontendParentName, "test", l, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, bind, err = clientTest.GetBind("created", configuration.FrontendParentName, "test", "")
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
	err = clientTest.DeleteBind("created", configuration.FrontendParentName, "test", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ := clientTest.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	_, _, err = clientTest.GetBind("created", configuration.FrontendParentName, "test", "")
	if err == nil {
		t.Error("DeleteBind failed, bind test still exists")
	}

	err = clientTest.DeleteBind("created", configuration.FrontendParentName, "test2", "", version)
	if err == nil {
		t.Error("Should throw error, non existent bind")
		version++
	}
}
