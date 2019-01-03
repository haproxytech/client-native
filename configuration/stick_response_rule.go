package configuration

import (
	"strconv"
	"strings"

	strfmt "github.com/go-openapi/strfmt"
	"github.com/haproxytech/models"
)

// GetStickResponseRules returns a struct with configuration version and an array of
// configured stick response rules in the specified backend. Returns error on fail.
func (c *LBCTLClient) GetStickResponseRules(backend string, transactionID string) (*models.GetStickResponseRulesOKBody, error) {
	if c.Cache.Enabled() {
		stickResRules, found := c.Cache.StickResponseRules.Get(backend, transactionID)
		if found {
			return &models.GetStickResponseRulesOKBody{Version: c.Cache.Version.Get(), Data: stickResRules}, nil
		}
	}
	stickResRulesString, err := c.executeLBCTL("l7-farm-stickrsp-dump", transactionID, backend)
	if err != nil {
		return nil, err
	}

	stickResRules := c.parseStickResponseRules(stickResRulesString)

	v, err := c.GetVersion()
	if err != nil {
		return nil, err
	}

	if c.Cache.Enabled() {
		c.Cache.StickResponseRules.SetAll(backend, transactionID, stickResRules)
	}
	return &models.GetStickResponseRulesOKBody{Version: v, Data: stickResRules}, nil
}

// GetStickResponseRule returns a struct with configuration version and a responseed stick response rule
// in the specified backend. Returns error on fail or if stick response rule does not exist.
func (c *LBCTLClient) GetStickResponseRule(id int64, backend string, transactionID string) (*models.GetStickResponseRuleOKBody, error) {
	if c.Cache.Enabled() {
		stickResRule, found := c.Cache.StickResponseRules.GetOne(id, backend, transactionID)
		if found {
			return &models.GetStickResponseRuleOKBody{Version: c.Cache.Version.Get(), Data: stickResRule}, nil
		}
	}
	stickResRuleStr, err := c.executeLBCTL("l7-farm-stickrsp-show", transactionID, backend, strconv.FormatInt(id, 10))
	if err != nil {
		return nil, err
	}
	stickResRule := &models.StickResponseRule{ID: id}

	c.parseObject(stickResRuleStr, stickResRule)

	v, err := c.GetVersion()
	if err != nil {
		return nil, err
	}

	return &models.GetStickResponseRuleOKBody{Version: v, Data: stickResRule}, nil
}

// DeleteStickResponseRule deletes a stick response rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *LBCTLClient) DeleteStickResponseRule(id int64, backend string, transactionID string, version int64) error {
	err := c.deleteObject(strconv.FormatInt(id, 10), "stickrsp", backend, "farm", transactionID, version)
	if err != nil {
		return err
	}
	if c.Cache.Enabled() {
		c.Cache.StickResponseRules.InvalidateBackend(transactionID, backend)
	}
	return nil
}

// CreateStickResponseRule creates a stick response rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *LBCTLClient) CreateStickResponseRule(backend string, data *models.StickResponseRule, transactionID string, version int64) error {
	if c.UseValidation() {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return NewConfError(ErrValidationError, validationErr.Error())
		}
	}
	err := c.createObject(strconv.FormatInt(data.ID, 10), "stickrsp", backend, "farm", data, nil, transactionID, version)
	if err != nil {
		return err
	}
	if c.Cache.Enabled() {
		c.Cache.StickResponseRules.InvalidateBackend(transactionID, backend)
	}
	return nil
}

// EditStickResponseRule edits a stick response rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *LBCTLClient) EditStickResponseRule(id int64, backend string, data *models.StickResponseRule, transactionID string, version int64) error {
	if c.UseValidation() {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return NewConfError(ErrValidationError, validationErr.Error())
		}
	}
	ondiskR, err := c.GetStickResponseRule(id, backend, transactionID)
	if err != nil {
		return err
	}

	err = c.editObject(strconv.FormatInt(data.ID, 10), "stickrsp", backend, "farm", data, ondiskR, nil, transactionID, version)
	if err != nil {
		return err
	}
	if c.Cache.Enabled() {
		c.Cache.StickResponseRules.InvalidateBackend(transactionID, backend)
	}
	return nil
}

func (c *LBCTLClient) parseStickResponseRules(response string) models.StickResponseRules {
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
