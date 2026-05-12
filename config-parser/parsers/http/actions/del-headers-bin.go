/*
Copyright 2026 HAProxy Technologies

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
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

// DelHeadersBin parses: { http-request | http-response | http-after-response }
// del-headers-bin <expr> [ -m <meth> ] [ { if | unless } <condition> ]
// where <meth> is one of str|beg|end|sub (regex is not supported).
type DelHeadersBin struct {
	Expr     string
	Method   string
	Cond     string
	CondTest string
	Comment  string
}

func (f *DelHeadersBin) Parse(parts []string, parserType types.ParserType, comment string) error {
	if comment != "" {
		f.Comment = comment
	}
	if len(parts) < 3 {
		return stderrors.New("not enough params")
	}
	f.Expr = parts[2]
	command, condition := common.SplitRequest(parts[3:])
	if len(command) > 0 {
		if command[0] != "-m" || len(command) < 2 {
			return stderrors.New("unknown params after expr")
		}
		f.Method = command[1]
	}
	if len(condition) > 1 {
		f.Cond = condition[0]
		f.CondTest = strings.Join(condition[1:], " ")
	}
	return nil
}

func (f *DelHeadersBin) String() string {
	var sb strings.Builder
	sb.WriteString("del-headers-bin ")
	sb.WriteString(f.Expr)
	if f.Method != "" {
		sb.WriteString(" -m ")
		sb.WriteString(f.Method)
	}
	if f.Cond != "" {
		sb.WriteString(" ")
		sb.WriteString(f.Cond)
		sb.WriteString(" ")
		sb.WriteString(f.CondTest)
	}
	return sb.String()
}

func (f *DelHeadersBin) GetComment() string {
	return f.Comment
}
