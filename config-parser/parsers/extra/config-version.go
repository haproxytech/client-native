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

package extra

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

type ConfigVersion struct {
	Name        string
	data        *types.ConfigVersion
	preComments []string // comments that appear before the actual line
}

func (p *ConfigVersion) Init() {
	p.Name = "# _version"
	p.data = nil
}

func (p *ConfigVersion) Get(createIfNotExist bool) (common.ParserData, error) {
	if p.data != nil {
		return p.data, nil
	} else if createIfNotExist {
		p.data = &types.ConfigVersion{
			Value: 1,
		}
		return p.data, nil
	}
	return nil, fmt.Errorf("no data")
}

// Parse see if we have version, since it is not haproxy keyword, it's in comments
func (p *ConfigVersion) Parse(line string, parts []string, comment string) (string, error) {
	if strings.HasPrefix(comment, "_version") {
		data := common.StringSplitIgnoreEmpty(comment, '=')
		if len(data) < 2 {
			return "", &errors.ParseError{Parser: "ConfigVersion", Line: line}
		}
		if version, err := strconv.ParseInt(data[1], 10, 64); err == nil {
			p.data = &types.ConfigVersion{
				Value: version,
			}
		}
		return "", nil
	}
	return "", &errors.ParseError{Parser: "ConfigVersion", Line: line}
}

func (p *ConfigVersion) Result() ([]common.ReturnResultLine, error) {
	if p.data == nil {
		return nil, errors.ErrFetch
	}

	return []common.ReturnResultLine{
		{
			Data:    fmt.Sprintf("# _version=%d", p.data.Value),
			Comment: "",
		},
	}, nil
}
