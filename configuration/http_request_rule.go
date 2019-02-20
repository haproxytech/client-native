package configuration

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/haproxytech/config-parser/common"

	"github.com/haproxytech/config-parser/parsers/http/actions"

	strfmt "github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/config-parser"
	parser_errors "github.com/haproxytech/config-parser/errors"
	"github.com/haproxytech/config-parser/types"
	"github.com/haproxytech/models"
)

// GetHTTPRequestRules returns a struct with configuration version and an array of
// configured http request rules in the specified parent. Returns error on fail.
func (c *Client) GetHTTPRequestRules(parentType, parentName string, transactionID string) (*models.GetHTTPRequestRulesOKBody, error) {
	if c.Cache.Enabled() {
		httpRules, found := c.Cache.HttpRequestRules.Get(parentName, parentType, transactionID)
		if found {
			return &models.GetHTTPRequestRulesOKBody{Version: c.Cache.Version.Get(transactionID), Data: httpRules}, nil
		}
	}
	if err := c.ConfigParser.LoadData(c.getTransactionFile(transactionID)); err != nil {
		return nil, err
	}

	httpRules, err := c.parseHTTPRequestRules(parentType, parentName)
	if err != nil {
		if err == parser_errors.SectionMissingErr {
			return nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("%s %s does not exist", parentType, parentName))
		}
		return nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return nil, err
	}

	if c.Cache.Enabled() {
		c.Cache.HttpRequestRules.SetAll(parentName, parentType, transactionID, httpRules)
	}
	return &models.GetHTTPRequestRulesOKBody{Version: v, Data: httpRules}, nil
}

// GetHTTPRequestRule returns a struct with configuration version and a requested http request rule
// in the specified parent. Returns error on fail or if http request rule does not exist.
func (c *Client) GetHTTPRequestRule(id int64, parentType, parentName string, transactionID string) (*models.GetHTTPRequestRuleOKBody, error) {
	if c.Cache.Enabled() {
		httpRule, found := c.Cache.HttpRequestRules.GetOne(id, parentName, parentType, transactionID)
		if found {
			return &models.GetHTTPRequestRuleOKBody{Version: c.Cache.Version.Get(transactionID), Data: httpRule}, nil
		}
	}
	if err := c.ConfigParser.LoadData(c.getTransactionFile(transactionID)); err != nil {
		return nil, err
	}

	var p parser.Section
	if parentType == "backend" {
		p = parser.Backends
	} else if parentType == "frontend" {
		p = parser.Frontends
	}

	data, err := c.ConfigParser.GetOne(p, parentName, "http-request", int(id))
	if err != nil {
		if err == parser_errors.SectionMissingErr {
			return nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("%s %s does not exist", parentType, parentName))
		}
		if err == parser_errors.FetchError {
			return nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("HTTP Request Rule %v does not exist in %s %s", id, parentType, parentName))
		}
		return nil, err
	}

	httpRule := parseHTTPRequestRule(data.(types.HTTPAction))
	httpRule.ID = &id

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return nil, err
	}

	if c.Cache.Enabled() {
		c.Cache.HttpRequestRules.Set(id, parentName, parentType, transactionID, httpRule)
	}

	return &models.GetHTTPRequestRuleOKBody{Version: v, Data: httpRule}, nil
}

