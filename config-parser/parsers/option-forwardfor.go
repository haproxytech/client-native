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

type OptionForwardFor struct {
	data        *types.OptionForwardFor
	preComments []string // comments that appear before the actual line
}

/*
option forwardfor [ except <network> ] [ header <name> ] [ if-none ]
*/
func (s *OptionForwardFor) Parse(line string, parts []string, comment string) (string, error) {
	if len(parts) > 1 && parts[0] == "option" && parts[1] == "forwardfor" {
		data := &types.OptionForwardFor{
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
			case "if-none":
				data.IfNone = true
			default:
				return "", errors.ErrInvalidData
			}
			index++
		}
		s.data = data
		return "", nil
	}
	return "", &errors.ParseError{Parser: "option forwardfor", Line: line}
}

func (s *OptionForwardFor) Result() ([]common.ReturnResultLine, error) {
	if s.data == nil {
		return nil, errors.ErrFetch
	}
	var sb strings.Builder
	sb.WriteString("option forwardfor")
	// option forwardfor [ except <network> ] [ header <name> ] [ if-none ]
	if s.data.Except != "" {
		sb.WriteString(" except ")
		sb.WriteString(s.data.Except)
	}
	if s.data.Header != "" {
		sb.WriteString(" header ")
		sb.WriteString(s.data.Header)
	}
	if s.data.IfNone {
		sb.WriteString(" if-none")
	}
	return []common.ReturnResultLine{
		{
			Data:    sb.String(),
			Comment: s.data.Comment,
		},
	}, nil
}
