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

type Option struct {
	Name        string
	name        string
	data        *types.SimpleOption
	preComments []string // comments that appear before the actual line
}

func (o *Option) Init() {
	if !strings.HasPrefix(o.Name, "option") {
		o.name = o.Name
		o.Name = fmt.Sprintf("option %s", o.Name)
	}
	o.data = nil
	o.preComments = nil
}

func (o *Option) Parse(line string, parts []string, comment string) (string, error) {
	if len(parts) > 1 && parts[0] == "option" && parts[1] == o.name {
		o.data = &types.SimpleOption{
			Comment: comment,
		}
		return "", nil
	}
	if len(parts) > 2 && parts[0] == "no" && parts[1] == "option" && parts[2] == o.name {
		o.data = &types.SimpleOption{
			NoOption: true,
			Comment:  comment,
		}
		return "", nil
	}
	return "", &errors.ParseError{Parser: fmt.Sprintf("option %s", o.name), Line: line}
}

func (o *Option) Result() ([]common.ReturnResultLine, error) {
	if o.data == nil {
		return nil, errors.ErrFetch
	}
	noOption := ""
	if o.data.NoOption {
		noOption = "no "
	}
	return []common.ReturnResultLine{
		{
			Data:    fmt.Sprintf("%soption %s", noOption, o.name),
			Comment: o.data.Comment,
		},
	}, nil
}
