/*
Copyright 2023 HAProxy Technologies

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

type ScIncGpc struct {
	Idx      string
	ID       string
	Cond     string
	CondTest string
	Comment  string
}

func (f *ScIncGpc) Parse(parts []string, parserType types.ParserType, comment string) error {
	if comment != "" {
		f.Comment = comment
	}
	var data string
	var command []string
	var minLen, requiredLen int
	switch parserType {
	case types.HTTP, types.QUIC:
		data = parts[1]
		command = parts[2:]
		minLen = 2
		requiredLen = 4
	case types.TCP:
		data = parts[2]
		command = parts[3:]
		minLen = 3
		requiredLen = 5
	}
	idIdx := strings.TrimPrefix(data, "sc-inc-gpc(")
	idIdx = strings.TrimRight(idIdx, ")")
	idIdxValues := strings.SplitN(idIdx, ",", 2)
	f.Idx, f.ID = idIdxValues[0], idIdxValues[1]
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

func (f *ScIncGpc) String() string {
	var result strings.Builder
	result.WriteString("sc-inc-gpc(")
	result.WriteString(f.Idx)
	result.WriteString(",")
	result.WriteString(f.ID)
	result.WriteString(")")
	if f.Cond != "" {
		result.WriteString(" ")
		result.WriteString(f.Cond)
		result.WriteString(" ")
		result.WriteString(f.CondTest)
	}
	return result.String()
}

func (f *ScIncGpc) GetComment() string {
	return f.Comment
}
