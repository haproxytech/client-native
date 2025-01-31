/*
Copyright 2025 HAProxy Technologies

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
	"slices"

	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

type OnLogStep struct {
	data        []types.OnLogStep
	preComments []string // comments that appear before the actual line
}

func (o *OnLogStep) GetParserName() string {
	return "on"
}

func (o *OnLogStep) Parse(line string, parts []string, comment string) (string, error) {
	if parts[0] == o.GetParserName() {
		data, err := o.parse(line, parts, comment)
		if err != nil {
			if _, ok := err.(*errors.ParseError); ok { //nolint:errorlint
				return "", err
			}
			return "", &errors.ParseError{Parser: "OnLogStep", Line: line}
		}
		o.data = append(o.data, *data)
		return "", nil
	}
	return "", &errors.ParseError{Parser: "OnLogStep", Line: line}
}

func (o *OnLogStep) parse(line string, parts []string, comment string) (*types.OnLogStep, error) {
	on := o.GetParserName()
	if parts[0] != on {
		return nil, &errors.ParseError{Parser: on, Line: line}
	}
	if len(parts) < 3 {
		return nil, &errors.ParseError{Parser: on, Line: line, Message: "Parse error: not enough arguments"}
	}

	var step types.OnLogStep
	step.Comment = comment

	step.Step = parts[1]
	if !o.isKnownStep(step.Step) {
		return nil, &errors.ParseError{Parser: on, Line: line, Message: "Parse error: invalid step name"}
	}

	if parts[2] == "drop" {
		if len(parts) > 3 {
			return nil, &errors.ParseError{Parser: on, Line: line, Message: "Parse error: nothing is allowed after 'drop'"}
		}
		step.Drop = true
		return &step, nil
	}

	parts = parts[2:]

	if len(parts) < 2 {
		return nil, &errors.ParseError{Parser: on, Line: line, Message: "Parse error: not enough arguments"}
	}

	for {
		switch parts[0] {
		case "format":
			step.Format = parts[1]
		case "sd":
			step.Sd = parts[1]
		default:
			return nil, &errors.ParseError{Parser: on, Line: line, Message: "Parse error: bad trailing argument"}
		}
		if len(parts) == 2 {
			break // success
		}
		if len(parts) > 2 {
			parts = parts[2:]
		} else {
			return nil, &errors.ParseError{Parser: on, Line: line, Message: "Parse error: bad trailing argument"}
		}
	}

	return &step, nil
}

func (o *OnLogStep) Result() ([]common.ReturnResultLine, error) {
	if len(o.data) == 0 {
		return nil, errors.ErrFetch
	}

	result := make([]common.ReturnResultLine, len(o.data))
	on := o.GetParserName()

	for i, step := range o.data {
		if step.Drop {
			result[i] = common.ReturnResultLine{
				Data:    common.SmartJoin(on, step.Step, "drop"),
				Comment: step.Comment,
			}
			continue
		}
		result[i] = common.ReturnResultLine{
			Data: common.SmartJoin(on, step.Step,
				o.maybe("format", step.Format), step.Format,
				o.maybe("sd", step.Sd), step.Sd),
			Comment: step.Comment,
		}
	}
	return result, nil
}

func (o *OnLogStep) maybe(keyword, value string) string {
	if len(value) > 0 {
		return keyword
	}
	return ""
}

func (o *OnLogStep) isKnownStep(step string) bool {
	return slices.Contains([]string{"accept", "any", "close", "connect", "error", "request", "response", "http-req", "http-res", "http-after-res", "quic-init", "tcp-req-conn", "tcp-req-cont", "tcp-req-sess"}, step)
}
