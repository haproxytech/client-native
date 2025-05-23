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
	"strings"

	"github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/client-native/v6/config-parser"
	"github.com/haproxytech/client-native/v6/config-parser/common"
	parser_errors "github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/parsers/actions"
	tcp_actions "github.com/haproxytech/client-native/v6/config-parser/parsers/tcp/actions"
	tcp_types "github.com/haproxytech/client-native/v6/config-parser/parsers/tcp/types"
	"github.com/haproxytech/client-native/v6/config-parser/types"

	"github.com/haproxytech/client-native/v6/configuration/options"
	"github.com/haproxytech/client-native/v6/misc"
	"github.com/haproxytech/client-native/v6/models"
)

type TCPResponseRule interface {
	GetTCPResponseRules(parentType, parentName string, transactionID string) (int64, models.TCPResponseRules, error)
	GetTCPResponseRule(id int64, parentType, parentName string, transactionID string) (int64, *models.TCPResponseRule, error)
	DeleteTCPResponseRule(id int64, parentType, parentName string, transactionID string, version int64) error
	CreateTCPResponseRule(id int64, parentType, parentName string, data *models.TCPResponseRule, transactionID string, version int64) error
	EditTCPResponseRule(id int64, parentType, parentName string, data *models.TCPResponseRule, transactionID string, version int64) error
	ReplaceTCPResponseRules(parentType, parentName string, data models.TCPResponseRules, transactionID string, version int64) error
}

// GetTCPResponseRules returns configuration version and an array of
// configured tcp response rules in the specified backend. Returns error on fail.
func (c *client) GetTCPResponseRules(parentType, parentName string, transactionID string) (int64, models.TCPResponseRules, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	tcpRules, err := ParseTCPResponseRules(parentType, parentName, p)
	if err != nil {
		return v, nil, c.HandleError("", parentType, parentName, "", false, err)
	}

	return v, tcpRules, nil
}

// GetTCPResponseRule returns configuration version and a requested tcp response rule
// in the specified backend. Returns error on fail or if tcp response rule does not exist.
func (c *client) GetTCPResponseRule(id int64, parentType, parentName string, transactionID string) (int64, *models.TCPResponseRule, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	section, parentName, err := getParserFromParent("tcp-response", parentType, parentName)
	if err != nil {
		return v, nil, err
	}

	data, err := p.GetOne(section, parentName, "tcp-response", int(id))
	if err != nil {
		return v, nil, c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, "", false, err)
	}

	tcpRule, parseErr := ParseTCPResponseRule(data.(types.TCPType))
	if parseErr != nil {
		return 0, nil, parseErr
	}

	return v, tcpRule, nil
}

