/*
Copyright 2026 HAProxy Technologies

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

type FilterAPI struct {
	data        *types.StringC
	preComments []string // comments that appear before the actual line
}

func (p *FilterAPI) Parse(line string, parts []string, comment string) (string, error) {
	if len(parts) != 2 || parts[0] != "filter-api" {
		return "", &errors.ParseError{Parser: "FilterAPI", Line: line}
	}
	switch parts[1] {
	case "default", "downgrade-to-legacy":
	default:
		return "", &errors.ParseError{Parser: "FilterAPI", Line: line}
	}
	p.data = &types.StringC{Value: parts[1], Comment: comment}
	return "", nil
}

func (p *FilterAPI) Result() ([]common.ReturnResultLine, error) {
	if p.data == nil {
		return nil, errors.ErrFetch
	}
	return []common.ReturnResultLine{{
		Data:    "filter-api " + p.data.Value,
		Comment: p.data.Comment,
	}}, nil
}
