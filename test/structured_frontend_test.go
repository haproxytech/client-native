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
	"github.com/haproxytech/client-native/v6/configuration"
	"github.com/haproxytech/client-native/v6/misc"
	"github.com/haproxytech/client-native/v6/models"
	"github.com/stretchr/testify/require"
)

func TestGetStructuredFrontends(t *testing.T) { //nolint:gocognit
	clientTest, filename, err := getTestClient()
	require.NoError(t, err)
	defer os.Remove(filename)
	version := int64(1)

	m := make(map[string]models.Frontends)
	v, frontends, err := clientTest.GetStructuredFrontends("")
	require.NoError(t, err)
	require.Equal(t, 2, len(frontends), "%v frontends returned, expected 2", len(frontends))
	require.Equal(t, version, v, "Version %v returned, expected %v", v, version)

	m[""] = frontends

	checkStructuredFrontends(t, m)
}

func checkStructuredFrontends(t *testing.T, got map[string]models.Frontends) {
	exp := frontendExpectation()
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

func TestGetStructuredFrontend(t *testing.T) {
	clientTest, filename, err := getTestClient()
	require.NoError(t, err)
	defer os.Remove(filename)
	version := int64(1)

	m := make(map[string]models.Frontends)

	v, f, err := clientTest.GetStructuredFrontend("test", "")
	require.NoError(t, err)
	require.Equal(t, version, v, "Version %v returned, expected %v", v, version)

	m["test"] = models.Frontends{f}
	checkStructuredFrontends(t, m)

	_, err = f.MarshalBinary()
	require.NoError(t, err)

	_, _, err = clientTest.GetStructuredFrontend("doesnotexist", "")
	require.Error(t, err, "Should throw error, non existent frontend")
}

func TestCreateEditDeleteStructuredFrontend(t *testing.T) {
	clientTest, filename, err := getTestClient()
	require.NoError(t, err)
	defer os.Remove(filename)
	version := int64(1)

	// TestCreateFrontend
	mConn := int64(3000)
	tOut := int64(2)
	clitcpkaCnt := int64(10)
	clitcpkaTimeout := int64(10000)
	statsRealm := "Haproxy Stats"
	f := &models.Frontend{
		FrontendBase: models.FrontendBase{
			From:                                "unnamed_defaults_1",
			Name:                                "created",
			Mode:                                "tcp",
			Maxconn:                             &mConn,
			Httplog:                             true,
			HTTPConnectionMode:                  "http-keep-alive",
			HTTPKeepAliveTimeout:                &tOut,
			Logasap:                             "disabled",
			UniqueIDFormat:                      "%{+X}o_%fi:%fp_%Ts_%rt:%pid",
			UniqueIDHeader:                      "X-Unique-Id",
			AcceptInvalidHTTPRequest:            "enabled",
			AcceptUnsafeViolationsInHTTPRequest: "enabled",
			DisableH2Upgrade:                    "enabled",
			ClitcpkaCnt:                         &clitcpkaCnt,
			ClitcpkaIdle:                        &clitcpkaTimeout,
			ClitcpkaIntvl:                       &clitcpkaTimeout,
			HTTPIgnoreProbes:                    "enabled",
			HTTPUseProxyHeader:                  "enabled",
			Httpslog:                            "enabled",
			IndependentStreams:                  "enabled",
			Nolinger:                            "enabled",
			Originalto: &models.Originalto{
				Enabled: misc.StringP("enabled"),
				Except:  "127.0.0.1",
				Header:  "X-Client-Dst",
			},
			SocketStats:         "enabled",
			TCPSmartAccept:      "enabled",
			DontlogNormal:       "enabled",
			HTTPNoDelay:         "enabled",
			SpliceAuto:          "enabled",
			SpliceRequest:       "enabled",
			SpliceResponse:      "enabled",
			IdleCloseOnResponse: "enabled",
			StatsOptions: &models.StatsOptions{
				StatsShowModules: true,
				StatsRealm:       true,
				StatsRealmRealm:  &statsRealm,
				StatsAuths: []*models.StatsAuth{
					{User: misc.StringP("user1"), Passwd: misc.StringP("pwd1")},
					{User: misc.StringP("user2"), Passwd: misc.StringP("pwd2")},
				},
			},
			Enabled: true,
		},
		ACLList: models.Acls{
			&models.ACL{
				ACLName:   "acl1",
				Criterion: "src_port",
				Value:     "0:101",
			},
		},
		CaptureList: models.Captures{
			&models.Capture{
				Length: 1,
				Type:   "request",
			},
		},
		FilterList: models.Filters{
			&models.Filter{
				TraceName:       "BEFORE-HTTP-COMP",
				TraceRndParsing: true,
				Type:            "trace",
			},
		},
		BackendSwitchingRuleList: models.BackendSwitchingRules{
			&models.BackendSwitchingRule{
				Name:     "bsr1",
				Cond:     "if",
				CondTest: "TRUE",
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
		TCPRequestRuleList: models.TCPRequestRules{
			&models.TCPRequestRule{
				Cond:     "if",
				CondTest: "FALSE",
				Action:   "accept",
				Type:     "connection",
			},
		},
		Binds: map[string]models.Bind{
			"192.168.1.1:9200": {
				BindParams: models.BindParams{Name: "192.168.1.1:9200"},
				Address:    "192.168.1.1",
				Port:       misc.Ptr[int64](9200),
			},
		},
		LogTargetList: models.LogTargets{
			&models.LogTarget{
				Address:  "192.169.0.1",
				Facility: "mail",
				Global:   true,
			},
		},
	}
	err = clientTest.CreateStructuredFrontend(f, "", version)
	require.NoError(t, err)
	version++

	v, frontend, err := clientTest.GetStructuredFrontend("created", "")
	require.NoError(t, err)
	require.True(t, frontend.Equal(*f), "frontend=%s - diff %v", frontend.Name, cmp.Diff(*frontend, *f))
	require.Equal(t, version, v, "Version %v returned, expected %v", v, version)

	err = clientTest.CreateStructuredFrontend(f, "", version)
	require.Error(t, err, "Should throw error frontend already exists")

	// TestEditFrontend
	mConn = int64(4000)
	f = &models.Frontend{
		FrontendBase: models.FrontendBase{
			Name:               "created",
			Mode:               "tcp",
			Maxconn:            &mConn,
			Backlog:            misc.Int64P(1024),
			Clflog:             true,
			HTTPConnectionMode: "httpclose",
			MonitorURI:         "/healthz",
			MonitorFail: &models.MonitorFail{
				Cond:     misc.StringP("if"),
				CondTest: misc.StringP("site_is_dead"),
			},
			Compression: &models.Compression{
				Offload: true,
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
			ClitcpkaCnt:        &clitcpkaCnt,
			ClitcpkaIdle:       &clitcpkaTimeout,
			ClitcpkaIntvl:      &clitcpkaTimeout,
			HTTPIgnoreProbes:   "disabled",
			HTTPUseProxyHeader: "disabled",
			Httpslog:           "disabled",
			IndependentStreams: "disabled",
			Nolinger:           "disabled",
			Originalto: &models.Originalto{
				Enabled: misc.StringP("enabled"),
				Except:  "127.0.0.1",
				Header:  "X-Client-Dst",
			},
			SocketStats:         "disabled",
			TCPSmartAccept:      "disabled",
			DontlogNormal:       "disabled",
			HTTPNoDelay:         "disabled",
			SpliceAuto:          "disabled",
			SpliceRequest:       "disabled",
			SpliceResponse:      "disabled",
			IdleCloseOnResponse: "disabled",
			StatsOptions: &models.StatsOptions{
				StatsShowModules: true,
				StatsRealm:       true,
				StatsRealmRealm:  &statsRealm,
				StatsAuths: []*models.StatsAuth{
					{User: misc.StringP("new_user1"), Passwd: misc.StringP("new_pwd1")},
					{User: misc.StringP("new_user2"), Passwd: misc.StringP("new_pwd2")},
				},
			},
			EmailAlert: &models.EmailAlert{
				From:    misc.StringP("srv01@example.com"),
				To:      misc.StringP("problems@example.com"),
				Level:   "warning",
				Mailers: misc.StringP("localmailer1"),
			},
		},
	}

	err = clientTest.EditStructuredFrontend("created", f, "", version)
	require.NoError(t, err)
	version++

	v, frontend, err = clientTest.GetStructuredFrontend("created", "")
	require.NoError(t, err)
	require.True(t, frontend.Equal(*f), "frontend=%s - diff %v", frontend.Name, cmp.Diff(*frontend, *f))
	require.Equal(t, version, v, "Version %v returned, expected %v", v, version)

	// TestDeleteFrontend
	err = clientTest.DeleteFrontend("created", "", version)
	require.NoError(t, err)
	version++

	v, _ = clientTest.GetVersion("")
	require.Equal(t, version, v, "Version %v returned, expected %v", v, version)

	err = clientTest.DeleteFrontend("created", "", 999999)
	require.Error(t, err, "Should throw error, non existent frontend")
	require.ErrorIs(t, err, configuration.ErrVersionMismatch, "Should throw configuration.ErrVersionMismatch error")

	_, _, err = clientTest.GetStructuredFrontend("created", "")
	require.Error(t, err, "DeleteFrontend failed, frontend test still exists")

	err = clientTest.DeleteFrontend("doesnotexist", "", version)
	require.Error(t, err, "Should throw error, non existent frontend")
}
