package configuration

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/haproxytech/models"
)

func TestGetListeners(t *testing.T) {
	listeners, err := client.GetListeners("test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if len(listeners.Data) != 2 {
		t.Errorf("%v listeners returned, expected 2", len(listeners.Data))
	}

	if listeners.Version != version {
		t.Errorf("Version %v returned, expected %v", listeners.Version, version)
	}

	for _, l := range listeners.Data {
		if l.Name != "webserv" && l.Name != "webserv2" {
			t.Errorf("Expected only webserv or webserv2 listeners, %v found", l.Name)
		}
		if l.Address != "192.168.1.1" {
			t.Errorf("%v: Address not 192.168.1.1: %v", l.Name, l.Address)
		}
		if *l.Port != 80 && *l.Port != 8080 {
			t.Errorf("%v: Port not 80 or 8080: %v", l.Name, *l.Port)
		}
	}

	lJSON, err := listeners.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	listeners, err = client.GetListeners("test2", "")
	if err != nil {
		t.Error(err.Error())
	}
	if len(listeners.Data) > 0 {
		t.Errorf("%v listeners returned, expected 0", len(listeners.Data))
	}

	if !t.Failed() {
		fmt.Println("GetListeners succesful\nResponse: \n" + string(lJSON) + "\n")
	}
}

func TestGetListener(t *testing.T) {
	listener, err := client.GetListener("webserv", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	l := listener.Data

	if listener.Version != version {
		t.Errorf("Version %v returned, expected %v", listener.Version, version)
	}

	if l.Name != "webserv" {
		t.Errorf("Expected only webserv or webserv2 listeners, %v found", l.Name)
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

	_, err = client.GetListener("webserv", "test2", "")
	if err == nil {
		t.Error("Should throw error, non existant listener")
	}

	if !t.Failed() {
		fmt.Println("GetListener succesful\nResponse: \n" + string(lJSON) + "\n")
	}
}

func TestCreateEditDeleteListener(t *testing.T) {
	// TestCreateListener
	port := int64(4300)
	l := &models.Listener{
		Name:           "created",
		Address:        "192.168.2.1",
		Port:           &port,
		Ssl:            "enabled",
		SslCertificate: "dummy.crt",
	}

	err := client.CreateListener("test", l, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	listener, err := client.GetListener("created", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(listener.Data, l) {
		fmt.Printf("Created listener: %v\n", listener.Data)
		fmt.Printf("Given listener: %v\n", l)
		t.Error("Created listener not equal to given listener")
	}

	if listener.Version != version {
		t.Errorf("Version %v returned, expected %v", listener.Version, version)
	}

	err = client.CreateListener("test", l, "", version)
	if err == nil {
		t.Error("Should throw error listener already exists")
		version++
	}

	if !t.Failed() {
		fmt.Println("CreateListener successful")
	}

	// TestEditListener
	port = int64(5300)
	tOut := int64(5)
	l = &models.Listener{
		Name:           "created",
		Address:        "192.168.3.1",
		Port:           &port,
		Transparent:    "enabled",
		TCPUserTimeout: &tOut,
	}

	err = client.EditListener("created", "test", l, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	listener, err = client.GetListener("created", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(listener.Data, l) {
		fmt.Printf("Edited listener: %v\n", listener.Data)
		fmt.Printf("Given lsitener: %v\n", l)
		t.Error("Edited listener not equal to given listener")
	}

	if listener.Version != version {
		t.Errorf("Version %v returned, expected %v", listener.Version, version)
	}

	if !t.Failed() {
		fmt.Println("EditListener successful")
	}

	// TestDeleteListener
	err = client.DeleteListener("created", "test", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ := client.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	_, err = client.GetListener("created", "test", "")
	if err == nil {
		t.Error("DeleteListener failed, listener test still exists")
	}

	err = client.DeleteListener("created", "test2", "", version)
	if err == nil {
		t.Error("Should throw error, non existant listener")
		version++
	}

	if !t.Failed() {
		fmt.Println("DeleteListener successful")
	}
}
