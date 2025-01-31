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

// Equal checks if two structs of type LogProfileStep are equal
//
//	var a, b LogProfileStep
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s LogProfileStep) Equal(t LogProfileStep, opts ...Options) bool {
	if s.Drop != t.Drop {
		return false
	}

	if s.Format != t.Format {
		return false
	}

	if s.Sd != t.Sd {
		return false
	}

	if s.Step != t.Step {
		return false
	}

	return true
}

// Diff checks if two structs of type LogProfileStep are equal
//
//	var a, b LogProfileStep
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s LogProfileStep) Diff(t LogProfileStep, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if s.Drop != t.Drop {
		diff["Drop"] = []interface{}{s.Drop, t.Drop}
	}

	if s.Format != t.Format {
		diff["Format"] = []interface{}{s.Format, t.Format}
	}

	if s.Sd != t.Sd {
		diff["Sd"] = []interface{}{s.Sd, t.Sd}
	}

	if s.Step != t.Step {
		diff["Step"] = []interface{}{s.Step, t.Step}
	}

	return diff
}
