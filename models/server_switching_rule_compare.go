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

// Equal checks if two structs of type ServerSwitchingRule are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b ServerSwitchingRule
//	equal := a.Equal(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b ServerSwitchingRule
//	equal := a.Equal(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s ServerSwitchingRule) Equal(t ServerSwitchingRule, opts ...Options) bool {
	opt := getOptions(opts...)

	if s.Cond != t.Cond {
		return false
	}

	if s.CondTest != t.CondTest {
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

	if s.TargetServer != t.TargetServer {
		return false
	}

	return true
}

// Diff checks if two structs of type ServerSwitchingRule are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b ServerSwitchingRule
//	diff := a.Diff(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b ServerSwitchingRule
//	diff := a.Diff(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s ServerSwitchingRule) Diff(t ServerSwitchingRule, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})
	if s.Cond != t.Cond {
		diff["Cond"] = []interface{}{s.Cond, t.Cond}
	}

	if s.CondTest != t.CondTest {
		diff["CondTest"] = []interface{}{s.CondTest, t.CondTest}
	}

	if !CheckSameNilAndLenMap[string](s.Metadata, t.Metadata, opt) {
		diff["Metadata"] = []interface{}{s.Metadata, t.Metadata}
	}

	for k, v := range s.Metadata {
		if !reflect.DeepEqual(t.Metadata[k], v) {
			diff["Metadata"] = []interface{}{s.Metadata, t.Metadata}
		}
	}

	if s.TargetServer != t.TargetServer {
		diff["TargetServer"] = []interface{}{s.TargetServer, t.TargetServer}
	}

	return diff
}
