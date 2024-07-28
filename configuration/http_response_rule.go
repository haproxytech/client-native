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

	strfmt "github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/config-parser/v5"
	"github.com/haproxytech/config-parser/v5/common"
	parser_errors "github.com/haproxytech/config-parser/v5/errors"
	"github.com/haproxytech/config-parser/v5/parsers/actions"
	http_actions "github.com/haproxytech/config-parser/v5/parsers/http/actions"
	"github.com/haproxytech/config-parser/v5/types"

	"github.com/haproxytech/client-native/v6/configuration/options"
	"github.com/haproxytech/client-native/v6/misc"
	"github.com/haproxytech/client-native/v6/models"
)

type HTTPResponseRule interface {
	GetHTTPResponseRules(parentType, parentName string, transactionID string) (int64, models.HTTPResponseRules, error)
	GetHTTPResponseRule(id int64, parentType, parentName string, transactionID string) (int64, *models.HTTPResponseRule, error)
	DeleteHTTPResponseRule(id int64, parentType string, parentName string, transactionID string, version int64) error
	CreateHTTPResponseRule(id int64, parentType string, parentName string, data *models.HTTPResponseRule, transactionID string, version int64) error
	EditHTTPResponseRule(id int64, parentType string, parentName string, data *models.HTTPResponseRule, transactionID string, version int64) error
	ReplaceHTTPResponseRules(parentType string, parentName string, data models.HTTPResponseRules, transactionID string, version int64) error
}

// GetHTTPResponseRules returns configuration version and an array of
// configured http response rules in the specified parent. Returns error on fail.
func (c *client) GetHTTPResponseRules(parentType, parentName string, transactionID string) (int64, models.HTTPResponseRules, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	httpRules, err := ParseHTTPResponseRules(parentType, parentName, p)
	if err != nil {
		return v, nil, c.HandleError("", parentType, parentName, "", false, err)
	}

	return v, httpRules, nil
}

// GetHTTPResponseRule returns configuration version and a response http response rule
// in the specified parent. Returns error on fail or if http response rule does not exist.
func (c *client) GetHTTPResponseRule(id int64, parentType, parentName string, transactionID string) (int64, *models.HTTPResponseRule, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	section, parentName, err := getParserFromParent("http-response", parentType, parentName)
	if err != nil {
		return v, nil, err
	}

	data, err := p.GetOne(section, parentName, "http-response", int(id))
	if err != nil {
		return v, nil, c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, "", false, err)
	}

	httpRule := ParseHTTPResponseRule(data.(types.Action))

	return v, httpRule, nil
}

