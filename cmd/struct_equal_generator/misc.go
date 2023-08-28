package main

import (
	"go/ast"
	"os"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Field struct {
	Name         string
	Type         string
	TypeInFile   string
	IsBasicType  bool
	IsComparable bool
	IsEmbedded   bool
	HasString    bool
	HasEqual     bool
	HasEqualOpt  bool
	IsArray      bool
	IsMap        bool
	SubType      *Field
}

type generateEqualAndDiffOptions struct {
	File              *os.File
	FileTest          *os.File
	CurrType          *ast.TypeSpec
	PackageName       string
	Fields            []Field
	Name              string
	Type              string
	IsPointer         bool
	NeedsOptions      bool
	NeedsOptionsIndex bool
	IsBasicType       bool
	IsComplex         bool
	IsComparable      bool
	Mode              string
}

func toTitle(s string) string {
	caser := cases.Title(language.Und)
	return caser.String(strings.TrimPrefix(s, "*"))
}

func toJSON(x any) string {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	b, err := json.Marshal(x)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func toCamelCase(s string) string {
	caser := cases.Title(language.Und)
	var result string
	words := strings.Split(s, "_")

	for _, word := range words {
		result += caser.String(word)
	}

	result = strings.ToLower(result[:1]) + result[1:]

	return result
}

// isComparable checks if a given type is comparable.
// It takes in the name of the type as a string and returns a boolean value.
func isComparable(typeName string) bool {
	return basicTypes[typeName]
}

func getFulPath(typeName string, imports map[string]string) string {
	parts := strings.Split(typeName, ".")
	if len(parts) > 1 {
		path, ok := imports[parts[0]]
		if ok {
			typeName = path + "." + parts[1]
		}
	}
	return typeName
}

func hasEqual(typeName string) bool {
	return hasEqualTypes[typeName]
}

func hasEqualOpt(typeName string) bool {
	return generatedTypes[typeName]
}

type getTypeStringResponse struct {
	Name         string
	TypeName     string
	IsBasicType  bool
	IsComplex    bool
	IsComparable bool
	HasStringer  bool
	HasEqual     bool
	HasEqualOpt  bool
	IsArray      bool
	IsMap        bool
	SubType      *getTypeStringResponse
}

func getTypeString(expr ast.Expr, imports map[string]string) getTypeStringResponse {
	switch t := expr.(type) {
	case *ast.ArrayType:
		res := getTypeString(t.Elt, imports)
		return getTypeStringResponse{
			Name:         "[]" + res.Name,
			IsComplex:    !basicTypes[res.Name],
			IsComparable: false,
			IsArray:      true,
			SubType: &getTypeStringResponse{
				Name:         res.Name,
				TypeName:     res.TypeName,
				IsBasicType:  res.IsBasicType,
				IsComplex:    res.IsComplex,
				IsComparable: res.IsComparable,
				IsArray:      res.IsArray,
				IsMap:        res.IsMap,
				HasStringer:  res.HasStringer,
				HasEqual:     res.HasEqual,
				HasEqualOpt:  res.HasEqualOpt,
				SubType:      res.SubType,
			},
		}
	case *ast.Ident:
		if t.Obj == nil {
			return getTypeStringResponse{
				Name:         t.Name,
				IsBasicType:  basicTypes[t.Name],
				IsComplex:    !basicTypes[t.Name],
				IsComparable: isComparable(t.Name),
				// HasStringer:  hasStringer(getFulPath(t.Name, imports)),
				HasEqual:    hasEqual(getFulPath(t.Name, imports)),
				HasEqualOpt: hasEqualOpt(getFulPath(t.Name, imports)),
			}
		}
		return getTypeStringResponse{
			Name:         t.Obj.Name,
			IsBasicType:  basicTypes[t.Name],
			IsComplex:    true,
			IsComparable: isComparable(t.Name),
			// HasStringer:  hasStringer(getFulPath(t.Name, imports)),
			HasEqual:    hasEqual(getFulPath(t.Name, imports)),
			HasEqualOpt: hasEqualOpt(getFulPath(t.Name, imports)),
		}
	case *ast.MapType:
		rKey := getTypeString(t.Key, imports)
		rValue := getTypeString(t.Value, imports)
		return getTypeStringResponse{
			Name:      "map[" + rKey.Name + "]" + rValue.Name,
			IsComplex: rValue.IsComplex,
			IsMap:     true,
		}
	case *ast.SelectorExpr:
		start := expr.Pos() - 1
		end := expr.End() - 1
		typeInSource := sourceOfFile[start:end]
		return getTypeStringResponse{
			Name:         t.Sel.Name,
			TypeName:     typeInSource,
			IsComparable: isComparable(getFulPath(typeInSource, imports)),
			////HasStringer:  hasStringer(getFulPath(typeInSource, imports)),
			// HasStringer: hasStringer2(getPackage(typeInSource, imports), typeInSource),
			HasEqual:    hasEqual(getFulPath(typeInSource, imports)),
			HasEqualOpt: hasEqualOpt(getFulPath(typeInSource, imports)),
		}
	case *ast.StarExpr:
		res := getTypeString(t.X, imports)
		res.Name = "*" + res.Name
		return res
	}
	return getTypeStringResponse{}
}

var generatedTypes = map[string]bool{} //nolint:gochecknoglobals

var hasEqualTypes = map[string]bool{ //nolint:gochecknoglobals
	"time.Time":                             true,
	"github.com/go-openapi/strfmt.DateTime": true,
	"github.com/go-openapi/strfmt.Date":     true,
}

var basicTypes = map[string]bool{ //nolint:gochecknoglobals
	"bool":       true,
	"byte":       true,
	"complex64":  true,
	"complex128": true,
	"error":      true,
	"float32":    true,
	"float64":    true,
	"int":        true,
	"int8":       true,
	"int16":      true,
	"int32":      true,
	"int64":      true,
	"rune":       true,
	"string":     true,
	"uint":       true,
	"uint8":      true,
	"uint16":     true,
	"uint32":     true,
	"uint64":     true,
	"uintptr":    true,
	// custom types
	"github.com/go-openapi/strfmt.Password": true,
}
