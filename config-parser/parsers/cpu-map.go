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

	"github.com/haproxytech/client-native/v5/config-parser/common"
	"github.com/haproxytech/client-native/v5/config-parser/errors"
	"github.com/haproxytech/client-native/v5/config-parser/types"
)

type CPUMap struct {
	data        []types.CPUMap
	preComments []string // comments that appear before the actual line
}

func (c *CPUMap) parse(line string, parts []string, comment string) (*types.CPUMap, error) {
	if len(parts) < 3 {
		return nil, &errors.ParseError{Parser: "CPUMap", Line: line, Message: "Parse error"}
	}
	cpuMap := &types.CPUMap{
		Process: parts[1],
		CPUSet:  strings.Join(parts[2:], " "),
		Comment: comment,
	}
	return cpuMap, nil
}

func (c *CPUMap) Result() ([]common.ReturnResultLine, error) {
	if len(c.data) == 0 {
		return nil, errors.ErrFetch
	}
	result := make([]common.ReturnResultLine, len(c.data))
	for index, cpuMap := range c.data {
		result[index] = common.ReturnResultLine{
			Data:    fmt.Sprintf("cpu-map %s %s", cpuMap.Process, cpuMap.CPUSet),
			Comment: cpuMap.Comment,
		}
	}
	return result, nil
}

func (c *CPUMap) Equal(b *CPUMap) bool {
	if b == nil {
		return false
	}
	if b.data == nil {
		return false
	}
	if len(c.data) != len(b.data) {
		return false
	}
	for _, cCPUMap := range c.data {
		found := false
		for _, bCPUMap := range b.data {
			if cCPUMap.Process == bCPUMap.Process {
				if cCPUMap.CPUSet != bCPUMap.CPUSet {
					return false
				}
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
