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
	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/params"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

type DefaultServer struct {
	data        []types.DefaultServer
	preComments []string // comments that appear before the actual line
}

func (h *DefaultServer) parse(line string, parts []string, comment string) (*types.DefaultServer, error) {
	if len(parts) >= 2 {
		sp, err := params.ParseServerOptions(parts[1:])
		if err != nil {
			return nil, err
		}
		data := &types.DefaultServer{
			Params:  sp,
			Comment: comment,
		}
		return data, nil
	}
	return nil, &errors.ParseError{Parser: "DefaultServer", Line: line}
}

func (h *DefaultServer) Result() ([]common.ReturnResultLine, error) {
	if len(h.data) == 0 {
		return nil, errors.ErrFetch
	}
	result := make([]common.ReturnResultLine, len(h.data))
	for index, req := range h.data {
		result[index] = common.ReturnResultLine{
			Data:    "default-server " + params.ServerOptionsString(req.Params),
			Comment: req.Comment,
		}
	}
	return result, nil
}
