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

// Equal checks if two structs of type ProcessInfo are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b ProcessInfo
//	equal := a.Equal(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b ProcessInfo
//	equal := a.Equal(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s ProcessInfo) Equal(t ProcessInfo, opts ...Options) bool {
	opt := getOptions(opts...)

	if s.Error != t.Error {
		return false
	}

	if s.Info == nil || t.Info == nil {
		if s.Info != nil || t.Info != nil {
			if opt.NilSameAsEmpty {
				empty := &ProcessInfoItem{}
				if s.Info == nil {
					if !(t.Info.Equal(*empty)) {
						return false
					}
				}
				if t.Info == nil {
					if !(s.Info.Equal(*empty)) {
						return false
					}
				}
			} else {
				return false
			}
		}
	} else if !s.Info.Equal(*t.Info, opt) {
		return false
	}

	if s.RuntimeAPI != t.RuntimeAPI {
		return false
	}

	return true
}

// Diff checks if two structs of type ProcessInfo are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b ProcessInfo
//	diff := a.Diff(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b ProcessInfo
//	diff := a.Diff(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s ProcessInfo) Diff(t ProcessInfo, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})
	if s.Error != t.Error {
		diff["Error"] = []interface{}{s.Error, t.Error}
	}

	if s.Info == nil || t.Info == nil {
		if s.Info != nil || t.Info != nil {
			if opt.NilSameAsEmpty {
				empty := &ProcessInfoItem{}
				if s.Info == nil {
					if !(t.Info.Equal(*empty)) {
						diff["Info"] = []interface{}{ValueOrNil(s.Info), ValueOrNil(t.Info)}
					}
				}
				if t.Info == nil {
					if !(s.Info.Equal(*empty)) {
						diff["Info"] = []interface{}{ValueOrNil(s.Info), ValueOrNil(t.Info)}
					}
				}
			} else {
				diff["Info"] = []interface{}{ValueOrNil(s.Info), ValueOrNil(t.Info)}
			}
		}
	} else if !s.Info.Equal(*t.Info, opt) {
		diff["Info"] = []interface{}{ValueOrNil(s.Info), ValueOrNil(t.Info)}
	}

	if s.RuntimeAPI != t.RuntimeAPI {
		diff["RuntimeAPI"] = []interface{}{s.RuntimeAPI, t.RuntimeAPI}
	}

	return diff
}
