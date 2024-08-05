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
	"strings"

	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

type StickTable struct {
	data        *types.StickTable
	preComments []string // comments that appear before the actual line
}

//nolint:gocognit
func (h *StickTable) parse(line string, parts []string, comment string) (*types.StickTable, error) {
	if len(parts) >= 3 && parts[0] == "stick-table" && parts[1] == "type" {
		index := 2
		data := &types.StickTable{
			Type:    parts[index],
			Comment: comment,
		}
		index++
		for index < len(parts) {
			switch parts[index] {
			case "len":
				index++
				if index == len(parts) {
					return nil, &errors.ParseError{Parser: "StickTable", Line: line}
				}
				data.Length = parts[index]
			case "size":
				index++
				if index == len(parts) {
					return nil, &errors.ParseError{Parser: "StickTable", Line: line}
				}
				data.Size = parts[index]
			case "expire":
				index++
				if index == len(parts) {
					return nil, &errors.ParseError{Parser: "StickTable", Line: line}
				}
				data.Expire = parts[index]
			case "nopurge":
				data.NoPurge = true
			case "peers":
				index++
				if index == len(parts) {
					return nil, &errors.ParseError{Parser: "StickTable", Line: line}
				}
				data.Peers = parts[index]
			case "srvkey":
				index++
				if index == len(parts) {
					return nil, &errors.ParseError{Parser: "StickTable", Line: line}
				}
				key := parts[index]
				if key != "addr" && key != "name" {
					return nil, &errors.ParseError{Parser: "StickTable", Line: line, Message: "invalid srvkey"}
				}
				data.SrvKey = key
			case "write-to":
				index++
				if index == len(parts) {
					return nil, &errors.ParseError{Parser: "StickTable", Line: line}
				}
				data.WriteTo = parts[index]
			case "store":
				index++
				if index == len(parts) {
					return nil, &errors.ParseError{Parser: "StickTable", Line: line}
				}
				data.Store = parts[index]
			default:
				return nil, &errors.ParseError{Parser: "StickTable", Line: line}
			}
			index++
		}
		return data, nil
	}
	return nil, &errors.ParseError{Parser: "StickTable", Line: line}
}

func (h *StickTable) Parse(line string, parts []string, comment string) (string, error) {
	if parts[0] == "stick-table" {
		data, err := h.parse(line, parts, comment)
		if err != nil {
			return "", &errors.ParseError{Parser: "StickTable", Line: line}
		}
		h.data = data
		return "", nil
	}
	return "", &errors.ParseError{Parser: "StickTable", Line: line}
}

func (h *StickTable) Result() ([]common.ReturnResultLine, error) {
	if h.data == nil {
		return nil, errors.ErrFetch
	}
	result := make([]common.ReturnResultLine, 1)
	req := h.data

	var data strings.Builder
	data.WriteString("stick-table type ")
	data.WriteString(req.Type)
	if req.Length != "" {
		data.WriteString(" len ")
		data.WriteString(req.Length)
	}
	if req.Size != "" {
		data.WriteString(" size ")
		data.WriteString(req.Size)
	}
	if req.Expire != "" {
		data.WriteString(" expire ")
		data.WriteString(req.Expire)
	}
	if req.NoPurge {
		data.WriteString(" nopurge")
	}
	if req.Peers != "" {
		data.WriteString(" peers ")
		data.WriteString(req.Peers)
	}
	if req.SrvKey != "" {
		data.WriteString(" srvkey ")
		data.WriteString(req.SrvKey)
	}
	if req.WriteTo != "" {
		data.WriteString(" write-to ")
		data.WriteString(req.WriteTo)
	}
	if req.Store != "" {
		data.WriteString(" store ")
		data.WriteString(req.Store)
	}
	result[0] = common.ReturnResultLine{
		Data:    data.String(),
		Comment: req.Comment,
	}
	return result, nil
}
