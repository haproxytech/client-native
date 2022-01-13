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

	strfmt "github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/config-parser/v4"
	"github.com/haproxytech/config-parser/v4/common"
	parser_errors "github.com/haproxytech/config-parser/v4/errors"
	actions "github.com/haproxytech/config-parser/v4/parsers/actions"
	tcp_actions "github.com/haproxytech/config-parser/v4/parsers/tcp/actions"
	"github.com/haproxytech/config-parser/v4/types"

	"github.com/haproxytech/client-native/v2/models"
)

// GetTCPChecks returns configuration version and an array of configured tcp-checks in the specified parent.
// Returns error on fail.
func (c *Client) GetTCPChecks(parentType, parentName string, transactionID string) (int64, models.TCPChecks, error) {
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
func (c *Client) GetTCPCheck(id int64, parentType string, parentName string, transactionID string) (int64, *models.TCPCheck, error) {
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
func (c *Client) DeleteTCPCheck(id int64, parentType string, parentName string, transactionID string, version int64) error {
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
func (c *Client) CreateTCPCheck(parentType string, parentName string, data *models.TCPCheck, transactionID string, version int64) error {
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
func (c *Client) EditTCPCheck(id int64, parentType string, parentName string, data *models.TCPCheck, transactionID string, version int64) error {
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
	section := parser.Global
	if t == "frontend" {
		section = parser.Frontends
	} else if t == "backend" {
		section = parser.Backends
	}

	checks := models.TCPChecks{}
	data, err := p.Get(section, pName, "tcp-check", false)
	if err != nil {
		if errors.Is(err, parser_errors.ErrFetch) {
			return checks, nil
		}
		return nil, err
	}
	items := data.([]types.Action)
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

func ParseTCPCheck(f types.Action) (check *models.TCPCheck, err error) { //nolint:gocyclo

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
			Action:       models.TCPCheckActionSendLf,
			Fmt:          v.Fmt,
			CheckComment: v.CheckComment,
		}
	case *tcp_actions.CheckSendBinary:
		check = &models.TCPCheck{
			Action:       models.TCPCheckActionSendBinary,
			HexString:    v.HexString,
			CheckComment: v.CheckComment,
		}
	case *tcp_actions.CheckSendBinaryLf:
		check = &models.TCPCheck{
			Action:       models.TCPCheckActionSendBinaryLf,
			HexFmt:       v.HexFmt,
			CheckComment: v.CheckComment,
		}
	case *actions.SetVarCheck:
		check = &models.TCPCheck{
			Action:   models.TCPCheckActionSetVar,
			VarScope: v.VarScope,
			VarName:  v.VarName,
			VarExpr:  strings.Join(v.Expr.Expr, " "),
		}
	case *tcp_actions.CheckSetVarFmt:
		check = &models.TCPCheck{
			Action:   models.TCPCheckActionSetVarFmt,
			VarScope: v.VarScope,
			VarName:  v.VarName,
			VarFmt:   v.Format,
		}
	case *actions.UnsetVarCheck:
		check = &models.TCPCheck{
			Action:   models.TCPCheckActionUnsetVar,
			VarScope: v.Scope,
			VarName:  v.Name,
		}
	}

	return check, nil
}

func SerializeTCPCheck(f models.TCPCheck) (action types.Action, err error) { //nolint:gocyclo

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
		return &actions.CheckExpect{
			MinRecv:         &f.MinRecv,
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
		}, nil
	case models.TCPCheckActionSend:
		return &tcp_actions.CheckSend{
			Data:         f.Data,
			CheckComment: f.CheckComment,
		}, nil
	case models.TCPCheckActionSendLf:
		return &tcp_actions.CheckSendLf{
			Fmt:          f.Fmt,
			CheckComment: f.CheckComment,
		}, nil
	case models.TCPCheckActionSendBinary:
		return &tcp_actions.CheckSendBinary{
			HexString:    f.HexString,
			CheckComment: f.CheckComment,
		}, nil
	case models.TCPCheckActionSendBinaryLf:
		return &tcp_actions.CheckSendBinaryLf{
			HexFmt:       f.HexFmt,
			CheckComment: f.CheckComment,
		}, nil
	case models.TCPCheckActionSetVar:
		return &actions.SetVarCheck{
			VarScope: f.VarScope,
			VarName:  f.VarName,
			Expr:     common.Expression{Expr: strings.Split(f.VarExpr, " ")},
		}, nil
	case models.TCPCheckActionSetVarFmt:
		return &tcp_actions.CheckSetVarFmt{
			VarScope: f.VarScope,
			VarName:  f.VarName,
			Format:   f.VarFmt,
		}, nil
	case models.TCPCheckActionUnsetVar:
		return &actions.UnsetVarCheck{
			Scope: f.VarScope,
			Name:  f.VarName,
		}, nil
	}

	return nil, NewConfError(ErrValidationError, "unsupported action in tcp_check")
}
