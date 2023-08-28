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

// Equal checks if two structs of type LogForward are equal
//
//	var a, b LogForward
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s LogForward) Equal(t LogForward, opts ...Options) bool {
	if !equalPointers(s.Backlog, t.Backlog) {
		return false
	}

	if !equalPointers(s.Maxconn, t.Maxconn) {
		return false
	}

	if s.Name != t.Name {
		return false
	}

	if !equalPointers(s.TimeoutClient, t.TimeoutClient) {
		return false
	}

	return true
}

// Diff checks if two structs of type LogForward are equal
//
//	var a, b LogForward
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s LogForward) Diff(t LogForward, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if !equalPointers(s.Backlog, t.Backlog) {
		diff["Backlog"] = []interface{}{ValueOrNil(s.Backlog), ValueOrNil(t.Backlog)}
	}

	if !equalPointers(s.Maxconn, t.Maxconn) {
		diff["Maxconn"] = []interface{}{ValueOrNil(s.Maxconn), ValueOrNil(t.Maxconn)}
	}

	if s.Name != t.Name {
		diff["Name"] = []interface{}{s.Name, t.Name}
	}

	if !equalPointers(s.TimeoutClient, t.TimeoutClient) {
		diff["TimeoutClient"] = []interface{}{ValueOrNil(s.TimeoutClient), ValueOrNil(t.TimeoutClient)}
	}

	return diff
}
