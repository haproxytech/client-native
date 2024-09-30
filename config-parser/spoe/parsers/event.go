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

	"github.com/haproxytech/client-native/v5/config-parser/common"
	"github.com/haproxytech/client-native/v5/config-parser/errors"
	"github.com/haproxytech/client-native/v5/config-parser/spoe/types"
)

type Event struct {
	data        *types.Event
	preComments []string // comments that appear before the actual line
}

func (e *Event) Parse(line string, parts []string, comment string) (string, error) {
	if len(parts) < 2 || parts[0] != "event" {
		return "", &errors.ParseError{Parser: "EventLines", Line: line}
	}
	e.data = &types.Event{
		Name: parts[1],
	}
	_, condition := common.SplitRequest(parts[2:])
	if len(condition) > 1 {
		e.data.Cond = condition[0]
		e.data.CondTest = strings.Join(condition[1:], " ")
	}
	return "", nil
}

func (e *Event) Result() ([]common.ReturnResultLine, error) {
	if e.data == nil {
		return nil, errors.ErrFetch
	}
	var sb strings.Builder
	sb.WriteString("event ")
	sb.WriteString(e.data.Name)

	if e.data.Cond != "" {
		sb.WriteString(" ")
		sb.WriteString(e.data.Cond)
		sb.WriteString(" ")
		sb.WriteString(e.data.CondTest)
	}

	return []common.ReturnResultLine{
		{
			Data:    sb.String(),
			Comment: e.data.Comment,
		},
	}, nil
}
