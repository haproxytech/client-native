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

// Userlist Userlist with all it's children resources
//
// swagger:model Userlist
type Userlist struct {
	UserlistBase `json:",inline"`

	// groups
	Groups map[string]Group `json:"groups,omitempty"`

	// users
	Users map[string]User `json:"users,omitempty"`
}

// UnmarshalJSON unmarshals this object from a JSON structure
func (m *Userlist) UnmarshalJSON(raw []byte) error {
	// AO0
	var aO0 UserlistBase
	if err := swag.ReadJSON(raw, &aO0); err != nil {
		return err
	}
	m.UserlistBase = aO0

	// AO1
	var dataAO1 struct {
		Groups map[string]Group `json:"groups,omitempty"`

		Users map[string]User `json:"users,omitempty"`
	}
	if err := swag.ReadJSON(raw, &dataAO1); err != nil {
		return err
	}

	m.Groups = dataAO1.Groups

	m.Users = dataAO1.Users

	return nil
}

// MarshalJSON marshals this object to a JSON structure
func (m Userlist) MarshalJSON() ([]byte, error) {
	_parts := make([][]byte, 0, 2)

	aO0, err := swag.WriteJSON(m.UserlistBase)
	if err != nil {
		return nil, err
	}
	_parts = append(_parts, aO0)
	var dataAO1 struct {
		Groups map[string]Group `json:"groups,omitempty"`

		Users map[string]User `json:"users,omitempty"`
	}

	dataAO1.Groups = m.Groups

	dataAO1.Users = m.Users

	jsonDataAO1, errAO1 := swag.WriteJSON(dataAO1)
	if errAO1 != nil {
		return nil, errAO1
	}
	_parts = append(_parts, jsonDataAO1)
	return swag.ConcatJSON(_parts...), nil
}

// Validate validates this userlist
func (m *Userlist) Validate(formats strfmt.Registry) error {
	var res []error

	// validation for a type composition with UserlistBase
	if err := m.UserlistBase.Validate(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateGroups(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateUsers(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Userlist) validateGroups(formats strfmt.Registry) error {

	if swag.IsZero(m.Groups) { // not required
		return nil
	}

	for k := range m.Groups {

		if err := validate.Required("groups"+"."+k, "body", m.Groups[k]); err != nil {
			return err
		}
		if val, ok := m.Groups[k]; ok {
			if err := val.Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("groups" + "." + k)
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("groups" + "." + k)
				}
				return err
			}
		}

	}

	return nil
}

func (m *Userlist) validateUsers(formats strfmt.Registry) error {

	if swag.IsZero(m.Users) { // not required
		return nil
	}

	for k := range m.Users {

		if err := validate.Required("users"+"."+k, "body", m.Users[k]); err != nil {
			return err
		}
		if val, ok := m.Users[k]; ok {
			if err := val.Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("users" + "." + k)
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("users" + "." + k)
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this userlist based on the context it is used
func (m *Userlist) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	// validation for a type composition with UserlistBase
	if err := m.UserlistBase.ContextValidate(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateGroups(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateUsers(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Userlist) contextValidateGroups(ctx context.Context, formats strfmt.Registry) error {

	for k := range m.Groups {

		if val, ok := m.Groups[k]; ok {
			if err := val.ContextValidate(ctx, formats); err != nil {
				return err
			}
		}

	}

	return nil
}

func (m *Userlist) contextValidateUsers(ctx context.Context, formats strfmt.Registry) error {

	for k := range m.Users {

		if val, ok := m.Users[k]; ok {
			if err := val.ContextValidate(ctx, formats); err != nil {
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *Userlist) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Userlist) UnmarshalBinary(b []byte) error {
	var res Userlist
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
