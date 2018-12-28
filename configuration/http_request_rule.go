package configuration

import (
	"fmt"
	"strconv"
	"strings"

	strfmt "github.com/go-openapi/strfmt"
	"github.com/haproxytech/models"
)

// GetHTTPRequestRules returns a struct with configuration version and an array of
// configured http request rules in the specified parent. Returns error on fail.
func (c *LBCTLClient) GetHTTPRequestRules(parentType, parentName string, transactionID string) (*models.GetHTTPRequestRulesOKBody, error) {
	lbctlType := typeToLbctlType(parentType)
	if lbctlType == "" {
		return nil, NewConfError(ErrValidationError, fmt.Sprintf("Parent type %v not recognized", parentType))
	}

	httpRulesStr, err := c.executeLBCTL("l7-"+lbctlType+"-httpreq-dump", transactionID, parentName)
	if err != nil {
		return nil, err
	}

	httpRules := c.parseHTTPRequestRules(httpRulesStr)

	v, err := c.GetVersion()
	if err != nil {
		return nil, err
	}

	return &models.GetHTTPRequestRulesOKBody{Version: v, Data: httpRules}, nil
}

// GetHTTPRequestRule returns a struct with configuration version and a requested http request rule
// in the specified parent. Returns error on fail or if http request rule does not exist.
func (c *LBCTLClient) GetHTTPRequestRule(id int64, parentType, parentName string, transactionID string) (*models.GetHTTPRequestRuleOKBody, error) {
	lbctlType := typeToLbctlType(parentType)
	if lbctlType == "" {
		return nil, NewConfError(ErrValidationError, fmt.Sprintf("Parent type %v not recognized", parentType))
	}

	httpRuleStr, err := c.executeLBCTL("l7-"+lbctlType+"-httpreq-show", transactionID, parentName, strconv.FormatInt(id, 10))
	if err != nil {
		return nil, err
	}
	httpRule := &models.HTTPRequestRule{ID: id}

	c.parseObject(httpRuleStr, httpRule)

	v, err := c.GetVersion()
	if err != nil {
		return nil, err
	}

	return &models.GetHTTPRequestRuleOKBody{Version: v, Data: httpRule}, nil
}

// DeleteHTTPRequestRule deletes a http request rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *LBCTLClient) DeleteHTTPRequestRule(id int64, parentType string, parentName string, transactionID string, version int64) error {
	lbctlType := typeToLbctlType(parentType)
	if lbctlType == "" {
		return NewConfError(ErrValidationError, fmt.Sprintf("Parent type %v not recognized", parentType))
	}

	return c.deleteObject(strconv.FormatInt(id, 10), "httpreq", parentName, lbctlType, transactionID, version)
}

// CreateHTTPRequestRule creates a http request rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *LBCTLClient) CreateHTTPRequestRule(parentType string, parentName string, data *models.HTTPRequestRule, transactionID string, version int64) error {
	if c.UseValidation() {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return NewConfError(ErrValidationError, validationErr.Error())
		}
	}

	lbctlType := typeToLbctlType(parentType)
	if lbctlType == "" {
		return NewConfError(ErrValidationError, fmt.Sprintf("Parent type %v not recognized", parentType))
	}

	return c.createObject(strconv.FormatInt(data.ID, 10), "httpreq", parentName, lbctlType, data, nil, transactionID, version)
}

// EditHTTPRequestRule edits a http request rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *LBCTLClient) EditHTTPRequestRule(id int64, parentType string, parentName string, data *models.HTTPRequestRule, transactionID string, version int64) error {
	if c.UseValidation() {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return NewConfError(ErrValidationError, validationErr.Error())
		}
	}

	lbctlType := typeToLbctlType(parentType)
	if lbctlType == "" {
		return NewConfError(ErrValidationError, fmt.Sprintf("Parent type %v not recognized", parentType))
	}

	ondiskR, err := c.GetHTTPRequestRule(id, parentType, parentName, transactionID)
	if err != nil {
		return err
	}

	return c.editObject(strconv.FormatInt(data.ID, 10), "httpreq", parentName, lbctlType, data, ondiskR, nil, transactionID, version)
}

func (c *LBCTLClient) parseHTTPRequestRules(response string) models.HTTPRequestRules {
	httpRules := make(models.HTTPRequestRules, 0, 1)
	for _, httpRulesStr := range strings.Split(response, "\n\n") {
		if strings.TrimSpace(httpRulesStr) == "" {
			continue
		}
		idStr, _ := splitHeaderLine(httpRulesStr)
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			id = 0
		}

		httpRulesObj := &models.HTTPRequestRule{ID: id}
		c.parseObject(httpRulesStr, httpRulesObj)
		httpRules = append(httpRules, httpRulesObj)
	}
	return httpRules
}
