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
	"fmt"

	"github.com/haproxytech/client-native/v5/config-parser/common"
	"github.com/haproxytech/client-native/v5/config-parser/errors"
	"github.com/haproxytech/client-native/v5/config-parser/types"
)

type Nameserver struct {
	data        []types.Nameserver
	preComments []string // comments that appear before the actual line
}

func (l *Nameserver) parse(line string, parts []string, comment string) (*types.Nameserver, error) {
	if len(parts) >= 3 {
		return &types.Nameserver{
			Name:    parts[1],
			Address: parts[2],
			Comment: comment,
		}, nil
	}
	return nil, &errors.ParseError{Parser: "Nameserver", Line: line}
}

func (l *Nameserver) Result() ([]common.ReturnResultLine, error) {
	if len(l.data) == 0 {
		return nil, errors.ErrFetch
	}
	result := make([]common.ReturnResultLine, len(l.data))
	for index, nameserver := range l.data {
		result[index] = common.ReturnResultLine{
			Data:    fmt.Sprintf("nameserver %s %s", nameserver.Name, nameserver.Address),
			Comment: nameserver.Comment,
		}
	}
	return result, nil
}
