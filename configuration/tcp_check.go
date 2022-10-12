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

	strfmt "github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/config-parser/v4"
	"github.com/haproxytech/config-parser/v4/common"
	parser_errors "github.com/haproxytech/config-parser/v4/errors"
	actions "github.com/haproxytech/config-parser/v4/parsers/actions"
	tcp_actions "github.com/haproxytech/config-parser/v4/parsers/tcp/actions"
	"github.com/haproxytech/config-parser/v4/types"

	"github.com/haproxytech/client-native/v5/misc"
	"github.com/haproxytech/client-native/v5/models"
)

type TCPCheck interface {
	GetTCPChecks(parentType, parentName string, transactionID string) (int64, models.TCPChecks, error)
	GetTCPCheck(id int64, parentType string, parentName string, transactionID string) (int64, *models.TCPCheck, error)
	DeleteTCPCheck(id int64, parentType string, parentName string, transactionID string, version int64) error
	CreateTCPCheck(parentType string, parentName string, data *models.TCPCheck, transactionID string, version int64) error
	EditTCPCheck(id int64, parentType string, parentName string, data *models.TCPCheck, transactionID string, version int64) error
}

// GetTCPChecks returns configuration version and an array of configured tcp-checks in the specified parent.
// Returns error on fail.
func (c *client) GetTCPChecks(parentType, parentName string, transactionID string) (int64, models.TCPChecks, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	tcpRules, err := ParseTCPChecks(parentType, parentName, p)
	if err != nil {
		return v, nil, c.HandleError("", parentType, parentName, "", false, err)
	}

	return v, tcpRules, nil
}

// GetTCPCheck returns configuration version and the requested tcp check in the specified parent.
// Returns error on fail or if tcp check does not exist
func (c *client) GetTCPCheck(id int64, parentType string, parentName string, transactionID string) (int64, *models.TCPCheck, error) {
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
	} else if parentType == "defaults" {
		section = parser.Defaults
		parentName = parser.DefaultSectionName
	}

	data, err := p.GetOne(section, parentName, "tcp-check", int(id))
	if err != nil {
		return v, nil, c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, "", false, err)
	}

	tcpCheck, err := ParseTCPCheck(data.(types.Action))
	if err != nil {
		return v, nil, err
	}
	tcpCheck.Index = &id
	return v, tcpCheck, nil
}

// DeleteTCPCheck deletes a tcp check in the configuration. One of version or transactionID is mandatory.
// Returns error on fail, nil on success.
func (c *client) DeleteTCPCheck(id int64, parentType string, parentName string, transactionID string, version int64) error {
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	var section parser.Section
	if parentType == "backend" {
		section = parser.Backends
	} else if parentType == "defaults" {
		section = parser.Defaults
		parentName = parser.DefaultSectionName
	}

	if err := p.Delete(section, parentName, "tcp-check", int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	if err := c.SaveData(p, t, transactionID == ""); err != nil {
		return err
	}
	return nil
}

// CreateTCPCheck creates a tcp check in the configuration. One of version or transationID is mandatory.
// Returns error on fail, nil on success.
func (c *client) CreateTCPCheck(parentType string, parentName string, data *models.TCPCheck, transactionID string, version int64) error {
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
	} else if parentType == "defaults" {
		section = parser.Defaults
		parentName = parser.DefaultSectionName
	}

	s, err := SerializeTCPCheck(*data)
	if err != nil {
		return err
	}

	if err := p.Insert(section, parentName, "tcp-check", s, int(*data.Index)); err != nil {
		return c.HandleError(strconv.FormatInt(*data.Index, 10), parentType, parentName, t, transactionID == "", err)
	}
	if err := c.SaveData(p, t, transactionID == ""); err != nil {
		return err
	}
	return nil
}

