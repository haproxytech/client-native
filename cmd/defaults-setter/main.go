package main

import (
	"fmt"
	"go/ast"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/go-openapi/swag/mangling"
	"golang.org/x/tools/go/packages"
	"gopkg.in/yaml.v3"
)

// Spec represents the structure of the OpenAPI specification
type Spec struct {
	Definitions map[string]Definition `yaml:"definitions"`
}

// Definition represents a model definition in the OpenAPI specification
type Definition struct {
	Properties map[string]Property `yaml:"properties"`
	GoType     string              `yaml:"x-go-name"` //nolint:tagliatelle
}

// Property represents a property of a model in the OpenAPI specification
type Property struct {
	Default any    `yaml:"default"`
	Type    string `yaml:"type"`
	GoType  string `yaml:"x-go-name"` //nolint:tagliatelle
}

type FieldInfo struct {
	IsPointer bool
}

type FileStruct struct {
	FileName string
	Fields   map[string]FieldInfo
}

func main() {
	if len(os.Args) != 3 {
		log.Fatalf("Usage: %s <specification_file> <models_directory>", os.Args[0])
	}
	specFile := os.Args[1]
	modelsDir := os.Args[2]
	// 1. Parse all OpenAPI specifications
	goStructsToChange := make(map[string]map[string]struct{})
	allDefinitions := make(map[string]Definition)
	spec, err := parseSpec(specFile)
	if err != nil {
		log.Fatalf("Failed to parse %s: %v", specFile, err)
	}
	for key, definition := range spec.Definitions {
		properties := make(map[string]Property)
		fields := make(map[string]struct{})
		for prop, value := range definition.Properties {
			if value.Default != nil {
				properties[prop] = value
				fields[toGoName(prop, "")] = struct{}{}
			}
		}
		if len(properties) > 0 {
			allDefinitions[key] = Definition{Properties: properties, GoType: definition.GoType}
			goStructsToChange[toGoName(key, definition.GoType)] = fields
		}
	}

	modelsMap := initModelsMap(goStructsToChange)

	// 2. Iterate through all definitions and find corresponding Go files
	for schemaName, definition := range allDefinitions {
		structName := toGoName(schemaName, definition.GoType)
		goType, ok := modelsMap[structName]
		if !ok {
			continue
		}

		// Find the file that defines this struct
		filePath := filepath.Join(modelsDir, goType.FileName)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			continue
		}

		// Check if the struct has a Validate function
		hasValidate, err := hasValidateFunc(filePath, structName)
		if err != nil || !hasValidate {
			continue
		}

		// Generate and insert default value assignments
		assignments := generateAssignments(definition.Properties, structName, goType)
		if len(assignments) == 0 {
			continue
		}

		log.Printf("Processing %s in %s", structName, filePath)
		if err := insertCode(filePath, structName, assignments, hasPointers(goType)); err != nil {
			log.Printf("Failed to insert code into %s for struct %s: %v", filePath, structName, err)
		}
	}
}

func hasPointers(fs FileStruct) bool {
	for _, v := range fs.Fields {
		if v.IsPointer {
			return true
		}
	}
	return false
}

// parseSpec reads and parses the main OpenAPI specification file
func parseSpec(filepath string) (*Spec, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	var spec Spec
	if err := yaml.Unmarshal(data, &spec); err != nil {
		return nil, err
	}
	return &spec, nil
}

// hasValidateFunc checks if a file contains a Validate function for a given struct.
func hasValidateFunc(filePath, structName string) (bool, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return false, err
	}
	validateRegex := regexp.MustCompile(fmt.Sprintf(`func \(m \*%s\) Validate\(formats strfmt\.Registry\) error`, structName))
	return validateRegex.Match(content), nil
}

