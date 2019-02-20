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
	if global.Data.RuntimeAPI != "/var/run/haproxy.sock" {
		t.Errorf("RuntimeAPI is %v, expected /var/run/haproxy.sock", global.Data.RuntimeAPI)
	}
	if global.Data.RuntimeAPILevel != "admin" {
		t.Errorf("RuntimeAPILevel is %v, expected admin", global.Data.RuntimeAPILevel)
	}
	if global.Data.RuntimeAPIMode != "0660" {
		t.Errorf("RuntimeAPIMode is %v, expected 0660", global.Data.RuntimeAPIMode)
	}
	if global.Data.Nbproc != 4 {
		t.Errorf("Nbproc is %v, expected 4", global.Data.Nbproc)
	}
	if global.Data.Maxconn != 2000 {
		t.Errorf("Maxconn is %v, expected 2000", global.Data.Maxconn)
	}
}

func TestPutGlobal(t *testing.T) {
	g := &models.Global{
		Daemon:                "enabled",
		RuntimeAPI:            "/var/run/haproxy.sock",
		RuntimeAPILevel:       "admin",
		Maxconn:               1000,
		SslDefaultBindCiphers: "test",
		SslDefaultBindOptions: "ssl-min-ver TLSv1.0 no-tls-tickets",
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
