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
}

func TestPutGlobal(t *testing.T) {
	tOut := int64(3600)
	n := "1/1"
	v := "0"
	a := "/var/run/haproxy.sock"
	f := "/etc/foo.lua"
	enabled := "enabled"
	g := &models.Global{
		Daemon: "enabled",
		CPUMaps: []*models.CPUMap{
			&models.CPUMap{
				Process: &n,
				CPUSet:  &v,
			},
		},
		RuntimeAPIs: []*models.RuntimeAPI{
			&models.RuntimeAPI{
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
		LuaLoads: []*models.LuaLoad{
			&models.LuaLoad{
				File: &f,
			},
		},
		LogSendHostname: &models.GlobalLogSendHostname{
			Enabled: &enabled,
			Param:   "something",
		},
		SslModeAsync: "disabled",
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
