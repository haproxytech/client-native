// Code generated by go-swagger; DO NOT EDIT.

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

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"encoding/json"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// QUICInitialRule QUIC Initial
//
// QUIC Initial configuration
// Example: {"type":"reject"}
//
// swagger:model QUICInitialRule
type QUICInitialRule struct {
	// cond
	// Enum: ["if","unless"]
	// +kubebuilder:validation:Enum="if","unless";
	Cond string `json:"cond,omitempty"`

	// cond test
	CondTest string `json:"cond_test,omitempty"`

	// type
	// Required: true
	// Enum: ["reject","accept","send-retry","dgram-drop"]
	// +kubebuilder:validation:Enum="reject","accept","send-retry","dgram-drop";
	Type string `json:"type"`
}

// Validate validates this q UI c initial rule
func (m *QUICInitialRule) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateCond(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateType(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

var qUiCInitialRuleTypeCondPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["if","unless"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		qUiCInitialRuleTypeCondPropEnum = append(qUiCInitialRuleTypeCondPropEnum, v)
	}
}

const (

	// QUICInitialRuleCondIf captures enum value "if"
	QUICInitialRuleCondIf string = "if"

	// QUICInitialRuleCondUnless captures enum value "unless"
	QUICInitialRuleCondUnless string = "unless"
)

// prop value enum
func (m *QUICInitialRule) validateCondEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, qUiCInitialRuleTypeCondPropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *QUICInitialRule) validateCond(formats strfmt.Registry) error {
	if swag.IsZero(m.Cond) { // not required
		return nil
	}

	// value enum
	if err := m.validateCondEnum("cond", "body", m.Cond); err != nil {
		return err
	}

	return nil
}

var qUiCInitialRuleTypeTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["reject","accept","send-retry","dgram-drop"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		qUiCInitialRuleTypeTypePropEnum = append(qUiCInitialRuleTypeTypePropEnum, v)
	}
}

const (

	// QUICInitialRuleTypeReject captures enum value "reject"
	QUICInitialRuleTypeReject string = "reject"

	// QUICInitialRuleTypeAccept captures enum value "accept"
	QUICInitialRuleTypeAccept string = "accept"

	// QUICInitialRuleTypeSendDashRetry captures enum value "send-retry"
	QUICInitialRuleTypeSendDashRetry string = "send-retry"

	// QUICInitialRuleTypeDgramDashDrop captures enum value "dgram-drop"
	QUICInitialRuleTypeDgramDashDrop string = "dgram-drop"
)

// prop value enum
func (m *QUICInitialRule) validateTypeEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, qUiCInitialRuleTypeTypePropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *QUICInitialRule) validateType(formats strfmt.Registry) error {

	if err := validate.RequiredString("type", "body", m.Type); err != nil {
		return err
	}

	// value enum
	if err := m.validateTypeEnum("type", "body", m.Type); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this q UI c initial rule based on context it is used
func (m *QUICInitialRule) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *QUICInitialRule) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *QUICInitialRule) UnmarshalBinary(b []byte) error {
	var res QUICInitialRule
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}