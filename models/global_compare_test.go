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

import (
	"encoding/json"
	"testing"

	"github.com/go-faker/faker/v4"

	jsoniter "github.com/json-iterator/go"
)

func TestGlobalEqual(t *testing.T) {
	samples := []struct {
		a, b Global
	}{}
	for i := 0; i < 2; i++ {
		var sample Global
		var result Global
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		byteJSON, err := json.Marshal(sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = json.Unmarshal(byteJSON, &result)
		if err != nil {
			t.Errorf(err.Error())
		}

		samples = append(samples, struct {
			a, b Global
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Equal(sample.b)
		if !result {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected Global to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestGlobalEqualFalse(t *testing.T) {
	samples := []struct {
		a, b Global
	}{}
	for i := 0; i < 2; i++ {
		var sample Global
		var result Global
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		result.Anonkey = Ptr(*sample.Anonkey + 1)
		result.BusyPolling = !sample.BusyPolling
		result.CloseSpreadTime = Ptr(*sample.CloseSpreadTime + 1)
		result.ExposeExperimentalDirectives = !sample.ExposeExperimentalDirectives
		result.ExternalCheck = !sample.ExternalCheck
		result.Gid = sample.Gid + 1
		result.Grace = Ptr(*sample.Grace + 1)
		result.H2WorkaroundBogusWebsocketClients = !sample.H2WorkaroundBogusWebsocketClients
		result.HardStopAfter = Ptr(*sample.HardStopAfter + 1)
		result.HttpclientRetries = sample.HttpclientRetries + 1
		result.HttpclientTimeoutConnect = Ptr(*sample.HttpclientTimeoutConnect + 1)
		result.InsecureForkWanted = !sample.InsecureForkWanted
		result.InsecureSetuidWanted = !sample.InsecureSetuidWanted
		result.LimitedQuic = !sample.LimitedQuic
		result.MasterWorker = !sample.MasterWorker
		result.MaxSpreadChecks = Ptr(*sample.MaxSpreadChecks + 1)
		result.Maxcompcpuusage = sample.Maxcompcpuusage + 1
		result.Maxcomprate = sample.Maxcomprate + 1
		result.Maxconn = sample.Maxconn + 1
		result.Maxconnrate = sample.Maxconnrate + 1
		result.Maxpipes = sample.Maxpipes + 1
		result.Maxsessrate = sample.Maxsessrate + 1
		result.Maxsslconn = sample.Maxsslconn + 1
		result.Maxsslrate = sample.Maxsslrate + 1
		result.Maxzlibmem = sample.Maxzlibmem + 1
		result.MworkerMaxReloads = Ptr(*sample.MworkerMaxReloads + 1)
		result.Nbproc = sample.Nbproc + 1
		result.Nbthread = sample.Nbthread + 1
		result.NoQuic = !sample.NoQuic
		result.Noepoll = !sample.Noepoll
		result.Noevports = !sample.Noevports
		result.Nogetaddrinfo = !sample.Nogetaddrinfo
		result.Nokqueue = !sample.Nokqueue
		result.Nopoll = !sample.Nopoll
		result.Noreuseport = !sample.Noreuseport
		result.Nosplice = !sample.Nosplice
		result.Pp2NeverSendLocal = !sample.Pp2NeverSendLocal
		result.PreallocFd = !sample.PreallocFd
		result.Quiet = !sample.Quiet
		result.SetDumpable = !sample.SetDumpable
		result.SpreadChecks = sample.SpreadChecks + 1
		result.SslSkipSelfIssuedCa = !sample.SslSkipSelfIssuedCa
		result.StatsMaxconn = Ptr(*sample.StatsMaxconn + 1)
		result.StatsTimeout = Ptr(*sample.StatsTimeout + 1)
		result.StrictLimits = !sample.StrictLimits
		result.ThreadGroups = sample.ThreadGroups + 1
		result.TuneSslDefaultDhParam = sample.TuneSslDefaultDhParam + 1
		result.UID = sample.UID + 1
		result.Ulimitn = sample.Ulimitn + 1
		result.ZeroWarning = !sample.ZeroWarning
		samples = append(samples, struct {
			a, b Global
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Equal(sample.b)
		if result {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected Global to be different, but it is not %s %s", a, b)
		}
	}
}

func TestGlobalDiff(t *testing.T) {
	samples := []struct {
		a, b Global
	}{}
	for i := 0; i < 2; i++ {
		var sample Global
		var result Global
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		byteJSON, err := json.Marshal(sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = json.Unmarshal(byteJSON, &result)
		if err != nil {
			t.Errorf(err.Error())
		}

		samples = append(samples, struct {
			a, b Global
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		if len(result) != 0 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected Global to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestGlobalDiffFalse(t *testing.T) {
	samples := []struct {
		a, b Global
	}{}
	for i := 0; i < 2; i++ {
		var sample Global
		var result Global
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		result.Anonkey = Ptr(*sample.Anonkey + 1)
		result.BusyPolling = !sample.BusyPolling
		result.CloseSpreadTime = Ptr(*sample.CloseSpreadTime + 1)
		result.ExposeExperimentalDirectives = !sample.ExposeExperimentalDirectives
		result.ExternalCheck = !sample.ExternalCheck
		result.Gid = sample.Gid + 1
		result.Grace = Ptr(*sample.Grace + 1)
		result.H2WorkaroundBogusWebsocketClients = !sample.H2WorkaroundBogusWebsocketClients
		result.HardStopAfter = Ptr(*sample.HardStopAfter + 1)
		result.HttpclientRetries = sample.HttpclientRetries + 1
		result.HttpclientTimeoutConnect = Ptr(*sample.HttpclientTimeoutConnect + 1)
		result.InsecureForkWanted = !sample.InsecureForkWanted
		result.InsecureSetuidWanted = !sample.InsecureSetuidWanted
		result.LimitedQuic = !sample.LimitedQuic
		result.MasterWorker = !sample.MasterWorker
		result.MaxSpreadChecks = Ptr(*sample.MaxSpreadChecks + 1)
		result.Maxcompcpuusage = sample.Maxcompcpuusage + 1
		result.Maxcomprate = sample.Maxcomprate + 1
		result.Maxconn = sample.Maxconn + 1
		result.Maxconnrate = sample.Maxconnrate + 1
		result.Maxpipes = sample.Maxpipes + 1
		result.Maxsessrate = sample.Maxsessrate + 1
		result.Maxsslconn = sample.Maxsslconn + 1
		result.Maxsslrate = sample.Maxsslrate + 1
		result.Maxzlibmem = sample.Maxzlibmem + 1
		result.MworkerMaxReloads = Ptr(*sample.MworkerMaxReloads + 1)
		result.Nbproc = sample.Nbproc + 1
		result.Nbthread = sample.Nbthread + 1
		result.NoQuic = !sample.NoQuic
		result.Noepoll = !sample.Noepoll
		result.Noevports = !sample.Noevports
		result.Nogetaddrinfo = !sample.Nogetaddrinfo
		result.Nokqueue = !sample.Nokqueue
		result.Nopoll = !sample.Nopoll
		result.Noreuseport = !sample.Noreuseport
		result.Nosplice = !sample.Nosplice
		result.Pp2NeverSendLocal = !sample.Pp2NeverSendLocal
		result.PreallocFd = !sample.PreallocFd
		result.Quiet = !sample.Quiet
		result.SetDumpable = !sample.SetDumpable
		result.SpreadChecks = sample.SpreadChecks + 1
		result.SslSkipSelfIssuedCa = !sample.SslSkipSelfIssuedCa
		result.StatsMaxconn = Ptr(*sample.StatsMaxconn + 1)
		result.StatsTimeout = Ptr(*sample.StatsTimeout + 1)
		result.StrictLimits = !sample.StrictLimits
		result.ThreadGroups = sample.ThreadGroups + 1
		result.TuneSslDefaultDhParam = sample.TuneSslDefaultDhParam + 1
		result.UID = sample.UID + 1
		result.Ulimitn = sample.Ulimitn + 1
		result.ZeroWarning = !sample.ZeroWarning
		samples = append(samples, struct {
			a, b Global
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		if len(result) != 113 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected Global to be different in 113 cases, but it is not (%d) %s %s", len(result), a, b)
		}
	}
}

func TestCPUMapEqual(t *testing.T) {
	samples := []struct {
		a, b CPUMap
	}{}
	for i := 0; i < 2; i++ {
		var sample CPUMap
		var result CPUMap
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		byteJSON, err := json.Marshal(sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = json.Unmarshal(byteJSON, &result)
		if err != nil {
			t.Errorf(err.Error())
		}

		samples = append(samples, struct {
			a, b CPUMap
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Equal(sample.b)
		if !result {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected CPUMap to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestCPUMapEqualFalse(t *testing.T) {
	samples := []struct {
		a, b CPUMap
	}{}
	for i := 0; i < 2; i++ {
		var sample CPUMap
		var result CPUMap
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		samples = append(samples, struct {
			a, b CPUMap
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Equal(sample.b)
		if result {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected CPUMap to be different, but it is not %s %s", a, b)
		}
	}
}

func TestCPUMapDiff(t *testing.T) {
	samples := []struct {
		a, b CPUMap
	}{}
	for i := 0; i < 2; i++ {
		var sample CPUMap
		var result CPUMap
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		byteJSON, err := json.Marshal(sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = json.Unmarshal(byteJSON, &result)
		if err != nil {
			t.Errorf(err.Error())
		}

		samples = append(samples, struct {
			a, b CPUMap
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		if len(result) != 0 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected CPUMap to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestCPUMapDiffFalse(t *testing.T) {
	samples := []struct {
		a, b CPUMap
	}{}
	for i := 0; i < 2; i++ {
		var sample CPUMap
		var result CPUMap
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		samples = append(samples, struct {
			a, b CPUMap
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		if len(result) != 2 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected CPUMap to be different in 2 cases, but it is not (%d) %s %s", len(result), a, b)
		}
	}
}

func TestGlobalDefaultPathEqual(t *testing.T) {
	samples := []struct {
		a, b GlobalDefaultPath
	}{}
	for i := 0; i < 2; i++ {
		var sample GlobalDefaultPath
		var result GlobalDefaultPath
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		byteJSON, err := json.Marshal(sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = json.Unmarshal(byteJSON, &result)
		if err != nil {
			t.Errorf(err.Error())
		}

		samples = append(samples, struct {
			a, b GlobalDefaultPath
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Equal(sample.b)
		if !result {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected GlobalDefaultPath to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestGlobalDefaultPathEqualFalse(t *testing.T) {
	samples := []struct {
		a, b GlobalDefaultPath
	}{}
	for i := 0; i < 2; i++ {
		var sample GlobalDefaultPath
		var result GlobalDefaultPath
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		samples = append(samples, struct {
			a, b GlobalDefaultPath
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Equal(sample.b)
		if result {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected GlobalDefaultPath to be different, but it is not %s %s", a, b)
		}
	}
}

func TestGlobalDefaultPathDiff(t *testing.T) {
	samples := []struct {
		a, b GlobalDefaultPath
	}{}
	for i := 0; i < 2; i++ {
		var sample GlobalDefaultPath
		var result GlobalDefaultPath
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		byteJSON, err := json.Marshal(sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = json.Unmarshal(byteJSON, &result)
		if err != nil {
			t.Errorf(err.Error())
		}

		samples = append(samples, struct {
			a, b GlobalDefaultPath
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		if len(result) != 0 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected GlobalDefaultPath to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestGlobalDefaultPathDiffFalse(t *testing.T) {
	samples := []struct {
		a, b GlobalDefaultPath
	}{}
	for i := 0; i < 2; i++ {
		var sample GlobalDefaultPath
		var result GlobalDefaultPath
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		samples = append(samples, struct {
			a, b GlobalDefaultPath
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		if len(result) != 2 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected GlobalDefaultPath to be different in 2 cases, but it is not (%d) %s %s", len(result), a, b)
		}
	}
}

func TestGlobalDeviceAtlasOptionsEqual(t *testing.T) {
	samples := []struct {
		a, b GlobalDeviceAtlasOptions
	}{}
	for i := 0; i < 2; i++ {
		var sample GlobalDeviceAtlasOptions
		var result GlobalDeviceAtlasOptions
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		byteJSON, err := json.Marshal(sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = json.Unmarshal(byteJSON, &result)
		if err != nil {
			t.Errorf(err.Error())
		}

		samples = append(samples, struct {
			a, b GlobalDeviceAtlasOptions
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Equal(sample.b)
		if !result {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected GlobalDeviceAtlasOptions to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestGlobalDeviceAtlasOptionsEqualFalse(t *testing.T) {
	samples := []struct {
		a, b GlobalDeviceAtlasOptions
	}{}
	for i := 0; i < 2; i++ {
		var sample GlobalDeviceAtlasOptions
		var result GlobalDeviceAtlasOptions
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		samples = append(samples, struct {
			a, b GlobalDeviceAtlasOptions
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Equal(sample.b)
		if result {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected GlobalDeviceAtlasOptions to be different, but it is not %s %s", a, b)
		}
	}
}

func TestGlobalDeviceAtlasOptionsDiff(t *testing.T) {
	samples := []struct {
		a, b GlobalDeviceAtlasOptions
	}{}
	for i := 0; i < 2; i++ {
		var sample GlobalDeviceAtlasOptions
		var result GlobalDeviceAtlasOptions
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		byteJSON, err := json.Marshal(sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = json.Unmarshal(byteJSON, &result)
		if err != nil {
			t.Errorf(err.Error())
		}

		samples = append(samples, struct {
			a, b GlobalDeviceAtlasOptions
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		if len(result) != 0 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected GlobalDeviceAtlasOptions to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestGlobalDeviceAtlasOptionsDiffFalse(t *testing.T) {
	samples := []struct {
		a, b GlobalDeviceAtlasOptions
	}{}
	for i := 0; i < 2; i++ {
		var sample GlobalDeviceAtlasOptions
		var result GlobalDeviceAtlasOptions
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		samples = append(samples, struct {
			a, b GlobalDeviceAtlasOptions
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		if len(result) != 4 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected GlobalDeviceAtlasOptions to be different in 4 cases, but it is not (%d) %s %s", len(result), a, b)
		}
	}
}

func TestGlobalFiftyOneDegreesOptionsEqual(t *testing.T) {
	samples := []struct {
		a, b GlobalFiftyOneDegreesOptions
	}{}
	for i := 0; i < 2; i++ {
		var sample GlobalFiftyOneDegreesOptions
		var result GlobalFiftyOneDegreesOptions
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		byteJSON, err := json.Marshal(sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = json.Unmarshal(byteJSON, &result)
		if err != nil {
			t.Errorf(err.Error())
		}

		samples = append(samples, struct {
			a, b GlobalFiftyOneDegreesOptions
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Equal(sample.b)
		if !result {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected GlobalFiftyOneDegreesOptions to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestGlobalFiftyOneDegreesOptionsEqualFalse(t *testing.T) {
	samples := []struct {
		a, b GlobalFiftyOneDegreesOptions
	}{}
	for i := 0; i < 2; i++ {
		var sample GlobalFiftyOneDegreesOptions
		var result GlobalFiftyOneDegreesOptions
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		result.CacheSize = sample.CacheSize + 1
		samples = append(samples, struct {
			a, b GlobalFiftyOneDegreesOptions
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Equal(sample.b)
		if result {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected GlobalFiftyOneDegreesOptions to be different, but it is not %s %s", a, b)
		}
	}
}

func TestGlobalFiftyOneDegreesOptionsDiff(t *testing.T) {
	samples := []struct {
		a, b GlobalFiftyOneDegreesOptions
	}{}
	for i := 0; i < 2; i++ {
		var sample GlobalFiftyOneDegreesOptions
		var result GlobalFiftyOneDegreesOptions
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		byteJSON, err := json.Marshal(sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = json.Unmarshal(byteJSON, &result)
		if err != nil {
			t.Errorf(err.Error())
		}

		samples = append(samples, struct {
			a, b GlobalFiftyOneDegreesOptions
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		if len(result) != 0 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected GlobalFiftyOneDegreesOptions to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestGlobalFiftyOneDegreesOptionsDiffFalse(t *testing.T) {
	samples := []struct {
		a, b GlobalFiftyOneDegreesOptions
	}{}
	for i := 0; i < 2; i++ {
		var sample GlobalFiftyOneDegreesOptions
		var result GlobalFiftyOneDegreesOptions
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		result.CacheSize = sample.CacheSize + 1
		samples = append(samples, struct {
			a, b GlobalFiftyOneDegreesOptions
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		if len(result) != 4 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected GlobalFiftyOneDegreesOptions to be different in 4 cases, but it is not (%d) %s %s", len(result), a, b)
		}
	}
}

func TestH1CaseAdjustEqual(t *testing.T) {
	samples := []struct {
		a, b H1CaseAdjust
	}{}
	for i := 0; i < 2; i++ {
		var sample H1CaseAdjust
		var result H1CaseAdjust
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		byteJSON, err := json.Marshal(sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = json.Unmarshal(byteJSON, &result)
		if err != nil {
			t.Errorf(err.Error())
		}

		samples = append(samples, struct {
			a, b H1CaseAdjust
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Equal(sample.b)
		if !result {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected H1CaseAdjust to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestH1CaseAdjustEqualFalse(t *testing.T) {
	samples := []struct {
		a, b H1CaseAdjust
	}{}
	for i := 0; i < 2; i++ {
		var sample H1CaseAdjust
		var result H1CaseAdjust
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		samples = append(samples, struct {
			a, b H1CaseAdjust
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Equal(sample.b)
		if result {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected H1CaseAdjust to be different, but it is not %s %s", a, b)
		}
	}
}

func TestH1CaseAdjustDiff(t *testing.T) {
	samples := []struct {
		a, b H1CaseAdjust
	}{}
	for i := 0; i < 2; i++ {
		var sample H1CaseAdjust
		var result H1CaseAdjust
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		byteJSON, err := json.Marshal(sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = json.Unmarshal(byteJSON, &result)
		if err != nil {
			t.Errorf(err.Error())
		}

		samples = append(samples, struct {
			a, b H1CaseAdjust
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		if len(result) != 0 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected H1CaseAdjust to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestH1CaseAdjustDiffFalse(t *testing.T) {
	samples := []struct {
		a, b H1CaseAdjust
	}{}
	for i := 0; i < 2; i++ {
		var sample H1CaseAdjust
		var result H1CaseAdjust
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		samples = append(samples, struct {
			a, b H1CaseAdjust
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		if len(result) != 2 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected H1CaseAdjust to be different in 2 cases, but it is not (%d) %s %s", len(result), a, b)
		}
	}
}

func TestGlobalLogSendHostnameEqual(t *testing.T) {
	samples := []struct {
		a, b GlobalLogSendHostname
	}{}
	for i := 0; i < 2; i++ {
		var sample GlobalLogSendHostname
		var result GlobalLogSendHostname
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		byteJSON, err := json.Marshal(sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = json.Unmarshal(byteJSON, &result)
		if err != nil {
			t.Errorf(err.Error())
		}

		samples = append(samples, struct {
			a, b GlobalLogSendHostname
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Equal(sample.b)
		if !result {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected GlobalLogSendHostname to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestGlobalLogSendHostnameEqualFalse(t *testing.T) {
	samples := []struct {
		a, b GlobalLogSendHostname
	}{}
	for i := 0; i < 2; i++ {
		var sample GlobalLogSendHostname
		var result GlobalLogSendHostname
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		samples = append(samples, struct {
			a, b GlobalLogSendHostname
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Equal(sample.b)
		if result {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected GlobalLogSendHostname to be different, but it is not %s %s", a, b)
		}
	}
}

func TestGlobalLogSendHostnameDiff(t *testing.T) {
	samples := []struct {
		a, b GlobalLogSendHostname
	}{}
	for i := 0; i < 2; i++ {
		var sample GlobalLogSendHostname
		var result GlobalLogSendHostname
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		byteJSON, err := json.Marshal(sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = json.Unmarshal(byteJSON, &result)
		if err != nil {
			t.Errorf(err.Error())
		}

		samples = append(samples, struct {
			a, b GlobalLogSendHostname
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		if len(result) != 0 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected GlobalLogSendHostname to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestGlobalLogSendHostnameDiffFalse(t *testing.T) {
	samples := []struct {
		a, b GlobalLogSendHostname
	}{}
	for i := 0; i < 2; i++ {
		var sample GlobalLogSendHostname
		var result GlobalLogSendHostname
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		samples = append(samples, struct {
			a, b GlobalLogSendHostname
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		if len(result) != 2 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected GlobalLogSendHostname to be different in 2 cases, but it is not (%d) %s %s", len(result), a, b)
		}
	}
}

func TestLuaLoadEqual(t *testing.T) {
	samples := []struct {
		a, b LuaLoad
	}{}
	for i := 0; i < 2; i++ {
		var sample LuaLoad
		var result LuaLoad
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		byteJSON, err := json.Marshal(sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = json.Unmarshal(byteJSON, &result)
		if err != nil {
			t.Errorf(err.Error())
		}

		samples = append(samples, struct {
			a, b LuaLoad
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Equal(sample.b)
		if !result {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected LuaLoad to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestLuaLoadEqualFalse(t *testing.T) {
	samples := []struct {
		a, b LuaLoad
	}{}
	for i := 0; i < 2; i++ {
		var sample LuaLoad
		var result LuaLoad
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		samples = append(samples, struct {
			a, b LuaLoad
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Equal(sample.b)
		if result {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected LuaLoad to be different, but it is not %s %s", a, b)
		}
	}
}

func TestLuaLoadDiff(t *testing.T) {
	samples := []struct {
		a, b LuaLoad
	}{}
	for i := 0; i < 2; i++ {
		var sample LuaLoad
		var result LuaLoad
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		byteJSON, err := json.Marshal(sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = json.Unmarshal(byteJSON, &result)
		if err != nil {
			t.Errorf(err.Error())
		}

		samples = append(samples, struct {
			a, b LuaLoad
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		if len(result) != 0 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected LuaLoad to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestLuaLoadDiffFalse(t *testing.T) {
	samples := []struct {
		a, b LuaLoad
	}{}
	for i := 0; i < 2; i++ {
		var sample LuaLoad
		var result LuaLoad
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		samples = append(samples, struct {
			a, b LuaLoad
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		if len(result) != 1 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected LuaLoad to be different in 1 cases, but it is not (%d) %s %s", len(result), a, b)
		}
	}
}

func TestLuaPrependPathEqual(t *testing.T) {
	samples := []struct {
		a, b LuaPrependPath
	}{}
	for i := 0; i < 2; i++ {
		var sample LuaPrependPath
		var result LuaPrependPath
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		byteJSON, err := json.Marshal(sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = json.Unmarshal(byteJSON, &result)
		if err != nil {
			t.Errorf(err.Error())
		}

		samples = append(samples, struct {
			a, b LuaPrependPath
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Equal(sample.b)
		if !result {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected LuaPrependPath to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestLuaPrependPathEqualFalse(t *testing.T) {
	samples := []struct {
		a, b LuaPrependPath
	}{}
	for i := 0; i < 2; i++ {
		var sample LuaPrependPath
		var result LuaPrependPath
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		samples = append(samples, struct {
			a, b LuaPrependPath
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Equal(sample.b)
		if result {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected LuaPrependPath to be different, but it is not %s %s", a, b)
		}
	}
}

func TestLuaPrependPathDiff(t *testing.T) {
	samples := []struct {
		a, b LuaPrependPath
	}{}
	for i := 0; i < 2; i++ {
		var sample LuaPrependPath
		var result LuaPrependPath
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		byteJSON, err := json.Marshal(sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = json.Unmarshal(byteJSON, &result)
		if err != nil {
			t.Errorf(err.Error())
		}

		samples = append(samples, struct {
			a, b LuaPrependPath
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		if len(result) != 0 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected LuaPrependPath to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestLuaPrependPathDiffFalse(t *testing.T) {
	samples := []struct {
		a, b LuaPrependPath
	}{}
	for i := 0; i < 2; i++ {
		var sample LuaPrependPath
		var result LuaPrependPath
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		samples = append(samples, struct {
			a, b LuaPrependPath
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		if len(result) != 2 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected LuaPrependPath to be different in 2 cases, but it is not (%d) %s %s", len(result), a, b)
		}
	}
}

func TestPresetEnvEqual(t *testing.T) {
	samples := []struct {
		a, b PresetEnv
	}{}
	for i := 0; i < 2; i++ {
		var sample PresetEnv
		var result PresetEnv
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		byteJSON, err := json.Marshal(sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = json.Unmarshal(byteJSON, &result)
		if err != nil {
			t.Errorf(err.Error())
		}

		samples = append(samples, struct {
			a, b PresetEnv
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Equal(sample.b)
		if !result {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected PresetEnv to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestPresetEnvEqualFalse(t *testing.T) {
	samples := []struct {
		a, b PresetEnv
	}{}
	for i := 0; i < 2; i++ {
		var sample PresetEnv
		var result PresetEnv
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		samples = append(samples, struct {
			a, b PresetEnv
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Equal(sample.b)
		if result {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected PresetEnv to be different, but it is not %s %s", a, b)
		}
	}
}

func TestPresetEnvDiff(t *testing.T) {
	samples := []struct {
		a, b PresetEnv
	}{}
	for i := 0; i < 2; i++ {
		var sample PresetEnv
		var result PresetEnv
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		byteJSON, err := json.Marshal(sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = json.Unmarshal(byteJSON, &result)
		if err != nil {
			t.Errorf(err.Error())
		}

		samples = append(samples, struct {
			a, b PresetEnv
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		if len(result) != 0 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected PresetEnv to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestPresetEnvDiffFalse(t *testing.T) {
	samples := []struct {
		a, b PresetEnv
	}{}
	for i := 0; i < 2; i++ {
		var sample PresetEnv
		var result PresetEnv
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		samples = append(samples, struct {
			a, b PresetEnv
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		if len(result) != 2 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected PresetEnv to be different in 2 cases, but it is not (%d) %s %s", len(result), a, b)
		}
	}
}

func TestRuntimeAPIEqual(t *testing.T) {
	samples := []struct {
		a, b RuntimeAPI
	}{}
	for i := 0; i < 2; i++ {
		var sample RuntimeAPI
		var result RuntimeAPI
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		byteJSON, err := json.Marshal(sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = json.Unmarshal(byteJSON, &result)
		if err != nil {
			t.Errorf(err.Error())
		}

		samples = append(samples, struct {
			a, b RuntimeAPI
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Equal(sample.b)
		if !result {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected RuntimeAPI to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestRuntimeAPIEqualFalse(t *testing.T) {
	samples := []struct {
		a, b RuntimeAPI
	}{}
	for i := 0; i < 2; i++ {
		var sample RuntimeAPI
		var result RuntimeAPI
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		samples = append(samples, struct {
			a, b RuntimeAPI
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Equal(sample.b)
		if result {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected RuntimeAPI to be different, but it is not %s %s", a, b)
		}
	}
}

func TestRuntimeAPIDiff(t *testing.T) {
	samples := []struct {
		a, b RuntimeAPI
	}{}
	for i := 0; i < 2; i++ {
		var sample RuntimeAPI
		var result RuntimeAPI
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		byteJSON, err := json.Marshal(sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = json.Unmarshal(byteJSON, &result)
		if err != nil {
			t.Errorf(err.Error())
		}

		samples = append(samples, struct {
			a, b RuntimeAPI
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		if len(result) != 0 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected RuntimeAPI to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestRuntimeAPIDiffFalse(t *testing.T) {
	samples := []struct {
		a, b RuntimeAPI
	}{}
	for i := 0; i < 2; i++ {
		var sample RuntimeAPI
		var result RuntimeAPI
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		samples = append(samples, struct {
			a, b RuntimeAPI
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		if len(result) != 2 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected RuntimeAPI to be different in 2 cases, but it is not (%d) %s %s", len(result), a, b)
		}
	}
}

func TestSetVarFmtEqual(t *testing.T) {
	samples := []struct {
		a, b SetVarFmt
	}{}
	for i := 0; i < 2; i++ {
		var sample SetVarFmt
		var result SetVarFmt
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		byteJSON, err := json.Marshal(sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = json.Unmarshal(byteJSON, &result)
		if err != nil {
			t.Errorf(err.Error())
		}

		samples = append(samples, struct {
			a, b SetVarFmt
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Equal(sample.b)
		if !result {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected SetVarFmt to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestSetVarFmtEqualFalse(t *testing.T) {
	samples := []struct {
		a, b SetVarFmt
	}{}
	for i := 0; i < 2; i++ {
		var sample SetVarFmt
		var result SetVarFmt
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		samples = append(samples, struct {
			a, b SetVarFmt
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Equal(sample.b)
		if result {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected SetVarFmt to be different, but it is not %s %s", a, b)
		}
	}
}

func TestSetVarFmtDiff(t *testing.T) {
	samples := []struct {
		a, b SetVarFmt
	}{}
	for i := 0; i < 2; i++ {
		var sample SetVarFmt
		var result SetVarFmt
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		byteJSON, err := json.Marshal(sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = json.Unmarshal(byteJSON, &result)
		if err != nil {
			t.Errorf(err.Error())
		}

		samples = append(samples, struct {
			a, b SetVarFmt
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		if len(result) != 0 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected SetVarFmt to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestSetVarFmtDiffFalse(t *testing.T) {
	samples := []struct {
		a, b SetVarFmt
	}{}
	for i := 0; i < 2; i++ {
		var sample SetVarFmt
		var result SetVarFmt
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		samples = append(samples, struct {
			a, b SetVarFmt
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		if len(result) != 2 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected SetVarFmt to be different in 2 cases, but it is not (%d) %s %s", len(result), a, b)
		}
	}
}

func TestSetVarEqual(t *testing.T) {
	samples := []struct {
		a, b SetVar
	}{}
	for i := 0; i < 2; i++ {
		var sample SetVar
		var result SetVar
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		byteJSON, err := json.Marshal(sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = json.Unmarshal(byteJSON, &result)
		if err != nil {
			t.Errorf(err.Error())
		}

		samples = append(samples, struct {
			a, b SetVar
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Equal(sample.b)
		if !result {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected SetVar to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestSetVarEqualFalse(t *testing.T) {
	samples := []struct {
		a, b SetVar
	}{}
	for i := 0; i < 2; i++ {
		var sample SetVar
		var result SetVar
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		samples = append(samples, struct {
			a, b SetVar
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Equal(sample.b)
		if result {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected SetVar to be different, but it is not %s %s", a, b)
		}
	}
}

func TestSetVarDiff(t *testing.T) {
	samples := []struct {
		a, b SetVar
	}{}
	for i := 0; i < 2; i++ {
		var sample SetVar
		var result SetVar
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		byteJSON, err := json.Marshal(sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = json.Unmarshal(byteJSON, &result)
		if err != nil {
			t.Errorf(err.Error())
		}

		samples = append(samples, struct {
			a, b SetVar
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		if len(result) != 0 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected SetVar to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestSetVarDiffFalse(t *testing.T) {
	samples := []struct {
		a, b SetVar
	}{}
	for i := 0; i < 2; i++ {
		var sample SetVar
		var result SetVar
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		samples = append(samples, struct {
			a, b SetVar
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		if len(result) != 2 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected SetVar to be different in 2 cases, but it is not (%d) %s %s", len(result), a, b)
		}
	}
}

func TestSetEnvEqual(t *testing.T) {
	samples := []struct {
		a, b SetEnv
	}{}
	for i := 0; i < 2; i++ {
		var sample SetEnv
		var result SetEnv
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		byteJSON, err := json.Marshal(sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = json.Unmarshal(byteJSON, &result)
		if err != nil {
			t.Errorf(err.Error())
		}

		samples = append(samples, struct {
			a, b SetEnv
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Equal(sample.b)
		if !result {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected SetEnv to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestSetEnvEqualFalse(t *testing.T) {
	samples := []struct {
		a, b SetEnv
	}{}
	for i := 0; i < 2; i++ {
		var sample SetEnv
		var result SetEnv
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		samples = append(samples, struct {
			a, b SetEnv
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Equal(sample.b)
		if result {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected SetEnv to be different, but it is not %s %s", a, b)
		}
	}
}

func TestSetEnvDiff(t *testing.T) {
	samples := []struct {
		a, b SetEnv
	}{}
	for i := 0; i < 2; i++ {
		var sample SetEnv
		var result SetEnv
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		byteJSON, err := json.Marshal(sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = json.Unmarshal(byteJSON, &result)
		if err != nil {
			t.Errorf(err.Error())
		}

		samples = append(samples, struct {
			a, b SetEnv
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		if len(result) != 0 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected SetEnv to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestSetEnvDiffFalse(t *testing.T) {
	samples := []struct {
		a, b SetEnv
	}{}
	for i := 0; i < 2; i++ {
		var sample SetEnv
		var result SetEnv
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		samples = append(samples, struct {
			a, b SetEnv
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		if len(result) != 2 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected SetEnv to be different in 2 cases, but it is not (%d) %s %s", len(result), a, b)
		}
	}
}

func TestSslEngineEqual(t *testing.T) {
	samples := []struct {
		a, b SslEngine
	}{}
	for i := 0; i < 2; i++ {
		var sample SslEngine
		var result SslEngine
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		byteJSON, err := json.Marshal(sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = json.Unmarshal(byteJSON, &result)
		if err != nil {
			t.Errorf(err.Error())
		}

		samples = append(samples, struct {
			a, b SslEngine
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Equal(sample.b)
		if !result {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected SslEngine to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestSslEngineEqualFalse(t *testing.T) {
	samples := []struct {
		a, b SslEngine
	}{}
	for i := 0; i < 2; i++ {
		var sample SslEngine
		var result SslEngine
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		samples = append(samples, struct {
			a, b SslEngine
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Equal(sample.b)
		if result {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected SslEngine to be different, but it is not %s %s", a, b)
		}
	}
}

func TestSslEngineDiff(t *testing.T) {
	samples := []struct {
		a, b SslEngine
	}{}
	for i := 0; i < 2; i++ {
		var sample SslEngine
		var result SslEngine
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		byteJSON, err := json.Marshal(sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = json.Unmarshal(byteJSON, &result)
		if err != nil {
			t.Errorf(err.Error())
		}

		samples = append(samples, struct {
			a, b SslEngine
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		if len(result) != 0 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected SslEngine to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestSslEngineDiffFalse(t *testing.T) {
	samples := []struct {
		a, b SslEngine
	}{}
	for i := 0; i < 2; i++ {
		var sample SslEngine
		var result SslEngine
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		samples = append(samples, struct {
			a, b SslEngine
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		if len(result) != 2 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected SslEngine to be different in 2 cases, but it is not (%d) %s %s", len(result), a, b)
		}
	}
}

func TestThreadGroupEqual(t *testing.T) {
	samples := []struct {
		a, b ThreadGroup
	}{}
	for i := 0; i < 2; i++ {
		var sample ThreadGroup
		var result ThreadGroup
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		byteJSON, err := json.Marshal(sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = json.Unmarshal(byteJSON, &result)
		if err != nil {
			t.Errorf(err.Error())
		}

		samples = append(samples, struct {
			a, b ThreadGroup
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Equal(sample.b)
		if !result {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected ThreadGroup to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestThreadGroupEqualFalse(t *testing.T) {
	samples := []struct {
		a, b ThreadGroup
	}{}
	for i := 0; i < 2; i++ {
		var sample ThreadGroup
		var result ThreadGroup
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		samples = append(samples, struct {
			a, b ThreadGroup
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Equal(sample.b)
		if result {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected ThreadGroup to be different, but it is not %s %s", a, b)
		}
	}
}

func TestThreadGroupDiff(t *testing.T) {
	samples := []struct {
		a, b ThreadGroup
	}{}
	for i := 0; i < 2; i++ {
		var sample ThreadGroup
		var result ThreadGroup
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		byteJSON, err := json.Marshal(sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = json.Unmarshal(byteJSON, &result)
		if err != nil {
			t.Errorf(err.Error())
		}

		samples = append(samples, struct {
			a, b ThreadGroup
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		if len(result) != 0 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected ThreadGroup to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestThreadGroupDiffFalse(t *testing.T) {
	samples := []struct {
		a, b ThreadGroup
	}{}
	for i := 0; i < 2; i++ {
		var sample ThreadGroup
		var result ThreadGroup
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		samples = append(samples, struct {
			a, b ThreadGroup
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		if len(result) != 2 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected ThreadGroup to be different in 2 cases, but it is not (%d) %s %s", len(result), a, b)
		}
	}
}

func TestGlobalTuneOptionsEqual(t *testing.T) {
	samples := []struct {
		a, b GlobalTuneOptions
	}{}
	for i := 0; i < 2; i++ {
		var sample GlobalTuneOptions
		var result GlobalTuneOptions
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		byteJSON, err := json.Marshal(sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = json.Unmarshal(byteJSON, &result)
		if err != nil {
			t.Errorf(err.Error())
		}

		samples = append(samples, struct {
			a, b GlobalTuneOptions
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Equal(sample.b)
		if !result {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected GlobalTuneOptions to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestGlobalTuneOptionsEqualFalse(t *testing.T) {
	samples := []struct {
		a, b GlobalTuneOptions
	}{}
	for i := 0; i < 2; i++ {
		var sample GlobalTuneOptions
		var result GlobalTuneOptions
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		result.BuffersLimit = Ptr(*sample.BuffersLimit + 1)
		result.BuffersReserve = sample.BuffersReserve + 1
		result.Bufsize = sample.Bufsize + 1
		result.CompMaxlevel = sample.CompMaxlevel + 1
		result.DisableZeroCopyForwarding = !sample.DisableZeroCopyForwarding
		result.FailAlloc = !sample.FailAlloc
		result.H2BeInitialWindowSize = sample.H2BeInitialWindowSize + 1
		result.H2BeMaxConcurrentStreams = sample.H2BeMaxConcurrentStreams + 1
		result.H2FeInitialWindowSize = sample.H2FeInitialWindowSize + 1
		result.H2FeMaxConcurrentStreams = sample.H2FeMaxConcurrentStreams + 1
		result.H2HeaderTableSize = sample.H2HeaderTableSize + 1
		result.H2InitialWindowSize = Ptr(*sample.H2InitialWindowSize + 1)
		result.H2MaxConcurrentStreams = sample.H2MaxConcurrentStreams + 1
		result.H2MaxFrameSize = sample.H2MaxFrameSize + 1
		result.HTTPCookielen = sample.HTTPCookielen + 1
		result.HTTPLogurilen = sample.HTTPLogurilen + 1
		result.HTTPMaxhdr = sample.HTTPMaxhdr + 1
		result.Idletimer = Ptr(*sample.Idletimer + 1)
		result.LuaBurstTimeout = Ptr(*sample.LuaBurstTimeout + 1)
		result.LuaForcedYield = sample.LuaForcedYield + 1
		result.LuaMaxmem = !sample.LuaMaxmem
		result.LuaServiceTimeout = Ptr(*sample.LuaServiceTimeout + 1)
		result.LuaSessionTimeout = Ptr(*sample.LuaSessionTimeout + 1)
		result.LuaTaskTimeout = Ptr(*sample.LuaTaskTimeout + 1)
		result.MaxChecksPerThread = Ptr(*sample.MaxChecksPerThread + 1)
		result.Maxaccept = sample.Maxaccept + 1
		result.Maxpollevents = sample.Maxpollevents + 1
		result.Maxrewrite = sample.Maxrewrite + 1
		result.MemoryHotSize = Ptr(*sample.MemoryHotSize + 1)
		result.PatternCacheSize = Ptr(*sample.PatternCacheSize + 1)
		result.PeersMaxUpdatesAtOnce = sample.PeersMaxUpdatesAtOnce + 1
		result.Pipesize = sample.Pipesize + 1
		result.PoolHighFdRatio = sample.PoolHighFdRatio + 1
		result.PoolLowFdRatio = sample.PoolLowFdRatio + 1
		result.QuicFrontendConnTxBuffersLimit = Ptr(*sample.QuicFrontendConnTxBuffersLimit + 1)
		result.QuicFrontendMaxIdleTimeout = Ptr(*sample.QuicFrontendMaxIdleTimeout + 1)
		result.QuicFrontendMaxStreamsBidi = Ptr(*sample.QuicFrontendMaxStreamsBidi + 1)
		result.QuicMaxFrameLoss = Ptr(*sample.QuicMaxFrameLoss + 1)
		result.QuicRetryThreshold = Ptr(*sample.QuicRetryThreshold + 1)
		result.RcvbufBackend = Ptr(*sample.RcvbufBackend + 1)
		result.RcvbufClient = Ptr(*sample.RcvbufClient + 1)
		result.RcvbufFrontend = Ptr(*sample.RcvbufFrontend + 1)
		result.RcvbufServer = Ptr(*sample.RcvbufServer + 1)
		result.RecvEnough = sample.RecvEnough + 1
		result.RunqueueDepth = sample.RunqueueDepth + 1
		result.SndbufBackend = Ptr(*sample.SndbufBackend + 1)
		result.SndbufClient = Ptr(*sample.SndbufClient + 1)
		result.SndbufFrontend = Ptr(*sample.SndbufFrontend + 1)
		result.SndbufServer = Ptr(*sample.SndbufServer + 1)
		result.SslCachesize = Ptr(*sample.SslCachesize + 1)
		result.SslCaptureBufferSize = Ptr(*sample.SslCaptureBufferSize + 1)
		result.SslCtxCacheSize = sample.SslCtxCacheSize + 1
		result.SslDefaultDhParam = sample.SslDefaultDhParam + 1
		result.SslForcePrivateCache = !sample.SslForcePrivateCache
		result.SslLifetime = Ptr(*sample.SslLifetime + 1)
		result.SslMaxrecord = Ptr(*sample.SslMaxrecord + 1)
		result.SslOcspUpdateMaxDelay = Ptr(*sample.SslOcspUpdateMaxDelay + 1)
		result.SslOcspUpdateMinDelay = Ptr(*sample.SslOcspUpdateMinDelay + 1)
		result.StickCounters = Ptr(*sample.StickCounters + 1)
		result.VarsGlobalMaxSize = Ptr(*sample.VarsGlobalMaxSize + 1)
		result.VarsProcMaxSize = Ptr(*sample.VarsProcMaxSize + 1)
		result.VarsReqresMaxSize = Ptr(*sample.VarsReqresMaxSize + 1)
		result.VarsSessMaxSize = Ptr(*sample.VarsSessMaxSize + 1)
		result.VarsTxnMaxSize = Ptr(*sample.VarsTxnMaxSize + 1)
		result.ZlibMemlevel = sample.ZlibMemlevel + 1
		result.ZlibWindowsize = sample.ZlibWindowsize + 1
		samples = append(samples, struct {
			a, b GlobalTuneOptions
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Equal(sample.b)
		if result {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected GlobalTuneOptions to be different, but it is not %s %s", a, b)
		}
	}
}

func TestGlobalTuneOptionsDiff(t *testing.T) {
	samples := []struct {
		a, b GlobalTuneOptions
	}{}
	for i := 0; i < 2; i++ {
		var sample GlobalTuneOptions
		var result GlobalTuneOptions
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		byteJSON, err := json.Marshal(sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = json.Unmarshal(byteJSON, &result)
		if err != nil {
			t.Errorf(err.Error())
		}

		samples = append(samples, struct {
			a, b GlobalTuneOptions
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		if len(result) != 0 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected GlobalTuneOptions to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestGlobalTuneOptionsDiffFalse(t *testing.T) {
	samples := []struct {
		a, b GlobalTuneOptions
	}{}
	for i := 0; i < 2; i++ {
		var sample GlobalTuneOptions
		var result GlobalTuneOptions
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		result.BuffersLimit = Ptr(*sample.BuffersLimit + 1)
		result.BuffersReserve = sample.BuffersReserve + 1
		result.Bufsize = sample.Bufsize + 1
		result.CompMaxlevel = sample.CompMaxlevel + 1
		result.DisableZeroCopyForwarding = !sample.DisableZeroCopyForwarding
		result.FailAlloc = !sample.FailAlloc
		result.H2BeInitialWindowSize = sample.H2BeInitialWindowSize + 1
		result.H2BeMaxConcurrentStreams = sample.H2BeMaxConcurrentStreams + 1
		result.H2FeInitialWindowSize = sample.H2FeInitialWindowSize + 1
		result.H2FeMaxConcurrentStreams = sample.H2FeMaxConcurrentStreams + 1
		result.H2HeaderTableSize = sample.H2HeaderTableSize + 1
		result.H2InitialWindowSize = Ptr(*sample.H2InitialWindowSize + 1)
		result.H2MaxConcurrentStreams = sample.H2MaxConcurrentStreams + 1
		result.H2MaxFrameSize = sample.H2MaxFrameSize + 1
		result.HTTPCookielen = sample.HTTPCookielen + 1
		result.HTTPLogurilen = sample.HTTPLogurilen + 1
		result.HTTPMaxhdr = sample.HTTPMaxhdr + 1
		result.Idletimer = Ptr(*sample.Idletimer + 1)
		result.LuaBurstTimeout = Ptr(*sample.LuaBurstTimeout + 1)
		result.LuaForcedYield = sample.LuaForcedYield + 1
		result.LuaMaxmem = !sample.LuaMaxmem
		result.LuaServiceTimeout = Ptr(*sample.LuaServiceTimeout + 1)
		result.LuaSessionTimeout = Ptr(*sample.LuaSessionTimeout + 1)
		result.LuaTaskTimeout = Ptr(*sample.LuaTaskTimeout + 1)
		result.MaxChecksPerThread = Ptr(*sample.MaxChecksPerThread + 1)
		result.Maxaccept = sample.Maxaccept + 1
		result.Maxpollevents = sample.Maxpollevents + 1
		result.Maxrewrite = sample.Maxrewrite + 1
		result.MemoryHotSize = Ptr(*sample.MemoryHotSize + 1)
		result.PatternCacheSize = Ptr(*sample.PatternCacheSize + 1)
		result.PeersMaxUpdatesAtOnce = sample.PeersMaxUpdatesAtOnce + 1
		result.Pipesize = sample.Pipesize + 1
		result.PoolHighFdRatio = sample.PoolHighFdRatio + 1
		result.PoolLowFdRatio = sample.PoolLowFdRatio + 1
		result.QuicFrontendConnTxBuffersLimit = Ptr(*sample.QuicFrontendConnTxBuffersLimit + 1)
		result.QuicFrontendMaxIdleTimeout = Ptr(*sample.QuicFrontendMaxIdleTimeout + 1)
		result.QuicFrontendMaxStreamsBidi = Ptr(*sample.QuicFrontendMaxStreamsBidi + 1)
		result.QuicMaxFrameLoss = Ptr(*sample.QuicMaxFrameLoss + 1)
		result.QuicRetryThreshold = Ptr(*sample.QuicRetryThreshold + 1)
		result.RcvbufBackend = Ptr(*sample.RcvbufBackend + 1)
		result.RcvbufClient = Ptr(*sample.RcvbufClient + 1)
		result.RcvbufFrontend = Ptr(*sample.RcvbufFrontend + 1)
		result.RcvbufServer = Ptr(*sample.RcvbufServer + 1)
		result.RecvEnough = sample.RecvEnough + 1
		result.RunqueueDepth = sample.RunqueueDepth + 1
		result.SndbufBackend = Ptr(*sample.SndbufBackend + 1)
		result.SndbufClient = Ptr(*sample.SndbufClient + 1)
		result.SndbufFrontend = Ptr(*sample.SndbufFrontend + 1)
		result.SndbufServer = Ptr(*sample.SndbufServer + 1)
		result.SslCachesize = Ptr(*sample.SslCachesize + 1)
		result.SslCaptureBufferSize = Ptr(*sample.SslCaptureBufferSize + 1)
		result.SslCtxCacheSize = sample.SslCtxCacheSize + 1
		result.SslDefaultDhParam = sample.SslDefaultDhParam + 1
		result.SslForcePrivateCache = !sample.SslForcePrivateCache
		result.SslLifetime = Ptr(*sample.SslLifetime + 1)
		result.SslMaxrecord = Ptr(*sample.SslMaxrecord + 1)
		result.SslOcspUpdateMaxDelay = Ptr(*sample.SslOcspUpdateMaxDelay + 1)
		result.SslOcspUpdateMinDelay = Ptr(*sample.SslOcspUpdateMinDelay + 1)
		result.StickCounters = Ptr(*sample.StickCounters + 1)
		result.VarsGlobalMaxSize = Ptr(*sample.VarsGlobalMaxSize + 1)
		result.VarsProcMaxSize = Ptr(*sample.VarsProcMaxSize + 1)
		result.VarsReqresMaxSize = Ptr(*sample.VarsReqresMaxSize + 1)
		result.VarsSessMaxSize = Ptr(*sample.VarsSessMaxSize + 1)
		result.VarsTxnMaxSize = Ptr(*sample.VarsTxnMaxSize + 1)
		result.ZlibMemlevel = sample.ZlibMemlevel + 1
		result.ZlibWindowsize = sample.ZlibWindowsize + 1
		samples = append(samples, struct {
			a, b GlobalTuneOptions
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		if len(result) != 75 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected GlobalTuneOptions to be different in 75 cases, but it is not (%d) %s %s", len(result), a, b)
		}
	}
}

func TestGlobalWurflOptionsEqual(t *testing.T) {
	samples := []struct {
		a, b GlobalWurflOptions
	}{}
	for i := 0; i < 2; i++ {
		var sample GlobalWurflOptions
		var result GlobalWurflOptions
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		byteJSON, err := json.Marshal(sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = json.Unmarshal(byteJSON, &result)
		if err != nil {
			t.Errorf(err.Error())
		}

		samples = append(samples, struct {
			a, b GlobalWurflOptions
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Equal(sample.b)
		if !result {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected GlobalWurflOptions to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestGlobalWurflOptionsEqualFalse(t *testing.T) {
	samples := []struct {
		a, b GlobalWurflOptions
	}{}
	for i := 0; i < 2; i++ {
		var sample GlobalWurflOptions
		var result GlobalWurflOptions
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		result.CacheSize = sample.CacheSize + 1
		samples = append(samples, struct {
			a, b GlobalWurflOptions
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Equal(sample.b)
		if result {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected GlobalWurflOptions to be different, but it is not %s %s", a, b)
		}
	}
}

func TestGlobalWurflOptionsDiff(t *testing.T) {
	samples := []struct {
		a, b GlobalWurflOptions
	}{}
	for i := 0; i < 2; i++ {
		var sample GlobalWurflOptions
		var result GlobalWurflOptions
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		byteJSON, err := json.Marshal(sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = json.Unmarshal(byteJSON, &result)
		if err != nil {
			t.Errorf(err.Error())
		}

		samples = append(samples, struct {
			a, b GlobalWurflOptions
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		if len(result) != 0 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected GlobalWurflOptions to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestGlobalWurflOptionsDiffFalse(t *testing.T) {
	samples := []struct {
		a, b GlobalWurflOptions
	}{}
	for i := 0; i < 2; i++ {
		var sample GlobalWurflOptions
		var result GlobalWurflOptions
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		result.CacheSize = sample.CacheSize + 1
		samples = append(samples, struct {
			a, b GlobalWurflOptions
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		if len(result) != 5 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected GlobalWurflOptions to be different in 5 cases, but it is not (%d) %s %s", len(result), a, b)
		}
	}
}
