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

// Equal checks if two structs of type SslCrl are equal
//
//	var a, b SslCrl
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s SslCrl) Equal(t SslCrl, opts ...Options) bool {
	if s.Description != t.Description {
		return false
	}

	if s.File != t.File {
		return false
	}

	if s.StorageName != t.StorageName {
		return false
	}

	return true
}

// Diff checks if two structs of type SslCrl are equal
//
//	var a, b SslCrl
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s SslCrl) Diff(t SslCrl, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if s.Description != t.Description {
		diff["Description"] = []interface{}{s.Description, t.Description}
	}

	if s.File != t.File {
		diff["File"] = []interface{}{s.File, t.File}
	}

	if s.StorageName != t.StorageName {
		diff["StorageName"] = []interface{}{s.StorageName, t.StorageName}
	}

	return diff
}
