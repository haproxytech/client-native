package configuration

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/haproxytech/models"
)

func TestGetBackends(t *testing.T) {
	backends, err := client.GetBackends()
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
		if b.AdvCheck != "http" {
			t.Errorf("%v: Adv check not http: %v", b.Name, b.AdvCheck)
		}
		if b.Protocol != "http" {
			t.Errorf("%v: Protocol not http: %v", b.Name, b.Protocol)
		}
		if b.Balance != "roundrobin" {
			t.Errorf("%v: Balance not roundrobin: %v", b.Name, b.Balance)
		}
		if b.Log != "enabled" {
			t.Errorf("%v: Log not enabled: %v", b.Name, b.Log)
		}
		if b.LogFormat != "http" {
			t.Errorf("%v: LogFormat not http: %v", b.Name, b.LogFormat)
		}
		if b.HTTPConnectionMode != "keep-alive" {
			t.Errorf("%v: HTTPConnectionMode not keep-alive: %v", b.Name, b.HTTPConnectionMode)
		}
		if b.HTTPXffHeaderInsert != "enabled" {
			t.Errorf("%v: HTTPXffHeaderInsert not enabled: %v", b.Name, b.HTTPXffHeaderInsert)
		}
		if b.AdvCheckHTTPMethod != "HEAD" {
			t.Errorf("%v: AdvCheckHTTPMethod not HEAD: %v", b.Name, b.AdvCheckHTTPMethod)
		}
		if b.AdvCheckHTTPURI != "/" {
			t.Errorf("%v: AdvCheckHTTPURI not /: %v", b.Name, b.AdvCheckHTTPURI)
		}
		if *b.CheckFall != 2 {
			t.Errorf("%v: CheckFall not 2: %v", b.Name, *b.CheckFall)
		}
		if *b.CheckRise != 4 {
			t.Errorf("%v: CheckRise not 4: %v", b.Name, *b.CheckRise)
		}
		if *b.CheckInterval != 5 {
			t.Errorf("%v: CheckInterval not 5: %v", b.Name, *b.CheckInterval)
		}
		if *b.CheckPort != 8888 {
			t.Errorf("%v: CheckPort not 8888: %v", b.Name, *b.CheckPort)
		}
		if b.ContinuousStatistics != "enabled" {
			t.Errorf("%v: ContinuousStatistics not enabled: %v", b.Name, b.ContinuousStatistics)
		}
		if b.HTTPCookie != "enabled" {
			t.Errorf("%v: HTTPCookie not enabled: %v", b.Name, b.HTTPCookie)
		}
		if b.HTTPCookieName != "BLA" {
			t.Errorf("%v: HTTPCookieName not BLA: %v", b.Name, b.HTTPCookieName)
		}
		if *b.CheckTimeout != 2 {
			t.Errorf("%v: CheckTimeout not 2: %v", b.Name, *b.CheckTimeout)
		}
		if *b.ServerInactivityTimeout != 3 {
			t.Errorf("%v: ServerInactivityTimeout not 3: %v", b.Name, *b.ServerInactivityTimeout)
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
	backend, err := client.GetBackend("test")
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
	if b.AdvCheck != "http" {
		t.Errorf("%v: Adv check not http: %v", b.Name, b.AdvCheck)
	}
	if b.Protocol != "http" {
		t.Errorf("%v: Protocol not http: %v", b.Name, b.Protocol)
	}
	if b.Balance != "roundrobin" {
		t.Errorf("%v: Balance not roundrobin: %v", b.Name, b.Balance)
	}
	if b.Log != "enabled" {
		t.Errorf("%v: Log not enabled: %v", b.Name, b.Log)
	}
	if b.LogFormat != "http" {
		t.Errorf("%v: LogFormat not http: %v", b.Name, b.LogFormat)
	}
	if b.HTTPConnectionMode != "keep-alive" {
		t.Errorf("%v: HTTPConnectionMode not keep-alive: %v", b.Name, b.HTTPConnectionMode)
	}
	if b.HTTPXffHeaderInsert != "enabled" {
		t.Errorf("%v: HTTPXffHeaderInsert not enabled: %v", b.Name, b.HTTPXffHeaderInsert)
	}
	if b.AdvCheckHTTPMethod != "HEAD" {
		t.Errorf("%v: AdvCheckHTTPMethod not HEAD: %v", b.Name, b.AdvCheckHTTPMethod)
	}
	if b.AdvCheckHTTPURI != "/" {
		t.Errorf("%v: AdvCheckHTTPURI not /: %v", b.Name, b.AdvCheckHTTPURI)
	}
	if *b.CheckFall != 2 {
		t.Errorf("%v: CheckFall not 2: %v", b.Name, *b.CheckFall)
	}
	if *b.CheckRise != 4 {
		t.Errorf("%v: CheckRise not 4: %v", b.Name, *b.CheckRise)
	}
	if *b.CheckInterval != 5 {
		t.Errorf("%v: CheckInterval not 5: %v", b.Name, *b.CheckInterval)
	}
	if *b.CheckPort != 8888 {
		t.Errorf("%v: CheckPort not 8888: %v", b.Name, *b.CheckPort)
	}
	if b.ContinuousStatistics != "enabled" {
		t.Errorf("%v: ContinuousStatistics not enabled: %v", b.Name, b.ContinuousStatistics)
	}
	if b.HTTPCookie != "enabled" {
		t.Errorf("%v: HTTPCookie not enabled: %v", b.Name, b.HTTPCookie)
	}
	if b.HTTPCookieName != "BLA" {
		t.Errorf("%v: HTTPCookieName not BLA: %v", b.Name, b.HTTPCookieName)
	}
	if *b.CheckTimeout != 2 {
		t.Errorf("%v: CheckTimeout not 2: %v", b.Name, *b.CheckTimeout)
	}
	if *b.ServerInactivityTimeout != 3 {
		t.Errorf("%v: ServerInactivityTimeout not 3: %v", b.Name, *b.ServerInactivityTimeout)
	}

	bJSON, err := b.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	_, err = client.GetBackend("doesnotexist")
	if err == nil {
		t.Error("Should throw error, non existant bck")
	}

	if !t.Failed() {
		fmt.Println("GetBackend succesful\nResponse: \n" + string(bJSON) + "\n")
	}
}

func TestCreateBackend(t *testing.T) {
	tOut := int64(5)
	b := &models.Backend{
		Name:               "created",
		Protocol:           "http",
		Balance:            "hash-uri",
		HTTPConnectionMode: "keep-alive",
		ConnectTimeout:     &tOut,
	}

	err := client.CreateBackend(b, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version = version + 1
	}

	backend, err := client.GetBackend("created")
	if err != nil {
		t.Error(err.Error())
	}

	bCreated := backend.Data

	if !reflect.DeepEqual(bCreated, b) {
		fmt.Printf("Created bck: %v\n", bCreated)
		fmt.Printf("Given bck: %v\n", b)
		t.Error("Created backend not equal to given backend")
	}

	if backend.Version != version {
		t.Errorf("Version %v returned, expected %v", backend.Version, version)
	}

	err = client.CreateBackend(b, "", version)
	if err == nil {
		t.Error("Should throw error bck already exists")
		version = version + 1
	}

	if !t.Failed() {
		fmt.Println("CreateBackend successful")
	}
}

func TestEditBackend(t *testing.T) {
	tOut := int64(3)
	b := &models.Backend{
		Name:               "created",
		Protocol:           "http",
		Balance:            "roundrobin",
		HTTPConnectionMode: "tunnel",
		ConnectTimeout:     &tOut,
	}

	err := client.EditBackend("created", b, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version = version + 1
	}

	backend, err := client.GetBackend("created")
	if err != nil {
		t.Error(err.Error())
	}
	bEdited := backend.Data

	if !reflect.DeepEqual(bEdited, b) {
		fmt.Printf("Edited bck: %v\n", bEdited)
		fmt.Printf("Given bck: %v\n", b)
		t.Error("Edited backend not equal to given backend")
	}

	if backend.Version != version {
		t.Errorf("Version %v returned, expected %v", backend.Version, version)
	}

	if !t.Failed() {
		fmt.Println("EditBackend successful")
	}
}

func TestDeleteBackend(t *testing.T) {
	err := client.DeleteBackend("test_2", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version = version + 1
	}

	if v, _ := client.GetVersion(); v != version {
		t.Error("Version not incremented")
	}

	err = client.DeleteBackend("test_2", "", 999999999)
	if err != nil {
		if err != ErrVersionMismatch {
			t.Error("DeleteBackend failed, should return version mismatch")
		}
	}

	_, err = client.GetBackend("test_2")
	if err == nil {
		t.Error("DeleteBackend failed, bck test still exists")
	}

	err = client.DeleteBackend("doesnotexist", "", version)
	if err == nil {
		t.Error("Should throw error, non existant bck")
		version = version + 1
	}

	if !t.Failed() {
		fmt.Println("DeleteBackend successful")
	}
}
