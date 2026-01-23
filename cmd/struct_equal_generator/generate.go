package main

import (
	"strings"
	"text/template"
)

func generateEqualAndDiff(opt generateEqualAndDiffOptions) error {
	funcMaps := template.FuncMap{
		"HasPrefix":  strings.HasPrefix,
		"TrimPrefix": strings.TrimPrefix,
		"Title":      toTitle,
		"CamelCase":  toCamelCase,
		"LowerCase":  toLowerCase,
		"JSON":       toJSON,
	}
	tmpl, err := template.New("generate.tmpl").Funcs(funcMaps).Parse(tmplEqualAndDiff)
	if err != nil {
		return err
	}

	data := map[string]any{
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

	functions := map[string]any{
		"Functions": []any{
			map[string]any{
				"Name": "Equal",
				"Data": data,
			},
			map[string]any{
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
	metadataIndex := -1
	for i, f := range opt.Fields {
		if f.Name == "Metadata" {
			metadataIndex = i
			break
		}
	}
	if metadataIndex > -1 {
		opt.Fields = append(opt.Fields[:metadataIndex], opt.Fields[metadataIndex+1:]...)
	}
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
	err = tmpl.Execute(opt.FileTest, map[string]any{
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
