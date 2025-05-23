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

// HTTPAfterResponseRule HTTP after Response Rule
//
// HAProxy HTTP after response rule configuration (corresponds to http-after-response directives)
// Example: {"cond":"unless","cond_test":"{ src 192.168.0.0/16 }","hdr_format":"max-age=31536000","hdr_name":"Strict-Transport-Security","type":"set-header"}
//
// swagger:model http_after_response_rule
type HTTPAfterResponseRule struct {
	// acl file
	// Pattern: ^[^\s]+$
	// +kubebuilder:validation:Pattern=`^[^\s]+$`
	ACLFile string `json:"acl_file,omitempty"`

	// acl keyfmt
	// Pattern: ^[^\s]+$
	// +kubebuilder:validation:Pattern=`^[^\s]+$`
	ACLKeyfmt string `json:"acl_keyfmt,omitempty"`

	// capture id
	CaptureID *int64 `json:"capture_id,omitempty"`

	// capture len
	CaptureLen int64 `json:"capture_len,omitempty"`

	// capture sample
	// Pattern: ^(?:[A-Za-z]+\("([A-Za-z\s]+)"\)|[A-Za-z]+)
	// +kubebuilder:validation:Pattern=`^(?:[A-Za-z]+\("([A-Za-z\s]+)"\)|[A-Za-z]+)`
	CaptureSample string `json:"capture_sample,omitempty"`

	// cond
	// Enum: ["if","unless"]
	// +kubebuilder:validation:Enum=if;unless;
	Cond string `json:"cond,omitempty"`

	// cond test
	CondTest string `json:"cond_test,omitempty"`

	// hdr format
	HdrFormat string `json:"hdr_format,omitempty"`

	// hdr match
	HdrMatch string `json:"hdr_match,omitempty"`

	// hdr method
	HdrMethod string `json:"hdr_method,omitempty"`

	// hdr name
	HdrName string `json:"hdr_name,omitempty"`

	// log level
	// Enum: ["emerg","alert","crit","err","warning","notice","info","debug","silent"]
	// +kubebuilder:validation:Enum=emerg;alert;crit;err;warning;notice;info;debug;silent;
	LogLevel string `json:"log_level,omitempty"`

	// map file
	// Pattern: ^[^\s]+$
	// +kubebuilder:validation:Pattern=`^[^\s]+$`
	MapFile string `json:"map_file,omitempty"`

	// map keyfmt
	// Pattern: ^[^\s]+$
	// +kubebuilder:validation:Pattern=`^[^\s]+$`
	MapKeyfmt string `json:"map_keyfmt,omitempty"`

	// map valuefmt
	// Pattern: ^[^\s]+$
	// +kubebuilder:validation:Pattern=`^[^\s]+$`
	MapValuefmt string `json:"map_valuefmt,omitempty"`

	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:Schemaless
	Metadata map[string]interface{} `json:"metadata,omitempty"`

	// sc expr
	ScExpr string `json:"sc_expr,omitempty"`

	// sc id
	ScID int64 `json:"sc_id,omitempty"`

	// sc idx
	ScIdx int64 `json:"sc_idx,omitempty"`

	// sc int
	ScInt *int64 `json:"sc_int,omitempty"`

	// status
	// Maximum: 999
	// Minimum: 100
	// +kubebuilder:validation:Maximum=999
	// +kubebuilder:validation:Minimum=100
	Status int64 `json:"status,omitempty"`

	// status reason
	StatusReason string `json:"status_reason,omitempty"`

	// strict mode
	// Enum: ["on","off"]
	// +kubebuilder:validation:Enum=on;off;
	StrictMode string `json:"strict_mode,omitempty"`

	// type
	// Required: true
	// Enum: ["add-header","allow","capture","del-acl","del-header","del-map","replace-header","replace-value","sc-add-gpc","sc-inc-gpc","sc-inc-gpc0","sc-inc-gpc1","sc-set-gpt","sc-set-gpt0","set-header","set-log-level","set-map","set-status","set-var","set-var-fmt","strict-mode","unset-var","do-log"]
	// +kubebuilder:validation:Enum=add-header;allow;capture;del-acl;del-header;del-map;replace-header;replace-value;sc-add-gpc;sc-inc-gpc;sc-inc-gpc0;sc-inc-gpc1;sc-set-gpt;sc-set-gpt0;set-header;set-log-level;set-map;set-status;set-var;set-var-fmt;strict-mode;unset-var;do-log;
	Type string `json:"type"`

	// var expr
	VarExpr string `json:"var_expr,omitempty"`

	// var format
	VarFormat string `json:"var_format,omitempty"`

	// var name
	// Pattern: ^[^\s]+$
	// +kubebuilder:validation:Pattern=`^[^\s]+$`
	VarName string `json:"var_name,omitempty"`

	// var scope
	// Pattern: ^[^\s]+$
	// +kubebuilder:validation:Pattern=`^[^\s]+$`
	VarScope string `json:"var_scope,omitempty"`
}