// DeleteHTTPRequestRule deletes a http request rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) DeleteHTTPRequestRule(id int64, parentType string, parentName string, transactionID string, version int64) error {
	t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	var p parser.Section
	if parentType == "backend" {
		p = parser.Backends
	} else if parentType == "frontend" {
		p = parser.Frontends
	}

	if err := c.ConfigParser.Delete(p, parentName, "http-request", int(id)); err != nil {
		if err == parser_errors.SectionMissingErr {
			return NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("%s %s does not exist", parentName, parentType))
		}
		if err == parser_errors.FetchError {
			return NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("HTTP Request Rule %v does not exist in %s %s", id, parentName, parentType))
		}
		return err
	}

	if err := c.saveData(t, transactionID); err != nil {
		return err
	}

	if c.Cache.Enabled() {
		c.Cache.HttpRequestRules.InvalidateParent(transactionID, parentName, parentType)
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

	t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	var p parser.Section
	if parentType == "backend" {
		p = parser.Backends
	} else if parentType == "frontend" {
		p = parser.Frontends
	}

	if err := c.ConfigParser.Insert(p, parentName, "http-request", serializeHTTPRequestRule(*data), int(*data.ID)); err != nil {
		if err == parser_errors.IndexOutOfRange {
			return c.errAndDeleteTransaction(NewConfError(ErrObjectIndexOutOfRange,
				fmt.Sprintf("HTTP Request Rule with id %v in %s %s out of range", int(*data.ID), parentType, parentName)), t, transactionID == "")
		}
		if err == parser_errors.SectionMissingErr {
			return NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("%s %s does not exist", parentType, parentName))
		}
		return c.errAndDeleteTransaction(err, t, transactionID == "")
	}

	if err := c.saveData(t, transactionID); err != nil {
		return err
	}
	if c.Cache.Enabled() {
		c.Cache.HttpRequestRules.InvalidateParent(transactionID, parentName, parentType)
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
	t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	var p parser.Section
	if parentType == "backend" {
		p = parser.Backends
	} else if parentType == "frontend" {
		p = parser.Frontends
	}

	if _, err := c.ConfigParser.GetOne(p, parentName, "http-request", int(id)); err != nil {
		if err == parser_errors.SectionMissingErr {
			return NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("%s %s does not exist", parentType, parentName))
		}
		if err == parser_errors.FetchError {
			return NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("HTTP Request Rule %v does not exist in %s %s", id, parentType, parentName))
		}
		return err
	}

	if err := c.ConfigParser.Set(p, parentName, "http-request", serializeHTTPRequestRule(*data), int(id)); err != nil {
		if err == parser_errors.IndexOutOfRange {
			return c.errAndDeleteTransaction(NewConfError(ErrObjectIndexOutOfRange,
				fmt.Sprintf("HTTP Request Rule with id %v in %s %s out of range", int(*data.ID), parentType, parentName)), t, transactionID == "")
		}
		return c.errAndDeleteTransaction(err, t, transactionID == "")
	}

	if err := c.saveData(t, transactionID); err != nil {
		return err
	}
	if c.Cache.Enabled() {
		c.Cache.HttpRequestRules.InvalidateParent(transactionID, parentName, parentType)
	}
	return nil
}

func (c *Client) parseHTTPRequestRules(t, pName string) (models.HTTPRequestRules, error) {
	p := parser.Global
	if t == "frontend" {
		p = parser.Frontends
	} else if t == "backend" {
		p = parser.Backends
	}

	httpReqRules := models.HTTPRequestRules{}
	data, err := c.ConfigParser.Get(p, pName, "http-request", false)
	if err != nil {
		if err == parser_errors.FetchError {
			return httpReqRules, nil
		}
		return nil, err
	}

	rules := data.([]types.HTTPAction)
	for i, r := range rules {
		id := int64(i)
		httpReqRule := parseHTTPRequestRule(r)
		if httpReqRule != nil {
			httpReqRule.ID = &id
			httpReqRules = append(httpReqRules, httpReqRule)
		}
	}
	return httpReqRules, nil
}

