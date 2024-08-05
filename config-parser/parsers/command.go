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
	"strings"

	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

type Command struct {
	data        *types.Command
	preComments []string
}

func (p *Command) Parse(line string, parts []string, comment string) (string, error) {
	if len(parts) < 2 {
		return "", &errors.ParseError{Parser: "Command", Line: line, Message: "missing command"}
	}

	if parts[0] != "command" {
		return "", &errors.ParseError{Parser: "Command", Line: line, Message: fmt.Sprintf("expected command, got %s", parts[0])}
	}

	p.data = &types.Command{
		Args:    strings.Join(parts[1:], " "),
		Comment: comment,
	}
	return "", nil
}

func (p *Command) Result() ([]common.ReturnResultLine, error) {
	if p.data == nil {
		return nil, errors.ErrFetch
	}

	var sb strings.Builder

	sb.WriteString("command ")
	sb.WriteString(p.data.Args)

	return []common.ReturnResultLine{
		{
			Data:    sb.String(),
			Comment: p.data.Comment,
		},
	}, nil
}
