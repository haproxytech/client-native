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

package actions

import (
	"fmt"
	"strings"

	"github.com/haproxytech/client-native/v5/config-parser/common"
	"github.com/haproxytech/client-native/v5/config-parser/types"
)

type SwitchMode struct {
	Proto    string
	Cond     string
	CondTest string
	Comment  string
}

func (f *SwitchMode) Parse(parts []string, parserType types.ParserType, comment string) error {
	if f.Comment != "" {
		f.Comment = comment
	}
	command, condition := common.SplitRequest(parts)
	switch len(command) {
	case 4:
	case 6:
		if command[4] != "proto" {
			return fmt.Errorf("invalid param: %s", command[4])
		}
		f.Proto = command[5]
	default:
		return fmt.Errorf("not enough params")
	}
	if command[3] != "http" {
		return fmt.Errorf("invalid param %s", command[3])
	}
	if len(condition) > 1 {
		f.Cond = condition[0]
		f.CondTest = strings.Join(condition[1:], " ")
	}
	return nil
}

func (f *SwitchMode) String() string {
	var result strings.Builder
	result.WriteString("switch-mode http")
	if f.Proto != "" {
		result.WriteString(" proto ")
		result.WriteString(f.Proto)
	}
	if f.Cond != "" {
		result.WriteString(" ")
		result.WriteString(f.Cond)
		result.WriteString(" ")
		result.WriteString(f.CondTest)
	}
	return result.String()
}

func (f *SwitchMode) GetComment() string {
	return f.Comment
}
