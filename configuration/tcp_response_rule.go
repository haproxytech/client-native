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

	"github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/config-parser/v3"
	parser_errors "github.com/haproxytech/config-parser/v3/errors"
	tcp_actions "github.com/haproxytech/config-parser/v3/parsers/tcp/actions"
	tcp_types "github.com/haproxytech/config-parser/v3/parsers/tcp/types"
	"github.com/haproxytech/config-parser/v3/types"

	"github.com/haproxytech/client-native/v2/misc"
	"github.com/haproxytech/client-native/v2/models"
)

// GetTCPResponseRules returns configuration version and an array of
// configured tcp response rules in the specified backend. Returns error on fail.
func (c *Client) GetTCPResponseRules(backend string, transactionID string) (int64, models.TCPResponseRules, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	tcpRules, err := ParseTCPResponseRules(backend, p)
	if err != nil {
		return v, nil, c.HandleError("", "backend", backend, "", false, err)
	}

	return v, tcpRules, nil
}

// GetTCPResponseRule returns configuration version and a requested tcp response rule
// in the specified backend. Returns error on fail or if tcp response rule does not exist.
func (c *Client) GetTCPResponseRule(id int64, backend string, transactionID string) (int64, *models.TCPResponseRule, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	data, err := p.GetOne(parser.Backends, backend, "tcp-response", int(id))
	if err != nil {
		return v, nil, c.HandleError(strconv.FormatInt(id, 10), "backend", backend, "", false, err)
	}

	tcpRule := ParseTCPResponseRule(data.(types.TCPType))
	tcpRule.Index = &id

	return v, tcpRule, nil
}

// DeleteTCPResponseRule deletes a tcp response rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) DeleteTCPResponseRule(id int64, backend string, transactionID string, version int64) error {
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	if err := p.Delete(parser.Backends, backend, "tcp-response", int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), "backend", backend, t, transactionID == "", err)
	}

	if err := c.SaveData(p, t, transactionID == ""); err != nil {
		return err
	}
	return nil
}

// CreateTCPResponseRule creates a tcp response rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) CreateTCPResponseRule(backend string, data *models.TCPResponseRule, transactionID string, version int64) error {
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

	if err := p.Insert(parser.Backends, backend, "tcp-response", SerializeTCPResponseRule(*data), int(*data.Index)); err != nil {
		return c.HandleError(strconv.FormatInt(*data.Index, 10), "backend", backend, t, transactionID == "", err)
	}

	if err := c.SaveData(p, t, transactionID == ""); err != nil {
		return err
	}
	return nil
}

// EditTCPResponseRule edits a tcp response rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) EditTCPResponseRule(id int64, backend string, data *models.TCPResponseRule, transactionID string, version int64) error {
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

	if _, err := p.GetOne(parser.Backends, backend, "tcp-response", int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(*data.Index, 10), "backend", backend, t, transactionID == "", err)
	}

	if err := p.Set(parser.Backends, backend, "tcp-response", SerializeTCPResponseRule(*data), int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(*data.Index, 10), "backend", backend, t, transactionID == "", err)
	}

	if err := c.SaveData(p, t, transactionID == ""); err != nil {
		return err
	}
	return nil
}

func ParseTCPResponseRules(backend string, p *parser.Parser) (models.TCPResponseRules, error) {
	tcpResRules := models.TCPResponseRules{}

	data, err := p.Get(parser.Backends, backend, "tcp-response", false)
	if err != nil {
		if errors.Is(err, parser_errors.ErrFetch) {
			return tcpResRules, nil
		}
		return nil, err
	}

	tRules := data.([]types.TCPType)
	for i, tRule := range tRules {
		id := int64(i)
		tcpResRule := ParseTCPResponseRule(tRule)
		if tcpResRule != nil {
			tcpResRule.Index = &id
			tcpResRules = append(tcpResRules, tcpResRule)
		}
	}
	return tcpResRules, nil
}

func ParseTCPResponseRule(t types.TCPType) *models.TCPResponseRule {
	switch v := t.(type) {
	case *tcp_types.InspectDelay:
		return &models.TCPResponseRule{
			Type:    models.TCPResponseRuleTypeInspectDelay,
			Timeout: misc.ParseTimeout(v.Timeout),
		}
	case *tcp_types.Content:
		switch a := v.Action.(type) {
		case *tcp_actions.Accept:
			return &models.TCPResponseRule{
				Type:     models.TCPResponseRuleTypeContent,
				Action:   a.String(),
				Cond:     v.Cond,
				CondTest: v.CondTest,
			}
		case *tcp_actions.Reject:
			return &models.TCPResponseRule{
				Type:     models.TCPResponseRuleTypeContent,
				Action:   a.String(),
				Cond:     v.Cond,
				CondTest: v.CondTest,
			}
		case *tcp_actions.Lua:
			return &models.TCPResponseRule{
				Type:      models.TCPResponseRuleTypeContent,
				Action:    models.TCPResponseRuleActionLua,
				LuaAction: a.Action,
				LuaParams: a.Params,
				Cond:      v.Cond,
				CondTest:  v.CondTest,
			}
		}
	}
	return nil
}

func SerializeTCPResponseRule(t models.TCPResponseRule) types.TCPType {
	switch t.Type {
	case models.TCPResponseRuleTypeContent:
		switch t.Action {
		case models.TCPResponseRuleActionAccept:
			return &tcp_types.Content{
				Action:   &tcp_actions.Accept{},
				Cond:     t.Cond,
				CondTest: t.CondTest,
			}
		case models.TCPResponseRuleActionReject:
			return &tcp_types.Content{
				Action:   &tcp_actions.Reject{},
				Cond:     t.Cond,
				CondTest: t.CondTest,
			}
		case models.TCPResponseRuleActionLua:
			return &tcp_types.Content{
				Action: &tcp_actions.Lua{
					Action: t.LuaAction,
					Params: t.LuaParams,
				},
				Cond:     t.Cond,
				CondTest: t.CondTest,
			}
		}
	case models.TCPResponseRuleTypeInspectDelay:
		if t.Timeout != nil {
			return &tcp_types.InspectDelay{
				Timeout: strconv.FormatInt(*t.Timeout, 10),
			}
		}
	}

	return nil
}
