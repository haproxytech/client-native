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
	"strconv"
	"strings"

	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

type OptionRedispatch struct {
	data        *types.OptionRedispatch
	preComments []string // comments that appear before the actual line
}

/*
option redispatch <interval>
*/
func (s *OptionRedispatch) Parse(line string, parts []string, comment string) (string, error) {
	if len(parts) > 1 && parts[0] == "option" && parts[1] == "redispatch" {
		data := &types.OptionRedispatch{
			Comment: comment,
		}
		if len(parts) > 2 {
			// see the interval
			if interval, err := strconv.ParseInt(parts[2], 10, 64); err == nil {
				data.Interval = &interval
			} else {
				return "", errors.ErrInvalidData
			}
		}
		s.data = data
		return "", nil
	}
	if len(parts) > 2 && parts[0] == "no" && parts[1] == "option" && parts[2] == "redispatch" {
		data := &types.OptionRedispatch{
			NoOption: true,
			Comment:  comment,
		}
		s.data = data
		return "", nil
	}
	return "", &errors.ParseError{Parser: "option redispatch", Line: line}
}

func (s *OptionRedispatch) Result() ([]common.ReturnResultLine, error) {
	if s.data == nil {
		return nil, errors.ErrFetch
	}
	var sb strings.Builder
	if s.data.NoOption {
		sb.WriteString("no ")
	}
	sb.WriteString("option redispatch")
	if !s.data.NoOption {
		if s.data.Interval != nil {
			sb.WriteString(" ")
			sb.WriteString(strconv.FormatInt(*s.data.Interval, 10))
		}
	}
	return []common.ReturnResultLine{
		{
			Data:    sb.String(),
			Comment: s.data.Comment,
		},
	}, nil
}
