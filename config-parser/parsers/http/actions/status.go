/*
Copyright 2022 HAProxy Technologies

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package actions

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/haproxytech/client-native/v5/config-parser/types"
)

type Status struct {
	Comment       string
	Status        *int64
	ContentType   string
	ContentFormat string
	Content       string
	Hdrs          []*Hdr
}

var allowedErrorStatusCodes = map[int64]struct{}{ //nolint:gochecknoglobals
	200: {},
	400: {},
	401: {},
	403: {},
	404: {},
	405: {},
	407: {},
	408: {},
	410: {},
	413: {},
	425: {},
	429: {},
	500: {},
	501: {},
	502: {},
	503: {},
	504: {},
}

func AllowedErrorStatusCode(code int64) bool {
	_, ok := allowedErrorStatusCodes[code]
	return ok
}

// Parse parses http-error status <code> [content-type <type>] [ { default-errorfiles | errorfile <file> | errorfiles <name> | file <file> | lf-file <file> | string <str> | lf-string <fmt> } ] [ hdr <name> <fmt> ]*
func (f *Status) Parse(parts []string, parserType types.ParserType, comment string) error {
	f.Hdrs = []*Hdr{}
	if comment != "" {
		f.Comment = comment
	}

	// Parsing specific to http-error status directive.
	if len(parts) < 3 {
		return fmt.Errorf("not enough params")
	}
	if parts[0] != "http-error" {
		return fmt.Errorf("unexpected keyword %s", parts[0])
	}
	if parts[1] != "status" {
		return fmt.Errorf("unsupported action %s", parts[1])
	}
	code, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse status code")
	}
	if !AllowedErrorStatusCode(code) {
		return fmt.Errorf("unsupported status code %d", code)
	}
	f.Status = &code

	// Process optional arguments.
	if command := parts[3:]; len(command) >= 1 {
		for i := 0; i < len(command); i++ {
			switch command[i] {
			case "content-type":
				i++
				f.ContentType = command[i]
			case "errorfile", "errorfiles", "file", "lf-file", "string", "lf-string":
				f.ContentFormat = command[i]
				i++
				f.Content = command[i]
			case "default-errorfiles":
				f.ContentFormat = command[i]
			case "hdr":
				hdr := Hdr{}
				if len(command) < i+3 {
					return fmt.Errorf("failed to parse return hdr")
				}
				i++
				hdr.Name = command[i]
				i++
				hdr.Fmt = command[i]
				f.Hdrs = append(f.Hdrs, &hdr)
			default:
				return fmt.Errorf("unsupported keyword %s", command[i])
			}
		}
	}
	return nil
}

func (f *Status) String() string {
	var result strings.Builder
	if f.Status != nil {
		result.WriteString("status ")
		result.WriteString(strconv.FormatInt(*f.Status, 10))
	}
	if f.ContentType != "" {
		result.WriteString(" content-type ")
		result.WriteString(f.ContentType)
	}
	if f.ContentFormat != "" {
		result.WriteString(" ")
		result.WriteString(f.ContentFormat)
		if f.Content != "" && f.ContentFormat != "default-errorfiles" {
			result.WriteString(" ")
			result.WriteString(f.Content)
		}
	}
	if IsPayload(f.ContentFormat) {
		for _, hdr := range f.Hdrs {
			result.WriteString(" hdr ")
			result.WriteString(hdr.Name)
			result.WriteString(" ")
			result.WriteString(hdr.Fmt)
		}
	}
	return result.String()
}

func (f *Status) GetComment() string {
	return f.Comment
}
