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

// Equal checks if two structs of type Table are equal
//
//	var a, b Table
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s Table) Equal(t Table, opts ...Options) bool {
	if !equalPointers(s.Expire, t.Expire) {
		return false
	}

	if s.Name != t.Name {
		return false
	}

	if s.NoPurge != t.NoPurge {
		return false
	}

	if s.Size != t.Size {
		return false
	}

	if s.Store != t.Store {
		return false
	}

	if s.Type != t.Type {
		return false
	}

	if !equalPointers(s.TypeLen, t.TypeLen) {
		return false
	}

	if !equalPointers(s.WriteTo, t.WriteTo) {
		return false
	}

	return true
}

// Diff checks if two structs of type Table are equal
//
//	var a, b Table
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s Table) Diff(t Table, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if !equalPointers(s.Expire, t.Expire) {
		diff["Expire"] = []interface{}{ValueOrNil(s.Expire), ValueOrNil(t.Expire)}
	}

	if s.Name != t.Name {
		diff["Name"] = []interface{}{s.Name, t.Name}
	}

	if s.NoPurge != t.NoPurge {
		diff["NoPurge"] = []interface{}{s.NoPurge, t.NoPurge}
	}

	if s.Size != t.Size {
		diff["Size"] = []interface{}{s.Size, t.Size}
	}

	if s.Store != t.Store {
		diff["Store"] = []interface{}{s.Store, t.Store}
	}

	if s.Type != t.Type {
		diff["Type"] = []interface{}{s.Type, t.Type}
	}

	if !equalPointers(s.TypeLen, t.TypeLen) {
		diff["TypeLen"] = []interface{}{ValueOrNil(s.TypeLen), ValueOrNil(t.TypeLen)}
	}

	if !equalPointers(s.WriteTo, t.WriteTo) {
		diff["WriteTo"] = []interface{}{ValueOrNil(s.WriteTo), ValueOrNil(t.WriteTo)}
	}

	return diff
}
