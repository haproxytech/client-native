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

// Equal checks if two structs of type TuneLuaOptions are equal
//
//	var a, b TuneLuaOptions
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s TuneLuaOptions) Equal(t TuneLuaOptions, opts ...Options) bool {
	if !equalPointers(s.BurstTimeout, t.BurstTimeout) {
		return false
	}

	if s.ForcedYield != t.ForcedYield {
		return false
	}

	if s.LogLoggers != t.LogLoggers {
		return false
	}

	if s.LogStderr != t.LogStderr {
		return false
	}

	if !equalPointers(s.Maxmem, t.Maxmem) {
		return false
	}

	if !equalPointers(s.ServiceTimeout, t.ServiceTimeout) {
		return false
	}

	if !equalPointers(s.SessionTimeout, t.SessionTimeout) {
		return false
	}

	if !equalPointers(s.TaskTimeout, t.TaskTimeout) {
		return false
	}

	return true
}

// Diff checks if two structs of type TuneLuaOptions are equal
//
//	var a, b TuneLuaOptions
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s TuneLuaOptions) Diff(t TuneLuaOptions, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if !equalPointers(s.BurstTimeout, t.BurstTimeout) {
		diff["BurstTimeout"] = []interface{}{ValueOrNil(s.BurstTimeout), ValueOrNil(t.BurstTimeout)}
	}

	if s.ForcedYield != t.ForcedYield {
		diff["ForcedYield"] = []interface{}{s.ForcedYield, t.ForcedYield}
	}

	if s.LogLoggers != t.LogLoggers {
		diff["LogLoggers"] = []interface{}{s.LogLoggers, t.LogLoggers}
	}

	if s.LogStderr != t.LogStderr {
		diff["LogStderr"] = []interface{}{s.LogStderr, t.LogStderr}
	}

	if !equalPointers(s.Maxmem, t.Maxmem) {
		diff["Maxmem"] = []interface{}{ValueOrNil(s.Maxmem), ValueOrNil(t.Maxmem)}
	}

	if !equalPointers(s.ServiceTimeout, t.ServiceTimeout) {
		diff["ServiceTimeout"] = []interface{}{ValueOrNil(s.ServiceTimeout), ValueOrNil(t.ServiceTimeout)}
	}

	if !equalPointers(s.SessionTimeout, t.SessionTimeout) {
		diff["SessionTimeout"] = []interface{}{ValueOrNil(s.SessionTimeout), ValueOrNil(t.SessionTimeout)}
	}

	if !equalPointers(s.TaskTimeout, t.TaskTimeout) {
		diff["TaskTimeout"] = []interface{}{ValueOrNil(s.TaskTimeout), ValueOrNil(t.TaskTimeout)}
	}

	return diff
}
