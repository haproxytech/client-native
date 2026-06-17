// Copyright 2026 HAProxy Technologies
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package configuration

import (
	"errors"
	"fmt"

	"github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/client-native/v6/config-parser"
	parser_errors "github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/types"
	"github.com/haproxytech/client-native/v6/misc"
	"github.com/haproxytech/client-native/v6/models"
)

type ForceBeSwitch interface {
	GetForceBeSwitches(parentType string, parentName string, transactionID string) (int64, models.ForceBeSwitches, error)
	GetForceBeSwitch(number int64, parentType string, parentName string, transactionID string) (int64, *models.ForceBeSwitch, error)
	DeleteForceBeSwitch(number int64, parentType string, parentName string, transactionID string, version int64) error
	CreateForceBeSwitch(parentType string, parentName string, data *models.ForceBeSwitch, transactionID string, version int64) error
	EditForceBeSwitch(number int64, parentType string, parentName string, data *models.ForceBeSwitch, transactionID string, version int64) error
}

func (c *client) GetForceBeSwitches(parentType string, parentName string, transactionID string) (int64, models.ForceBeSwitches, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	rules, err := ParseForceBeSwitches(parentType, parentName, p)
	if err != nil {
		return v, nil, c.HandleError("", parentType, parentName, "", false, err)
	}

	return v, rules, nil
}

func (c *client) GetForceBeSwitch(number int64, parentType string, parentName string, transactionID string) (int64, *models.ForceBeSwitch, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	rule := GetForceBeSwitchByNumber(number, parentType, parentName, p)
	if rule == nil {
		return v, nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("ForceBeSwitch %d does not exist in %s %s", number, parentName, parentType))
	}

	return v, rule, nil
}

func (c *client) DeleteForceBeSwitch(number int64, parentType string, parentName string, transactionID string, version int64) error {
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	rule := GetForceBeSwitchByNumber(number, parentType, parentName, p)
	if rule == nil {
		e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("ForceBeSwitch %d does not exist in %s %s", number, parentName, parentType))
		return c.HandleError("force-be-switch", parentType, parentName, t, transactionID == "", e)
	}

	if err := p.Delete(parser.Section(parentType), parentName, "force-be-switch", int(number)); err != nil {
		return c.HandleError("force-be-switch", parentType, parentName, t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

func (c *client) CreateForceBeSwitch(parentType string, parentName string, data *models.ForceBeSwitch, transactionID string, version int64) error {
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

	if err := p.Insert(parser.Section(parentType), parentName, "force-be-switch", SerializeForceBeSwitch(*data)); err != nil {
		return c.HandleError("force-be-switch", parentType, parentName, t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

func (c *client) EditForceBeSwitch(number int64, parentType string, parentName string, data *models.ForceBeSwitch, transactionID string, version int64) error {
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

	if err := p.Set(parser.Section(parentType), parentName, "force-be-switch", SerializeForceBeSwitch(*data), int(number)); err != nil {
		return c.HandleError("force-be-switch", parentType, parentName, t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

func ParseForceBeSwitches(parentType string, parentName string, p parser.Parser) (models.ForceBeSwitches, error) {
	var rules models.ForceBeSwitches

	data, err := p.Get(parser.Section(parentType), parentName, "force-be-switch")
	if err != nil {
		if errors.Is(err, parser_errors.ErrFetch) {
			return rules, nil
		}
		return nil, err
	}

	ondiskRules := data.([]types.ForceBeSwitch) //nolint:forcetypeassert
	rules = make(models.ForceBeSwitches, 0, len(ondiskRules))
	for _, ondiskRule := range ondiskRules {
		r := ParseForceBeSwitch(ondiskRule)
		if r != nil {
			rules = append(rules, r)
		}
	}
	return rules, nil
}

func ParseForceBeSwitch(ondisk types.ForceBeSwitch) *models.ForceBeSwitch {
	cond := ondisk.Cond
	condTest := ondisk.CondTest
	return &models.ForceBeSwitch{
		Cond:     &cond,
		CondTest: &condTest,
		Metadata: misc.ParseMetadata(ondisk.Comment),
	}
}

func SerializeForceBeSwitch(m models.ForceBeSwitch) types.ForceBeSwitch {
	comment, err := misc.SerializeMetadata(m.Metadata)
	if err != nil {
		comment = ""
	}
	var cond, condTest string
	if m.Cond != nil {
		cond = *m.Cond
	}
	if m.CondTest != nil {
		condTest = *m.CondTest
	}
	return types.ForceBeSwitch{
		Cond:     cond,
		CondTest: condTest,
		Comment:  comment,
	}
}

func GetForceBeSwitchByNumber(number int64, parentType string, parentName string, p parser.Parser) *models.ForceBeSwitch {
	rules, err := ParseForceBeSwitches(parentType, parentName, p)
	if err != nil || int64(len(rules)) < number+1 {
		return nil
	}
	return rules[number]
}
