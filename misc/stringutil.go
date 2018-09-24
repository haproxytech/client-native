package misc

import (
	"strings"
)

// StringInSlice checks if a string is in a list of strings
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// CamelCase turns snake case to camel case string
func CamelCase(fieldName string, initCase bool) string {
	s := strings.Trim(fieldName, " ")
	n := ""
	capNext := initCase
	for _, v := range s {
		if v >= 'A' && v <= 'Z' {
			n += string(v)
		}
		if v >= '0' && v <= '9' {
			n += string(v)
		}
		if v >= 'a' && v <= 'z' {
			if capNext {
				n += strings.ToUpper(string(v))
			} else {
				n += string(v)
			}
		}
		if v == '_' || v == ' ' || v == '-' {
			capNext = true
		} else {
			capNext = false
		}
	}
	n = strings.Replace(n, "Http", "HTTP", -1)
	n = strings.Replace(n, "Uri", "URI", -1)
	n = strings.Replace(n, "http", "HTTP", -1)
	n = strings.Replace(n, "tcp", "TCP", -1)
	n = strings.Replace(n, "Tcp", "TCP", -1)
	n = strings.Replace(n, "Id", "ID", -1)
	return n
}

// SnakeCase turns camel case to snake case string
func SnakeCase(fieldName string) string {
	fieldName = strings.Trim(fieldName, " ")
	n := ""
	for i, v := range fieldName {
		// treat acronyms as words, eg for JSONData -> JSON is a whole word
		nextCaseIsChanged := false
		if i+1 < len(fieldName) {
			next := fieldName[i+1]
			if (v >= 'A' && v <= 'Z' && next >= 'a' && next <= 'z') || (v >= 'a' && v <= 'z' && next >= 'A' && next <= 'Z') {
				nextCaseIsChanged = true
			}
		}

		if i > 0 && n[len(n)-1] != '_' && nextCaseIsChanged {
			// add underscore if next letter case type is changed
			if v >= 'A' && v <= 'Z' {
				n += "_" + string(v)
			} else if v >= 'a' && v <= 'z' {
				n += string(v) + "_"
			}
		} else if v == ' ' {
			// replace spaces with underscores
			n += "_"
		} else {
			n = n + string(v)
		}
	}
	n = strings.ToLower(n)
	// special case
	n = strings.Replace(n, "httpuri", "http_uri", -1)
	return n
}
