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
	"fmt"
	"strconv"
	"strings"

	"github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/config-parser/v4"
	"github.com/haproxytech/config-parser/v4/common"
	parser_errors "github.com/haproxytech/config-parser/v4/errors"
	"github.com/haproxytech/config-parser/v4/parsers/actions"
	tcp_actions "github.com/haproxytech/config-parser/v4/parsers/tcp/actions"
	tcp_types "github.com/haproxytech/config-parser/v4/parsers/tcp/types"
	"github.com/haproxytech/config-parser/v4/types"

	"github.com/haproxytech/client-native/v4/misc"
	"github.com/haproxytech/client-native/v4/models"
)

type TCPRequestRule interface {
	GetTCPRequestRules(parentType, parentName string, transactionID string) (int64, models.TCPRequestRules, error)
	GetTCPRequestRule(id int64, parentType, parentName string, transactionID string) (int64, *models.TCPRequestRule, error)
	DeleteTCPRequestRule(id int64, parentType string, parentName string, transactionID string, version int64) error
	CreateTCPRequestRule(parentType string, parentName string, data *models.TCPRequestRule, transactionID string, version int64) error
	EditTCPRequestRule(id int64, parentType string, parentName string, data *models.TCPRequestRule, transactionID string, version int64) error
}

// GetTCPRequestRules returns configuration version and an array of
// configured TCP request rules in the specified parent. Returns error on fail.
func (c *client) GetTCPRequestRules(parentType, parentName string, transactionID string) (int64, models.TCPRequestRules, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	tcpRules, err := ParseTCPRequestRules(parentType, parentName, p)
	if err != nil {
		return v, nil, c.HandleError("", parentType, parentName, "", false, err)
	}

	return v, tcpRules, nil
}

// GetTCPRequestRule returns configuration version and a requested tcp request rule
// in the specified parent. Returns error on fail or if http request rule does not exist.
func (c *client) GetTCPRequestRule(id int64, parentType, parentName string, transactionID string) (int64, *models.TCPRequestRule, error) {
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
		return v, nil, c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, "", false, err)
	}

	tcpRule, err := ParseTCPRequestRule(data.(types.TCPType))
	if err != nil {
		return v, nil, err
	}

	tcpRule.Index = &id
	return v, tcpRule, nil
}

// DeleteTCPRequestRule deletes a tcp request rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) DeleteTCPRequestRule(id int64, parentType string, parentName string, transactionID string, version int64) error {
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
		return c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	if err := c.SaveData(p, t, transactionID == ""); err != nil {
		return err
	}

	return nil
}

// CreateTCPRequestRule creates a tcp request rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) CreateTCPRequestRule(parentType string, parentName string, data *models.TCPRequestRule, transactionID string, version int64) error {
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

	var section parser.Section
	if parentType == "backend" {
		section = parser.Backends
	} else if parentType == "frontend" {
		section = parser.Frontends
	}

	s, err := SerializeTCPRequestRule(*data)
	if err != nil {
		return err
	}

	if err := p.Insert(section, parentName, "tcp-request", s, int(*data.Index)); err != nil {
		return c.HandleError(strconv.FormatInt(*data.Index, 10), parentType, parentName, t, transactionID == "", err)
	}

	if err := c.SaveData(p, t, transactionID == ""); err != nil {
		return err
	}
	return nil
}

