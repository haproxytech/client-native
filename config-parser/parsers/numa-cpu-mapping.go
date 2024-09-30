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

	"github.com/haproxytech/client-native/v5/config-parser/common"
	"github.com/haproxytech/client-native/v5/config-parser/errors"
	"github.com/haproxytech/client-native/v5/config-parser/types"
)

type NumaCPUMapping struct {
	data        *types.NumaCPUMapping
	preComments []string // comments that appear before the the actual line
}

func (n *NumaCPUMapping) Parse(line string, parts []string, comment string) (string, error) {
	if len(parts) > 0 && parts[0] == "numa-cpu-mapping" {
		n.data = &types.NumaCPUMapping{
			NoOption: false,
			Comment:  comment,
		}
		return "", nil
	} else if len(parts) > 1 && parts[0] == "no" && parts[1] == "numa-cpu-mapping" {
		n.data = &types.NumaCPUMapping{
			NoOption: true,
			Comment:  comment,
		}
		return "", nil
	}
	return "", &errors.ParseError{Parser: "NumaCPUMapping", Line: line}
}

func (n *NumaCPUMapping) Result() ([]common.ReturnResultLine, error) {
	if n.data == nil {
		return nil, errors.ErrFetch
	}
	noOption := ""
	if n.data.NoOption {
		noOption = "no "
	}
	return []common.ReturnResultLine{
		{
			Data:    fmt.Sprintf("%s%s", noOption, "numa-cpu-mapping"),
			Comment: n.data.Comment,
		},
	}, nil
}
