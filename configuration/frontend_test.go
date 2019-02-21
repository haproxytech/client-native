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
		if f.Mode != "http" {
			t.Errorf("%v: Mode not http: %v", f.Name, f.Mode)
		}
		if f.Dontlognull != "enabled" {
			t.Errorf("%v: Dontlognull not enabled: %v", f.Name, f.Dontlognull)
		}
		if f.HTTPConnectionMode != "httpclose" {
			t.Errorf("%v: HTTPConnectionMode not httpclose: %v", f.Name, f.HTTPConnectionMode)
		}
		if f.Contstats != "enabled" {
			t.Errorf("%v: Contstats not enabled: %v", f.Name, f.Contstats)
		}
		if *f.HTTPRequestTimeout != 2000 {
			t.Errorf("%v: HTTPRequestTimeout not 2: %v", f.Name, *f.HTTPRequestTimeout)
		}
		if *f.HTTPKeepAliveTimeout != 3000 {
			t.Errorf("%v: HTTPKeepAliveTimeout not 3: %v", f.Name, *f.HTTPKeepAliveTimeout)
		}
		if f.DefaultBackend != "test" && f.DefaultBackend != "test_2" {
			t.Errorf("%v: DefaultFarm not test or test_2: %v", f.Name, f.DefaultBackend)
		}
		if *f.Maxconn != 2000 {
			t.Errorf("%v: Maxconn not 2000: %v", f.Name, *f.Maxconn)
		}
		if *f.ClientTimeout != 4000 {
			t.Errorf("%v: ClientTimeout not 4: %v", f.Name, *f.ClientTimeout)
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
	if f.Mode != "http" {
		t.Errorf("%v: Mode not http: %v", f.Name, f.Mode)
	}
	if f.Dontlognull != "enabled" {
		t.Errorf("%v: Dontlognull not enabled: %v", f.Name, f.Dontlognull)
	}
	if f.HTTPConnectionMode != "httpclose" {
		t.Errorf("%v: HTTPConnectionMode not httpclose: %v", f.Name, f.HTTPConnectionMode)
	}
	if f.Contstats != "enabled" {
		t.Errorf("%v: Contstats not enabled: %v", f.Name, f.Contstats)
	}
	if *f.HTTPRequestTimeout != 2000 {
		t.Errorf("%v: HTTPRequestTimeout not 2000: %v", f.Name, *f.HTTPRequestTimeout)
	}
	if *f.HTTPKeepAliveTimeout != 3000 {
		t.Errorf("%v: HTTPKeepAliveTimeout not 3000: %v", f.Name, *f.HTTPKeepAliveTimeout)
	}
	if f.DefaultBackend != "test" {
		t.Errorf("%v: DefaultBackend not test: %v", f.Name, f.DefaultBackend)
	}
	if *f.Maxconn != 2000 {
		t.Errorf("%v: Maxconn not 2000: %v", f.Name, *f.Maxconn)
	}
	if *f.ClientTimeout != 4000 {
		t.Errorf("%v: ClientTimeout not 4000: %v", f.Name, *f.ClientTimeout)
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
		Mode:                 "tcp",
		Maxconn:              &mConn,
		HTTPConnectionMode:   "http-keep-alive",
		HTTPKeepAliveTimeout: &tOut,
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
		Mode:               "tcp",
		Maxconn:            &mConn,
		HTTPConnectionMode: "http-tunnel",
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
