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
	"fmt"
	"strconv"
	"strings"

	strfmt "github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/config-parser/v4"
	"github.com/haproxytech/config-parser/v4/common"
	parser_errors "github.com/haproxytech/config-parser/v4/errors"
	actions "github.com/haproxytech/config-parser/v4/parsers/actions"
	http_actions "github.com/haproxytech/config-parser/v4/parsers/http/actions"
	"github.com/haproxytech/config-parser/v4/types"

	"github.com/haproxytech/client-native/v3/misc"
	"github.com/haproxytech/client-native/v3/models"
)

// GetHTTPChecks returns configuration version and an array of configured http-checks in the specified parent.
// Returns error on fail.
func (c *Client) GetHTTPChecks(parentType, parentName string, transactionID string) (int64, models.HTTPCheckRules, error) {
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
func (c *Client) GetHTTPCheck(id int64, parentType string, parentName string, transactionID string) (int64, *models.HTTPCheckRule, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}
	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	var section parser.Section
	if parentType == "backend" {
		section = parser.Backends
	} else if parentType == "defaults" {
		section = parser.Defaults
		parentName = parser.DefaultSectionName
	}

	data, err := p.GetOne(section, parentName, "http-check", int(id))
	if err != nil {
		return v, nil, c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, "", false, err)
	}

	httpCheck, err := ParseHTTPCheck(data.(types.Action))
	if err != nil {
		return v, nil, err
	}
	httpCheck.Index = &id
	return v, httpCheck, nil
}

