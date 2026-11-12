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
	"strings"

	"github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/models"
)

type spoe struct{}

func (spoe) Parse(parts []string, config *models.FilterConfig) error {
	for index := 0; index < len(parts); index++ {
		switch parts[index] {
		case "engine":
			index++
			if index == len(parts) {
				return errors.ErrInvalidData
			}
			config.SpoeEngine = parts[index]
		case "config":
			index++
			if index == len(parts) {
				return errors.ErrInvalidData
			}
			config.SpoeConfig = parts[index]
		default:
			return errors.ErrInvalidData
		}
	}
	if config.SpoeConfig == "" {
		return errors.ErrInvalidData
	}
	return nil
}

func (spoe) Result(config models.FilterConfig) string {
	var result strings.Builder
	if config.SpoeEngine != "" {
		result.WriteString(" engine ")
		result.WriteString(config.SpoeEngine)
	}
	result.WriteString(" config ")
	result.WriteString(config.SpoeConfig)
	return result.String()
}
