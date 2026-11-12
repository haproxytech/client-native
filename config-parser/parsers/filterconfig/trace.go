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

package filterconfig

import (
	"strconv"
	"strings"

	"github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/models"
)

type trace struct{}

func (trace) Parse(parts []string, config *models.FilterConfig) error {
	for index := 0; index < len(parts); index++ {
		switch parts[index] {
		case "quiet":
			config.TraceQuiet = true
		case "random-parsing":
			config.TraceRndParsing = true
		case "random-forwarding":
			config.TraceRndForwarding = true
		case "hexdump":
			config.TraceHexdump = true
		case "max-fwd":
			index++
			if index == len(parts) {
				return errors.ErrInvalidData
			}
			value, err := strconv.ParseInt(parts[index], 10, 64)
			if err != nil {
				return errors.ErrInvalidData
			}
			config.TraceMaxFwd = &value
		default:
			return errors.ErrInvalidData
		}
	}
	return nil
}

func (trace) Result(config models.FilterConfig) string {
	var result strings.Builder
	if config.TraceQuiet {
		result.WriteString(" quiet")
	}
	if config.TraceRndParsing {
		result.WriteString(" random-parsing")
	}
	if config.TraceRndForwarding {
		result.WriteString(" random-forwarding")
	}
	if config.TraceMaxFwd != nil {
		result.WriteString(" max-fwd ")
		result.WriteString(strconv.FormatInt(*config.TraceMaxFwd, 10))
	}
	if config.TraceHexdump {
		result.WriteString(" hexdump")
	}
	return result.String()
}
