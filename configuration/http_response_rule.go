package configuration

import (
	"fmt"
	"strconv"
	"strings"

	strfmt "github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/config-parser"
	"github.com/haproxytech/config-parser/common"
	parser_errors "github.com/haproxytech/config-parser/errors"
	"github.com/haproxytech/config-parser/parsers/http/actions"
	"github.com/haproxytech/config-parser/types"
	"github.com/haproxytech/models"
)

// GetHTTPResponseRules returns a struct with configuration version and an array of
// configured http response rules in the specified parent. Returns error on fail.
func (c *Client) GetHTTPResponseRules(parentType, parentName string, transactionID string) (*models.GetHTTPResponseRulesOKBody, error) {
	if c.Cache.Enabled() {
		httpRules, found := c.Cache.HttpResponseRules.Get(parentName, parentType, transactionID)
		if found {
			return &models.GetHTTPResponseRulesOKBody{Version: c.Cache.Version.Get(transactionID), Data: httpRules}, nil
		}
	}
	if err := c.ConfigParser.LoadData(c.getTransactionFile(transactionID)); err != nil {
		return nil, err
	}

	httpRules, err := c.parseHTTPResponseRules(parentType, parentName)
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
		c.Cache.HttpResponseRules.SetAll(parentName, parentType, transactionID, httpRules)
	}
	return &models.GetHTTPResponseRulesOKBody{Version: v, Data: httpRules}, nil
}

// GetHTTPResponseRule returns a struct with configuration version and a responseed http response rule
// in the specified parent. Returns error on fail or if http response rule does not exist.
func (c *Client) GetHTTPResponseRule(id int64, parentType, parentName string, transactionID string) (*models.GetHTTPResponseRuleOKBody, error) {
	if c.Cache.Enabled() {
		httpRule, found := c.Cache.HttpResponseRules.GetOne(id, parentName, parentType, transactionID)
		if found {
			return &models.GetHTTPResponseRuleOKBody{Version: c.Cache.Version.Get(transactionID), Data: httpRule}, nil
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

	data, err := c.ConfigParser.GetOne(p, parentName, "http-response", int(id))
	if err != nil {
		if err == parser_errors.SectionMissingErr {
			return nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("%s %s does not exist", parentType, parentName))
		}
		if err == parser_errors.FetchError {
			return nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("HTTP Response Rule %v does not exist in %s %s", id, parentType, parentName))
		}
		return nil, err
	}

	httpRule := parseHTTPResponseRule(data.(types.HTTPAction))
	httpRule.ID = &id

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return nil, err
	}

	if c.Cache.Enabled() {
		c.Cache.HttpResponseRules.Set(id, parentName, parentType, transactionID, httpRule)
	}

	return &models.GetHTTPResponseRuleOKBody{Version: v, Data: httpRule}, nil
}

// DeleteHTTPResponseRule deletes a http response rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) DeleteHTTPResponseRule(id int64, parentType string, parentName string, transactionID string, version int64) error {
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

	if err := c.ConfigParser.Delete(p, parentName, "http-response", int(id)); err != nil {
		if err == parser_errors.SectionMissingErr {
			return NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("%s %s does not exist", parentName, parentType))
		}
		if err == parser_errors.FetchError {
			return NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("HTTP Response Rule %v does not exist in %s %s", id, parentName, parentType))
		}
		return err
	}

	if err := c.saveData(t, transactionID); err != nil {
		return err
	}

	if c.Cache.Enabled() {
		c.Cache.HttpResponseRules.InvalidateParent(transactionID, parentName, parentType)
	}
	return nil
}

// CreateHTTPResponseRule creates a http response rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) CreateHTTPResponseRule(parentType string, parentName string, data *models.HTTPResponseRule, transactionID string, version int64) error {
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

	if err := c.ConfigParser.Insert(p, parentName, "http-response", serializeHTTPResponseRule(*data), int(*data.ID)); err != nil {
		if err == parser_errors.IndexOutOfRange {
			return c.errAndDeleteTransaction(NewConfError(ErrObjectIndexOutOfRange,
				fmt.Sprintf("HTTP Response Rule with id %v in %s %s out of range", int(*data.ID), parentType, parentName)), t, transactionID == "")
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
		c.Cache.HttpResponseRules.InvalidateParent(transactionID, parentName, parentType)
	}
	return nil
}

// EditHTTPResponseRule edits a http response rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) EditHTTPResponseRule(id int64, parentType string, parentName string, data *models.HTTPResponseRule, transactionID string, version int64) error {
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

	if _, err := c.ConfigParser.GetOne(p, parentName, "http-response", int(id)); err != nil {
		if err == parser_errors.SectionMissingErr {
			return NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("%s %s does not exist", parentType, parentName))
		}
		if err == parser_errors.FetchError {
			return NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("HTTP Response Rule %v does not exist in %s %s", id, parentType, parentName))
		}
		return err
	}

	if err := c.ConfigParser.Set(p, parentName, "http-response", serializeHTTPResponseRule(*data), int(id)); err != nil {
		if err == parser_errors.IndexOutOfRange {
			return c.errAndDeleteTransaction(NewConfError(ErrObjectIndexOutOfRange,
				fmt.Sprintf("HTTP Response Rule with id %v in %s %s out of range", int(*data.ID), parentType, parentName)), t, transactionID == "")
		}
		return c.errAndDeleteTransaction(err, t, transactionID == "")
	}

	if err := c.saveData(t, transactionID); err != nil {
		return err
	}
	if c.Cache.Enabled() {
		c.Cache.HttpResponseRules.InvalidateParent(transactionID, parentName, parentType)
	}
	return nil
}

