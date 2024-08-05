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
	"slices"
	"strconv"
	"strings"

	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

type Table struct {
	data        []types.Table
	preComments []string // comments that appear before the actual line
}

func (t *Table) parse(line string, parts []string, comment string) (*types.Table, error) { //nolint:gocognit
	indexSize := slices.IndexFunc(parts, func(s string) bool { return s == "size" })
	if indexSize >= 4 && len(parts) >= 6 && parts[0] == "table" && parts[2] == "type" {
		index := 3
		data := &types.Table{
			Name:    parts[1],
			Type:    parts[index],
			Comment: comment,
		}
		index++
		for index < len(parts) {
			switch parts[index] {
			case "len":
				index++
				if index == len(parts) {
					return nil, &errors.ParseError{Parser: "Table", Line: line}
				}
				length, err := strconv.Atoi(parts[index])
				if err != nil {
					return nil, err
				}

				length64 := int64(length)
				data.TypeLen = length64

			case "size":
				index++
				if index == len(parts) {
					return nil, &errors.ParseError{Parser: "Table", Line: line}
				}
				data.Size = parts[index]
			case "expire":
				index++
				if index == len(parts) {
					return nil, &errors.ParseError{Parser: "Table", Line: line}
				}
				data.Expire = parts[index]
			case "write-to":
				index++
				if index == len(parts) {
					return nil, &errors.ParseError{Parser: "Table", Line: line}
				}
				data.WriteTo = parts[index]
			case "nopurge":
				data.NoPurge = true
			case "store":
				index++
				if index == len(parts) {
					return nil, &errors.ParseError{Parser: "Table", Line: line}
				}
				if len(data.Store) == 0 {
					data.Store = parts[index]
				} else {
					data.Store = fmt.Sprintf("%s,%s", data.Store, parts[index])
				}
			default:
				return nil, &errors.ParseError{Parser: "Table", Line: line}
			}
			index++
		}
		return data, nil
	}
	return nil, &errors.ParseError{Parser: "Table", Line: line}
}

func (t *Table) Result() ([]common.ReturnResultLine, error) {
	dataLength := len(t.data)
	if dataLength == 0 {
		return nil, errors.ErrFetch
	}
	result := make([]common.ReturnResultLine, dataLength)
	for index, table := range t.data {

		var data strings.Builder
		data.WriteString("table ")
		data.WriteString(table.Name)
		data.WriteString(" type ")
		data.WriteString(table.Type)
		if table.TypeLen != 0 {
			data.WriteString(" len ")
			data.WriteString(strconv.FormatInt(table.TypeLen, 10))
		}
		if table.Size != "" {
			data.WriteString(" size ")
			data.WriteString(table.Size)
		}
		if table.Expire != "" {
			data.WriteString(" expire ")
			data.WriteString(table.Expire)
		}
		if table.WriteTo != "" {
			data.WriteString(" write-to ")
			data.WriteString(table.WriteTo)
		}
		if table.NoPurge {
			data.WriteString(" nopurge")
		}
		if table.Store != "" {
			data.WriteString(" store ")
			data.WriteString(table.Store)
		}

		result[index] = common.ReturnResultLine{
			Data:    data.String(),
			Comment: table.Comment,
		}
	}
	return result, nil
}
