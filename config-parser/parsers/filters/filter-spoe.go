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

package filters

import (
	"strings"

	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/errors"
)

type Spoe struct { // filter spoe [engine <name>] config <file>
	Engine  string
	Config  string
	Comment string
}

func (f *Spoe) Parse(parts []string, comment string) error {
	if comment != "" {
		f.Comment = comment
	}
	index := 2
	for index < len(parts) {
		switch parts[index] {
		case "engine":
			index++
			if index == len(parts) {
				return errors.ErrInvalidData
			}
			f.Engine = parts[index]
		case "config":
			index++
			if index == len(parts) {
				return errors.ErrInvalidData
			}
			f.Config = parts[index]
		default:
			return errors.ErrInvalidData
		}
		index++
	}
	if f.Config == "" {
		return errors.ErrInvalidData
	}
	return nil
}

func (f *Spoe) Result() common.ReturnResultLine {
	var result strings.Builder
	result.WriteString("filter spoe")
	if f.Engine != "" {
		result.WriteString(" engine ")
		result.WriteString(f.Engine)
	}
	result.WriteString(" config ")
	result.WriteString(f.Config)
	return common.ReturnResultLine{
		Data:    result.String(),
		Comment: f.Comment,
	}
}
