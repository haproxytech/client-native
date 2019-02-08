package configuration

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/haproxytech/models"
)

func TestGetFrontends(t *testing.T) {
	frontends, err := client.GetFrontends("")
	if err != nil {
		t.Error(err.Error())
	}

	if len(frontends.Data) != 2 {
		t.Errorf("%v frontends returned, expected 2", len(frontends.Data))
	}

	if frontends.Version != version {
		t.Errorf("Version %v returned, expected %v", frontends.Version, version)
	}

	for _, f := range frontends.Data {
		if f.Name != "test" && f.Name != "test_2" {
			t.Errorf("Expected only test or test_2 frontend, %v found", f.Name)
		}
		if f.Protocol != "http" {
			t.Errorf("%v: Protocol not http: %v", f.Name, f.Protocol)
		}
		if f.Log != "enabled" {
			t.Errorf("%v: Log not enabled: %v", f.Name, f.Log)
		}
		if f.LogFormat != "http" {
			t.Errorf("%v: LogFormat not http: %v", f.Name, f.LogFormat)
		}
		if f.LogIgnoreNull != "enabled" {
			t.Errorf("%v: LogIgnoreNull not enabled: %v", f.Name, f.LogIgnoreNull)
		}
		if f.HTTPConnectionMode != "passive-close" {
			t.Errorf("%v: HTTPConnectionMode not passive-close: %v", f.Name, f.HTTPConnectionMode)
		}
		if f.ContinuousStatistics != "enabled" {
			t.Errorf("%v: ContinuousStatistics not enabled: %v", f.Name, f.ContinuousStatistics)
		}
		if *f.HTTPRequestTimeout != 2 {
			t.Errorf("%v: HTTPRequestTimeout not 2: %v", f.Name, *f.HTTPRequestTimeout)
		}
		if *f.HTTPKeepaliveTimeout != 3 {
			t.Errorf("%v: HTTPKeepaliveTimeout not 3: %v", f.Name, *f.HTTPKeepaliveTimeout)
		}
		if f.DefaultFarm != "test" && f.DefaultFarm != "test_2" {
			t.Errorf("%v: DefaultFarm not test or test_2: %v", f.Name, f.DefaultFarm)
		}
		if *f.MaxConnections != 2000 {
			t.Errorf("%v: MaxConnections not 2000: %v", f.Name, *f.MaxConnections)
		}
		if *f.ClientInactivityTimeout != 4 {
			t.Errorf("%v: ClientInactivityTimeout not 4: %v", f.Name, *f.ClientInactivityTimeout)
		}
	}

	fJSON, err := frontends.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	if !t.Failed() {
		fmt.Println("GetFrontends succesful\nResponse: \n" + string(fJSON) + "\n")
	}
}

func TestGetFrontend(t *testing.T) {
	frontend, err := client.GetFrontend("test", "")
	if err != nil {
		t.Error(err.Error())
	}

	f := frontend.Data

	if frontend.Version != version {
		t.Errorf("Version %v returned, expected %v", frontend.Version, version)
	}

	if f.Name != "test" {
		t.Errorf("Expected only test, %v found", f.Name)
	}
	if f.Name != "test" {
		t.Errorf("Expected only test, %v found", f.Name)
	}
	if f.Protocol != "http" {
		t.Errorf("%v: Protocol not http: %v", f.Name, f.Protocol)
	}
	if f.Log != "enabled" {
		t.Errorf("%v: Log not enabled: %v", f.Name, f.Log)
	}
	if f.LogFormat != "http" {
		t.Errorf("%v: LogFormat not http: %v", f.Name, f.LogFormat)
	}
	if f.LogIgnoreNull != "enabled" {
		t.Errorf("%v: LogIgnoreNull not enabled: %v", f.Name, f.LogIgnoreNull)
	}
	if f.HTTPConnectionMode != "passive-close" {
		t.Errorf("%v: HTTPConnectionMode not passive-close: %v", f.Name, f.HTTPConnectionMode)
	}
	if f.ContinuousStatistics != "enabled" {
		t.Errorf("%v: ContinuousStatistics not enabled: %v", f.Name, f.ContinuousStatistics)
	}
	if *f.HTTPRequestTimeout != 2 {
		t.Errorf("%v: HTTPRequestTimeout not 2: %v", f.Name, *f.HTTPRequestTimeout)
	}
	if *f.HTTPKeepaliveTimeout != 3 {
		t.Errorf("%v: HTTPKeepaliveTimeout not 3: %v", f.Name, *f.HTTPKeepaliveTimeout)
	}
	if f.DefaultFarm != "test" {
		t.Errorf("%v: DefaultFarm not test: %v", f.Name, f.DefaultFarm)
	}
	if *f.MaxConnections != 2000 {
		t.Errorf("%v: MaxConnections not 2000: %v", f.Name, *f.MaxConnections)
	}
	if *f.ClientInactivityTimeout != 4 {
		t.Errorf("%v: ClientInactivityTimeout not 4: %v", f.Name, *f.ClientInactivityTimeout)
	}

	fJSON, err := f.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	_, err = client.GetFrontend("doesnotexist", "")
	if err == nil {
		t.Error("Should throw error, non existant frontend")
	}

	if !t.Failed() {
		fmt.Println("GetFrontend succesful\nResponse: \n" + string(fJSON) + "\n")
	}
}

func TestCreateEditDeleteFrontend(t *testing.T) {
	// TestCreateFrontend
	mConn := int64(3000)
	tOut := int64(2)
	f := &models.Frontend{
		Name:                 "created",
		Protocol:             "tcp",
		MaxConnections:       &mConn,
		HTTPConnectionMode:   "keep-alive",
		HTTPKeepaliveTimeout: &tOut,
	}

	err := client.CreateFrontend(f, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	frontend, err := client.GetFrontend("created", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(frontend.Data, f) {
		fmt.Printf("Created frontend: %v\n", frontend.Data)
		fmt.Printf("Given frontend: %v\n", f)
		t.Error("Created frontend not equal to given frontend")
	}

	if frontend.Version != version {
		t.Errorf("Version %v returned, expected %v", frontend.Version, version)
	}

	err = client.CreateFrontend(f, "", version)
	if err == nil {
		t.Error("Should throw error frontend already exists")
		version++
	}

	if !t.Failed() {
		fmt.Println("CreateFrontend successful")
	}

	// TestEditBackend
	mConn = int64(4000)
	f = &models.Frontend{
		Name:               "created",
		Protocol:           "tcp",
		MaxConnections:     &mConn,
		HTTPConnectionMode: "tunnel",
	}

	err = client.EditFrontend("created", f, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	frontend, err = client.GetFrontend("created", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(frontend.Data, f) {
		fmt.Printf("Edited frontend: %v\n", frontend.Data)
		fmt.Printf("Given frontend: %v\n", f)
		t.Error("Edited frontend not equal to given frontend")
	}

	if frontend.Version != version {
		t.Errorf("Version %v returned, expected %v", frontend.Version, version)
	}

	if !t.Failed() {
		fmt.Println("EditFrontend successful")
	}

	// TestDeleteBackend
	err = client.DeleteFrontend("created", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ := client.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	err = client.DeleteFrontend("created", "", 999999)
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
	_, err = client.GetFrontend("created", "")
	if err == nil {
		t.Error("DeleteFrontend failed, frontend test still exists")
	}

	err = client.DeleteFrontend("doesnotexist", "", version)
	if err == nil {
		t.Error("Should throw error, non existant frontend")
		version++
	}

	if !t.Failed() {
		fmt.Println("DeleteFrontend successful")
	}
}
