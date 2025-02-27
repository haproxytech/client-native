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
	"math/rand"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/go-faker/faker/v4/pkg/options"
	"github.com/go-openapi/strfmt"

	jsoniter "github.com/json-iterator/go"
)

func TestProcessInfoItemEqual(t *testing.T) {
	samples := []struct {
		a, b ProcessInfoItem
	}{}
	for i := 0; i < 2; i++ {
		var sample ProcessInfoItem
		var result ProcessInfoItem
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
			a, b ProcessInfoItem
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
			t.Errorf("Expected ProcessInfoItem to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestProcessInfoItemEqualFalse(t *testing.T) {
	samples := []struct {
		a, b ProcessInfoItem
	}{}
	for i := 0; i < 2; i++ {
		var sample ProcessInfoItem
		var result ProcessInfoItem
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		result.ActivePeers = Ptr(*sample.ActivePeers + 1)
		result.BusyPolling = Ptr(*sample.BusyPolling + 1)
		result.BytesOutRate = Ptr(*sample.BytesOutRate + 1)
		result.CompressBpsIn = Ptr(*sample.CompressBpsIn + 1)
		result.CompressBpsOut = Ptr(*sample.CompressBpsOut + 1)
		result.CompressBpsRateLim = Ptr(*sample.CompressBpsRateLim + 1)
		result.ConnRate = Ptr(*sample.ConnRate + 1)
		result.ConnRateLimit = Ptr(*sample.ConnRateLimit + 1)
		result.ConnectedPeers = Ptr(*sample.ConnectedPeers + 1)
		result.CumConns = Ptr(*sample.CumConns + 1)
		result.CumReq = Ptr(*sample.CumReq + 1)
		result.CumSslConns = Ptr(*sample.CumSslConns + 1)
		result.CurrConns = Ptr(*sample.CurrConns + 1)
		result.CurrSslConns = Ptr(*sample.CurrSslConns + 1)
		result.DroppedLogs = Ptr(*sample.DroppedLogs + 1)
		result.FailedResolutions = Ptr(*sample.FailedResolutions + 1)
		result.HardMaxConn = Ptr(*sample.HardMaxConn + 1)
		result.IdlePct = Ptr(*sample.IdlePct + 1)
		result.Jobs = Ptr(*sample.Jobs + 1)
		result.Listeners = Ptr(*sample.Listeners + 1)
		result.MaxConn = Ptr(*sample.MaxConn + 1)
		result.MaxConnRate = Ptr(*sample.MaxConnRate + 1)
		result.MaxPipes = Ptr(*sample.MaxPipes + 1)
		result.MaxSessRate = Ptr(*sample.MaxSessRate + 1)
		result.MaxSock = Ptr(*sample.MaxSock + 1)
		result.MaxSslConns = Ptr(*sample.MaxSslConns + 1)
		result.MaxSslRate = Ptr(*sample.MaxSslRate + 1)
		result.MaxZlibMemUsage = Ptr(*sample.MaxZlibMemUsage + 1)
		result.MemMaxMb = Ptr(*sample.MemMaxMb + 1)
		result.Nbthread = Ptr(*sample.Nbthread + 1)
		result.Pid = Ptr(*sample.Pid + 1)
		result.PipesFree = Ptr(*sample.PipesFree + 1)
		result.PipesUsed = Ptr(*sample.PipesUsed + 1)
		result.PoolAllocMb = Ptr(*sample.PoolAllocMb + 1)
		result.PoolFailed = Ptr(*sample.PoolFailed + 1)
		result.PoolUsedMb = Ptr(*sample.PoolUsedMb + 1)
		result.ProcessNum = Ptr(*sample.ProcessNum + 1)
		result.Processes = Ptr(*sample.Processes + 1)
		result.ReleaseDate = strfmt.Date(time.Now().AddDate(rand.Intn(10), rand.Intn(12), rand.Intn(28)))
		result.RunQueue = Ptr(*sample.RunQueue + 1)
		result.SessRate = Ptr(*sample.SessRate + 1)
		result.SessRateLimit = Ptr(*sample.SessRateLimit + 1)
		result.SslBackendKeyRate = Ptr(*sample.SslBackendKeyRate + 1)
		result.SslBackendMaxKeyRate = Ptr(*sample.SslBackendMaxKeyRate + 1)
		result.SslCacheLookups = Ptr(*sample.SslCacheLookups + 1)
		result.SslCacheMisses = Ptr(*sample.SslCacheMisses + 1)
		result.SslFrontendKeyRate = Ptr(*sample.SslFrontendKeyRate + 1)
		result.SslFrontendMaxKeyRate = Ptr(*sample.SslFrontendMaxKeyRate + 1)
		result.SslFrontendSessionReuse = Ptr(*sample.SslFrontendSessionReuse + 1)
		result.SslRate = Ptr(*sample.SslRate + 1)
		result.SslRateLimit = Ptr(*sample.SslRateLimit + 1)
		result.Stopping = Ptr(*sample.Stopping + 1)
		result.Tasks = Ptr(*sample.Tasks + 1)
		result.TotalBytesOut = Ptr(*sample.TotalBytesOut + 1)
		result.Ulimitn = Ptr(*sample.Ulimitn + 1)
		result.Unstoppable = Ptr(*sample.Unstoppable + 1)
		result.Uptime = Ptr(*sample.Uptime + 1)
		result.ZlibMemUsage = Ptr(*sample.ZlibMemUsage + 1)
		samples = append(samples, struct {
			a, b ProcessInfoItem
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
			t.Errorf("Expected ProcessInfoItem to be different, but it is not %s %s", a, b)
		}
	}
}

func TestProcessInfoItemDiff(t *testing.T) {
	samples := []struct {
		a, b ProcessInfoItem
	}{}
	for i := 0; i < 2; i++ {
		var sample ProcessInfoItem
		var result ProcessInfoItem
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
			a, b ProcessInfoItem
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
			t.Errorf("Expected ProcessInfoItem to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestProcessInfoItemDiffFalse(t *testing.T) {
	samples := []struct {
		a, b ProcessInfoItem
	}{}
	for i := 0; i < 2; i++ {
		var sample ProcessInfoItem
		var result ProcessInfoItem
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		result.ActivePeers = Ptr(*sample.ActivePeers + 1)
		result.BusyPolling = Ptr(*sample.BusyPolling + 1)
		result.BytesOutRate = Ptr(*sample.BytesOutRate + 1)
		result.CompressBpsIn = Ptr(*sample.CompressBpsIn + 1)
		result.CompressBpsOut = Ptr(*sample.CompressBpsOut + 1)
		result.CompressBpsRateLim = Ptr(*sample.CompressBpsRateLim + 1)
		result.ConnRate = Ptr(*sample.ConnRate + 1)
		result.ConnRateLimit = Ptr(*sample.ConnRateLimit + 1)
		result.ConnectedPeers = Ptr(*sample.ConnectedPeers + 1)
		result.CumConns = Ptr(*sample.CumConns + 1)
		result.CumReq = Ptr(*sample.CumReq + 1)
		result.CumSslConns = Ptr(*sample.CumSslConns + 1)
		result.CurrConns = Ptr(*sample.CurrConns + 1)
		result.CurrSslConns = Ptr(*sample.CurrSslConns + 1)
		result.DroppedLogs = Ptr(*sample.DroppedLogs + 1)
		result.FailedResolutions = Ptr(*sample.FailedResolutions + 1)
		result.HardMaxConn = Ptr(*sample.HardMaxConn + 1)
		result.IdlePct = Ptr(*sample.IdlePct + 1)
		result.Jobs = Ptr(*sample.Jobs + 1)
		result.Listeners = Ptr(*sample.Listeners + 1)
		result.MaxConn = Ptr(*sample.MaxConn + 1)
		result.MaxConnRate = Ptr(*sample.MaxConnRate + 1)
		result.MaxPipes = Ptr(*sample.MaxPipes + 1)
		result.MaxSessRate = Ptr(*sample.MaxSessRate + 1)
		result.MaxSock = Ptr(*sample.MaxSock + 1)
		result.MaxSslConns = Ptr(*sample.MaxSslConns + 1)
		result.MaxSslRate = Ptr(*sample.MaxSslRate + 1)
		result.MaxZlibMemUsage = Ptr(*sample.MaxZlibMemUsage + 1)
		result.MemMaxMb = Ptr(*sample.MemMaxMb + 1)
		result.Nbthread = Ptr(*sample.Nbthread + 1)
		result.Pid = Ptr(*sample.Pid + 1)
		result.PipesFree = Ptr(*sample.PipesFree + 1)
		result.PipesUsed = Ptr(*sample.PipesUsed + 1)
		result.PoolAllocMb = Ptr(*sample.PoolAllocMb + 1)
		result.PoolFailed = Ptr(*sample.PoolFailed + 1)
		result.PoolUsedMb = Ptr(*sample.PoolUsedMb + 1)
		result.ProcessNum = Ptr(*sample.ProcessNum + 1)
		result.Processes = Ptr(*sample.Processes + 1)
		result.ReleaseDate = strfmt.Date(time.Now().AddDate(rand.Intn(10), rand.Intn(12), rand.Intn(28)))
		result.RunQueue = Ptr(*sample.RunQueue + 1)
		result.SessRate = Ptr(*sample.SessRate + 1)
		result.SessRateLimit = Ptr(*sample.SessRateLimit + 1)
		result.SslBackendKeyRate = Ptr(*sample.SslBackendKeyRate + 1)
		result.SslBackendMaxKeyRate = Ptr(*sample.SslBackendMaxKeyRate + 1)
		result.SslCacheLookups = Ptr(*sample.SslCacheLookups + 1)
		result.SslCacheMisses = Ptr(*sample.SslCacheMisses + 1)
		result.SslFrontendKeyRate = Ptr(*sample.SslFrontendKeyRate + 1)
		result.SslFrontendMaxKeyRate = Ptr(*sample.SslFrontendMaxKeyRate + 1)
		result.SslFrontendSessionReuse = Ptr(*sample.SslFrontendSessionReuse + 1)
		result.SslRate = Ptr(*sample.SslRate + 1)
		result.SslRateLimit = Ptr(*sample.SslRateLimit + 1)
		result.Stopping = Ptr(*sample.Stopping + 1)
		result.Tasks = Ptr(*sample.Tasks + 1)
		result.TotalBytesOut = Ptr(*sample.TotalBytesOut + 1)
		result.Ulimitn = Ptr(*sample.Ulimitn + 1)
		result.Unstoppable = Ptr(*sample.Unstoppable + 1)
		result.Uptime = Ptr(*sample.Uptime + 1)
		result.ZlibMemUsage = Ptr(*sample.ZlibMemUsage + 1)
		samples = append(samples, struct {
			a, b ProcessInfoItem
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		if len(result) != 60 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Error(err)
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Error(err)
			}
			t.Errorf("Expected ProcessInfoItem to be different in 60 cases, but it is not (%d) %s %s", len(result), a, b)
		}
	}
}
