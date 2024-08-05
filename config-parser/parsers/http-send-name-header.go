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
	"strings"

	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

type HTTPSendNameHeader struct {
	data        *types.HTTPSendNameHeader
	preComments []string // comments that appear before the actual line
}

func (m *HTTPSendNameHeader) Parse(line string, parts []string, comment string) (string, error) {
	if parts[0] != "http-send-name-header" {
		return "", &errors.ParseError{Parser: "HTTPSendNameHeader", Line: line}
	}
	m.data = &types.HTTPSendNameHeader{}
	if len(parts) > 1 {
		m.data.Name = parts[1]
	}
	return "", nil
}

func (m *HTTPSendNameHeader) Result() ([]common.ReturnResultLine, error) {
	if m.data == nil {
		return nil, errors.ErrFetch
	}
	var sb strings.Builder
	sb.WriteString("http-send-name-header")
	if m.data.Name != "" {
		sb.WriteString(" ")
		sb.WriteString(m.data.Name)
	}
	return []common.ReturnResultLine{
		{
			Data:    sb.String(),
			Comment: m.data.Comment,
		},
	}, nil
}
