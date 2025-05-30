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

//nolint:dupl
package configuration

import (
	goerrors "errors"
	"strconv"

	"github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/client-native/v6/config-parser"
	parser_errors "github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/types"

	"github.com/haproxytech/client-native/v6/misc"
	"github.com/haproxytech/client-native/v6/models"
)

// GetBackendSwitchingRules returns configuration version and an array of
// configured backend switching rules in the specified frontend. Returns error on fail.
func (c *client) GetBackendSwitchingRules(frontend string, transactionID string) (int64, models.BackendSwitchingRules, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	bckRules, err := ParseBackendSwitchingRules(frontend, p)
	if err != nil {
		return v, nil, c.HandleError("", FrontendParentName, frontend, "", false, err)
	}

	return v, bckRules, nil
}

// GetBackendSwitchingRule returns configuration version and a requested backend switching rule
// in the specified frontend. Returns error on fail or if backend switching rule does not exist.
func (c *client) GetBackendSwitchingRule(id int64, frontend string, transactionID string) (int64, *models.BackendSwitchingRule, error) {
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
		return v, nil, c.HandleError(strconv.FormatInt(id, 10), FrontendParentName, frontend, "", false, err)
	}

	bckRule := ParseBackendSwitchingRule(data.(types.UseBackend))

	return v, bckRule, nil
}

// DeleteBackendSwitchingRule deletes a backend switching rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) DeleteBackendSwitchingRule(id int64, frontend string, transactionID string, version int64) error {
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	if err := p.Delete(parser.Frontends, frontend, "use_backend", int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), FrontendParentName, frontend, t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

// CreateBackendSwitchingRule creates a backend switching rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) CreateBackendSwitchingRule(id int64, frontend string, data *models.BackendSwitchingRule, transactionID string, version int64) error {
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

	if err := p.Insert(parser.Frontends, frontend, "use_backend", SerializeBackendSwitchingRule(*data), int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), FrontendParentName, frontend, t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

// EditBackendSwitchingRule edits a backend switching rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) EditBackendSwitchingRule(id int64, frontend string, data *models.BackendSwitchingRule, transactionID string, version int64) error {
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

	if _, err := p.GetOne(parser.Frontends, frontend, "use_backend", int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), FrontendParentName, frontend, t, transactionID == "", err)
	}

	if err := p.Set(parser.Frontends, frontend, "use_backend", SerializeBackendSwitchingRule(*data), int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), FrontendParentName, frontend, t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

// ReplaceBackendSwitchingRules replaces all ACL lines in configuration for a frontend.
// One of version or transactionID is mandatory.
// Returns error on fail, nil on success.
func (c *client) ReplaceBackendSwitchingRules(frontend string, data models.BackendSwitchingRules, transactionID string, version int64) error {
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

	bsRules, err := ParseBackendSwitchingRules(frontend, p)
	if err != nil {
		return c.HandleError("", FrontendParentName, frontend, "", false, err)
	}

	for i := range bsRules {
		// Always delete index 0
		if err := p.Delete(parser.Frontends, frontend, "use_backend", 0); err != nil {
			return c.HandleError(strconv.FormatInt(int64(i), 10), FrontendParentName, frontend, t, transactionID == "", err)
		}
	}

	for i, newbrRule := range data {
		if err := p.Insert(parser.Frontends, frontend, "use_backend", SerializeBackendSwitchingRule(*newbrRule), i); err != nil {
			return c.HandleError(strconv.FormatInt(int64(i), 10), FrontendParentName, frontend, t, transactionID == "", err)
		}
	}

	return c.SaveData(p, t, transactionID == "")
}

func ParseBackendSwitchingRules(frontend string, p parser.Parser) (models.BackendSwitchingRules, error) {
	var br models.BackendSwitchingRules

	data, err := p.Get(parser.Frontends, frontend, "use_backend", false)
	if err != nil {
		if goerrors.Is(err, parser_errors.ErrFetch) {
			return br, nil
		}
		return nil, err
	}

	bRules, ok := data.([]types.UseBackend)
	if !ok {
		return nil, misc.CreateTypeAssertError("[]types.DeclareCapture")
	}
	for _, bRule := range bRules {
		b := ParseBackendSwitchingRule(bRule)
		if b != nil {
			br = append(br, b)
		}
	}
	return br, nil
}

func ParseBackendSwitchingRule(ub types.UseBackend) *models.BackendSwitchingRule {
	return &models.BackendSwitchingRule{
		Name:     ub.Name,
		Cond:     ub.Cond,
		CondTest: ub.CondTest,
		Metadata: parseMetadata(ub.Comment),
	}
}

func SerializeBackendSwitchingRule(bRule models.BackendSwitchingRule) types.UseBackend {
	comment, err := serializeMetadata(bRule.Metadata)
	if err != nil {
		comment = ""
	}
	return types.UseBackend{
		Name:     bRule.Name,
		Cond:     bRule.Cond,
		CondTest: bRule.CondTest,
		Comment:  comment,
	}
}
