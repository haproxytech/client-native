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
	"strings"

	"github.com/haproxytech/client-native/v5/config-parser/common"
	"github.com/haproxytech/client-native/v5/config-parser/errors"
	"github.com/haproxytech/client-native/v5/config-parser/types"
)

type CompressionAlgoRes struct {
	data        *types.StringSliceC
	preComments []string // comments that appear before the actual line
}

func (c *CompressionAlgoRes) Parse(line string, parts []string, comment string) (string, error) {
	if len(parts) < 3 {
		return "", &errors.ParseError{Parser: "CompressionAlgoRes", Line: line, Message: "Parse error"}
	}
	c.data = &types.StringSliceC{
		Value:   parts[2:],
		Comment: comment,
	}
	return "", nil
}

func (c *CompressionAlgoRes) Result() ([]common.ReturnResultLine, error) {
	if c.data == nil || len(c.data.Value) == 0 {
		return nil, errors.ErrFetch
	}
	var result strings.Builder
	result.WriteString("compression algo-res")

	for _, typereq := range c.data.Value {
		result.WriteString(" ")
		result.WriteString(typereq)
	}

	return []common.ReturnResultLine{
		{
			Data:    result.String(),
			Comment: c.data.Comment,
		},
	}, nil
}
