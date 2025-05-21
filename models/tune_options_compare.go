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

// Equal checks if two structs of type TuneOptions are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b TuneOptions
//	equal := a.Equal(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b TuneOptions
//	equal := a.Equal(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s TuneOptions) Equal(t TuneOptions, opts ...Options) bool {
	opt := getOptions(opts...)

	if s.AppletZeroCopyForwarding != t.AppletZeroCopyForwarding {
		return false
	}

	if s.CompMaxlevel != t.CompMaxlevel {
		return false
	}

	if s.DisableFastForward != t.DisableFastForward {
		return false
	}

	if s.DisableZeroCopyForwarding != t.DisableZeroCopyForwarding {
		return false
	}

	if !equalComparableSlice(s.EpollMaskEvents, t.EpollMaskEvents, opt) {
		return false
	}

	if s.EventsMaxEventsAtOnce != t.EventsMaxEventsAtOnce {
		return false
	}

	if s.FailAlloc != t.FailAlloc {
		return false
	}

	if s.FdEdgeTriggered != t.FdEdgeTriggered {
		return false
	}

	if s.H1ZeroCopyFwdRecv != t.H1ZeroCopyFwdRecv {
		return false
	}

	if s.H1ZeroCopyFwdSend != t.H1ZeroCopyFwdSend {
		return false
	}

	if !equalPointers(s.H2BeGlitchesThreshold, t.H2BeGlitchesThreshold) {
		return false
	}

	if s.H2BeInitialWindowSize != t.H2BeInitialWindowSize {
		return false
	}

	if s.H2BeMaxConcurrentStreams != t.H2BeMaxConcurrentStreams {
		return false
	}

	if !equalPointers(s.H2BeRxbuf, t.H2BeRxbuf) {
		return false
	}

	if !equalPointers(s.H2FeGlitchesThreshold, t.H2FeGlitchesThreshold) {
		return false
	}

	if s.H2FeInitialWindowSize != t.H2FeInitialWindowSize {
		return false
	}

	if s.H2FeMaxConcurrentStreams != t.H2FeMaxConcurrentStreams {
		return false
	}

	if !equalPointers(s.H2FeMaxTotalStreams, t.H2FeMaxTotalStreams) {
		return false
	}

	if !equalPointers(s.H2FeRxbuf, t.H2FeRxbuf) {
		return false
	}

	if s.H2HeaderTableSize != t.H2HeaderTableSize {
		return false
	}

	if !equalPointers(s.H2InitialWindowSize, t.H2InitialWindowSize) {
		return false
	}

	if s.H2MaxConcurrentStreams != t.H2MaxConcurrentStreams {
		return false
	}

	if s.H2MaxFrameSize != t.H2MaxFrameSize {
		return false
	}

	if s.H2ZeroCopyFwdSend != t.H2ZeroCopyFwdSend {
		return false
	}

	if s.HTTPCookielen != t.HTTPCookielen {
		return false
	}

	if s.HTTPLogurilen != t.HTTPLogurilen {
		return false
	}

	if s.HTTPMaxhdr != t.HTTPMaxhdr {
		return false
	}

	if s.IdlePoolShared != t.IdlePoolShared {
		return false
	}

	if !equalPointers(s.Idletimer, t.Idletimer) {
		return false
	}

	if s.ListenerDefaultShards != t.ListenerDefaultShards {
		return false
	}

	if s.ListenerMultiQueue != t.ListenerMultiQueue {
		return false
	}

	if !equalPointers(s.MaxChecksPerThread, t.MaxChecksPerThread) {
		return false
	}

	if !equalPointers(s.MaxRulesAtOnce, t.MaxRulesAtOnce) {
		return false
	}

	if s.Maxaccept != t.Maxaccept {
		return false
	}

	if s.Maxpollevents != t.Maxpollevents {
		return false
	}

	if s.Maxrewrite != t.Maxrewrite {
		return false
	}

	if !equalPointers(s.MemoryHotSize, t.MemoryHotSize) {
		return false
	}

	if !equalPointers(s.NotsentLowatClient, t.NotsentLowatClient) {
		return false
	}

	if !equalPointers(s.NotsentLowatServer, t.NotsentLowatServer) {
		return false
	}

	if !equalPointers(s.PatternCacheSize, t.PatternCacheSize) {
		return false
	}

	if s.PeersMaxUpdatesAtOnce != t.PeersMaxUpdatesAtOnce {
		return false
	}

	if s.PoolHighFdRatio != t.PoolHighFdRatio {
		return false
	}

	if s.PoolLowFdRatio != t.PoolLowFdRatio {
		return false
	}

	if s.PtZeroCopyForwarding != t.PtZeroCopyForwarding {
		return false
	}

	if !equalPointers(s.ReniceRuntime, t.ReniceRuntime) {
		return false
	}

	if !equalPointers(s.ReniceStartup, t.ReniceStartup) {
		return false
	}

	if !equalPointers(s.RingQueues, t.RingQueues) {
		return false
	}

	if s.RunqueueDepth != t.RunqueueDepth {
		return false
	}

	if s.SchedLowLatency != t.SchedLowLatency {
		return false
	}

	if !equalPointers(s.StickCounters, t.StickCounters) {
		return false
	}

	return true
}

