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

func frontendExpectation() map[string]models.Frontends {
	initStructuredExpected()
	res := StructuredToFrontendMap()
	// Add individual entries
	for _, vs := range res {
		for _, v := range vs {
			key := v.Name
			res[key] = models.Frontends{v}
		}
	}
	return res
}

func TestGetFrontends(t *testing.T) { //nolint:gocognit
	m := make(map[string]models.Frontends)
	v, frontends, err := clientTest.GetFrontends("")
	if err != nil {
		t.Error(err.Error())
	}

	if len(frontends) != 2 {
		t.Errorf("%v frontends returned, expected 2", len(frontends))
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}
	m[""] = frontends

	checkFrontends(t, m)
}

func checkFrontends(t *testing.T, got map[string]models.Frontends) {
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

func TestGetFrontend(t *testing.T) {
	m := make(map[string]models.Frontends)

	v, f, err := clientTest.GetFrontend("test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}
	m["test"] = models.Frontends{f}
	checkFrontends(t, m)

	_, err = f.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	_, _, err = clientTest.GetFrontend("doesnotexist", "")
	if err == nil {
		t.Error("Should throw error, non existent frontend")
	}
}

func TestCreateEditDeleteFrontend(t *testing.T) {
	// TestCreateFrontend
	mConn := int64(3000)
	tOut := int64(2)
	clitcpkaCnt := int64(10)
	clitcpkaTimeout := int64(10000)
	statsRealm := "Haproxy Stats"
	f := &models.Frontend{
		From:                     "unnamed_defaults_1",
		Name:                     "created",
		Mode:                     "tcp",
		Maxconn:                  &mConn,
		Httplog:                  true,
		HTTPConnectionMode:       "http-keep-alive",
		HTTPKeepAliveTimeout:     &tOut,
		BindProcess:              "4",
		Logasap:                  "disabled",
		UniqueIDFormat:           "%{+X}o_%fi:%fp_%Ts_%rt:%pid",
		UniqueIDHeader:           "X-Unique-Id",
		AcceptInvalidHTTPRequest: "enabled",
		DisableH2Upgrade:         "enabled",
		ClitcpkaCnt:              &clitcpkaCnt,
		ClitcpkaIdle:             &clitcpkaTimeout,
		ClitcpkaIntvl:            &clitcpkaTimeout,
		HTTPIgnoreProbes:         "enabled",
		HTTPUseProxyHeader:       "enabled",
		Httpslog:                 "enabled",
		IndependentStreams:       "enabled",
		Nolinger:                 "enabled",
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
	}

	err := clientTest.CreateFrontend(f, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, frontend, err := clientTest.GetFrontend("created", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(frontend, f) {
		fmt.Printf("Created frontend: %v\n", frontend)
		fmt.Printf("Given frontend: %v\n", f)
		t.Error("Created frontend not equal to given frontend")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	err = clientTest.CreateFrontend(f, "", version)
	if err == nil {
		t.Error("Should throw error frontend already exists")
		version++
	}

	// TestEditFrontend
	mConn = int64(4000)
	f = &models.Frontend{
		Name:               "created",
		Mode:               "tcp",
		Maxconn:            &mConn,
		Backlog:            misc.Int64P(1024),
		Clflog:             true,
		HTTPConnectionMode: "httpclose",
		BindProcess:        "3",
		MonitorURI:         "/healthz",
		MonitorFail: &models.MonitorFail{
			Cond:     misc.StringP("if"),
			CondTest: misc.StringP("site_is_dead"),
		},
		Compression: &models.Compression{
			Offload: true,
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
	}

	err = clientTest.EditFrontend("created", f, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, frontend, err = clientTest.GetFrontend("created", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(frontend, f) {
		fmt.Printf("Edited frontend: %v\n", frontend)
		fmt.Printf("Given frontend: %v\n", f)
		t.Error("Edited frontend not equal to given frontend")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	// TestDeleteFrontend
	err = clientTest.DeleteFrontend("created", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ := clientTest.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	err = clientTest.DeleteFrontend("created", "", 999999)
	if err != nil {
		if confErr, ok := err.(*configuration.ConfError); ok {
			if !confErr.Is(configuration.ErrVersionMismatch) {
				t.Error("Should throw configuration.ErrVersionMismatch error")
			}
		} else {
			t.Error("Should throw configuration.ErrVersionMismatch error")
		}
	}
	_, _, err = clientTest.GetFrontend("created", "")
	if err == nil {
		t.Error("DeleteFrontend failed, frontend test still exists")
	}

	err = clientTest.DeleteFrontend("doesnotexist", "", version)
	if err == nil {
		t.Error("Should throw error, non existent frontend")
		version++
	}
}
