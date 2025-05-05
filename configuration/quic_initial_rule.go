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
	goerrors "errors"
	"strconv"

	"github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/client-native/v6/config-parser"
	parser_errors "github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/parsers/actions"
	quic_actions "github.com/haproxytech/client-native/v6/config-parser/parsers/quic/actions"
	"github.com/haproxytech/client-native/v6/config-parser/types"

	"github.com/haproxytech/client-native/v6/misc"
	"github.com/haproxytech/client-native/v6/models"
)

type QUICInitialRule interface {
	GetQUICInitialRules(parentType, parentName string, transactionID string) (int64, models.QUICInitialRules, error)
	GetQUICInitialRule(id int64, parentType, parentName string, transactionID string) (int64, *models.QUICInitialRule, error)
	DeleteQUICInitialRule(id int64, parentType string, parentName string, transactionID string, version int64) error
	CreateQUICInitialRule(id int64, parentType string, parentName string, data *models.QUICInitialRule, transactionID string, version int64) error
	EditQUICInitialRule(id int64, parentType string, parentName string, data *models.QUICInitialRule, transactionID string, version int64) error
	ReplaceQUICInitialRules(parentType string, parentName string, data models.QUICInitialRules, transactionID string, version int64) error
}

// GetQUICInitialRules returns configuration version and an array of configured quic initial rules in the specified parent.
// Returns error on fail.
func (c *client) GetQUICInitialRules(parentType, parentName string, transactionID string) (int64, models.QUICInitialRules, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	httpRules, err := ParseQUICInitialRules(parentType, parentName, p)
	if err != nil {
		return v, nil, c.HandleError("", parentType, parentName, "", false, err)
	}

	return v, httpRules, nil
}

// GetQUICInitialRule returns configuration version and a response quic initial rule in the specified parent.
// Returns error on fail or if quic initial rule does not exist.
func (c *client) GetQUICInitialRule(id int64, parentType, parentName string, transactionID string) (int64, *models.QUICInitialRule, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	section, parentName, err := getParserFromParent("quic-initial", parentType, parentName)
	if err != nil {
		return v, nil, err
	}

	data, err := p.GetOne(section, parentName, "quic-initial", int(id))
	if err != nil {
		return v, nil, c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, "", false, err)
	}

	httpRule, err := ParseQUICInitialRule(data.(types.Action))
	if err != nil {
		return v, nil, err
	}

	return v, httpRule, nil
}

