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
	"strings"

	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/parsers/actions"
	httpActions "github.com/haproxytech/client-native/v6/config-parser/parsers/http/actions"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

type Checks struct {
	Name        string
	Mode        string
	data        []types.Action
	preComments []string // comments that appear before the actual line
}

func (h *Checks) Init() {
	h.Name = "http-check"
	h.data = []types.Action{}
}

func (h *Checks) parseHTTPCheck(request types.Action, parts []string, comment string) error {
	err := request.Parse(parts, types.HTTP, comment)
	if err != nil {
		return &errors.ParseError{Parser: "HTTPCheck", Line: "", Message: err.Error()}
	}
	h.data = append(h.data, request)
	return nil
}

func (h *Checks) Parse(line string, parts []string, comment string) (string, error) {
	if len(parts) < 2 {
		return "", &errors.ParseError{Parser: "HTTPCheck", Line: line, Message: "http-check type not provided"}
	}

	if parts[0] != h.Name {
		return "", &errors.ParseError{Parser: "HTTPCheck", Line: line, Message: "name is not http-check"}
	}

	if h.Mode == "frontend" {
		return "", &errors.ParseError{Parser: "HTTPCheck", Line: line, Message: "http-check cannot be used in frontend section"}
	}

	var err error

	switch {
	case parts[1] == "comment":
		err = h.parseHTTPCheck(&httpActions.CheckComment{}, parts, comment)
	case parts[1] == "connect":
		err = h.parseHTTPCheck(&actions.CheckConnect{}, parts, comment)
	case parts[1] == "disable-on-404":
		err = h.parseHTTPCheck(&httpActions.CheckDisableOn404{}, parts, comment)
	case parts[1] == "expect":
		err = h.parseHTTPCheck(&actions.CheckExpect{}, parts, comment)
	case parts[1] == "send":
		err = h.parseHTTPCheck(&httpActions.CheckSend{}, parts, comment)
	case parts[1] == "send-state":
		err = h.parseHTTPCheck(&httpActions.CheckSendState{}, parts, comment)
	case strings.HasPrefix(parts[1], "set-var("):
		err = h.parseHTTPCheck(&actions.SetVarCheck{}, parts, comment)
	case strings.HasPrefix(parts[1], "set-var-fmt("):
		err = h.parseHTTPCheck(&actions.SetVarFmtCheck{}, parts, comment)
	case strings.HasPrefix(parts[1], "unset-var("):
		err = h.parseHTTPCheck(&actions.UnsetVarCheck{}, parts, comment)
	default:
		err = &errors.ParseError{Parser: "HTTPCheck", Line: line, Message: "invalid http-check type provided"}
	}
	return "", err
}

func (h *Checks) Result() ([]common.ReturnResultLine, error) {
	if len(h.data) == 0 {
		return nil, errors.ErrFetch
	}
	result := make([]common.ReturnResultLine, len(h.data))
	for index, req := range h.data {
		result[index] = common.ReturnResultLine{
			Data:    "http-check " + req.String(),
			Comment: req.GetComment(),
		}
	}
	return result, nil
}
