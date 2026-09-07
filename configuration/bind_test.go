package configuration

import (
	"testing"

	"github.com/haproxytech/client-native/v6/config-parser/params"
	"github.com/haproxytech/client-native/v6/configuration/options"
	"github.com/haproxytech/client-native/v6/models"

	"github.com/stretchr/testify/require"
)

func TestParseBindParamsMultipleCrtList(t *testing.T) {
	bindOptions := []params.BindOption{
		&params.BindOptionValue{Name: "crt-list", Value: "list1.txt"},
		&params.BindOptionValue{Name: "crt-list", Value: "list2.txt"},
	}

	b, _ := parseBindParams(bindOptions)

	require.Equal(t, "list1.txt:list2.txt", b.CrtList)
}

func TestSerializeBindParamsMultipleCrtList(t *testing.T) {
	b := models.BindParams{CrtList: "list1.txt:list2.txt"}

	bindOptions := serializeBindParams(b, "", "", &options.ConfigurationOptions{})

	var crtLists []string
	for _, o := range bindOptions {
		if v, ok := o.(*params.BindOptionValue); ok && v.Name == "crt-list" {
			crtLists = append(crtLists, v.Value)
		}
	}
	require.Equal(t, []string{"list1.txt", "list2.txt"}, crtLists)
}
