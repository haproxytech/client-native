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

package stats

import (
	"fmt"
	"strings"

	"github.com/haproxytech/client-native/v5/config-parser/common"
)

type Auth struct {
	User     string
	Password string
	Comment  string
}

func (a *Auth) Parse(parts []string, comment string) error {
	if len(parts) < 3 {
		return fmt.Errorf("not enough params")
	}

	if comment != "" {
		a.Comment = comment
	}
	split := common.StringSplitIgnoreEmpty(strings.Join(parts[2:], " "), ':')
	if len(split) < 2 {
		return fmt.Errorf("wrong format for user & password")
	}
	a.User = split[0]
	a.Password = split[1]
	return nil
}

func (a *Auth) String() string {
	if a.User != "" && a.Password != "" {
		return fmt.Sprintf("auth %s:%s", a.User, a.Password)
	}
	return "auth"
}

func (a *Auth) GetComment() string {
	return a.Comment
}
