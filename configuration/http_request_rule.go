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
	parser "github.com/haproxytech/client-native/v6/config-parser"
	"github.com/haproxytech/client-native/v6/config-parser/common"
	parser_errors "github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/parsers/actions"
	http_actions "github.com/haproxytech/client-native/v6/config-parser/parsers/http/actions"
	"github.com/haproxytech/client-native/v6/config-parser/types"

	"github.com/haproxytech/client-native/v6/configuration/options"
	"github.com/haproxytech/client-native/v6/misc"
	"github.com/haproxytech/client-native/v6/models"
)

type HTTPRequestRule interface {
	GetHTTPRequestRules(parentType, parentName string, transactionID string) (int64, models.HTTPRequestRules, error)
	GetHTTPRequestRule(id int64, parentType, parentName string, transactionID string) (int64, *models.HTTPRequestRule, error)
	DeleteHTTPRequestRule(id int64, parentType string, parentName string, transactionID string, version int64) error
	CreateHTTPRequestRule(id int64, parentType string, parentName string, data *models.HTTPRequestRule, transactionID string, version int64) error
	EditHTTPRequestRule(id int64, parentType string, parentName string, data *models.HTTPRequestRule, transactionID string, version int64) error
	ReplaceHTTPRequestRules(parentType string, parentName string, data models.HTTPRequestRules, transactionID string, version int64) error
}

// GetHTTPRequestRules returns configuration version and an array of
// configured http request rules in the specified parent. Returns error on fail.
func (c *client) GetHTTPRequestRules(parentType, parentName string, transactionID string) (int64, models.HTTPRequestRules, error) {
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
func (c *client) GetHTTPRequestRule(id int64, parentType, parentName string, transactionID string) (int64, *models.HTTPRequestRule, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	section, parentName, err := getParserFromParent("http-request", parentType, parentName)
	if err != nil {
		return 0, nil, err
	}

	data, err := p.GetOne(section, parentName, "http-request", int(id))
	if err != nil {
		return v, nil, c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, "", false, err)
	}

	httpRule, err := ParseHTTPRequestRule(data.(types.Action))
	if err != nil {
		return v, nil, err
	}

	return v, httpRule, nil
}

