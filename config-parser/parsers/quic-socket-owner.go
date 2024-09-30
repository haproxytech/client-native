/*
Copyright 2023 HAProxy Technologies

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

type QuicSocketOwner struct {
	data        *types.QuicSocketOwner
	preComments []string
}

func (p *QuicSocketOwner) Parse(line string, parts []string, comment string) (string, error) {
	if len(parts) == 2 && parts[0] == "tune.quic.socket-owner" {
		p.data = &types.QuicSocketOwner{}
		switch parts[1] {
		case "connection", "listener":
			p.data.Owner = parts[1]
		default:
			return "", &errors.ParseError{Parser: "tune.quic.socket-owner", Line: line}
		}
		return "", nil
	}
	return "", &errors.ParseError{Parser: "tune.quic.socket-owner", Line: line}
}

func (p *QuicSocketOwner) Result() ([]common.ReturnResultLine, error) {
	if p.data == nil || len(p.data.Owner) == 0 {
		return nil, errors.ErrFetch
	}
	data := fmt.Sprintf("tune.quic.socket-owner %s", p.data.Owner)
	return []common.ReturnResultLine{
		{
			Data: data,
		},
	}, nil
}
