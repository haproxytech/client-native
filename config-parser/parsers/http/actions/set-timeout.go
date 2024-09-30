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
	"github.com/haproxytech/client-native/v5/config-parser/errors"
	"github.com/haproxytech/client-native/v5/config-parser/types"
)

type SetTimeout struct {
	Type     string
	Timeout  string
	Cond     string
	CondTest string
	Comment  string
}

func (f *SetTimeout) Parse(parts []string, parserType types.ParserType, comment string) error {
	if comment != "" {
		f.Comment = comment
	}
	if len(parts) >= 4 {
		command, condition := common.SplitRequest(parts[2:])
		if len(command) < 2 {
			return errors.ErrInvalidData
		}
		if command[0] != "server" && command[0] != "tunnel" && command[0] != "client" {
			return fmt.Errorf("unknown timeout type: %s", command[0])
		}
		f.Type = command[0]
		f.Timeout = command[1]
		if len(condition) > 1 {
			f.Cond = condition[0]
			f.CondTest = strings.Join(condition[1:], " ")
		}
		return nil

	}
	return fmt.Errorf("not enough params")
}

func (f *SetTimeout) String() string {
	var result strings.Builder
	result.WriteString("set-timeout ")
	result.WriteString(f.Type)
	result.WriteString(" ")
	result.WriteString(f.Timeout)
	if f.Cond != "" {
		result.WriteString(" ")
		result.WriteString(f.Cond)
		result.WriteString(" ")
		result.WriteString(f.CondTest)
	}
	return result.String()
}

func (f *SetTimeout) GetComment() string {
	return f.Comment
}
