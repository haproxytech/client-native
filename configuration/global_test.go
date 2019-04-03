package configuration

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/haproxytech/models"
)

func TestGetGlobal(t *testing.T) {
	global, err := client.GetGlobalConfiguration("")
	if err != nil {
		t.Error(err.Error())
	}

	if global.Version != version {
		t.Errorf("Version %v returned, expected %v", global.Version, version)
	}

	if global.Data.Daemon != "enabled" {
		t.Errorf("Daemon is %v, expected enabled", global.Data.Daemon)
	}
	if len(global.Data.RuntimeApis) == 1 {
		if *global.Data.RuntimeApis[0].Address != "/var/run/haproxy.sock" {
			t.Errorf("RuntimeAPI.Address is %v, expected /var/run/haproxy.sock", *global.Data.RuntimeApis[0].Address)
		}
		if global.Data.RuntimeApis[0].Level != "admin" {
			t.Errorf("RuntimeAPI.Level is %v, expected admin", global.Data.RuntimeApis[0].Level)
		}
		if global.Data.RuntimeApis[0].Mode != "0660" {
			t.Errorf("RuntimeAPI.Mode is %v, expected 0660", global.Data.RuntimeApis[0].Mode)
		}
	} else {
		t.Errorf("RuntimeAPI is not set")
	}
	if global.Data.Nbproc != 4 {
		t.Errorf("Nbproc is %v, expected 4", global.Data.Nbproc)
	}
	if global.Data.Maxconn != 2000 {
		t.Errorf("Maxconn is %v, expected 2000", global.Data.Maxconn)
	}
}

func TestPutGlobal(t *testing.T) {
	tOut := int64(3600)
	n := "1/1"
	v := "0"
	a := "/var/run/haproxy.sock"
	g := &models.Global{
		Daemon: "enabled",
		CPUMaps: []*models.GlobalCPUMapsItems{
			&models.GlobalCPUMapsItems{
				Process: &n,
				CPUSet:  &v,
			},
		},
		RuntimeApis: []*models.GlobalRuntimeApisItems{
			&models.GlobalRuntimeApisItems{
				Address: &a,
				Level:   "admin",
			},
		},
		Maxconn:               1000,
		SslDefaultBindCiphers: "test",
		SslDefaultBindOptions: "ssl-min-ver TLSv1.0 no-tls-tickets",
		StatsTimeout:          &tOut,
		TuneSslDefaultDhParam: 1024,
	}

	err := client.PushGlobalConfiguration(g, "", version)

	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	global, err := client.GetGlobalConfiguration("")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(global.Data, g) {
		fmt.Printf("Created global config: %v\n", global.Data)
		fmt.Printf("Given global config: %v\n", g)
		t.Error("Created global config not equal to given global config")
	}

	if global.Version != version {
		t.Error("Version not incremented!")
	}

	err = client.PushGlobalConfiguration(g, "", 55)

	if err == nil {
		t.Error("Should have returned version conflict.")
	}
}
