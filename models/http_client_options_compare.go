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

// Equal checks if two structs of type HTTPClientOptions are equal
//
//	var a, b HTTPClientOptions
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s HTTPClientOptions) Equal(t HTTPClientOptions, opts ...Options) bool {
	if s.ResolversDisabled != t.ResolversDisabled {
		return false
	}

	if s.ResolversID != t.ResolversID {
		return false
	}

	if s.ResolversPrefer != t.ResolversPrefer {
		return false
	}

	if s.Retries != t.Retries {
		return false
	}

	if s.SslCaFile != t.SslCaFile {
		return false
	}

	if !equalPointers(s.SslVerify, t.SslVerify) {
		return false
	}

	if !equalPointers(s.TimeoutConnect, t.TimeoutConnect) {
		return false
	}

	return true
}

// Diff checks if two structs of type HTTPClientOptions are equal
//
//	var a, b HTTPClientOptions
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s HTTPClientOptions) Diff(t HTTPClientOptions, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if s.ResolversDisabled != t.ResolversDisabled {
		diff["ResolversDisabled"] = []interface{}{s.ResolversDisabled, t.ResolversDisabled}
	}

	if s.ResolversID != t.ResolversID {
		diff["ResolversID"] = []interface{}{s.ResolversID, t.ResolversID}
	}

	if s.ResolversPrefer != t.ResolversPrefer {
		diff["ResolversPrefer"] = []interface{}{s.ResolversPrefer, t.ResolversPrefer}
	}

	if s.Retries != t.Retries {
		diff["Retries"] = []interface{}{s.Retries, t.Retries}
	}

	if s.SslCaFile != t.SslCaFile {
		diff["SslCaFile"] = []interface{}{s.SslCaFile, t.SslCaFile}
	}

	if !equalPointers(s.SslVerify, t.SslVerify) {
		diff["SslVerify"] = []interface{}{ValueOrNil(s.SslVerify), ValueOrNil(t.SslVerify)}
	}

	if !equalPointers(s.TimeoutConnect, t.TimeoutConnect) {
		diff["TimeoutConnect"] = []interface{}{ValueOrNil(s.TimeoutConnect), ValueOrNil(t.TimeoutConnect)}
	}

	return diff
}
