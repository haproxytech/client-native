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

func TestStickTableEntryEqual(t *testing.T) {
	samples := []struct {
		a, b StickTableEntry
	}{}
	for i := 0; i < 2; i++ {
		var sample StickTableEntry
		var result StickTableEntry
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
			a, b StickTableEntry
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
			t.Errorf("Expected StickTableEntry to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestStickTableEntryEqualFalse(t *testing.T) {
	samples := []struct {
		a, b StickTableEntry
	}{}
	for i := 0; i < 2; i++ {
		var sample StickTableEntry
		var result StickTableEntry
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		result.BytesInCnt = Ptr(*sample.BytesInCnt + 1)
		result.BytesInRate = Ptr(*sample.BytesInRate + 1)
		result.BytesOutCnt = Ptr(*sample.BytesOutCnt + 1)
		result.BytesOutRate = Ptr(*sample.BytesOutRate + 1)
		result.ConnCnt = Ptr(*sample.ConnCnt + 1)
		result.ConnCur = Ptr(*sample.ConnCur + 1)
		result.ConnRate = Ptr(*sample.ConnRate + 1)
		result.Exp = Ptr(*sample.Exp + 1)
		result.GlitchCnt = Ptr(*sample.GlitchCnt + 1)
		result.GlitchRate = Ptr(*sample.GlitchRate + 1)
		result.Gpc0 = Ptr(*sample.Gpc0 + 1)
		result.Gpc0Rate = Ptr(*sample.Gpc0Rate + 1)
		result.Gpc1 = Ptr(*sample.Gpc1 + 1)
		result.Gpc1Rate = Ptr(*sample.Gpc1Rate + 1)
		result.Gpt0 = Ptr(*sample.Gpt0 + 1)
		result.HTTPErrCnt = Ptr(*sample.HTTPErrCnt + 1)
		result.HTTPErrRate = Ptr(*sample.HTTPErrRate + 1)
		result.HTTPFailCnt = Ptr(*sample.HTTPFailCnt + 1)
		result.HTTPFailRate = Ptr(*sample.HTTPFailRate + 1)
		result.HTTPReqCnt = Ptr(*sample.HTTPReqCnt + 1)
		result.HTTPReqRate = Ptr(*sample.HTTPReqRate + 1)
		result.ServerID = Ptr(*sample.ServerID + 1)
		result.SessCnt = Ptr(*sample.SessCnt + 1)
		result.SessRate = Ptr(*sample.SessRate + 1)
		result.Use = !sample.Use
		samples = append(samples, struct {
			a, b StickTableEntry
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
			t.Errorf("Expected StickTableEntry to be different, but it is not %s %s", a, b)
		}
	}
}

func TestStickTableEntryDiff(t *testing.T) {
	samples := []struct {
		a, b StickTableEntry
	}{}
	for i := 0; i < 2; i++ {
		var sample StickTableEntry
		var result StickTableEntry
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
			a, b StickTableEntry
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
			t.Errorf("Expected StickTableEntry to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestStickTableEntryDiffFalse(t *testing.T) {
	samples := []struct {
		a, b StickTableEntry
	}{}
	for i := 0; i < 2; i++ {
		var sample StickTableEntry
		var result StickTableEntry
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		result.BytesInCnt = Ptr(*sample.BytesInCnt + 1)
		result.BytesInRate = Ptr(*sample.BytesInRate + 1)
		result.BytesOutCnt = Ptr(*sample.BytesOutCnt + 1)
		result.BytesOutRate = Ptr(*sample.BytesOutRate + 1)
		result.ConnCnt = Ptr(*sample.ConnCnt + 1)
		result.ConnCur = Ptr(*sample.ConnCur + 1)
		result.ConnRate = Ptr(*sample.ConnRate + 1)
		result.Exp = Ptr(*sample.Exp + 1)
		result.GlitchCnt = Ptr(*sample.GlitchCnt + 1)
		result.GlitchRate = Ptr(*sample.GlitchRate + 1)
		result.Gpc0 = Ptr(*sample.Gpc0 + 1)
		result.Gpc0Rate = Ptr(*sample.Gpc0Rate + 1)
		result.Gpc1 = Ptr(*sample.Gpc1 + 1)
		result.Gpc1Rate = Ptr(*sample.Gpc1Rate + 1)
		result.Gpt0 = Ptr(*sample.Gpt0 + 1)
		result.HTTPErrCnt = Ptr(*sample.HTTPErrCnt + 1)
		result.HTTPErrRate = Ptr(*sample.HTTPErrRate + 1)
		result.HTTPFailCnt = Ptr(*sample.HTTPFailCnt + 1)
		result.HTTPFailRate = Ptr(*sample.HTTPFailRate + 1)
		result.HTTPReqCnt = Ptr(*sample.HTTPReqCnt + 1)
		result.HTTPReqRate = Ptr(*sample.HTTPReqRate + 1)
		result.ServerID = Ptr(*sample.ServerID + 1)
		result.SessCnt = Ptr(*sample.SessCnt + 1)
		result.SessRate = Ptr(*sample.SessRate + 1)
		result.Use = !sample.Use
		samples = append(samples, struct {
			a, b StickTableEntry
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		if len(result) != 30 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Error(err)
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Error(err)
			}
			t.Errorf("Expected StickTableEntry to be different in 30 cases, but it is not (%d) %s %s", len(result), a, b)
		}
	}
}

func TestStickTableEntryGpcEqual(t *testing.T) {
	samples := []struct {
		a, b StickTableEntryGpc
	}{}
	for i := 0; i < 2; i++ {
		var sample StickTableEntryGpc
		var result StickTableEntryGpc
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
			a, b StickTableEntryGpc
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
			t.Errorf("Expected StickTableEntryGpc to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestStickTableEntryGpcEqualFalse(t *testing.T) {
	samples := []struct {
		a, b StickTableEntryGpc
	}{}
	for i := 0; i < 2; i++ {
		var sample StickTableEntryGpc
		var result StickTableEntryGpc
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		result.Idx = sample.Idx + 1
		result.Value = Ptr(*sample.Value + 1)
		samples = append(samples, struct {
			a, b StickTableEntryGpc
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
			t.Errorf("Expected StickTableEntryGpc to be different, but it is not %s %s", a, b)
		}
	}
}

func TestStickTableEntryGpcDiff(t *testing.T) {
	samples := []struct {
		a, b StickTableEntryGpc
	}{}
	for i := 0; i < 2; i++ {
		var sample StickTableEntryGpc
		var result StickTableEntryGpc
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
			a, b StickTableEntryGpc
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
			t.Errorf("Expected StickTableEntryGpc to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestStickTableEntryGpcDiffFalse(t *testing.T) {
	samples := []struct {
		a, b StickTableEntryGpc
	}{}
	for i := 0; i < 2; i++ {
		var sample StickTableEntryGpc
		var result StickTableEntryGpc
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		result.Idx = sample.Idx + 1
		result.Value = Ptr(*sample.Value + 1)
		samples = append(samples, struct {
			a, b StickTableEntryGpc
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
			t.Errorf("Expected StickTableEntryGpc to be different in 2 cases, but it is not (%d) %s %s", len(result), a, b)
		}
	}
}

func TestStickTableEntryGpcRateEqual(t *testing.T) {
	samples := []struct {
		a, b StickTableEntryGpcRate
	}{}
	for i := 0; i < 2; i++ {
		var sample StickTableEntryGpcRate
		var result StickTableEntryGpcRate
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
			a, b StickTableEntryGpcRate
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
			t.Errorf("Expected StickTableEntryGpcRate to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestStickTableEntryGpcRateEqualFalse(t *testing.T) {
	samples := []struct {
		a, b StickTableEntryGpcRate
	}{}
	for i := 0; i < 2; i++ {
		var sample StickTableEntryGpcRate
		var result StickTableEntryGpcRate
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		result.Idx = sample.Idx + 1
		result.Value = Ptr(*sample.Value + 1)
		samples = append(samples, struct {
			a, b StickTableEntryGpcRate
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
			t.Errorf("Expected StickTableEntryGpcRate to be different, but it is not %s %s", a, b)
		}
	}
}

func TestStickTableEntryGpcRateDiff(t *testing.T) {
	samples := []struct {
		a, b StickTableEntryGpcRate
	}{}
	for i := 0; i < 2; i++ {
		var sample StickTableEntryGpcRate
		var result StickTableEntryGpcRate
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
			a, b StickTableEntryGpcRate
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
			t.Errorf("Expected StickTableEntryGpcRate to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestStickTableEntryGpcRateDiffFalse(t *testing.T) {
	samples := []struct {
		a, b StickTableEntryGpcRate
	}{}
	for i := 0; i < 2; i++ {
		var sample StickTableEntryGpcRate
		var result StickTableEntryGpcRate
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		result.Idx = sample.Idx + 1
		result.Value = Ptr(*sample.Value + 1)
		samples = append(samples, struct {
			a, b StickTableEntryGpcRate
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
			t.Errorf("Expected StickTableEntryGpcRate to be different in 2 cases, but it is not (%d) %s %s", len(result), a, b)
		}
	}
}

func TestStickTableEntryGptEqual(t *testing.T) {
	samples := []struct {
		a, b StickTableEntryGpt
	}{}
	for i := 0; i < 2; i++ {
		var sample StickTableEntryGpt
		var result StickTableEntryGpt
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
			a, b StickTableEntryGpt
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
			t.Errorf("Expected StickTableEntryGpt to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestStickTableEntryGptEqualFalse(t *testing.T) {
	samples := []struct {
		a, b StickTableEntryGpt
	}{}
	for i := 0; i < 2; i++ {
		var sample StickTableEntryGpt
		var result StickTableEntryGpt
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		result.Idx = sample.Idx + 1
		result.Value = Ptr(*sample.Value + 1)
		samples = append(samples, struct {
			a, b StickTableEntryGpt
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
			t.Errorf("Expected StickTableEntryGpt to be different, but it is not %s %s", a, b)
		}
	}
}

func TestStickTableEntryGptDiff(t *testing.T) {
	samples := []struct {
		a, b StickTableEntryGpt
	}{}
	for i := 0; i < 2; i++ {
		var sample StickTableEntryGpt
		var result StickTableEntryGpt
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
			a, b StickTableEntryGpt
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
			t.Errorf("Expected StickTableEntryGpt to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestStickTableEntryGptDiffFalse(t *testing.T) {
	samples := []struct {
		a, b StickTableEntryGpt
	}{}
	for i := 0; i < 2; i++ {
		var sample StickTableEntryGpt
		var result StickTableEntryGpt
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		result.Idx = sample.Idx + 1
		result.Value = Ptr(*sample.Value + 1)
		samples = append(samples, struct {
			a, b StickTableEntryGpt
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
			t.Errorf("Expected StickTableEntryGpt to be different in 2 cases, but it is not (%d) %s %s", len(result), a, b)
		}
	}
}
