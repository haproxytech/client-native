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
	parser "github.com/haproxytech/client-native/v6/config-parser"
	"github.com/haproxytech/client-native/v6/config-parser/common"
	parser_errors "github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/parsers/actions"
	http_actions "github.com/haproxytech/client-native/v6/config-parser/parsers/http/actions"
	tcp_actions "github.com/haproxytech/client-native/v6/config-parser/parsers/tcp/actions"
	tcp_types "github.com/haproxytech/client-native/v6/config-parser/parsers/tcp/types"
	"github.com/haproxytech/client-native/v6/config-parser/types"

	"github.com/haproxytech/client-native/v6/configuration/options"
	"github.com/haproxytech/client-native/v6/misc"
	"github.com/haproxytech/client-native/v6/models"
)

type TCPRequestRule interface {
	GetTCPRequestRules(parentType, parentName string, transactionID string) (int64, models.TCPRequestRules, error)
	GetTCPRequestRule(id int64, parentType, parentName string, transactionID string) (int64, *models.TCPRequestRule, error)
	DeleteTCPRequestRule(id int64, parentType string, parentName string, transactionID string, version int64) error
	CreateTCPRequestRule(id int64, parentType string, parentName string, data *models.TCPRequestRule, transactionID string, version int64) error
	EditTCPRequestRule(id int64, parentType string, parentName string, data *models.TCPRequestRule, transactionID string, version int64) error
	ReplaceTCPRequestRules(parentType string, parentName string, data models.TCPRequestRules, transactionID string, version int64) error
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

	section, parentName, err := getParserFromParent("tcp-request", parentType, parentName)
	if err != nil {
		return v, nil, err
	}

	data, err := p.GetOne(section, parentName, "tcp-request", int(id))
	if err != nil {
		return v, nil, c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, "", false, err)
	}

	tcpRule, err := ParseTCPRequestRule(data.(types.TCPType))
	if err != nil {
		return v, nil, err
	}

	return v, tcpRule, nil
}

