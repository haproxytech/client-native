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

// Equal checks if two structs of type HTTPErrorRule are equal
//
// By default empty maps and slices are equal to nil:
//  var a, b HTTPErrorRule
//  equal := a.Equal(b)
// For more advanced use case you can configure these options (default values are shown):
//  var a, b HTTPErrorRule
//  equal := a.Equal(b,Options{
//  	NilSameAsEmpty: true,

//		SkipIndex: true,
//	})
func (s HTTPErrorRule) Equal(t HTTPErrorRule, opts ...Options) bool {
	opt := getOptions(opts...)

	if !CheckSameNilAndLen(s.ReturnHeaders, t.ReturnHeaders, opt) {
		return false
	}
	for i := range s.ReturnHeaders {
		if !s.ReturnHeaders[i].Equal(*t.ReturnHeaders[i], opt) {
			return false
		}
	}

	if !opt.SkipIndex && !equalPointers(s.Index, t.Index) {
		return false
	}

	if s.ReturnContent != t.ReturnContent {
		return false
	}

	if s.ReturnContentFormat != t.ReturnContentFormat {
		return false
	}

	if !equalPointers(s.ReturnContentType, t.ReturnContentType) {
		return false
	}

	if s.Status != t.Status {
		return false
	}

	if s.Type != t.Type {
		return false
	}

	return true
}

// Diff checks if two structs of type HTTPErrorRule are equal
//
// By default empty arrays, maps and slices are equal to nil:
//  var a, b HTTPErrorRule
//  diff := a.Diff(b)
// For more advanced use case you can configure the options (default values are shown):
//  var a, b HTTPErrorRule
//  equal := a.Diff(b,Options{
//  	NilSameAsEmpty: true,

//		SkipIndex: true,
//	})
func (s HTTPErrorRule) Diff(t HTTPErrorRule, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})
	if !CheckSameNilAndLen(s.ReturnHeaders, t.ReturnHeaders, opt) {
		diff["ReturnHeaders"] = []interface{}{s.ReturnHeaders, t.ReturnHeaders}
	} else {
		diff2 := make(map[string][]interface{})
		for i := range s.ReturnHeaders {
			diffSub := s.ReturnHeaders[i].Diff(*t.ReturnHeaders[i], opt)
			if len(diffSub) > 0 {
				diff2[strconv.Itoa(i)] = []interface{}{diffSub}
			}
		}
		if len(diff2) > 0 {
			diff["ReturnHeaders"] = []interface{}{diff2}
		}
	}

	if !opt.SkipIndex && !equalPointers(s.Index, t.Index) {
		diff["Index"] = []interface{}{ValueOrNil(s.Index), ValueOrNil(t.Index)}
	}

	if s.ReturnContent != t.ReturnContent {
		diff["ReturnContent"] = []interface{}{s.ReturnContent, t.ReturnContent}
	}

	if s.ReturnContentFormat != t.ReturnContentFormat {
		diff["ReturnContentFormat"] = []interface{}{s.ReturnContentFormat, t.ReturnContentFormat}
	}

	if !equalPointers(s.ReturnContentType, t.ReturnContentType) {
		diff["ReturnContentType"] = []interface{}{ValueOrNil(s.ReturnContentType), ValueOrNil(t.ReturnContentType)}
	}

	if s.Status != t.Status {
		diff["Status"] = []interface{}{s.Status, t.Status}
	}

	if s.Type != t.Type {
		diff["Type"] = []interface{}{s.Type, t.Type}
	}

	return diff
}
