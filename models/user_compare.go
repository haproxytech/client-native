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

// Equal checks if two structs of type User are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b User
//	equal := a.Equal(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b User
//	equal := a.Equal(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s User) Equal(t User, opts ...Options) bool {
	opt := getOptions(opts...)

	if s.Groups != t.Groups {
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

	if s.Password != t.Password {
		return false
	}

	if !equalPointers(s.SecurePassword, t.SecurePassword) {
		return false
	}

	if s.Username != t.Username {
		return false
	}

	return true
}

// Diff checks if two structs of type User are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b User
//	diff := a.Diff(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b User
//	diff := a.Diff(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s User) Diff(t User, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})
	if s.Groups != t.Groups {
		diff["Groups"] = []interface{}{s.Groups, t.Groups}
	}

	if !CheckSameNilAndLenMap[string](s.Metadata, t.Metadata, opt) {
		diff["Metadata"] = []interface{}{s.Metadata, t.Metadata}
	}

	for k, v := range s.Metadata {
		if !reflect.DeepEqual(t.Metadata[k], v) {
			diff["Metadata"] = []interface{}{s.Metadata, t.Metadata}
		}
	}

	if s.Password != t.Password {
		diff["Password"] = []interface{}{s.Password, t.Password}
	}

	if !equalPointers(s.SecurePassword, t.SecurePassword) {
		diff["SecurePassword"] = []interface{}{ValueOrNil(s.SecurePassword), ValueOrNil(t.SecurePassword)}
	}

	if s.Username != t.Username {
		diff["Username"] = []interface{}{s.Username, t.Username}
	}

	return diff
}
