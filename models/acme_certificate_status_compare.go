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

// Equal checks if two structs of type AcmeCertificateStatus are equal
//
//	var a, b AcmeCertificateStatus
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s AcmeCertificateStatus) Equal(t AcmeCertificateStatus, opts ...Options) bool {
	if s.AcmeSection != t.AcmeSection {
		return false
	}

	if s.Certificate != t.Certificate {
		return false
	}

	if s.ExpiriesIn != t.ExpiriesIn {
		return false
	}

	if !s.ExpiryDate.Equal(t.ExpiryDate) {
		return false
	}

	if s.RenewalIn != t.RenewalIn {
		return false
	}

	if !s.ScheduledRenewal.Equal(t.ScheduledRenewal) {
		return false
	}

	if s.State != t.State {
		return false
	}

	return true
}

// Diff checks if two structs of type AcmeCertificateStatus are equal
//
//	var a, b AcmeCertificateStatus
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s AcmeCertificateStatus) Diff(t AcmeCertificateStatus, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if s.AcmeSection != t.AcmeSection {
		diff["AcmeSection"] = []interface{}{s.AcmeSection, t.AcmeSection}
	}

	if s.Certificate != t.Certificate {
		diff["Certificate"] = []interface{}{s.Certificate, t.Certificate}
	}

	if s.ExpiriesIn != t.ExpiriesIn {
		diff["ExpiriesIn"] = []interface{}{s.ExpiriesIn, t.ExpiriesIn}
	}

	if !s.ExpiryDate.Equal(t.ExpiryDate) {
		diff["ExpiryDate"] = []interface{}{s.ExpiryDate, t.ExpiryDate}
	}

	if s.RenewalIn != t.RenewalIn {
		diff["RenewalIn"] = []interface{}{s.RenewalIn, t.RenewalIn}
	}

	if !s.ScheduledRenewal.Equal(t.ScheduledRenewal) {
		diff["ScheduledRenewal"] = []interface{}{s.ScheduledRenewal, t.ScheduledRenewal}
	}

	if s.State != t.State {
		diff["State"] = []interface{}{s.State, t.State}
	}

	return diff
}
