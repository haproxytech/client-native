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

func TestServerParamsEqual(t *testing.T) {
	samples := []struct {
		a, b ServerParams
	}{}
	for i := 0; i < 2; i++ {
		var sample ServerParams
		var result ServerParams
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
			a, b ServerParams
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
			t.Errorf("Expected ServerParams to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestServerParamsEqualFalse(t *testing.T) {
	samples := []struct {
		a, b ServerParams
	}{}
	for i := 0; i < 2; i++ {
		var sample ServerParams
		var result ServerParams
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		result.AgentInter = Ptr(*sample.AgentInter + 1)
		result.AgentPort = Ptr(*sample.AgentPort + 1)
		result.Allow0rtt = !sample.Allow0rtt
		result.Downinter = Ptr(*sample.Downinter + 1)
		result.ErrorLimit = sample.ErrorLimit + 1
		result.Fall = Ptr(*sample.Fall + 1)
		result.Fastinter = Ptr(*sample.Fastinter + 1)
		result.HealthCheckPort = Ptr(*sample.HealthCheckPort + 1)
		result.IdlePing = Ptr(*sample.IdlePing + 1)
		result.Inter = Ptr(*sample.Inter + 1)
		result.LogBufsize = Ptr(*sample.LogBufsize + 1)
		result.MaxReuse = Ptr(*sample.MaxReuse + 1)
		result.Maxconn = Ptr(*sample.Maxconn + 1)
		result.Maxqueue = Ptr(*sample.Maxqueue + 1)
		result.Minconn = Ptr(*sample.Minconn + 1)
		result.PoolLowConn = Ptr(*sample.PoolLowConn + 1)
		result.PoolMaxConn = Ptr(*sample.PoolMaxConn + 1)
		result.PoolPurgeDelay = Ptr(*sample.PoolPurgeDelay + 1)
		result.Rise = Ptr(*sample.Rise + 1)
		result.Shard = sample.Shard + 1
		result.Slowstart = Ptr(*sample.Slowstart + 1)
		result.StrictMaxconn = !sample.StrictMaxconn
		result.TCPUt = Ptr(*sample.TCPUt + 1)
		result.Weight = Ptr(*sample.Weight + 1)
		samples = append(samples, struct {
			a, b ServerParams
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
			t.Errorf("Expected ServerParams to be different, but it is not %s %s", a, b)
		}
	}
}

func TestServerParamsDiff(t *testing.T) {
	samples := []struct {
		a, b ServerParams
	}{}
	for i := 0; i < 2; i++ {
		var sample ServerParams
		var result ServerParams
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
			a, b ServerParams
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
			t.Errorf("Expected ServerParams to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestServerParamsDiffFalse(t *testing.T) {
	samples := []struct {
		a, b ServerParams
	}{}
	for i := 0; i < 2; i++ {
		var sample ServerParams
		var result ServerParams
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		result.AgentInter = Ptr(*sample.AgentInter + 1)
		result.AgentPort = Ptr(*sample.AgentPort + 1)
		result.Allow0rtt = !sample.Allow0rtt
		result.Downinter = Ptr(*sample.Downinter + 1)
		result.ErrorLimit = sample.ErrorLimit + 1
		result.Fall = Ptr(*sample.Fall + 1)
		result.Fastinter = Ptr(*sample.Fastinter + 1)
		result.HealthCheckPort = Ptr(*sample.HealthCheckPort + 1)
		result.IdlePing = Ptr(*sample.IdlePing + 1)
		result.Inter = Ptr(*sample.Inter + 1)
		result.LogBufsize = Ptr(*sample.LogBufsize + 1)
		result.MaxReuse = Ptr(*sample.MaxReuse + 1)
		result.Maxconn = Ptr(*sample.Maxconn + 1)
		result.Maxqueue = Ptr(*sample.Maxqueue + 1)
		result.Minconn = Ptr(*sample.Minconn + 1)
		result.PoolLowConn = Ptr(*sample.PoolLowConn + 1)
		result.PoolMaxConn = Ptr(*sample.PoolMaxConn + 1)
		result.PoolPurgeDelay = Ptr(*sample.PoolPurgeDelay + 1)
		result.Rise = Ptr(*sample.Rise + 1)
		result.Shard = sample.Shard + 1
		result.Slowstart = Ptr(*sample.Slowstart + 1)
		result.StrictMaxconn = !sample.StrictMaxconn
		result.TCPUt = Ptr(*sample.TCPUt + 1)
		result.Weight = Ptr(*sample.Weight + 1)
		samples = append(samples, struct {
			a, b ServerParams
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		if len(result) != 103 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Error(err)
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Error(err)
			}
			t.Errorf("Expected ServerParams to be different in 103 cases, but it is not (%d) %s %s", len(result), a, b)
		}
	}
}

func TestServerParamsSetProxyV2TlvFmtEqual(t *testing.T) {
	samples := []struct {
		a, b ServerParamsSetProxyV2TlvFmt
	}{}
	for i := 0; i < 2; i++ {
		var sample ServerParamsSetProxyV2TlvFmt
		var result ServerParamsSetProxyV2TlvFmt
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
			a, b ServerParamsSetProxyV2TlvFmt
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
			t.Errorf("Expected ServerParamsSetProxyV2TlvFmt to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestServerParamsSetProxyV2TlvFmtEqualFalse(t *testing.T) {
	samples := []struct {
		a, b ServerParamsSetProxyV2TlvFmt
	}{}
	for i := 0; i < 2; i++ {
		var sample ServerParamsSetProxyV2TlvFmt
		var result ServerParamsSetProxyV2TlvFmt
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		samples = append(samples, struct {
			a, b ServerParamsSetProxyV2TlvFmt
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
			t.Errorf("Expected ServerParamsSetProxyV2TlvFmt to be different, but it is not %s %s", a, b)
		}
	}
}

func TestServerParamsSetProxyV2TlvFmtDiff(t *testing.T) {
	samples := []struct {
		a, b ServerParamsSetProxyV2TlvFmt
	}{}
	for i := 0; i < 2; i++ {
		var sample ServerParamsSetProxyV2TlvFmt
		var result ServerParamsSetProxyV2TlvFmt
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
			a, b ServerParamsSetProxyV2TlvFmt
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
			t.Errorf("Expected ServerParamsSetProxyV2TlvFmt to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestServerParamsSetProxyV2TlvFmtDiffFalse(t *testing.T) {
	samples := []struct {
		a, b ServerParamsSetProxyV2TlvFmt
	}{}
	for i := 0; i < 2; i++ {
		var sample ServerParamsSetProxyV2TlvFmt
		var result ServerParamsSetProxyV2TlvFmt
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		samples = append(samples, struct {
			a, b ServerParamsSetProxyV2TlvFmt
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		if len(result) != 2 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Error(err)
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Error(err)
			}
			t.Errorf("Expected ServerParamsSetProxyV2TlvFmt to be different in 2 cases, but it is not (%d) %s %s", len(result), a, b)
		}
	}
}
