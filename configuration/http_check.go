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
	"errors"
	"strconv"
	"strings"

	strfmt "github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/client-native/v6/config-parser"
	"github.com/haproxytech/client-native/v6/config-parser/common"
	parser_errors "github.com/haproxytech/client-native/v6/config-parser/errors"
	actions "github.com/haproxytech/client-native/v6/config-parser/parsers/actions"
	http_actions "github.com/haproxytech/client-native/v6/config-parser/parsers/http/actions"
	"github.com/haproxytech/client-native/v6/config-parser/types"

	"github.com/haproxytech/client-native/v6/misc"
	"github.com/haproxytech/client-native/v6/models"
)

type HTTPCheck interface {
	GetHTTPChecks(parentType, parentName string, transactionID string) (int64, models.HTTPChecks, error)
	GetHTTPCheck(id int64, parentType string, parentName string, transactionID string) (int64, *models.HTTPCheck, error)
	DeleteHTTPCheck(id int64, parentType string, parentName string, transactionID string, version int64) error
	CreateHTTPCheck(id int64, parentType string, parentName string, data *models.HTTPCheck, transactionID string, version int64) error
	EditHTTPCheck(id int64, parentType string, parentName string, data *models.HTTPCheck, transactionID string, version int64) error
	ReplaceHTTPChecks(parentType string, parentName string, data models.HTTPChecks, transactionID string, version int64) error
}

// GetHTTPChecks returns configuration version and an array of configured http-checks in the specified parent.
// Returns error on fail.
func (c *client) GetHTTPChecks(parentType, parentName string, transactionID string) (int64, models.HTTPChecks, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	httpChecks, err := ParseHTTPChecks(parentType, parentName, p)
	if err != nil {
		return v, nil, c.HandleError("", parentType, parentName, "", false, err)
	}

	return v, httpChecks, nil
}

// GetHTTPCheck returns configuration version and the requested http check in the specified parent.
// Returns error on fail or if http check does not exist
func (c *client) GetHTTPCheck(id int64, parentType string, parentName string, transactionID string) (int64, *models.HTTPCheck, error) {
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
	case BackendParentName:
		section = parser.Backends
	case DefaultsParentName:
		section = parser.Defaults
		if parentName == "" {
			parentName = parser.DefaultSectionName
		}
	}

	data, err := p.GetOne(section, parentName, "http-check", int(id))
	if err != nil {
		return v, nil, c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, "", false, err)
	}

	httpCheck, err := ParseHTTPCheck(data.(types.Action))
	if err != nil {
		return v, nil, err
	}
	return v, httpCheck, nil
}

