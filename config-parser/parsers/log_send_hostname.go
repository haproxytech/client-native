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

type LogSendHostName struct {
	data        *types.StringC
	preComments []string // comments that appear before the actual line
}

func (p *LogSendHostName) Parse(line string, parts []string, comment string) (string, error) {
	if parts[0] != "log-send-hostname" {
		return "", &errors.ParseError{Parser: "log-send-hostname", Line: line}
	}
	p.data = &types.StringC{
		Comment: comment,
	}
	if len(parts) > 1 {
		p.data.Value = parts[1]
	}
	return "", nil
}

func (p *LogSendHostName) Result() ([]common.ReturnResultLine, error) {
	if p.data == nil {
		return nil, errors.ErrFetch
	}
	data := "log-send-hostname"
	if p.data.Value != "" {
		data = fmt.Sprintf("log-send-hostname %s", p.data.Value)
	}
	return []common.ReturnResultLine{
		{
			Data:    data,
			Comment: p.data.Comment,
		},
	}, nil
}
