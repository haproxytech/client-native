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

package actions

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

type SetStatus struct { // http-request redirect location <loc> [code <code>] [<option>] [<condition>]
	Status   string
	Reason   string
	Cond     string
	CondTest string
	Comment  string
}

func (f *SetStatus) Parse(parts []string, parserType types.ParserType, comment string) error {
	if comment != "" {
		f.Comment = comment
	}
	if len(parts) >= 3 {
		command, condition := common.SplitRequest(parts[2:])
		if len(command) < 1 {
			return errors.ErrInvalidData
		}
		f.Status = command[0]
		if value, err := strconv.Atoi(f.Status); err != nil {
			return fmt.Errorf("expecting an integer value, %w", err)
		} else if value < 100 || value > 999 {
			return fmt.Errorf(" expects an integer status code between 100 and 999, got %d", value)
		}
		index := 1

		if len(command) >= 3 && command[index] == "reason" {
			index++
			f.Reason = command[index]
		}
		if len(condition) > 1 {
			f.Cond = condition[0]
			f.CondTest = strings.Join(condition[1:], " ")
		}
		return nil
	}
	return fmt.Errorf("not enough params")
}

func (f *SetStatus) String() string {
	var result strings.Builder
	result.WriteString("set-status ")
	result.WriteString(f.Status)
	if f.Reason != "" {
		result.WriteString(" reason ")
		result.WriteString(f.Reason)
	}
	if f.Cond != "" {
		result.WriteString(" ")
		result.WriteString(f.Cond)
		result.WriteString(" ")
		result.WriteString(f.CondTest)
	}
	return result.String()
}

func (f *SetStatus) GetComment() string {
	return f.Comment
}
