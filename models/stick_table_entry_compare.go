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

// Equal checks if two structs of type StickTableEntry are equal
//
//	var a, b StickTableEntry
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s StickTableEntry) Equal(t StickTableEntry, opts ...Options) bool {
	if !equalPointers(s.BytesInCnt, t.BytesInCnt) {
		return false
	}

	if !equalPointers(s.BytesInRate, t.BytesInRate) {
		return false
	}

	if !equalPointers(s.BytesOutCnt, t.BytesOutCnt) {
		return false
	}

	if !equalPointers(s.BytesOutRate, t.BytesOutRate) {
		return false
	}

	if !equalPointers(s.ConnCnt, t.ConnCnt) {
		return false
	}

	if !equalPointers(s.ConnCur, t.ConnCur) {
		return false
	}

	if !equalPointers(s.ConnRate, t.ConnRate) {
		return false
	}

	if !equalPointers(s.Exp, t.Exp) {
		return false
	}

	if !equalPointers(s.Gpc0, t.Gpc0) {
		return false
	}

	if !equalPointers(s.Gpc0Rate, t.Gpc0Rate) {
		return false
	}

	if !equalPointers(s.Gpc1, t.Gpc1) {
		return false
	}

	if !equalPointers(s.Gpc1Rate, t.Gpc1Rate) {
		return false
	}

	if !equalPointers(s.Gpt0, t.Gpt0) {
		return false
	}

	if !equalPointers(s.HTTPErrCnt, t.HTTPErrCnt) {
		return false
	}

	if !equalPointers(s.HTTPErrRate, t.HTTPErrRate) {
		return false
	}

	if !equalPointers(s.HTTPReqCnt, t.HTTPReqCnt) {
		return false
	}

	if !equalPointers(s.HTTPReqRate, t.HTTPReqRate) {
		return false
	}

	if s.ID != t.ID {
		return false
	}

	if s.Key != t.Key {
		return false
	}

	if !equalPointers(s.ServerID, t.ServerID) {
		return false
	}

	if !equalPointers(s.SessCnt, t.SessCnt) {
		return false
	}

	if !equalPointers(s.SessRate, t.SessRate) {
		return false
	}

	if s.Use != t.Use {
		return false
	}

	return true
}

// Diff checks if two structs of type StickTableEntry are equal
//
//	var a, b StickTableEntry
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s StickTableEntry) Diff(t StickTableEntry, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if !equalPointers(s.BytesInCnt, t.BytesInCnt) {
		diff["BytesInCnt"] = []interface{}{s.BytesInCnt, t.BytesInCnt}
	}

	if !equalPointers(s.BytesInRate, t.BytesInRate) {
		diff["BytesInRate"] = []interface{}{s.BytesInRate, t.BytesInRate}
	}

	if !equalPointers(s.BytesOutCnt, t.BytesOutCnt) {
		diff["BytesOutCnt"] = []interface{}{s.BytesOutCnt, t.BytesOutCnt}
	}

	if !equalPointers(s.BytesOutRate, t.BytesOutRate) {
		diff["BytesOutRate"] = []interface{}{s.BytesOutRate, t.BytesOutRate}
	}

	if !equalPointers(s.ConnCnt, t.ConnCnt) {
		diff["ConnCnt"] = []interface{}{s.ConnCnt, t.ConnCnt}
	}

	if !equalPointers(s.ConnCur, t.ConnCur) {
		diff["ConnCur"] = []interface{}{s.ConnCur, t.ConnCur}
	}

	if !equalPointers(s.ConnRate, t.ConnRate) {
		diff["ConnRate"] = []interface{}{s.ConnRate, t.ConnRate}
	}

	if !equalPointers(s.Exp, t.Exp) {
		diff["Exp"] = []interface{}{s.Exp, t.Exp}
	}

	if !equalPointers(s.Gpc0, t.Gpc0) {
		diff["Gpc0"] = []interface{}{s.Gpc0, t.Gpc0}
	}

	if !equalPointers(s.Gpc0Rate, t.Gpc0Rate) {
		diff["Gpc0Rate"] = []interface{}{s.Gpc0Rate, t.Gpc0Rate}
	}

	if !equalPointers(s.Gpc1, t.Gpc1) {
		diff["Gpc1"] = []interface{}{s.Gpc1, t.Gpc1}
	}

	if !equalPointers(s.Gpc1Rate, t.Gpc1Rate) {
		diff["Gpc1Rate"] = []interface{}{s.Gpc1Rate, t.Gpc1Rate}
	}

	if !equalPointers(s.Gpt0, t.Gpt0) {
		diff["Gpt0"] = []interface{}{s.Gpt0, t.Gpt0}
	}

	if !equalPointers(s.HTTPErrCnt, t.HTTPErrCnt) {
		diff["HTTPErrCnt"] = []interface{}{s.HTTPErrCnt, t.HTTPErrCnt}
	}

	if !equalPointers(s.HTTPErrRate, t.HTTPErrRate) {
		diff["HTTPErrRate"] = []interface{}{s.HTTPErrRate, t.HTTPErrRate}
	}

	if !equalPointers(s.HTTPReqCnt, t.HTTPReqCnt) {
		diff["HTTPReqCnt"] = []interface{}{s.HTTPReqCnt, t.HTTPReqCnt}
	}

	if !equalPointers(s.HTTPReqRate, t.HTTPReqRate) {
		diff["HTTPReqRate"] = []interface{}{s.HTTPReqRate, t.HTTPReqRate}
	}

	if s.ID != t.ID {
		diff["ID"] = []interface{}{s.ID, t.ID}
	}

	if s.Key != t.Key {
		diff["Key"] = []interface{}{s.Key, t.Key}
	}

	if !equalPointers(s.ServerID, t.ServerID) {
		diff["ServerID"] = []interface{}{s.ServerID, t.ServerID}
	}

	if !equalPointers(s.SessCnt, t.SessCnt) {
		diff["SessCnt"] = []interface{}{s.SessCnt, t.SessCnt}
	}

	if !equalPointers(s.SessRate, t.SessRate) {
		diff["SessRate"] = []interface{}{s.SessRate, t.SessRate}
	}

	if s.Use != t.Use {
		diff["Use"] = []interface{}{s.Use, t.Use}
	}

	return diff
}
