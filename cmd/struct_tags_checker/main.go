package main

import (
	_ "embed"
	"go/token"
	"log"
	"os"

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
							if len(field.Names) == 0 {
								// This is an embedded struct.
								if field.Tag == nil {
									log.Printf("Embedded struct found with no tag: %s: %+v\n", fileName, field.Type)
									field.Tag = &dst.BasicLit{
										Kind:  token.STRING,
										Value: "`json:\",inline\"`",
									}
									continue
								}
							}
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

	err = decorator.Fprint(outputFile, f)
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	return nil
}
