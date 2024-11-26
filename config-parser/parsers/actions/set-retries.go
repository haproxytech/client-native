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
	stderrors "errors"
	"strings"

	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

type SetRetries struct {
	Expr     common.Expression
	Cond     string
	CondTest string
	Comment  string
}

func (f *SetRetries) Parse(parts []string, parserType types.ParserType, comment string) error {
	if comment != "" {
		f.Comment = comment
	}
	if len(parts) < 3 {
		return stderrors.New("not enough params")
	}
	var command []string
	switch parserType {
	case types.HTTP:
		command = parts[2:]
	case types.TCP:
		command = parts[3:]
	}
	command, condition := common.SplitRequest(command)
	if len(command) == 0 {
		return errors.ErrInvalidData
	}
	expr := common.Expression{}
	if expr.Parse(command) != nil {
		return stderrors.New("not enough params")
	}
	f.Expr = expr
	if len(condition) > 1 {
		f.Cond = condition[0]
		f.CondTest = strings.Join(condition[1:], " ")
	}
	return nil
}

func (f *SetRetries) String() string {
	var result strings.Builder
	result.WriteString("set-retries ")
	result.WriteString(f.Expr.String())
	if f.Cond != "" {
		result.WriteString(" ")
		result.WriteString(f.Cond)
		result.WriteString(" ")
		result.WriteString(f.CondTest)
	}
	return result.String()
}

func (f *SetRetries) GetComment() string {
	return f.Comment
}
