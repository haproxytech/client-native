package configuration

import (
	"fmt"

	strfmt "github.com/go-openapi/strfmt"

	parser "github.com/haproxytech/config-parser"
	parser_errors "github.com/haproxytech/config-parser/errors"
	"github.com/haproxytech/config-parser/types"
	"github.com/haproxytech/models"
)

// GetBackendSwitchingRules returns a struct with configuration version and an array of
// configured backend switching rules in the specified frontend. Returns error on fail.
func (c *Client) GetBackendSwitchingRules(frontend string, transactionID string) (*models.GetBackendSwitchingRulesOKBody, error) {
	if c.Cache.Enabled() {
		bckRules, found := c.Cache.BackendSwitchingRules.Get(frontend, transactionID)
		if found {
			return &models.GetBackendSwitchingRulesOKBody{Version: c.Cache.Version.Get(transactionID), Data: bckRules}, nil
		}
	}

	if err := c.ConfigParser.LoadData(c.getTransactionFile(transactionID)); err != nil {
		return nil, err
	}

	bckRules, err := c.parseBackendSwitchingRules(frontend)
	if err != nil {
		if err == parser_errors.SectionMissingErr {
			return nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Frontend %s does not exist", frontend))
		}
		return nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return nil, err
	}

	if c.Cache.Enabled() {
		c.Cache.BackendSwitchingRules.SetAll(frontend, transactionID, bckRules)
	}
	return &models.GetBackendSwitchingRulesOKBody{Version: v, Data: bckRules}, nil
}

// GetBackendSwitchingRule returns a struct with configuration version and a requested backend switching rule
// in the specified frontend. Returns error on fail or if backend switching rule does not exist.
func (c *Client) GetBackendSwitchingRule(id int64, frontend string, transactionID string) (*models.GetBackendSwitchingRuleOKBody, error) {
	if c.Cache.Enabled() {
		bckRule, found := c.Cache.BackendSwitchingRules.GetOne(id, frontend, transactionID)
		if found {
			return &models.GetBackendSwitchingRuleOKBody{Version: c.Cache.Version.Get(transactionID), Data: bckRule}, nil
		}
	}
	if err := c.ConfigParser.LoadData(c.getTransactionFile(transactionID)); err != nil {
		return nil, err
	}

	data, err := c.ConfigParser.GetOne(parser.Frontends, frontend, "use_backend", int(id))
	if err != nil {
		if err == parser_errors.SectionMissingErr {
			return nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Frontend %s does not exist", frontend))
		}
		if err == parser_errors.FetchError {
			return nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Backend Switching Rule %v does not exist in frontend %s", id, frontend))
		}
		return nil, err
	}

	bckRule := parseBackendSwitchingRule(data.(types.UseBackend))
	bckRule.ID = &id

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return nil, err
	}

	if c.Cache.Enabled() {
		c.Cache.BackendSwitchingRules.Set(id, frontend, transactionID, bckRule)
	}

	return &models.GetBackendSwitchingRuleOKBody{Version: v, Data: bckRule}, nil
}

// DeleteBackendSwitchingRule deletes a backend switching rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) DeleteBackendSwitchingRule(id int64, frontend string, transactionID string, version int64) error {
	t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	if err := c.ConfigParser.Delete(parser.Frontends, frontend, "use_backend", int(id)); err != nil {
		if err == parser_errors.SectionMissingErr {
			return NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Frontend %s does not exist", frontend))
		}
		if err == parser_errors.FetchError {
			return NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Backend Switching Rule %v does not exist in frontend %s", id, frontend))
		}
		return err
	}

	if err := c.saveData(t, transactionID); err != nil {
		return err
	}
	if c.Cache.Enabled() {
		c.Cache.BackendSwitchingRules.InvalidateFrontend(transactionID, frontend)
	}
	return nil
}

// CreateBackendSwitchingRule creates a backend switching rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) CreateBackendSwitchingRule(frontend string, data *models.BackendSwitchingRule, transactionID string, version int64) error {
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

	if err := c.ConfigParser.Insert(parser.Frontends, frontend, "use_backend", serializeBackendSwitchingRule(*data), int(*data.ID)); err != nil {
		if err == parser_errors.IndexOutOfRange {
			return c.errAndDeleteTransaction(NewConfError(ErrObjectIndexOutOfRange,
				fmt.Sprintf("Backend Switching Rule with id %v in frontend %s out of range", int(*data.ID), frontend)), t, transactionID == "")
		}
		if err == parser_errors.SectionMissingErr {
			return NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Frontend %s does not exist", frontend))
		}
		return c.errAndDeleteTransaction(err, t, transactionID == "")
	}

	if err := c.ConfigParser.Save(c.getTransactionFile(t)); err != nil {
		return c.errAndDeleteTransaction(err, t, transactionID == "")
	}

	if err := c.saveData(t, transactionID); err != nil {
		return err
	}

	if c.Cache.Enabled() {
		c.Cache.BackendSwitchingRules.InvalidateFrontend(transactionID, frontend)
	}
	return nil
}

// EditBackendSwitchingRule edits a backend switching rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) EditBackendSwitchingRule(id int64, frontend string, data *models.BackendSwitchingRule, transactionID string, version int64) error {
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

	if _, err := c.ConfigParser.GetOne(parser.Frontends, frontend, "use_backend", int(id)); err != nil {
		if err == parser_errors.SectionMissingErr {
			return NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Frontend %s does not exist", frontend))
		}
		if err == parser_errors.FetchError {
			return NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Backend Switching Rule %v does not exist in frontend %s", id, frontend))
		}
		return err
	}

	if err := c.ConfigParser.Set(parser.Frontends, frontend, "use_backend", serializeBackendSwitchingRule(*data), int(id)); err != nil {
		if err == parser_errors.IndexOutOfRange {
			return c.errAndDeleteTransaction(NewConfError(ErrObjectIndexOutOfRange,
				fmt.Sprintf("Backend Switching Rule with id %v in frontend %s out of range", int(*data.ID), frontend)), t, transactionID == "")
		}
		return c.errAndDeleteTransaction(err, t, transactionID == "")
	}

	if err := c.saveData(t, transactionID); err != nil {
		return err
	}

	if c.Cache.Enabled() {
		c.Cache.BackendSwitchingRules.InvalidateFrontend(transactionID, frontend)
	}
	return nil
}

func (c *Client) parseBackendSwitchingRules(frontend string) (models.BackendSwitchingRules, error) {
	br := models.BackendSwitchingRules{}

	data, err := c.ConfigParser.Get(parser.Frontends, frontend, "use_backend", false)
	if err != nil {
		if err == parser_errors.FetchError {
			return br, nil
		}
		return nil, err
	}

	bRules := data.([]types.UseBackend)
	for i, bRule := range bRules {
		id := int64(i)
		b := parseBackendSwitchingRule(bRule)
		if b != nil {
			b.ID = &id
			br = append(br, b)
		}
	}
	return br, nil
}

func parseBackendSwitchingRule(ub types.UseBackend) *models.BackendSwitchingRule {
	return &models.BackendSwitchingRule{
		Name:     ub.Name,
		Cond:     ub.Cond,
		CondTest: ub.CondTest,
	}
}

func serializeBackendSwitchingRule(bRule models.BackendSwitchingRule) types.UseBackend {
	return types.UseBackend{
		Name:     bRule.Name,
		Cond:     bRule.Cond,
		CondTest: bRule.CondTest,
	}
}
