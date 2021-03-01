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
	"strconv"
	"strings"

	"github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/config-parser/v3"
	"github.com/haproxytech/config-parser/v3/common"
	parser_errors "github.com/haproxytech/config-parser/v3/errors"
	"github.com/haproxytech/config-parser/v3/parsers/http/actions"
	"github.com/haproxytech/config-parser/v3/types"

	"github.com/haproxytech/client-native/v2/models"
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
		return v, nil, c.HandleError("", parentType, parentName, "", false, err)
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
		return v, nil, c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, "", false, err)
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
		return c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	if err := c.SaveData(p, t, transactionID == ""); err != nil {
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
		return c.HandleError(strconv.FormatInt(*data.Index, 10), parentType, parentName, t, transactionID == "", err)
	}

	if err := c.SaveData(p, t, transactionID == ""); err != nil {
		return err
	}
	return nil
}

// EditHTTPRequestRule edits a http request rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
// nolint:dupl
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

	if _, err = p.GetOne(section, parentName, "http-request", int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	s, err := SerializeHTTPRequestRule(*data)
	if err != nil {
		return err
	}

	if err := p.Set(section, parentName, "http-request", s, int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	if err := c.SaveData(p, t, transactionID == ""); err != nil {
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
		if errors.Is(err, parser_errors.ErrFetch) {
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

func ParseHTTPRequestRule(f types.HTTPAction) (rule *models.HTTPRequestRule, err error) { //nolint:gocyclo
	switch v := f.(type) {
	case *actions.Allow:
		rule = &models.HTTPRequestRule{
			Type:     "allow",
			Cond:     v.Cond,
			CondTest: v.CondTest,
		}
	case *actions.Deny:
		var denyPtr *int64
		var ds int64
		if ds, err = strconv.ParseInt(v.DenyStatus, 10, 64); err == nil {
			denyPtr = &ds
		}
		rule = &models.HTTPRequestRule{
			Type:       "deny",
			Cond:       v.Cond,
			CondTest:   v.CondTest,
			DenyStatus: denyPtr,
		}
	case *actions.Auth:
		rule = &models.HTTPRequestRule{
			Type:      "auth",
			AuthRealm: v.Realm,
			Cond:      v.Cond,
			CondTest:  v.CondTest,
		}
	case *actions.Redirect:
		var codePtr *int64
		var code int64
		if code, err = strconv.ParseInt(v.Code, 10, 64); err == nil {
			codePtr = &code
		}
		rule = &models.HTTPRequestRule{
			Type:        "redirect",
			RedirType:   v.Type,
			RedirValue:  v.Value,
			RedirOption: v.Option,
			Cond:        v.Cond,
			CondTest:    v.CondTest,
			RedirCode:   codePtr,
		}

	case *actions.Tarpit:
		var dsPtr *int64
		var ds int64
		if ds, err = strconv.ParseInt(v.DenyStatus, 10, 64); err == nil {
			dsPtr = &ds
		}
		rule = &models.HTTPRequestRule{
			Type:       "tarpit",
			Cond:       v.Cond,
			CondTest:   v.CondTest,
			DenyStatus: dsPtr,
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
	case *actions.ReplacePath:
		rule = &models.HTTPRequestRule{
			Type:      "replace-path",
			PathMatch: v.MatchRegex,
			PathFmt:   v.ReplaceFmt,
			Cond:      v.Cond,
			CondTest:  v.CondTest,
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
	case *actions.SetMap:
		rule = &models.HTTPRequestRule{
			Type:        "set-map",
			MapFile:     v.FileName,
			MapKeyfmt:   v.KeyFmt,
			MapValuefmt: v.ValueFmt,
			Cond:        v.Cond,
			CondTest:    v.CondTest,
		}
	case *actions.DelMap:
		rule = &models.HTTPRequestRule{
			Type:      "del-map",
			MapFile:   v.FileName,
			MapKeyfmt: v.KeyFmt,
			Cond:      v.Cond,
			CondTest:  v.CondTest,
		}
	case *actions.CacheUse:
		rule = &models.HTTPRequestRule{
			Type:      "cache-use",
			CacheName: v.Name,
			Cond:      v.Cond,
			CondTest:  v.CondTest,
		}
	case *actions.DisableL7Retry:
		rule = &models.HTTPRequestRule{
			Type:     "disable-l7-retry",
			Cond:     v.Cond,
			CondTest: v.CondTest,
		}
	case *actions.EarlyHint:
		rule = &models.HTTPRequestRule{
			Type:       "early-hint",
			HintName:   v.Name,
			HintFormat: v.Fmt,
			Cond:       v.Cond,
			CondTest:   v.CondTest,
		}
	case *actions.ReplaceURI:
		rule = &models.HTTPRequestRule{
			Type:     "replace-uri",
			URIMatch: v.MatchRegex,
			URIFmt:   v.ReplaceFmt,
			Cond:     v.Cond,
			CondTest: v.CondTest,
		}
	case *actions.ScIncGpc0:
		ID, _ := strconv.ParseInt(v.ID, 10, 64)
		rule = &models.HTTPRequestRule{
			Type:     "sc-inc-gpc0",
			ScID:     ID,
			Cond:     v.Cond,
			CondTest: v.CondTest,
		}
	case *actions.ScIncGpc1:
		ID, _ := strconv.ParseInt(v.ID, 10, 64)
		rule = &models.HTTPRequestRule{
			Type:     "sc-inc-gpc1",
			ScID:     ID,
			Cond:     v.Cond,
			CondTest: v.CondTest,
		}
	case *actions.DoResolve:
		rule = &models.HTTPRequestRule{
			Type:      "do-resolve",
			VarName:   v.Var,
			Resolvers: v.Resolvers,
			Protocol:  v.Protocol,
			Expr:      v.Expr.String(),
			Cond:      v.Cond,
			CondTest:  v.CondTest,
		}
	case *actions.SetDst:
		rule = &models.HTTPRequestRule{
			Type:     "set-dst",
			Expr:     v.Expr.String(),
			Cond:     v.Cond,
			CondTest: v.CondTest,
		}
	case *actions.SetDstPort:
		rule = &models.HTTPRequestRule{
			Type:     "set-dst-port",
			Expr:     v.Expr.String(),
			Cond:     v.Cond,
			CondTest: v.CondTest,
		}
	case *actions.ScSetGpt0:
		if v.Int == nil && len(v.Expr.Expr) == 0 {
			return nil, NewConfError(ErrValidationError, "sc-set-gpt0 int or expr has to be set")
		}
		if v.Int != nil && len(v.Expr.Expr) > 0 {
			return nil, NewConfError(ErrValidationError, "sc-set-gpt0 int and expr are exclusive")
		}
		ID, _ := strconv.ParseInt(v.ID, 10, 64)
		rule = &models.HTTPRequestRule{
			Type:     "sc-set-gpt0",
			ScID:     ID,
			ScExpr:   strings.Join(v.Expr.Expr, " "),
			ScInt:    v.Int,
			Cond:     v.Cond,
			CondTest: v.CondTest,
		}
	case *actions.SetMark:
		rule = &models.HTTPRequestRule{
			Type:      "set-mark",
			MarkValue: v.Value,
			Cond:      v.Cond,
			CondTest:  v.CondTest,
		}
	case *actions.SetNice:
		nice, _ := strconv.ParseInt(v.Value, 10, 64)
		rule = &models.HTTPRequestRule{
			Type:      "set-nice",
			NiceValue: nice,
			Cond:      v.Cond,
			CondTest:  v.CondTest,
		}
	case *actions.SetMethod:
		rule = &models.HTTPRequestRule{
			Type:      "set-method",
			MethodFmt: v.Fmt,
			Cond:      v.Cond,
			CondTest:  v.CondTest,
		}
	case *actions.SetPriorityClass:
		rule = &models.HTTPRequestRule{
			Type:     "set-priority-class",
			Expr:     strings.Join(v.Expr.Expr, " "),
			Cond:     v.Cond,
			CondTest: v.CondTest,
		}
	case *actions.SetPriorityOffset:
		rule = &models.HTTPRequestRule{
			Type:     "set-priority-offset",
			Expr:     strings.Join(v.Expr.Expr, " "),
			Cond:     v.Cond,
			CondTest: v.CondTest,
		}
	case *actions.SetSrc:
		rule = &models.HTTPRequestRule{
			Type:     "set-src",
			Expr:     v.Expr.String(),
			Cond:     v.Cond,
			CondTest: v.CondTest,
		}
	case *actions.SetSrcPort:
		rule = &models.HTTPRequestRule{
			Type:     "set-src-port",
			Expr:     v.Expr.String(),
			Cond:     v.Cond,
			CondTest: v.CondTest,
		}
	case *actions.WaitForHandshake:
		rule = &models.HTTPRequestRule{
			Type:     "wait-for-handshake",
			Cond:     v.Cond,
			CondTest: v.CondTest,
		}
	case *actions.SetTos:
		rule = &models.HTTPRequestRule{
			Type:     "set-tos",
			TosValue: v.Value,
			Cond:     v.Cond,
			CondTest: v.CondTest,
		}
	case *actions.SilentDrop:
		rule = &models.HTTPRequestRule{
			Type:     "silent-drop",
			Cond:     v.Cond,
			CondTest: v.CondTest,
		}
	case *actions.UnsetVar:
		rule = &models.HTTPRequestRule{
			Type:     "unset-var",
			VarName:  v.Name,
			VarScope: v.Scope,
			Cond:     v.Cond,
			CondTest: v.CondTest,
		}
	case *actions.StrictMode:
		rule = &models.HTTPRequestRule{
			Type:       "strict-mode",
			StrictMode: v.Mode,
			Cond:       v.Cond,
			CondTest:   v.CondTest,
		}
	case *actions.Lua:
		rule = &models.HTTPRequestRule{
			Type:      "lua",
			LuaAction: v.Action,
			LuaParams: v.Params,
			Cond:      v.Cond,
			CondTest:  v.CondTest,
		}
	case *actions.UseService:
		rule = &models.HTTPRequestRule{
			Type:        "use-service",
			ServiceName: v.Name,
			Cond:        v.Cond,
			CondTest:    v.CondTest,
		}
	case *actions.Return:
		rule = &models.HTTPRequestRule{
			Cond:                v.Cond,
			CondTest:            v.CondTest,
			ReturnHeaders:       actionHdr2ModelHdr(v.Hdrs),
			ReturnContent:       v.Content,
			ReturnContentFormat: v.ContentFormat,
			ReturnContentType:   &v.ContentType,
			ReturnStatusCode:    v.Status,
			Type:                "return",
		}
	}

	return rule, err
}

func actionHdr2ModelHdr(hdrs []*actions.Hdr) []*models.HTTPRequestRuleReturnHdrsItems0 {
	if len(hdrs) == 0 {
		return nil
	}
	headers := []*models.HTTPRequestRuleReturnHdrsItems0{}
	for _, h := range hdrs {
		hdr := models.HTTPRequestRuleReturnHdrsItems0{
			Fmt:  &h.Fmt,
			Name: &h.Name,
		}
		headers = append(headers, &hdr)
	}
	return headers
}

func modelHdr2ActionHdr(hdrs []*models.HTTPRequestRuleReturnHdrsItems0) []*actions.Hdr {
	if len(hdrs) == 0 {
		return nil
	}
	headers := []*actions.Hdr{}
	for _, h := range hdrs {
		hdr := actions.Hdr{
			Name: *h.Name,
			Fmt:  *h.Fmt,
		}
		headers = append(headers, &hdr)
	}
	return headers
}

func SerializeHTTPRequestRule(f models.HTTPRequestRule) (rule types.HTTPAction, err error) { //nolint:gocyclo,gocognit
	switch f.Type {
	case "allow":
		rule = &actions.Allow{
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "deny":
		ds := ""
		if f.DenyStatus != nil {
			ds = strconv.FormatInt(*f.DenyStatus, 10)
		}
		rule = &actions.Deny{
			DenyStatus: ds,
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
		code := ""
		if f.RedirCode != nil {
			code = strconv.FormatInt(*f.RedirCode, 10)
		}
		rule = &actions.Redirect{
			Type:     f.RedirType,
			Value:    f.RedirValue,
			Code:     code,
			Option:   f.RedirOption,
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "tarpit":
		ds := ""
		if f.DenyStatus != nil {
			ds = strconv.FormatInt(*f.DenyStatus, 10)
		}
		rule = &actions.Tarpit{
			DenyStatus: ds,
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
	case "replace-path":
		rule = &actions.ReplacePath{
			MatchRegex: f.PathMatch,
			ReplaceFmt: f.PathFmt,
			Cond:       f.Cond,
			CondTest:   f.CondTest,
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
			Key:      f.TrackSc0Key,
			Table:    f.TrackSc0Table,
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "track-sc1":
		rule = &actions.TrackSc1{
			Key:      f.TrackSc1Key,
			Table:    f.TrackSc1Table,
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "track-sc2":
		rule = &actions.TrackSc2{
			Key:      f.TrackSc2Key,
			Table:    f.TrackSc2Table,
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "set-map":
		rule = &actions.SetMap{
			FileName: f.MapFile,
			KeyFmt:   f.MapKeyfmt,
			ValueFmt: f.MapValuefmt,
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "del-map":
		rule = &actions.DelMap{
			FileName: f.MapFile,
			KeyFmt:   f.MapKeyfmt,
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "cache-use":
		rule = &actions.CacheUse{
			Name:     f.CacheName,
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "disable-l7-retry":
		rule = &actions.DisableL7Retry{
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "early-hint":
		rule = &actions.EarlyHint{
			Name:     f.HintName,
			Fmt:      f.HintFormat,
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "replace-uri":
		rule = &actions.ReplaceURI{
			ReplaceFmt: f.URIFmt,
			MatchRegex: f.URIMatch,
			Cond:       f.Cond,
			CondTest:   f.CondTest,
		}
	case "sc-inc-gpc0":
		rule = &actions.ScIncGpc0{
			ID:       strconv.FormatInt(f.ScID, 10),
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "sc-inc-gpc1":
		rule = &actions.ScIncGpc1{
			ID:       strconv.FormatInt(f.ScID, 10),
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "do-resolve":
		rule = &actions.DoResolve{
			Var:       f.VarName,
			Resolvers: f.Resolvers,
			Protocol:  f.Protocol,
			Expr:      common.Expression{Expr: strings.Split(f.Expr, " ")},
			Cond:      f.Cond,
			CondTest:  f.CondTest,
		}
	case "set-dst":
		rule = &actions.SetDst{
			Expr:     common.Expression{Expr: strings.Split(f.Expr, " ")},
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "set-dst-port":
		rule = &actions.SetDstPort{
			Expr:     common.Expression{Expr: strings.Split(f.Expr, " ")},
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "set-method":
		rule = &actions.SetMethod{
			Fmt:      f.MethodFmt,
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "set-priority-class":
		rule = &actions.SetPriorityClass{
			Expr:     common.Expression{Expr: strings.Split(f.Expr, " ")},
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "sc-set-gpt0":
		if len(f.ScExpr) > 0 && f.ScInt != nil {
			return nil, NewConfError(ErrValidationError, "sc-set-gpt0 int and expr are exclusive")
		}
		if len(f.ScExpr) == 0 && f.ScInt == nil {
			return nil, NewConfError(ErrValidationError, "sc-set-gpt0 int or expr has to be set")
		}
		rule = &actions.ScSetGpt0{
			ID:       strconv.FormatInt(f.ScID, 10),
			Int:      f.ScInt,
			Expr:     common.Expression{Expr: strings.Split(f.ScExpr, " ")},
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "set-mark":
		rule = &actions.SetMark{
			Value:    f.MarkValue,
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "set-nice":
		rule = &actions.SetNice{
			Value:    strconv.FormatInt(f.NiceValue, 10),
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "set-priority-offset":
		rule = &actions.SetPriorityOffset{
			Expr:     common.Expression{Expr: strings.Split(f.Expr, " ")},
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "set-src":
		rule = &actions.SetSrc{
			Expr:     common.Expression{Expr: strings.Split(f.Expr, " ")},
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "set-src-port":
		rule = &actions.SetSrcPort{
			Expr:     common.Expression{Expr: strings.Split(f.Expr, " ")},
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "wait-for-handshake":
		rule = &actions.WaitForHandshake{
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "set-tos":
		rule = &actions.SetTos{
			Value:    f.TosValue,
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "silent-drop":
		rule = &actions.SilentDrop{
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
	case "strict-mode":
		rule = &actions.StrictMode{
			Mode:     f.StrictMode,
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "lua":
		rule = &actions.Lua{
			Action:   f.LuaAction,
			Params:   f.LuaParams,
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "use-service":
		rule = &actions.UseService{
			Name:     f.ServiceName,
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "return":
		rule = &actions.Return{
			Status:        f.ReturnStatusCode,
			ContentType:   *f.ReturnContentType,
			ContentFormat: f.ReturnContentFormat,
			Content:       f.ReturnContent,
			Hdrs:          modelHdr2ActionHdr(f.ReturnHeaders),
			Cond:          f.Cond,
			CondTest:      f.CondTest,
		}
		if !actions.IsPayload(f.ReturnContentFormat) {
			if ok := actions.AllowedErrorCode(*f.ReturnStatusCode); !ok {
				return rule, NewConfError(ErrValidationError, "invalid Status Code for error type response")
			}
		}
	}

	return rule, err
}
