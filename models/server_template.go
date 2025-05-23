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

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// ServerTemplate Server template
//
// Set a template to initialize servers with shared parameters.
// Example: {"fqdn":"google.com","num_or_range":"1-3","port":80,"prefix":"srv"}
//
// swagger:model server_template
type ServerTemplate struct {
	ServerParams `json:",inline"`

	// fqdn
	// Required: true
	Fqdn string `json:"fqdn"`

	// id
	ID *int64 `json:"id,omitempty"`

	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:Schemaless
	Metadata map[string]interface{} `json:"metadata,omitempty"`

	// num or range
	// Required: true
	NumOrRange string `json:"num_or_range"`

	// port
	// Maximum: 65535
	// Minimum: 1
	// +kubebuilder:validation:Maximum=65535
	// +kubebuilder:validation:Minimum=1
	Port *int64 `json:"port,omitempty"`

	// prefix
	// Required: true
	// Pattern: ^[^\s]+$
	// +kubebuilder:validation:Pattern=`^[^\s]+$`
	Prefix string `json:"prefix"`
}

// UnmarshalJSON unmarshals this object from a JSON structure
func (m *ServerTemplate) UnmarshalJSON(raw []byte) error {
	// AO0
	var aO0 ServerParams
	if err := swag.ReadJSON(raw, &aO0); err != nil {
		return err
	}
	m.ServerParams = aO0

	// AO1
	var dataAO1 struct {
		Fqdn string `json:"fqdn"`

		ID *int64 `json:"id,omitempty"`

		Metadata map[string]interface{} `json:"metadata,omitempty"`

		NumOrRange string `json:"num_or_range"`

		Port *int64 `json:"port,omitempty"`

		Prefix string `json:"prefix"`
	}
	if err := swag.ReadJSON(raw, &dataAO1); err != nil {
		return err
	}

	m.Fqdn = dataAO1.Fqdn

	m.ID = dataAO1.ID

	m.Metadata = dataAO1.Metadata

	m.NumOrRange = dataAO1.NumOrRange

	m.Port = dataAO1.Port

	m.Prefix = dataAO1.Prefix

	return nil
}

// MarshalJSON marshals this object to a JSON structure
func (m ServerTemplate) MarshalJSON() ([]byte, error) {
	_parts := make([][]byte, 0, 2)

	aO0, err := swag.WriteJSON(m.ServerParams)
	if err != nil {
		return nil, err
	}
	_parts = append(_parts, aO0)
	var dataAO1 struct {
		Fqdn string `json:"fqdn"`

		ID *int64 `json:"id,omitempty"`

		Metadata map[string]interface{} `json:"metadata,omitempty"`

		NumOrRange string `json:"num_or_range"`

		Port *int64 `json:"port,omitempty"`

		Prefix string `json:"prefix"`
	}

	dataAO1.Fqdn = m.Fqdn

	dataAO1.ID = m.ID

	dataAO1.Metadata = m.Metadata

	dataAO1.NumOrRange = m.NumOrRange

	dataAO1.Port = m.Port

	dataAO1.Prefix = m.Prefix

	jsonDataAO1, errAO1 := swag.WriteJSON(dataAO1)
	if errAO1 != nil {
		return nil, errAO1
	}
	_parts = append(_parts, jsonDataAO1)
	return swag.ConcatJSON(_parts...), nil
}

// Validate validates this server template
func (m *ServerTemplate) Validate(formats strfmt.Registry) error {
	var res []error

	// validation for a type composition with ServerParams
	if err := m.ServerParams.Validate(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateFqdn(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateNumOrRange(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePort(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePrefix(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ServerTemplate) validateFqdn(formats strfmt.Registry) error {

	if err := validate.RequiredString("fqdn", "body", m.Fqdn); err != nil {
		return err
	}

	return nil
}

func (m *ServerTemplate) validateNumOrRange(formats strfmt.Registry) error {

	if err := validate.RequiredString("num_or_range", "body", m.NumOrRange); err != nil {
		return err
	}

	return nil
}

func (m *ServerTemplate) validatePort(formats strfmt.Registry) error {

	if swag.IsZero(m.Port) { // not required
		return nil
	}

	if err := validate.MinimumInt("port", "body", *m.Port, 1, false); err != nil {
		return err
	}

	if err := validate.MaximumInt("port", "body", *m.Port, 65535, false); err != nil {
		return err
	}

	return nil
}

func (m *ServerTemplate) validatePrefix(formats strfmt.Registry) error {

	if err := validate.RequiredString("prefix", "body", m.Prefix); err != nil {
		return err
	}

	if err := validate.Pattern("prefix", "body", m.Prefix, `^[^\s]+$`); err != nil {
		return err
	}

	return nil
}

// ContextValidate validate this server template based on the context it is used
func (m *ServerTemplate) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	// validation for a type composition with ServerParams
	if err := m.ServerParams.ContextValidate(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// MarshalBinary interface implementation
func (m *ServerTemplate) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ServerTemplate) UnmarshalBinary(b []byte) error {
	var res ServerTemplate
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
