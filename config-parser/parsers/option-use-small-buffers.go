/*
Copyright 2026 HAProxy Technologies

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

type OptionUseSmallBuffers struct {
	data        *types.OptionUseSmallBuffers
	preComments []string // comments that appear before the actual line
}

func (o *OptionUseSmallBuffers) Parse(line string, parts []string, comment string) (string, error) {
	idx := 0
	var noOption bool
	if len(parts) > 0 && parts[0] == "no" {
		noOption = true
		idx++
	}
	if len(parts) < idx+2 || parts[idx] != "option" || parts[idx+1] != "use-small-buffers" {
		return "", &errors.ParseError{Parser: "option use-small-buffers", Line: line}
	}
	out := &types.OptionUseSmallBuffers{
		NoOption: noOption,
		Comment:  comment,
	}
	for _, mode := range parts[idx+2:] {
		switch mode {
		case "queue":
			out.Queue = true
		case "l7-retries":
			out.L7Retries = true
		case "check":
			out.Check = true
		default:
			return "", &errors.ParseError{Parser: "option use-small-buffers", Line: line, Message: "unknown mode " + mode}
		}
	}
	o.data = out
	return "", nil
}

func (o *OptionUseSmallBuffers) Result() ([]common.ReturnResultLine, error) {
	if o.data == nil {
		return nil, errors.ErrFetch
	}
	var sb strings.Builder
	if o.data.NoOption {
		sb.WriteString("no ")
	}
	sb.WriteString("option use-small-buffers")
	if o.data.Queue {
		sb.WriteString(" queue")
	}
	if o.data.L7Retries {
		sb.WriteString(" l7-retries")
	}
	if o.data.Check {
		sb.WriteString(" check")
	}
	return []common.ReturnResultLine{
		{
			Data:    sb.String(),
			Comment: o.data.Comment,
		},
	}, nil
}
