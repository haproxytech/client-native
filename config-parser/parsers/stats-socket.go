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

type Socket struct {
	data        []types.Socket
	preComments []string // comments that appear before the actual line
}

func (l *Socket) parse(line string, parts []string, comment string) (*types.Socket, error) {
	if len(parts) < 3 {
		return nil, &errors.ParseError{Parser: "SocketSingle", Line: line, Message: "Parse error"}
	}
	paramsBindOptions, _ := params.ParseBindOptions(parts[3:])
	socket := &types.Socket{
		Path:    parts[2],
		Params:  paramsBindOptions,
		Comment: comment,
	}
	// s.value = elements[1:]
	return socket, nil
}

func (l *Socket) Result() ([]common.ReturnResultLine, error) {
	if len(l.data) == 0 {
		return nil, errors.ErrFetch
	}
	result := make([]common.ReturnResultLine, len(l.data))
	for index, socket := range l.data {
		var sb strings.Builder
		sb.WriteString("stats socket ")
		sb.WriteString(socket.Path)
		params := params.BindOptionsString(socket.Params)
		if params != "" {
			sb.WriteString(" ")
			sb.WriteString(params)
		}
		result[index] = common.ReturnResultLine{
			Data:    sb.String(),
			Comment: socket.Comment,
		}
	}
	return result, nil
}
