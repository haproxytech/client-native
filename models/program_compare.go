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

// Equal checks if two structs of type Program are equal
//
//	var a, b Program
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s Program) Equal(t Program, opts ...Options) bool {
	if !equalPointers(s.Command, t.Command) {
		return false
	}

	if s.Group != t.Group {
		return false
	}

	if s.Name != t.Name {
		return false
	}

	if s.StartOnReload != t.StartOnReload {
		return false
	}

	if s.User != t.User {
		return false
	}

	return true
}

// Diff checks if two structs of type Program are equal
//
//	var a, b Program
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s Program) Diff(t Program, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if !equalPointers(s.Command, t.Command) {
		diff["Command"] = []interface{}{ValueOrNil(s.Command), ValueOrNil(t.Command)}
	}

	if s.Group != t.Group {
		diff["Group"] = []interface{}{s.Group, t.Group}
	}

	if s.Name != t.Name {
		diff["Name"] = []interface{}{s.Name, t.Name}
	}

	if s.StartOnReload != t.StartOnReload {
		diff["StartOnReload"] = []interface{}{s.StartOnReload, t.StartOnReload}
	}

	if s.User != t.User {
		diff["User"] = []interface{}{s.User, t.User}
	}

	return diff
}
