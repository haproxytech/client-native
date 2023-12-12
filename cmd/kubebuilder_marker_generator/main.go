package main

import (
	_ "embed"
	"go/token"
	"log"
	"os"
	"strings"

	"github.com/sirkon/dst"
	"github.com/sirkon/dst/decorator"
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

func generate(fileName string) error { //nolint:gocognit,unparam
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
							for _, comment := range comments {
								if strings.HasPrefix(comment, "// Enum: [") {
									field.Decorations().Before = dst.NewLine
									newComment := `// +kubebuilder:validation:Enum=`
									comment = strings.TrimPrefix(comment, "// Enum: [")
									comment = strings.TrimSuffix(comment, "]")
									for _, enum := range strings.Split(comment, " ") {
										enum = strings.TrimSpace(enum)
										newComment += enum
										newComment += ";"
									}
									field.Decorations().Start.Append(newComment)
									log.Printf("Adding comment for: %s: %+v %s\n", fileName, field.Type, newComment)
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
