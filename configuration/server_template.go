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

	"github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/config-parser/v5"
	parser_errors "github.com/haproxytech/config-parser/v5/errors"
	"github.com/haproxytech/config-parser/v5/params"
	"github.com/haproxytech/config-parser/v5/types"

	"github.com/haproxytech/client-native/v5/misc"
	"github.com/haproxytech/client-native/v5/models"
)

type ServerTemplate interface {
	GetServerTemplates(backend string, transactionID string) (int64, models.ServerTemplates, error)
	GetServerTemplate(prefix string, backend string, transactionID string) (int64, *models.ServerTemplate, error)
	DeleteServerTemplate(prefix string, backend string, transactionID string, version int64) error
	CreateServerTemplate(backend string, data *models.ServerTemplate, transactionID string, version int64) error
	EditServerTemplate(prefix string, backend string, data *models.ServerTemplate, transactionID string, version int64) error
}

// GetServerTemplates returns configuration version and an array of
// configured server templates in the specified backend. Returns error on fail.
func (c *client) GetServerTemplates(backend string, transactionID string) (int64, models.ServerTemplates, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	templates, err := ParseServerTemplates(backend, p)
	if err != nil {
		return v, nil, c.HandleError("", "backend", backend, "", false, err)
	}

	return v, templates, nil
}

// GetServerTemplate returns configuration version and a requested server template
// in the specified backend. Returns error on fail or if server template does not exist.
func (c *client) GetServerTemplate(prefix string, backend string, transactionID string) (int64, *models.ServerTemplate, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	template, _ := GetServerTemplateByPrefix(prefix, backend, p)
	if template == nil {
		return v, nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Server template %s does not exist in backend %s", prefix, backend))
	}

	return v, template, nil
}

// DeleteServerTemplate deletes a server template in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) DeleteServerTemplate(prefix string, backend string, transactionID string, version int64) error {
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	template, i := GetServerTemplateByPrefix(prefix, backend, p)
	if template == nil {
		e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Server template %s does not exist in backend %s", prefix, backend))
		return c.HandleError(prefix, "backend", backend, t, transactionID == "", e)
	}

	if err := p.Delete(parser.Backends, backend, "server-template", i); err != nil {
		return c.HandleError(prefix, "backend", backend, t, transactionID == "", err)
	}

	if err := c.SaveData(p, t, transactionID == ""); err != nil {
		return err
	}

	return nil
}

// CreateServerTemplate creates a server template in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) CreateServerTemplate(backend string, data *models.ServerTemplate, transactionID string, version int64) error {
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

	template, _ := GetServerTemplateByPrefix(data.Prefix, backend, p)
	if template != nil {
		e := NewConfError(ErrObjectAlreadyExists, fmt.Sprintf("Server template %s already exists in backend %s", data.Prefix, backend))
		return c.HandleError(data.Prefix, "backend", backend, t, transactionID == "", e)
	}

	if err := p.Insert(parser.Backends, backend, "server-template", SerializeServerTemplate(*data), -1); err != nil {
		return c.HandleError(data.Prefix, "backend", backend, t, transactionID == "", err)
	}

	if err := c.SaveData(p, t, transactionID == ""); err != nil {
		return err
	}
	return nil
}

// EditServerTemplate edits a server template in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) EditServerTemplate(prefix string, backend string, data *models.ServerTemplate, transactionID string, version int64) error {
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

	template, i := GetServerTemplateByPrefix(prefix, backend, p)
	if template == nil {
		e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Server template %v does not exist in backend %s", prefix, backend))
		return c.HandleError(data.Prefix, "backend", backend, t, transactionID == "", e)
	}

	if err := p.Set(parser.Backends, backend, "server-template", SerializeServerTemplate(*data), i); err != nil {
		return c.HandleError(data.Prefix, "backend", backend, t, transactionID == "", err)
	}

	if err := c.SaveData(p, t, transactionID == ""); err != nil {
		return err
	}
	return nil
}

func ParseServerTemplates(backend string, p parser.Parser) (models.ServerTemplates, error) {
	templates := models.ServerTemplates{}

	data, err := p.Get(parser.Backends, backend, "server-template", false)
	if err != nil {
		if errors.Is(err, parser_errors.ErrFetch) {
			return templates, nil
		}
		return nil, err
	}

	ondiskServerTemplates, ok := data.([]types.ServerTemplate)
	if !ok {
		return nil, misc.CreateTypeAssertError("server-template")
	}
	for _, ondiskServerTemplate := range ondiskServerTemplates {
		template := ParseServerTemplate(ondiskServerTemplate)
		if template != nil {
			templates = append(templates, template)
		}
	}
	return templates, nil
}

func ParseServerTemplate(ondiskServerTemplate types.ServerTemplate) *models.ServerTemplate {
	template := &models.ServerTemplate{
		Prefix:     ondiskServerTemplate.Prefix,
		NumOrRange: ondiskServerTemplate.NumOrRange,
		Fqdn:       ondiskServerTemplate.Fqdn,
		Port:       &ondiskServerTemplate.Port,
	}
	parseServerParams(ondiskServerTemplate.Params, &template.ServerParams)
	return template
}

func SerializeServerTemplate(s models.ServerTemplate) types.ServerTemplate {
	template := types.ServerTemplate{
		Prefix:     s.Prefix,
		NumOrRange: s.NumOrRange,
		Fqdn:       s.Fqdn,
		Params:     []params.ServerOption{},
	}
	if s.Port != nil {
		template.Port = *s.Port
	}
	template.Params = serializeServerParams(s.ServerParams)
	return template
}

func GetServerTemplateByPrefix(prefix string, backend string, p parser.Parser) (*models.ServerTemplate, int) {
	templates, err := ParseServerTemplates(backend, p)
	if err != nil {
		return nil, 0
	}
	for i, template := range templates {
		if template.Prefix == prefix {
			return template, i
		}
	}
	return nil, 0
}
