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

// Equal checks if two structs of type EnvironmentOptions are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b EnvironmentOptions
//	equal := a.Equal(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b EnvironmentOptions
//	equal := a.Equal(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s EnvironmentOptions) Equal(t EnvironmentOptions, opts ...Options) bool {
	opt := getOptions(opts...)

	if !CheckSameNilAndLen(s.PresetEnvs, t.PresetEnvs, opt) {
		return false
	} else {
		for i := range s.PresetEnvs {
			if !s.PresetEnvs[i].Equal(*t.PresetEnvs[i], opt) {
				return false
			}
		}
	}

	if !CheckSameNilAndLen(s.SetEnvs, t.SetEnvs, opt) {
		return false
	} else {
		for i := range s.SetEnvs {
			if !s.SetEnvs[i].Equal(*t.SetEnvs[i], opt) {
				return false
			}
		}
	}

	if s.Resetenv != t.Resetenv {
		return false
	}

	if s.Unsetenv != t.Unsetenv {
		return false
	}

	return true
}

// Diff checks if two structs of type EnvironmentOptions are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b EnvironmentOptions
//	diff := a.Diff(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b EnvironmentOptions
//	diff := a.Diff(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s EnvironmentOptions) Diff(t EnvironmentOptions, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})
	if !CheckSameNilAndLen(s.PresetEnvs, t.PresetEnvs, opt) {
		diff["PresetEnvs"] = []interface{}{s.PresetEnvs, t.PresetEnvs}
	} else {
		diff2 := make(map[string][]interface{})
		for i := range s.PresetEnvs {
			if !s.PresetEnvs[i].Equal(*t.PresetEnvs[i], opt) {
				diffSub := s.PresetEnvs[i].Diff(*t.PresetEnvs[i], opt)
				if len(diffSub) > 0 {
					diff2[strconv.Itoa(i)] = []interface{}{diffSub}
				}
			}
		}
		if len(diff2) > 0 {
			diff["PresetEnvs"] = []interface{}{diff2}
		}
	}

	if !CheckSameNilAndLen(s.SetEnvs, t.SetEnvs, opt) {
		diff["SetEnvs"] = []interface{}{s.SetEnvs, t.SetEnvs}
	} else {
		diff2 := make(map[string][]interface{})
		for i := range s.SetEnvs {
			if !s.SetEnvs[i].Equal(*t.SetEnvs[i], opt) {
				diffSub := s.SetEnvs[i].Diff(*t.SetEnvs[i], opt)
				if len(diffSub) > 0 {
					diff2[strconv.Itoa(i)] = []interface{}{diffSub}
				}
			}
		}
		if len(diff2) > 0 {
			diff["SetEnvs"] = []interface{}{diff2}
		}
	}

	if s.Resetenv != t.Resetenv {
		diff["Resetenv"] = []interface{}{s.Resetenv, t.Resetenv}
	}

	if s.Unsetenv != t.Unsetenv {
		diff["Unsetenv"] = []interface{}{s.Unsetenv, t.Unsetenv}
	}

	return diff
}

// Equal checks if two structs of type PresetEnv are equal
//
//	var a, b PresetEnv
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s PresetEnv) Equal(t PresetEnv, opts ...Options) bool {
	if !equalPointers(s.Name, t.Name) {
		return false
	}

	if !equalPointers(s.Value, t.Value) {
		return false
	}

	return true
}

// Diff checks if two structs of type PresetEnv are equal
//
//	var a, b PresetEnv
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s PresetEnv) Diff(t PresetEnv, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if !equalPointers(s.Name, t.Name) {
		diff["Name"] = []interface{}{ValueOrNil(s.Name), ValueOrNil(t.Name)}
	}

	if !equalPointers(s.Value, t.Value) {
		diff["Value"] = []interface{}{ValueOrNil(s.Value), ValueOrNil(t.Value)}
	}

	return diff
}

// Equal checks if two structs of type SetEnv are equal
//
//	var a, b SetEnv
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s SetEnv) Equal(t SetEnv, opts ...Options) bool {
	if !equalPointers(s.Name, t.Name) {
		return false
	}

	if !equalPointers(s.Value, t.Value) {
		return false
	}

	return true
}

// Diff checks if two structs of type SetEnv are equal
//
//	var a, b SetEnv
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s SetEnv) Diff(t SetEnv, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if !equalPointers(s.Name, t.Name) {
		diff["Name"] = []interface{}{ValueOrNil(s.Name), ValueOrNil(t.Name)}
	}

	if !equalPointers(s.Value, t.Value) {
		diff["Value"] = []interface{}{ValueOrNil(s.Value), ValueOrNil(t.Value)}
	}

	return diff
}
