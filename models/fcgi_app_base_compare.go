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
	"reflect"
	"strconv"
)

// Equal checks if two structs of type FCGIAppBase are equal
//
// By default empty maps and slices are equal to nil:
//  var a, b FCGIAppBase
//  equal := a.Equal(b)
// For more advanced use case you can configure these options (default values are shown):
//  var a, b FCGIAppBase
//  equal := a.Equal(b,Options{
//  	NilSameAsEmpty: true,

//		SkipIndex: true,
//	})
func (s FCGIAppBase) Equal(t FCGIAppBase, opts ...Options) bool {
	opt := getOptions(opts...)

	if !equalPointers(s.Docroot, t.Docroot) {
		return false
	}

	if s.GetValues != t.GetValues {
		return false
	}

	if s.Index != t.Index {
		return false
	}

	if s.KeepConn != t.KeepConn {
		return false
	}

	if !CheckSameNilAndLen(s.LogStderrs, t.LogStderrs, opt) {
		return false
	} else {
		for i := range s.LogStderrs {
			if !s.LogStderrs[i].Equal(*t.LogStderrs[i], opt) {
				return false
			}
		}
	}

	if s.MaxReqs != t.MaxReqs {
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

	if s.MpxsConns != t.MpxsConns {
		return false
	}

	if s.Name != t.Name {
		return false
	}

	if !CheckSameNilAndLen(s.PassHeaders, t.PassHeaders, opt) {
		return false
	} else {
		for i := range s.PassHeaders {
			if !s.PassHeaders[i].Equal(*t.PassHeaders[i], opt) {
				return false
			}
		}
	}

	if s.PathInfo != t.PathInfo {
		return false
	}

	if !CheckSameNilAndLen(s.SetParams, t.SetParams, opt) {
		return false
	} else {
		for i := range s.SetParams {
			if !s.SetParams[i].Equal(*t.SetParams[i], opt) {
				return false
			}
		}
	}

	return true
}

// Diff checks if two structs of type FCGIAppBase are equal
//
// By default empty maps and slices are equal to nil:
//  var a, b FCGIAppBase
//  diff := a.Diff(b)
// For more advanced use case you can configure these options (default values are shown):
//  var a, b FCGIAppBase
//  diff := a.Diff(b,Options{
//  	NilSameAsEmpty: true,

//		SkipIndex: true,
//	})
func (s FCGIAppBase) Diff(t FCGIAppBase, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})
	if !equalPointers(s.Docroot, t.Docroot) {
		diff["Docroot"] = []interface{}{ValueOrNil(s.Docroot), ValueOrNil(t.Docroot)}
	}

	if s.GetValues != t.GetValues {
		diff["GetValues"] = []interface{}{s.GetValues, t.GetValues}
	}

	if s.Index != t.Index {
		diff["Index"] = []interface{}{s.Index, t.Index}
	}

	if s.KeepConn != t.KeepConn {
		diff["KeepConn"] = []interface{}{s.KeepConn, t.KeepConn}
	}

	if !CheckSameNilAndLen(s.LogStderrs, t.LogStderrs, opt) {
		diff["LogStderrs"] = []interface{}{s.LogStderrs, t.LogStderrs}
	} else {
		diff2 := make(map[string][]interface{})
		for i := range s.LogStderrs {
			if !s.LogStderrs[i].Equal(*t.LogStderrs[i], opt) {
				diffSub := s.LogStderrs[i].Diff(*t.LogStderrs[i], opt)
				if len(diffSub) > 0 {
					diff2[strconv.Itoa(i)] = []interface{}{diffSub}
				}
			}
		}
		if len(diff2) > 0 {
			diff["LogStderrs"] = []interface{}{diff2}
		}
	}

	if s.MaxReqs != t.MaxReqs {
		diff["MaxReqs"] = []interface{}{s.MaxReqs, t.MaxReqs}
	}

	if !CheckSameNilAndLenMap[string](s.Metadata, t.Metadata, opt) {
		diff["Metadata"] = []interface{}{s.Metadata, t.Metadata}
	}

	for k, v := range s.Metadata {
		if !reflect.DeepEqual(t.Metadata[k], v) {
			diff["Metadata"] = []interface{}{s.Metadata, t.Metadata}
		}
	}

	if s.MpxsConns != t.MpxsConns {
		diff["MpxsConns"] = []interface{}{s.MpxsConns, t.MpxsConns}
	}

	if s.Name != t.Name {
		diff["Name"] = []interface{}{s.Name, t.Name}
	}

	if !CheckSameNilAndLen(s.PassHeaders, t.PassHeaders, opt) {
		diff["PassHeaders"] = []interface{}{s.PassHeaders, t.PassHeaders}
	} else {
		diff2 := make(map[string][]interface{})
		for i := range s.PassHeaders {
			if !s.PassHeaders[i].Equal(*t.PassHeaders[i], opt) {
				diffSub := s.PassHeaders[i].Diff(*t.PassHeaders[i], opt)
				if len(diffSub) > 0 {
					diff2[strconv.Itoa(i)] = []interface{}{diffSub}
				}
			}
		}
		if len(diff2) > 0 {
			diff["PassHeaders"] = []interface{}{diff2}
		}
	}

	if s.PathInfo != t.PathInfo {
		diff["PathInfo"] = []interface{}{s.PathInfo, t.PathInfo}
	}

	if !CheckSameNilAndLen(s.SetParams, t.SetParams, opt) {
		diff["SetParams"] = []interface{}{s.SetParams, t.SetParams}
	} else {
		diff2 := make(map[string][]interface{})
		for i := range s.SetParams {
			if !s.SetParams[i].Equal(*t.SetParams[i], opt) {
				diffSub := s.SetParams[i].Diff(*t.SetParams[i], opt)
				if len(diffSub) > 0 {
					diff2[strconv.Itoa(i)] = []interface{}{diffSub}
				}
			}
		}
		if len(diff2) > 0 {
			diff["SetParams"] = []interface{}{diff2}
		}
	}

	return diff
}
