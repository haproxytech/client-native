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
	"strconv"
	"strings"

	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

type BindProcess struct {
	data        *types.BindProcess
	preComments []string // comments that appear before the actual line
}

func (p *BindProcess) Parse(line string, parts []string, comment string) (string, error) {
	if parts[0] != "bind-process" || len(parts) < 2 {
		return "", &errors.ParseError{Parser: "BindProcess", Line: line, Message: "Parse error"}
	}

	switch parts[1] {
	case "all", "odd", "even":
		p.data = &types.BindProcess{
			Comment: comment,
			Process: parts[1],
		}
		return "", nil
	}

	for _, d := range parts[1:] {
		if strings.Contains(d, "-") {
			if len(parts) > 2 {
				return "", &errors.ParseError{Parser: "BindProcess", Line: line}
			}
			split := common.StringSplitIgnoreEmpty(d, '-')
			for _, s := range split {
				_, err := strconv.Atoi(s)
				if err != nil {
					return "", &errors.ParseError{Parser: "BindProcess", Line: line, Message: err.Error()}
				}
			}
		} else {
			_, err := strconv.Atoi(d)
			if err != nil {
				return "", &errors.ParseError{Parser: "BindProcess", Line: line, Message: err.Error()}
			}
		}
	}

	p.data = &types.BindProcess{
		Comment: comment,
		Process: strings.Join(parts[1:], " "),
	}
	return "", nil
}

func (p *BindProcess) Result() ([]common.ReturnResultLine, error) {
	if p.data == nil {
		return nil, errors.ErrFetch
	}

	return []common.ReturnResultLine{
		{
			Data:    fmt.Sprintf("bind-process %s", p.data.Process),
			Comment: p.data.Comment,
		},
	}, nil
}
