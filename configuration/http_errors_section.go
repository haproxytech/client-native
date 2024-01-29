// Copyright 2022 HAProxy Technologies
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

package configuration

import (
	"fmt"
	"strconv"

	"github.com/go-openapi/strfmt"
	"github.com/haproxytech/client-native/v6/models"
	parser "github.com/haproxytech/config-parser/v5"
	"github.com/haproxytech/config-parser/v5/types"
)

type HTTPErrorsSection interface {
	GetHTTPErrorsSections(transactionID string) (int64, models.HTTPErrorsSections, error)
	GetHTTPErrorsSection(name string, transactionID string) (int64, *models.HTTPErrorsSection, error)
	DeleteHTTPErrorsSection(name string, transactionID string, version int64) error
	CreateHTTPErrorsSection(data *models.HTTPErrorsSection, transactionID string, version int64) error
	EditHTTPErrorsSection(name string, data *models.HTTPErrorsSection, transactionID string, version int64) error
}

// GetHTTPErrorsSections returns all http-error sections in the configuration.
func (c *client) GetHTTPErrorsSections(transactionID string) (int64, models.HTTPErrorsSections, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	names, err := p.SectionsGet(parser.HTTPErrors)
	if err != nil {
		return v, nil, err
	}

	sections := make(models.HTTPErrorsSections, 0, len(names))

	for _, name := range names {
		section, parseErr := ParseHTTPErrorsSection(p, name)
		if parseErr != nil {
			continue
		}
		sections = append(sections, section)
	}

	return v, sections, nil
}

// GetHTTPErrorsSection returns a single http-error section with a given name.
func (c *client) GetHTTPErrorsSection(name, transactionID string) (int64, *models.HTTPErrorsSection, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	if !c.checkSectionExists(parser.HTTPErrors, name, p) {
		return v, nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("HTTPErrorsSection %s does not exist", name))
	}

	section, parseErr := ParseHTTPErrorsSection(p, name)
	if parseErr != nil {
		return 0, nil, parseErr
	}

	return v, section, nil
}

// DeleteHTTPErrorsSection deletes a single http-error section with a given name.
func (c *client) DeleteHTTPErrorsSection(name, transactionID string, version int64) error {
	return c.deleteSection(parser.HTTPErrors, name, transactionID, version)
}

// CreateHTTPErrorsSection adds a new http-errors section.
func (c *client) CreateHTTPErrorsSection(data *models.HTTPErrorsSection, transactionID string, version int64) error {
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

	if c.checkSectionExists(parser.HTTPErrors, data.Name, p) {
		e := NewConfError(ErrObjectAlreadyExists, fmt.Sprintf("%s %s already exists", parser.HTTPErrors, data.Name))

		return c.HandleError(data.Name, "", "", t, transactionID == "", e)
	}

	if err = p.SectionsCreate(parser.HTTPErrors, data.Name); err != nil {
		return c.HandleError(data.Name, "", "", t, transactionID == "", err)
	}

	if err = SerializeHTTPErrorsSection(p, data); err != nil {
		return err
	}

	return c.SaveData(p, t, transactionID == "")
}

// EditHTTPErrorsSection replaces a single http-errors section with a given name.
func (c *client) EditHTTPErrorsSection(name string, data *models.HTTPErrorsSection, transactionID string, version int64) error {
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

	if !c.checkSectionExists(parser.HTTPErrors, name, p) {
		e := NewConfError(ErrObjectAlreadyExists, fmt.Sprintf("%s %s does not exists", parser.HTTPErrors, data.Name))
		return c.HandleError(data.Name, "", "", t, transactionID == "", e)
	}

	if err = p.SectionsDelete(parser.HTTPErrors, name); err != nil {
		return c.HandleError(name, "", "", t, transactionID == "", err)
	}

	if err = p.SectionsCreate(parser.HTTPErrors, name); err != nil {
		return c.HandleError(name, "", "", t, transactionID == "", err)
	}

	if err = SerializeHTTPErrorsSection(p, data); err != nil {
		return err
	}

	return c.SaveData(p, t, transactionID == "")
}

// SerializeProgramSection saves a single http-errors section's data in the configuration.
func SerializeHTTPErrorsSection(p parser.Parser, data *models.HTTPErrorsSection) error {
	if data == nil {
		return fmt.Errorf("empty http-errors section")
	}

	for _, ef := range data.ErrorFiles {
		if err := p.Set(parser.HTTPErrors, data.Name, "errorfile", types.ErrorFile{Code: strconv.Itoa(int(ef.Code)), File: ef.File}); err != nil {
			return err
		}
	}

	return nil
}

// ParseHTTPErrorsSection parses a single http-errors section with a given name.
func ParseHTTPErrorsSection(p parser.Parser, name string) (*models.HTTPErrorsSection, error) {
	section := models.HTTPErrorsSection{
		Name: name,
	}

	// Parse errorfile entries for the section with the given name.
	// This will return an error if no valid entries were found.
	data, err := p.Get(parser.HTTPErrors, name, "errorfile")
	if err != nil {
		return nil, err
	}

	// Parse and convert errorfile entries.
	if errorFiles, ok := data.([]types.ErrorFile); ok {
		section.ErrorFiles = make([]*models.Errorfile, 0, len(errorFiles))
		for _, ef := range errorFiles {
			code, err := strconv.ParseInt(ef.Code, 10, 64)
			if err != nil {
				continue
			}
			section.ErrorFiles = append(section.ErrorFiles, &models.Errorfile{Code: code, File: ef.File})
		}
	} else {
		return nil, fmt.Errorf("http-errors section %s: unexpected type %T returned by parser", name, data)
	}

	return &section, nil
}
