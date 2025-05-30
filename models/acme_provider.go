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

// AcmeProvider ACME Provider
//
// # Define an ACME provider to generate certificates automatically
//
// swagger:model acme_provider
type AcmeProvider struct {

	// Path where the the ACME account key is stored
	AccountKey string `json:"account_key,omitempty"`

	// Number of bits to generate an RSA certificate
	// Minimum: 1024
	// +kubebuilder:validation:Minimum=1024
	Bits *int64 `json:"bits,omitempty"`

	// ACME challenge type. Only HTTP-01 and DNS-01 are supported.
	// Enum: ["HTTP-01","DNS-01"]
	// +kubebuilder:validation:Enum=HTTP-01;DNS-01;
	Challenge string `json:"challenge,omitempty"`

	// Contact email for the ACME account
	Contact string `json:"contact,omitempty"`

	// Curves used with the ECDSA key type
	Curves string `json:"curves,omitempty"`

	// URL to the ACME provider's directory. For example:
	// https://acme-staging-v02.api.letsencrypt.org/directory
	//
	// Required: true
	// Pattern: ^https://[^\s]+$
	// +kubebuilder:validation:Pattern=`^https://[^\s]+$`
	Directory string `json:"directory"`

	// Type of key to generate
	// Enum: ["RSA","ECDSA"]
	// +kubebuilder:validation:Enum=RSA;ECDSA;
	Keytype string `json:"keytype,omitempty"`

	// The map which will be used to store the ACME token (key) and thumbprint
	Map string `json:"map,omitempty"`

	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:Schemaless
	Metadata map[string]interface{} `json:"metadata,omitempty"`

	// ACME provider's name
	// Required: true
	Name string `json:"name"`
}

// Validate validates this acme provider
func (m *AcmeProvider) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateBits(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateChallenge(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateDirectory(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateKeytype(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateName(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *AcmeProvider) validateBits(formats strfmt.Registry) error {
	if swag.IsZero(m.Bits) { // not required
		return nil
	}

	if err := validate.MinimumInt("bits", "body", *m.Bits, 1024, false); err != nil {
		return err
	}

	return nil
}

var acmeProviderTypeChallengePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["HTTP-01","DNS-01"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		acmeProviderTypeChallengePropEnum = append(acmeProviderTypeChallengePropEnum, v)
	}
}

const (

	// AcmeProviderChallengeHTTPDash01 captures enum value "HTTP-01"
	AcmeProviderChallengeHTTPDash01 string = "HTTP-01"

	// AcmeProviderChallengeDNSDash01 captures enum value "DNS-01"
	AcmeProviderChallengeDNSDash01 string = "DNS-01"
)

// prop value enum
func (m *AcmeProvider) validateChallengeEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, acmeProviderTypeChallengePropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *AcmeProvider) validateChallenge(formats strfmt.Registry) error {
	if swag.IsZero(m.Challenge) { // not required
		return nil
	}

	// value enum
	if err := m.validateChallengeEnum("challenge", "body", m.Challenge); err != nil {
		return err
	}

	return nil
}

func (m *AcmeProvider) validateDirectory(formats strfmt.Registry) error {

	if err := validate.RequiredString("directory", "body", m.Directory); err != nil {
		return err
	}

	if err := validate.Pattern("directory", "body", m.Directory, `^https://[^\s]+$`); err != nil {
		return err
	}

	return nil
}

var acmeProviderTypeKeytypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["RSA","ECDSA"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		acmeProviderTypeKeytypePropEnum = append(acmeProviderTypeKeytypePropEnum, v)
	}
}

const (

	// AcmeProviderKeytypeRSA captures enum value "RSA"
	AcmeProviderKeytypeRSA string = "RSA"

	// AcmeProviderKeytypeECDSA captures enum value "ECDSA"
	AcmeProviderKeytypeECDSA string = "ECDSA"
)

// prop value enum
func (m *AcmeProvider) validateKeytypeEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, acmeProviderTypeKeytypePropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *AcmeProvider) validateKeytype(formats strfmt.Registry) error {
	if swag.IsZero(m.Keytype) { // not required
		return nil
	}

	// value enum
	if err := m.validateKeytypeEnum("keytype", "body", m.Keytype); err != nil {
		return err
	}

	return nil
}

func (m *AcmeProvider) validateName(formats strfmt.Registry) error {

	if err := validate.RequiredString("name", "body", m.Name); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this acme provider based on context it is used
func (m *AcmeProvider) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *AcmeProvider) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *AcmeProvider) UnmarshalBinary(b []byte) error {
	var res AcmeProvider
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
