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

func TestGlobalBaseEqual(t *testing.T) {
	samples := []struct {
		a, b GlobalBase
	}{}
	for i := 0; i < 2; i++ {
		var sample GlobalBase
		var result GlobalBase
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
			a, b GlobalBase
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
			t.Errorf("Expected GlobalBase to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestGlobalBaseEqualFalse(t *testing.T) {
	samples := []struct {
		a, b GlobalBase
	}{}
	for i := 0; i < 2; i++ {
		var sample GlobalBase
		var result GlobalBase
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		result.CloseSpreadTime = Ptr(*sample.CloseSpreadTime + 1)
		result.Daemon = !sample.Daemon
		result.ExposeDeprecatedDirectives = !sample.ExposeDeprecatedDirectives
		result.ExposeExperimentalDirectives = !sample.ExposeExperimentalDirectives
		result.ExternalCheck = !sample.ExternalCheck
		result.ForceCfgParserPause = Ptr(*sample.ForceCfgParserPause + 1)
		result.Gid = sample.Gid + 1
		result.Grace = Ptr(*sample.Grace + 1)
		result.H1AcceptPayloadWithAnyMethod = !sample.H1AcceptPayloadWithAnyMethod
		result.H1DoNotCloseOnInsecureTransferEncoding = !sample.H1DoNotCloseOnInsecureTransferEncoding
		result.H2WorkaroundBogusWebsocketClients = !sample.H2WorkaroundBogusWebsocketClients
		result.HardStopAfter = Ptr(*sample.HardStopAfter + 1)
		result.InsecureForkWanted = !sample.InsecureForkWanted
		result.InsecureSetuidWanted = !sample.InsecureSetuidWanted
		result.LimitedQuic = !sample.LimitedQuic
		result.MasterWorker = !sample.MasterWorker
		result.MworkerMaxReloads = Ptr(*sample.MworkerMaxReloads + 1)
		result.Nbthread = sample.Nbthread + 1
		result.NoQuic = !sample.NoQuic
		result.Pp2NeverSendLocal = !sample.Pp2NeverSendLocal
		result.PreallocFd = !sample.PreallocFd
		result.SetDumpable = !sample.SetDumpable
		result.ShmStatsFileMaxObjects = Ptr(*sample.ShmStatsFileMaxObjects + 1)
		result.StatsMaxconn = Ptr(*sample.StatsMaxconn + 1)
		result.StatsTimeout = Ptr(*sample.StatsTimeout + 1)
		result.StrictLimits = !sample.StrictLimits
		result.ThreadGroups = sample.ThreadGroups + 1
		result.UID = sample.UID + 1
		result.Ulimitn = sample.Ulimitn + 1
		result.WarnBlockedTrafficAfter = Ptr(*sample.WarnBlockedTrafficAfter + 1)
		samples = append(samples, struct {
			a, b GlobalBase
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
			t.Errorf("Expected GlobalBase to be different, but it is not %s %s", a, b)
		}
	}
}

func TestGlobalBaseDiff(t *testing.T) {
	samples := []struct {
		a, b GlobalBase
	}{}
	for i := 0; i < 2; i++ {
		var sample GlobalBase
		var result GlobalBase
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
			a, b GlobalBase
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
			t.Errorf("Expected GlobalBase to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestGlobalBaseDiffFalse(t *testing.T) {
	samples := []struct {
		a, b GlobalBase
	}{}
	for i := 0; i < 2; i++ {
		var sample GlobalBase
		var result GlobalBase
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		result.CloseSpreadTime = Ptr(*sample.CloseSpreadTime + 1)
		result.Daemon = !sample.Daemon
		result.ExposeDeprecatedDirectives = !sample.ExposeDeprecatedDirectives
		result.ExposeExperimentalDirectives = !sample.ExposeExperimentalDirectives
		result.ExternalCheck = !sample.ExternalCheck
		result.ForceCfgParserPause = Ptr(*sample.ForceCfgParserPause + 1)
		result.Gid = sample.Gid + 1
		result.Grace = Ptr(*sample.Grace + 1)
		result.H1AcceptPayloadWithAnyMethod = !sample.H1AcceptPayloadWithAnyMethod
		result.H1DoNotCloseOnInsecureTransferEncoding = !sample.H1DoNotCloseOnInsecureTransferEncoding
		result.H2WorkaroundBogusWebsocketClients = !sample.H2WorkaroundBogusWebsocketClients
		result.HardStopAfter = Ptr(*sample.HardStopAfter + 1)
		result.InsecureForkWanted = !sample.InsecureForkWanted
		result.InsecureSetuidWanted = !sample.InsecureSetuidWanted
		result.LimitedQuic = !sample.LimitedQuic
		result.MasterWorker = !sample.MasterWorker
		result.MworkerMaxReloads = Ptr(*sample.MworkerMaxReloads + 1)
		result.Nbthread = sample.Nbthread + 1
		result.NoQuic = !sample.NoQuic
		result.Pp2NeverSendLocal = !sample.Pp2NeverSendLocal
		result.PreallocFd = !sample.PreallocFd
		result.SetDumpable = !sample.SetDumpable
		result.ShmStatsFileMaxObjects = Ptr(*sample.ShmStatsFileMaxObjects + 1)
		result.StatsMaxconn = Ptr(*sample.StatsMaxconn + 1)
		result.StatsTimeout = Ptr(*sample.StatsTimeout + 1)
		result.StrictLimits = !sample.StrictLimits
		result.ThreadGroups = sample.ThreadGroups + 1
		result.UID = sample.UID + 1
		result.Ulimitn = sample.Ulimitn + 1
		result.WarnBlockedTrafficAfter = Ptr(*sample.WarnBlockedTrafficAfter + 1)
		samples = append(samples, struct {
			a, b GlobalBase
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		listDiffFields := GetListOfDiffFields(result)
		if len(listDiffFields) != 74 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Error(err)
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Error(err)
			}
			t.Errorf("Expected GlobalBase to be different in 74 cases, but it is not (%d) %s %s", len(result), a, b)
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
			a, b CPUMap
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
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
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
				t.Error(err)
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Error(err)
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
			a, b CPUMap
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
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		samples = append(samples, struct {
			a, b CPUMap
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		listDiffFields := GetListOfDiffFields(result)
		if len(listDiffFields) != 2 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Error(err)
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Error(err)
			}
			t.Errorf("Expected CPUMap to be different in 2 cases, but it is not (%d) %s %s", len(result), a, b)
		}
	}
}

func TestCPUSetEqual(t *testing.T) {
	samples := []struct {
		a, b CPUSet
	}{}
	for i := 0; i < 2; i++ {
		var sample CPUSet
		var result CPUSet
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
			a, b CPUSet
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
			t.Errorf("Expected CPUSet to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestCPUSetEqualFalse(t *testing.T) {
	samples := []struct {
		a, b CPUSet
	}{}
	for i := 0; i < 2; i++ {
		var sample CPUSet
		var result CPUSet
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		samples = append(samples, struct {
			a, b CPUSet
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
			t.Errorf("Expected CPUSet to be different, but it is not %s %s", a, b)
		}
	}
}

func TestCPUSetDiff(t *testing.T) {
	samples := []struct {
		a, b CPUSet
	}{}
	for i := 0; i < 2; i++ {
		var sample CPUSet
		var result CPUSet
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
			a, b CPUSet
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
			t.Errorf("Expected CPUSet to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestCPUSetDiffFalse(t *testing.T) {
	samples := []struct {
		a, b CPUSet
	}{}
	for i := 0; i < 2; i++ {
		var sample CPUSet
		var result CPUSet
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		samples = append(samples, struct {
			a, b CPUSet
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		listDiffFields := GetListOfDiffFields(result)
		if len(listDiffFields) != 2 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Error(err)
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Error(err)
			}
			t.Errorf("Expected CPUSet to be different in 2 cases, but it is not (%d) %s %s", len(result), a, b)
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
			a, b GlobalDefaultPath
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
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
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
				t.Error(err)
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Error(err)
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
			a, b GlobalDefaultPath
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
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		samples = append(samples, struct {
			a, b GlobalDefaultPath
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		listDiffFields := GetListOfDiffFields(result)
		if len(listDiffFields) != 2 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Error(err)
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Error(err)
			}
			t.Errorf("Expected GlobalDefaultPath to be different in 2 cases, but it is not (%d) %s %s", len(result), a, b)
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
			a, b H1CaseAdjust
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
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
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
				t.Error(err)
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Error(err)
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
			a, b H1CaseAdjust
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
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		samples = append(samples, struct {
			a, b H1CaseAdjust
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		listDiffFields := GetListOfDiffFields(result)
		if len(listDiffFields) != 2 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Error(err)
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Error(err)
			}
			t.Errorf("Expected H1CaseAdjust to be different in 2 cases, but it is not (%d) %s %s", len(result), a, b)
		}
	}
}

func TestGlobalHardenEqual(t *testing.T) {
	samples := []struct {
		a, b GlobalHarden
	}{}
	for i := 0; i < 2; i++ {
		var sample GlobalHarden
		var result GlobalHarden
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
			a, b GlobalHarden
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
			t.Errorf("Expected GlobalHarden to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestGlobalHardenEqualFalse(t *testing.T) {
	samples := []struct {
		a, b GlobalHarden
	}{}
	for i := 0; i < 2; i++ {
		var sample GlobalHarden
		var result GlobalHarden
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		samples = append(samples, struct {
			a, b GlobalHarden
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
			t.Errorf("Expected GlobalHarden to be different, but it is not %s %s", a, b)
		}
	}
}

func TestGlobalHardenDiff(t *testing.T) {
	samples := []struct {
		a, b GlobalHarden
	}{}
	for i := 0; i < 2; i++ {
		var sample GlobalHarden
		var result GlobalHarden
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
			a, b GlobalHarden
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
			t.Errorf("Expected GlobalHarden to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestGlobalHardenDiffFalse(t *testing.T) {
	samples := []struct {
		a, b GlobalHarden
	}{}
	for i := 0; i < 2; i++ {
		var sample GlobalHarden
		var result GlobalHarden
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		samples = append(samples, struct {
			a, b GlobalHarden
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		listDiffFields := GetListOfDiffFields(result)
		if len(listDiffFields) != 1 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Error(err)
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Error(err)
			}
			t.Errorf("Expected GlobalHarden to be different in 1 cases, but it is not (%d) %s %s", len(result), a, b)
		}
	}
}

func TestGlobalHardenRejectPrivilegedPortsEqual(t *testing.T) {
	samples := []struct {
		a, b GlobalHardenRejectPrivilegedPorts
	}{}
	for i := 0; i < 2; i++ {
		var sample GlobalHardenRejectPrivilegedPorts
		var result GlobalHardenRejectPrivilegedPorts
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
			a, b GlobalHardenRejectPrivilegedPorts
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
			t.Errorf("Expected GlobalHardenRejectPrivilegedPorts to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestGlobalHardenRejectPrivilegedPortsEqualFalse(t *testing.T) {
	samples := []struct {
		a, b GlobalHardenRejectPrivilegedPorts
	}{}
	for i := 0; i < 2; i++ {
		var sample GlobalHardenRejectPrivilegedPorts
		var result GlobalHardenRejectPrivilegedPorts
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		samples = append(samples, struct {
			a, b GlobalHardenRejectPrivilegedPorts
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
			t.Errorf("Expected GlobalHardenRejectPrivilegedPorts to be different, but it is not %s %s", a, b)
		}
	}
}

func TestGlobalHardenRejectPrivilegedPortsDiff(t *testing.T) {
	samples := []struct {
		a, b GlobalHardenRejectPrivilegedPorts
	}{}
	for i := 0; i < 2; i++ {
		var sample GlobalHardenRejectPrivilegedPorts
		var result GlobalHardenRejectPrivilegedPorts
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
			a, b GlobalHardenRejectPrivilegedPorts
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
			t.Errorf("Expected GlobalHardenRejectPrivilegedPorts to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestGlobalHardenRejectPrivilegedPortsDiffFalse(t *testing.T) {
	samples := []struct {
		a, b GlobalHardenRejectPrivilegedPorts
	}{}
	for i := 0; i < 2; i++ {
		var sample GlobalHardenRejectPrivilegedPorts
		var result GlobalHardenRejectPrivilegedPorts
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		samples = append(samples, struct {
			a, b GlobalHardenRejectPrivilegedPorts
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		listDiffFields := GetListOfDiffFields(result)
		if len(listDiffFields) != 2 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Error(err)
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Error(err)
			}
			t.Errorf("Expected GlobalHardenRejectPrivilegedPorts to be different in 2 cases, but it is not (%d) %s %s", len(result), a, b)
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
			a, b GlobalLogSendHostname
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
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
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
				t.Error(err)
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Error(err)
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
			a, b GlobalLogSendHostname
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
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		samples = append(samples, struct {
			a, b GlobalLogSendHostname
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		listDiffFields := GetListOfDiffFields(result)
		if len(listDiffFields) != 2 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Error(err)
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Error(err)
			}
			t.Errorf("Expected GlobalLogSendHostname to be different in 2 cases, but it is not (%d) %s %s", len(result), a, b)
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
			a, b RuntimeAPI
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
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
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
				t.Error(err)
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Error(err)
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
			a, b RuntimeAPI
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
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		samples = append(samples, struct {
			a, b RuntimeAPI
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		listDiffFields := GetListOfDiffFields(result)
		if len(listDiffFields) != 3 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Error(err)
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Error(err)
			}
			t.Errorf("Expected RuntimeAPI to be different in 3 cases, but it is not (%d) %s %s", len(result), a, b)
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
			a, b SetVarFmt
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
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
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
				t.Error(err)
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Error(err)
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
			a, b SetVarFmt
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
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		samples = append(samples, struct {
			a, b SetVarFmt
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		listDiffFields := GetListOfDiffFields(result)
		if len(listDiffFields) != 2 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Error(err)
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Error(err)
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
			a, b SetVar
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
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
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
				t.Error(err)
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Error(err)
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
			a, b SetVar
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
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		samples = append(samples, struct {
			a, b SetVar
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		listDiffFields := GetListOfDiffFields(result)
		if len(listDiffFields) != 2 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Error(err)
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Error(err)
			}
			t.Errorf("Expected SetVar to be different in 2 cases, but it is not (%d) %s %s", len(result), a, b)
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
			a, b ThreadGroup
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
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
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
				t.Error(err)
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Error(err)
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
			a, b ThreadGroup
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
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		samples = append(samples, struct {
			a, b ThreadGroup
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		listDiffFields := GetListOfDiffFields(result)
		if len(listDiffFields) != 2 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Error(err)
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Error(err)
			}
			t.Errorf("Expected ThreadGroup to be different in 2 cases, but it is not (%d) %s %s", len(result), a, b)
		}
	}
}
