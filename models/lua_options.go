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
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// LuaOptions lua options
//
// swagger:model lua_options
type LuaOptions struct {

	// load per thread
	LoadPerThread string `json:"load_per_thread,omitempty"`

	// loads
	Loads []*LuaLoad `json:"loads,omitempty"`

	// prepend path
	PrependPath []*LuaPrependPath `json:"prepend_path,omitempty"`
}

// Validate validates this lua options
func (m *LuaOptions) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateLoads(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePrependPath(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *LuaOptions) validateLoads(formats strfmt.Registry) error {
	if swag.IsZero(m.Loads) { // not required
		return nil
	}

	for i := 0; i < len(m.Loads); i++ {
		if swag.IsZero(m.Loads[i]) { // not required
			continue
		}

		if m.Loads[i] != nil {
			if err := m.Loads[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("loads" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("loads" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *LuaOptions) validatePrependPath(formats strfmt.Registry) error {
	if swag.IsZero(m.PrependPath) { // not required
		return nil
	}

	for i := 0; i < len(m.PrependPath); i++ {
		if swag.IsZero(m.PrependPath[i]) { // not required
			continue
		}

		if m.PrependPath[i] != nil {
			if err := m.PrependPath[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("prepend_path" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("prepend_path" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this lua options based on the context it is used
func (m *LuaOptions) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateLoads(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidatePrependPath(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *LuaOptions) contextValidateLoads(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Loads); i++ {

		if m.Loads[i] != nil {

			if swag.IsZero(m.Loads[i]) { // not required
				return nil
			}

			if err := m.Loads[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("loads" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("loads" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *LuaOptions) contextValidatePrependPath(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.PrependPath); i++ {

		if m.PrependPath[i] != nil {

			if swag.IsZero(m.PrependPath[i]) { // not required
				return nil
			}

			if err := m.PrependPath[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("prepend_path" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("prepend_path" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *LuaOptions) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *LuaOptions) UnmarshalBinary(b []byte) error {
	var res LuaOptions
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// LuaLoad lua load
//
// swagger:model LuaLoad
type LuaLoad struct {
	// file
	// Required: true
	// Pattern: ^[^\s]+$
	// +kubebuilder:validation:Pattern=`^[^\s]+$`
	File *string `json:"file"`
}

// Validate validates this lua load
func (m *LuaLoad) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateFile(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *LuaLoad) validateFile(formats strfmt.Registry) error {

	if err := validate.Required("file", "body", m.File); err != nil {
		return err
	}

	if err := validate.Pattern("file", "body", *m.File, `^[^\s]+$`); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this lua load based on context it is used
func (m *LuaLoad) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *LuaLoad) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *LuaLoad) UnmarshalBinary(b []byte) error {
	var res LuaLoad
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// LuaPrependPath lua prepend path
//
// swagger:model LuaPrependPath
type LuaPrependPath struct {
	// path
	// Required: true
	// Pattern: ^[^\s]+$
	// +kubebuilder:validation:Pattern=`^[^\s]+$`
	Path *string `json:"path"`

	// type
	// Enum: ["path","cpath"]
	// +kubebuilder:validation:Enum=path;cpath;
	Type string `json:"type,omitempty"`
}

// Validate validates this lua prepend path
func (m *LuaPrependPath) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validatePath(formats); err != nil {
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

func (m *LuaPrependPath) validatePath(formats strfmt.Registry) error {

	if err := validate.Required("path", "body", m.Path); err != nil {
		return err
	}

	if err := validate.Pattern("path", "body", *m.Path, `^[^\s]+$`); err != nil {
		return err
	}

	return nil
}

var luaPrependPathTypeTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["path","cpath"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		luaPrependPathTypeTypePropEnum = append(luaPrependPathTypeTypePropEnum, v)
	}
}

const (

	// LuaPrependPathTypePath captures enum value "path"
	LuaPrependPathTypePath string = "path"

	// LuaPrependPathTypeCpath captures enum value "cpath"
	LuaPrependPathTypeCpath string = "cpath"
)

// prop value enum
func (m *LuaPrependPath) validateTypeEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, luaPrependPathTypeTypePropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *LuaPrependPath) validateType(formats strfmt.Registry) error {
	if swag.IsZero(m.Type) { // not required
		return nil
	}

	// value enum
	if err := m.validateTypeEnum("type", "body", m.Type); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this lua prepend path based on context it is used
func (m *LuaPrependPath) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *LuaPrependPath) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *LuaPrependPath) UnmarshalBinary(b []byte) error {
	var res LuaPrependPath
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
