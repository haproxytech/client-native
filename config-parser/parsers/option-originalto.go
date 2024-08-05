/*
Copyright 2022 HAProxy Technologies

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

type OptionOriginalTo struct {
	data        *types.OptionOriginalTo
	preComments []string // comments that appear before the actual line
}

/*
option originalto [ except <network> ] [ header <name> ]
*/
func (s *OptionOriginalTo) Parse(line string, parts []string, comment string) (string, error) {
	if len(parts) > 1 && parts[0] == "option" && parts[1] == "originalto" {
		data := &types.OptionOriginalTo{
			Comment: comment,
		}
		index := 2
		for index < len(parts) {
			switch parts[index] {
			case "except":
				index++
				if index == len(parts) {
					return "", errors.ErrInvalidData
				}
				data.Except = parts[index]
			case "header":
				index++
				if index == len(parts) {
					return "", errors.ErrInvalidData
				}
				data.Header = parts[index]
			default:
				return "", errors.ErrInvalidData
			}
			index++
		}
		s.data = data
		return "", nil
	}
	return "", &errors.ParseError{Parser: "option originalto", Line: line}
}

func (s *OptionOriginalTo) Result() ([]common.ReturnResultLine, error) {
	if s.data == nil {
		return nil, errors.ErrFetch
	}
	var sb strings.Builder
	sb.WriteString("option originalto")
	// option originalto [ except <network> ] [ header <name> ]
	if s.data.Except != "" {
		sb.WriteString(" except ")
		sb.WriteString(s.data.Except)
	}
	if s.data.Header != "" {
		sb.WriteString(" header ")
		sb.WriteString(s.data.Header)
	}
	return []common.ReturnResultLine{
		{
			Data:    sb.String(),
			Comment: s.data.Comment,
		},
	}, nil
}
