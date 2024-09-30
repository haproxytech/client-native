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

type Admin struct {
	Cond     string
	CondTest string
	Comment  string
}

func (a *Admin) Parse(parts []string, comment string) error {
	if len(parts) < 4 {
		return fmt.Errorf("not enough params")
	}

	if comment != "" {
		a.Comment = comment
	}
	_, condition := common.SplitRequest(parts[2:])
	if len(condition) > 1 {
		a.Cond = condition[0]
		a.CondTest = strings.Join(condition[1:], " ")
	}
	return nil
}

func (a *Admin) String() string {
	if a.Cond != "" {
		return fmt.Sprintf("admin %s %s", a.Cond, a.CondTest)
	}
	return "admin"
}

func (a *Admin) GetComment() string {
	return a.Comment
}
