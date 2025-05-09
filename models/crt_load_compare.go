// Code generated with struct_equal_generator; DO NOT EDIT.

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

import "reflect"

// Equal checks if two structs of type CrtLoad are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b CrtLoad
//	equal := a.Equal(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b CrtLoad
//	equal := a.Equal(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s CrtLoad) Equal(t CrtLoad, opts ...Options) bool {
	opt := getOptions(opts...)

	if !equalComparableSlice(s.Domains, t.Domains, opt) {
		return false
	}

	if s.Acme != t.Acme {
		return false
	}

	if s.Alias != t.Alias {
		return false
	}

	if s.Certificate != t.Certificate {
		return false
	}

	if s.Issuer != t.Issuer {
		return false
	}

	if s.Key != t.Key {
		return false
	}

	if !CheckSameNilAndLenMap[string](s.Metadata, t.Metadata, opt) {
		return false
	}

	for k, v := range s.Metadata {
		if !reflect.DeepEqual(t.Metadata[k], v) {
			return false
		}
	}

	if s.Ocsp != t.Ocsp {
		return false
	}

	if s.OcspUpdate != t.OcspUpdate {
		return false
	}

	if s.Sctl != t.Sctl {
		return false
	}

	return true
}

// Diff checks if two structs of type CrtLoad are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b CrtLoad
//	diff := a.Diff(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b CrtLoad
//	diff := a.Diff(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s CrtLoad) Diff(t CrtLoad, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})
	if !equalComparableSlice(s.Domains, t.Domains, opt) {
		diff["Domains"] = []interface{}{s.Domains, t.Domains}
	}

	if s.Acme != t.Acme {
		diff["Acme"] = []interface{}{s.Acme, t.Acme}
	}

	if s.Alias != t.Alias {
		diff["Alias"] = []interface{}{s.Alias, t.Alias}
	}

	if s.Certificate != t.Certificate {
		diff["Certificate"] = []interface{}{s.Certificate, t.Certificate}
	}

	if s.Issuer != t.Issuer {
		diff["Issuer"] = []interface{}{s.Issuer, t.Issuer}
	}

	if s.Key != t.Key {
		diff["Key"] = []interface{}{s.Key, t.Key}
	}

	if !CheckSameNilAndLenMap[string](s.Metadata, t.Metadata, opt) {
		diff["Metadata"] = []interface{}{s.Metadata, t.Metadata}
	}

	for k, v := range s.Metadata {
		if !reflect.DeepEqual(t.Metadata[k], v) {
			diff["Metadata"] = []interface{}{s.Metadata, t.Metadata}
		}
	}

	if s.Ocsp != t.Ocsp {
		diff["Ocsp"] = []interface{}{s.Ocsp, t.Ocsp}
	}

	if s.OcspUpdate != t.OcspUpdate {
		diff["OcspUpdate"] = []interface{}{s.OcspUpdate, t.OcspUpdate}
	}

	if s.Sctl != t.Sctl {
		diff["Sctl"] = []interface{}{s.Sctl, t.Sctl}
	}

	return diff
}
