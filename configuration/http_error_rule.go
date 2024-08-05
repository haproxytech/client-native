// Copyright 2022 HAProxy Technologies
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
	goerrors "errors"
	"fmt"
	"strconv"

	strfmt "github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/client-native/v6/config-parser"
	parser_errors "github.com/haproxytech/client-native/v6/config-parser/errors"
	http_actions "github.com/haproxytech/client-native/v6/config-parser/parsers/http/actions"
	"github.com/haproxytech/client-native/v6/config-parser/types"

	"github.com/haproxytech/client-native/v6/misc"
	"github.com/haproxytech/client-native/v6/models"
)

type HTTPErrorRule interface {
	GetHTTPErrorRules(parentType, parentName string, transactionID string) (int64, models.HTTPErrorRules, error)
	GetHTTPErrorRule(id int64, parentType, parentName string, transactionID string) (int64, *models.HTTPErrorRule, error)
	DeleteHTTPErrorRule(id int64, parentType string, parentName string, transactionID string, version int64) error
	CreateHTTPErrorRule(id int64, parentType string, parentName string, data *models.HTTPErrorRule, transactionID string, version int64) error
	EditHTTPErrorRule(id int64, parentType string, parentName string, data *models.HTTPErrorRule, transactionID string, version int64) error
	ReplaceHTTPErrorRules(parentType string, parentName string, data models.HTTPErrorRules, transactionID string, version int64) error
}

// GetHTTPErrorRules returns configuration version and an array of
// configured http error rules in the specified parent. Returns error on fail.
func (c *client) GetHTTPErrorRules(parentType, parentName string, transactionID string) (int64, models.HTTPErrorRules, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	httpRules, err := ParseHTTPErrorRules(parentType, parentName, p)
	if err != nil {
		return v, nil, c.HandleError("", parentType, parentName, "", false, err)
	}

	return v, httpRules, nil
}

// GetHTTPErrorRule returns configuration version and a http error rule
// in the specified parent. Returns error on fail or if http error rule does not exist.
func (c *client) GetHTTPErrorRule(id int64, parentType, parentName string, transactionID string) (int64, *models.HTTPErrorRule, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	var section parser.Section
	switch parentType {
	case DefaultsParentName:
		section = parser.Defaults
		if parentName == "" {
			parentName = parser.DefaultSectionName
		}
	case FrontendParentName:
		section = parser.Frontends
	case BackendParentName:
		section = parser.Backends
	}

	data, err := p.GetOne(section, parentName, "http-error", int(id))
	if err != nil {
		return v, nil, c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, "", false, err)
	}

	httpRule := ParseHTTPErrorRule(data.(types.Action))

	return v, httpRule, nil
}

