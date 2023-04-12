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

// Equal checks if two structs of type SslCrtListEntry are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b SslCrtListEntry
//	equal := a.Equal(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b SslCrtListEntry
//	equal := a.Equal(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s SslCrtListEntry) Equal(t SslCrtListEntry, opts ...Options) bool {
	opt := getOptions(opts...)

	if !equalComparableSlice(s.SNIFilter, t.SNIFilter, opt) {
		return false
	}

	if s.SSLBindConfig != t.SSLBindConfig {
		return false
	}

	if s.File != t.File {
		return false
	}

	if s.LineNumber != t.LineNumber {
		return false
	}

	return true
}

// Diff checks if two structs of type SslCrtListEntry are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b SslCrtListEntry
//	diff := a.Diff(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b SslCrtListEntry
//	diff := a.Diff(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s SslCrtListEntry) Diff(t SslCrtListEntry, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})
	if !equalComparableSlice(s.SNIFilter, t.SNIFilter, opt) {
		diff["SNIFilter"] = []interface{}{s.SNIFilter, t.SNIFilter}
	}

	if s.SSLBindConfig != t.SSLBindConfig {
		diff["SSLBindConfig"] = []interface{}{s.SSLBindConfig, t.SSLBindConfig}
	}

	if s.File != t.File {
		diff["File"] = []interface{}{s.File, t.File}
	}

	if s.LineNumber != t.LineNumber {
		diff["LineNumber"] = []interface{}{s.LineNumber, t.LineNumber}
	}

	return diff
}
