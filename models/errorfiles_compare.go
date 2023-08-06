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

// Equal checks if two structs of type Errorfiles are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b Errorfiles
//	equal := a.Equal(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b Errorfiles
//	equal := a.Equal(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s Errorfiles) Equal(t Errorfiles, opts ...Options) bool {
	opt := getOptions(opts...)

	if !equalComparableSlice(s.Codes, t.Codes, opt) {
		return false
	}

	if s.Name != t.Name {
		return false
	}

	return true
}

// Diff checks if two structs of type Errorfiles are equal
//
// By default empty arrays, maps and slices are equal to nil:
//
//	var a, b Errorfiles
//	diff := a.Diff(b)
//
// For more advanced use case you can configure the options (default values are shown):
//
//	var a, b Errorfiles
//	equal := a.Diff(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s Errorfiles) Diff(t Errorfiles, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})
	if !equalComparableSlice(s.Codes, t.Codes, opt) {
		diff["Codes"] = []interface{}{s.Codes, t.Codes}
	}

	if s.Name != t.Name {
		diff["Name"] = []interface{}{s.Name, t.Name}
	}

	return diff
}
