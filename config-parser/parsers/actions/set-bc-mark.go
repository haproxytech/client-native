/*
Copyright 2024 HAProxy Technologies

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
	"math"
	"strings"

	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

type SetBcMark struct {
	Expr     common.Expression
	Cond     string
	CondTest string
	Comment  string
}

func (f *SetBcMark) Parse(parts []string, parserType types.ParserType, comment string) error {
	if comment != "" {
		f.Comment = comment
	}
	if len(parts) < 3 {
		return stderrors.New("not enough params")
	}
	var command []string
	switch parserType {
	case types.HTTP, types.QUIC:
		command = parts[2:]
	case types.TCP:
		command = parts[3:]
	}
	command, condition := common.SplitRequest(command)

	if len(command) > 0 {
		if len(command) == 1 && !validateUnsignedNumber(command[0], math.MaxUint32) {
			return fmt.Errorf("number '%s' is not a valid mark value", command[0])
		}
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

func (f *SetBcMark) String() string {
	return common.SmartJoin("set-bc-mark", f.Expr.String(), f.Cond, f.CondTest)
}

func (f *SetBcMark) GetComment() string {
	return f.Comment
}
