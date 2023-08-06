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

// Equal checks if two structs of type Reload are equal
//
//	var a, b Reload
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s Reload) Equal(t Reload, opts ...Options) bool {
	if s.ID != t.ID {
		return false
	}

	if s.ReloadTimestamp != t.ReloadTimestamp {
		return false
	}

	if s.Response != t.Response {
		return false
	}

	if s.Status != t.Status {
		return false
	}

	return true
}

// Diff checks if two structs of type Reload are equal
//
//	var a, b Reload
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s Reload) Diff(t Reload, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if s.ID != t.ID {
		diff["ID"] = []interface{}{s.ID, t.ID}
	}

	if s.ReloadTimestamp != t.ReloadTimestamp {
		diff["ReloadTimestamp"] = []interface{}{s.ReloadTimestamp, t.ReloadTimestamp}
	}

	if s.Response != t.Response {
		diff["Response"] = []interface{}{s.Response, t.Response}
	}

	if s.Status != t.Status {
		diff["Status"] = []interface{}{s.Status, t.Status}
	}

	return diff
}
