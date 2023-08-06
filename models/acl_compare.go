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
// By default empty maps and slices are equal to nil:
//
//	var a, b ACL
//	equal := a.Equal(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b ACL
//	equal := a.Equal(b,Options{
//		SkipIndex: true,
//	})
func (s ACL) Equal(t ACL, opts ...Options) bool {
	opt := getOptions(opts...)

	if s.ACLName != t.ACLName {
		return false
	}

	if s.Criterion != t.Criterion {
		return false
	}

	if !opt.SkipIndex && !equalPointers(s.Index, t.Index) {
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
// For more advanced use case you can configure the options (default values are shown):
//
//	var a, b ACL
//	equal := a.Diff(b,Options{
//		SkipIndex: true,
//	})
func (s ACL) Diff(t ACL, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})
	if s.ACLName != t.ACLName {
		diff["ACLName"] = []interface{}{s.ACLName, t.ACLName}
	}

	if s.Criterion != t.Criterion {
		diff["Criterion"] = []interface{}{s.Criterion, t.Criterion}
	}

	if !opt.SkipIndex && !equalPointers(s.Index, t.Index) {
		diff["Index"] = []interface{}{s.Index, t.Index}
	}

	if s.Value != t.Value {
		diff["Value"] = []interface{}{s.Value, t.Value}
	}

	return diff
}
