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

// Equal checks if two structs of type EmailAlert are equal
//
//	var a, b EmailAlert
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s EmailAlert) Equal(t EmailAlert, opts ...Options) bool {
	if !equalPointers(s.From, t.From) {
		return false
	}

	if s.Level != t.Level {
		return false
	}

	if !equalPointers(s.Mailers, t.Mailers) {
		return false
	}

	if s.Myhostname != t.Myhostname {
		return false
	}

	if !equalPointers(s.To, t.To) {
		return false
	}

	return true
}

// Diff checks if two structs of type EmailAlert are equal
//
//	var a, b EmailAlert
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s EmailAlert) Diff(t EmailAlert, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if !equalPointers(s.From, t.From) {
		diff["From"] = []interface{}{s.From, t.From}
	}

	if s.Level != t.Level {
		diff["Level"] = []interface{}{s.Level, t.Level}
	}

	if !equalPointers(s.Mailers, t.Mailers) {
		diff["Mailers"] = []interface{}{s.Mailers, t.Mailers}
	}

	if s.Myhostname != t.Myhostname {
		diff["Myhostname"] = []interface{}{s.Myhostname, t.Myhostname}
	}

	if !equalPointers(s.To, t.To) {
		diff["To"] = []interface{}{s.To, t.To}
	}

	return diff
}
