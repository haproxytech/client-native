package test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/haproxytech/client-native/v5/models"
	"github.com/stretchr/testify/require"
)

func serverTemplateExpectation() map[string]models.ServerTemplates {
	initStructuredExpected()
	res := StructuredToServerTemplateMap()
	// Add individual entries
	for k, vs := range res {
		for _, v := range vs {
			key := fmt.Sprintf("%s/%s", k, v.Prefix)
			res[key] = models.ServerTemplates{v}
		}
	}
	return res
}

func TestGetServerTemplates(t *testing.T) { //nolint:gocognit,gocyclo
	mt := make(map[string]models.ServerTemplates)
	v, templates, err := clientTest.GetServerTemplates("test", "")
	if err != nil {
		t.Error(err.Error())
	}
	mt["backend/test"] = templates

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	checkServerTemplates(t, mt)
}

func checkServerTemplates(t *testing.T, got map[string]models.ServerTemplates) {
	exp := serverTemplateExpectation()
	for k, v := range got {
		want, ok := exp[k]
		require.True(t, ok, "k=%s", k)
		require.Equal(t, len(want), len(v), "k=%s", k)
		for _, g := range v {
			for _, w := range want {
				if g.Prefix == w.Prefix {
					require.True(t, g.Equal(*w), "k=%s - diff %v", k, cmp.Diff(*g, *w))
					break
				}
			}
		}
	}
}

func TestGetServerTemplate(t *testing.T) {
	m := make(map[string]models.ServerTemplates)
	v, template, err := clientTest.GetServerTemplate("srv", "test", "")
	if err != nil {
		t.Error(err.Error())
	}
	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}
	m["backend/test/srv"] = models.ServerTemplates{template}
	_, _, err = clientTest.GetServerTemplate("test2", "example", "")
	if err == nil {
		t.Error("Should throw error, non existent server template")
	}
	checkServerTemplates(t, m)
}

func TestGetServerTemplateSecond(t *testing.T) {
	m := make(map[string]models.ServerTemplates)

	v, template, err := clientTest.GetServerTemplate("site", "test", "")
	if err != nil {
		t.Error(err.Error())
	}
	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}
	m["backend/test/site"] = models.ServerTemplates{template}

	_, err = template.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}
}

func TestGetServerTemplateThird(t *testing.T) {
	m := make(map[string]models.ServerTemplates)

	v, template, err := clientTest.GetServerTemplate("website", "test", "")
	if err != nil {
		t.Error(err.Error())
	}
	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}
	m["backend/test/website"] = models.ServerTemplates{template}

	_, err = template.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	checkServerTemplates(t, m)
}

func TestGetServerTemplateFourth(t *testing.T) {
	m := make(map[string]models.ServerTemplates)

	v, template, err := clientTest.GetServerTemplate("test", "test", "")
	if err != nil {
		t.Error(err.Error())
	}
	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}
	m["backend/test/test"] = models.ServerTemplates{template}

	_, err = template.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}
	checkServerTemplates(t, m)
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
