package main

import (
	_ "embed"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"strings"
	"text/template"
)

func main() {
	var packageName string
	args := Args{}
	err := args.Parse()
	if err != nil {
		log.Panic(err)
	}
	for _, fileName := range args.Files {
		packageName, err = generate(fileName, args)
		if err != nil {
			log.Panic(err)
		}
	}
	err = createUtilsFile(packageName, args)
	if err != nil {
		log.Panic(err)
	}
}

var sourceofFile string //nolint:gochecknoglobals

//nolint:gocognit
func generate(fileName string, args Args) (string, error) {
	fset := token.NewFileSet()
	var packageName string

	src, err := os.ReadFile(fileName)
	if err != nil {
		return packageName, err
	}
	sourceofFile = string(src)

	// node, err := parser.ParseFile(fset, fileName, nil, parser.ParseComments)
	node, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		return packageName, err
	}
	generatedFileName := strings.TrimSuffix(fileName, ".go") + "_compare.go"
	generatedFileNameTest := strings.TrimSuffix(fileName, ".go") + "_compare_test.go"

	_ = os.Truncate(generatedFileName, 0)
	file, err := os.OpenFile(generatedFileName, os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return packageName, err
	}
	defer file.Close()

	_ = os.Truncate(generatedFileNameTest, 0)
	fileTest, err := os.OpenFile(generatedFileNameTest, os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return packageName, err
	}
	defer fileTest.Close()

	// Adding the header to the generated file
	tmpl, err := template.New("generate.tmpl").Parse(tmplHeader)
	// ParseFiles(path.Join(templatePath))
	if err != nil {
		return packageName, err
	}

	err = tmpl.Execute(file, map[string]interface{}{
		"Package": node.Name.String(),
		"Licence": args.Licence,
	})
	if err != nil {
		return packageName, err
	}

	// Adding the header to the generated file
	tmpl2, err := template.New("generate.tmpl").Parse(tmplHeader)
	// ParseFiles(path.Join(templatePath))
	if err != nil {
		return packageName, err
	}

	err = tmpl2.Execute(fileTest, map[string]interface{}{
		"Package": node.Name.String(),
		"Licence": args.Licence,
	})
	if err != nil {
		return packageName, err
	}

	packageName = node.Name.String()

	hasTests := false
	// For each declaration in the node's declarations list, check if it is a generic declaration.
	// If it is, get the type specification of each generic type.
	for _, f := range node.Decls {
		g, ok := f.(*ast.GenDecl)
		if !ok {
			continue
		}

		for _, spec := range g.Specs {
			currSpecType, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			switch currType := currSpecType.Type.(type) {
			case *ast.StructType:
				var fields []Field
				needsOptions := false
				needsOptionsIndex := false
				for _, field := range currType.Fields.List {
					if len(field.Names) > 0 {
						res := getTypeString(field.Type)
						fields = append(fields, Field{
							Name:         field.Names[0].Name,
							Type:         res.Name,
							TypeInFile:   res.TypeName,
							IsBasicType:  res.IsBasicType,
							IsComparable: res.IsComparable,
						})
						if field.Names[0].Name == "Index" {
							needsOptionsIndex = true
						}
						if strings.HasPrefix(res.Name, "[]") {
							needsOptions = true
						}
						if strings.HasPrefix(res.Name, "map") {
							needsOptions = true
						}
						needsOptions = needsOptions || res.IsComplex
					}
					// For embedded struct
					if len(field.Names) == 0 && field.Type != nil {
						res := getTypeString(field.Type)
						fields = append(fields, Field{
							Name:         res.Name,
							IsEmbedded:   true,
							IsComparable: res.IsComparable,
						})
						if res.Name == "Index" {
							needsOptionsIndex = true
						}
						needsOptions = true
					}
				}
				hasTests = true
				err = generateEqualAndDiff(generateEqualAndDiffOptions{
					PackageName:       packageName,
					File:              file,
					FileTest:          fileTest,
					Name:              currSpecType.Name.Name,
					CurrType:          currSpecType,
					Fields:            fields,
					NeedsOptions:      needsOptions,
					NeedsOptionsIndex: needsOptionsIndex,
					Mode:              "struct",
				})
				if err != nil {
					log.Panic(err)
				}
			case *ast.Ident:
				hasTests = true
				err = generateEqualAndDiff(generateEqualAndDiffOptions{
					PackageName:  packageName,
					File:         file,
					FileTest:     fileTest,
					Name:         currSpecType.Name.Name,
					NeedsOptions: false,
					Mode:         "ident",
				})
				if err != nil {
					log.Panic(err)
				}
			case *ast.ArrayType:
				res := getTypeString(currType.Elt)
				// needsOptions := !res.IsBasicType
				needsOptions := true
				needsOptionsIndex := false
				if res.Name == "Index" {
					needsOptionsIndex = true
				}
				err = generateEqualAndDiff(generateEqualAndDiffOptions{
					PackageName:       packageName,
					File:              file,
					FileTest:          fileTest,
					Name:              currSpecType.Name.Name,
					Type:              res.Name,
					CurrType:          currSpecType,
					IsBasicType:       res.IsBasicType,
					IsComplex:         res.IsComplex,
					IsComparable:      res.IsComparable,
					IsPointer:         strings.HasPrefix(res.Name, "*"),
					NeedsOptions:      needsOptions,
					NeedsOptionsIndex: needsOptionsIndex,
					Mode:              "array",
				})
				if err != nil {
					log.Panic(err)
				}
			}
		}
	}

	if hasTests {
		err = fmtFile(generatedFileNameTest)
		if err != nil {
			return packageName, err
		}
	} else {
		os.Remove(generatedFileNameTest)
	}
	// Format the file
	err = fmtFile(generatedFileName)
	if err != nil {
		return packageName, err
	}
	return packageName, nil
}
