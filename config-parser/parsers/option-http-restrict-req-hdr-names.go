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

type OptionHTTPRestrictReqHdrNames struct {
	data        *types.OptionHTTPRestrictReqHdrNames
	preComments []string // comments that appear before the actual line
}

func (o *OptionHTTPRestrictReqHdrNames) Parse(line string, parts []string, comment string) (string, error) {
	if len(parts) != 3 {
		return "", &errors.ParseError{Parser: "option http-restrict-req-hdr-names", Line: line}
	}
	o.data = &types.OptionHTTPRestrictReqHdrNames{Comment: comment}
	switch parts[2] {
	case "preserve", "delete", "reject":
		o.data.Policy = parts[2]
	default:
		return "", &errors.ParseError{Parser: "option http-restrict-req-hdr-names", Line: line}
	}
	return "", nil
}

func (o *OptionHTTPRestrictReqHdrNames) Result() ([]common.ReturnResultLine, error) {
	if o.data == nil {
		return nil, errors.ErrFetch
	}
	return []common.ReturnResultLine{
		{
			Data:    fmt.Sprintf("option http-restrict-req-hdr-names %s", o.data.Policy),
			Comment: o.data.Comment,
		},
	}, nil
}
