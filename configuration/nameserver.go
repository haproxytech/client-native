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
	"fmt"
	"strconv"
	"strings"

	strfmt "github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/config-parser/v4"
	parser_errors "github.com/haproxytech/config-parser/v4/errors"
	"github.com/haproxytech/config-parser/v4/types"

	"github.com/haproxytech/client-native/v5/misc"
	"github.com/haproxytech/client-native/v5/models"
)

type Nameserver interface {
	GetNameservers(resolverSection string, transactionID string) (int64, models.Nameservers, error)
	GetNameserver(name string, resolverSection string, transactionID string) (int64, *models.Nameserver, error)
	DeleteNameserver(name string, resolverSection string, transactionID string, version int64) error
	CreateNameserver(resolverSection string, data *models.Nameserver, transactionID string, version int64) error
	EditNameserver(name string, resolverSection string, data *models.Nameserver, transactionID string, version int64) error
}

// GetNameservers returns configuration version and an array of
// configured namservers in the specified resolvers section. Returns error on fail.
func (c *client) GetNameservers(resolverSection string, transactionID string) (int64, models.Nameservers, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	nameservers, err := ParseNameservers(resolverSection, p)
	if err != nil {
		return v, nil, c.HandleError("", "resolvers", resolverSection, "", false, err)
	}

	return v, nameservers, nil
}

// GetNameserver returns configuration version and a requested nameserver
// in the specified resolvers section. Returns error on fail or if nameserver does not exist.
func (c *client) GetNameserver(name string, resolverSection string, transactionID string) (int64, *models.Nameserver, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	nameserver, _ := GetNameserverByName(name, resolverSection, p)
	if nameserver == nil {
		return v, nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Nameserver %s does not exist in resolvers section %s", name, resolverSection))
	}

	return v, nameserver, nil
}

// DeleteNameserver deletes an nameserver in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) DeleteNameserver(name string, resolverSection string, transactionID string, version int64) error {
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	nameserver, i := GetNameserverByName(name, resolverSection, p)
	if nameserver == nil {
		e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Nameserver %s does not exist in resolvers section %s", name, resolverSection))
		return c.HandleError(name, "resolvers", resolverSection, t, transactionID == "", e)
	}

	if err := p.Delete(parser.Resolvers, resolverSection, "nameserver", i); err != nil {
		return c.HandleError(name, "resolvers", resolverSection, t, transactionID == "", err)
	}

	if err := c.SaveData(p, t, transactionID == ""); err != nil {
		return err
	}
	return nil
}

// CreateNameserver creates a nameserver in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) CreateNameserver(resolverSection string, data *models.Nameserver, transactionID string, version int64) error {
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

	nameserver, _ := GetNameserverByName(data.Name, resolverSection, p)
	if nameserver != nil {
		e := NewConfError(ErrObjectAlreadyExists, fmt.Sprintf("Nameserver %s already exists in resolvers section %s", data.Name, resolverSection))
		return c.HandleError(data.Name, "resolvers", resolverSection, t, transactionID == "", e)
	}

	if err := p.Insert(parser.Resolvers, resolverSection, "nameserver", SerializeNameserver(*data), -1); err != nil {
		return c.HandleError(data.Name, "resolvers", resolverSection, t, transactionID == "", err)
	}

	if err := c.SaveData(p, t, transactionID == ""); err != nil {
		return err
	}

	return nil
}

// EditNameserver edits a nameserver in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) EditNameserver(name string, resolverSection string, data *models.Nameserver, transactionID string, version int64) error {
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

	nameserver, i := GetNameserverByName(name, resolverSection, p)
	if nameserver == nil {
		e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Nameserver %v does not exist in resolvers section %s", name, resolverSection))
		return c.HandleError(data.Name, "resolvers", resolverSection, t, transactionID == "", e)
	}

	if err := p.Set(parser.Resolvers, resolverSection, "nameserver", SerializeNameserver(*data), i); err != nil {
		return c.HandleError(data.Name, "resolvers", resolverSection, t, transactionID == "", err)
	}

	if err := c.SaveData(p, t, transactionID == ""); err != nil {
		return err
	}

	return nil
}

func ParseNameservers(resolverSection string, p parser.Parser) (models.Nameservers, error) {
	nameserver := models.Nameservers{}

	data, err := p.Get(parser.Resolvers, resolverSection, "nameserver", false)
	if err != nil {
		if errors.Is(err, parser_errors.ErrFetch) {
			return nameserver, nil
		}
		return nil, err
	}

	nameservers, ok := data.([]types.Nameserver)
	if !ok {
		return nil, misc.CreateTypeAssertError("nameserver")
	}
	for _, e := range nameservers {
		pe := ParseNameserver(e)
		if pe != nil {
			nameserver = append(nameserver, pe)
		}
	}
	return nameserver, nil
}

func ParseNameserver(p types.Nameserver) *models.Nameserver {
	parts := strings.Split(p.Address, ":")
	if len(parts) < 1 {
		return nil
	}
	ip := parts[0]
	if len(parts) > 1 {
		port, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return nil
		}
		return &models.Nameserver{
			Address: &ip,
			Port:    &port,
			Name:    p.Name,
		}
	}
	return &models.Nameserver{
		Address: &ip,
		Port:    nil,
		Name:    p.Name,
	}
}

func SerializeNameserver(pe models.Nameserver) types.Nameserver {
	addr := ""
	if pe.Address != nil {
		addr = *pe.Address
	}
	if pe.Port != nil {
		addr = fmt.Sprintf("%s:%d", addr, *pe.Port)
	}
	return types.Nameserver{
		Address: addr,
		Name:    pe.Name,
	}
}

func GetNameserverByName(name string, resolverSection string, p parser.Parser) (*models.Nameserver, int) {
	nameservers, err := ParseNameservers(resolverSection, p)
	if err != nil {
		return nil, 0
	}

	for i, b := range nameservers {
		if b.Name == name {
			return b, i
		}
	}
	return nil, 0
}
