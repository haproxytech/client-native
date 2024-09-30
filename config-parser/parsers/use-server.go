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
	"github.com/haproxytech/client-native/v5/config-parser/types"
)

type UseServer struct {
	data        []types.UseServer
	preComments []string // comments that appear before the actual line
}

func (l *UseServer) Parse(line string, parts []string, comment string) (string, error) {
	if len(parts) > 3 && parts[0] == "use-server" {
		data := types.UseServer{
			Name:     parts[1],
			CondTest: strings.Join(parts[3:], " "),
			Comment:  comment,
		}
		switch parts[2] {
		case "if", "unless":
			data.Cond = parts[2]
		default:
			return "", &errors.ParseError{Parser: "UseServer", Line: line}
		}
		l.data = append(l.data, data)
		return "", nil
	}
	return "", &errors.ParseError{Parser: "UseServer", Line: line}
}

func (l *UseServer) Result() ([]common.ReturnResultLine, error) {
	if len(l.data) == 0 {
		return nil, errors.ErrFetch
	}
	// use-server
	result := make([]common.ReturnResultLine, len(l.data))
	for index, data := range l.data {
		result[index] = common.ReturnResultLine{
			Data:    fmt.Sprintf("use-server %s %s %s", data.Name, data.Cond, data.CondTest),
			Comment: data.Comment,
		}
	}
	return result, nil
}
