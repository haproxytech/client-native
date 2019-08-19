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
	"strconv"
	"strings"

	"github.com/haproxytech/client-native/misc"
	"github.com/haproxytech/config-parser/parsers/tcp/actions"

	strfmt "github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/config-parser"
	parser_errors "github.com/haproxytech/config-parser/errors"
	"github.com/haproxytech/config-parser/types"
	"github.com/haproxytech/models"
)

// GetTCPRequestRules returns configuration version and an array of
// configured TCP request rules in the specified parent. Returns error on fail.
func (c *Client) GetTCPRequestRules(parentType, parentName string, transactionID string) (int64, models.TCPRequestRules, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	tcpRules, err := c.parseTCPRequestRules(parentType, parentName, p)
	if err != nil {
		return v, nil, c.handleError("", parentType, parentName, "", false, err)
	}

	return v, tcpRules, nil
}

// GetTCPRequestRule returns configuration version and a requested tcp request rule
// in the specified parent. Returns error on fail or if http request rule does not exist.
func (c *Client) GetTCPRequestRule(id int64, parentType, parentName string, transactionID string) (int64, *models.TCPRequestRule, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	var section parser.Section
	if parentType == "backend" {
		section = parser.Backends
	} else if parentType == "frontend" {
		section = parser.Frontends
	}

	data, err := p.GetOne(section, parentName, "tcp-request", int(id))
	if err != nil {
		return v, nil, c.handleError(strconv.FormatInt(id, 10), parentType, parentName, "", false, err)
	}

	tcpRule := parseTCPRequestRule(data.(types.TCPAction))
	tcpRule.ID = &id

	return v, tcpRule, nil
}

// DeleteTCPRequestRule deletes a tcp request rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) DeleteTCPRequestRule(id int64, parentType string, parentName string, transactionID string, version int64) error {
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

	if err := p.Delete(section, parentName, "tcp-request", int(id)); err != nil {
		return c.handleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	if err := c.saveData(p, t, transactionID == ""); err != nil {
		return err
	}

	return nil
}

// CreateTCPRequestRule creates a tcp request rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) CreateTCPRequestRule(parentType string, parentName string, data *models.TCPRequestRule, transactionID string, version int64) error {
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

	if err := p.Insert(section, parentName, "tcp-request", serializeTCPRequestRule(*data), int(*data.ID)); err != nil {
		return c.handleError(strconv.FormatInt(*data.ID, 10), parentType, parentName, t, transactionID == "", err)
	}

	if err := c.saveData(p, t, transactionID == ""); err != nil {
		return err
	}
	return nil
}

// EditTCPRequestRule edits a tcp request rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) EditTCPRequestRule(id int64, parentType string, parentName string, data *models.TCPRequestRule, transactionID string, version int64) error {
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

	if _, err := p.GetOne(section, parentName, "tcp-request", int(id)); err != nil {
		return c.handleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	if err := p.Set(section, parentName, "tcp-request", serializeTCPRequestRule(*data), int(id)); err != nil {
		return c.handleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	if err := c.saveData(p, t, transactionID == ""); err != nil {
		return err
	}
	return nil
}

func (c *Client) parseTCPRequestRules(t, pName string, p *parser.Parser) (models.TCPRequestRules, error) {
	section := parser.Global
	if t == "frontend" {
		section = parser.Frontends
	} else if t == "backend" {
		section = parser.Backends
	}

	tcpReqRules := models.TCPRequestRules{}
	data, err := p.Get(section, pName, "tcp-request", false)
	if err != nil {
		if err == parser_errors.ErrFetch {
			return tcpReqRules, nil
		}
		return nil, err
	}

	rules := data.([]types.TCPAction)
	for i, r := range rules {
		id := int64(i)
		tcpReqRule := parseTCPRequestRule(r)
		if tcpReqRule != nil {
			tcpReqRule.ID = &id
			tcpReqRules = append(tcpReqRules, tcpReqRule)
		}
	}
	return tcpReqRules, nil
}

func parseTCPRequestRule(f types.TCPAction) *models.TCPRequestRule {
	switch v := f.(type) {
	case *actions.Connection:
		r := &models.TCPRequestRule{
			Type:     "connection",
			Cond:     v.Cond,
			CondTest: v.CondTest,
		}
		if strings.Join(v.Action, " ") == "accept" {
			r.Action = "accept"
		} else if strings.Join(v.Action, " ") == "reject" {
			r.Action = "reject"
		} else {
			return nil
		}
		return r
	case *actions.Content:
		r := &models.TCPRequestRule{
			Type:     "content",
			Cond:     v.Cond,
			CondTest: v.CondTest,
		}
		if strings.Join(v.Action, " ") == "accept" {
			r.Action = "accept"
		} else if strings.Join(v.Action, " ") == "reject" {
			r.Action = "reject"
		} else {
			return nil
		}
		return r
	case *actions.InspectDelay:
		return &models.TCPRequestRule{
			Type:    "inspect-delay",
			Timeout: misc.ParseTimeout(v.Timeout),
		}
	case *actions.Session:
		r := &models.TCPRequestRule{
			Type:     "session",
			Cond:     v.Cond,
			CondTest: v.CondTest,
		}
		if strings.Join(v.Action, " ") == "accept" {
			r.Action = "accept"
		} else if strings.Join(v.Action, " ") == "reject" {
			r.Action = "reject"
		} else {
			return nil
		}
		return r
	}
	return nil
}

func serializeTCPRequestRule(f models.TCPRequestRule) types.TCPAction {
	switch f.Type {
	case "connection":
		return &actions.Connection{
			Action:   []string{f.Action},
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "content":
		return &actions.Content{
			Action:   []string{f.Action},
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	case "inspect-delay":
		if f.Timeout != nil {
			return &actions.InspectDelay{
				Timeout: strconv.FormatInt(*f.Timeout, 10),
			}
		}
	case "session":
		return &actions.Session{
			Action:   []string{f.Action},
			Cond:     f.Cond,
			CondTest: f.CondTest,
		}
	}
	return nil
}
