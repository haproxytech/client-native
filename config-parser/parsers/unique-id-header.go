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
	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

type UniqueIDHeader struct {
	data        *types.UniqueIDHeader
	preComments []string // comments that appear before the actual line
}

func (s *UniqueIDHeader) Parse(line string, parts []string, comment string) (string, error) {
	if len(parts) != 2 {
		return "", &errors.ParseError{Parser: "unique-id-header", Line: line}
	}
	if parts[0] != "unique-id-header" {
		return "", &errors.ParseError{Parser: "unique-id-header", Line: line}
	}
	s.data = &types.UniqueIDHeader{
		Name:    parts[1],
		Comment: comment,
	}
	return "", nil
}

func (s *UniqueIDHeader) Result() ([]common.ReturnResultLine, error) {
	if s.data == nil {
		return nil, errors.ErrFetch
	}
	return []common.ReturnResultLine{
		{
			Data:    "unique-id-header " + s.data.Name,
			Comment: s.data.Comment,
		},
	}, nil
}
