/*
Copyright 2024 HAProxy Technologies

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

type Trace struct {
	data        []types.Trace
	preComments []string // comments that appear before the actual line
}

func (t *Trace) parse(line string, parts []string, comment string) (*types.Trace, error) {
	if parts[0] != t.GetParserName() {
		return nil, &errors.ParseError{Parser: t.GetParserName(), Line: line}
	}
	if len(parts) < 2 {
		return nil, &errors.ParseError{Parser: t.GetParserName(), Line: line, Message: "Parse error: not enough arguments"}
	}

	return &types.Trace{
		Params:  parts[1:],
		Comment: comment,
	}, nil
}

func (t *Trace) Result() ([]common.ReturnResultLine, error) {
	if len(t.data) == 0 {
		return nil, errors.ErrFetch
	}

	result := make([]common.ReturnResultLine, len(t.data))

	for i, trace := range t.data {
		result[i] = common.ReturnResultLine{
			Data:    t.GetParserName() + " " + common.SmartJoin(trace.Params...),
			Comment: trace.Comment,
		}
	}
	return result, nil
}
