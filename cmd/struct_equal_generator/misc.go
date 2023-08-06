package main

import (
	"go/ast"
	"os"
	"reflect"
	"strings"

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
	// Get the type object using reflection.
	typeObj := reflect.TypeOf(typeName)

	// Check if the type is comparable.
	return typeObj.Comparable()
}

type getTypeStringResponse struct {
	Name         string
	TypeName     string
	IsBasicType  bool
	IsComplex    bool
	IsComparable bool
}

func getTypeString(expr ast.Expr) getTypeStringResponse {
	switch t := expr.(type) {
	case *ast.ArrayType:
		res := getTypeString(t.Elt)
		return getTypeStringResponse{
			Name:         "[]" + res.Name,
			IsComplex:    !basicTypes[res.Name],
			IsComparable: true,
		}
	case *ast.Ident:
		if t.Obj == nil {
			return getTypeStringResponse{
				Name:         t.Name,
				IsBasicType:  basicTypes[t.Name],
				IsComplex:    !basicTypes[t.Name],
				IsComparable: isComparable(t.Name),
			}
		}
		return getTypeStringResponse{
			Name:         t.Obj.Name,
			IsBasicType:  basicTypes[t.Name],
			IsComplex:    true,
			IsComparable: isComparable(t.Name),
		}
	case *ast.MapType:
		rKey := getTypeString(t.Key)
		rValue := getTypeString(t.Value)
		return getTypeStringResponse{
			Name:      "map[" + rKey.Name + "]" + rValue.Name,
			IsComplex: rValue.IsComplex,
		}
	case *ast.SelectorExpr:
		start := expr.Pos() - 1
		end := expr.End() - 1
		typeInSource := sourceofFile[start:end]
		return getTypeStringResponse{Name: t.Sel.Name, TypeName: typeInSource}
	case *ast.StarExpr:
		res := getTypeString(t.X)
		res.Name = "*" + res.Name
		return res
	}
	return getTypeStringResponse{}
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
}
