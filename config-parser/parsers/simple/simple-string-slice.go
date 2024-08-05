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

package simple

import (
	"fmt"
	"strings"

	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

type StringSlice struct {
	Name        string
	data        *types.StringSliceC
	preComments []string // comments that appear before the actual line
}

func (s *StringSlice) Parse(line string, parts []string, comment string) (string, error) {
	if parts[0] != s.Name {
		return "", &errors.ParseError{Parser: s.Name, Line: line}
	}
	if len(parts) < 2 {
		return "", &errors.ParseError{Parser: s.Name, Line: line, Message: "Parse error"}
	}

	s.data = &types.StringSliceC{
		Value:   parts[1:],
		Comment: comment,
	}
	return "", nil
}

func (s *StringSlice) Result() ([]common.ReturnResultLine, error) {
	if s.data == nil {
		return nil, errors.ErrFetch
	}
	return []common.ReturnResultLine{
		{
			Data:    fmt.Sprintf("%s %s", s.Name, strings.Join(s.data.Value, " ")),
			Comment: s.data.Comment,
		},
	}, nil
}
