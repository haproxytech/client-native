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

// AddHeadersBin parses: { http-request | http-response | http-after-response }
// add-headers-bin <expr> [ prefix <str> ] [ { if | unless } <condition> ]
type AddHeadersBin struct {
	Expr     string
	Prefix   string
	Cond     string
	CondTest string
	Comment  string
}

func (f *AddHeadersBin) Parse(parts []string, parserType types.ParserType, comment string) error {
	if comment != "" {
		f.Comment = comment
	}
	if len(parts) < 3 {
		return stderrors.New("not enough params")
	}
	f.Expr = parts[2]
	command, condition := common.SplitRequest(parts[3:])
	if len(command) > 0 {
		if command[0] != "prefix" || len(command) < 2 {
			return stderrors.New("unknown params after expr")
		}
		f.Prefix = command[1]
	}
	if len(condition) > 1 {
		f.Cond = condition[0]
		f.CondTest = strings.Join(condition[1:], " ")
	}
	return nil
}

func (f *AddHeadersBin) String() string {
	var sb strings.Builder
	sb.WriteString("add-headers-bin ")
	sb.WriteString(f.Expr)
	if f.Prefix != "" {
		sb.WriteString(" prefix ")
		sb.WriteString(f.Prefix)
	}
	if f.Cond != "" {
		sb.WriteString(" ")
		sb.WriteString(f.Cond)
		sb.WriteString(" ")
		sb.WriteString(f.CondTest)
	}
	return sb.String()
}

func (f *AddHeadersBin) GetComment() string {
	return f.Comment
}
