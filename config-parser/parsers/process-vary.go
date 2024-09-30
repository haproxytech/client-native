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
	"github.com/haproxytech/client-native/v5/config-parser/common"
	"github.com/haproxytech/client-native/v5/config-parser/errors"
	"github.com/haproxytech/client-native/v5/config-parser/types"
)

type ProcessVary struct {
	data        *types.ProcessVary
	preComments []string // comments that appear before the actual line
}

func (p *ProcessVary) parse(line string, parts []string, comment string) (*types.ProcessVary, error) {
	if len(parts) != 2 {
		return nil, &errors.ParseError{Parser: "ProcessVary", Line: line, Message: "Parse error"}
	}
	processVary := &types.ProcessVary{
		Comment: comment,
	}
	switch parts[1] {
	case "on":
		processVary.On = true
	case "off":
		processVary.On = false
	default:
		return nil, &errors.ParseError{Parser: "ProcessVary", Line: line, Message: "Parse error"}
	}
	return processVary, nil
}

func (p *ProcessVary) Parse(line string, parts []string, comment string) (string, error) {
	if parts[0] == "process-vary" {
		data, err := p.parse(line, parts, comment)
		if err != nil {
			return "", &errors.ParseError{Parser: "ProcessVary", Line: line}
		}
		p.data = data
		return "", nil
	}
	return "", &errors.ParseError{Parser: "ProcessVary", Line: line}
}

func (p *ProcessVary) Result() ([]common.ReturnResultLine, error) {
	if p.data == nil {
		return nil, errors.ErrFetch
	}
	result := make([]common.ReturnResultLine, 1)

	if p.data.On {
		result[0] = common.ReturnResultLine{
			Data:    "process-vary on",
			Comment: p.data.Comment,
		}
		return result, nil
	}
	result[0] = common.ReturnResultLine{
		Data:    "process-vary off",
		Comment: p.data.Comment,
	}
	return result, nil
}
