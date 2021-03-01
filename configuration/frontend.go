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
	parser "github.com/haproxytech/config-parser/v3"

	"github.com/haproxytech/client-native/v2/models"
)

// GetFrontends returns configuration version and an array of
// configured frontends. Returns error on fail.
func (c *Client) GetFrontends(transactionID string) (int64, models.Frontends, error) {
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
func (c *Client) GetFrontend(name string, transactionID string) (int64, *models.Frontend, error) {
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
func (c *Client) DeleteFrontend(name string, transactionID string, version int64) error {
	if err := c.deleteSection(parser.Frontends, name, transactionID, version); err != nil {
		return err
	}
	return nil
}

// EditFrontend edits a frontend in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) EditFrontend(name string, data *models.Frontend, transactionID string, version int64) error {
	if c.UseValidation {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return NewConfError(ErrValidationError, validationErr.Error())
		}
	}

	if err := c.editSection(parser.Frontends, name, data, transactionID, version); err != nil {
		return err
	}

	return nil
}

// CreateFrontend creates a frontend in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) CreateFrontend(data *models.Frontend, transactionID string, version int64) error {
	if c.UseValidation {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return NewConfError(ErrValidationError, validationErr.Error())
		}
	}

	if err := c.createSection(parser.Frontends, data.Name, data, transactionID, version); err != nil {
		return err
	}

	return nil
}
