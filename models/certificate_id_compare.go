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

// Equal checks if two structs of type CertificateID are equal
//
//	var a, b CertificateID
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s CertificateID) Equal(t CertificateID, opts ...Options) bool {
	if s.HashAlgorithm != t.HashAlgorithm {
		return false
	}

	if s.IssuerKeyHash != t.IssuerKeyHash {
		return false
	}

	if s.IssuerNameHash != t.IssuerNameHash {
		return false
	}

	if s.SerialNumber != t.SerialNumber {
		return false
	}

	return true
}

// Diff checks if two structs of type CertificateID are equal
//
//	var a, b CertificateID
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s CertificateID) Diff(t CertificateID, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if s.HashAlgorithm != t.HashAlgorithm {
		diff["HashAlgorithm"] = []interface{}{s.HashAlgorithm, t.HashAlgorithm}
	}

	if s.IssuerKeyHash != t.IssuerKeyHash {
		diff["IssuerKeyHash"] = []interface{}{s.IssuerKeyHash, t.IssuerKeyHash}
	}

	if s.IssuerNameHash != t.IssuerNameHash {
		diff["IssuerNameHash"] = []interface{}{s.IssuerNameHash, t.IssuerNameHash}
	}

	if s.SerialNumber != t.SerialNumber {
		diff["SerialNumber"] = []interface{}{s.SerialNumber, t.SerialNumber}
	}

	return diff
}
