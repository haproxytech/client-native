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
	"fmt"
	"strings"

	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

type EmailAlert struct {
	data        []types.EmailAlert
	preComments []string // comments that appear before the actual line
}

func (e *EmailAlert) parse(line string, parts []string, comment string) (*types.EmailAlert, error) {
	if len(parts) != 3 {
		return nil, &errors.ParseError{Parser: "EmailAlert", Line: line}
	}

	attr, value := parts[1], parts[2]

	switch attr {
	case "from", "to":
		if !strings.Contains(value, "@") {
			return nil, &errors.ParseError{Parser: "EmailAlert", Line: line, Message: fmt.Sprintf("invalid email address: '%s'", value)}
		}
	case "level":
		// Must be a valid syslog severity level.
		if _, exists := logAllowedLevels[value]; !exists {
			return nil, &errors.ParseError{Parser: "EmailAlert", Line: line, Message: fmt.Sprintf("invalid email-alert log level '%s'", value)}
		}
	case "mailers", "myhostname":
	default:
		return nil, &errors.ParseError{Parser: "EmailAlert", Line: line, Message: fmt.Sprintf("unknown email-alert attribute '%s'", attr)}
	}

	data := &types.EmailAlert{
		Attribute: attr,
		Value:     value,
		Comment:   comment,
	}
	return data, nil
}

func (e *EmailAlert) Result() ([]common.ReturnResultLine, error) {
	if len(e.data) == 0 {
		return nil, errors.ErrFetch
	}
	result := make([]common.ReturnResultLine, len(e.data))
	for index, ea := range e.data {
		result[index] = common.ReturnResultLine{
			Data:    fmt.Sprintf("email-alert %s %s", ea.Attribute, ea.Value),
			Comment: ea.Comment,
		}
	}
	return result, nil
}
