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

// Equal checks if two structs of type SslOcspResponse are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b SslOcspResponse
//	equal := a.Equal(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b SslOcspResponse
//	equal := a.Equal(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s SslOcspResponse) Equal(t SslOcspResponse, opts ...Options) bool {
	opt := getOptions(opts...)

	if s.Base64Response != t.Base64Response {
		return false
	}

	if s.OcspResponseStatus != t.OcspResponseStatus {
		return false
	}

	if !s.ProducedAt.Equal(t.ProducedAt) {
		return false
	}

	if !equalComparableSlice(s.ResponderID, t.ResponderID, opt) {
		return false
	}

	if s.ResponseType != t.ResponseType {
		return false
	}

	if s.Responses == nil || t.Responses == nil {
		if s.Responses != nil || t.Responses != nil {
			if opt.NilSameAsEmpty {
				empty := &OCSPResponses{}
				if s.Responses == nil {
					if !(t.Responses.Equal(*empty)) {
						return false
					}
				}
				if t.Responses == nil {
					if !(s.Responses.Equal(*empty)) {
						return false
					}
				}
			} else {
				return false
			}
		}
	} else if !s.Responses.Equal(*t.Responses, opt) {
		return false
	}

	if s.Version != t.Version {
		return false
	}

	return true
}

// Diff checks if two structs of type SslOcspResponse are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b SslOcspResponse
//	diff := a.Diff(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b SslOcspResponse
//	diff := a.Diff(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s SslOcspResponse) Diff(t SslOcspResponse, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})
	if s.Base64Response != t.Base64Response {
		diff["Base64Response"] = []interface{}{s.Base64Response, t.Base64Response}
	}

	if s.OcspResponseStatus != t.OcspResponseStatus {
		diff["OcspResponseStatus"] = []interface{}{s.OcspResponseStatus, t.OcspResponseStatus}
	}

	if !s.ProducedAt.Equal(t.ProducedAt) {
		diff["ProducedAt"] = []interface{}{s.ProducedAt, t.ProducedAt}
	}

	if !equalComparableSlice(s.ResponderID, t.ResponderID, opt) {
		diff["ResponderID"] = []interface{}{s.ResponderID, t.ResponderID}
	}

	if s.ResponseType != t.ResponseType {
		diff["ResponseType"] = []interface{}{s.ResponseType, t.ResponseType}
	}

	if s.Responses == nil || t.Responses == nil {
		if s.Responses != nil || t.Responses != nil {
			if opt.NilSameAsEmpty {
				empty := &OCSPResponses{}
				if s.Responses == nil {
					if !(t.Responses.Equal(*empty)) {
						diff["Responses"] = []interface{}{ValueOrNil(s.Responses), ValueOrNil(t.Responses)}
					}
				}
				if t.Responses == nil {
					if !(s.Responses.Equal(*empty)) {
						diff["Responses"] = []interface{}{ValueOrNil(s.Responses), ValueOrNil(t.Responses)}
					}
				}
			} else {
				diff["Responses"] = []interface{}{ValueOrNil(s.Responses), ValueOrNil(t.Responses)}
			}
		}
	} else if !s.Responses.Equal(*t.Responses, opt) {
		diff["Responses"] = []interface{}{ValueOrNil(s.Responses), ValueOrNil(t.Responses)}
	}

	if s.Version != t.Version {
		diff["Version"] = []interface{}{s.Version, t.Version}
	}

	return diff
}

// Equal checks if two structs of type OCSPResponses are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b OCSPResponses
//	equal := a.Equal(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b OCSPResponses
//	equal := a.Equal(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s OCSPResponses) Equal(t OCSPResponses, opts ...Options) bool {
	opt := getOptions(opts...)

	if s.CertStatus != t.CertStatus {
		return false
	}

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

	if !s.NextUpdate.Equal(t.NextUpdate) {
		return false
	}

	if s.RevocationReason != t.RevocationReason {
		return false
	}

	if !s.ThisUpdate.Equal(t.ThisUpdate) {
		return false
	}

	return true
}

// Diff checks if two structs of type OCSPResponses are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b OCSPResponses
//	diff := a.Diff(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b OCSPResponses
//	diff := a.Diff(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s OCSPResponses) Diff(t OCSPResponses, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})
	if s.CertStatus != t.CertStatus {
		diff["CertStatus"] = []interface{}{s.CertStatus, t.CertStatus}
	}

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

	if !s.NextUpdate.Equal(t.NextUpdate) {
		diff["NextUpdate"] = []interface{}{s.NextUpdate, t.NextUpdate}
	}

	if s.RevocationReason != t.RevocationReason {
		diff["RevocationReason"] = []interface{}{s.RevocationReason, t.RevocationReason}
	}

	if !s.ThisUpdate.Equal(t.ThisUpdate) {
		diff["ThisUpdate"] = []interface{}{s.ThisUpdate, t.ThisUpdate}
	}

	return diff
}
