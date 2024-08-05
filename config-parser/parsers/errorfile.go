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

	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

// 200, 400, 401, 403, 404, 405, 407, 408, 410,
// 413, 425, 429, 500, 501, 502, 503, and 504
var errorFileAllowedCode = map[string]struct{}{ //nolint:gochecknoglobals
	"200": {},
	"400": {},
	"401": {},
	"403": {},
	"404": {},
	"405": {},
	"407": {},
	"408": {},
	"410": {},
	"413": {},
	"425": {},
	"429": {},
	"500": {},
	"501": {},
	"502": {},
	"503": {},
	"504": {},
}

type ErrorFile struct {
	data        []types.ErrorFile
	preComments []string // comments that appear before the actual line
}

func (l *ErrorFile) Init() {
	l.data = []types.ErrorFile{}
	l.preComments = []string{}
}

func (l *ErrorFile) parse(line string, parts []string, comment string) (*types.ErrorFile, error) {
	if len(parts) < 3 {
		return nil, &errors.ParseError{Parser: "ErrorFile", Line: line}
	}
	errorfile := &types.ErrorFile{
		File:    parts[2],
		Comment: comment,
	}
	code := parts[1]
	if _, ok := errorFileAllowedCode[code]; !ok {
		return nil, &errors.ParseError{Parser: "ErrorFile", Line: line}
	}
	errorfile.Code = code
	return errorfile, nil
}

func (l *ErrorFile) Result() ([]common.ReturnResultLine, error) {
	if len(l.data) == 0 {
		return nil, errors.ErrFetch
	}
	result := make([]common.ReturnResultLine, len(l.data))
	for index, data := range l.data {
		result[index] = common.ReturnResultLine{
			Data:    fmt.Sprintf("errorfile %s %s", data.Code, data.File),
			Comment: data.Comment,
		}
	}
	return result, nil
}
