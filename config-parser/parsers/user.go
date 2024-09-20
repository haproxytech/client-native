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

type User struct {
	data        []types.User
	preComments []string // comments that appear before the actual line
}

func (l *User) parse(line string, parts []string, comment string) (*types.User, error) {
	if len(parts) >= 2 {
		user := types.User{
			Name:    parts[1],
			Comment: comment,
		}
		// see if we have password
		index := 3
		if len(parts) > index {
			if parts[2] == "password" {
				user.Password = parts[3]
				index += 2
			}
			if parts[2] == "insecure-password" {
				user.Password = parts[3]
				user.IsInsecure = true
				index += 2
			}
		}
		if len(parts) > index {
			user.Groups = common.StringSplitIgnoreEmpty(parts[index], ',')
		}
		return &user, nil
	}
	return nil, &errors.ParseError{Parser: "User", Line: line}
}

func (l *User) Result() ([]common.ReturnResultLine, error) {
	if len(l.data) == 0 {
		return nil, errors.ErrFetch
	}
	result := make([]common.ReturnResultLine, len(l.data))
	for index, user := range l.data {
		pwd := ""
		if user.Password != "" {
			if user.IsInsecure {
				pwd = " insecure-password " + user.Password
			} else {
				pwd = " password " + user.Password
			}
		}
		groups := ""
		if len(user.Groups) > 0 {
			var s strings.Builder
			s.WriteString(" groups ")
			first := true
			for _, user := range user.Groups {
				if !first {
					s.WriteString(",")
				} else {
					first = false
				}
				s.WriteString(user)
			}
			groups = s.String()
		}
		result[index] = common.ReturnResultLine{
			Data:    fmt.Sprintf("user %s%s%s", user.Name, pwd, groups),
			Comment: user.Comment,
		}
	}
	return result, nil
}