func (c *Client) parseHTTPResponseRules(t, pName string) (models.HTTPResponseRules, error) {
	p := parser.Global
	if t == "frontend" {
		p = parser.Frontends
	} else if t == "backend" {
		p = parser.Backends
	}

	httpResRules := models.HTTPResponseRules{}
	data, err := c.ConfigParser.Get(p, pName, "http-response", false)
	if err != nil {
		if err == parser_errors.FetchError {
			return httpResRules, nil
		}
		return nil, err
	}

	rules := data.([]types.HTTPAction)
	for i, r := range rules {
		id := int64(i)
		httpResRule := parseHTTPResponseRule(r)
		if httpResRule != nil {
			httpResRule.ID = &id
			httpResRules = append(httpResRules, httpResRule)
		}
	}
	return httpResRules, nil
}

func parseHTTPResponseRule(f types.HTTPAction) *models.HTTPResponseRule {
	switch v := f.(type) {
	case *actions.Allow:
		return &models.HTTPResponseRule{
			Type:     "allow",
			Cond:     v.Cond,
			CondTest: v.CondTest,
		}
	case *actions.Deny:
		return &models.HTTPResponseRule{
			Type:     "deny",
			Cond:     v.Cond,
			CondTest: v.CondTest,
		}
	case *actions.Redirect:
		code, _ := strconv.ParseInt(v.Code, 10, 64)
		r := &models.HTTPResponseRule{
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
	case *actions.AddHeader:
		return &models.HTTPResponseRule{
			Type:      "add-header",
			HdrName:   v.Name,
			HdrFormat: v.Fmt,
			Cond:      v.Cond,
			CondTest:  v.CondTest,
		}
	case *actions.SetHeader:
		return &models.HTTPResponseRule{
			Type:      "set-header",
			HdrName:   v.Name,
			HdrFormat: v.Fmt,
			Cond:      v.Cond,
			CondTest:  v.CondTest,
		}
	case *actions.DelHeader:
		return &models.HTTPResponseRule{
			Type:     "del-header",
			HdrName:  v.Name,
			Cond:     v.Cond,
			CondTest: v.CondTest,
		}
	case *actions.ReplaceHeader:
		return &models.HTTPResponseRule{
			Type:      "replace-header",
			HdrName:   v.Name,
			HdrFormat: v.ReplaceFmt,
			HdrMatch:  v.MatchRegex,
			Cond:      v.Cond,
			CondTest:  v.CondTest,
		}
	case *actions.ReplaceValue:
		return &models.HTTPResponseRule{
			Type:      "replace-value",
			HdrName:   v.Name,
			HdrFormat: v.ReplaceFmt,
			HdrMatch:  v.MatchRegex,
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
	case *actions.SetVar:
		return &models.HTTPResponseRule{
			Type:     "set-var",
			VarName:  v.VarName,
			VarExpr:  strings.Join(v.Expr.Expr, " "),
			VarScope: v.VarScope,
			Cond:     v.Cond,
			CondTest: v.CondTest,
		}
	case *actions.SetStatus:
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
	case *actions.AddAcl:
		return &models.HTTPResponseRule{
			Type:      "add-acl",
			ACLFile:   v.FileName,
			ACLKeyfmt: v.KeyFmt,
			Cond:      v.Cond,
			CondTest:  v.CondTest,
		}
	case *actions.DelAcl:
		return &models.HTTPResponseRule{
			Type:      "del-acl",
			ACLFile:   v.FileName,
			ACLKeyfmt: v.KeyFmt,
			Cond:      v.Cond,
			CondTest:  v.CondTest,
		}
	case *actions.SendSpoeGroup:
		return &models.HTTPResponseRule{
			Type:       "send-spoe-group",
			SpoeEngine: v.Engine,
			SpoeGroup:  v.Group,
			Cond:       v.Cond,
			CondTest:   v.CondTest,
		}
	}
	return nil
}

func serializeHTTPResponseRule(f models.HTTPResponseRule) types.HTTPAction {
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
	case "redirect":
		return &actions.Redirect{
			Type:     f.RedirType,
			Value:    f.RedirValue,
			Code:     strconv.FormatInt(f.RedirCode, 10),
			Option:   f.RedirOption,
			Cond:     f.Cond,
			CondTest: f.CondTest,
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
	case "set-status":
		return &actions.SetStatus{
			Status:   strconv.FormatInt(f.Status, 10),
			Reason:   f.StatusReason,
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
