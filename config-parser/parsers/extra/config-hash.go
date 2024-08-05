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
	"strings"

	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

type ConfigHash struct {
	Name string
	// Mode string
	data        *types.ConfigHash
	preComments []string // comments that appear before the actual line
}

func (h *ConfigHash) Init() {
	h.Name = "# _md5hash"
	h.data = nil
}

func (h *ConfigHash) Get(createIfNotExist bool) (common.ParserData, error) {
	if h.data != nil {
		return h.data, nil
	} else if createIfNotExist {
		h.data = &types.ConfigHash{
			Value: "",
		}
		return h.data, nil
	}
	return nil, fmt.Errorf("no data")
}

// Parse see if we have version, since it is not haproxy keyword, it's in comments
func (h *ConfigHash) Parse(line string, parts []string, comment string) (string, error) {
	if strings.HasPrefix(comment, "_md5hash") {
		data := common.StringSplitIgnoreEmpty(comment, '=')
		if len(data) < 2 {
			return "", &errors.ParseError{Parser: "ConfigHash", Line: line}
		}
		h.data = &types.ConfigHash{
			Value: data[1],
		}

		return "", nil
	}
	return "", &errors.ParseError{Parser: "ConfigHash", Line: line}
}

func (h *ConfigHash) Result() ([]common.ReturnResultLine, error) {
	return nil, errors.ErrFetch
}
