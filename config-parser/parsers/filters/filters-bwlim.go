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

type BandwidthLimit struct {
	Comment       string
	Attribute     string
	Name          string
	DefaultLimit  string
	DefaultPeriod string
	Limit         string
	Key           string
	Table         *string
	MinSize       *string
}

func (b *BandwidthLimit) Parse(parts []string, comment string) error {
	b.Comment = comment

	switch {
	case len(parts) < 7 || len(parts[1:])%2 == 1:
		return fmt.Errorf("missing required options")
	case len(parts) > 11:
		return fmt.Errorf("unsupported extra options")
	}

	for index, part := range parts {
		switch index {
		case 1:
			b.Attribute = part
		case 2:
			b.Name = part
		case 3, 5, 7, 9:
			if err := b.parsePart(parts, index); err != nil {
				return err
			}
		case 0, 4, 6, 8, 10: // keyword or values, can be ignored
			continue
		default:
			return fmt.Errorf("unexpected options: %s", strings.Join(parts[index:], " "))
		}
	}

	switch {
	case len(b.Limit) > 0 && len(b.Key) == 0:
		return fmt.Errorf("missing required key options")
	case len(b.Key) > 0 && len(b.Limit) == 0:
		return fmt.Errorf("missing required key limit")
	}

	return nil
}

func (b *BandwidthLimit) parsePart(parts []string, index int) error {
	part := parts[index]
	switch part {
	case "default-limit":
		b.DefaultLimit = parts[index+1]
	case "limit":
		b.Limit = parts[index+1]
	case "key":
		b.Key = parts[index+1]
	case "default-period":
		b.DefaultPeriod = parts[index+1]
	case "min-size":
		v := parts[index+1]
		b.MinSize = &v
	case "table":
		v := parts[index+1]
		b.Table = &v
	default:
		return fmt.Errorf("unsupported option: %s", part)
	}

	return nil
}

func (b *BandwidthLimit) Result() common.ReturnResultLine {
	var result strings.Builder

	result.WriteString("filter ")
	result.WriteString(b.Attribute)
	result.WriteString(" ")
	result.WriteString(b.Name)

	if len(b.Key) > 0 && len(b.Limit) > 0 {
		result.WriteString(" limit ")
		result.WriteString(b.Limit)
		result.WriteString(" key ")
		result.WriteString(b.Key)

		if table := b.Table; table != nil {
			result.WriteString(" table ")
			result.WriteString(*table)
		}
	} else {
		result.WriteString(" default-limit ")
		result.WriteString(b.DefaultLimit)
		result.WriteString(" default-period ")
		result.WriteString(b.DefaultPeriod)
	}

	if b.MinSize != nil {
		result.WriteString(" min-size ")
		result.WriteString(*b.MinSize)
	}

	return common.ReturnResultLine{
		Data:    result.String(),
		Comment: b.Comment,
	}
}
