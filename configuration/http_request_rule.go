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
	"strconv"
	"strings"

	"github.com/haproxytech/config-parser/v2/common"

	"github.com/haproxytech/config-parser/v2/parsers/http/actions"

	strfmt "github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/config-parser/v2"
	parser_errors "github.com/haproxytech/config-parser/v2/errors"
	"github.com/haproxytech/config-parser/v2/types"
	"github.com/haproxytech/models"
)

// GetHTTPRequestRules returns configuration version and an array of
// configured http request rules in the specified parent. Returns error on fail.
func (c *Client) GetHTTPRequestRules(parentType, parentName string, transactionID string) (int64, models.HTTPRequestRules, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	httpRules, err := ParseHTTPRequestRules(parentType, parentName, p)
	if err != nil {
		return v, nil, c.handleError("", parentType, parentName, "", false, err)
	}

	return v, httpRules, nil
}

// GetHTTPRequestRule returns configuration version and a requested http request rule
// in the specified parent. Returns error on fail or if http request rule does not exist.
func (c *Client) GetHTTPRequestRule(id int64, parentType, parentName string, transactionID string) (int64, *models.HTTPRequestRule, error) {
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

	data, err := p.GetOne(section, parentName, "http-request", int(id))
	if err != nil {
		return v, nil, c.handleError(strconv.FormatInt(id, 10), parentType, parentName, "", false, err)
	}

	httpRule, err := ParseHTTPRequestRule(data.(types.HTTPAction))
	if err != nil {
		return v, nil, err
	}
	httpRule.Index = &id

	return v, httpRule, nil
}

// DeleteHTTPRequestRule deletes a http request rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) DeleteHTTPRequestRule(id int64, parentType string, parentName string, transactionID string, version int64) error {
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

	if err := p.Delete(section, parentName, "http-request", int(id)); err != nil {
		return c.handleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	if err := c.saveData(p, t, transactionID == ""); err != nil {
		return err
	}
	return nil
}

// CreateHTTPRequestRule creates a http request rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) CreateHTTPRequestRule(parentType string, parentName string, data *models.HTTPRequestRule, transactionID string, version int64) error {
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
	} else if parentType == "frontend" {
		section = parser.Frontends
	}

	s, err := SerializeHTTPRequestRule(*data)
	if err != nil {
		return err
	}

	if err := p.Insert(section, parentName, "http-request", s, int(*data.Index)); err != nil {
		return c.handleError(strconv.FormatInt(*data.Index, 10), parentType, parentName, t, transactionID == "", err)
	}

	if err := c.saveData(p, t, transactionID == ""); err != nil {
		return err
	}
	return nil
}

