// Copyright 2025 HAProxy Technologies
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
	parser "github.com/haproxytech/client-native/v6/config-parser"
	"github.com/haproxytech/client-native/v6/models"
)

type StructuredAcmeProvider interface {
	GetStructuredAcmeProviders(transactionID string) (int64, models.AcmeProviders, error)
	GetStructuredAcmeProvider(name string, transactionID string) (int64, *models.AcmeProvider, error)
	CreateStructuredAcmeProvider(data *models.AcmeProvider, transactionID string, version int64) error
	EditStructuredAcmeProvider(name string, data *models.AcmeProvider, transactionID string, version int64) error
}

func (c *client) GetStructuredAcmeProvider(name string, transactionID string) (int64, *models.AcmeProvider, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	if !p.SectionExists(parser.Acme, name) {
		return v, nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("acme provider %s does not exist", name))
	}

	f, err := ParseAcmeProvider(p, name)

	return v, f, err
}

func (c *client) GetStructuredAcmeProviders(transactionID string) (int64, models.AcmeProviders, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	acmes, err := parseStructuredAcmeProviders(p)
	if err != nil {
		return 0, nil, err
	}

	return v, acmes, nil
}

func (c *client) EditStructuredAcmeProvider(name string, data *models.AcmeProvider, transactionID string, version int64) error {
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

	if !p.SectionExists(parser.Acme, name) {
		e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("%s %s does not exist", parser.Acme, name))
		return c.HandleError(name, "", "", t, transactionID == "", e)
	}

	if err = p.SectionsDelete(parser.Acme, name); err != nil {
		return c.HandleError(name, "", "", t, transactionID == "", err)
	}

	if err = serializeStructuredAcmeProvider(StructuredToParserArgs{
		TID:         transactionID,
		Parser:      &p,
		Options:     &c.ConfigurationOptions,
		HandleError: c.HandleError,
	}, data); err != nil {
		return err
	}
	return c.SaveData(p, t, transactionID == "")
}

func (c *client) CreateStructuredAcmeProvider(data *models.AcmeProvider, transactionID string, version int64) error {
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

	if p.SectionExists(parser.Acme, data.Name) {
		e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("%s %s already exist", parser.Acme, data.Name))
		return c.HandleError(data.Name, "", "", t, transactionID == "", e)
	}

	if err = serializeStructuredAcmeProvider(StructuredToParserArgs{
		TID:         transactionID,
		Parser:      &p,
		Options:     &c.ConfigurationOptions,
		HandleError: c.HandleError,
	}, data); err != nil {
		return err
	}
	return c.SaveData(p, t, transactionID == "")
}

func parseStructuredAcmeProviders(p parser.Parser) (models.AcmeProviders, error) {
	names, err := p.SectionsGet(parser.Acme)
	if err != nil {
		return nil, err
	}
	acmes := make(models.AcmeProviders, len(names))
	for i, name := range names {
		ap, err := ParseAcmeProvider(p, name)
		if err != nil {
			return nil, err
		}
		acmes[i] = ap
	}
	return acmes, nil
}

func serializeStructuredAcmeProvider(a StructuredToParserArgs, ap *models.AcmeProvider) error {
	p := *a.Parser

	err := p.SectionsCreate(parser.Acme, ap.Name)
	if err != nil {
		return a.HandleError(ap.Name, "", "", a.TID, a.TID == "", err)
	}

	err = SerializeAcmeProvider(p, ap)
	if err != nil {
		return a.HandleError(ap.Name, "", "", a.TID, a.TID == "", err)
	}

	return nil
}
