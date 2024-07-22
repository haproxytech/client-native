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

	"github.com/go-openapi/strfmt"
	"github.com/haproxytech/client-native/v6/models"
	parser "github.com/haproxytech/config-parser/v5"
)

type StructuredMailersSection interface {
	GetStructuredMailersSections(transactionID string) (int64, models.MailersSections, error)
	GetStructuredMailersSection(name, transactionID string) (int64, *models.MailersSection, error)
	CreateStructuredMailersSection(data *models.MailersSection, transactionID string, version int64) error
	EditStructuredMailersSection(name string, data *models.MailersSection, transactionID string, version int64) error
}

// GetStructuredMailersSection returns configuration version and a requested frontend with all its child resources.
// Returns error on fail or if frontend does not exist.
func (c *client) GetStructuredMailersSection(name string, transactionID string) (int64, *models.MailersSection, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	if !c.checkSectionExists(parser.Mailers, name, p) {
		return v, nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("MailersSection %s does not exist", name))
	}

	f, err := parseMailersSectionsSection(name, p)

	return v, f, err
}

func (c *client) GetStructuredMailersSections(transactionID string) (int64, models.MailersSections, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	frontends, err := parseMailersSectionsSections(p)
	if err != nil {
		return 0, nil, err
	}

	return v, frontends, nil
}

// EditStructuredMailersSection replaces a frontend and all it's child resources in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) EditStructuredMailersSection(name string, data *models.MailersSection, transactionID string, version int64) error {
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

	if !c.checkSectionExists(parser.Mailers, name, p) {
		e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("%s %s does not exist", parser.Mailers, name))
		return c.HandleError(name, "", "", t, transactionID == "", e)
	}

	if err = p.SectionsDelete(parser.Mailers, name); err != nil {
		return c.HandleError(name, "", "", t, transactionID == "", err)
	}

	if err = serializeMailersSectionSection(StructuredToParserArgs{
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

// CreateStructuredMailersSection creates a frontend and all it's child resources in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) CreateStructuredMailersSection(data *models.MailersSection, transactionID string, version int64) error {
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
		e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("%s %s already exist", parser.Mailers, data.Name))
		return c.HandleError(data.Name, "", "", t, transactionID == "", e)
	}

	if err = serializeMailersSectionSection(StructuredToParserArgs{
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

func parseMailersSectionsSections(p parser.Parser) (models.MailersSections, error) {
	names, err := p.SectionsGet(parser.Mailers)
	if err != nil {
		return nil, err
	}
	mss := []*models.MailersSection{}
	for _, name := range names {
		f, err := parseMailersSectionsSection(name, p)
		if err != nil {
			return nil, err
		}
		mss = append(mss, f)
	}
	return mss, nil
}

func parseMailersSectionsSection(name string, p parser.Parser) (*models.MailersSection, error) {
	ms := &models.MailersSection{MailersSectionBase: models.MailersSectionBase{Name: name}}
	if err := ParseMailersSection(p, ms); err != nil {
		return nil, err
	}

	// Mailer entries
	mailerEntries, err := ParseMailerEntries(name, p)
	if err != nil {
		return nil, err
	}
	mailerEntriesa, errmea := namedResourceArrayToMap(mailerEntries)
	if errmea != nil {
		return nil, errmea
	}
	ms.MailerEntries = mailerEntriesa

	return ms, nil
}

func serializeMailersSectionSection(a StructuredToParserArgs, ms *models.MailersSection) error {
	p := *a.Parser
	var err error
	err = p.SectionsCreate(parser.Mailers, ms.Name)
	if err != nil {
		return err
	}
	if err = SerializeMailersSection(p, ms, a.Options); err != nil {
		return err
	}
	for _, mailerEntry := range ms.MailerEntries {
		if err = p.Insert(parser.Mailers, ms.Name, "mailer", SerializeMailerEntry(mailerEntry), -1); err != nil {
			return a.HandleError(mailerEntry.Name, "mailer", ms.Name, a.TID, a.TID == "", err)
		}
	}

	return nil
}
