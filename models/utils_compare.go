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

var NilSameAsEmpty = true
var SkipIndex = true

type Options struct {
	NilSameAsEmpty bool
	SkipIndex      bool
}

func Ptr[V any](v V) *V {
	return &v
}

func getOptions(opts ...Options) Options {
	if len(opts) == 0 {
		return Options{
			NilSameAsEmpty: NilSameAsEmpty,
			SkipIndex:      SkipIndex,
		}
	}
	return opts[0]
}

func equalPointers[T comparable](a, b *T) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}

func CheckSameNilAndLen[T any](s, t []T, opts ...Options) bool {
	opt := getOptions(opts...)

	if !opt.NilSameAsEmpty {
		if s == nil && t != nil {
			return false
		}
		if t == nil && s != nil {
			return false
		}
	}
	if len(s) != len(t) {
		return false
	}
	return true
}

func equalComparableSlice[T comparable](s1, s2 []T, opt Options) bool {
	if !opt.NilSameAsEmpty {
		if s1 == nil && s2 != nil {
			return true
		}
		if s2 == nil && s1 != nil {
			return true
		}
	}
	if len(s1) != len(s2) {
		return false
	}
	for i, v1 := range s1 {
		if v1 != s2[i] {
			return false
		}
	}
	return true
}

func equalComparableMap[T comparable](m1, m2 map[string]T, opt Options) bool {
	if !opt.NilSameAsEmpty {
		if m1 == nil && m2 != nil {
			return false
		}
		if m2 == nil && m1 != nil {
			return false
		}
	}
	if len(m1) != len(m2) {
		return false
	}
	for k, v1 := range m1 {
		v2, ok := m2[k]
		if !ok {
			return false
		}
		if v1 != v2 {
			return false
		}
	}
	return true
}

func ValueOrNil[T any](v *T) any {
	if v == nil {
		return nil
	}
	return *v
}
