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
	"errors"
	"strconv"

	strfmt "github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/config-parser/v5"
	parser_errors "github.com/haproxytech/config-parser/v5/errors"
	"github.com/haproxytech/config-parser/v5/types"

	"github.com/haproxytech/client-native/v6/misc"
	"github.com/haproxytech/client-native/v6/models"
)

// GetServerSwitchingRules returns configuration version and an array of
// configured server switching rules in the specified backend. Returns error on fail.
func (c *client) GetServerSwitchingRules(backend string, transactionID string) (int64, models.ServerSwitchingRules, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	srvRules, err := ParseServerSwitchingRules(backend, p)
	if err != nil {
		return v, nil, c.HandleError("", BackendParentName, backend, "", false, err)
	}

	return v, srvRules, nil
}

// GetServerSwitchingRule returns configuration version and a requested server switching rule
// in the specified backend. Returns error on fail or if server switching rule does not exist.
func (c *client) GetServerSwitchingRule(id int64, backend string, transactionID string) (int64, *models.ServerSwitchingRule, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	data, err := p.GetOne(parser.Backends, backend, "use-server", int(id))
	if err != nil {
		return v, nil, c.HandleError(strconv.FormatInt(id, 10), BackendParentName, backend, "", false, err)
	}

	srvRule := ParseServerSwitchingRule(data.(types.UseServer))

	return v, srvRule, nil
}

// DeleteServerSwitchingRule deletes a server switching rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) DeleteServerSwitchingRule(id int64, backend string, transactionID string, version int64) error {
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	if err := p.Delete(parser.Backends, backend, "use-server", int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), BackendParentName, backend, t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

// CreateServerSwitchingRule creates a server switching rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) CreateServerSwitchingRule(id int64, backend string, data *models.ServerSwitchingRule, transactionID string, version int64) error {
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

	if err := p.Insert(parser.Backends, backend, "use-server", SerializeServerSwitchingRule(*data), int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), BackendParentName, backend, t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

// EditServerSwitchingRule edits a server switching rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) EditServerSwitchingRule(id int64, backend string, data *models.ServerSwitchingRule, transactionID string, version int64) error {
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

	if _, err := p.GetOne(parser.Backends, backend, "use-server", int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), BackendParentName, backend, t, transactionID == "", err)
	}

	if err := p.Set(parser.Backends, backend, "use-server", SerializeServerSwitchingRule(*data), int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), BackendParentName, backend, t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

// ReplaceServerSwitchingRules replaces all ACL lines in configuration for a backend.
// One of version or transactionID is mandatory.
// Returns error on fail, nil on success.
func (c *client) ReplaceServerSwitchingRules(backend string, data models.ServerSwitchingRules, transactionID string, version int64) error {
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

	ssRules, err := ParseServerSwitchingRules(backend, p)
	if err != nil {
		return c.HandleError("", BackendParentName, backend, "", false, err)
	}

	for i := range ssRules {
		// Always delete index 0
		if err := p.Delete(parser.Backends, backend, "use-server", 0); err != nil {
			return c.HandleError(strconv.FormatInt(int64(i), 10), BackendParentName, backend, t, transactionID == "", err)
		}
	}

	for i, newssRule := range data {
		if err := p.Insert(parser.Backends, backend, "use-server", SerializeServerSwitchingRule(*newssRule), i); err != nil {
			return c.HandleError(strconv.FormatInt(int64(i), 10), BackendParentName, backend, t, transactionID == "", err)
		}
	}

	return c.SaveData(p, t, transactionID == "")
}

func ParseServerSwitchingRules(backend string, p parser.Parser) (models.ServerSwitchingRules, error) {
	var sr models.ServerSwitchingRules

	data, err := p.Get(parser.Backends, backend, "use-server", false)
	if err != nil {
		if errors.Is(err, parser_errors.ErrFetch) {
			return sr, nil
		}
		return nil, err
	}

	sRules, ok := data.([]types.UseServer)
	if !ok {
		return nil, misc.CreateTypeAssertError("use-server")
	}
	for _, sRule := range sRules {
		s := ParseServerSwitchingRule(sRule)
		if s != nil {
			sr = append(sr, s)
		}
	}
	return sr, nil
}

func ParseServerSwitchingRule(us types.UseServer) *models.ServerSwitchingRule {
	return &models.ServerSwitchingRule{
		TargetServer: us.Name,
		Cond:         us.Cond,
		CondTest:     us.CondTest,
	}
}

func SerializeServerSwitchingRule(sRule models.ServerSwitchingRule) types.UseServer {
	return types.UseServer{
		Name:     sRule.TargetServer,
		Cond:     sRule.Cond,
		CondTest: sRule.CondTest,
	}
}