// Diff checks if two structs of type TuneOptions are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b TuneOptions
//	diff := a.Diff(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b TuneOptions
//	diff := a.Diff(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s TuneOptions) Diff(t TuneOptions, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})
	if s.AppletZeroCopyForwarding != t.AppletZeroCopyForwarding {
		diff["AppletZeroCopyForwarding"] = []interface{}{s.AppletZeroCopyForwarding, t.AppletZeroCopyForwarding}
	}

	if s.CompMaxlevel != t.CompMaxlevel {
		diff["CompMaxlevel"] = []interface{}{s.CompMaxlevel, t.CompMaxlevel}
	}

	if s.DisableFastForward != t.DisableFastForward {
		diff["DisableFastForward"] = []interface{}{s.DisableFastForward, t.DisableFastForward}
	}

	if s.DisableZeroCopyForwarding != t.DisableZeroCopyForwarding {
		diff["DisableZeroCopyForwarding"] = []interface{}{s.DisableZeroCopyForwarding, t.DisableZeroCopyForwarding}
	}

	if !equalComparableSlice(s.EpollMaskEvents, t.EpollMaskEvents, opt) {
		diff["EpollMaskEvents"] = []interface{}{s.EpollMaskEvents, t.EpollMaskEvents}
	}

	if s.EventsMaxEventsAtOnce != t.EventsMaxEventsAtOnce {
		diff["EventsMaxEventsAtOnce"] = []interface{}{s.EventsMaxEventsAtOnce, t.EventsMaxEventsAtOnce}
	}

	if s.FailAlloc != t.FailAlloc {
		diff["FailAlloc"] = []interface{}{s.FailAlloc, t.FailAlloc}
	}

	if s.FdEdgeTriggered != t.FdEdgeTriggered {
		diff["FdEdgeTriggered"] = []interface{}{s.FdEdgeTriggered, t.FdEdgeTriggered}
	}

	if s.H1ZeroCopyFwdRecv != t.H1ZeroCopyFwdRecv {
		diff["H1ZeroCopyFwdRecv"] = []interface{}{s.H1ZeroCopyFwdRecv, t.H1ZeroCopyFwdRecv}
	}

	if s.H1ZeroCopyFwdSend != t.H1ZeroCopyFwdSend {
		diff["H1ZeroCopyFwdSend"] = []interface{}{s.H1ZeroCopyFwdSend, t.H1ZeroCopyFwdSend}
	}

	if !equalPointers(s.H2BeGlitchesThreshold, t.H2BeGlitchesThreshold) {
		diff["H2BeGlitchesThreshold"] = []interface{}{ValueOrNil(s.H2BeGlitchesThreshold), ValueOrNil(t.H2BeGlitchesThreshold)}
	}

	if s.H2BeInitialWindowSize != t.H2BeInitialWindowSize {
		diff["H2BeInitialWindowSize"] = []interface{}{s.H2BeInitialWindowSize, t.H2BeInitialWindowSize}
	}

	if s.H2BeMaxConcurrentStreams != t.H2BeMaxConcurrentStreams {
		diff["H2BeMaxConcurrentStreams"] = []interface{}{s.H2BeMaxConcurrentStreams, t.H2BeMaxConcurrentStreams}
	}

	if !equalPointers(s.H2BeRxbuf, t.H2BeRxbuf) {
		diff["H2BeRxbuf"] = []interface{}{ValueOrNil(s.H2BeRxbuf), ValueOrNil(t.H2BeRxbuf)}
	}

	if !equalPointers(s.H2FeGlitchesThreshold, t.H2FeGlitchesThreshold) {
		diff["H2FeGlitchesThreshold"] = []interface{}{ValueOrNil(s.H2FeGlitchesThreshold), ValueOrNil(t.H2FeGlitchesThreshold)}
	}

	if s.H2FeInitialWindowSize != t.H2FeInitialWindowSize {
		diff["H2FeInitialWindowSize"] = []interface{}{s.H2FeInitialWindowSize, t.H2FeInitialWindowSize}
	}

	if s.H2FeMaxConcurrentStreams != t.H2FeMaxConcurrentStreams {
		diff["H2FeMaxConcurrentStreams"] = []interface{}{s.H2FeMaxConcurrentStreams, t.H2FeMaxConcurrentStreams}
	}

	if !equalPointers(s.H2FeMaxTotalStreams, t.H2FeMaxTotalStreams) {
		diff["H2FeMaxTotalStreams"] = []interface{}{ValueOrNil(s.H2FeMaxTotalStreams), ValueOrNil(t.H2FeMaxTotalStreams)}
	}

	if !equalPointers(s.H2FeRxbuf, t.H2FeRxbuf) {
		diff["H2FeRxbuf"] = []interface{}{ValueOrNil(s.H2FeRxbuf), ValueOrNil(t.H2FeRxbuf)}
	}

	if s.H2HeaderTableSize != t.H2HeaderTableSize {
		diff["H2HeaderTableSize"] = []interface{}{s.H2HeaderTableSize, t.H2HeaderTableSize}
	}

	if !equalPointers(s.H2InitialWindowSize, t.H2InitialWindowSize) {
		diff["H2InitialWindowSize"] = []interface{}{ValueOrNil(s.H2InitialWindowSize), ValueOrNil(t.H2InitialWindowSize)}
	}

	if s.H2MaxConcurrentStreams != t.H2MaxConcurrentStreams {
		diff["H2MaxConcurrentStreams"] = []interface{}{s.H2MaxConcurrentStreams, t.H2MaxConcurrentStreams}
	}

	if s.H2MaxFrameSize != t.H2MaxFrameSize {
		diff["H2MaxFrameSize"] = []interface{}{s.H2MaxFrameSize, t.H2MaxFrameSize}
	}

	if s.H2ZeroCopyFwdSend != t.H2ZeroCopyFwdSend {
		diff["H2ZeroCopyFwdSend"] = []interface{}{s.H2ZeroCopyFwdSend, t.H2ZeroCopyFwdSend}
	}

	if s.HTTPCookielen != t.HTTPCookielen {
		diff["HTTPCookielen"] = []interface{}{s.HTTPCookielen, t.HTTPCookielen}
	}

	if s.HTTPLogurilen != t.HTTPLogurilen {
		diff["HTTPLogurilen"] = []interface{}{s.HTTPLogurilen, t.HTTPLogurilen}
	}

	if s.HTTPMaxhdr != t.HTTPMaxhdr {
		diff["HTTPMaxhdr"] = []interface{}{s.HTTPMaxhdr, t.HTTPMaxhdr}
	}

	if s.IdlePoolShared != t.IdlePoolShared {
		diff["IdlePoolShared"] = []interface{}{s.IdlePoolShared, t.IdlePoolShared}
	}

	if !equalPointers(s.Idletimer, t.Idletimer) {
		diff["Idletimer"] = []interface{}{ValueOrNil(s.Idletimer), ValueOrNil(t.Idletimer)}
	}

	if s.ListenerDefaultShards != t.ListenerDefaultShards {
		diff["ListenerDefaultShards"] = []interface{}{s.ListenerDefaultShards, t.ListenerDefaultShards}
	}

	if s.ListenerMultiQueue != t.ListenerMultiQueue {
		diff["ListenerMultiQueue"] = []interface{}{s.ListenerMultiQueue, t.ListenerMultiQueue}
	}

	if !equalPointers(s.MaxChecksPerThread, t.MaxChecksPerThread) {
		diff["MaxChecksPerThread"] = []interface{}{ValueOrNil(s.MaxChecksPerThread), ValueOrNil(t.MaxChecksPerThread)}
	}

	if !equalPointers(s.MaxRulesAtOnce, t.MaxRulesAtOnce) {
		diff["MaxRulesAtOnce"] = []interface{}{ValueOrNil(s.MaxRulesAtOnce), ValueOrNil(t.MaxRulesAtOnce)}
	}

	if s.Maxaccept != t.Maxaccept {
		diff["Maxaccept"] = []interface{}{s.Maxaccept, t.Maxaccept}
	}

	if s.Maxpollevents != t.Maxpollevents {
		diff["Maxpollevents"] = []interface{}{s.Maxpollevents, t.Maxpollevents}
	}

	if s.Maxrewrite != t.Maxrewrite {
		diff["Maxrewrite"] = []interface{}{s.Maxrewrite, t.Maxrewrite}
	}

	if !equalPointers(s.MemoryHotSize, t.MemoryHotSize) {
		diff["MemoryHotSize"] = []interface{}{ValueOrNil(s.MemoryHotSize), ValueOrNil(t.MemoryHotSize)}
	}

	if !equalPointers(s.NotsentLowatClient, t.NotsentLowatClient) {
		diff["NotsentLowatClient"] = []interface{}{ValueOrNil(s.NotsentLowatClient), ValueOrNil(t.NotsentLowatClient)}
	}

	if !equalPointers(s.NotsentLowatServer, t.NotsentLowatServer) {
		diff["NotsentLowatServer"] = []interface{}{ValueOrNil(s.NotsentLowatServer), ValueOrNil(t.NotsentLowatServer)}
	}

	if !equalPointers(s.PatternCacheSize, t.PatternCacheSize) {
		diff["PatternCacheSize"] = []interface{}{ValueOrNil(s.PatternCacheSize), ValueOrNil(t.PatternCacheSize)}
	}

	if s.PeersMaxUpdatesAtOnce != t.PeersMaxUpdatesAtOnce {
		diff["PeersMaxUpdatesAtOnce"] = []interface{}{s.PeersMaxUpdatesAtOnce, t.PeersMaxUpdatesAtOnce}
	}

	if s.PoolHighFdRatio != t.PoolHighFdRatio {
		diff["PoolHighFdRatio"] = []interface{}{s.PoolHighFdRatio, t.PoolHighFdRatio}
	}

	if s.PoolLowFdRatio != t.PoolLowFdRatio {
		diff["PoolLowFdRatio"] = []interface{}{s.PoolLowFdRatio, t.PoolLowFdRatio}
	}

	if s.PtZeroCopyForwarding != t.PtZeroCopyForwarding {
		diff["PtZeroCopyForwarding"] = []interface{}{s.PtZeroCopyForwarding, t.PtZeroCopyForwarding}
	}

	if !equalPointers(s.ReniceRuntime, t.ReniceRuntime) {
		diff["ReniceRuntime"] = []interface{}{ValueOrNil(s.ReniceRuntime), ValueOrNil(t.ReniceRuntime)}
	}

	if !equalPointers(s.ReniceStartup, t.ReniceStartup) {
		diff["ReniceStartup"] = []interface{}{ValueOrNil(s.ReniceStartup), ValueOrNil(t.ReniceStartup)}
	}

	if !equalPointers(s.RingQueues, t.RingQueues) {
		diff["RingQueues"] = []interface{}{ValueOrNil(s.RingQueues), ValueOrNil(t.RingQueues)}
	}

	if s.RunqueueDepth != t.RunqueueDepth {
		diff["RunqueueDepth"] = []interface{}{s.RunqueueDepth, t.RunqueueDepth}
	}

	if s.SchedLowLatency != t.SchedLowLatency {
		diff["SchedLowLatency"] = []interface{}{s.SchedLowLatency, t.SchedLowLatency}
	}

	if !equalPointers(s.StickCounters, t.StickCounters) {
		diff["StickCounters"] = []interface{}{ValueOrNil(s.StickCounters), ValueOrNil(t.StickCounters)}
	}

	return diff
}
