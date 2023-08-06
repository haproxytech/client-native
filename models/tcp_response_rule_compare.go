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

// Equal checks if two structs of type TCPResponseRule are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b TCPResponseRule
//	equal := a.Equal(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b TCPResponseRule
//	equal := a.Equal(b,Options{
//		SkipIndex: true,
//	})
func (s TCPResponseRule) Equal(t TCPResponseRule, opts ...Options) bool {
	opt := getOptions(opts...)

	if s.Action != t.Action {
		return false
	}

	if s.BandwidthLimitLimit != t.BandwidthLimitLimit {
		return false
	}

	if s.BandwidthLimitName != t.BandwidthLimitName {
		return false
	}

	if s.BandwidthLimitPeriod != t.BandwidthLimitPeriod {
		return false
	}

	if s.Cond != t.Cond {
		return false
	}

	if s.CondTest != t.CondTest {
		return false
	}

	if s.Expr != t.Expr {
		return false
	}

	if !opt.SkipIndex && !equalPointers(s.Index, t.Index) {
		return false
	}

	if s.LogLevel != t.LogLevel {
		return false
	}

	if s.LuaAction != t.LuaAction {
		return false
	}

	if s.LuaParams != t.LuaParams {
		return false
	}

	if s.MarkValue != t.MarkValue {
		return false
	}

	if s.NiceValue != t.NiceValue {
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

	if s.SpoeEngine != t.SpoeEngine {
		return false
	}

	if s.SpoeGroup != t.SpoeGroup {
		return false
	}

	if !equalPointers(s.Timeout, t.Timeout) {
		return false
	}

	if s.TosValue != t.TosValue {
		return false
	}

	if s.Type != t.Type {
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

// Diff checks if two structs of type TCPResponseRule are equal
//
//	var a, b TCPResponseRule
//	diff := a.Diff(b)
//
// For more advanced use case you can configure the options (default values are shown):
//
//	var a, b TCPResponseRule
//	equal := a.Diff(b,Options{
//		SkipIndex: true,
//	})
func (s TCPResponseRule) Diff(t TCPResponseRule, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})
	if s.Action != t.Action {
		diff["Action"] = []interface{}{s.Action, t.Action}
	}

	if s.BandwidthLimitLimit != t.BandwidthLimitLimit {
		diff["BandwidthLimitLimit"] = []interface{}{s.BandwidthLimitLimit, t.BandwidthLimitLimit}
	}

	if s.BandwidthLimitName != t.BandwidthLimitName {
		diff["BandwidthLimitName"] = []interface{}{s.BandwidthLimitName, t.BandwidthLimitName}
	}

	if s.BandwidthLimitPeriod != t.BandwidthLimitPeriod {
		diff["BandwidthLimitPeriod"] = []interface{}{s.BandwidthLimitPeriod, t.BandwidthLimitPeriod}
	}

	if s.Cond != t.Cond {
		diff["Cond"] = []interface{}{s.Cond, t.Cond}
	}

	if s.CondTest != t.CondTest {
		diff["CondTest"] = []interface{}{s.CondTest, t.CondTest}
	}

	if s.Expr != t.Expr {
		diff["Expr"] = []interface{}{s.Expr, t.Expr}
	}

	if !opt.SkipIndex && !equalPointers(s.Index, t.Index) {
		diff["Index"] = []interface{}{s.Index, t.Index}
	}

	if s.LogLevel != t.LogLevel {
		diff["LogLevel"] = []interface{}{s.LogLevel, t.LogLevel}
	}

	if s.LuaAction != t.LuaAction {
		diff["LuaAction"] = []interface{}{s.LuaAction, t.LuaAction}
	}

	if s.LuaParams != t.LuaParams {
		diff["LuaParams"] = []interface{}{s.LuaParams, t.LuaParams}
	}

	if s.MarkValue != t.MarkValue {
		diff["MarkValue"] = []interface{}{s.MarkValue, t.MarkValue}
	}

	if s.NiceValue != t.NiceValue {
		diff["NiceValue"] = []interface{}{s.NiceValue, t.NiceValue}
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
		diff["ScInt"] = []interface{}{s.ScInt, t.ScInt}
	}

	if s.SpoeEngine != t.SpoeEngine {
		diff["SpoeEngine"] = []interface{}{s.SpoeEngine, t.SpoeEngine}
	}

	if s.SpoeGroup != t.SpoeGroup {
		diff["SpoeGroup"] = []interface{}{s.SpoeGroup, t.SpoeGroup}
	}

	if !equalPointers(s.Timeout, t.Timeout) {
		diff["Timeout"] = []interface{}{s.Timeout, t.Timeout}
	}

	if s.TosValue != t.TosValue {
		diff["TosValue"] = []interface{}{s.TosValue, t.TosValue}
	}

	if s.Type != t.Type {
		diff["Type"] = []interface{}{s.Type, t.Type}
	}

	if s.VarName != t.VarName {
		diff["VarName"] = []interface{}{s.VarName, t.VarName}
	}

	if s.VarScope != t.VarScope {
		diff["VarScope"] = []interface{}{s.VarScope, t.VarScope}
	}

	return diff
}
