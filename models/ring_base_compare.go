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

// Equal checks if two structs of type RingBase are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b RingBase
//	equal := a.Equal(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b RingBase
//	equal := a.Equal(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s RingBase) Equal(t RingBase, opts ...Options) bool {
	opt := getOptions(opts...)

	if s.Description != t.Description {
		return false
	}

	if s.Format != t.Format {
		return false
	}

	if !equalPointers(s.Maxlen, t.Maxlen) {
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

// Diff checks if two structs of type RingBase are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b RingBase
//	diff := a.Diff(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b RingBase
//	diff := a.Diff(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s RingBase) Diff(t RingBase, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

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
