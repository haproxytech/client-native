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
	"strconv"
	"strings"

	"github.com/haproxytech/client-native/v5/config-parser/common"
	"github.com/haproxytech/client-native/v5/config-parser/errors"
	"github.com/haproxytech/client-native/v5/config-parser/types"
)

type SilentDrop struct {
	RstTTL   int64
	Cond     string
	CondTest string
	Comment  string
}

func (f *SilentDrop) Parse(parts []string, parserType types.ParserType, comment string) error {
	if comment != "" {
		f.Comment = comment
	}
	if len(parts) < 2 {
		return stderrors.New("not enough params")
	}
	var command []string
	switch parserType {
	case types.HTTP:
		command = parts[2:]
	case types.TCP:
		command = parts[3:]
	}
	command, condition := common.SplitRequest(command)
	if len(condition) > 1 {
		f.Cond = condition[0]
		f.CondTest = strings.Join(condition[1:], " ")
	}
	if len(command) > 0 && command[0] == "rst-ttl" {
		if len(command) <= 1 {
			return stderrors.New("missing rst-ttl value")
		}
		rstTTL, err := strconv.ParseInt(command[1], 10, 64)
		if err != nil {
			return &errors.ParseError{Parser: "SilentDrop", Message: err.Error()}
		}
		f.RstTTL = rstTTL
	}
	return nil
}

func (f *SilentDrop) String() string {
	var result strings.Builder
	result.WriteString("silent-drop")
	if f.RstTTL > 0 {
		result.WriteString(" rst-ttl ")
		result.WriteString(strconv.FormatInt(f.RstTTL, 10))
	}
	if f.Cond != "" {
		result.WriteString(" ")
		result.WriteString(f.Cond)
		result.WriteString(" ")
		result.WriteString(f.CondTest)
	}
	return result.String()
}

func (f *SilentDrop) GetComment() string {
	return f.Comment
}
