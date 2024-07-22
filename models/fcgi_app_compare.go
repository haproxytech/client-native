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

// Equal checks if two structs of type FCGIApp are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b FCGIApp
//	equal := a.Equal(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b FCGIApp
//	equal := a.Equal(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s FCGIApp) Equal(t FCGIApp, opts ...Options) bool {
	opt := getOptions(opts...)

	if !s.FCGIAppBase.Equal(t.FCGIAppBase, opt) {
		return false
	}

	if !s.ACLList.Equal(t.ACLList, opt) {
		return false
	}

	return true
}

// Diff checks if two structs of type FCGIApp are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b FCGIApp
//	diff := a.Diff(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b FCGIApp
//	diff := a.Diff(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s FCGIApp) Diff(t FCGIApp, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})

	if !s.FCGIAppBase.Equal(t.FCGIAppBase, opt) {
		diff["FCGIAppBase"] = []interface{}{s.FCGIAppBase, t.FCGIAppBase}
	}

	if !s.ACLList.Equal(t.ACLList, opt) {
		diff["ACLList"] = []interface{}{s.ACLList, t.ACLList}
	}

	return diff
}
