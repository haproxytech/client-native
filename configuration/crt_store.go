// Copyright 2024 HAProxy Technologies
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
	"strings"

	strfmt "github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/client-native/v6/config-parser"
	parser_errors "github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/types"
	"github.com/haproxytech/client-native/v6/misc"
	"github.com/haproxytech/client-native/v6/models"
)

type CrtStore interface {
	GetCrtStores(transactionID string) (int64, models.CrtStores, error)
	GetCrtStore(name, transactionID string) (int64, *models.CrtStore, error)
	DeleteCrtStore(name, transactionID string, version int64) error
	CreateCrtStore(data *models.CrtStore, transactionID string, version int64) error
	EditCrtStore(name string, data *models.CrtStore, transactionID string, version int64) error
}

func (c *client) GetCrtStores(transactionID string) (int64, models.CrtStores, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	names, err := p.SectionsGet(parser.CrtStore)
	if err != nil {
		return v, nil, err
	}

	stores := make([]*models.CrtStore, len(names))

	for i, name := range names {
		_, store, err := c.GetCrtStore(name, transactionID)
		if err != nil {
			return v, nil, err
		}
		stores[i] = store
	}

	return v, stores, nil
}

func (c *client) GetCrtStore(name, transactionID string) (int64, *models.CrtStore, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	if !p.SectionExists(parser.CrtStore, name) {
		return v, nil, NewConfError(ErrObjectDoesNotExist,
			fmt.Sprintf("%s section '%s' does not exist", CrtStoreParentName, name))
	}

	store, err := ParseCrtStore(p, name)
	if err != nil {
		return 0, nil, err
	}

	return v, store, nil
}

func (c *client) DeleteCrtStore(name, transactionID string, version int64) error {
	return c.deleteSection(parser.CrtStore, name, transactionID, version)
}

func (c *client) CreateCrtStore(data *models.CrtStore, transactionID string, version int64) error {
	if c.UseModelsValidation {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return NewConfError(ErrValidationError, validationErr.Error())
		}
	}

	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return c.HandleError(data.Name, "", "", t, transactionID == "", err)
	}

	if p.SectionExists(parser.CrtStore, data.Name) {
		e := NewConfError(ErrObjectAlreadyExists, fmt.Sprintf("%s %s already exists", parser.CrtStore, data.Name))
		return c.HandleError(data.Name, "", "", t, transactionID == "", e)
	}

	if err = p.SectionsCreate(parser.CrtStore, data.Name); err != nil {
		return c.HandleError(data.Name, "", "", t, transactionID == "", err)
	}

	if err = SerializeCrtStore(p, data); err != nil {
		return c.HandleError(data.Name, "", "", t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

func (c *client) EditCrtStore(name string, data *models.CrtStore, transactionID string, version int64) error { //nolint:revive
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

	if !p.SectionExists(parser.CrtStore, data.Name) {
		e := NewConfError(ErrObjectAlreadyExists, fmt.Sprintf("%s %s does not exists", parser.CrtStore, data.Name))
		return c.HandleError(data.Name, "", "", t, transactionID == "", e)
	}

	if err = SerializeCrtStore(p, data); err != nil {
		return err
	}

	return c.SaveData(p, t, transactionID == "")
}

func ParseCrtStore(p parser.Parser, name string) (*models.CrtStore, error) {
	store := &models.CrtStore{Name: name}

	if data, err := p.SectionGet(parser.CrtStore, name); err == nil {
		d, ok := data.(types.Section)
		if ok {
			store.Metadata = misc.ParseMetadata(d.Comment)
		}
	}

	// get optional crt-base
	crtBase, err := p.Get(parser.CrtStore, name, "crt-base", false)
	if err != nil {
		if errors.Is(err, parser_errors.ErrFetch) {
			return store, nil
		}
		return nil, err
	}
	sc, ok := crtBase.(*types.StringC)
	if !ok {
		return nil, misc.CreateTypeAssertError("crt-base")
	}
	store.CrtBase = sc.Value

	// get optional key-base
	keyBase, err := p.Get(parser.CrtStore, name, "key-base", false)
	if err != nil {
		if errors.Is(err, parser_errors.ErrFetch) {
			return store, nil
		}
		return nil, err
	}
	sc, ok = keyBase.(*types.StringC)
	if !ok {
		return nil, misc.CreateTypeAssertError("key-base")
	}
	store.KeyBase = sc.Value

	// get optional loads
	loads, err := p.Get(parser.CrtStore, name, "load", false)
	if err != nil {
		if errors.Is(err, parser_errors.ErrFetch) {
			return store, nil
		}
		return nil, err
	}
	tloads, ok := loads.([]types.LoadCert)
	if !ok {
		return nil, misc.CreateTypeAssertError("load crt")
	}
	store.Loads = make(models.CrtLoads, len(tloads))
	for i, l := range tloads {
		domains := strings.Split(l.Domains, ",")
		if len(domains) == 1 && domains[0] == "" {
			domains = nil
		}
		mload := &models.CrtLoad{
			Acme:        l.Acme,
			Alias:       l.Alias,
			Certificate: l.Certificate,
			Domains:     domains,
			Issuer:      l.Issuer,
			Key:         l.Key,
			Ocsp:        l.Ocsp,
			Sctl:        l.Sctl,
			Metadata:    misc.ParseMetadata(l.Comment),
		}
		if l.OcspUpdate != nil {
			if *l.OcspUpdate {
				mload.OcspUpdate = models.CrtLoadOcspUpdateEnabled
			} else {
				mload.OcspUpdate = models.CrtLoadOcspUpdateDisabled
			}
		}
		store.Loads[i] = mload
	}

	return store, nil
}

func SerializeCrtStore(p parser.Parser, store *models.CrtStore) error {
	if store == nil {
		return fmt.Errorf("empty %s section", CrtStoreParentName)
	}

	if store.Metadata != nil {
		comment, err := misc.SerializeMetadata(store.Metadata)
		if err != nil {
			return err
		}
		if err := p.SectionCommentSet(parser.CrtStore, store.Name, comment); err != nil {
			return err
		}
	}

	crtBase := types.StringC{Value: store.CrtBase}
	if err := p.Set(parser.CrtStore, store.Name, "crt-base", crtBase); err != nil {
		return err
	}

	keyBase := types.StringC{Value: store.KeyBase}
	if err := p.Set(parser.CrtStore, store.Name, "key-base", keyBase); err != nil {
		return err
	}

	if len(store.Loads) > 0 {
		for i, load := range store.Loads {
			err := p.Insert(parser.CrtStore, store.Name, "load", SerializeCrtLoad(load), i)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
