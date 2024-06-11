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

// Equal checks if two structs of type LuaOptions are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b LuaOptions
//	equal := a.Equal(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b LuaOptions
//	equal := a.Equal(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s LuaOptions) Equal(t LuaOptions, opts ...Options) bool {
	opt := getOptions(opts...)

	if s.LoadPerThread != t.LoadPerThread {
		return false
	}

	if !CheckSameNilAndLen(s.Loads, t.Loads, opt) {
		return false
	} else {
		for i := range s.Loads {
			if !s.Loads[i].Equal(*t.Loads[i], opt) {
				return false
			}
		}
	}

	if !CheckSameNilAndLen(s.PrependPath, t.PrependPath, opt) {
		return false
	} else {
		for i := range s.PrependPath {
			if !s.PrependPath[i].Equal(*t.PrependPath[i], opt) {
				return false
			}
		}
	}

	return true
}

// Diff checks if two structs of type LuaOptions are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b LuaOptions
//	diff := a.Diff(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b LuaOptions
//	diff := a.Diff(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s LuaOptions) Diff(t LuaOptions, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})
	if s.LoadPerThread != t.LoadPerThread {
		diff["LoadPerThread"] = []interface{}{s.LoadPerThread, t.LoadPerThread}
	}

	if !CheckSameNilAndLen(s.Loads, t.Loads, opt) {
		diff["Loads"] = []interface{}{s.Loads, t.Loads}
	} else {
		diff2 := make(map[string][]interface{})
		for i := range s.Loads {
			if !s.Loads[i].Equal(*t.Loads[i], opt) {
				diffSub := s.Loads[i].Diff(*t.Loads[i], opt)
				if len(diffSub) > 0 {
					diff2[strconv.Itoa(i)] = []interface{}{diffSub}
				}
			}
		}
		if len(diff2) > 0 {
			diff["Loads"] = []interface{}{diff2}
		}
	}

	if !CheckSameNilAndLen(s.PrependPath, t.PrependPath, opt) {
		diff["PrependPath"] = []interface{}{s.PrependPath, t.PrependPath}
	} else {
		diff2 := make(map[string][]interface{})
		for i := range s.PrependPath {
			if !s.PrependPath[i].Equal(*t.PrependPath[i], opt) {
				diffSub := s.PrependPath[i].Diff(*t.PrependPath[i], opt)
				if len(diffSub) > 0 {
					diff2[strconv.Itoa(i)] = []interface{}{diffSub}
				}
			}
		}
		if len(diff2) > 0 {
			diff["PrependPath"] = []interface{}{diff2}
		}
	}

	return diff
}

// Equal checks if two structs of type LuaLoad are equal
//
//	var a, b LuaLoad
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s LuaLoad) Equal(t LuaLoad, opts ...Options) bool {
	if !equalPointers(s.File, t.File) {
		return false
	}

	return true
}

// Diff checks if two structs of type LuaLoad are equal
//
//	var a, b LuaLoad
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s LuaLoad) Diff(t LuaLoad, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if !equalPointers(s.File, t.File) {
		diff["File"] = []interface{}{ValueOrNil(s.File), ValueOrNil(t.File)}
	}

	return diff
}

// Equal checks if two structs of type LuaPrependPath are equal
//
//	var a, b LuaPrependPath
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s LuaPrependPath) Equal(t LuaPrependPath, opts ...Options) bool {
	if !equalPointers(s.Path, t.Path) {
		return false
	}

	if s.Type != t.Type {
		return false
	}

	return true
}

// Diff checks if two structs of type LuaPrependPath are equal
//
//	var a, b LuaPrependPath
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s LuaPrependPath) Diff(t LuaPrependPath, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if !equalPointers(s.Path, t.Path) {
		diff["Path"] = []interface{}{ValueOrNil(s.Path), ValueOrNil(t.Path)}
	}

	if s.Type != t.Type {
		diff["Type"] = []interface{}{s.Type, t.Type}
	}

	return diff
}
