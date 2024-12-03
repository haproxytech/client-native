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
package quic

import (
	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/parsers/actions"
	quic_actions "github.com/haproxytech/client-native/v6/config-parser/parsers/quic/actions"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

type Initial struct {
	Name string
	// Mode string
	data        []types.Action
	preComments []string // comments that appear before the actual line
}

func (h *Initial) Init() {
	h.Name = "quic-initial"
	h.data = []types.Action{}
	h.preComments = []string{}
}

func (h *Initial) Parse(line string, parts []string, comment string) (string, error) {
	var err error
	if len(parts) == 0 {
		return "", &errors.ParseError{Parser: "QUICInitial", Line: line, Message: "missing attribute"}
	}

	if parts[0] != h.Name {
		return "", &errors.ParseError{Parser: "QUICInitial", Line: line, Message: "expected attribute http-after-response"}
	}

	if len(parts) == 1 {
		return "", &errors.ParseError{Parser: "QUICInitial", Line: line, Message: "expected action for http-after-response"}
	}

	switch parts[1] {
	case "reject":
		err = h.ParseQuicInitial(&actions.Reject{}, parts, comment)
	case "accept":
		err = h.ParseQuicInitial(&actions.Accept{}, parts, comment)
	case "send-retry":
		err = h.ParseQuicInitial(&quic_actions.SendRetry{}, parts, comment)
	case "dgram-drop":
		err = h.ParseQuicInitial(&quic_actions.DgramDrop{}, parts, comment)
	default:
		return "", &errors.ParseError{Parser: "QuicInitial", Line: line}
	}

	return "", err
}

func (h *Initial) Result() ([]common.ReturnResultLine, error) {
	if len(h.data) == 0 {
		return nil, errors.ErrFetch
	}

	result := make([]common.ReturnResultLine, len(h.data))
	for index, req := range h.data {
		result[index] = common.ReturnResultLine{
			Data:    "quic-initial " + req.String(),
			Comment: req.GetComment(),
		}
	}
	return result, nil
}

func (h *Initial) ParseQuicInitial(request types.Action, parts []string, comment string) error {
	err := request.Parse(parts, types.QUIC, comment)
	if err != nil {
		return &errors.ParseError{Parser: "QuicInitial", Line: ""}
	}

	h.data = append(h.data, request)

	return nil
}
