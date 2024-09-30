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
	"github.com/haproxytech/client-native/v5/config-parser/params"
	"github.com/haproxytech/client-native/v5/config-parser/types"
)

type Server struct {
	data        []types.Server
	preComments []string // comments that appear before the actual line
}

func (h *Server) parse(line string, parts []string, comment string) (*types.Server, error) {
	if len(parts) >= 3 {
		data := &types.Server{
			Name:    parts[1],
			Address: parts[2],
			Params:  params.ParseServerOptions(parts[3:]),
			Comment: comment,
		}
		return data, nil
	}
	return nil, &errors.ParseError{Parser: "Server", Line: line}
}

func (h *Server) Result() ([]common.ReturnResultLine, error) {
	if len(h.data) == 0 {
		return nil, errors.ErrFetch
	}
	result := make([]common.ReturnResultLine, len(h.data))
	for index, req := range h.data {
		var sb strings.Builder
		sb.WriteString("server ")
		sb.WriteString(req.Name)
		sb.WriteString(" ")
		sb.WriteString(req.Address)
		params := params.ServerOptionsString(req.Params)
		if params != "" {
			sb.WriteString(" ")
			sb.WriteString(params)
		}

		result[index] = common.ReturnResultLine{
			Data:    sb.String(),
			Comment: req.Comment,
		}
	}
	return result, nil
}
