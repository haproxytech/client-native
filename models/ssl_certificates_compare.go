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
	"fmt"
)

// Equal checks if two structs of type SslCertificates are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b SslCertificates
//	equal := a.Equal(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b SslCertificates
//	equal := a.Equal(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s SslCertificates) Equal(t SslCertificates, opts ...Options) bool {
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
	for i, v := range s {
		if !v.Equal(*t[i], opt) {
			return false
		}

	}
	return true
}

// Diff checks if two structs of type SslCertificates are equal
//
// By default empty arrays, maps and slices are equal to nil:
//
//	var a, b SslCertificates
//	diff := a.Diff(b)
//
// For more advanced use case you can configure the options (default values are shown):
//
//	var a, b SslCertificates
//	equal := a.Diff(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s SslCertificates) Diff(t SslCertificates, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})
	if !opt.NilSameAsEmpty {
		if s == nil && t != nil {
			diff["SslCertificates"] = []interface{}{s, t}
			return diff
		}
		if t == nil && s != nil {
			diff["SslCertificates"] = []interface{}{s, t}
			return diff
		}
	}
	if len(s) != len(t) {
		diff["SslCertificates"] = []interface{}{s, t}
		return diff
	}
	for i, v := range s {
		if !v.Equal(*t[i], opt) {
			diff[fmt.Sprintf("SslCertificates[%d]", i)] = []interface{}{v, t[i]}
		}

	}
	return diff
}
