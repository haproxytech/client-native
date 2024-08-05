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
//

package configuration

import (
	"errors"
	"fmt"

	strfmt "github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/client-native/v6/config-parser"
	parser_errors "github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/types"
	"github.com/haproxytech/client-native/v6/configuration/options"
	"github.com/haproxytech/client-native/v6/misc"
	"github.com/haproxytech/client-native/v6/models"
)

type MailersSection interface {
	GetMailersSections(transactionID string) (int64, models.MailersSections, error)
	GetMailersSection(name, transactionID string) (int64, *models.MailersSection, error)
	DeleteMailersSection(name, transactionID string, version int64) error
	CreateMailersSection(data *models.MailersSection, transactionID string, version int64) error
	EditMailersSection(name string, data *models.MailersSection, transactionID string, version int64) error
}

// GetMailersSections returns the configuration version and the list of
// all the mailers sections.
func (c *client) GetMailersSections(transactionID string) (int64, models.MailersSections, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	names, err := p.SectionsGet(parser.Mailers)
	if err != nil {
		return v, nil, err
	}

	mailersSections := make([]*models.MailersSection, len(names))

	for i, name := range names {
		_, ms, err := c.GetMailersSection(name, transactionID)
		if err != nil {
			return v, nil, err
		}
		mailersSections[i] = ms
	}

	return v, mailersSections, nil
}

// GetMailersSection returns a single Mailers section by name.
func (c *client) GetMailersSection(name, transactionID string) (int64, *models.MailersSection, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	if !c.checkSectionExists(parser.Mailers, name, p) {
		return v, nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("mailers section '%s' does not exist", name))
	}

	ms := &models.MailersSection{MailersSectionBase: models.MailersSectionBase{Name: name}}
	if err = ParseMailersSection(p, ms); err != nil {
		return 0, nil, err
	}
	return v, ms, err
}

// DeleteMailersSection deletes the named mailers section from configuration.
// Returns an error on failure, nil on success.
func (c *client) DeleteMailersSection(name, transactionID string, version int64) error {
	return c.deleteSection(parser.Mailers, name, transactionID, version)
}

// CreateMailersSection creates a mailers section in configuration. One of version or transactionID is
// mandatory. Returns an error on failure, nil on success.
func (c *client) CreateMailersSection(data *models.MailersSection, transactionID string, version int64) error {
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

	if c.checkSectionExists(parser.Mailers, data.Name, p) {
		e := NewConfError(ErrObjectAlreadyExists, fmt.Sprintf("%s %s already exists", parser.Mailers, data.Name))
		return c.HandleError(data.Name, "", "", t, transactionID == "", e)
	}

	if err = p.SectionsCreate(parser.Mailers, data.Name); err != nil {
		return c.HandleError(data.Name, "", "", t, transactionID == "", err)
	}

	if err = SerializeMailersSection(p, data, &c.ConfigurationOptions); err != nil {
		return err
	}

	return c.SaveData(p, t, transactionID == "")
}

func (c *client) EditMailersSection(name string, data *models.MailersSection, transactionID string, version int64) error { //nolint:revive
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

	if !c.checkSectionExists(parser.Mailers, data.Name, p) {
		e := NewConfError(ErrObjectAlreadyExists, fmt.Sprintf("%s %s does not exists", parser.Mailers, data.Name))
		return c.HandleError(data.Name, "", "", t, transactionID == "", e)
	}

	if err = SerializeMailersSection(p, data, &c.ConfigurationOptions); err != nil {
		return err
	}

	return c.SaveData(p, t, transactionID == "")
}

func ParseMailersSection(p parser.Parser, ms *models.MailersSection) error {
	// Get the optional "timeout mail" attribute
	timeout, err := p.Get(parser.Mailers, ms.Name, "timeout mail", false)
	if err != nil {
		if errors.Is(err, parser_errors.ErrFetch) {
			return nil
		}
		return err
	}

	if timeout != nil {
		to, ok := timeout.(*types.StringC)
		if !ok {
			return misc.CreateTypeAssertError("timeout mail")
		}
		ms.Timeout = misc.ParseTimeout(to.Value)
	}

	return nil
}

func SerializeMailersSection(p parser.Parser, data *models.MailersSection, opt *options.ConfigurationOptions) error {
	if data == nil {
		return fmt.Errorf("empty mailers section")
	}

	if data.Timeout != nil {
		t := types.StringC{Value: misc.SerializeTime(*data.Timeout, opt.PreferredTimeSuffix)}
		if err := p.Set(parser.Mailers, data.Name, "timeout mail", t); err != nil {
			return err
		}
	}

	return nil
}
