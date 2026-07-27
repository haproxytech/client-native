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

package filters

import (
	"strings"

	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/errors"
)

// Lua represents a Lua filter declared with "filter lua.<name> [args...]".
// The filter keyword is the name registered from Lua via core.register_filter().
type Lua struct {
	Name    string
	Args    []string
	Comment string
}

func (f *Lua) Parse(parts []string, comment string) error {
	if comment != "" {
		f.Comment = comment
	}
	if len(parts) < 2 {
		return errors.ErrInvalidData
	}
	name := strings.TrimPrefix(parts[1], "lua.")
	if name == "" || name == parts[1] {
		return errors.ErrInvalidData
	}
	f.Name = name
	if len(parts) > 2 {
		f.Args = append([]string{}, parts[2:]...)
	}
	return nil
}

func (f *Lua) Result() common.ReturnResultLine {
	var result strings.Builder
	result.WriteString("filter lua.")
	result.WriteString(f.Name)
	for _, arg := range f.Args {
		result.WriteString(" ")
		result.WriteString(arg)
	}
	return common.ReturnResultLine{
		Data:    result.String(),
		Comment: f.Comment,
	}
}
