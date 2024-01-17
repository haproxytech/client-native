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
	"strconv"

	strfmt "github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/config-parser/v5"
	parser_errors "github.com/haproxytech/config-parser/v5/errors"
	"github.com/haproxytech/config-parser/v5/types"

	"github.com/haproxytech/client-native/v5/misc"
	"github.com/haproxytech/client-native/v5/models"
)

type LogTarget interface {
	GetLogTargets(parentType, parentName string, transactionID string) (int64, models.LogTargets, error)
	GetLogTarget(id int64, parentType, parentName string, transactionID string) (int64, *models.LogTarget, error)
	DeleteLogTarget(id int64, parentType string, parentName string, transactionID string, version int64) error
	CreateLogTarget(parentType string, parentName string, data *models.LogTarget, transactionID string, version int64) error
	EditLogTarget(id int64, parentType string, parentName string, data *models.LogTarget, transactionID string, version int64) error
}

// GetLogTargets returns configuration version and an array of
// configured log targets in the specified parent. Returns error on fail.
func (c *client) GetLogTargets(parentType, parentName string, transactionID string) (int64, models.LogTargets, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	logTargets, err := ParseLogTargets(parentType, parentName, p)
	if err != nil {
		return v, nil, c.HandleError("", parentType, parentName, "", false, err)
	}

	return v, logTargets, nil
}

// GetLogTarget returns configuration version and a requested log target
// in the specified parent. Returns error on fail or if log target does not exist.
func (c *client) GetLogTarget(id int64, parentType, parentName string, transactionID string) (int64, *models.LogTarget, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	var section parser.Section
	section, parentName = logTargetSectionType(parentType, parentName)

	data, err := p.GetOne(section, parentName, "log", int(id))
	if err != nil {
		return v, nil, c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, "", false, err)
	}

	logTarget := ParseLogTarget(data.(types.Log))
	logTarget.Index = &id

	return v, logTarget, nil
}

// DeleteLogTarget deletes a log target in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) DeleteLogTarget(id int64, parentType string, parentName string, transactionID string, version int64) error {
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	var section parser.Section
	section, parentName = logTargetSectionType(parentType, parentName)

	if err := p.Delete(section, parentName, "log", int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

// CreateLogTarget creates a log target in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) CreateLogTarget(parentType string, parentName string, data *models.LogTarget, transactionID string, version int64) error {
	if c.UseModelsValidation {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return NewConfError(ErrValidationError, validationErr.Error())
		}
	}

	// additional validation
	if data.SampleRange != "" && data.SampleSize == 0 {
		return NewConfError(ErrValidationError, "sample_range set without sample_size")
	}
	if data.SampleSize != 0 && data.SampleRange == "" {
		return NewConfError(ErrValidationError, "sample_size set without sample_range")
	}

	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	var section parser.Section
	section, parentName = logTargetSectionType(parentType, parentName)

	if err := p.Insert(section, parentName, "log", SerializeLogTarget(*data), int(*data.Index)); err != nil {
		return c.HandleError(strconv.FormatInt(*data.Index, 10), parentType, parentName, t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

// EditLogTarget edits a log target in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) EditLogTarget(id int64, parentType string, parentName string, data *models.LogTarget, transactionID string, version int64) error {
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

	// additional validation
	if data.SampleRange != "" && data.SampleSize == 0 {
		return NewConfError(ErrValidationError, "sample_range set without sample_size")
	}
	if data.SampleSize != 0 && data.SampleRange == "" {
		return NewConfError(ErrValidationError, "sample_size set without sample_range")
	}

	var section parser.Section
	section, parentName = logTargetSectionType(parentType, parentName)

	if _, err := p.GetOne(section, parentName, "log", int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	if err := p.Set(section, parentName, "log", SerializeLogTarget(*data), int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

func ParseLogTargets(t, pName string, p parser.Parser) (models.LogTargets, error) {
	var section parser.Section
	section, pName = logTargetSectionType(t, pName)

	logTargets := models.LogTargets{}
	data, err := p.Get(section, pName, "log", false)
	if err != nil {
		if errors.Is(err, parser_errors.ErrFetch) {
			return logTargets, nil
		}
		return nil, err
	}

	targets, ok := data.([]types.Log)
	if !ok {
		return nil, misc.CreateTypeAssertError("log targets")
	}
	for i, l := range targets {
		id := int64(i)
		logTarget := ParseLogTarget(l)
		if logTarget != nil {
			logTarget.Index = &id
			logTargets = append(logTargets, logTarget)
		}
	}
	return logTargets, nil
}

func ParseLogTarget(l types.Log) *models.LogTarget {
	return &models.LogTarget{
		Address:     l.Address,
		Facility:    l.Facility,
		Format:      l.Format,
		Global:      l.Global,
		Length:      l.Length,
		Level:       l.Level,
		Minlevel:    l.MinLevel,
		Nolog:       l.NoLog,
		SampleRange: l.SampleRange,
		SampleSize:  l.SampleSize,
	}
}

func SerializeLogTarget(l models.LogTarget) types.Log {
	return types.Log{
		Address:     l.Address,
		Facility:    l.Facility,
		Format:      l.Format,
		Global:      l.Global,
		Length:      l.Length,
		Level:       l.Level,
		MinLevel:    l.Minlevel,
		NoLog:       l.Nolog,
		SampleRange: l.SampleRange,
		SampleSize:  l.SampleSize,
	}
}

func logTargetSectionType(parentType string, parentName string) (parser.Section, string) {
	var section parser.Section
	switch parentType {
	case GlobalParentName:
		section = parser.Global
		parentName = parser.GlobalSectionName
	case DefaultsParentName:
		section = parser.Defaults
		if parentName == "" {
			parentName = parser.DefaultSectionName
		}
	case BackendParentName:
		section = parser.Backends
	case FrontendParentName:
		section = parser.Frontends
	case LogForwardParentName:
		section = parser.LogForward
	case PeersParentName:
		section = parser.Peers
	}
	return section, parentName
}
