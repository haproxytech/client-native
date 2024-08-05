/*
Copyright 2022 HAProxy Technologies

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

type H1CaseAdjust struct {
	data        []types.H1CaseAdjust
	preComments []string // comments that appear before the actual line
}

func (ca *H1CaseAdjust) parse(line string, parts []string, comment string) (*types.H1CaseAdjust, error) {
	if len(parts) != 3 {
		return nil, &errors.ParseError{Parser: "H1CaseAdjust", Line: line, Message: "Parse error"}
	}

	h1CaseAdjust := &types.H1CaseAdjust{
		From:    parts[1],
		To:      parts[2],
		Comment: comment,
	}
	return h1CaseAdjust, nil
}

func (ca *H1CaseAdjust) Result() ([]common.ReturnResultLine, error) {
	if len(ca.data) == 0 {
		return nil, errors.ErrFetch
	}
	result := make([]common.ReturnResultLine, len(ca.data))
	for index, h1CaseAdjust := range ca.data {
		result[index] = common.ReturnResultLine{
			Data:    fmt.Sprintf("h1-case-adjust %s %s", h1CaseAdjust.From, h1CaseAdjust.To),
			Comment: h1CaseAdjust.Comment,
		}
	}
	return result, nil
}
