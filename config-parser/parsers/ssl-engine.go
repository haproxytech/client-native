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

type SslEngine struct {
	data        []types.SslEngine
	preComments []string
}

func (s *SslEngine) parse(line string, parts []string, comment string) (*types.SslEngine, error) {
	if len(parts) < 2 || len(parts) > 3 || parts[0] != "ssl-engine" {
		return nil, &errors.ParseError{Parser: "ssl-engine", Line: line}
	}
	sslEngine := &types.SslEngine{
		Name: parts[1],
		Algorithms: func() []string {
			var res []string
			if len(parts) == 3 {
				res = strings.Split(parts[2], ",")
			}
			return res
		}(),
		Comment: comment,
	}
	return sslEngine, nil
}

func (s *SslEngine) Result() ([]common.ReturnResultLine, error) {
	if len(s.data) == 0 {
		return nil, errors.ErrFetch
	}

	result := make([]common.ReturnResultLine, len(s.data))
	for index, data := range s.data {

		var sb strings.Builder
		sb.WriteString("ssl-engine")
		sb.WriteString(" ")
		sb.WriteString(data.Name)

		if len(data.Algorithms) > 0 {
			sb.WriteString(" ")
			sb.WriteString(strings.Join(data.Algorithms, ","))
		}

		result[index] = common.ReturnResultLine{
			Data:    sb.String(),
			Comment: data.Comment,
		}
	}
	return result, nil
}
