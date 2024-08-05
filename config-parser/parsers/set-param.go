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

	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

type SetParam struct {
	data        []types.SetParam
	preComments []string
}

func (p *SetParam) Parse(line string, parts []string, comment string) (string, error) {
	switch len(parts) {
	case 1, 2:
		return "", &errors.ParseError{Parser: "SetParam", Line: line, Message: "Missing required values"}
	case 4:
		return "", &errors.ParseError{Parser: "SetParam", Line: line, Message: "Missing ACL criterion"}
	}

	data := types.SetParam{
		Name:   parts[1],
		Format: parts[2],
	}

	if len(parts) > 3 {
		switch parts[3] {
		case "if", "unless":
			data.Criterion = parts[3]
			data.Value = strings.Join(parts[4:], " ")
		default:
			return "", &errors.ParseError{Parser: "SetParam", Line: line, Message: "Unexpected ACL criterion"}
		}
	}

	p.data = append(p.data, data)

	return "", nil
}

func (p *SetParam) Result() ([]common.ReturnResultLine, error) {
	if len(p.data) == 0 {
		return nil, errors.ErrFetch
	}

	lines := make([]common.ReturnResultLine, 0, len(p.data))

	for _, data := range p.data {
		var sb strings.Builder

		sb.WriteString("set-param ")
		sb.WriteString(data.Name)
		sb.WriteString(" ")
		sb.WriteString(data.Format)

		if criterion := data.Criterion; len(criterion) > 0 {
			sb.WriteString(" ")
			sb.WriteString(criterion)
			sb.WriteString(" ")
			sb.WriteString(data.Value)
		}

		lines = append(lines, common.ReturnResultLine{
			Data:    sb.String(),
			Comment: data.Comment,
		})
	}

	return lines, nil
}
