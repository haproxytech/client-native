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
	"strings"

	"github.com/haproxytech/client-native/v5/config-parser/common"
	"github.com/haproxytech/client-native/v5/config-parser/errors"
	"github.com/haproxytech/client-native/v5/config-parser/types"
)

type UnProcessed struct {
	Name        string
	data        []types.UnProcessed
	preComments []string // comments that appear before the actual line
}

func (u *UnProcessed) Init() {
	u.Name = ""
	u.data = []types.UnProcessed{}
}

func (u *UnProcessed) Parse(line string, parts []string, comment string) (string, error) {
	u.data = append(u.data, types.UnProcessed{
		Value: strings.TrimSpace(line),
	})
	return "", nil
}

func (u *UnProcessed) Result() ([]common.ReturnResultLine, error) {
	if len(u.data) == 0 {
		return nil, errors.ErrFetch
	}
	result := make([]common.ReturnResultLine, len(u.data))
	for index, d := range u.data {
		result[index] = common.ReturnResultLine{
			Data:    d.Value,
			Comment: "",
		}
	}
	return result, nil
}
