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

	"github.com/haproxytech/client-native/v5/config-parser/common"
	"github.com/haproxytech/client-native/v5/config-parser/errors"
	"github.com/haproxytech/client-native/v5/config-parser/params"
	"github.com/haproxytech/client-native/v5/config-parser/types"
)

type Persist struct {
	data        *types.Persist
	preComments []string // comments that appear before the actual line
}

func (p *Persist) parsePersistParams(pb params.PersistParams, line string, parts []string) (*types.Persist, error) {
	if len(parts) >= 2 {
		b, err := pb.Parse(parts[1:])
		if err != nil {
			return nil, &errors.ParseError{Parser: "Persist", Line: line}
		}
		data := &types.Persist{
			Params: b,
		}
		return data, nil
	}

	return nil, &errors.ParseError{Parser: "Persist", Line: line}
}

func (p *Persist) Parse(line string, parts []string, comment string) (string, error) {
	var err error
	if parts[0] != "persist" {
		return "", &errors.ParseError{Parser: "Persist", Line: line}
	}
	if len(parts) != 2 {
		return "", &errors.ParseError{Parser: "Persist", Line: line}
	}

	data := &types.Persist{
		Comment: comment,
	}
	var pb *types.Persist

	switch parts[1] {
	case "rdp-cookie":
		data.Type = parts[1]
	default:
		switch {
		case strings.HasPrefix(parts[1], "rdp-cookie(") && strings.HasSuffix(parts[1], ")"):
			pb, err = p.parsePersistParams(&params.PersistRdpCookie{}, line, parts)
			data.Type = "rdp-cookie"
		default:
			return "", &errors.ParseError{Parser: "Persist", Line: line}
		}
	}
	if err != nil {
		return "", &errors.ParseError{Parser: "Persist", Line: line}
	}

	if pb != nil && pb.Params != nil {
		data.Params = pb.Params
	}

	p.data = data
	return "", nil
}

func (p *Persist) Result() ([]common.ReturnResultLine, error) {
	if p.data == nil {
		return nil, errors.ErrFetch
	}

	params := ""
	if p.data.Params != nil {
		params = p.data.Params.String()
	}

	return []common.ReturnResultLine{
		{
			Data:    fmt.Sprintf("persist %s%s", p.data.Type, params),
			Comment: p.data.Comment,
		},
	}, nil
}
