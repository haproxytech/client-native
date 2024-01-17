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
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/haproxytech/client-native/v5/configuration"
	"github.com/haproxytech/client-native/v5/misc"
	"github.com/haproxytech/client-native/v5/models"
	"github.com/stretchr/testify/require"
)

func backendExpectation() map[string]models.Backends {
	initStructuredExpected()
	res := StructuredToBackendMap()
	// Add individual entries
	for _, vs := range res {
		for _, v := range vs {
			key := v.Name
			res[key] = models.Backends{v}
		}
	}
	return res
}

func TestGetBackends(t *testing.T) { //nolint:gocognit,gocyclo
	m := make(map[string]models.Backends)
	v, backends, err := clientTest.GetBackends("")
	if err != nil {
		t.Error(err.Error())
	}
	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}
	m[""] = backends

	checkBackends(t, m)
}

func checkBackends(t *testing.T, got map[string]models.Backends) {
	exp := backendExpectation()
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

func TestGetBackend(t *testing.T) {
	m := make(map[string]models.Backends)

	v, b, err := clientTest.GetBackend("test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}
	m["test"] = models.Backends{b}

	_, err = b.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	_, _, err = clientTest.GetBackend("doesnotexist", "")
	if err == nil {
		t.Error("Should throw error, non existent bck")
	}

	checkBackends(t, m)
}

func TestCreateEditDeleteBackend(t *testing.T) {
	// TestCreateBackend
	tOut := int64(5)
	cookieName := "BLA"
	balanceAlgorithm := "uri"
	srvtcpkaCnt := int64(10)
	srvtcpkaTimeout := int64(10000)
	statsRealm := "Haproxy Stats"
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
				{Value: "dom1"},
				{Value: "dom2"},
			},
			Attrs:    []*models.Attr{},
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
		HashType: &models.HashType{
			Method:   "map-based",
			Function: "crc32",
		},
		DefaultServer: &models.DefaultServer{
			ServerParams: models.ServerParams{
				Fall:       &tOut,
				Inter:      &tOut,
				LogBufsize: misc.Int64P(123),
			},
		},
		HTTPConnectionMode:   "http-keep-alive",
		HTTPKeepAlive:        "enabled",
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
		AcceptInvalidHTTPResponse: "enabled",
		Compression: &models.Compression{
			Offload: true,
		},
		LogHealthChecks:    "enabled",
		Checkcache:         "enabled",
		IndependentStreams: "enabled",
		Nolinger:           "enabled",
		Originalto: &models.Originalto{
			Enabled: misc.StringP("enabled"),
			Except:  "127.0.0.1",
			Header:  "X-Client-Dst",
		},
		Persist:          "enabled",
		PreferLastServer: "enabled",
		SpopCheck:        "enabled",
		TCPSmartConnect:  "enabled",
		Transparent:      "enabled",
		SpliceAuto:       "enabled",
		SpliceRequest:    "enabled",
		SpliceResponse:   "enabled",
		SrvtcpkaCnt:      &srvtcpkaCnt,
		SrvtcpkaIdle:     &srvtcpkaTimeout,
		SrvtcpkaIntvl:    &srvtcpkaTimeout,
		StatsOptions: &models.StatsOptions{
			StatsShowModules: true,
			StatsRealm:       true,
			StatsRealmRealm:  &statsRealm,
			StatsAuths: []*models.StatsAuth{
				{User: misc.StringP("user1"), Passwd: misc.StringP("pwd1")},
				{User: misc.StringP("user2"), Passwd: misc.StringP("pwd2")},
			},
			StatsHTTPRequests: []*models.StatsHTTPRequest{
				{Type: misc.StringP("allow"), Cond: "if", CondTest: "something"},
				{Type: misc.StringP("auth"), Realm: "haproxy\\ stats"},
			},
		},
		EmailAlert: &models.EmailAlert{
			From:    misc.StringP("prod01@example.com"),
			To:      misc.StringP("sre@example.com"),
			Level:   "warning",
			Mailers: misc.StringP("localmailer1"),
		},
		ErrorFilesFromHTTPErrors: []*models.Errorfiles{
			{Name: "test_errors", Codes: []int64{400}},
			{Name: "test_errors_all"},
		},
		Disabled: true,
	}

	err := clientTest.CreateBackend(b, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, backend, err := clientTest.GetBackend("created", "")
	if err != nil {
		t.Error(err.Error())
	}

	var givenJSONB []byte
	givenJSONB, err = b.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	var ondiskJSONB []byte
	ondiskJSONB, err = backend.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	if string(givenJSONB) != string(ondiskJSONB) {
		fmt.Printf("Created backend: %v\n", string(ondiskJSONB))
		fmt.Printf("Given backend: %v\n", string(givenJSONB))
		t.Error("Created backend not equal to given backend")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	err = clientTest.CreateBackend(b, "", version)
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
			From: "test_defaults",
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
			// Use non deprecated option only
			HTTPConnectionMode: "httpclose",
			ConnectTimeout:     &tOut,
			StickTable: &models.ConfigStickTable{
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
			EmailAlert: &models.EmailAlert{
				From:    misc.StringP("prod01@example.com"),
				To:      misc.StringP("sre@example.com"),
				Level:   "warning",
				Mailers: misc.StringP("localmailer1"),
			},
			Originalto: &models.Originalto{
				Enabled: misc.StringP("enabled"),
				Except:  "127.0.0.1",
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
			StickTable:     &models.ConfigStickTable{},
			AdvCheck:       "pgsql-check",
			PgsqlCheckParams: &models.PgsqlCheckParams{
				Username: "user",
			},
			EmailAlert: &models.EmailAlert{
				From:    misc.StringP("prod01@example.com"),
				To:      misc.StringP("sre@example.com"),
				Level:   "warning",
				Mailers: misc.StringP("localmailer1"),
			},
			// Use deprecated option only
			Httpclose: "enabled",
			Originalto: &models.Originalto{
				Enabled: misc.StringP("enabled"),
				Header:  "X-Client-Dst",
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
			StickTable:     &models.ConfigStickTable{},
			AdvCheck:       "httpchk",
			HttpchkParams: &models.HttpchkParams{
				Method: "HEAD",
				URI:    "/",
			},
			HTTPCheck: &models.HTTPCheck{
				Type:    "send",
				Method:  "OPTIONS",
				URI:     "/",
				Version: "HTTP/1.1",
				Index:   misc.Int64P(0),
			},
			Checkcache:         "disabled",
			IndependentStreams: "disabled",
			Nolinger:           "disabled",
			Persist:            "disabled",
			PreferLastServer:   "disabled",
			SpopCheck:          "disabled",
			TCPSmartConnect:    "disabled",
			Transparent:        "disabled",
			SpliceAuto:         "disabled",
			SpliceRequest:      "disabled",
			SpliceResponse:     "disabled",
			SrvtcpkaCnt:        &srvtcpkaCnt,
			SrvtcpkaIdle:       &srvtcpkaTimeout,
			SrvtcpkaIntvl:      &srvtcpkaTimeout,
			StatsOptions: &models.StatsOptions{
				StatsShowModules: true,
				StatsRealm:       true,
				StatsRealmRealm:  &statsRealm,
				StatsAuths: []*models.StatsAuth{
					{User: misc.StringP("new_user1"), Passwd: misc.StringP("new_pwd1")},
					{User: misc.StringP("new_user2"), Passwd: misc.StringP("new_pwd2")},
				},
				StatsHTTPRequests: []*models.StatsHTTPRequest{
					{Type: misc.StringP("allow"), Cond: "if", CondTest: "something_else"},
					{Type: misc.StringP("auth"), Realm: "haproxy\\ stats2"},
				},
			},
			EmailAlert: &models.EmailAlert{
				From:    misc.StringP("prod01@example.com"),
				To:      misc.StringP("sre@example.com"),
				Level:   "warning",
				Mailers: misc.StringP("localmailer1"),
			},
			Originalto: &models.Originalto{
				Enabled: misc.StringP("enabled"),
				Except:  "127.0.0.1",
				Header:  "X-Client-Dst",
			},
		},
	}

	for i, backend := range backends {
		if errB := testBackendUpdate(backend, t); errB != nil {
			t.Errorf("failed update for backend %d: %v", i, errB)
		}
	}

	// TestDeleteBackend
	err = clientTest.DeleteBackend("created", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ := clientTest.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	err = clientTest.DeleteBackend("created", "", 999999999)
	if err != nil {
		if confErr, ok := err.(*configuration.ConfError); ok {
			if !confErr.Is(configuration.ErrVersionMismatch) {
				t.Error("Should throw configuration.ErrVersionMismatch error")
			}
		} else {
			t.Error("Should throw configuration.ErrVersionMismatch error")
		}
	}

	_, _, err = clientTest.GetBackend("created", "")
	if err == nil {
		t.Error("DeleteBackend failed, bck test still exists")
	}

	err = clientTest.DeleteBackend("doesnotexist", "", version)
	if err == nil {
		t.Error("Should throw error, non existent bck")
		version++
	}
}

func testBackendUpdate(b *models.Backend, t *testing.T) error {
	err := clientTest.EditBackend("created", b, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, backend, err := clientTest.GetBackend("created", "")
	if err != nil {
		return err
	}

	if !compareBackends(backend, b, t) {
		fmt.Printf("Edited bck: %+v\n", backend)
		fmt.Printf("Given bck: %+v\n", b)
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

	if !cmp.Equal(x.HTTPCheck, y.HTTPCheck, cmpopts.EquateEmpty()) {
		t.Errorf("Diff in HTTPChecK %s", cmp.Diff(x.HTTPCheck, y.HTTPCheck, cmpopts.EquateEmpty()))
		return false
	}

	x.HTTPCheck = nil
	y.HTTPCheck = nil

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

	if !reflect.DeepEqual(x.Originalto, y.Originalto) {
		return false
	}

	x.Originalto = nil
	y.Originalto = nil

	// Due to deprecated fields Httpclose, HTTPKeepAlive, HTTPServerClose
	// in favor of HTTPConnectionMode
	// If HTTPConnectionMode is set in original backend
	// - Httpclose, HTTPKeepAlive, HTTPServerClose will be set in updated backend, even if not present in original backend
	// If HTTPConnectionMode is unset in original backend:
	// - it will be set in updated backend
	switch y.HTTPConnectionMode {
	case "http-keep-alive":
		if x.HTTPKeepAlive != "enabled" {
			return false
		}
		x.HTTPKeepAlive = ""
	case "http-server-close":
		if x.HTTPServerClose != "enabled" {
			return false
		}
		x.HTTPServerClose = ""
	case "httpclose":
		if x.Httpclose != "enabled" {
			return false
		}
		x.Httpclose = ""
	case "":
		x.HTTPConnectionMode = ""
	}

	return reflect.DeepEqual(x, y)
}

func TestCreateEditDeleteBackendHTTPConnectionMode(t *testing.T) {
	// TestCreateBackend
	tOut := int64(5)

	// Backend with HTTPConnectionMode only
	b := &models.Backend{
		Name: "special-httpconnectionmode",
		Mode: "http",
		DefaultServer: &models.DefaultServer{
			ServerParams: models.ServerParams{
				Fall:  &tOut,
				Inter: &tOut,
			},
		},
		HTTPConnectionMode: "http-keep-alive",
	}

	err := clientTest.CreateBackend(b, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, backend, err := clientTest.GetBackend("special-httpconnectionmode", "")
	if err != nil {
		t.Error(err.Error())
	}

	if backend.HTTPConnectionMode != "http-keep-alive" {
		t.Errorf("Created backend is not correct for HTTPConnectionMode: %s", backend.HTTPConnectionMode)
	}
	if backend.HTTPKeepAlive != "enabled" {
		t.Errorf("Created backend is not correct for HTTPKeepAlive: %s", backend.HTTPConnectionMode)
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	err = clientTest.CreateBackend(b, "", version)
	if err == nil {
		t.Error("Should throw error bck already exists")
		version++
	}

	type testinput struct {
		backend                    *models.Backend
		expectedHTTPConnectionMode string
		expectedHTTPKeepAlive      string
		expectedHttpclose          string
		exptectedHTTPServerClose   string
	}
	// TestEditBackend
	inputs := []testinput{
		{
			// Update HTTPConnectionMode
			backend: &models.Backend{
				Name: "special-httpconnectionmode",
				Mode: "http",
				DefaultServer: &models.DefaultServer{
					ServerParams: models.ServerParams{
						Fall:  &tOut,
						Inter: &tOut,
					},
				},
				HTTPConnectionMode: "httpclose",
			},
			expectedHTTPConnectionMode: "httpclose",
			expectedHTTPKeepAlive:      "",
			exptectedHTTPServerClose:   "",
			expectedHttpclose:          "enabled",
		},
		{
			// Use only deprecated option
			backend: &models.Backend{
				Name: "special-httpconnectionmode",
				Mode: "http",
				DefaultServer: &models.DefaultServer{
					ServerParams: models.ServerParams{
						Fall:  &tOut,
						Inter: &tOut,
					},
				},
				HTTPServerClose: "enabled",
			},
			expectedHTTPConnectionMode: "http-server-close",
			expectedHTTPKeepAlive:      "",
			exptectedHTTPServerClose:   "enabled",
			expectedHttpclose:          "",
		},
		{
			// Use both - Priority on HTTPConnection
			backend: &models.Backend{
				Name: "special-httpconnectionmode",
				Mode: "http",
				DefaultServer: &models.DefaultServer{
					ServerParams: models.ServerParams{
						Fall:  &tOut,
						Inter: &tOut,
					},
				},
				HTTPConnectionMode: "http-keep-alive",
				HTTPServerClose:    "enabled",
			},
			expectedHTTPConnectionMode: "http-keep-alive",
			expectedHTTPKeepAlive:      "enabled",
			exptectedHTTPServerClose:   "",
			expectedHttpclose:          "",
		},
		{
			// no option with deprecated option
			backend: &models.Backend{
				Name: "special-httpconnectionmode",
				Mode: "http",
				DefaultServer: &models.DefaultServer{
					ServerParams: models.ServerParams{
						Fall:  &tOut,
						Inter: &tOut,
					},
				},
				HTTPServerClose: "disabled",
			},
			expectedHTTPConnectionMode: "", // not possible to set no option with this field
			expectedHTTPKeepAlive:      "",
			exptectedHTTPServerClose:   "disabled",
			expectedHttpclose:          "",
		},
		{
			// set back with HTTPConnectionMode
			backend: &models.Backend{
				Name: "special-httpconnectionmode",
				Mode: "http",
				DefaultServer: &models.DefaultServer{
					ServerParams: models.ServerParams{
						Fall:  &tOut,
						Inter: &tOut,
					},
				},
				HTTPConnectionMode: "httpclose",
			},
			expectedHTTPConnectionMode: "httpclose",
			expectedHTTPKeepAlive:      "",
			exptectedHTTPServerClose:   "",
			expectedHttpclose:          "enabled",
		},
		{
			// remove option
			backend: &models.Backend{
				Name: "special-httpconnectionmode",
				Mode: "http",
				DefaultServer: &models.DefaultServer{
					ServerParams: models.ServerParams{
						Fall:  &tOut,
						Inter: &tOut,
					},
				},
				HTTPConnectionMode: "",
			},
			expectedHTTPConnectionMode: "",
			expectedHTTPKeepAlive:      "",
			exptectedHTTPServerClose:   "",
			expectedHttpclose:          "",
		},
	}

	for i, input := range inputs {
		err := clientTest.EditBackend("special-httpconnectionmode", input.backend, "", version)
		if err != nil {
			t.Error(err.Error())
		} else {
			version++
		}

		_, backend, err := clientTest.GetBackend("special-httpconnectionmode", "")
		if err != nil {
			t.Error(err.Error())
		}

		if backend.HTTPConnectionMode != input.expectedHTTPConnectionMode {
			t.Errorf("Updated backend %d is not correct for HTTPConnectionMode: %s", i, backend.HTTPConnectionMode)
		}
		if backend.HTTPKeepAlive != input.expectedHTTPKeepAlive {
			t.Errorf("Updated backend %d is not correct for HTTPKeepAlive: %s", i, backend.HTTPConnectionMode)
		}
		if backend.HTTPServerClose != input.exptectedHTTPServerClose {
			t.Errorf("Updated backend  %d is not correct for HTTPServerClose: %s", i, backend.HTTPServerClose)
		}
		if backend.Httpclose != input.expectedHttpclose {
			t.Errorf("Updated backend %d is not correct for Httpclose: %s", i, backend.Httpclose)
		}

	}

	// TestDeleteBackend
	err = clientTest.DeleteBackend("special-httpconnectionmode", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}
}
