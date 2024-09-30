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
	"fmt"
	"strings"

	"github.com/haproxytech/client-native/v5/config-parser/common"
	"github.com/haproxytech/client-native/v5/config-parser/errors"
	"github.com/haproxytech/client-native/v5/config-parser/types"
)

type DefaultPath struct {
	data        *types.DefaultPath
	preComments []string // comments that appear before the actual line
}

/*
default-path { current | config | parent | origin <path> }
*/
func (d *DefaultPath) Parse(line string, parts []string, comment string) (string, error) {
	if len(parts) > 1 && parts[0] == "default-path" &&
		(parts[1] == "current" || parts[1] == "config" || parts[1] == "parent" || parts[1] == "origin") {
		data := &types.DefaultPath{
			Type:    parts[1],
			Comment: comment,
		}
		if parts[1] == "origin" {
			if len(parts) != 3 {
				return "", errors.ErrInvalidData
			}
			data.Path = parts[2]
		}
		d.data = data
		return "", nil
	}
	return "", &errors.ParseError{Parser: "default-path", Line: line}
}

func (d *DefaultPath) Result() ([]common.ReturnResultLine, error) {
	if d.data == nil {
		return nil, errors.ErrFetch
	}
	var sb strings.Builder
	fmt.Fprint(&sb, "default-path ", d.data.Type)
	if d.data.Type == "origin" {
		fmt.Fprint(&sb, " ", d.data.Path)
	}
	return []common.ReturnResultLine{
		{
			Data:    sb.String(),
			Comment: d.data.Comment,
		},
	}, nil
}
