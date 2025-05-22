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

// Equal checks if two structs of type TuneQuicOptions are equal
//
//	var a, b TuneQuicOptions
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s TuneQuicOptions) Equal(t TuneQuicOptions, opts ...Options) bool {
	if !equalPointers(s.FrontendConnTxBuffersLimit, t.FrontendConnTxBuffersLimit) {
		return false
	}

	if !equalPointers(s.FrontendMaxIdleTimeout, t.FrontendMaxIdleTimeout) {
		return false
	}

	if !equalPointers(s.FrontendMaxStreamsBidi, t.FrontendMaxStreamsBidi) {
		return false
	}

	if !equalPointers(s.FrontendMaxTxMemory, t.FrontendMaxTxMemory) {
		return false
	}

	if !equalPointers(s.MaxFrameLoss, t.MaxFrameLoss) {
		return false
	}

	if !equalPointers(s.ReorderRatio, t.ReorderRatio) {
		return false
	}

	if !equalPointers(s.RetryThreshold, t.RetryThreshold) {
		return false
	}

	if s.SocketOwner != t.SocketOwner {
		return false
	}

	if s.ZeroCopyFwdSend != t.ZeroCopyFwdSend {
		return false
	}

	return true
}

// Diff checks if two structs of type TuneQuicOptions are equal
//
//	var a, b TuneQuicOptions
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s TuneQuicOptions) Diff(t TuneQuicOptions, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if !equalPointers(s.FrontendConnTxBuffersLimit, t.FrontendConnTxBuffersLimit) {
		diff["FrontendConnTxBuffersLimit"] = []interface{}{ValueOrNil(s.FrontendConnTxBuffersLimit), ValueOrNil(t.FrontendConnTxBuffersLimit)}
	}

	if !equalPointers(s.FrontendMaxIdleTimeout, t.FrontendMaxIdleTimeout) {
		diff["FrontendMaxIdleTimeout"] = []interface{}{ValueOrNil(s.FrontendMaxIdleTimeout), ValueOrNil(t.FrontendMaxIdleTimeout)}
	}

	if !equalPointers(s.FrontendMaxStreamsBidi, t.FrontendMaxStreamsBidi) {
		diff["FrontendMaxStreamsBidi"] = []interface{}{ValueOrNil(s.FrontendMaxStreamsBidi), ValueOrNil(t.FrontendMaxStreamsBidi)}
	}

	if !equalPointers(s.FrontendMaxTxMemory, t.FrontendMaxTxMemory) {
		diff["FrontendMaxTxMemory"] = []interface{}{ValueOrNil(s.FrontendMaxTxMemory), ValueOrNil(t.FrontendMaxTxMemory)}
	}

	if !equalPointers(s.MaxFrameLoss, t.MaxFrameLoss) {
		diff["MaxFrameLoss"] = []interface{}{ValueOrNil(s.MaxFrameLoss), ValueOrNil(t.MaxFrameLoss)}
	}

	if !equalPointers(s.ReorderRatio, t.ReorderRatio) {
		diff["ReorderRatio"] = []interface{}{ValueOrNil(s.ReorderRatio), ValueOrNil(t.ReorderRatio)}
	}

	if !equalPointers(s.RetryThreshold, t.RetryThreshold) {
		diff["RetryThreshold"] = []interface{}{ValueOrNil(s.RetryThreshold), ValueOrNil(t.RetryThreshold)}
	}

	if s.SocketOwner != t.SocketOwner {
		diff["SocketOwner"] = []interface{}{s.SocketOwner, t.SocketOwner}
	}

	if s.ZeroCopyFwdSend != t.ZeroCopyFwdSend {
		diff["ZeroCopyFwdSend"] = []interface{}{s.ZeroCopyFwdSend, t.ZeroCopyFwdSend}
	}

	return diff
}
