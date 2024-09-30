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

//nolint:dupl
package actions

import (
	"fmt"
	"strings"

	"github.com/haproxytech/client-native/v5/config-parser/common"
	"github.com/haproxytech/client-native/v5/config-parser/errors"
	"github.com/haproxytech/client-native/v5/config-parser/types"
)

type SetTos struct {
	Value    string
	Cond     string
	CondTest string
	Comment  string
}

func (f *SetTos) Parse(parts []string, parserType types.ParserType, comment string) error {
	if comment != "" {
		f.Comment = comment
	}

	var command []string
	switch parserType {
	case types.HTTP:
		command = parts[2:]
	case types.TCP:
		command = parts[3:]
	}
	if len(parts) >= 1 {
		var condition []string
		command, condition = common.SplitRequest(command)
		if len(command) == 0 {
			return errors.ErrInvalidData
		}
		f.Value = command[0]
		if len(condition) > 1 {
			f.Cond = condition[0]
			f.CondTest = strings.Join(condition[1:], " ")
		}
		return nil
	}
	return fmt.Errorf("not enough params")
}

func (f *SetTos) String() string {
	var result strings.Builder
	result.WriteString("set-tos ")
	result.WriteString(f.Value)
	if f.Cond != "" {
		result.WriteString(" ")
		result.WriteString(f.Cond)
		result.WriteString(" ")
		result.WriteString(f.CondTest)
	}
	return result.String()
}

func (f *SetTos) GetComment() string {
	return f.Comment
}
