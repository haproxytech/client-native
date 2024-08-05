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

type Capture struct { // http-request capture sample [ len <length> | id <id> ] [ { if | unless } <condition> ]
	Sample   string
	Len      *int64 // Has to be > 0.
	SlotID   *int64
	Cond     string
	CondTest string
	Comment  string
}

func (f *Capture) Parse(parts []string, parserType types.ParserType, comment string) error {
	if comment != "" {
		f.Comment = comment
	}
	if len(parts) >= 5 {
		command, condition := common.SplitRequest(parts[2:])
		if len(command) < 3 {
			return errors.ErrInvalidData
		}
		i, err := strconv.ParseInt(command[2], 10, 64)
		if err != nil {
			return errors.ErrInvalidData
		}
		f.Sample = command[0]
		switch command[1] {
		case "len":
			// Response only takes id.
			if parts[0] == "http-response" {
				return errors.ErrInvalidData
			}
			f.Len = &i
		case "id":
			f.SlotID = &i
		default:
			return errors.ErrInvalidData
		}
		if len(condition) > 1 {
			f.Cond = condition[0]
			f.CondTest = strings.Join(condition[1:], " ")
		}
		return nil
	}
	return fmt.Errorf("not enough params")
}

func (f *Capture) String() string {
	var result strings.Builder
	result.WriteString("capture ")
	result.WriteString(f.Sample)
	result.WriteString(" ")
	if f.SlotID != nil {
		result.WriteString("id")
		result.WriteString(" ")
		result.WriteString(strconv.FormatInt(*f.SlotID, 10))
	} else if f.Len != nil {
		result.WriteString("len")
		result.WriteString(" ")
		result.WriteString(strconv.FormatInt(*f.Len, 10))
	}
	if f.Cond != "" {
		result.WriteString(" ")
		result.WriteString(f.Cond)
		result.WriteString(" ")
		result.WriteString(f.CondTest)
	}
	return result.String()
}

func (f *Capture) GetComment() string {
	return f.Comment
}
