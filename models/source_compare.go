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

// Equal checks if two structs of type Source are equal
//
//	var a, b Source
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s Source) Equal(t Source, opts ...Options) bool {
	if !equalPointers(s.Address, t.Address) {
		return false
	}

	if s.AddressSecond != t.AddressSecond {
		return false
	}

	if s.Hdr != t.Hdr {
		return false
	}

	if s.Interface != t.Interface {
		return false
	}

	if s.Occ != t.Occ {
		return false
	}

	if s.Port != t.Port {
		return false
	}

	if s.PortSecond != t.PortSecond {
		return false
	}

	if s.Usesrc != t.Usesrc {
		return false
	}

	return true
}

// Diff checks if two structs of type Source are equal
//
//	var a, b Source
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s Source) Diff(t Source, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if !equalPointers(s.Address, t.Address) {
		diff["Address"] = []interface{}{ValueOrNil(s.Address), ValueOrNil(t.Address)}
	}

	if s.AddressSecond != t.AddressSecond {
		diff["AddressSecond"] = []interface{}{s.AddressSecond, t.AddressSecond}
	}

	if s.Hdr != t.Hdr {
		diff["Hdr"] = []interface{}{s.Hdr, t.Hdr}
	}

	if s.Interface != t.Interface {
		diff["Interface"] = []interface{}{s.Interface, t.Interface}
	}

	if s.Occ != t.Occ {
		diff["Occ"] = []interface{}{s.Occ, t.Occ}
	}

	if s.Port != t.Port {
		diff["Port"] = []interface{}{s.Port, t.Port}
	}

	if s.PortSecond != t.PortSecond {
		diff["PortSecond"] = []interface{}{s.PortSecond, t.PortSecond}
	}

	if s.Usesrc != t.Usesrc {
		diff["Usesrc"] = []interface{}{s.Usesrc, t.Usesrc}
	}

	return diff
}
