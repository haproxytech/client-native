package configuration

import (
	"strconv"
	"strings"

	strfmt "github.com/go-openapi/strfmt"
	"github.com/haproxytech/models"
)

// GetBackendSwitchingRules returns a struct with configuration version and an array of
// configured backend switching rules in the specified frontend. Returns error on fail.
func (c *LBCTLClient) GetBackendSwitchingRules(frontend string, transactionID string) (*models.GetBackendSwitchingRulesOKBody, error) {
	bckRulesString, err := c.executeLBCTL("l7-service-usefarm-dump", transactionID, frontend)
	if err != nil {
		return nil, err
	}

	bckRules := c.parseBackendSwitchingRules(bckRulesString)

	v, err := c.GetVersion()
	if err != nil {
		return nil, err
	}

	return &models.GetBackendSwitchingRulesOKBody{Version: v, Data: bckRules}, nil
}

// GetBackendSwitchingRule returns a struct with configuration version and a requested backend switching rule
// in the specified frontend. Returns error on fail or if backend switching rule does not exist.
func (c *LBCTLClient) GetBackendSwitchingRule(id int64, frontend string, transactionID string) (*models.GetBackendSwitchingRuleOKBody, error) {
	bckRuleStr, err := c.executeLBCTL("l7-service-usefarm-show", transactionID, frontend, strconv.FormatInt(id, 10))
	if err != nil {
		return nil, err
	}
	bckRule := &models.BackendSwitchingRule{ID: id}

	c.parseObject(bckRuleStr, bckRule)

	v, err := c.GetVersion()
	if err != nil {
		return nil, err
	}

	return &models.GetBackendSwitchingRuleOKBody{Version: v, Data: bckRule}, nil
}

// DeleteBackendSwitchingRule deletes a backend switching rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *LBCTLClient) DeleteBackendSwitchingRule(id int64, frontend string, transactionID string, version int64) error {
	return c.deleteObject(strconv.FormatInt(id, 10), "usefarm", frontend, "service", transactionID, version)
}

// CreateBackendSwitchingRule creates a backend switching rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *LBCTLClient) CreateBackendSwitchingRule(frontend string, data *models.BackendSwitchingRule, transactionID string, version int64) error {
	validationErr := data.Validate(strfmt.Default)
	if validationErr != nil {
		return NewConfError(ErrValidationError, validationErr.Error())
	}
	return c.createObject(strconv.FormatInt(data.ID, 10), "usefarm", frontend, "service", data, nil, transactionID, version)
}

// EditBackendSwitchingRule edits a backend switching rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *LBCTLClient) EditBackendSwitchingRule(id int64, frontend string, data *models.BackendSwitchingRule, transactionID string, version int64) error {
	validationErr := data.Validate(strfmt.Default)
	if validationErr != nil {
		return NewConfError(ErrValidationError, validationErr.Error())
	}
	ondiskBr, err := c.GetBackendSwitchingRule(id, frontend, transactionID)
	if err != nil {
		return err
	}

	return c.editObject(strconv.FormatInt(data.ID, 10), "usefarm", frontend, "service", data, ondiskBr, nil, transactionID, version)
}

func (c *LBCTLClient) parseBackendSwitchingRules(response string) models.BackendSwitchingRules {
	bckRules := make(models.BackendSwitchingRules, 0, 1)
	for _, bckRulesStr := range strings.Split(response, "\n\n") {
		if strings.TrimSpace(bckRulesStr) == "" {
			continue
		}
		idStr, _ := splitHeaderLine(bckRulesStr)
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			id = 0
		}

		bckRulesObj := &models.BackendSwitchingRule{ID: id}
		c.parseObject(bckRulesStr, bckRulesObj)
		bckRules = append(bckRules, bckRulesObj)
	}
	return bckRules
}
