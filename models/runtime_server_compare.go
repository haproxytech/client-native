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

// Equal checks if two structs of type RuntimeServer are equal
//
//	var a, b RuntimeServer
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s RuntimeServer) Equal(t RuntimeServer, opts ...Options) bool {
	if s.Address != t.Address {
		return false
	}

	if s.AdminState != t.AdminState {
		return false
	}

	if s.ID != t.ID {
		return false
	}

	if s.Name != t.Name {
		return false
	}

	if s.OperationalState != t.OperationalState {
		return false
	}

	if !equalPointers(s.Port, t.Port) {
		return false
	}

	return true
}

// Diff checks if two structs of type RuntimeServer are equal
//
//	var a, b RuntimeServer
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s RuntimeServer) Diff(t RuntimeServer, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if s.Address != t.Address {
		diff["Address"] = []interface{}{s.Address, t.Address}
	}

	if s.AdminState != t.AdminState {
		diff["AdminState"] = []interface{}{s.AdminState, t.AdminState}
	}

	if s.ID != t.ID {
		diff["ID"] = []interface{}{s.ID, t.ID}
	}

	if s.Name != t.Name {
		diff["Name"] = []interface{}{s.Name, t.Name}
	}

	if s.OperationalState != t.OperationalState {
		diff["OperationalState"] = []interface{}{s.OperationalState, t.OperationalState}
	}

	if !equalPointers(s.Port, t.Port) {
		diff["Port"] = []interface{}{s.Port, t.Port}
	}

	return diff
}
