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
	"fmt"
	"strconv"

	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

type StatsMaxconn struct {
	data        *types.Int64C
	preComments []string // comments that appear before the actual line
}

func (s *StatsMaxconn) Parse(line string, parts []string, comment string) (string, error) {
	if len(parts) > 1 && parts[0] == "stats" && parts[1] == "maxconn" {
		if len(parts) < 3 {
			return "", &errors.ParseError{Parser: "StatsMaxconn", Line: line, Message: "Parse error"}
		}
		connections, err := strconv.ParseInt(parts[2], 10, 64)
		if err != nil {
			return "", &errors.ParseError{Parser: "StatsMaxconn", Line: line, Message: "Parse error"}
		}
		s.data = &types.Int64C{
			Value:   connections,
			Comment: comment,
		}
		return "", nil
	}
	return "", &errors.ParseError{Parser: "StatsMaxconn", Line: line}
}

func (s *StatsMaxconn) Result() ([]common.ReturnResultLine, error) {
	if s.data == nil {
		return nil, errors.ErrFetch
	}
	return []common.ReturnResultLine{
		{
			Data:    fmt.Sprintf("stats maxconn %d", s.data.Value),
			Comment: s.data.Comment,
		},
	}, nil
}
