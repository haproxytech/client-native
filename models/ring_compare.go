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

// Equal checks if two structs of type Ring are equal
//
//	var a, b Ring
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s Ring) Equal(t Ring, opts ...Options) bool {
	if s.Description != t.Description {
		return false
	}

	if s.Format != t.Format {
		return false
	}

	if !equalPointers(s.Maxlen, t.Maxlen) {
		return false
	}

	if s.Name != t.Name {
		return false
	}

	if !equalPointers(s.Size, t.Size) {
		return false
	}

	if !equalPointers(s.TimeoutConnect, t.TimeoutConnect) {
		return false
	}

	if !equalPointers(s.TimeoutServer, t.TimeoutServer) {
		return false
	}

	return true
}

// Diff checks if two structs of type Ring are equal
//
//	var a, b Ring
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s Ring) Diff(t Ring, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if s.Description != t.Description {
		diff["Description"] = []interface{}{s.Description, t.Description}
	}

	if s.Format != t.Format {
		diff["Format"] = []interface{}{s.Format, t.Format}
	}

	if !equalPointers(s.Maxlen, t.Maxlen) {
		diff["Maxlen"] = []interface{}{ValueOrNil(s.Maxlen), ValueOrNil(t.Maxlen)}
	}

	if s.Name != t.Name {
		diff["Name"] = []interface{}{s.Name, t.Name}
	}

	if !equalPointers(s.Size, t.Size) {
		diff["Size"] = []interface{}{ValueOrNil(s.Size), ValueOrNil(t.Size)}
	}

	if !equalPointers(s.TimeoutConnect, t.TimeoutConnect) {
		diff["TimeoutConnect"] = []interface{}{ValueOrNil(s.TimeoutConnect), ValueOrNil(t.TimeoutConnect)}
	}

	if !equalPointers(s.TimeoutServer, t.TimeoutServer) {
		diff["TimeoutServer"] = []interface{}{ValueOrNil(s.TimeoutServer), ValueOrNil(t.TimeoutServer)}
	}

	return diff
}