// DeleteHTTPCheck deletes a http check in the configuration. One of version or transactionID is mandatory.
// Returns error on fail, nil on success.
func (c *Client) DeleteHTTPCheck(id int64, parentType string, parentName string, transactionID string, version int64) error {
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	var section parser.Section
	if parentType == "backend" {
		section = parser.Backends
	} else if parentType == "defaults" {
		section = parser.Defaults
		parentName = parser.DefaultSectionName
	}

	if err := p.Delete(section, parentName, "http-check", int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	if err := c.SaveData(p, t, transactionID == ""); err != nil {
		return err
	}
	return nil
}

// CreateHTTPCheck creates a http check in the configuration. One of version or transationID is mandatory.
// Returns error on fail, nil on success.
func (c *Client) CreateHTTPCheck(parentType string, parentName string, data *models.HTTPCheckRule, transactionID string, version int64) error {
	if c.UseValidation {
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
	if parentType == "backend" {
		section = parser.Backends
	} else if parentType == "defaults" {
		section = parser.Defaults
		parentName = parser.DefaultSectionName
	}

	check, err := SerializeHTTPCheck(*data)
	if err != nil {
		return err
	}

	if err := p.Insert(section, parentName, "http-check", check, int(*data.Index)); err != nil {
		return c.HandleError(strconv.FormatInt(*data.Index, 10), parentType, parentName, t, transactionID == "", err)
	}
	if err := c.SaveData(p, t, transactionID == ""); err != nil {
		return err
	}
	return nil
}

// EditHTTPCheck edits a http check in the configuration. One of version or transactionID is mandatory.
// Returns error on fail, nil on success.
// nolint:dupl
func (c *Client) EditHTTPCheck(id int64, parentType string, parentName string, data *models.HTTPCheckRule, transactionID string, version int64) error {
	if c.UseValidation {
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
	if parentType == "backend" {
		section = parser.Backends
	} else if parentType == "defaults" {
		section = parser.Defaults
		parentName = parser.DefaultSectionName
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
	if err := c.SaveData(p, t, transactionID == ""); err != nil {
		return err
	}
	return nil
}

func ParseHTTPChecks(t, pName string, p parser.Parser) (models.HTTPCheckRules, error) {
	var section parser.Section
	switch t {
	case "defaults":
		section = parser.Defaults
		pName = parser.DefaultSectionName
	case "backend":
		section = parser.Backends
	default:
		return nil, NewConfError(ErrValidationError, fmt.Sprintf("unsupported section in http_check: %s", t))
	}

	checks := models.HTTPCheckRules{}
	data, err := p.Get(section, pName, "http-check", false)
	if err != nil {
		if errors.Is(err, parser_errors.ErrFetch) {
			return checks, nil
		}
		return nil, err
	}
	items := data.([]types.Action)
	for i, c := range items {
		id := int64(i)
		check, err := ParseHTTPCheck(c)
		if err == nil {
			check.Index = &id
			checks = append(checks, check)
		}
	}
	return checks, nil
}

func ParseHTTPCheck(f types.Action) (check *models.HTTPCheckRule, err error) {
	switch v := f.(type) {
	case *http_actions.CheckComment:
		check = &models.HTTPCheckRule{
			Action:       models.HTTPCheckRuleActionComment,
			CheckComment: v.LogMessage,
		}
	case *actions.CheckConnect:
		check = &models.HTTPCheckRule{
			Action:       models.HTTPCheckRuleActionConnect,
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
		check = &models.HTTPCheckRule{
			Action:          models.HTTPCheckRuleActionExpect,
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
			check.MinRecv = *v.MinRecv
		}
	case *http_actions.CheckDisableOn404:
		check = &models.HTTPCheckRule{
			Action: models.HTTPCheckRuleActionDisableOn404,
		}
	case *http_actions.CheckSend:
		check = &models.HTTPCheckRule{
			Action:        models.HTTPCheckRuleActionSend,
			Method:        v.Method,
			URI:           v.URI,
			URILogFormat:  v.URILogFormat,
			Version:       v.Version,
			Body:          v.Body,
			BodyLogFormat: v.BodyLogFormat,
			CheckComment:  v.CheckComment,
		}
		headers := []*models.ReturnHeader{}
		for _, h := range v.Header {
			name := h.Name
			value := h.Format
			header := &models.ReturnHeader{
				Name: &name,
				Fmt:  &value,
			}
			headers = append(headers, header)
		}
		check.CheckHeaders = headers
	case *http_actions.CheckSendState:
		check = &models.HTTPCheckRule{
			Action: models.HTTPCheckRuleActionSendState,
		}
	case *actions.SetVarCheck:
		check = &models.HTTPCheckRule{
			Action:   models.HTTPCheckRuleActionSetVar,
			VarScope: v.VarScope,
			VarName:  v.VarName,
			VarExpr:  strings.Join(v.Expr.Expr, " "),
		}
	case *actions.SetVarFmtCheck:
		check = &models.HTTPCheckRule{
			Action:   models.HTTPCheckRuleActionSetVarFmt,
			VarScope: v.VarScope,
			VarName:  v.VarName,
			VarExpr:  strings.Join(v.Format.Expr, " "),
		}
	case *actions.UnsetVarCheck:
		check = &models.HTTPCheckRule{
			Action:   models.HTTPCheckRuleActionUnsetVar,
			VarScope: v.Scope,
			VarName:  v.Name,
		}
	}

	return check, nil
}

func SerializeHTTPCheck(f models.HTTPCheckRule) (action types.Action, err error) { //nolint:ireturn
	switch f.Action {
	case models.HTTPCheckRuleActionComment:
		return &http_actions.CheckComment{
			LogMessage: f.CheckComment,
		}, nil
	case models.TCPCheckActionConnect:
		return &actions.CheckConnect{
			Port:         f.PortString,
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
	case models.HTTPCheckRuleActionExpect:
		return &actions.CheckExpect{
			MinRecv:         &f.MinRecv,
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
	case models.HTTPCheckRuleActionDisableOn404:
		return &http_actions.CheckDisableOn404{}, nil
	case models.HTTPCheckRuleActionSend:
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
	case models.HTTPCheckRuleActionSendState:
		return &http_actions.CheckSendState{}, nil
	case models.HTTPCheckRuleActionSetVar:
		return &actions.SetVarCheck{
			VarScope: f.VarScope,
			VarName:  f.VarName,
			Expr:     common.Expression{Expr: strings.Split(f.VarExpr, " ")},
		}, nil
	case models.HTTPCheckRuleActionSetVarFmt:
		return &actions.SetVarFmtCheck{
			VarScope: f.VarScope,
			VarName:  f.VarName,
			Format:   common.Expression{Expr: strings.Split(f.VarFormat, " ")},
		}, nil
	case models.HTTPCheckRuleActionUnsetVar:
		return &actions.UnsetVarCheck{
			Scope: f.VarScope,
			Name:  f.VarName,
		}, nil
	}

	return nil, NewConfError(ErrValidationError, "unsupported action in tcp_check")
}