func parseHTTPRequestRule(f types.HTTPAction) *models.HTTPRequestRule {
	switch v := f.(type) {
	case *actions.Allow:
		return &models.HTTPRequestRule{
			Type:     "allow",
			Cond:     v.Cond,
			CondTest: v.CondTest,
		}
	case *actions.Deny:
		return &models.HTTPRequestRule{
			Type:     "deny",
			Cond:     v.Cond,
			CondTest: v.CondTest,
		}
	case *actions.Auth:
		return &models.HTTPRequestRule{
			Type:      "auth",
			AuthRealm: v.Realm,
			Cond:      v.Cond,
			CondTest:  v.CondTest,
		}
	case *actions.Redirect:
		code, _ := strconv.ParseInt(v.Code, 10, 64)
		r := &models.HTTPRequestRule{
			Type:        "redirect",
			RedirType:   v.Type,
			RedirValue:  v.Value,
			RedirOption: v.Option,
			Cond:        v.Cond,
			CondTest:    v.CondTest,
		}
		if code != 0 {
			r.RedirCode = code
		}
		return r
	case *actions.Tarpit:
		s, _ := strconv.ParseInt(v.DenyStatus, 10, 64)
		r := &models.HTTPRequestRule{
			Type:     "tarpit",
			Cond:     v.Cond,
			CondTest: v.CondTest,
		}
		if s != 0 {
			r.DenyStatus = s
		}
		return r
	case *actions.AddHeader:
		return &models.HTTPRequestRule{
			Type:      "add-header",
			HdrName:   v.Name,
			HdrFormat: v.Fmt,
			Cond:      v.Cond,
			CondTest:  v.CondTest,
		}
	case *actions.SetHeader:
		return &models.HTTPRequestRule{
			Type:      "set-header",
			HdrName:   v.Name,
			HdrFormat: v.Fmt,
			Cond:      v.Cond,
			CondTest:  v.CondTest,
		}
	case *actions.DelHeader:
		return &models.HTTPRequestRule{
			Type:     "del-header",
			HdrName:  v.Name,
			Cond:     v.Cond,
			CondTest: v.CondTest,
		}
	case *actions.ReplaceHeader:
		return &models.HTTPRequestRule{
			Type:      "replace-header",
			HdrName:   v.Name,
			HdrFormat: v.ReplaceFmt,
			HdrMatch:  v.MatchRegex,
			Cond:      v.Cond,
			CondTest:  v.CondTest,
		}
	case *actions.SetLogLevel:
		return &models.HTTPRequestRule{
			Type:     "set-log-level",
			LogLevel: v.Level,
			Cond:     v.Cond,
			CondTest: v.CondTest,
		}
	case *actions.SetVar:
		return &models.HTTPRequestRule{
			Type:     "set-var",
			VarName:  v.VarName,
			VarExpr:  strings.Join(v.Expr.Expr, " "),
			VarScope: v.VarScope,
			Cond:     v.Cond,
			CondTest: v.CondTest,
		}
	case *actions.ReplaceValue:
		return &models.HTTPRequestRule{
			Type:      "replace-value",
			HdrName:   v.Name,
			HdrMatch:  v.MatchRegex,
			HdrFormat: v.ReplaceFmt,
			Cond:      v.Cond,
			CondTest:  v.CondTest,
		}
	case *actions.AddAcl:
		return &models.HTTPRequestRule{
			Type:      "add-acl",
			ACLFile:   v.FileName,
			ACLKeyfmt: v.KeyFmt,
			Cond:      v.Cond,
			CondTest:  v.CondTest,
		}
	case *actions.DelAcl:
		return &models.HTTPRequestRule{
			Type:      "del-acl",
			ACLFile:   v.FileName,
			ACLKeyfmt: v.KeyFmt,
			Cond:      v.Cond,
			CondTest:  v.CondTest,
		}
	case *actions.SendSpoeGroup:
		return &models.HTTPRequestRule{
			Type:       "send-spoe-group",
			SpoeEngine: v.Engine,
			SpoeGroup:  v.Group,
			Cond:       v.Cond,
			CondTest:   v.CondTest,
		}
	}
	return nil
}

func serializeHTTPRequestRule(f models.HTTPRequestRule) types.HTTPAction {
	switch f.Type {
	case "allow":
		return &actions.Allow{
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "deny":
		return &actions.Deny{
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "auth":
		return &actions.Auth{
			Realm:    f.AuthRealm,
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "redirect":
		return &actions.Redirect{
			Type:     f.RedirType,
			Value:    f.RedirValue,
			Code:     strconv.FormatInt(f.RedirCode, 10),
			Option:   f.RedirOption,
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "tarpit":
		return &actions.Tarpit{
			DenyStatus: strconv.FormatInt(f.DenyStatus, 10),
			Cond:       f.Cond,
			CondTest:   f.CondTest,
		}
	case "add-header":
		return &actions.AddHeader{
			Name:     f.HdrName,
			Fmt:      f.HdrFormat,
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "set-header":
		return &actions.SetHeader{
			Name:     f.HdrName,
			Fmt:      f.HdrFormat,
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "del-header":
		return &actions.DelHeader{
			Name:     f.HdrName,
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "replace-header":
		return &actions.ReplaceHeader{
			Name:       f.HdrName,
			ReplaceFmt: f.HdrFormat,
			MatchRegex: f.HdrMatch,
			Cond:       f.Cond,
			CondTest:   f.CondTest,
		}
	case "replace-value":
		return &actions.ReplaceValue{
			Name:       f.HdrName,
			ReplaceFmt: f.HdrFormat,
			MatchRegex: f.HdrMatch,
			Cond:       f.Cond,
			CondTest:   f.CondTest,
		}
	case "set-log-level":
		return &actions.SetLogLevel{
			Level:    f.LogLevel,
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "set-var":
		return &actions.SetVar{
			Expr:     common.Expression{Expr: strings.Split(f.VarExpr, " ")},
			VarName:  f.VarName,
			VarScope: f.VarScope,
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "add-acl":
		return &actions.AddAcl{
			FileName: f.ACLFile,
			KeyFmt:   f.ACLKeyfmt,
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "del-acl":
		return &actions.DelAcl{
			FileName: f.ACLFile,
			KeyFmt:   f.ACLKeyfmt,
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "send-spoe-group":
		return &actions.SendSpoeGroup{
			Engine:   f.SpoeEngine,
			Group:    f.SpoeGroup,
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	}
	return nil
}
