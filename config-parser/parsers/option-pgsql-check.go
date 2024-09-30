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

type OptionPgsqlCheck struct {
	data        *types.OptionPgsqlCheck
	preComments []string // comments that appear before the actual line
}

/*
option pgsql-check [ user <username> ]
*/
func (s *OptionPgsqlCheck) Parse(line string, parts []string, comment string) (string, error) {
	if len(parts) != 4 {
		return "", errors.ErrInvalidData
	}
	if parts[0] == "option" && parts[1] == "pgsql-check" {
		data := &types.OptionPgsqlCheck{
			Comment: comment,
		}

		if parts[2] != "user" {
			return "", errors.ErrInvalidData
		}
		data.User = parts[3]

		s.data = data
		return "", nil
	}
	return "", &errors.ParseError{Parser: "option pgsql-check", Line: line}
}

func (s *OptionPgsqlCheck) Result() ([]common.ReturnResultLine, error) {
	if s.data == nil {
		return nil, errors.ErrFetch
	}
	var sb strings.Builder
	sb.WriteString("option pgsql-check")
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
