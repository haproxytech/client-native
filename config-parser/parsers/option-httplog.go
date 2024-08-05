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

type OptionHTTPLog struct {
	data        *types.OptionHTTPLog
	preComments []string // comments that appear before the actual line
}

func (o *OptionHTTPLog) Parse(line string, parts []string, comment string) (string, error) {
	if len(parts) > 2 && parts[0] == "option" && parts[1] == "httplog" && parts[2] == "clf" {
		o.data = &types.OptionHTTPLog{
			Comment: comment,
			Clf:     true,
		}
		return "", nil
	}
	if len(parts) > 1 && parts[0] == "option" && parts[1] == "httplog" {
		o.data = &types.OptionHTTPLog{
			Comment: comment,
		}
		return "", nil
	}
	if len(parts) > 3 && parts[0] == "no" && parts[1] == "option" && parts[2] == "httplog" && parts[3] == "clf" {
		o.data = &types.OptionHTTPLog{
			NoOption: true,
			Comment:  comment,
			Clf:      true,
		}
		return "", nil
	}
	if len(parts) > 2 && parts[0] == "no" && parts[1] == "option" && parts[2] == "httplog" {
		o.data = &types.OptionHTTPLog{
			NoOption: true,
			Comment:  comment,
		}
		return "", nil
	}
	return "", &errors.ParseError{Parser: "option httplog", Line: line}
}

func (o *OptionHTTPLog) Result() ([]common.ReturnResultLine, error) {
	if o.data == nil {
		return nil, errors.ErrFetch
	}
	clf := ""
	if o.data.Clf {
		clf = " clf"
	}
	noOption := ""
	if o.data.NoOption {
		noOption = "no "
	}
	return []common.ReturnResultLine{
		{
			Data:    fmt.Sprintf("%soption httplog%s", noOption, clf),
			Comment: o.data.Comment,
		},
	}, nil
}
