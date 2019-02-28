package configuration

import (
	"strconv"

	strfmt "github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/config-parser"
	parser_errors "github.com/haproxytech/config-parser/errors"
	"github.com/haproxytech/config-parser/types"
	"github.com/haproxytech/models"
)

// GetServerSwitchingRules returns a struct with configuration version and an array of
// configured server switching rules in the specified backend. Returns error on fail.
func (c *Client) GetServerSwitchingRules(backend string, transactionID string) (*models.GetServerSwitchingRulesOKBody, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return nil, err
	}

	srvRules, err := c.parseServerSwitchingRules(backend, p)
	if err != nil {
		return nil, c.handleError("", "backend", backend, "", false, err)
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return nil, err
	}
	return &models.GetServerSwitchingRulesOKBody{Version: v, Data: srvRules}, nil
}

// GetServerSwitchingRule returns a struct with configuration version and a requested server switching rule
// in the specified backend. Returns error on fail or if server switching rule does not exist.
func (c *Client) GetServerSwitchingRule(id int64, backend string, transactionID string) (*models.GetServerSwitchingRuleOKBody, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return nil, err
	}

	data, err := p.GetOne(parser.Backends, backend, "use-server", int(id))
	if err != nil {
		return nil, c.handleError(strconv.FormatInt(id, 10), "backend", backend, "", false, err)
	}

	srvRule := parseServerSwitchingRule(data.(types.UseServer))
	srvRule.ID = &id

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return nil, err
	}

	return &models.GetServerSwitchingRuleOKBody{Version: v, Data: srvRule}, nil
}

// DeleteServerSwitchingRule deletes a server switching rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) DeleteServerSwitchingRule(id int64, backend string, transactionID string, version int64) error {
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	if err := p.Delete(parser.Backends, backend, "use-server", int(id)); err != nil {
		return c.handleError(strconv.FormatInt(id, 10), "backend", backend, t, transactionID == "", err)
	}

	if err := c.saveData(p, t, transactionID == ""); err != nil {
		return err
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
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	if err := p.Insert(parser.Backends, backend, "use-server", serializeServerSwitchingRule(*data), int(*data.ID)); err != nil {
		return c.handleError(strconv.FormatInt(*data.ID, 10), "backend", backend, t, transactionID == "", err)
	}

	if err := c.saveData(p, t, transactionID == ""); err != nil {
		return err
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
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	if _, err := p.GetOne(parser.Backends, backend, "use-server", int(id)); err != nil {
		return c.handleError(strconv.FormatInt(*data.ID, 10), "backend", backend, t, transactionID == "", err)
	}

	if err := p.Set(parser.Backends, backend, "use-server", serializeServerSwitchingRule(*data), int(id)); err != nil {
		return c.handleError(strconv.FormatInt(*data.ID, 10), "backend", backend, t, transactionID == "", err)
	}

	if err := c.saveData(p, t, transactionID == ""); err != nil {
		return err
	}
	return nil
}

func (c *Client) parseServerSwitchingRules(backend string, p *parser.Parser) (models.ServerSwitchingRules, error) {
	sr := models.ServerSwitchingRules{}

	data, err := p.Get(parser.Backends, backend, "use-server", false)
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
