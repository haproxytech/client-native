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

// Equal checks if two structs of type ConfigStickTable are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b ConfigStickTable
//	equal := a.Equal(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b ConfigStickTable
//	equal := a.Equal(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s ConfigStickTable) Equal(t ConfigStickTable, opts ...Options) bool {
	opt := getOptions(opts...)

	if !equalPointers(s.Expire, t.Expire) {
		return false
	}

	if !equalPointers(s.Keylen, t.Keylen) {
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

	if s.Nopurge != t.Nopurge {
		return false
	}

	if s.Peers != t.Peers {
		return false
	}

	if !equalPointers(s.Size, t.Size) {
		return false
	}

	if !equalPointers(s.Srvkey, t.Srvkey) {
		return false
	}

	if s.Store != t.Store {
		return false
	}

	if s.Type != t.Type {
		return false
	}

	if !equalPointers(s.WriteTo, t.WriteTo) {
		return false
	}

	return true
}

// Diff checks if two structs of type ConfigStickTable are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b ConfigStickTable
//	diff := a.Diff(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b ConfigStickTable
//	diff := a.Diff(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s ConfigStickTable) Diff(t ConfigStickTable, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})
	if !equalPointers(s.Expire, t.Expire) {
		diff["Expire"] = []interface{}{ValueOrNil(s.Expire), ValueOrNil(t.Expire)}
	}

	if !equalPointers(s.Keylen, t.Keylen) {
		diff["Keylen"] = []interface{}{ValueOrNil(s.Keylen), ValueOrNil(t.Keylen)}
	}

	if !CheckSameNilAndLenMap[string](s.Metadata, t.Metadata, opt) {
		diff["Metadata"] = []interface{}{s.Metadata, t.Metadata}
	}

	for k, v := range s.Metadata {
		if !reflect.DeepEqual(t.Metadata[k], v) {
			diff["Metadata"] = []interface{}{s.Metadata, t.Metadata}
		}
	}

	if s.Nopurge != t.Nopurge {
		diff["Nopurge"] = []interface{}{s.Nopurge, t.Nopurge}
	}

	if s.Peers != t.Peers {
		diff["Peers"] = []interface{}{s.Peers, t.Peers}
	}

	if !equalPointers(s.Size, t.Size) {
		diff["Size"] = []interface{}{ValueOrNil(s.Size), ValueOrNil(t.Size)}
	}

	if !equalPointers(s.Srvkey, t.Srvkey) {
		diff["Srvkey"] = []interface{}{ValueOrNil(s.Srvkey), ValueOrNil(t.Srvkey)}
	}

	if s.Store != t.Store {
		diff["Store"] = []interface{}{s.Store, t.Store}
	}

	if s.Type != t.Type {
		diff["Type"] = []interface{}{s.Type, t.Type}
	}

	if !equalPointers(s.WriteTo, t.WriteTo) {
		diff["WriteTo"] = []interface{}{ValueOrNil(s.WriteTo), ValueOrNil(t.WriteTo)}
	}

	return diff
}
