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

type FilterConfig interface {
	GetFilterConfigs(parentType, parentName string, transactionID string) (int64, models.FilterConfigs, error)
	GetFilterConfig(id int64, parentType, parentName string, transactionID string) (int64, *models.FilterConfig, error)
	DeleteFilterConfig(id int64, parentType string, parentName string, transactionID string, version int64) error
	CreateFilterConfig(id int64, parentType string, parentName string, data *models.FilterConfig, transactionID string, version int64) error
	EditFilterConfig(id int64, parentType string, parentName string, data *models.FilterConfig, transactionID string, version int64) error
	ReplaceFilterConfigs(parentType string, parentName string, data models.FilterConfigs, transactionID string, version int64) error
}

func (c *client) GetFilterConfigs(parentType, parentName string, transactionID string) (int64, models.FilterConfigs, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	configs, err := ParseFilterConfigs(parentType, parentName, p)
	if err != nil {
		return v, nil, c.HandleError("", parentType, parentName, "", false, err)
	}

	return v, configs, nil
}

func (c *client) GetFilterConfig(id int64, parentType, parentName string, transactionID string) (int64, *models.FilterConfig, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	section, parentName, err := getParserFromParent(p, "filter-config", parentType, parentName)
	if err != nil {
		return v, nil, err
	}

	data, err := p.GetOne(section, parentName, "filter-config", int(id))
	if err != nil {
		return v, nil, c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, "", false, err)
	}

	config, ok := data.(models.FilterConfig)
	if !ok {
		return v, nil, misc.CreateTypeAssertError("filter-config")
	}

	return v, &config, nil
}

func (c *client) DeleteFilterConfig(id int64, parentType string, parentName string, transactionID string, version int64) error {
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	section, parentName, err := getParserFromParent(p, "filter-config", parentType, parentName)
	if err != nil {
		return err
	}

	if err := p.Delete(section, parentName, "filter-config", int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

func (c *client) CreateFilterConfig(id int64, parentType string, parentName string, data *models.FilterConfig, transactionID string, version int64) error {
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

	section, parentName, err := getParserFromParent(p, "filter-config", parentType, parentName)
	if err != nil {
		return err
	}

	if err := p.Insert(section, parentName, "filter-config", *data, int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

func (c *client) EditFilterConfig(id int64, parentType string, parentName string, data *models.FilterConfig, transactionID string, version int64) error {
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

	section, parentName, err := getParserFromParent(p, "filter-config", parentType, parentName)
	if err != nil {
		return err
	}

	if _, err := p.GetOne(section, parentName, "filter-config", int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	if err := p.Set(section, parentName, "filter-config", *data, int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

func (c *client) ReplaceFilterConfigs(parentType string, parentName string, data models.FilterConfigs, transactionID string, version int64) error {
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

	section, parentName, err := getParserFromParent(p, "filter-config", parentType, parentName)
	if err != nil {
		return err
	}

	configs, err := ParseFilterConfigs(parentType, parentName, p)
	if err != nil {
		return c.HandleError("", parentType, parentName, "", false, err)
	}

	for i := range configs {
		if err := p.Delete(section, parentName, "filter-config", 0); err != nil {
			return c.HandleError(strconv.FormatInt(int64(i), 10), parentType, parentName, t, transactionID == "", err)
		}
	}

	for i, config := range data {
		if err := p.Insert(section, parentName, "filter-config", *config, i); err != nil {
			return c.HandleError(strconv.FormatInt(int64(i), 10), parentType, parentName, t, transactionID == "", err)
		}
	}

	return c.SaveData(p, t, transactionID == "")
}

func ParseFilterConfigs(parentType, parentName string, p parser.Parser) (models.FilterConfigs, error) {
	section, parentName, err := getParserFromParent(p, "filter-config", parentType, parentName)
	if err != nil {
		return nil, err
	}

	data, err := p.Get(section, parentName, "filter-config", false)
	if err != nil {
		if errors.Is(err, parser_errors.ErrFetch) {
			return nil, nil
		}
		return nil, err
	}

	configs, ok := data.([]models.FilterConfig)
	if !ok {
		return nil, misc.CreateTypeAssertError("filter-config")
	}
	result := make(models.FilterConfigs, 0, len(configs))
	for i := range configs {
		config := configs[i]
		result = append(result, &config)
	}
	return result, nil
}
