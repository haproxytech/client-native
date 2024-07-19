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

// Equal checks if two structs of type Userlist are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b Userlist
//	equal := a.Equal(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b Userlist
//	equal := a.Equal(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s Userlist) Equal(t Userlist, opts ...Options) bool {
	opt := getOptions(opts...)

	if !s.UserlistBase.Equal(t.UserlistBase, opt) {
		return false
	}

	if !CheckSameNilAndLenMap[string, Group](s.Groups, t.Groups, opt) {
		return false
	}

	for k, v := range s.Groups {
		if !t.Groups[k].Equal(v, opt) {
			return false
		}
	}

	if !CheckSameNilAndLenMap[string, User](s.Users, t.Users, opt) {
		return false
	}

	for k, v := range s.Users {
		if !t.Users[k].Equal(v, opt) {
			return false
		}
	}

	return true
}

// Diff checks if two structs of type Userlist are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b Userlist
//	diff := a.Diff(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b Userlist
//	diff := a.Diff(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s Userlist) Diff(t Userlist, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})

	if !s.UserlistBase.Equal(t.UserlistBase, opt) {
		diff["UserlistBase"] = []interface{}{s.UserlistBase, t.UserlistBase}
	}

	if !CheckSameNilAndLenMap[string, Group](s.Groups, t.Groups, opt) {
		diff["Groups"] = []interface{}{s.Groups, t.Groups}
	}

	for k, v := range s.Groups {
		if !t.Groups[k].Equal(v, opt) {
			diff["Groups"] = []interface{}{s.Groups, t.Groups}
		}
	}

	if !CheckSameNilAndLenMap[string, User](s.Users, t.Users, opt) {
		diff["Users"] = []interface{}{s.Users, t.Users}
	}

	for k, v := range s.Users {
		if !t.Users[k].Equal(v, opt) {
			diff["Users"] = []interface{}{s.Users, t.Users}
		}
	}

	return diff
}
