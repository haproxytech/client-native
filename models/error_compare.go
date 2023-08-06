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

// Equal checks if two structs of type Error are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b Error
//	equal := a.Equal(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b Error
//	equal := a.Equal(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s Error) Equal(t Error, opts ...Options) bool {
	opt := getOptions(opts...)

	if !equalPointers(s.Code, t.Code) {
		return false
	}

	if !equalPointers(s.Message, t.Message) {
		return false
	}

	if !equalComparableMap(s.Error, t.Error, opt) {
		return false
	}

	return true
}

// Diff checks if two structs of type Error are equal
//
// By default empty arrays, maps and slices are equal to nil:
//
//	var a, b Error
//	diff := a.Diff(b)
//
// For more advanced use case you can configure the options (default values are shown):
//
//	var a, b Error
//	equal := a.Diff(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s Error) Diff(t Error, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})
	if !equalPointers(s.Code, t.Code) {
		diff["Code"] = []interface{}{s.Code, t.Code}
	}

	if !equalPointers(s.Message, t.Message) {
		diff["Message"] = []interface{}{s.Message, t.Message}
	}

	if !equalComparableMap(s.Error, t.Error, opt) {
		diff["Error"] = []interface{}{s.Error, t.Error}
	}

	return diff
}
