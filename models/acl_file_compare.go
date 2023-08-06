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

// Equal checks if two structs of type ACLFile are equal
//
//	var a, b ACLFile
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s ACLFile) Equal(t ACLFile, opts ...Options) bool {
	if s.Description != t.Description {
		return false
	}

	if s.ID != t.ID {
		return false
	}

	if s.StorageName != t.StorageName {
		return false
	}

	return true
}

// Diff checks if two structs of type ACLFile are equal
//
//	var a, b ACLFile
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s ACLFile) Diff(t ACLFile, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if s.Description != t.Description {
		diff["Description"] = []interface{}{s.Description, t.Description}
	}

	if s.ID != t.ID {
		diff["ID"] = []interface{}{s.ID, t.ID}
	}

	if s.StorageName != t.StorageName {
		diff["StorageName"] = []interface{}{s.StorageName, t.StorageName}
	}

	return diff
}
