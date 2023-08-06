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

// Equal checks if two structs of type FCGILogStderr are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b FCGILogStderr
//	equal := a.Equal(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b FCGILogStderr
//	equal := a.Equal(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s FCGILogStderr) Equal(t FCGILogStderr, opts ...Options) bool {
	opt := getOptions(opts...)

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

	if s.Len != t.Len {
		return false
	}

	if s.Level != t.Level {
		return false
	}

	if s.Minlevel != t.Minlevel {
		return false
	}

	if !s.Sample.Equal(*t.Sample, opt) {
		return false
	}

	return true
}

// Diff checks if two structs of type FCGILogStderr are equal
//
// By default empty arrays, maps and slices are equal to nil:
//
//	var a, b FCGILogStderr
//	diff := a.Diff(b)
//
// For more advanced use case you can configure the options (default values are shown):
//
//	var a, b FCGILogStderr
//	equal := a.Diff(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s FCGILogStderr) Diff(t FCGILogStderr, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

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

	if s.Len != t.Len {
		diff["Len"] = []interface{}{s.Len, t.Len}
	}

	if s.Level != t.Level {
		diff["Level"] = []interface{}{s.Level, t.Level}
	}

	if s.Minlevel != t.Minlevel {
		diff["Minlevel"] = []interface{}{s.Minlevel, t.Minlevel}
	}

	if !s.Sample.Equal(*t.Sample, opt) {
		diff["Sample"] = []interface{}{s.Sample, t.Sample}
	}

	return diff
}

// Equal checks if two structs of type FCGILogStderrSample are equal
//
//	var a, b FCGILogStderrSample
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s FCGILogStderrSample) Equal(t FCGILogStderrSample, opts ...Options) bool {
	if !equalPointers(s.Ranges, t.Ranges) {
		return false
	}

	if !equalPointers(s.Size, t.Size) {
		return false
	}

	return true
}

// Diff checks if two structs of type FCGILogStderrSample are equal
//
//	var a, b FCGILogStderrSample
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s FCGILogStderrSample) Diff(t FCGILogStderrSample, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if !equalPointers(s.Ranges, t.Ranges) {
		diff["Ranges"] = []interface{}{s.Ranges, t.Ranges}
	}

	if !equalPointers(s.Size, t.Size) {
		diff["Size"] = []interface{}{s.Size, t.Size}
	}

	return diff
}
