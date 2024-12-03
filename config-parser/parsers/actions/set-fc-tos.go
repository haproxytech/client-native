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

package actions

import (
	stderrors "errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

type SetFcTos struct {
	Expr     common.Expression
	Cond     string
	CondTest string
	Comment  string
}

func (f *SetFcTos) Parse(parts []string, parserType types.ParserType, comment string) error {
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
		if len(command) == 1 && !validateUnsignedNumber(command[0], math.MaxUint8) {
			return fmt.Errorf("number '%s' is not a valid TOS value", command[0])
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

func (f *SetFcTos) String() string {
	return common.SmartJoin("set-fc-tos", f.Expr.String(), f.Cond, f.CondTest)
}

func (f *SetFcTos) GetComment() string {
	return f.Comment
}

// Test if the given string is an unsigned integer between zero and "max".
// The number can be in decimal or hexadecimal (0x).
// If the parsing failed, assume the string was an Expr and return true.
func validateUnsignedNumber(text string, maximum int64) bool {
	var n int64
	var err error
	if strings.HasPrefix(text, "0x") {
		n, err = strconv.ParseInt(text, 16, 64)
	} else {
		n, err = strconv.ParseInt(text, 10, 64)
	}
	if err != nil {
		// Assume it was an expression, not a number.
		return true
	}
	return n >= 0 && n <= maximum
}
