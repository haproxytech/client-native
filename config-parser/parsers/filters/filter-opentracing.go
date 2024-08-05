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

type Opentracing struct {
	Comment string
	ID      string
	Config  string
}

func (o *Opentracing) Parse(parts []string, comment string) error {
	o.Comment = comment

	if len(parts) < 3 || len(parts)%2 == 1 {
		return fmt.Errorf("missing required options")
	}

	for index, part := range parts {
		switch index {
		case 0, 1:
			continue
		case 2, 4:
			switch part {
			case "id":
				o.ID = parts[index+1]
			case "config":
				o.Config = parts[index+1]
			default:
				return fmt.Errorf("unsupported option: %s", part)
			}
		case 3, 5: // values, can be ignored
			continue
		default:
			return fmt.Errorf("unexpected options: %s", strings.Join(parts[index:], " "))
		}
	}

	return nil
}

func (o *Opentracing) Result() common.ReturnResultLine {
	var result strings.Builder

	result.WriteString("filter opentracing")

	if o.ID != "" {
		result.WriteString(" id ")
		result.WriteString(o.ID)
	}

	if o.Config != "" {
		result.WriteString(" config ")
		result.WriteString(o.Config)
	}

	return common.ReturnResultLine{
		Data:    result.String(),
		Comment: o.Comment,
	}
}
