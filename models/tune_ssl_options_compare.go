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

// Equal checks if two structs of type TuneSslOptions are equal
//
//	var a, b TuneSslOptions
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s TuneSslOptions) Equal(t TuneSslOptions, opts ...Options) bool {
	if !equalPointers(s.Cachesize, t.Cachesize) {
		return false
	}

	if !equalPointers(s.CaptureBufferSize, t.CaptureBufferSize) {
		return false
	}

	if s.CtxCacheSize != t.CtxCacheSize {
		return false
	}

	if s.DefaultDhParam != t.DefaultDhParam {
		return false
	}

	if s.ForcePrivateCache != t.ForcePrivateCache {
		return false
	}

	if s.Keylog != t.Keylog {
		return false
	}

	if !equalPointers(s.Lifetime, t.Lifetime) {
		return false
	}

	if !equalPointers(s.Maxrecord, t.Maxrecord) {
		return false
	}

	if !equalPointers(s.OcspUpdateMaxDelay, t.OcspUpdateMaxDelay) {
		return false
	}

	if !equalPointers(s.OcspUpdateMinDelay, t.OcspUpdateMinDelay) {
		return false
	}

	return true
}

// Diff checks if two structs of type TuneSslOptions are equal
//
//	var a, b TuneSslOptions
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s TuneSslOptions) Diff(t TuneSslOptions, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if !equalPointers(s.Cachesize, t.Cachesize) {
		diff["Cachesize"] = []interface{}{ValueOrNil(s.Cachesize), ValueOrNil(t.Cachesize)}
	}

	if !equalPointers(s.CaptureBufferSize, t.CaptureBufferSize) {
		diff["CaptureBufferSize"] = []interface{}{ValueOrNil(s.CaptureBufferSize), ValueOrNil(t.CaptureBufferSize)}
	}

	if s.CtxCacheSize != t.CtxCacheSize {
		diff["CtxCacheSize"] = []interface{}{s.CtxCacheSize, t.CtxCacheSize}
	}

	if s.DefaultDhParam != t.DefaultDhParam {
		diff["DefaultDhParam"] = []interface{}{s.DefaultDhParam, t.DefaultDhParam}
	}

	if s.ForcePrivateCache != t.ForcePrivateCache {
		diff["ForcePrivateCache"] = []interface{}{s.ForcePrivateCache, t.ForcePrivateCache}
	}

	if s.Keylog != t.Keylog {
		diff["Keylog"] = []interface{}{s.Keylog, t.Keylog}
	}

	if !equalPointers(s.Lifetime, t.Lifetime) {
		diff["Lifetime"] = []interface{}{ValueOrNil(s.Lifetime), ValueOrNil(t.Lifetime)}
	}

	if !equalPointers(s.Maxrecord, t.Maxrecord) {
		diff["Maxrecord"] = []interface{}{ValueOrNil(s.Maxrecord), ValueOrNil(t.Maxrecord)}
	}

	if !equalPointers(s.OcspUpdateMaxDelay, t.OcspUpdateMaxDelay) {
		diff["OcspUpdateMaxDelay"] = []interface{}{ValueOrNil(s.OcspUpdateMaxDelay), ValueOrNil(t.OcspUpdateMaxDelay)}
	}

	if !equalPointers(s.OcspUpdateMinDelay, t.OcspUpdateMinDelay) {
		diff["OcspUpdateMinDelay"] = []interface{}{ValueOrNil(s.OcspUpdateMinDelay), ValueOrNil(t.OcspUpdateMinDelay)}
	}

	return diff
}
