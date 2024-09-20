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
	stderrors "errors"
	"fmt"
	"strings"

	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

type SetPath struct {
	Fmt      string
	Cond     string
	CondTest string
	Comment  string
}

// Parse parses http-request set-path <fmt> [ { if | unless } <condition> ]
func (f *SetPath) Parse(parts []string, parserType types.ParserType, comment string) error {
	if comment != "" {
		f.Comment = comment
	}

	if len(parts) >= 3 {
		command, condition := common.SplitRequest(parts[2:])
		f.Fmt = strings.Join(command, " ")
		if len(condition) > 1 {
			f.Cond = condition[0]
			f.CondTest = strings.Join(condition[1:], " ")
		}

		return nil
	}

	return stderrors.New("not enough params")
}

func (f *SetPath) String() string {
	condition := ""

	if f.Cond != "" {
		condition = fmt.Sprintf(" %s %s", f.Cond, f.CondTest)
	}

	return fmt.Sprintf("set-path %s%s", f.Fmt, condition)
}

func (f *SetPath) GetComment() string {
	return f.Comment
}
