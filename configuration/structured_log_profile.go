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

type StructuredLogProfile interface {
	GetStructuredLogProfiles(transactionID string) (int64, models.LogProfiles, error)
	GetStructuredLogProfile(name string, transactionID string) (int64, *models.LogProfile, error)
	CreateStructuredLogProfile(data *models.LogProfile, transactionID string, version int64) error
	EditStructuredLogProfile(name string, data *models.LogProfile, transactionID string, version int64) error
}

func (c *client) GetStructuredLogProfile(name string, transactionID string) (int64, *models.LogProfile, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	if !p.SectionExists(parser.LogProfile, name) {
		return v, nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("LogProfile %s does not exist", name))
	}

	f, err := ParseLogProfile(p, name)

	return v, f, err
}

func (c *client) GetStructuredLogProfiles(transactionID string) (int64, models.LogProfiles, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	rings, err := parseLogProfilesSections(p)
	if err != nil {
		return 0, nil, err
	}

	return v, rings, nil
}

func (c *client) EditStructuredLogProfile(name string, data *models.LogProfile, transactionID string, version int64) error {
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

	if !p.SectionExists(parser.LogProfile, name) {
		e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("%s %s does not exist", parser.LogProfile, name))
		return c.HandleError(name, "", "", t, transactionID == "", e)
	}

	if err = p.SectionsDelete(parser.LogProfile, name); err != nil {
		return c.HandleError(name, "", "", t, transactionID == "", err)
	}

	if err = serializeLogProfileSection(StructuredToParserArgs{
		TID:         transactionID,
		Parser:      &p,
		Options:     &c.ConfigurationOptions,
		HandleError: c.HandleError,
	}, data); err != nil {
		return err
	}
	return c.SaveData(p, t, transactionID == "")
}

func (c *client) CreateStructuredLogProfile(data *models.LogProfile, transactionID string, version int64) error {
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

	if p.SectionExists(parser.LogProfile, data.Name) {
		e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("%s %s already exist", parser.LogProfile, data.Name))
		return c.HandleError(data.Name, "", "", t, transactionID == "", e)
	}

	if err = serializeLogProfileSection(StructuredToParserArgs{
		TID:         transactionID,
		Parser:      &p,
		Options:     &c.ConfigurationOptions,
		HandleError: c.HandleError,
	}, data); err != nil {
		return err
	}
	return c.SaveData(p, t, transactionID == "")
}

func parseLogProfilesSections(p parser.Parser) (models.LogProfiles, error) {
	names, err := p.SectionsGet(parser.LogProfile)
	if err != nil {
		return nil, err
	}
	profiles := make(models.LogProfiles, len(names))
	for i, name := range names {
		lp, err := ParseLogProfile(p, name)
		if err != nil {
			return nil, err
		}
		profiles[i] = lp
	}
	return profiles, nil
}

func serializeLogProfileSection(a StructuredToParserArgs, r *models.LogProfile) error {
	p := *a.Parser

	err := p.SectionsCreate(parser.LogProfile, r.Name)
	if err != nil {
		return err
	}

	return SerializeLogProfile(p, r)
}
