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

type TypePgsqlCheck struct {
	data        *types.TypePgsqlCheck
	preComments []string // comments that appear before the actual line
}

/*
type pgsql-check [ user <username> ]
*/
func (s *TypePgsqlCheck) Parse(line string, parts []string, comment string) (string, error) {
	if len(parts) != 4 {
		return "", errors.ErrInvalidData
	}
	if parts[0] == "type" && parts[1] == "pgsql-check" {
		data := &types.TypePgsqlCheck{
			Comment: comment,
		}

		if parts[2] != "user" {
			return "", errors.ErrInvalidData
		}
		data.User = parts[3]

		s.data = data
		return "", nil
	}
	return "", &errors.ParseError{Parser: "type pgsql-check", Line: line}
}

func (s *TypePgsqlCheck) Result() ([]common.ReturnResultLine, error) {
	if s.data == nil {
		return nil, errors.ErrFetch
	}
	var sb strings.Builder
	sb.WriteString("type pgsql-check")
	if s.data.User != "" {
		sb.WriteString(" user ")
		sb.WriteString(s.data.User)
	}
	return []common.ReturnResultLine{
		{
			Data:    sb.String(),
			Comment: s.data.Comment,
		},
	}, nil
}
