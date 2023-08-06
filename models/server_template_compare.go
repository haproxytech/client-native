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

// Equal checks if two structs of type ServerTemplate are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b ServerTemplate
//	equal := a.Equal(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b ServerTemplate
//	equal := a.Equal(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s ServerTemplate) Equal(t ServerTemplate, opts ...Options) bool {
	opt := getOptions(opts...)

	if !s.ServerParams.Equal(t.ServerParams, opt) {
		return false
	}

	if s.Fqdn != t.Fqdn {
		return false
	}

	if !equalPointers(s.ID, t.ID) {
		return false
	}

	if s.NumOrRange != t.NumOrRange {
		return false
	}

	if !equalPointers(s.Port, t.Port) {
		return false
	}

	if s.Prefix != t.Prefix {
		return false
	}

	return true
}

// Diff checks if two structs of type ServerTemplate are equal
//
// By default empty arrays, maps and slices are equal to nil:
//
//	var a, b ServerTemplate
//	diff := a.Diff(b)
//
// For more advanced use case you can configure the options (default values are shown):
//
//	var a, b ServerTemplate
//	equal := a.Diff(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s ServerTemplate) Diff(t ServerTemplate, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})
	if !s.ServerParams.Equal(t.ServerParams, opt) {
		diff["ServerParams"] = []interface{}{s.ServerParams, t.ServerParams}
	}

	if s.Fqdn != t.Fqdn {
		diff["Fqdn"] = []interface{}{s.Fqdn, t.Fqdn}
	}

	if !equalPointers(s.ID, t.ID) {
		diff["ID"] = []interface{}{s.ID, t.ID}
	}

	if s.NumOrRange != t.NumOrRange {
		diff["NumOrRange"] = []interface{}{s.NumOrRange, t.NumOrRange}
	}

	if !equalPointers(s.Port, t.Port) {
		diff["Port"] = []interface{}{s.Port, t.Port}
	}

	if s.Prefix != t.Prefix {
		diff["Prefix"] = []interface{}{s.Prefix, t.Prefix}
	}

	return diff
}
