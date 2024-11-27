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
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"

	"github.com/haproxytech/client-native/v6/configuration"
	"github.com/haproxytech/client-native/v6/misc"
	"github.com/haproxytech/client-native/v6/models"
)

func TestGetStructuredBackends(t *testing.T) { //nolint:gocognit,gocyclo
	clientTest, filename, err := getTestClient()
	require.NoError(t, err)
	defer os.Remove(filename)
	version := int64(1)

	m := make(map[string]models.Backends)
	v, backends, err := clientTest.GetStructuredBackends("")
	require.NoError(t, err)
	require.Equal(t, version, v, "Version %v returned, expected %v", v, version)
	m[""] = backends

	checkStructuredBackends(t, m)
}

func checkStructuredBackends(t *testing.T, got map[string]models.Backends) {
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

func TestGetStructuredBackend(t *testing.T) {
	clientTest, filename, err := getTestClient()
	require.NoError(t, err)
	defer os.Remove(filename)
	version := int64(1)

	m := make(map[string]models.Backends)

	v, b, err := clientTest.GetStructuredBackend("test", "")
	require.NoError(t, err)

	require.Equal(t, version, v, "Version %v returned, expected %v", v, version)
	m["test"] = models.Backends{b}

	_, err = b.MarshalBinary()
	require.NoError(t, err)

	_, _, err = clientTest.GetStructuredBackend("doesnotexist", "")
	require.Error(t, err, "Should throw error, non existent bck")

	checkStructuredBackends(t, m)
}

func TestCreateEditDeleteStructuredBackend(t *testing.T) {
	clientTest, filename, err := getTestClient()
	require.NoError(t, err)
	defer os.Remove(filename)
	version := int64(1)

	// TestCreateBackend
	tOut := int64(5)
	cookieName := "BLA"
	balanceAlgorithm := "uri"
	srvtcpkaCnt := int64(10)
	srvtcpkaTimeout := int64(10000)
	statsRealm := "Haproxy Stats"
	b := &models.Backend{
		BackendBase: models.BackendBase{
			Name: "created",
			Mode: "http",
			Balance: &models.Balance{
				Algorithm: &balanceAlgorithm,
				URILen:    100,
				URIDepth:  250,
			},
			Cookie: &models.Cookie{
				Domains: []*models.Domain{
					{Value: "dom1"},
					{Value: "dom2"},
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
			AcceptInvalidHTTPResponse:            "enabled",
			AcceptUnsafeViolationsInHTTPResponse: "enabled",
			Compression: &models.Compression{
				Offload:   true,
				Direction: "both",
				TypesReq: []string{
					"text/html",
					"text/plain",
				},
				TypesRes: []string{
					"text/plain",
				},
				AlgoReq: "deflate",
				AlgosRes: []string{
					"deflate",
					"gzip",
				},
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
		},
		ACLList: models.Acls{
			&models.ACL{
				ACLName:   "acl1",
				Criterion: "src_port",
				Value:     "0:101",
			},
		},
		FilterList: models.Filters{
			&models.Filter{
				TraceName:       "BEFORE-HTTP-COMP",
				TraceRndParsing: true,
				Type:            "trace",
			},
		},
		HTTPAfterResponseRuleList: models.HTTPAfterResponseRules{
			&models.HTTPAfterResponseRule{
				Cond:      "if",
				CondTest:  "FALSE",
				MapFile:   "map.lst",
				MapKeyfmt: "%[src]",
				Type:      "del-map",
			},
		},
		HTTPCheckList: models.HTTPChecks{
			&models.HTTPCheck{
				Type: "send-state",
			},
		},
		HTTPErrorRuleList: models.HTTPErrorRules{
			&models.HTTPErrorRule{
				ReturnContent:       "/test/503",
				ReturnContentFormat: "file",
				ReturnContentType:   misc.Ptr("\"application/json\""),
				Status:              503,
				Type:                "status",
			},
		},
		HTTPRequestRuleList: models.HTTPRequestRules{
			&models.HTTPRequestRule{
				Cond:      "if",
				CondTest:  "FALSE",
				MapFile:   "map.lst",
				MapKeyfmt: "%[src]",
				Type:      "del-map",
			},
		},
		HTTPResponseRuleList: models.HTTPResponseRules{
			&models.HTTPResponseRule{
				Cond:      "if",
				CondTest:  "FALSE",
				MapFile:   "map.lst",
				MapKeyfmt: "%[src]",
				Type:      "del-map",
			},
		},
		ServerSwitchingRuleList: models.ServerSwitchingRules{
			&models.ServerSwitchingRule{
				TargetServer: "srv1",
				Cond:         "if",
				CondTest:     "TRUE",
			},
		},
		StickRuleList: models.StickRules{
			&models.StickRule{
				Type:     "match",
				Pattern:  "src",
				Cond:     "if",
				CondTest: "TRUE",
			},
		},
		TCPCheckRuleList: models.TCPChecks{
			&models.TCPCheck{
				Pattern: "src",
				Match:   "string",
				Action:  "expect",
			},
		},
		TCPRequestRuleList: models.TCPRequestRules{
			&models.TCPRequestRule{
				Cond:     "if",
				CondTest: "FALSE",
				Action:   "accept",
				Type:     "connection",
			},
		},
		TCPResponseRuleList: models.TCPResponseRules{
			&models.TCPResponseRule{
				Cond:     "if",
				CondTest: "FALSE",
				Action:   "accept",
				Type:     "content",
			},
		},
		LogTargetList: models.LogTargets{
			&models.LogTarget{
				Address:  "192.169.0.1",
				Facility: "mail",
				Global:   true,
			},
		},
		Servers: map[string]models.Server{
			"webserv": {
				Address: "192.168.1.1",
				Name:    "webserv",
				ID:      misc.Ptr[int64](1234),
				Port:    misc.Ptr[int64](9200),
			},
		},
		ServerTemplates: map[string]models.ServerTemplate{
			"webserv": {
				Fqdn:       "google.com",
				NumOrRange: "1-10",
				Prefix:     "webserv",
				Port:       misc.Ptr[int64](9200),
			},
		},
	}

	err = clientTest.CreateStructuredBackend(b, "", version)
	require.NoError(t, err)
	version++

	v, backend, err := clientTest.GetStructuredBackend("created", "")
	require.NoError(t, err)

	require.True(t, backend.Equal(*b), "backend=%s - diff %v", backend.Name, cmp.Diff(*backend, *b))

	require.Equal(t, version, v, "Version %v returned, expected %v", v, version)

	err = clientTest.CreateStructuredBackend(b, "", version)
	require.Error(t, err, "Should throw error bck already exists")

	// TestEditBackend
	tOut = int64(3)
	e := int64(1200000)
	kl := int64(128)
	s := int64(25600)
	backends := []*models.Backend{
		{
			BackendBase: models.BackendBase{
				From: "test_defaults",
				Name: "created",
				Mode: "http",
				Balance: &models.Balance{
					Algorithm: &balanceAlgorithm,
					URILen:    10,
					URIDepth:  25,
				},
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
			}},
		{
			BackendBase: models.BackendBase{
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
				ForcePersistList: []*models.ForcePersist{
					{Cond: misc.StringP("unless"), CondTest: misc.StringP("invalid_src")},
					{Cond: misc.StringP("if"), CondTest: misc.StringP("auth_ok")},
				},
				IgnorePersistList: []*models.IgnorePersist{
					{Cond: misc.StringP("if"), CondTest: misc.StringP("host_www")},
					{Cond: misc.StringP("unless"), CondTest: misc.StringP("missing_cl")},
				},
			}},
	}

	for _, bck := range backends {
		err := clientTest.EditStructuredBackend("created", bck, "", version)
		require.NoError(t, err)
		version++

		v, backend, err := clientTest.GetStructuredBackend("created", "")
		require.NoError(t, err)
		require.True(t, backend.Equal(*bck), "backend=%s - diff %v", backend.Name, cmp.Diff(*backend, *bck))
		require.Equal(t, version, v, "Version %v returned, expected %v", v, version)
	}

	// TestDeleteBackend
	err = clientTest.DeleteBackend("created", "", version)
	require.NoError(t, err)
	version++

	v, _ = clientTest.GetVersion("")
	require.Equal(t, version, v, "Version %v returned, expected %v", v, version)

	err = clientTest.DeleteBackend("created", "", 999999999)
	require.Error(t, err, "Should throw error, non existent frontend")
	require.ErrorIs(t, err, configuration.ErrVersionMismatch, "Should throw configuration.ErrVersionMismatch error")

	_, _, err = clientTest.GetStructuredBackend("created", "")
	require.Error(t, err, "DeleteBackend failed, bck test still exists")

	err = clientTest.DeleteBackend("doesnotexist", "", version)
	require.Error(t, err, "Should throw error, non existent bck")
}

func TestCreateEditDeleteStructuredBackendHTTPConnectionMode(t *testing.T) {
	clientTest, filename, err := getTestClient()
	require.NoError(t, err)
	defer os.Remove(filename)
	version := int64(1)

	// TestCreateBackend
	tOut := int64(5)

	// Backend with HTTPConnectionMode only
	b := &models.Backend{
		BackendBase: models.BackendBase{
			Name: "special-httpconnectionmode",
			Mode: "http",
			DefaultServer: &models.DefaultServer{
				ServerParams: models.ServerParams{
					Fall:  &tOut,
					Inter: &tOut,
				},
			},
			HTTPConnectionMode: "http-keep-alive",
		},
	}

	err = clientTest.CreateStructuredBackend(b, "", version)
	require.NoError(t, err)
	version++

	v, backend, err := clientTest.GetStructuredBackend("special-httpconnectionmode", "")
	require.NoError(t, err)
	require.Equal(t, version, v, "Version %v returned, expected %v", v, version)
	require.Equal(t, "http-keep-alive", backend.HTTPConnectionMode, "Created backend is not correct for HTTPConnectionMode: %s", backend.HTTPConnectionMode)
	require.Equal(t, version, v, "Version %v returned, expected %v", v, version)

	err = clientTest.CreateStructuredBackend(b, "", version)
	require.Error(t, err, "Should throw error bck already exists")

	type testinput struct {
		backend                    *models.Backend
		expectedHTTPConnectionMode string
	}
	// TestEditBackend
	inputs := []testinput{
		{
			// Update HTTPConnectionMode
			backend: &models.Backend{
				BackendBase: models.BackendBase{
					Name: "special-httpconnectionmode",
					Mode: "http",
					DefaultServer: &models.DefaultServer{
						ServerParams: models.ServerParams{
							Fall:  &tOut,
							Inter: &tOut,
						},
					},
					HTTPConnectionMode: "httpclose",
				}},
			expectedHTTPConnectionMode: "httpclose",
		},
		{
			// Use both - Priority on HTTPConnection
			backend: &models.Backend{
				BackendBase: models.BackendBase{
					Name: "special-httpconnectionmode",
					Mode: "http",
					DefaultServer: &models.DefaultServer{
						ServerParams: models.ServerParams{
							Fall:  &tOut,
							Inter: &tOut,
						},
					},
					HTTPConnectionMode: "http-keep-alive",
				}},
			expectedHTTPConnectionMode: "http-keep-alive",
		},
		{
			// remove option
			backend: &models.Backend{
				BackendBase: models.BackendBase{
					Name: "special-httpconnectionmode",
					Mode: "http",
					DefaultServer: &models.DefaultServer{
						ServerParams: models.ServerParams{
							Fall:  &tOut,
							Inter: &tOut,
						},
					},
					HTTPConnectionMode: "",
				}},
			expectedHTTPConnectionMode: "",
		},
	}

	for i, input := range inputs {
		err := clientTest.EditStructuredBackend("special-httpconnectionmode", input.backend, "", version)
		require.NoError(t, err)
		version++

		_, backend, err := clientTest.GetStructuredBackend("special-httpconnectionmode", "")
		require.NoError(t, err)
		require.Equal(t, input.expectedHTTPConnectionMode, backend.HTTPConnectionMode, "Updated backend %d is not correct for HTTPConnectionMode: %s", i, backend.HTTPConnectionMode)
	}

	// TestDeleteBackend
	err = clientTest.DeleteBackend("special-httpconnectionmode", "", version)
	require.NoError(t, err)
	version++
}
