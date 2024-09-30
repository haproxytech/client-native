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
	"strings"

	"github.com/haproxytech/client-native/v5/config-parser/common"
	"github.com/haproxytech/client-native/v5/config-parser/errors"
	"github.com/haproxytech/client-native/v5/config-parser/types"
)

type SetVarFmt struct {
	data        []types.SetVarFmt
	preComments []string // comments that appear before the actual line
}

func (t *SetVarFmt) parse(line string, parts []string, comment string) (*types.SetVarFmt, error) {
	if len(parts) < 3 {
		return nil, &errors.ParseError{Parser: "SetVarFmt", Line: line}
	}
	data := &types.SetVarFmt{
		Name:    parts[1],
		Format:  parts[2],
		Comment: comment,
	}
	return data, nil
}

func (t *SetVarFmt) Result() ([]common.ReturnResultLine, error) {
	if len(t.data) == 0 {
		return nil, errors.ErrFetch
	}
	result := make([]common.ReturnResultLine, len(t.data))
	for index, f := range t.data {
		var sb strings.Builder
		sb.WriteString("set-var-fmt")
		sb.WriteString(" ")
		sb.WriteString(f.Name)
		sb.WriteString(" ")
		sb.WriteString(f.Format)
		result[index] = common.ReturnResultLine{
			Data:    sb.String(),
			Comment: f.Comment,
		}
	}
	return result, nil
}
