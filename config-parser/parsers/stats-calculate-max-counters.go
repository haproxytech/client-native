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

	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

type StatsCalculateMaxCounters struct {
	data        *types.StringC
	preComments []string // comments that appear before the actual line
}

func (s *StatsCalculateMaxCounters) Parse(line string, parts []string, comment string) (string, error) {
	if len(parts) > 1 && parts[0] == "stats" && parts[1] == "calculate-max-counters" {
		if len(parts) < 3 {
			return "", &errors.ParseError{Parser: "StatsCalculateMaxCounters", Line: line, Message: "Parse error"}
		}
		if parts[2] != "on" && parts[2] != "off" {
			return "", &errors.ParseError{Parser: "StatsCalculateMaxCounters", Line: line, Message: fmt.Sprintf("invalid value %q (allowed: on, off)", parts[2])}
		}
		s.data = &types.StringC{
			Value:   parts[2],
			Comment: comment,
		}
		return "", nil
	}
	return "", &errors.ParseError{Parser: "StatsCalculateMaxCounters", Line: line}
}

func (s *StatsCalculateMaxCounters) Result() ([]common.ReturnResultLine, error) {
	if s.data == nil {
		return nil, errors.ErrFetch
	}
	return []common.ReturnResultLine{
		{
			Data:    "stats calculate-max-counters " + s.data.Value,
			Comment: s.data.Comment,
		},
	}, nil
}
