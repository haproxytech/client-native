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

	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/params"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

type DgramBind struct {
	data        []types.DgramBind
	preComments []string // comments that appear before the actual line
}

func (h *DgramBind) parse(line string, parts []string, comment string) (*types.DgramBind, error) {
	if len(parts) >= 2 {
		data := &types.DgramBind{
			Path:    parts[1],
			Comment: comment,
		}
		if len(parts) > 2 {
			pr, err := params.ParseDgramBindOptions(parts[2:])
			if err != nil {
				return nil, err
			}
			data.Params = pr
		}
		return data, nil
	}
	return nil, &errors.ParseError{Parser: "DgramBindLines", Line: line}
}

func (h *DgramBind) Result() ([]common.ReturnResultLine, error) {
	if len(h.data) == 0 {
		return nil, errors.ErrFetch
	}
	result := make([]common.ReturnResultLine, len(h.data))
	for index, req := range h.data {
		var sb strings.Builder
		sb.WriteString("dgram-bind ")
		sb.WriteString(req.Path)
		options := params.DgramBindOptionsString(req.Params)
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
