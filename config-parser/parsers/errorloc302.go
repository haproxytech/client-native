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

type ErrorLoc302 struct {
	data        *types.ErrorLoc302
	preComments []string // comments that appear before the actual line
}

func (l *ErrorLoc302) Parse(line string, parts []string, comment string) (string, error) {
	if len(parts) < 3 {
		return "", &errors.ParseError{Parser: "ErrorLoc302", Line: line}
	}
	errorLoc := &types.ErrorLoc302{
		URL:     parts[2],
		Comment: comment,
	}
	code := parts[1]
	if _, ok := errorFileAllowedCode[code]; !ok {
		return "", &errors.ParseError{Parser: "ErrorLoc302", Line: line}
	}
	errorLoc.Code = code
	l.data = errorLoc
	return "", nil
}

func (l *ErrorLoc302) Result() ([]common.ReturnResultLine, error) {
	if l.data == nil {
		return nil, errors.ErrFetch
	}
	return []common.ReturnResultLine{
		{
			Data:    fmt.Sprintf("errorloc302 %s %s", l.data.Code, l.data.URL),
			Comment: l.data.Comment,
		},
	}, nil
}
