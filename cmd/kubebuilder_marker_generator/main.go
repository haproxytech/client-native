package main

import (
	_ "embed"
	"fmt"
	"go/token"
	"log"
	"os"
	"strings"

	"github.com/sirkon/dst"
	"github.com/sirkon/dst/decorator"
)

const (
	kubebuilderValidationMarker = "// +kubebuilder:validation:"
)

func main() {
	// tool to add `json:",inline"` for embedded structs

	args := Args{}
	err := args.Parse()
	if err != nil {
		log.Panic(err)
	}
	for _, fileName := range args.Files {
		err = generate(fileName)
		if err != nil {
			log.Panic(err)
		}
	}
}

func generate(fileName string) error { //nolint:gocognit
	f, err := decorator.ParseFile(token.NewFileSet(), fileName, nil, 0)
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
							// Remove // +kubebuilder:validation
							comments = cleanup(comments, kubebuilderValidationMarker)
							field.Decorations().Start.Replace(comments...)
							// Then do the job
							for _, comment := range comments {
								if strings.HasPrefix(comment, "// Enum: [") {
									field.Decorations().Before = dst.NewLine
									newComment := kubebuilderValidationMarker + `Enum=`
									comment = strings.TrimPrefix(comment, "// Enum: [")
									comment = strings.TrimSuffix(comment, "]")
									// We must keep empty strings:
									// For example in Globals HttpclientSslVerify: // Enum: [ none required]
									// from swagger: enum: ["", "none", "required"]
									for _, enum := range strings.Split(comment, " ") {
										if enum == "" {
											newComment += `""`
										}
										newComment += enum
										newComment += ";"
									}
									field.Decorations().Start.Append(newComment)
									log.Printf("Adding comment for: %s: %s %s\n", fileName, field.Names[0].Name, newComment)
								}
								if strings.HasPrefix(comment, "// Pattern: ") {
									addSimpleMarker(field, fileName, comment, "Pattern", "string")
								}
								if strings.HasPrefix(comment, "// Maximum: ") {
									addSimpleMarker(field, fileName, comment, "Maximum", "raw")
								}
								if strings.HasPrefix(comment, "// Minimum: ") {
									addSimpleMarker(field, fileName, comment, "Minimum", "raw")
								}
								if strings.HasPrefix(comment, "// Format: ") {
									addSimpleMarker(field, fileName, comment, "Format", "raw")
								}
								if strings.HasPrefix(comment, "// Required: true") {
									if len(comments) == 2 && comments[0] == "// index" {
										field.Decorations().Before = dst.NewLine
										field.Decorations().Start.Append("// +kubebuilder:validation:Optional")
									}
								}
							}
							// if len(field.Names) > 0 {
							// log.Printf("Comments before the field %s: %v\n", field.Names[0].Name, comments)
							// }
						}
					}
				}
			}
		}
	}

	outputFile, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	// decorator.Print(f)

	err = decorator.Fprint(outputFile, f)
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	return nil
}

// addSimpleMarker adds a simple kubebuilder marker to the given field.
//
// Parameters:
//   - field: the field to add the marker to.
//   - fileName: the name of the file.
//   - markerValue: the marker value.
//   - validationType: the type of validation (for example, "Pattern", "Maximum", "Minimum").
//   - validationContentType: the content type of the validation:
//     "string" to add “ around the comment
//     "raw" to add the comment as is
func addSimpleMarker(field *dst.Field, fileName, markerValue, validationType, validationContentType string) {
	field.Decorations().Before = dst.NewLine
	markerValue = strings.TrimPrefix(markerValue, "// "+validationType+": ")
	var marker string
	switch validationContentType {
	case "string":
		marker = fmt.Sprintf("%s%s=`%v`", kubebuilderValidationMarker, validationType, markerValue)
	case "raw":
		marker = fmt.Sprintf("%s%s=%v", kubebuilderValidationMarker, validationType, markerValue)
	default:
		log.Printf("Unknown validation content type: %s", validationContentType)
	}
	field.Decorations().Start.Append(marker)
	log.Printf("Adding comment for: %s: %+v %s\n", fileName, field.Names[0].Name, marker)
}

// cleanup removes all strings from a slice of strings that have a specific prefix.
//
// comments is the slice of strings containing the comments to clean up.
// prefixToRemove is the prefix of the comments to remove.
// It returns a new slice of strings without the comments that have the specified prefix.
func cleanup(comments []string, prefixToRemove string) []string {
	var res []string
	for _, comment := range comments {
		if !strings.HasPrefix(comment, prefixToRemove) {
			res = append(res, comment)
		}
	}
	return res
}
