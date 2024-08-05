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
	"fmt"
	"strings"

	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

type SetVarFmtCheck struct {
	VarScope string
	VarName  string
	Format   common.Expression
	Comment  string
}

func (f *SetVarFmtCheck) Parse(parts []string, parserType types.ParserType, comment string) error {
	if comment != "" {
		f.Comment = comment
	}
	if len(parts) < 3 {
		return fmt.Errorf("not enough params")
	}
	var data string
	var command []string

	data = parts[1]
	command = parts[2:]

	data = strings.TrimPrefix(data, "set-var-fmt(")
	data = strings.TrimRight(data, ")")
	d := strings.SplitN(data, ".", 2)
	f.VarScope = d[0]
	f.VarName = d[1]
	command, condition := common.SplitRequest(command)
	if len(command) > 0 {
		format := common.Expression{}
		err := format.Parse(command)
		if err != nil {
			return fmt.Errorf("not enough params")
		}
		f.Format = format
	} else {
		return fmt.Errorf("not enough params")
	}
	if len(condition) > 1 {
		return errors.ErrInvalidData
	}

	return nil
}

func (f *SetVarFmtCheck) String() string {
	return fmt.Sprintf("set-var-fmt(%s.%s) %s", f.VarScope, f.VarName, f.Format.String())
}

func (f *SetVarFmtCheck) GetComment() string {
	return f.Comment
}
