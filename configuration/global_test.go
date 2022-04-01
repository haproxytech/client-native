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

	"github.com/haproxytech/client-native/v3/models"
)

func TestGetGlobal(t *testing.T) {
	v, global, err := clientTest.GetGlobalConfiguration("")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	if global.Daemon != "enabled" {
		t.Errorf("Daemon is %v, expected enabled", global.Daemon)
	}
	if len(global.RuntimeAPIs) == 1 {
		if *global.RuntimeAPIs[0].Address != "/var/run/haproxy.sock" {
			t.Errorf("RuntimeAPI.Address is %v, expected /var/run/haproxy.sock", *global.RuntimeAPIs[0].Address)
		}
		if global.RuntimeAPIs[0].Level != "admin" {
			t.Errorf("RuntimeAPI.Level is %v, expected admin", global.RuntimeAPIs[0].Level)
		}
		if global.RuntimeAPIs[0].Mode != "0660" {
			t.Errorf("RuntimeAPI.Mode is %v, expected 0660", global.RuntimeAPIs[0].Mode)
		}
	} else {
		t.Errorf("RuntimeAPI is not set")
	}
	if global.CaBase != "/etc/ssl/certs" {
		t.Errorf("CaBase is %v, expected /etc/ssl/certs", global.CaBase)
	}
	if global.CrtBase != "/etc/ssl/private" {
		t.Errorf("CrtBase is %v, expected /etc/ssl/private", global.CrtBase)
	}
	if global.Nbproc != 4 {
		t.Errorf("Nbproc is %v, expected 4", global.Nbproc)
	}
	if global.Maxconn != 2000 {
		t.Errorf("Maxconn is %v, expected 2000", global.Maxconn)
	}
	if global.ExternalCheck != true {
		t.Errorf("ExternalCheck is false, expected true")
	}
	if len(global.LuaPrependPath) == 2 {
		if *global.LuaPrependPath[0].Path != "/usr/share/haproxy-lua/?/init.lua" {
			t.Errorf("LuaPrependPath.Path is %v, expected /usr/share/haproxy-lua/?/init.lua", *global.LuaPrependPath[0].Path)
		}
		if *global.LuaPrependPath[1].Path != "/usr/share/haproxy-lua/?.lua" {
			t.Errorf("LuaPrependPath.Path is %v, expected /usr/share/haproxy-lua/?.lua", global.LuaPrependPath[1].Path)
		}
		typ := "cpath"
		if global.LuaPrependPath[1].Type != typ {
			t.Errorf("LuaPrependPath.Type is %v, expected cpath", global.LuaPrependPath[1].Type)
		}
	} else {
		t.Errorf("%v LuaPrependPath returned, expected 2", len(global.LuaPrependPath))
	}
	if len(global.LuaLoads) == 2 {
		if *global.LuaLoads[0].File != "/etc/foo.lua" {
			t.Errorf("LuaLoad.File is %v, expected /etc/foo.lua", *global.LuaLoads[0].File)
		}
		if *global.LuaLoads[1].File != "/etc/bar.lua" {
			t.Errorf("LuaLoad.File is %v, expected /etc/bar.lua", global.LuaLoads[1].File)
		}
	} else {
		t.Errorf("%v LuaLoads returned, expected 2", len(global.LuaLoads))
	}
	if len(global.H1CaseAdjusts) == 2 {
		if *global.H1CaseAdjusts[0].From != "host" {
			t.Errorf("H1CaseAdjusts[0].From is %v, expected host", *global.H1CaseAdjusts[0].From)
		}
		if *global.H1CaseAdjusts[0].To != "Host" {
			t.Errorf("H1CaseAdjusts[0].To is %v, expected Host", *global.H1CaseAdjusts[0].To)
		}
		if *global.H1CaseAdjusts[1].From != "content-type" {
			t.Errorf("H1CaseAdjusts[1].From is %v, expected content-type", *global.H1CaseAdjusts[1].From)
		}
		if *global.H1CaseAdjusts[1].To != "Content-Type" {
			t.Errorf("H1CaseAdjusts[1].To is %v, expected Content-Type", *global.H1CaseAdjusts[1].To)
		}
	}
	if global.H1CaseAdjustFile != "/etc/headers.adjust" {
		t.Errorf("H1CaseAdjustFile is %v, expected /etc/headers.adjust", global.H1CaseAdjustFile)
	}
	if global.UID != 1 {
		t.Errorf("UID is %v, expected 1", global.UID)
	}
	if global.Gid != 1 {
		t.Errorf("Gid is %v, expected 1", global.Gid)
	}
	if *global.TuneOptions.BuffersLimit != 11 {
		t.Errorf("BuffersLimit is %v, expected 11", global.TuneOptions.BuffersLimit)
	}
	if global.TuneOptions.BuffersReserve != 12 {
		t.Errorf("BuffersReserve is %v, expected 12", global.TuneOptions.BuffersReserve)
	}
	if global.TuneOptions.Bufsize != 13 {
		t.Errorf("Bufsize is %v, expected 13", global.TuneOptions.Bufsize)
	}
	if global.TuneOptions.CompMaxlevel != 14 {
		t.Errorf("CompMaxlevel is %v, expected 14", global.TuneOptions.CompMaxlevel)
	}
	if global.TuneOptions.H2HeaderTableSize != 15 {
		t.Errorf("H2HeaderTableSize is %v, expected 15", global.TuneOptions.H2HeaderTableSize)
	}
	if *global.TuneOptions.H2InitialWindowSize != 16 {
		t.Errorf("H2InitialWindowSize is %v, expected 16", global.TuneOptions.H2InitialWindowSize)
	}
	if global.TuneOptions.H2MaxConcurrentStreams != 17 {
		t.Errorf("H2MaxConcurrentStreams is %v, expected 17", global.TuneOptions.H2MaxConcurrentStreams)
	}
	if global.TuneOptions.H2MaxFrameSize != 18 {
		t.Errorf("H2MaxFrameSize is %v, expected 18", global.TuneOptions.H2MaxFrameSize)
	}
	if global.TuneOptions.HTTPCookielen != 19 {
		t.Errorf("HTTPCookielen is %v, expected 19", global.TuneOptions.HTTPCookielen)
	}
	if global.TuneOptions.HTTPLogurilen != 20 {
		t.Errorf("HTTPLogurilen is %v, expected 20", global.TuneOptions.HTTPLogurilen)
	}
	if global.TuneOptions.HTTPMaxhdr != 21 {
		t.Errorf("HTTPMaxhdr is %v, expected 21", global.TuneOptions.HTTPMaxhdr)
	}
	if *global.TuneOptions.Idletimer != 22 {
		t.Errorf("Idletimer is %v, expected 22", global.TuneOptions.Idletimer)
	}
	if global.TuneOptions.LuaForcedYield != 23 {
		t.Errorf("LuaForcedYield is %v, expected 23", global.TuneOptions.LuaForcedYield)
	}
	if global.TuneOptions.LuaMaxmem != true {
		t.Errorf("Maxzlibmem is false, expected true")
	}
	if *global.TuneOptions.LuaSessionTimeout != 25 {
		t.Errorf("LuaSessionTimeout is %v, expected 25", global.TuneOptions.LuaSessionTimeout)
	}
	if *global.TuneOptions.LuaTaskTimeout != 26 {
		t.Errorf("LuaTaskTimeout is %v, expected 26", global.TuneOptions.LuaTaskTimeout)
	}
	if *global.TuneOptions.LuaServiceTimeout != 27 {
		t.Errorf("LuaServiceTimeout is %v, expected 27", global.TuneOptions.LuaServiceTimeout)
	}
	if global.TuneOptions.Maxaccept != 28 {
		t.Errorf("Maxaccept is %v, expected 28", global.TuneOptions.Maxaccept)
	}
	if global.TuneOptions.Maxpollevents != 29 {
		t.Errorf("Maxpollevents is %v, expected 29", global.TuneOptions.Maxpollevents)
	}
	if global.TuneOptions.Maxrewrite != 30 {
		t.Errorf("Maxrewrite is %v, expected 30", global.TuneOptions.Maxrewrite)
	}
	if *global.TuneOptions.PatternCacheSize != 31 {
		t.Errorf("PatternCacheSize is %v, expected 31", global.TuneOptions.PatternCacheSize)
	}
	if global.TuneOptions.Pipesize != 32 {
		t.Errorf("Pipesize is %v, expected 32", global.TuneOptions.Pipesize)
	}
	if global.TuneOptions.PoolHighFdRatio != 33 {
		t.Errorf("PoolHighFdRatio is %v, expected 33", global.TuneOptions.PoolHighFdRatio)
	}
	if global.TuneOptions.PoolLowFdRatio != 34 {
		t.Errorf("PoolLowFdRatio is %v, expected 34", global.TuneOptions.PoolLowFdRatio)
	}
	if *global.TuneOptions.RcvbufClient != 35 {
		t.Errorf("RcvbufClient is %v, expected 35", global.TuneOptions.RcvbufClient)
	}
	if *global.TuneOptions.RcvbufServer != 36 {
		t.Errorf("RcvbufServer is %v, expected 36", global.TuneOptions.RcvbufServer)
	}
	if global.TuneOptions.RecvEnough != 37 {
		t.Errorf("RecvEnough is %v, expected 37", global.TuneOptions.RecvEnough)
	}
	if global.TuneOptions.RunqueueDepth != 38 {
		t.Errorf("RunqueueDepth is %v, expected 38", global.TuneOptions.RunqueueDepth)
	}
	if *global.TuneOptions.SndbufClient != 39 {
		t.Errorf("SndbufClient is %v, expected 39", global.TuneOptions.SndbufClient)
	}
	if *global.TuneOptions.SndbufServer != 40 {
		t.Errorf("SndbufServer is %v, expected 40", global.TuneOptions.SndbufServer)
	}
	if *global.TuneOptions.SslCachesize != 41 {
		t.Errorf("SslCachesize is %v, expected 41", global.TuneOptions.SslCachesize)
	}
	if global.TuneOptions.SslKeylog != "enabled" {
		t.Errorf("SslKeylog is %v, expected enabled", global.TuneOptions.SslKeylog)
	}
	if *global.TuneOptions.SslLifetime != 43 {
		t.Errorf("SslLifetime is %v, expected 43", global.TuneOptions.SslLifetime)
	}
	if *global.TuneOptions.SslMaxrecord != 44 {
		t.Errorf("SslMaxrecord is %v, expected 44", global.TuneOptions.SslMaxrecord)
	}
	if global.TuneOptions.SslDefaultDhParam != 45 {
		t.Errorf("SslDefaultDhParam is %v, expected 45", global.TuneOptions.SslDefaultDhParam)
	}
	if global.TuneOptions.SslCtxCacheSize != 46 {
		t.Errorf("SslCtxCacheSize is %v, expected 46", global.TuneOptions.SslCtxCacheSize)
	}
	if *global.TuneOptions.SslCaptureBufferSize != 47 {
		t.Errorf("SslCaptureBufferSize is %v, expected 47", global.TuneOptions.SslCaptureBufferSize)
	}
	if *global.TuneOptions.VarsGlobalMaxSize != 49 {
		t.Errorf("VarsGlobalMaxSize is %v, expected 49", global.TuneOptions.VarsGlobalMaxSize)
	}
	if *global.TuneOptions.VarsProcMaxSize != 50 {
		t.Errorf("VarsProcMaxSize is %v, expected 50", global.TuneOptions.VarsProcMaxSize)
	}
	if *global.TuneOptions.VarsReqresMaxSize != 51 {
		t.Errorf("VarsReqresMaxSize is %v, expected 51", global.TuneOptions.VarsReqresMaxSize)
	}
	if *global.TuneOptions.VarsSessMaxSize != 52 {
		t.Errorf("VarsSessMaxSize is %v, expected 52", global.TuneOptions.VarsSessMaxSize)
	}
	if *global.TuneOptions.VarsTxnMaxSize != 53 {
		t.Errorf("VarsTxnMaxSize is %v, expected 53", global.TuneOptions.VarsTxnMaxSize)
	}
	if *&global.TuneOptions.ZlibMemlevel != 54 {
		t.Errorf("ZlibMemlevel is %v, expected 54", global.TuneOptions.ZlibMemlevel)
	}
	if *&global.TuneOptions.ZlibWindowsize != 55 {
		t.Errorf("ZlibWindowsize is %v, expected 55", global.TuneOptions.ZlibWindowsize)
	}
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
		Maxconn:               1000,
		SslDefaultBindCiphers: "test",
		SslDefaultBindOptions: "ssl-min-ver TLSv1.0 no-tls-tickets",
		StatsTimeout:          &tOut,
		TuneSslDefaultDhParam: 1024,
		ExternalCheck:         false,
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
		TuneOptions: &models.GlobalTuneOptions{},
		UID:         1234,
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

	if !reflect.DeepEqual(global, g) {
		fmt.Printf("Created global config: %v\n", global)
		fmt.Printf("Given global config: %v\n", g)
		t.Error("Created global config not equal to given global config")
	}

	if ver != version {
		t.Error("Version not incremented!")
	}

	err = clientTest.PushGlobalConfiguration(g, "", 55)

	if err == nil {
		t.Error("Should have returned version conflict.")
	}
}
