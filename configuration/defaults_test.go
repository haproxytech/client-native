package configuration

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/haproxytech/client-native/v2/models"
)

func TestGetDefaults(t *testing.T) { //nolint:gocognit,gocyclo
	v, d, err := client.GetDefaultsConfiguration("")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	if *d.Balance.Algorithm != "roundrobin" {
		t.Errorf("Balance.Algorithm not roundrobin: %v", d.Balance.Algorithm)
	}
	if d.Mode != "http" {
		t.Errorf("Mode not http: %v", d.Mode)
	}
	if d.MonitorURI != "/monitor" {
		t.Errorf("MonitorURI not /monitor: %v", d.MonitorURI)
	}
	if d.BindProcess != "1-4" {
		t.Errorf("BindProcess not 1-4: %v", d.BindProcess)
	}
	if d.Clitcpka != "enabled" {
		t.Errorf("Clitcpka not enabled: %v", d.Clitcpka)
	}
	if d.Dontlognull != "enabled" {
		t.Errorf("Dontlognull not enabled: %v", d.Dontlognull)
	}
	if d.HTTPUseHtx != "enabled" {
		t.Errorf("HTTPUseHtx not enabled: %v", d.HTTPUseHtx)
	}
	if !d.Httplog {
		t.Errorf("Httplog not enabled: %v", d.Httplog)
	}
	if d.HTTPConnectionMode != "httpclose" {
		t.Errorf("HTTPConnectionMode not httpclose: %v", d.HTTPConnectionMode)
	}
	if d.DefaultBackend != "test" {
		t.Errorf("DefaultBackend not test: %v", d.DefaultBackend)
	}
	if *d.Maxconn != 2000 {
		t.Errorf("Maxconn not 2000: %v", *d.Maxconn)
	}
	if *d.ClientTimeout != 4000 {
		t.Errorf("ClientTimeout not 4000: %v", *d.ClientTimeout)
	}
	if *d.CheckTimeout != 2000 {
		t.Errorf("CheckTimeout not 2000: %v", *d.CheckTimeout)
	}
	if *d.ConnectTimeout != 5000 {
		t.Errorf("ConnectTimeout not 5000: %v", *d.ConnectTimeout)
	}
	if *d.QueueTimeout != 900 {
		t.Errorf("QueueTimeout not 900: %v", *d.QueueTimeout)
	}
	if *d.ServerTimeout != 2000 {
		t.Errorf("ServerTimeout not 2000: %v", *d.ServerTimeout)
	}
	if *d.HTTPRequestTimeout != 2000 {
		t.Errorf("HTTPRequestTimeout not 2000: %v", *d.HTTPRequestTimeout)
	}
	if *d.HTTPKeepAliveTimeout != 3000 {
		t.Errorf("HTTPKeepAliveTimeout not 3000: %v", *d.HTTPKeepAliveTimeout)
	}
	if *d.DefaultServer.Fall != 2000 {
		t.Errorf("DefaultServer.Fall not 2000: %v", *d.DefaultServer.Fall)
	}
	if *d.DefaultServer.Rise != 4000 {
		t.Errorf("DefaultServer.Rise not 4000: %v", *d.DefaultServer.Rise)
	}
	if *d.DefaultServer.Inter != 5000 {
		t.Errorf("DefaultServer.Inter not 5000: %v", *d.DefaultServer.Inter)
	}
	if *d.DefaultServer.Port != 8888 {
		t.Errorf("DefaultServer.Port not 8888: %v", *d.DefaultServer.Port)
	}
	if len(d.ErrorFiles) != 3 {
		t.Errorf("ErrorFiles not 3: %v", len(d.ErrorFiles))
	} else {
		for _, ef := range d.ErrorFiles {
			if ef.Code == 403 {
				if ef.File != "/test/403.html" {
					t.Errorf("File for %v not 403: %v", ef.Code, ef.File)
				}
			}
			if ef.Code == 500 {
				if ef.File != "/test/500.html" {
					t.Errorf("File for %v not 500: %v", ef.Code, ef.File)
				}
			}
			if ef.Code == 429 {
				if ef.File != "/test/429.html" {
					t.Errorf("File for %v not 429: %v", ef.Code, ef.File)
				}
			}
		}
	}
	if d.ExternalCheck != "enabled" {
		t.Errorf("ExternalCheck not enabled: %v", d.ExternalCheck)
	}
	if d.ExternalCheckPath != "/bin" {
		t.Errorf("ExternalCheckPath not /bin: %v", d.ExternalCheckPath)
	}
	if d.ExternalCheckCommand != "/bin/true" {
		t.Errorf("ExternalCheckCommand not /bin/true: %v", d.ExternalCheckCommand)
	}
}

func TestPushDefaults(t *testing.T) {
	tOut := int64(6000)
	tOutS := int64(200)
	balanceAlgorithm := "leastconn"
	d := &models.Defaults{
		Clitcpka:       "disabled",
		BindProcess:    "1-4",
		DefaultBackend: "test2",
		ErrorFiles: []*models.Errorfile{
			&models.Errorfile{
				Code: 400,
				File: "/test/400.html",
			},
			&models.Errorfile{
				Code: 403,
				File: "/test/403.html",
			},
			&models.Errorfile{
				Code: 429,
				File: "/test/429.html",
			},
			&models.Errorfile{
				Code: 500,
				File: "/test/500.html",
			},
		},
		CheckTimeout:   &tOutS,
		ConnectTimeout: &tOut,
		ServerTimeout:  &tOutS,
		QueueTimeout:   &tOutS,
		Mode:           "tcp",
		MonitorURI:     "/healthz",
		HTTPUseHtx:     "enabled",
		Balance: &models.Balance{
			Algorithm: &balanceAlgorithm,
		},
		ExternalCheck:        "",
		ExternalCheckPath:    "/bin",
		ExternalCheckCommand: "/bin/flase",
		Logasap:              "disabled",
		Allbackups:           "enabled",
	}

	err := client.PushDefaultsConfiguration(d, "", version)

	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	ver, defaults, err := client.GetDefaultsConfiguration("")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(defaults, d) {
		fmt.Printf("Created defaults config: %v\n", defaults)
		fmt.Printf("Given defaults config: %v\n", d)
		t.Error("Created defaults config not equal to given defaults config")
	}

	if ver != version {
		t.Error("Version not incremented!")
	}

	err = client.PushDefaultsConfiguration(d, "", 1055)

	if err == nil {
		t.Error("Should have returned version conflict.")
	}
}
