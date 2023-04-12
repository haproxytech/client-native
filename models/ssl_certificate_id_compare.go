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

// Equal checks if two structs of type SslCertificateID are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b SslCertificateID
//	equal := a.Equal(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b SslCertificateID
//	equal := a.Equal(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s SslCertificateID) Equal(t SslCertificateID, opts ...Options) bool {
	opt := getOptions(opts...)

	if s.CertificateID == nil || t.CertificateID == nil {
		if s.CertificateID != nil || t.CertificateID != nil {
			if opt.NilSameAsEmpty {
				empty := &CertificateID{}
				if s.CertificateID == nil {
					if !(t.CertificateID.Equal(*empty)) {
						return false
					}
				}
				if t.CertificateID == nil {
					if !(s.CertificateID.Equal(*empty)) {
						return false
					}
				}
			} else {
				return false
			}
		}
	} else if !s.CertificateID.Equal(*t.CertificateID, opt) {
		return false
	}

	if s.CertificateIDKey != t.CertificateIDKey {
		return false
	}

	if s.CertificatePath != t.CertificatePath {
		return false
	}

	return true
}

// Diff checks if two structs of type SslCertificateID are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b SslCertificateID
//	diff := a.Diff(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b SslCertificateID
//	diff := a.Diff(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s SslCertificateID) Diff(t SslCertificateID, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})

	if s.CertificateID == nil || t.CertificateID == nil {
		if s.CertificateID != nil || t.CertificateID != nil {
			if opt.NilSameAsEmpty {
				empty := &CertificateID{}
				if s.CertificateID == nil {
					if !(t.CertificateID.Equal(*empty)) {
						diff["CertificateID"] = []interface{}{ValueOrNil(s.CertificateID), ValueOrNil(t.CertificateID)}
					}
				}
				if t.CertificateID == nil {
					if !(s.CertificateID.Equal(*empty)) {
						diff["CertificateID"] = []interface{}{ValueOrNil(s.CertificateID), ValueOrNil(t.CertificateID)}
					}
				}
			} else {
				diff["CertificateID"] = []interface{}{ValueOrNil(s.CertificateID), ValueOrNil(t.CertificateID)}
			}
		}
	} else if !s.CertificateID.Equal(*t.CertificateID, opt) {
		diff["CertificateID"] = []interface{}{ValueOrNil(s.CertificateID), ValueOrNil(t.CertificateID)}
	}

	if s.CertificateIDKey != t.CertificateIDKey {
		diff["CertificateIDKey"] = []interface{}{s.CertificateIDKey, t.CertificateIDKey}
	}

	if s.CertificatePath != t.CertificatePath {
		diff["CertificatePath"] = []interface{}{s.CertificatePath, t.CertificatePath}
	}

	return diff
}
