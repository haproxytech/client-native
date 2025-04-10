/*
Copyright 2025 HAProxy Technologies

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
	"github.com/haproxytech/client-native/v6/config-parser/params"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

type SSLFrontUse struct {
	data        []types.SSLFrontUse
	preComments []string // comments that appear before the actual line
}

func (h *SSLFrontUse) parse(line string, parts []string, comment string) (*types.SSLFrontUse, error) {
	// Minimum: ssl-f-use <key> <val>
	if len(parts) < 3 {
		return nil, &errors.ParseError{Parser: h.GetParserName(), Line: line}
	}
	opt, err := params.ParseSSLBindOptions(parts[1:])
	if err != nil {
		return nil, &errors.ParseError{Parser: h.GetParserName(), Line: line}
	}
	return &types.SSLFrontUse{Comment: comment, Params: opt}, nil
}

func (h *SSLFrontUse) Result() ([]common.ReturnResultLine, error) {
	if len(h.data) == 0 {
		return nil, errors.ErrFetch
	}
	result := make([]common.ReturnResultLine, len(h.data))
	for i, line := range h.data {
		result[i] = common.ReturnResultLine{
			Data:    common.SmartJoin(h.GetParserName(), params.SSLBindOptionsString(line.Params)),
			Comment: line.Comment,
		}
	}
	return result, nil
}
