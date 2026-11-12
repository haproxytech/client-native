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
	"github.com/haproxytech/client-native/v6/misc"
	"github.com/haproxytech/client-native/v6/models"
)

type FilterRule struct {
	data        []models.FilterRule
	preComments []string // comments that appear before the actual line
}

func (p *FilterRule) Parse(line string, parts []string, comment string) (string, error) {
	rule, err := parseFilterRule(parts, comment)
	if err != nil {
		return "", &errors.ParseError{Parser: "FilterRule", Line: line}
	}
	p.data = append(p.data, *rule)
	return "", nil
}

func parseFilterRule(parts []string, comment string) (*models.FilterRule, error) {
	if len(parts) < 2 {
		return nil, errors.ErrInvalidData
	}

	rule := &models.FilterRule{}
	switch parts[0] {
	case "filter-enable":
		rule.Type = "enable"
	case "filter-disable":
		rule.Type = "disable"
	default:
		return nil, errors.ErrInvalidData
	}

	for _, token := range parts[1:] {
		selector, err := parseFilterSelector(token)
		if err != nil {
			return nil, err
		}
		rule.Filters = append(rule.Filters, selector)
	}
	if comment != "" {
		rule.Metadata = misc.ParseMetadata(comment)
	}
	return rule, nil
}

func parseFilterSelector(value string) (*models.FilterSelector, error) {
	class, id, hasID := strings.Cut(value, "/")
	if !isFilterIdentifier(class) {
		return nil, errors.ErrInvalidData
	}
	if hasID && !isFilterIdentifier(id) {
		return nil, errors.ErrInvalidData
	}
	return &models.FilterSelector{Class: class, ID: id}, nil
}

func filterSelectorString(selector *models.FilterSelector) string {
	if selector == nil {
		return ""
	}
	if selector.ID == "" {
		return selector.Class
	}
	return selector.Class + "/" + selector.ID
}

func (p *FilterRule) Result() ([]common.ReturnResultLine, error) {
	if len(p.data) == 0 {
		return nil, errors.ErrFetch
	}
	result := make([]common.ReturnResultLine, len(p.data))
	for index, rule := range p.data {
		var line strings.Builder
		line.WriteString("filter-")
		line.WriteString(rule.Type)
		for _, selector := range rule.Filters {
			line.WriteString(" ")
			line.WriteString(filterSelectorString(selector))
		}

		comment := ""
		if rule.Metadata != nil {
			serializedComment, err := misc.SerializeMetadata(rule.Metadata)
			if err == nil {
				comment = serializedComment
			}
		}
		result[index] = common.ReturnResultLine{
			Data:    line.String(),
			Comment: comment,
		}
	}
	return result, nil
}
