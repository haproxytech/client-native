package configuration

import (
	"fmt"
	"strconv"
	"strings"

	strfmt "github.com/go-openapi/strfmt"
	"github.com/haproxytech/models"
)

// GetTCPContentRules returns a struct with configuration version and an array of
// configured tcp content rules in the specified parent. Returns error on fail.
func (c *LBCTLClient) GetTCPContentRules(parentType, parentName, ruleType, transactionID string) (*models.GetTCPContentRulesOKBody, error) {
	lbctlType := typeToLbctlType(parentType)
	if lbctlType == "" {
		return nil, NewConfError(ErrValidationError, fmt.Sprintf("Parent type %v not recognized", parentType))
	}

	var lbctlRType string
	switch ruleType {
	case "request":
		lbctlRType = "tcpreqcont"
	case "response":
		lbctlRType = "tcprspcont"
		if parentType == "frontend" {
			return nil, NewConfError(ErrValidationError, "Rule type cannot be response for frontend parent")
		}
	}
	if lbctlRType == "" {
		return nil, NewConfError(ErrValidationError, fmt.Sprintf("Rule type %v not recognized", ruleType))
	}

	tcpRulesStr, err := c.executeLBCTL("l7-"+lbctlType+"-"+lbctlRType+"-dump", "", parentName)
	if err != nil {
		return nil, err
	}

	tcpRules := c.parseTCPContentRules(tcpRulesStr)

	v, err := c.GetVersion()
	if err != nil {
		return nil, err
	}

	return &models.GetTCPContentRulesOKBody{Version: v, Data: tcpRules}, nil
}

// GetTCPContentRule returns a struct with configuration version and a requested tcp content rule
// in the specified parent. Returns error on fail or if tcp content rule does not exist.
func (c *LBCTLClient) GetTCPContentRule(id int64, parentType, parentName, ruleType, transactionID string) (*models.GetTCPContentRuleOKBody, error) {
	lbctlType := typeToLbctlType(parentType)
	if lbctlType == "" {
		return nil, NewConfError(ErrValidationError, fmt.Sprintf("Parent type %v not recognized", parentType))
	}

	var lbctlRType string
	switch ruleType {
	case "request":
		lbctlRType = "tcpreqcont"
	case "response":
		lbctlRType = "tcprspcont"
		if parentType == "frontend" {
			return nil, NewConfError(ErrValidationError, "Rule type cannot be response for frontend parent")
		}
	}
	if lbctlRType == "" {
		return nil, NewConfError(ErrValidationError, fmt.Sprintf("Rule type %v not recognized", ruleType))
	}

	tcpRuleStr, err := c.executeLBCTL("l7-"+lbctlType+"-"+lbctlRType+"-show", "", parentName, strconv.FormatInt(id, 10))
	if err != nil {
		return nil, err
	}
	tcpRule := &models.TCPRule{ID: id}

	c.parseObject(tcpRuleStr, tcpRule)

	v, err := c.GetVersion()
	if err != nil {
		return nil, err
	}

	return &models.GetTCPContentRuleOKBody{Version: v, Data: tcpRule}, nil
}

// DeleteTCPContentRule deletes a tcp content rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *LBCTLClient) DeleteTCPContentRule(id int64, parentType, parentName, ruleType string, transactionID string, version int64) error {
	lbctlType := typeToLbctlType(parentType)
	if lbctlType == "" {
		return NewConfError(ErrValidationError, fmt.Sprintf("Parent type %v not recognized", parentType))
	}

	var lbctlRType string
	switch ruleType {
	case "request":
		lbctlRType = "tcpreqcont"
	case "response":
		lbctlRType = "tcprspcont"
		if parentType == "frontend" {
			return NewConfError(ErrValidationError, "Rule type cannot be response for frontend parent")
		}
	}
	if lbctlRType == "" {
		return NewConfError(ErrValidationError, fmt.Sprintf("Rule type %v not recognized", ruleType))
	}

	return c.deleteObject(strconv.FormatInt(id, 10), lbctlRType, parentName, lbctlType, transactionID, version)
}

// CreateTCPContentRule creates a tcp content rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *LBCTLClient) CreateTCPContentRule(parentType, parentName, ruleType string, data *models.TCPRule, transactionID string, version int64) error {
	if c.UseValidation() {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return NewConfError(ErrValidationError, validationErr.Error())
		}
	}
	var lbctlRType string
	switch ruleType {
	case "request":
		lbctlRType = "tcpreqcont"
	case "response":
		lbctlRType = "tcprspcont"
		if parentType == "frontend" {
			return NewConfError(ErrValidationError, "Rule type cannot be response for frontend parent")
		}
	}
	if lbctlRType == "" {
		return NewConfError(ErrValidationError, fmt.Sprintf("Rule type %v not recognized", ruleType))
	}

	lbctlType := typeToLbctlType(parentType)
	if lbctlType == "" {
		return NewConfError(ErrValidationError, fmt.Sprintf("Parent type %v not recognized", parentType))
	}

	return c.createObject(strconv.FormatInt(data.ID, 10), lbctlRType, parentName, lbctlType, data, nil, transactionID, version)
}

// EditTCPContentRule edits a tcp content rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *LBCTLClient) EditTCPContentRule(id int64, parentType, parentName, ruleType string, data *models.TCPRule, transactionID string, version int64) error {
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

	var lbctlRType string
	switch ruleType {
	case "request":
		lbctlRType = "tcpreqcont"
	case "response":
		lbctlRType = "tcprspcont"
		if parentType == "frontend" {
			return NewConfError(ErrValidationError, "Rule type cannot be response for frontend parent")
		}

	}
	if lbctlRType == "" {
		return NewConfError(ErrValidationError, fmt.Sprintf("Rule type %v not recognized", ruleType))
	}

	ondiskBr, err := c.GetTCPContentRule(id, parentType, parentName, ruleType, transactionID)
	if err != nil {
		return err
	}

	return c.editObject(strconv.FormatInt(data.ID, 10), lbctlRType, parentName, lbctlType, data, ondiskBr, nil, transactionID, version)
}

func (c *LBCTLClient) parseTCPContentRules(response string) models.TCPRules {
	tcpRules := make(models.TCPRules, 0, 1)
	for _, tcpRulesStr := range strings.Split(response, "\n\n") {
		if strings.TrimSpace(tcpRulesStr) == "" {
			continue
		}
		idStr, _ := splitHeaderLine(tcpRulesStr)
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			id = 0
		}

		tcpRulesObj := &models.TCPRule{ID: id}
		c.parseObject(tcpRulesStr, tcpRulesObj)
		tcpRules = append(tcpRules, tcpRulesObj)
	}
	return tcpRules
}
