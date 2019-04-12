package configuration

import (
	"strconv"

	strfmt "github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/config-parser"
	parser_errors "github.com/haproxytech/config-parser/errors"
	"github.com/haproxytech/config-parser/types"
	"github.com/haproxytech/models"
)

// GetACLs returns a struct with configuration version and an array of
// configured ACL lines in the specified parent. Returns error on fail.
func (c *Client) GetACLs(parentType, parentName string, transactionID string) (*models.GetAclsOKBody, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return nil, err
	}

	acls, err := c.parseACLs(parentType, parentName, p)
	if err != nil {
		return nil, c.handleError("", parentType, parentName, "", false, err)
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return nil, err
	}

	return &models.GetAclsOKBody{Version: v, Data: acls}, nil
}

// GetACL returns a struct with configuration version and a requested ACL line
// in the specified parent. Returns error on fail or if ACL line does not exist.
func (c *Client) GetACL(id int64, parentType, parentName string, transactionID string) (*models.GetACLOKBody, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return nil, err
	}

	var section parser.Section
	if parentType == "backend" {
		section = parser.Backends
	} else if parentType == "frontend" {
		section = parser.Frontends
	}

	data, err := p.GetOne(section, parentName, "acl", int(id))
	if err != nil {
		return nil, c.handleError(strconv.FormatInt(id, 10), parentType, parentName, "", false, err)
	}

	acl := parseACL(data.(types.Acl))
	acl.ID = &id

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return nil, err
	}

	return &models.GetACLOKBody{Version: v, Data: acl}, nil
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
		if err == parser_errors.FetchError {
			return acls, nil
		}
		return nil, err
	}

	aclLines := data.([]types.Acl)
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

func parseACL(f types.Acl) *models.ACL {
	return &models.ACL{
		ACLName:   &f.Name,
		Criterion: &f.Criterion,
		Value:     &f.Value,
	}
}

func serializeACL(f models.ACL) types.Acl {
	acl := types.Acl{}
	if f.ACLName != nil {
		acl.Name = *f.ACLName
	}
	if f.Criterion != nil {
		acl.Criterion = *f.Criterion
	}
	if f.Value != nil {
		acl.Value = *f.Value
	}
	return acl
}
