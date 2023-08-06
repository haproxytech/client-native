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

// Equal checks if two structs of type SslCertificate are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b SslCertificate
//	equal := a.Equal(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b SslCertificate
//	equal := a.Equal(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s SslCertificate) Equal(t SslCertificate, opts ...Options) bool {
	opt := getOptions(opts...)

	if s.Description != t.Description {
		return false
	}

	if !equalComparableSlice(s.Domains, t.Domains, opt) {
		return false
	}

	if s.File != t.File {
		return false
	}

	if !equalComparableSlice(s.IPAddresses, t.IPAddresses, opt) {
		return false
	}

	if !equalComparableSlice(s.Issuers, t.Issuers, opt) {
		return false
	}

	if !s.NotAfter.Equal(t.NotAfter) {
		return false
	}

	if !s.NotBefore.Equal(t.NotBefore) {
		return false
	}

	if s.Size != t.Size {
		return false
	}

	if s.StorageName != t.StorageName {
		return false
	}

	return true
}

// Diff checks if two structs of type SslCertificate are equal
//
// By default empty arrays, maps and slices are equal to nil:
//
//	var a, b SslCertificate
//	diff := a.Diff(b)
//
// For more advanced use case you can configure the options (default values are shown):
//
//	var a, b SslCertificate
//	equal := a.Diff(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s SslCertificate) Diff(t SslCertificate, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})
	if s.Description != t.Description {
		diff["Description"] = []interface{}{s.Description, t.Description}
	}

	if !equalComparableSlice(s.Domains, t.Domains, opt) {
		diff["Domains"] = []interface{}{s.Domains, t.Domains}
	}

	if s.File != t.File {
		diff["File"] = []interface{}{s.File, t.File}
	}

	if !equalComparableSlice(s.IPAddresses, t.IPAddresses, opt) {
		diff["IPAddresses"] = []interface{}{s.IPAddresses, t.IPAddresses}
	}

	if !equalComparableSlice(s.Issuers, t.Issuers, opt) {
		diff["Issuers"] = []interface{}{s.Issuers, t.Issuers}
	}

	if !s.NotAfter.Equal(t.NotAfter) {
		diff["NotAfter"] = []interface{}{s.NotAfter, t.NotAfter}
	}

	if !s.NotBefore.Equal(t.NotBefore) {
		diff["NotBefore"] = []interface{}{s.NotBefore, t.NotBefore}
	}

	if s.Size != t.Size {
		diff["Size"] = []interface{}{s.Size, t.Size}
	}

	if s.StorageName != t.StorageName {
		diff["StorageName"] = []interface{}{s.StorageName, t.StorageName}
	}

	return diff
}
