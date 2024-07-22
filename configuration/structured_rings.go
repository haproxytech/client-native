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

	"github.com/go-openapi/strfmt"
	"github.com/haproxytech/client-native/v6/models"
	parser "github.com/haproxytech/config-parser/v5"
)

type StructuredRing interface {
	GetStructuredRings(transactionID string) (int64, models.Rings, error)
	GetStructuredRing(name string, transactionID string) (int64, *models.Ring, error)
	CreateStructuredRing(data *models.Ring, transactionID string, version int64) error
	EditStructuredRing(name string, data *models.Ring, transactionID string, version int64) error
}

// GetStructuredRing returns configuration version and a requested ring with all its child resources.
// Returns error on fail or if ring does not exist.
func (c *client) GetStructuredRing(name string, transactionID string) (int64, *models.Ring, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	if !c.checkSectionExists(parser.Ring, name, p) {
		return v, nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Ring %s does not exist", name))
	}

	f, err := parseRingsSection(name, p)

	return v, f, err
}

func (c *client) GetStructuredRings(transactionID string) (int64, models.Rings, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	rings, err := parseRingsSections(p)
	if err != nil {
		return 0, nil, err
	}

	return v, rings, nil
}

// EditStructuredRing replaces a ring and all it's child resources in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) EditStructuredRing(name string, data *models.Ring, transactionID string, version int64) error {
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

	if !c.checkSectionExists(parser.Ring, name, p) {
		e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("%s %s does not exist", parser.Ring, name))
		return c.HandleError(name, "", "", t, transactionID == "", e)
	}

	if err = p.SectionsDelete(parser.Ring, name); err != nil {
		return c.HandleError(name, "", "", t, transactionID == "", err)
	}

	if err = serializeRingSection(StructuredToParserArgs{
		TID:                transactionID,
		Parser:             &p,
		Options:            &c.ConfigurationOptions,
		HandleError:        c.HandleError,
		CheckSectionExists: c.checkSectionExists,
	}, data); err != nil {
		return err
	}
	return c.SaveData(p, t, transactionID == "")
}

// CreateStructuredRing creates a ring and all it's child resources in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) CreateStructuredRing(data *models.Ring, transactionID string, version int64) error {
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

	if c.checkSectionExists(parser.Ring, data.Name, p) {
		e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("%s %s already exist", parser.Ring, data.Name))
		return c.HandleError(data.Name, "", "", t, transactionID == "", e)
	}

	if err = serializeRingSection(StructuredToParserArgs{
		TID:                transactionID,
		Parser:             &p,
		Options:            &c.ConfigurationOptions,
		HandleError:        c.HandleError,
		CheckSectionExists: c.checkSectionExists,
	}, data); err != nil {
		return err
	}
	return c.SaveData(p, t, transactionID == "")
}

func parseRingsSections(p parser.Parser) (models.Rings, error) {
	names, err := p.SectionsGet(parser.Ring)
	if err != nil {
		return nil, err
	}
	rings := []*models.Ring{}
	for _, name := range names {
		f, err := parseRingsSection(name, p)
		if err != nil {
			return nil, err
		}
		rings = append(rings, f)
	}
	return rings, nil
}

func parseRingsSection(name string, p parser.Parser) (*models.Ring, error) {
	r := &models.Ring{RingBase: models.RingBase{Name: name}}
	if err := ParseRingSection(p, r); err != nil {
		return nil, err
	}
	// nameservers
	servers, err := ParseServers(RingParentName, name, p)
	if err != nil {
		return nil, err
	}
	serversa, errsa := namedResourceArrayToMap(servers)
	if errsa != nil {
		return nil, errsa
	}
	r.Servers = serversa
	return r, nil
}

func serializeRingSection(a StructuredToParserArgs, r *models.Ring) error {
	p := *a.Parser
	var err error
	err = p.SectionsCreate(parser.Ring, r.Name)
	if err != nil {
		return err
	}
	if err = SerializeRingSection(p, r, a.Options); err != nil {
		return err
	}
	for _, ns := range r.Servers {
		if err = p.Insert(parser.Ring, r.Name, "server", SerializeServer(ns, a.Options), -1); err != nil {
			return a.HandleError(ns.Name, RingParentName, r.Name, a.TID, a.TID == "", err)
		}
	}

	return nil
}
