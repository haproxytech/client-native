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

// HTTPErrorRule HTTP Error Rule
//
// HAProxy HTTP error rule configuration (corresponds to http-error directives)
// Example: {"index":0,"status":425,"type":"status"}
//
// swagger:model http_error_rule
type HTTPErrorRule struct {

	// return headers
	ReturnHeaders []*ReturnHeader `json:"return_hdrs,omitempty"`

	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:Schemaless
	Metadata map[string]interface{} `json:"metadata,omitempty"`

	// return content
	ReturnContent string `json:"return_content,omitempty"`

	// return content format
	// Enum: ["default-errorfiles","errorfile","errorfiles","file","lf-file","string","lf-string"]
	// +kubebuilder:validation:Enum=default-errorfiles;errorfile;errorfiles;file;lf-file;string;lf-string;
	ReturnContentFormat string `json:"return_content_format,omitempty"`

	// return content type
	ReturnContentType *string `json:"return_content_type,omitempty"`

	// status
	// Required: true
	// Enum: [200,400,401,403,404,405,407,408,410,413,425,429,500,501,502,503,504]
	// +kubebuilder:validation:Enum=200;400;401;403;404;405;407;408;410;413;425;429;500;501;502;503;504;
	Status int64 `json:"status"`

	// type
	// Required: true
	// Enum: ["status"]
	// +kubebuilder:validation:Enum=status;
	Type string `json:"type"`
}

// Validate validates this http error rule
func (m *HTTPErrorRule) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateReturnHeaders(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateReturnContentFormat(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateStatus(formats); err != nil {
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

func (m *HTTPErrorRule) validateReturnHeaders(formats strfmt.Registry) error {
	if swag.IsZero(m.ReturnHeaders) { // not required
		return nil
	}

	for i := 0; i < len(m.ReturnHeaders); i++ {
		if swag.IsZero(m.ReturnHeaders[i]) { // not required
			continue
		}

		if m.ReturnHeaders[i] != nil {
			if err := m.ReturnHeaders[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("return_hdrs" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("return_hdrs" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

var httpErrorRuleTypeReturnContentFormatPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["default-errorfiles","errorfile","errorfiles","file","lf-file","string","lf-string"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		httpErrorRuleTypeReturnContentFormatPropEnum = append(httpErrorRuleTypeReturnContentFormatPropEnum, v)
	}
}

const (

	// HTTPErrorRuleReturnContentFormatDefaultDashErrorfiles captures enum value "default-errorfiles"
	HTTPErrorRuleReturnContentFormatDefaultDashErrorfiles string = "default-errorfiles"

	// HTTPErrorRuleReturnContentFormatErrorfile captures enum value "errorfile"
	HTTPErrorRuleReturnContentFormatErrorfile string = "errorfile"

	// HTTPErrorRuleReturnContentFormatErrorfiles captures enum value "errorfiles"
	HTTPErrorRuleReturnContentFormatErrorfiles string = "errorfiles"

	// HTTPErrorRuleReturnContentFormatFile captures enum value "file"
	HTTPErrorRuleReturnContentFormatFile string = "file"

	// HTTPErrorRuleReturnContentFormatLfDashFile captures enum value "lf-file"
	HTTPErrorRuleReturnContentFormatLfDashFile string = "lf-file"

	// HTTPErrorRuleReturnContentFormatString captures enum value "string"
	HTTPErrorRuleReturnContentFormatString string = "string"

	// HTTPErrorRuleReturnContentFormatLfDashString captures enum value "lf-string"
	HTTPErrorRuleReturnContentFormatLfDashString string = "lf-string"
)

// prop value enum
func (m *HTTPErrorRule) validateReturnContentFormatEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, httpErrorRuleTypeReturnContentFormatPropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *HTTPErrorRule) validateReturnContentFormat(formats strfmt.Registry) error {
	if swag.IsZero(m.ReturnContentFormat) { // not required
		return nil
	}

	// value enum
	if err := m.validateReturnContentFormatEnum("return_content_format", "body", m.ReturnContentFormat); err != nil {
		return err
	}

	return nil
}

var httpErrorRuleTypeStatusPropEnum []interface{}

func init() {
	var res []int64
	if err := json.Unmarshal([]byte(`[200,400,401,403,404,405,407,408,410,413,425,429,500,501,502,503,504]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		httpErrorRuleTypeStatusPropEnum = append(httpErrorRuleTypeStatusPropEnum, v)
	}
}

// prop value enum
func (m *HTTPErrorRule) validateStatusEnum(path, location string, value int64) error {
	if err := validate.EnumCase(path, location, value, httpErrorRuleTypeStatusPropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *HTTPErrorRule) validateStatus(formats strfmt.Registry) error {

	if err := validate.Required("status", "body", int64(m.Status)); err != nil {
		return err
	}

	// value enum
	if err := m.validateStatusEnum("status", "body", m.Status); err != nil {
		return err
	}

	return nil
}

var httpErrorRuleTypeTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["status"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		httpErrorRuleTypeTypePropEnum = append(httpErrorRuleTypeTypePropEnum, v)
	}
}

const (

	// HTTPErrorRuleTypeStatus captures enum value "status"
	HTTPErrorRuleTypeStatus string = "status"
)

// prop value enum
func (m *HTTPErrorRule) validateTypeEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, httpErrorRuleTypeTypePropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *HTTPErrorRule) validateType(formats strfmt.Registry) error {

	if err := validate.RequiredString("type", "body", m.Type); err != nil {
		return err
	}

	// value enum
	if err := m.validateTypeEnum("type", "body", m.Type); err != nil {
		return err
	}

	return nil
}

// ContextValidate validate this http error rule based on the context it is used
func (m *HTTPErrorRule) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateReturnHeaders(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *HTTPErrorRule) contextValidateReturnHeaders(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.ReturnHeaders); i++ {

		if m.ReturnHeaders[i] != nil {

			if swag.IsZero(m.ReturnHeaders[i]) { // not required
				return nil
			}

			if err := m.ReturnHeaders[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("return_hdrs" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("return_hdrs" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *HTTPErrorRule) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *HTTPErrorRule) UnmarshalBinary(b []byte) error {
	var res HTTPErrorRule
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
