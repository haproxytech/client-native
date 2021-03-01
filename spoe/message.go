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
	spoe_types "github.com/haproxytech/config-parser/v3/spoe/types"
	"github.com/haproxytech/config-parser/v3/types"

	conf "github.com/haproxytech/client-native/v2/configuration"
	"github.com/haproxytech/client-native/v2/models"
)

// GetMessages returns configuration version and an array of
// configured messages. Returns error on fail.
func (c *SingleSpoe) GetMessages(scope, transactionID string) (int64, models.SpoeMessages, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}
	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	aNames, err := p.SectionsGet(scope, parser.SPOEMessage)
	if err != nil {
		return v, nil, err
	}

	messages := models.SpoeMessages{}
	for _, name := range aNames {
		_, a, err := c.GetMessage(scope, name, transactionID)
		if err == nil {
			messages = append(messages, a)
		}
	}

	return v, messages, nil
}

// GetMessage returns configuration version and a requested message.
// Returns error on fail or if message does not exist.
func (c *SingleSpoe) GetMessage(scope, name, transactionID string) (int64, *models.SpoeMessage, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}
	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	if !c.checkSectionExists(scope, parser.SPOEMessage, name, p) {
		return v, nil, conf.NewConfError(conf.ErrObjectDoesNotExist, fmt.Sprintf("message %s does not exist", name))
	}

	message := &models.SpoeMessage{Name: &name}

	data, err := p.Get(scope, parser.SPOEMessage, name, "acl", true)
	if err != nil {
		return v, nil, err
	}
	if acls, ok := data.([]types.ACL); ok {
		for i, a := range acls {
			indx := int64(i)
			acl := &models.ACL{
				ACLName:   a.Name,
				Value:     a.Value,
				Criterion: a.Criterion,
				Index:     &indx,
			}
			message.ACL = append(message.ACL, acl)
		}
	}

	data, err = p.Get(scope, parser.SPOEMessage, name, "args", true)
	if err != nil {
		return v, nil, err
	}
	if d, ok := data.(*types.StringC); ok {
		if d.Value != "" {
			message.Args = d.Value
		}
	}

	data, err = p.Get(scope, parser.SPOEMessage, name, "event", true)
	if err != nil {
		return v, nil, err
	}
	if d, ok := data.(*spoe_types.Event); ok {
		message.Event = &models.SpoeMessageEvent{
			Name:     &d.Name,
			Cond:     d.Cond,
			CondTest: d.CondTest,
		}
	}

	return v, message, nil
}

// DeleteMessage deletes an message in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *SingleSpoe) DeleteMessage(scope, name, transactionID string, version int64) error {
	if err := c.deleteSection(scope, parser.SPOEMessage, name, transactionID, version); err != nil {
		return err
	}
	return nil
}

// CreateMessage creates a message in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *SingleSpoe) CreateMessage(scope string, data *models.SpoeMessage, transactionID string, version int64) error {
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

	if c.checkSectionExists(scope, parser.SPOEMessage, *data.Name, p) {
		e := conf.NewConfError(conf.ErrObjectAlreadyExists, fmt.Sprintf("%s %s already exists", parser.SPOEMessage, *data.Name))
		return c.Transaction.HandleError(*data.Name, "", "", t, transactionID == "", e)
	}

	if err = p.SectionsCreate(scope, parser.SPOEMessage, *data.Name); err != nil {
		return err
	}

	err = c.createEditMessage(scope, data, t, transactionID, p)
	if err != nil {
		return err
	}

	if err := c.Transaction.SaveData(p, t, transactionID == ""); err != nil {
		return err
	}

	return nil
}

// EditMessage edits a message in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *SingleSpoe) EditMessage(scope string, data *models.SpoeMessage, name, transactionID string, version int64) error {
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

	if !c.checkSectionExists(scope, parser.SPOEMessage, *data.Name, p) {
		e := conf.NewConfError(conf.ErrObjectAlreadyExists, fmt.Sprintf("%s %s does not exists", parser.SPOEMessage, *data.Name))
		return c.Transaction.HandleError(*data.Name, "", "", t, transactionID == "", e)
	}

	err = c.createEditMessage(scope, data, t, transactionID, p)
	if err != nil {
		return err
	}

	if err := c.Transaction.SaveData(p, t, transactionID == ""); err != nil {
		return err
	}

	return nil
}

func (c *SingleSpoe) createEditMessage(scope string, data *models.SpoeMessage, t string, transactionID string, p *spoe.Parser) error {
	if data == nil {
		return fmt.Errorf("spoe message not initialized")
	}
	name := *data.Name

	if len(data.ACL) > 0 {
		acls := []types.ACL{}
		for _, d := range data.ACL {
			acl := types.ACL{
				Criterion: d.Criterion,
				Name:      d.ACLName,
				Value:     d.Value,
			}
			acls = append(acls, acl)
		}
		if err := p.Set(scope, parser.SPOEMessage, name, "acl", acls); err != nil {
			return c.Transaction.HandleError("acl", "", "", t, transactionID == "", err)
		}
	} else if err := p.Set(scope, parser.SPOEMessage, name, "acl", nil); err != nil {
		return err
	}

	if data.Args != "" {
		d := &types.StringC{Value: data.Args}
		if err := p.Set(scope, parser.SPOEMessage, name, "args", d); err != nil {
			return c.Transaction.HandleError("args", "", "", t, transactionID == "", err)
		}
	} else if err := p.Set(scope, parser.SPOEMessage, name, "args", nil); err != nil {
		return err
	}

	if data.Event != nil {
		d := &spoe_types.Event{
			Cond:     data.Event.Cond,
			CondTest: data.Event.CondTest,
			Name:     *data.Event.Name,
		}
		if err := p.Set(scope, parser.SPOEMessage, name, "event", d); err != nil {
			return c.Transaction.HandleError("event", "", "", t, transactionID == "", err)
		}
	} else if err := p.Set(scope, parser.SPOEMessage, name, "event", nil); err != nil {
		return err
	}

	return nil
}
