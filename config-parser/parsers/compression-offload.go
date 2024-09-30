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
	"github.com/haproxytech/client-native/v5/config-parser/common"
	"github.com/haproxytech/client-native/v5/config-parser/errors"
	"github.com/haproxytech/client-native/v5/config-parser/types"
)

type CompressionOffload struct {
	data        *types.Enabled
	preComments []string // comments that appear before the actual line
}

func (c *CompressionOffload) Parse(line string, parts []string, comment string) (string, error) {
	if len(parts) == 2 && parts[0] == "compression" && parts[1] == "offload" {
		c.data = &types.Enabled{
			Comment: comment,
		}
		return "", nil
	}
	return "", &errors.ParseError{Parser: "CompressionOffload", Line: line, Message: "Parse error"}
}

func (c *CompressionOffload) Result() ([]common.ReturnResultLine, error) {
	if c.data == nil {
		return nil, errors.ErrFetch
	}
	return []common.ReturnResultLine{
		{
			Data:    "compression offload",
			Comment: c.data.Comment,
		},
	}, nil
}
