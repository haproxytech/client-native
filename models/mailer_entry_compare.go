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

import "reflect"

// Equal checks if two structs of type MailerEntry are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b MailerEntry
//	equal := a.Equal(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b MailerEntry
//	equal := a.Equal(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s MailerEntry) Equal(t MailerEntry, opts ...Options) bool {
	opt := getOptions(opts...)

	if s.Address != t.Address {
		return false
	}

	if !CheckSameNilAndLenMap[string](s.Metadata, t.Metadata, opt) {
		return false
	}

	for k, v := range s.Metadata {
		if !reflect.DeepEqual(t.Metadata[k], v) {
			return false
		}
	}

	if s.Name != t.Name {
		return false
	}

	if s.Port != t.Port {
		return false
	}

	return true
}

// Diff checks if two structs of type MailerEntry are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b MailerEntry
//	diff := a.Diff(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b MailerEntry
//	diff := a.Diff(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s MailerEntry) Diff(t MailerEntry, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})
	if s.Address != t.Address {
		diff["Address"] = []interface{}{s.Address, t.Address}
	}

	if !CheckSameNilAndLenMap[string](s.Metadata, t.Metadata, opt) {
		diff["Metadata"] = []interface{}{s.Metadata, t.Metadata}
	}

	for k, v := range s.Metadata {
		if !reflect.DeepEqual(t.Metadata[k], v) {
			diff["Metadata"] = []interface{}{s.Metadata, t.Metadata}
		}
	}

	if s.Name != t.Name {
		diff["Name"] = []interface{}{s.Name, t.Name}
	}

	if s.Port != t.Port {
		diff["Port"] = []interface{}{s.Port, t.Port}
	}

	return diff
}
