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

// Equal checks if two structs of type SslProviders are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b SslProviders
//	equal := a.Equal(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b SslProviders
//	equal := a.Equal(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s SslProviders) Equal(t SslProviders, opts ...Options) bool {
	opt := getOptions(opts...)

	if !equalComparableSlice(s.Providers, t.Providers, opt) {
		return false
	}

	return true
}

// Diff checks if two structs of type SslProviders are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b SslProviders
//	diff := a.Diff(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b SslProviders
//	diff := a.Diff(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s SslProviders) Diff(t SslProviders, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})
	if !equalComparableSlice(s.Providers, t.Providers, opt) {
		diff["Providers"] = []interface{}{s.Providers, t.Providers}
	}

	return diff
}
