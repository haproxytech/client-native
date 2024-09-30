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
	"strings"

	"github.com/haproxytech/client-native/v5/config-parser/common"
	"github.com/haproxytech/client-native/v5/config-parser/types"
)

type NormalizeURI struct {
	Normalizer string
	Full       bool
	Strict     bool
	Cond       string
	CondTest   string
	Comment    string
}

func (f *NormalizeURI) Parse(parts []string, parserType types.ParserType, comment string) error {
	if comment != "" {
		f.Comment = comment
	}
	if len(parts) >= 3 {
		command, condition := common.SplitRequest(parts[2:])
		if len(command) > 0 {
			switch command[0] {
			case "path-strip-dotdot":
				f.Normalizer = command[0]
				if len(command) > 1 && command[1] == "full" {
					f.Full = true
				}
			case "percent-decode-unreserved", "percent-to-uppercase":
				f.Normalizer = command[0]
				if len(command) > 1 && command[1] == "strict" {
					f.Strict = true
				}
			case "fragment-encode", "fragment-strip", "path-merge-slashes", "path-strip-dot", "query-sort-by-name":
				f.Normalizer = command[0]
				if len(command) > 1 {
					return fmt.Errorf("unsupported keyword for %s: %s", command[0], command[1:])
				}
			default:
				return fmt.Errorf("unrecognized normalizer %s", command[0])
			}
			if len(condition) > 1 {
				f.Cond = condition[0]
				f.CondTest = strings.Join(condition[1:], " ")
			}
			return nil
		}
	}

	return fmt.Errorf("not enough params")
}

func (f *NormalizeURI) String() string {
	var result strings.Builder
	result.WriteString("normalize-uri ")
	result.WriteString(f.Normalizer)
	if f.Strict {
		result.WriteString(" strict")
	}
	if f.Full {
		result.WriteString(" full")
	}
	if f.Cond != "" {
		result.WriteString(" ")
		result.WriteString(f.Cond)
		result.WriteString(" ")
		result.WriteString(f.CondTest)
	}
	return result.String()
}

func (f *NormalizeURI) GetComment() string {
	return f.Comment
}
