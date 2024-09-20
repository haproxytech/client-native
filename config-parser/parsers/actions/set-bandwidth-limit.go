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

type SetBandwidthLimit struct {
	Limit    common.Expression
	Period   common.Expression
	Name     string
	Cond     string
	CondTest string
	Comment  string
}

func (f *SetBandwidthLimit) Parse(parts []string, parserType types.ParserType, comment string) error {
	if comment != "" {
		f.Comment = comment
	}

	var command []string
	switch parserType {
	case types.HTTP:
		if len(parts) < 3 {
			return stderrors.New("not enough params")
		}
		command = parts[2:]
	case types.TCP:
		if len(parts) < 4 {
			return stderrors.New("not enough params")
		}
		command = parts[3:]
	}
	command, condition := common.SplitRequest(command)
	f.Name = command[0]

	for i := 1; i < len(command); i++ {
		var expr []string
		if len(command) < i+2 {
			return stderrors.New("not enough params")
		}
		el := command[i]
		switch el {
		case "limit":
			expr, i = f.parseExpr(command, i+1, "period")
			if len(expr) == 0 {
				return stderrors.New("not enough params")
			}
			f.Limit.Expr = expr
		case "period":
			expr, i = f.parseExpr(command, i+1, "limit")
			if len(expr) == 0 {
				return stderrors.New("not enough params")
			}
			f.Period.Expr = expr
		default:
			return fmt.Errorf("invalid param %s", el)
		}
	}

	if len(condition) > 1 {
		f.Cond = condition[0]
		f.CondTest = strings.Join(condition[1:], " ")
	}
	return nil
}

func (f *SetBandwidthLimit) parseExpr(parts []string, start int, end string) ([]string, int) {
	for i := start; i < len(parts); i++ {
		if parts[i] == end {
			return parts[start:i], i - 1
		}
	}
	return parts[start:], len(parts) - 1
}

func (f *SetBandwidthLimit) String() string {
	var result strings.Builder
	result.WriteString("set-bandwidth-limit ")
	result.WriteString(f.Name)
	if len(f.Limit.Expr) != 0 {
		result.WriteString(" limit ")
		result.WriteString(f.Limit.String())
	}
	if len(f.Period.Expr) != 0 {
		result.WriteString(" period ")
		result.WriteString(f.Period.String())
	}
	if f.Cond != "" {
		result.WriteString(" ")
		result.WriteString(f.Cond)
		result.WriteString(" ")
		result.WriteString(f.CondTest)
	}
	return result.String()
}

func (f *SetBandwidthLimit) GetComment() string {
	return f.Comment
}
