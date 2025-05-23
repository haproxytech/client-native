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

// EmailAlert Email Alert
//
// Send emails for important log messages.
//
// swagger:model email_alert
type EmailAlert struct {
	// from
	// Required: true
	// Pattern: ^\S+@\S+$
	// +kubebuilder:validation:Pattern=`^\S+@\S+$`
	From *string `json:"from"`

	// level
	// Enum: ["emerg","alert","crit","err","warning","notice","info","debug"]
	// +kubebuilder:validation:Enum=emerg;alert;crit;err;warning;notice;info;debug;
	Level string `json:"level,omitempty"`

	// mailers
	// Required: true
	Mailers *string `json:"mailers"`

	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:Schemaless
	Metadata map[string]interface{} `json:"metadata,omitempty"`

	// myhostname
	Myhostname string `json:"myhostname,omitempty"`

	// to
	// Required: true
	// Pattern: ^\S+@\S+$
	// +kubebuilder:validation:Pattern=`^\S+@\S+$`
	To *string `json:"to"`
}

// Validate validates this email alert
func (m *EmailAlert) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateFrom(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateLevel(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateMailers(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTo(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *EmailAlert) validateFrom(formats strfmt.Registry) error {

	if err := validate.Required("from", "body", m.From); err != nil {
		return err
	}

	if err := validate.Pattern("from", "body", *m.From, `^\S+@\S+$`); err != nil {
		return err
	}

	return nil
}

var emailAlertTypeLevelPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["emerg","alert","crit","err","warning","notice","info","debug"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		emailAlertTypeLevelPropEnum = append(emailAlertTypeLevelPropEnum, v)
	}
}

const (

	// EmailAlertLevelEmerg captures enum value "emerg"
	EmailAlertLevelEmerg string = "emerg"

	// EmailAlertLevelAlert captures enum value "alert"
	EmailAlertLevelAlert string = "alert"

	// EmailAlertLevelCrit captures enum value "crit"
	EmailAlertLevelCrit string = "crit"

	// EmailAlertLevelErr captures enum value "err"
	EmailAlertLevelErr string = "err"

	// EmailAlertLevelWarning captures enum value "warning"
	EmailAlertLevelWarning string = "warning"

	// EmailAlertLevelNotice captures enum value "notice"
	EmailAlertLevelNotice string = "notice"

	// EmailAlertLevelInfo captures enum value "info"
	EmailAlertLevelInfo string = "info"

	// EmailAlertLevelDebug captures enum value "debug"
	EmailAlertLevelDebug string = "debug"
)

// prop value enum
func (m *EmailAlert) validateLevelEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, emailAlertTypeLevelPropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *EmailAlert) validateLevel(formats strfmt.Registry) error {
	if swag.IsZero(m.Level) { // not required
		return nil
	}

	// value enum
	if err := m.validateLevelEnum("level", "body", m.Level); err != nil {
		return err
	}

	return nil
}

func (m *EmailAlert) validateMailers(formats strfmt.Registry) error {

	if err := validate.Required("mailers", "body", m.Mailers); err != nil {
		return err
	}

	return nil
}

func (m *EmailAlert) validateTo(formats strfmt.Registry) error {

	if err := validate.Required("to", "body", m.To); err != nil {
		return err
	}

	if err := validate.Pattern("to", "body", *m.To, `^\S+@\S+$`); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this email alert based on context it is used
func (m *EmailAlert) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *EmailAlert) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *EmailAlert) UnmarshalBinary(b []byte) error {
	var res EmailAlert
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
