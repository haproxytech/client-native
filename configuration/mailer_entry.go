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
	"github.com/haproxytech/client-native/v5/misc"
	"github.com/haproxytech/client-native/v5/models"
	parser "github.com/haproxytech/config-parser/v5"
	parser_errors "github.com/haproxytech/config-parser/v5/errors"
	"github.com/haproxytech/config-parser/v5/types"
)

type MailerEntry interface {
	GetMailerEntries(mailersSection string, transactionID string) (int64, models.MailerEntries, error)
	GetMailerEntry(name string, mailersSection string, transactionID string) (int64, *models.MailerEntry, error)
	DeleteMailerEntry(name string, mailersSection string, transactionID string, version int64) error
	CreateMailerEntry(mailersSection string, data *models.MailerEntry, transactionID string, version int64) error
	EditMailerEntry(name string, mailersSection string, data *models.MailerEntry, transactionID string, version int64) error
}

// GetMailersEntries returns configuration version and an array of
// mailer entries in the specified mailers section. Returns error on fail.
func (c *client) GetMailerEntries(mailersSection string, transactionID string) (int64, models.MailerEntries, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	mailerEntries, err := ParseMailerEntries(mailersSection, p)
	if err != nil {
		return v, nil, c.HandleError("", "mailers", mailersSection, "", false, err)
	}

	return v, mailerEntries, nil
}

// GetMailerEntry returns configuration version and a requested mailer entry
// in the specified mailers section. Returns error on fail or if bind does not exist.
func (c *client) GetMailerEntry(name string, mailersSection string, transactionID string) (int64, *models.MailerEntry, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	mailerEntry, _ := GetMailerEntryByName(name, mailersSection, p)
	if mailerEntry == nil {
		return v, nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("MailerEntry %s does not exist in mailers section %s", name, mailersSection))
	}

	return v, mailerEntry, nil
}

// DeleteMailerEntry deletes an mailer entry in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) DeleteMailerEntry(name string, mailersSection string, transactionID string, version int64) error {
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	mailerEntry, i := GetMailerEntryByName(name, mailersSection, p)
	if mailerEntry == nil {
		e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("MailerEntry %s does not exist in mailers section %s", name, mailersSection))
		return c.HandleError(name, "mailers", mailersSection, t, transactionID == "", e)
	}

	if err := p.Delete(parser.Mailers, mailersSection, "mailer", i); err != nil {
		return c.HandleError(name, "mailers", mailersSection, t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

// CreateMailerEntry creates a mailer entry in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) CreateMailerEntry(mailersSection string, data *models.MailerEntry, transactionID string, version int64) error {
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

	mailerEntry, _ := GetMailerEntryByName(data.Name, mailersSection, p)
	if mailerEntry != nil {
		e := NewConfError(ErrObjectAlreadyExists, fmt.Sprintf("MailerEntry %s already exists in peer section %s", data.Name, mailersSection))
		return c.HandleError(data.Name, "mailers", mailersSection, t, transactionID == "", e)
	}

	if err := p.Insert(parser.Mailers, mailersSection, "mailer", SerializeMailerEntry(*data), -1); err != nil {
		return c.HandleError(data.Name, "mailers", mailersSection, t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

// EditMailerEntry edits a mailer entry in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) EditMailerEntry(name string, mailersSection string, data *models.MailerEntry, transactionID string, version int64) error {
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

	mailerEntry, i := GetMailerEntryByName(name, mailersSection, p)
	if mailerEntry == nil {
		e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("MailerEntry %s does not exist in peer section %s", name, mailersSection))
		return c.HandleError(data.Name, "mailers", mailersSection, t, transactionID == "", e)
	}

	if err := p.Set(parser.Mailers, mailersSection, "mailer", SerializeMailerEntry(*data), i); err != nil {
		return c.HandleError(data.Name, "mailers", mailersSection, t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

func ParseMailerEntries(mailersSection string, p parser.Parser) (models.MailerEntries, error) {
	mailerEntry := models.MailerEntries{}

	data, err := p.Get(parser.Mailers, mailersSection, "mailer", false)
	if err != nil {
		if errors.Is(err, parser_errors.ErrFetch) {
			return mailerEntry, nil
		}
		return nil, err
	}

	mailerEntries, ok := data.([]types.Mailer)
	if !ok {
		return nil, misc.CreateTypeAssertError("mailer")
	}
	for _, e := range mailerEntries {
		me := ParseMailerEntry(e)
		if me != nil {
			mailerEntry = append(mailerEntry, me)
		}
	}

	return mailerEntry, nil
}

func ParseMailerEntry(m types.Mailer) *models.MailerEntry {
	return &models.MailerEntry{
		Name:    m.Name,
		Address: m.IP,
		Port:    m.Port,
	}
}

func SerializeMailerEntry(me models.MailerEntry) types.Mailer {
	return types.Mailer{
		Name: me.Name,
		IP:   me.Address,
		Port: me.Port,
	}
}

func GetMailerEntryByName(name, section string, p parser.Parser) (*models.MailerEntry, int) {
	mailerEntries, err := ParseMailerEntries(section, p)
	if err != nil {
		return nil, 0
	}

	for i, e := range mailerEntries {
		if e.Name == name {
			return e, i
		}
	}
	return nil, 0
}
