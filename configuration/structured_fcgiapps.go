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
	"strconv"

	"github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/client-native/v6/config-parser"
	"github.com/haproxytech/client-native/v6/models"
)

type StructuredFCGIApp interface {
	GetStructuredFCGIApplications(transactionID string) (int64, models.FCGIApps, error)
	GetStructuredFCGIApplication(name string, transactionID string) (int64, *models.FCGIApp, error)
	EditStructuredFCGIApplication(name string, data *models.FCGIApp, transactionID string, version int64) error
	CreateStructuredFCGIApplication(data *models.FCGIApp, transactionID string, version int64) error
}

// GetStructuredFCGIApp returns configuration version and a requested fcgiapp with all its child resources.
// Returns error on fail or if fcgiapp does not exist.
func (c *client) GetStructuredFCGIApplication(name string, transactionID string) (int64, *models.FCGIApp, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	if !c.checkSectionExists(parser.FCGIApp, name, p) {
		return v, nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("FCGIApp %s does not exist", name))
	}

	f, err := parseFCGIAppsSection(name, p)

	return v, f, err
}

func (c *client) GetStructuredFCGIApplications(transactionID string) (int64, models.FCGIApps, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	fcgiapps, err := parseFCGIAppsSections(p)
	if err != nil {
		return 0, nil, err
	}

	return v, fcgiapps, nil
}

// EditStructuredFCGIApp replaces a fcgiapp and all it's child resources in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) EditStructuredFCGIApplication(name string, data *models.FCGIApp, transactionID string, version int64) error {
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

	if !c.checkSectionExists(parser.FCGIApp, name, p) {
		e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("%s %s does not exist", parser.FCGIApp, name))
		return c.HandleError(name, "", "", t, transactionID == "", e)
	}

	if err = p.SectionsDelete(parser.FCGIApp, name); err != nil {
		return c.HandleError(name, "", "", t, transactionID == "", err)
	}

	if err = serializeFCGIAppSection(StructuredToParserArgs{
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

// CreateStructuredFCGIApp creates a fcgiapp and all it's child resources in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) CreateStructuredFCGIApplication(data *models.FCGIApp, transactionID string, version int64) error {
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

	if c.checkSectionExists(parser.FCGIApp, data.Name, p) {
		e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("%s %s already exist", parser.FCGIApp, data.Name))
		return c.HandleError(data.Name, "", "", t, transactionID == "", e)
	}

	if err = serializeFCGIAppSection(StructuredToParserArgs{
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

func parseFCGIAppsSections(p parser.Parser) (models.FCGIApps, error) {
	names, err := p.SectionsGet(parser.FCGIApp)
	if err != nil {
		return nil, err
	}
	rings := []*models.FCGIApp{}
	for _, name := range names {
		f, err := parseFCGIAppsSection(name, p)
		if err != nil {
			return nil, err
		}
		rings = append(rings, f)
	}
	return rings, nil
}

func parseFCGIAppsSection(name string, p parser.Parser) (*models.FCGIApp, error) {
	f, err := ParseFCGIApp(p, name)
	if err != nil {
		return nil, err
	}

	acls, err := ParseACLs(FCGIAppParentName, name, p)
	if err != nil {
		return nil, err
	}
	f.ACLList = acls

	return f, nil
}

func serializeFCGIAppSection(a StructuredToParserArgs, f *models.FCGIApp) error {
	p := *a.Parser
	var err error
	err = p.SectionsCreate(parser.FCGIApp, f.Name)
	if err != nil {
		return err
	}
	if err = SerializeFCGIAppSection(p, f); err != nil {
		return err
	}
	for i, acl := range f.ACLList {
		if err = p.Insert(parser.FCGIApp, f.Name, "acl", SerializeACL(*acl), i); err != nil {
			return a.HandleError(strconv.FormatInt(int64(i), 10), "fcgi", f.Name, a.TID, a.TID == "", err)
		}
	}

	return nil
}
