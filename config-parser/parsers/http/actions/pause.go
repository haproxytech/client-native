/*
Copyright 2025 HAProxy Technologies

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
	"strings"

	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

type Pause struct {
	Pause    common.Expression
	Cond     string
	CondTest string
	Comment  string
}

func (f *Pause) Parse(parts []string, parserType types.ParserType, comment string) error {
	if comment != "" {
		f.Comment = comment
	}
	if len(parts) >= 3 {
		command, condition := common.SplitRequest(parts[2:])
		if len(command) == 0 {
			return errors.ErrInvalidData
		}
		pauseExpr := common.Expression{}
		if err := pauseExpr.Parse(command); err != nil {
			return stderrors.New("not enough params")
		}
		f.Pause = pauseExpr
		if len(condition) > 1 {
			f.Cond = condition[0]
			f.CondTest = strings.Join(condition[1:], " ")
		}
		return nil
	}
	return stderrors.New("not enough params")
}

func (f *Pause) String() string {
	var result strings.Builder
	result.WriteString("pause ")
	result.WriteString(f.Pause.String())
	if f.Cond != "" {
		result.WriteString(" ")
		result.WriteString(f.Cond)
		result.WriteString(" ")
		result.WriteString(f.CondTest)
	}
	return result.String()
}

func (f *Pause) GetComment() string {
	return f.Comment
}
