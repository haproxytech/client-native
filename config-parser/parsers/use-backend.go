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
	"fmt"
	"strings"

	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

type UseBackend struct {
	data        []types.UseBackend
	preComments []string // comments that appear before the actual line
}

func (h *UseBackend) parse(line string, parts []string, comment string) (*types.UseBackend, error) {
	if len(parts) >= 2 {
		_, condition := common.SplitRequest(parts[2:])
		data := &types.UseBackend{
			Name:    parts[1],
			Comment: comment,
		}
		if len(condition) > 1 {
			data.Cond = condition[0]
			data.CondTest = strings.Join(condition[1:], " ")
		}
		return data, nil
	}
	return nil, &errors.ParseError{Parser: "UseBackend", Line: line}
}

func (h *UseBackend) Result() ([]common.ReturnResultLine, error) {
	if len(h.data) == 0 {
		return nil, errors.ErrFetch
	}
	result := make([]common.ReturnResultLine, len(h.data))
	for index, req := range h.data {
		condition := ""
		if req.Cond != "" {
			condition = fmt.Sprintf(" %s %s", req.Cond, req.CondTest)
		}
		result[index] = common.ReturnResultLine{
			Data:    fmt.Sprintf("use_backend %s%s", req.Name, condition),
			Comment: req.Comment,
		}
	}
	return result, nil
}
