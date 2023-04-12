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

// Equal checks if two structs of type SslOcspUpdate are equal
//
//	var a, b SslOcspUpdate
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s SslOcspUpdate) Equal(t SslOcspUpdate, opts ...Options) bool {
	if s.CertID != t.CertID {
		return false
	}

	if s.Failures != t.Failures {
		return false
	}

	if s.LastUpdate != t.LastUpdate {
		return false
	}

	if s.LastUpdateStatus != t.LastUpdateStatus {
		return false
	}

	if s.LastUpdateStatusStr != t.LastUpdateStatusStr {
		return false
	}

	if s.NextUpdate != t.NextUpdate {
		return false
	}

	if s.Path != t.Path {
		return false
	}

	if s.Successes != t.Successes {
		return false
	}

	return true
}

// Diff checks if two structs of type SslOcspUpdate are equal
//
//	var a, b SslOcspUpdate
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s SslOcspUpdate) Diff(t SslOcspUpdate, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if s.CertID != t.CertID {
		diff["CertID"] = []interface{}{s.CertID, t.CertID}
	}

	if s.Failures != t.Failures {
		diff["Failures"] = []interface{}{s.Failures, t.Failures}
	}

	if s.LastUpdate != t.LastUpdate {
		diff["LastUpdate"] = []interface{}{s.LastUpdate, t.LastUpdate}
	}

	if s.LastUpdateStatus != t.LastUpdateStatus {
		diff["LastUpdateStatus"] = []interface{}{s.LastUpdateStatus, t.LastUpdateStatus}
	}

	if s.LastUpdateStatusStr != t.LastUpdateStatusStr {
		diff["LastUpdateStatusStr"] = []interface{}{s.LastUpdateStatusStr, t.LastUpdateStatusStr}
	}

	if s.NextUpdate != t.NextUpdate {
		diff["NextUpdate"] = []interface{}{s.NextUpdate, t.NextUpdate}
	}

	if s.Path != t.Path {
		diff["Path"] = []interface{}{s.Path, t.Path}
	}

	if s.Successes != t.Successes {
		diff["Successes"] = []interface{}{s.Successes, t.Successes}
	}

	return diff
}
