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

const ExternalCheckKeyword = "external-check"

type ExternalCheckCommand struct {
	data        *types.ExternalCheckCommand
	preComments []string // comments that appear before the actual line
}

// external-check command <command>
func (s *ExternalCheckCommand) Parse(line string, parts []string, comment string) (string, error) {
	if len(parts) == 3 && parts[0] == "external-check" && parts[1] == "command" {
		s.data = &types.ExternalCheckCommand{
			Command: parts[2],
			Comment: comment,
		}
		return "", nil
	}
	return "", &errors.ParseError{Parser: "external-check command", Line: line}
}

func (s *ExternalCheckCommand) Result() ([]common.ReturnResultLine, error) {
	if s.data == nil {
		return nil, errors.ErrFetch
	}
	var data string
	if s.data.Command != "" {
		data = fmt.Sprintf("external-check command %s", s.data.Command)
	}
	return []common.ReturnResultLine{
		{
			Data:    data,
			Comment: s.data.Comment,
		},
	}, nil
}
