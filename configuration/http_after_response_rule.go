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
	goerrors "errors"
	"strconv"
	"strings"

	"github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/config-parser/v4"
	"github.com/haproxytech/config-parser/v4/common"
	parser_errors "github.com/haproxytech/config-parser/v4/errors"
	"github.com/haproxytech/config-parser/v4/parsers/actions"
	http_actions "github.com/haproxytech/config-parser/v4/parsers/http/actions"
	"github.com/haproxytech/config-parser/v4/types"

	"github.com/haproxytech/client-native/v4/misc"
	"github.com/haproxytech/client-native/v4/models"
)

type HTTPAfterResponseRule interface {
	GetHTTPAfterResponseRules(parentType, parentName string, transactionID string) (int64, models.HTTPAfterResponseRules, error)
	GetHTTPAfterResponseRule(id int64, parentType, parentName string, transactionID string) (int64, *models.HTTPAfterResponseRule, error)
	DeleteHTTPAfterResponseRule(id int64, parentType string, parentName string, transactionID string, version int64) error
	CreateHTTPAfterResponseRule(parentType string, parentName string, data *models.HTTPAfterResponseRule, transactionID string, version int64) error
	EditHTTPAfterResponseRule(id int64, parentType string, parentName string, data *models.HTTPAfterResponseRule, transactionID string, version int64) error
}

// GetHTTPAfterResponseRules returns configuration version and an array of configured http response rules in the specified parent.
// Returns error on fail.
func (c *client) GetHTTPAfterResponseRules(parentType, parentName string, transactionID string) (int64, models.HTTPAfterResponseRules, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	httpRules, err := ParseHTTPAfterRules(parentType, parentName, p)
	if err != nil {
		return v, nil, c.HandleError("", parentType, parentName, "", false, err)
	}

	return v, httpRules, nil
}

// GetHTTPAfterResponseRule returns configuration version and a response http response rule in the specified parent.
// Returns error on fail or if http response rule does not exist.
func (c *client) GetHTTPAfterResponseRule(id int64, parentType, parentName string, transactionID string) (int64, *models.HTTPAfterResponseRule, error) {
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
	} else if parentType == "frontend" {
		section = parser.Frontends
	}

	data, err := p.GetOne(section, parentName, "http-after-response", int(id))
	if err != nil {
		return v, nil, c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, "", false, err)
	}

	httpRule := ParseHTTPAfterRule(data.(types.Action))
	httpRule.Index = &id

	return v, httpRule, nil
}

