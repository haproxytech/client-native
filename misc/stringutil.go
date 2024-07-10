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
	"fmt"
	"math/rand"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var chars = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ") //nolint:gochecknoglobals

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
	n = strings.ReplaceAll(n, "Http", "HTTP")
	n = strings.ReplaceAll(n, "Uri", "URI")
	n = strings.ReplaceAll(n, "http", "HTTP")
	n = strings.ReplaceAll(n, "tcp", "TCP")
	n = strings.ReplaceAll(n, "Tcp", "TCP")
	n = strings.ReplaceAll(n, "Id", "ID")
	n = strings.ReplaceAll(n, "Tls", "TLS")
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

		switch {
		case i > 0 && n[len(n)-1] != '_' && nextCaseIsChanged:
			// add underscore if next letter case type is changed
			if v >= 'A' && v <= 'Z' {
				n += "_" + string(v)
			} else if v >= 'a' && v <= 'z' {
				n += string(v) + "_"
			}
		case v == ' ':
			// replace spaces with underscores
			n += "_"
		default:
			n += string(v)
		}
	}
	n = strings.ToLower(n)
	// special case
	n = strings.ReplaceAll(n, "httpuri", "http_uri")
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

		switch {
		case i > 0 && n[len(n)-1] != '-' && nextCaseIsChanged:
			// add underscore if next letter case type is changed
			if v >= 'A' && v <= 'Z' {
				n += "-" + string(v)
			} else if v >= 'a' && v <= 'z' {
				n += string(v) + "-"
			}
		case v == ' ':
			// replace spaces with underscores
			n += "-"
		default:
			n += string(v)
		}
	}
	n = strings.ToLower(n)
	// special case
	n = strings.ReplaceAll(n, "httpuri", "http-uri")
	return n
}

// ParseTimeout returns the number of milliseconds in a timeout string.
func ParseTimeout(tOut string) *int64 {
	return parseTimeout(tOut, 1)
}

func ParseTimeoutDefaultSeconds(tOut string) *int64 {
	return parseTimeout(tOut, 1000)
}

func parseTimeout(tOut string, defaultMultiplier int64) *int64 {
	var v int64
	var err error
	switch {
	case strings.HasSuffix(tOut, "us"):
		v, err = strconv.ParseInt(strings.TrimSuffix(tOut, "us"), 10, 64)
		if v >= 1000 {
			v /= 1000
		} else if v > 0 {
			v = 1
		}
	case strings.HasSuffix(tOut, "ms"):
		v, err = strconv.ParseInt(strings.TrimSuffix(tOut, "ms"), 10, 64)
	case strings.HasSuffix(tOut, "s"):
		v, err = strconv.ParseInt(strings.TrimSuffix(tOut, "s"), 10, 64)
		v *= 1000
	case strings.HasSuffix(tOut, "m"):
		v, err = strconv.ParseInt(strings.TrimSuffix(tOut, "m"), 10, 64)
		v = v * 1000 * 60
	case strings.HasSuffix(tOut, "h"):
		v, err = strconv.ParseInt(strings.TrimSuffix(tOut, "h"), 10, 64)
		v = v * 1000 * 60 * 60
	case strings.HasSuffix(tOut, "d"):
		v, err = strconv.ParseInt(strings.TrimSuffix(tOut, "d"), 10, 64)
		v = v * 1000 * 60 * 60 * 24
	default:
		v, err = strconv.ParseInt(tOut, 10, 64)
		v *= defaultMultiplier
	}
	if err != nil || v < 0 {
		return nil
	}
	return &v
}

func ParseSize(size string) *int64 {
	var v int64
	var err error
	switch {
	case strings.HasSuffix(size, "k"):
		v, err = strconv.ParseInt(strings.TrimSuffix(size, "k"), 10, 64)
		v *= 1024
	case strings.HasSuffix(size, "K"):
		v, err = strconv.ParseInt(strings.TrimSuffix(size, "K"), 10, 64)
		v *= 1024
	case strings.HasSuffix(size, "m"):
		v, err = strconv.ParseInt(strings.TrimSuffix(size, "m"), 10, 64)
		v = v * 1024 * 1024
	case strings.HasSuffix(size, "M"):
		v, err = strconv.ParseInt(strings.TrimSuffix(size, "M"), 10, 64)
		v = v * 1024 * 1024
	case strings.HasSuffix(size, "g"):
		v, err = strconv.ParseInt(strings.TrimSuffix(size, "g"), 10, 64)
		v = v * 1024 * 1024 * 1024
	case strings.HasSuffix(size, "G"):
		v, err = strconv.ParseInt(strings.TrimSuffix(size, "G"), 10, 64)
		v = v * 1024 * 1024 * 1024
	default:
		v, err = strconv.ParseInt(size, 10, 64)
	}
	if err != nil {
		return nil
	}
	return &v
}

func StringP(s string) *string {
	return &s
}

func Int64P(i int) *int64 {
	ret := int64(i)
	return &ret
}

func BoolP(b bool) *bool {
	return &b
}

func RandomString(n int) string {
	b := make([]rune, n)
	size := len(chars)
	for i := range b {
		b[i] = chars[rand.Intn(size)] //nolint:gosec
	}
	return string(b)
}

// SanitizeFilename collapses paths and replaces most non-alphanumeric characters with underscores
func SanitizeFilename(name string) string {
	var ext string

	// save the extension if possible
	ext = filepath.Ext(name)
	name = name[:len(name)-len(ext)]

	if name != "" {
		// collapse paths
		name = filepath.Clean(name)
	}
	// leave all alphanumeric and 3 additional ones
	// # _ -
	reg := regexp.MustCompile(`[^a-zA-Z0-9#_\\-]+`)
	name = reg.ReplaceAllString(name, "_")

	if ext != "" {
		ext = reg.ReplaceAllString(ext[1:], "_")
		if name != "" {
			return fmt.Sprintf("%s.%s", name, ext)
		}
		return fmt.Sprintf("_%s", ext)
	}

	return name
}
