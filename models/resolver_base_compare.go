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

import "reflect"

// Equal checks if two structs of type ResolverBase are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b ResolverBase
//	equal := a.Equal(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b ResolverBase
//	equal := a.Equal(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s ResolverBase) Equal(t ResolverBase, opts ...Options) bool {
	opt := getOptions(opts...)

	if s.AcceptedPayloadSize != t.AcceptedPayloadSize {
		return false
	}

	if !equalPointers(s.HoldNx, t.HoldNx) {
		return false
	}

	if !equalPointers(s.HoldObsolete, t.HoldObsolete) {
		return false
	}

	if !equalPointers(s.HoldOther, t.HoldOther) {
		return false
	}

	if !equalPointers(s.HoldRefused, t.HoldRefused) {
		return false
	}

	if !equalPointers(s.HoldTimeout, t.HoldTimeout) {
		return false
	}

	if !equalPointers(s.HoldValid, t.HoldValid) {
		return false
	}

	if !CheckSameNilAndLenMap[string](s.Metadata, t.Metadata, opt) {
		return false
	}

	for k, v := range s.Metadata {
		if !reflect.DeepEqual(t.Metadata[k], v) {
			return false
		}
	}

	if s.Name != t.Name {
		return false
	}

	if s.ParseResolvConf != t.ParseResolvConf {
		return false
	}

	if s.ResolveRetries != t.ResolveRetries {
		return false
	}

	if s.TimeoutResolve != t.TimeoutResolve {
		return false
	}

	if s.TimeoutRetry != t.TimeoutRetry {
		return false
	}

	return true
}

// Diff checks if two structs of type ResolverBase are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b ResolverBase
//	diff := a.Diff(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b ResolverBase
//	diff := a.Diff(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s ResolverBase) Diff(t ResolverBase, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})
	if s.AcceptedPayloadSize != t.AcceptedPayloadSize {
		diff["AcceptedPayloadSize"] = []interface{}{s.AcceptedPayloadSize, t.AcceptedPayloadSize}
	}

	if !equalPointers(s.HoldNx, t.HoldNx) {
		diff["HoldNx"] = []interface{}{ValueOrNil(s.HoldNx), ValueOrNil(t.HoldNx)}
	}

	if !equalPointers(s.HoldObsolete, t.HoldObsolete) {
		diff["HoldObsolete"] = []interface{}{ValueOrNil(s.HoldObsolete), ValueOrNil(t.HoldObsolete)}
	}

	if !equalPointers(s.HoldOther, t.HoldOther) {
		diff["HoldOther"] = []interface{}{ValueOrNil(s.HoldOther), ValueOrNil(t.HoldOther)}
	}

	if !equalPointers(s.HoldRefused, t.HoldRefused) {
		diff["HoldRefused"] = []interface{}{ValueOrNil(s.HoldRefused), ValueOrNil(t.HoldRefused)}
	}

	if !equalPointers(s.HoldTimeout, t.HoldTimeout) {
		diff["HoldTimeout"] = []interface{}{ValueOrNil(s.HoldTimeout), ValueOrNil(t.HoldTimeout)}
	}

	if !equalPointers(s.HoldValid, t.HoldValid) {
		diff["HoldValid"] = []interface{}{ValueOrNil(s.HoldValid), ValueOrNil(t.HoldValid)}
	}

	if !CheckSameNilAndLenMap[string](s.Metadata, t.Metadata, opt) {
		diff["Metadata"] = []interface{}{s.Metadata, t.Metadata}
	}

	for k, v := range s.Metadata {
		if !reflect.DeepEqual(t.Metadata[k], v) {
			diff["Metadata"] = []interface{}{s.Metadata, t.Metadata}
		}
	}

	if s.Name != t.Name {
		diff["Name"] = []interface{}{s.Name, t.Name}
	}

	if s.ParseResolvConf != t.ParseResolvConf {
		diff["ParseResolvConf"] = []interface{}{s.ParseResolvConf, t.ParseResolvConf}
	}

	if s.ResolveRetries != t.ResolveRetries {
		diff["ResolveRetries"] = []interface{}{s.ResolveRetries, t.ResolveRetries}
	}

	if s.TimeoutResolve != t.TimeoutResolve {
		diff["TimeoutResolve"] = []interface{}{s.TimeoutResolve, t.TimeoutResolve}
	}

	if s.TimeoutRetry != t.TimeoutRetry {
		diff["TimeoutRetry"] = []interface{}{s.TimeoutRetry, t.TimeoutRetry}
	}

	return diff
}
