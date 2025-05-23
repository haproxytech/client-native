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

// Equal checks if two structs of type DebugOptions are equal
//
//	var a, b DebugOptions
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s DebugOptions) Equal(t DebugOptions, opts ...Options) bool {
	if !equalPointers(s.Anonkey, t.Anonkey) {
		return false
	}

	if s.Quiet != t.Quiet {
		return false
	}

	if !equalPointers(s.StressLevel, t.StressLevel) {
		return false
	}

	if s.ZeroWarning != t.ZeroWarning {
		return false
	}

	return true
}

// Diff checks if two structs of type DebugOptions are equal
//
//	var a, b DebugOptions
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s DebugOptions) Diff(t DebugOptions, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if !equalPointers(s.Anonkey, t.Anonkey) {
		diff["Anonkey"] = []interface{}{ValueOrNil(s.Anonkey), ValueOrNil(t.Anonkey)}
	}

	if s.Quiet != t.Quiet {
		diff["Quiet"] = []interface{}{s.Quiet, t.Quiet}
	}

	if !equalPointers(s.StressLevel, t.StressLevel) {
		diff["StressLevel"] = []interface{}{ValueOrNil(s.StressLevel), ValueOrNil(t.StressLevel)}
	}

	if s.ZeroWarning != t.ZeroWarning {
		diff["ZeroWarning"] = []interface{}{s.ZeroWarning, t.ZeroWarning}
	}

	return diff
}
