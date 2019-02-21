package configuration

import (
	"fmt"

	strfmt "github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/config-parser"
	parser_errors "github.com/haproxytech/config-parser/errors"
	"github.com/haproxytech/config-parser/types"
	"github.com/haproxytech/models"
)

// GetLogTargets returns a struct with configuration version and an array of
// configured log targets in the specified parent. Returns error on fail.
func (c *Client) GetLogTargets(parentType, parentName string, transactionID string) (*models.GetLogTargetsOKBody, error) {
	if c.Cache.Enabled() {
		logTargets, found := c.Cache.LogTargets.Get(parentName, parentType, transactionID)
		if found {
			return &models.GetLogTargetsOKBody{Version: c.Cache.Version.Get(transactionID), Data: logTargets}, nil
		}
	}
	if err := c.ConfigParser.LoadData(c.getTransactionFile(transactionID)); err != nil {
		return nil, err
	}

	logTargets, err := c.parseLogTargets(parentType, parentName)
	if err != nil {
		if err == parser_errors.SectionMissingErr {
			return nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("%s %s does not exist", parentType, parentName))
		}
		return nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return nil, err
	}

	if c.Cache.Enabled() {
		c.Cache.LogTargets.SetAll(parentName, parentType, transactionID, logTargets)
	}
	return &models.GetLogTargetsOKBody{Version: v, Data: logTargets}, nil
}

// GetLogTarget returns a struct with configuration version and a requested log target
// in the specified parent. Returns error on fail or if log target does not exist.
func (c *Client) GetLogTarget(id int64, parentType, parentName string, transactionID string) (*models.GetLogTargetOKBody, error) {
	if c.Cache.Enabled() {
		logTarget, found := c.Cache.LogTargets.GetOne(id, parentName, parentType, transactionID)
		if found {
			return &models.GetLogTargetOKBody{Version: c.Cache.Version.Get(transactionID), Data: logTarget}, nil
		}
	}
	if err := c.ConfigParser.LoadData(c.getTransactionFile(transactionID)); err != nil {
		return nil, err
	}

	var p parser.Section
	if parentType == "backend" {
		p = parser.Backends
	} else if parentType == "frontend" {
		p = parser.Frontends
	}

	data, err := c.ConfigParser.GetOne(p, parentName, "log", int(id))
	if err != nil {
		if err == parser_errors.SectionMissingErr {
			return nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("%s %s does not exist", parentType, parentName))
		}
		if err == parser_errors.FetchError {
			return nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Log Target %v does not exist in %s %s", id, parentType, parentName))
		}
		return nil, err
	}

	logTarget := parseLogTarget(data.(types.Log))
	logTarget.ID = &id

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return nil, err
	}

	if c.Cache.Enabled() {
		c.Cache.LogTargets.Set(id, parentName, parentType, transactionID, logTarget)
	}

	return &models.GetLogTargetOKBody{Version: v, Data: logTarget}, nil
}

// DeleteLogTarget deletes a log target in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) DeleteLogTarget(id int64, parentType string, parentName string, transactionID string, version int64) error {
	t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	var p parser.Section
	if parentType == "backend" {
		p = parser.Backends
	} else if parentType == "frontend" {
		p = parser.Frontends
	}

	if err := c.ConfigParser.Delete(p, parentName, "log", int(id)); err != nil {
		if err == parser_errors.SectionMissingErr {
			return NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("%s %s does not exist", parentName, parentType))
		}
		if err == parser_errors.FetchError {
			return NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Log Target %v does not exist in %s %s", id, parentName, parentType))
		}
		return err
	}

	if err := c.saveData(t, transactionID); err != nil {
		return err
	}

	if c.Cache.Enabled() {
		c.Cache.LogTargets.InvalidateParent(transactionID, parentName, parentType)
	}
	return nil
}

