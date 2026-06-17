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

type FilterSequence struct {
	data        []types.FilterSequence
	preComments []string // comments that appear before the actual line
}

func (f *FilterSequence) Parse(line string, parts []string, comment string) (string, error) {
	if parts[0] != "filter-sequence" {
		return "", &errors.ParseError{Parser: "FilterSequence", Line: line}
	}
	if len(parts) < 3 {
		return "", &errors.ParseError{Parser: "FilterSequence", Line: line, Message: "Parse error"}
	}
	direction := parts[1]
	if direction != "request" && direction != "response" {
		return "", &errors.ParseError{Parser: "FilterSequence", Line: line, Message: "direction must be request or response"}
	}
	filtersList := strings.Split(parts[2], ",")
	cleaned := filtersList[:0]
	for _, name := range filtersList {
		name = strings.TrimSpace(name)
		if name != "" {
			cleaned = append(cleaned, name)
		}
	}
	if len(cleaned) == 0 {
		return "", &errors.ParseError{Parser: "FilterSequence", Line: line, Message: "missing filter list"}
	}
	f.data = append(f.data, types.FilterSequence{
		Direction: direction,
		Filters:   cleaned,
		Comment:   comment,
	})
	return "", nil
}

func (f *FilterSequence) Result() ([]common.ReturnResultLine, error) {
	if len(f.data) == 0 {
		return nil, errors.ErrFetch
	}
	result := make([]common.ReturnResultLine, len(f.data))
	for i, entry := range f.data {
		result[i] = common.ReturnResultLine{
			Data:    "filter-sequence " + entry.Direction + " " + strings.Join(entry.Filters, ","),
			Comment: entry.Comment,
		}
	}
	return result, nil
}
