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

// Equal checks if two structs of type Cache are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b Cache
//	equal := a.Equal(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b Cache
//	equal := a.Equal(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s Cache) Equal(t Cache, opts ...Options) bool {
	opt := getOptions(opts...)

	if s.MaxAge != t.MaxAge {
		return false
	}

	if s.MaxObjectSize != t.MaxObjectSize {
		return false
	}

	if s.MaxSecondaryEntries != t.MaxSecondaryEntries {
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

	if !equalPointers(s.Name, t.Name) {
		return false
	}

	if !equalPointers(s.ProcessVary, t.ProcessVary) {
		return false
	}

	if s.TotalMaxSize != t.TotalMaxSize {
		return false
	}

	return true
}

// Diff checks if two structs of type Cache are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b Cache
//	diff := a.Diff(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b Cache
//	diff := a.Diff(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s Cache) Diff(t Cache, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})
	if s.MaxAge != t.MaxAge {
		diff["MaxAge"] = []interface{}{s.MaxAge, t.MaxAge}
	}

	if s.MaxObjectSize != t.MaxObjectSize {
		diff["MaxObjectSize"] = []interface{}{s.MaxObjectSize, t.MaxObjectSize}
	}

	if s.MaxSecondaryEntries != t.MaxSecondaryEntries {
		diff["MaxSecondaryEntries"] = []interface{}{s.MaxSecondaryEntries, t.MaxSecondaryEntries}
	}

	if !CheckSameNilAndLenMap[string](s.Metadata, t.Metadata, opt) {
		diff["Metadata"] = []interface{}{s.Metadata, t.Metadata}
	}

	for k, v := range s.Metadata {
		if !reflect.DeepEqual(t.Metadata[k], v) {
			diff["Metadata"] = []interface{}{s.Metadata, t.Metadata}
		}
	}

	if !equalPointers(s.Name, t.Name) {
		diff["Name"] = []interface{}{ValueOrNil(s.Name), ValueOrNil(t.Name)}
	}

	if !equalPointers(s.ProcessVary, t.ProcessVary) {
		diff["ProcessVary"] = []interface{}{ValueOrNil(s.ProcessVary), ValueOrNil(t.ProcessVary)}
	}

	if s.TotalMaxSize != t.TotalMaxSize {
		diff["TotalMaxSize"] = []interface{}{s.TotalMaxSize, t.TotalMaxSize}
	}

	return diff
}
