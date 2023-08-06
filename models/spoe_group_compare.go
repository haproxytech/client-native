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

// Equal checks if two structs of type SpoeGroup are equal
//
//	var a, b SpoeGroup
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s SpoeGroup) Equal(t SpoeGroup, opts ...Options) bool {
	if s.Messages != t.Messages {
		return false
	}

	if !equalPointers(s.Name, t.Name) {
		return false
	}

	return true
}

// Diff checks if two structs of type SpoeGroup are equal
//
//	var a, b SpoeGroup
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s SpoeGroup) Diff(t SpoeGroup, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if s.Messages != t.Messages {
		diff["Messages"] = []interface{}{s.Messages, t.Messages}
	}

	if !equalPointers(s.Name, t.Name) {
		diff["Name"] = []interface{}{s.Name, t.Name}
	}

	return diff
}
