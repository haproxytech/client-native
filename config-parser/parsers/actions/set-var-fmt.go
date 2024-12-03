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

//nolint:dupl
package actions

import (
	stderrors "errors"
	"fmt"
	"strings"

	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

type SetVarFmt struct {
	VarScope string
	VarName  string
	Fmt      common.Expression
	Cond     string
	CondTest string
	Comment  string
}

func (f *SetVarFmt) Parse(parts []string, parserType types.ParserType, comment string) error {
	if comment != "" {
		f.Comment = comment
	}
	if len(parts) < 3 {
		return stderrors.New("not enough params")
	}
	var data string
	var command []string
	switch parserType {
	case types.HTTP, types.QUIC:
		data = parts[1]
		command = parts[2:]
	case types.TCP:
		data = parts[2]
		command = parts[3:]
	}
	data = strings.TrimPrefix(data, "set-var-fmt(")
	data = strings.TrimRight(data, ")")
	d := strings.SplitN(data, ".", 2)
	f.VarScope = d[0]
	f.VarName = d[1]
	command, condition := common.SplitRequest(command)
	if len(command) > 0 {
		expr := common.Expression{}
		err := expr.Parse(command)
		if err != nil {
			return stderrors.New("not enough params")
		}
		f.Fmt = expr
	}
	if len(condition) > 1 {
		f.Cond = condition[0]
		f.CondTest = strings.Join(condition[1:], " ")
	}
	return nil
}

func (f *SetVarFmt) String() string {
	condition := ""
	if f.Cond != "" {
		condition = fmt.Sprintf(" %s %s", f.Cond, f.CondTest)
	}
	return fmt.Sprintf("set-var-fmt(%s.%s) %s%s", f.VarScope, f.VarName, f.Fmt.String(), condition)
}

func (f *SetVarFmt) GetComment() string {
	return f.Comment
}