// EditTCPCheck edits a tcp check in the configuration. One of version or transactionID is mandatory.
// Returns error on fail, nil on success.
// nolint:dupl
func (c *client) EditTCPCheck(id int64, parentType string, parentName string, data *models.TCPCheck, transactionID string, version int64) error {
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
	} else if parentType == "defaults" {
		section = parser.Defaults
		parentName = parser.DefaultSectionName
	}

	if _, err = p.GetOne(section, parentName, "tcp-check", int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	s, err := SerializeTCPCheck(*data)
	if err != nil {
		return err
	}

	if err := p.Set(section, parentName, "tcp-check", s, int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}
	if err := c.SaveData(p, t, transactionID == ""); err != nil {
		return err
	}
	return nil
}

func ParseTCPChecks(t, pName string, p parser.Parser) (models.TCPChecks, error) {
	var section parser.Section
	switch t {
	case "defaults":
		section = parser.Defaults
		pName = parser.DefaultSectionName
	case "backend":
		section = parser.Backends
	default:
		return nil, NewConfError(ErrValidationError, fmt.Sprintf("unsupported section in tcp_check: %s", t))
	}

	checks := models.TCPChecks{}
	data, err := p.Get(section, pName, "tcp-check", false)
	if err != nil {
		if errors.Is(err, parser_errors.ErrFetch) {
			return checks, nil
		}
		return nil, err
	}
	items, ok := data.([]types.Action)
	if !ok {
		return nil, misc.CreateTypeAssertError("tcp-check")
	}
	for i, c := range items {
		id := int64(i)
		check, err := ParseTCPCheck(c)
		if err == nil {
			check.Index = &id
			checks = append(checks, check)
		}
	}
	return checks, nil
}

func ParseTCPCheck(f types.Action) (check *models.TCPCheck, err error) {
	switch v := f.(type) {
	case *tcp_actions.CheckComment:
		check = &models.TCPCheck{
			Action:     models.TCPCheckActionComment,
			LogMessage: v.LogMessage,
		}
	case *actions.CheckConnect:
		check = &models.TCPCheck{
			Action:       models.TCPCheckActionConnect,
			PortString:   v.Port,
			Addr:         v.Addr,
			Sni:          v.SNI,
			Alpn:         v.ALPN,
			Proto:        v.Proto,
			CheckComment: v.CheckComment,
			Default:      v.Default,
			SendProxy:    v.SendProxy,
			ViaSocks4:    v.ViaSOCKS4,
			Ssl:          v.SSL,
			Linger:       v.Linger,
		}
	case *actions.CheckExpect:
		check = &models.TCPCheck{
			Action:          models.TCPCheckActionExpect,
			CheckComment:    v.CheckComment,
			OkStatus:        v.OKStatus,
			ErrorStatus:     v.ErrorStatus,
			ToutStatus:      v.TimeoutStatus,
			OnSuccess:       v.OnSuccess,
			OnError:         v.OnError,
			StatusCode:      v.StatusCode,
			ExclamationMark: v.ExclamationMark,
			Match:           v.Match,
			Pattern:         v.Pattern,
		}
		if v.MinRecv != nil {
			check.MinRecv = *v.MinRecv
		}
	case *tcp_actions.CheckSend:
		check = &models.TCPCheck{
			Action:       models.TCPCheckActionSend,
			Data:         v.Data,
			CheckComment: v.CheckComment,
		}
	case *tcp_actions.CheckSendLf:
		check = &models.TCPCheck{
			Action:       models.TCPCheckActionSendDashLf,
			Fmt:          v.Fmt,
			CheckComment: v.CheckComment,
		}
	case *tcp_actions.CheckSendBinary:
		check = &models.TCPCheck{
			Action:       models.TCPCheckActionSendDashBinary,
			HexString:    v.HexString,
			CheckComment: v.CheckComment,
		}
	case *tcp_actions.CheckSendBinaryLf:
		check = &models.TCPCheck{
			Action:       models.TCPCheckActionSendDashBinaryDashLf,
			HexFmt:       v.HexFmt,
			CheckComment: v.CheckComment,
		}
	case *actions.SetVarCheck:
		check = &models.TCPCheck{
			Action:   models.TCPCheckActionSetDashVar,
			VarScope: v.VarScope,
			VarName:  v.VarName,
			VarExpr:  strings.Join(v.Expr.Expr, " "),
		}
	case *actions.SetVarFmtCheck:
		check = &models.TCPCheck{
			Action:   models.TCPCheckActionSetDashVarDashFmt,
			VarScope: v.VarScope,
			VarName:  v.VarName,
			VarFmt:   strings.Join(v.Format.Expr, " "),
		}
	case *actions.UnsetVarCheck:
		check = &models.TCPCheck{
			Action:   models.TCPCheckActionUnsetDashVar,
			VarScope: v.Scope,
			VarName:  v.Name,
		}
	}

	return check, nil
}

