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

package spoe

import (
	"fmt"

	"github.com/go-openapi/strfmt"

	conf "github.com/haproxytech/client-native/v2/configuration"
	"github.com/haproxytech/client-native/v2/models"
)

// GetScopes returns configuration version and an array of
// configured scopes Returns error on fail.
func (c *SingleSpoe) GetScopes(transactionID string) (int64, models.SpoeScopes, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}
	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	sNames := p.Parsers

	scopes := models.SpoeScopes{}
	for name := range sNames {
		if p.IsScope(name) {
			name := models.SpoeScope(name)
			scopes = append(scopes, name)
		}
	}
	return v, scopes, nil
}

// GetScope returns configuration version and a requested scope
// Returns error on fail or if scope does not exist.
func (c *SingleSpoe) GetScope(name, transactionID string) (int64, *models.SpoeScope, error) {
	v, scopes, err := c.GetScopes(transactionID)
	if err != nil {
		return 0, nil, err
	}
	for _, s := range scopes {
		if s == models.SpoeScope(name) {
			return v, &s, nil
		}
	}
	return v, nil, conf.NewConfError(conf.ErrObjectDoesNotExist, fmt.Sprintf("scope %s does not exist", name))
}

// DeleteScope deletes a scope in configuration
// One of version or transactionID is mandatory
// Returns error on fail, nil on success.
func (c *SingleSpoe) DeleteScope(name, transactionID string, version int64) error {
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	_, _, err = c.GetScope(name, transactionID)
	if err != nil {
		e := conf.NewConfError(conf.ErrObjectDoesNotExist, fmt.Sprintf("scope %s does not exist", name))
		return c.Transaction.HandleError(name, "", "", t, transactionID == "", e)
	}

	if err := p.ScopeDelete(name); err != nil {
		return c.Transaction.HandleError(name, "", "", t, transactionID == "", err)
	}

	if err := c.Transaction.SaveData(p, t, transactionID == ""); err != nil {
		return err
	}

	return nil
}

// CreateScope creates a scope in configuration
// One of version or transactionID is mandatory
// Returns error on fail, nil on success.
func (c *SingleSpoe) CreateScope(data *models.SpoeScope, transactionID string, version int64) error {
	if c.Transaction.UseValidation {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return conf.NewConfError(conf.ErrValidationError, validationErr.Error())
		}
	}

	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	scope := string(*data)
	_, s, _ := c.GetScope(scope, transactionID)
	if s != nil {
		e := conf.NewConfError(conf.ErrObjectAlreadyExists, fmt.Sprintf("scope %s already exists", scope))
		return c.Transaction.HandleError(scope, "", "", t, transactionID == "", e)
	}

	if err := p.ScopeCreate(scope); err != nil {
		return err
	}

	if err := c.Transaction.SaveData(p, t, transactionID == ""); err != nil {
		return err
	}

	return nil
}
