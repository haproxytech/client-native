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
	"slices"
	"strconv"

	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

var earlyHintsValues = []string{ //nolint:gochecknoglobals
	"on",
	"off",
	"only",
}

type EarlyHints struct {
	data        *types.EarlyHints
	preComments []string // comments that appear before the actual line
}

func (p *EarlyHints) Parse(line string, parts []string, comment string) (string, error) {
	if parts[0] != "early-hints" {
		return "", &errors.ParseError{Parser: "EarlyHints", Line: line}
	}
	if len(parts) < 2 {
		return "", &errors.ParseError{Parser: "EarlyHints", Line: line, Message: "Parse error"}
	}

	mode := parts[1]
	if !slices.Contains(earlyHintsValues, mode) {
		return "", &errors.ParseError{Parser: "EarlyHints", Line: line, Message: fmt.Sprintf("invalid mode %q (allowed: %v)", mode, earlyHintsValues)}
	}

	var ratio *int64
	if len(parts) > 2 {
		if parts[2] != "ratio" {
			return "", &errors.ParseError{Parser: "EarlyHints", Line: line, Message: fmt.Sprintf("invalid argument %q (allowed: ratio)", parts[2])}
		}
		if len(parts) != 4 {
			return "", &errors.ParseError{Parser: "EarlyHints", Line: line, Message: "Parse error"}
		}
		parsedRatio, err := strconv.ParseInt(parts[3], 10, 64)
		if err != nil || parsedRatio < 1 || parsedRatio > 99 {
			return "", &errors.ParseError{Parser: "EarlyHints", Line: line, Message: fmt.Sprintf("invalid ratio %q (allowed: 1-99)", parts[3])}
		}
		ratio = new(parsedRatio)
	}

	p.data = &types.EarlyHints{
		Mode:    mode,
		Ratio:   ratio,
		Comment: comment,
	}
	return "", nil
}

func (p *EarlyHints) Result() ([]common.ReturnResultLine, error) {
	if p.data == nil {
		return nil, errors.ErrFetch
	}
	data := "early-hints " + p.data.Mode
	if p.data.Ratio != nil {
		data += " ratio " + strconv.FormatInt(*p.data.Ratio, 10)
	}
	return []common.ReturnResultLine{
		{
			Data:    data,
			Comment: p.data.Comment,
		},
	}, nil
}
