/*
Copyright 2026 HAProxy Technologies

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

type MaxThreadsPerGroup struct {
	data        *types.Int64C
	preComments []string // comments that appear before the actual line
}

func (m *MaxThreadsPerGroup) Parse(line string, parts []string, comment string) (string, error) {
	if parts[0] == "max-threads-per-group" {
		if len(parts) < 2 {
			return "", &errors.ParseError{Parser: "MaxThreadsPerGroup", Line: line, Message: "Parse error"}
		}
		num, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return "", &errors.ParseError{Parser: "MaxThreadsPerGroup", Line: line, Message: err.Error()}
		}
		m.data = &types.Int64C{
			Value:   num,
			Comment: comment,
		}
		return "", nil
	}
	return "", &errors.ParseError{Parser: "max-threads-per-group", Line: line}
}

func (m *MaxThreadsPerGroup) Result() ([]common.ReturnResultLine, error) {
	if m.data == nil {
		return nil, errors.ErrFetch
	}
	return []common.ReturnResultLine{
		{
			Data:    fmt.Sprintf("max-threads-per-group %d", m.data.Value),
			Comment: m.data.Comment,
		},
	}, nil
}
