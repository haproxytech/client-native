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
	v, global, err := client.GetGlobalConfiguration("")
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
				Level:   "admin",
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
	}

	err := client.PushGlobalConfiguration(g, "", version)

	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	ver, global, err := client.GetGlobalConfiguration("")
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

	err = client.PushGlobalConfiguration(g, "", 55)

	if err == nil {
		t.Error("Should have returned version conflict.")
	}
}
