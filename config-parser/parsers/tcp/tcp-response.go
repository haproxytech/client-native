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

type Responses struct {
	Name        string
	Mode        string // frontent, backend
	data        []types.TCPType
	preComments []string // comments that appear before the actual line
}

func (h *Responses) Init() {
	h.Name = "tcp-response"
	h.data = []types.TCPType{}
}

func (h *Responses) ParseTCPRequest(request types.TCPType, parts []string, comment string) error {
	err := request.Parse(parts, comment)
	if err != nil {
		return &errors.ParseError{Parser: "Responses", Line: ""}
	}
	h.data = append(h.data, request)
	return nil
}

func (h *Responses) Parse(line string, parts []string, comment string) (string, error) {
	if len(parts) >= 2 && parts[0] == "tcp-response" {
		var err error
		switch parts[1] {
		case "content":
			err = h.ParseTCPRequest(&tcptypes.Content{}, parts, comment)
		case "inspect-delay":
			err = h.ParseTCPRequest(&tcptypes.InspectDelay{}, parts, comment)
		default:
			return "", &errors.ParseError{Parser: "Responses", Line: line}
		}
		return "", err
	}
	return "", &errors.ParseError{Parser: "Responses", Line: line}
}

func (h *Responses) Result() ([]common.ReturnResultLine, error) {
	if len(h.data) == 0 {
		return nil, errors.ErrFetch
	}
	result := make([]common.ReturnResultLine, len(h.data))
	for index, req := range h.data {
		result[index] = common.ReturnResultLine{
			Data:    "tcp-response " + req.String(),
			Comment: req.GetComment(),
		}
	}
	return result, nil
}