// DeleteTCPResponseRule deletes a tcp response rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) DeleteTCPResponseRule(id int64, parentType, parentName string, transactionID string, version int64) error {
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}
	section, parentName, err := getParserFromParent("tcp-response", parentType, parentName)
	if err != nil {
		return err
	}

	if err := p.Delete(section, parentName, "tcp-response", int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

// CreateTCPResponseRule creates a tcp response rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) CreateTCPResponseRule(id int64, parentType, parentName string, data *models.TCPResponseRule, transactionID string, version int64) error {
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

	section, parentName, err := getParserFromParent("tcp-response", parentType, parentName)
	if err != nil {
		return err
	}

	tcpRule, serializeErr := SerializeTCPResponseRule(*data, &c.ConfigurationOptions)
	if serializeErr != nil {
		return serializeErr
	}
	if err := p.Insert(section, parentName, "tcp-response", tcpRule, int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

// EditTCPResponseRule edits a tcp response rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) EditTCPResponseRule(id int64, parentType, parentName string, data *models.TCPResponseRule, transactionID string, version int64) error {
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
	section, parentName, err := getParserFromParent("tcp-response", parentType, parentName)
	if err != nil {
		return err
	}

	if _, err := p.GetOne(section, parentName, "tcp-response", int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	tcpRule, serializeErr := SerializeTCPResponseRule(*data, &c.ConfigurationOptions)
	if serializeErr != nil {
		return serializeErr
	}
	if err := p.Set(section, parentName, "tcp-response", tcpRule, int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

// ReplaceTCPResponseRules replaces all TCP Response Rule lines in configuration for a parentType/parentName.
// One of version or transactionID is mandatory.
// Returns error on fail, nil on success.
func (c *client) ReplaceTCPResponseRules(parentType, parentName string, data models.TCPResponseRules, transactionID string, version int64) error {
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
	section, parentName, err := getParserFromParent("tcp-response", parentType, parentName)
	if err != nil {
		return err
	}

	tcpResponseRules, err := ParseTCPResponseRules(parentType, parentName, p)
	if err != nil {
		return c.HandleError("", parentType, parentName, "", false, err)
	}

	for i := range tcpResponseRules {
		// Always delete index 0
		if err := p.Delete(section, parentName, "tcp-response", 0); err != nil {
			return c.HandleError(strconv.FormatInt(int64(i), 10), parentType, parentName, t, transactionID == "", err)
		}
	}

	for i, newTCPResponseRule := range data {
		s, err := SerializeTCPResponseRule(*newTCPResponseRule, &c.ConfigurationOptions)
		if err != nil {
			return err
		}
		if err := p.Insert(section, parentName, "tcp-response", s, i); err != nil {
			return c.HandleError(strconv.FormatInt(int64(i), 10), parentType, parentName, t, transactionID == "", err)
		}
	}

	return c.SaveData(p, t, transactionID == "")
}

func ParseTCPResponseRules(parentType, parentName string, p parser.Parser) (models.TCPResponseRules, error) {
	section, pName, err := getParserFromParent("tcp-response", parentType, parentName)
	if err != nil {
		return nil, err
	}

	var tcpResRules models.TCPResponseRules

	data, err := p.Get(section, pName, "tcp-response", false)
	if err != nil {
		if errors.Is(err, parser_errors.ErrFetch) {
			return tcpResRules, nil
		}
		return nil, err
	}

	tRules, ok := data.([]types.TCPType)
	if !ok {
		return nil, misc.CreateTypeAssertError("tcp response")
	}
	for _, tRule := range tRules {
		tcpResRule, parseErr := ParseTCPResponseRule(tRule)
		if parseErr != nil {
			return nil, parseErr
		}
		if tcpResRule != nil {
			tcpResRules = append(tcpResRules, tcpResRule)
		}
	}
	return tcpResRules, nil
}

//nolint:maintidx,gocognit
func ParseTCPResponseRule(t types.TCPType) (*models.TCPResponseRule, error) {
	switch v := t.(type) {
	case *tcp_types.InspectDelay:
		return &models.TCPResponseRule{
			Type:     models.TCPResponseRuleTypeInspectDashDelay,
			Timeout:  misc.ParseTimeout(v.Timeout),
			Metadata: parseMetadata(v.Comment),
		}, nil
	case *tcp_types.Content:
		switch a := v.Action.(type) {
		case *actions.Accept:
			return &models.TCPResponseRule{
				Type:     models.TCPResponseRuleTypeContent,
				Action:   models.TCPResponseRuleActionAccept,
				Cond:     a.Cond,
				CondTest: a.CondTest,
				Metadata: parseMetadata(v.Comment),
			}, nil
		case *actions.Reject:
			return &models.TCPResponseRule{
				Type:     models.TCPResponseRuleTypeContent,
				Action:   models.TCPResponseRuleActionReject,
				Cond:     a.Cond,
				CondTest: a.CondTest,
				Metadata: parseMetadata(v.Comment),
			}, nil
		case *actions.Lua:
			return &models.TCPResponseRule{
				Type:      models.TCPResponseRuleTypeContent,
				Action:    models.TCPResponseRuleActionLua,
				LuaAction: a.Action,
				LuaParams: a.Params,
				Cond:      a.Cond,
				CondTest:  a.CondTest,
				Metadata:  parseMetadata(v.Comment),
			}, nil
		case *actions.SetBandwidthLimit:
			return &models.TCPResponseRule{
				Type:                 models.TCPResponseRuleTypeContent,
				Action:               models.TCPResponseRuleActionSetDashBandwidthDashLimit,
				BandwidthLimitName:   a.Name,
				BandwidthLimitLimit:  a.Limit.String(),
				BandwidthLimitPeriod: a.Period.String(),
				Cond:                 a.Cond,
				CondTest:             a.CondTest,
				Metadata:             parseMetadata(v.Comment),
			}, nil
		case *tcp_actions.Close:
			return &models.TCPResponseRule{
				Type:     models.TCPResponseRuleTypeContent,
				Action:   models.TCPResponseRuleActionClose,
				Cond:     a.Cond,
				CondTest: a.CondTest,
				Metadata: parseMetadata(v.Comment),
			}, nil
		case *actions.ScAddGpc:
			if a.Int == nil && len(a.Expr.Expr) == 0 {
				return nil, NewConfError(ErrValidationError, "sc-add-gpc int or expr has to be set")
			}
			if a.Int != nil && len(a.Expr.Expr) > 0 {
				return nil, NewConfError(ErrValidationError, "sc-add-gpc int and expr are exclusive")
			}
			ID, _ := strconv.ParseInt(a.ID, 10, 64)
			Idx, _ := strconv.ParseInt(a.Idx, 10, 64)
			return &models.TCPResponseRule{
				Type:     models.TCPResponseRuleTypeContent,
				Action:   models.TCPResponseRuleActionScDashAddDashGpc,
				ScID:     ID,
				ScIdx:    Idx,
				Expr:     strings.Join(a.Expr.Expr, " "),
				ScInt:    a.Int,
				Cond:     a.Cond,
				CondTest: a.CondTest,
				Metadata: parseMetadata(v.Comment),
			}, nil
		case *actions.ScIncGpc:
			ID, _ := strconv.ParseInt(a.ID, 10, 64)
			Idx, _ := strconv.ParseInt(a.Idx, 10, 64)
			return &models.TCPResponseRule{
				Type:     models.TCPResponseRuleTypeContent,
				Action:   models.TCPResponseRuleActionScDashIncDashGpc,
				ScID:     ID,
				ScIdx:    Idx,
				Cond:     a.Cond,
				CondTest: a.CondTest,
				Metadata: parseMetadata(v.Comment),
			}, nil
		case *actions.ScIncGpc0:
			ID, _ := strconv.ParseInt(a.ID, 10, 64)
			return &models.TCPResponseRule{
				Type:     models.TCPResponseRuleTypeContent,
				Action:   models.TCPResponseRuleActionScDashIncDashGpc0,
				ScID:     ID,
				Cond:     a.Cond,
				CondTest: a.CondTest,
				Metadata: parseMetadata(v.Comment),
			}, nil
		case *actions.ScIncGpc1:
			ID, _ := strconv.ParseInt(a.ID, 10, 64)
			return &models.TCPResponseRule{
				Type:     models.TCPResponseRuleTypeContent,
				Action:   models.TCPResponseRuleActionScDashIncDashGpc1,
				ScID:     ID,
				Cond:     a.Cond,
				CondTest: a.CondTest,
				Metadata: parseMetadata(v.Comment),
			}, nil
		case *actions.ScSetGpt:
			if a.Int == nil && len(a.Expr.Expr) == 0 {
				return nil, NewConfError(ErrValidationError, "sc-set-gpt: int or expr has to be set")
			}
			if a.Int != nil && len(a.Expr.Expr) > 0 {
				return nil, NewConfError(ErrValidationError, "sc-set-gpt: int and expr are exclusive")
			}
			scID, err := strconv.ParseInt(a.ScID, 10, 64)
			if err != nil {
				return nil, NewConfError(ErrValidationError, "sc-set-gpt: failed to parse sc-id as an int")
			}
			return &models.TCPResponseRule{
				Type:     models.TCPResponseRuleTypeContent,
				Action:   models.TCPResponseRuleActionScDashSetDashGpt,
				ScID:     scID,
				ScIdx:    a.Idx,
				Expr:     strings.Join(a.Expr.Expr, " "),
				ScInt:    a.Int,
				Cond:     a.Cond,
				CondTest: a.CondTest,
				Metadata: parseMetadata(v.Comment),
			}, nil
		case *actions.ScSetGpt0:
			if a.Int == nil && len(a.Expr.Expr) == 0 {
				return nil, NewConfError(ErrValidationError, "sc-set-gpt0 int or expr has to be set")
			}
			if a.Int != nil && len(a.Expr.Expr) > 0 {
				return nil, NewConfError(ErrValidationError, "sc-set-gpt0 int and expr are exclusive")
			}
			ID, _ := strconv.ParseInt(a.ID, 10, 64)
			return &models.TCPResponseRule{
				Type:     models.TCPResponseRuleTypeContent,
				Action:   models.TCPResponseRuleActionScDashSetDashGpt0,
				ScID:     ID,
				Expr:     strings.Join(a.Expr.Expr, " "),
				ScInt:    a.Int,
				Cond:     a.Cond,
				CondTest: a.CondTest,
				Metadata: parseMetadata(v.Comment),
			}, nil
		case *actions.SendSpoeGroup:
			return &models.TCPResponseRule{
				Type:       models.TCPResponseRuleTypeContent,
				Action:     models.TCPResponseRuleActionSendDashSpoeDashGroup,
				SpoeEngine: a.Engine,
				SpoeGroup:  a.Group,
				Cond:       a.Cond,
				CondTest:   a.CondTest,
				Metadata:   parseMetadata(v.Comment),
			}, nil
		case *actions.SetLogLevel:
			return &models.TCPResponseRule{
				Type:     models.TCPResponseRuleTypeContent,
				Action:   models.TCPResponseRuleActionSetDashLogDashLevel,
				LogLevel: a.Level,
				Cond:     a.Cond,
				CondTest: a.CondTest,
				Metadata: parseMetadata(v.Comment),
			}, nil
		case *actions.SetMark:
			return &models.TCPResponseRule{
				Type:      models.TCPResponseRuleTypeContent,
				Action:    models.TCPResponseRuleActionSetDashMark,
				MarkValue: a.Value,
				Cond:      a.Cond,
				CondTest:  a.CondTest,
				Metadata:  parseMetadata(v.Comment),
			}, nil
		case *actions.SetNice:
			nice, err := strconv.ParseInt(a.Value, 10, 64)
			if err != nil {
				return nil, err
			}
			return &models.TCPResponseRule{
				Type:      models.TCPResponseRuleTypeContent,
				Action:    models.TCPResponseRuleActionSetDashNice,
				NiceValue: nice,
				Cond:      a.Cond,
				CondTest:  a.CondTest,
				Metadata:  parseMetadata(v.Comment),
			}, nil
		case *actions.SetTos:
			return &models.TCPResponseRule{
				Type:     models.TCPResponseRuleTypeContent,
				Action:   models.TCPResponseRuleActionSetDashTos,
				TosValue: a.Value,
				Cond:     a.Cond,
				CondTest: a.CondTest,
				Metadata: parseMetadata(v.Comment),
			}, nil
		case *actions.SetFcMark:
			return &models.TCPResponseRule{
				Type:     models.TCPResponseRuleTypeContent,
				Action:   models.TCPResponseRuleActionSetDashFcDashMark,
				Expr:     a.Expr.String(),
				Cond:     a.Cond,
				CondTest: a.CondTest,
				Metadata: parseMetadata(v.Comment),
			}, nil
		case *actions.SetFcTos:
			return &models.TCPResponseRule{
				Type:     models.TCPResponseRuleTypeContent,
				Action:   models.TCPResponseRuleActionSetDashFcDashTos,
				Expr:     a.Expr.String(),
				Cond:     a.Cond,
				CondTest: a.CondTest,
				Metadata: parseMetadata(v.Comment),
			}, nil
		case *actions.SilentDrop:
			return &models.TCPResponseRule{
				Type:     models.TCPResponseRuleTypeContent,
				Action:   models.TCPResponseRuleActionSilentDashDrop,
				RstTTL:   a.RstTTL,
				Cond:     a.Cond,
				CondTest: a.CondTest,
				Metadata: parseMetadata(v.Comment),
			}, nil
		case *actions.SetVar:
			return &models.TCPResponseRule{
				Action:   models.TCPResponseRuleActionSetDashVar,
				VarScope: a.VarScope,
				VarName:  a.VarName,
				Expr:     a.Expr.String(),
				Cond:     a.Cond,
				CondTest: a.CondTest,
				Metadata: parseMetadata(v.Comment),
			}, nil
		case *actions.SetVarFmt:
			return &models.TCPResponseRule{
				Action:    models.TCPResponseRuleActionSetDashVarDashFmt,
				VarName:   a.VarName,
				VarFormat: strings.Join(a.Fmt.Expr, " "),
				VarScope:  a.VarScope,
				Cond:      a.Cond,
				CondTest:  a.CondTest,
				Metadata:  parseMetadata(v.Comment),
			}, nil
		case *actions.UnsetVar:
			return &models.TCPResponseRule{
				Type:     models.TCPResponseRuleTypeContent,
				Action:   models.TCPResponseRuleActionUnsetDashVar,
				VarName:  a.Name,
				VarScope: a.Scope,
				Cond:     a.Cond,
				CondTest: a.CondTest,
				Metadata: parseMetadata(v.Comment),
			}, nil
		case *actions.DoLog:
			return &models.TCPResponseRule{
				Type:     models.TCPResponseRuleTypeContent,
				Action:   models.TCPResponseRuleActionDoDashLog,
				Cond:     a.Cond,
				CondTest: a.CondTest,
				Metadata: parseMetadata(v.Comment),
			}, nil
		}
	}
	return nil, NewConfError(ErrValidationError, "invalid action")
}

func SerializeTCPResponseRule(t models.TCPResponseRule, opt *options.ConfigurationOptions) (types.TCPType, error) { //nolint:maintidx
	comment, err := serializeMetadata(t.Metadata)
	if err != nil {
		comment = ""
	}
	switch t.Type {
	case models.TCPResponseRuleTypeContent:
		switch t.Action {
		case models.TCPResponseRuleActionAccept:
			return &tcp_types.Content{
				Action: &actions.Accept{
					Cond:     t.Cond,
					CondTest: t.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPResponseRuleActionReject:
			return &tcp_types.Content{
				Action: &actions.Reject{
					Cond:     t.Cond,
					CondTest: t.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPResponseRuleActionLua:
			return &tcp_types.Content{
				Action: &actions.Lua{
					Action:   t.LuaAction,
					Params:   t.LuaParams,
					Cond:     t.Cond,
					CondTest: t.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPResponseRuleActionSetDashBandwidthDashLimit:
			return &tcp_types.Content{
				Action: &actions.SetBandwidthLimit{
					Name:     t.BandwidthLimitName,
					Limit:    common.Expression{Expr: strings.Split(t.BandwidthLimitLimit, " ")},
					Period:   common.Expression{Expr: strings.Split(t.BandwidthLimitPeriod, " ")},
					Cond:     t.Cond,
					CondTest: t.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPResponseRuleActionClose:
			return &tcp_types.Content{
				Action: &tcp_actions.Close{
					Cond:     t.Cond,
					CondTest: t.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPResponseRuleActionScDashAddDashGpc:
			if len(t.Expr) > 0 && t.ScInt != nil {
				return nil, NewConfError(ErrValidationError, "sc-add-gpc int and expr are exclusive")
			}
			if len(t.Expr) == 0 && t.ScInt == nil {
				return nil, NewConfError(ErrValidationError, "sc-add-gpc int or expr has to be set")
			}
			return &tcp_types.Content{
				Action: &actions.ScAddGpc{
					ID:       strconv.FormatInt(t.ScID, 10),
					Idx:      strconv.FormatInt(t.ScIdx, 10),
					Int:      t.ScInt,
					Expr:     common.Expression{Expr: strings.Split(t.Expr, " ")},
					Cond:     t.Cond,
					CondTest: t.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPResponseRuleActionScDashIncDashGpc:
			return &tcp_types.Content{
				Action: &actions.ScIncGpc{
					ID:       strconv.FormatInt(t.ScID, 10),
					Idx:      strconv.FormatInt(t.ScIdx, 10),
					Cond:     t.Cond,
					CondTest: t.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPResponseRuleActionScDashIncDashGpc0:
			return &tcp_types.Content{
				Action: &actions.ScIncGpc0{
					ID:       strconv.FormatInt(t.ScID, 10),
					Cond:     t.Cond,
					CondTest: t.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPResponseRuleActionScDashIncDashGpc1:
			return &tcp_types.Content{
				Action: &actions.ScIncGpc1{
					ID:       strconv.FormatInt(t.ScID, 10),
					Cond:     t.Cond,
					CondTest: t.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPResponseRuleActionScDashSetDashGpt:
			if len(t.Expr) > 0 && t.ScInt != nil {
				return nil, NewConfError(ErrValidationError, "sc-set-gpt: int and expr are exclusive")
			}
			if len(t.Expr) == 0 && t.ScInt == nil {
				return nil, NewConfError(ErrValidationError, "sc-set-gpt: int or expr has to be set")
			}
			return &tcp_types.Content{
				Action: &actions.ScSetGpt{
					ScID:     strconv.FormatInt(t.ScID, 10),
					Idx:      t.ScIdx,
					Int:      t.ScInt,
					Expr:     common.Expression{Expr: strings.Split(t.Expr, " ")},
					Cond:     t.Cond,
					CondTest: t.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPResponseRuleActionScDashSetDashGpt0:
			if len(t.Expr) > 0 && t.ScInt != nil {
				return nil, NewConfError(ErrValidationError, "sc-set-gpt0 int and expr are exclusive")
			}
			if len(t.Expr) == 0 && t.ScInt == nil {
				return nil, NewConfError(ErrValidationError, "sc-set-gpt0 int or expr has to be set")
			}
			return &tcp_types.Content{
				Action: &actions.ScSetGpt0{
					ID:       strconv.FormatInt(t.ScID, 10),
					Int:      t.ScInt,
					Expr:     common.Expression{Expr: strings.Split(t.Expr, " ")},
					Cond:     t.Cond,
					CondTest: t.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPResponseRuleActionSendDashSpoeDashGroup:
			return &tcp_types.Content{
				Action: &actions.SendSpoeGroup{
					Engine:   t.SpoeEngine,
					Group:    t.SpoeGroup,
					Cond:     t.Cond,
					CondTest: t.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPResponseRuleActionSetDashLogDashLevel:
			return &tcp_types.Content{
				Action: &actions.SetLogLevel{
					Level:    t.LogLevel,
					Cond:     t.Cond,
					CondTest: t.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPResponseRuleActionSetDashMark:
			return &tcp_types.Content{
				Action: &actions.SetMark{
					Value:    t.MarkValue,
					Cond:     t.Cond,
					CondTest: t.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPResponseRuleActionSetDashNice:
			return &tcp_types.Content{
				Action: &actions.SetNice{
					Value:    strconv.FormatInt(t.NiceValue, 10),
					Cond:     t.Cond,
					CondTest: t.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPResponseRuleActionSetDashTos:
			return &tcp_types.Content{
				Action: &actions.SetTos{
					Value:    t.TosValue,
					Cond:     t.Cond,
					CondTest: t.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPResponseRuleActionSilentDashDrop:
			return &tcp_types.Content{
				Action: &actions.SilentDrop{
					RstTTL:   t.RstTTL,
					Cond:     t.Cond,
					CondTest: t.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPResponseRuleActionSetDashVarDashFmt:
			return &tcp_types.Content{
				Action: &actions.SetVarFmt{
					Fmt:      common.Expression{Expr: strings.Split(t.VarFormat, " ")},
					VarName:  t.VarName,
					VarScope: t.VarScope,
					Cond:     t.Cond,
					CondTest: t.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPResponseRuleActionSetDashVar:
			return &tcp_types.Content{
				Action: &actions.SetVar{
					VarName:  t.VarName,
					VarScope: t.VarScope,
					Expr:     common.Expression{Expr: strings.Split(t.Expr, " ")},
					Cond:     t.Cond,
					CondTest: t.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPResponseRuleActionUnsetDashVar:
			return &tcp_types.Content{
				Action: &actions.UnsetVar{
					Name:     t.VarName,
					Scope:    t.VarScope,
					Cond:     t.Cond,
					CondTest: t.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPResponseRuleActionSetDashFcDashMark:
			return &tcp_types.Content{
				Action: &actions.SetFcMark{
					Expr:     common.Expression{Expr: strings.Split(t.Expr+t.MarkValue, " ")},
					Cond:     t.Cond,
					CondTest: t.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPResponseRuleActionSetDashFcDashTos:
			return &tcp_types.Content{
				Action: &actions.SetFcTos{
					Expr:     common.Expression{Expr: strings.Split(t.Expr+t.TosValue, " ")},
					Cond:     t.Cond,
					CondTest: t.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPResponseRuleActionDoDashLog:
			return &tcp_types.Content{
				Action: &actions.DoLog{
					Cond:     t.Cond,
					CondTest: t.CondTest,
				},
				Comment: comment,
			}, nil
		}
	case models.TCPResponseRuleTypeInspectDashDelay:
		if t.Timeout != nil {
			return &tcp_types.InspectDelay{
				Timeout: misc.SerializeTime(*t.Timeout, opt.PreferredTimeSuffix),
				Comment: comment,
			}, nil
		}
	}

	return nil, NewConfError(ErrValidationError, "invalid action")
}
