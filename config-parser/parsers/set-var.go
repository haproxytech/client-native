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

package parsers

import (
	stderrors "errors"
	"strings"

	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

type SetVar struct {
	data        []types.SetVar
	preComments []string // comments that appear before the actual line
}

func (p *SetVar) parse(line string, parts []string, comment string) (*types.SetVar, error) {
	if len(parts) < 3 {
		return nil, &errors.ParseError{Parser: "SetVar", Line: line}
	}
	expr := common.Expression{}
	if expr.Parse(parts[2:]) != nil {
		return nil, stderrors.New("not enough params")
	}
	data := &types.SetVar{
		Name:    parts[1],
		Expr:    expr,
		Comment: comment,
	}
	return data, nil
}

func (p *SetVar) Result() ([]common.ReturnResultLine, error) {
	if len(p.data) == 0 {
		return nil, errors.ErrFetch
	}
	result := make([]common.ReturnResultLine, len(p.data))
	for index, tg := range p.data {
		var sb strings.Builder
		sb.WriteString("set-var")
		sb.WriteString(" ")
		sb.WriteString(tg.Name)
		sb.WriteString(" ")
		sb.WriteString(tg.Expr.String())
		result[index] = common.ReturnResultLine{
			Data:    sb.String(),
			Comment: tg.Comment,
		}
	}
	return result, nil
}
