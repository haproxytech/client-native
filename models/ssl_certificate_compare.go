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
//	var a, b SslCertificate
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s SslCertificate) Equal(t SslCertificate, opts ...Options) bool {
	if s.Description != t.Description {
		return false
	}

	if s.Domains != t.Domains {
		return false
	}

	if s.File != t.File {
		return false
	}

	if s.IPAddresses != t.IPAddresses {
		return false
	}

	if s.Issuers != t.Issuers {
		return false
	}

	if !s.NotAfter.Equal(*t.NotAfter) {
		return false
	}

	if !s.NotBefore.Equal(*t.NotBefore) {
		return false
	}

	if !equalPointers(s.Size, t.Size) {
		return false
	}

	if s.StorageName != t.StorageName {
		return false
	}

	return true
}

// Diff checks if two structs of type SslCertificate are equal
//
//	var a, b SslCertificate
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s SslCertificate) Diff(t SslCertificate, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if s.Description != t.Description {
		diff["Description"] = []interface{}{s.Description, t.Description}
	}

	if s.Domains != t.Domains {
		diff["Domains"] = []interface{}{s.Domains, t.Domains}
	}

	if s.File != t.File {
		diff["File"] = []interface{}{s.File, t.File}
	}

	if s.IPAddresses != t.IPAddresses {
		diff["IPAddresses"] = []interface{}{s.IPAddresses, t.IPAddresses}
	}

	if s.Issuers != t.Issuers {
		diff["Issuers"] = []interface{}{s.Issuers, t.Issuers}
	}

	if !s.NotAfter.Equal(*t.NotAfter) {
		diff["NotAfter"] = []interface{}{s.NotAfter, t.NotAfter}
	}

	if !s.NotBefore.Equal(*t.NotBefore) {
		diff["NotBefore"] = []interface{}{s.NotBefore, t.NotBefore}
	}

	if !equalPointers(s.Size, t.Size) {
		diff["Size"] = []interface{}{s.Size, t.Size}
	}

	if s.StorageName != t.StorageName {
		diff["StorageName"] = []interface{}{s.StorageName, t.StorageName}
	}

	return diff
}
