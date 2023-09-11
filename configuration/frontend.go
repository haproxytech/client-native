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

	"github.com/haproxytech/client-native/v5/models"
)

type Frontend interface {
	GetFrontends(transactionID string) (int64, models.Frontends, error)
	GetFrontend(name string, transactionID string) (int64, *models.Frontend, error)
	DeleteFrontend(name string, transactionID string, version int64) error
	EditFrontend(name string, data *models.Frontend, transactionID string, version int64) error
	CreateFrontend(data *models.Frontend, transactionID string, version int64) error
}

// GetFrontends returns configuration version and an array of
// configured frontends. Returns error on fail.
func (c *client) GetFrontends(transactionID string) (int64, models.Frontends, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	fNames, err := p.SectionsGet(parser.Frontends)
	if err != nil {
		return v, nil, err
	}

	frontends := []*models.Frontend{}
	for _, name := range fNames {
		f := &models.Frontend{Name: name}
		if err := ParseSection(f, parser.Frontends, name, p); err != nil {
			continue
		}
		frontends = append(frontends, f)
	}

	return v, frontends, nil
}

// GetFrontend returns configuration version and a requested frontend.
// Returns error on fail or if frontend does not exist.
func (c *client) GetFrontend(name string, transactionID string) (int64, *models.Frontend, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	if !c.checkSectionExists(parser.Frontends, name, p) {
		return v, nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Frontend %s does not exist", name))
	}

	frontend := &models.Frontend{Name: name}
	if err := ParseSection(frontend, parser.Frontends, name, p); err != nil {
		return v, nil, err
	}

	return v, frontend, nil
}

// DeleteFrontend deletes a frontend in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) DeleteFrontend(name string, transactionID string, version int64) error {
	return c.deleteSection(parser.Frontends, name, transactionID, version)
}

// EditFrontend edits a frontend in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) EditFrontend(name string, data *models.Frontend, transactionID string, version int64) error {
	if c.UseModelsValidation {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return NewConfError(ErrValidationError, validationErr.Error())
		}
	}

	return c.editSection(parser.Frontends, name, data, transactionID, version)
}

// CreateFrontend creates a frontend in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) CreateFrontend(data *models.Frontend, transactionID string, version int64) error {
	if c.UseModelsValidation {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return NewConfError(ErrValidationError, validationErr.Error())
		}
	}

	return c.createSection(parser.Frontends, data.Name, data, transactionID, version)
}
