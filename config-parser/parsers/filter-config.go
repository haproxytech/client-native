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
	"strings"

	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/errors"
	filter_config "github.com/haproxytech/client-native/v6/config-parser/parsers/filterconfig"
	"github.com/haproxytech/client-native/v6/misc"
	"github.com/haproxytech/client-native/v6/models"
)

type FilterConfig struct {
	data        []models.FilterConfig
	preComments []string // comments that appear before the actual line
}

func (p *FilterConfig) Parse(line string, parts []string, comment string) (string, error) {
	config, err := p.parse(parts, comment)
	if err != nil {
		return "", &errors.ParseError{Parser: "FilterConfig", Line: line}
	}
	p.data = append(p.data, *config)
	return "", nil
}

func (p *FilterConfig) parse(parts []string, comment string) (*models.FilterConfig, error) {
	if len(parts) < 2 || parts[0] != "filter-config" {
		return nil, errors.ErrInvalidData
	}

	class, ok := filter_config.Get(parts[1])
	if !ok {
		return nil, errors.ErrInvalidData
	}

	config := &models.FilterConfig{Class: parts[1]}
	index := 2
	if filter_config.HasID(config.Class) {
		if index+1 >= len(parts) || parts[index] != "id" || !isFilterIdentifier(parts[index+1]) {
			return nil, errors.ErrInvalidData
		}
		config.ID = parts[index+1]
		index += 2
	}

	if index < len(parts) && parts[index] == "enabled" {
		config.Enabled = true
		index++
	}

	if err := class.Parse(parts[index:], config); err != nil {
		return nil, err
	}

	if comment != "" {
		config.Metadata = misc.ParseMetadata(comment)
	}
	return config, nil
}

func (p *FilterConfig) Result() ([]common.ReturnResultLine, error) {
	if len(p.data) == 0 {
		return nil, errors.ErrFetch
	}
	result := make([]common.ReturnResultLine, len(p.data))
	for index, config := range p.data {
		result[index] = filterConfigResult(config)
	}
	return result, nil
}

func filterConfigResult(config models.FilterConfig) common.ReturnResultLine {
	var result strings.Builder
	result.WriteString("filter-config ")
	result.WriteString(config.Class)
	if config.ID != "" {
		result.WriteString(" id ")
		result.WriteString(config.ID)
	}
	if config.Enabled {
		result.WriteString(" enabled")
	}

	if class, ok := filter_config.Get(config.Class); ok {
		result.WriteString(class.Result(config))
	}

	comment := ""
	if config.Metadata != nil {
		serializedComment, err := misc.SerializeMetadata(config.Metadata)
		if err == nil {
			comment = serializedComment
		}
	}
	return common.ReturnResultLine{
		Data:    result.String(),
		Comment: comment,
	}
}

func isFilterIdentifier(value string) bool {
	if value == "" {
		return false
	}
	for _, r := range value {
		switch {
		case r >= 'a' && r <= 'z':
		case r >= 'A' && r <= 'Z':
		case r >= '0' && r <= '9':
		case r == '.', r == '_', r == '-':
		default:
			return false
		}
	}
	return true
}
