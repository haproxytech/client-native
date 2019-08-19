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
	"strconv"

	strfmt "github.com/go-openapi/strfmt"

	parser "github.com/haproxytech/config-parser"
	parser_errors "github.com/haproxytech/config-parser/errors"
	"github.com/haproxytech/config-parser/types"
	"github.com/haproxytech/models"
)

// GetBackendSwitchingRules returns configuration version and an array of
// configured backend switching rules in the specified frontend. Returns error on fail.
func (c *Client) GetBackendSwitchingRules(frontend string, transactionID string) (int64, models.BackendSwitchingRules, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	bckRules, err := c.parseBackendSwitchingRules(frontend, p)
	if err != nil {
		return v, nil, c.handleError("", "frontend", frontend, "", false, err)
	}

	return v, bckRules, nil
}

// GetBackendSwitchingRule returns configuration version and a requested backend switching rule
// in the specified frontend. Returns error on fail or if backend switching rule does not exist.
func (c *Client) GetBackendSwitchingRule(id int64, frontend string, transactionID string) (int64, *models.BackendSwitchingRule, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	data, err := p.GetOne(parser.Frontends, frontend, "use_backend", int(id))
	if err != nil {
		return v, nil, c.handleError(strconv.FormatInt(id, 10), "frontend", frontend, "", false, err)
	}

	bckRule := parseBackendSwitchingRule(data.(types.UseBackend))
	bckRule.ID = &id

	return v, bckRule, nil
}

// DeleteBackendSwitchingRule deletes a backend switching rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) DeleteBackendSwitchingRule(id int64, frontend string, transactionID string, version int64) error {
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	if err := p.Delete(parser.Frontends, frontend, "use_backend", int(id)); err != nil {
		return c.handleError(strconv.FormatInt(id, 10), "frontend", frontend, t, transactionID == "", err)
	}

	if err := c.saveData(p, t, transactionID == ""); err != nil {
		return err
	}
	return nil
}

// CreateBackendSwitchingRule creates a backend switching rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) CreateBackendSwitchingRule(frontend string, data *models.BackendSwitchingRule, transactionID string, version int64) error {
	if c.UseValidation {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return NewConfError(ErrValidationError, validationErr.Error())
		}
	}

	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	if err := p.Insert(parser.Frontends, frontend, "use_backend", serializeBackendSwitchingRule(*data), int(*data.ID)); err != nil {
		return c.handleError(strconv.FormatInt(*data.ID, 10), "frontend", frontend, t, transactionID == "", err)
	}

	if err := c.saveData(p, t, transactionID == ""); err != nil {
		return err
	}

	return nil
}

// EditBackendSwitchingRule edits a backend switching rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) EditBackendSwitchingRule(id int64, frontend string, data *models.BackendSwitchingRule, transactionID string, version int64) error {
	if c.UseValidation {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return NewConfError(ErrValidationError, validationErr.Error())
		}
	}
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	if _, err := p.GetOne(parser.Frontends, frontend, "use_backend", int(id)); err != nil {
		return c.handleError(strconv.FormatInt(id, 10), "frontend", frontend, t, transactionID == "", err)
	}

	if err := p.Set(parser.Frontends, frontend, "use_backend", serializeBackendSwitchingRule(*data), int(id)); err != nil {
		return c.handleError(strconv.FormatInt(id, 10), "frontend", frontend, t, transactionID == "", err)
	}

	if err := c.saveData(p, t, transactionID == ""); err != nil {
		return err
	}

	return nil
}

func (c *Client) parseBackendSwitchingRules(frontend string, p *parser.Parser) (models.BackendSwitchingRules, error) {
	br := models.BackendSwitchingRules{}

	data, err := p.Get(parser.Frontends, frontend, "use_backend", false)
	if err != nil {
		if err == parser_errors.ErrFetch {
			return br, nil
		}
		return nil, err
	}

	bRules := data.([]types.UseBackend)
	for i, bRule := range bRules {
		id := int64(i)
		b := parseBackendSwitchingRule(bRule)
		if b != nil {
			b.ID = &id
			br = append(br, b)
		}
	}
	return br, nil
}

func parseBackendSwitchingRule(ub types.UseBackend) *models.BackendSwitchingRule {
	return &models.BackendSwitchingRule{
		Name:     ub.Name,
		Cond:     ub.Cond,
		CondTest: ub.CondTest,
	}
}

func serializeBackendSwitchingRule(bRule models.BackendSwitchingRule) types.UseBackend {
	return types.UseBackend{
		Name:     bRule.Name,
		Cond:     bRule.Cond,
		CondTest: bRule.CondTest,
	}
}