// DeleteQUICInitialRule deletes a quic initial rule in configuration. One of version or transactionID is mandatory.
// Returns error on fail, nil on success.
func (c *client) DeleteQUICInitialRule(id int64, parentType string, parentName string, transactionID string, version int64) error {
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	section, parentName, err := getParserFromParent("quic-initial", parentType, parentName)
	if err != nil {
		return err
	}

	if err := p.Delete(section, parentName, "quic-initial", int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

// CreateQUICInitialRule creates a quic initial rule in configuration. One of version or transactionID is mandatory.
// Returns error on fail, nil on success.
func (c *client) CreateQUICInitialRule(id int64, parentType string, parentName string, data *models.QUICInitialRule, transactionID string, version int64) error {
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

	section, parentName, err := getParserFromParent("quic-initial", parentType, parentName)
	if err != nil {
		return err
	}

	s, err := SerializeQUICInitialRule(*data)
	if err != nil {
		return err
	}
	if err := p.Insert(section, parentName, "quic-initial", s, int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

// EditQUICInitialRule edits a quic initial rule in configuration. One of version or transactionID is mandatory.
// Returns error on fail, nil on success.
func (c *client) EditQUICInitialRule(id int64, parentType string, parentName string, data *models.QUICInitialRule, transactionID string, version int64) error {
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

	section, parentName, err := getParserFromParent("quic-initial", parentType, parentName)
	if err != nil {
		return err
	}

	if _, err = p.GetOne(section, parentName, "quic-initial", int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	s, err := SerializeQUICInitialRule(*data)
	if err != nil {
		return err
	}
	if err := p.Set(section, parentName, "quic-initial", s, int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

// ReplaceQUICInitialRules replaces all quic initial rules lines in configuration for a parentType/parentName.
// One of version or transactionID is mandatory.
// Returns error on fail, nil on success.
//
//nolint:dupl
func (c *client) ReplaceQUICInitialRules(parentType string, parentName string, data models.QUICInitialRules, transactionID string, version int64) error {
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

	section, parentName, err := getParserFromParent("quic-initial", parentType, parentName)
	if err != nil {
		return err
	}

	qiRules, err := ParseQUICInitialRules(parentType, parentName, p)
	if err != nil {
		return c.HandleError("", parentType, parentName, "", false, err)
	}

	for i := range qiRules {
		// Always delete index 0
		if err := p.Delete(section, parentName, "quic-initial", 0); err != nil {
			return c.HandleError(strconv.FormatInt(int64(i), 10), parentType, parentName, t, transactionID == "", err)
		}
	}

	for i, newQiRule := range data {
		s, err := SerializeQUICInitialRule(*newQiRule)
		if err != nil {
			return err
		}
		if err := p.Insert(section, parentName, "quic-initial", s, i); err != nil {
			return c.HandleError(strconv.FormatInt(int64(i), 10), parentType, parentName, t, transactionID == "", err)
		}
	}

	return c.SaveData(p, t, transactionID == "")
}

func ParseQUICInitialRules(t, pName string, p parser.Parser) (models.QUICInitialRules, error) {
	section, pName, err := getParserFromParent("quic-initial", t, pName)
	if err != nil {
		return nil, err
	}

	var httpResRules models.QUICInitialRules
	data, err := p.Get(section, pName, "quic-initial", false)
	if err != nil {
		if goerrors.Is(err, parser_errors.ErrFetch) {
			return httpResRules, nil
		}
		return nil, err
	}

	rules, ok := data.([]types.Action)
	if !ok {
		return nil, misc.CreateTypeAssertError("quic-initial")
	}
	for _, r := range rules {
		httpResRule, err := ParseQUICInitialRule(r)
		if err == nil && httpResRule != nil {
			httpResRules = append(httpResRules, httpResRule)
		}
	}
	return httpResRules, nil
}

func ParseQUICInitialRule(f types.Action) (*models.QUICInitialRule, error) {
	switch v := f.(type) {
	case *actions.Accept:
		return &models.QUICInitialRule{
			Type:     "accept",
			Cond:     v.Cond,
			CondTest: v.CondTest,
			Metadata: parseMetadata(v.Comment),
		}, nil
	case *actions.Reject:
		return &models.QUICInitialRule{
			Type:     "reject",
			Cond:     v.Cond,
			CondTest: v.CondTest,
			Metadata: parseMetadata(v.Comment),
		}, nil
	case *quic_actions.DgramDrop:
		return &models.QUICInitialRule{
			Type:     "dgram-drop",
			Cond:     v.Cond,
			CondTest: v.CondTest,
			Metadata: parseMetadata(v.Comment),
		}, nil
	case *quic_actions.SendRetry:
		return &models.QUICInitialRule{
			Type:     "send-retry",
			Cond:     v.Cond,
			CondTest: v.CondTest,
			Metadata: parseMetadata(v.Comment),
		}, nil
	}
	return nil, nil //nolint:nilnil
}

func SerializeQUICInitialRule(f models.QUICInitialRule) (types.Action, error) {
	comment, err := serializeMetadata(f.Metadata)
	if err != nil {
		return nil, err
	}
	var rule types.Action
	switch f.Type {
	case "accept":
		rule = &actions.Accept{
			Cond:     f.Cond,
			CondTest: f.CondTest,
			Comment:  comment,
		}
	case "reject":
		rule = &actions.Reject{
			Cond:     f.Cond,
			CondTest: f.CondTest,
			Comment:  comment,
		}
	case "dgram-drop":
		rule = &quic_actions.DgramDrop{
			Cond:     f.Cond,
			CondTest: f.CondTest,
			Comment:  comment,
		}
	case "send-retry":
		rule = &quic_actions.SendRetry{
			Cond:     f.Cond,
			CondTest: f.CondTest,
			Comment:  comment,
		}
	}
	return rule, nil
}
