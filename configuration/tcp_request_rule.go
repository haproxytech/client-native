package configuration

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/haproxytech/config-parser/parsers/tcp/actions"

	strfmt "github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/config-parser"
	parser_errors "github.com/haproxytech/config-parser/errors"
	"github.com/haproxytech/config-parser/types"
	"github.com/haproxytech/models"
)

// GetTCPRequestRules returns a struct with configuration version and an array of
// configured TCP request rules in the specified parent. Returns error on fail.
func (c *Client) GetTCPRequestRules(parentType, parentName string, transactionID string) (*models.GetTCPRequestRulesOKBody, error) {
	if c.Cache.Enabled() {
		tcpRules, found := c.Cache.TcpRequestRules.Get(parentName, parentType, transactionID)
		if found {
			return &models.GetTCPRequestRulesOKBody{Version: c.Cache.Version.Get(transactionID), Data: tcpRules}, nil
		}
	}
	if err := c.ConfigParser.LoadData(c.getTransactionFile(transactionID)); err != nil {
		return nil, err
	}

	tcpRules, err := c.parseTCPRequestRules(parentType, parentName)
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
		c.Cache.TcpRequestRules.SetAll(parentName, parentType, transactionID, tcpRules)
	}
	return &models.GetTCPRequestRulesOKBody{Version: v, Data: tcpRules}, nil
}

// GetTCPRequestRule returns a struct with configuration version and a requested tcp request rule
// in the specified parent. Returns error on fail or if http request rule does not exist.
func (c *Client) GetTCPRequestRule(id int64, parentType, parentName string, transactionID string) (*models.GetTCPRequestRuleOKBody, error) {
	if c.Cache.Enabled() {
		tcpRule, found := c.Cache.TcpRequestRules.GetOne(id, parentName, parentType, transactionID)
		if found {
			return &models.GetTCPRequestRuleOKBody{Version: c.Cache.Version.Get(transactionID), Data: tcpRule}, nil
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

	data, err := c.ConfigParser.GetOne(p, parentName, "tcp-request", int(id))
	if err != nil {
		if err == parser_errors.SectionMissingErr {
			return nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("%s %s does not exist", parentType, parentName))
		}
		if err == parser_errors.FetchError {
			return nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("TCP Request Rule %v does not exist in %s %s", id, parentType, parentName))
		}
		return nil, err
	}

	tcpRule := parseTCPRequestRule(data.(types.TCPAction))
	tcpRule.ID = &id

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return nil, err
	}

	if c.Cache.Enabled() {
		c.Cache.TcpRequestRules.Set(id, parentName, parentType, transactionID, tcpRule)
	}

	return &models.GetTCPRequestRuleOKBody{Version: v, Data: tcpRule}, nil
}

// DeleteTCPRequestRule deletes a tcp request rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) DeleteTCPRequestRule(id int64, parentType string, parentName string, transactionID string, version int64) error {
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

	if err := c.ConfigParser.Delete(p, parentName, "tcp-request", int(id)); err != nil {
		if err == parser_errors.SectionMissingErr {
			return NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("%s %s does not exist", parentName, parentType))
		}
		if err == parser_errors.FetchError {
			return NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("TCP Request Rule %v does not exist in %s %s", id, parentName, parentType))
		}
		return err
	}

	if err := c.saveData(t, transactionID); err != nil {
		return err
	}
	if c.Cache.Enabled() {
		c.Cache.TcpRequestRules.InvalidateParent(transactionID, parentName, parentType)
	}
	return nil
}

// CreateTCPRequestRule creates a tcp request rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) CreateTCPRequestRule(parentType string, parentName string, data *models.TCPRequestRule, transactionID string, version int64) error {
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

	if err := c.ConfigParser.Insert(p, parentName, "tcp-request", serializeTCPRequestRule(*data), int(*data.ID)); err != nil {
		if err == parser_errors.IndexOutOfRange {
			return c.errAndDeleteTransaction(NewConfError(ErrObjectIndexOutOfRange,
				fmt.Sprintf("TCP Request Rule with id %v in %s %s out of range", int(*data.ID), parentType, parentName)), t, transactionID == "")
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
		c.Cache.TcpRequestRules.InvalidateParent(transactionID, parentName, parentType)
	}
	return nil
}

