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
	"github.com/haproxytech/client-native/v5/config-parser/common"
	"github.com/haproxytech/client-native/v5/config-parser/errors"
	"github.com/haproxytech/client-native/v5/config-parser/types"
)

type MasterWorker struct {
	data        *types.Enabled
	preComments []string // comments that appear before the actual line
}

func (m *MasterWorker) Parse(line string, parts []string, comment string) (string, error) {
	if parts[0] == "master-worker" {
		m.data = &types.Enabled{
			Comment: comment,
		}
		return "", nil
	}
	return "", &errors.ParseError{Parser: "MasterWorker", Line: line}
}

func (m *MasterWorker) Result() ([]common.ReturnResultLine, error) {
	if m.data == nil {
		return nil, errors.ErrFetch
	}
	return []common.ReturnResultLine{
		{
			Data:    "master-worker",
			Comment: m.data.Comment,
		},
	}, nil
}
