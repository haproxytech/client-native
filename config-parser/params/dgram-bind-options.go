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

var dgramBindOptionFactoryMethods = map[string]func() DgramBindOption{ //nolint:gochecknoglobals
	"transparent": func() DgramBindOption { return &BindOptionWord{Name: "transparent"} },
	"interface":   func() DgramBindOption { return &BindOptionValue{Name: "interface"} },
	"namespace":   func() DgramBindOption { return &BindOptionValue{Name: "namespace"} },
	"name":        func() DgramBindOption { return &BindOptionValue{Name: "name"} },
}

func getDgramBindOption(option string) DgramBindOption {
	if factoryMethod, found := dgramBindOptionFactoryMethods[option]; found {
		return factoryMethod()
	}
	return nil
}

// Parse ...
func ParseDgramBindOptions(options []string) ([]DgramBindOption, error) {
	result := []DgramBindOption{}
	currentIndex := 0
	for currentIndex < len(options) {
		dgramBindOption := getDgramBindOption(options[currentIndex])
		if dgramBindOption == nil {
			currentIndex++
			continue
		}
		size, err := dgramBindOption.Parse(options, currentIndex)
		if err != nil {
			return nil, err
		}
		result = append(result, dgramBindOption)
		currentIndex += size
	}
	return result, nil
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
