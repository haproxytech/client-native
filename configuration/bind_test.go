package configuration

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/haproxytech/models"
)

func TestGetBinds(t *testing.T) {
	binds, err := client.GetBinds("test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if len(binds.Data) != 2 {
		t.Errorf("%v binds returned, expected 2", len(binds.Data))
	}

	if binds.Version != version {
		t.Errorf("Version %v returned, expected %v", binds.Version, version)
	}

	for _, l := range binds.Data {
		if l.Name != "webserv" && l.Name != "webserv2" {
			t.Errorf("Expected only webserv or webserv2 binds, %v found", l.Name)
		}
		if l.Address != "192.168.1.1" {
			t.Errorf("%v: Address not 192.168.1.1: %v", l.Name, l.Address)
		}
		if *l.Port != 80 && *l.Port != 8080 {
			t.Errorf("%v: Port not 80 or 8080: %v", l.Name, *l.Port)
		}
	}

	lJSON, err := binds.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	binds, err = client.GetBinds("test_2", "")
	if err != nil {
		t.Error(err.Error())
	}
	if len(binds.Data) > 0 {
		t.Errorf("%v binds returned, expected 0", len(binds.Data))
	}

	if !t.Failed() {
		fmt.Println("GetBinds succesful\nResponse: \n" + string(lJSON) + "\n")
	}
}

func TestGetBind(t *testing.T) {
	bind, err := client.GetBind("webserv", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	l := bind.Data

	if bind.Version != version {
		t.Errorf("Version %v returned, expected %v", bind.Version, version)
	}

	if l.Name != "webserv" {
		t.Errorf("Expected only webserv or webserv2 binds, %v found", l.Name)
	}
	if l.Address != "192.168.1.1" {
		t.Errorf("%v: Address not 192.168.1.1: %v", l.Name, l.Address)
	}
	if *l.Port != 80 {
		t.Errorf("%v: Port not 80 or 8080: %v", l.Name, *l.Port)
	}

	lJSON, err := l.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	_, err = client.GetBind("webserv", "test_2", "")
	if err == nil {
		t.Error("Should throw error, non existant bind")
	}

	if !t.Failed() {
		fmt.Println("GetBind succesful\nResponse: \n" + string(lJSON) + "\n")
	}
}

func TestCreateEditDeleteBind(t *testing.T) {
	// TestCreateBind
	port := int64(4300)
	l := &models.Bind{
		Name:           "created",
		Address:        "192.168.2.1",
		Port:           &port,
		Ssl:            true,
		SslCertificate: "dummy.crt",
	}

	err := client.CreateBind("test", l, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	bind, err := client.GetBind("created", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(bind.Data, l) {
		fmt.Printf("Created bind: %v\n", bind.Data)
		fmt.Printf("Given bind: %v\n", l)
		t.Error("Created bind not equal to given bind")
	}

	if bind.Version != version {
		t.Errorf("Version %v returned, expected %v", bind.Version, version)
	}

	err = client.CreateBind("test", l, "", version)
	if err == nil {
		t.Error("Should throw error bind already exists")
		version++
	}

	if !t.Failed() {
		fmt.Println("CreateBind successful")
	}

	// TestEditBind
	port = int64(5300)
	tOut := int64(5)
	l = &models.Bind{
		Name:           "created",
		Address:        "192.168.3.1",
		Port:           &port,
		Transparent:    true,
		TCPUserTimeout: &tOut,
	}

	err = client.EditBind("created", "test", l, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	bind, err = client.GetBind("created", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(bind.Data, l) {
		fmt.Printf("Edited bind: %v\n", bind.Data)
		fmt.Printf("Given lsitener: %v\n", l)
		t.Error("Edited bind not equal to given bind")
	}

	if bind.Version != version {
		t.Errorf("Version %v returned, expected %v", bind.Version, version)
	}

	if !t.Failed() {
		fmt.Println("EditBind successful")
	}

	// TestDeleteBind
	err = client.DeleteBind("created", "test", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ := client.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	_, err = client.GetBind("created", "test", "")
	if err == nil {
		t.Error("DeleteBind failed, bind test still exists")
	}

	err = client.DeleteBind("created", "test2", "", version)
	if err == nil {
		t.Error("Should throw error, non existant bind")
		version++
	}

	if !t.Failed() {
		fmt.Println("DeleteBind successful")
	}
}
