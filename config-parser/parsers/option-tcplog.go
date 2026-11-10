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

type OptionTCPLog struct {
	data        *types.OptionTCPLog
	preComments []string // comments that appear before the actual line
}

func (o *OptionTCPLog) Parse(line string, parts []string, comment string) (string, error) {
	switch {
	case len(parts) == 2 && parts[0] == "option" && parts[1] == "tcplog":
		o.data = &types.OptionTCPLog{
			Comment: comment,
		}
		return "", nil
	case len(parts) == 3 && parts[0] == "option" && parts[1] == "tcplog" && parts[2] == "clf":
		o.data = &types.OptionTCPLog{
			Comment: comment,
			Clf:     true,
		}
		return "", nil
	case len(parts) == 3 && parts[0] == "no" && parts[1] == "option" && parts[2] == "tcplog":
		o.data = &types.OptionTCPLog{
			NoOption: true,
			Comment:  comment,
		}
		return "", nil
	case len(parts) == 4 && parts[0] == "no" && parts[1] == "option" && parts[2] == "tcplog" && parts[3] == "clf":
		o.data = &types.OptionTCPLog{
			NoOption: true,
			Comment:  comment,
			Clf:      true,
		}
		return "", nil
	}
	return "", &errors.ParseError{Parser: "option tcplog", Line: line}
}

func (o *OptionTCPLog) Result() ([]common.ReturnResultLine, error) {
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
			Data:    fmt.Sprintf("%soption tcplog%s", noOption, clf),
			Comment: o.data.Comment,
		},
	}, nil
}
