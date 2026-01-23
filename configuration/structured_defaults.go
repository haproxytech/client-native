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
	"errors"
	"fmt"
	"strconv"

	"github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/client-native/v6/config-parser"
	parserErrors "github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/types"
	"github.com/haproxytech/client-native/v6/models"
)

type StructuredDefaults interface {
	GetStructuredDefaultsConfiguration(transactionID string) (int64, *models.Defaults, error)
	PushStructuredDefaultsConfiguration(data *models.Defaults, transactionID string, version int64) error
	GetStructuredDefaultsSections(transactionID string) (int64, models.DefaultsSections, error)
	GetStructuredDefaultsSection(name string, transactionID string) (int64, *models.Defaults, error)
	EditStructuredDefaultsSection(name string, data *models.Defaults, transactionID string, version int64) error
	CreateStructuredDefaultsSection(data *models.Defaults, transactionID string, version int64) error
}

// GetStructuredDefaultsConfiguration returns configuration version and a
// struct representing Defaults configuration with all child resources
func (c *client) GetStructuredDefaultsConfiguration(transactionID string) (int64, *models.Defaults, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	d, err := parseDefaultsSection(parser.DefaultSectionName, p)
	if err != nil {
		return 0, nil, err
	}

	return v, d, nil
}

// GetStructuredDefaults returns configuration version and a requested defaults section with all its child resources.
// Returns error on fail or if default does not exist.
func (c *client) GetStructuredDefaultsSection(name string, transactionID string) (int64, *models.Defaults, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	if !c.checkSectionExists(parser.Defaults, name, p) {
		return v, nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Defaults %s does not exist", name))
	}

	f, err := parseDefaultsSection(name, p)

	return v, f, err
}

func (c *client) GetStructuredDefaultsSections(transactionID string) (int64, models.DefaultsSections, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	defaults, err := parseDefaultsSections(p)
	if err != nil {
		return 0, nil, err
	}

	return v, defaults, nil
}

// PushStructuredDefaultsConfiguration replaces a defaults section and all it's child resources in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) PushStructuredDefaultsConfiguration(data *models.Defaults, transactionID string, version int64) error {
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

	if !c.checkSectionExists(parser.Defaults, parser.DefaultSectionName, p) {
		e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("%s %s does not exist", parser.Defaults, parser.DefaultSectionName))
		return c.HandleError(parser.DefaultSectionName, "", "", t, transactionID == "", e)
	}

	if err = p.SectionsDelete(parser.Defaults, parser.DefaultSectionName); err != nil {
		return c.HandleError(parser.DefaultSectionName, "", "", t, transactionID == "", err)
	}

	data.Name = parser.DefaultSectionName
	if err = serializeDefaultsSection(StructuredToParserArgs{
		TID:                transactionID,
		Parser:             &p,
		Options:            &c.ConfigurationOptions,
		HandleError:        c.HandleError,
		CheckSectionExists: c.checkSectionExists,
	}, data); err != nil {
		return err
	}
	return c.SaveData(p, t, transactionID == "")
}

// EditStructuredDefaultsSection replaces a defaults section and all it's child resources in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) EditStructuredDefaultsSection(name string, data *models.Defaults, transactionID string, version int64) error {
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

	if !c.checkSectionExists(parser.Defaults, name, p) {
		e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("%s %s does not exist", parser.Defaults, name))
		return c.HandleError(name, "", "", t, transactionID == "", e)
	}

	if err = p.SectionsDelete(parser.Defaults, name); err != nil {
		return c.HandleError(name, "", "", t, transactionID == "", err)
	}

	if err = serializeDefaultsSection(StructuredToParserArgs{
		TID:                transactionID,
		Parser:             &p,
		Options:            &c.ConfigurationOptions,
		HandleError:        c.HandleError,
		CheckSectionExists: c.checkSectionExists,
	}, data); err != nil {
		return err
	}
	return c.SaveData(p, t, transactionID == "")
}

