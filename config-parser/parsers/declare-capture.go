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

type DeclareCapture struct {
	data        []types.DeclareCapture
	preComments []string // comments that appear before the actual line
}

func (dc *DeclareCapture) parse(line string, parts []string, comment string) (*types.DeclareCapture, error) {
	if len(parts) != 5 {
		return nil, &errors.ParseError{Parser: "DeclareCapture", Line: line, Message: "Parse error"}
	}

	if parts[3] != "len" {
		return nil, &errors.ParseError{Parser: "DeclareCapture", Line: line, Message: "Parse error"}
	}

	length, err := strconv.ParseInt(parts[4], 10, 64)
	if err != nil {
		return nil, &errors.ParseError{Parser: "DeclareCapture", Line: line, Message: "Parse error"}
	}
	declareCapture := &types.DeclareCapture{
		Type:    parts[2],
		Length:  length,
		Comment: comment,
	}
	return declareCapture, nil
}

func (dc *DeclareCapture) Result() ([]common.ReturnResultLine, error) {
	if len(dc.data) == 0 {
		return nil, errors.ErrFetch
	}
	result := make([]common.ReturnResultLine, len(dc.data))
	for index, declareCapture := range dc.data {
		result[index] = common.ReturnResultLine{
			Data:    fmt.Sprintf("declare capture %s len %d", declareCapture.Type, declareCapture.Length),
			Comment: declareCapture.Comment,
		}
	}
	return result, nil
}
