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

//go:build equal

package models

import (
	"encoding/json"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/go-faker/faker/v4/pkg/options"

	jsoniter "github.com/json-iterator/go"
)

func TestTuneOptionsEqual(t *testing.T) {
	samples := []struct {
		a, b TuneOptions
	}{}
	for i := 0; i < 2; i++ {
		var sample TuneOptions
		var result TuneOptions
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		byteJSON, err := json.Marshal(sample)
		if err != nil {
			t.Error(err)
		}
		err = json.Unmarshal(byteJSON, &result)
		if err != nil {
			t.Error(err)
		}

		samples = append(samples, struct {
			a, b TuneOptions
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Equal(sample.b)
		if !result {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Error(err)
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Error(err)
			}
			t.Errorf("Expected TuneOptions to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestTuneOptionsEqualFalse(t *testing.T) {
	samples := []struct {
		a, b TuneOptions
	}{}
	for i := 0; i < 2; i++ {
		var sample TuneOptions
		var result TuneOptions
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		result.CompMaxlevel = sample.CompMaxlevel + 1
		result.DisableFastForward = !sample.DisableFastForward
		result.DisableZeroCopyForwarding = !sample.DisableZeroCopyForwarding
		result.EventsMaxEventsAtOnce = sample.EventsMaxEventsAtOnce + 1
		result.FailAlloc = !sample.FailAlloc
		result.H2BeGlitchesThreshold = Ptr(*sample.H2BeGlitchesThreshold + 1)
		result.H2BeInitialWindowSize = sample.H2BeInitialWindowSize + 1
		result.H2BeMaxConcurrentStreams = sample.H2BeMaxConcurrentStreams + 1
		result.H2BeRxbuf = Ptr(*sample.H2BeRxbuf + 1)
		result.H2FeGlitchesThreshold = Ptr(*sample.H2FeGlitchesThreshold + 1)
		result.H2FeInitialWindowSize = sample.H2FeInitialWindowSize + 1
		result.H2FeMaxConcurrentStreams = sample.H2FeMaxConcurrentStreams + 1
		result.H2FeMaxTotalStreams = Ptr(*sample.H2FeMaxTotalStreams + 1)
		result.H2FeRxbuf = Ptr(*sample.H2FeRxbuf + 1)
		result.H2HeaderTableSize = sample.H2HeaderTableSize + 1
		result.H2InitialWindowSize = Ptr(*sample.H2InitialWindowSize + 1)
		result.H2MaxConcurrentStreams = sample.H2MaxConcurrentStreams + 1
		result.H2MaxFrameSize = sample.H2MaxFrameSize + 1
		result.HTTPCookielen = sample.HTTPCookielen + 1
		result.HTTPLogurilen = sample.HTTPLogurilen + 1
		result.HTTPMaxhdr = sample.HTTPMaxhdr + 1
		result.Idletimer = Ptr(*sample.Idletimer + 1)
		result.MaxChecksPerThread = Ptr(*sample.MaxChecksPerThread + 1)
		result.MaxRulesAtOnce = Ptr(*sample.MaxRulesAtOnce + 1)
		result.Maxaccept = sample.Maxaccept + 1
		result.Maxpollevents = sample.Maxpollevents + 1
		result.Maxrewrite = sample.Maxrewrite + 1
		result.MemoryHotSize = Ptr(*sample.MemoryHotSize + 1)
		result.NotsentLowatClient = Ptr(*sample.NotsentLowatClient + 1)
		result.PatternCacheSize = Ptr(*sample.PatternCacheSize + 1)
		result.PeersMaxUpdatesAtOnce = sample.PeersMaxUpdatesAtOnce + 1
		result.PoolHighFdRatio = sample.PoolHighFdRatio + 1
		result.PoolLowFdRatio = sample.PoolLowFdRatio + 1
		result.ReniceRuntime = Ptr(*sample.ReniceRuntime + 1)
		result.ReniceStartup = Ptr(*sample.ReniceStartup + 1)
		result.RingQueues = Ptr(*sample.RingQueues + 1)
		result.RunqueueDepth = sample.RunqueueDepth + 1
		result.StickCounters = Ptr(*sample.StickCounters + 1)
		samples = append(samples, struct {
			a, b TuneOptions
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Equal(sample.b)
		if result {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Error(err)
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Error(err)
			}
			t.Errorf("Expected TuneOptions to be different, but it is not %s %s", a, b)
		}
	}
}

func TestTuneOptionsDiff(t *testing.T) {
	samples := []struct {
		a, b TuneOptions
	}{}
	for i := 0; i < 2; i++ {
		var sample TuneOptions
		var result TuneOptions
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		byteJSON, err := json.Marshal(sample)
		if err != nil {
			t.Error(err)
		}
		err = json.Unmarshal(byteJSON, &result)
		if err != nil {
			t.Error(err)
		}

		samples = append(samples, struct {
			a, b TuneOptions
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		if len(result) != 0 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Error(err)
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Error(err)
			}
			t.Errorf("Expected TuneOptions to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestTuneOptionsDiffFalse(t *testing.T) {
	samples := []struct {
		a, b TuneOptions
	}{}
	for i := 0; i < 2; i++ {
		var sample TuneOptions
		var result TuneOptions
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		result.CompMaxlevel = sample.CompMaxlevel + 1
		result.DisableFastForward = !sample.DisableFastForward
		result.DisableZeroCopyForwarding = !sample.DisableZeroCopyForwarding
		result.EventsMaxEventsAtOnce = sample.EventsMaxEventsAtOnce + 1
		result.FailAlloc = !sample.FailAlloc
		result.H2BeGlitchesThreshold = Ptr(*sample.H2BeGlitchesThreshold + 1)
		result.H2BeInitialWindowSize = sample.H2BeInitialWindowSize + 1
		result.H2BeMaxConcurrentStreams = sample.H2BeMaxConcurrentStreams + 1
		result.H2BeRxbuf = Ptr(*sample.H2BeRxbuf + 1)
		result.H2FeGlitchesThreshold = Ptr(*sample.H2FeGlitchesThreshold + 1)
		result.H2FeInitialWindowSize = sample.H2FeInitialWindowSize + 1
		result.H2FeMaxConcurrentStreams = sample.H2FeMaxConcurrentStreams + 1
		result.H2FeMaxTotalStreams = Ptr(*sample.H2FeMaxTotalStreams + 1)
		result.H2FeRxbuf = Ptr(*sample.H2FeRxbuf + 1)
		result.H2HeaderTableSize = sample.H2HeaderTableSize + 1
		result.H2InitialWindowSize = Ptr(*sample.H2InitialWindowSize + 1)
		result.H2MaxConcurrentStreams = sample.H2MaxConcurrentStreams + 1
		result.H2MaxFrameSize = sample.H2MaxFrameSize + 1
		result.HTTPCookielen = sample.HTTPCookielen + 1
		result.HTTPLogurilen = sample.HTTPLogurilen + 1
		result.HTTPMaxhdr = sample.HTTPMaxhdr + 1
		result.Idletimer = Ptr(*sample.Idletimer + 1)
		result.MaxChecksPerThread = Ptr(*sample.MaxChecksPerThread + 1)
		result.MaxRulesAtOnce = Ptr(*sample.MaxRulesAtOnce + 1)
		result.Maxaccept = sample.Maxaccept + 1
		result.Maxpollevents = sample.Maxpollevents + 1
		result.Maxrewrite = sample.Maxrewrite + 1
		result.MemoryHotSize = Ptr(*sample.MemoryHotSize + 1)
		result.NotsentLowatClient = Ptr(*sample.NotsentLowatClient + 1)
		result.PatternCacheSize = Ptr(*sample.PatternCacheSize + 1)
		result.PeersMaxUpdatesAtOnce = sample.PeersMaxUpdatesAtOnce + 1
		result.PoolHighFdRatio = sample.PoolHighFdRatio + 1
		result.PoolLowFdRatio = sample.PoolLowFdRatio + 1
		result.ReniceRuntime = Ptr(*sample.ReniceRuntime + 1)
		result.ReniceStartup = Ptr(*sample.ReniceStartup + 1)
		result.RingQueues = Ptr(*sample.RingQueues + 1)
		result.RunqueueDepth = sample.RunqueueDepth + 1
		result.StickCounters = Ptr(*sample.StickCounters + 1)
		samples = append(samples, struct {
			a, b TuneOptions
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		if len(result) != 49 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Error(err)
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Error(err)
			}
			t.Errorf("Expected TuneOptions to be different in 49 cases, but it is not (%d) %s %s", len(result), a, b)
		}
	}
}
