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

// Equal checks if two structs of type Compression are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b Compression
//	equal := a.Equal(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b Compression
//	equal := a.Equal(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s Compression) Equal(t Compression, opts ...Options) bool {
	opt := getOptions(opts...)

	if s.AlgoReq != t.AlgoReq {
		return false
	}

	if !equalComparableSlice(s.Algorithms, t.Algorithms, opt) {
		return false
	}

	if !equalComparableSlice(s.AlgosRes, t.AlgosRes, opt) {
		return false
	}

	if s.Direction != t.Direction {
		return false
	}

	if s.MinsizeReq != t.MinsizeReq {
		return false
	}

	if s.MinsizeRes != t.MinsizeRes {
		return false
	}

	if s.Offload != t.Offload {
		return false
	}

	if !equalComparableSlice(s.Types, t.Types, opt) {
		return false
	}

	if !equalComparableSlice(s.TypesReq, t.TypesReq, opt) {
		return false
	}

	if !equalComparableSlice(s.TypesRes, t.TypesRes, opt) {
		return false
	}

	return true
}

// Diff checks if two structs of type Compression are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b Compression
//	diff := a.Diff(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b Compression
//	diff := a.Diff(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s Compression) Diff(t Compression, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})
	if s.AlgoReq != t.AlgoReq {
		diff["AlgoReq"] = []interface{}{s.AlgoReq, t.AlgoReq}
	}

	if !equalComparableSlice(s.Algorithms, t.Algorithms, opt) {
		diff["Algorithms"] = []interface{}{s.Algorithms, t.Algorithms}
	}

	if !equalComparableSlice(s.AlgosRes, t.AlgosRes, opt) {
		diff["AlgosRes"] = []interface{}{s.AlgosRes, t.AlgosRes}
	}

	if s.Direction != t.Direction {
		diff["Direction"] = []interface{}{s.Direction, t.Direction}
	}

	if s.MinsizeReq != t.MinsizeReq {
		diff["MinsizeReq"] = []interface{}{s.MinsizeReq, t.MinsizeReq}
	}

	if s.MinsizeRes != t.MinsizeRes {
		diff["MinsizeRes"] = []interface{}{s.MinsizeRes, t.MinsizeRes}
	}

	if s.Offload != t.Offload {
		diff["Offload"] = []interface{}{s.Offload, t.Offload}
	}

	if !equalComparableSlice(s.Types, t.Types, opt) {
		diff["Types"] = []interface{}{s.Types, t.Types}
	}

	if !equalComparableSlice(s.TypesReq, t.TypesReq, opt) {
		diff["TypesReq"] = []interface{}{s.TypesReq, t.TypesReq}
	}

	if !equalComparableSlice(s.TypesRes, t.TypesRes, opt) {
		diff["TypesRes"] = []interface{}{s.TypesRes, t.TypesRes}
	}

	return diff
}
