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

// Equal checks if two structs of type TuneVarsOptions are equal
//
//	var a, b TuneVarsOptions
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s TuneVarsOptions) Equal(t TuneVarsOptions, opts ...Options) bool {
	if !equalPointers(s.GlobalMaxSize, t.GlobalMaxSize) {
		return false
	}

	if !equalPointers(s.ProcMaxSize, t.ProcMaxSize) {
		return false
	}

	if !equalPointers(s.ReqresMaxSize, t.ReqresMaxSize) {
		return false
	}

	if !equalPointers(s.SessMaxSize, t.SessMaxSize) {
		return false
	}

	if !equalPointers(s.TxnMaxSize, t.TxnMaxSize) {
		return false
	}

	return true
}

// Diff checks if two structs of type TuneVarsOptions are equal
//
//	var a, b TuneVarsOptions
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s TuneVarsOptions) Diff(t TuneVarsOptions, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if !equalPointers(s.GlobalMaxSize, t.GlobalMaxSize) {
		diff["GlobalMaxSize"] = []interface{}{ValueOrNil(s.GlobalMaxSize), ValueOrNil(t.GlobalMaxSize)}
	}

	if !equalPointers(s.ProcMaxSize, t.ProcMaxSize) {
		diff["ProcMaxSize"] = []interface{}{ValueOrNil(s.ProcMaxSize), ValueOrNil(t.ProcMaxSize)}
	}

	if !equalPointers(s.ReqresMaxSize, t.ReqresMaxSize) {
		diff["ReqresMaxSize"] = []interface{}{ValueOrNil(s.ReqresMaxSize), ValueOrNil(t.ReqresMaxSize)}
	}

	if !equalPointers(s.SessMaxSize, t.SessMaxSize) {
		diff["SessMaxSize"] = []interface{}{ValueOrNil(s.SessMaxSize), ValueOrNil(t.SessMaxSize)}
	}

	if !equalPointers(s.TxnMaxSize, t.TxnMaxSize) {
		diff["TxnMaxSize"] = []interface{}{ValueOrNil(s.TxnMaxSize), ValueOrNil(t.TxnMaxSize)}
	}

	return diff
}