// generateAssignments generates the Go code for default value assignments
func generateAssignments(properties map[string]Property, structName string, fileStruct FileStruct) []string {
	var assignments []string

	// Extract and sort the keys
	keys := make([]string, 0, len(properties))
	for name := range properties {
		keys = append(keys, name)
	}
	sort.Strings(keys)

	for _, name := range keys {
		prop := properties[name]
		if prop.Default != nil {
			goFieldName := toGoName(name, prop.GoType)

			field, ok := fileStruct.Fields[goFieldName]
			if !ok {
				log.Printf("Field %s not found in struct %s", goFieldName, structName)
				continue
			}

			assignment := strings.Builder{}
			assignment.WriteString("\tif swag.IsZero(m.")
			assignment.WriteString(goFieldName)
			assignment.WriteString(") {\n\t\t")
			if field.IsPointer {
				switch prop.Type {
				case "string":
					assignment.WriteString("m.")
					assignment.WriteString(goFieldName)
					assignment.WriteString(` = misc.Ptr("`)
					assignment.WriteString(fmt.Sprintf("%v", prop.Default))
					assignment.WriteString(`")`)
				case "integer":
					assignment.WriteString("m.")
					assignment.WriteString(goFieldName)
					assignment.WriteString(" = misc.Ptr(int64(")
					assignment.WriteString(fmt.Sprintf("%v", prop.Default))
					assignment.WriteString("))")
				case "number":
					assignment.WriteString("m.")
					assignment.WriteString(goFieldName)
					assignment.WriteString(" = misc.Ptr(float64(")
					assignment.WriteString(fmt.Sprintf("%v", prop.Default))
					assignment.WriteString("))")
				case "boolean":
					assignment.WriteString("m.")
					assignment.WriteString(goFieldName)
					assignment.WriteString(" = misc.Ptr(")
					assignment.WriteString(fmt.Sprintf("%v", prop.Default))
					assignment.WriteString(")")
				default:
					log.Printf("unhandled pointer type for field %s", goFieldName)
					assignment.WriteString("m.")
					assignment.WriteString(goFieldName)
					assignment.WriteString(" = misc.Ptr(")
					assignment.WriteString(fmt.Sprintf("%v", prop.Default))
					assignment.WriteString(")")
				}
			} else {
				switch prop.Type {
				case "string":
					assignment.WriteString("m.")
					assignment.WriteString(goFieldName)
					assignment.WriteString(` = "`)
					assignment.WriteString(fmt.Sprintf("%v", prop.Default))
					assignment.WriteString(`"`)
				default:
					assignment.WriteString("m.")
					assignment.WriteString(goFieldName)
					assignment.WriteString(" = ")
					assignment.WriteString(fmt.Sprintf("%v", prop.Default))
				}
			}
			assignment.WriteString("\n\t}")
			assignments = append(assignments, assignment.String())
		}
	}
	return assignments
}

