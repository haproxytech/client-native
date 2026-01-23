package main

import (
	_ "embed"
	"go/token"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/sirkon/dst"
	"github.com/sirkon/dst/decorator"
)

const (
	serverParamFileName      = "server_params.go"
	commentEnumFieldsToReset = `// Enum: ["enabled","disabled"]`
)

func main() {
	var allEnumFields []string
	var packageName string
	args := Args{}
	err := args.Parse()
	if err != nil {
		log.Panic(err)
	}
	for _, fileName := range args.Files {
		packageName, allEnumFields = enumFieldsEnabledDisabled(fileName)
		if err != nil {
			log.Panic(err)
		}

		// check that all enum
		// Enum: [enabled disabled]"
		// fields have an entry in the ServerParamsPrepareForRuntimeMap
		missingFields, errMissingField := checkMissingEnumFields(allEnumFields)
		if errMissingField != nil {
			// Exit with error if any new enum field is found and its behavior is not defined in ServerParamsPrepareForRuntimeMap
			log.Printf("There are some fields in models/server_params.go that are Enum[enabled disabled] but missing in `ServerParamsPrepareForRuntimeMap`")
			log.Printf("  File location `cmd/server_params_runtime/server_parans_runtime_fields_behavior.go`")
			log.Printf("ACTION: Please add them to `ServerParamsPrepareForRuntimeMap`")
			log.Printf("Missing fields %v", missingFields)
			log.Printf("\t For doc, read cmd/server_params_runtime/README.md")

			os.Exit(1)
		}

		// Generate the reset function using the template
		tmpl, errTmpl := template.New("generate.tmpl").Parse(tmplResetEnumDisabledFields)
		// ParseFiles(path.Join(templatePath))
		if errTmpl != nil {
			log.Panic(errTmpl)
		}

		generatedFileName := strings.TrimSuffix(fileName, ".go") + "_prepare_for_runtime.go"
		_ = os.Truncate(generatedFileName, 0)
		file, errFile := os.OpenFile(generatedFileName, os.O_CREATE|os.O_WRONLY, 0o600)
		if errFile != nil {
			log.Panic(errFile)
		}
		defer file.Close()

		doNotSendDisabledFields := listEmptyDisabledFields(allEnumFields)
		doNotSendEnabledFields := listEmtpyEnabledFields(allEnumFields)

		errTmpl = tmpl.Execute(file, map[string]any{
			"PrepareFieldsForRuntimeAddServer": FuncPrepareFieldsForRuntimeAddServer,
			"DoNotSendDisabledFields":          doNotSendDisabledFields,
			"DoNotSendDisabledFieldsFunc":      FuncDoNotSendDisabledFields,
			"DoNotSendEnabledFields":           doNotSendEnabledFields,
			"DoNotSendEnabledFieldsFunc":       FuncDoNotSendEnabledFields,
			"Package":                          packageName,
			"Licence":                          args.Licence,
		})
		if errTmpl != nil {
			log.Panic(errTmpl)
		}

		errFmt := fmtFile(generatedFileName)
		if errFmt != nil {
			log.Panic(errFmt)
		}
	}
}

func enumFieldsEnabledDisabled(filename string) (string, []string) { //nolint:gocognit
	var fieldsWithComment []string
	f, err := decorator.ParseFile(token.NewFileSet(), filename, nil, 0)
	if err != nil {
		log.Fatal(err)
	}

	for _, decl := range f.Decls {
		if genDecl, ok := decl.(*dst.GenDecl); ok && genDecl.Tok == token.TYPE {
			for _, spec := range genDecl.Specs {
				if typeSpec, ok := spec.(*dst.TypeSpec); ok {
					if structType, ok := typeSpec.Type.(*dst.StructType); ok {
						for _, field := range structType.Fields.List {
							comments := field.Decorations().Start.All()

							for _, comment := range comments {
								if comment == commentEnumFieldsToReset {
									// Add field name to the list
									if len(field.Names) > 0 {
										fieldsWithComment = append(fieldsWithComment, field.Names[0].Name)
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return f.Name.Name, fieldsWithComment
}
