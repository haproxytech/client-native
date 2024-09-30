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
	"github.com/haproxytech/client-native/v5/config-parser/types"
)

type ConfigSnippet struct {
	data        *types.StringSliceC
	active      bool
	preComments []string // comments that appear before the actual line
}

func (p *ConfigSnippet) Parse(line string, parts []string, comment string) (string, error) {
	if comment != "" {
		commentParts := strings.Fields(comment)
		if len(commentParts) > 1 || commentParts[0] == "##_config-snippet_###" {
			switch commentParts[1] {
			case "BEGIN":
				p.active = true
				p.data = &types.StringSliceC{}
				return "snippet_beg", nil
			case "END":
				p.active = false
				return "snippet_end", nil
			default:
				p.data.Value = append(p.data.Value, strings.TrimSpace(line))
				return "", nil
			}
		}
	}
	if p.active {
		p.data.Value = append(p.data.Value, strings.TrimSpace(line))
		return "", nil
	}
	return "", &errors.ParseError{Parser: "ConfigSnippet", Line: line}
}

func (p *ConfigSnippet) Result() ([]common.ReturnResultLine, error) {
	if p.data == nil {
		return nil, errors.ErrFetch
	}
	return []common.ReturnResultLine{
		{
			Data:    "###_config-snippet_### BEGIN\n  " + strings.Join(p.data.Value, "\n  ") + "\n  ###_config-snippet_### END",
			Comment: p.data.Comment,
		},
	}, nil
}