func SerializeTCPCheck(f models.TCPCheck) (action types.Action, err error) { //nolint:ireturn
	switch f.Action {
	case models.TCPCheckActionComment:
		return &tcp_actions.CheckComment{
			LogMessage: f.LogMessage,
		}, nil
	case models.TCPCheckActionConnect:
		return &actions.CheckConnect{
			Port:         f.PortString,
			Addr:         f.Addr,
			SNI:          f.Sni,
			ALPN:         f.Alpn,
			Proto:        f.Proto,
			CheckComment: f.CheckComment,
			Default:      f.Default,
			SendProxy:    f.SendProxy,
			ViaSOCKS4:    f.ViaSocks4,
			SSL:          f.Ssl,
			Linger:       f.Linger,
		}, nil
	case models.TCPCheckActionExpect:
		action := &actions.CheckExpect{
			Match:           f.Match,
			OKStatus:        f.OkStatus,
			ErrorStatus:     f.ErrorStatus,
			CheckComment:    f.CheckComment,
			TimeoutStatus:   f.ToutStatus,
			OnSuccess:       f.OnSuccess,
			OnError:         f.OnError,
			StatusCode:      f.StatusCode,
			ExclamationMark: f.ExclamationMark,
			Pattern:         f.Pattern,
		}
		if f.MinRecv > 0 {
			action.MinRecv = &f.MinRecv
		}
		return action, nil
	case models.TCPCheckActionSend:
		return &tcp_actions.CheckSend{
			Data:         f.Data,
			CheckComment: f.CheckComment,
		}, nil
	case models.TCPCheckActionSendDashLf:
		return &tcp_actions.CheckSendLf{
			Fmt:          f.Fmt,
			CheckComment: f.CheckComment,
		}, nil
	case models.TCPCheckActionSendDashBinary:
		return &tcp_actions.CheckSendBinary{
			HexString:    f.HexString,
			CheckComment: f.CheckComment,
		}, nil
	case models.TCPCheckActionSendDashBinaryDashLf:
		return &tcp_actions.CheckSendBinaryLf{
			HexFmt:       f.HexFmt,
			CheckComment: f.CheckComment,
		}, nil
	case models.TCPCheckActionSetDashVar:
		return &actions.SetVarCheck{
			VarScope: f.VarScope,
			VarName:  f.VarName,
			Expr:     common.Expression{Expr: strings.Split(f.VarExpr, " ")},
		}, nil
	case models.TCPCheckActionSetDashVarDashFmt:
		return &actions.SetVarFmtCheck{
			VarScope: f.VarScope,
			VarName:  f.VarName,
			Format:   common.Expression{Expr: strings.Split(f.VarFmt, " ")},
		}, nil
	case models.TCPCheckActionUnsetDashVar:
		return &actions.UnsetVarCheck{
			Scope: f.VarScope,
			Name:  f.VarName,
		}, nil
	}

	return nil, NewConfError(ErrValidationError, "unsupported action in tcp_check")
}