// DeleteHTTPResponseRule deletes a http response rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) DeleteHTTPResponseRule(id int64, parentType string, parentName string, transactionID string, version int64) error {
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	section, parentName, err := getParserFromParent("http-response", parentType, parentName)
	if err != nil {
		return err
	}

	if err := p.Delete(section, parentName, "http-response", int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

// CreateHTTPResponseRule creates a http response rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) CreateHTTPResponseRule(id int64, parentType string, parentName string, data *models.HTTPResponseRule, transactionID string, version int64) error {
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

	section, parentName, err := getParserFromParent("http-response", parentType, parentName)
	if err != nil {
		return err
	}

	s, err := SerializeHTTPResponseRule(*data, &c.ConfigurationOptions)
	if err != nil {
		return err
	}
	if err := p.Insert(section, parentName, "http-response", s, int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

// EditHTTPResponseRule edits a http response rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) EditHTTPResponseRule(id int64, parentType string, parentName string, data *models.HTTPResponseRule, transactionID string, version int64) error {
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

	section, parentName, err := getParserFromParent("http-response", parentType, parentName)
	if err != nil {
		return err
	}

	if _, err = p.GetOne(section, parentName, "http-response", int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	s, err := SerializeHTTPResponseRule(*data, &c.ConfigurationOptions)
	if err != nil {
		return err
	}
	if err := p.Set(section, parentName, "http-response", s, int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

func ParseHTTPResponseRules(t, pName string, p parser.Parser) (models.HTTPResponseRules, error) {
	section, pName, err := getParserFromParent("http-response", t, pName)
	if err != nil {
		return nil, err
	}

	httpResRules := models.HTTPResponseRules{}

	data, err := p.Get(section, pName, "http-response", false)
	if err != nil {
		if goerrors.Is(err, parser_errors.ErrFetch) {
			return httpResRules, nil
		}
		return nil, err
	}

	rules, ok := data.([]types.Action)
	if !ok {
		return nil, misc.CreateTypeAssertError("http-response")
	}
	for _, r := range rules {
		httpResRule := ParseHTTPResponseRule(r)
		if httpResRule != nil {
			httpResRules = append(httpResRules, httpResRule)
		}
	}
	return httpResRules, nil
}

// ReplaceHTTPResponseRules replaces all HTTP Response Rule lines in configuration for a parentType/parentName.
// One of version or transactionID is mandatory.
// Returns error on fail, nil on success.
//
//nolint:dupl
func (c *client) ReplaceHTTPResponseRules(parentType string, parentName string, data models.HTTPResponseRules, transactionID string, version int64) error {
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

	section, parentName, err := getParserFromParent("http-response", parentType, parentName)
	if err != nil {
		return err
	}

	httpResponseRules, err := ParseHTTPResponseRules(parentType, parentName, p)
	if err != nil {
		return c.HandleError("", parentType, parentName, "", false, err)
	}

	for i := range httpResponseRules {
		// Always delete index 0
		if err := p.Delete(section, parentName, "http-response", 0); err != nil {
			return c.HandleError(strconv.FormatInt(int64(i), 10), parentType, parentName, t, transactionID == "", err)
		}
	}

	for i, newHTTPResponseRule := range data {
		s, err := SerializeHTTPResponseRule(*newHTTPResponseRule, &c.ConfigurationOptions)
		if err != nil {
			return err
		}
		if err := p.Insert(section, parentName, "http-response", s, i); err != nil {
			return c.HandleError(strconv.FormatInt(int64(i), 10), parentType, parentName, t, transactionID == "", err)
		}
	}

	return c.SaveData(p, t, transactionID == "")
}

func ParseHTTPResponseRule(f types.Action) *models.HTTPResponseRule { //nolint:maintidx,cyclop,gocyclo
	switch v := f.(type) {
	case *http_actions.AddACL:
		return &models.HTTPResponseRule{
			Type:      "add-acl",
			ACLFile:   v.FileName,
			ACLKeyfmt: v.KeyFmt,
			Cond:      v.Cond,
			CondTest:  v.CondTest,
		}
	case *http_actions.AddHeader:
		return &models.HTTPResponseRule{
			Type:      "add-header",
			HdrName:   v.Name,
			HdrFormat: v.Fmt,
			Cond:      v.Cond,
			CondTest:  v.CondTest,
		}
	case *http_actions.Allow:
		return &models.HTTPResponseRule{
			Type:     "allow",
			Cond:     v.Cond,
			CondTest: v.CondTest,
		}
	case *http_actions.CacheStore:
		return &models.HTTPResponseRule{
			Type:      "cache-store",
			CacheName: v.Name,
			Cond:      v.Cond,
			CondTest:  v.CondTest,
		}
	case *http_actions.Capture:
		return &models.HTTPResponseRule{
			Type:          "capture",
			CaptureSample: v.Sample,
			Cond:          v.Cond,
			CondTest:      v.CondTest,
			CaptureID:     v.SlotID,
		}
	case *http_actions.DelACL:
		return &models.HTTPResponseRule{
			Type:      "del-acl",
			ACLFile:   v.FileName,
			ACLKeyfmt: v.KeyFmt,
			Cond:      v.Cond,
			CondTest:  v.CondTest,
		}
	case *http_actions.DelHeader:
		return &models.HTTPResponseRule{
			Type:      "del-header",
			HdrName:   v.Name,
			Cond:      v.Cond,
			CondTest:  v.CondTest,
			HdrMethod: v.Method,
		}
	case *http_actions.DelMap:
		return &models.HTTPResponseRule{
			Type:      "del-map",
			MapFile:   v.FileName,
			MapKeyfmt: v.KeyFmt,
			Cond:      v.Cond,
			CondTest:  v.CondTest,
		}
	case *http_actions.Deny:
		return &models.HTTPResponseRule{
			Type:                "deny",
			Cond:                v.Cond,
			CondTest:            v.CondTest,
			ReturnHeaders:       actionHdr2ModelHdr(v.Hdrs),
			ReturnContent:       v.Content,
			ReturnContentFormat: v.ContentFormat,
			ReturnContentType:   &v.ContentType,
			DenyStatus:          v.Status,
		}
	case *actions.Lua:
		return &models.HTTPResponseRule{
			Type:      "lua",
			LuaAction: v.Action,
			LuaParams: v.Params,
			Cond:      v.Cond,
			CondTest:  v.CondTest,
		}
	case *http_actions.Redirect:
		var codePtr *int64
		if code, err := strconv.ParseInt(v.Code, 10, 64); err == nil {
			codePtr = &code
		}
		return &models.HTTPResponseRule{
			Type:        "redirect",
			RedirType:   v.Type,
			RedirValue:  v.Value,
			RedirOption: v.Option,
			Cond:        v.Cond,
			CondTest:    v.CondTest,
			RedirCode:   codePtr,
		}
	case *http_actions.ReplaceHeader:
		return &models.HTTPResponseRule{
			Type:      "replace-header",
			HdrName:   v.Name,
			HdrFormat: v.ReplaceFmt,
			HdrMatch:  v.MatchRegex,
			Cond:      v.Cond,
			CondTest:  v.CondTest,
		}
	case *http_actions.ReplaceValue:
		return &models.HTTPResponseRule{
			Type:      "replace-value",
			HdrName:   v.Name,
			HdrFormat: v.ReplaceFmt,
			HdrMatch:  v.MatchRegex,
			Cond:      v.Cond,
			CondTest:  v.CondTest,
		}
	case *http_actions.Return:
		return &models.HTTPResponseRule{
			Type:                "return",
			Cond:                v.Cond,
			CondTest:            v.CondTest,
			ReturnHeaders:       actionHdr2ModelHdr(v.Hdrs),
			ReturnContent:       v.Content,
			ReturnContentFormat: v.ContentFormat,
			ReturnContentType:   &v.ContentType,
			ReturnStatusCode:    v.Status,
		}
	case *actions.ScAddGpc:
		if (v.Int == nil && len(v.Expr.Expr) == 0) || (v.Int != nil && len(v.Expr.Expr) > 0) {
			return nil
		}
		ID, _ := strconv.ParseInt(v.ID, 10, 64)
		Idx, _ := strconv.ParseInt(v.Idx, 10, 64)
		return &models.HTTPResponseRule{
			Type:     "sc-add-gpc",
			ScID:     ID,
			ScIdx:    Idx,
			ScExpr:   strings.Join(v.Expr.Expr, " "),
			ScInt:    v.Int,
			Cond:     v.Cond,
			CondTest: v.CondTest,
		}
	case *actions.ScIncGpc:
		ID, _ := strconv.ParseInt(v.ID, 10, 64)
		Idx, _ := strconv.ParseInt(v.Idx, 10, 64)
		return &models.HTTPResponseRule{
			Type:     "sc-inc-gpc",
			ScID:     ID,
			ScIdx:    Idx,
			Cond:     v.Cond,
			CondTest: v.CondTest,
		}
	case *actions.ScIncGpc0:
		ID, _ := strconv.ParseInt(v.ID, 10, 64)
		return &models.HTTPResponseRule{
			Type:     "sc-inc-gpc0",
			ScID:     ID,
			Cond:     v.Cond,
			CondTest: v.CondTest,
		}
	case *actions.ScIncGpc1:
		ID, _ := strconv.ParseInt(v.ID, 10, 64)
		return &models.HTTPResponseRule{
			Type:     "sc-inc-gpc1",
			ScID:     ID,
			Cond:     v.Cond,
			CondTest: v.CondTest,
		}
	case *actions.ScSetGpt:
		if (v.Int == nil && len(v.Expr.Expr) == 0) || (v.Int != nil && len(v.Expr.Expr) > 0) {
			return nil
		}
		scID, err := strconv.ParseInt(v.ScID, 10, 64)
		if err != nil {
			return nil
		}
		return &models.HTTPResponseRule{
			Type:     "sc-set-gpt",
			ScID:     scID,
			ScIdx:    v.Idx,
			ScExpr:   strings.Join(v.Expr.Expr, " "),
			ScInt:    v.Int,
			Cond:     v.Cond,
			CondTest: v.CondTest,
		}
	case *actions.ScSetGpt0:
		if (v.Int == nil && len(v.Expr.Expr) == 0) || (v.Int != nil && len(v.Expr.Expr) > 0) {
			return nil
		}
		ID, _ := strconv.ParseInt(v.ID, 10, 64)
		return &models.HTTPResponseRule{
			Type:     "sc-set-gpt0",
			ScID:     ID,
			ScExpr:   strings.Join(v.Expr.Expr, " "),
			ScInt:    v.Int,
			Cond:     v.Cond,
			CondTest: v.CondTest,
		}
	case *actions.SendSpoeGroup:
		return &models.HTTPResponseRule{
			Type:       "send-spoe-group",
			SpoeEngine: v.Engine,
			SpoeGroup:  v.Group,
			Cond:       v.Cond,
			CondTest:   v.CondTest,
		}
	case *http_actions.SetHeader:
		return &models.HTTPResponseRule{
			Type:      "set-header",
			HdrName:   v.Name,
			HdrFormat: v.Fmt,
			Cond:      v.Cond,
			CondTest:  v.CondTest,
		}
	case *actions.SetLogLevel:
		return &models.HTTPResponseRule{
			Type:     "set-log-level",
			LogLevel: v.Level,
			Cond:     v.Cond,
			CondTest: v.CondTest,
		}
	case *http_actions.SetMap:
		return &models.HTTPResponseRule{
			Type:        "set-map",
			MapFile:     v.FileName,
			MapKeyfmt:   v.KeyFmt,
			MapValuefmt: v.ValueFmt,
			Cond:        v.Cond,
			CondTest:    v.CondTest,
		}
	case *actions.SetMark:
		return &models.HTTPResponseRule{
			Type:      "set-mark",
			MarkValue: v.Value,
			Cond:      v.Cond,
			CondTest:  v.CondTest,
		}
	case *actions.SetNice:
		nice, _ := strconv.ParseInt(v.Value, 10, 64)
		return &models.HTTPResponseRule{
			Type:      "set-nice",
			NiceValue: nice,
			Cond:      v.Cond,
			CondTest:  v.CondTest,
		}
	case *http_actions.SetStatus:
		status, _ := strconv.ParseInt(v.Status, 10, 64)
		r := &models.HTTPResponseRule{
			Type:         "set-status",
			StatusReason: v.Reason,
			Cond:         v.Cond,
			CondTest:     v.CondTest,
		}
		if status != 0 {
			r.Status = status
		}
		return r
	case *http_actions.SetTimeout:
		return &models.HTTPResponseRule{
			Type:        models.HTTPResponseRuleTypeSetDashTimeout,
			Timeout:     v.Timeout,
			TimeoutType: v.Type,
			Cond:        v.Cond,
			CondTest:    v.CondTest,
		}
	case *actions.SetTos:
		return &models.HTTPResponseRule{
			Type:     "set-tos",
			TosValue: v.Value,
			Cond:     v.Cond,
			CondTest: v.CondTest,
		}
	case *actions.SetVar:
		return &models.HTTPResponseRule{
			Type:     "set-var",
			VarName:  v.VarName,
			VarExpr:  strings.Join(v.Expr.Expr, " "),
			VarScope: v.VarScope,
			Cond:     v.Cond,
			CondTest: v.CondTest,
		}
	case *actions.SetVarFmt:
		return &models.HTTPResponseRule{
			Type:      "set-var-fmt",
			VarName:   v.VarName,
			VarFormat: strings.Join(v.Fmt.Expr, " "),
			VarScope:  v.VarScope,
			Cond:      v.Cond,
			CondTest:  v.CondTest,
		}
	case *actions.SilentDrop:
		return &models.HTTPResponseRule{
			Type:     "silent-drop",
			Cond:     v.Cond,
			CondTest: v.CondTest,
		}
	case *http_actions.StrictMode:
		return &models.HTTPResponseRule{
			Type:       "strict-mode",
			StrictMode: v.Mode,
			Cond:       v.Cond,
			CondTest:   v.CondTest,
		}
	case *actions.TrackSc:
		return &models.HTTPResponseRule{
			Type:                string(actions.TrackScType),
			TrackScKey:          v.Key,
			TrackScTable:        v.Table,
			Cond:                v.Cond,
			CondTest:            v.CondTest,
			TrackScStickCounter: &v.StickCounter,
		}
	case *actions.UnsetVar:
		return &models.HTTPResponseRule{
			Type:     "unset-var",
			VarName:  v.Name,
			VarScope: v.Scope,
			Cond:     v.Cond,
			CondTest: v.CondTest,
		}
	case *http_actions.WaitForBody:
		return &models.HTTPResponseRule{
			Type:        "wait-for-body",
			WaitTime:    misc.ParseTimeout(v.Time),
			WaitAtLeast: misc.ParseSize(v.AtLeast),
		}
	case *actions.SetBandwidthLimit:
		return &models.HTTPResponseRule{
			Type:                 "set-bandwidth-limit",
			BandwidthLimitName:   v.Name,
			BandwidthLimitLimit:  v.Limit.String(),
			BandwidthLimitPeriod: v.Period.String(),
			Cond:                 v.Cond,
			CondTest:             v.CondTest,
		}
	case *actions.SetFcMark:
		return &models.HTTPResponseRule{
			Type:     models.HTTPResponseRuleTypeSetDashFcDashMark,
			Expr:     v.Expr.String(),
			Cond:     v.Cond,
			CondTest: v.CondTest,
		}
	case *actions.SetFcTos:
		return &models.HTTPResponseRule{
			Type:     models.HTTPResponseRuleTypeSetDashFcDashTos,
			Expr:     v.Expr.String(),
			Cond:     v.Cond,
			CondTest: v.CondTest,
		}
	}
	return nil
}

func SerializeHTTPResponseRule(f models.HTTPResponseRule, opt *options.ConfigurationOptions) (rule types.Action, err error) { //nolint:gocyclo,ireturn,cyclop,maintidx,gocognit
	switch f.Type {
	case "add-acl":
		rule = &http_actions.AddACL{
			FileName: f.ACLFile,
			KeyFmt:   f.ACLKeyfmt,
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
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
	case "cache-store":
		rule = &http_actions.CacheStore{
			Name:     f.CacheName,
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "capture":
		rule = &http_actions.Capture{
			Sample:   f.CaptureSample,
			Cond:     f.Cond,
			CondTest: f.CondTest,
			SlotID:   f.CaptureID,
		}
	case "del-acl":
		rule = &http_actions.DelACL{
			FileName: f.ACLFile,
			KeyFmt:   f.ACLKeyfmt,
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
	case "del-map":
		rule = &http_actions.DelMap{
			FileName: f.MapFile,
			KeyFmt:   f.MapKeyfmt,
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "deny":
		contentType := ""
		if f.ReturnContentType != nil {
			contentType = *f.ReturnContentType
		}
		rule = &http_actions.Deny{
			Status:        f.DenyStatus,
			ContentType:   contentType,
			ContentFormat: f.ReturnContentFormat,
			Content:       f.ReturnContent,
			Hdrs:          modelHdr2ActionHdr(f.ReturnHeaders),
			Cond:          f.Cond,
			CondTest:      f.CondTest,
		}
		if !http_actions.IsPayload(f.ReturnContentFormat) && f.DenyStatus != nil {
			if ok := http_actions.AllowedErrorCode(*f.DenyStatus); !ok {
				return rule, NewConfError(ErrValidationError, "invalid Status Code for error type response")
			}
		}
	case "lua":
		rule = &actions.Lua{
			Action:   f.LuaAction,
			Params:   f.LuaParams,
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "redirect":
		code := ""
		if f.RedirCode != nil {
			code = strconv.FormatInt(*f.RedirCode, 10)
		}
		rule = &http_actions.Redirect{
			Type:     f.RedirType,
			Value:    f.RedirValue,
			Code:     code,
			Option:   f.RedirOption,
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
	case "return":
		contentType := ""
		if f.ReturnContentType != nil {
			contentType = *f.ReturnContentType
		}
		rule = &http_actions.Return{
			Status:        f.ReturnStatusCode,
			ContentType:   contentType,
			ContentFormat: f.ReturnContentFormat,
			Content:       f.ReturnContent,
			Hdrs:          modelHdr2ActionHdr(f.ReturnHeaders),
			Cond:          f.Cond,
			CondTest:      f.CondTest,
		}
		if !http_actions.IsPayload(f.ReturnContentFormat) && f.ReturnStatusCode != nil {
			if ok := http_actions.AllowedErrorCode(*f.ReturnStatusCode); !ok {
				return rule, NewConfError(ErrValidationError, "invalid Status Code for error type response")
			}
		}
	case "sc-add-gpc":
		rule = &actions.ScAddGpc{
			ID:       strconv.FormatInt(f.ScID, 10),
			Idx:      strconv.FormatInt(f.ScIdx, 10),
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "sc-inc-gpc":
		rule = &actions.ScAddGpc{
			ID:       strconv.FormatInt(f.ScID, 10),
			Idx:      strconv.FormatInt(f.ScIdx, 10),
			Cond:     f.Cond,
			CondTest: f.CondTest,
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
	case "sc-set-gpt":
		if len(f.ScExpr) > 0 && f.ScInt != nil {
			return nil, NewConfError(ErrValidationError, "sc-set-gpt: int and expr are exclusive")
		}
		if len(f.ScExpr) == 0 && f.ScInt == nil {
			return nil, NewConfError(ErrValidationError, "sc-set-gpt: int or expr has to be set")
		}
		rule = &actions.ScSetGpt{
			ScID:     strconv.FormatInt(f.ScID, 10),
			Idx:      f.ScIdx,
			Int:      f.ScInt,
			Expr:     common.Expression{Expr: strings.Split(f.ScExpr, " ")},
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
	case "send-spoe-group":
		rule = &actions.SendSpoeGroup{
			Engine:   f.SpoeEngine,
			Group:    f.SpoeGroup,
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "set-header":
		rule = &http_actions.SetHeader{
			Name:     f.HdrName,
			Fmt:      f.HdrFormat,
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "set-log-level":
		rule = &actions.SetLogLevel{
			Level:    f.LogLevel,
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "set-map":
		rule = &http_actions.SetMap{
			FileName: f.MapFile,
			KeyFmt:   f.MapKeyfmt,
			ValueFmt: f.MapValuefmt,
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
	case "set-status":
		rule = &http_actions.SetStatus{
			Status:   strconv.FormatInt(f.Status, 10),
			Reason:   f.StatusReason,
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "set-timeout":
		rule = &http_actions.SetTimeout{
			Timeout:  f.Timeout,
			Type:     f.TimeoutType,
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "set-tos":
		rule = &actions.SetTos{
			Value:    f.TosValue,
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
	case "set-var-fmt":
		rule = &actions.SetVarFmt{
			Fmt:      common.Expression{Expr: strings.Split(f.VarFormat, " ")},
			VarName:  f.VarName,
			VarScope: f.VarScope,
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "silent-drop":
		rule = &actions.SilentDrop{
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "strict-mode":
		rule = &http_actions.StrictMode{
			Mode:     f.StrictMode,
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "track-sc":
		if f.TrackScStickCounter == nil {
			return nil, NewConfError(ErrValidationError, "track_sc_stick_counter must be set")
		}
		rule = &actions.TrackSc{
			StickCounter: *f.TrackScStickCounter,
			Type:         actions.TrackScType,
			Key:          f.TrackScKey,
			Table:        f.TrackScTable,
			Cond:         f.Cond,
			CondTest:     f.CondTest,
		}
	case "unset-var":
		rule = &actions.UnsetVar{
			Name:     f.VarName,
			Scope:    f.VarScope,
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "wait-for-body":
		rule = &http_actions.WaitForBody{
			Time:     misc.SerializeTime(*f.WaitTime, opt.PreferredTimeSuffix),
			AtLeast:  misc.SerializeSize(*f.WaitAtLeast),
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "set-bandwidth-limit":
		rule = &actions.SetBandwidthLimit{
			Name:     f.BandwidthLimitName,
			Limit:    common.Expression{Expr: strings.Split(f.BandwidthLimitLimit, " ")},
			Period:   common.Expression{Expr: strings.Split(f.BandwidthLimitPeriod, " ")},
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "set-fc-mark":
		rule = &actions.SetFcMark{
			Expr:     common.Expression{Expr: strings.Split(f.Expr+f.MarkValue, " ")},
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "set-fc-tos":
		rule = &actions.SetFcTos{
			Expr:     common.Expression{Expr: strings.Split(f.Expr+f.TosValue, " ")},
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	}
	return rule, err
}
