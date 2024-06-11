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

// Equal checks if two structs of type ACL are equal
//
//	var a, b ACL
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s ACL) Equal(t ACL, opts ...Options) bool {
	if s.ACLName != t.ACLName {
		return false
	}

	if s.Criterion != t.Criterion {
		return false
	}

	if s.Value != t.Value {
		return false
	}

	return true
}

// Diff checks if two structs of type ACL are equal
//
//	var a, b ACL
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s ACL) Diff(t ACL, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if s.ACLName != t.ACLName {
		diff["ACLName"] = []interface{}{s.ACLName, t.ACLName}
	}

	if s.Criterion != t.Criterion {
		diff["Criterion"] = []interface{}{s.Criterion, t.Criterion}
	}

	if s.Value != t.Value {
		diff["Value"] = []interface{}{s.Value, t.Value}
	}

	return diff
}
