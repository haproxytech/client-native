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
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/haproxytech/client-native/v6/misc"
	"github.com/haproxytech/client-native/v6/models"
	"github.com/stretchr/testify/require"
)

func globalExpcectation() *models.Global {
	initStructuredExpected()
	res := StructuredToGlobalMap()
	return &res
}

func TestGetGlobal(t *testing.T) {
	v, global, err := clientTest.GetGlobalConfiguration("")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	checkGlobal(t, global)
}

func checkGlobal(t *testing.T, global *models.Global) {
	want := globalExpcectation()
	require.True(t, global.Equal(*want), "diff %v", cmp.Diff(*global, *want))
}

func TestPutGlobal(t *testing.T) {
	tOut := int64(3600)
	n := "1/1"
	v := "0"
	a := "/var/run/haproxy.sock"
	f := "/etc/foo.lua"
	luaPrependPath := "/usr/share/haproxy-lua/?/init.lua"
	enabled := "enabled"
	g := &models.Global{
		Daemon: "enabled",
		CPUMaps: []*models.CPUMap{
			{
				Process: &n,
				CPUSet:  &v,
			},
		},
		RuntimeAPIs: []*models.RuntimeAPI{
			{
				Address: &a,
				BindParams: models.BindParams{
					Level: "admin",
				},
			},
		},
		Maxconn:                1000,
		SslDefaultBindCiphers:  "test",
		SslDefaultBindOptions:  "ssl-min-ver TLSv1.0 no-tls-tickets",
		SslDefaultServerCurves: "secp384r1",
		StatsTimeout:           &tOut,
		TuneSslDefaultDhParam:  1024,
		ExternalCheck:          false,
		LuaPrependPath: []*models.LuaPrependPath{
			{
				Path: &luaPrependPath,
				Type: "cpath",
			},
		},
		LuaLoads: []*models.LuaLoad{
			{
				File: &f,
			},
		},
		LogSendHostname: &models.GlobalLogSendHostname{
			Enabled: &enabled,
			Param:   "something",
		},
		TuneOptions: &models.GlobalTuneOptions{
			DisableZeroCopyForwarding:      true,
			EventsMaxEventsAtOnce:          50,
			H1ZeroCopyFwdRecv:              "disabled",
			H1ZeroCopyFwdSend:              "disabled",
			H2ZeroCopyFwdSend:              "disabled",
			LuaLogLoggers:                  "disabled",
			LuaLogStderr:                   "disabled",
			MaxChecksPerThread:             misc.Int64P(20),
			PeersMaxUpdatesAtOnce:          100,
			PtZeroCopyForwarding:           "disabled",
			QuicFrontendConnTxBuffersLimit: nil,
			QuicFrontendMaxIdleTimeout:     misc.Int64P(5000),
			QuicSocketOwner:                "listener",
			RcvbufBackend:                  misc.Int64P(8192),
			RcvbufFrontend:                 misc.Int64P(4096),
			SndbufBackend:                  misc.Int64P(1234),
			SndbufFrontend:                 misc.Int64P(5678),
			SslOcspUpdateMaxDelay:          misc.Int64P(48),
			SslOcspUpdateMinDelay:          misc.Int64P(49),
			StickCounters:                  misc.Int64P(50),
		},
		HttpclientResolversDisabled: "disabled",
		HttpclientResolversPrefer:   "ipv6",
		HttpclientResolversID:       "my2",
		HttpclientRetries:           5,
		HttpclientSslCaFile:         "my_ca_file.ca",
		HttpclientTimeoutConnect:    misc.Int64P(5000),
		HttpclientSslVerify:         misc.StringP(""),
		UID:                         1234,
		WurflOptions:                &models.GlobalWurflOptions{},
		DeviceAtlasOptions:          &models.GlobalDeviceAtlasOptions{},
		FiftyOneDegreesOptions:      &models.GlobalFiftyOneDegreesOptions{},
		StatsMaxconn:                misc.Int64P(30),
		Anonkey:                     misc.Int64P(40),
		NumaCPUMapping:              "disabled",
		DefaultPath: &models.GlobalDefaultPath{
			Type: "origin",
			Path: "/some/other/path",
		},
		NoQuic:                        false,
		ClusterSecret:                 "",
		SslDefaultServerSigalgs:       "ECDSA+SHA256",
		SslDefaultServerClientSigalgs: "ECDSA+SHA256",
		SslPropquery:                  "foo",
		SslProvider:                   "my_provider",
		SslProviderPath:               "providers/",
		Setcap:                        "none",
		LimitedQuic:                   true,
	}

	err := clientTest.PushGlobalConfiguration(g, "", version)

	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	ver, global, err := clientTest.GetGlobalConfiguration("")
	if err != nil {
		t.Error(err.Error())
	}

	var givenJSON []byte
	givenJSON, err = g.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	var onDiskJSON []byte
	onDiskJSON, err = global.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	if string(givenJSON) != string(onDiskJSON) {
		fmt.Printf("Created global: %v\n", string(onDiskJSON))
		fmt.Printf("Given global: %v\n", string(givenJSON))
		t.Error("Created global not equal to given global")
	}

	if ver != version {
		t.Error("Version not incremented!")
	}

	err = clientTest.PushGlobalConfiguration(g, "", 55)

	if err == nil {
		t.Error("Should have returned version conflict.")
	}
}
