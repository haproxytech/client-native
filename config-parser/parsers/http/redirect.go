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

package http

import (
	"github.com/haproxytech/client-native/v5/config-parser/common"
	"github.com/haproxytech/client-native/v5/config-parser/errors"
	"github.com/haproxytech/client-native/v5/config-parser/parsers/http/actions"
	"github.com/haproxytech/client-native/v5/config-parser/types"
)

type Redirect struct {
	Name        string
	data        []types.Action
	preComments []string // comments that appear before the actual line
}

func (h *Redirect) Init() {
	h.Name = "redirect"
	h.data = []types.Action{}
}

func (h *Redirect) ParseHTTPResponse(response types.Action, parts []string, comment string) error {
	err := response.Parse(parts, types.HTTP, comment)
	if err != nil {
		return &errors.ParseError{Parser: "Redirect", Line: ""}
	}
	h.data = append(h.data, response)
	return nil
}

func (h *Redirect) Parse(line string, parts []string, comment string) (string, error) {
	if len(parts) >= 2 && parts[0] == "redirect" {
		adjusted := append([]string{""}, parts...)
		err := h.ParseHTTPResponse(&actions.Redirect{}, adjusted, comment)
		if err != nil {
			return "", err
		}
		return "", nil
	}
	return "", &errors.ParseError{Parser: "Redirect", Line: line}
}

func (h *Redirect) Result() ([]common.ReturnResultLine, error) {
	if len(h.data) == 0 {
		return nil, errors.ErrFetch
	}
	result := make([]common.ReturnResultLine, len(h.data))
	for index, res := range h.data {
		result[index] = common.ReturnResultLine{
			Data:    res.String(),
			Comment: res.GetComment(),
		}
	}
	return result, nil
}
