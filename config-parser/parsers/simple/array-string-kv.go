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

	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

type ArrayKeyValue struct {
	Name        string
	data        []types.StringKeyValueC
	preComments []string // comments that appear before the actual line
}

func (p *ArrayKeyValue) Parse(line string, parts []string, comment string) (string, error) {
	if parts[0] == p.Name {
		p.data = append(p.data, types.StringKeyValueC{
			Key:     parts[1],
			Value:   parts[2],
			Comment: comment,
		})
		return "", nil
	}
	return "", &errors.ParseError{Parser: "ArrayKeyValue", Line: line}
}

func (p *ArrayKeyValue) Result() ([]common.ReturnResultLine, error) {
	if len(p.data) == 0 {
		return nil, errors.ErrFetch
	}

	result := make([]common.ReturnResultLine, len(p.data))
	for index, req := range p.data {
		result[index] = common.ReturnResultLine{
			Data:    fmt.Sprintf("%s %s %s", p.Name, req.Key, req.Value),
			Comment: req.Comment,
		}
	}
	return result, nil
}
