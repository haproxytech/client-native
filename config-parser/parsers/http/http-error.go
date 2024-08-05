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

package http

import (
	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/errors"
	httpActions "github.com/haproxytech/client-native/v6/config-parser/parsers/http/actions"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

//nolint:revive
type HTTPErrors struct {
	Name        string
	data        []types.Action
	preComments []string // comments that appear before the actual line
}

func (h *HTTPErrors) Init() {
	h.Name = "http-error"
	h.data = []types.Action{}
	h.preComments = []string{}
}

func (h *HTTPErrors) ParseHTTPError(httpErr types.Action, parts []string, comment string) error {
	err := httpErr.Parse(parts, types.HTTP, comment)
	if err != nil {
		return &errors.ParseError{Parser: "HTTPErrorLines", Line: ""}
	}
	h.data = append(h.data, httpErr)
	return nil
}

func (h *HTTPErrors) Parse(line string, parts []string, comment string) (string, error) {
	if len(parts) < 3 || parts[0] != "http-error" {
		return "", &errors.ParseError{Parser: "HTTPErrorLines", Line: line}
	}
	return "", h.ParseHTTPError(&httpActions.Status{}, parts, comment)
}

func (h *HTTPErrors) Result() ([]common.ReturnResultLine, error) {
	if len(h.data) == 0 {
		return nil, errors.ErrFetch
	}
	result := make([]common.ReturnResultLine, len(h.data))
	for index, res := range h.data {
		result[index] = common.ReturnResultLine{
			Data:    "http-error " + res.String(),
			Comment: res.GetComment(),
		}
	}
	return result, nil
}