// insertCode inserts the generated assignments into the Validate function
func insertCode(filepath, structName string, assignments []string, hasPointers bool) error { //nolint:gocognit
	content, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}
	lines := strings.Split(string(content), "\n")

	// Find import block and add misc import if it doesn't exist
	importBlockStart := -1
	miscImportExists := false
	for i, line := range lines {
		if strings.Contains(line, "\"github.com/haproxytech/client-native/v6/misc\"") {
			miscImportExists = true
			break
		}
		if strings.HasPrefix(line, "import (") {
			importBlockStart = i
		}
	}

	if !miscImportExists && importBlockStart != -1 && hasPointers {
		importLine := "\t\"github.com/haproxytech/client-native/v6/misc\""
		lines = append(lines[:importBlockStart+1], append([]string{importLine}, lines[importBlockStart+1:]...)...)
	}

	// Find Validate function and remove existing default assignments
	validateFuncStart := -1
	validateFuncRegex := regexp.MustCompile(fmt.Sprintf(`func \(m \*%s\) Validate\(formats strfmt\.Registry\) error \{`, structName))
	for i, line := range lines {
		if validateFuncRegex.MatchString(line) {
			validateFuncStart = i
			break
		}
	}

	if validateFuncStart != -1 {
		var newLines []string
		inValidateFunc := false
		for i, line := range lines {
			if i == validateFuncStart {
				inValidateFunc = true
			}
			if inValidateFunc && strings.HasPrefix(line, "func ") {
				inValidateFunc = false
			}
			// Remove existing default assignments
			isDefaultAssignment := false
			for _, assignment := range assignments {
				parts := strings.Split(assignment, " = ")
				if len(parts) > 0 && strings.Contains(line, parts[0]+" = ") {
					isDefaultAssignment = true
					break
				}
			}
			if !isDefaultAssignment {
				newLines = append(newLines, line)
			}
		}
		lines = newLines

		// Re-find the start of the validate function
		for i, line := range lines {
			if validateFuncRegex.MatchString(line) {
				validateFuncStart = i
				break
			}
		}

		var newAssignments []string //nolint: prealloc
		newAssignments = append(newAssignments, assignments...)
		newAssignments = append(newAssignments, "")

		lines = append(lines[:validateFuncStart+1], append(newAssignments, lines[validateFuncStart+1:]...)...)
	} else {
		return fmt.Errorf("could not find insertion point in %s for struct %s", filepath, structName)
	}

	if err := os.WriteFile(filepath, []byte(strings.Join(lines, "\n")), 0o600); err != nil {
		return err
	}

	return nil
}

func initModelsMap(structs map[string]map[string]struct{}) map[string]FileStruct { //nolint:gocognit
	retMap := make(map[string]FileStruct)
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedSyntax,
	}
	pkgs, err := packages.Load(cfg, "./models")
	if err != nil {
		log.Fatalf("failed to load models package: %v", err)
	}
	if len(pkgs) != 1 {
		log.Fatalf("expected 1 package, found %d", len(pkgs))
	}
	pkg := pkgs[0]

	for _, file := range pkg.Syntax {
		fileName := pkg.Fset.File(file.Pos()).Name()
		if strings.HasSuffix(fileName, "_compare.go") || strings.HasSuffix(fileName, "_compare_test.go") {
			continue
		}

		ast.Inspect(file, func(n ast.Node) bool {
			if typeSpec, ok := n.(*ast.TypeSpec); ok {
				if structType, ok := typeSpec.Type.(*ast.StructType); ok {
					structName := typeSpec.Name.Name
					if _, ok := structs[structName]; ok {
						// Check for Validate method
						if hasValidateMethod(pkg, structName) {
							fields := make(map[string]FieldInfo)
							for _, field := range structType.Fields.List {
								if len(field.Names) > 0 {
									fieldName := field.Names[0].Name
									if _, ok := structs[structName][fieldName]; ok {
										_, isPointer := field.Type.(*ast.StarExpr)
										fields[fieldName] = FieldInfo{IsPointer: isPointer}
									}
								}
							}
							retMap[structName] = FileStruct{
								FileName: filepath.Base(fileName),
								Fields:   fields,
							}
						}
					}
				}
			}
			return true
		})
	}
	return retMap
}

func hasValidateMethod(pkg *packages.Package, structName string) bool {
	for _, file := range pkg.Syntax {
		for _, decl := range file.Decls {
			if fn, ok := decl.(*ast.FuncDecl); ok {
				if fn.Name.Name == "Validate" && fn.Recv != nil && len(fn.Recv.List) > 0 {
					if starExpr, ok := fn.Recv.List[0].Type.(*ast.StarExpr); ok {
						if ident, ok := starExpr.X.(*ast.Ident); ok {
							if ident.Name == structName {
								return true
							}
						}
					}
				}
			}
		}
	}
	return false
}

func toGoName(key string, xGoType string) string {
	if xGoType != "" {
		return xGoType
	}
	mangler := mangling.NewNameMangler(mangling.WithAdditionalInitialisms("QUIC", "FCGI", "VRRP"))
	return mangler.ToGoName(key)
}
