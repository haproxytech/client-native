// Copyright 2019 HAProxy Technologies
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
	"fmt"

	strfmt "github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/config-parser/v5"

	"github.com/haproxytech/client-native/v6/models"
)

type Defaults interface {
	GetDefaultsConfiguration(transactionID string) (int64, *models.Defaults, error)
	PushDefaultsConfiguration(data *models.Defaults, transactionID string, version int64) error
	GetDefaultsSections(transactionID string) (int64, models.DefaultsSections, error)
	GetDefaultsSection(name string, transactionID string) (int64, *models.Defaults, error)
	DeleteDefaultsSection(name string, transactionID string, version int64) error
	EditDefaultsSection(name string, data *models.Defaults, transactionID string, version int64) error
	CreateDefaultsSection(data *models.Defaults, transactionID string, version int64) error
}

// GetDefaultsConfiguration returns configuration version and a
// struct representing Defaults configuration
func (c *client) GetDefaultsConfiguration(transactionID string) (int64, *models.Defaults, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	d := &models.Defaults{}
	_ = ParseSection(&d.DefaultsBase, parser.Defaults, parser.DefaultSectionName, p)

	return v, d, nil
}

// PushDefaultsConfiguration pushes a Defaults config struct to global
// config file
func (c *client) PushDefaultsConfiguration(data *models.Defaults, transactionID string, version int64) error {
	if c.UseModelsValidation {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return NewConfError(ErrValidationError, validationErr.Error())
		}
	}

	return c.editSection(parser.Defaults, parser.DefaultSectionName, &data.DefaultsBase, transactionID, version)
}

// GetDefaultsSections returns configuration version and an array of
// configured defaults sections. Returns error on fail.
func (c *client) GetDefaultsSections(transactionID string) (int64, models.DefaultsSections, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	dNames, err := p.SectionsGet(parser.Defaults)
	if err != nil {
		return v, nil, err
	}

	defaults := []*models.Defaults{}
	for _, name := range dNames {
		d := &models.Defaults{DefaultsBase: models.DefaultsBase{Name: name}}
		if err := ParseSection(&d.DefaultsBase, parser.Defaults, name, p); err != nil {
			continue
		}
		defaults = append(defaults, d)
	}

	return v, defaults, nil
}

// GetDefaultsSection returns configuration version and a requested defaults section.
// Returns error on fail or if defaults section does not exist.
func (c *client) GetDefaultsSection(name string, transactionID string) (int64, *models.Defaults, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	if !c.checkSectionExists(parser.Defaults, name, p) {
		return v, nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Defaults section %s does not exist", name))
	}

	defaults := &models.Defaults{DefaultsBase: models.DefaultsBase{Name: name}}
	if err := ParseSection(&defaults.DefaultsBase, parser.Defaults, name, p); err != nil {
		return v, nil, err
	}

	return v, defaults, nil
}

// DeleteDefaultsSection deletes a defaults section in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) DeleteDefaultsSection(name string, transactionID string, version int64) error {
	return c.deleteSection(parser.Defaults, name, transactionID, version)
}

// EditDefaultsSection edits a defaults section in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) EditDefaultsSection(name string, data *models.Defaults, transactionID string, version int64) error {
	if c.UseModelsValidation {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return NewConfError(ErrValidationError, validationErr.Error())
		}
	}

	return c.editSection(parser.Defaults, name, &data.DefaultsBase, transactionID, version)
}

// CreateDefaultsSection creates a defaults section in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) CreateDefaultsSection(data *models.Defaults, transactionID string, version int64) error {
	if c.UseModelsValidation {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return NewConfError(ErrValidationError, validationErr.Error())
		}
	}

	return c.createSection(parser.Defaults, data.Name, &data.DefaultsBase, transactionID, version)
}
