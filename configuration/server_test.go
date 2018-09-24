package configuration

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/haproxytech/models"
)

func TestGetServers(t *testing.T) {
	servers, err := client.GetServers("test")
	if err != nil {
		t.Error(err.Error())
	}

	if len(servers.Data) != 2 {
		t.Errorf("%v servers returned, expected 2", len(servers.Data))
	}

	if servers.Version != version {
		t.Errorf("Version %v returned, expected %v", servers.Version, version)
	}

	for _, s := range servers.Data {
		if s.Name != "webserv" && s.Name != "webserv2" {
			t.Errorf("Expected only webserv or webserv2 servers, %v found", s.Name)
		}
		if s.Address != "192.168.1.1" {
			t.Errorf("%v: Address not 192.168.1.1: %v", s.Name, s.Address)
		}
		if *s.Port != 9300 && *s.Port != 9200 {
			t.Errorf("%v: Port not 9300 or 9200: %v", s.Name, *s.Port)
		}
		if s.Ssl != "enabled" {
			t.Errorf("%v: Ssl not enabled: %v", s.Name, s.Ssl)
		}
		if s.HTTPCookieID != "BLAH" {
			t.Errorf("%v: HTTPCookieID not BLAH: %v", s.Name, s.HTTPCookieID)
		}
		if *s.MaxConnections != 1000 {
			t.Errorf("%v: MaxConnections not 1000: %v", s.Name, *s.MaxConnections)
		}
		if *s.Weight != 10 {
			t.Errorf("%v: Weight not 10: %v", s.Name, *s.Weight)
		}
	}

	sJSON, err := servers.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	servers, err = client.GetServers("test2")
	if err != nil {
		t.Error(err.Error())
	}
	if len(servers.Data) > 0 {
		t.Errorf("%v servers returned, expected 0", len(servers.Data))
	}

	if !t.Failed() {
		fmt.Println("GetServers succesful\nResponse: \n" + string(sJSON) + "\n")
	}
}

func TestGetServer(t *testing.T) {
	server, err := client.GetServer("webserv", "test")
	if err != nil {
		t.Error(err.Error())
	}

	s := server.Data

	if server.Version != version {
		t.Errorf("Version %v returned, expected %v", server.Version, version)
	}

	if s.Name != "webserv" {
		t.Errorf("Expected only webserv, %v found", s.Name)
	}
	if s.Address != "192.168.1.1" {
		t.Errorf("%v: Address not 192.168.1.1: %v", s.Name, s.Address)
	}
	if *s.Port != 9200 {
		t.Errorf("%v: Port not 9200: %v", s.Name, *s.Port)
	}
	if s.Ssl != "enabled" {
		t.Errorf("%v: Ssl not enabled: %v", s.Name, s.Ssl)
	}
	if s.HTTPCookieID != "BLAH" {
		t.Errorf("%v: HTTPCookieID not BLAH: %v", s.Name, s.HTTPCookieID)
	}
	if *s.MaxConnections != 1000 {
		t.Errorf("%v: MaxConnections not 1000: %v", s.Name, *s.MaxConnections)
	}
	if *s.Weight != 10 {
		t.Errorf("%v: Weight not 10: %v", s.Name, *s.Weight)
	}

	sJSON, err := s.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	_, err = client.GetServer("webserv", "test2")
	if err == nil {
		t.Error("Should throw error, non existant server")
	}

	if !t.Failed() {
		fmt.Println("GetServer succesful\nResponse: \n" + string(sJSON) + "\n")
	}
}

func TestCreateEditDeleteServer(t *testing.T) {
	// TestCreateServer
	port := int64(4300)
	s := &models.Server{
		Name:        "created",
		Address:     "192.168.2.1",
		Port:        &port,
		Sorry:       "enabled",
		Check:       "enabled",
		Maintenance: "enabled",
	}

	err := client.CreateServer("test", s, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version = version + 1
	}

	server, err := client.GetServer("created", "test")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(server.Data, s) {
		fmt.Printf("Created server: %v\n", server.Data)
		fmt.Printf("Given server: %v\n", s)
		t.Error("Created server not equal to given server")
	}

	if server.Version != version {
		t.Errorf("Version %v returned, expected %v", server.Version, version)
	}

	err = client.CreateServer("test", s, "", version)
	if err == nil {
		t.Error("Should throw error server already exists")
		version = version + 1
	}

	if !t.Failed() {
		fmt.Println("CreateServer successful")
	}

	// TestEditServer
	port = int64(5300)
	s = &models.Server{
		Name:    "created",
		Address: "192.168.3.1",
		Port:    &port,
	}

	err = client.EditServer("created", "test", s, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version = version + 1
	}

	server, err = client.GetServer("created", "test")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(server.Data, s) {
		fmt.Printf("Edited server: %v\n", server.Data)
		fmt.Printf("Given server: %v\n", s)
		t.Error("Edited server not equal to given server")
	}

	if server.Version != version {
		t.Errorf("Version %v returned, expected %v", server.Version, version)
	}

	if !t.Failed() {
		fmt.Println("EditServer successful")
	}

	// TestDeleteServer
	err = client.DeleteServer("created", "test", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version = version + 1
	}

	if v, _ := client.GetVersion(); v != version {
		t.Error("Version not incremented")
	}

	_, err = client.GetServer("created", "test")
	if err == nil {
		t.Error("DeleteServer failed, server test still exists")
	}

	err = client.DeleteServer("created", "test2", "", version)
	if err == nil {
		t.Error("Should throw error, non existant server")
		version = version + 1
	}

	if !t.Failed() {
		fmt.Println("DeleteServer successful")
	}
}
