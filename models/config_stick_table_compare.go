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

// Equal checks if two structs of type ConfigStickTable are equal
//
//	var a, b ConfigStickTable
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s ConfigStickTable) Equal(t ConfigStickTable, opts ...Options) bool {
	if !equalPointers(s.Expire, t.Expire) {
		return false
	}

	if !equalPointers(s.Keylen, t.Keylen) {
		return false
	}

	if s.Nopurge != t.Nopurge {
		return false
	}

	if s.Peers != t.Peers {
		return false
	}

	if !equalPointers(s.Size, t.Size) {
		return false
	}

	if s.Store != t.Store {
		return false
	}

	if s.Type != t.Type {
		return false
	}

	return true
}

// Diff checks if two structs of type ConfigStickTable are equal
//
//	var a, b ConfigStickTable
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s ConfigStickTable) Diff(t ConfigStickTable, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if !equalPointers(s.Expire, t.Expire) {
		diff["Expire"] = []interface{}{s.Expire, t.Expire}
	}

	if !equalPointers(s.Keylen, t.Keylen) {
		diff["Keylen"] = []interface{}{s.Keylen, t.Keylen}
	}

	if s.Nopurge != t.Nopurge {
		diff["Nopurge"] = []interface{}{s.Nopurge, t.Nopurge}
	}

	if s.Peers != t.Peers {
		diff["Peers"] = []interface{}{s.Peers, t.Peers}
	}

	if !equalPointers(s.Size, t.Size) {
		diff["Size"] = []interface{}{s.Size, t.Size}
	}

	if s.Store != t.Store {
		diff["Store"] = []interface{}{s.Store, t.Store}
	}

	if s.Type != t.Type {
		diff["Type"] = []interface{}{s.Type, t.Type}
	}

	return diff
}
