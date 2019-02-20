package configuration

import (
	"fmt"

	"github.com/haproxytech/config-parser/parsers/filters"

	strfmt "github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/config-parser"
	parser_errors "github.com/haproxytech/config-parser/errors"
	"github.com/haproxytech/config-parser/types"
	"github.com/haproxytech/models"
)

// GetFilters returns a struct with configuration version and an array of
// configured filters in the specified parent. Returns error on fail.
func (c *Client) GetFilters(parentType, parentName string, transactionID string) (*models.GetFiltersOKBody, error) {
	if c.Cache.Enabled() {
		filters, found := c.Cache.Filters.Get(parentName, parentType, transactionID)
		if found {
			return &models.GetFiltersOKBody{Version: c.Cache.Version.Get(transactionID), Data: filters}, nil
		}
	}

	if err := c.ConfigParser.LoadData(c.getTransactionFile(transactionID)); err != nil {
		return nil, err
	}

	filters, err := c.parseFilters(parentType, parentName)
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
		c.Cache.Filters.SetAll(parentName, parentType, transactionID, filters)
	}
	return &models.GetFiltersOKBody{Version: v, Data: filters}, nil
}

// GetFilter returns a struct with configuration version and a requested filter
// in the specified parent. Returns error on fail or if filter does not exist.
func (c *Client) GetFilter(id int64, parentType, parentName string, transactionID string) (*models.GetFilterOKBody, error) {
	if c.Cache.Enabled() {
		filter, found := c.Cache.Filters.GetOne(id, parentName, parentType, transactionID)
		if found {
			return &models.GetFilterOKBody{Version: c.Cache.Version.Get(transactionID), Data: filter}, nil
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

	data, err := c.ConfigParser.GetOne(p, parentName, "filter", int(id))
	if err != nil {
		if err == parser_errors.SectionMissingErr {
			return nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("%s %s does not exist", parentType, parentName))
		}
		if err == parser_errors.FetchError {
			return nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Filter %v does not exist in %s %s", id, parentType, parentName))
		}
		return nil, err
	}

	filter := parseFilter(data.(types.Filter))
	filter.ID = &id

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return nil, err
	}

	if c.Cache.Enabled() {
		c.Cache.Filters.Set(id, parentName, parentType, transactionID, filter)
	}

	return &models.GetFilterOKBody{Version: v, Data: filter}, nil
}

// DeleteFilter deletes a filter in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) DeleteFilter(id int64, parentType string, parentName string, transactionID string, version int64) error {
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

	if err := c.ConfigParser.Delete(p, parentName, "filter", int(id)); err != nil {
		if err == parser_errors.SectionMissingErr {
			return NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("%s %s does not exist", parentName, parentType))
		}
		if err == parser_errors.FetchError {
			return NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Filter %v does not exist in %s %s", id, parentName, parentType))
		}
		return err
	}

	if err := c.saveData(t, transactionID); err != nil {
		return err
	}

	if c.Cache.Enabled() {
		c.Cache.Filters.InvalidateParent(transactionID, parentName, parentType)
	}
	return nil
}

// CreateFilter creates a filter in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) CreateFilter(parentType string, parentName string, data *models.Filter, transactionID string, version int64) error {
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

	if err := c.ConfigParser.Insert(p, parentName, "filter", serializeFilter(*data), int(*data.ID)); err != nil {
		if err == parser_errors.IndexOutOfRange {
			return c.errAndDeleteTransaction(NewConfError(ErrObjectIndexOutOfRange,
				fmt.Sprintf("Filter with id %v in %s %s out of range", int(*data.ID), parentType, parentName)), t, transactionID == "")
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
		c.Cache.Filters.InvalidateParent(transactionID, parentName, parentType)
	}
	return nil
}

// EditFilter edits a filter in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) EditFilter(id int64, parentType string, parentName string, data *models.Filter, transactionID string, version int64) error {
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

	if _, err := c.ConfigParser.GetOne(p, parentName, "filter", int(id)); err != nil {
		if err == parser_errors.SectionMissingErr {
			return NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("%s %s does not exist", parentType, parentName))
		}
		if err == parser_errors.FetchError {
			return NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Filter %v does not exist in %s %s", id, parentType, parentName))
		}
		return err
	}

	if err := c.ConfigParser.Set(p, parentName, "filter", serializeFilter(*data), int(id)); err != nil {
		if err == parser_errors.IndexOutOfRange {
			return c.errAndDeleteTransaction(NewConfError(ErrObjectIndexOutOfRange,
				fmt.Sprintf("Filter with id %v in %s %s out of range", int(*data.ID), parentType, parentName)), t, transactionID == "")
		}
		return c.errAndDeleteTransaction(err, t, transactionID == "")
	}

	if err := c.saveData(t, transactionID); err != nil {
		return err
	}

	if c.Cache.Enabled() {
		c.Cache.Filters.InvalidateParent(transactionID, parentName, parentType)
	}
	return nil
}

func (c *Client) parseFilters(t, pName string) (models.Filters, error) {
	p := parser.Global
	if t == "frontend" {
		p = parser.Frontends
	} else if t == "backend" {
		p = parser.Backends
	}

	f := models.Filters{}
	data, err := c.ConfigParser.Get(p, pName, "filter", false)
	if err != nil {
		if err == parser_errors.FetchError {
			return f, nil
		}
		return nil, err
	}

	filters := data.([]types.Filter)
	for i, filter := range filters {
		id := int64(i)
		mFilter := parseFilter(filter)
		if mFilter != nil {
			mFilter.ID = &id
			f = append(f, mFilter)
		}
	}
	return f, nil
}

func parseFilter(f types.Filter) *models.Filter {
	switch v := f.(type) {
	case *filters.Trace:
		return &models.Filter{
			Type:               "trace",
			TraceName:          v.Name,
			TraceHexdump:       v.Hexdump,
			TraceRndForwarding: v.RandomForwarding,
			TraceRndParsing:    v.RandomParsing,
		}
	case *filters.Compression:
		return &models.Filter{
			Type: "compression",
		}
	case *filters.Spoe:
		return &models.Filter{
			Type:       "spoe",
			SpoeConfig: v.Config,
			SpoeEngine: v.Engine,
		}
	case *filters.Cache:
		return &models.Filter{
			Type:      "cache",
			CacheName: v.Name,
		}
	}
	return nil
}

func serializeFilter(f models.Filter) types.Filter {
	switch f.Type {
	case "trace":
		return &filters.Trace{
			Name:             f.TraceName,
			Hexdump:          f.TraceHexdump,
			RandomForwarding: f.TraceRndForwarding,
			RandomParsing:    f.TraceRndParsing,
		}
	case "compression":
		return &filters.Compression{
			Enabled: true,
		}
	case "spoe":
		return &filters.Spoe{
			Config: f.SpoeConfig,
			Engine: f.SpoeEngine,
		}
	case "cache":
		return &filters.Cache{
			Name: f.CacheName,
		}
	}
	return nil
}
