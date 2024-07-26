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
	"errors"
	"strconv"

	strfmt "github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/config-parser/v5"
	parser_errors "github.com/haproxytech/config-parser/v5/errors"
	"github.com/haproxytech/config-parser/v5/types"

	"github.com/haproxytech/client-native/v6/misc"
	"github.com/haproxytech/client-native/v6/models"
)

type StickRule interface {
	GetStickRules(backend string, transactionID string) (int64, models.StickRules, error)
	GetStickRule(id int64, backend string, transactionID string) (int64, *models.StickRule, error)
	DeleteStickRule(id int64, backend string, transactionID string, version int64) error
	CreateStickRule(id int64, backend string, data *models.StickRule, transactionID string, version int64) error
	EditStickRule(id int64, backend string, data *models.StickRule, transactionID string, version int64) error
	ReplaceStickRules(backend string, data models.StickRules, transactionID string, version int64) error
}

// GetStickRules returns configuration version and an array of
// configured stick rules in the specified backend. Returns error on fail.
func (c *client) GetStickRules(backend string, transactionID string) (int64, models.StickRules, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	sRules, err := ParseStickRules(backend, p)
	if err != nil {
		return v, nil, c.HandleError("", BackendParentName, backend, "", false, err)
	}

	return v, sRules, nil
}

// GetStickRule returns configuration version and a requested stick rule
// in the specified backend. Returns error on fail or if stick rule does not exist.
func (c *client) GetStickRule(id int64, backend string, transactionID string) (int64, *models.StickRule, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	data, err := p.GetOne(parser.Backends, backend, "stick", int(id))
	if err != nil {
		return v, nil, c.HandleError(strconv.FormatInt(id, 10), BackendParentName, backend, "", false, err)
	}

	sRule := ParseStickRule(data.(types.Stick))

	return v, sRule, nil
}

// DeleteStickRule deletes a stick rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) DeleteStickRule(id int64, backend string, transactionID string, version int64) error {
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	if err := p.Delete(parser.Backends, backend, "stick", int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), BackendParentName, backend, t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

// CreateStickRule creates a stick rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) CreateStickRule(id int64, backend string, data *models.StickRule, transactionID string, version int64) error {
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

	if err := p.Insert(parser.Backends, backend, "stick", SerializeStickRule(*data), int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), BackendParentName, backend, t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

// EditStickRule edits a stick rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) EditStickRule(id int64, backend string, data *models.StickRule, transactionID string, version int64) error {
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

	if _, err := p.GetOne(parser.Backends, backend, "stick", int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), BackendParentName, backend, t, transactionID == "", err)
	}

	if err := p.Set(parser.Backends, backend, "stick", SerializeStickRule(*data), int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), BackendParentName, backend, t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

// ReplaceStickRules replaces all Stick rule lines in configuration for a backend.
// One of version or transactionID is mandatory.
// Returns error on fail, nil on success.
func (c *client) ReplaceStickRules(backend string, data models.StickRules, transactionID string, version int64) error {
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

	stickRules, err := ParseStickRules(backend, p)
	if err != nil {
		return c.HandleError("", BackendParentName, backend, "", false, err)
	}

	for i := range stickRules {
		// Always delete index 0
		if err := p.Delete(parser.Backends, backend, "stick", 0); err != nil {
			return c.HandleError(strconv.FormatInt(int64(i), 10), BackendParentName, backend, t, transactionID == "", err)
		}
	}

	for i, newStickRule := range data {
		if err := p.Insert(parser.Backends, backend, "stick", SerializeStickRule(*newStickRule), i); err != nil {
			return c.HandleError(strconv.FormatInt(int64(i), 10), BackendParentName, backend, t, transactionID == "", err)
		}
	}

	return c.SaveData(p, t, transactionID == "")
}

func ParseStickRules(backend string, p parser.Parser) (models.StickRules, error) {
	var sr models.StickRules

	data, err := p.Get(parser.Backends, backend, "stick", false)
	if err != nil {
		if errors.Is(err, parser_errors.ErrFetch) {
			return sr, nil
		}
		return nil, err
	}

	sRules, ok := data.([]types.Stick)
	if !ok {
		return nil, misc.CreateTypeAssertError("stick rules")
	}
	for _, sRule := range sRules {
		s := ParseStickRule(sRule)
		if s != nil {
			sr = append(sr, s)
		}
	}
	return sr, nil
}

func ParseStickRule(s types.Stick) *models.StickRule {
	return &models.StickRule{
		Type:     s.Type,
		Table:    s.Table,
		Pattern:  s.Pattern,
		Cond:     s.Cond,
		CondTest: s.CondTest,
	}
}

func SerializeStickRule(sRule models.StickRule) types.Stick {
	sr := types.Stick{
		Type:     sRule.Type,
		Table:    sRule.Table,
		Pattern:  sRule.Pattern,
		Cond:     sRule.Cond,
		CondTest: sRule.CondTest,
	}
	return sr
}
