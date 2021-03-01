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

// GetBackends returns configuration version and an array of
// configured backends. Returns error on fail.
func (c *Client) GetBackends(transactionID string) (int64, models.Backends, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	bNames, err := p.SectionsGet(parser.Backends)
	if err != nil {
		return v, nil, err
	}

	backends := []*models.Backend{}
	for _, name := range bNames {
		b := &models.Backend{Name: name}
		if err := ParseSection(b, parser.Backends, name, p); err != nil {
			continue
		}
		backends = append(backends, b)
	}

	return v, backends, nil
}

// GetBackend returns configuration version and a requested backend.
// Returns error on fail or if backend does not exist.
func (c *Client) GetBackend(name string, transactionID string) (int64, *models.Backend, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	if !c.checkSectionExists(parser.Backends, name, p) {
		return v, nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Backend %s does not exist", name))
	}

	backend := &models.Backend{Name: name}
	if err := ParseSection(backend, parser.Backends, name, p); err != nil {
		return v, nil, err
	}

	return v, backend, nil
}

// DeleteBackend deletes a backend in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) DeleteBackend(name string, transactionID string, version int64) error {
	if err := c.deleteSection(parser.Backends, name, transactionID, version); err != nil {
		return err
	}
	return nil
}

// CreateBackend creates a backend in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) CreateBackend(data *models.Backend, transactionID string, version int64) error {
	if c.UseValidation {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return NewConfError(ErrValidationError, validationErr.Error())
		}
	}
	if err := c.createSection(parser.Backends, data.Name, data, transactionID, version); err != nil {
		return err
	}
	return nil
}

// EditBackend edits a backend in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) EditBackend(name string, data *models.Backend, transactionID string, version int64) error {
	if c.UseValidation {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return NewConfError(ErrValidationError, validationErr.Error())
		}
	}
	if err := c.editSection(parser.Backends, name, data, transactionID, version); err != nil {
		return err
	}
	return nil
}