// EditTCPRequestRule edits a tcp request rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) EditTCPRequestRule(id int64, parentType string, parentName string, data *models.TCPRequestRule, transactionID string, version int64) error {
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

	if _, err := c.ConfigParser.GetOne(p, parentName, "tcp-request", int(id)); err != nil {
		if err == parser_errors.SectionMissingErr {
			return NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("%s %s does not exist", parentType, parentName))
		}
		if err == parser_errors.FetchError {
			return NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("TCP Request Rule %v does not exist in %s %s", id, parentType, parentName))
		}
		return err
	}

	if err := c.ConfigParser.Set(p, parentName, "tcp-request", serializeTCPRequestRule(*data), int(id)); err != nil {
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
		c.Cache.TcpRequestRules.InvalidateParent(transactionID, parentName, parentType)
	}
	return nil
}

func (c *Client) parseTCPRequestRules(t, pName string) (models.TCPRequestRules, error) {
	p := parser.Global
	if t == "frontend" {
		p = parser.Frontends
	} else if t == "backend" {
		p = parser.Backends
	}

	tcpReqRules := models.TCPRequestRules{}
	data, err := c.ConfigParser.Get(p, pName, "tcp-request", false)
	if err != nil {
		if err == parser_errors.FetchError {
			return tcpReqRules, nil
		}
		return nil, err
	}

	rules := data.([]types.TCPAction)
	for i, r := range rules {
		id := int64(i)
		tcpReqRule := parseTCPRequestRule(r)
		if tcpReqRule != nil {
			tcpReqRule.ID = &id
			tcpReqRules = append(tcpReqRules, tcpReqRule)
		}
	}
	return tcpReqRules, nil
}

func parseTCPRequestRule(f types.TCPAction) *models.TCPRequestRule {
	switch v := f.(type) {
	case *actions.Connection:
		r := &models.TCPRequestRule{
			Type:     "connection",
			Cond:     v.Cond,
			CondTest: v.CondTest,
		}
		if strings.Join(v.Action, " ") == "accept" {
			r.Action = "accept"
		} else if strings.Join(v.Action, " ") == "reject" {
			r.Action = "reject"
		} else {
			return nil
		}
		return r
	case *actions.Content:
		r := &models.TCPRequestRule{
			Type:     "content",
			Cond:     v.Cond,
			CondTest: v.CondTest,
		}
		if strings.Join(v.Action, " ") == "accept" {
			r.Action = "accept"
		} else if strings.Join(v.Action, " ") == "reject" {
			r.Action = "reject"
		} else {
			return nil
		}
		return r
	case *actions.InspectDelay:
		t, _ := strconv.ParseInt(v.Timeout, 10, 64)
		return &models.TCPRequestRule{
			Timeout: &t,
		}
	case *actions.Session:
		r := &models.TCPRequestRule{
			Type:     "session",
			Cond:     v.Cond,
			CondTest: v.CondTest,
		}
		if strings.Join(v.Action, " ") == "accept" {
			r.Action = "accept"
		} else if strings.Join(v.Action, " ") == "reject" {
			r.Action = "reject"
		} else {
			return nil
		}
		return r
	}
	return nil
}

func serializeTCPRequestRule(f models.TCPRequestRule) types.TCPAction {
	switch f.Type {
	case "connection":
		return &actions.Connection{
			Action:   []string{f.Action},
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "content":
		return &actions.Content{
			Action:   []string{f.Action},
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "inspect-delay":
		return &actions.InspectDelay{
			Timeout: strconv.FormatInt(*f.Timeout, 10),
		}
	case "session":
		return &actions.Session{
			Action:   []string{f.Action},
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	}
	return nil
}
