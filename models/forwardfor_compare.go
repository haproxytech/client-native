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

// Equal checks if two structs of type Forwardfor are equal
//
//	var a, b Forwardfor
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s Forwardfor) Equal(t Forwardfor, opts ...Options) bool {
	if !equalPointers(s.Enabled, t.Enabled) {
		return false
	}

	if s.Except != t.Except {
		return false
	}

	if s.Header != t.Header {
		return false
	}

	if s.Ifnone != t.Ifnone {
		return false
	}

	return true
}

// Diff checks if two structs of type Forwardfor are equal
//
//	var a, b Forwardfor
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s Forwardfor) Diff(t Forwardfor, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if !equalPointers(s.Enabled, t.Enabled) {
		diff["Enabled"] = []interface{}{ValueOrNil(s.Enabled), ValueOrNil(t.Enabled)}
	}

	if s.Except != t.Except {
		diff["Except"] = []interface{}{s.Except, t.Except}
	}

	if s.Header != t.Header {
		diff["Header"] = []interface{}{s.Header, t.Header}
	}

	if s.Ifnone != t.Ifnone {
		diff["Ifnone"] = []interface{}{s.Ifnone, t.Ifnone}
	}

	return diff
}
