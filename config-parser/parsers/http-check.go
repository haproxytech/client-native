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
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

type HTTPCheckV2 struct {
	data        []types.HTTPCheckV2
	preComments []string // comments that appear before the actual line
}

func (h *HTTPCheckV2) parse(line string, parts []string, comment string) (*types.HTTPCheckV2, error) {
	if len(parts) < 2 {
		return nil, &errors.ParseError{Parser: "HttpCheck", Line: line}
	}

	hc := &types.HTTPCheckV2{
		Comment: comment,
		Type:    parts[1],
	}
	if len(parts) == 2 {
		return hc, nil
	}

	if len(parts) >= 3 {
		if parts[2] == "!" {
			hc.ExclamationMark = true
			hc.Match = parts[3]
			hc.Pattern = strings.Join(parts[4:], " ")
			return hc, nil
		}
		hc.Match = parts[2]
		hc.Pattern = strings.Join(parts[3:], " ")
		return hc, nil
	}
	return nil, &errors.ParseError{Parser: "HttpCheck", Line: line}
}

func (h *HTTPCheckV2) Result() ([]common.ReturnResultLine, error) {
	if len(h.data) == 0 {
		return nil, errors.ErrFetch
	}

	result := make([]common.ReturnResultLine, len(h.data))
	for index, c := range h.data {
		var sb strings.Builder
		sb.WriteString("http-check")
		if c.Type != "" {
			sb.WriteString(" ")
			sb.WriteString(c.Type)
		}
		if c.ExclamationMark {
			sb.WriteString(" !")
		}
		if c.Match != "" {
			sb.WriteString(" ")
			sb.WriteString(c.Match)
		}
		if c.Pattern != "" {
			sb.WriteString(" ")
			sb.WriteString(c.Pattern)
		}
		result[index] = common.ReturnResultLine{
			Data:    sb.String(),
			Comment: c.Comment,
		}
	}
	return result, nil
}
