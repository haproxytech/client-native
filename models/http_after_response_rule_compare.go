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

// Equal checks if two structs of type HTTPAfterResponseRule are equal
//
//	var a, b HTTPAfterResponseRule
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s HTTPAfterResponseRule) Equal(t HTTPAfterResponseRule, opts ...Options) bool {
	if s.ACLFile != t.ACLFile {
		return false
	}

	if s.ACLKeyfmt != t.ACLKeyfmt {
		return false
	}

	if !equalPointers(s.CaptureID, t.CaptureID) {
		return false
	}

	if s.CaptureLen != t.CaptureLen {
		return false
	}

	if s.CaptureSample != t.CaptureSample {
		return false
	}

	if s.Cond != t.Cond {
		return false
	}

	if s.CondTest != t.CondTest {
		return false
	}

	if s.HdrFormat != t.HdrFormat {
		return false
	}

	if s.HdrMatch != t.HdrMatch {
		return false
	}

	if s.HdrMethod != t.HdrMethod {
		return false
	}

	if s.HdrName != t.HdrName {
		return false
	}

	if s.LogLevel != t.LogLevel {
		return false
	}

	if s.MapFile != t.MapFile {
		return false
	}

	if s.MapKeyfmt != t.MapKeyfmt {
		return false
	}

	if s.MapValuefmt != t.MapValuefmt {
		return false
	}

	if s.ScExpr != t.ScExpr {
		return false
	}

	if s.ScID != t.ScID {
		return false
	}

	if s.ScIdx != t.ScIdx {
		return false
	}

	if !equalPointers(s.ScInt, t.ScInt) {
		return false
	}

	if s.Status != t.Status {
		return false
	}

	if s.StatusReason != t.StatusReason {
		return false
	}

	if s.StrictMode != t.StrictMode {
		return false
	}

	if s.Type != t.Type {
		return false
	}

	if s.VarExpr != t.VarExpr {
		return false
	}

	if s.VarFormat != t.VarFormat {
		return false
	}

	if s.VarName != t.VarName {
		return false
	}

	if s.VarScope != t.VarScope {
		return false
	}

	return true
}

// Diff checks if two structs of type HTTPAfterResponseRule are equal
//
//	var a, b HTTPAfterResponseRule
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s HTTPAfterResponseRule) Diff(t HTTPAfterResponseRule, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if s.ACLFile != t.ACLFile {
		diff["ACLFile"] = []interface{}{s.ACLFile, t.ACLFile}
	}

	if s.ACLKeyfmt != t.ACLKeyfmt {
		diff["ACLKeyfmt"] = []interface{}{s.ACLKeyfmt, t.ACLKeyfmt}
	}

	if !equalPointers(s.CaptureID, t.CaptureID) {
		diff["CaptureID"] = []interface{}{ValueOrNil(s.CaptureID), ValueOrNil(t.CaptureID)}
	}

	if s.CaptureLen != t.CaptureLen {
		diff["CaptureLen"] = []interface{}{s.CaptureLen, t.CaptureLen}
	}

	if s.CaptureSample != t.CaptureSample {
		diff["CaptureSample"] = []interface{}{s.CaptureSample, t.CaptureSample}
	}

	if s.Cond != t.Cond {
		diff["Cond"] = []interface{}{s.Cond, t.Cond}
	}

	if s.CondTest != t.CondTest {
		diff["CondTest"] = []interface{}{s.CondTest, t.CondTest}
	}

	if s.HdrFormat != t.HdrFormat {
		diff["HdrFormat"] = []interface{}{s.HdrFormat, t.HdrFormat}
	}

	if s.HdrMatch != t.HdrMatch {
		diff["HdrMatch"] = []interface{}{s.HdrMatch, t.HdrMatch}
	}

	if s.HdrMethod != t.HdrMethod {
		diff["HdrMethod"] = []interface{}{s.HdrMethod, t.HdrMethod}
	}

	if s.HdrName != t.HdrName {
		diff["HdrName"] = []interface{}{s.HdrName, t.HdrName}
	}

	if s.LogLevel != t.LogLevel {
		diff["LogLevel"] = []interface{}{s.LogLevel, t.LogLevel}
	}

	if s.MapFile != t.MapFile {
		diff["MapFile"] = []interface{}{s.MapFile, t.MapFile}
	}

	if s.MapKeyfmt != t.MapKeyfmt {
		diff["MapKeyfmt"] = []interface{}{s.MapKeyfmt, t.MapKeyfmt}
	}

	if s.MapValuefmt != t.MapValuefmt {
		diff["MapValuefmt"] = []interface{}{s.MapValuefmt, t.MapValuefmt}
	}

	if s.ScExpr != t.ScExpr {
		diff["ScExpr"] = []interface{}{s.ScExpr, t.ScExpr}
	}

	if s.ScID != t.ScID {
		diff["ScID"] = []interface{}{s.ScID, t.ScID}
	}

	if s.ScIdx != t.ScIdx {
		diff["ScIdx"] = []interface{}{s.ScIdx, t.ScIdx}
	}

	if !equalPointers(s.ScInt, t.ScInt) {
		diff["ScInt"] = []interface{}{ValueOrNil(s.ScInt), ValueOrNil(t.ScInt)}
	}

	if s.Status != t.Status {
		diff["Status"] = []interface{}{s.Status, t.Status}
	}

	if s.StatusReason != t.StatusReason {
		diff["StatusReason"] = []interface{}{s.StatusReason, t.StatusReason}
	}

	if s.StrictMode != t.StrictMode {
		diff["StrictMode"] = []interface{}{s.StrictMode, t.StrictMode}
	}

	if s.Type != t.Type {
		diff["Type"] = []interface{}{s.Type, t.Type}
	}

	if s.VarExpr != t.VarExpr {
		diff["VarExpr"] = []interface{}{s.VarExpr, t.VarExpr}
	}

	if s.VarFormat != t.VarFormat {
		diff["VarFormat"] = []interface{}{s.VarFormat, t.VarFormat}
	}

	if s.VarName != t.VarName {
		diff["VarName"] = []interface{}{s.VarName, t.VarName}
	}

	if s.VarScope != t.VarScope {
		diff["VarScope"] = []interface{}{s.VarScope, t.VarScope}
	}

	return diff
}
