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

type TypeHttpchk struct {
	data        *types.TypeHttpchk
	preComments []string // comments that appear before the actual line
}

/*
type httpchk <uri>
type httpchk <method> <uri>
type httpchk <method> <uri> <version>
*/
func (s *TypeHttpchk) Parse(line string, parts []string, comment string) (string, error) {
	if len(parts) > 1 && parts[0] == "type" && parts[1] == "httpchk" {
		switch len(parts) {
		case 2:
			s.data = &types.TypeHttpchk{
				Comment: comment,
			}
		case 3:
			s.data = &types.TypeHttpchk{
				URI:     parts[2],
				Comment: comment,
			}
		case 4:
			s.data = &types.TypeHttpchk{
				Method:  parts[2],
				URI:     parts[3],
				Comment: comment,
			}
		case 5:
			s.data = &types.TypeHttpchk{
				Method:  parts[2],
				URI:     parts[3],
				Version: parts[4],
				Comment: comment,
			}
		case 6:
			s.data = &types.TypeHttpchk{
				Method:  parts[2],
				URI:     parts[3],
				Version: parts[4],
				Host:    parts[5],
				Comment: comment,
			}
		default: // > 6
			return "", &errors.ParseError{Parser: "type httpchk", Line: line}
		}
		return "", nil
	}
	return "", &errors.ParseError{Parser: "type httpchk", Line: line}
}

func (s *TypeHttpchk) Result() ([]common.ReturnResultLine, error) {
	if s.data == nil {
		return nil, errors.ErrFetch
	}
	var sb strings.Builder
	sb.WriteString("type httpchk")
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
	if s.data.Host != "" {
		sb.WriteString(" ")
		sb.WriteString(s.data.Host)
	}
	return []common.ReturnResultLine{
		{
			Data:    sb.String(),
			Comment: s.data.Comment,
		},
	}, nil
}
