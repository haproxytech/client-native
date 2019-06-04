// Copyright 2019 HAProxy Technologies
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package misc

import (
	"strconv"
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
	// special cases
	n = strings.Replace(n, "Http", "HTTP", -1)
	n = strings.Replace(n, "Uri", "URI", -1)
	n = strings.Replace(n, "http", "HTTP", -1)
	n = strings.Replace(n, "tcp", "TCP", -1)
	n = strings.Replace(n, "Tcp", "TCP", -1)
	n = strings.Replace(n, "Id", "ID", -1)
	n = strings.Replace(n, "Tls", "TLS", -1)
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

// DashCase turns camel case to snake case string
func DashCase(fieldName string) string {
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

		if i > 0 && n[len(n)-1] != '-' && nextCaseIsChanged {
			// add underscore if next letter case type is changed
			if v >= 'A' && v <= 'Z' {
				n += "-" + string(v)
			} else if v >= 'a' && v <= 'z' {
				n += string(v) + "-"
			}
		} else if v == ' ' {
			// replace spaces with underscores
			n += "-"
		} else {
			n = n + string(v)
		}
	}
	n = strings.ToLower(n)
	// special case
	n = strings.Replace(n, "httpuri", "http-uri", -1)
	return n
}

func ParseTimeout(tOut string) *int64 {
	var v int64
	if strings.HasSuffix(tOut, "ms") {
		v, _ = strconv.ParseInt(strings.TrimSuffix(tOut, "ms"), 10, 64)
	} else if strings.HasSuffix(tOut, "s") {
		v, _ = strconv.ParseInt(strings.TrimSuffix(tOut, "s"), 10, 64)
		v = v * 1000
	} else if strings.HasSuffix(tOut, "m") {
		v, _ = strconv.ParseInt(strings.TrimSuffix(tOut, "m"), 10, 64)
		v = v * 1000 * 60
	} else if strings.HasSuffix(tOut, "h") {
		v, _ = strconv.ParseInt(strings.TrimSuffix(tOut, "h"), 10, 64)
		v = v * 1000 * 60 * 60
	} else if strings.HasSuffix(tOut, "d") {
		v, _ = strconv.ParseInt(strings.TrimSuffix(tOut, "d"), 10, 64)
		v = v * 1000 * 60 * 60 * 24
	} else {
		v, _ = strconv.ParseInt(tOut, 10, 64)
	}
	if v != 0 {
		return &v
	}
	return nil
}

func ParseSize(size string) *int64 {
	var v int64
	if strings.HasSuffix(size, "k") {
		v, _ = strconv.ParseInt(strings.TrimSuffix(size, "k"), 10, 64)
		v = v * 1024
	} else if strings.HasSuffix(size, "m") {
		v, _ = strconv.ParseInt(strings.TrimSuffix(size, "m"), 10, 64)
		v = v * 1024 * 1024
	} else if strings.HasSuffix(size, "g") {
		v, _ = strconv.ParseInt(strings.TrimSuffix(size, "g"), 10, 64)
		v = v * 1024 * 1024 * 1024
	} else {
		v, _ = strconv.ParseInt(size, 10, 64)
	}
	if v != 0 {
		return &v
	}
	return nil
}
