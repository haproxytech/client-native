// Copyright 2025 HAProxy Technologies
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
	"sort"

	"github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/client-native/v6/config-parser"
	convert "github.com/haproxytech/client-native/v6/configuration/convert/v2v3"
	"github.com/haproxytech/client-native/v6/models"
)

type StructuredCrtStore interface {
	GetStructuredCrtStores(transactionID string) (int64, models.CrtStores, error)
	GetStructuredCrtStore(name string, transactionID string) (int64, *models.CrtStore, error)
	EditStructuredCrtStore(name string, data *models.CrtStore, transactionID string, version int64) error
	CreateStructuredCrtStore(data *models.CrtStore, transactionID string, version int64) error
}

func (c *client) GetStructuredCrtStore(name string, transactionID string) (int64, *models.CrtStore, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	if !p.SectionExists(parser.CrtStore, name) {
		return v, nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("CrtStore '%s' does not exist", name))
	}

	capt, err := parseCrtStoreSection(name, p)

	return v, capt, err
}

func (c *client) GetStructuredCrtStores(transactionID string) (int64, models.CrtStores, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	crtStores, err := parseCrtStoreSections(p)
	if err != nil {
		return 0, nil, err
	}

	return v, crtStores, nil
}

func (c *client) EditStructuredCrtStore(name string, data *models.CrtStore, transactionID string, version int64) error {
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

	if !p.SectionExists(parser.CrtStore, name) {
		e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("%s %s does not exist", parser.CrtStore, name))
		return c.HandleError(name, "", "", t, transactionID == "", e)
	}

	if err = p.SectionsDelete(parser.CrtStore, name); err != nil {
		return c.HandleError(name, "", "", t, transactionID == "", err)
	}

	if err = serializeCrtStoreSection(StructuredToParserArgs{
		TID:         transactionID,
		Parser:      &p,
		Options:     &c.ConfigurationOptions,
		HandleError: c.HandleError,
	}, data); err != nil {
		return err
	}
	return c.SaveData(p, t, transactionID == "")
}

func (c *client) CreateStructuredCrtStore(data *models.CrtStore, transactionID string, version int64) error {
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

	if p.SectionExists(parser.CrtStore, data.Name) {
		e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("%s %s already exist", parser.CrtStore, data.Name))
		return c.HandleError(data.Name, "", "", t, transactionID == "", e)
	}

	if err = serializeCrtStoreSection(StructuredToParserArgs{
		TID:         transactionID,
		Parser:      &p,
		Options:     &c.ConfigurationOptions,
		HandleError: c.HandleError,
	}, data); err != nil {
		return err
	}
	return c.SaveData(p, t, transactionID == "")
}

func serializeCrtStoreSection(a StructuredToParserArgs, r *models.CrtStore) error {
	p := *a.Parser

	if err := p.SectionsCreate(parser.CrtStore, r.Name); err != nil {
		return err
	}
	if err := SerializeCrtStore(p, r); err != nil {
		return err
	}
	crtLoadCerts := make([]string, 0, len(r.CrtLoads))
	for name := range r.CrtLoads {
		crtLoadCerts = append(crtLoadCerts, name)
	}
	sort.Strings(crtLoadCerts)
	for _, name := range crtLoadCerts {
		crtLoad := r.CrtLoads[name]
		if err := p.Insert(parser.CrtStore, r.Name, "load", SerializeCrtLoad(&crtLoad), -1); err != nil {
			return a.HandleError(crtLoad.Certificate, CrtStoreParentName, r.Name, a.TID, a.TID == "", err)
		}
	}
	return nil
}

func parseCrtStoreSections(p parser.Parser) (models.CrtStores, error) {
	names, err := p.SectionsGet(parser.CrtStore)
	if err != nil {
		return nil, err
	}
	crtStores := []*models.CrtStore{}
	for _, name := range names {
		c, err := parseCrtStoreSection(name, p)
		if err != nil {
			return nil, err
		}
		crtStores = append(crtStores, c)
	}
	return crtStores, nil
}

func parseCrtStoreSection(name string, p parser.Parser) (*models.CrtStore, error) {
	store := &models.CrtStore{CrtStoreBase: models.CrtStoreBase{Name: name}}
	if err := ParseCrtStore(p, store); err != nil {
		return nil, err
	}

	// crtLoads
	loads, err := ParseCrtLoads(name, p)
	if err != nil {
		return nil, err
	}
	crtLoads, err := convert.NamedResourceArrayToMap(loads)
	if err != nil {
		return nil, err
	}
	store.CrtLoads = crtLoads
	return store, nil
}
