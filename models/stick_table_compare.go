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

import (
	"strconv"
)

// Equal checks if two structs of type StickTable are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b StickTable
//	equal := a.Equal(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b StickTable
//	equal := a.Equal(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s StickTable) Equal(t StickTable, opts ...Options) bool {
	opt := getOptions(opts...)

	if !CheckSameNilAndLen(s.Fields, t.Fields, opt) {
		return false
	} else {
		for i := range s.Fields {
			if !s.Fields[i].Equal(*t.Fields[i], opt) {
				return false
			}
		}
	}

	if s.Name != t.Name {
		return false
	}

	if !equalPointers(s.Size, t.Size) {
		return false
	}

	if s.Type != t.Type {
		return false
	}

	if !equalPointers(s.Used, t.Used) {
		return false
	}

	return true
}

// Diff checks if two structs of type StickTable are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b StickTable
//	diff := a.Diff(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b StickTable
//	diff := a.Diff(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s StickTable) Diff(t StickTable, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})
	if !CheckSameNilAndLen(s.Fields, t.Fields, opt) {
		diff["Fields"] = []interface{}{s.Fields, t.Fields}
	} else {
		diff2 := make(map[string][]interface{})
		for i := range s.Fields {
			if !s.Fields[i].Equal(*t.Fields[i], opt) {
				diffSub := s.Fields[i].Diff(*t.Fields[i], opt)
				if len(diffSub) > 0 {
					diff2[strconv.Itoa(i)] = []interface{}{diffSub}
				}
			}
		}
		if len(diff2) > 0 {
			diff["Fields"] = []interface{}{diff2}
		}
	}

	if s.Name != t.Name {
		diff["Name"] = []interface{}{s.Name, t.Name}
	}

	if !equalPointers(s.Size, t.Size) {
		diff["Size"] = []interface{}{ValueOrNil(s.Size), ValueOrNil(t.Size)}
	}

	if s.Type != t.Type {
		diff["Type"] = []interface{}{s.Type, t.Type}
	}

	if !equalPointers(s.Used, t.Used) {
		diff["Used"] = []interface{}{ValueOrNil(s.Used), ValueOrNil(t.Used)}
	}

	return diff
}

// Equal checks if two structs of type StickTableField are equal
//
//	var a, b StickTableField
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s StickTableField) Equal(t StickTableField, opts ...Options) bool {
	if s.Field != t.Field {
		return false
	}

	if s.Idx != t.Idx {
		return false
	}

	if s.Period != t.Period {
		return false
	}

	if s.Type != t.Type {
		return false
	}

	return true
}

// Diff checks if two structs of type StickTableField are equal
//
//	var a, b StickTableField
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s StickTableField) Diff(t StickTableField, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if s.Field != t.Field {
		diff["Field"] = []interface{}{s.Field, t.Field}
	}

	if s.Idx != t.Idx {
		diff["Idx"] = []interface{}{s.Idx, t.Idx}
	}

	if s.Period != t.Period {
		diff["Period"] = []interface{}{s.Period, t.Period}
	}

	if s.Type != t.Type {
		diff["Type"] = []interface{}{s.Type, t.Type}
	}

	return diff
}
