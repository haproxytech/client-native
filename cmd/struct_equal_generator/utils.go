package main

import (
	_ "embed"
	"os"
	"path"
	"text/template"
)

//go:embed utils.tmpl
var tmplUtils string

func createUtilsFile(packageName string, args Args) error {
	fileName := path.Join(args.Directory, "utils_compare.go")
	_ = os.Truncate(fileName, 0)
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	defer file.Close()

	// Adding the header to the generated file
	tmpl, err := template.New("generate.tmpl").Parse(tmplHeader)
	// ParseFiles(path.Join(templatePath))
	if err != nil {
		return err
	}

	err = tmpl.Execute(file, map[string]interface{}{
		"Package": packageName,
		"License": args.License,
	})
	if err != nil {
		return err
	}

	funcMaps := template.FuncMap{}
	tmpl, err = template.New("generate.tmpl").Funcs(funcMaps).Parse(tmplUtils)
	if err != nil {
		return err
	}
	err = tmpl.Execute(file, map[string]interface{}{})
	if err != nil {
		return err
	}
	err = file.Close()
	if err != nil {
		return err
	}
	return fmtFile(fileName)
}
