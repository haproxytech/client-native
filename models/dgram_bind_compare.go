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

// Equal checks if two structs of type DgramBind are equal
//
//	var a, b DgramBind
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s DgramBind) Equal(t DgramBind, opts ...Options) bool {
	if s.Address != t.Address {
		return false
	}

	if s.Interface != t.Interface {
		return false
	}

	if s.Name != t.Name {
		return false
	}

	if s.Namespace != t.Namespace {
		return false
	}

	if !equalPointers(s.Port, t.Port) {
		return false
	}

	if !equalPointers(s.PortRangeEnd, t.PortRangeEnd) {
		return false
	}

	if s.Transparent != t.Transparent {
		return false
	}

	return true
}

// Diff checks if two structs of type DgramBind are equal
//
//	var a, b DgramBind
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s DgramBind) Diff(t DgramBind, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if s.Address != t.Address {
		diff["Address"] = []interface{}{s.Address, t.Address}
	}

	if s.Interface != t.Interface {
		diff["Interface"] = []interface{}{s.Interface, t.Interface}
	}

	if s.Name != t.Name {
		diff["Name"] = []interface{}{s.Name, t.Name}
	}

	if s.Namespace != t.Namespace {
		diff["Namespace"] = []interface{}{s.Namespace, t.Namespace}
	}

	if !equalPointers(s.Port, t.Port) {
		diff["Port"] = []interface{}{ValueOrNil(s.Port), ValueOrNil(t.Port)}
	}

	if !equalPointers(s.PortRangeEnd, t.PortRangeEnd) {
		diff["PortRangeEnd"] = []interface{}{ValueOrNil(s.PortRangeEnd), ValueOrNil(t.PortRangeEnd)}
	}

	if s.Transparent != t.Transparent {
		diff["Transparent"] = []interface{}{s.Transparent, t.Transparent}
	}

	return diff
}
