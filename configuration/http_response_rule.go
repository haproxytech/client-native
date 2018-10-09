package configuration

import (
	"fmt"
	"strconv"
	"strings"

	strfmt "github.com/go-openapi/strfmt"
	"github.com/haproxytech/models"
)

// GetHTTPResponseRules returns a struct with configuration version and an array of
// configured http response rules in the specified parent. Returns error on fail.
func (c *LBCTLConfigurationClient) GetHTTPResponseRules(parentType, parentName string) (*models.GetHTTPResponseRulesOKBody, error) {
	lbctlType := typeToLbctlType(parentType)
	if lbctlType == "" {
		return nil, NewConfError(ErrValidationError, fmt.Sprintf("Parent type %v not recognized", parentType))
	}

	httpRulesStr, err := c.executeLBCTL("l7-"+lbctlType+"-httprsp-dump", "", parentName)
	if err != nil {
		return nil, err
	}

	httpRules := c.parseHTTPResponseRules(httpRulesStr)

	v, err := c.GetVersion()
	if err != nil {
		return nil, err
	}

	return &models.GetHTTPResponseRulesOKBody{Version: v, Data: httpRules}, nil
}

// GetHTTPResponseRule returns a struct with configuration version and a responseed http response rule
// in the specified parent. Returns error on fail or if http response rule does not exist.
func (c *LBCTLConfigurationClient) GetHTTPResponseRule(id int64, parentType, parentName string) (*models.GetHTTPResponseRuleOKBody, error) {
	lbctlType := typeToLbctlType(parentType)
	if lbctlType == "" {
		return nil, NewConfError(ErrValidationError, fmt.Sprintf("Parent type %v not recognized", parentType))
	}

	httpRuleStr, err := c.executeLBCTL("l7-"+lbctlType+"-httprsp-show", "", parentName, strconv.FormatInt(id, 10))
	if err != nil {
		return nil, err
	}
	httpRule := &models.HTTPResponseRule{ID: id}

	c.parseObject(httpRuleStr, httpRule)

	v, err := c.GetVersion()
	if err != nil {
		return nil, err
	}

	return &models.GetHTTPResponseRuleOKBody{Version: v, Data: httpRule}, nil
}

// DeleteHTTPResponseRule deletes a http response rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *LBCTLConfigurationClient) DeleteHTTPResponseRule(id int64, parentType string, parentName string, transactionID string, version int64) error {
	lbctlType := typeToLbctlType(parentType)
	if lbctlType == "" {
		return NewConfError(ErrValidationError, fmt.Sprintf("Parent type %v not recognized", parentType))
	}

	return c.deleteObject(strconv.FormatInt(id, 10), "httprsp", parentName, lbctlType, transactionID, version)
}

// CreateHTTPResponseRule creates a http response rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *LBCTLConfigurationClient) CreateHTTPResponseRule(parentType string, parentName string, data *models.HTTPResponseRule, transactionID string, version int64) error {
	validationErr := data.Validate(strfmt.Default)
	if validationErr != nil {
		return NewConfError(ErrValidationError, validationErr.Error())
	}

	lbctlType := typeToLbctlType(parentType)
	if lbctlType == "" {
		return NewConfError(ErrValidationError, fmt.Sprintf("Parent type %v not recognized", parentType))
	}

	return c.createObject(strconv.FormatInt(data.ID, 10), "httprsp", parentName, lbctlType, data, nil, transactionID, version)
}

// EditHTTPResponseRule edits a http response rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *LBCTLConfigurationClient) EditHTTPResponseRule(id int64, parentType string, parentName string, data *models.HTTPResponseRule, transactionID string, version int64) error {
	validationErr := data.Validate(strfmt.Default)
	if validationErr != nil {
		return NewConfError(ErrValidationError, validationErr.Error())
	}

	lbctlType := typeToLbctlType(parentType)
	if lbctlType == "" {
		return NewConfError(ErrValidationError, fmt.Sprintf("Parent type %v not recognized", parentType))
	}

	ondiskR, err := c.GetHTTPResponseRule(id, parentType, parentName)
	if err != nil {
		return err
	}

	return c.editObject(strconv.FormatInt(data.ID, 10), "httprsp", parentName, lbctlType, data, ondiskR, nil, transactionID, version)
}

func (c *LBCTLConfigurationClient) parseHTTPResponseRules(response string) models.HTTPResponseRules {
	httpRules := make(models.HTTPResponseRules, 0, 1)
	for _, httpRulesStr := range strings.Split(response, "\n\n") {
		if strings.TrimSpace(httpRulesStr) == "" {
			continue
		}
		idStr, _ := splitHeaderLine(httpRulesStr)
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			id = 0
		}

		httpRulesObj := &models.HTTPResponseRule{ID: id}
		c.parseObject(httpRulesStr, httpRulesObj)
		httpRules = append(httpRules, httpRulesObj)
	}
	return httpRules
}
