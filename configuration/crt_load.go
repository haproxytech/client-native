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

type CrtLoad interface {
	GetCrtLoads(crtStore, transactionID string) (int64, models.CrtLoads, error)
	GetCrtLoad(certificate, crtStore, transactionID string) (int64, *models.CrtLoad, error)
	DeleteCrtLoad(certificate, crtStore, transactionID string, version int64) error
	CreateCrtLoad(crtStore string, data *models.CrtLoad, transactionID string, version int64) error
	EditCrtLoad(certificate, crtStore string, data *models.CrtLoad, transactionID string, version int64) error
}

func (c *client) GetCrtLoads(crtStore, transactionID string) (int64, models.CrtLoads, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	loads, err := ParseCrtLoads(crtStore, p)
	if err != nil {
		return v, nil, c.HandleError("", "crt_store", crtStore, "", false, err)
	}

	return v, loads, nil
}

func (c *client) GetCrtLoad(certificate, crtStore, transactionID string) (int64, *models.CrtLoad, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}
	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	load, _ := GetCrtLoadByName(certificate, crtStore, p)
	if load == nil {
		return 0, nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("CrtLoad %v does not exist in %s crt_store", certificate, crtStore))
	}

	return v, load, nil
}

func (c *client) DeleteCrtLoad(certificate, crtStore, transactionID string, version int64) error {
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	load, i := GetCrtLoadByName(certificate, crtStore, p)
	if load == nil {
		e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("CrtLoad %v does not exist in %s crt_store", certificate, crtStore))
		return c.HandleError(certificate, crtStore, "crt_store", t, transactionID == "", e)
	}

	if err = p.Delete(parser.CrtStore, crtStore, "load", i); err != nil {
		return c.HandleError(certificate, CrtStoreParentName, crtStore, t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
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

	load, _ := GetCrtLoadByName(data.Certificate, crtStore, p)
	if load != nil {
		e := NewConfError(ErrObjectAlreadyExists, fmt.Sprintf("crtLoad %s already exists in %s crt_store", data.Certificate, crtStore))
		return c.HandleError(data.Certificate, "crt_store", crtStore, t, transactionID == "", e)
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

	load, i := GetCrtLoadByName(certificate, crtStore, p)
	if load == nil {
		e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("CrtLoad %v does not exist in %s crt_store", certificate, crtStore))
		return c.HandleError(data.Certificate, crtStore, "crt_store", t, transactionID == "", e)
	}

	if err := p.Set(parser.CrtStore, crtStore, "load", SerializeCrtLoad(data), i); err != nil {
		return c.HandleError(data.Certificate, crtStore, "crt_store", t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

func GetCrtLoadByName(certificate, crtStore string, p parser.Parser) (*models.CrtLoad, int) {
	loads, err := ParseCrtLoads(crtStore, p)
	if err != nil {
		return nil, 0
	}

	for i, l := range loads {
		if l.Certificate == certificate {
			return l, i
		}
	}
	return nil, 0
}

func ParseCrtLoads(parent string, p parser.Parser) (models.CrtLoads, error) {
	var crtLoads models.CrtLoads

	data, err := p.Get(parser.CrtStore, parent, "load", false)
	if err != nil {
		if errors.Is(err, parser_errors.ErrFetch) {
			return crtLoads, nil
		}
		return nil, err
	}

	ondiskCrtLoads, ok := data.([]types.LoadCert)
	if !ok {
		return nil, misc.CreateTypeAssertError("load crt")
	}

	for _, ondiskCrtLoad := range ondiskCrtLoads {
		c := ParseCrtLoad(ondiskCrtLoad)
		if c != nil {
			crtLoads = append(crtLoads, c)
		}
	}
	return crtLoads, nil
}

func ParseCrtLoad(load types.LoadCert) *models.CrtLoad {
	domains := strings.Split(load.Domains, ",")
	if len(domains) == 1 && domains[0] == "" {
		domains = nil
	}

	l := &models.CrtLoad{
		Acme:        load.Acme,
		Alias:       load.Alias,
		Certificate: load.Certificate,
		Domains:     domains,
		Issuer:      load.Issuer,
		Key:         load.Key,
		Ocsp:        load.Ocsp,
		Sctl:        load.Sctl,
		Metadata:    misc.ParseMetadata(load.Comment),
	}
	if load.OcspUpdate != nil {
		if *load.OcspUpdate {
			l.OcspUpdate = models.CrtLoadOcspUpdateEnabled
		} else {
			l.OcspUpdate = models.CrtLoadOcspUpdateDisabled
		}
	}
	return l
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
	comment, err := misc.SerializeMetadata(load.Metadata)
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
