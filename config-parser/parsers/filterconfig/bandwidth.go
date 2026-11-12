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
	"github.com/haproxytech/client-native/v6/misc"
	"github.com/haproxytech/client-native/v6/models"
)

type bandwidth struct{}

func (bandwidth) Parse(parts []string, config *models.FilterConfig) error {
	var hasLimit, hasKey, hasDefaultLimit, hasDefaultPeriod bool
	for index := 0; index < len(parts); index++ {
		option := parts[index]
		index++
		if index == len(parts) {
			return errors.ErrInvalidData
		}
		value := parts[index]
		switch option {
		case "limit":
			limit := misc.ParseSize(value)
			if limit == nil {
				return errors.ErrInvalidData
			}
			config.Limit = *limit
			hasLimit = true
		case "key":
			config.Key = value
			hasKey = true
		case "default-limit":
			limit := misc.ParseSize(value)
			if limit == nil {
				return errors.ErrInvalidData
			}
			config.DefaultLimit = *limit
			hasDefaultLimit = true
		case "default-period":
			period := misc.ParseTimeout(value)
			if period == nil {
				return errors.ErrInvalidData
			}
			config.DefaultPeriod = *period
			hasDefaultPeriod = true
		case "table":
			config.Table = value
		case "min-size":
			size := misc.ParseSize(value)
			if size == nil {
				return errors.ErrInvalidData
			}
			config.MinSize = *size
		default:
			return errors.ErrInvalidData
		}
	}

	shared := hasLimit && hasKey
	perStream := hasDefaultLimit && hasDefaultPeriod
	if hasLimit != hasKey || hasDefaultLimit != hasDefaultPeriod {
		return errors.ErrInvalidData
	}
	if shared == perStream {
		return errors.ErrInvalidData
	}
	if config.Table != "" && !shared {
		return errors.ErrInvalidData
	}
	return nil
}

func (bandwidth) Result(config models.FilterConfig) string {
	var result strings.Builder
	if config.Key != "" {
		result.WriteString(" limit ")
		result.WriteString(misc.SerializeSize(config.Limit))
		result.WriteString(" key ")
		result.WriteString(config.Key)
		if config.Table != "" {
			result.WriteString(" table ")
			result.WriteString(config.Table)
		}
	} else {
		result.WriteString(" default-limit ")
		result.WriteString(misc.SerializeSize(config.DefaultLimit))
		result.WriteString(" default-period ")
		result.WriteString(misc.SerializeTime(config.DefaultPeriod, string(misc.Nearest)))
	}
	if config.MinSize != 0 {
		result.WriteString(" min-size ")
		result.WriteString(misc.SerializeSize(config.MinSize))
	}
	return result.String()
}