// Validate validates this http after response rule
func (m *HTTPAfterResponseRule) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateACLFile(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateACLKeyfmt(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateCaptureSample(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateCond(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateLogLevel(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateMapFile(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateMapKeyfmt(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateMapValuefmt(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateStatus(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateStrictMode(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateType(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateVarName(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateVarScope(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *HTTPAfterResponseRule) validateACLFile(formats strfmt.Registry) error {
	if swag.IsZero(m.ACLFile) { // not required
		return nil
	}

	if err := validate.Pattern("acl_file", "body", m.ACLFile, `^[^\s]+$`); err != nil {
		return err
	}

	return nil
}

func (m *HTTPAfterResponseRule) validateACLKeyfmt(formats strfmt.Registry) error {
	if swag.IsZero(m.ACLKeyfmt) { // not required
		return nil
	}

	if err := validate.Pattern("acl_keyfmt", "body", m.ACLKeyfmt, `^[^\s]+$`); err != nil {
		return err
	}

	return nil
}

func (m *HTTPAfterResponseRule) validateCaptureSample(formats strfmt.Registry) error {
	if swag.IsZero(m.CaptureSample) { // not required
		return nil
	}

	if err := validate.Pattern("capture_sample", "body", m.CaptureSample, `^(?:[A-Za-z]+\("([A-Za-z\s]+)"\)|[A-Za-z]+)`); err != nil {
		return err
	}

	return nil
}

var httpAfterResponseRuleTypeCondPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["if","unless"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		httpAfterResponseRuleTypeCondPropEnum = append(httpAfterResponseRuleTypeCondPropEnum, v)
	}
}

const (

	// HTTPAfterResponseRuleCondIf captures enum value "if"
	HTTPAfterResponseRuleCondIf string = "if"

	// HTTPAfterResponseRuleCondUnless captures enum value "unless"
	HTTPAfterResponseRuleCondUnless string = "unless"
)

// prop value enum
func (m *HTTPAfterResponseRule) validateCondEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, httpAfterResponseRuleTypeCondPropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *HTTPAfterResponseRule) validateCond(formats strfmt.Registry) error {
	if swag.IsZero(m.Cond) { // not required
		return nil
	}

	// value enum
	if err := m.validateCondEnum("cond", "body", m.Cond); err != nil {
		return err
	}

	return nil
}

var httpAfterResponseRuleTypeLogLevelPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["emerg","alert","crit","err","warning","notice","info","debug","silent"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		httpAfterResponseRuleTypeLogLevelPropEnum = append(httpAfterResponseRuleTypeLogLevelPropEnum, v)
	}
}

const (

	// HTTPAfterResponseRuleLogLevelEmerg captures enum value "emerg"
	HTTPAfterResponseRuleLogLevelEmerg string = "emerg"

	// HTTPAfterResponseRuleLogLevelAlert captures enum value "alert"
	HTTPAfterResponseRuleLogLevelAlert string = "alert"

	// HTTPAfterResponseRuleLogLevelCrit captures enum value "crit"
	HTTPAfterResponseRuleLogLevelCrit string = "crit"

	// HTTPAfterResponseRuleLogLevelErr captures enum value "err"
	HTTPAfterResponseRuleLogLevelErr string = "err"

	// HTTPAfterResponseRuleLogLevelWarning captures enum value "warning"
	HTTPAfterResponseRuleLogLevelWarning string = "warning"

	// HTTPAfterResponseRuleLogLevelNotice captures enum value "notice"
	HTTPAfterResponseRuleLogLevelNotice string = "notice"

	// HTTPAfterResponseRuleLogLevelInfo captures enum value "info"
	HTTPAfterResponseRuleLogLevelInfo string = "info"

	// HTTPAfterResponseRuleLogLevelDebug captures enum value "debug"
	HTTPAfterResponseRuleLogLevelDebug string = "debug"

	// HTTPAfterResponseRuleLogLevelSilent captures enum value "silent"
	HTTPAfterResponseRuleLogLevelSilent string = "silent"
)

// prop value enum
func (m *HTTPAfterResponseRule) validateLogLevelEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, httpAfterResponseRuleTypeLogLevelPropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *HTTPAfterResponseRule) validateLogLevel(formats strfmt.Registry) error {
	if swag.IsZero(m.LogLevel) { // not required
		return nil
	}

	// value enum
	if err := m.validateLogLevelEnum("log_level", "body", m.LogLevel); err != nil {
		return err
	}

	return nil
}

func (m *HTTPAfterResponseRule) validateMapFile(formats strfmt.Registry) error {
	if swag.IsZero(m.MapFile) { // not required
		return nil
	}

	if err := validate.Pattern("map_file", "body", m.MapFile, `^[^\s]+$`); err != nil {
		return err
	}

	return nil
}

func (m *HTTPAfterResponseRule) validateMapKeyfmt(formats strfmt.Registry) error {
	if swag.IsZero(m.MapKeyfmt) { // not required
		return nil
	}

	if err := validate.Pattern("map_keyfmt", "body", m.MapKeyfmt, `^[^\s]+$`); err != nil {
		return err
	}

	return nil
}

func (m *HTTPAfterResponseRule) validateMapValuefmt(formats strfmt.Registry) error {
	if swag.IsZero(m.MapValuefmt) { // not required
		return nil
	}

	if err := validate.Pattern("map_valuefmt", "body", m.MapValuefmt, `^[^\s]+$`); err != nil {
		return err
	}

	return nil
}

func (m *HTTPAfterResponseRule) validateStatus(formats strfmt.Registry) error {
	if swag.IsZero(m.Status) { // not required
		return nil
	}

	if err := validate.MinimumInt("status", "body", m.Status, 100, false); err != nil {
		return err
	}

	if err := validate.MaximumInt("status", "body", m.Status, 999, false); err != nil {
		return err
	}

	return nil
}

var httpAfterResponseRuleTypeStrictModePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["on","off"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		httpAfterResponseRuleTypeStrictModePropEnum = append(httpAfterResponseRuleTypeStrictModePropEnum, v)
	}
}

const (

	// HTTPAfterResponseRuleStrictModeOn captures enum value "on"
	HTTPAfterResponseRuleStrictModeOn string = "on"

	// HTTPAfterResponseRuleStrictModeOff captures enum value "off"
	HTTPAfterResponseRuleStrictModeOff string = "off"
)

// prop value enum
func (m *HTTPAfterResponseRule) validateStrictModeEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, httpAfterResponseRuleTypeStrictModePropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *HTTPAfterResponseRule) validateStrictMode(formats strfmt.Registry) error {
	if swag.IsZero(m.StrictMode) { // not required
		return nil
	}

	// value enum
	if err := m.validateStrictModeEnum("strict_mode", "body", m.StrictMode); err != nil {
		return err
	}

	return nil
}

var httpAfterResponseRuleTypeTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["add-header","allow","capture","del-acl","del-header","del-map","replace-header","replace-value","sc-add-gpc","sc-inc-gpc","sc-inc-gpc0","sc-inc-gpc1","sc-set-gpt","sc-set-gpt0","set-header","set-log-level","set-map","set-status","set-var","set-var-fmt","strict-mode","unset-var","do-log"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		httpAfterResponseRuleTypeTypePropEnum = append(httpAfterResponseRuleTypeTypePropEnum, v)
	}
}

const (

	// HTTPAfterResponseRuleTypeAddDashHeader captures enum value "add-header"
	HTTPAfterResponseRuleTypeAddDashHeader string = "add-header"

	// HTTPAfterResponseRuleTypeAllow captures enum value "allow"
	HTTPAfterResponseRuleTypeAllow string = "allow"

	// HTTPAfterResponseRuleTypeCapture captures enum value "capture"
	HTTPAfterResponseRuleTypeCapture string = "capture"

	// HTTPAfterResponseRuleTypeDelDashACL captures enum value "del-acl"
	HTTPAfterResponseRuleTypeDelDashACL string = "del-acl"

	// HTTPAfterResponseRuleTypeDelDashHeader captures enum value "del-header"
	HTTPAfterResponseRuleTypeDelDashHeader string = "del-header"

	// HTTPAfterResponseRuleTypeDelDashMap captures enum value "del-map"
	HTTPAfterResponseRuleTypeDelDashMap string = "del-map"

	// HTTPAfterResponseRuleTypeReplaceDashHeader captures enum value "replace-header"
	HTTPAfterResponseRuleTypeReplaceDashHeader string = "replace-header"

	// HTTPAfterResponseRuleTypeReplaceDashValue captures enum value "replace-value"
	HTTPAfterResponseRuleTypeReplaceDashValue string = "replace-value"

	// HTTPAfterResponseRuleTypeScDashAddDashGpc captures enum value "sc-add-gpc"
	HTTPAfterResponseRuleTypeScDashAddDashGpc string = "sc-add-gpc"

	// HTTPAfterResponseRuleTypeScDashIncDashGpc captures enum value "sc-inc-gpc"
	HTTPAfterResponseRuleTypeScDashIncDashGpc string = "sc-inc-gpc"

	// HTTPAfterResponseRuleTypeScDashIncDashGpc0 captures enum value "sc-inc-gpc0"
	HTTPAfterResponseRuleTypeScDashIncDashGpc0 string = "sc-inc-gpc0"

	// HTTPAfterResponseRuleTypeScDashIncDashGpc1 captures enum value "sc-inc-gpc1"
	HTTPAfterResponseRuleTypeScDashIncDashGpc1 string = "sc-inc-gpc1"

	// HTTPAfterResponseRuleTypeScDashSetDashGpt captures enum value "sc-set-gpt"
	HTTPAfterResponseRuleTypeScDashSetDashGpt string = "sc-set-gpt"

	// HTTPAfterResponseRuleTypeScDashSetDashGpt0 captures enum value "sc-set-gpt0"
	HTTPAfterResponseRuleTypeScDashSetDashGpt0 string = "sc-set-gpt0"

	// HTTPAfterResponseRuleTypeSetDashHeader captures enum value "set-header"
	HTTPAfterResponseRuleTypeSetDashHeader string = "set-header"

	// HTTPAfterResponseRuleTypeSetDashLogDashLevel captures enum value "set-log-level"
	HTTPAfterResponseRuleTypeSetDashLogDashLevel string = "set-log-level"

	// HTTPAfterResponseRuleTypeSetDashMap captures enum value "set-map"
	HTTPAfterResponseRuleTypeSetDashMap string = "set-map"

	// HTTPAfterResponseRuleTypeSetDashStatus captures enum value "set-status"
	HTTPAfterResponseRuleTypeSetDashStatus string = "set-status"

	// HTTPAfterResponseRuleTypeSetDashVar captures enum value "set-var"
	HTTPAfterResponseRuleTypeSetDashVar string = "set-var"

	// HTTPAfterResponseRuleTypeSetDashVarDashFmt captures enum value "set-var-fmt"
	HTTPAfterResponseRuleTypeSetDashVarDashFmt string = "set-var-fmt"

	// HTTPAfterResponseRuleTypeStrictDashMode captures enum value "strict-mode"
	HTTPAfterResponseRuleTypeStrictDashMode string = "strict-mode"

	// HTTPAfterResponseRuleTypeUnsetDashVar captures enum value "unset-var"
	HTTPAfterResponseRuleTypeUnsetDashVar string = "unset-var"

	// HTTPAfterResponseRuleTypeDoDashLog captures enum value "do-log"
	HTTPAfterResponseRuleTypeDoDashLog string = "do-log"
)

// prop value enum
func (m *HTTPAfterResponseRule) validateTypeEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, httpAfterResponseRuleTypeTypePropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *HTTPAfterResponseRule) validateType(formats strfmt.Registry) error {

	if err := validate.RequiredString("type", "body", m.Type); err != nil {
		return err
	}

	// value enum
	if err := m.validateTypeEnum("type", "body", m.Type); err != nil {
		return err
	}

	return nil
}

func (m *HTTPAfterResponseRule) validateVarName(formats strfmt.Registry) error {
	if swag.IsZero(m.VarName) { // not required
		return nil
	}

	if err := validate.Pattern("var_name", "body", m.VarName, `^[^\s]+$`); err != nil {
		return err
	}

	return nil
}

func (m *HTTPAfterResponseRule) validateVarScope(formats strfmt.Registry) error {
	if swag.IsZero(m.VarScope) { // not required
		return nil
	}

	if err := validate.Pattern("var_scope", "body", m.VarScope, `^[^\s]+$`); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this http after response rule based on context it is used
func (m *HTTPAfterResponseRule) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *HTTPAfterResponseRule) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *HTTPAfterResponseRule) UnmarshalBinary(b []byte) error {
	var res HTTPAfterResponseRule
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
