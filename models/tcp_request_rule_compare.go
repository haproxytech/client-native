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

// Equal checks if two structs of type TCPRequestRule are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b TCPRequestRule
//	equal := a.Equal(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b TCPRequestRule
//	equal := a.Equal(b,Options{
//		SkipIndex: true,
//	})
func (s TCPRequestRule) Equal(t TCPRequestRule, opts ...Options) bool {
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

	if s.Expr != t.Expr {
		return false
	}

	if s.GptValue != t.GptValue {
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

	if s.ResolveProtocol != t.ResolveProtocol {
		return false
	}

	if s.ResolveResolvers != t.ResolveResolvers {
		return false
	}

	if s.ResolveVar != t.ResolveVar {
		return false
	}

	if s.ScIdx != t.ScIdx {
		return false
	}

	if s.ScIncID != t.ScIncID {
		return false
	}

	if !equalPointers(s.ScInt, t.ScInt) {
		return false
	}

	if s.ServerName != t.ServerName {
		return false
	}

	if s.ServiceName != t.ServiceName {
		return false
	}

	if s.SpoeEngineName != t.SpoeEngineName {
		return false
	}

	if s.SpoeGroupName != t.SpoeGroupName {
		return false
	}

	if s.SwitchModeProto != t.SwitchModeProto {
		return false
	}

	if !equalPointers(s.Timeout, t.Timeout) {
		return false
	}

	if s.TosValue != t.TosValue {
		return false
	}

	if s.TrackKey != t.TrackKey {
		return false
	}

	if !equalPointers(s.TrackStickCounter, t.TrackStickCounter) {
		return false
	}

	if s.TrackTable != t.TrackTable {
		return false
	}

	if s.Type != t.Type {
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

// Diff checks if two structs of type TCPRequestRule are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b TCPRequestRule
//	diff := a.Diff(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b TCPRequestRule
//	diff := a.Diff(b,Options{
//		SkipIndex: true,
//	})
func (s TCPRequestRule) Diff(t TCPRequestRule, opts ...Options) map[string][]interface{} {
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

	if s.Expr != t.Expr {
		diff["Expr"] = []interface{}{s.Expr, t.Expr}
	}

	if s.GptValue != t.GptValue {
		diff["GptValue"] = []interface{}{s.GptValue, t.GptValue}
	}

	if !opt.SkipIndex && !equalPointers(s.Index, t.Index) {
		diff["Index"] = []interface{}{ValueOrNil(s.Index), ValueOrNil(t.Index)}
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

	if s.ResolveProtocol != t.ResolveProtocol {
		diff["ResolveProtocol"] = []interface{}{s.ResolveProtocol, t.ResolveProtocol}
	}

	if s.ResolveResolvers != t.ResolveResolvers {
		diff["ResolveResolvers"] = []interface{}{s.ResolveResolvers, t.ResolveResolvers}
	}

	if s.ResolveVar != t.ResolveVar {
		diff["ResolveVar"] = []interface{}{s.ResolveVar, t.ResolveVar}
	}

	if s.ScIdx != t.ScIdx {
		diff["ScIdx"] = []interface{}{s.ScIdx, t.ScIdx}
	}

	if s.ScIncID != t.ScIncID {
		diff["ScIncID"] = []interface{}{s.ScIncID, t.ScIncID}
	}

	if !equalPointers(s.ScInt, t.ScInt) {
		diff["ScInt"] = []interface{}{ValueOrNil(s.ScInt), ValueOrNil(t.ScInt)}
	}

	if s.ServerName != t.ServerName {
		diff["ServerName"] = []interface{}{s.ServerName, t.ServerName}
	}

	if s.ServiceName != t.ServiceName {
		diff["ServiceName"] = []interface{}{s.ServiceName, t.ServiceName}
	}

	if s.SpoeEngineName != t.SpoeEngineName {
		diff["SpoeEngineName"] = []interface{}{s.SpoeEngineName, t.SpoeEngineName}
	}

	if s.SpoeGroupName != t.SpoeGroupName {
		diff["SpoeGroupName"] = []interface{}{s.SpoeGroupName, t.SpoeGroupName}
	}

	if s.SwitchModeProto != t.SwitchModeProto {
		diff["SwitchModeProto"] = []interface{}{s.SwitchModeProto, t.SwitchModeProto}
	}

	if !equalPointers(s.Timeout, t.Timeout) {
		diff["Timeout"] = []interface{}{ValueOrNil(s.Timeout), ValueOrNil(t.Timeout)}
	}

	if s.TosValue != t.TosValue {
		diff["TosValue"] = []interface{}{s.TosValue, t.TosValue}
	}

	if s.TrackKey != t.TrackKey {
		diff["TrackKey"] = []interface{}{s.TrackKey, t.TrackKey}
	}

	if !equalPointers(s.TrackStickCounter, t.TrackStickCounter) {
		diff["TrackStickCounter"] = []interface{}{ValueOrNil(s.TrackStickCounter), ValueOrNil(t.TrackStickCounter)}
	}

	if s.TrackTable != t.TrackTable {
		diff["TrackTable"] = []interface{}{s.TrackTable, t.TrackTable}
	}

	if s.Type != t.Type {
		diff["Type"] = []interface{}{s.Type, t.Type}
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
