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

// Equal checks if two structs of type Cache are equal
//
//	var a, b Cache
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s Cache) Equal(t Cache, opts ...Options) bool {
	if s.MaxAge != t.MaxAge {
		return false
	}

	if s.MaxObjectSize != t.MaxObjectSize {
		return false
	}

	if s.MaxSecondaryEntries != t.MaxSecondaryEntries {
		return false
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
//	var a, b Cache
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s Cache) Diff(t Cache, opts ...Options) map[string][]interface{} {
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

	if !equalPointers(s.Name, t.Name) {
		diff["Name"] = []interface{}{s.Name, t.Name}
	}

	if !equalPointers(s.ProcessVary, t.ProcessVary) {
		diff["ProcessVary"] = []interface{}{s.ProcessVary, t.ProcessVary}
	}

	if s.TotalMaxSize != t.TotalMaxSize {
		diff["TotalMaxSize"] = []interface{}{s.TotalMaxSize, t.TotalMaxSize}
	}

	return diff
}
