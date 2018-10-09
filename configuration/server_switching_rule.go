package configuration

import (
	"strconv"
	"strings"

	strfmt "github.com/go-openapi/strfmt"
	"github.com/haproxytech/models"
)

// GetServerSwitchingRules returns a struct with configuration version and an array of
// configured server switching rules in the specified backend. Returns error on fail.
func (c *LBCTLConfigurationClient) GetServerSwitchingRules(backend string) (*models.GetServerSwitchingRulesOKBody, error) {
	srvRulesString, err := c.executeLBCTL("l7-farm-useserver-dump", "", backend)
	if err != nil {
		return nil, err
	}

	srvRules := c.parseServerSwitchingRules(srvRulesString)

	v, err := c.GetVersion()
	if err != nil {
		return nil, err
	}

	return &models.GetServerSwitchingRulesOKBody{Version: v, Data: srvRules}, nil
}

// GetServerSwitchingRule returns a struct with configuration version and a requested server switching rule
// in the specified backend. Returns error on fail or if server switching rule does not exist.
func (c *LBCTLConfigurationClient) GetServerSwitchingRule(id int64, backend string) (*models.GetServerSwitchingRuleOKBody, error) {
	srvRuleStr, err := c.executeLBCTL("l7-farm-useserver-show", "", backend, strconv.FormatInt(id, 10))
	if err != nil {
		return nil, err
	}
	srvRule := &models.ServerSwitchingRule{ID: id}

	c.parseObject(srvRuleStr, srvRule)

	v, err := c.GetVersion()
	if err != nil {
		return nil, err
	}

	return &models.GetServerSwitchingRuleOKBody{Version: v, Data: srvRule}, nil
}

// DeleteServerSwitchingRule deletes a server switching rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *LBCTLConfigurationClient) DeleteServerSwitchingRule(id int64, backend string, transactionID string, version int64) error {
	return c.deleteObject(strconv.FormatInt(id, 10), "useserver", backend, "farm", transactionID, version)
}

// CreateServerSwitchingRule creates a server switching rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *LBCTLConfigurationClient) CreateServerSwitchingRule(backend string, data *models.ServerSwitchingRule, transactionID string, version int64) error {
	validationErr := data.Validate(strfmt.Default)
	if validationErr != nil {
		return NewConfError(ErrValidationError, validationErr.Error())
	}
	return c.createObject(strconv.FormatInt(data.ID, 10), "useserver", backend, "farm", data, nil, transactionID, version)
}

// EditServerSwitchingRule edits a server switching rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *LBCTLConfigurationClient) EditServerSwitchingRule(id int64, backend string, data *models.ServerSwitchingRule, transactionID string, version int64) error {
	validationErr := data.Validate(strfmt.Default)
	if validationErr != nil {
		return NewConfError(ErrValidationError, validationErr.Error())
	}
	ondiskSr, err := c.GetServerSwitchingRule(id, backend)
	if err != nil {
		return err
	}

	return c.editObject(strconv.FormatInt(data.ID, 10), "useserver", backend, "farm", data, ondiskSr, nil, transactionID, version)
}

func (c *LBCTLConfigurationClient) parseServerSwitchingRules(response string) models.ServerSwitchingRules {
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