// DeleteHTTPRequestRule deletes a http request rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) DeleteHTTPRequestRule(id int64, parentType string, parentName string, transactionID string, version int64) error {
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	section, parentName, err := getParserFromParent("http-request", parentType, parentName)
	if err != nil {
		return err
	}

	if err := p.Delete(section, parentName, "http-request", int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

// CreateHTTPRequestRule creates a http request rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) CreateHTTPRequestRule(id int64, parentType string, parentName string, data *models.HTTPRequestRule, transactionID string, version int64) error {
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

	section, parentName, err := getParserFromParent("http-request", parentType, parentName)
	if err != nil {
		return err
	}

	s, err := SerializeHTTPRequestRule(*data, &c.ConfigurationOptions)
	if err != nil {
		return err
	}

	if err := p.Insert(section, parentName, "http-request", s, int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

// EditHTTPRequestRule edits a http request rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) EditHTTPRequestRule(id int64, parentType string, parentName string, data *models.HTTPRequestRule, transactionID string, version int64) error {
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

	section, parentName, err := getParserFromParent("http-request", parentType, parentName)
	if err != nil {
		return err
	}

	if _, err = p.GetOne(section, parentName, "http-request", int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	s, err := SerializeHTTPRequestRule(*data, &c.ConfigurationOptions)
	if err != nil {
		return err
	}

	if err := p.Set(section, parentName, "http-request", s, int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

// ReplaceHTTPRequestRules replaces all HTTP Request Rule lines in configuration for a parentType/parentName.
// One of version or transactionID is mandatory.
// Returns error on fail, nil on success.
//
//nolint:dupl
func (c *client) ReplaceHTTPRequestRules(parentType string, parentName string, data models.HTTPRequestRules, transactionID string, version int64) error {
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

	section, parentName, err := getParserFromParent("http-request", parentType, parentName)
	if err != nil {
		return err
	}

	httpRequestRules, err := ParseHTTPRequestRules(parentType, parentName, p)
	if err != nil {
		return c.HandleError("", parentType, parentName, "", false, err)
	}

	for i := range httpRequestRules {
		// Always delete index 0
		if err := p.Delete(section, parentName, "http-request", 0); err != nil {
			return c.HandleError(strconv.FormatInt(int64(i), 10), parentType, parentName, t, transactionID == "", err)
		}
	}

	for i, newHTTPRequestRule := range data {
		s, err := SerializeHTTPRequestRule(*newHTTPRequestRule, &c.ConfigurationOptions)
		if err != nil {
			return err
		}
		if err := p.Insert(section, parentName, "http-request", s, i); err != nil {
			return c.HandleError(strconv.FormatInt(int64(i), 10), parentType, parentName, t, transactionID == "", err)
		}
	}

	return c.SaveData(p, t, transactionID == "")
}

func ParseHTTPRequestRules(t, pName string, p parser.Parser) (models.HTTPRequestRules, error) {
	section, pName, err := getParserFromParent("http-request", t, pName)
	if err != nil {
		return nil, err
	}

	var httpReqRules models.HTTPRequestRules
	data, err := p.Get(section, pName, "http-request", false)
	if err != nil {
		if errors.Is(err, parser_errors.ErrFetch) {
			return httpReqRules, nil
		}
		return nil, err
	}

	rules, ok := data.([]types.Action)
	if !ok {
		return nil, misc.CreateTypeAssertError("http request")
	}
	for _, r := range rules {
		httpReqRule, err := ParseHTTPRequestRule(r)
		if err == nil {
			httpReqRules = append(httpReqRules, httpReqRule)
		}
	}
	return httpReqRules, nil
}

func ParseHTTPRequestRule(f types.Action) (*models.HTTPRequestRule, error) { //nolint:gocyclo,cyclop,maintidx,gocognit
	var rule *models.HTTPRequestRule
	switch v := f.(type) {
	case *http_actions.AddACL:
		rule = &models.HTTPRequestRule{
			Type:      "add-acl",
			ACLFile:   v.FileName,
			ACLKeyfmt: v.KeyFmt,
			Cond:      v.Cond,
			CondTest:  v.CondTest,
			Metadata:  parseMetadata(v.Comment),
		}
	case *http_actions.AddHeader:
		rule = &models.HTTPRequestRule{
			Type:      "add-header",
			HdrName:   v.Name,
			HdrFormat: v.Fmt,
			Cond:      v.Cond,
			CondTest:  v.CondTest,
			Metadata:  parseMetadata(v.Comment),
		}
	case *http_actions.Allow:
		rule = &models.HTTPRequestRule{
			Type:     "allow",
			Cond:     v.Cond,
			CondTest: v.CondTest,
			Metadata: parseMetadata(v.Comment),
		}
	case *http_actions.Auth:
		rule = &models.HTTPRequestRule{
			Type:      "auth",
			AuthRealm: v.Realm,
			Cond:      v.Cond,
			CondTest:  v.CondTest,
			Metadata:  parseMetadata(v.Comment),
		}
	case *http_actions.CacheUse:
		rule = &models.HTTPRequestRule{
			Type:      "cache-use",
			CacheName: v.Name,
			Cond:      v.Cond,
			CondTest:  v.CondTest,
			Metadata:  parseMetadata(v.Comment),
		}
	case *http_actions.Capture:
		if (v.SlotID == nil && v.Len == nil) || (v.SlotID != nil && v.Len != nil) {
			return nil, NewConfError(ErrValidationError, "capture len can't be zero")
		}
		rule = &models.HTTPRequestRule{
			Type:          "capture",
			CaptureSample: v.Sample,
			Cond:          v.Cond,
			CondTest:      v.CondTest,
			Metadata:      parseMetadata(v.Comment),
		}
		if v.SlotID != nil {
			rule.CaptureID = v.SlotID
		}
		if v.Len != nil {
			rule.CaptureLen = *v.Len
		}
	case *http_actions.DelACL:
		rule = &models.HTTPRequestRule{
			Type:      "del-acl",
			ACLFile:   v.FileName,
			ACLKeyfmt: v.KeyFmt,
			Cond:      v.Cond,
			CondTest:  v.CondTest,
			Metadata:  parseMetadata(v.Comment),
		}
	case *http_actions.DelHeader:
		rule = &models.HTTPRequestRule{
			Type:      "del-header",
			HdrName:   v.Name,
			HdrMethod: v.Method,
			Cond:      v.Cond,
			CondTest:  v.CondTest,
			Metadata:  parseMetadata(v.Comment),
		}
	case *http_actions.DelMap:
		rule = &models.HTTPRequestRule{
			Type:      "del-map",
			MapFile:   v.FileName,
			MapKeyfmt: v.KeyFmt,
			Cond:      v.Cond,
			CondTest:  v.CondTest,
			Metadata:  parseMetadata(v.Comment),
		}
	case *http_actions.Deny:
		var returnContentTypePtr *string
		if v.ContentType != "" {
			returnContentTypePtr = &v.ContentType
		}
		rule = &models.HTTPRequestRule{
			Type:                "deny",
			Cond:                v.Cond,
			CondTest:            v.CondTest,
			ReturnHeaders:       actionHdr2ModelHdr(v.Hdrs),
			ReturnContent:       v.Content,
			ReturnContentFormat: v.ContentFormat,
			ReturnContentType:   returnContentTypePtr,
			DenyStatus:          v.Status,
			Metadata:            parseMetadata(v.Comment),
		}
	case *http_actions.DisableL7Retry:
		rule = &models.HTTPRequestRule{
			Type:     "disable-l7-retry",
			Cond:     v.Cond,
			CondTest: v.CondTest,
			Metadata: parseMetadata(v.Comment),
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
			Metadata:  parseMetadata(v.Comment),
		}
	case *http_actions.EarlyHint:
		rule = &models.HTTPRequestRule{
			Type:       "early-hint",
			HintName:   v.Name,
			HintFormat: v.Fmt,
			Cond:       v.Cond,
			CondTest:   v.CondTest,
			Metadata:   parseMetadata(v.Comment),
		}
	case *actions.Lua:
		rule = &models.HTTPRequestRule{
			Type:      "lua",
			LuaAction: v.Action,
			LuaParams: v.Params,
			Cond:      v.Cond,
			CondTest:  v.CondTest,
			Metadata:  parseMetadata(v.Comment),
		}
	case *http_actions.NormalizeURI:
		rule = &models.HTTPRequestRule{
			Type:             models.HTTPRequestRuleTypeNormalizeDashURI,
			Normalizer:       v.Normalizer,
			NormalizerFull:   v.Full,
			NormalizerStrict: v.Strict,
			Cond:             v.Cond,
			CondTest:         v.CondTest,
			Metadata:         parseMetadata(v.Comment),
		}
	case *http_actions.Pause:
		rule = &models.HTTPRequestRule{
			Type:     models.HTTPRequestRuleTypePause,
			Expr:     v.Pause.String(),
			Cond:     v.Cond,
			CondTest: v.CondTest,
			Metadata: parseMetadata(v.Comment),
		}
	case *http_actions.Redirect:
		var codePtr *int64
		if v.Code != "" {
			if code, err := strconv.ParseInt(v.Code, 10, 64); err == nil {
				codePtr = &code
			} else {
				return nil, err
			}
		}
		rule = &models.HTTPRequestRule{
			Type:        "redirect",
			RedirType:   v.Type,
			RedirValue:  v.Value,
			RedirOption: v.Option,
			Cond:        v.Cond,
			CondTest:    v.CondTest,
			RedirCode:   codePtr,
			Metadata:    parseMetadata(v.Comment),
		}
	case *actions.Reject:
		rule = &models.HTTPRequestRule{
			Type:     models.HTTPRequestRuleTypeReject,
			Cond:     v.Cond,
			CondTest: v.CondTest,
			Metadata: parseMetadata(v.Comment),
		}
	case *http_actions.ReplaceHeader:
		rule = &models.HTTPRequestRule{
			Type:      "replace-header",
			HdrName:   v.Name,
			HdrFormat: v.ReplaceFmt,
			HdrMatch:  v.MatchRegex,
			Cond:      v.Cond,
			CondTest:  v.CondTest,
			Metadata:  parseMetadata(v.Comment),
		}
	case *http_actions.ReplacePath:
		rule = &models.HTTPRequestRule{
			Type:      "replace-path",
			PathMatch: v.MatchRegex,
			PathFmt:   v.ReplaceFmt,
			Cond:      v.Cond,
			CondTest:  v.CondTest,
			Metadata:  parseMetadata(v.Comment),
		}
	case *http_actions.ReplacePathQ:
		rule = &models.HTTPRequestRule{
			Type:      "replace-pathq",
			PathMatch: v.MatchRegex,
			PathFmt:   v.ReplaceFmt,
			Cond:      v.Cond,
			CondTest:  v.CondTest,
			Metadata:  parseMetadata(v.Comment),
		}
	case *http_actions.ReplaceURI:
		rule = &models.HTTPRequestRule{
			Type:     "replace-uri",
			URIMatch: v.MatchRegex,
			URIFmt:   v.ReplaceFmt,
			Cond:     v.Cond,
			CondTest: v.CondTest,
			Metadata: parseMetadata(v.Comment),
		}
	case *http_actions.ReplaceValue:
		rule = &models.HTTPRequestRule{
			Type:      "replace-value",
			HdrName:   v.Name,
			HdrMatch:  v.MatchRegex,
			HdrFormat: v.ReplaceFmt,
			Cond:      v.Cond,
			CondTest:  v.CondTest,
			Metadata:  parseMetadata(v.Comment),
		}
	case *http_actions.Return:
		var returnContentTypePtr *string
		if v.ContentType != "" {
			returnContentTypePtr = &v.ContentType
		}
		rule = &models.HTTPRequestRule{
			Cond:                v.Cond,
			CondTest:            v.CondTest,
			ReturnHeaders:       actionHdr2ModelHdr(v.Hdrs),
			ReturnContent:       v.Content,
			ReturnContentFormat: v.ContentFormat,
			ReturnContentType:   returnContentTypePtr,
			ReturnStatusCode:    v.Status,
			Type:                "return",
			Metadata:            parseMetadata(v.Comment),
		}
	case *actions.ScAddGpc:
		if v.Int == nil && len(v.Expr.Expr) == 0 {
			return nil, NewConfError(ErrValidationError, "sc-add-gpc int or expr has to be set")
		}
		if v.Int != nil && len(v.Expr.Expr) > 0 {
			return nil, NewConfError(ErrValidationError, "sc-add-gpc int and expr are exclusive")
		}
		ID, _ := strconv.ParseInt(v.ID, 10, 64)
		Idx, _ := strconv.ParseInt(v.Idx, 10, 64)
		rule = &models.HTTPRequestRule{
			Type:     "sc-add-gpc",
			ScID:     ID,
			ScIdx:    Idx,
			ScExpr:   strings.Join(v.Expr.Expr, " "),
			ScInt:    v.Int,
			Cond:     v.Cond,
			CondTest: v.CondTest,
			Metadata: parseMetadata(v.Comment),
		}
	case *actions.ScIncGpc:
		ID, _ := strconv.ParseInt(v.ID, 10, 64)
		Idx, _ := strconv.ParseInt(v.Idx, 10, 64)
		rule = &models.HTTPRequestRule{
			Type:     "sc-inc-gpc",
			ScID:     ID,
			ScIdx:    Idx,
			Cond:     v.Cond,
			CondTest: v.CondTest,
			Metadata: parseMetadata(v.Comment),
		}
	case *actions.ScIncGpc0:
		ID, _ := strconv.ParseInt(v.ID, 10, 64)
		rule = &models.HTTPRequestRule{
			Type:     "sc-inc-gpc0",
			ScID:     ID,
			Cond:     v.Cond,
			CondTest: v.CondTest,
			Metadata: parseMetadata(v.Comment),
		}
	case *actions.ScIncGpc1:
		ID, _ := strconv.ParseInt(v.ID, 10, 64)
		rule = &models.HTTPRequestRule{
			Type:     "sc-inc-gpc1",
			ScID:     ID,
			Cond:     v.Cond,
			CondTest: v.CondTest,
			Metadata: parseMetadata(v.Comment),
		}
	case *actions.ScSetGpt:
		if v.Int == nil && len(v.Expr.Expr) == 0 {
			return nil, NewConfError(ErrValidationError, "sc-set-gpt: int or expr has to be set")
		}
		if v.Int != nil && len(v.Expr.Expr) > 0 {
			return nil, NewConfError(ErrValidationError, "sc-set-gpt: int and expr are exclusive")
		}
		scID, errp := strconv.ParseInt(v.ScID, 10, 64)
		if errp != nil {
			return nil, NewConfError(ErrValidationError, "sc-set-gpt: failed to parse sc-id an an int")
		}
		rule = &models.HTTPRequestRule{
			Type:     "sc-set-gpt",
			ScID:     scID,
			ScIdx:    v.Idx,
			ScExpr:   strings.Join(v.Expr.Expr, " "),
			ScInt:    v.Int,
			Cond:     v.Cond,
			CondTest: v.CondTest,
			Metadata: parseMetadata(v.Comment),
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
			Metadata: parseMetadata(v.Comment),
		}
	case *actions.SetDst:
		rule = &models.HTTPRequestRule{
			Type:     "set-dst",
			Expr:     v.Expr.String(),
			Cond:     v.Cond,
			CondTest: v.CondTest,
			Metadata: parseMetadata(v.Comment),
		}
	case *actions.SetDstPort:
		rule = &models.HTTPRequestRule{
			Type:     "set-dst-port",
			Expr:     v.Expr.String(),
			Cond:     v.Cond,
			CondTest: v.CondTest,
			Metadata: parseMetadata(v.Comment),
		}
	case *http_actions.SetHeader:
		rule = &models.HTTPRequestRule{
			Type:      "set-header",
			HdrName:   v.Name,
			HdrFormat: v.Fmt,
			Cond:      v.Cond,
			CondTest:  v.CondTest,
			Metadata:  parseMetadata(v.Comment),
		}
	case *actions.SetLogLevel:
		rule = &models.HTTPRequestRule{
			Type:     "set-log-level",
			LogLevel: v.Level,
			Cond:     v.Cond,
			CondTest: v.CondTest,
			Metadata: parseMetadata(v.Comment),
		}
	case *http_actions.SetMap:
		rule = &models.HTTPRequestRule{
			Type:        "set-map",
			MapFile:     v.FileName,
			MapKeyfmt:   v.KeyFmt,
			MapValuefmt: v.ValueFmt,
			Cond:        v.Cond,
			CondTest:    v.CondTest,
			Metadata:    parseMetadata(v.Comment),
		}
	case *actions.SetMark:
		rule = &models.HTTPRequestRule{
			Type:      "set-mark",
			MarkValue: v.Value,
			Cond:      v.Cond,
			CondTest:  v.CondTest,
			Metadata:  parseMetadata(v.Comment),
		}
	case *http_actions.SetMethod:
		rule = &models.HTTPRequestRule{
			Type:      "set-method",
			MethodFmt: v.Fmt,
			Cond:      v.Cond,
			CondTest:  v.CondTest,
			Metadata:  parseMetadata(v.Comment),
		}
	case *actions.SetNice:
		nice, _ := strconv.ParseInt(v.Value, 10, 64)
		rule = &models.HTTPRequestRule{
			Type:      "set-nice",
			NiceValue: nice,
			Cond:      v.Cond,
			CondTest:  v.CondTest,
			Metadata:  parseMetadata(v.Comment),
		}
	case *http_actions.SetPath:
		rule = &models.HTTPRequestRule{
			Type:     "set-path",
			PathFmt:  v.Fmt,
			Cond:     v.Cond,
			CondTest: v.CondTest,
			Metadata: parseMetadata(v.Comment),
		}
	case *http_actions.SetPathQ:
		rule = &models.HTTPRequestRule{
			Type:     "set-pathq",
			PathFmt:  v.Fmt,
			Cond:     v.Cond,
			CondTest: v.CondTest,
			Metadata: parseMetadata(v.Comment),
		}
	case *actions.SetPriorityClass:
		rule = &models.HTTPRequestRule{
			Type:     "set-priority-class",
			Expr:     strings.Join(v.Expr.Expr, " "),
			Cond:     v.Cond,
			CondTest: v.CondTest,
			Metadata: parseMetadata(v.Comment),
		}
	case *actions.SetPriorityOffset:
		rule = &models.HTTPRequestRule{
			Type:     "set-priority-offset",
			Expr:     strings.Join(v.Expr.Expr, " "),
			Cond:     v.Cond,
			CondTest: v.CondTest,
			Metadata: parseMetadata(v.Comment),
		}
	case *http_actions.SetQuery:
		rule = &models.HTTPRequestRule{
			Type:      "set-query",
			HdrFormat: v.Fmt,
			Cond:      v.Cond,
			CondTest:  v.CondTest,
			Metadata:  parseMetadata(v.Comment),
		}
	case *http_actions.SetSrc:
		rule = &models.HTTPRequestRule{
			Type:     "set-src",
			Expr:     v.Expr.String(),
			Cond:     v.Cond,
			CondTest: v.CondTest,
			Metadata: parseMetadata(v.Comment),
		}
	case *actions.SetSrcPort:
		rule = &models.HTTPRequestRule{
			Type:     "set-src-port",
			Expr:     v.Expr.String(),
			Cond:     v.Cond,
			CondTest: v.CondTest,
			Metadata: parseMetadata(v.Comment),
		}
	case *http_actions.SetTimeout:
		rule = &models.HTTPRequestRule{
			Type:        models.HTTPRequestRuleTypeSetDashTimeout,
			Timeout:     v.Timeout,
			TimeoutType: v.Type,
			Cond:        v.Cond,
			CondTest:    v.CondTest,
			Metadata:    parseMetadata(v.Comment),
		}
	case *actions.SetTos:
		rule = &models.HTTPRequestRule{
			Type:     "set-tos",
			TosValue: v.Value,
			Cond:     v.Cond,
			CondTest: v.CondTest,
			Metadata: parseMetadata(v.Comment),
		}
	case *http_actions.SetURI:
		rule = &models.HTTPRequestRule{
			Type:     "set-uri",
			URIFmt:   v.Fmt,
			Cond:     v.Cond,
			CondTest: v.CondTest,
			Metadata: parseMetadata(v.Comment),
		}
	case *actions.SetVar:
		rule = &models.HTTPRequestRule{
			Type:     "set-var",
			VarName:  v.VarName,
			VarExpr:  strings.Join(v.Expr.Expr, " "),
			VarScope: v.VarScope,
			Cond:     v.Cond,
			CondTest: v.CondTest,
			Metadata: parseMetadata(v.Comment),
		}
	case *actions.SetVarFmt:
		rule = &models.HTTPRequestRule{
			Type:      "set-var-fmt",
			VarName:   v.VarName,
			VarFormat: strings.Join(v.Fmt.Expr, " "),
			VarScope:  v.VarScope,
			Cond:      v.Cond,
			CondTest:  v.CondTest,
			Metadata:  parseMetadata(v.Comment),
		}
	case *actions.SendSpoeGroup:
		rule = &models.HTTPRequestRule{
			Type:       "send-spoe-group",
			SpoeEngine: v.Engine,
			SpoeGroup:  v.Group,
			Cond:       v.Cond,
			CondTest:   v.CondTest,
			Metadata:   parseMetadata(v.Comment),
		}
	case *actions.SilentDrop:
		rule = &models.HTTPRequestRule{
			RstTTL:   v.RstTTL,
			Type:     "silent-drop",
			Cond:     v.Cond,
			CondTest: v.CondTest,
			Metadata: parseMetadata(v.Comment),
		}
	case *http_actions.StrictMode:
		rule = &models.HTTPRequestRule{
			Type:       "strict-mode",
			StrictMode: v.Mode,
			Cond:       v.Cond,
			CondTest:   v.CondTest,
			Metadata:   parseMetadata(v.Comment),
		}
	case *http_actions.Tarpit:
		rule = &models.HTTPRequestRule{
			Type:       "tarpit",
			Cond:       v.Cond,
			CondTest:   v.CondTest,
			DenyStatus: v.Status,
			Metadata:   parseMetadata(v.Comment),
		}
	case *actions.TrackSc:
		rule = &models.HTTPRequestRule{
			Type:                string(actions.TrackScType),
			TrackScKey:          v.Key,
			TrackScTable:        v.Table,
			Cond:                v.Cond,
			CondTest:            v.CondTest,
			TrackScStickCounter: &v.StickCounter,
			Metadata:            parseMetadata(v.Comment),
		}
	case *actions.UnsetVar:
		rule = &models.HTTPRequestRule{
			Type:     "unset-var",
			VarName:  v.Name,
			VarScope: v.Scope,
			Cond:     v.Cond,
			CondTest: v.CondTest,
			Metadata: parseMetadata(v.Comment),
		}
	case *actions.UseService:
		rule = &models.HTTPRequestRule{
			Type:        "use-service",
			ServiceName: v.Name,
			Cond:        v.Cond,
			CondTest:    v.CondTest,
			Metadata:    parseMetadata(v.Comment),
		}
	case *http_actions.WaitForBody:
		rule = &models.HTTPRequestRule{
			Type:        "wait-for-body",
			WaitTime:    misc.ParseTimeout(v.Time),
			WaitAtLeast: misc.ParseSize(v.AtLeast),
			Cond:        v.Cond,
			CondTest:    v.CondTest,
			Metadata:    parseMetadata(v.Comment),
		}
	case *http_actions.WaitForHandshake:
		rule = &models.HTTPRequestRule{
			Type:     "wait-for-handshake",
			Cond:     v.Cond,
			CondTest: v.CondTest,
			Metadata: parseMetadata(v.Comment),
		}
	case *actions.SetBandwidthLimit:
		rule = &models.HTTPRequestRule{
			Type:                 "set-bandwidth-limit",
			BandwidthLimitName:   v.Name,
			BandwidthLimitLimit:  v.Limit.String(),
			BandwidthLimitPeriod: v.Period.String(),
			Cond:                 v.Cond,
			CondTest:             v.CondTest,
			Metadata:             parseMetadata(v.Comment),
		}
	case *actions.SetBcMark:
		rule = &models.HTTPRequestRule{
			Type:     models.HTTPRequestRuleTypeSetDashBcDashMark,
			Expr:     v.Expr.String(),
			Cond:     v.Cond,
			CondTest: v.CondTest,
			Metadata: parseMetadata(v.Comment),
		}
	case *actions.SetBcTos:
		rule = &models.HTTPRequestRule{
			Type:     models.HTTPRequestRuleTypeSetDashBcDashTos,
			Expr:     v.Expr.String(),
			Cond:     v.Cond,
			CondTest: v.CondTest,
			Metadata: parseMetadata(v.Comment),
		}
	case *actions.SetFcMark:
		rule = &models.HTTPRequestRule{
			Type:     models.HTTPRequestRuleTypeSetDashFcDashMark,
			Expr:     v.Expr.String(),
			Cond:     v.Cond,
			CondTest: v.CondTest,
			Metadata: parseMetadata(v.Comment),
		}
	case *actions.SetFcTos:
		rule = &models.HTTPRequestRule{
			Type:     models.HTTPRequestRuleTypeSetDashFcDashTos,
			Expr:     v.Expr.String(),
			Cond:     v.Cond,
			CondTest: v.CondTest,
			Metadata: parseMetadata(v.Comment),
		}
	case *actions.SetRetries:
		rule = &models.HTTPRequestRule{
			Type:     models.HTTPRequestRuleTypeSetDashRetries,
			Expr:     v.Expr.String(),
			Cond:     v.Cond,
			CondTest: v.CondTest,
			Metadata: parseMetadata(v.Comment),
		}
	case *actions.DoLog:
		rule = &models.HTTPRequestRule{
			Type:     models.HTTPRequestRuleTypeDoDashLog,
			Cond:     v.Cond,
			CondTest: v.CondTest,
			Metadata: parseMetadata(v.Comment),
		}
	}

	return rule, nil
}

func SerializeHTTPRequestRule(f models.HTTPRequestRule, opt *options.ConfigurationOptions) (types.Action, error) { //nolint:gocyclo,gocognit,ireturn,cyclop,maintidx
	comment, err := serializeMetadata(f.Metadata)
	if err != nil {
		return nil, err
	}
	var rule types.Action
	switch f.Type {
	case "add-acl":
		rule = &http_actions.AddACL{
			FileName: f.ACLFile,
			KeyFmt:   f.ACLKeyfmt,
			Cond:     f.Cond,
			CondTest: f.CondTest,
			Comment:  comment,
		}
	case "add-header":
		rule = &http_actions.AddHeader{
			Name:     f.HdrName,
			Fmt:      f.HdrFormat,
			Cond:     f.Cond,
			CondTest: f.CondTest,
			Comment:  comment,
		}
	case "allow":
		rule = &http_actions.Allow{
			Cond:     f.Cond,
			CondTest: f.CondTest,
			Comment:  comment,
		}
	case "auth":
		rule = &http_actions.Auth{
			Realm:    f.AuthRealm,
			Cond:     f.Cond,
			CondTest: f.CondTest,
			Comment:  comment,
		}
	case "cache-use":
		rule = &http_actions.CacheUse{
			Name:     f.CacheName,
			Cond:     f.Cond,
			CondTest: f.CondTest,
			Comment:  comment,
		}
	case "capture":
		if f.CaptureLen > 0 && f.CaptureID != nil {
			return nil, NewConfError(ErrValidationError, "capture len and id are exclusive")
		}
		if f.CaptureLen == 0 && f.CaptureID == nil {
			return nil, NewConfError(ErrValidationError, "capture len has to be greater than 0 or capture_id has to be set")
		}
		r := &http_actions.Capture{
			Sample:   f.CaptureSample,
			Cond:     f.Cond,
			CondTest: f.CondTest,
			Comment:  comment,
		}
		if f.CaptureLen > 0 {
			r.Len = &f.CaptureLen
		} else {
			r.SlotID = f.CaptureID
		}
		rule = r
	case "del-acl":
		rule = &http_actions.DelACL{
			FileName: f.ACLFile,
			KeyFmt:   f.ACLKeyfmt,
			Cond:     f.Cond,
			CondTest: f.CondTest,
			Comment:  comment,
		}
	case "del-header":
		rule = &http_actions.DelHeader{
			Name:     f.HdrName,
			Method:   f.HdrMethod,
			Cond:     f.Cond,
			CondTest: f.CondTest,
			Comment:  comment,
		}
	case "del-map":
		rule = &http_actions.DelMap{
			FileName: f.MapFile,
			KeyFmt:   f.MapKeyfmt,
			Cond:     f.Cond,
			CondTest: f.CondTest,
			Comment:  comment,
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
			Comment:       comment,
		}
		if !http_actions.IsPayload(f.ReturnContentFormat) && f.DenyStatus != nil {
			if ok := http_actions.AllowedErrorCode(*f.DenyStatus); !ok {
				return rule, NewConfError(ErrValidationError, "invalid Status Code for error type response")
			}
		}
	case "disable-l7-retry":
		rule = &http_actions.DisableL7Retry{
			Cond:     f.Cond,
			CondTest: f.CondTest,
			Comment:  comment,
		}
	case "do-resolve":
		rule = &actions.DoResolve{
			Var:       f.VarName,
			Resolvers: f.Resolvers,
			Protocol:  f.Protocol,
			Expr:      common.Expression{Expr: strings.Split(f.Expr, " ")},
			Cond:      f.Cond,
			CondTest:  f.CondTest,
			Comment:   comment,
		}
	case "early-hint":
		rule = &http_actions.EarlyHint{
			Name:     f.HintName,
			Fmt:      f.HintFormat,
			Cond:     f.Cond,
			CondTest: f.CondTest,
			Comment:  comment,
		}
	case "lua":
		rule = &actions.Lua{
			Action:   f.LuaAction,
			Params:   f.LuaParams,
			Cond:     f.Cond,
			CondTest: f.CondTest,
			Comment:  comment,
		}
	case "normalize-uri":
		rule = &http_actions.NormalizeURI{
			Normalizer: f.Normalizer,
			Strict:     f.NormalizerStrict,
			Full:       f.NormalizerFull,
			Cond:       f.Cond,
			CondTest:   f.CondTest,
			Comment:    comment,
		}
	case models.HTTPRequestRuleTypePause:
		rule = &http_actions.Pause{
			Pause:    common.Expression{Expr: strings.Split(f.Expr, " ")},
			Cond:     f.Cond,
			CondTest: f.CondTest,
			Comment:  comment,
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
			Comment:  comment,
		}
	case "reject":
		rule = &actions.Reject{
			Cond:     f.Cond,
			CondTest: f.CondTest,
			Comment:  comment,
		}
	case "replace-header":
		rule = &http_actions.ReplaceHeader{
			Name:       f.HdrName,
			ReplaceFmt: f.HdrFormat,
			MatchRegex: f.HdrMatch,
			Cond:       f.Cond,
			CondTest:   f.CondTest,
			Comment:    comment,
		}
	case "replace-path":
		rule = &http_actions.ReplacePath{
			MatchRegex: f.PathMatch,
			ReplaceFmt: f.PathFmt,
			Cond:       f.Cond,
			CondTest:   f.CondTest,
			Comment:    comment,
		}
	case "replace-pathq":
		rule = &http_actions.ReplacePathQ{
			MatchRegex: f.PathMatch,
			ReplaceFmt: f.PathFmt,
			Cond:       f.Cond,
			CondTest:   f.CondTest,
			Comment:    comment,
		}
	case "replace-uri":
		rule = &http_actions.ReplaceURI{
			ReplaceFmt: f.URIFmt,
			MatchRegex: f.URIMatch,
			Cond:       f.Cond,
			CondTest:   f.CondTest,
			Comment:    comment,
		}
	case "replace-value":
		rule = &http_actions.ReplaceValue{
			Name:       f.HdrName,
			ReplaceFmt: f.HdrFormat,
			MatchRegex: f.HdrMatch,
			Cond:       f.Cond,
			CondTest:   f.CondTest,
			Comment:    comment,
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
			Comment:       comment,
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
			Comment:  comment,
		}
	case "sc-inc-gpc":
		rule = &actions.ScAddGpc{
			ID:       strconv.FormatInt(f.ScID, 10),
			Idx:      strconv.FormatInt(f.ScIdx, 10),
			Cond:     f.Cond,
			CondTest: f.CondTest,
			Comment:  comment,
		}
	case "sc-inc-gpc0":
		rule = &actions.ScIncGpc0{
			ID:       strconv.FormatInt(f.ScID, 10),
			Cond:     f.Cond,
			CondTest: f.CondTest,
			Comment:  comment,
		}
	case "sc-inc-gpc1":
		rule = &actions.ScIncGpc1{
			ID:       strconv.FormatInt(f.ScID, 10),
			Cond:     f.Cond,
			CondTest: f.CondTest,
			Comment:  comment,
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
			Comment:  comment,
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
			Comment:  comment,
		}
	case "set-dst":
		rule = &actions.SetDst{
			Expr:     common.Expression{Expr: strings.Split(f.Expr, " ")},
			Cond:     f.Cond,
			CondTest: f.CondTest,
			Comment:  comment,
		}
	case "set-dst-port":
		rule = &actions.SetDstPort{
			Expr:     common.Expression{Expr: strings.Split(f.Expr, " ")},
			Cond:     f.Cond,
			CondTest: f.CondTest,
			Comment:  comment,
		}
	case "set-header":
		rule = &http_actions.SetHeader{
			Name:     f.HdrName,
			Fmt:      f.HdrFormat,
			Cond:     f.Cond,
			CondTest: f.CondTest,
			Comment:  comment,
		}
	case "set-log-level":
		rule = &actions.SetLogLevel{
			Level:    f.LogLevel,
			Cond:     f.Cond,
			CondTest: f.CondTest,
			Comment:  comment,
		}
	case "set-map":
		rule = &http_actions.SetMap{
			FileName: f.MapFile,
			KeyFmt:   f.MapKeyfmt,
			ValueFmt: f.MapValuefmt,
			Cond:     f.Cond,
			CondTest: f.CondTest,
			Comment:  comment,
		}
	case "set-mark":
		rule = &actions.SetMark{
			Value:    f.MarkValue,
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "set-method":
		rule = &http_actions.SetMethod{
			Fmt:      f.MethodFmt,
			Cond:     f.Cond,
			CondTest: f.CondTest,
			Comment:  comment,
		}
	case "set-nice":
		rule = &actions.SetNice{
			Value:    strconv.FormatInt(f.NiceValue, 10),
			Cond:     f.Cond,
			CondTest: f.CondTest,
			Comment:  comment,
		}
	case "set-path":
		rule = &http_actions.SetPath{
			Fmt:      f.PathFmt,
			Cond:     f.Cond,
			CondTest: f.CondTest,
			Comment:  comment,
		}
	case "set-pathq":
		rule = &http_actions.SetPathQ{
			Fmt:      f.PathFmt,
			Cond:     f.Cond,
			CondTest: f.CondTest,
			Comment:  comment,
		}
	case "set-priority-class":
		rule = &actions.SetPriorityClass{
			Expr:     common.Expression{Expr: strings.Split(f.Expr, " ")},
			Cond:     f.Cond,
			CondTest: f.CondTest,
			Comment:  comment,
		}
	case "set-priority-offset":
		rule = &actions.SetPriorityOffset{
			Expr:     common.Expression{Expr: strings.Split(f.Expr, " ")},
			Cond:     f.Cond,
			CondTest: f.CondTest,
			Comment:  comment,
		}
	case "set-query":
		rule = &http_actions.SetQuery{
			Fmt:      f.HdrFormat,
			Cond:     f.Cond,
			CondTest: f.CondTest,
			Comment:  comment,
		}
	case "set-src":
		rule = &http_actions.SetSrc{
			Expr:     common.Expression{Expr: strings.Split(f.Expr, " ")},
			Cond:     f.Cond,
			CondTest: f.CondTest,
			Comment:  comment,
		}
	case "set-src-port":
		rule = &actions.SetSrcPort{
			Expr:     common.Expression{Expr: strings.Split(f.Expr, " ")},
			Cond:     f.Cond,
			CondTest: f.CondTest,
			Comment:  comment,
		}
	case "set-timeout":
		rule = &http_actions.SetTimeout{
			Timeout:  f.Timeout,
			Type:     f.TimeoutType,
			Cond:     f.Cond,
			CondTest: f.CondTest,
			Comment:  comment,
		}
	case "set-tos":
		rule = &actions.SetTos{
			Value:    f.TosValue,
			Cond:     f.Cond,
			CondTest: f.CondTest,
			Comment:  comment,
		}
	case "set-uri":
		rule = &http_actions.SetURI{
			Fmt:      f.URIFmt,
			Cond:     f.Cond,
			CondTest: f.CondTest,
			Comment:  comment,
		}
	case "set-var":
		rule = &actions.SetVar{
			Expr:     common.Expression{Expr: strings.Split(f.VarExpr, " ")},
			VarName:  f.VarName,
			VarScope: f.VarScope,
			Cond:     f.Cond,
			CondTest: f.CondTest,
			Comment:  comment,
		}
	case "set-var-fmt":
		rule = &actions.SetVarFmt{
			Fmt:      common.Expression{Expr: strings.Split(f.VarFormat, " ")},
			VarName:  f.VarName,
			VarScope: f.VarScope,
			Cond:     f.Cond,
			CondTest: f.CondTest,
			Comment:  comment,
		}
	case "send-spoe-group":
		rule = &actions.SendSpoeGroup{
			Engine:   f.SpoeEngine,
			Group:    f.SpoeGroup,
			Cond:     f.Cond,
			CondTest: f.CondTest,
			Comment:  comment,
		}
	case "silent-drop":
		rule = &actions.SilentDrop{
			RstTTL:   f.RstTTL,
			Cond:     f.Cond,
			CondTest: f.CondTest,
			Comment:  comment,
		}
	case "strict-mode":
		rule = &http_actions.StrictMode{
			Mode:     f.StrictMode,
			Cond:     f.Cond,
			CondTest: f.CondTest,
			Comment:  comment,
		}
	case "tarpit":
		rule = &http_actions.Tarpit{
			Status:   f.DenyStatus,
			Cond:     f.Cond,
			CondTest: f.CondTest,
			Comment:  comment,
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
			Comment:      comment,
		}
	case "unset-var":
		rule = &actions.UnsetVar{
			Name:     f.VarName,
			Scope:    f.VarScope,
			Cond:     f.Cond,
			CondTest: f.CondTest,
			Comment:  comment,
		}
	case "use-service":
		rule = &actions.UseService{
			Name:     f.ServiceName,
			Cond:     f.Cond,
			CondTest: f.CondTest,
			Comment:  comment,
		}
	case "wait-for-body":
		r := &http_actions.WaitForBody{
			Cond:     f.Cond,
			CondTest: f.CondTest,
			Comment:  comment,
		}
		if f.WaitTime != nil {
			r.Time = misc.SerializeTime(*f.WaitTime, opt.PreferredTimeSuffix)
		}
		if f.WaitAtLeast != nil {
			r.AtLeast = misc.SerializeSize(*f.WaitAtLeast)
		}
		rule = r
	case "wait-for-handshake":
		rule = &http_actions.WaitForHandshake{
			Cond:     f.Cond,
			CondTest: f.CondTest,
			Comment:  comment,
		}
	case "set-bandwidth-limit":
		rule = &actions.SetBandwidthLimit{
			Name:     f.BandwidthLimitName,
			Limit:    common.Expression{Expr: strings.Split(f.BandwidthLimitLimit, " ")},
			Period:   common.Expression{Expr: strings.Split(f.BandwidthLimitPeriod, " ")},
			Cond:     f.Cond,
			CondTest: f.CondTest,
			Comment:  comment,
		}
	case "set-bc-mark":
		rule = &actions.SetBcMark{
			Expr:     common.Expression{Expr: strings.Split(f.Expr+f.MarkValue, " ")},
			Cond:     f.Cond,
			CondTest: f.CondTest,
			Comment:  comment,
		}
	case "set-bc-tos":
		rule = &actions.SetBcTos{
			Expr:     common.Expression{Expr: strings.Split(f.Expr+f.TosValue, " ")},
			Cond:     f.Cond,
			CondTest: f.CondTest,
			Comment:  comment,
		}
	case "set-fc-mark":
		rule = &actions.SetFcMark{
			Expr:     common.Expression{Expr: strings.Split(f.Expr+f.MarkValue, " ")},
			Cond:     f.Cond,
			CondTest: f.CondTest,
			Comment:  comment,
		}
	case "set-fc-tos":
		rule = &actions.SetFcTos{
			Expr:     common.Expression{Expr: strings.Split(f.Expr+f.TosValue, " ")},
			Cond:     f.Cond,
			CondTest: f.CondTest,
			Comment:  comment,
		}
	case "set-retries":
		rule = &actions.SetRetries{
			Expr:     common.Expression{Expr: strings.Split(f.Expr, " ")},
			Cond:     f.Cond,
			CondTest: f.CondTest,
			Comment:  comment,
		}
	case "do-log":
		rule = &actions.DoLog{
			Cond:     f.Cond,
			CondTest: f.CondTest,
			Comment:  comment,
		}
	}
	return rule, nil
}
