package configuration

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/haproxytech/models"
)

func TestGetBackends(t *testing.T) {
	backends, err := client.GetBackends("")
	if err != nil {
		t.Error(err.Error())
	}

	if len(backends.Data) != 2 {
		t.Errorf("%v backends returned, expected 2", len(backends.Data))
	}

	if backends.Version != version {
		t.Errorf("Version %v returned, expected %v", backends.Version, version)
	}

	for _, b := range backends.Data {
		if b.Name != "test" && b.Name != "test_2" {
			t.Errorf("Expected only test or test_2 backend, %v found", b.Name)
		}
		if b.Httpchk.Method != "HEAD" {
			t.Errorf("%v: Httpchk.Method not HEAD: %v", b.Name, b.Httpchk.Method)
		}
		if b.Httpchk.URI != "/" {
			t.Errorf("%v: Httpchk.URI not HEAD: %v", b.Name, b.Httpchk.URI)
		}
		if b.Mode != "http" {
			t.Errorf("%v: Mode not http: %v", b.Name, b.Mode)
		}
		if b.Balance.Algorithm != "roundrobin" {
			t.Errorf("%v: Balance.Algorithm not roundrobin: %v", b.Name, b.Balance.Algorithm)
		}
		if b.Log != true {
			t.Errorf("%v: Log not true: %v", b.Name, b.Log)
		}
		if b.LogFormat != "http" {
			t.Errorf("%v: LogFormat not http: %v", b.Name, b.LogFormat)
		}
		if b.HTTPConnectionMode != "http-keep-alive" {
			t.Errorf("%v: HTTPConnectionMode not http-keep-alive: %v", b.Name, b.HTTPConnectionMode)
		}
		if b.Forwardfor != "enabled" {
			t.Errorf("%v: Forwardfor not enabled: %v", b.Name, b.Forwardfor)
		}
		if *b.DefaultServer.Fall != 2000 {
			t.Errorf("%v: DefaultServer.Fall not 2000: %v", b.Name, *b.DefaultServer.Fall)
		}
		if *b.DefaultServer.Rise != 4000 {
			t.Errorf("%v: DefaultServer.Rise not 4000: %v", b.Name, *b.DefaultServer.Rise)
		}
		if *b.DefaultServer.Inter != 5000 {
			t.Errorf("%v: DefaultServer.Inter not 5000: %v", b.Name, *b.DefaultServer.Inter)
		}
		if *b.DefaultServer.Port != 8888 {
			t.Errorf("%v: DefaultServer.Port not 8888: %v", b.Name, *b.DefaultServer.Port)
		}
		if b.Contstats != "enabled" {
			t.Errorf("%v: ContinuousStatistics not enabled: %v", b.Name, b.Contstats)
		}
		if b.Cookie != "BLA" {
			t.Errorf("%v: HTTPCookieName not BLA: %v", b.Name, b.Cookie)
		}
		if *b.CheckTimeout != 2000 {
			t.Errorf("%v: CheckTimeout not 2000: %v", b.Name, *b.CheckTimeout)
		}
		if *b.ServerTimeout != 3000 {
			t.Errorf("%v: ServerTimeout not 3000: %v", b.Name, *b.ServerTimeout)
		}
	}

	bJSON, err := backends.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	if !t.Failed() {
		fmt.Println("GetBackends succesful\nResponse: \n" + string(bJSON) + "\n")
	}
}

