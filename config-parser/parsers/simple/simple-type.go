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

package simple

import (
	"fmt"
	"strings"

	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

type Type struct {
	Name        string
	name        string
	data        *types.SimpleType
	preComments []string // comments that appear before the actual line
}

func (o *Type) Init() {
	if !strings.HasPrefix(o.Name, "type") {
		o.name = o.Name
		o.Name = "type " + o.Name
	}
	o.data = nil
	o.preComments = nil
}

func (o *Type) Parse(line string, parts []string, comment string) (string, error) {
	if len(parts) > 1 && parts[0] == "type" && parts[1] == o.name {
		o.data = &types.SimpleType{
			Comment: comment,
		}
		return "", nil
	}
	return "", &errors.ParseError{Parser: "type " + o.name, Line: line}
}

func (o *Type) Result() ([]common.ReturnResultLine, error) {
	if o.data == nil {
		return nil, errors.ErrFetch
	}
	noType := ""
	if o.data.NoType {
		noType = "no "
	}
	return []common.ReturnResultLine{
		{
			Data:    fmt.Sprintf("%stype %s", noType, o.name),
			Comment: o.data.Comment,
		},
	}, nil
}
