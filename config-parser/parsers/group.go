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
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

type Group struct {
	data        []types.Group
	preComments []string // comments that appear before the actual line
}

func (l *Group) parse(line string, parts []string, comment string) (*types.Group, error) {
	if len(parts) >= 2 {
		group := &types.Group{
			Name:    parts[1],
			Comment: comment,
		}
		if len(parts) > 3 && parts[2] == "users" {
			group.Users = common.StringSplitIgnoreEmpty(parts[3], ',')
		}
		return group, nil
	}
	return nil, &errors.ParseError{Parser: "Group", Line: line}
}

func (l *Group) Result() ([]common.ReturnResultLine, error) {
	if len(l.data) == 0 {
		return nil, errors.ErrFetch
	}
	result := make([]common.ReturnResultLine, len(l.data))
	for index, group := range l.data {
		users := ""
		if len(group.Users) > 0 {
			var s strings.Builder
			s.WriteString(" users ")
			first := true
			for _, user := range group.Users {
				if !first {
					s.WriteString(",")
				} else {
					first = false
				}
				s.WriteString(user)
			}
			users = s.String()
		}
		result[index] = common.ReturnResultLine{
			Data:    fmt.Sprintf("group %s%s", group.Name, users),
			Comment: group.Comment,
		}
	}
	return result, nil
}
