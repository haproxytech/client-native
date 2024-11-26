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
	"fmt"
	"strings"

	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

type DoLog struct {
	Cond     string
	CondTest string
	Comment  string
}

func (f *DoLog) Parse(parts []string, parserType types.ParserType, comment string) error {
	if comment != "" {
		f.Comment = comment
	}
	var command []string
	var minLen, requiredLen int
	switch parserType {
	case types.HTTP:
		command = parts[1:]
		minLen = 2
		requiredLen = 4
	case types.TCP:
		command = parts[2:]
		minLen = 3
		requiredLen = 5
	}
	if len(parts) == minLen {
		return nil
	}
	if len(parts) < requiredLen {
		return stderrors.New("not enough params")
	}
	_, condition := common.SplitRequest(command)
	if len(condition) > 1 {
		f.Cond = condition[0]
		f.CondTest = strings.Join(condition[1:], " ")
	}
	return nil
}

func (f *DoLog) String() string {
	if f.Cond != "" {
		return fmt.Sprintf("do-log %s %s", f.Cond, f.CondTest)
	}
	return "do-log"
}

func (f *DoLog) GetComment() string {
	return f.Comment
}
