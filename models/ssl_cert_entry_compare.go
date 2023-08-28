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

// Equal checks if two structs of type SslCertEntry are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b SslCertEntry
//	equal := a.Equal(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b SslCertEntry
//	equal := a.Equal(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s SslCertEntry) Equal(t SslCertEntry, opts ...Options) bool {
	opt := getOptions(opts...)

	if s.Algorithm != t.Algorithm {
		return false
	}

	if s.ChainIssuer != t.ChainIssuer {
		return false
	}

	if s.ChainSubject != t.ChainSubject {
		return false
	}

	if s.Issuer != t.Issuer {
		return false
	}

	if !s.NotAfter.Equal(t.NotAfter) {
		return false
	}

	if !s.NotBefore.Equal(t.NotBefore) {
		return false
	}

	if s.Serial != t.Serial {
		return false
	}

	if s.Sha1FingerPrint != t.Sha1FingerPrint {
		return false
	}

	if s.Status != t.Status {
		return false
	}

	if s.StorageName != t.StorageName {
		return false
	}

	if s.Subject != t.Subject {
		return false
	}

	if !equalComparableSlice(s.SubjectAlternativeNames, t.SubjectAlternativeNames, opt) {
		return false
	}

	return true
}

// Diff checks if two structs of type SslCertEntry are equal
//
// By default empty arrays, maps and slices are equal to nil:
//
//	var a, b SslCertEntry
//	diff := a.Diff(b)
//
// For more advanced use case you can configure the options (default values are shown):
//
//	var a, b SslCertEntry
//	equal := a.Diff(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s SslCertEntry) Diff(t SslCertEntry, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})
	if s.Algorithm != t.Algorithm {
		diff["Algorithm"] = []interface{}{s.Algorithm, t.Algorithm}
	}

	if s.ChainIssuer != t.ChainIssuer {
		diff["ChainIssuer"] = []interface{}{s.ChainIssuer, t.ChainIssuer}
	}

	if s.ChainSubject != t.ChainSubject {
		diff["ChainSubject"] = []interface{}{s.ChainSubject, t.ChainSubject}
	}

	if s.Issuer != t.Issuer {
		diff["Issuer"] = []interface{}{s.Issuer, t.Issuer}
	}

	if !s.NotAfter.Equal(t.NotAfter) {
		diff["NotAfter"] = []interface{}{s.NotAfter, t.NotAfter}
	}

	if !s.NotBefore.Equal(t.NotBefore) {
		diff["NotBefore"] = []interface{}{s.NotBefore, t.NotBefore}
	}

	if s.Serial != t.Serial {
		diff["Serial"] = []interface{}{s.Serial, t.Serial}
	}

	if s.Sha1FingerPrint != t.Sha1FingerPrint {
		diff["Sha1FingerPrint"] = []interface{}{s.Sha1FingerPrint, t.Sha1FingerPrint}
	}

	if s.Status != t.Status {
		diff["Status"] = []interface{}{s.Status, t.Status}
	}

	if s.StorageName != t.StorageName {
		diff["StorageName"] = []interface{}{s.StorageName, t.StorageName}
	}

	if s.Subject != t.Subject {
		diff["Subject"] = []interface{}{s.Subject, t.Subject}
	}

	if !CheckSameNilAndLen(s.SubjectAlternativeNames, t.SubjectAlternativeNames, opt) {
		diff["SubjectAlternativeNames"] = []interface{}{s.SubjectAlternativeNames, t.SubjectAlternativeNames}
	} else {
		diff2 := make(map[string][]interface{})
		for i := range s.SubjectAlternativeNames {
			if s.SubjectAlternativeNames[i] != t.SubjectAlternativeNames[i] {
				diff2[strconv.Itoa(i)] = []interface{}{s.SubjectAlternativeNames[i], t.SubjectAlternativeNames[i]}
			}
		}
		if len(diff2) > 0 {
			diff["SubjectAlternativeNames"] = []interface{}{diff2}
		}
	}

	return diff
}
