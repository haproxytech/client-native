package configuration

import (
	"fmt"
	"strconv"
	"strings"

	strfmt "github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/config-parser"
	parser_errors "github.com/haproxytech/config-parser/errors"
	"github.com/haproxytech/config-parser/parsers/tcp/actions"
	"github.com/haproxytech/config-parser/types"
	"github.com/haproxytech/models"
)

// GetTCPResponseRules returns a struct with configuration version and an array of
// configured tcp response rules in the specified backend. Returns error on fail.
func (c *Client) GetTCPResponseRules(backend string, transactionID string) (*models.GetTCPResponseRulesOKBody, error) {
	if c.Cache.Enabled() {
		tcpRules, found := c.Cache.TcpResponseRules.Get(backend, transactionID)
		if found {
			return &models.GetTCPResponseRulesOKBody{Version: c.Cache.Version.Get(transactionID), Data: tcpRules}, nil
		}
	}
	if err := c.ConfigParser.LoadData(c.getTransactionFile(transactionID)); err != nil {
		return nil, err
	}

	tcpRules, err := c.parseTCPResponseRules(backend)
	if err != nil {
		if err == parser_errors.SectionMissingErr {
			return nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Backend %s does not exist", backend))
		}
		return nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return nil, err
	}

	if c.Cache.Enabled() {
		c.Cache.TcpResponseRules.SetAll(backend, transactionID, tcpRules)
	}
	return &models.GetTCPResponseRulesOKBody{Version: v, Data: tcpRules}, nil
}

// GetTCPResponseRule returns a struct with configuration version and a requested tcp response rule
// in the specified backend. Returns error on fail or if tcp response rule does not exist.
func (c *Client) GetTCPResponseRule(id int64, backend string, transactionID string) (*models.GetTCPResponseRuleOKBody, error) {
	if c.Cache.Enabled() {
		tcpRule, found := c.Cache.TcpResponseRules.GetOne(id, backend, transactionID)
		if found {
			return &models.GetTCPResponseRuleOKBody{Version: c.Cache.Version.Get(transactionID), Data: tcpRule}, nil
		}
	}
	if err := c.ConfigParser.LoadData(c.getTransactionFile(transactionID)); err != nil {
		return nil, err
	}

	data, err := c.ConfigParser.GetOne(parser.Backends, backend, "tcp-response", int(id))
	if err != nil {
		if err == parser_errors.SectionMissingErr {
			return nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Backend %s does not exist", backend))
		}
		if err == parser_errors.FetchError {
			return nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("TCP Response Rule %v does not exist in backend %s", id, backend))
		}
		return nil, err
	}

	tcpRule := parseTCPResponseRule(data.(types.TCPAction))
	tcpRule.ID = &id

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return nil, err
	}

	if c.Cache.Enabled() {
		c.Cache.TcpResponseRules.Set(id, backend, transactionID, tcpRule)
	}

	return &models.GetTCPResponseRuleOKBody{Version: v, Data: tcpRule}, nil
}

// DeleteTCPResponseRule deletes a tcp response rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) DeleteTCPResponseRule(id int64, backend string, transactionID string, version int64) error {
	t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	if err := c.ConfigParser.Delete(parser.Backends, backend, "tcp-response", int(id)); err != nil {
		if err == parser_errors.SectionMissingErr {
			return NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Backend %s does not exist", backend))
		}
		if err == parser_errors.FetchError {
			return NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("TCP Response Rule %v does not exist in backend %s", id, backend))
		}
		return err
	}

	if err := c.saveData(t, transactionID); err != nil {
		return err
	}
	if c.Cache.Enabled() {
		c.Cache.TcpResponseRules.InvalidateBackend(transactionID, backend)
	}
	return nil
}

// CreateTCPResponseRule creates a tcp response rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) CreateTCPResponseRule(backend string, data *models.TCPResponseRule, transactionID string, version int64) error {
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

	if err := c.ConfigParser.Insert(parser.Backends, backend, "tcp-response", serializeTCPResponseRule(*data), int(*data.ID)); err != nil {
		if err == parser_errors.IndexOutOfRange {
			return c.errAndDeleteTransaction(NewConfError(ErrObjectIndexOutOfRange,
				fmt.Sprintf("TCP Response Rule with id %v in backend %s out of range", int(*data.ID), backend)), t, transactionID == "")
		}
		if err == parser_errors.SectionMissingErr {
			return NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Backend %s does not exist", backend))
		}
		return c.errAndDeleteTransaction(err, t, transactionID == "")
	}

	if err := c.saveData(t, transactionID); err != nil {
		return err
	}
	if c.Cache.Enabled() {
		c.Cache.TcpResponseRules.InvalidateBackend(transactionID, backend)
	}
	return nil
}

// EditTCPResponseRule edits a tcp response rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) EditTCPResponseRule(id int64, backend string, data *models.TCPResponseRule, transactionID string, version int64) error {
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

	if _, err := c.ConfigParser.GetOne(parser.Backends, backend, "tcp-response", int(id)); err != nil {
		if err == parser_errors.SectionMissingErr {
			return NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Backend %s does not exist", backend))
		}
		if err == parser_errors.FetchError {
			return NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("TCP Response Rule %v does not exist in backend %s", id, backend))
		}
		return err
	}

	if err := c.ConfigParser.Set(parser.Backends, backend, "tcp-response", serializeTCPResponseRule(*data), int(id)); err != nil {
		if err == parser_errors.IndexOutOfRange {
			return c.errAndDeleteTransaction(NewConfError(ErrObjectIndexOutOfRange,
				fmt.Sprintf("TCP Response Rule with id %v in backend %s out of range", int(*data.ID), backend)), t, transactionID == "")
		}
		return c.errAndDeleteTransaction(err, t, transactionID == "")
	}

	if err := c.saveData(t, transactionID); err != nil {
		return err
	}
	if c.Cache.Enabled() {
		c.Cache.TcpResponseRules.InvalidateBackend(transactionID, backend)
	}
	return nil
}

func (c *Client) parseTCPResponseRules(backend string) (models.TCPResponseRules, error) {
	tcpResRules := models.TCPResponseRules{}

	data, err := c.ConfigParser.Get(parser.Backends, backend, "tcp-response", false)
	if err != nil {
		if err == parser_errors.FetchError {
			return tcpResRules, nil
		}
		return nil, err
	}

	tRules := data.([]types.TCPAction)
	for i, tRule := range tRules {
		id := int64(i)
		tcpResRule := parseTCPResponseRule(tRule)
		if tcpResRule != nil {
			tcpResRule.ID = &id
			tcpResRules = append(tcpResRules, tcpResRule)
		}
	}
	return tcpResRules, nil
}

func parseTCPResponseRule(t types.TCPAction) *models.TCPResponseRule {
	switch v := t.(type) {
	case *actions.Content:
		r := &models.TCPResponseRule{
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
		tOut, _ := strconv.ParseInt(v.Timeout, 10, 64)
		return &models.TCPResponseRule{
			Timeout: &tOut,
		}
	}
	return nil
}

func serializeTCPResponseRule(t models.TCPResponseRule) types.TCPAction {
	switch t.Type {
	case "content":
		return &actions.Content{
			Action:   []string{t.Action},
			Cond:     t.Cond,
			CondTest: t.CondTest,
		}
	case "inspect-delay":
		return &actions.InspectDelay{
			Timeout: strconv.FormatInt(*t.Timeout, 10),
		}
	}
	return nil
}