// CreateStructuredDefaultsSection creates a default and all it's child resources in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) CreateStructuredDefaultsSection(data *models.Defaults, transactionID string, version int64) error {
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

	if c.checkSectionExists(parser.Defaults, data.Name, p) {
		e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("%s %s already exist", parser.Defaults, data.Name))
		return c.HandleError(data.Name, "", "", t, transactionID == "", e)
	}

	if err = serializeDefaultsSection(StructuredToParserArgs{
		TID:                transactionID,
		Parser:             &p,
		Options:            &c.ConfigurationOptions,
		HandleError:        c.HandleError,
		CheckSectionExists: c.checkSectionExists,
	}, data); err != nil {
		return err
	}
	return c.SaveData(p, t, transactionID == "")
}

func parseDefaultsSections(p parser.Parser) (models.DefaultsSections, error) {
	names, err := p.SectionsGet(parser.Defaults)
	if err != nil {
		return nil, err
	}
	defaults := []*models.Defaults{}
	for _, name := range names {
		f, err := parseDefaultsSection(name, p)
		if err != nil {
			return nil, err
		}
		defaults = append(defaults, f)
	}
	return defaults, nil
}

func parseDefaultsSection(name string, p parser.Parser) (*models.Defaults, error) {
	d := &models.Defaults{DefaultsBase: models.DefaultsBase{Name: name}}
	if err := ParseSection(&d.DefaultsBase, parser.Defaults, name, p); err != nil {
		return nil, err
	}

	hchecks, err := ParseHTTPChecks(DefaultsParentName, name, p)
	if err != nil {
		return nil, err
	}
	d.HTTPCheckList = hchecks
	errorRules, err := ParseHTTPErrorRules(DefaultsParentName, name, p)
	if err != nil {
		return nil, err
	}
	d.HTTPErrorRuleList = errorRules
	tchecks, err := ParseTCPChecks(DefaultsParentName, name, p)
	if err != nil {
		return nil, err
	}
	d.TCPCheckRuleList = tchecks
	lt, err := ParseLogTargets(DefaultsParentName, name, p)
	if err != nil {
		return nil, err
	}
	d.LogTargetList = lt

	return d, nil
}

func serializeDefaultsSection(a StructuredToParserArgs, d *models.Defaults) error {
	p := *a.Parser
	var err error

	err = p.SectionsCreate(parser.Defaults, d.Name)
	if err != nil {
		if !errors.Is(err, parserErrors.ErrSectionAlreadyExists) || parser.DefaultSectionName != d.Name {
			return err
		}
	}
	if err = CreateEditSection(&d.DefaultsBase, parser.Defaults, d.Name, p, a.Options); err != nil {
		return a.HandleError(d.Name, "", "", a.TID, a.TID == "", err)
	}
	for i, log := range d.LogTargetList {
		if err = p.Insert(parser.Defaults, d.Name, "log", SerializeLogTarget(*log), i); err != nil {
			return a.HandleError(strconv.FormatInt(int64(i), 10), DefaultsParentName, "", a.TID, a.TID == "", err)
		}
	}
	for i, httpCheck := range d.HTTPCheckList {
		var s types.Action
		s, err = SerializeHTTPCheck(*httpCheck)
		if err != nil {
			return err
		}
		if err = p.Insert(parser.Defaults, d.Name, "http-check", s, i); err != nil {
			return a.HandleError(strconv.FormatInt(int64(i), 10), DefaultsParentName, "", a.TID, a.TID == "", err)
		}
	}
	for i, httpErrorRule := range d.HTTPErrorRuleList {
		var s types.Action
		s, err = SerializeHTTPErrorRule(*httpErrorRule)
		if err != nil {
			return err
		}
		if err = p.Insert(parser.Defaults, d.Name, "http-error", s, i); err != nil {
			return a.HandleError(strconv.FormatInt(int64(i), 10), DefaultsParentName, "", a.TID, a.TID == "", err)
		}
	}
	for i, tcpCheck := range d.TCPCheckRuleList {
		var s types.Action
		s, err = SerializeTCPCheck(*tcpCheck)
		if err != nil {
			return err
		}
		if err = p.Insert(parser.Defaults, d.Name, "tcp-check", s, i); err != nil {
			return a.HandleError(strconv.FormatInt(int64(i), 10), DefaultsParentName, "", a.TID, a.TID == "", err)
		}
	}

	return nil
}
