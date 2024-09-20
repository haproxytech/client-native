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
	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

type StatsTimeout struct {
	data        *types.StringC
	preComments []string // comments that appear before the actual line
}

func (s *StatsTimeout) Parse(line string, parts []string, comment string) (string, error) {
	if len(parts) > 1 && parts[0] == "stats" && parts[1] == "timeout" {
		if len(parts) < 3 {
			return "", &errors.ParseError{Parser: "StatsTimeout", Line: line, Message: "Parse error"}
		}
		s.data = &types.StringC{
			Value:   parts[2],
			Comment: comment,
		}
		return "", nil
	}
	return "", &errors.ParseError{Parser: "StatsTimeout", Line: line}
}

func (s *StatsTimeout) Result() ([]common.ReturnResultLine, error) {
	if s.data == nil {
		return nil, errors.ErrFetch
	}
	return []common.ReturnResultLine{
		{
			Data:    "stats timeout " + s.data.Value,
			Comment: s.data.Comment,
		},
	}, nil
}
