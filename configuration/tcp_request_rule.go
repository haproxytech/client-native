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
	parser "github.com/haproxytech/config-parser/v3"
	"github.com/haproxytech/config-parser/v3/common"
	parser_errors "github.com/haproxytech/config-parser/v3/errors"
	tcp_actions "github.com/haproxytech/config-parser/v3/parsers/tcp/actions"
	tcp_types "github.com/haproxytech/config-parser/v3/parsers/tcp/types"
	"github.com/haproxytech/config-parser/v3/types"

	"github.com/haproxytech/client-native/v2/misc"
	"github.com/haproxytech/client-native/v2/models"
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

	tcpRules, err := ParseTCPRequestRules(parentType, parentName, p)
	if err != nil {
		return v, nil, c.HandleError("", parentType, parentName, "", false, err)
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
		return c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	if err := c.SaveData(p, t, transactionID == ""); err != nil {
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

func ParseTCPRequestRules(t, pName string, p *parser.Parser) (models.TCPRequestRules, error) {
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

	rules := data.([]types.TCPType)
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

func ParseTCPRequestRule(f types.TCPType) (rule *models.TCPRequestRule, err error) { //nolint:gocyclo
	switch v := f.(type) {
	case *tcp_types.InspectDelay:
		return &models.TCPRequestRule{
			Type:    models.TCPRequestRuleTypeInspectDelay,
			Timeout: misc.ParseTimeout(v.Timeout),
		}, nil

	case *tcp_types.Connection:
		rule = &models.TCPRequestRule{
			Type:     models.TCPRequestRuleTypeConnection,
			Cond:     v.Cond,
			CondTest: v.CondTest,
		}

		switch a := v.Action.(type) {
		case *tcp_actions.Accept:
			rule.Action = models.TCPRequestRuleActionAccept
		case *tcp_actions.Reject:
			rule.Action = models.TCPRequestRuleActionReject
		case *tcp_actions.ExpectProxy:
			rule.Action = "expect-proxy layer4"
		case *tcp_actions.ExpectNetscalerCip:
			rule.Action = "expect-netscaler-cip layer4"
		case *tcp_actions.Capture:
			rule.Action = models.TCPRequestRuleActionCapture
			rule.Expr = a.Expr.String()
			rule.CaptureLen = a.Len
		case *tcp_actions.TrackSc0:
			rule.Action = models.TCPRequestRuleActionTrackSc0
			rule.TrackKey = a.Key
			if a.Table != "" {
				rule.TrackTable = a.Table
			}
		case *tcp_actions.TrackSc1:
			rule.Action = models.TCPRequestRuleActionTrackSc1
			rule.TrackKey = a.Key
			if a.Table != "" {
				rule.TrackTable = a.Table
			}
		case *tcp_actions.TrackSc2:
			rule.Action = models.TCPRequestRuleActionTrackSc2
			rule.TrackKey = a.Key
			if a.Table != "" {
				rule.TrackTable = a.Table
			}
		case *tcp_actions.ScIncGpc0:
			rule.Action = models.TCPRequestRuleActionScSetGpt0
			rule.ScIncID = a.ScID
		case *tcp_actions.ScIncGpc1:
			rule.Action = models.TCPRequestRuleActionScIncGpc1
			rule.ScIncID = a.ScID
		case *tcp_actions.ScSetGpt0:
			rule.Action = models.TCPRequestRuleActionScSetGpt0
			rule.ScIncID = a.ScID
			rule.Expr = a.Value
		case *tcp_actions.SetSrc:
			rule.Action = models.TCPRequestRuleActionSetSrc
			rule.Expr = a.Expr.String()
		case *tcp_actions.Lua:
			rule.Action = models.TCPRequestRuleActionLua
			rule.LuaAction = a.Action
			rule.LuaParams = a.Params
		default:
			return nil, NewConfError(ErrValidationError, "unsupported action in tcp_request_rule")
		}

		return rule, nil
	case *tcp_types.Content:
		rule = &models.TCPRequestRule{
			Type:     models.TCPRequestRuleTypeContent,
			Cond:     v.Cond,
			CondTest: v.CondTest,
		}

		switch a := v.Action.(type) {
		case *tcp_actions.Accept:
			rule.Action = models.TCPRequestRuleActionAccept
		case *tcp_actions.DoResolve:
			rule.Action = models.TCPRequestRuleActionDoResolve
			rule.VarName = a.Var
			rule.ResolveResolvers = a.Resolvers
			rule.ResolveProtocol = a.Protocol
			rule.Expr = a.Expr.String()
		case *tcp_actions.Reject:
			rule.Action = models.TCPRequestRuleActionReject
		case *tcp_actions.Capture:
			rule.Action = models.TCPRequestRuleActionCapture
			rule.Expr = a.Expr.String()
			rule.CaptureLen = a.Len
		case *tcp_actions.SetPriorityClass:
			rule.Action = "set-priority-class"
			rule.Expr = a.Expr.String()
		case *tcp_actions.SetPriorityOffset:
			rule.Action = "set-priority-offset"
			rule.Expr = a.Expr.String()
		case *tcp_actions.TrackSc0:
			rule.Action = models.TCPRequestRuleActionTrackSc0
			rule.TrackKey = a.Key
			if a.Table != "" {
				rule.TrackTable = a.Table
			}
		case *tcp_actions.TrackSc1:
			rule.Action = models.TCPRequestRuleActionTrackSc1
			rule.TrackKey = a.Key
			if a.Table != "" {
				rule.TrackTable = a.Table
			}
		case *tcp_actions.TrackSc2:
			rule.Action = models.TCPRequestRuleActionTrackSc2
			rule.TrackKey = a.Key
			if a.Table != "" {
				rule.TrackTable = a.Table
			}
		case *tcp_actions.ScIncGpc0:
			rule.Action = models.TCPRequestRuleActionScIncGpc0
			rule.ScIncID = a.ScID
		case *tcp_actions.ScIncGpc1:
			rule.Action = models.TCPRequestRuleActionScIncGpc1
			rule.ScIncID = a.ScID
		case *tcp_actions.ScSetGpt0:
			rule.Action = models.TCPRequestRuleActionScSetGpt0
			rule.ScIncID = a.ScID
			rule.GptValue = a.Value
		case *tcp_actions.SetDst:
			rule.Action = models.TCPRequestRuleActionSetDst
			rule.Expr = a.Expr.String()
		case *tcp_actions.SetDstPort:
			rule.Action = models.TCPRequestRuleActionSetDstPort
			rule.Expr = a.Expr.String()
		case *tcp_actions.SetVar:
			rule.Action = models.TCPRequestRuleActionSetVar
			rule.VarScope = a.VarScope
			rule.VarName = a.VarName
			rule.Expr = a.Expr.String()
		case *tcp_actions.UnsetVar:
			rule.Action = models.TCPRequestRuleActionUnsetVar
			rule.VarScope = a.VarScope
			rule.VarName = a.VarName
		case *tcp_actions.SilentDrop:
			rule.Action = models.TCPRequestRuleActionSilentDrop
		case *tcp_actions.SendSpoeGroup:
			rule.Action = models.TCPRequestRuleActionSendSpoeGroup
			rule.SpoeEngineName = a.Engine
			rule.SpoeGroupName = a.Group
		case *tcp_actions.UseService:
			rule.Action = models.TCPRequestRuleActionUseService
			rule.ServiceName = a.ServiceName
		case *tcp_actions.Lua:
			rule.Action = models.TCPRequestRuleActionLua
			rule.LuaAction = a.Action
			rule.LuaParams = a.Params
		default:
			return nil, NewConfError(ErrValidationError, "unsupported action in tcp_request_rule")
		}
	case *tcp_types.Session:
		rule = &models.TCPRequestRule{
			Type:     models.TCPRequestRuleTypeSession,
			Cond:     v.Cond,
			CondTest: v.CondTest,
		}
		switch a := v.Action.(type) {
		case *tcp_actions.Accept:
			rule.Action = models.TCPRequestRuleActionAccept
		case *tcp_actions.Reject:
			rule.Action = models.TCPRequestRuleActionAccept
		case *tcp_actions.TrackSc0:
			rule.Action = models.TCPRequestRuleActionTrackSc0
			rule.TrackKey = a.Key
			if a.Table != "" {
				rule.TrackTable = a.Table
			}
		case *tcp_actions.TrackSc1:
			rule.Action = models.TCPRequestRuleActionTrackSc1
			rule.TrackKey = a.Key
			if a.Table != "" {
				rule.TrackTable = a.Table
			}
		case *tcp_actions.TrackSc2:
			rule.Action = models.TCPRequestRuleActionTrackSc2
			rule.TrackKey = a.Key
			if a.Table != "" {
				rule.TrackTable = a.Table
			}
		case *tcp_actions.ScIncGpc0:
			rule.Action = models.TCPRequestRuleActionScIncGpc0
			rule.ScIncID = a.ScID
		case *tcp_actions.ScIncGpc1:
			rule.Action = models.TCPRequestRuleActionScIncGpc1
			rule.ScIncID = a.ScID
		case *tcp_actions.ScSetGpt0:
			rule.Action = models.TCPRequestRuleActionScSetGpt0
			rule.ScIncID = a.ScID
			rule.GptValue = a.Value
		case *tcp_actions.SetVar:
			rule.Action = models.TCPRequestRuleActionSetVar
			rule.VarScope = a.VarScope
			rule.VarName = a.VarName
			rule.Expr = a.Expr.String()
		case *tcp_actions.UnsetVar:
			rule.Action = models.TCPRequestRuleActionUnsetVar
			rule.VarScope = a.VarScope
			rule.VarName = a.VarName
		case *tcp_actions.SilentDrop:
			rule.Action = models.TCPRequestRuleActionSilentDrop
		default:
			return nil, NewConfError(ErrValidationError, "unsupported action in tcp_request_rule")
		}
	default:
		return nil, NewConfError(ErrValidationError, "unsupported action in tcp_request_rule")
	}

	return rule, nil
}

func SerializeTCPRequestRule(f models.TCPRequestRule) (rule types.TCPType, err error) { //nolint:gocyclo
	switch f.Type {
	case models.TCPRequestRuleTypeConnection:
		switch f.Action {
		case models.TCPRequestRuleActionAccept:
			return &tcp_types.Connection{
				Action:   &tcp_actions.Accept{},
				Cond:     f.Cond,
				CondTest: f.CondTest,
			}, nil
		case models.TCPRequestRuleActionReject:
			return &tcp_types.Connection{
				Action:   &tcp_actions.Reject{},
				Cond:     f.Cond,
				CondTest: f.CondTest,
			}, nil
		case "expect-proxy layer4":
			return &tcp_types.Connection{
				Action:   &tcp_actions.ExpectProxy{},
				Cond:     f.Cond,
				CondTest: f.CondTest,
			}, nil
		case "expect-netscaler-cip layer4":
			return &tcp_types.Connection{
				Action:   &tcp_actions.ExpectNetscalerCip{},
				Cond:     f.Cond,
				CondTest: f.CondTest,
			}, nil
		case models.TCPRequestRuleActionCapture:
			return &tcp_types.Connection{
				Action: &tcp_actions.Capture{
					Expr: common.Expression{Expr: strings.Split(f.Expr, " ")},
					Len:  f.CaptureLen,
				},
				Cond:     f.Cond,
				CondTest: f.CondTest,
			}, nil
		case models.TCPRequestRuleActionTrackSc0:
			return &tcp_types.Connection{
				Action: &tcp_actions.TrackSc0{
					Key:   f.TrackKey,
					Table: f.TrackTable,
				},
				Cond:     f.Cond,
				CondTest: f.CondTest,
			}, nil
		case models.TCPRequestRuleActionTrackSc1:
			return &tcp_types.Connection{
				Action: &tcp_actions.TrackSc1{
					Key:   f.TrackKey,
					Table: f.TrackTable,
				},
				Cond:     f.Cond,
				CondTest: f.CondTest,
			}, nil
		case models.TCPRequestRuleActionTrackSc2:
			return &tcp_types.Connection{
				Action: &tcp_actions.TrackSc2{
					Key:   f.TrackKey,
					Table: f.TrackTable,
				},
				Cond:     f.Cond,
				CondTest: f.CondTest,
			}, nil
		case models.TCPRequestRuleActionScIncGpc0:
			return &tcp_types.Connection{
				Action: &tcp_actions.ScIncGpc0{
					ScID: f.ScIncID,
				},
				Cond:     f.Cond,
				CondTest: f.CondTest,
			}, nil
		case models.TCPRequestRuleActionScIncGpc1:
			return &tcp_types.Connection{
				Action: &tcp_actions.ScIncGpc1{
					ScID: f.ScIncID,
				},
				Cond:     f.Cond,
				CondTest: f.CondTest,
			}, nil
		case models.TCPRequestRuleActionLua:
			return &tcp_types.Connection{
				Action: &tcp_actions.Lua{
					Action: f.LuaAction,
					Params: f.LuaParams,
				},
				Cond:     f.Cond,
				CondTest: f.CondTest,
			}, nil
		}
		return nil, NewConfError(ErrValidationError, "unsupported action in tcp_request_rule")
	case models.TCPRequestRuleTypeContent:
		switch f.Action {
		case models.TCPRequestRuleActionAccept:
			return &tcp_types.Content{
				Action:   &tcp_actions.Accept{},
				Cond:     f.Cond,
				CondTest: f.CondTest,
			}, nil
		case models.TCPRequestRuleActionDoResolve:
			return &tcp_types.Content{
				Action: &tcp_actions.DoResolve{
					Var:       f.VarName,
					Resolvers: f.ResolveResolvers,
					Protocol:  f.ResolveProtocol,
					Expr:      common.Expression{Expr: strings.Split(f.Expr, " ")},
				},
				Cond:     f.Cond,
				CondTest: f.CondTest,
			}, nil
		case models.TCPRequestRuleActionReject:
			return &tcp_types.Content{
				Action:   &tcp_actions.Reject{},
				Cond:     f.Cond,
				CondTest: f.CondTest,
			}, nil
		case models.TCPRequestRuleActionCapture:
			return &tcp_types.Content{
				Action: &tcp_actions.Capture{
					Expr: common.Expression{Expr: strings.Split(f.Expr, " ")},
					Len:  f.CaptureLen,
				},
				Cond:     f.Cond,
				CondTest: f.CondTest,
			}, nil
		case "set-priority-class":
			return &tcp_types.Content{
				Action: &tcp_actions.SetPriorityClass{
					Expr: common.Expression{Expr: strings.Split(f.Expr, " ")},
				},
				Cond:     f.Cond,
				CondTest: f.CondTest,
			}, nil
		case "set-priority-offset":
			return &tcp_types.Content{
				Action: &tcp_actions.SetPriorityOffset{
					Expr: common.Expression{Expr: strings.Split(f.Expr, " ")},
				},
				Cond:     f.Cond,
				CondTest: f.CondTest,
			}, nil
		case models.TCPRequestRuleActionTrackSc0:
			return &tcp_types.Content{
				Action: &tcp_actions.TrackSc0{
					Key:   f.TrackKey,
					Table: f.TrackTable,
				},
				Cond:     f.Cond,
				CondTest: f.CondTest,
			}, nil
		case models.TCPRequestRuleActionTrackSc1:
			return &tcp_types.Content{
				Action: &tcp_actions.TrackSc1{
					Key:   f.TrackKey,
					Table: f.TrackTable,
				},
				Cond:     f.Cond,
				CondTest: f.CondTest,
			}, nil
		case models.TCPRequestRuleActionTrackSc2:
			return &tcp_types.Content{
				Action: &tcp_actions.TrackSc2{
					Key:   f.TrackKey,
					Table: f.TrackTable,
				},
				Cond:     f.Cond,
				CondTest: f.CondTest,
			}, nil
		case models.TCPRequestRuleActionScIncGpc0:
			return &tcp_types.Content{
				Action: &tcp_actions.ScIncGpc0{
					ScID: f.ScIncID,
				},
				Cond:     f.Cond,
				CondTest: f.CondTest,
			}, nil
		case models.TCPRequestRuleActionScIncGpc1:
			return &tcp_types.Content{
				Action: &tcp_actions.ScIncGpc1{
					ScID: f.ScIncID,
				},
				Cond:     f.Cond,
				CondTest: f.CondTest,
			}, nil
		case models.TCPRequestRuleActionSetDst:
			return &tcp_types.Content{
				Action: &tcp_actions.SetDst{
					Expr: common.Expression{Expr: strings.Split(f.Expr, " ")},
				},
				Cond:     f.Cond,
				CondTest: f.CondTest,
			}, nil
		case models.TCPRequestRuleActionSetDstPort:
			return &tcp_types.Content{
				Action: &tcp_actions.SetDstPort{
					Expr: common.Expression{Expr: strings.Split(f.Expr, " ")},
				},
				Cond:     f.Cond,
				CondTest: f.CondTest,
			}, nil
		case models.TCPRequestRuleActionSetVar:
			return &tcp_types.Content{
				Action: &tcp_actions.SetVar{
					VarName:  f.VarName,
					VarScope: f.VarScope,
					Expr:     common.Expression{Expr: strings.Split(f.Expr, " ")},
				},
				Cond:     f.Cond,
				CondTest: f.CondTest,
			}, nil
		case models.TCPRequestRuleActionUnsetVar:
			return &tcp_types.Content{
				Action: &tcp_actions.UnsetVar{
					VarName:  f.VarName,
					VarScope: f.VarScope,
				},
				Cond:     f.Cond,
				CondTest: f.CondTest,
			}, nil
		case models.TCPRequestRuleActionSilentDrop:
			return &tcp_types.Content{
				Action:   &tcp_actions.SilentDrop{},
				Cond:     f.Cond,
				CondTest: f.CondTest,
			}, nil
		case models.TCPRequestRuleActionSendSpoeGroup:
			return &tcp_types.Content{
				Action: &tcp_actions.SendSpoeGroup{
					Engine: f.SpoeEngineName,
					Group:  f.SpoeGroupName,
				},
				Cond:     f.Cond,
				CondTest: f.CondTest,
			}, nil
		case models.TCPRequestRuleActionUseService:
			return &tcp_types.Content{
				Action: &tcp_actions.UseService{
					ServiceName: f.ServiceName,
				},
				Cond:     f.Cond,
				CondTest: f.CondTest,
			}, nil
		case models.TCPRequestRuleActionLua:
			return &tcp_types.Content{
				Action: &tcp_actions.Lua{
					Action: f.LuaAction,
					Params: f.LuaParams,
				},
				Cond:     f.Cond,
				CondTest: f.CondTest,
			}, nil
		}
		return nil, NewConfError(ErrValidationError, "unsupported action in tcp_request_rule")
	case models.TCPRequestRuleTypeSession:
		switch f.Action {
		case models.TCPRequestRuleActionAccept:
			return &tcp_types.Session{
				Action:   &tcp_actions.Accept{},
				Cond:     f.Cond,
				CondTest: f.CondTest,
			}, nil
		case models.TCPRequestRuleActionReject:
			return &tcp_types.Session{
				Action:   &tcp_actions.Reject{},
				Cond:     f.Cond,
				CondTest: f.CondTest,
			}, nil
		case models.TCPRequestRuleActionTrackSc0:
			return &tcp_types.Session{
				Action: &tcp_actions.TrackSc0{
					Key:   f.TrackKey,
					Table: f.TrackTable,
				},
				Cond:     f.Cond,
				CondTest: f.CondTest,
			}, nil
		case models.TCPRequestRuleActionTrackSc1:
			return &tcp_types.Session{
				Action: &tcp_actions.TrackSc1{
					Key:   f.TrackKey,
					Table: f.TrackTable,
				},
				Cond:     f.Cond,
				CondTest: f.CondTest,
			}, nil
		case models.TCPRequestRuleActionTrackSc2:
			return &tcp_types.Session{
				Action: &tcp_actions.TrackSc2{
					Key:   f.TrackKey,
					Table: f.TrackTable,
				},
				Cond:     f.Cond,
				CondTest: f.CondTest,
			}, nil
		case models.TCPRequestRuleActionScIncGpc0:
			return &tcp_types.Session{
				Action: &tcp_actions.ScIncGpc0{
					ScID: f.ScIncID,
				},
				Cond:     f.Cond,
				CondTest: f.CondTest,
			}, nil
		case models.TCPRequestRuleActionScIncGpc1:
			return &tcp_types.Session{
				Action: &tcp_actions.ScIncGpc1{
					ScID: f.ScIncID,
				},
				Cond:     f.Cond,
				CondTest: f.CondTest,
			}, nil
		case "sc-inc-gpt0":
			return &tcp_types.Session{
				Action: &tcp_actions.ScSetGpt0{
					ScID:  f.ScIncID,
					Value: f.GptValue,
				},
				Cond:     f.Cond,
				CondTest: f.CondTest,
			}, nil
		case models.TCPRequestRuleActionSetVar:
			return &tcp_types.Session{
				Action: &tcp_actions.SetVar{
					VarName:  f.VarName,
					VarScope: f.VarScope,
					Expr:     common.Expression{Expr: strings.Split(f.Expr, " ")},
				},
				Cond:     f.Cond,
				CondTest: f.CondTest,
			}, nil
		case models.TCPRequestRuleActionUnsetVar:
			return &tcp_types.Session{
				Action: &tcp_actions.UnsetVar{
					VarName:  f.VarName,
					VarScope: f.VarScope,
				},
				Cond:     f.Cond,
				CondTest: f.CondTest,
			}, nil
		case models.TCPRequestRuleActionSilentDrop:
			return &tcp_types.Session{
				Action:   &tcp_actions.SilentDrop{},
				Cond:     f.Cond,
				CondTest: f.CondTest,
			}, nil
		}
		return nil, NewConfError(ErrValidationError, "unsupported action in tcp_request_rule")
	case models.TCPRequestRuleTypeInspectDelay:
		if f.Timeout == nil {
			return nil, NewConfError(ErrValidationError, "unsupported action in tcp_request_rule")
		}
		return &tcp_types.InspectDelay{
			Timeout: strconv.FormatInt(*f.Timeout, 10),
		}, nil
	}

	return nil, NewConfError(ErrValidationError, "unsupported action in tcp_request_rule")
}
