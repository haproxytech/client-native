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

// OcspUpdateOptions ocsp update options
//
// swagger:model ocsp_update_options
type OcspUpdateOptions struct {

	// disable
	Disable *bool `json:"disable,omitempty"`

	// httpproxy
	Httpproxy *OcspUpdateOptionsHttpproxy `json:"httpproxy,omitempty"`

	// Sets the maximum interval between two automatic updates of the same OCSP response.This time is expressed in seconds
	Maxdelay *int64 `json:"maxdelay,omitempty"`

	// Sets the minimum interval between two automatic updates of the same OCSP response. This time is expressed in seconds
	Mindelay *int64 `json:"mindelay,omitempty"`

	// mode
	// Enum: ["enabled","disabled"]
	// +kubebuilder:validation:Enum=enabled;disabled;
	Mode string `json:"mode,omitempty"`
}

// Validate validates this ocsp update options
func (m *OcspUpdateOptions) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateHttpproxy(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateMode(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *OcspUpdateOptions) validateHttpproxy(formats strfmt.Registry) error {
	if swag.IsZero(m.Httpproxy) { // not required
		return nil
	}

	if m.Httpproxy != nil {
		if err := m.Httpproxy.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("httpproxy")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("httpproxy")
			}
			return err
		}
	}

	return nil
}

var ocspUpdateOptionsTypeModePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["enabled","disabled"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		ocspUpdateOptionsTypeModePropEnum = append(ocspUpdateOptionsTypeModePropEnum, v)
	}
}

const (

	// OcspUpdateOptionsModeEnabled captures enum value "enabled"
	OcspUpdateOptionsModeEnabled string = "enabled"

	// OcspUpdateOptionsModeDisabled captures enum value "disabled"
	OcspUpdateOptionsModeDisabled string = "disabled"
)

// prop value enum
func (m *OcspUpdateOptions) validateModeEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, ocspUpdateOptionsTypeModePropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *OcspUpdateOptions) validateMode(formats strfmt.Registry) error {
	if swag.IsZero(m.Mode) { // not required
		return nil
	}

	// value enum
	if err := m.validateModeEnum("mode", "body", m.Mode); err != nil {
		return err
	}

	return nil
}

// ContextValidate validate this ocsp update options based on the context it is used
func (m *OcspUpdateOptions) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateHttpproxy(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *OcspUpdateOptions) contextValidateHttpproxy(ctx context.Context, formats strfmt.Registry) error {

	if m.Httpproxy != nil {

		if swag.IsZero(m.Httpproxy) { // not required
			return nil
		}

		if err := m.Httpproxy.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("httpproxy")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("httpproxy")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *OcspUpdateOptions) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *OcspUpdateOptions) UnmarshalBinary(b []byte) error {
	var res OcspUpdateOptions
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// OcspUpdateOptionsHttpproxy ocsp update options httpproxy
//
// swagger:model OcspUpdateOptionsHttpproxy
type OcspUpdateOptionsHttpproxy struct {
	// address
	// Example: 127.0.0.1
	// Pattern: ^[^\s]+$
	// +kubebuilder:validation:Pattern=`^[^\s]+$`
	Address string `json:"address,omitempty"`

	// port
	// Example: 80
	// Maximum: 65535
	// Minimum: 1
	// +kubebuilder:validation:Maximum=65535
	// +kubebuilder:validation:Minimum=1
	Port *int64 `json:"port,omitempty"`
}

// Validate validates this ocsp update options httpproxy
func (m *OcspUpdateOptionsHttpproxy) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAddress(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePort(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *OcspUpdateOptionsHttpproxy) validateAddress(formats strfmt.Registry) error {
	if swag.IsZero(m.Address) { // not required
		return nil
	}

	if err := validate.Pattern("httpproxy"+"."+"address", "body", m.Address, `^[^\s]+$`); err != nil {
		return err
	}

	return nil
}

func (m *OcspUpdateOptionsHttpproxy) validatePort(formats strfmt.Registry) error {
	if swag.IsZero(m.Port) { // not required
		return nil
	}

	if err := validate.MinimumInt("httpproxy"+"."+"port", "body", *m.Port, 1, false); err != nil {
		return err
	}

	if err := validate.MaximumInt("httpproxy"+"."+"port", "body", *m.Port, 65535, false); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this ocsp update options httpproxy based on context it is used
func (m *OcspUpdateOptionsHttpproxy) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *OcspUpdateOptionsHttpproxy) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *OcspUpdateOptionsHttpproxy) UnmarshalBinary(b []byte) error {
	var res OcspUpdateOptionsHttpproxy
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
