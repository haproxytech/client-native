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

// Equal checks if two structs of type TuneBufferOptions are equal
//
//	var a, b TuneBufferOptions
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s TuneBufferOptions) Equal(t TuneBufferOptions, opts ...Options) bool {
	if !equalPointers(s.BuffersLimit, t.BuffersLimit) {
		return false
	}

	if s.BuffersReserve != t.BuffersReserve {
		return false
	}

	if s.Bufsize != t.Bufsize {
		return false
	}

	if s.Pipesize != t.Pipesize {
		return false
	}

	if !equalPointers(s.RcvbufBackend, t.RcvbufBackend) {
		return false
	}

	if !equalPointers(s.RcvbufClient, t.RcvbufClient) {
		return false
	}

	if !equalPointers(s.RcvbufFrontend, t.RcvbufFrontend) {
		return false
	}

	if !equalPointers(s.RcvbufServer, t.RcvbufServer) {
		return false
	}

	if s.RecvEnough != t.RecvEnough {
		return false
	}

	if !equalPointers(s.SndbufBackend, t.SndbufBackend) {
		return false
	}

	if !equalPointers(s.SndbufClient, t.SndbufClient) {
		return false
	}

	if !equalPointers(s.SndbufFrontend, t.SndbufFrontend) {
		return false
	}

	if !equalPointers(s.SndbufServer, t.SndbufServer) {
		return false
	}

	return true
}

// Diff checks if two structs of type TuneBufferOptions are equal
//
//	var a, b TuneBufferOptions
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s TuneBufferOptions) Diff(t TuneBufferOptions, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if !equalPointers(s.BuffersLimit, t.BuffersLimit) {
		diff["BuffersLimit"] = []interface{}{ValueOrNil(s.BuffersLimit), ValueOrNil(t.BuffersLimit)}
	}

	if s.BuffersReserve != t.BuffersReserve {
		diff["BuffersReserve"] = []interface{}{s.BuffersReserve, t.BuffersReserve}
	}

	if s.Bufsize != t.Bufsize {
		diff["Bufsize"] = []interface{}{s.Bufsize, t.Bufsize}
	}

	if s.Pipesize != t.Pipesize {
		diff["Pipesize"] = []interface{}{s.Pipesize, t.Pipesize}
	}

	if !equalPointers(s.RcvbufBackend, t.RcvbufBackend) {
		diff["RcvbufBackend"] = []interface{}{ValueOrNil(s.RcvbufBackend), ValueOrNil(t.RcvbufBackend)}
	}

	if !equalPointers(s.RcvbufClient, t.RcvbufClient) {
		diff["RcvbufClient"] = []interface{}{ValueOrNil(s.RcvbufClient), ValueOrNil(t.RcvbufClient)}
	}

	if !equalPointers(s.RcvbufFrontend, t.RcvbufFrontend) {
		diff["RcvbufFrontend"] = []interface{}{ValueOrNil(s.RcvbufFrontend), ValueOrNil(t.RcvbufFrontend)}
	}

	if !equalPointers(s.RcvbufServer, t.RcvbufServer) {
		diff["RcvbufServer"] = []interface{}{ValueOrNil(s.RcvbufServer), ValueOrNil(t.RcvbufServer)}
	}

	if s.RecvEnough != t.RecvEnough {
		diff["RecvEnough"] = []interface{}{s.RecvEnough, t.RecvEnough}
	}

	if !equalPointers(s.SndbufBackend, t.SndbufBackend) {
		diff["SndbufBackend"] = []interface{}{ValueOrNil(s.SndbufBackend), ValueOrNil(t.SndbufBackend)}
	}

	if !equalPointers(s.SndbufClient, t.SndbufClient) {
		diff["SndbufClient"] = []interface{}{ValueOrNil(s.SndbufClient), ValueOrNil(t.SndbufClient)}
	}

	if !equalPointers(s.SndbufFrontend, t.SndbufFrontend) {
		diff["SndbufFrontend"] = []interface{}{ValueOrNil(s.SndbufFrontend), ValueOrNil(t.SndbufFrontend)}
	}

	if !equalPointers(s.SndbufServer, t.SndbufServer) {
		diff["SndbufServer"] = []interface{}{ValueOrNil(s.SndbufServer), ValueOrNil(t.SndbufServer)}
	}

	return diff
}
