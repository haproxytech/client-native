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

// Equal checks if two structs of type MysqlCheckParams are equal
//
//	var a, b MysqlCheckParams
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s MysqlCheckParams) Equal(t MysqlCheckParams, opts ...Options) bool {
	if s.ClientVersion != t.ClientVersion {
		return false
	}

	if s.Username != t.Username {
		return false
	}

	return true
}

// Diff checks if two structs of type MysqlCheckParams are equal
//
//	var a, b MysqlCheckParams
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s MysqlCheckParams) Diff(t MysqlCheckParams, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if s.ClientVersion != t.ClientVersion {
		diff["ClientVersion"] = []interface{}{s.ClientVersion, t.ClientVersion}
	}

	if s.Username != t.Username {
		diff["Username"] = []interface{}{s.Username, t.Username}
	}

	return diff
}
