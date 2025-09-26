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

// Equal checks if two structs of type AcmeProvider are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b AcmeProvider
//	equal := a.Equal(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b AcmeProvider
//	equal := a.Equal(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s AcmeProvider) Equal(t AcmeProvider, opts ...Options) bool {
	opt := getOptions(opts...)

	if s.AccountKey != t.AccountKey {
		return false
	}

	if s.AcmeProvider != t.AcmeProvider {
		return false
	}

	if !equalComparableMap(s.AcmeVars, t.AcmeVars, opt) {
		return false
	}

	if !equalPointers(s.Bits, t.Bits) {
		return false
	}

	if s.Challenge != t.Challenge {
		return false
	}

	if s.Contact != t.Contact {
		return false
	}

	if s.Curves != t.Curves {
		return false
	}

	if s.Directory != t.Directory {
		return false
	}

	if s.Keytype != t.Keytype {
		return false
	}

	if s.Map != t.Map {
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

	if s.Name != t.Name {
		return false
	}

	return true
}

// Diff checks if two structs of type AcmeProvider are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b AcmeProvider
//	diff := a.Diff(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b AcmeProvider
//	diff := a.Diff(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s AcmeProvider) Diff(t AcmeProvider, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})
	if s.AccountKey != t.AccountKey {
		diff["AccountKey"] = []interface{}{s.AccountKey, t.AccountKey}
	}

	if s.AcmeProvider != t.AcmeProvider {
		diff["AcmeProvider"] = []interface{}{s.AcmeProvider, t.AcmeProvider}
	}

	if !equalComparableMap(s.AcmeVars, t.AcmeVars, opt) {
		diff["AcmeVars"] = []interface{}{s.AcmeVars, t.AcmeVars}
	}

	if !equalPointers(s.Bits, t.Bits) {
		diff["Bits"] = []interface{}{ValueOrNil(s.Bits), ValueOrNil(t.Bits)}
	}

	if s.Challenge != t.Challenge {
		diff["Challenge"] = []interface{}{s.Challenge, t.Challenge}
	}

	if s.Contact != t.Contact {
		diff["Contact"] = []interface{}{s.Contact, t.Contact}
	}

	if s.Curves != t.Curves {
		diff["Curves"] = []interface{}{s.Curves, t.Curves}
	}

	if s.Directory != t.Directory {
		diff["Directory"] = []interface{}{s.Directory, t.Directory}
	}

	if s.Keytype != t.Keytype {
		diff["Keytype"] = []interface{}{s.Keytype, t.Keytype}
	}

	if s.Map != t.Map {
		diff["Map"] = []interface{}{s.Map, t.Map}
	}

	if !CheckSameNilAndLenMap[string](s.Metadata, t.Metadata, opt) {
		diff["Metadata"] = []interface{}{s.Metadata, t.Metadata}
	}

	for k, v := range s.Metadata {
		if !reflect.DeepEqual(t.Metadata[k], v) {
			diff["Metadata"] = []interface{}{s.Metadata, t.Metadata}
		}
	}

	if s.Name != t.Name {
		diff["Name"] = []interface{}{s.Name, t.Name}
	}

	return diff
}
