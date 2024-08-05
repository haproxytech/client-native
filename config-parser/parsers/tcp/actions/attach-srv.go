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
	"fmt"
	"strings"

	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

type AttachSrv struct {
	Server   string
	Name     common.Expression
	Cond     string
	CondTest string
	Comment  string
}

// tcp-request session attach-srv <srv> [name <expr>] [{ if | unless } <condition>]

func (f *AttachSrv) Parse(parts []string, parserType types.ParserType, comment string) error {
	if f.Comment != "" {
		f.Comment = comment
	}

	n := len(parts)
	if n < 4 {
		return fmt.Errorf("not enough params")
	}

	f.Server = parts[3]

	if n == 4 {
		// Nothing more to parse.
		return nil
	}

	i := 4

	if parts[i] == "name" {
		if n < 6 {
			return fmt.Errorf("not enough params")
		}
		expr := common.Expression{}
		err := expr.Parse([]string{parts[i+1]})
		if err != nil {
			return fmt.Errorf("invalid expression after '%s': %w", parts[i], err)
		}
		f.Name = expr
		i += 2
		if i == n {
			return nil
		}
	}

	_, condition := common.SplitRequest(parts[i:])
	if len(condition) > 1 {
		f.Cond = condition[0]
		f.CondTest = strings.Join(condition[1:], " ")
	} else {
		return fmt.Errorf("invalid condition after '%s'", parts[i])
	}

	return nil
}

func (f *AttachSrv) String() string {
	var result strings.Builder
	result.WriteString("attach-srv ")
	result.WriteString(f.Server)
	if name := f.Name.String(); name != "" {
		result.WriteString(" name ")
		result.WriteString(name)
	}
	if f.Cond != "" {
		result.WriteString(" ")
		result.WriteString(f.Cond)
		result.WriteString(" ")
		result.WriteString(f.CondTest)
	}
	return result.String()
}

func (f *AttachSrv) GetComment() string {
	return f.Comment
}
