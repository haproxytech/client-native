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

	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

type DelHeader struct {
	Name     string
	Method   string
	Cond     string
	CondTest string
	Comment  string
}

// Parse parses { http-request | http-response } del-header <name> [ { if | unless } <condition> ]
func (f *DelHeader) Parse(parts []string, parserType types.ParserType, comment string) error {
	if comment != "" {
		f.Comment = comment
	}

	if len(parts) >= 3 {
		command, condition := common.SplitRequest(parts[3:])
		f.Name = parts[2]
		if len(command) > 1 && command[0] == "-m" {
			f.Method = command[1]
		} else if len(command) > 0 {
			return fmt.Errorf("unknown params after name")
		}
		if len(condition) > 1 {
			f.Cond = condition[0]
			f.CondTest = strings.Join(condition[1:], " ")
		}
		return nil
	}
	return fmt.Errorf("not enough params")
}

func (f *DelHeader) String() string {
	sb := &strings.Builder{}

	sb.WriteString("del-header ")
	sb.WriteString(f.Name)
	if f.Method != "" {
		sb.WriteString(fmt.Sprintf(" -m %s", f.Method))
	}
	if f.Cond != "" {
		sb.WriteString(fmt.Sprintf(" %s %s", f.Cond, f.CondTest))
	}
	return sb.String()
}

func (f *DelHeader) GetComment() string {
	return f.Comment
}
