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

// Equal checks if two structs of type LogTarget are equal
//
//	var a, b LogTarget
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s LogTarget) Equal(t LogTarget, opts ...Options) bool {
	if s.Address != t.Address {
		return false
	}

	if s.Facility != t.Facility {
		return false
	}

	if s.Format != t.Format {
		return false
	}

	if s.Global != t.Global {
		return false
	}

	if s.Length != t.Length {
		return false
	}

	if s.Level != t.Level {
		return false
	}

	if s.Minlevel != t.Minlevel {
		return false
	}

	if s.Nolog != t.Nolog {
		return false
	}

	if s.SampleRange != t.SampleRange {
		return false
	}

	if s.SampleSize != t.SampleSize {
		return false
	}

	return true
}

// Diff checks if two structs of type LogTarget are equal
//
//	var a, b LogTarget
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s LogTarget) Diff(t LogTarget, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if s.Address != t.Address {
		diff["Address"] = []interface{}{s.Address, t.Address}
	}

	if s.Facility != t.Facility {
		diff["Facility"] = []interface{}{s.Facility, t.Facility}
	}

	if s.Format != t.Format {
		diff["Format"] = []interface{}{s.Format, t.Format}
	}

	if s.Global != t.Global {
		diff["Global"] = []interface{}{s.Global, t.Global}
	}

	if s.Length != t.Length {
		diff["Length"] = []interface{}{s.Length, t.Length}
	}

	if s.Level != t.Level {
		diff["Level"] = []interface{}{s.Level, t.Level}
	}

	if s.Minlevel != t.Minlevel {
		diff["Minlevel"] = []interface{}{s.Minlevel, t.Minlevel}
	}

	if s.Nolog != t.Nolog {
		diff["Nolog"] = []interface{}{s.Nolog, t.Nolog}
	}

	if s.SampleRange != t.SampleRange {
		diff["SampleRange"] = []interface{}{s.SampleRange, t.SampleRange}
	}

	if s.SampleSize != t.SampleSize {
		diff["SampleSize"] = []interface{}{s.SampleSize, t.SampleSize}
	}

	return diff
}
