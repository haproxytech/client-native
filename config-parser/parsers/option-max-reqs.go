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
	"strconv"
	"strings"

	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

type OptionMaxReqs struct {
	data        *types.OptionMaxReqs
	preComments []string
}

func (p *OptionMaxReqs) Parse(line string, parts []string, comment string) (string, error) {
	if len(parts) != 3 {
		return "", &errors.ParseError{Parser: "OptionMaxReqs", Line: line, Message: "Missing required values"}
	}

	v, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return "", &errors.ParseError{Parser: "OptionMaxReqs", Line: line, Message: "Expecting a valid integer for option reqs"}
	}

	p.data = &types.OptionMaxReqs{
		Reqs:    v,
		Comment: comment,
	}

	return "", nil
}

func (p *OptionMaxReqs) Result() ([]common.ReturnResultLine, error) {
	if p.data == nil {
		return nil, errors.ErrFetch
	}

	var sb strings.Builder

	sb.WriteString("option max-reqs ")
	sb.WriteString(strconv.FormatInt(p.data.Reqs, 10))

	return []common.ReturnResultLine{
		{
			Data:    sb.String(),
			Comment: p.data.Comment,
		},
	}, nil
}
