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

// Table table
//
// swagger:model table
type Table struct {

	// expire
	// Pattern: ^\d+(ms|s|m|h|d)?$
	Expire *string `json:"expire,omitempty"`

	// name
	Name string `json:"name,omitempty"`

	// no purge
	NoPurge bool `json:"no_purge,omitempty"`

	// size
	// Pattern: ^\d+(k|K|m|M|g|G)?$
	Size string `json:"size,omitempty"`

	// store
	Store string `json:"store,omitempty"`

	// type
	// Enum: [ip integer string binary]
	Type string `json:"type,omitempty"`

	// type len
	TypeLen *int64 `json:"type_len,omitempty"`
}

// Validate validates this table
func (m *Table) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateExpire(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSize(formats); err != nil {
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

func (m *Table) validateExpire(formats strfmt.Registry) error {
	if swag.IsZero(m.Expire) { // not required
		return nil
	}

	if err := validate.Pattern("expire", "body", *m.Expire, `^\d+(ms|s|m|h|d)?$`); err != nil {
		return err
	}

	return nil
}

func (m *Table) validateSize(formats strfmt.Registry) error {
	if swag.IsZero(m.Size) { // not required
		return nil
	}

	if err := validate.Pattern("size", "body", m.Size, `^\d+(k|K|m|M|g|G)?$`); err != nil {
		return err
	}

	return nil
}

var tableTypeTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["ip","integer","string","binary"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		tableTypeTypePropEnum = append(tableTypeTypePropEnum, v)
	}
}

const (

	// TableTypeIP captures enum value "ip"
	TableTypeIP string = "ip"

	// TableTypeInteger captures enum value "integer"
	TableTypeInteger string = "integer"

	// TableTypeString captures enum value "string"
	TableTypeString string = "string"

	// TableTypeBinary captures enum value "binary"
	TableTypeBinary string = "binary"
)

// prop value enum
func (m *Table) validateTypeEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, tableTypeTypePropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *Table) validateType(formats strfmt.Registry) error {
	if swag.IsZero(m.Type) { // not required
		return nil
	}

	// value enum
	if err := m.validateTypeEnum("type", "body", m.Type); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this table based on context it is used
func (m *Table) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *Table) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Table) UnmarshalBinary(b []byte) error {
	var res Table
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}