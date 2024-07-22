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
	"fmt"
	"strconv"

	"github.com/go-openapi/strfmt"
	"github.com/haproxytech/client-native/v6/models"
	parser "github.com/haproxytech/config-parser/v5"
)

type StructuredLogForward interface {
	GetStructuredLogForwards(transactionID string) (int64, models.LogForwards, error)
	GetStructuredLogForward(name string, transactionID string) (int64, *models.LogForward, error)
	CreateStructuredLogForward(data *models.LogForward, transactionID string, version int64) error
	EditStructuredLogForward(name string, data *models.LogForward, transactionID string, version int64) error
}

// GetStructuredLogForward returns configuration version and a requested log_forward with all its child resources.
// Returns error on fail or if log_forward does not exist.
func (c *client) GetStructuredLogForward(name string, transactionID string) (int64, *models.LogForward, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	if !c.checkSectionExists(parser.LogForward, name, p) {
		return v, nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("LogForward %s does not exist", name))
	}

	f, err := parseLogForwardsSection(name, p)

	return v, f, err
}

func (c *client) GetStructuredLogForwards(transactionID string) (int64, models.LogForwards, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	logForwards, err := parseLogForwardsSections(p)
	if err != nil {
		return 0, nil, err
	}

	return v, logForwards, nil
}

// EditStructuredLogForward replaces a log_forward and all it's child resources in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) EditStructuredLogForward(name string, data *models.LogForward, transactionID string, version int64) error {
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

	if !c.checkSectionExists(parser.LogForward, name, p) {
		e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("%s %s does not exist", parser.LogForward, name))
		return c.HandleError(name, "", "", t, transactionID == "", e)
	}

	if err = p.SectionsDelete(parser.LogForward, name); err != nil {
		return c.HandleError(name, "", "", t, transactionID == "", err)
	}

	if err = serializeLogForwardSection(StructuredToParserArgs{
		TID:                transactionID,
		Parser:             &p,
		Options:            &c.ConfigurationOptions,
		HandleError:        c.HandleError,
		CheckSectionExists: c.checkSectionExists,
	}, data); err != nil {
		return err
	}
	return c.SaveData(p, t, transactionID == "")
}

// CreateStructuredLogForward creates a log_forward and all it's child resources in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) CreateStructuredLogForward(data *models.LogForward, transactionID string, version int64) error {
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

	if c.checkSectionExists(parser.LogForward, data.Name, p) {
		e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("%s %s already exist", parser.LogForward, data.Name))
		return c.HandleError(data.Name, "", "", t, transactionID == "", e)
	}

	if err = serializeLogForwardSection(StructuredToParserArgs{
		TID:                transactionID,
		Parser:             &p,
		Options:            &c.ConfigurationOptions,
		HandleError:        c.HandleError,
		CheckSectionExists: c.checkSectionExists,
	}, data); err != nil {
		return err
	}
	return c.SaveData(p, t, transactionID == "")
}

func parseLogForwardsSections(p parser.Parser) (models.LogForwards, error) {
	names, err := p.SectionsGet(parser.LogForward)
	if err != nil {
		return nil, err
	}
	logForwards := []*models.LogForward{}
	for _, name := range names {
		f, err := parseLogForwardsSection(name, p)
		if err != nil {
			return nil, err
		}
		logForwards = append(logForwards, f)
	}
	return logForwards, nil
}

func parseLogForwardsSection(name string, p parser.Parser) (*models.LogForward, error) {
	lf := &models.LogForward{LogForwardBase: models.LogForwardBase{Name: name}}
	if err := ParseLogForward(p, lf); err != nil {
		return nil, err
	}
	// binds
	b, err := ParseBinds(LogForwardParentName, name, p)
	if err != nil {
		return nil, err
	}
	ba, errba := namedResourceArrayToMap(b)
	if errba != nil {
		return nil, errba
	}
	lf.Binds = ba

	// dgram binds
	d, err := ParseDgramBinds(name, p)
	if err != nil {
		return nil, err
	}
	da, errda := namedResourceArrayToMap(d)
	if errda != nil {
		return nil, errda
	}
	lf.DgramBinds = da

	// log targets
	logTargets, err := ParseLogTargets(LogForwardParentName, name, p)
	if err != nil {
		return nil, err
	}
	lf.LogTargetList = logTargets

	return lf, nil
}

func serializeLogForwardSection(a StructuredToParserArgs, lf *models.LogForward) error {
	p := *a.Parser
	var err error

	err = p.SectionsCreate(parser.LogForward, lf.Name)
	if err != nil {
		return err
	}
	if err = SerializeLogForwardSection(p, lf, a.Options); err != nil {
		return err
	}
	for _, dgramBind := range lf.DgramBinds {
		if err = p.Insert(parser.LogForward, lf.Name, "dgram-bind", SerializeDgramBind(dgramBind), -1); err != nil {
			return a.HandleError(dgramBind.Name, "log-forward", lf.Name, a.TID, a.TID == "", err)
		}
	}
	for _, bind := range lf.Binds {
		if err = p.Insert(parser.LogForward, lf.Name, "bind", SerializeBind(bind), -1); err != nil {
			return a.HandleError(bind.Name, "log-forward", lf.Name, a.TID, a.TID == "", err)
		}
	}
	for i, logTarget := range lf.LogTargetList {
		if err = p.Insert(parser.LogForward, lf.Name, "log", SerializeLogTarget(*logTarget), i); err != nil {
			return a.HandleError(strconv.FormatInt(int64(i), 10), "log-forward", lf.Name, a.TID, a.TID == "", err)
		}
	}

	return nil
}
