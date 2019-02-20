package configuration

import (
	"fmt"

	strfmt "github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/config-parser"
	parser_errors "github.com/haproxytech/config-parser/errors"
	"github.com/haproxytech/config-parser/types"
	"github.com/haproxytech/models"
)

// GetStickRules returns a struct with configuration version and an array of
// configured stick rules in the specified backend. Returns error on fail.
func (c *Client) GetStickRules(backend string, transactionID string) (*models.GetStickRulesOKBody, error) {
	if c.Cache.Enabled() {
		stickRules, found := c.Cache.StickRules.Get(backend, transactionID)
		if found {
			return &models.GetStickRulesOKBody{Version: c.Cache.Version.Get(transactionID), Data: stickRules}, nil
		}
	}
	if err := c.ConfigParser.LoadData(c.getTransactionFile(transactionID)); err != nil {
		return nil, err
	}

	sRules, err := c.parseStickRules(backend)
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
		c.Cache.StickRules.SetAll(backend, transactionID, sRules)
	}
	return &models.GetStickRulesOKBody{Version: v, Data: sRules}, nil
}

// GetStickRule returns a struct with configuration version and a requested stick rule
// in the specified backend. Returns error on fail or if stick rule does not exist.
func (c *Client) GetStickRule(id int64, backend string, transactionID string) (*models.GetStickRuleOKBody, error) {
	if c.Cache.Enabled() {
		stickRule, found := c.Cache.StickRules.GetOne(id, backend, transactionID)
		if found {
			return &models.GetStickRuleOKBody{Version: c.Cache.Version.Get(transactionID), Data: stickRule}, nil
		}
	}
	if err := c.ConfigParser.LoadData(c.getTransactionFile(transactionID)); err != nil {
		return nil, err
	}

	data, err := c.ConfigParser.GetOne(parser.Backends, backend, "stick", int(id))
	if err != nil {
		if err == parser_errors.SectionMissingErr {
			return nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Backend %s does not exist", backend))
		}
		if err == parser_errors.FetchError {
			return nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Stick Rule %v does not exist in backend %s", id, backend))
		}
		return nil, err
	}

	sRule := parseStickRule(data.(types.Stick))
	sRule.ID = &id

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return nil, err
	}

	if c.Cache.Enabled() {
		c.Cache.StickRules.Set(id, backend, transactionID, sRule)
	}

	return &models.GetStickRuleOKBody{Version: v, Data: sRule}, nil
}

// DeleteStickRule deletes a stick rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) DeleteStickRule(id int64, backend string, transactionID string, version int64) error {
	t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	if err := c.ConfigParser.Delete(parser.Backends, backend, "stick", int(id)); err != nil {
		if err == parser_errors.SectionMissingErr {
			return NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Backend %s does not exist", backend))
		}
		if err == parser_errors.FetchError {
			return NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Stick Rule %v does not exist in backend %s", id, backend))
		}
		return err
	}

	if err := c.saveData(t, transactionID); err != nil {
		return err
	}
	if c.Cache.Enabled() {
		c.Cache.StickRules.InvalidateBackend(transactionID, backend)
	}
	return nil
}

// CreateStickRule creates a stick rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) CreateStickRule(backend string, data *models.StickRule, transactionID string, version int64) error {
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

	if err := c.ConfigParser.Insert(parser.Backends, backend, "stick", serializeStickRule(*data), int(*data.ID)); err != nil {
		if err == parser_errors.IndexOutOfRange {
			return c.errAndDeleteTransaction(NewConfError(ErrObjectIndexOutOfRange,
				fmt.Sprintf("Stick Rule with id %v in backend %s out of range", int(*data.ID), backend)), t, transactionID == "")
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
		c.Cache.StickRules.InvalidateBackend(transactionID, backend)
	}
	return nil
}

// EditStickRule edits a stick rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) EditStickRule(id int64, backend string, data *models.StickRule, transactionID string, version int64) error {
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

	if _, err := c.ConfigParser.GetOne(parser.Backends, backend, "stick", int(id)); err != nil {
		if err == parser_errors.SectionMissingErr {
			return NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Backend %s does not exist", backend))
		}
		if err == parser_errors.FetchError {
			return NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Stick Rule %v does not exist in backend %s", id, backend))
		}
		return err
	}

	if err := c.ConfigParser.Set(parser.Backends, backend, "stick", serializeStickRule(*data), int(id)); err != nil {
		if err == parser_errors.IndexOutOfRange {
			return c.errAndDeleteTransaction(NewConfError(ErrObjectIndexOutOfRange,
				fmt.Sprintf("Stick Rule with id %v in backend %s out of range", int(*data.ID), backend)), t, transactionID == "")
		}
		return c.errAndDeleteTransaction(err, t, transactionID == "")
	}

	if err := c.saveData(t, transactionID); err != nil {
		return err
	}
	if c.Cache.Enabled() {
		c.Cache.StickRules.InvalidateBackend(transactionID, backend)
	}
	return nil
}

func (c *Client) parseStickRules(backend string) (models.StickRules, error) {
	sr := models.StickRules{}

	data, err := c.ConfigParser.Get(parser.Backends, backend, "stick", false)
	if err != nil {
		if err == parser_errors.FetchError {
			return sr, nil
		}
		return nil, err
	}

	sRules := data.([]types.Stick)
	for i, sRule := range sRules {
		id := int64(i)
		s := parseStickRule(sRule)
		if s != nil {
			s.ID = &id
			sr = append(sr, s)
		}
	}
	return sr, nil
}

func parseStickRule(s types.Stick) *models.StickRule {
	return &models.StickRule{
		Type:     s.Type,
		Table:    s.Table,
		Pattern:  s.Pattern,
		Cond:     s.Cond,
		CondTest: s.CondTest,
	}
}

func serializeStickRule(sRule models.StickRule) types.Stick {
	return types.Stick{
		Type:     sRule.Type,
		Table:    sRule.Table,
		Pattern:  sRule.Pattern,
		Cond:     sRule.Cond,
		CondTest: sRule.CondTest,
	}
}
