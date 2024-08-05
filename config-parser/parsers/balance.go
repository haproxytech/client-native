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
	"strings"

	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/params"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

type Balance struct {
	data        *types.Balance
	preComments []string // comments that appear before the actual line
}

func (p *Balance) parseBalanceParams(pb params.BalanceParams, line string, parts []string) (*types.Balance, error) {
	if len(parts) >= 2 {
		b, err := pb.Parse(parts[1:])
		if err != nil {
			return nil, &errors.ParseError{Parser: "Balance", Line: line}
		}
		data := &types.Balance{
			Params: b,
		}
		return data, nil
	}

	return nil, &errors.ParseError{Parser: "Balance", Line: line}
}

func (p *Balance) Parse(line string, parts []string, comment string) (string, error) {
	if parts[0] == "balance" {
		if len(parts) < 2 {
			return "", &errors.ParseError{Parser: "Balance", Line: line}
		}

		data := &types.Balance{
			Comment: comment,
		}

		var err error
		var pb *types.Balance

		switch parts[1] {
		case "roundrobin", "static-rr", "leastconn", "first", "source", "random", "rdp-cookie":
			data.Algorithm = parts[1]
		case "uri":
			pb, err = p.parseBalanceParams(&params.BalanceURI{}, line, parts)
			data.Algorithm = parts[1]
		case "url_param":
			pb, err = p.parseBalanceParams(&params.BalanceURLParam{}, line, parts)
			data.Algorithm = parts[1]
		default:
			switch {
			case parts[1] == "hash":
				pb, err = p.parseBalanceParams(&params.BalanceHash{}, line, parts)
				data.Algorithm = parts[1]
			case strings.HasPrefix(parts[1], "random(") && strings.HasSuffix(parts[1], ")"):
				pb, err = p.parseBalanceParams(&params.BalanceRandom{}, line, parts)
				data.Algorithm = "random"
			case strings.HasPrefix(parts[1], "rdp-cookie(") && strings.HasSuffix(parts[1], ")"):
				pb, err = p.parseBalanceParams(&params.BalanceRdpCookie{}, line, parts)
				data.Algorithm = "rdp-cookie"
			case strings.HasPrefix(parts[1], "hdr(") && strings.HasSuffix(parts[1], ")"):
				pb, err = p.parseBalanceParams(&params.BalanceHdr{}, line, parts)
				data.Algorithm = "hdr"
			default:
				return "", &errors.ParseError{Parser: "Balance", Line: line}
			}
		}

		if err != nil {
			return "", &errors.ParseError{Parser: "Balance", Line: line}
		}

		if pb != nil && pb.Params != nil {
			data.Params = pb.Params
		}

		p.data = data
		return "", nil
	}
	return "", &errors.ParseError{Parser: "Balance", Line: line}
}

func (p *Balance) Result() ([]common.ReturnResultLine, error) {
	if p.data == nil {
		return nil, errors.ErrFetch
	}

	params := ""
	if p.data.Params != nil {
		params = p.data.Params.String()
	}

	return []common.ReturnResultLine{
		{
			Data:    fmt.Sprintf("balance %s%s", p.data.Algorithm, params),
			Comment: p.data.Comment,
		},
	}, nil
}
