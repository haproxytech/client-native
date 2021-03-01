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

package spoe

import (
	"fmt"

	"github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/config-parser/v3"
	"github.com/haproxytech/config-parser/v3/spoe"
	"github.com/haproxytech/config-parser/v3/types"

	conf "github.com/haproxytech/client-native/v2/configuration"
	"github.com/haproxytech/client-native/v2/models"
)

// GetGroups returns configuration version and an array of
// configured messages. Returns error on fail.
func (c *SingleSpoe) GetGroups(scope, transactionID string) (int64, models.SpoeGroups, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}
	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	aNames, err := p.SectionsGet(scope, parser.SPOEGroup)
	if err != nil {
		return v, nil, err
	}

	messages := models.SpoeGroups{}
	for _, name := range aNames {
		_, a, err := c.GetGroup(scope, name, transactionID)
		if err == nil {
			messages = append(messages, a)
		}
	}

	return v, messages, nil
}

// GetGroup returns configuration version and a requested group
// Returns an error on fail or if group does not exist.
func (c *SingleSpoe) GetGroup(scope, name, transactionID string) (int64, *models.SpoeGroup, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}
	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	if !c.checkSectionExists(scope, parser.SPOEGroup, name, p) {
		return v, nil, conf.NewConfError(conf.ErrObjectDoesNotExist, fmt.Sprintf("group %s does not exist", name))
	}

	group := &models.SpoeGroup{Name: &name}

	data, err := p.Get(scope, parser.SPOEGroup, name, "messages", true)
	if err != nil {
		return v, nil, err
	}
	if d, ok := data.(*types.StringC); ok {
		group.Messages = d.Value
	}

	return v, group, nil
}

// DeleteGroup deletes an group in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *SingleSpoe) DeleteGroup(scope, name, transactionID string, version int64) error {
	if err := c.deleteSection(scope, parser.SPOEGroup, name, transactionID, version); err != nil {
		return err
	}
	return nil
}

// CreateGroup creates a group in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *SingleSpoe) CreateGroup(scope string, data *models.SpoeGroup, transactionID string, version int64) error {
	if c.Transaction.UseValidation {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return conf.NewConfError(conf.ErrValidationError, validationErr.Error())
		}
	}

	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	name := *data.Name
	if c.checkSectionExists(scope, parser.SPOEGroup, name, p) {
		e := conf.NewConfError(conf.ErrObjectAlreadyExists, fmt.Sprintf("%s %s already exists", parser.SPOEGroup, name))
		return c.Transaction.HandleError(name, "", "", t, transactionID == "", e)
	}

	if err = p.SectionsCreate(scope, parser.SPOEGroup, name); err != nil {
		return err
	}

	err = c.createEditGroup(scope, data, t, transactionID, p)
	if err != nil {
		return err
	}

	if err := c.Transaction.SaveData(p, t, transactionID == ""); err != nil {
		return err
	}

	return nil
}

// EditMessage edits a group in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *SingleSpoe) EditGroup(scope string, data *models.SpoeGroup, name, transactionID string, version int64) error {
	if c.Transaction.UseValidation {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return conf.NewConfError(conf.ErrValidationError, validationErr.Error())
		}
	}

	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	if !c.checkSectionExists(scope, parser.SPOEGroup, *data.Name, p) {
		e := conf.NewConfError(conf.ErrObjectAlreadyExists, fmt.Sprintf("%s %s does not exists", parser.SPOEGroup, *data.Name))
		return c.Transaction.HandleError(*data.Name, "", "", t, transactionID == "", e)
	}

	err = c.createEditGroup(scope, data, t, transactionID, p)
	if err != nil {
		return err
	}

	if err := c.Transaction.SaveData(p, t, transactionID == ""); err != nil {
		return err
	}

	return nil
}

func (c *SingleSpoe) createEditGroup(scope string, data *models.SpoeGroup, t string, transactionID string, p *spoe.Parser) error {
	if data == nil {
		return fmt.Errorf("spoe group not initialized")
	}
	name := *data.Name

	if data.Messages != "" {
		d := &types.StringC{Value: data.Messages}
		if err := p.Set(scope, parser.SPOEGroup, name, "messages", d); err != nil {
			return c.Transaction.HandleError("messages", "", "", t, transactionID == "", err)
		}
	} else if err := p.Set(scope, parser.SPOEGroup, name, "messages", nil); err != nil {
		return err
	}

	return nil
}
