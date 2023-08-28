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

// Equal checks if two structs of type User are equal
//
//	var a, b User
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s User) Equal(t User, opts ...Options) bool {
	if s.Groups != t.Groups {
		return false
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
//	var a, b User
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s User) Diff(t User, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if s.Groups != t.Groups {
		diff["Groups"] = []interface{}{s.Groups, t.Groups}
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