// CreateLogTarget creates a log target in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) CreateLogTarget(parentType string, parentName string, data *models.LogTarget, transactionID string, version int64) error {
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

	var p parser.Section
	if parentType == "backend" {
		p = parser.Backends
	} else if parentType == "frontend" {
		p = parser.Frontends
	}

	if err := c.ConfigParser.Insert(p, parentName, "log", serializeLogTarget(*data), int(*data.ID)); err != nil {
		if err == parser_errors.IndexOutOfRange {
			return c.errAndDeleteTransaction(NewConfError(ErrObjectIndexOutOfRange,
				fmt.Sprintf("Log Target with id %v in %s %s out of range", int(*data.ID), parentType, parentName)), t, transactionID == "")
		}
		if err == parser_errors.SectionMissingErr {
			return NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("%s %s does not exist", parentType, parentName))
		}
		return c.errAndDeleteTransaction(err, t, transactionID == "")
	}

	if err := c.saveData(t, transactionID); err != nil {
		return err
	}
	if c.Cache.Enabled() {
		c.Cache.LogTargets.InvalidateParent(transactionID, parentName, parentType)
	}
	return nil
}

// EditLogTarget edits a log target in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) EditLogTarget(id int64, parentType string, parentName string, data *models.LogTarget, transactionID string, version int64) error {
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

	var p parser.Section
	if parentType == "backend" {
		p = parser.Backends
	} else if parentType == "frontend" {
		p = parser.Frontends
	}

	if _, err := c.ConfigParser.GetOne(p, parentName, "log", int(id)); err != nil {
		if err == parser_errors.SectionMissingErr {
			return NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("%s %s does not exist", parentType, parentName))
		}
		if err == parser_errors.FetchError {
			return NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Log Target %v does not exist in %s %s", id, parentType, parentName))
		}
		return err
	}

	if err := c.ConfigParser.Set(p, parentName, "log", serializeLogTarget(*data), int(id)); err != nil {
		if err == parser_errors.IndexOutOfRange {
			return c.errAndDeleteTransaction(NewConfError(ErrObjectIndexOutOfRange,
				fmt.Sprintf("Log Target with id %v in %s %s out of range", int(*data.ID), parentType, parentName)), t, transactionID == "")
		}
		return c.errAndDeleteTransaction(err, t, transactionID == "")
	}

	if err := c.saveData(t, transactionID); err != nil {
		return err
	}
	if c.Cache.Enabled() {
		c.Cache.LogTargets.InvalidateParent(transactionID, parentName, parentType)
	}
	return nil
}

func (c *Client) parseLogTargets(t, pName string) (models.LogTargets, error) {
	p := parser.Global
	if t == "frontend" {
		p = parser.Frontends
	} else if t == "backend" {
		p = parser.Backends
	}

	logTargets := models.LogTargets{}
	data, err := c.ConfigParser.Get(p, pName, "log", false)
	if err != nil {
		if err == parser_errors.FetchError {
			return logTargets, nil
		}
		return nil, err
	}

	targets := data.([]types.Log)
	for i, l := range targets {
		id := int64(i)
		logTarget := parseLogTarget(l)
		if logTarget != nil {
			logTarget.ID = &id
			logTargets = append(logTargets, logTarget)
		}
	}
	return logTargets, nil
}

func parseLogTarget(l types.Log) *models.LogTarget {
	return &models.LogTarget{
		Address:  l.Address,
		Facility: l.Facility,
		Format:   l.Format,
		Global:   l.Global,
		Length:   l.Length,
		Level:    l.Level,
		Minlevel: l.MinLevel,
		Nolog:    l.NoLog,
	}
}

func serializeLogTarget(l models.LogTarget) types.Log {
	return types.Log{
		Address:  l.Address,
		Facility: l.Facility,
		Format:   l.Format,
		Global:   l.Global,
		Length:   l.Length,
		Level:    l.Level,
		MinLevel: l.Minlevel,
		NoLog:    l.Nolog,
	}
}
