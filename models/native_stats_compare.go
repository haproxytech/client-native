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

import (
	"strconv"
)

// Equal checks if two structs of type NativeStats are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b NativeStats
//	equal := a.Equal(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b NativeStats
//	equal := a.Equal(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s NativeStats) Equal(t NativeStats, opts ...Options) bool {
	opt := getOptions(opts...)

	if s.Error != t.Error {
		return false
	}

	if s.RuntimeAPI != t.RuntimeAPI {
		return false
	}

	if !CheckSameNilAndLen(s.Stats, t.Stats, opt) {
		return false
	} else {
		for i := range s.Stats {
			if !s.Stats[i].Equal(*t.Stats[i], opt) {
				return false
			}
		}
	}

	return true
}

// Diff checks if two structs of type NativeStats are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b NativeStats
//	diff := a.Diff(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b NativeStats
//	diff := a.Diff(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s NativeStats) Diff(t NativeStats, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})
	if s.Error != t.Error {
		diff["Error"] = []interface{}{s.Error, t.Error}
	}

	if s.RuntimeAPI != t.RuntimeAPI {
		diff["RuntimeAPI"] = []interface{}{s.RuntimeAPI, t.RuntimeAPI}
	}

	if !CheckSameNilAndLen(s.Stats, t.Stats, opt) {
		diff["Stats"] = []interface{}{s.Stats, t.Stats}
	} else {
		diff2 := make(map[string][]interface{})
		for i := range s.Stats {
			if !s.Stats[i].Equal(*t.Stats[i], opt) {
				diffSub := s.Stats[i].Diff(*t.Stats[i], opt)
				if len(diffSub) > 0 {
					diff2[strconv.Itoa(i)] = []interface{}{diffSub}
				}
			}
		}
		if len(diff2) > 0 {
			diff["Stats"] = []interface{}{diff2}
		}
	}

	return diff
}
