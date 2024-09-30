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

package params

import (
	"strings"
)

// DgramBindOption ...
type DgramBindOption interface {
	Parse(options []string, currentIndex int) (int, error)
	Valid() bool
	String() string
}

func getDgramBindOptions() []DgramBindOption {
	return []DgramBindOption{
		&BindOptionWord{Name: "transparent"},
		&BindOptionValue{Name: "interface"},
		&BindOptionValue{Name: "namespace"},
		&BindOptionValue{Name: "name"},
	}
}

// Parse ...
func ParseDgramBindOptions(options []string) []DgramBindOption {
	result := []DgramBindOption{}
	currentIndex := 0
	for currentIndex < len(options) {
		found := false
		for _, parser := range getDgramBindOptions() {
			if size, err := parser.Parse(options, currentIndex); err == nil {
				result = append(result, parser)
				found = true
				currentIndex += size
			}
		}
		if !found {
			currentIndex++
		}
	}
	return result
}

func DgramBindOptionsString(options []DgramBindOption) string {
	var sb strings.Builder
	first := true
	for _, parser := range options {
		if parser.Valid() {
			if !first {
				sb.WriteString(" ")
			} else {
				first = false
			}
			sb.WriteString(parser.String())
		}
	}
	return sb.String()
}
