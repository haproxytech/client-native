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

// Equal checks if two structs of type Balance are equal
//
//	var a, b Balance
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s Balance) Equal(t Balance, opts ...Options) bool {
	if !equalPointers(s.Algorithm, t.Algorithm) {
		return false
	}

	if s.HashExpression != t.HashExpression {
		return false
	}

	if s.HdrName != t.HdrName {
		return false
	}

	if s.HdrUseDomainOnly != t.HdrUseDomainOnly {
		return false
	}

	if s.RandomDraws != t.RandomDraws {
		return false
	}

	if s.RdpCookieName != t.RdpCookieName {
		return false
	}

	if s.URIDepth != t.URIDepth {
		return false
	}

	if s.URILen != t.URILen {
		return false
	}

	if s.URIPathOnly != t.URIPathOnly {
		return false
	}

	if s.URIWhole != t.URIWhole {
		return false
	}

	if s.URLParam != t.URLParam {
		return false
	}

	if s.URLParamCheckPost != t.URLParamCheckPost {
		return false
	}

	if s.URLParamMaxWait != t.URLParamMaxWait {
		return false
	}

	return true
}

// Diff checks if two structs of type Balance are equal
//
//	var a, b Balance
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s Balance) Diff(t Balance, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if !equalPointers(s.Algorithm, t.Algorithm) {
		diff["Algorithm"] = []interface{}{ValueOrNil(s.Algorithm), ValueOrNil(t.Algorithm)}
	}

	if s.HashExpression != t.HashExpression {
		diff["HashExpression"] = []interface{}{s.HashExpression, t.HashExpression}
	}

	if s.HdrName != t.HdrName {
		diff["HdrName"] = []interface{}{s.HdrName, t.HdrName}
	}

	if s.HdrUseDomainOnly != t.HdrUseDomainOnly {
		diff["HdrUseDomainOnly"] = []interface{}{s.HdrUseDomainOnly, t.HdrUseDomainOnly}
	}

	if s.RandomDraws != t.RandomDraws {
		diff["RandomDraws"] = []interface{}{s.RandomDraws, t.RandomDraws}
	}

	if s.RdpCookieName != t.RdpCookieName {
		diff["RdpCookieName"] = []interface{}{s.RdpCookieName, t.RdpCookieName}
	}

	if s.URIDepth != t.URIDepth {
		diff["URIDepth"] = []interface{}{s.URIDepth, t.URIDepth}
	}

	if s.URILen != t.URILen {
		diff["URILen"] = []interface{}{s.URILen, t.URILen}
	}

	if s.URIPathOnly != t.URIPathOnly {
		diff["URIPathOnly"] = []interface{}{s.URIPathOnly, t.URIPathOnly}
	}

	if s.URIWhole != t.URIWhole {
		diff["URIWhole"] = []interface{}{s.URIWhole, t.URIWhole}
	}

	if s.URLParam != t.URLParam {
		diff["URLParam"] = []interface{}{s.URLParam, t.URLParam}
	}

	if s.URLParamCheckPost != t.URLParamCheckPost {
		diff["URLParamCheckPost"] = []interface{}{s.URLParamCheckPost, t.URLParamCheckPost}
	}

	if s.URLParamMaxWait != t.URLParamMaxWait {
		diff["URLParamMaxWait"] = []interface{}{s.URLParamMaxWait, t.URLParamMaxWait}
	}

	return diff
}
