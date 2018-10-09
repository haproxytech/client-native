package configuration

import (
	"fmt"
	"strconv"
	"strings"

	strfmt "github.com/go-openapi/strfmt"
	"github.com/haproxytech/models"
)

// GetFilters returns a struct with configuration version and an array of
// configured filters in the specified parent. Returns error on fail.
func (c *LBCTLConfigurationClient) GetFilters(parentType, parentName string) (*models.GetFiltersOKBody, error) {
	lbctlType := typeToLbctlType(parentType)
	if lbctlType == "" {
		return nil, NewConfError(ErrValidationError, fmt.Sprintf("Parent type %v not recognized", parentType))
	}

	filtersStr, err := c.executeLBCTL("l7-"+lbctlType+"-filter-dump", "", parentName)
	if err != nil {
		return nil, err
	}

	filters := c.parseFilters(filtersStr)

	v, err := c.GetVersion()
	if err != nil {
		return nil, err
	}

	return &models.GetFiltersOKBody{Version: v, Data: filters}, nil
}

// GetFilter returns a struct with configuration version and a requested filter
// in the specified parent. Returns error on fail or if filter does not exist.
func (c *LBCTLConfigurationClient) GetFilter(id int64, parentType, parentName string) (*models.GetFilterOKBody, error) {
	lbctlType := typeToLbctlType(parentType)
	if lbctlType == "" {
		return nil, NewConfError(ErrValidationError, fmt.Sprintf("Parent type %v not recognized", parentType))
	}

	filterStr, err := c.executeLBCTL("l7-"+lbctlType+"-filter-show", "", parentName, strconv.FormatInt(id, 10))
	if err != nil {
		return nil, err
	}
	filter := &models.Filter{ID: id}

	c.parseObject(filterStr, filter)

	v, err := c.GetVersion()
	if err != nil {
		return nil, err
	}

	return &models.GetFilterOKBody{Version: v, Data: filter}, nil
}

// DeleteFilter deletes a filter in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *LBCTLConfigurationClient) DeleteFilter(id int64, parentType string, parentName string, transactionID string, version int64) error {
	lbctlType := typeToLbctlType(parentType)
	if lbctlType == "" {
		return NewConfError(ErrValidationError, fmt.Sprintf("Parent type %v not recognized", parentType))
	}

	return c.deleteObject(strconv.FormatInt(id, 10), "filter", parentName, lbctlType, transactionID, version)
}

// CreateFilter creates a filter in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *LBCTLConfigurationClient) CreateFilter(parentType string, parentName string, data *models.Filter, transactionID string, version int64) error {
	validationErr := data.Validate(strfmt.Default)
	if validationErr != nil {
		return NewConfError(ErrValidationError, validationErr.Error())
	}

	lbctlType := typeToLbctlType(parentType)
	if lbctlType == "" {
		return NewConfError(ErrValidationError, fmt.Sprintf("Parent type %v not recognized", parentType))
	}

	return c.createObject(strconv.FormatInt(data.ID, 10), "filter", parentName, lbctlType, data, nil, transactionID, version)
}

// EditFilter edits a filter in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *LBCTLConfigurationClient) EditFilter(id int64, parentType string, parentName string, data *models.Filter, transactionID string, version int64) error {
	validationErr := data.Validate(strfmt.Default)
	if validationErr != nil {
		return NewConfError(ErrValidationError, validationErr.Error())
	}

	lbctlType := typeToLbctlType(parentType)
	if lbctlType == "" {
		return NewConfError(ErrValidationError, fmt.Sprintf("Parent type %v not recognized", parentType))
	}

	ondiskF, err := c.GetFilter(id, parentType, parentName)
	if err != nil {
		return err
	}

	return c.editObject(strconv.FormatInt(data.ID, 10), "filter", parentName, lbctlType, data, ondiskF, nil, transactionID, version)
}

func (c *LBCTLConfigurationClient) parseFilters(response string) models.Filters {
	filters := make(models.Filters, 0, 1)
	for _, filtersStr := range strings.Split(response, "\n\n") {
		if strings.TrimSpace(filtersStr) == "" {
			continue
		}
		idStr, _ := splitHeaderLine(filtersStr)
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			id = 0
		}

		filtersObj := &models.Filter{ID: id}
		c.parseObject(filtersStr, filtersObj)
		filters = append(filters, filtersObj)
	}
	return filters
}
