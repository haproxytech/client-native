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
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

type Mode struct {
	data        *types.StringC
	preComments []string // comments that appear before the actual line
}

func (p *Mode) Parse(line string, parts []string, comment string) (string, error) {
	if parts[0] == "mode" {
		if len(parts) < 2 {
			return "", &errors.ParseError{Parser: "Mode", Line: line, Message: "Parse error"}
		}
		if parts[1] == "http" || parts[1] == "tcp" || parts[1] == "log" {
			p.data = &types.StringC{
				Value:   parts[1],
				Comment: comment,
			}
			return "", nil
		}
		return "", &errors.ParseError{Parser: "Mode", Line: line}
	}
	return "", &errors.ParseError{Parser: "Mode", Line: line}
}

func (p *Mode) Result() ([]common.ReturnResultLine, error) {
	if p.data == nil {
		return nil, errors.ErrFetch
	}
	return []common.ReturnResultLine{
		{
			Data:    "mode " + p.data.Value,
			Comment: p.data.Comment,
		},
	}, nil
}
