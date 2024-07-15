// Copyright 2019 HAProxy Technologies
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package configuration

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/config-parser/v5"
	parser_errors "github.com/haproxytech/config-parser/v5/errors"
	"github.com/haproxytech/config-parser/v5/parsers/filters"
	"github.com/haproxytech/config-parser/v5/types"

	"github.com/haproxytech/client-native/v6/misc"
	"github.com/haproxytech/client-native/v6/models"
)

type Filter interface {
	GetFilters(parentType, parentName string, transactionID string) (int64, models.Filters, error)
	GetFilter(id int64, parentType, parentName string, transactionID string) (int64, *models.Filter, error)
	DeleteFilter(id int64, parentType string, parentName string, transactionID string, version int64) error
	CreateFilter(id int64, parentType string, parentName string, data *models.Filter, transactionID string, version int64) error
	EditFilter(id int64, parentType string, parentName string, data *models.Filter, transactionID string, version int64) error
	ReplaceFilters(parentType string, parentName string, data models.Filters, transactionID string, version int64) error
}

// GetFilters returns configuration version and an array of
// configured filters in the specified parent. Returns error on fail.
func (c *client) GetFilters(parentType, parentName string, transactionID string) (int64, models.Filters, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	filters, err := ParseFilters(parentType, parentName, p)
	if err != nil {
		return v, nil, c.HandleError("", parentType, parentName, "", false, err)
	}

	return v, filters, nil
}

// GetFilter returns configuration version and a requested filter
// in the specified parent. Returns error on fail or if filter does not exist.
func (c *client) GetFilter(id int64, parentType, parentName string, transactionID string) (int64, *models.Filter, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	var section parser.Section
	if parentType == BackendParentName {
		section = parser.Backends
	} else if parentType == FrontendParentName {
		section = parser.Frontends
	}

	data, err := p.GetOne(section, parentName, "filter", int(id))
	if err != nil {
		return v, nil, c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, "", false, err)
	}

	filter := ParseFilter(data.(types.Filter))

	return v, filter, nil
}

