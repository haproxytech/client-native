/*
Copyright 2026 HAProxy Technologies

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
	"slices"

	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

var cpuAffinityValues = []string{ //nolint:gochecknoglobals
	"per-core",
	"per-group",
	"auto",
	"per-thread",
	"per-ccx",
}

var cpuAffinityPerGroupArgs = []string{ //nolint:gochecknoglobals
	"auto",
	"loose",
}

type CPUAffinity struct {
	data        *types.CPUAffinity
	preComments []string // comments that appear before the actual line
}

func (c *CPUAffinity) Parse(line string, parts []string, comment string) (string, error) {
	if parts[0] != "cpu-affinity" {
		return "", &errors.ParseError{Parser: "CPUAffinity", Line: line}
	}
	if len(parts) < 2 {
		return "", &errors.ParseError{Parser: "CPUAffinity", Line: line, Message: "Parse error"}
	}
	affinity := parts[1]
	if !slices.Contains(cpuAffinityValues, affinity) {
		return "", &errors.ParseError{Parser: "CPUAffinity", Line: line, Message: fmt.Sprintf("invalid affinity %q (allowed: %v)", affinity, cpuAffinityValues)}
	}
	var argument string
	if len(parts) > 2 {
		if affinity != "per-group" {
			return "", &errors.ParseError{Parser: "CPUAffinity", Line: line, Message: fmt.Sprintf("affinity %q does not take an argument", affinity)}
		}
		if !slices.Contains(cpuAffinityPerGroupArgs, parts[2]) {
			return "", &errors.ParseError{Parser: "CPUAffinity", Line: line, Message: fmt.Sprintf("invalid per-group argument %q (allowed: %v)", parts[2], cpuAffinityPerGroupArgs)}
		}
		argument = parts[2]
	}
	c.data = &types.CPUAffinity{
		Affinity: affinity,
		Argument: argument,
		Comment:  comment,
	}
	return "", nil
}

func (c *CPUAffinity) Result() ([]common.ReturnResultLine, error) {
	if c.data == nil {
		return nil, errors.ErrFetch
	}
	data := "cpu-affinity " + c.data.Affinity
	if c.data.Argument != "" {
		data += " " + c.data.Argument
	}
	return []common.ReturnResultLine{
		{
			Data:    data,
			Comment: c.data.Comment,
		},
	}, nil
}
