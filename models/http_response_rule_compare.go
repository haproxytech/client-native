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

// Equal checks if two structs of type HTTPResponseRule are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b HTTPResponseRule
//	equal := a.Equal(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b HTTPResponseRule
//	equal := a.Equal(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s HTTPResponseRule) Equal(t HTTPResponseRule, opts ...Options) bool {
	opt := getOptions(opts...)

	if !CheckSameNilAndLen(s.ReturnHeaders, t.ReturnHeaders, opt) {
		return false
	} else {
		for i := range s.ReturnHeaders {
			if !s.ReturnHeaders[i].Equal(*t.ReturnHeaders[i], opt) {
				return false
			}
		}
	}

	if s.ACLFile != t.ACLFile {
		return false
	}

	if s.ACLKeyfmt != t.ACLKeyfmt {
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

	if s.CacheName != t.CacheName {
		return false
	}

	if !equalPointers(s.CaptureID, t.CaptureID) {
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

	if !equalPointers(s.DenyStatus, t.DenyStatus) {
		return false
	}

	if s.Expr != t.Expr {
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

	if s.LuaAction != t.LuaAction {
		return false
	}

	if s.LuaParams != t.LuaParams {
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

	if s.MarkValue != t.MarkValue {
		return false
	}

	if s.NiceValue != t.NiceValue {
		return false
	}

	if !equalPointers(s.RedirCode, t.RedirCode) {
		return false
	}

	if s.RedirOption != t.RedirOption {
		return false
	}

	if s.RedirType != t.RedirType {
		return false
	}

	if s.RedirValue != t.RedirValue {
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

	if !equalPointers(s.ReturnStatusCode, t.ReturnStatusCode) {
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

	if s.Status != t.Status {
		return false
	}

	if s.StatusReason != t.StatusReason {
		return false
	}

	if s.StrictMode != t.StrictMode {
		return false
	}

	if s.Timeout != t.Timeout {
		return false
	}

	if s.TimeoutType != t.TimeoutType {
		return false
	}

	if s.TosValue != t.TosValue {
		return false
	}

	if s.TrackScKey != t.TrackScKey {
		return false
	}

	if !equalPointers(s.TrackScStickCounter, t.TrackScStickCounter) {
		return false
	}

	if s.TrackScTable != t.TrackScTable {
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

	if !equalPointers(s.WaitAtLeast, t.WaitAtLeast) {
		return false
	}

	if !equalPointers(s.WaitTime, t.WaitTime) {
		return false
	}

	return true
}

// Diff checks if two structs of type HTTPResponseRule are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b HTTPResponseRule
//	diff := a.Diff(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b HTTPResponseRule
//	diff := a.Diff(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s HTTPResponseRule) Diff(t HTTPResponseRule, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})
	if !CheckSameNilAndLen(s.ReturnHeaders, t.ReturnHeaders, opt) {
		diff["ReturnHeaders"] = []interface{}{s.ReturnHeaders, t.ReturnHeaders}
	} else {
		diff2 := make(map[string][]interface{})
		for i := range s.ReturnHeaders {
			if !s.ReturnHeaders[i].Equal(*t.ReturnHeaders[i], opt) {
				diffSub := s.ReturnHeaders[i].Diff(*t.ReturnHeaders[i], opt)
				if len(diffSub) > 0 {
					diff2[strconv.Itoa(i)] = []interface{}{diffSub}
				}
			}
		}
		if len(diff2) > 0 {
			diff["ReturnHeaders"] = []interface{}{diff2}
		}
	}

	if s.ACLFile != t.ACLFile {
		diff["ACLFile"] = []interface{}{s.ACLFile, t.ACLFile}
	}

	if s.ACLKeyfmt != t.ACLKeyfmt {
		diff["ACLKeyfmt"] = []interface{}{s.ACLKeyfmt, t.ACLKeyfmt}
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

	if s.CacheName != t.CacheName {
		diff["CacheName"] = []interface{}{s.CacheName, t.CacheName}
	}

	if !equalPointers(s.CaptureID, t.CaptureID) {
		diff["CaptureID"] = []interface{}{ValueOrNil(s.CaptureID), ValueOrNil(t.CaptureID)}
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

	if !equalPointers(s.DenyStatus, t.DenyStatus) {
		diff["DenyStatus"] = []interface{}{ValueOrNil(s.DenyStatus), ValueOrNil(t.DenyStatus)}
	}

	if s.Expr != t.Expr {
		diff["Expr"] = []interface{}{s.Expr, t.Expr}
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

	if s.LuaAction != t.LuaAction {
		diff["LuaAction"] = []interface{}{s.LuaAction, t.LuaAction}
	}

	if s.LuaParams != t.LuaParams {
		diff["LuaParams"] = []interface{}{s.LuaParams, t.LuaParams}
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

	if s.MarkValue != t.MarkValue {
		diff["MarkValue"] = []interface{}{s.MarkValue, t.MarkValue}
	}

	if s.NiceValue != t.NiceValue {
		diff["NiceValue"] = []interface{}{s.NiceValue, t.NiceValue}
	}

	if !equalPointers(s.RedirCode, t.RedirCode) {
		diff["RedirCode"] = []interface{}{ValueOrNil(s.RedirCode), ValueOrNil(t.RedirCode)}
	}

	if s.RedirOption != t.RedirOption {
		diff["RedirOption"] = []interface{}{s.RedirOption, t.RedirOption}
	}

	if s.RedirType != t.RedirType {
		diff["RedirType"] = []interface{}{s.RedirType, t.RedirType}
	}

	if s.RedirValue != t.RedirValue {
		diff["RedirValue"] = []interface{}{s.RedirValue, t.RedirValue}
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

	if !equalPointers(s.ReturnStatusCode, t.ReturnStatusCode) {
		diff["ReturnStatusCode"] = []interface{}{ValueOrNil(s.ReturnStatusCode), ValueOrNil(t.ReturnStatusCode)}
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

	if s.SpoeEngine != t.SpoeEngine {
		diff["SpoeEngine"] = []interface{}{s.SpoeEngine, t.SpoeEngine}
	}

	if s.SpoeGroup != t.SpoeGroup {
		diff["SpoeGroup"] = []interface{}{s.SpoeGroup, t.SpoeGroup}
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

	if s.Timeout != t.Timeout {
		diff["Timeout"] = []interface{}{s.Timeout, t.Timeout}
	}

	if s.TimeoutType != t.TimeoutType {
		diff["TimeoutType"] = []interface{}{s.TimeoutType, t.TimeoutType}
	}

	if s.TosValue != t.TosValue {
		diff["TosValue"] = []interface{}{s.TosValue, t.TosValue}
	}

	if s.TrackScKey != t.TrackScKey {
		diff["TrackScKey"] = []interface{}{s.TrackScKey, t.TrackScKey}
	}

	if !equalPointers(s.TrackScStickCounter, t.TrackScStickCounter) {
		diff["TrackScStickCounter"] = []interface{}{ValueOrNil(s.TrackScStickCounter), ValueOrNil(t.TrackScStickCounter)}
	}

	if s.TrackScTable != t.TrackScTable {
		diff["TrackScTable"] = []interface{}{s.TrackScTable, t.TrackScTable}
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

	if !equalPointers(s.WaitAtLeast, t.WaitAtLeast) {
		diff["WaitAtLeast"] = []interface{}{ValueOrNil(s.WaitAtLeast), ValueOrNil(t.WaitAtLeast)}
	}

	if !equalPointers(s.WaitTime, t.WaitTime) {
		diff["WaitTime"] = []interface{}{ValueOrNil(s.WaitTime), ValueOrNil(t.WaitTime)}
	}

	return diff
}
