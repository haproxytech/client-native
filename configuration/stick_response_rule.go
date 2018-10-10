package configuration

import (
	"strconv"
	"strings"

	strfmt "github.com/go-openapi/strfmt"
	"github.com/haproxytech/models"
)

// GetStickResponseRules returns a struct with configuration version and an array of
// configured stick response rules in the specified backend. Returns error on fail.
func (c *LBCTLConfigurationClient) GetStickResponseRules(backend string, transactionID string) (*models.GetStickResponseRulesOKBody, error) {
	stickReqRulesString, err := c.executeLBCTL("l7-farm-stickrsp-dump", transactionID, backend)
	if err != nil {
		return nil, err
	}

	stickReqRules := c.parseStickResponseRules(stickReqRulesString)

	v, err := c.GetVersion()
	if err != nil {
		return nil, err
	}

	return &models.GetStickResponseRulesOKBody{Version: v, Data: stickReqRules}, nil
}

// GetStickResponseRule returns a struct with configuration version and a responseed stick response rule
// in the specified backend. Returns error on fail or if stick response rule does not exist.
func (c *LBCTLConfigurationClient) GetStickResponseRule(id int64, backend string, transactionID string) (*models.GetStickResponseRuleOKBody, error) {
	stickReqRuleStr, err := c.executeLBCTL("l7-farm-stickrsp-show", transactionID, backend, strconv.FormatInt(id, 10))
	if err != nil {
		return nil, err
	}
	stickReqRule := &models.StickResponseRule{ID: id}

	c.parseObject(stickReqRuleStr, stickReqRule)

	v, err := c.GetVersion()
	if err != nil {
		return nil, err
	}

	return &models.GetStickResponseRuleOKBody{Version: v, Data: stickReqRule}, nil
}

// DeleteStickResponseRule deletes a stick response rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *LBCTLConfigurationClient) DeleteStickResponseRule(id int64, backend string, transactionID string, version int64) error {
	return c.deleteObject(strconv.FormatInt(id, 10), "stickrsp", backend, "farm", transactionID, version)
}

// CreateStickResponseRule creates a stick response rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *LBCTLConfigurationClient) CreateStickResponseRule(backend string, data *models.StickResponseRule, transactionID string, version int64) error {
	validationErr := data.Validate(strfmt.Default)
	if validationErr != nil {
		return NewConfError(ErrValidationError, validationErr.Error())
	}
	return c.createObject(strconv.FormatInt(data.ID, 10), "stickrsp", backend, "farm", data, nil, transactionID, version)
}

// EditStickResponseRule edits a stick response rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *LBCTLConfigurationClient) EditStickResponseRule(id int64, backend string, data *models.StickResponseRule, transactionID string, version int64) error {
	validationErr := data.Validate(strfmt.Default)
	if validationErr != nil {
		return NewConfError(ErrValidationError, validationErr.Error())
	}
	ondiskR, err := c.GetStickResponseRule(id, backend, transactionID)
	if err != nil {
		return err
	}

	return c.editObject(strconv.FormatInt(data.ID, 10), "stickrsp", backend, "farm", data, ondiskR, nil, transactionID, version)
}

func (c *LBCTLConfigurationClient) parseStickResponseRules(response string) models.StickResponseRules {
	stickReqRules := make(models.StickResponseRules, 0, 1)
	for _, stickReqRulesStr := range strings.Split(response, "\n\n") {
		if strings.TrimSpace(stickReqRulesStr) == "" {
			continue
		}
		idStr, _ := splitHeaderLine(stickReqRulesStr)
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			id = 0
		}

		stickReqRulesObj := &models.StickResponseRule{ID: id}
		c.parseObject(stickReqRulesStr, stickReqRulesObj)
		stickReqRules = append(stickReqRules, stickReqRulesObj)
	}
	return stickReqRules
}
