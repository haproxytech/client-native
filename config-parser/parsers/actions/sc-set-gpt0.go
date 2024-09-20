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
	"strconv"
	"strings"

	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

type ScSetGpt0 struct {
	ID       string
	Int      *int64
	Expr     common.Expression
	Cond     string
	CondTest string
	Comment  string
}

func (f *ScSetGpt0) Parse(parts []string, parserType types.ParserType, comment string) error {
	if comment != "" {
		f.Comment = comment
	}
	if len(parts) < 3 {
		return stderrors.New("not enough params")
	}
	var data string
	var command []string
	switch parserType {
	case types.HTTP:
		data = parts[1]
		command = parts[2:]
	case types.TCP:
		data = parts[2]
		command = parts[3:]
	}
	f.ID = strings.TrimPrefix(data, "sc-set-gpt0(")
	f.ID = strings.TrimRight(f.ID, ")")
	command, condition := common.SplitRequest(command)
	if len(command) < 1 {
		return errors.ErrInvalidData
	}
	i, err := strconv.ParseInt(command[0], 10, 64)
	if err == nil {
		f.Int = &i
	} else {
		expr := common.Expression{}
		err := expr.Parse(command)
		if err != nil {
			return stderrors.New("not enough params")
		}
		f.Expr = expr
	}
	if len(condition) > 1 {
		f.Cond = condition[0]
		f.CondTest = strings.Join(condition[1:], " ")
	}
	return nil
}

func (f *ScSetGpt0) String() string {
	var result strings.Builder
	result.WriteString("sc-set-gpt0(")
	result.WriteString(f.ID)
	result.WriteString(") ")
	if f.Int != nil {
		result.WriteString(strconv.FormatInt(*f.Int, 10))
	} else {
		result.WriteString(f.Expr.String())
	}
	if f.Cond != "" {
		result.WriteString(" ")
		result.WriteString(f.Cond)
		result.WriteString(" ")
		result.WriteString(f.CondTest)
	}
	return result.String()
}

func (f *ScSetGpt0) GetComment() string {
	return f.Comment
}