func TestGetBackend(t *testing.T) {
	backend, err := client.GetBackend("test", "")
	if err != nil {
		t.Error(err.Error())
	}

	b := backend.Data

	if backend.Version != version {
		t.Errorf("Version %v returned, expected %v", backend.Version, version)
	}

	if b.Name != "test" {
		t.Errorf("Expected only test, %v found", b.Name)
	}
	if b.Httpchk.Method != "HEAD" {
		t.Errorf("%v: Httpchk.Method not HEAD: %v", b.Name, b.Httpchk.Method)
	}
	if b.Httpchk.URI != "/" {
		t.Errorf("%v: Httpchk.URI not HEAD: %v", b.Name, b.Httpchk.URI)
	}
	if b.Mode != "http" {
		t.Errorf("%v: Mode not http: %v", b.Name, b.Mode)
	}
	if b.Balance.Algorithm != "roundrobin" {
		t.Errorf("%v: Balance.Algorithm not roundrobin: %v", b.Name, b.Balance.Algorithm)
	}
	if b.Log != true {
		t.Errorf("%v: Log not true: %v", b.Name, b.Log)
	}
	if b.LogFormat != "http" {
		t.Errorf("%v: LogFormat not http: %v", b.Name, b.LogFormat)
	}
	if b.HTTPConnectionMode != "http-keep-alive" {
		t.Errorf("%v: HTTPConnectionMode not http-keep-alive: %v", b.Name, b.HTTPConnectionMode)
	}
	if b.Forwardfor != "enabled" {
		t.Errorf("%v: Forwardfor not enabled: %v", b.Name, b.Forwardfor)
	}
	if *b.DefaultServer.Fall != 2000 {
		t.Errorf("%v: DefaultServer.Fall not 2000: %v", b.Name, *b.DefaultServer)
	}
	if *b.DefaultServer.Rise != 4000 {
		t.Errorf("%v: DefaultServer.Rise not 4000: %v", b.Name, *b.DefaultServer.Rise)
	}
	if *b.DefaultServer.Inter != 5000 {
		t.Errorf("%v: DefaultServer.Inter not 5000: %v", b.Name, *b.DefaultServer.Inter)
	}
	if *b.DefaultServer.Port != 8888 {
		t.Errorf("%v: DefaultServer.Port not 8888: %v", b.Name, *b.DefaultServer.Port)
	}
	if b.Contstats != "enabled" {
		t.Errorf("%v: ContinuousStatistics not enabled: %v", b.Name, b.Contstats)
	}
	if b.Cookie != "BLA" {
		t.Errorf("%v: HTTPCookieName not BLA: %v", b.Name, b.Cookie)
	}
	if *b.CheckTimeout != 2000 {
		t.Errorf("%v: CheckTimeout not 2000: %v", b.Name, *b.CheckTimeout)
	}
	if *b.ServerTimeout != 3000 {
		t.Errorf("%v: ServerTimeout not 3000: %v", b.Name, *b.ServerTimeout)
	}

	bJSON, err := b.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	_, err = client.GetBackend("doesnotexist", "")
	if err == nil {
		t.Error("Should throw error, non existant bck")
	}

	if !t.Failed() {
		fmt.Println("GetBackend succesful\nResponse: \n" + string(bJSON) + "\n")
	}
}

func TestCreateEditDeleteBackend(t *testing.T) {
	// TestCreateBackend
	tOut := int64(5)
	b := &models.Backend{
		Name:               "created",
		Mode:               "http",
		Balance:            &models.BackendBalance{Algorithm: "uri"},
		HTTPConnectionMode: "http-keep-alive",
		ConnectTimeout:     &tOut,
	}

	err := client.CreateBackend(b, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	backend, err := client.GetBackend("created", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(backend.Data, b) {
		fmt.Printf("Created bck: %v\n", backend.Data)
		fmt.Printf("Given bck: %v\n", b)
		t.Error("Created backend not equal to given backend")
	}

	if backend.Version != version {
		t.Errorf("Version %v returned, expected %v", backend.Version, version)
	}

	err = client.CreateBackend(b, "", version)
	if err == nil {
		t.Error("Should throw error bck already exists")
		version++
	}

	if !t.Failed() {
		fmt.Println("CreateBackend successful")
	}

	// TestEditBackend
	tOut = int64(3)
	b = &models.Backend{
		Name:               "created",
		Mode:               "http",
		Balance:            &models.BackendBalance{Algorithm: "roundrobin"},
		HTTPConnectionMode: "http-tunnel",
		ConnectTimeout:     &tOut,
	}

	err = client.EditBackend("created", b, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	backend, err = client.GetBackend("created", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(backend.Data, b) {
		fmt.Printf("Edited bck: %v\n", backend.Data)
		fmt.Printf("Given bck: %v\n", b)
		t.Error("Edited backend not equal to given backend")
	}

	if backend.Version != version {
		t.Errorf("Version %v returned, expected %v", backend.Version, version)
	}

	if !t.Failed() {
		fmt.Println("EditBackend successful")
	}

	// TestDeleteBackend
	err = client.DeleteBackend("created", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ := client.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	err = client.DeleteBackend("created", "", 999999999)
	if err != nil {
		switch err.(type) {
		case *ConfError:
			if err.(*ConfError).Code() != ErrVersionMismatch {
				t.Error("Should throw ErrVersionMismatch error")
			}
		default:
			t.Error("Should throw ErrVersionMismatch error")
		}
	}

	_, err = client.GetBackend("created", "")
	if err == nil {
		t.Error("DeleteBackend failed, bck test still exists")
	}

	err = client.DeleteBackend("doesnotexist", "", version)
	if err == nil {
		t.Error("Should throw error, non existant bck")
		version++
	}

	if !t.Failed() {
		fmt.Println("DeleteBackend successful")
	}
}
