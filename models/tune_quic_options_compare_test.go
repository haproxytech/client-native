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

func TestTuneQuicOptionsEqual(t *testing.T) {
	samples := []struct {
		a, b TuneQuicOptions
	}{}
	for i := 0; i < 2; i++ {
		var sample TuneQuicOptions
		var result TuneQuicOptions
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
			a, b TuneQuicOptions
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
			t.Errorf("Expected TuneQuicOptions to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestTuneQuicOptionsEqualFalse(t *testing.T) {
	samples := []struct {
		a, b TuneQuicOptions
	}{}
	for i := 0; i < 2; i++ {
		var sample TuneQuicOptions
		var result TuneQuicOptions
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		result.FrontendConnTxBuffersLimit = Ptr(*sample.FrontendConnTxBuffersLimit + 1)
		result.FrontendMaxIdleTimeout = Ptr(*sample.FrontendMaxIdleTimeout + 1)
		result.FrontendMaxStreamsBidi = Ptr(*sample.FrontendMaxStreamsBidi + 1)
		result.MaxFrameLoss = Ptr(*sample.MaxFrameLoss + 1)
		result.ReorderRatio = Ptr(*sample.ReorderRatio + 1)
		result.RetryThreshold = Ptr(*sample.RetryThreshold + 1)
		samples = append(samples, struct {
			a, b TuneQuicOptions
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
			t.Errorf("Expected TuneQuicOptions to be different, but it is not %s %s", a, b)
		}
	}
}

func TestTuneQuicOptionsDiff(t *testing.T) {
	samples := []struct {
		a, b TuneQuicOptions
	}{}
	for i := 0; i < 2; i++ {
		var sample TuneQuicOptions
		var result TuneQuicOptions
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
			a, b TuneQuicOptions
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
			t.Errorf("Expected TuneQuicOptions to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestTuneQuicOptionsDiffFalse(t *testing.T) {
	samples := []struct {
		a, b TuneQuicOptions
	}{}
	for i := 0; i < 2; i++ {
		var sample TuneQuicOptions
		var result TuneQuicOptions
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		result.FrontendConnTxBuffersLimit = Ptr(*sample.FrontendConnTxBuffersLimit + 1)
		result.FrontendMaxIdleTimeout = Ptr(*sample.FrontendMaxIdleTimeout + 1)
		result.FrontendMaxStreamsBidi = Ptr(*sample.FrontendMaxStreamsBidi + 1)
		result.MaxFrameLoss = Ptr(*sample.MaxFrameLoss + 1)
		result.ReorderRatio = Ptr(*sample.ReorderRatio + 1)
		result.RetryThreshold = Ptr(*sample.RetryThreshold + 1)
		samples = append(samples, struct {
			a, b TuneQuicOptions
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		if len(result) != 8 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected TuneQuicOptions to be different in 8 cases, but it is not (%d) %s %s", len(result), a, b)
		}
	}
}
