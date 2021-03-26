package configuration

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/haproxytech/client-native/v2/models"
)

func TestGetServerTemplates(t *testing.T) { //nolint:gocognit,gocyclo

	v, templates, err := client.GetServerTemplates("test", "")

	if err != nil {
		t.Error(err.Error())
	}

	if len(templates) != 3 {
		t.Errorf("%v server templates returned, expected 3", len(templates))
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	for _, template := range templates {
		if template.Check != "enabled" {
			t.Errorf("%v: Check not enabled: %v", template.Prefix, template.Check)
		}
	}

}

func TestGetServerTemplate(t *testing.T) {

	v, template, err := client.GetServerTemplate("srv", "test", "")

	if err != nil {
		t.Error(err.Error())
	}
	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}
	if template.Prefix != "srv" {
		t.Errorf("Expected only srv, %v found", template.Prefix)
	}
	if template.NumOrRange != "1-3" {
		t.Errorf("Expected 1-3, %v found", template.NumOrRange)
	}
	if template.Fqdn != "google.com:80" {
		t.Errorf("Expected google.com:80, %v found", template.Fqdn)
	}
	if template.Check != "enabled" {
		t.Errorf("Expected check to be enabled")
	}
	_, err = template.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}
	_, _, err = client.GetServerTemplate("test2", "example", "")
	if err == nil {
		t.Error("Should throw error, non existant server template")
	}
}

func TestGetServerTemplateSecond(t *testing.T) {

	v, template, err := client.GetServerTemplate("site", "test", "")

	if err != nil {
		t.Error(err.Error())
	}
	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}
	if template.Prefix != "site" {
		t.Errorf("Expected only site, %v found", template.Prefix)
	}
	if template.NumOrRange != "1-10" {
		t.Errorf("Expected 1-10, %v found", template.NumOrRange)
	}
	if template.Fqdn != "google.com:8080" {
		t.Errorf("Expected google.com:8080, %v found", template.Fqdn)
	}
	if template.Check != "enabled" {
		t.Errorf("Expected check to be enabled")
	}
	if template.Backup != "enabled" {
		t.Errorf("Expected backup to be enabled")
	}
	_, err = template.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}
}

func TestGetServerTemplateThird(t *testing.T) {

	v, template, err := client.GetServerTemplate("website", "test", "")

	if err != nil {
		t.Error(err.Error())
	}
	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}
	if template.Prefix != "website" {
		t.Errorf("Expected only website, %v found", template.Prefix)
	}
	if template.NumOrRange != "10-100" {
		t.Errorf("Expected 10-100, %v found", template.NumOrRange)
	}
	if template.Fqdn != "google.com:443" {
		t.Errorf("Expected google.com:443, %v found", template.Fqdn)
	}
	if template.Check != "enabled" {
		t.Errorf("Expected check to be enabled")
	}
	if template.Backup != "disabled" {
		t.Errorf("Expected backup to be disabled")
	}
	_, err = template.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}
}

func TestCreateEditDeleteServerTemplate(t *testing.T) {

	// TestCreateServerTemplate
	template := &models.ServerTemplate{
		Prefix:     "dev",
		NumOrRange: "1-10",
		Fqdn:       "site.com:80",
		Check:      "enabled",
	}

	err := client.CreateServerTemplate("test", template, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, serverTemplate, err := client.GetServerTemplate("dev", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(serverTemplate, template) {
		fmt.Printf("Created server template: %v\n", serverTemplate)
		fmt.Printf("Given server template: %v\n", template)
		t.Error("Created server template is not equal to the given server")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	err = client.CreateServerTemplate("test", template, "", version)
	if err == nil {
		t.Error("Should throw error server already exists")
		version++
	}

	// EditServerTemplate
	template = &models.ServerTemplate{
		Prefix:     "dev",
		NumOrRange: "11-20",
		Fqdn:       "site.com:8080",
		Check:      "disabled",
	}

	err = client.EditServerTemplate("dev", "test", template, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, serverTemplate, err = client.GetServerTemplate("dev", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(serverTemplate, template) {
		fmt.Printf("Created server template: %v\n", serverTemplate)
		fmt.Printf("Given server template: %v\n", template)
		t.Error("Created server template is not equal to the given server")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	// TestDeleteServerTemplate
	err = client.DeleteServerTemplate("dev", "test", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ := client.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	_, _, err = client.GetServerTemplate("dev", "test", "")
	if err == nil {
		t.Error("DeleteServerTemplate failed, server test still exists")
	}

	err = client.DeleteServerTemplate("dev", "test2", "", version)
	if err == nil {
		t.Error("Should throw error, non existant server")
		version++
	}
}
