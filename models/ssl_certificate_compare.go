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
	"github.com/go-openapi/strfmt"
)

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

	if s.Algorithm != t.Algorithm {
		return false
	}

	if s.AuthorityKeyID != t.AuthorityKeyID {
		return false
	}

	if s.ChainIssuer != t.ChainIssuer {
		return false
	}

	if s.ChainSubject != t.ChainSubject {
		return false
	}

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

	if s.NotAfter == nil || t.NotAfter == nil {
		if s.NotAfter != nil || t.NotAfter != nil {
			if opt.NilSameAsEmpty {
				empty := &strfmt.DateTime{}
				if s.NotAfter == nil {
					if !(t.NotAfter.Equal(*empty)) {
						return false
					}
				}
				if t.NotAfter == nil {
					if !(s.NotAfter.Equal(*empty)) {
						return false
					}
				}
			} else {
				return false
			}
		}
	} else if !s.NotAfter.Equal(*t.NotAfter) {
		return false
	}

	if s.NotBefore == nil || t.NotBefore == nil {
		if s.NotBefore != nil || t.NotBefore != nil {
			if opt.NilSameAsEmpty {
				empty := &strfmt.DateTime{}
				if s.NotBefore == nil {
					if !(t.NotBefore.Equal(*empty)) {
						return false
					}
				}
				if t.NotBefore == nil {
					if !(s.NotBefore.Equal(*empty)) {
						return false
					}
				}
			} else {
				return false
			}
		}
	} else if !s.NotBefore.Equal(*t.NotBefore) {
		return false
	}

	if s.Serial != t.Serial {
		return false
	}

	if s.Sha1FingerPrint != t.Sha1FingerPrint {
		return false
	}

	if s.Sha256FingerPrint != t.Sha256FingerPrint {
		return false
	}

	if !equalPointers(s.Size, t.Size) {
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

	if s.SubjectAlternativeNames != t.SubjectAlternativeNames {
		return false
	}

	if s.SubjectKeyID != t.SubjectKeyID {
		return false
	}

	return true
}

// Diff checks if two structs of type SslCertificate are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b SslCertificate
//	diff := a.Diff(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b SslCertificate
//	diff := a.Diff(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s SslCertificate) Diff(t SslCertificate, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})
	if s.Algorithm != t.Algorithm {
		diff["Algorithm"] = []interface{}{s.Algorithm, t.Algorithm}
	}

	if s.AuthorityKeyID != t.AuthorityKeyID {
		diff["AuthorityKeyID"] = []interface{}{s.AuthorityKeyID, t.AuthorityKeyID}
	}

	if s.ChainIssuer != t.ChainIssuer {
		diff["ChainIssuer"] = []interface{}{s.ChainIssuer, t.ChainIssuer}
	}

	if s.ChainSubject != t.ChainSubject {
		diff["ChainSubject"] = []interface{}{s.ChainSubject, t.ChainSubject}
	}

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

	if s.NotAfter == nil || t.NotAfter == nil {
		if s.NotAfter != nil || t.NotAfter != nil {
			if opt.NilSameAsEmpty {
				empty := &strfmt.DateTime{}
				if s.NotAfter == nil {
					if !(t.NotAfter.Equal(*empty)) {
						diff["NotAfter"] = []interface{}{ValueOrNil(s.NotAfter), ValueOrNil(t.NotAfter)}
					}
				}
				if t.NotAfter == nil {
					if !(s.NotAfter.Equal(*empty)) {
						diff["NotAfter"] = []interface{}{ValueOrNil(s.NotAfter), ValueOrNil(t.NotAfter)}
					}
				}
			} else {
				diff["NotAfter"] = []interface{}{ValueOrNil(s.NotAfter), ValueOrNil(t.NotAfter)}
			}
		}
	} else if !s.NotAfter.Equal(*t.NotAfter) {
		diff["NotAfter"] = []interface{}{ValueOrNil(s.NotAfter), ValueOrNil(t.NotAfter)}
	}

	if s.NotBefore == nil || t.NotBefore == nil {
		if s.NotBefore != nil || t.NotBefore != nil {
			if opt.NilSameAsEmpty {
				empty := &strfmt.DateTime{}
				if s.NotBefore == nil {
					if !(t.NotBefore.Equal(*empty)) {
						diff["NotBefore"] = []interface{}{ValueOrNil(s.NotBefore), ValueOrNil(t.NotBefore)}
					}
				}
				if t.NotBefore == nil {
					if !(s.NotBefore.Equal(*empty)) {
						diff["NotBefore"] = []interface{}{ValueOrNil(s.NotBefore), ValueOrNil(t.NotBefore)}
					}
				}
			} else {
				diff["NotBefore"] = []interface{}{ValueOrNil(s.NotBefore), ValueOrNil(t.NotBefore)}
			}
		}
	} else if !s.NotBefore.Equal(*t.NotBefore) {
		diff["NotBefore"] = []interface{}{ValueOrNil(s.NotBefore), ValueOrNil(t.NotBefore)}
	}

	if s.Serial != t.Serial {
		diff["Serial"] = []interface{}{s.Serial, t.Serial}
	}

	if s.Sha1FingerPrint != t.Sha1FingerPrint {
		diff["Sha1FingerPrint"] = []interface{}{s.Sha1FingerPrint, t.Sha1FingerPrint}
	}

	if s.Sha256FingerPrint != t.Sha256FingerPrint {
		diff["Sha256FingerPrint"] = []interface{}{s.Sha256FingerPrint, t.Sha256FingerPrint}
	}

	if !equalPointers(s.Size, t.Size) {
		diff["Size"] = []interface{}{ValueOrNil(s.Size), ValueOrNil(t.Size)}
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

	if s.SubjectAlternativeNames != t.SubjectAlternativeNames {
		diff["SubjectAlternativeNames"] = []interface{}{s.SubjectAlternativeNames, t.SubjectAlternativeNames}
	}

	if s.SubjectKeyID != t.SubjectKeyID {
		diff["SubjectKeyID"] = []interface{}{s.SubjectKeyID, t.SubjectKeyID}
	}

	return diff
}
