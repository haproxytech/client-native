// Copyright 2026 HAProxy Technologies
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package configuration

import (
	"errors"
	"strconv"

	"github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/client-native/v6/config-parser"
	parser_errors "github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/misc"
	"github.com/haproxytech/client-native/v6/models"
)

type FilterRule interface {
	GetFilterRules(parentType, parentName string, transactionID string) (int64, models.FilterRules, error)
	GetFilterRule(id int64, parentType, parentName string, transactionID string) (int64, *models.FilterRule, error)
	DeleteFilterRule(id int64, parentType string, parentName string, transactionID string, version int64) error
	CreateFilterRule(id int64, parentType string, parentName string, data *models.FilterRule, transactionID string, version int64) error
	EditFilterRule(id int64, parentType string, parentName string, data *models.FilterRule, transactionID string, version int64) error
	ReplaceFilterRules(parentType string, parentName string, data models.FilterRules, transactionID string, version int64) error
}

func (c *client) GetFilterRules(parentType, parentName string, transactionID string) (int64, models.FilterRules, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	rules, err := ParseFilterRules(parentType, parentName, p)
	if err != nil {
		return v, nil, c.HandleError("", parentType, parentName, "", false, err)
	}

	return v, rules, nil
}

func (c *client) GetFilterRule(id int64, parentType, parentName string, transactionID string) (int64, *models.FilterRule, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	section, parentName, err := getParserFromParent(p, "filter-rule", parentType, parentName)
	if err != nil {
		return v, nil, err
	}

	data, err := p.GetOne(section, parentName, "filter-rule", int(id))
	if err != nil {
		return v, nil, c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, "", false, err)
	}

	rule, ok := data.(models.FilterRule)
	if !ok {
		return v, nil, misc.CreateTypeAssertError("filter-rule")
	}

	return v, &rule, nil
}

func (c *client) DeleteFilterRule(id int64, parentType string, parentName string, transactionID string, version int64) error {
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	section, parentName, err := getParserFromParent(p, "filter-rule", parentType, parentName)
	if err != nil {
		return err
	}

	if err := p.Delete(section, parentName, "filter-rule", int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

func (c *client) CreateFilterRule(id int64, parentType string, parentName string, data *models.FilterRule, transactionID string, version int64) error {
	if c.UseModelsValidation {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return NewConfError(ErrValidationError, validationErr.Error())
		}
	}

	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	section, parentName, err := getParserFromParent(p, "filter-rule", parentType, parentName)
	if err != nil {
		return err
	}

	if err := p.Insert(section, parentName, "filter-rule", *data, int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

func (c *client) EditFilterRule(id int64, parentType string, parentName string, data *models.FilterRule, transactionID string, version int64) error {
	if c.UseModelsValidation {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return NewConfError(ErrValidationError, validationErr.Error())
		}
	}

	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	section, parentName, err := getParserFromParent(p, "filter-rule", parentType, parentName)
	if err != nil {
		return err
	}

	if _, err := p.GetOne(section, parentName, "filter-rule", int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	if err := p.Set(section, parentName, "filter-rule", *data, int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

func (c *client) ReplaceFilterRules(parentType string, parentName string, data models.FilterRules, transactionID string, version int64) error {
	if c.UseModelsValidation {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return NewConfError(ErrValidationError, validationErr.Error())
		}
	}

	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	section, parentName, err := getParserFromParent(p, "filter-rule", parentType, parentName)
	if err != nil {
		return err
	}

	rules, err := ParseFilterRules(parentType, parentName, p)
	if err != nil {
		return c.HandleError("", parentType, parentName, "", false, err)
	}

	for i := range rules {
		if err := p.Delete(section, parentName, "filter-rule", 0); err != nil {
			return c.HandleError(strconv.FormatInt(int64(i), 10), parentType, parentName, t, transactionID == "", err)
		}
	}

	for i, rule := range data {
		if err := p.Insert(section, parentName, "filter-rule", *rule, i); err != nil {
			return c.HandleError(strconv.FormatInt(int64(i), 10), parentType, parentName, t, transactionID == "", err)
		}
	}

	return c.SaveData(p, t, transactionID == "")
}

func ParseFilterRules(parentType, parentName string, p parser.Parser) (models.FilterRules, error) {
	section, parentName, err := getParserFromParent(p, "filter-rule", parentType, parentName)
	if err != nil {
		return nil, err
	}

	data, err := p.Get(section, parentName, "filter-rule", false)
	if err != nil {
		if errors.Is(err, parser_errors.ErrFetch) {
			return nil, nil
		}
		return nil, err
	}

	rules, ok := data.([]models.FilterRule)
	if !ok {
		return nil, misc.CreateTypeAssertError("filter-rule")
	}
	result := make(models.FilterRules, 0, len(rules))
	for i := range rules {
		rule := rules[i]
		result = append(result, &rule)
	}
	return result, nil
}