// DeleteHTTPCheck deletes a http check in the configuration. One of version or transactionID is mandatory.
// Returns error on fail, nil on success.
func (c *client) DeleteHTTPCheck(id int64, parentType string, parentName string, transactionID string, version int64) error {
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	var section parser.Section
	switch parentType {
	case BackendParentName:
		section = parser.Backends
	case DefaultsParentName:
		section = parser.Defaults
		if parentName == "" {
			parentName = parser.DefaultSectionName
		}
	}

	if err := p.Delete(section, parentName, "http-check", int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

// CreateHTTPCheck creates a http check in the configuration. One of version or transationID is mandatory.
// Returns error on fail, nil on success.
func (c *client) CreateHTTPCheck(id int64, parentType string, parentName string, data *models.HTTPCheck, transactionID string, version int64) error {
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
	case BackendParentName:
		section = parser.Backends
	case DefaultsParentName:
		section = parser.Defaults
		if parentName == "" {
			parentName = parser.DefaultSectionName
		}
	}

	check, err := SerializeHTTPCheck(*data)
	if err != nil {
		return err
	}

	if err := p.Insert(section, parentName, "http-check", check, int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}
	return c.SaveData(p, t, transactionID == "")
}

// EditHTTPCheck edits a http check in the configuration. One of version or transactionID is mandatory.
// Returns error on fail, nil on success.
func (c *client) EditHTTPCheck(id int64, parentType string, parentName string, data *models.HTTPCheck, transactionID string, version int64) error { //nolint:dupl
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
	case BackendParentName:
		section = parser.Backends
	case DefaultsParentName:
		section = parser.Defaults
		if parentName == "" {
			parentName = parser.DefaultSectionName
		}
	}

	if _, err = p.GetOne(section, parentName, "http-check", int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	check, err := SerializeHTTPCheck(*data)
	if err != nil {
		return err
	}

	if err := p.Set(section, parentName, "http-check", check, int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}
	return c.SaveData(p, t, transactionID == "")
}

// ReplaceHTTPChecks replaces all HTTP check lines in configuration for a parentType/parentName.
// One of version or transactionID is mandatory.
// Returns error on fail, nil on success.
func (c *client) ReplaceHTTPChecks(parentType string, parentName string, data models.HTTPChecks, transactionID string, version int64) error {
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

	httpChecks, err := ParseHTTPChecks(parentType, parentName, p)
	if err != nil {
		return c.HandleError("", parentType, parentName, "", false, err)
	}

	var section parser.Section
	switch parentType {
	case BackendParentName:
		section = parser.Backends
	case DefaultsParentName:
		section = parser.Defaults
		if parentName == "" {
			parentName = parser.DefaultSectionName
		}
	}

	for i := range httpChecks {
		// Always delete index 0
		if err := p.Delete(section, parentName, "http-check", 0); err != nil {
			return c.HandleError(strconv.FormatInt(int64(i), 10), parentType, parentName, t, transactionID == "", err)
		}
	}

	for i, newHTTPCheck := range data {
		s, err := SerializeHTTPCheck(*newHTTPCheck)
		if err != nil {
			return err
		}
		if err := p.Insert(section, parentName, "http-check", s, i); err != nil {
			return c.HandleError(strconv.FormatInt(int64(i), 10), parentType, parentName, t, transactionID == "", err)
		}
	}

	return c.SaveData(p, t, transactionID == "")
}

func ParseHTTPChecks(t, pName string, p parser.Parser) (models.HTTPChecks, error) {
	var section parser.Section
	switch t {
	case DefaultsParentName:
		section = parser.Defaults
		if pName == "" {
			pName = parser.DefaultSectionName
		}
	case BackendParentName:
		section = parser.Backends
	default:
		return nil, NewConfError(ErrValidationError, "unsupported section in http_error: "+t)
	}

	var checks models.HTTPChecks
	data, err := p.Get(section, pName, "http-check", false)
	if err != nil {
		if errors.Is(err, parser_errors.ErrFetch) {
			return checks, nil
		}
		return nil, err
	}
	items, ok := data.([]types.Action)
	if !ok {
		return nil, misc.CreateTypeAssertError("http-check")
	}
	for _, c := range items {
		check, err := ParseHTTPCheck(c)
		if err == nil {
			checks = append(checks, check)
		}
	}
	return checks, nil
}

func ParseHTTPCheck(f types.Action) (*models.HTTPCheck, error) {
	var check *models.HTTPCheck
	switch v := f.(type) {
	case *http_actions.CheckComment:
		check = &models.HTTPCheck{
			Type:         models.HTTPCheckTypeComment,
			CheckComment: v.LogMessage,
		}
	case *actions.CheckConnect:
		check = &models.HTTPCheck{
			Type:         models.HTTPCheckTypeConnect,
			Addr:         v.Addr,
			Sni:          v.SNI,
			Alpn:         v.ALPN,
			Proto:        v.Proto,
			CheckComment: v.CheckComment,
			Default:      v.Default,
			SendProxy:    v.SendProxy,
			ViaSocks4:    v.ViaSOCKS4,
			Ssl:          v.SSL,
			Linger:       v.Linger,
		}
		if v.Port != "" {
			portInt, err := strconv.ParseInt(v.Port, 10, 64)
			if err == nil {
				check.Port = misc.Int64P(int(portInt))
			} else {
				check.PortString = v.Port
			}
		}
	case *actions.CheckExpect:
		check = &models.HTTPCheck{
			Type:            models.HTTPCheckTypeExpect,
			CheckComment:    v.CheckComment,
			OkStatus:        v.OKStatus,
			ErrorStatus:     v.ErrorStatus,
			ToutStatus:      v.TimeoutStatus,
			OnSuccess:       v.OnSuccess,
			OnError:         v.OnError,
			StatusCode:      v.StatusCode,
			ExclamationMark: v.ExclamationMark,
			Match:           v.Match,
			Pattern:         v.Pattern,
		}
		if v.MinRecv != nil {
			check.MinRecv = v.MinRecv
		}
	case *http_actions.CheckDisableOn404:
		check = &models.HTTPCheck{
			Type: models.HTTPCheckTypeDisableDashOnDash404,
		}
	case *http_actions.CheckSend:
		check = &models.HTTPCheck{
			Type:          models.HTTPCheckTypeSend,
			Method:        v.Method,
			URI:           v.URI,
			URILogFormat:  v.URILogFormat,
			Version:       v.Version,
			Body:          v.Body,
			BodyLogFormat: v.BodyLogFormat,
			CheckComment:  v.CheckComment,
		}
		headers := make([]*models.ReturnHeader, len(v.Header))
		for i, h := range v.Header {
			name := h.Name
			value := h.Format
			header := &models.ReturnHeader{
				Name: &name,
				Fmt:  &value,
			}
			headers[i] = header
		}
		check.CheckHeaders = headers
	case *http_actions.CheckSendState:
		check = &models.HTTPCheck{
			Type: models.HTTPCheckTypeSendDashState,
		}
	case *actions.SetVarCheck:
		check = &models.HTTPCheck{
			Type:     models.HTTPCheckTypeSetDashVar,
			VarScope: v.VarScope,
			VarName:  v.VarName,
			VarExpr:  strings.Join(v.Expr.Expr, " "),
		}
	case *actions.SetVarFmtCheck:
		check = &models.HTTPCheck{
			Type:     models.HTTPCheckTypeSetDashVarDashFmt,
			VarScope: v.VarScope,
			VarName:  v.VarName,
			VarExpr:  strings.Join(v.Format.Expr, " "),
		}
	case *actions.UnsetVarCheck:
		check = &models.HTTPCheck{
			Type:     models.HTTPCheckTypeUnsetDashVar,
			VarScope: v.Scope,
			VarName:  v.Name,
		}
	}

	return check, nil
}

func SerializeHTTPCheck(f models.HTTPCheck) (types.Action, error) { //nolint:ireturn
	switch f.Type {
	case models.HTTPCheckTypeComment:
		return &http_actions.CheckComment{
			LogMessage: f.CheckComment,
		}, nil
	case models.HTTPCheckTypeConnect:
		port := f.PortString
		if f.Port != nil {
			port = strconv.FormatInt(*f.Port, 10)
		}
		return &actions.CheckConnect{
			Port:         port,
			Addr:         f.Addr,
			SNI:          f.Sni,
			ALPN:         f.Alpn,
			Proto:        f.Proto,
			CheckComment: f.CheckComment,
			Default:      f.Default,
			SendProxy:    f.SendProxy,
			ViaSOCKS4:    f.ViaSocks4,
			SSL:          f.Ssl,
			Linger:       f.Linger,
		}, nil
	case models.HTTPCheckTypeExpect:
		return &actions.CheckExpect{
			MinRecv:         f.MinRecv,
			Match:           f.Match,
			OKStatus:        f.OkStatus,
			ErrorStatus:     f.ErrorStatus,
			CheckComment:    f.CheckComment,
			TimeoutStatus:   f.ToutStatus,
			OnSuccess:       f.OnSuccess,
			OnError:         f.OnError,
			StatusCode:      f.StatusCode,
			ExclamationMark: f.ExclamationMark,
			Pattern:         f.Pattern,
		}, nil
	case models.HTTPCheckTypeDisableDashOnDash404:
		return &http_actions.CheckDisableOn404{}, nil
	case models.HTTPCheckTypeSend:
		action := &http_actions.CheckSend{
			Method:        f.Method,
			URI:           f.URI,
			URILogFormat:  f.URILogFormat,
			Version:       f.Version,
			Body:          f.Body,
			BodyLogFormat: f.BodyLogFormat,
			CheckComment:  f.CheckComment,
		}
		headers := []http_actions.CheckSendHeader{}
		for _, h := range f.CheckHeaders {
			if h == nil || h.Name == nil || h.Fmt == nil {
				continue
			}
			header := http_actions.CheckSendHeader{
				Name:   *h.Name,
				Format: *h.Fmt,
			}
			headers = append(headers, header)
		}
		action.Header = headers
		return action, nil
	case models.HTTPCheckTypeSendDashState:
		return &http_actions.CheckSendState{}, nil
	case models.HTTPCheckTypeSetDashVar:
		return &actions.SetVarCheck{
			VarScope: f.VarScope,
			VarName:  f.VarName,
			Expr:     common.Expression{Expr: strings.Split(f.VarExpr, " ")},
		}, nil
	case models.HTTPCheckTypeSetDashVarDashFmt:
		return &actions.SetVarFmtCheck{
			VarScope: f.VarScope,
			VarName:  f.VarName,
			Format:   common.Expression{Expr: strings.Split(f.VarFormat, " ")},
		}, nil
	case models.HTTPCheckTypeUnsetDashVar:
		return &actions.UnsetVarCheck{
			Scope: f.VarScope,
			Name:  f.VarName,
		}, nil
	}

	return nil, NewConfError(ErrValidationError, "unsupported action in http_check")
}