// DeleteHTTPErrorRule deletes a http error rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) DeleteHTTPErrorRule(id int64, parentType string, parentName string, transactionID string, version int64) error {
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	var section parser.Section
	switch parentType {
	case DefaultsParentName:
		section = parser.Defaults
		if parentName == "" {
			parentName = parser.DefaultSectionName
		}
	case FrontendParentName:
		section = parser.Frontends
	case BackendParentName:
		section = parser.Backends
	}

	if err := p.Delete(section, parentName, "http-error", int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

// CreateHTTPErrorRule creates a http error rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) CreateHTTPErrorRule(id int64, parentType string, parentName string, data *models.HTTPErrorRule, transactionID string, version int64) error {
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

	var section parser.Section
	switch parentType {
	case DefaultsParentName:
		section = parser.Defaults
		if parentName == "" {
			parentName = parser.DefaultSectionName
		}
	case FrontendParentName:
		section = parser.Frontends
	case BackendParentName:
		section = parser.Backends
	}

	s, err := SerializeHTTPErrorRule(*data)
	if err != nil {
		return err
	}
	if err := p.Insert(section, parentName, "http-error", s, int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

// EditHTTPErrorRule edits a http error rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) EditHTTPErrorRule(id int64, parentType string, parentName string, data *models.HTTPErrorRule, transactionID string, version int64) error {
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

	var section parser.Section
	switch parentType {
	case DefaultsParentName:
		section = parser.Defaults
		if parentName == "" {
			parentName = parser.DefaultSectionName
		}
	case FrontendParentName:
		section = parser.Frontends
	case BackendParentName:
		section = parser.Backends
	}

	if _, err = p.GetOne(section, parentName, "http-error", int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	s, err := SerializeHTTPErrorRule(*data)
	if err != nil {
		return err
	}
	if err := p.Set(section, parentName, "http-error", s, int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

// ReplaceHTTPErrorRules replaces all HTTP Error Rule lines in configuration for a parentType/parentName.
// One of version or transactionID is mandatory.
// Returns error on fail, nil on success.
func (c *client) ReplaceHTTPErrorRules(parentType string, parentName string, data models.HTTPErrorRules, transactionID string, version int64) error {
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

	httpErrorRules, err := ParseHTTPErrorRules(parentType, parentName, p)
	if err != nil {
		return c.HandleError("", parentType, parentName, "", false, err)
	}

	var section parser.Section
	switch parentType {
	case DefaultsParentName:
		section = parser.Defaults
		if parentName == "" {
			parentName = parser.DefaultSectionName
		}
	case FrontendParentName:
		section = parser.Frontends
	case BackendParentName:
		section = parser.Backends
	}

	for i := range httpErrorRules {
		// Always delete index 0
		if err := p.Delete(section, parentName, "http-error", 0); err != nil {
			return c.HandleError(strconv.FormatInt(int64(i), 10), parentType, parentName, t, transactionID == "", err)
		}
	}

	for i, newHTTPErrorRule := range data {
		s, err := SerializeHTTPErrorRule(*newHTTPErrorRule)
		if err != nil {
			return err
		}
		if err := p.Insert(section, parentName, "http-error", s, i); err != nil {
			return c.HandleError(strconv.FormatInt(int64(i), 10), parentType, parentName, t, transactionID == "", err)
		}
	}

	return c.SaveData(p, t, transactionID == "")
}

func ParseHTTPErrorRules(t, pName string, p parser.Parser) (models.HTTPErrorRules, error) {
	var section parser.Section
	switch t {
	case DefaultsParentName:
		section = parser.Defaults
		if pName == "" {
			pName = parser.DefaultSectionName
		}
	case FrontendParentName:
		section = parser.Frontends
	case BackendParentName:
		section = parser.Backends
	default:
		return nil, NewConfError(ErrValidationError, fmt.Sprintf("unsupported section in http_error: %s", t))
	}

	var httpErrRules models.HTTPErrorRules
	data, err := p.Get(section, pName, "http-error", false)
	if err != nil {
		if goerrors.Is(err, parser_errors.ErrFetch) {
			return httpErrRules, nil
		}
		return nil, err
	}

	rules, ok := data.([]types.Action)
	if !ok {
		return nil, misc.CreateTypeAssertError("http-error")
	}
	for _, r := range rules {
		httpResRule := ParseHTTPErrorRule(r)
		if httpResRule != nil {
			httpErrRules = append(httpErrRules, httpResRule)
		}
	}
	return httpErrRules, nil
}

func ParseHTTPErrorRule(f types.Action) *models.HTTPErrorRule {
	switch v := f.(type) {
	case *http_actions.Status:
		return &models.HTTPErrorRule{
			Type:                "status",
			Status:              *v.Status,
			ReturnHeaders:       actionHdr2ModelHdr(v.Hdrs),
			ReturnContent:       v.Content,
			ReturnContentFormat: v.ContentFormat,
			ReturnContentType:   &v.ContentType,
		}
	default:
		return nil
	}
}

func SerializeHTTPErrorRule(f models.HTTPErrorRule) (rule types.Action, err error) { //nolint:ireturn
	if f.Type != "status" {
		return nil, NewConfError(ErrValidationError, fmt.Sprintf("unsupported action %s in http_error", f.Type))
	}

	contentType := ""
	if f.ReturnContentType != nil {
		contentType = *f.ReturnContentType
	}
	rule = &http_actions.Status{
		Status:        &f.Status,
		ContentType:   contentType,
		ContentFormat: f.ReturnContentFormat,
		Content:       f.ReturnContent,
		Hdrs:          modelHdr2ActionHdr(f.ReturnHeaders),
	}

	if !http_actions.AllowedErrorStatusCode(f.Status) {
		return rule, NewConfError(ErrValidationError, fmt.Sprintf("unsupported status code %d in http_error", f.Status))
	}

	return rule, nil
}
