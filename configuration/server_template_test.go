package configuration

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/haproxytech/client-native/v5/models"
)

func TestGetServerTemplates(t *testing.T) { //nolint:gocognit,gocyclo

	v, templates, err := clientTest.GetServerTemplates("test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if len(templates) != 4 {
		t.Errorf("%v server templates returned, expected 4", len(templates))
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
	v, template, err := clientTest.GetServerTemplate("srv", "test", "")
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
	if template.Fqdn != "google.com" {
		t.Errorf("Expected google.com, %v found", template.Fqdn)
	}
	if *template.Port != 80 {
		t.Errorf(" Expected 80, %v found", template.Port)
	}
	if template.Check != "enabled" {
		t.Errorf("Expected check to be enabled")
	}
	_, err = template.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}
	_, _, err = clientTest.GetServerTemplate("test2", "example", "")
	if err == nil {
		t.Error("Should throw error, non existent server template")
	}
}

func TestGetServerTemplateSecond(t *testing.T) {
	v, template, err := clientTest.GetServerTemplate("site", "test", "")
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
	if template.Fqdn != "google.com" {
		t.Errorf("Expected google.com, %v found", template.Fqdn)
	}
	if *template.Port != 8080 {
		t.Errorf("Expected 8080, %v found", *template.Port)
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
	v, template, err := clientTest.GetServerTemplate("website", "test", "")
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
	if template.Fqdn != "google.com" {
		t.Errorf("Expected google.com, %v found", template.Fqdn)
	}
	if *template.Port != 443 {
		t.Errorf("Expected 443, %v found", *template.Port)
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

func TestGetServerTemplateFourth(t *testing.T) {
	v, template, err := clientTest.GetServerTemplate("test", "test", "")
	if err != nil {
		t.Error(err.Error())
	}
	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}
	if template.Prefix != "test" {
		t.Errorf("Expected only website, %v found", template.Prefix)
	}
	if template.NumOrRange != "5" {
		t.Errorf("Expected 5, %v found", template.NumOrRange)
	}
	if template.Fqdn != "test.com" {
		t.Errorf("Expected test.com, %v found", template.Fqdn)
	}
	if *template.Port != 0 {
		t.Errorf("Expected 0, %v found", *template.Port)
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

func TestCreateEditDeleteServerTemplate(t *testing.T) {
	port := int64(80)

	// TestCreateServerTemplate
	template := &models.ServerTemplate{
		Prefix:     "dev",
		NumOrRange: "1-10",
		Fqdn:       "site.com",
		Port:       &port,
		ServerParams: models.ServerParams{
			Check: "enabled",
		},
	}

	err := clientTest.CreateServerTemplate("test", template, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, serverTemplate, err := clientTest.GetServerTemplate("dev", "test", "")
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

	err = clientTest.CreateServerTemplate("test", template, "", version)
	if err == nil {
		t.Error("Should throw error server already exists")
		version++
	}

	port = 8080

	// EditServerTemplate
	template = &models.ServerTemplate{
		Prefix:     "dev",
		NumOrRange: "11-20",
		Fqdn:       "site.com",
		Port:       &port,
		ServerParams: models.ServerParams{
			Check: "disabled",
		},
	}

	err = clientTest.EditServerTemplate("dev", "test", template, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, serverTemplate, err = clientTest.GetServerTemplate("dev", "test", "")
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
	err = clientTest.DeleteServerTemplate("dev", "test", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ := clientTest.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	_, _, err = clientTest.GetServerTemplate("dev", "test", "")
	if err == nil {
		t.Error("DeleteServerTemplate failed, server test still exists")
	}

	err = clientTest.DeleteServerTemplate("dev", "test2", "", version)
	if err == nil {
		t.Error("Should throw error, non existent server")
		version++
	}
}
