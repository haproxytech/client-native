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
	"strings"

	strfmt "github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/client-native/v6/config-parser"
	parser_errors "github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/types"
	"github.com/haproxytech/client-native/v6/models"
)

type CrtLoad interface {
	GetCrtLoads(crtStore, transactionID string) (int64, models.CrtLoads, error)
	GetCrtLoad(certificate, crtStore, transactionID string) (int64, *models.CrtLoad, error)
	DeleteCrtLoad(certificate, crtStore, transactionID string, version int64) error
	CreateCrtLoad(crtStore string, data *models.CrtLoad, transactionID string, version int64) error
	EditCrtLoad(certificate, crtStore string, data *models.CrtLoad, transactionID string, version int64) error
}

func (c *client) GetCrtLoads(crtStore, transactionID string) (int64, models.CrtLoads, error) {
	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	_, store, err := c.GetCrtStore(crtStore, transactionID)
	if err != nil {
		return v, nil, c.HandleError("", CrtStoreParentName, crtStore, transactionID, transactionID == "", err)
	}

	return v, store.Loads, nil
}

func (c *client) GetCrtLoad(certificate, crtStore, transactionID string) (int64, *models.CrtLoad, error) {
	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	_, store, err := c.GetCrtStore(crtStore, transactionID)
	if err != nil {
		return v, nil, c.HandleError("", CrtStoreParentName, crtStore, transactionID, transactionID == "", err)
	}

	for _, load := range store.Loads {
		if load.Certificate == certificate {
			return v, load, nil
		}
	}

	return v, nil, c.HandleError(certificate, CrtStoreParentName, crtStore, transactionID, transactionID == "", parser_errors.ErrFetch)
}

func (c *client) DeleteCrtLoad(certificate, crtStore, transactionID string, version int64) error {
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	store, err := ParseCrtStore(p, crtStore)
	if err != nil {
		return c.HandleError(certificate, CrtStoreParentName, crtStore, t, transactionID == "", err)
	}

	for i, load := range store.Loads {
		if load.Certificate == certificate {
			err = p.Delete(parser.CrtStore, crtStore, "load", i)
			if err != nil {
				return c.HandleError(certificate, CrtStoreParentName, crtStore, t, transactionID == "", err)
			}
			return c.SaveData(p, t, transactionID == "")
		}
	}

	return c.HandleError(certificate, CrtStoreParentName, crtStore, t, transactionID == "", parser_errors.ErrFetch)
}

func (c *client) CreateCrtLoad(crtStore string, data *models.CrtLoad, transactionID string, version int64) error {
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

	store, err := ParseCrtStore(p, crtStore)
	if err != nil {
		return c.HandleError("", CrtStoreParentName, crtStore, t, transactionID == "", err)
	}

	for _, load := range store.Loads {
		if load.Certificate == data.Certificate {
			return c.HandleError(data.Certificate, CrtStoreParentName, crtStore, t, transactionID == "", parser_errors.ErrFetch)
		}
	}

	if err = p.Insert(parser.CrtStore, crtStore, "load", SerializeCrtLoad(data), -1); err != nil {
		return c.HandleError(data.Certificate, CrtStoreParentName, crtStore, t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

func (c *client) EditCrtLoad(certificate, crtStore string, data *models.CrtLoad, transactionID string, version int64) error {
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

	store, err := ParseCrtStore(p, crtStore)
	if err != nil {
		return c.HandleError("", CrtStoreParentName, crtStore, t, transactionID == "", err)
	}

	for i, load := range store.Loads {
		if load.Certificate == certificate {
			err = p.Set(parser.CrtStore, crtStore, "load", SerializeCrtLoad(data), i)
			if err != nil {
				return c.HandleError(certificate, CrtStoreParentName, crtStore, t, transactionID == "", err)
			}
			return c.SaveData(p, t, transactionID == "")
		}
	}

	return c.HandleError(certificate, CrtStoreParentName, crtStore, t, transactionID == "", parser_errors.ErrFetch)
}

func SerializeCrtLoad(load *models.CrtLoad) *types.LoadCert {
	t := &types.LoadCert{
		Acme:        load.Acme,
		Alias:       load.Alias,
		Certificate: load.Certificate,
		Domains:     strings.Join(load.Domains, ","),
		Key:         load.Key,
		Issuer:      load.Issuer,
		Ocsp:        load.Ocsp,
		Sctl:        load.Sctl,
	}
	comment, err := serializeMetadata(load.Metadata)
	if err == nil {
		t.Comment = comment
	}
	if load.OcspUpdate != "" {
		t.OcspUpdate = new(bool)
		if load.OcspUpdate == models.CrtLoadOcspUpdateEnabled {
			*t.OcspUpdate = true
		}
	}
	return t
}
