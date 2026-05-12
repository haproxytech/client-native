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
	"strconv"
	"strings"

	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/errors"
)

type Trace struct { // filter trace [name <name>] [random-parsing] [random-forwarding] [max-fwd <max>] [hexdump]
	Name             string
	RandomParsing    bool
	RandomForwarding bool
	MaxFwd           *int64
	Hexdump          bool
	Comment          string
}

func (f *Trace) Parse(parts []string, comment string) error {
	// we have filter trace [name <name>] [random-parsing] [random-forwarding] [max-fwd <max>] [hexdump]
	if comment != "" {
		f.Comment = comment
	}
	if len(parts) < 2 {
		return errors.ErrInvalidData
	}
	index := 2
	for index < len(parts) {
		switch parts[index] {
		case "name":
			index++
			if index == len(parts) {
				return errors.ErrInvalidData
			}
			f.Name = parts[index]
		case "random-parsing":
			f.RandomParsing = true
		case "random-forwarding":
			f.RandomForwarding = true
		case "max-fwd":
			index++
			if index == len(parts) {
				return errors.ErrInvalidData
			}
			n, err := strconv.ParseInt(parts[index], 10, 64)
			if err != nil {
				return errors.ErrInvalidData
			}
			f.MaxFwd = &n
		case "hexdump":
			f.Hexdump = true
		default:
			return errors.ErrInvalidData
		}
		index++
	}
	return nil
}

func (f *Trace) Result() common.ReturnResultLine {
	var result strings.Builder
	result.WriteString("filter trace")
	if f.Name != "" {
		result.WriteString(" name ")
		result.WriteString(f.Name)
	}
	if f.RandomParsing {
		result.WriteString(" random-parsing")
	}
	if f.RandomForwarding {
		result.WriteString(" random-forwarding")
	}
	if f.MaxFwd != nil {
		result.WriteString(" max-fwd ")
		result.WriteString(strconv.FormatInt(*f.MaxFwd, 10))
	}
	if f.Hexdump {
		result.WriteString(" hexdump")
	}
	return common.ReturnResultLine{
		Data:    result.String(),
		Comment: f.Comment,
	}
}
