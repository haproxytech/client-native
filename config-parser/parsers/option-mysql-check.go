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

type OptionMysqlCheck struct {
	data        *types.OptionMysqlCheck
	preComments []string // comments that appear before the actual line
}

/*
option mysql-check [ user <username> [ post-41 ] ]
*/
func (s *OptionMysqlCheck) Parse(line string, parts []string, comment string) (string, error) {
	if len(parts) > 1 && parts[0] == "option" && parts[1] == "mysql-check" {
		data := &types.OptionMysqlCheck{
			Comment: comment,
		}
		if len(parts) > 2 {
			if len(parts) < 6 {
				if parts[2] != "user" {
					return "", errors.ErrInvalidData
				}
				if len(parts) < 4 {
					return "", errors.ErrInvalidData
				}
				data.User = parts[3]
				if len(parts) == 5 {
					if parts[4] != "post-41" && parts[4] != "pre-41" {
						return "", errors.ErrInvalidData
					}
					data.ClientVersion = parts[4]
				}
			} else {
				return "", errors.ErrInvalidData
			}
		}
		s.data = data
		return "", nil
	}
	return "", &errors.ParseError{Parser: "option mysql-check", Line: line}
}

func (s *OptionMysqlCheck) Result() ([]common.ReturnResultLine, error) {
	if s.data == nil {
		return nil, errors.ErrFetch
	}
	var sb strings.Builder
	sb.WriteString("option mysql-check")
	if s.data.User != "" {
		sb.WriteString(" user ")
		sb.WriteString(s.data.User)
	}
	if s.data.ClientVersion != "" {
		sb.WriteString(" ")
		sb.WriteString(s.data.ClientVersion)
	}
	return []common.ReturnResultLine{
		{
			Data:    sb.String(),
			Comment: s.data.Comment,
		},
	}, nil
}
