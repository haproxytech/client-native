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

package extra

import (
	"github.com/haproxytech/client-native/v5/config-parser/common"
	"github.com/haproxytech/client-native/v5/config-parser/errors"
	"github.com/haproxytech/client-native/v5/config-parser/types"
)

type Comments struct {
	Name        string
	data        []types.Comments
	preComments []string
}

func (p *Comments) Init() {
	p.Name = "#"
	p.data = []types.Comments{}
}

func (p *Comments) Parse(line string, parts []string, comment string) (string, error) {
	if line[0] == '#' {
		p.data = append(p.data, types.Comments{
			Value: comment,
		})
		return "", nil
	}
	return "", &errors.ParseError{Parser: "Comments", Line: line}
}

func (p *Comments) Result() ([]common.ReturnResultLine, error) {
	if len(p.data) == 0 {
		return nil, errors.ErrFetch
	}
	result := make([]common.ReturnResultLine, len(p.data))
	for index, comment := range p.data {
		result[index] = common.ReturnResultLine{
			Data:    "# " + comment.Value,
			Comment: "",
		}
	}
	return result, nil
}
