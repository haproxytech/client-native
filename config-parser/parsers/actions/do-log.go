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
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

type DoLog struct {
	Profile  string
	Cond     string
	CondTest string
	Comment  string
}

func (f *DoLog) Parse(parts []string, parserType types.ParserType, comment string) error {
	if comment != "" {
		f.Comment = comment
	}
	var command []string
	var minLen int
	switch parserType {
	case types.HTTP, types.QUIC:
		command = parts[1:]
		minLen = 2
	case types.TCP:
		command = parts[2:]
		minLen = 3
	}
	if len(parts) == minLen {
		return nil
	}
	// command[0] is "do-log". The remaining tokens are an optional
	// "profile <name>" plus an optional "if|unless <cond>".
	rest := command[1:]
	if len(rest) >= 2 && rest[0] == "profile" {
		f.Profile = rest[1]
		rest = rest[2:]
	}
	if len(rest) == 0 {
		return nil
	}
	if len(rest) < 2 {
		return stderrors.New("not enough params")
	}
	_, condition := common.SplitRequest(rest)
	if len(condition) > 1 {
		f.Cond = condition[0]
		f.CondTest = strings.Join(condition[1:], " ")
	}
	return nil
}

func (f *DoLog) String() string {
	var b strings.Builder
	b.WriteString("do-log")
	if f.Profile != "" {
		b.WriteString(" profile ")
		b.WriteString(f.Profile)
	}
	if f.Cond != "" {
		b.WriteString(" ")
		b.WriteString(f.Cond)
		b.WriteString(" ")
		b.WriteString(f.CondTest)
	}
	return b.String()
}

func (f *DoLog) GetComment() string {
	return f.Comment
}
