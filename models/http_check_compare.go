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

// Equal checks if two structs of type HTTPCheck are equal
//
// By default empty maps and slices are equal to nil:
//  var a, b HTTPCheck
//  equal := a.Equal(b)
// For more advanced use case you can configure these options (default values are shown):
//  var a, b HTTPCheck
//  equal := a.Equal(b,Options{
//  	NilSameAsEmpty: true,

//		SkipIndex: true,
//	})
func (s HTTPCheck) Equal(t HTTPCheck, opts ...Options) bool {
	opt := getOptions(opts...)

	if !CheckSameNilAndLen(s.CheckHeaders, t.CheckHeaders, opt) {
		return false
	}
	for i := range s.CheckHeaders {
		if !s.CheckHeaders[i].Equal(*t.CheckHeaders[i], opt) {
			return false
		}
	}

	if s.Addr != t.Addr {
		return false
	}

	if s.Alpn != t.Alpn {
		return false
	}

	if s.Body != t.Body {
		return false
	}

	if s.BodyLogFormat != t.BodyLogFormat {
		return false
	}

	if s.CheckComment != t.CheckComment {
		return false
	}

	if s.Default != t.Default {
		return false
	}

	if s.ErrorStatus != t.ErrorStatus {
		return false
	}

	if s.ExclamationMark != t.ExclamationMark {
		return false
	}

	if !opt.SkipIndex && !equalPointers(s.Index, t.Index) {
		return false
	}

	if s.Linger != t.Linger {
		return false
	}

	if s.Match != t.Match {
		return false
	}

	if s.Method != t.Method {
		return false
	}

	if !equalPointers(s.MinRecv, t.MinRecv) {
		return false
	}

	if s.OkStatus != t.OkStatus {
		return false
	}

	if s.OnError != t.OnError {
		return false
	}

	if s.OnSuccess != t.OnSuccess {
		return false
	}

	if s.Pattern != t.Pattern {
		return false
	}

	if !equalPointers(s.Port, t.Port) {
		return false
	}

	if s.PortString != t.PortString {
		return false
	}

	if s.Proto != t.Proto {
		return false
	}

	if s.SendProxy != t.SendProxy {
		return false
	}

	if s.Sni != t.Sni {
		return false
	}

	if s.Ssl != t.Ssl {
		return false
	}

	if s.StatusCode != t.StatusCode {
		return false
	}

	if s.ToutStatus != t.ToutStatus {
		return false
	}

	if s.Type != t.Type {
		return false
	}

	if s.URI != t.URI {
		return false
	}

	if s.URILogFormat != t.URILogFormat {
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

	if s.Version != t.Version {
		return false
	}

	if s.ViaSocks4 != t.ViaSocks4 {
		return false
	}

	return true
}

// Diff checks if two structs of type HTTPCheck are equal
//
// By default empty arrays, maps and slices are equal to nil:
//  var a, b HTTPCheck
//  diff := a.Diff(b)
// For more advanced use case you can configure the options (default values are shown):
//  var a, b HTTPCheck
//  equal := a.Diff(b,Options{
//  	NilSameAsEmpty: true,

//		SkipIndex: true,
//	})
func (s HTTPCheck) Diff(t HTTPCheck, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})
	if len(s.CheckHeaders) != len(t.CheckHeaders) {
		diff["CheckHeaders"] = []interface{}{s.CheckHeaders, t.CheckHeaders}
	} else {
		diff2 := make(map[string][]interface{})
		for i := range s.CheckHeaders {
			diffSub := s.CheckHeaders[i].Diff(*t.CheckHeaders[i], opt)
			if len(diffSub) > 0 {
				diff2[strconv.Itoa(i)] = []interface{}{diffSub}
			}
		}
		if len(diff2) > 0 {
			diff["CheckHeaders"] = []interface{}{diff2}
		}
	}

	if s.Addr != t.Addr {
		diff["Addr"] = []interface{}{s.Addr, t.Addr}
	}

	if s.Alpn != t.Alpn {
		diff["Alpn"] = []interface{}{s.Alpn, t.Alpn}
	}

	if s.Body != t.Body {
		diff["Body"] = []interface{}{s.Body, t.Body}
	}

	if s.BodyLogFormat != t.BodyLogFormat {
		diff["BodyLogFormat"] = []interface{}{s.BodyLogFormat, t.BodyLogFormat}
	}

	if s.CheckComment != t.CheckComment {
		diff["CheckComment"] = []interface{}{s.CheckComment, t.CheckComment}
	}

	if s.Default != t.Default {
		diff["Default"] = []interface{}{s.Default, t.Default}
	}

	if s.ErrorStatus != t.ErrorStatus {
		diff["ErrorStatus"] = []interface{}{s.ErrorStatus, t.ErrorStatus}
	}

	if s.ExclamationMark != t.ExclamationMark {
		diff["ExclamationMark"] = []interface{}{s.ExclamationMark, t.ExclamationMark}
	}

	if !opt.SkipIndex && !equalPointers(s.Index, t.Index) {
		diff["Index"] = []interface{}{s.Index, t.Index}
	}

	if s.Linger != t.Linger {
		diff["Linger"] = []interface{}{s.Linger, t.Linger}
	}

	if s.Match != t.Match {
		diff["Match"] = []interface{}{s.Match, t.Match}
	}

	if s.Method != t.Method {
		diff["Method"] = []interface{}{s.Method, t.Method}
	}

	if !equalPointers(s.MinRecv, t.MinRecv) {
		diff["MinRecv"] = []interface{}{s.MinRecv, t.MinRecv}
	}

	if s.OkStatus != t.OkStatus {
		diff["OkStatus"] = []interface{}{s.OkStatus, t.OkStatus}
	}

	if s.OnError != t.OnError {
		diff["OnError"] = []interface{}{s.OnError, t.OnError}
	}

	if s.OnSuccess != t.OnSuccess {
		diff["OnSuccess"] = []interface{}{s.OnSuccess, t.OnSuccess}
	}

	if s.Pattern != t.Pattern {
		diff["Pattern"] = []interface{}{s.Pattern, t.Pattern}
	}

	if !equalPointers(s.Port, t.Port) {
		diff["Port"] = []interface{}{s.Port, t.Port}
	}

	if s.PortString != t.PortString {
		diff["PortString"] = []interface{}{s.PortString, t.PortString}
	}

	if s.Proto != t.Proto {
		diff["Proto"] = []interface{}{s.Proto, t.Proto}
	}

	if s.SendProxy != t.SendProxy {
		diff["SendProxy"] = []interface{}{s.SendProxy, t.SendProxy}
	}

	if s.Sni != t.Sni {
		diff["Sni"] = []interface{}{s.Sni, t.Sni}
	}

	if s.Ssl != t.Ssl {
		diff["Ssl"] = []interface{}{s.Ssl, t.Ssl}
	}

	if s.StatusCode != t.StatusCode {
		diff["StatusCode"] = []interface{}{s.StatusCode, t.StatusCode}
	}

	if s.ToutStatus != t.ToutStatus {
		diff["ToutStatus"] = []interface{}{s.ToutStatus, t.ToutStatus}
	}

	if s.Type != t.Type {
		diff["Type"] = []interface{}{s.Type, t.Type}
	}

	if s.URI != t.URI {
		diff["URI"] = []interface{}{s.URI, t.URI}
	}

	if s.URILogFormat != t.URILogFormat {
		diff["URILogFormat"] = []interface{}{s.URILogFormat, t.URILogFormat}
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

	if s.Version != t.Version {
		diff["Version"] = []interface{}{s.Version, t.Version}
	}

	if s.ViaSocks4 != t.ViaSocks4 {
		diff["ViaSocks4"] = []interface{}{s.ViaSocks4, t.ViaSocks4}
	}

	return diff
}
