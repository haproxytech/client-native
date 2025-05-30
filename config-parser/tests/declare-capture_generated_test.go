// Code generated by go generate; DO NOT EDIT.
/*
Copyright 2019 HAProxy Technologies

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

package tests

import (
	"fmt"
	"strings"
	"testing"

	"github.com/haproxytech/client-native/v6/config-parser/parsers"
)

func TestDeclareCapture(t *testing.T) {
	tests := map[string]bool{
		"declare capture request len 1":             true,
		"declare capture response len 2":            true,
		"declare capture":                           false,
		"declare capture request":                   false,
		"declare capture response":                  false,
		"declare capture something":                 false,
		"declare capture request len":               false,
		"declare capture response 1":                false,
		"declare capture response a":                false,
		"declare capture response len":              false,
		"declare capture response length":           false,
		"declare capture response length 1":         false,
		"declare capture response something":        false,
		"declare capture response something 1":      false,
		"declare capture request len abcd":          false,
		"declare capture response len abcd":         false,
		"declare capture request len 1 if true":     false,
		"declare capture response len 2 if true":    false,
		"declare capture request len abcd if true":  false,
		"declare capture response len abcd if true": false,
		"---":     false,
		"--- ---": false,
	}
	parser := &parsers.DeclareCapture{}
	for command, shouldPass := range tests {
		t.Run(command, func(t *testing.T) {
			line := strings.TrimSpace(command)
			lines := strings.SplitN(line, "\n", -1)
			var err error
			parser.Init()
			if len(lines) > 1 {
				for _, line = range lines {
					line = strings.TrimSpace(line)
					if err = ProcessLine(line, parser); err != nil {
						break
					}
				}
			} else {
				err = ProcessLine(line, parser)
			}
			if shouldPass {
				if err != nil {
					t.Error(err)
					return
				}
				result, err := parser.Result()
				if err != nil {
					t.Error(err)
					return
				}
				var returnLine string
				if result[0].Comment == "" {
					returnLine = result[0].Data
				} else {
					returnLine = fmt.Sprintf("%s # %s", result[0].Data, result[0].Comment)
				}
				if command != returnLine {
					t.Errorf("error: has [%s] expects [%s]", returnLine, command)
				}
			} else {
				if err == nil {
					t.Errorf("error: did not throw error for line [%s]", line)
				}
				_, parseErr := parser.Result()
				if parseErr == nil {
					t.Errorf("error: did not throw error on result for line [%s]", line)
				}
			}
		})
	}
}