// DeleteTCPRequestRule deletes a tcp request rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) DeleteTCPRequestRule(id int64, parentType string, parentName string, transactionID string, version int64) error {
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}
	section, parentName, err := getParserFromParent("tcp-request", parentType, parentName)
	if err != nil {
		return err
	}
	if err := p.Delete(section, parentName, "tcp-request", int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

// CreateTCPRequestRule creates a tcp request rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) CreateTCPRequestRule(id int64, parentType string, parentName string, data *models.TCPRequestRule, transactionID string, version int64) error {
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

	section, parentName, err := getParserFromParent("tcp-request", parentType, parentName)
	if err != nil {
		return err
	}

	s, err := SerializeTCPRequestRule(*data, &c.ConfigurationOptions)
	if err != nil {
		return err
	}

	if err := p.Insert(section, parentName, "tcp-request", s, int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

// EditTCPRequestRule edits a tcp request rule in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
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
	section, parentName, err := getParserFromParent("tcp-request", parentType, parentName)
	if err != nil {
		return err
	}

	if _, err = p.GetOne(section, parentName, "tcp-request", int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	s, err := SerializeTCPRequestRule(*data, &c.ConfigurationOptions)
	if err != nil {
		return err
	}

	if err := p.Set(section, parentName, "tcp-request", s, int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

// ReplaceTCPRequestRules replaces all TCP Request Rule lines in configuration for a parentType/parentName.
// One of version or transactionID is mandatory.
// Returns error on fail, nil on success.
//
//nolint:dupl
func (c *client) ReplaceTCPRequestRules(parentType string, parentName string, data models.TCPRequestRules, transactionID string, version int64) error {
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

	section, parentName, err := getParserFromParent("tcp-request", parentType, parentName)
	if err != nil {
		return err
	}

	tcpRequestRules, err := ParseTCPRequestRules(parentType, parentName, p)
	if err != nil {
		return c.HandleError("", parentType, parentName, "", false, err)
	}

	for i := range tcpRequestRules {
		// Always delete index 0
		if err := p.Delete(section, parentName, "tcp-request", 0); err != nil {
			return c.HandleError(strconv.FormatInt(int64(i), 10), parentType, parentName, t, transactionID == "", err)
		}
	}

	for i, newTCPRequestRule := range data {
		s, err := SerializeTCPRequestRule(*newTCPRequestRule, &c.ConfigurationOptions)
		if err != nil {
			return err
		}
		if err := p.Insert(section, parentName, "tcp-request", s, i); err != nil {
			return c.HandleError(strconv.FormatInt(int64(i), 10), parentType, parentName, t, transactionID == "", err)
		}
	}

	return c.SaveData(p, t, transactionID == "")
}

func ParseTCPRequestRules(t, pName string, p parser.Parser) (models.TCPRequestRules, error) {
	section, pName, err := getParserFromParent("tcp-request", t, pName)
	if err != nil {
		return nil, err
	}

	var tcpReqRules models.TCPRequestRules
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
	for _, r := range rules {
		tcpReqRule, err := ParseTCPRequestRule(r)
		if err == nil {
			tcpReqRules = append(tcpReqRules, tcpReqRule)
		}
	}
	return tcpReqRules, nil
}

func ParseTCPRequestRule(f types.TCPType) (*models.TCPRequestRule, error) { //nolint:gocyclo,cyclop,maintidx
	var rule *models.TCPRequestRule
	switch v := f.(type) {
	case *tcp_types.InspectDelay:
		return &models.TCPRequestRule{
			Type:     models.TCPRequestRuleTypeInspectDashDelay,
			Timeout:  misc.ParseTimeout(v.Timeout),
			Metadata: parseMetadata(v.Comment),
		}, nil

	case *tcp_types.Connection:
		rule = &models.TCPRequestRule{
			Type:     models.TCPRequestRuleTypeConnection,
			Metadata: parseMetadata(v.Comment),
		}

		switch a := v.Action.(type) {
		case *actions.Accept:
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
			rule.Action = models.TCPRequestRuleActionTrackDashSc
			rule.TrackKey = a.Key
			if a.Table != "" {
				rule.TrackTable = a.Table
			}
			rule.TrackStickCounter = &a.StickCounter
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.ScAddGpc:
			rule.Action = models.TCPRequestRuleActionScDashAddDashGpc
			rule.ScIncID = a.ID
			rule.ScIdx = a.Idx
			rule.ScInt = a.Int
			rule.Expr = a.Expr.String()
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.ScIncGpc:
			rule.Action = models.TCPRequestRuleActionScDashIncDashGpc
			rule.ScIncID = a.ID
			rule.ScIdx = a.Idx
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
		case *actions.ScSetGpt:
			rule.Action = models.TCPRequestRuleActionScDashSetDashGpt
			rule.ScIncID = a.ScID
			rule.ScIdx = strconv.FormatInt(a.Idx, 10)
			rule.ScInt = a.Int
			rule.Expr = a.Expr.String()
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.ScSetGpt0:
			rule.Action = models.TCPRequestRuleActionScDashSetDashGpt0
			rule.ScIncID = a.ID
			rule.ScInt = a.Int
			rule.Expr = a.Expr.String()
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.SetFcMark:
			rule.Action = models.TCPRequestRuleActionSetDashFcDashMark
			rule.Expr = a.Expr.String()
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.SetFcTos:
			rule.Action = models.TCPRequestRuleActionSetDashFcDashTos
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
			rule.RstTTL = a.RstTTL
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.Lua:
			rule.Action = models.TCPRequestRuleActionLua
			rule.LuaAction = a.Action
			rule.LuaParams = a.Params
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.SetMark:
			rule.Action = models.TCPRequestRuleActionSetDashMark
			rule.MarkValue = a.Value
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.SetSrcPort:
			rule.Action = models.TCPRequestRuleActionSetDashSrcDashPort
			rule.Expr = a.Expr.String()
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.SetTos:
			rule.Action = models.TCPRequestRuleActionSetDashTos
			rule.TosValue = a.Value
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
		case *actions.SetVarFmt:
			rule.Action = models.TCPRequestRuleActionSetDashVarDashFmt
			rule.VarName = a.VarName
			rule.VarFormat = strings.Join(a.Fmt.Expr, " ")
			rule.VarScope = a.VarScope
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.UnsetVar:
			rule.Action = models.TCPRequestRuleActionUnsetDashVar
			rule.VarScope = a.Scope
			rule.VarName = a.Name
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.DoLog:
			rule.Action = models.TCPRequestRuleActionDoDashLog
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		default:
			return nil, NewConfError(ErrValidationError, fmt.Sprintf("unsupported action '%s' in tcp_request_rule", a))
		}

		return rule, nil
	case *tcp_types.Content:
		rule = &models.TCPRequestRule{
			Type:     models.TCPRequestRuleTypeContent,
			Metadata: parseMetadata(v.Comment),
		}

		switch a := v.Action.(type) {
		case *actions.Accept:
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
			rule.Action = models.TCPRequestRuleActionSetDashPriorityDashClass
			rule.Expr = a.Expr.String()
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.SetPriorityOffset:
			rule.Action = models.TCPRequestRuleActionSetDashPriorityDashOffset
			rule.Expr = a.Expr.String()
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.TrackSc:
			rule.Action = models.TCPRequestRuleActionTrackDashSc
			rule.TrackKey = a.Key
			if a.Table != "" {
				rule.TrackTable = a.Table
			}
			rule.TrackStickCounter = &a.StickCounter
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.ScAddGpc:
			rule.Action = models.TCPRequestRuleActionScDashAddDashGpc
			rule.ScIncID = a.ID
			rule.ScIdx = a.Idx
			rule.ScInt = a.Int
			rule.Expr = a.Expr.String()
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.ScIncGpc:
			rule.Action = models.TCPRequestRuleActionScDashIncDashGpc
			rule.ScIncID = a.ID
			rule.ScIdx = a.Idx
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
		case *actions.ScSetGpt:
			rule.Action = models.TCPRequestRuleActionScDashSetDashGpt
			rule.ScIncID = a.ScID
			rule.ScIdx = strconv.FormatInt(a.Idx, 10)
			rule.ScInt = a.Int
			rule.Expr = a.Expr.String()
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.ScSetGpt0:
			rule.Action = models.TCPRequestRuleActionScDashSetDashGpt0
			rule.ScIncID = a.ID
			rule.ScInt = a.Int
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
		case *tcp_actions.SetSrc:
			rule.Action = models.TCPRequestRuleActionSetDashSrc
			rule.Expr = a.Expr.String()
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.SetSrcPort:
			rule.Action = models.TCPRequestRuleActionSetDashSrcDashPort
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
			rule.RstTTL = a.RstTTL
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
		case *actions.SetBandwidthLimit:
			rule.Action = models.TCPRequestRuleActionSetDashBandwidthDashLimit
			rule.BandwidthLimitName = a.Name
			rule.BandwidthLimitLimit = a.Limit.String()
			rule.BandwidthLimitPeriod = a.Period.String()
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.SetMark:
			rule.Action = models.TCPRequestRuleActionSetDashMark
			rule.MarkValue = a.Value
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.SetBcMark:
			rule.Action = models.TCPRequestRuleActionSetDashBcDashMark
			rule.Expr = a.Expr.String()
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.SetBcTos:
			rule.Action = models.TCPRequestRuleActionSetDashBcDashTos
			rule.Expr = a.Expr.String()
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.SetFcMark:
			rule.Action = models.TCPRequestRuleActionSetDashFcDashMark
			rule.Expr = a.Expr.String()
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.SetFcTos:
			rule.Action = models.TCPRequestRuleActionSetDashFcDashTos
			rule.Expr = a.Expr.String()
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.SetTos:
			rule.Action = models.TCPRequestRuleActionSetDashTos
			rule.TosValue = a.Value
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.SetVarFmt:
			rule.Action = models.TCPRequestRuleActionSetDashVarDashFmt
			rule.VarName = a.VarName
			rule.VarFormat = strings.Join(a.Fmt.Expr, " ")
			rule.VarScope = a.VarScope
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.SetNice:
			nice, err := strconv.ParseInt(a.Value, 10, 64)
			if err != nil {
				return nil, err
			}
			rule.Action = models.TCPRequestRuleActionSetDashNice
			rule.NiceValue = nice
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.SetLogLevel:
			rule.Action = models.TCPRequestRuleActionSetDashLogDashLevel
			rule.LogLevel = a.Level
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *tcp_actions.SwitchMode:
			rule.Action = models.TCPRequestRuleActionSwitchDashMode
			rule.SwitchModeProto = a.Proto
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.SetRetries:
			rule.Action = models.TCPRequestRuleActionSetDashRetries
			rule.Expr = a.Expr.String()
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.DoLog:
			rule.Action = models.TCPRequestRuleActionDoDashLog
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		default:
			return nil, NewConfError(ErrValidationError, fmt.Sprintf("unsupported action '%s' in tcp_request_rule", a))
		}
	case *tcp_types.Session:
		rule = &models.TCPRequestRule{
			Type:     models.TCPRequestRuleTypeSession,
			Metadata: parseMetadata(v.Comment),
		}
		switch a := v.Action.(type) {
		case *actions.Accept:
			rule.Action = models.TCPRequestRuleActionAccept
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *tcp_actions.AttachSrv:
			rule.Action = models.TCPRequestRuleActionAttachDashSrv
			rule.ServerName = a.Server
			rule.Expr = a.Name.String()
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.Reject:
			rule.Action = models.TCPRequestRuleActionAccept
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.TrackSc:
			rule.Action = models.TCPRequestRuleActionTrackDashSc
			rule.TrackKey = a.Key
			if a.Table != "" {
				rule.TrackTable = a.Table
			}
			rule.TrackStickCounter = &a.StickCounter
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.ScAddGpc:
			rule.Action = models.TCPRequestRuleActionScDashAddDashGpc
			rule.ScIncID = a.ID
			rule.ScIdx = a.Idx
			rule.ScInt = a.Int
			rule.Expr = a.Expr.String()
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.ScIncGpc:
			rule.Action = models.TCPRequestRuleActionScDashIncDashGpc
			rule.ScIncID = a.ID
			rule.ScIdx = a.Idx
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
		case *actions.ScSetGpt:
			rule.Action = models.TCPRequestRuleActionScDashSetDashGpt
			rule.ScIncID = a.ScID
			rule.ScIdx = strconv.FormatInt(a.Idx, 10)
			rule.ScInt = a.Int
			rule.Expr = a.Expr.String()
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
		case *actions.SetMark:
			rule.Action = models.TCPRequestRuleActionSetDashMark
			rule.MarkValue = a.Value
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *tcp_actions.SetSrc:
			rule.Action = models.TCPRequestRuleActionSetDashSrc
			rule.Expr = a.Expr.String()
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.SetSrcPort:
			rule.Action = models.TCPRequestRuleActionSetDashSrcDashPort
			rule.Expr = a.Expr.String()
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.SetTos:
			rule.Action = models.TCPRequestRuleActionSetDashTos
			rule.TosValue = a.Value
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.SetFcMark:
			rule.Action = models.TCPRequestRuleActionSetDashFcDashMark
			rule.Expr = a.Expr.String()
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.SetFcTos:
			rule.Action = models.TCPRequestRuleActionSetDashFcDashTos
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
		case *actions.SetVarFmt:
			rule.Action = models.TCPRequestRuleActionSetDashVarDashFmt
			rule.VarName = a.VarName
			rule.VarFormat = strings.Join(a.Fmt.Expr, " ")
			rule.VarScope = a.VarScope
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.SilentDrop:
			rule.Action = models.TCPRequestRuleActionSilentDashDrop
			rule.RstTTL = a.RstTTL
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		case *actions.DoLog:
			rule.Action = models.TCPRequestRuleActionDoDashLog
			rule.Cond = a.Cond
			rule.CondTest = a.CondTest
		default:
			return nil, NewConfError(ErrValidationError, fmt.Sprintf("unsupported action '%s' in tcp_request_rule", a))
		}
	default:
		return nil, NewConfError(ErrValidationError, fmt.Sprintf("unsupported action '%s' in tcp_request_rule", v))
	}
	return rule, nil
}

func SerializeTCPRequestRule(f models.TCPRequestRule, opt *options.ConfigurationOptions) (types.TCPType, error) { //nolint:gocyclo,cyclop,maintidx
	comment, err := serializeMetadata(f.Metadata)
	if err != nil {
		comment = ""
	}
	switch f.Type {
	case models.TCPRequestRuleTypeConnection:
		switch f.Action {
		case models.TCPRequestRuleActionAccept:
			return &tcp_types.Connection{
				Action: &actions.Accept{
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionReject:
			return &tcp_types.Connection{
				Action: &actions.Reject{
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionExpectDashProxy:
			return &tcp_types.Connection{
				Action: &tcp_actions.ExpectProxy{
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
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
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionTrackDashSc:
			if f.TrackStickCounter == nil {
				return nil, NewConfError(ErrValidationError, "track_sc_stick_counter must be set")
			}
			return &tcp_types.Connection{
				Action: &actions.TrackSc{
					Type:         actions.TrackScType,
					StickCounter: *f.TrackStickCounter,
					Key:          f.TrackKey,
					Table:        f.TrackTable,
					Cond:         f.Cond,
					CondTest:     f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionScDashAddDashGpc:
			return &tcp_types.Connection{
				Action: &actions.ScAddGpc{
					ID:       f.ScIncID,
					Idx:      f.ScIdx,
					Int:      f.ScInt,
					Expr:     common.Expression{Expr: strings.Split(f.Expr, " ")},
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionScDashIncDashGpc:
			return &tcp_types.Connection{
				Action: &actions.ScIncGpc{
					ID:       f.ScIncID,
					Idx:      f.ScIdx,
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionScDashIncDashGpc0:
			return &tcp_types.Connection{
				Action: &actions.ScIncGpc0{
					ID:       f.ScIncID,
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionScDashIncDashGpc1:
			return &tcp_types.Connection{
				Action: &actions.ScIncGpc1{
					ID:       f.ScIncID,
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionScDashSetDashGpt:
			idx, _ := strconv.ParseInt(f.ScIdx, 10, 64)
			return &tcp_types.Connection{
				Action: &actions.ScSetGpt{
					ScID:     f.ScIncID,
					Idx:      idx,
					Int:      f.ScInt,
					Expr:     common.Expression{Expr: strings.Split(f.Expr, " ")},
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionScDashSetDashGpt0:
			return &tcp_types.Connection{
				Action: &actions.ScSetGpt0{
					ID:       f.ScIncID,
					Expr:     common.Expression{Expr: strings.Split(f.Expr, " ")},
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionSilentDashDrop:
			return &tcp_types.Connection{
				Action: &actions.SilentDrop{
					RstTTL:   f.RstTTL,
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionLua:
			return &tcp_types.Connection{
				Action: &actions.Lua{
					Action:   f.LuaAction,
					Params:   f.LuaParams,
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionSetDashMark:
			return &tcp_types.Connection{
				Action: &actions.SetMark{
					Value:    f.MarkValue,
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionSetDashSrcDashPort:
			return &tcp_types.Connection{
				Action: &actions.SetSrcPort{
					Expr:     common.Expression{Expr: strings.Split(f.Expr, " ")},
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionSetDashTos:
			return &tcp_types.Connection{
				Action: &actions.SetTos{
					Value:    f.TosValue,
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionSetDashVar:
			return &tcp_types.Connection{
				Action: &actions.SetVar{
					VarName:  f.VarName,
					VarScope: f.VarScope,
					Expr:     common.Expression{Expr: strings.Split(f.Expr, " ")},
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionUnsetDashVar:
			return &tcp_types.Connection{
				Action: &actions.UnsetVar{
					Name:     f.VarName,
					Scope:    f.VarScope,
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionSetDashVarDashFmt:
			return &tcp_types.Connection{
				Action: &actions.SetVarFmt{
					Fmt:      common.Expression{Expr: strings.Split(f.VarFormat, " ")},
					VarName:  f.VarName,
					VarScope: f.VarScope,
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionSetDashSrc:
			return &tcp_types.Connection{
				Action: &http_actions.SetSrc{
					Expr:     common.Expression{Expr: strings.Split(f.Expr, " ")},
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionSetDashDst:
			return &tcp_types.Connection{
				Action: &actions.SetDst{
					Expr:     common.Expression{Expr: strings.Split(f.Expr, " ")},
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionSetDashDstDashPort:
			return &tcp_types.Connection{
				Action: &actions.SetDstPort{
					Expr:     common.Expression{Expr: strings.Split(f.Expr, " ")},
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionSetDashFcDashMark:
			return &tcp_types.Connection{
				Action: &actions.SetFcMark{
					Expr:     common.Expression{Expr: strings.Split(f.Expr+f.MarkValue, " ")},
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionSetDashFcDashTos:
			return &tcp_types.Connection{
				Action: &actions.SetFcTos{
					Expr:     common.Expression{Expr: strings.Split(f.Expr+f.TosValue, " ")},
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionDoDashLog:
			return &tcp_types.Connection{
				Action: &actions.DoLog{
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		}
		return nil, NewConfError(ErrValidationError, fmt.Sprintf("unsupported action '%s' in tcp_request_rule", f.Action))
	case models.TCPRequestRuleTypeContent:
		switch f.Action {
		case models.TCPRequestRuleActionAccept:
			return &tcp_types.Content{
				Action: &actions.Accept{
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
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
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionReject:
			return &tcp_types.Content{
				Action: &actions.Reject{
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionCapture:
			return &tcp_types.Content{
				Action: &tcp_actions.Capture{
					Expr:     common.Expression{Expr: strings.Split(f.Expr, " ")},
					Len:      f.CaptureLen,
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionSetDashPriorityDashClass:
			return &tcp_types.Content{
				Action: &actions.SetPriorityClass{
					Expr:     common.Expression{Expr: strings.Split(f.Expr, " ")},
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionSetDashPriorityDashOffset:
			return &tcp_types.Content{
				Action: &actions.SetPriorityOffset{
					Expr:     common.Expression{Expr: strings.Split(f.Expr, " ")},
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionTrackDashSc:
			if f.TrackStickCounter == nil {
				return nil, NewConfError(ErrValidationError, "track_sc_stick_counter must be set")
			}
			return &tcp_types.Content{
				Action: &actions.TrackSc{
					Type:         actions.TrackScType,
					StickCounter: *f.TrackStickCounter,
					Key:          f.TrackKey,
					Table:        f.TrackTable,
					Cond:         f.Cond,
					CondTest:     f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionScDashAddDashGpc:
			return &tcp_types.Content{
				Action: &actions.ScAddGpc{
					ID:       f.ScIncID,
					Idx:      f.ScIdx,
					Int:      f.ScInt,
					Expr:     common.Expression{Expr: strings.Split(f.Expr, " ")},
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionScDashIncDashGpc:
			return &tcp_types.Content{
				Action: &actions.ScIncGpc{
					ID:       f.ScIncID,
					Idx:      f.ScIdx,
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionScDashIncDashGpc0:
			return &tcp_types.Content{
				Action: &actions.ScIncGpc0{
					ID:       f.ScIncID,
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionScDashIncDashGpc1:
			return &tcp_types.Content{
				Action: &actions.ScIncGpc1{
					ID:       f.ScIncID,
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionScDashSetDashGpt:
			idx, _ := strconv.ParseInt(f.ScIdx, 10, 64)
			return &tcp_types.Content{
				Action: &actions.ScSetGpt{
					ScID:     f.ScIncID,
					Idx:      idx,
					Int:      f.ScInt,
					Expr:     common.Expression{Expr: strings.Split(f.Expr, " ")},
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionScDashSetDashGpt0:
			return &tcp_types.Content{
				Action: &actions.ScSetGpt0{
					ID:       f.ScIncID,
					Expr:     common.Expression{Expr: strings.Split(f.Expr, " ")},
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionSetDashDst:
			return &tcp_types.Content{
				Action: &actions.SetDst{
					Expr:     common.Expression{Expr: strings.Split(f.Expr, " ")},
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionSetDashDstDashPort:
			return &tcp_types.Content{
				Action: &actions.SetDstPort{
					Expr:     common.Expression{Expr: strings.Split(f.Expr, " ")},
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionSetDashSrc:
			return &tcp_types.Content{
				Action: &http_actions.SetSrc{
					Expr:     common.Expression{Expr: strings.Split(f.Expr, " ")},
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
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
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionUnsetDashVar:
			return &tcp_types.Content{
				Action: &actions.UnsetVar{
					Name:     f.VarName,
					Scope:    f.VarScope,
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionSilentDashDrop:
			return &tcp_types.Content{
				Action: &actions.SilentDrop{
					RstTTL:   f.RstTTL,
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionSendDashSpoeDashGroup:
			return &tcp_types.Content{
				Action: &actions.SendSpoeGroup{
					Engine:   f.SpoeEngineName,
					Group:    f.SpoeGroupName,
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionUseDashService:
			return &tcp_types.Content{
				Action: &actions.UseService{
					Name:     f.ServiceName,
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionLua:
			return &tcp_types.Content{
				Action: &actions.Lua{
					Action:   f.LuaAction,
					Params:   f.LuaParams,
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionSetDashBandwidthDashLimit:
			return &tcp_types.Content{
				Action: &actions.SetBandwidthLimit{
					Name:     f.BandwidthLimitName,
					Limit:    common.Expression{Expr: strings.Split(f.BandwidthLimitLimit, " ")},
					Period:   common.Expression{Expr: strings.Split(f.BandwidthLimitPeriod, " ")},
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionSetDashMark:
			return &tcp_types.Content{
				Action: &actions.SetMark{
					Value:    f.MarkValue,
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionSetDashSrcDashPort:
			return &tcp_types.Content{
				Action: &actions.SetSrcPort{
					Expr:     common.Expression{Expr: strings.Split(f.Expr, " ")},
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionSetDashTos:
			return &tcp_types.Content{
				Action: &actions.SetTos{
					Value:    f.TosValue,
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionSetDashVarDashFmt:
			return &tcp_types.Content{
				Action: &actions.SetVarFmt{
					Fmt:      common.Expression{Expr: strings.Split(f.VarFormat, " ")},
					VarName:  f.VarName,
					VarScope: f.VarScope,
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionSetDashNice:
			return &tcp_types.Content{
				Action: &actions.SetNice{
					Value:    strconv.FormatInt(f.NiceValue, 10),
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionSetDashLogDashLevel:
			return &tcp_types.Content{
				Action: &actions.SetLogLevel{
					Level:    f.LogLevel,
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionSwitchDashMode:
			return &tcp_types.Content{
				Action: &tcp_actions.SwitchMode{
					Proto:    f.SwitchModeProto,
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionSetDashBcDashMark:
			return &tcp_types.Content{
				Action: &actions.SetBcMark{
					Expr:     common.Expression{Expr: strings.Split(f.Expr+f.MarkValue, " ")},
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionSetDashBcDashTos:
			return &tcp_types.Content{
				Action: &actions.SetBcTos{
					Expr:     common.Expression{Expr: strings.Split(f.Expr+f.TosValue, " ")},
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionSetDashFcDashMark:
			return &tcp_types.Content{
				Action: &actions.SetFcMark{
					Expr:     common.Expression{Expr: strings.Split(f.Expr+f.MarkValue, " ")},
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionSetDashFcDashTos:
			return &tcp_types.Content{
				Action: &actions.SetFcTos{
					Expr:     common.Expression{Expr: strings.Split(f.Expr+f.TosValue, " ")},
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionSetDashRetries:
			return &tcp_types.Content{
				Action: &actions.SetRetries{
					Expr:     common.Expression{Expr: strings.Split(f.Expr, " ")},
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionDoDashLog:
			return &tcp_types.Content{
				Action: &actions.DoLog{
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		}
		return nil, NewConfError(ErrValidationError, fmt.Sprintf("unsupported action '%s' in tcp_request_rule", f.Action))
	case models.TCPRequestRuleTypeSession:
		switch f.Action {
		case models.TCPRequestRuleActionAccept:
			return &tcp_types.Session{
				Action: &actions.Accept{
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionAttachDashSrv:
			return &tcp_types.Session{
				Action: &tcp_actions.AttachSrv{
					Server:   f.ServerName,
					Name:     common.Expression{Expr: []string{f.Expr}},
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionReject:
			return &tcp_types.Session{
				Action: &actions.Reject{
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionTrackDashSc:
			if f.TrackStickCounter == nil {
				return nil, NewConfError(ErrValidationError, "track_sc_stick_counter must be set")
			}
			return &tcp_types.Session{
				Action: &actions.TrackSc{
					Type:         actions.TrackScType,
					StickCounter: *f.TrackStickCounter,
					Key:          f.TrackKey,
					Table:        f.TrackTable,
					Cond:         f.Cond,
					CondTest:     f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionScDashAddDashGpc:
			return &tcp_types.Connection{
				Action: &actions.ScAddGpc{
					ID:       f.ScIncID,
					Idx:      f.ScIdx,
					Int:      f.ScInt,
					Expr:     common.Expression{Expr: strings.Split(f.Expr, " ")},
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionScDashIncDashGpc:
			return &tcp_types.Connection{
				Action: &actions.ScIncGpc{
					ID:       f.ScIncID,
					Idx:      f.ScIdx,
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionScDashIncDashGpc0:
			return &tcp_types.Session{
				Action: &actions.ScIncGpc0{
					ID:       f.ScIncID,
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionScDashIncDashGpc1:
			return &tcp_types.Session{
				Action: &actions.ScIncGpc1{
					ID:       f.ScIncID,
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionScDashSetDashGpt:
			idx, _ := strconv.ParseInt(f.ScIdx, 10, 64)
			return &tcp_types.Session{
				Action: &actions.ScSetGpt{
					ScID:     f.ScIncID,
					Idx:      idx,
					Int:      f.ScInt,
					Expr:     common.Expression{Expr: strings.Split(f.Expr, " ")},
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionScDashSetDashGpt0:
			return &tcp_types.Session{
				Action: &actions.ScSetGpt0{
					ID:       f.ScIncID,
					Expr:     common.Expression{Expr: strings.Split(f.Expr, " ")},
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionSetDashDst:
			return &tcp_types.Session{
				Action: &actions.SetDst{
					Expr:     common.Expression{Expr: strings.Split(f.Expr, " ")},
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionSetDashDstDashPort:
			return &tcp_types.Session{
				Action: &actions.SetDstPort{
					Expr:     common.Expression{Expr: strings.Split(f.Expr, " ")},
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionSetDashSrc:
			return &tcp_types.Session{
				Action: &http_actions.SetSrc{
					Expr:     common.Expression{Expr: strings.Split(f.Expr, " ")},
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionSetDashSrcDashPort:
			return &tcp_types.Session{
				Action: &actions.SetSrcPort{
					Expr:     common.Expression{Expr: strings.Split(f.Expr, " ")},
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionSetDashMark:
			return &tcp_types.Session{
				Action: &actions.SetDstPort{
					Expr:     common.Expression{Expr: strings.Split(f.Expr, " ")},
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionSetDashTos:
			return &tcp_types.Session{
				Action: &actions.SetTos{
					Value:    f.TosValue,
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
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
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionUnsetDashVar:
			return &tcp_types.Session{
				Action: &actions.UnsetVar{
					Name:     f.VarName,
					Scope:    f.VarScope,
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionSetDashVarDashFmt:
			return &tcp_types.Session{
				Action: &actions.SetVarFmt{
					Fmt:      common.Expression{Expr: strings.Split(f.VarFormat, " ")},
					VarName:  f.VarName,
					VarScope: f.VarScope,
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionSilentDashDrop:
			return &tcp_types.Session{
				Action: &actions.SilentDrop{
					RstTTL:   f.RstTTL,
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionSetDashFcDashMark:
			return &tcp_types.Session{
				Action: &actions.SetFcMark{
					Expr:     common.Expression{Expr: strings.Split(f.Expr+f.MarkValue, " ")},
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionSetDashFcDashTos:
			return &tcp_types.Session{
				Action: &actions.SetFcTos{
					Expr:     common.Expression{Expr: strings.Split(f.Expr+f.TosValue, " ")},
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		case models.TCPRequestRuleActionDoDashLog:
			return &tcp_types.Session{
				Action: &actions.DoLog{
					Cond:     f.Cond,
					CondTest: f.CondTest,
				},
				Comment: comment,
			}, nil
		}
		return nil, NewConfError(ErrValidationError, fmt.Sprintf("unsupported action '%s' in tcp_request_rule", f.Action))
	case models.TCPRequestRuleTypeInspectDashDelay:
		if f.Timeout == nil {
			return nil, NewConfError(ErrValidationError, fmt.Sprintf("unsupported action '%s' in tcp_request_rule", f.Type))
		}
		return &tcp_types.InspectDelay{
			Timeout: misc.SerializeTime(*f.Timeout, opt.PreferredTimeSuffix),
			Comment: comment,
		}, nil
	}
	return nil, NewConfError(ErrValidationError, fmt.Sprintf("unsupported action '%s' in tcp_request_rule", f.Type))
}
