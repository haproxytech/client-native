package main

import (
	"strings"
	"text/template"
)

func generateEqualAndDiff(opt generateEqualAndDiffOptions) error {
	funcMaps := template.FuncMap{
		"HasPrefix": strings.HasPrefix,
		"Title":     toTitle,
		"CamelCase": toCamelCase,
		"LowerCase": toLowerCase,
		"JSON":      toJSON,
	}
	tmpl, err := template.New("generate.tmpl").Funcs(funcMaps).Parse(tmplEqualAndDiff)
	if err != nil {
		return err
	}

	data := map[string]interface{}{
		"Mode":              opt.Mode,
		"Name":              opt.Name,
		"Type":              opt.Type,
		"Fields":            opt.Fields,
		"NeedsOptions":      opt.NeedsOptions,
		"NeedsOptionsIndex": opt.NeedsOptionsIndex,
		"IsBasicType":       opt.IsBasicType,
		"IsComplex":         opt.IsComplex,
		"IsComparable":      opt.IsComparable,
		"IsPointer":         opt.IsPointer,
	}

	functions := map[string]interface{}{
		"Functions": []interface{}{
			map[string]interface{}{
				"Name": "Equal",
				"Data": data,
			},
			map[string]interface{}{
				"Name": "Diff",
				"Data": data,
			},
		},
	}

	err = tmpl.Execute(opt.File, functions)
	if err != nil {
		return err
	}
	return generateCompareTests(opt)
}

func generateCompareTests(opt generateEqualAndDiffOptions) error {
	if opt.Mode == "array" {
		return nil
	}
	funcMaps := template.FuncMap{
		"HasPrefix": strings.HasPrefix,
		"Title":     toTitle,
		"CamelCase": toCamelCase,
	}
	tmpl, err := template.New("generate.tmpl").Funcs(funcMaps).Parse(tmplCompareTest)
	if err != nil {
		return err
	}
	hasIndex := false
	for _, file := range opt.Fields {
		if file.Name == "Index" && file.Type == "*int64" {
			hasIndex = true
			break
		}
	}
	err = tmpl.Execute(opt.FileTest, map[string]interface{}{
		"TestType":    []string{"Equal", "Diff"},
		"Name":        opt.Name,
		"HasIndex":    hasIndex,
		"Fields":      opt.Fields,
		"FieldCount":  len(opt.Fields),
		"PackageName": opt.PackageName,
	})
	if err != nil {
		return err
	}
	return nil
}
