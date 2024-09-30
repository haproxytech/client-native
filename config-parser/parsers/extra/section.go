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

	"github.com/haproxytech/client-native/v5/config-parser/common"
	"github.com/haproxytech/client-native/v5/config-parser/errors"
	"github.com/haproxytech/client-native/v5/config-parser/types"
)

type Section struct {
	Name        string
	data        *types.Section
	preComments []string // comments that appear before the actual line
}

func (s *Section) Init() {
	s.data = &types.Section{}
}

// Parse see if we have section name
func (s *Section) Parse(line string, parts []string, comment string) (string, error) {
	if parts[0] == s.Name {
		if len(parts) > 1 {
			s.data.Name = parts[1]
		}
		if len(parts) > 3 && parts[2] == "from" {
			s.data.FromDefaults = parts[3]
		}
		return s.Name, nil
	}
	return "", &errors.ParseError{Parser: "Section", Line: line}
}

func (s *Section) Result() ([]common.ReturnResultLine, error) {
	return nil, fmt.Errorf("not valid")
}
