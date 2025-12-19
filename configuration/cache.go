// Copyright 2022 HAProxy Technologies
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

	"github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/client-native/v6/config-parser"
	"github.com/haproxytech/client-native/v6/config-parser/common"
	parser_errors "github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/types"

	"github.com/haproxytech/client-native/v6/misc"
	"github.com/haproxytech/client-native/v6/models"
)

type Cache interface {
	GetCaches(transactionID string) (int64, models.Caches, error)
	GetCache(name string, transactionID string) (int64, *models.Cache, error)
	DeleteCache(name string, transactionID string, version int64) error
	EditCache(name string, data *models.Cache, transactionID string, version int64) error
	CreateCache(data *models.Cache, transactionID string, version int64) error
}

// GetCaches returns configuration version and an array of
// configured caches. Returns error on fail.
func (c *client) GetCaches(transactionID string) (int64, models.Caches, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	fNames, err := p.SectionsGet(parser.Cache)
	if err != nil {
		return v, nil, err
	}

	var cache *models.Cache
	caches := []*models.Cache{}
	for _, name := range fNames {
		if v, cache, err = c.GetCache(name, transactionID); err == nil {
			caches = append(caches, cache)
		}
	}

	return v, caches, nil
}

// GetCache returns configuration version and a requested cache.
// Returns error on fail or if cache does not exist.
func (c *client) GetCache(name string, transactionID string) (int64, *models.Cache, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	if !p.SectionExists(parser.Cache, name) {
		return v, nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Cache %s does not exist", name))
	}

	cache := &models.Cache{Name: misc.StringP(name)}
	if err = ParseCacheSection(p, cache); err != nil {
		return 0, nil, err
	}

	return v, cache, nil
}

// DeleteCache deletes a cache in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) DeleteCache(name string, transactionID string, version int64) error {
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	if !p.SectionExists(parser.Cache, name) {
		e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("%s %s does not exist", parser.Cache, name))
		return c.HandleError(name, "", "", t, transactionID == "", e)
	}

	if err := p.SectionsDelete(parser.Cache, name); err != nil {
		return c.HandleError(name, "", "", t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

// EditCache edits a cache in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) EditCache(name string, data *models.Cache, transactionID string, version int64) error {
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

	if !p.SectionExists(parser.Cache, name) {
		e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("%s %s does not exist", parser.Cache, name))
		return c.HandleError(name, "", "", t, transactionID == "", e)
	}

	if err = SerializeCacheSection(p, data); err != nil {
		return err
	}

	return c.SaveData(p, t, transactionID == "")
}

// CreateCache creates a cache in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) CreateCache(data *models.Cache, transactionID string, version int64) error {
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

	if p.SectionExists(parser.Cache, *data.Name) {
		e := NewConfError(ErrObjectAlreadyExists, fmt.Sprintf("%s %s already exists", parser.Cache, *data.Name))
		return c.HandleError(*data.Name, "", "", t, transactionID == "", e)
	}

	if err = p.SectionsCreate(parser.Cache, *data.Name); err != nil {
		return c.HandleError(*data.Name, "", "", t, transactionID == "", err)
	}

	if err = SerializeCacheSection(p, data); err != nil {
		return err
	}

	return c.SaveData(p, t, transactionID == "")
}

func ParseCacheSection(p parser.Parser, cache *models.Cache) error {
	var err error
	var data common.ParserData
	name := *cache.Name

	if data, err := p.SectionGet(parser.Cache, name); err == nil {
		d, ok := data.(*types.Section)
		if ok && d != nil {
			cache.Metadata = misc.ParseMetadata(d.Comment)
		}
	}

	if data, err = p.Get(parser.Cache, name, "total-max-size", false); err == nil {
		d, ok := data.(*types.Int64C)
		if ok && d != nil {
			cache.TotalMaxSize = d.Value
		}
	}
	if data, err = p.Get(parser.Cache, name, "max-object-size", false); err == nil {
		d, ok := data.(*types.Int64C)
		if ok && d != nil {
			cache.MaxObjectSize = d.Value
		}
	}
	if data, err = p.Get(parser.Cache, name, "max-age", false); err == nil {
		d, ok := data.(*types.Int64C)
		if ok && d != nil {
			cache.MaxAge = d.Value
		}
	}
	if data, err = p.Get(parser.Cache, name, "max-secondary-entries", false); err == nil {
		d, ok := data.(*types.Int64C)
		if ok && d != nil {
			cache.MaxSecondaryEntries = d.Value
		}
	}
	if data, err = p.Get(parser.Cache, name, "process-vary", false); err == nil {
		d, ok := data.(*types.ProcessVary)
		if ok && d != nil {
			cache.ProcessVary = misc.BoolP(d.On)
		}
	}

	if errors.Is(err, parser_errors.ErrFetch) {
		return nil
	}
	return err
}

func SerializeCacheSection(p parser.Parser, data *models.Cache) error {
	var err error

	if data.Metadata != nil {
		comment, err := misc.SerializeMetadata(data.Metadata)
		if err != nil {
			return err
		}
		if err := p.SectionCommentSet(parser.Cache, *data.Name, comment); err != nil {
			return err
		}
	}
	if data.TotalMaxSize == 0 {
		if err = p.Set(parser.Cache, *data.Name, "total-max-size", nil); err != nil {
			return err
		}
	} else {
		n := types.Int64C{Value: data.TotalMaxSize}
		if err = p.Set(parser.Cache, *data.Name, "total-max-size", n); err != nil {
			return err
		}
	}
	if data.MaxObjectSize == 0 {
		if err = p.Set(parser.Cache, *data.Name, "max-object-size", nil); err != nil {
			return err
		}
	} else {
		n := types.Int64C{Value: data.MaxObjectSize}
		if err = p.Set(parser.Cache, *data.Name, "max-object-size", n); err != nil {
			return err
		}
	}
	if data.MaxAge == 0 {
		if err = p.Set(parser.Cache, *data.Name, "max-age", nil); err != nil {
			return err
		}
	} else {
		n := types.Int64C{Value: data.MaxAge}
		if err = p.Set(parser.Cache, *data.Name, "max-age", n); err != nil {
			return err
		}
	}
	if data.MaxSecondaryEntries == 0 {
		if err = p.Set(parser.Cache, *data.Name, "max-secondary-entries", nil); err != nil {
			return err
		}
	} else {
		n := types.Int64C{Value: data.MaxSecondaryEntries}
		if err = p.Set(parser.Cache, *data.Name, "max-secondary-entries", n); err != nil {
			return err
		}
	}
	if data.ProcessVary == nil {
		if err = p.Set(parser.Cache, *data.Name, "process-vary", nil); err != nil {
			return err
		}
	} else {
		n := types.ProcessVary{On: *data.ProcessVary}
		if err = p.Set(parser.Cache, *data.Name, "process-vary", n); err != nil {
			return err
		}
	}

	return err
}
