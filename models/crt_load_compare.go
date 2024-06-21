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

// Equal checks if two structs of type CrtLoad are equal
//
//	var a, b CrtLoad
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s CrtLoad) Equal(t CrtLoad, opts ...Options) bool {
	if s.Alias != t.Alias {
		return false
	}

	if s.Certificate != t.Certificate {
		return false
	}

	if s.Issuer != t.Issuer {
		return false
	}

	if s.Key != t.Key {
		return false
	}

	if s.Ocsp != t.Ocsp {
		return false
	}

	if s.OcspUpdate != t.OcspUpdate {
		return false
	}

	if s.Sctl != t.Sctl {
		return false
	}

	return true
}

// Diff checks if two structs of type CrtLoad are equal
//
//	var a, b CrtLoad
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s CrtLoad) Diff(t CrtLoad, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if s.Alias != t.Alias {
		diff["Alias"] = []interface{}{s.Alias, t.Alias}
	}

	if s.Certificate != t.Certificate {
		diff["Certificate"] = []interface{}{s.Certificate, t.Certificate}
	}

	if s.Issuer != t.Issuer {
		diff["Issuer"] = []interface{}{s.Issuer, t.Issuer}
	}

	if s.Key != t.Key {
		diff["Key"] = []interface{}{s.Key, t.Key}
	}

	if s.Ocsp != t.Ocsp {
		diff["Ocsp"] = []interface{}{s.Ocsp, t.Ocsp}
	}

	if s.OcspUpdate != t.OcspUpdate {
		diff["OcspUpdate"] = []interface{}{s.OcspUpdate, t.OcspUpdate}
	}

	if s.Sctl != t.Sctl {
		diff["Sctl"] = []interface{}{s.Sctl, t.Sctl}
	}

	return diff
}
