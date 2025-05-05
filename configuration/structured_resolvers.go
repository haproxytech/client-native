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
	parser "github.com/haproxytech/client-native/v6/config-parser"
	"github.com/haproxytech/client-native/v6/models"
)

type StructuredResolver interface {
	GetStructuredResolvers(transactionID string) (int64, models.Resolvers, error)
	GetStructuredResolver(name string, transactionID string) (int64, *models.Resolver, error)
	EditStructuredResolver(name string, data *models.Resolver, transactionID string, version int64) error
	CreateStructuredResolver(data *models.Resolver, transactionID string, version int64) error
}

// GetStructuredResolver returns configuration version and a requested resolver with all its child resources.
// Returns error on fail or if resolver does not exist.
func (c *client) GetStructuredResolver(name string, transactionID string) (int64, *models.Resolver, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	if !p.SectionExists(parser.Resolvers, name) {
		return v, nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Resolver %s does not exist", name))
	}

	f, err := parseResolversSection(name, p)

	return v, f, err
}

func (c *client) GetStructuredResolvers(transactionID string) (int64, models.Resolvers, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	resolvers, err := parseResolversSections(p)
	if err != nil {
		return 0, nil, err
	}

	return v, resolvers, nil
}

// EditStructuredResolver replaces a resolver and all it's child resources in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) EditStructuredResolver(name string, data *models.Resolver, transactionID string, version int64) error {
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

	if !p.SectionExists(parser.Resolvers, name) {
		e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("%s %s does not exist", parser.Resolvers, name))
		return c.HandleError(name, "", "", t, transactionID == "", e)
	}

	if err = p.SectionsDelete(parser.Resolvers, name); err != nil {
		return c.HandleError(name, "", "", t, transactionID == "", err)
	}

	if err = serializeResolverSection(StructuredToParserArgs{
		TID:         transactionID,
		Parser:      &p,
		Options:     &c.ConfigurationOptions,
		HandleError: c.HandleError,
	}, data); err != nil {
		return err
	}
	return c.SaveData(p, t, transactionID == "")
}

// CreateStructuredResolver creates a resolver and all it's child resources in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) CreateStructuredResolver(data *models.Resolver, transactionID string, version int64) error {
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

	if p.SectionExists(parser.Resolvers, data.Name) {
		e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("%s %s already exist", parser.Resolvers, data.Name))
		return c.HandleError(data.Name, "", "", t, transactionID == "", e)
	}

	if err = serializeResolverSection(StructuredToParserArgs{
		TID:         transactionID,
		Parser:      &p,
		Options:     &c.ConfigurationOptions,
		HandleError: c.HandleError,
	}, data); err != nil {
		return err
	}
	return c.SaveData(p, t, transactionID == "")
}

func parseResolversSections(p parser.Parser) (models.Resolvers, error) {
	names, err := p.SectionsGet(parser.Resolvers)
	if err != nil {
		return nil, err
	}
	resolvers := []*models.Resolver{}
	for _, name := range names {
		f, err := parseResolversSection(name, p)
		if err != nil {
			return nil, err
		}
		resolvers = append(resolvers, f)
	}
	return resolvers, nil
}

func parseResolversSection(name string, p parser.Parser) (*models.Resolver, error) {
	r := &models.Resolver{ResolverBase: models.ResolverBase{Name: name}}
	if err := ParseResolverSection(p, r); err != nil {
		return nil, err
	}
	// nameservers
	ns, err := ParseNameservers(name, p)
	if err != nil {
		return nil, err
	}
	nsa, errNsa := namedResourceArrayToMap[models.Nameserver](ns)
	if errNsa != nil {
		return nil, errNsa
	}
	r.Nameservers = nsa
	return r, nil
}

func serializeResolverSection(a StructuredToParserArgs, r *models.Resolver) error {
	p := *a.Parser
	var err error
	err = p.SectionsCreate(parser.Resolvers, r.Name)
	if err != nil {
		return err
	}
	if err = SerializeResolverSection(p, r, a.Options); err != nil {
		return err
	}
	for _, ns := range r.Nameservers {
		if err = p.Insert(parser.Resolvers, r.Name, "nameserver", SerializeNameserver(ns), -1); err != nil {
			return a.HandleError(ns.Name, ResolverParentName, r.Name, a.TID, a.TID == "", err)
		}
	}

	return nil
}
