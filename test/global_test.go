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
	require.NoError(t, err, "failed to get global configuration")
	require.Equal(t, version, v, "version mismatch")
	checkGlobal(t, global)
}

func checkGlobal(t *testing.T, global *models.Global) {
	want := globalExpcectation()
	require.True(t, global.GlobalBase.Equal(want.GlobalBase), "diff %v", cmp.Diff(global.GlobalBase, want.GlobalBase))
}

func TestPutGlobal(t *testing.T) {
	g := &models.Global{
		GlobalBase: getGlobalBase(),
	}

	err := clientTest.PushGlobalConfiguration(g, "", version)
	require.NoError(t, err)
	version++

	ver, global, err := clientTest.GetGlobalConfiguration("")
	require.NoError(t, err)

	require.True(t, g.GlobalBase.Equal(global.GlobalBase, models.Options{NilSameAsEmpty: false}), "diff %v", cmp.Diff(g.GlobalBase, global.GlobalBase))
	if ver != version {
		t.Fatal("Version not incremented!")
	}
}

func TestPutEmptyGlobal(t *testing.T) {
	g := &models.Global{}

	err := clientTest.PushGlobalConfiguration(g, "", version)
	require.NoError(t, err)
	version++

	ver, global, err := clientTest.GetGlobalConfiguration("")
	require.NoError(t, err)

	require.True(t, g.GlobalBase.Equal(global.GlobalBase, models.Options{NilSameAsEmpty: false}), "diff %v", cmp.Diff(g.GlobalBase, global.GlobalBase))
	if ver != version {
		t.Fatal("Version not incremented!")
	}

	g = globalExpcectation()
	err = clientTest.PushGlobalConfiguration(g, "", version)
	require.NoError(t, err)
	version++

	ver, global, err = clientTest.GetGlobalConfiguration("")
	require.NoError(t, err)

	require.True(t, g.GlobalBase.Equal(global.GlobalBase, models.Options{NilSameAsEmpty: false}), "diff %v", cmp.Diff(g.GlobalBase, global.GlobalBase))
	if ver != version {
		t.Fatal("Version not incremented!")
	}
}

func getGlobalBase() models.GlobalBase {
	tOut := int64(3600)
	n := "1/1"
	v := "0"
	a := "/var/run/haproxy.sock"
	f := "/etc/foo.lua"
	luaPrependPath := "/usr/share/haproxy-lua/?/init.lua"
	enabled := "enabled"
	return models.GlobalBase{
		Daemon: true,
		CPUMaps: []*models.CPUMap{
			{
				Process: &n,
				CPUSet:  &v,
			},
		},
		CPUPolicy: models.GlobalBaseCPUPolicyEfficiency,
		CPUSets: []*models.CPUSet{
			{
				Directive: misc.StringP("reset"),
			},
			{
				Directive: misc.StringP("drop-cpu"),
				Set:       "1,3",
			},
			{
				Directive: misc.StringP("only-thread"),
				Set:       "4-9",
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
		PerformanceOptions: &models.PerformanceOptions{
			Maxconn:         1000,
			ProfilingMemory: "enabled",
			ThreadHardLimit: misc.Int64P(100),
		},
		SslOptions: &models.SslOptions{
			DefaultBindCiphers:         "test",
			DefaultBindOptions:         "ssl-min-ver TLSv1.0 no-tls-tickets",
			DefaultServerCurves:        "secp384r1",
			DefaultServerSigalgs:       "ECDSA+SHA256",
			DefaultServerClientSigalgs: "ECDSA+SHA256",
			Propquery:                  "foo",
			Provider:                   "my_provider",
			ProviderPath:               "providers/",
			SecurityLevel:              misc.Int64P(2),
		},
		StatsTimeout:  &tOut,
		ExternalCheck: false,
		LuaOptions: &models.LuaOptions{
			PrependPath: []*models.LuaPrependPath{
				{
					Path: &luaPrependPath,
					Type: "cpath",
				},
			},
			Loads: []*models.LuaLoad{
				{
					File: &f,
				},
			},
		},
		LogSendHostname: &models.GlobalLogSendHostname{
			Enabled: &enabled,
			Param:   "something",
		},
		// TuneLuaOptions: &models.TuneLuaOptions{
		// 	LogLoggers: "disabled",
		// 	LogStderr:  "disabled",
		// },
		TuneBufferOptions: &models.TuneBufferOptions{
			RcvbufBackend:  misc.Int64P(8192),
			RcvbufFrontend: misc.Int64P(4096),
			SndbufBackend:  misc.Int64P(1234),
			SndbufFrontend: misc.Int64P(5678),
		},
		TuneQuicOptions: &models.TuneQuicOptions{
			FrontendConnTxBuffersLimit: nil,
			FrontendMaxIdleTimeout:     misc.Int64P(5000),
			SocketOwner:                "listener",
		},
		TuneSslOptions: &models.TuneSslOptions{
			OcspUpdateMaxDelay: misc.Int64P(48),
			OcspUpdateMinDelay: misc.Int64P(49),
		},
		TuneOptions: &models.TuneOptions{
			DisableZeroCopyForwarding: true,
			EventsMaxEventsAtOnce:     50,
			H1ZeroCopyFwdRecv:         "disabled",
			H1ZeroCopyFwdSend:         "disabled",
			H2ZeroCopyFwdSend:         "disabled",
			MaxChecksPerThread:        misc.Int64P(20),
			PeersMaxUpdatesAtOnce:     100,
			PtZeroCopyForwarding:      "disabled",
			StickCounters:             misc.Int64P(50),
		},
		HTTPClientOptions: &models.HTTPClientOptions{
			ResolversDisabled: "disabled",
			ResolversPrefer:   "ipv6",
			ResolversID:       "my2",
			Retries:           5,
			SslCaFile:         "my_ca_file.ca",
			TimeoutConnect:    misc.Int64P(5000),
			SslVerify:         misc.StringP(""),
		},
		DebugOptions: &models.DebugOptions{
			Anonkey: misc.Int64P(40),
		},
		UID:            1234,
		StatsMaxconn:   misc.Int64P(30),
		NumaCPUMapping: "disabled",
		DefaultPath: &models.GlobalDefaultPath{
			Type: "origin",
			Path: "/some/other/path",
		},
		NoQuic:        false,
		ClusterSecret: "",
		Setcap:        "none",
		LimitedQuic:   true,
		Harden: &models.GlobalHarden{
			RejectPrivilegedPorts: &models.GlobalHardenRejectPrivilegedPorts{
				Quic: "enabled",
			},
		},
		HTTPErrCodes: []*models.HTTPCodes{{Value: misc.StringP("100-150 -123 +599")}},
	}
}
