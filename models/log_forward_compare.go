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
// By default empty maps and slices are equal to nil:
//
//	var a, b LogForward
//	equal := a.Equal(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b LogForward
//	equal := a.Equal(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s LogForward) Equal(t LogForward, opts ...Options) bool {
	opt := getOptions(opts...)

	if !s.LogForwardBase.Equal(t.LogForwardBase, opt) {
		return false
	}

	if !s.LogTargetList.Equal(t.LogTargetList, opt) {
		return false
	}

	if !CheckSameNilAndLenMap[string, Bind](s.Binds, t.Binds, opt) {
		return false
	}

	for k, v := range s.Binds {
		if !t.Binds[k].Equal(v, opt) {
			return false
		}
	}

	if !CheckSameNilAndLenMap[string, DgramBind](s.DgramBinds, t.DgramBinds, opt) {
		return false
	}

	for k, v := range s.DgramBinds {
		if !t.DgramBinds[k].Equal(v, opt) {
			return false
		}
	}

	return true
}

// Diff checks if two structs of type LogForward are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b LogForward
//	diff := a.Diff(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b LogForward
//	diff := a.Diff(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s LogForward) Diff(t LogForward, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})

	if !s.LogForwardBase.Equal(t.LogForwardBase, opt) {
		diff["LogForwardBase"] = []interface{}{s.LogForwardBase, t.LogForwardBase}
	}

	if !s.LogTargetList.Equal(t.LogTargetList, opt) {
		diff["LogTargetList"] = []interface{}{s.LogTargetList, t.LogTargetList}
	}

	if !CheckSameNilAndLenMap[string, Bind](s.Binds, t.Binds, opt) {
		diff["Binds"] = []interface{}{s.Binds, t.Binds}
	}

	for k, v := range s.Binds {
		if !t.Binds[k].Equal(v, opt) {
			diff["Binds"] = []interface{}{s.Binds, t.Binds}
		}
	}

	if !CheckSameNilAndLenMap[string, DgramBind](s.DgramBinds, t.DgramBinds, opt) {
		diff["DgramBinds"] = []interface{}{s.DgramBinds, t.DgramBinds}
	}

	for k, v := range s.DgramBinds {
		if !t.DgramBinds[k].Equal(v, opt) {
			diff["DgramBinds"] = []interface{}{s.DgramBinds, t.DgramBinds}
		}
	}

	return diff
}
