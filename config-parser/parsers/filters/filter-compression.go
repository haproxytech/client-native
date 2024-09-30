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

	"github.com/haproxytech/client-native/v5/config-parser/common"
)

type Compression struct {
	Enabled bool
	Comment string
}

func (f *Compression) Parse(parts []string, comment string) error {
	if comment != "" {
		f.Comment = comment
	}
	if len(parts) > 2 {
		return fmt.Errorf("unexpected extra args: %s", strings.Join(parts[2:], " "))
	}
	f.Enabled = true
	return nil
}

func (f *Compression) Result() common.ReturnResultLine {
	var result strings.Builder
	result.WriteString("filter compression")
	return common.ReturnResultLine{
		Data:    result.String(),
		Comment: f.Comment,
	}
}
