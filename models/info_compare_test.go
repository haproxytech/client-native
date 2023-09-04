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
	"math/rand"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/go-openapi/strfmt"

	jsoniter "github.com/json-iterator/go"
)

func TestInfoEqual(t *testing.T) {
	samples := []struct {
		a, b Info
	}{}
	for i := 0; i < 2; i++ {
		var sample Info
		var result Info
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
			a, b Info
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
			t.Errorf("Expected Info to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestInfoEqualFalse(t *testing.T) {
	samples := []struct {
		a, b Info
	}{}
	for i := 0; i < 2; i++ {
		var sample Info
		var result Info
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		samples = append(samples, struct {
			a, b Info
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
			t.Errorf("Expected Info to be different, but it is not %s %s", a, b)
		}
	}
}

func TestInfoDiff(t *testing.T) {
	samples := []struct {
		a, b Info
	}{}
	for i := 0; i < 2; i++ {
		var sample Info
		var result Info
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
			a, b Info
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
			t.Errorf("Expected Info to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestInfoDiffFalse(t *testing.T) {
	samples := []struct {
		a, b Info
	}{}
	for i := 0; i < 2; i++ {
		var sample Info
		var result Info
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		samples = append(samples, struct {
			a, b Info
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
			t.Errorf("Expected Info to be different in 2 cases, but it is not (%d) %s %s", len(result), a, b)
		}
	}
}

func TestInfoAPIEqual(t *testing.T) {
	samples := []struct {
		a, b InfoAPI
	}{}
	for i := 0; i < 2; i++ {
		var sample InfoAPI
		var result InfoAPI
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
			a, b InfoAPI
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
			t.Errorf("Expected InfoAPI to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestInfoAPIEqualFalse(t *testing.T) {
	samples := []struct {
		a, b InfoAPI
	}{}
	for i := 0; i < 2; i++ {
		var sample InfoAPI
		var result InfoAPI
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		result.BuildDate = strfmt.DateTime(time.Now().AddDate(rand.Intn(10), rand.Intn(12), rand.Intn(28)))
		samples = append(samples, struct {
			a, b InfoAPI
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
			t.Errorf("Expected InfoAPI to be different, but it is not %s %s", a, b)
		}
	}
}

func TestInfoAPIDiff(t *testing.T) {
	samples := []struct {
		a, b InfoAPI
	}{}
	for i := 0; i < 2; i++ {
		var sample InfoAPI
		var result InfoAPI
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
			a, b InfoAPI
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
			t.Errorf("Expected InfoAPI to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestInfoAPIDiffFalse(t *testing.T) {
	samples := []struct {
		a, b InfoAPI
	}{}
	for i := 0; i < 2; i++ {
		var sample InfoAPI
		var result InfoAPI
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		result.BuildDate = strfmt.DateTime(time.Now().AddDate(rand.Intn(10), rand.Intn(12), rand.Intn(28)))
		samples = append(samples, struct {
			a, b InfoAPI
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
			t.Errorf("Expected InfoAPI to be different in 2 cases, but it is not (%d) %s %s", len(result), a, b)
		}
	}
}

func TestInfoSystemEqual(t *testing.T) {
	samples := []struct {
		a, b InfoSystem
	}{}
	for i := 0; i < 2; i++ {
		var sample InfoSystem
		var result InfoSystem
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
			a, b InfoSystem
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
			t.Errorf("Expected InfoSystem to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestInfoSystemEqualFalse(t *testing.T) {
	samples := []struct {
		a, b InfoSystem
	}{}
	for i := 0; i < 2; i++ {
		var sample InfoSystem
		var result InfoSystem
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		result.Time = sample.Time + 1
		result.Uptime = Ptr(*sample.Uptime + 1)
		samples = append(samples, struct {
			a, b InfoSystem
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
			t.Errorf("Expected InfoSystem to be different, but it is not %s %s", a, b)
		}
	}
}

func TestInfoSystemDiff(t *testing.T) {
	samples := []struct {
		a, b InfoSystem
	}{}
	for i := 0; i < 2; i++ {
		var sample InfoSystem
		var result InfoSystem
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
			a, b InfoSystem
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
			t.Errorf("Expected InfoSystem to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestInfoSystemDiffFalse(t *testing.T) {
	samples := []struct {
		a, b InfoSystem
	}{}
	for i := 0; i < 2; i++ {
		var sample InfoSystem
		var result InfoSystem
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		result.Time = sample.Time + 1
		result.Uptime = Ptr(*sample.Uptime + 1)
		samples = append(samples, struct {
			a, b InfoSystem
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		if len(result) != 6 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected InfoSystem to be different in 6 cases, but it is not (%d) %s %s", len(result), a, b)
		}
	}
}

func TestInfoSystemCPUInfoEqual(t *testing.T) {
	samples := []struct {
		a, b InfoSystemCPUInfo
	}{}
	for i := 0; i < 2; i++ {
		var sample InfoSystemCPUInfo
		var result InfoSystemCPUInfo
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
			a, b InfoSystemCPUInfo
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
			t.Errorf("Expected InfoSystemCPUInfo to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestInfoSystemCPUInfoEqualFalse(t *testing.T) {
	samples := []struct {
		a, b InfoSystemCPUInfo
	}{}
	for i := 0; i < 2; i++ {
		var sample InfoSystemCPUInfo
		var result InfoSystemCPUInfo
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		result.NumCpus = sample.NumCpus + 1
		samples = append(samples, struct {
			a, b InfoSystemCPUInfo
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
			t.Errorf("Expected InfoSystemCPUInfo to be different, but it is not %s %s", a, b)
		}
	}
}

func TestInfoSystemCPUInfoDiff(t *testing.T) {
	samples := []struct {
		a, b InfoSystemCPUInfo
	}{}
	for i := 0; i < 2; i++ {
		var sample InfoSystemCPUInfo
		var result InfoSystemCPUInfo
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
			a, b InfoSystemCPUInfo
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
			t.Errorf("Expected InfoSystemCPUInfo to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestInfoSystemCPUInfoDiffFalse(t *testing.T) {
	samples := []struct {
		a, b InfoSystemCPUInfo
	}{}
	for i := 0; i < 2; i++ {
		var sample InfoSystemCPUInfo
		var result InfoSystemCPUInfo
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		result.NumCpus = sample.NumCpus + 1
		samples = append(samples, struct {
			a, b InfoSystemCPUInfo
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
			t.Errorf("Expected InfoSystemCPUInfo to be different in 2 cases, but it is not (%d) %s %s", len(result), a, b)
		}
	}
}

func TestInfoSystemMemInfoEqual(t *testing.T) {
	samples := []struct {
		a, b InfoSystemMemInfo
	}{}
	for i := 0; i < 2; i++ {
		var sample InfoSystemMemInfo
		var result InfoSystemMemInfo
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
			a, b InfoSystemMemInfo
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
			t.Errorf("Expected InfoSystemMemInfo to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestInfoSystemMemInfoEqualFalse(t *testing.T) {
	samples := []struct {
		a, b InfoSystemMemInfo
	}{}
	for i := 0; i < 2; i++ {
		var sample InfoSystemMemInfo
		var result InfoSystemMemInfo
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		result.DataplaneapiMemory = sample.DataplaneapiMemory + 1
		result.FreeMemory = sample.FreeMemory + 1
		result.TotalMemory = sample.TotalMemory + 1
		samples = append(samples, struct {
			a, b InfoSystemMemInfo
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
			t.Errorf("Expected InfoSystemMemInfo to be different, but it is not %s %s", a, b)
		}
	}
}

func TestInfoSystemMemInfoDiff(t *testing.T) {
	samples := []struct {
		a, b InfoSystemMemInfo
	}{}
	for i := 0; i < 2; i++ {
		var sample InfoSystemMemInfo
		var result InfoSystemMemInfo
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
			a, b InfoSystemMemInfo
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
			t.Errorf("Expected InfoSystemMemInfo to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestInfoSystemMemInfoDiffFalse(t *testing.T) {
	samples := []struct {
		a, b InfoSystemMemInfo
	}{}
	for i := 0; i < 2; i++ {
		var sample InfoSystemMemInfo
		var result InfoSystemMemInfo
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		result.DataplaneapiMemory = sample.DataplaneapiMemory + 1
		result.FreeMemory = sample.FreeMemory + 1
		result.TotalMemory = sample.TotalMemory + 1
		samples = append(samples, struct {
			a, b InfoSystemMemInfo
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		if len(result) != 3 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected InfoSystemMemInfo to be different in 3 cases, but it is not (%d) %s %s", len(result), a, b)
		}
	}
}