// DeleteHTTPAfterResponseRule deletes a http response rule in configuration. One of version or transactionID is mandatory.
// Returns error on fail, nil on success.
func (c *client) DeleteHTTPAfterResponseRule(id int64, parentType string, parentName string, transactionID string, version int64) error {
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	var section parser.Section
	if parentType == "backend" {
		section = parser.Backends
	} else if parentType == "frontend" {
		section = parser.Frontends
	}

	if err := p.Delete(section, parentName, "http-after-response", int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	if err := c.SaveData(p, t, transactionID == ""); err != nil {
		return err
	}

	return nil
}

// CreateHTTPAfterResponseRule creates a http response rule in configuration. One of version or transactionID is mandatory.
// Returns error on fail, nil on success.
func (c *client) CreateHTTPAfterResponseRule(parentType string, parentName string, data *models.HTTPAfterResponseRule, transactionID string, version int64) error {
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
	if parentType == "backend" {
		section = parser.Backends
	} else if parentType == "frontend" {
		section = parser.Frontends
	}

	s, err := SerializeHTTPAfterRule(*data)
	if err != nil {
		return err
	}
	if err := p.Insert(section, parentName, "http-after-response", s, int(*data.Index)); err != nil {
		return c.HandleError(strconv.FormatInt(*data.Index, 10), parentType, parentName, t, transactionID == "", err)
	}

	if err := c.SaveData(p, t, transactionID == ""); err != nil {
		return err
	}
	return nil
}

// EditHTTPAfterResponseRule edits a http response rule in configuration. One of version or transactionID is mandatory.
// Returns error on fail, nil on success.
//
//nolint:dupl
func (c *client) EditHTTPAfterResponseRule(id int64, parentType string, parentName string, data *models.HTTPAfterResponseRule, transactionID string, version int64) error {
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
	if parentType == "backend" {
		section = parser.Backends
	} else if parentType == "frontend" {
		section = parser.Frontends
	}

	if _, err = p.GetOne(section, parentName, "http-after-response", int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	s, err := SerializeHTTPAfterRule(*data)
	if err != nil {
		return err
	}
	if err := p.Set(section, parentName, "http-after-response", s, int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	if err := c.SaveData(p, t, transactionID == ""); err != nil {
		return err
	}
	return nil
}

func ParseHTTPAfterRules(t, pName string, p parser.Parser) (models.HTTPAfterResponseRules, error) {
	section := parser.Global
	if t == "frontend" {
		section = parser.Frontends
	} else if t == "backend" {
		section = parser.Backends
	}

	httpResRules := models.HTTPAfterResponseRules{}
	data, err := p.Get(section, pName, "http-after-response", false)
	if err != nil {
		if goerrors.Is(err, parser_errors.ErrFetch) {
			return httpResRules, nil
		}
		return nil, err
	}

	rules, ok := data.([]types.Action)
	if !ok {
		return nil, misc.CreateTypeAssertError("http-after-response")
	}
	for i, r := range rules {
		id := int64(i)
		httpResRule := ParseHTTPAfterRule(r)
		if httpResRule != nil {
			httpResRule.Index = &id
			httpResRules = append(httpResRules, httpResRule)
		}
	}
	return httpResRules, nil
}

func ParseHTTPAfterRule(f types.Action) *models.HTTPAfterResponseRule {
	switch v := f.(type) {
	case *http_actions.AddHeader:
		return &models.HTTPAfterResponseRule{
			Type:      "add-header",
			HdrName:   v.Name,
			HdrFormat: v.Fmt,
			Cond:      v.Cond,
			CondTest:  v.CondTest,
		}
	case *http_actions.Allow:
		return &models.HTTPAfterResponseRule{
			Type:     "allow",
			Cond:     v.Cond,
			CondTest: v.CondTest,
		}
	case *http_actions.DelHeader:
		return &models.HTTPAfterResponseRule{
			Type:      "del-header",
			HdrName:   v.Name,
			Cond:      v.Cond,
			CondTest:  v.CondTest,
			HdrMethod: v.Method,
		}
	case *http_actions.ReplaceHeader:
		return &models.HTTPAfterResponseRule{
			Type:      "replace-header",
			HdrName:   v.Name,
			HdrFormat: v.ReplaceFmt,
			HdrMatch:  v.MatchRegex,
			Cond:      v.Cond,
			CondTest:  v.CondTest,
		}
	case *http_actions.ReplaceValue:
		return &models.HTTPAfterResponseRule{
			Type:      "replace-value",
			HdrName:   v.Name,
			HdrFormat: v.ReplaceFmt,
			HdrMatch:  v.MatchRegex,
			Cond:      v.Cond,
			CondTest:  v.CondTest,
		}
	case *http_actions.SetHeader:
		return &models.HTTPAfterResponseRule{
			Type:      "set-header",
			HdrName:   v.Name,
			HdrFormat: v.Fmt,
			Cond:      v.Cond,
			CondTest:  v.CondTest,
		}
	case *http_actions.SetStatus:
		status, _ := strconv.ParseInt(v.Status, 10, 64)
		r := &models.HTTPAfterResponseRule{
			Type:         "set-status",
			StatusReason: v.Reason,
			Cond:         v.Cond,
			CondTest:     v.CondTest,
		}
		if status != 0 {
			r.Status = status
		}
		return r
	case *actions.SetVar:
		return &models.HTTPAfterResponseRule{
			Type:     "set-var",
			VarName:  v.VarName,
			VarExpr:  strings.Join(v.Expr.Expr, " "),
			VarScope: v.VarScope,
			Cond:     v.Cond,
			CondTest: v.CondTest,
		}
	case *http_actions.StrictMode:
		return &models.HTTPAfterResponseRule{
			Type:       "strict-mode",
			StrictMode: v.Mode,
			Cond:       v.Cond,
			CondTest:   v.CondTest,
		}
	case *actions.UnsetVar:
		return &models.HTTPAfterResponseRule{
			Type:     "unset-var",
			VarName:  v.Name,
			VarScope: v.Scope,
			Cond:     v.Cond,
			CondTest: v.CondTest,
		}
	}
	return nil
}

func SerializeHTTPAfterRule(f models.HTTPAfterResponseRule) (rule types.Action, err error) { //nolint:ireturn
	switch f.Type {
	case "add-header":
		rule = &http_actions.AddHeader{
			Name:     f.HdrName,
			Fmt:      f.HdrFormat,
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "allow":
		rule = &http_actions.Allow{
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "del-header":
		rule = &http_actions.DelHeader{
			Name:     f.HdrName,
			Method:   f.HdrMethod,
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "replace-header":
		rule = &http_actions.ReplaceHeader{
			Name:       f.HdrName,
			ReplaceFmt: f.HdrFormat,
			MatchRegex: f.HdrMatch,
			Cond:       f.Cond,
			CondTest:   f.CondTest,
		}
	case "replace-value":
		rule = &http_actions.ReplaceValue{
			Name:       f.HdrName,
			ReplaceFmt: f.HdrFormat,
			MatchRegex: f.HdrMatch,
			Cond:       f.Cond,
			CondTest:   f.CondTest,
		}
	case "set-header":
		rule = &http_actions.SetHeader{
			Name:     f.HdrName,
			Fmt:      f.HdrFormat,
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "set-status":
		rule = &http_actions.SetStatus{
			Status:   strconv.FormatInt(f.Status, 10),
			Reason:   f.StatusReason,
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "set-var":
		rule = &actions.SetVar{
			Expr:     common.Expression{Expr: strings.Split(f.VarExpr, " ")},
			VarName:  f.VarName,
			VarScope: f.VarScope,
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "strict-mode":
		rule = &http_actions.StrictMode{
			Mode:     f.StrictMode,
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "unset-var":
		rule = &actions.UnsetVar{
			Name:     f.VarName,
			Scope:    f.VarScope,
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	}
	return rule, err
}
