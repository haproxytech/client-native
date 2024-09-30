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

type OptionHttpchk struct {
	data        *types.OptionHttpchk
	preComments []string // comments that appear before the actual line
}

/*
option httpchk <uri>
option httpchk <method> <uri>
option httpchk <method> <uri> <version>
*/
func (s *OptionHttpchk) Parse(line string, parts []string, comment string) (string, error) {
	if len(parts) > 1 && parts[0] == "option" && parts[1] == "httpchk" {
		switch len(parts) {
		case 2:
			s.data = &types.OptionHttpchk{
				Comment: comment,
			}
		case 3:
			s.data = &types.OptionHttpchk{
				URI:     parts[2],
				Comment: comment,
			}
		case 4:
			s.data = &types.OptionHttpchk{
				Method:  parts[2],
				URI:     parts[3],
				Comment: comment,
			}
		case 5:
			s.data = &types.OptionHttpchk{
				Method:  parts[2],
				URI:     parts[3],
				Version: parts[4],
				Comment: comment,
			}
		default: // > 5
			s.data = &types.OptionHttpchk{
				Method:  parts[2],
				URI:     parts[3],
				Version: strings.Join(parts[4:], " "),
				Comment: comment,
			}
		}
		return "", nil
	}
	return "", &errors.ParseError{Parser: "option httpchk", Line: line}
}

func (s *OptionHttpchk) Result() ([]common.ReturnResultLine, error) {
	if s.data == nil {
		return nil, errors.ErrFetch
	}
	var sb strings.Builder
	sb.WriteString("option httpchk")
	if s.data.Method != "" {
		sb.WriteString(" ")
		sb.WriteString(s.data.Method)
	}
	if s.data.URI != "" {
		sb.WriteString(" ")
		sb.WriteString(s.data.URI)
	}
	if s.data.Version != "" {
		sb.WriteString(" ")
		sb.WriteString(s.data.Version)
	}
	return []common.ReturnResultLine{
		{
			Data:    sb.String(),
			Comment: s.data.Comment,
		},
	}, nil
}
