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
	"fmt"
	"strings"

	"github.com/haproxytech/client-native/v6/config-parser/common"
)

type FcgiApp struct {
	Comment string
	Name    string
}

func (f *FcgiApp) Parse(parts []string, comment string) error {
	f.Comment = comment

	switch len(parts) {
	case 3:
		f.Name = parts[2]
	case 2:
		return fmt.Errorf("no FastCGI application name")
	default:
		return fmt.Errorf("unsupported extra options: %s", strings.Join(parts[2:], " "))
	}

	return nil
}

func (f *FcgiApp) Result() common.ReturnResultLine {
	var result strings.Builder
	result.WriteString("filter fcgi-app")

	if f.Name != "" {
		result.WriteString(" ")
		result.WriteString(f.Name)
	}

	return common.ReturnResultLine{
		Data:    result.String(),
		Comment: f.Comment,
	}
}
