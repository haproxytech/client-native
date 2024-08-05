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

type Timeout struct {
	Name        string
	name        string
	data        *types.SimpleTimeout
	preComments []string // comments that appear before the actual line
}

func (t *Timeout) Init() {
	if !strings.HasPrefix(t.Name, "timeout") {
		t.name = t.Name
		t.Name = fmt.Sprintf("timeout %s", t.Name)
	}
	t.data = nil
	t.preComments = []string{}
}

func (t *Timeout) Parse(line string, parts []string, comment string) (string, error) {
	if len(parts) > 2 && parts[0] == "timeout" && parts[1] == t.name {
		t.data = &types.SimpleTimeout{
			Value:   parts[2],
			Comment: comment,
		}
		return "", nil
	}
	return "", &errors.ParseError{Parser: fmt.Sprintf("timeout %s", t.name), Line: line}
}

func (t *Timeout) Result() ([]common.ReturnResultLine, error) {
	if t.data == nil {
		return nil, errors.ErrFetch
	}
	return []common.ReturnResultLine{
		{
			Data:    fmt.Sprintf("timeout %s %s", t.name, t.data.Value),
			Comment: t.data.Comment,
		},
	}, nil
}
