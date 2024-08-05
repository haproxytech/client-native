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

	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

type ErrorFiles struct {
	data        []types.ErrorFiles
	preComments []string // comments that appear before the actual line
}

func (e *ErrorFiles) Init() {
	e.data = []types.ErrorFiles{}
	e.preComments = []string{}
}

func (e *ErrorFiles) parse(line string, parts []string, comment string) (*types.ErrorFiles, error) {
	if len(parts) < 2 {
		return nil, &errors.ParseError{Parser: "ErrorFiles", Line: line}
	}
	errorfiles := &types.ErrorFiles{
		Name:    parts[1],
		Comment: comment,
	}
	if len(parts) > 2 {
		codes := make([]int64, len(parts)-2)
		for i, code := range parts[2:] {
			if _, ok := errorFileAllowedCode[code]; !ok {
				return nil, &errors.ParseError{Parser: "ErrorFiles", Line: line}
			}
			intCode, err := strconv.ParseInt(code, 10, 0)
			if err != nil {
				return nil, &errors.ParseError{Parser: "ErrorFiles", Line: line}
			}
			codes[i] = intCode
		}
		errorfiles.Codes = codes
	}
	return errorfiles, nil
}

func (e *ErrorFiles) Result() ([]common.ReturnResultLine, error) {
	if len(e.data) == 0 {
		return nil, errors.ErrFetch
	}
	result := make([]common.ReturnResultLine, len(e.data))
	for index, data := range e.data {

		result[index] = common.ReturnResultLine{
			Data:    fmt.Sprintf("errorfiles %s", data.Name),
			Comment: data.Comment,
		}
		if len(data.Codes) > 0 {
			for _, code := range data.Codes {
				result[index].Data = fmt.Sprintf("%s %s", result[index].Data, strconv.FormatInt(code, 10))
			}
		}
	}
	return result, nil
}
