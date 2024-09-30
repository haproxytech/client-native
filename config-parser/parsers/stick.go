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
	"github.com/haproxytech/client-native/v5/config-parser/types"
)

type Stick struct {
	data        []types.Stick
	preComments []string // comments that appear before the actual line
}

func (h *Stick) Parse(line string, parts []string, comment string) (string, error) {
	if len(parts) >= 2 && parts[0] == "stick" {
		command, condition := common.SplitRequest(parts[2:])
		data := types.Stick{
			Pattern: command[0],
			Comment: comment,
		}
		if len(command) > 2 {
			data.Table = command[2]
		}
		if len(condition) > 1 {
			data.Cond = condition[0]
			data.CondTest = strings.Join(condition[1:], " ")
		}
		switch parts[1] {
		case "match", "on", "store-request", "store-response":
			data.Type = parts[1]
		default:
			return "", &errors.ParseError{Parser: "Stick", Line: line}
		}
		h.data = append(h.data, data)
		return "", nil
	}
	return "", &errors.ParseError{Parser: "Stick", Line: line}
}

func (h *Stick) Result() ([]common.ReturnResultLine, error) {
	if len(h.data) == 0 {
		return nil, errors.ErrFetch
	}
	result := make([]common.ReturnResultLine, len(h.data))
	for index, req := range h.data {
		var data strings.Builder
		data.WriteString("stick ")
		data.WriteString(req.Type)
		data.WriteString(" ")
		data.WriteString(req.Pattern)
		if req.Table != "" {
			data.WriteString(" table ")
			data.WriteString(req.Table)
		}
		if req.Cond != "" {
			data.WriteString(" ")
			data.WriteString(req.Cond)
			data.WriteString(" ")
			data.WriteString(req.CondTest)
		}
		result[index] = common.ReturnResultLine{
			Data:    data.String(),
			Comment: req.Comment,
		}
	}
	return result, nil
}
