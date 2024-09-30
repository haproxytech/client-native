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

type IgnorePersist struct {
	data        []types.IgnorePersist
	preComments []string // comments that appear before the actual line
}

func (m *IgnorePersist) parse(line string, parts []string, comment string) (*types.IgnorePersist, error) {
	if len(parts) != 3 {
		return nil, &errors.ParseError{Parser: "IgnorePersist", Line: line}
	}
	if parts[0] == "ignore-persist" {
		if parts[1] != "if" && parts[1] != "unless" {
			return nil, &errors.ParseError{Parser: "IgnorePersist", Line: line}
		}
		data := &types.IgnorePersist{
			Cond:     parts[1],
			CondTest: parts[2],
			Comment:  comment,
		}
		return data, nil
	}
	return nil, &errors.ParseError{Parser: "IgnorePersist", Line: line}
}

func (m *IgnorePersist) Result() ([]common.ReturnResultLine, error) {
	if len(m.data) == 0 {
		return nil, errors.ErrFetch
	}
	result := make([]common.ReturnResultLine, len(m.data))
	for i := range m.data {
		result[i] = common.ReturnResultLine{
			Data:    fmt.Sprintf("ignore-persist %s %s", m.data[i].Cond, m.data[i].CondTest),
			Comment: m.data[i].Comment,
		}
	}
	return result, nil
}
