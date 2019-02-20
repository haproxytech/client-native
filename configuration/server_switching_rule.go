package configuration

import (
	"fmt"

	strfmt "github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/config-parser"
	parser_errors "github.com/haproxytech/config-parser/errors"
	"github.com/haproxytech/config-parser/types"
	"github.com/haproxytech/models"
)

// GetServerSwitchingRules returns a struct with configuration version and an array of
// configured server switching rules in the specified backend. Returns error on fail.
func (c *Client) GetServerSwitchingRules(backend string, transactionID string) (*models.GetServerSwitchingRulesOKBody, error) {
	if c.Cache.Enabled() {
		srvRules, found := c.Cache.ServerSwitchingRules.Get(backend, transactionID)
		if found {
			return &models.GetServerSwitchingRulesOKBody{Version: c.Cache.Version.Get(transactionID), Data: srvRules}, nil
		}
	}
	if err := c.ConfigParser.LoadData(c.getTransactionFile(transactionID)); err != nil {
		return nil, err
	}

	srvRules, err := c.parseServerSwitchingRules(backend)
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
		c.Cache.ServerSwitchingRules.SetAll(backend, transactionID, srvRules)
	}
	return &models.GetServerSwitchingRulesOKBody{Version: v, Data: srvRules}, nil
}

// GetServerSwitchingRule returns a struct with configuration version and a requested server switching rule
// in the specified backend. Returns error on fail or if server switching rule does not exist.
func (c *Client) GetServerSwitchingRule(id int64, backend string, transactionID string) (*models.GetServerSwitchingRuleOKBody, error) {
	if c.Cache.Enabled() {
		srvRule, found := c.Cache.ServerSwitchingRules.GetOne(id, backend, transactionID)
		if found {
			return &models.GetServerSwitchingRuleOKBody{Version: c.Cache.Version.Get(transactionID), Data: srvRule}, nil
		}
	}
	if err := c.ConfigParser.LoadData(c.getTransactionFile(transactionID)); err != nil {
		return nil, err
	}

	data, err := c.ConfigParser.GetOne(parser.Backends, backend, "use-server", int(id))
	if err != nil {
		if err == parser_errors.SectionMissingErr {
			return nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Backend %s does not exist", backend))
		}
		if err == parser_errors.FetchError {
			return nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Server Switching Rule %v does not exist in backend %s", id, backend))
		}
		return nil, err
	}

	srvRule := parseServerSwitchingRule(data.(types.UseServer))
	srvRule.ID = &id

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return nil, err
	}

	if c.Cache.Enabled() {
		c.Cache.ServerSwitchingRules.Set(id, backend, transactionID, srvRule)
	}

	return &models.GetServerSwitchingRuleOKBody{Version: v, Data: srvRule}, nil
}

// DeleteServerSwitchingRule deletes a server switching rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) DeleteServerSwitchingRule(id int64, backend string, transactionID string, version int64) error {
	t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	if err := c.ConfigParser.Delete(parser.Backends, backend, "use-server", int(id)); err != nil {
		if err == parser_errors.SectionMissingErr {
			return NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Backend %s does not exist", backend))
		}
		if err == parser_errors.FetchError {
			return NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Server Switching Rule %v does not exist in backend %s", id, backend))
		}
		return err
	}

	if err := c.saveData(t, transactionID); err != nil {
		return err
	}
	if c.Cache.Enabled() {
		c.Cache.ServerSwitchingRules.InvalidateBackend(transactionID, backend)
	}
	return nil
}

// CreateServerSwitchingRule creates a server switching rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) CreateServerSwitchingRule(backend string, data *models.ServerSwitchingRule, transactionID string, version int64) error {
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

	if err := c.ConfigParser.Insert(parser.Backends, backend, "use-server", serializeServerSwitchingRule(*data), int(*data.ID)); err != nil {
		if err == parser_errors.IndexOutOfRange {
			return c.errAndDeleteTransaction(NewConfError(ErrObjectIndexOutOfRange,
				fmt.Sprintf("Server Switching Rule with id %v in backend %s out of range", int(*data.ID), backend)), t, transactionID == "")
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
		c.Cache.ServerSwitchingRules.InvalidateBackend(transactionID, backend)
	}
	return nil
}

// EditServerSwitchingRule edits a server switching rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) EditServerSwitchingRule(id int64, backend string, data *models.ServerSwitchingRule, transactionID string, version int64) error {
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

	if _, err := c.ConfigParser.GetOne(parser.Backends, backend, "use-server", int(id)); err != nil {
		if err == parser_errors.SectionMissingErr {
			return NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Backend %s does not exist", backend))
		}
		if err == parser_errors.FetchError {
			return NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Server Switching Rule %v does not exist in backend %s", id, backend))
		}
		return err
	}

	if err := c.ConfigParser.Set(parser.Backends, backend, "use-server", serializeServerSwitchingRule(*data), int(id)); err != nil {
		if err == parser_errors.IndexOutOfRange {
			return c.errAndDeleteTransaction(NewConfError(ErrObjectIndexOutOfRange,
				fmt.Sprintf("Server Switching Rule with id %v in backend %s out of range", int(*data.ID), backend)), t, transactionID == "")
		}
		return c.errAndDeleteTransaction(err, t, transactionID == "")
	}

	if err := c.saveData(t, transactionID); err != nil {
		return err
	}
	if c.Cache.Enabled() {
		c.Cache.ServerSwitchingRules.InvalidateBackend(transactionID, backend)
	}
	return nil
}

func (c *Client) parseServerSwitchingRules(backend string) (models.ServerSwitchingRules, error) {
	sr := models.ServerSwitchingRules{}

	data, err := c.ConfigParser.Get(parser.Backends, backend, "use-server", false)
	if err != nil {
		if err == parser_errors.FetchError {
			return sr, nil
		}
		return nil, err
	}

	sRules := data.([]types.UseServer)
	for i, sRule := range sRules {
		id := int64(i)
		s := parseServerSwitchingRule(sRule)
		if s != nil {
			s.ID = &id
			sr = append(sr, s)
		}
	}
	return sr, nil
}

func parseServerSwitchingRule(us types.UseServer) *models.ServerSwitchingRule {
	return &models.ServerSwitchingRule{
		TargetServer: us.Name,
		Cond:         us.Cond,
		CondTest:     us.CondTest,
	}
}

func serializeServerSwitchingRule(sRule models.ServerSwitchingRule) types.UseServer {
	return types.UseServer{
		Name:     sRule.TargetServer,
		Cond:     sRule.Cond,
		CondTest: sRule.CondTest,
	}
}
