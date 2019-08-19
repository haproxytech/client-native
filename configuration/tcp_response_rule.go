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

	strfmt "github.com/go-openapi/strfmt"
	"github.com/haproxytech/client-native/misc"
	parser "github.com/haproxytech/config-parser"
	parser_errors "github.com/haproxytech/config-parser/errors"
	"github.com/haproxytech/config-parser/parsers/tcp/actions"
	"github.com/haproxytech/config-parser/types"
	"github.com/haproxytech/models"
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

	tcpRules, err := c.parseTCPResponseRules(backend, p)
	if err != nil {
		return v, nil, c.handleError("", "backend", backend, "", false, err)
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
		return v, nil, c.handleError(strconv.FormatInt(id, 10), "backend", backend, "", false, err)
	}

	tcpRule := parseTCPResponseRule(data.(types.TCPAction))
	tcpRule.ID = &id

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
		return c.handleError(strconv.FormatInt(id, 10), "backend", backend, t, transactionID == "", err)
	}

	if err := c.saveData(p, t, transactionID == ""); err != nil {
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

	if err := p.Insert(parser.Backends, backend, "tcp-response", serializeTCPResponseRule(*data), int(*data.ID)); err != nil {
		return c.handleError(strconv.FormatInt(*data.ID, 10), "backend", backend, t, transactionID == "", err)
	}

	if err := c.saveData(p, t, transactionID == ""); err != nil {
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
		return c.handleError(strconv.FormatInt(*data.ID, 10), "backend", backend, t, transactionID == "", err)
	}

	if err := p.Set(parser.Backends, backend, "tcp-response", serializeTCPResponseRule(*data), int(id)); err != nil {
		return c.handleError(strconv.FormatInt(*data.ID, 10), "backend", backend, t, transactionID == "", err)
	}

	if err := c.saveData(p, t, transactionID == ""); err != nil {
		return err
	}
	return nil
}

func (c *Client) parseTCPResponseRules(backend string, p *parser.Parser) (models.TCPResponseRules, error) {
	tcpResRules := models.TCPResponseRules{}

	data, err := p.Get(parser.Backends, backend, "tcp-response", false)
	if err != nil {
		if err == parser_errors.ErrFetch {
			return tcpResRules, nil
		}
		return nil, err
	}

	tRules := data.([]types.TCPAction)
	for i, tRule := range tRules {
		id := int64(i)
		tcpResRule := parseTCPResponseRule(tRule)
		if tcpResRule != nil {
			tcpResRule.ID = &id
			tcpResRules = append(tcpResRules, tcpResRule)
		}
	}
	return tcpResRules, nil
}

func parseTCPResponseRule(t types.TCPAction) *models.TCPResponseRule {
	switch v := t.(type) {
	case *actions.Content:
		r := &models.TCPResponseRule{
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
		return &models.TCPResponseRule{
			Type:    "inspect-delay",
			Timeout: misc.ParseTimeout(v.Timeout),
		}
	}
	return nil
}

func serializeTCPResponseRule(t models.TCPResponseRule) types.TCPAction {
	switch t.Type {
	case "content":
		return &actions.Content{
			Action:   []string{t.Action},
			Cond:     t.Cond,
			CondTest: t.CondTest,
		}
	case "inspect-delay":
		if t.Timeout != nil {
			return &actions.InspectDelay{
				Timeout: strconv.FormatInt(*t.Timeout, 10),
			}
		}
	}
	return nil
}
