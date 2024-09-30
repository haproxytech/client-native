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
	"github.com/haproxytech/client-native/v5/config-parser/params"
	"github.com/haproxytech/client-native/v5/config-parser/types"
)

type Bind struct {
	data        []types.Bind
	preComments []string // comments that appear before the actual line
}

func (h *Bind) parse(line string, parts []string, comment string) (*types.Bind, error) {
	if len(parts) >= 2 {
		data := &types.Bind{
			Path:    parts[1],
			Comment: comment,
		}
		if len(parts) > 2 {
			paramsBindOptions, err := params.ParseBindOptions(parts[2:])
			if err != nil {
				return nil, &errors.ParseError{Parser: "BindLines", Line: line}
			}
			data.Params = paramsBindOptions
		}
		return data, nil
	}
	return nil, &errors.ParseError{Parser: "BindLines", Line: line}
}

func (h *Bind) Result() ([]common.ReturnResultLine, error) {
	if len(h.data) == 0 {
		return nil, errors.ErrFetch
	}
	result := make([]common.ReturnResultLine, len(h.data))
	for index, req := range h.data {
		var sb strings.Builder
		sb.WriteString("bind ")
		sb.WriteString(req.Path)
		options := params.BindOptionsString(req.Params)
		if options != "" {
			sb.WriteString(" ")
			sb.WriteString(options)
		}
		result[index] = common.ReturnResultLine{
			Data:    sb.String(),
			Comment: req.Comment,
		}
	}
	return result, nil
}
