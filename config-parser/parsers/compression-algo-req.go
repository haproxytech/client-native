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

package parsers

import (
	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

type CompressionAlgoReq struct {
	data        *types.StringC
	preComments []string // comments that appear before the actual line
}

func (c *CompressionAlgoReq) Parse(line string, parts []string, comment string) (string, error) {
	if len(parts) < 3 {
		return "", &errors.ParseError{Parser: "CompressionAlgoReq", Line: line, Message: "Parse error"}
	}
	c.data = &types.StringC{
		Value:   parts[2],
		Comment: comment,
	}
	return "", nil
}

func (c *CompressionAlgoReq) Result() ([]common.ReturnResultLine, error) {
	if c.data == nil {
		return nil, errors.ErrFetch
	}
	return []common.ReturnResultLine{
		{
			Data:    "compression algo-req " + c.data.Value,
			Comment: c.data.Comment,
		},
	}, nil
}
