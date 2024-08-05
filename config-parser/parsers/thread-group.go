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
	"strings"

	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

type ThreadGroup struct {
	data        []types.ThreadGroup
	preComments []string // comments that appear before the actual line
}

func (t *ThreadGroup) parse(line string, parts []string, comment string) (*types.ThreadGroup, error) {
	if len(parts) >= 3 {
		data := &types.ThreadGroup{
			Group:      parts[1],
			NumOrRange: parts[2],
			Comment:    comment,
		}
		return data, nil
	}
	return nil, &errors.ParseError{Parser: "ThreadGroup", Line: line}
}

func (t *ThreadGroup) Result() ([]common.ReturnResultLine, error) {
	if len(t.data) == 0 {
		return nil, errors.ErrFetch
	}
	result := make([]common.ReturnResultLine, len(t.data))
	for index, tg := range t.data {
		var sb strings.Builder
		sb.WriteString("thread-group")
		sb.WriteString(" ")
		sb.WriteString(tg.Group)
		sb.WriteString(" ")
		sb.WriteString(tg.NumOrRange)
		result[index] = common.ReturnResultLine{
			Data:    sb.String(),
			Comment: tg.Comment,
		}
	}
	return result, nil
}
