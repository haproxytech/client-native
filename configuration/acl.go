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

// GetACLs returns configuration version and an array of
// configured ACL lines in the specified parent. Returns error on fail.
func (c *Client) GetACLs(parentType, parentName string, transactionID string) (int64, models.Acls, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	acls, err := c.parseACLs(parentType, parentName, p)
	if err != nil {
		return v, nil, c.handleError("", parentType, parentName, "", false, err)
	}

	return v, acls, nil
}

// GetACL returns configuration version and a requested ACL line
// in the specified parent. Returns error on fail or if ACL line does not exist.
func (c *Client) GetACL(id int64, parentType, parentName string, transactionID string) (int64, *models.ACL, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	var section parser.Section
	if parentType == "backend" {
		section = parser.Backends
	} else if parentType == "frontend" {
		section = parser.Frontends
	}

	data, err := p.GetOne(section, parentName, "acl", int(id))
	if err != nil {
		return v, nil, c.handleError(strconv.FormatInt(id, 10), parentType, parentName, "", false, err)
	}

	acl := parseACL(data.(types.ACL))
	acl.ID = &id

	return v, acl, nil
}

// DeleteACL deletes a ACL line in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) DeleteACL(id int64, parentType string, parentName string, transactionID string, version int64) error {
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	var section parser.Section
	if parentType == "backend" {
		section = parser.Backends
	} else if parentType == "frontend" {
		section = parser.Frontends
	}

	if err := p.Delete(section, parentName, "acl", int(id)); err != nil {
		return c.handleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	if err := c.saveData(p, t, transactionID == ""); err != nil {
		return err
	}
	return nil
}

// CreateACL creates a ACL line in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) CreateACL(parentType string, parentName string, data *models.ACL, transactionID string, version int64) error {
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

	var section parser.Section
	if parentType == "backend" {
		section = parser.Backends
	} else if parentType == "frontend" {
		section = parser.Frontends
	}

	if err := p.Insert(section, parentName, "acl", serializeACL(*data), int(*data.ID)); err != nil {
		return c.handleError(strconv.FormatInt(*data.ID, 10), parentType, parentName, t, transactionID == "", err)
	}

	if err := c.saveData(p, t, transactionID == ""); err != nil {
		return err
	}
	return nil
}

// EditACL edits a ACL line in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) EditACL(id int64, parentType string, parentName string, data *models.ACL, transactionID string, version int64) error {
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

	var section parser.Section
	if parentType == "backend" {
		section = parser.Backends
	} else if parentType == "frontend" {
		section = parser.Frontends
	}

	if _, err := p.GetOne(section, parentName, "acl", int(id)); err != nil {
		return c.handleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	if err := p.Set(section, parentName, "acl", serializeACL(*data), int(id)); err != nil {
		return c.handleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	if err := c.saveData(p, t, transactionID == ""); err != nil {
		return err
	}
	return nil
}

func (c *Client) parseACLs(t, pName string, p *parser.Parser) (models.Acls, error) {
	section := parser.Global
	if t == "frontend" {
		section = parser.Frontends
	} else if t == "backend" {
		section = parser.Backends
	}

	acls := models.Acls{}
	data, err := p.Get(section, pName, "acl", false)
	if err != nil {
		if err == parser_errors.ErrFetch {
			return acls, nil
		}
		return nil, err
	}

	aclLines := data.([]types.ACL)
	for i, r := range aclLines {
		id := int64(i)
		acl := parseACL(r)
		if acl != nil {
			acl.ID = &id
			acls = append(acls, acl)
		}
	}
	return acls, nil
}

func parseACL(f types.ACL) *models.ACL {
	return &models.ACL{
		ACLName:   f.Name,
		Criterion: f.Criterion,
		Value:     f.Value,
	}
}

func serializeACL(f models.ACL) types.ACL {
	return types.ACL{
		Name:      f.ACLName,
		Criterion: f.Criterion,
		Value:     f.Value,
	}
}
