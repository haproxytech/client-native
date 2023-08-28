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

// Equal checks if two structs of type Cookie are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b Cookie
//	equal := a.Equal(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b Cookie
//	equal := a.Equal(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s Cookie) Equal(t Cookie, opts ...Options) bool {
	opt := getOptions(opts...)

	if !CheckSameNilAndLen(s.Attrs, t.Attrs, opt) {
		return false
	}
	for i := range s.Attrs {
		if !s.Attrs[i].Equal(*t.Attrs[i], opt) {
			return false
		}
	}

	if !CheckSameNilAndLen(s.Domains, t.Domains, opt) {
		return false
	}
	for i := range s.Domains {
		if !s.Domains[i].Equal(*t.Domains[i], opt) {
			return false
		}
	}

	if s.Dynamic != t.Dynamic {
		return false
	}

	if s.Httponly != t.Httponly {
		return false
	}

	if s.Indirect != t.Indirect {
		return false
	}

	if s.Maxidle != t.Maxidle {
		return false
	}

	if s.Maxlife != t.Maxlife {
		return false
	}

	if !equalPointers(s.Name, t.Name) {
		return false
	}

	if s.Nocache != t.Nocache {
		return false
	}

	if s.Postonly != t.Postonly {
		return false
	}

	if s.Preserve != t.Preserve {
		return false
	}

	if s.Secure != t.Secure {
		return false
	}

	if s.Type != t.Type {
		return false
	}

	return true
}

// Diff checks if two structs of type Cookie are equal
//
// By default empty arrays, maps and slices are equal to nil:
//
//	var a, b Cookie
//	diff := a.Diff(b)
//
// For more advanced use case you can configure the options (default values are shown):
//
//	var a, b Cookie
//	equal := a.Diff(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s Cookie) Diff(t Cookie, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})
	if !CheckSameNilAndLen(s.Attrs, t.Attrs, opt) {
		diff["Attrs"] = []interface{}{s.Attrs, t.Attrs}
	} else {
		diff2 := make(map[string][]interface{})
		for i := range s.Attrs {
			diffSub := s.Attrs[i].Diff(*t.Attrs[i], opt)
			if len(diffSub) > 0 {
				diff2[strconv.Itoa(i)] = []interface{}{diffSub}
			}
		}
		if len(diff2) > 0 {
			diff["Attrs"] = []interface{}{diff2}
		}
	}

	if !CheckSameNilAndLen(s.Domains, t.Domains, opt) {
		diff["Domains"] = []interface{}{s.Domains, t.Domains}
	} else {
		diff2 := make(map[string][]interface{})
		for i := range s.Domains {
			diffSub := s.Domains[i].Diff(*t.Domains[i], opt)
			if len(diffSub) > 0 {
				diff2[strconv.Itoa(i)] = []interface{}{diffSub}
			}
		}
		if len(diff2) > 0 {
			diff["Domains"] = []interface{}{diff2}
		}
	}

	if s.Dynamic != t.Dynamic {
		diff["Dynamic"] = []interface{}{s.Dynamic, t.Dynamic}
	}

	if s.Httponly != t.Httponly {
		diff["Httponly"] = []interface{}{s.Httponly, t.Httponly}
	}

	if s.Indirect != t.Indirect {
		diff["Indirect"] = []interface{}{s.Indirect, t.Indirect}
	}

	if s.Maxidle != t.Maxidle {
		diff["Maxidle"] = []interface{}{s.Maxidle, t.Maxidle}
	}

	if s.Maxlife != t.Maxlife {
		diff["Maxlife"] = []interface{}{s.Maxlife, t.Maxlife}
	}

	if !equalPointers(s.Name, t.Name) {
		diff["Name"] = []interface{}{ValueOrNil(s.Name), ValueOrNil(t.Name)}
	}

	if s.Nocache != t.Nocache {
		diff["Nocache"] = []interface{}{s.Nocache, t.Nocache}
	}

	if s.Postonly != t.Postonly {
		diff["Postonly"] = []interface{}{s.Postonly, t.Postonly}
	}

	if s.Preserve != t.Preserve {
		diff["Preserve"] = []interface{}{s.Preserve, t.Preserve}
	}

	if s.Secure != t.Secure {
		diff["Secure"] = []interface{}{s.Secure, t.Secure}
	}

	if s.Type != t.Type {
		diff["Type"] = []interface{}{s.Type, t.Type}
	}

	return diff
}

// Equal checks if two structs of type Attr are equal
//
//	var a, b Attr
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s Attr) Equal(t Attr, opts ...Options) bool {
	if s.Value != t.Value {
		return false
	}

	return true
}

// Diff checks if two structs of type Attr are equal
//
//	var a, b Attr
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s Attr) Diff(t Attr, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if s.Value != t.Value {
		diff["Value"] = []interface{}{s.Value, t.Value}
	}

	return diff
}

// Equal checks if two structs of type Domain are equal
//
//	var a, b Domain
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s Domain) Equal(t Domain, opts ...Options) bool {
	if s.Value != t.Value {
		return false
	}

	return true
}

// Diff checks if two structs of type Domain are equal
//
//	var a, b Domain
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s Domain) Diff(t Domain, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if s.Value != t.Value {
		diff["Value"] = []interface{}{s.Value, t.Value}
	}

	return diff
}
