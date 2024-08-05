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

type HashType struct {
	data        *types.HashType
	preComments []string // comments that appear before the actual line
}

func (p *HashType) Parse(line string, parts []string, comment string) (string, error) {
	if parts[0] == "hash-type" {
		if len(parts) < 2 {
			return "", &errors.ParseError{Parser: "HashType", Line: line, Message: "Parse error"}
		}
		data := types.HashType{
			Comment: comment,
		}
		index := 1
		for ; index < len(parts); index++ {
			switch parts[index] {
			case "map-based", "consistent":
				data.Method = parts[index]
			case "sdbm", "djb2", "wt6", "crc32", "none":
				data.Function = parts[index]
			case "avalanche":
				data.Modifier = parts[index]
			default:
				return "", &errors.ParseError{Parser: "HashType", Line: line, Message: "Parse error"}
			}
		}
		p.data = &data
		return "", nil
	}
	return "", &errors.ParseError{Parser: "HashType", Line: line}
}

func (p *HashType) Result() ([]common.ReturnResultLine, error) {
	if p.data == nil {
		return nil, errors.ErrFetch
	}
	var sb strings.Builder
	sb.WriteString("hash-type")
	if p.data.Method != "" {
		sb.WriteString(" ")
		sb.WriteString(p.data.Method)
	}
	if p.data.Function != "" {
		sb.WriteString(" ")
		sb.WriteString(p.data.Function)
	}
	if p.data.Modifier != "" {
		sb.WriteString(" ")
		sb.WriteString(p.data.Modifier)
	}
	return []common.ReturnResultLine{
		{
			Data:    sb.String(),
			Comment: p.data.Comment,
		},
	}, nil
}
