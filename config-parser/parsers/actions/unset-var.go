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
	"strings"

	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

type UnsetVar struct {
	Scope    string
	Name     string
	Cond     string
	CondTest string
	Comment  string
}

func (f *UnsetVar) Parse(parts []string, parserType types.ParserType, comment string) error {
	if comment != "" {
		f.Comment = comment
	}
	if len(parts) < 2 {
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
	data = strings.TrimPrefix(data, "unset-var(")
	data = strings.TrimRight(data, ")")
	d := strings.SplitN(data, ".", 2)
	if len(d) < 2 || len(d[1]) == 0 {
		return errors.ErrInvalidData
	}
	f.Scope = d[0]
	f.Name = d[1]
	_, condition := common.SplitRequest(command)
	if len(condition) > 1 {
		f.Cond = condition[0]
		f.CondTest = strings.Join(condition[1:], " ")
	}
	return nil
}

func (f *UnsetVar) String() string {
	var result strings.Builder
	result.WriteString("unset-var(")
	result.WriteString(f.Scope)
	result.WriteString(".")
	result.WriteString(f.Name)
	result.WriteString(")")
	if f.Cond != "" {
		result.WriteString(" ")
		result.WriteString(f.Cond)
		result.WriteString(" ")
		result.WriteString(f.CondTest)
	}
	return result.String()
}

func (f *UnsetVar) GetComment() string {
	return f.Comment
}
