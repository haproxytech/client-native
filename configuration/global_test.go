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

	"github.com/haproxytech/models"
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
	if len(global.RuntimeApis) == 1 {
		if *global.RuntimeApis[0].Address != "/var/run/haproxy.sock" {
			t.Errorf("RuntimeAPI.Address is %v, expected /var/run/haproxy.sock", *global.RuntimeApis[0].Address)
		}
		if global.RuntimeApis[0].Level != "admin" {
			t.Errorf("RuntimeAPI.Level is %v, expected admin", global.RuntimeApis[0].Level)
		}
		if global.RuntimeApis[0].Mode != "0660" {
			t.Errorf("RuntimeAPI.Mode is %v, expected 0660", global.RuntimeApis[0].Mode)
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
}

func TestPutGlobal(t *testing.T) {
	tOut := int64(3600)
	n := "1/1"
	v := "0"
	a := "/var/run/haproxy.sock"
	g := &models.Global{
		Daemon: "enabled",
		CPUMaps: []*models.CPUMap{
			&models.CPUMap{
				Process: &n,
				CPUSet:  &v,
			},
		},
		RuntimeApis: []*models.RuntimeAPI{
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
