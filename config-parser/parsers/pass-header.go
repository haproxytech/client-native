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
	"unicode"

	"github.com/haproxytech/client-native/v5/config-parser/common"
	"github.com/haproxytech/client-native/v5/config-parser/errors"
	"github.com/haproxytech/client-native/v5/config-parser/types"
)

type PassHeader struct {
	data        []types.PassHeader
	preComments []string
}

func (p *PassHeader) Parse(line string, parts []string, comment string) (string, error) {
	switch len(parts) {
	case 1:
		return "", &errors.ParseError{Parser: "PassHeader", Line: line, Message: "Missing header name"}
	case 3:
		return "", &errors.ParseError{Parser: "PassHeader", Line: line, Message: "Missing ACL condition"}
	}

	data := types.PassHeader{
		Name:    parts[1],
		Comment: comment,
	}

	switch {
	case !unicode.IsLetter([]rune(data.Name[0:1])[0]):
		return "", &errors.ParseError{Parser: "PassHeader", Line: line, Message: "Header name cannot start with non letter"}
	case !unicode.IsLetter([]rune(data.Name[len(data.Name)-1 : len(data.Name)])[0]):
		return "", &errors.ParseError{Parser: "PassHeader", Line: line, Message: "Header name cannot start with non letter"}
	}

	if len(parts) > 2 {
		switch parts[2] {
		case "if", "unless":
			data.Criterion = parts[2]
			data.Value = strings.Join(parts[3:], " ")
		default:
			return "", &errors.ParseError{Parser: "PassHeader", Line: line, Message: "Unexpected ACL criterion"}
		}
	}

	p.data = append(p.data, data)

	return "", nil
}

func (p *PassHeader) Result() ([]common.ReturnResultLine, error) {
	if len(p.data) == 0 {
		return nil, errors.ErrFetch
	}

	lines := make([]common.ReturnResultLine, 0, len(p.data))

	for _, data := range p.data {
		var sb strings.Builder

		sb.WriteString("pass-header ")
		sb.WriteString(data.Name)

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
