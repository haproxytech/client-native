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

// Equal checks if two structs of type HttpchkParams are equal
//
//	var a, b HttpchkParams
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s HttpchkParams) Equal(t HttpchkParams, opts ...Options) bool {
	if s.Method != t.Method {
		return false
	}

	if s.URI != t.URI {
		return false
	}

	if s.Version != t.Version {
		return false
	}

	return true
}

// Diff checks if two structs of type HttpchkParams are equal
//
//	var a, b HttpchkParams
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s HttpchkParams) Diff(t HttpchkParams, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if s.Method != t.Method {
		diff["Method"] = []interface{}{s.Method, t.Method}
	}

	if s.URI != t.URI {
		diff["URI"] = []interface{}{s.URI, t.URI}
	}

	if s.Version != t.Version {
		diff["Version"] = []interface{}{s.Version, t.Version}
	}

	return diff
}
