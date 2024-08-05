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

package tcp

import (
	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/errors"
	tcptypes "github.com/haproxytech/client-native/v6/config-parser/parsers/tcp/types"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

type Requests struct {
	Name        string
	Mode        string // frontent, backend, listen, defaults
	data        []types.TCPType
	preComments []string // comments that appear before the actual line
}

func (h *Requests) Init() {
	h.Name = "tcp-request"
	h.data = []types.TCPType{}
}

func (h *Requests) ParseTCPRequest(request types.TCPType, parts []string, comment string) error {
	err := request.Parse(parts, comment)
	if err != nil {
		return &errors.ParseError{Parser: "TCPRequest", Line: ""}
	}
	h.data = append(h.data, request)

	return nil
}

func (h *Requests) Parse(line string, parts []string, comment string) (string, error) {
	var err error
	if parts[0] != "tcp-request" {
		return "", &errors.ParseError{Parser: "TCPRequest", Line: line}
	}
	if len(parts) < 2 {
		return "", &errors.ParseError{Parser: "TCPRequest", Line: line}
	}

	switch parts[1] {
	case "connection":
		if h.Mode == "backend" {
			return "", &errors.ParseError{Parser: "TCPRequest", Line: line}
		}
		err = h.ParseTCPRequest(&tcptypes.Connection{}, parts, comment)
	case "session":
		if h.Mode == "backend" {
			return "", &errors.ParseError{Parser: "TCPRequest", Line: line}
		}
		err = h.ParseTCPRequest(&tcptypes.Session{}, parts, comment)
	case "content":
		err = h.ParseTCPRequest(&tcptypes.Content{}, parts, comment)
	case "inspect-delay":
		err = h.ParseTCPRequest(&tcptypes.InspectDelay{}, parts, comment)
	default:
		return "", &errors.ParseError{Parser: "TCPRequest", Line: line}
	}

	return "", err
}

func (h *Requests) Result() ([]common.ReturnResultLine, error) {
	if len(h.data) == 0 {
		return nil, errors.ErrFetch
	}
	result := make([]common.ReturnResultLine, len(h.data))
	for index, req := range h.data {
		result[index] = common.ReturnResultLine{
			Data:    "tcp-request " + req.String(),
			Comment: req.GetComment(),
		}
	}
	return result, nil
}
