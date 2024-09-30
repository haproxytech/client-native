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
	"fmt"
	"strconv"
	"strings"

	"github.com/haproxytech/client-native/v5/config-parser/common"
	"github.com/haproxytech/client-native/v5/config-parser/types"
)

type Capture struct {
	Expr     common.Expression
	Len      int64
	Cond     string
	CondTest string
	Comment  string
}

func (f *Capture) Parse(parts []string, parserType types.ParserType, comment string) error {
	if f.Comment != "" {
		f.Comment = comment
	}
	if len(parts) < 4 {
		return fmt.Errorf("not enough params")
	}
	expr := common.Expression{}

	err := expr.Parse([]string{parts[3]})
	if err != nil {
		return fmt.Errorf("invalid expression")
	}

	f.Expr = expr

	if ln, err := strconv.ParseInt(parts[5], 10, 64); err == nil {
		f.Len = ln
	} else {
		return fmt.Errorf("invalid value for len")
	}
	_, condition := common.SplitRequest(parts[5:])
	if len(condition) > 1 {
		f.Cond = condition[0]
		f.CondTest = strings.Join(condition[1:], " ")
	}
	return nil
}

func (f *Capture) String() string {
	var result strings.Builder
	result.WriteString("capture")
	result.WriteString(" ")
	result.WriteString(f.Expr.String())
	result.WriteString(" ")
	result.WriteString("len ")
	result.WriteString(strconv.FormatInt(f.Len, 10))
	if f.Cond != "" {
		result.WriteString(" ")
		result.WriteString(f.Cond)
		result.WriteString(" ")
		result.WriteString(f.CondTest)
	}
	return result.String()
}

func (f *Capture) GetComment() string {
	return f.Comment
}
