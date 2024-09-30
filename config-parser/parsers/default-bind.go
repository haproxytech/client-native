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

type DefaultBind struct {
	data        *types.DefaultBind
	preComments []string // comments that appear before the actual line
}

func (d *DefaultBind) Parse(line string, parts []string, comment string) (string, error) {
	if parts[0] == "default-bind" && len(parts) > 1 {
		paramsBindOptions, _ := params.ParseBindOptions(parts[1:])
		d.data = &types.DefaultBind{
			Params:  paramsBindOptions,
			Comment: comment,
		}
		return "", nil
	}
	return "", &errors.ParseError{Parser: "DefaultBind", Line: line}
}

func (d *DefaultBind) Result() ([]common.ReturnResultLine, error) {
	if d.data == nil {
		return nil, errors.ErrFetch
	}
	var result strings.Builder
	result.WriteString("default-bind")
	options := params.BindOptionsString(d.data.Params)
	if options != "" {
		result.WriteString(" ")
		result.WriteString(options)
	}
	return []common.ReturnResultLine{
		{
			Data:    result.String(),
			Comment: d.data.Comment,
		},
	}, nil
}
