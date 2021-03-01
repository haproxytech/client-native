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

func TestGetBackends(t *testing.T) { //nolint:gocognit,gocyclo
	v, backends, err := client.GetBackends("")
	if err != nil {
		t.Error(err.Error())
	}

	if len(backends) != 2 {
		t.Errorf("%v backends returned, expected 2", len(backends))
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	for _, b := range backends {
		if b.Name != "test" && b.Name != "test_2" {
			t.Errorf("Expected only test or test_2 backend, %v found", b.Name)
		}
		if b.BindProcess != "all" {
			t.Errorf("%v: BindProcess not all: %v", b.Name, b.BindProcess)
		}
		if b.AdvCheck != "httpchk" {
			t.Errorf("%v: AdvCheck.Method not HEAD: %v", b.Name, b.AdvCheck)
		}
		if b.HttpchkParams.Method != "HEAD" {
			t.Errorf("%v: HttpchkParams.Method not HEAD: %v", b.Name, b.HttpchkParams.Method)
		}
		if b.HttpchkParams.URI != "/" {
			t.Errorf("%v: HttpchkParams.URI not HEAD: %v", b.Name, b.HttpchkParams.URI)
		}
		if b.Mode != "http" {
			t.Errorf("%v: Mode not http: %v", b.Name, b.Mode)
		}
		if *b.Balance.Algorithm != "roundrobin" {
			t.Errorf("%v: Balance.Algorithm not roundrobin: %v", b.Name, b.Balance.Algorithm)
		}
		if b.HashType.Method != "consistent" {
			t.Errorf("%v: HashType.Method not consistent: %v", b.Name, b.HashType.Method)
		}
		if b.HashType.Function != "sdbm" {
			t.Errorf("%v: HashType.Function not sdbm: %v", b.Name, b.HashType.Function)
		}
		if b.HashType.Modifier != "avalanche" {
			t.Errorf("%v: HashType.Modifier not avalanche: %v", b.Name, b.HashType.Modifier)
		}
		if b.HTTPConnectionMode != "http-keep-alive" {
			t.Errorf("%v: HTTPConnectionMode not http-keep-alive: %v", b.Name, b.HTTPConnectionMode)
		}
		if *b.Forwardfor.Enabled != "enabled" {
			t.Errorf("%v: Forwardfor not enabled: %v", b.Name, b.Forwardfor)
		}
		if *b.DefaultServer.Fall != 2000 {
			t.Errorf("%v: DefaultServer.Fall not 2000: %v", b.Name, *b.DefaultServer.Fall)
		}
		if *b.DefaultServer.Rise != 4000 {
			t.Errorf("%v: DefaultServer.Rise not 4000: %v", b.Name, *b.DefaultServer.Rise)
		}
		if *b.DefaultServer.Inter != 5000 {
			t.Errorf("%v: DefaultServer.Inter not 5000: %v", b.Name, *b.DefaultServer.Inter)
		}
		if *b.DefaultServer.Port != 8888 {
			t.Errorf("%v: DefaultServer.Port not 8888: %v", b.Name, *b.DefaultServer.Port)
		}
		if *b.Cookie.Name != "BLA" {
			t.Errorf("%v: HTTPCookie Name not BLA: %v", b.Name, b.Cookie)
		}
		if b.Cookie.Type != "rewrite" {
			t.Errorf("%v: HTTPCookie Type not rewrite %v", b.Name, b.Cookie.Type)
		}
		if !b.Cookie.Httponly {
			t.Errorf("%v: HTTPCookie Httponly not true %v", b.Name, b.Cookie.Httponly)
		}
		if !b.Cookie.Nocache {
			t.Errorf("%v: HTTPCookie Nocache not false %v", b.Name, b.Cookie.Nocache)
		}
		if *b.CheckTimeout != 2000 {
			t.Errorf("%v: CheckTimeout not 2000: %v", b.Name, *b.CheckTimeout)
		}
		if *b.ServerTimeout != 3000 {
			t.Errorf("%v: ServerTimeout not 3000: %v", b.Name, *b.ServerTimeout)
		}
	}
}

func TestGetBackend(t *testing.T) {
	v, b, err := client.GetBackend("test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	if b.Name != "test" {
		t.Errorf("Expected only test, %v found", b.Name)
	}
	if b.BindProcess != "all" {
		t.Errorf("%v: BindProcess not all: %v", b.Name, b.BindProcess)
	}
	if b.AdvCheck != "httpchk" {
		t.Errorf("%v: AdvCheck.Method not HEAD: %v", b.Name, b.AdvCheck)
	}
	if b.HttpchkParams.Method != "HEAD" {
		t.Errorf("%v: HttpchkParams.Method not HEAD: %v", b.Name, b.HttpchkParams.Method)
	}
	if b.HttpchkParams.URI != "/" {
		t.Errorf("%v: HttpchkParams.URI not HEAD: %v", b.Name, b.HttpchkParams.URI)
	}
	if b.Mode != "http" {
		t.Errorf("%v: Mode not http: %v", b.Name, b.Mode)
	}
	if *b.Balance.Algorithm != "roundrobin" {
		t.Errorf("%v: Balance.Algorithm not roundrobin: %v", b.Name, b.Balance.Algorithm)
	}
	if b.HashType.Method != "consistent" {
		t.Errorf("%v: HashType.Method not consistent: %v", b.Name, b.HashType.Method)
	}
	if b.HashType.Function != "sdbm" {
		t.Errorf("%v: HashType.Function not sdbm: %v", b.Name, b.HashType.Function)
	}
	if b.HashType.Modifier != "avalanche" {
		t.Errorf("%v: HashType.Modifier not avalanche: %v", b.Name, b.HashType.Modifier)
	}
	if b.HTTPConnectionMode != "http-keep-alive" {
		t.Errorf("%v: HTTPConnectionMode not http-keep-alive: %v", b.Name, b.HTTPConnectionMode)
	}
	if *b.Forwardfor.Enabled != "enabled" {
		t.Errorf("%v: Forwardfor not enabled: %v", b.Name, b.Forwardfor)
	}
	if *b.DefaultServer.Fall != 2000 {
		t.Errorf("%v: DefaultServer.Fall not 2000: %v", b.Name, *b.DefaultServer)
	}
	if *b.DefaultServer.Rise != 4000 {
		t.Errorf("%v: DefaultServer.Rise not 4000: %v", b.Name, *b.DefaultServer.Rise)
	}
	if *b.DefaultServer.Inter != 5000 {
		t.Errorf("%v: DefaultServer.Inter not 5000: %v", b.Name, *b.DefaultServer.Inter)
	}
	if *b.DefaultServer.Port != 8888 {
		t.Errorf("%v: DefaultServer.Port not 8888: %v", b.Name, *b.DefaultServer.Port)
	}
	if *b.Cookie.Name != "BLA" {
		t.Errorf("%v: HTTPCookie Name not BLA: %v", b.Name, b.Cookie)
	}
	if b.Cookie.Type != "rewrite" {
		t.Errorf("%v: HTTPCookie Type not rewrite %v", b.Name, b.Cookie.Type)
	}
	if !b.Cookie.Httponly {
		t.Errorf("%v: HTTPCookie Httponly not true %v", b.Name, b.Cookie.Httponly)
	}
	if !b.Cookie.Nocache {
		t.Errorf("%v: HTTPCookie Nocache not false %v", b.Name, b.Cookie.Nocache)
	}
	if *b.CheckTimeout != 2000 {
		t.Errorf("%v: CheckTimeout not 2000: %v", b.Name, *b.CheckTimeout)
	}
	if *b.ServerTimeout != 3000 {
		t.Errorf("%v: ServerTimeout not 3000: %v", b.Name, *b.ServerTimeout)
	}

	_, err = b.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	_, _, err = client.GetBackend("doesnotexist", "")
	if err == nil {
		t.Error("Should throw error, non existent bck")
	}
}

func TestCreateEditDeleteBackend(t *testing.T) {
	// TestCreateBackend
	tOut := int64(5)
	cookieName := "BLA"
	balanceAlgorithm := "uri"
	b := &models.Backend{
		Name: "created",
		Mode: "http",
		Balance: &models.Balance{
			Algorithm: &balanceAlgorithm,
			URILen:    100,
			URIDepth:  250,
		},
		BindProcess: "4",
		Cookie: &models.Cookie{
			Domains: []*models.Domain{
				&models.Domain{Value: "dom1"},
				&models.Domain{Value: "dom2"},
			},
			Dynamic:  true,
			Httponly: true,
			Indirect: true,
			Maxidle:  5,
			Maxlife:  20,
			Name:     &cookieName,
			Nocache:  true,
			Postonly: true,
			Preserve: false,
			Secure:   false,
			Type:     "prefix",
		},
		HashType: &models.BackendHashType{
			Method:   "map-based",
			Function: "crc32",
		},
		DefaultServer: &models.DefaultServer{
			Fall:  &tOut,
			Inter: &tOut,
		},
		HTTPConnectionMode:   "http-keep-alive",
		ConnectTimeout:       &tOut,
		ExternalCheck:        "enabled",
		ExternalCheckCommand: "/bin/false",
		ExternalCheckPath:    "/bin",
		Allbackups:           "enabled",
		AdvCheck:             "smtpchk",
		SmtpchkParams: &models.SmtpchkParams{
			Hello:  "HELO",
			Domain: "example.com",
		},
	}

	err := client.CreateBackend(b, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, backend, err := client.GetBackend("created", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !compareBackends(backend, b, t) {
		fmt.Printf("Created bck: %v\n", backend)
		fmt.Printf("Given bck: %v\n", b)
		t.Error("Created backend not equal to given backend")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	err = client.CreateBackend(b, "", version)
	if err == nil {
		t.Error("Should throw error bck already exists")
		version++
	}
	// TestEditBackend
	tOut = int64(3)
	e := int64(1200000)
	kl := int64(128)
	s := int64(25600)
	backends := []*models.Backend{
		{
			Name: "created",
			Mode: "http",
			Balance: &models.Balance{
				Algorithm: &balanceAlgorithm,
				URILen:    10,
				URIDepth:  25,
			},
			BindProcess: "3",
			Cookie: &models.Cookie{
				Domains: []*models.Domain{
					{Value: "dom1"},
					{Value: "dom3"},
				},
				Dynamic:  true,
				Httponly: true,
				Indirect: false,
				Maxidle:  150,
				Maxlife:  100,
				Name:     &cookieName,
				Nocache:  false,
				Postonly: false,
				Preserve: true,
				Secure:   true,
				Type:     "rewrite",
			},
			HTTPConnectionMode: "httpclose",
			ConnectTimeout:     &tOut,
			StickTable: &models.BackendStickTable{
				Expire: &e,
				Keylen: &kl,
				Size:   &s,
				Store:  "gpc0,http_req_rate(40s)",
				Type:   "string",
			},
			AdvCheck: "mysql-check",
			MysqlCheckParams: &models.MysqlCheckParams{
				Username:      "user",
				ClientVersion: "pre-41",
			},
		},
		{
			Name: "created",
			Mode: "http",
			Balance: &models.Balance{
				Algorithm: &balanceAlgorithm,
			},
			Cookie: &models.Cookie{
				Domains: []*models.Domain{
					{Value: "dom1"},
					{Value: "dom2"},
				},
				Name: &cookieName,
			},
			ConnectTimeout: &tOut,
			StickTable:     &models.BackendStickTable{},
			AdvCheck:       "pgsql-check",
			PgsqlCheckParams: &models.PgsqlCheckParams{
				Username: "user",
			},
		},
		{
			Name: "created",
			Mode: "http",
			Balance: &models.Balance{
				Algorithm: &balanceAlgorithm,
			},
			Cookie: &models.Cookie{
				Domains: []*models.Domain{
					{Value: "dom4"},
					{Value: "dom5"},
				},
				Name: &cookieName,
			},
			ConnectTimeout: &tOut,
			StickTable:     &models.BackendStickTable{},
			AdvCheck:       "httpchk",
			HttpchkParams: &models.HttpchkParams{
				Method: "HEAD",
				URI:    "/",
			},
		},
	}

	for i, backend := range backends {
		if errB := testBackendUpdate(backend, t); errB != nil {
			t.Errorf("failed update for backend %d: %v", i, err)
		}
	}

	// TestDeleteBackend
	err = client.DeleteBackend("created", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ := client.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	err = client.DeleteBackend("created", "", 999999999)
	if err != nil {
		if confErr, ok := err.(*ConfError); ok {
			if confErr.Code() != ErrVersionMismatch {
				t.Error("Should throw ErrVersionMismatch error")
			}
		} else {
			t.Error("Should throw ErrVersionMismatch error")
		}
	}

	_, _, err = client.GetBackend("created", "")
	if err == nil {
		t.Error("DeleteBackend failed, bck test still exists")
	}

	err = client.DeleteBackend("doesnotexist", "", version)
	if err == nil {
		t.Error("Should throw error, non existent bck")
		version++
	}
}

func testBackendUpdate(b *models.Backend, t *testing.T) error {
	err := client.EditBackend("created", b, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, backend, err := client.GetBackend("created", "")
	if err != nil {
		return err
	}

	if !compareBackends(backend, b, t) {
		fmt.Printf("Edited bck: %v\n", backend)
		fmt.Printf("Given bck: %v\n", b)
		return fmt.Errorf("edited backend not equal to given backend")
	}

	if v != version {
		return fmt.Errorf("version %v returned, expected %v", v, version)
	}
	return nil
}

func compareBackends(x, y *models.Backend, t *testing.T) bool { //nolint:gocognit,gocyclo
	if *x.Balance.Algorithm != *y.Balance.Algorithm {
		return false
	}
	if x.Balance.HdrName != y.Balance.HdrName {
		return false
	}
	if x.Balance.HdrUseDomainOnly != y.Balance.HdrUseDomainOnly {
		return false
	}
	if x.Balance.RandomDraws != y.Balance.RandomDraws {
		return false
	}
	if x.Balance.RdpCookieName != y.Balance.RdpCookieName {
		return false
	}
	if x.Balance.URIDepth != y.Balance.URIDepth {
		return false
	}
	if x.Balance.URILen != y.Balance.URILen {
		return false
	}
	if x.Balance.URIWhole != y.Balance.URIWhole {
		return false
	}
	if x.Balance.URLParam != y.Balance.URLParam {
		return false
	}
	if x.Balance.URLParamCheckPost != y.Balance.URLParamCheckPost {
		return false
	}
	if x.Balance.URLParamMaxWait != y.Balance.URLParamMaxWait {
		return false
	}

	x.Balance = nil
	y.Balance = nil

	if *x.Cookie.Name != *y.Cookie.Name {
		return false
	}
	if len(x.Cookie.Domains) != len(y.Cookie.Domains) {
		return false
	}
	if x.Cookie.Domains[0].Value != y.Cookie.Domains[0].Value {
		return false
	}
	if x.Cookie.Dynamic != y.Cookie.Dynamic {
		return false
	}
	if x.Cookie.Httponly != y.Cookie.Httponly {
		return false
	}
	if x.Cookie.Indirect != y.Cookie.Indirect {
		return false
	}
	if x.Cookie.Maxidle != y.Cookie.Maxidle {
		return false
	}
	if x.Cookie.Maxlife != y.Cookie.Maxlife {
		return false
	}
	if x.Cookie.Nocache != y.Cookie.Nocache {
		return false
	}
	if x.Cookie.Postonly != y.Cookie.Postonly {
		return false
	}
	if x.Cookie.Preserve != y.Cookie.Preserve {
		return false
	}
	if x.Cookie.Secure != y.Cookie.Secure {
		return false
	}
	if x.Cookie.Type != y.Cookie.Type {
		return false
	}

	x.Cookie = nil
	y.Cookie = nil

	if x.BindProcess != y.BindProcess {
		return false
	}

	if !reflect.DeepEqual(x.DefaultServer, y.DefaultServer) {
		return false
	}

	x.DefaultServer = nil
	y.DefaultServer = nil

	if !reflect.DeepEqual(x.HttpchkParams, y.HttpchkParams) {
		return false
	}

	x.HttpchkParams = nil
	y.HttpchkParams = nil

	if !reflect.DeepEqual(x.StickTable, y.StickTable) {
		return false
	}

	x.StickTable = nil
	y.StickTable = nil

	if !reflect.DeepEqual(x.Redispatch, y.Redispatch) {
		return false
	}

	x.Redispatch = nil
	y.Redispatch = nil

	if !reflect.DeepEqual(x.Forwardfor, y.Forwardfor) {
		return false
	}

	x.Forwardfor = nil
	y.Forwardfor = nil

	if !reflect.DeepEqual(x.SmtpchkParams, y.SmtpchkParams) {
		return false
	}

	x.SmtpchkParams = nil
	y.SmtpchkParams = nil

	if !reflect.DeepEqual(x.MysqlCheckParams, y.MysqlCheckParams) {
		return false
	}

	x.MysqlCheckParams = nil
	y.MysqlCheckParams = nil

	if !reflect.DeepEqual(x.PgsqlCheckParams, y.PgsqlCheckParams) {
		return false
	}

	x.PgsqlCheckParams = nil
	y.PgsqlCheckParams = nil

	return reflect.DeepEqual(x, y)
}
