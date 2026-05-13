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

type TypeSmtpchk struct {
	data        *types.TypeSmtpchk
	preComments []string // comments that appear before the actual line
}

/*
type smtpchk <hello> <domain>
*/
func (s *TypeSmtpchk) Parse(line string, parts []string, comment string) (string, error) {
	if len(parts) > 1 && parts[0] == "type" && parts[1] == "smtpchk" {
		data := &types.TypeSmtpchk{
			Comment: comment,
		}
		if len(parts) > 2 {
			if len(parts) != 4 {
				return "", errors.ErrInvalidData
			}
			data.Hello = parts[2]
			data.Domain = parts[3]
			if data.Hello != "EHLO" {
				data.Hello = "HELO"
			}
		}
		s.data = data
		return "", nil
	}
	if len(parts) > 2 && parts[0] == "no" && parts[1] == "type" && parts[2] == "smtpchk" {
		data := &types.TypeSmtpchk{
			NoType:  true,
			Comment: comment,
		}
		s.data = data
		return "", nil
	}
	return "", &errors.ParseError{Parser: "type smtpchk", Line: line}
}

func (s *TypeSmtpchk) Result() ([]common.ReturnResultLine, error) {
	if s.data == nil {
		return nil, errors.ErrFetch
	}
	var sb strings.Builder
	if s.data.NoType {
		sb.WriteString("no ")
	}
	sb.WriteString("type smtpchk")
	if s.data.Hello != "" && !s.data.NoType {
		sb.WriteString(" ")
		sb.WriteString(s.data.Hello)
		sb.WriteString(" ")
		sb.WriteString(s.data.Domain)
	}
	return []common.ReturnResultLine{
		{
			Data:    sb.String(),
			Comment: s.data.Comment,
		},
	}, nil
}