// DeleteFilter deletes a filter in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) DeleteFilter(id int64, parentType string, parentName string, transactionID string, version int64) error {
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	var section parser.Section
	if parentType == BackendParentName {
		section = parser.Backends
	} else if parentType == FrontendParentName {
		section = parser.Frontends
	}

	if err := p.Delete(section, parentName, "filter", int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

// CreateFilter creates a filter in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) CreateFilter(id int64, parentType string, parentName string, data *models.Filter, transactionID string, version int64) error {
	if c.UseModelsValidation {
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
	if parentType == BackendParentName {
		section = parser.Backends
	} else if parentType == FrontendParentName {
		section = parser.Frontends
	}

	if err := p.Insert(section, parentName, "filter", SerializeFilter(*data), int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

// EditFilter edits a filter in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) EditFilter(id int64, parentType string, parentName string, data *models.Filter, transactionID string, version int64) error {
	if c.UseModelsValidation {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return NewConfError(ErrValidationError, validationErr.Error())
		}
	}
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	section, err := getFilterParserFromParent(parentType)
	if err != nil {
		return err
	}

	if _, err := p.GetOne(section, parentName, "filter", int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	if err := p.Set(section, parentName, "filter", SerializeFilter(*data), int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

// ReplaceFilters replaces all Filter lines in configuration for a parentType/parentName.
// One of version or transactionID is mandatory.
// Returns error on fail, nil on success.
func (c *client) ReplaceFilters(parentType string, parentName string, data models.Filters, transactionID string, version int64) error {
	if c.UseModelsValidation {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return NewConfError(ErrValidationError, validationErr.Error())
		}
	}
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	section, err := getFilterParserFromParent(parentType)
	if err != nil {
		return err
	}

	filters, err := ParseFilters(parentType, parentName, p)
	if err != nil {
		return c.HandleError("", parentType, parentName, "", false, err)
	}

	for i := range filters {
		// Always delete index 0
		if err := p.Delete(section, parentName, "filter", 0); err != nil {
			return c.HandleError(strconv.FormatInt(int64(i), 10), parentType, parentName, t, transactionID == "", err)
		}
	}

	for i, newFilter := range data {
		if err := p.Insert(section, parentName, "filter", SerializeFilter(*newFilter), i); err != nil {
			return c.HandleError(strconv.FormatInt(int64(i), 10), parentType, parentName, t, transactionID == "", err)
		}
	}

	return c.SaveData(p, t, transactionID == "")
}

func getFilterParserFromParent(parent string) (parser.Section, error) {
	switch parent {
	case BackendParentName:
		return parser.Backends, nil
	case FrontendParentName:
		return parser.Frontends, nil
	default:
		return "", fmt.Errorf("unsupported parent: %s", parent)
	}
}

func ParseFilters(t, pName string, p parser.Parser) (models.Filters, error) {
	section, err := getFilterParserFromParent(t)
	if err != nil {
		return nil, err
	}

	f := models.Filters{}
	data, err := p.Get(section, pName, "filter", false)
	if err != nil {
		if errors.Is(err, parser_errors.ErrFetch) {
			return f, nil
		}
		return nil, err
	}

	filters, ok := data.([]types.Filter)
	if !ok {
		return nil, misc.CreateTypeAssertError("filter")
	}
	for _, filter := range filters {
		mFilter := ParseFilter(filter)
		if mFilter != nil {
			f = append(f, mFilter)
		}
	}
	return f, nil
}

func ParseFilter(f types.Filter) *models.Filter {
	switch v := f.(type) {
	case *filters.BandwidthLimit:
		filter := &models.Filter{
			BandwidthLimitName: v.Name,
			Type:               v.Attribute,
		}
		if len(v.Limit) > 0 && len(v.Key) > 0 {
			filter.Key = v.Key
			limit := misc.ParseSize(v.Limit)
			if limit != nil {
				filter.Limit = *limit
			}

			if table := v.Table; table != nil {
				filter.Table = *table
			}
		} else {
			defaultLimit := misc.ParseSize(v.DefaultLimit)
			if defaultLimit != nil {
				filter.DefaultLimit = *defaultLimit
			}
			defaultPeriod := misc.ParseTimeout(v.DefaultPeriod)
			if defaultPeriod != nil {
				filter.DefaultPeriod = *defaultPeriod
			}
		}
		if minSize := v.MinSize; minSize != nil {
			minSizeValue := misc.ParseSize(*v.MinSize)
			if minSizeValue != nil {
				filter.MinSize = *minSizeValue
			}
		}
		return filter
	case *filters.FcgiApp:
		return &models.Filter{
			Type:    FCGIAppParentName,
			AppName: v.Name,
		}
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

func SerializeFilter(f models.Filter) types.Filter {
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
	case "fcgi-app":
		return &filters.FcgiApp{
			Name: f.AppName,
		}
	case "bwlim-in":
		return &filters.BandwidthLimit{
			Attribute:     "bwlim-in",
			Name:          f.BandwidthLimitName,
			DefaultLimit:  misc.SerializeSize(f.DefaultLimit),
			DefaultPeriod: misc.SerializeTime(f.DefaultPeriod),
			Limit:         misc.SerializeSize(f.Limit),
			Key:           f.Key,
			Table:         &f.Table,
		}
	case "bwlim-out":
		return &filters.BandwidthLimit{
			Attribute:     "bwlim-out",
			Name:          f.BandwidthLimitName,
			DefaultLimit:  misc.SerializeSize(f.DefaultLimit),
			DefaultPeriod: misc.SerializeTime(f.DefaultPeriod),
			Limit:         misc.SerializeSize(f.Limit),
			Key:           f.Key,
			Table:         &f.Table,
		}
	}
	return nil
}
