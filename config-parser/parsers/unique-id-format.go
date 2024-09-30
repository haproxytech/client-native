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
	"fmt"
	"strings"

	"github.com/haproxytech/client-native/v5/config-parser/common"
	"github.com/haproxytech/client-native/v5/config-parser/errors"
	"github.com/haproxytech/client-native/v5/config-parser/types"
)

type UniqueIDFormat struct {
	data        *types.UniqueIDFormat
	preComments []string // comments that appear before the actual line
}

func (p *UniqueIDFormat) Parse(line string, parts []string, comment string) (string, error) {
	if len(parts) < 2 {
		return "", &errors.ParseError{Parser: "unique-id-format", Line: line}
	}
	if parts[0] != "unique-id-format" {
		return "", &errors.ParseError{Parser: "unique-id-format", Line: line}
	}
	p.data = &types.UniqueIDFormat{
		LogFormat: strings.Join(parts[1:], " "),
		Comment:   comment,
	}
	return "", nil
}

func (p *UniqueIDFormat) Result() ([]common.ReturnResultLine, error) {
	if p.data == nil {
		return nil, errors.ErrFetch
	}
	return []common.ReturnResultLine{
		{
			Data:    fmt.Sprintf("unique-id-format %s", p.data.LogFormat),
			Comment: p.data.Comment,
		},
	}, nil
}
