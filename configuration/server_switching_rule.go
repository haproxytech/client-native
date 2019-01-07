package configuration

import (
	"strconv"
	"strings"

	strfmt "github.com/go-openapi/strfmt"
	"github.com/haproxytech/models"
)

// GetServerSwitchingRules returns a struct with configuration version and an array of
// configured server switching rules in the specified backend. Returns error on fail.
func (c *Client) GetServerSwitchingRules(backend string, transactionID string) (*models.GetServerSwitchingRulesOKBody, error) {
	if c.Cache.Enabled() {
		srvRules, found := c.Cache.ServerSwitchingRules.Get(backend, transactionID)
		if found {
			return &models.GetServerSwitchingRulesOKBody{Version: c.Cache.Version.Get(), Data: srvRules}, nil
		}
	}
	srvRulesString, err := c.executeLBCTL("l7-farm-useserver-dump", transactionID, backend)
	if err != nil {
		return nil, err
	}

	srvRules := c.parseServerSwitchingRules(srvRulesString)

	v, err := c.GetVersion()
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
			return &models.GetServerSwitchingRuleOKBody{Version: c.Cache.Version.Get(), Data: srvRule}, nil
		}
	}
	srvRuleStr, err := c.executeLBCTL("l7-farm-useserver-show", transactionID, backend, strconv.FormatInt(id, 10))
	if err != nil {
		return nil, err
	}
	srvRule := &models.ServerSwitchingRule{ID: id}

	c.parseObject(srvRuleStr, srvRule)

	v, err := c.GetVersion()
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
	err := c.deleteObject(strconv.FormatInt(id, 10), "useserver", backend, "farm", transactionID, version)
	if err != nil {
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
	err := c.createObject(strconv.FormatInt(data.ID, 10), "useserver", backend, "farm", data, nil, transactionID, version)
	if err != nil {
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
	ondiskSr, err := c.GetServerSwitchingRule(id, backend, transactionID)
	if err != nil {
		return err
	}

	err = c.editObject(strconv.FormatInt(data.ID, 10), "useserver", backend, "farm", data, ondiskSr, nil, transactionID, version)
	if err != nil {
		return err
	}
	if c.Cache.Enabled() {
		c.Cache.ServerSwitchingRules.InvalidateBackend(transactionID, backend)
	}
	return nil
}

func (c *Client) parseServerSwitchingRules(response string) models.ServerSwitchingRules {
	srvRules := make(models.ServerSwitchingRules, 0, 1)
	for _, srvRulesStr := range strings.Split(response, "\n\n") {
		if strings.TrimSpace(srvRulesStr) == "" {
			continue
		}
		idStr, _ := splitHeaderLine(srvRulesStr)
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			id = 0
		}

		srvRulesObj := &models.ServerSwitchingRule{ID: id}
		c.parseObject(srvRulesStr, srvRulesObj)
		srvRules = append(srvRules, srvRulesObj)
	}
	return srvRules
}
