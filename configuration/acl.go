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

	"github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/config-parser/v5"
	parser_errors "github.com/haproxytech/config-parser/v5/errors"
	"github.com/haproxytech/config-parser/v5/types"

	"github.com/haproxytech/client-native/v6/models"
)

type ACL interface {
	GetACLs(parentType, parentName string, transactionID string, aclName ...string) (int64, models.Acls, error)
	GetACL(id int64, parentType, parentName string, transactionID string) (int64, *models.ACL, error)
	DeleteACL(id int64, parentType string, parentName string, transactionID string, version int64) error
	CreateACL(id int64, parentType string, parentName string, data *models.ACL, transactionID string, version int64) error
	EditACL(id int64, parentType string, parentName string, data *models.ACL, transactionID string, version int64) error
	ReplaceAcls(parentType string, parentName string, data models.Acls, transactionID string, version int64) error
}

// GetACLs returns configuration version and an array of
// configured ACL lines in the specified parent. Returns error on fail.
func (c *client) GetACLs(parentType, parentName string, transactionID string, aclName ...string) (int64, models.Acls, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	section, parentName, err := getParserFromParent("acl", parentType, parentName)
	if err != nil {
		return 0, nil, err
	}

	acls, err := ParseACLs(section, parentName, p, aclName...)
	if err != nil {
		return v, nil, c.HandleError("", parentType, parentName, "", false, err)
	}

	return v, acls, nil
}

// GetACL returns configuration version and a requested ACL line
// in the specified parent. Returns error on fail or if ACL line does not exist.
func (c *client) GetACL(id int64, parentType, parentName string, transactionID string) (int64, *models.ACL, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	section, parentName, err := getParserFromParent("acl", parentType, parentName)
	if err != nil {
		return v, nil, err
	}

	data, err := p.GetOne(section, parentName, "acl", int(id))
	if err != nil {
		return v, nil, c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, "", false, err)
	}

	acl := ParseACL(data.(types.ACL))

	return v, acl, nil
}

// DeleteACL deletes a ACL line in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) DeleteACL(id int64, parentType string, parentName string, transactionID string, version int64) error {
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	section, parentName, err := getParserFromParent("acl", parentType, parentName)
	if err != nil {
		return err
	}

	if err := p.Delete(section, parentName, "acl", int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

// CreateACL creates a ACL line in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) CreateACL(id int64, parentType string, parentName string, data *models.ACL, transactionID string, version int64) error {
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

	section, parentName, err := getParserFromParent("acl", parentType, parentName)
	if err != nil {
		return err
	}

	if err := p.Insert(section, parentName, "acl", SerializeACL(*data), int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

// EditACL edits a ACL line in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) EditACL(id int64, parentType string, parentName string, data *models.ACL, transactionID string, version int64) error {
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

	section, parentName, err := getParserFromParent("acl", parentType, parentName)
	if err != nil {
		return err
	}

	if _, err := p.GetOne(section, parentName, "acl", int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	if err := p.Set(section, parentName, "acl", SerializeACL(*data), int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

// ReplaceAcls replaces all ACL lines in configuration for a parentType/parentName.
// One of version or transactionID is mandatory.
// Returns error on fail, nil on success.
func (c *client) ReplaceAcls(parentType string, parentName string, data models.Acls, transactionID string, version int64) error {
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

	section, parentName, err := getParserFromParent("acl", parentType, parentName)
	if err != nil {
		return err
	}

	acls, err := ParseACLs(section, parentName, p)
	if err != nil {
		return c.HandleError("", parentType, parentName, "", false, err)
	}

	for i := range acls {
		// Always delete index 0
		if err := p.Delete(section, parentName, "acl", 0); err != nil {
			return c.HandleError(strconv.FormatInt(int64(i), 10), parentType, parentName, t, transactionID == "", err)
		}
	}

	for i, newACL := range data {
		if err := p.Insert(section, parentName, "acl", SerializeACL(*newACL), i); err != nil {
			return c.HandleError(strconv.FormatInt(int64(i), 10), parentType, parentName, t, transactionID == "", err)
		}
	}

	return c.SaveData(p, t, transactionID == "")
}

func ParseACLs(section parser.Section, name string, p parser.Parser, aclName ...string) (models.Acls, error) {
	acls := models.Acls{}
	data, err := p.Get(section, name, "acl", false)
	if err != nil {
		if errors.Is(err, parser_errors.ErrFetch) {
			return acls, nil
		}
		return nil, err
	}

	aclLines, ok := data.([]types.ACL)
	if !ok {
		return nil, errors.New("type assert error []types.ACL")
	}
	lACL := len(aclName)
	for _, r := range aclLines {
		acl := ParseACL(r)
		if acl != nil {
			if lACL > 0 && aclName[0] == acl.ACLName {
				acls = append(acls, acl)
			} else if lACL == 0 {
				acls = append(acls, acl)
			}
		}
	}
	return acls, nil
}

func ParseACL(f types.ACL) *models.ACL {
	return &models.ACL{
		ACLName:   f.Name,
		Criterion: f.Criterion,
		Value:     f.Value,
	}
}

func SerializeACL(f models.ACL) types.ACL {
	return types.ACL{
		Name:      f.ACLName,
		Criterion: f.Criterion,
		Value:     f.Value,
	}
}
