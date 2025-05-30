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

// Equal checks if two structs of type SslCrlEntry are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b SslCrlEntry
//	equal := a.Equal(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b SslCrlEntry
//	equal := a.Equal(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s SslCrlEntry) Equal(t SslCrlEntry, opts ...Options) bool {
	opt := getOptions(opts...)

	if s.Issuer != t.Issuer {
		return false
	}

	if !s.LastUpdate.Equal(t.LastUpdate) {
		return false
	}

	if !s.NextUpdate.Equal(t.NextUpdate) {
		return false
	}

	if !CheckSameNilAndLen(s.RevokedCertificates, t.RevokedCertificates, opt) {
		return false
	} else {
		for i := range s.RevokedCertificates {
			if !s.RevokedCertificates[i].Equal(*t.RevokedCertificates[i], opt) {
				return false
			}
		}
	}

	if s.SignatureAlgorithm != t.SignatureAlgorithm {
		return false
	}

	if s.Status != t.Status {
		return false
	}

	if s.StorageName != t.StorageName {
		return false
	}

	if s.Version != t.Version {
		return false
	}

	return true
}

// Diff checks if two structs of type SslCrlEntry are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b SslCrlEntry
//	diff := a.Diff(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b SslCrlEntry
//	diff := a.Diff(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s SslCrlEntry) Diff(t SslCrlEntry, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})
	if s.Issuer != t.Issuer {
		diff["Issuer"] = []interface{}{s.Issuer, t.Issuer}
	}

	if !s.LastUpdate.Equal(t.LastUpdate) {
		diff["LastUpdate"] = []interface{}{s.LastUpdate, t.LastUpdate}
	}

	if !s.NextUpdate.Equal(t.NextUpdate) {
		diff["NextUpdate"] = []interface{}{s.NextUpdate, t.NextUpdate}
	}

	if !CheckSameNilAndLen(s.RevokedCertificates, t.RevokedCertificates, opt) {
		diff["RevokedCertificates"] = []interface{}{s.RevokedCertificates, t.RevokedCertificates}
	} else {
		diff2 := make(map[string][]interface{})
		for i := range s.RevokedCertificates {
			if !s.RevokedCertificates[i].Equal(*t.RevokedCertificates[i], opt) {
				diffSub := s.RevokedCertificates[i].Diff(*t.RevokedCertificates[i], opt)
				if len(diffSub) > 0 {
					diff2[strconv.Itoa(i)] = []interface{}{diffSub}
				}
			}
		}
		if len(diff2) > 0 {
			diff["RevokedCertificates"] = []interface{}{diff2}
		}
	}

	if s.SignatureAlgorithm != t.SignatureAlgorithm {
		diff["SignatureAlgorithm"] = []interface{}{s.SignatureAlgorithm, t.SignatureAlgorithm}
	}

	if s.Status != t.Status {
		diff["Status"] = []interface{}{s.Status, t.Status}
	}

	if s.StorageName != t.StorageName {
		diff["StorageName"] = []interface{}{s.StorageName, t.StorageName}
	}

	if s.Version != t.Version {
		diff["Version"] = []interface{}{s.Version, t.Version}
	}

	return diff
}

// Equal checks if two structs of type RevokedCertificates are equal
//
//	var a, b RevokedCertificates
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s RevokedCertificates) Equal(t RevokedCertificates, opts ...Options) bool {

	if !s.RevocationDate.Equal(t.RevocationDate) {
		return false
	}

	if s.SerialNumber != t.SerialNumber {
		return false
	}

	return true
}

// Diff checks if two structs of type RevokedCertificates are equal
//
//	var a, b RevokedCertificates
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s RevokedCertificates) Diff(t RevokedCertificates, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})

	if !s.RevocationDate.Equal(t.RevocationDate) {
		diff["RevocationDate"] = []interface{}{s.RevocationDate, t.RevocationDate}
	}

	if s.SerialNumber != t.SerialNumber {
		diff["SerialNumber"] = []interface{}{s.SerialNumber, t.SerialNumber}
	}

	return diff
}