// EditHTTPRequestRule edits a http request rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) EditHTTPRequestRule(id int64, parentType string, parentName string, data *models.HTTPRequestRule, transactionID string, version int64) error {
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
	} else if parentType == "frontend" {
		section = parser.Frontends
	}

	if _, err := p.GetOne(section, parentName, "http-request", int(id)); err != nil {
		return c.handleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	s, err := SerializeHTTPRequestRule(*data)
	if err != nil {
		return err
	}

	if err := p.Set(section, parentName, "http-request", s, int(id)); err != nil {
		return c.handleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	if err := c.saveData(p, t, transactionID == ""); err != nil {
		return err
	}
	return nil
}

func ParseHTTPRequestRules(t, pName string, p *parser.Parser) (models.HTTPRequestRules, error) {
	section := parser.Global
	if t == "frontend" {
		section = parser.Frontends
	} else if t == "backend" {
		section = parser.Backends
	}

	httpReqRules := models.HTTPRequestRules{}
	data, err := p.Get(section, pName, "http-request", false)
	if err != nil {
		if err == parser_errors.ErrFetch {
			return httpReqRules, nil
		}
		return nil, err
	}

	rules := data.([]types.HTTPAction)
	for i, r := range rules {
		id := int64(i)
		httpReqRule, err := ParseHTTPRequestRule(r)
		if err == nil {
			httpReqRule.Index = &id
			httpReqRules = append(httpReqRules, httpReqRule)
		}
	}
	return httpReqRules, nil
}

func ParseHTTPRequestRule(f types.HTTPAction) (rule *models.HTTPRequestRule, err error) {
	switch v := f.(type) {
	case *actions.Allow:
		rule = &models.HTTPRequestRule{
			Type:     "allow",
			Cond:     v.Cond,
			CondTest: v.CondTest,
		}
	case *actions.Deny:
		s, _ := strconv.ParseInt(v.DenyStatus, 10, 64)
		rule = &models.HTTPRequestRule{
			Type:     "deny",
			Cond:     v.Cond,
			CondTest: v.CondTest,
		}
		if s != 0 {
			rule.DenyStatus = s
		}
	case *actions.Auth:
		rule = &models.HTTPRequestRule{
			Type:      "auth",
			AuthRealm: v.Realm,
			Cond:      v.Cond,
			CondTest:  v.CondTest,
		}
	case *actions.Redirect:
		code, _ := strconv.ParseInt(v.Code, 10, 64)
		rule = &models.HTTPRequestRule{
			Type:        "redirect",
			RedirType:   v.Type,
			RedirValue:  v.Value,
			RedirOption: v.Option,
			Cond:        v.Cond,
			CondTest:    v.CondTest,
		}
		if code != 0 {
			rule.RedirCode = code
		}
	case *actions.Tarpit:
		s, _ := strconv.ParseInt(v.DenyStatus, 10, 64)
		rule = &models.HTTPRequestRule{
			Type:     "tarpit",
			Cond:     v.Cond,
			CondTest: v.CondTest,
		}
		if s != 0 {
			rule.DenyStatus = s
		}
	case *actions.AddHeader:
		rule = &models.HTTPRequestRule{
			Type:      "add-header",
			HdrName:   v.Name,
			HdrFormat: v.Fmt,
			Cond:      v.Cond,
			CondTest:  v.CondTest,
		}
	case *actions.SetHeader:
		rule = &models.HTTPRequestRule{
			Type:      "set-header",
			HdrName:   v.Name,
			HdrFormat: v.Fmt,
			Cond:      v.Cond,
			CondTest:  v.CondTest,
		}
	case *actions.SetQuery:
		rule = &models.HTTPRequestRule{
			Type:      "set-query",
			HdrFormat: v.Fmt,
			Cond:      v.Cond,
			CondTest:  v.CondTest,
		}
	case *actions.SetURI:
		rule = &models.HTTPRequestRule{
			Type:      "set-uri",
			HdrFormat: v.Fmt,
			Cond:      v.Cond,
			CondTest:  v.CondTest,
		}
	case *actions.DelHeader:
		rule = &models.HTTPRequestRule{
			Type:     "del-header",
			HdrName:  v.Name,
			Cond:     v.Cond,
			CondTest: v.CondTest,
		}
	case *actions.ReplaceHeader:
		rule = &models.HTTPRequestRule{
			Type:      "replace-header",
			HdrName:   v.Name,
			HdrFormat: v.ReplaceFmt,
			HdrMatch:  v.MatchRegex,
			Cond:      v.Cond,
			CondTest:  v.CondTest,
		}
	case *actions.SetLogLevel:
		rule = &models.HTTPRequestRule{
			Type:     "set-log-level",
			LogLevel: v.Level,
			Cond:     v.Cond,
			CondTest: v.CondTest,
		}
	case *actions.SetPath:
		rule = &models.HTTPRequestRule{
			Type:     "set-path",
			PathFmt:  v.Fmt,
			Cond:     v.Cond,
			CondTest: v.CondTest,
		}
	case *actions.SetVar:
		rule = &models.HTTPRequestRule{
			Type:     "set-var",
			VarName:  v.VarName,
			VarExpr:  strings.Join(v.Expr.Expr, " "),
			VarScope: v.VarScope,
			Cond:     v.Cond,
			CondTest: v.CondTest,
		}
	case *actions.ReplaceValue:
		rule = &models.HTTPRequestRule{
			Type:      "replace-value",
			HdrName:   v.Name,
			HdrMatch:  v.MatchRegex,
			HdrFormat: v.ReplaceFmt,
			Cond:      v.Cond,
			CondTest:  v.CondTest,
		}
	case *actions.AddACL:
		rule = &models.HTTPRequestRule{
			Type:      "add-acl",
			ACLFile:   v.FileName,
			ACLKeyfmt: v.KeyFmt,
			Cond:      v.Cond,
			CondTest:  v.CondTest,
		}
	case *actions.DelACL:
		rule = &models.HTTPRequestRule{
			Type:      "del-acl",
			ACLFile:   v.FileName,
			ACLKeyfmt: v.KeyFmt,
			Cond:      v.Cond,
			CondTest:  v.CondTest,
		}
	case *actions.SendSpoeGroup:
		rule = &models.HTTPRequestRule{
			Type:       "send-spoe-group",
			SpoeEngine: v.Engine,
			SpoeGroup:  v.Group,
			Cond:       v.Cond,
			CondTest:   v.CondTest,
		}
	case *actions.Capture:
		if (v.SlotID == nil && v.Len == nil) || (v.SlotID != nil && v.Len != nil) {
			return nil, NewConfError(ErrValidationError, "capture len can't be zero")
		}
		rule = &models.HTTPRequestRule{
			Type:          "capture",
			CaptureSample: v.Sample,
			Cond:          v.Cond,
			CondTest:      v.CondTest,
		}
		if v.SlotID != nil {
			rule.CaptureID = v.SlotID
		}
		if v.Len != nil {
			rule.CaptureLen = *v.Len
		}
	case *actions.TrackSc0:
		rule = &models.HTTPRequestRule{
			Type:          "track-sc0",
			TrackSc0Key:   v.Key,
			TrackSc0Table: v.Table,
			Cond:          v.Cond,
			CondTest:      v.CondTest,
		}
	case *actions.TrackSc1:
		rule = &models.HTTPRequestRule{
			Type:          "track-sc1",
			TrackSc1Key:   v.Key,
			TrackSc1Table: v.Table,
			Cond:          v.Cond,
			CondTest:      v.CondTest,
		}
	case *actions.TrackSc2:
		rule = &models.HTTPRequestRule{
			Type:          "track-sc2",
			TrackSc2Key:   v.Key,
			TrackSc2Table: v.Table,
			Cond:          v.Cond,
			CondTest:      v.CondTest,
		}
	}
	return rule, err
}

func SerializeHTTPRequestRule(f models.HTTPRequestRule) (rule types.HTTPAction, err error) {
	switch f.Type {
	case "allow":
		rule = &actions.Allow{
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "deny":
		rule = &actions.Deny{
			DenyStatus: strconv.FormatInt(f.DenyStatus, 10),
			Cond:       f.Cond,
			CondTest:   f.CondTest,
		}
	case "auth":
		rule = &actions.Auth{
			Realm:    f.AuthRealm,
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "redirect":
		rule = &actions.Redirect{
			Type:     f.RedirType,
			Value:    f.RedirValue,
			Code:     strconv.FormatInt(f.RedirCode, 10),
			Option:   f.RedirOption,
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "tarpit":
		rule = &actions.Tarpit{
			DenyStatus: strconv.FormatInt(f.DenyStatus, 10),
			Cond:       f.Cond,
			CondTest:   f.CondTest,
		}
	case "add-header":
		rule = &actions.AddHeader{
			Name:     f.HdrName,
			Fmt:      f.HdrFormat,
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "set-header":
		rule = &actions.SetHeader{
			Name:     f.HdrName,
			Fmt:      f.HdrFormat,
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "set-query":
		rule = &actions.SetQuery{
			Fmt:      f.HdrFormat,
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "set-uri":
		rule = &actions.SetURI{
			Fmt:      f.HdrFormat,
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "del-header":
		rule = &actions.DelHeader{
			Name:     f.HdrName,
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "replace-header":
		rule = &actions.ReplaceHeader{
			Name:       f.HdrName,
			ReplaceFmt: f.HdrFormat,
			MatchRegex: f.HdrMatch,
			Cond:       f.Cond,
			CondTest:   f.CondTest,
		}
	case "replace-value":
		rule = &actions.ReplaceValue{
			Name:       f.HdrName,
			ReplaceFmt: f.HdrFormat,
			MatchRegex: f.HdrMatch,
			Cond:       f.Cond,
			CondTest:   f.CondTest,
		}
	case "set-log-level":
		rule = &actions.SetLogLevel{
			Level:    f.LogLevel,
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "set-path":
		rule = &actions.SetPath{
			Fmt:      f.PathFmt,
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
	case "add-acl":
		rule = &actions.AddACL{
			FileName: f.ACLFile,
			KeyFmt:   f.ACLKeyfmt,
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "del-acl":
		rule = &actions.DelACL{
			FileName: f.ACLFile,
			KeyFmt:   f.ACLKeyfmt,
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "send-spoe-group":
		rule = &actions.SendSpoeGroup{
			Engine:   f.SpoeEngine,
			Group:    f.SpoeGroup,
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "capture":
		if f.CaptureLen > 0 && f.CaptureID != nil {
			return nil, NewConfError(ErrValidationError, "capture len and id are exclusive")
		}
		if f.CaptureLen == 0 && f.CaptureID == nil {
			return nil, NewConfError(ErrValidationError, "capture len has to be greater than 0 or capture_id has to be set")
		}
		r := &actions.Capture{
			Sample:   f.CaptureSample,
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
		if f.CaptureLen > 0 {
			r.Len = &f.CaptureLen
		} else {
			r.SlotID = f.CaptureID
		}
		rule = r
	case "track-sc0":
		rule = &actions.TrackSc0{
			Key:   f.TrackSc0Key,
			Table: f.TrackSc0Table,
		}
	case "track-sc1":
		rule = &actions.TrackSc1{
			Key:   f.TrackSc1Key,
			Table: f.TrackSc1Table,
		}
	case "track-sc2":
		rule = &actions.TrackSc2{
			Key:      f.TrackSc2Key,
			Table:    f.TrackSc2Table,
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	}
	return rule, err
}
