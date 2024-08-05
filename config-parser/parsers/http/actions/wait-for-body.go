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

	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

type WaitForBody struct {
	Time     string
	AtLeast  string
	Cond     string
	CondTest string
	Comment  string
}

func (f *WaitForBody) Parse(parts []string, parserType types.ParserType, comment string) error {
	if comment != "" {
		f.Comment = comment
	}
	if len(parts) >= 4 {
		if parts[2] != "time" {
			return errors.ErrInvalidData
		}
		command, condition := common.SplitRequest(parts[3:])
		if len(command) < 1 || len(command) == 2 {
			return fmt.Errorf("not enough params")
		}
		f.Time = command[0]
		if len(command) > 2 {
			if command[1] != "at-least" {
				return fmt.Errorf("unknown param: %s", command[1])
			}
			f.AtLeast = command[2]
		}

		if len(condition) > 1 {
			f.Cond = condition[0]
			f.CondTest = strings.Join(condition[1:], " ")
		}
		return nil
	}
	return fmt.Errorf("not enough params")
}

func (f *WaitForBody) String() string {
	var result strings.Builder
	result.WriteString("wait-for-body time ")
	result.WriteString(f.Time)
	if f.AtLeast != "" {
		result.WriteString(" at-least ")
		result.WriteString(f.AtLeast)
	}
	if f.Cond != "" {
		result.WriteString(" ")
		result.WriteString(f.Cond)
		result.WriteString(" ")
		result.WriteString(f.CondTest)
	}
	return result.String()
}

func (f *WaitForBody) GetComment() string {
	return f.Comment
}