// EditTCPRequestRule edits a tcp request rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
// nolint:dupl
func (c *client) EditTCPRequestRule(id int64, parentType string, parentName string, data *models.TCPRequestRule, transactionID string, version int64) error {
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

	var section parser.Section
	if parentType == "backend" {
		section = parser.Backends
	} else if parentType == "frontend" {
		section = parser.Frontends
	}

	if _, err = p.GetOne(section, parentName, "tcp-request", int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	s, err := SerializeTCPRequestRule(*data)
	if err != nil {
		return err
	}

	if err := p.Set(section, parentName, "tcp-request", s, int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	if err := c.SaveData(p, t, transactionID == ""); err != nil {
		return err
	}
	return nil
}

func ParseTCPRequestRules(t, pName string, p parser.Parser) (models.TCPRequestRules, error) {
	section := parser.Global
	if t == "frontend" {
		section = parser.Frontends
	} else if t == "backend" {
		section = parser.Backends
	}

	tcpReqRules := models.TCPRequestRules{}
	data, err := p.Get(section, pName, "tcp-request", false)
	if err != nil {
		if errors.Is(err, parser_errors.ErrFetch) {
			return tcpReqRules, nil
		}
		return nil, err
	}

	rules, ok := data.([]types.TCPType)
	if !ok {
		return nil, misc.CreateTypeAssertError("tcp request")
	}
	for i, r := range rules {
		id := int64(i)
		tcpReqRule, err := ParseTCPRequestRule(r)
		if err == nil {
			tcpReqRule.Index = &id
			tcpReqRules = append(tcpReqRules, tcpReqRule)
		}
	}
	return tcpReqRules, nil
}

func ParseTCPRequestRule(f types.TCPType) (rule *models.TCPRequestRule, err error) { //nolint:gocyclo,cyclop,maintidx
	switch v := f.(type) {
	case *tcp_types.InspectDelay:
		return &models.TCPRequestRule{
			Type:    models.TCPRequestRuleTypeInspectDashDelay,
			Timeout: misc.ParseTimeout(v.Timeout),
		}, nil

	case *tcp_types.Connection:
		rule = &models.TCPRequestRule{
			Type: models.TCPRequestRuleTypeConnection,
		}

		switch a := v.Action.(type) {
		case *tcp_actions.Accept:
			rule.Action = models.TCPRequestRuleActionAccept
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.Reject:
			rule.Action = models.TCPRequestRuleActionReject
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *tcp_actions.ExpectProxy:
			rule.Action = models.TCPRequestRuleActionExpectDashProxy
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *tcp_actions.ExpectNetscalerCip:
			rule.Action = models.TCPRequestRuleActionExpectDashNetscalerDashCip
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *tcp_actions.Capture:
			rule.Action = models.TCPRequestRuleActionCapture
			rule.Expr = a.Expr.String()
			rule.CaptureLen = a.Len
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.TrackSc:
			rule.TrackKey = a.Key
			switch a.Type {
			case actions.TrackSc0:
				rule.Action = models.TCPRequestRuleActionTrackDashSc0
			case actions.TrackSc1:
				rule.Action = models.TCPRequestRuleActionTrackDashSc1
			case actions.TrackSc2:
				rule.Action = models.TCPRequestRuleActionTrackDashSc2
			}
			if a.Table != "" {
				rule.TrackTable = a.Table
			}
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.ScIncGpc0:
			rule.Action = models.TCPRequestRuleActionScDashSetDashGpt0
			rule.ScIncID = a.ID
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.ScIncGpc1:
			rule.Action = models.TCPRequestRuleActionScDashIncDashGpc1
			rule.ScIncID = a.ID
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.ScSetGpt0:
			rule.Action = models.TCPRequestRuleActionScDashSetDashGpt0
			rule.ScIncID = a.ID
			rule.Expr = a.Expr.String()
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *tcp_actions.SetSrc:
			rule.Action = models.TCPRequestRuleActionSetDashSrc
			rule.Expr = a.Expr.String()
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.SilentDrop:
			rule.Action = models.TCPRequestRuleActionSilentDashDrop
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.Lua:
			rule.Action = models.TCPRequestRuleActionLua
			rule.LuaAction = a.Action
			rule.LuaParams = a.Params
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		default:
			return nil, NewConfError(ErrValidationError, fmt.Sprintf("unsupported action '%T' in tcp_request_rule", a))
		}

		return rule, nil
	case *tcp_types.Content:
		rule = &models.TCPRequestRule{
			Type: models.TCPRequestRuleTypeContent,
		}

		switch a := v.Action.(type) {
		case *tcp_actions.Accept:
			rule.Action = models.TCPRequestRuleActionAccept
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.DoResolve:
			rule.Action = models.TCPRequestRuleActionDoDashResolve
			rule.VarName = a.Var
			rule.ResolveResolvers = a.Resolvers
			rule.ResolveProtocol = a.Protocol
			rule.Expr = a.Expr.String()
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.Reject:
			rule.Action = models.TCPRequestRuleActionReject
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *tcp_actions.Capture:
			rule.Action = models.TCPRequestRuleActionCapture
			rule.Expr = a.Expr.String()
			rule.CaptureLen = a.Len
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.SetPriorityClass:
			rule.Action = "set-priority-class"
			rule.Expr = a.Expr.String()
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.SetPriorityOffset:
			rule.Action = "set-priority-offset"
			rule.Expr = a.Expr.String()
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.TrackSc:
			switch a.Type {
			case actions.TrackSc0:
				rule.Action = models.TCPRequestRuleActionTrackDashSc0
			case actions.TrackSc1:
				rule.Action = models.TCPRequestRuleActionTrackDashSc1
			case actions.TrackSc2:
				rule.Action = models.TCPRequestRuleActionTrackDashSc2
			}
			rule.TrackKey = a.Key
			if a.Table != "" {
				rule.TrackTable = a.Table
			}
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.ScIncGpc0:
			rule.Action = models.TCPRequestRuleActionScDashIncDashGpc0
			rule.ScIncID = a.ID
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.ScIncGpc1:
			rule.Action = models.TCPRequestRuleActionScDashIncDashGpc1
			rule.ScIncID = a.ID
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.ScSetGpt0:
			rule.Action = models.TCPRequestRuleActionScDashSetDashGpt0
			rule.ScIncID = a.ID
			rule.GptValue = a.Expr.String()
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.SetDst:
			rule.Action = models.TCPRequestRuleActionSetDashDst
			rule.Expr = a.Expr.String()
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.SetDstPort:
			rule.Action = models.TCPRequestRuleActionSetDashDstDashPort
			rule.Expr = a.Expr.String()
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.SetVar:
			rule.Action = models.TCPRequestRuleActionSetDashVar
			rule.VarScope = a.VarScope
			rule.VarName = a.VarName
			rule.Expr = a.Expr.String()
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.UnsetVar:
			rule.Action = models.TCPRequestRuleActionUnsetDashVar
			rule.VarScope = a.Scope
			rule.VarName = a.Name
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.SilentDrop:
			rule.Action = models.TCPRequestRuleActionSilentDashDrop
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.SendSpoeGroup:
			rule.Action = models.TCPRequestRuleActionSendDashSpoeDashGroup
			rule.SpoeEngineName = a.Engine
			rule.SpoeGroupName = a.Group
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.UseService:
			rule.Action = models.TCPRequestRuleActionUseDashService
			rule.ServiceName = a.Name
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.Lua:
			rule.Action = models.TCPRequestRuleActionLua
			rule.LuaAction = a.Action
			rule.LuaParams = a.Params
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		default:
			return nil, NewConfError(ErrValidationError, fmt.Sprintf("unsupported action '%T' in tcp_request_rule", a))
		}
	case *tcp_types.Session:
		rule = &models.TCPRequestRule{
			Type: models.TCPRequestRuleTypeSession,
		}
		switch a := v.Action.(type) {
		case *tcp_actions.Accept:
			rule.Action = models.TCPRequestRuleActionAccept
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.Reject:
			rule.Action = models.TCPRequestRuleActionAccept
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.TrackSc:
			switch a.Type {
			case actions.TrackSc0:
				rule.Action = models.TCPRequestRuleActionTrackDashSc0
			case actions.TrackSc1:
				rule.Action = models.TCPRequestRuleActionTrackDashSc1
			case actions.TrackSc2:
				rule.Action = models.TCPRequestRuleActionTrackDashSc2
			}
			rule.TrackKey = a.Key
			if a.Table != "" {
				rule.TrackTable = a.Table
			}
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.ScIncGpc0:
			rule.Action = models.TCPRequestRuleActionScDashIncDashGpc0
			rule.ScIncID = a.ID
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.ScIncGpc1:
			rule.Action = models.TCPRequestRuleActionScDashIncDashGpc1
			rule.ScIncID = a.ID
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.ScSetGpt0:
			rule.Action = models.TCPRequestRuleActionScDashSetDashGpt0
			rule.ScIncID = a.ID
			rule.GptValue = a.Expr.String()
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.SetVar:
			rule.Action = models.TCPRequestRuleActionSetDashVar
			rule.VarScope = a.VarScope
			rule.VarName = a.VarName
			rule.Expr = a.Expr.String()
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.UnsetVar:
			rule.Action = models.TCPRequestRuleActionUnsetDashVar
			rule.VarScope = a.Scope
			rule.VarName = a.Name
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.SilentDrop:
			rule.Action = models.TCPRequestRuleActionSilentDashDrop
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		default:
			return nil, NewConfError(ErrValidationError, fmt.Sprintf("unsupported action '%T' in tcp_request_rule", a))
		}
	default:
		return nil, NewConfError(ErrValidationError, fmt.Sprintf("unsupported action '%T' in tcp_request_rule", v))
	}
	return rule, nil
}

func SerializeTCPRequestRule(f models.TCPRequestRule) (rule types.TCPType, err error) { //nolint:gocyclo,cyclop,maintidx
	switch f.Type {
	case models.TCPRequestRuleTypeConnection:
		switch f.Action {
		case models.TCPRequestRuleActionAccept:
			return &tcp_types.Connection{
				Action: &tcp_actions.Accept{
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
			}, nil
		case models.TCPRequestRuleActionReject:
			return &tcp_types.Connection{
				Action: &actions.Reject{
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
			}, nil
		case models.TCPRequestRuleActionExpectDashProxy:
			return &tcp_types.Connection{
				Action: &tcp_actions.ExpectProxy{
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
			}, nil
		case models.TCPRequestRuleActionExpectDashNetscalerDashCip:
			return &tcp_types.Connection{
				Action: &tcp_actions.ExpectNetscalerCip{
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
			}, nil
		case models.TCPRequestRuleActionCapture:
			return &tcp_types.Connection{
				Action: &tcp_actions.Capture{
					Expr:     common.Expression{Expr: strings.Split(f.Expr, " ")},
					Len:      f.CaptureLen,
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
			}, nil
		case models.TCPRequestRuleActionTrackDashSc0:
			return &tcp_types.Connection{
				Action: &actions.TrackSc{
					Type:     actions.TrackSc0,
					Key:      f.TrackKey,
					Table:    f.TrackTable,
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
			}, nil
		case models.TCPRequestRuleActionTrackDashSc1:
			return &tcp_types.Connection{
				Action: &actions.TrackSc{
					Type:     actions.TrackSc1,
					Key:      f.TrackKey,
					Table:    f.TrackTable,
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
			}, nil
		case models.TCPRequestRuleActionTrackDashSc2:
			return &tcp_types.Connection{
				Action: &actions.TrackSc{
					Type:     actions.TrackSc2,
					Key:      f.TrackKey,
					Table:    f.TrackTable,
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
			}, nil
		case models.TCPRequestRuleActionScDashIncDashGpc0:
			return &tcp_types.Connection{
				Action: &actions.ScIncGpc0{
					ID:       f.ScIncID,
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
			}, nil
		case models.TCPRequestRuleActionScDashIncDashGpc1:
			return &tcp_types.Connection{
				Action: &actions.ScIncGpc1{
					ID:       f.ScIncID,
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
			}, nil
		case models.TCPRequestRuleActionSilentDashDrop:
			return &tcp_types.Connection{
				Action: &actions.SilentDrop{
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
			}, nil
		case models.TCPRequestRuleActionLua:
			return &tcp_types.Connection{
				Action: &actions.Lua{
					Action:   f.LuaAction,
					Params:   f.LuaParams,
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
			}, nil
		}
		return nil, NewConfError(ErrValidationError, fmt.Sprintf("unsupported action '%T' in tcp_request_rule", f.Action))
	case models.TCPRequestRuleTypeContent:
		switch f.Action {
		case models.TCPRequestRuleActionAccept:
			return &tcp_types.Content{
				Action: &tcp_actions.Accept{
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
			}, nil
		case models.TCPRequestRuleActionDoDashResolve:
			return &tcp_types.Content{
				Action: &actions.DoResolve{
					Var:       f.VarName,
					Resolvers: f.ResolveResolvers,
					Protocol:  f.ResolveProtocol,
					Expr:      common.Expression{Expr: strings.Split(f.Expr, " ")},
					Cond:      f.Cond,
					CondTest:  f.CondTest,
				},
			}, nil
		case models.TCPRequestRuleActionReject:
			return &tcp_types.Content{
				Action: &actions.Reject{
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
			}, nil
		case models.TCPRequestRuleActionCapture:
			return &tcp_types.Content{
				Action: &tcp_actions.Capture{
					Expr:     common.Expression{Expr: strings.Split(f.Expr, " ")},
					Len:      f.CaptureLen,
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
			}, nil
		case "set-priority-class":
			return &tcp_types.Content{
				Action: &actions.SetPriorityClass{
					Expr:     common.Expression{Expr: strings.Split(f.Expr, " ")},
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
			}, nil
		case "set-priority-offset":
			return &tcp_types.Content{
				Action: &actions.SetPriorityOffset{
					Expr:     common.Expression{Expr: strings.Split(f.Expr, " ")},
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
			}, nil
		case models.TCPRequestRuleActionTrackDashSc0:
			return &tcp_types.Content{
				Action: &actions.TrackSc{
					Type:     actions.TrackSc0,
					Key:      f.TrackKey,
					Table:    f.TrackTable,
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
			}, nil
		case models.TCPRequestRuleActionTrackDashSc1:
			return &tcp_types.Content{
				Action: &actions.TrackSc{
					Type:     actions.TrackSc1,
					Key:      f.TrackKey,
					Table:    f.TrackTable,
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
			}, nil
		case models.TCPRequestRuleActionTrackDashSc2:
			return &tcp_types.Content{
				Action: &actions.TrackSc{
					Type:     actions.TrackSc2,
					Key:      f.TrackKey,
					Table:    f.TrackTable,
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
			}, nil
		case models.TCPRequestRuleActionScDashIncDashGpc0:
			return &tcp_types.Content{
				Action: &actions.ScIncGpc0{
					ID:       f.ScIncID,
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
			}, nil
		case models.TCPRequestRuleActionScDashIncDashGpc1:
			return &tcp_types.Content{
				Action: &actions.ScIncGpc1{
					ID:       f.ScIncID,
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
			}, nil
		case models.TCPRequestRuleActionSetDashDst:
			return &tcp_types.Content{
				Action: &actions.SetDst{
					Expr:     common.Expression{Expr: strings.Split(f.Expr, " ")},
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
			}, nil
		case models.TCPRequestRuleActionSetDashDstDashPort:
			return &tcp_types.Content{
				Action: &actions.SetDstPort{
					Expr:     common.Expression{Expr: strings.Split(f.Expr, " ")},
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
			}, nil
		case models.TCPRequestRuleActionSetDashVar:
			return &tcp_types.Content{
				Action: &actions.SetVar{
					VarName:  f.VarName,
					VarScope: f.VarScope,
					Expr:     common.Expression{Expr: strings.Split(f.Expr, " ")},
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
			}, nil
		case models.TCPRequestRuleActionUnsetDashVar:
			return &tcp_types.Content{
				Action: &actions.UnsetVar{
					Name:     f.VarName,
					Scope:    f.VarScope,
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
			}, nil
		case models.TCPRequestRuleActionSilentDashDrop:
			return &tcp_types.Content{
				Action: &actions.SilentDrop{
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
			}, nil
		case models.TCPRequestRuleActionSendDashSpoeDashGroup:
			return &tcp_types.Content{
				Action: &actions.SendSpoeGroup{
					Engine:   f.SpoeEngineName,
					Group:    f.SpoeGroupName,
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
			}, nil
		case models.TCPRequestRuleActionUseDashService:
			return &tcp_types.Content{
				Action: &actions.UseService{
					Name:     f.ServiceName,
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
			}, nil
		case models.TCPRequestRuleActionLua:
			return &tcp_types.Content{
				Action: &actions.Lua{
					Action:   f.LuaAction,
					Params:   f.LuaParams,
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
			}, nil
		}
		return nil, NewConfError(ErrValidationError, fmt.Sprintf("unsupported action '%T' in tcp_request_rule", f.Action))
	case models.TCPRequestRuleTypeSession:
		switch f.Action {
		case models.TCPRequestRuleActionAccept:
			return &tcp_types.Session{
				Action: &tcp_actions.Accept{
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
			}, nil
		case models.TCPRequestRuleActionReject:
			return &tcp_types.Session{
				Action: &actions.Reject{
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
			}, nil
		case models.TCPRequestRuleActionTrackDashSc0:
			return &tcp_types.Session{
				Action: &actions.TrackSc{
					Type:     actions.TrackSc0,
					Key:      f.TrackKey,
					Table:    f.TrackTable,
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
			}, nil
		case models.TCPRequestRuleActionTrackDashSc1:
			return &tcp_types.Session{
				Action: &actions.TrackSc{
					Type:     actions.TrackSc1,
					Key:      f.TrackKey,
					Table:    f.TrackTable,
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
			}, nil
		case models.TCPRequestRuleActionTrackDashSc2:
			return &tcp_types.Session{
				Action: &actions.TrackSc{
					Type:     actions.TrackSc2,
					Key:      f.TrackKey,
					Table:    f.TrackTable,
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
			}, nil
		case models.TCPRequestRuleActionScDashIncDashGpc0:
			return &tcp_types.Session{
				Action: &actions.ScIncGpc0{
					ID:       f.ScIncID,
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
			}, nil
		case models.TCPRequestRuleActionScDashIncDashGpc1:
			return &tcp_types.Session{
				Action: &actions.ScIncGpc1{
					ID:       f.ScIncID,
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
			}, nil
		case "sc-inc-gpt0":
			return &tcp_types.Session{
				Action: &actions.ScSetGpt0{
					ID:       f.ScIncID,
					Expr:     common.Expression{Expr: []string{f.GptValue}},
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
			}, nil
		case models.TCPRequestRuleActionSetDashVar:
			return &tcp_types.Session{
				Action: &actions.SetVar{
					VarName:  f.VarName,
					VarScope: f.VarScope,
					Expr:     common.Expression{Expr: strings.Split(f.Expr, " ")},
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
			}, nil
		case models.TCPRequestRuleActionUnsetDashVar:
			return &tcp_types.Session{
				Action: &actions.UnsetVar{
					Name:     f.VarName,
					Scope:    f.VarScope,
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
			}, nil
		case models.TCPRequestRuleActionSilentDashDrop:
			return &tcp_types.Session{
				Action: &actions.SilentDrop{
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
			}, nil
		}
		return nil, NewConfError(ErrValidationError, fmt.Sprintf("unsupported action '%T' in tcp_request_rule", f.Action))
	case models.TCPRequestRuleTypeInspectDashDelay:
		if f.Timeout == nil {
			return nil, NewConfError(ErrValidationError, fmt.Sprintf("unsupported action '%T' in tcp_request_rule", f.Type))
		}
		return &tcp_types.InspectDelay{
			Timeout: strconv.FormatInt(*f.Timeout, 10),
		}, nil
	}
	return nil, NewConfError(ErrValidationError, fmt.Sprintf("unsupported action '%T' in tcp_request_rule", f.Type))
}
