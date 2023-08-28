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

// Equal checks if two structs of type MonitorFail are equal
//
//	var a, b MonitorFail
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s MonitorFail) Equal(t MonitorFail, opts ...Options) bool {
	if !equalPointers(s.Cond, t.Cond) {
		return false
	}

	if !equalPointers(s.CondTest, t.CondTest) {
		return false
	}

	return true
}

// Diff checks if two structs of type MonitorFail are equal
//
//	var a, b MonitorFail
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s MonitorFail) Diff(t MonitorFail, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if !equalPointers(s.Cond, t.Cond) {
		diff["Cond"] = []interface{}{ValueOrNil(s.Cond), ValueOrNil(t.Cond)}
	}

	if !equalPointers(s.CondTest, t.CondTest) {
		diff["CondTest"] = []interface{}{ValueOrNil(s.CondTest), ValueOrNil(t.CondTest)}
	}

	return diff
}
