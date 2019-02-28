package configuration

import (
	"strconv"

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
	p, err := c.GetParser(transactionID)
	if err != nil {
		return nil, err
	}

	filters, err := c.parseFilters(parentType, parentName, p)
	if err != nil {
		return nil, c.handleError("", parentType, parentName, "", false, err)
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return nil, err
	}

	return &models.GetFiltersOKBody{Version: v, Data: filters}, nil
}

// GetFilter returns a struct with configuration version and a requested filter
// in the specified parent. Returns error on fail or if filter does not exist.
func (c *Client) GetFilter(id int64, parentType, parentName string, transactionID string) (*models.GetFilterOKBody, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return nil, err
	}

	var section parser.Section
	if parentType == "backend" {
		section = parser.Backends
	} else if parentType == "frontend" {
		section = parser.Frontends
	}

	data, err := p.GetOne(section, parentName, "filter", int(id))
	if err != nil {
		return nil, c.handleError(strconv.FormatInt(id, 10), parentType, parentName, "", false, err)
	}

	filter := parseFilter(data.(types.Filter))
	filter.ID = &id

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return nil, err
	}

	return &models.GetFilterOKBody{Version: v, Data: filter}, nil
}

// DeleteFilter deletes a filter in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) DeleteFilter(id int64, parentType string, parentName string, transactionID string, version int64) error {
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	var section parser.Section
	if parentType == "backend" {
		section = parser.Backends
	} else if parentType == "frontend" {
		section = parser.Frontends
	}

	if err := p.Delete(section, parentName, "filter", int(id)); err != nil {
		return c.handleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	if err := c.saveData(p, t, transactionID == ""); err != nil {
		return err
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

	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	var section parser.Section
	if parentType == "backend" {
		section = parser.Backends
	} else if parentType == "frontend" {
		section = parser.Frontends
	}

	if err := p.Insert(section, parentName, "filter", serializeFilter(*data), int(*data.ID)); err != nil {
		return c.handleError(strconv.FormatInt(*data.ID, 10), parentType, parentName, t, transactionID == "", err)
	}

	if err := c.saveData(p, t, transactionID == ""); err != nil {
		return err
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
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	var section parser.Section
	if parentType == "backend" {
		section = parser.Backends
	} else if parentType == "frontend" {
		section = parser.Frontends
	}

	if _, err := p.GetOne(section, parentName, "filter", int(id)); err != nil {
		return c.handleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	if err := p.Set(section, parentName, "filter", serializeFilter(*data), int(id)); err != nil {
		return c.handleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	if err := c.saveData(p, t, transactionID == ""); err != nil {
		return err
	}

	return nil
}

func (c *Client) parseFilters(t, pName string, p *parser.Parser) (models.Filters, error) {
	section := parser.Global
	if t == "frontend" {
		section = parser.Frontends
	} else if t == "backend" {
		section = parser.Backends
	}

	f := models.Filters{}
	data, err := p.Get(section, pName, "filter", false)
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
