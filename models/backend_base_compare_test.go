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
	"github.com/go-faker/faker/v4/pkg/options"

	jsoniter "github.com/json-iterator/go"
)

func TestBackendBaseEqual(t *testing.T) {
	samples := []struct {
		a, b BackendBase
	}{}
	for i := 0; i < 2; i++ {
		var sample BackendBase
		var result BackendBase
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
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
			a, b BackendBase
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
			t.Errorf("Expected BackendBase to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestBackendBaseEqualFalse(t *testing.T) {
	samples := []struct {
		a, b BackendBase
	}{}
	for i := 0; i < 2; i++ {
		var sample BackendBase
		var result BackendBase
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Errorf(err.Error())
		}
		result.CheckTimeout = Ptr(*sample.CheckTimeout + 1)
		result.ConnectTimeout = Ptr(*sample.ConnectTimeout + 1)
		result.Disabled = !sample.Disabled
		result.Enabled = !sample.Enabled
		result.Fullconn = Ptr(*sample.Fullconn + 1)
		result.HashBalanceFactor = Ptr(*sample.HashBalanceFactor + 1)
		result.HTTPKeepAliveTimeout = Ptr(*sample.HTTPKeepAliveTimeout + 1)
		result.HTTPRequestTimeout = Ptr(*sample.HTTPRequestTimeout + 1)
		result.ID = Ptr(*sample.ID + 1)
		result.MaxKeepAliveQueue = Ptr(*sample.MaxKeepAliveQueue + 1)
		result.QueueTimeout = Ptr(*sample.QueueTimeout + 1)
		result.Retries = Ptr(*sample.Retries + 1)
		result.ServerFinTimeout = Ptr(*sample.ServerFinTimeout + 1)
		result.ServerTimeout = Ptr(*sample.ServerTimeout + 1)
		result.SrvtcpkaCnt = Ptr(*sample.SrvtcpkaCnt + 1)
		result.SrvtcpkaIdle = Ptr(*sample.SrvtcpkaIdle + 1)
		result.SrvtcpkaIntvl = Ptr(*sample.SrvtcpkaIntvl + 1)
		result.TarpitTimeout = Ptr(*sample.TarpitTimeout + 1)
		result.TunnelTimeout = Ptr(*sample.TunnelTimeout + 1)
		samples = append(samples, struct {
			a, b BackendBase
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
			t.Errorf("Expected BackendBase to be different, but it is not %s %s", a, b)
		}
	}
}

func TestBackendBaseDiff(t *testing.T) {
	samples := []struct {
		a, b BackendBase
	}{}
	for i := 0; i < 2; i++ {
		var sample BackendBase
		var result BackendBase
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
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
			a, b BackendBase
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
			t.Errorf("Expected BackendBase to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestBackendBaseDiffFalse(t *testing.T) {
	samples := []struct {
		a, b BackendBase
	}{}
	for i := 0; i < 2; i++ {
		var sample BackendBase
		var result BackendBase
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Errorf(err.Error())
		}
		result.CheckTimeout = Ptr(*sample.CheckTimeout + 1)
		result.ConnectTimeout = Ptr(*sample.ConnectTimeout + 1)
		result.Disabled = !sample.Disabled
		result.Enabled = !sample.Enabled
		result.Fullconn = Ptr(*sample.Fullconn + 1)
		result.HashBalanceFactor = Ptr(*sample.HashBalanceFactor + 1)
		result.HTTPKeepAliveTimeout = Ptr(*sample.HTTPKeepAliveTimeout + 1)
		result.HTTPRequestTimeout = Ptr(*sample.HTTPRequestTimeout + 1)
		result.ID = Ptr(*sample.ID + 1)
		result.MaxKeepAliveQueue = Ptr(*sample.MaxKeepAliveQueue + 1)
		result.QueueTimeout = Ptr(*sample.QueueTimeout + 1)
		result.Retries = Ptr(*sample.Retries + 1)
		result.ServerFinTimeout = Ptr(*sample.ServerFinTimeout + 1)
		result.ServerTimeout = Ptr(*sample.ServerTimeout + 1)
		result.SrvtcpkaCnt = Ptr(*sample.SrvtcpkaCnt + 1)
		result.SrvtcpkaIdle = Ptr(*sample.SrvtcpkaIdle + 1)
		result.SrvtcpkaIntvl = Ptr(*sample.SrvtcpkaIntvl + 1)
		result.TarpitTimeout = Ptr(*sample.TarpitTimeout + 1)
		result.TunnelTimeout = Ptr(*sample.TunnelTimeout + 1)
		samples = append(samples, struct {
			a, b BackendBase
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		if len(result) != 86 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected BackendBase to be different in 86 cases, but it is not (%d) %s %s", len(result), a, b)
		}
	}
}

func TestBackendForcePersistEqual(t *testing.T) {
	samples := []struct {
		a, b BackendForcePersist
	}{}
	for i := 0; i < 2; i++ {
		var sample BackendForcePersist
		var result BackendForcePersist
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
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
			a, b BackendForcePersist
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
			t.Errorf("Expected BackendForcePersist to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestBackendForcePersistEqualFalse(t *testing.T) {
	samples := []struct {
		a, b BackendForcePersist
	}{}
	for i := 0; i < 2; i++ {
		var sample BackendForcePersist
		var result BackendForcePersist
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Errorf(err.Error())
		}
		samples = append(samples, struct {
			a, b BackendForcePersist
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
			t.Errorf("Expected BackendForcePersist to be different, but it is not %s %s", a, b)
		}
	}
}

func TestBackendForcePersistDiff(t *testing.T) {
	samples := []struct {
		a, b BackendForcePersist
	}{}
	for i := 0; i < 2; i++ {
		var sample BackendForcePersist
		var result BackendForcePersist
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
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
			a, b BackendForcePersist
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
			t.Errorf("Expected BackendForcePersist to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestBackendForcePersistDiffFalse(t *testing.T) {
	samples := []struct {
		a, b BackendForcePersist
	}{}
	for i := 0; i < 2; i++ {
		var sample BackendForcePersist
		var result BackendForcePersist
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Errorf(err.Error())
		}
		samples = append(samples, struct {
			a, b BackendForcePersist
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
			t.Errorf("Expected BackendForcePersist to be different in 2 cases, but it is not (%d) %s %s", len(result), a, b)
		}
	}
}

func TestForcePersistEqual(t *testing.T) {
	samples := []struct {
		a, b ForcePersist
	}{}
	for i := 0; i < 2; i++ {
		var sample ForcePersist
		var result ForcePersist
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
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
			a, b ForcePersist
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
			t.Errorf("Expected ForcePersist to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestForcePersistEqualFalse(t *testing.T) {
	samples := []struct {
		a, b ForcePersist
	}{}
	for i := 0; i < 2; i++ {
		var sample ForcePersist
		var result ForcePersist
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Errorf(err.Error())
		}
		samples = append(samples, struct {
			a, b ForcePersist
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
			t.Errorf("Expected ForcePersist to be different, but it is not %s %s", a, b)
		}
	}
}

func TestForcePersistDiff(t *testing.T) {
	samples := []struct {
		a, b ForcePersist
	}{}
	for i := 0; i < 2; i++ {
		var sample ForcePersist
		var result ForcePersist
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
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
			a, b ForcePersist
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
			t.Errorf("Expected ForcePersist to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestForcePersistDiffFalse(t *testing.T) {
	samples := []struct {
		a, b ForcePersist
	}{}
	for i := 0; i < 2; i++ {
		var sample ForcePersist
		var result ForcePersist
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Errorf(err.Error())
		}
		samples = append(samples, struct {
			a, b ForcePersist
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
			t.Errorf("Expected ForcePersist to be different in 2 cases, but it is not (%d) %s %s", len(result), a, b)
		}
	}
}

func TestBackendIgnorePersistEqual(t *testing.T) {
	samples := []struct {
		a, b BackendIgnorePersist
	}{}
	for i := 0; i < 2; i++ {
		var sample BackendIgnorePersist
		var result BackendIgnorePersist
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
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
			a, b BackendIgnorePersist
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
			t.Errorf("Expected BackendIgnorePersist to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestBackendIgnorePersistEqualFalse(t *testing.T) {
	samples := []struct {
		a, b BackendIgnorePersist
	}{}
	for i := 0; i < 2; i++ {
		var sample BackendIgnorePersist
		var result BackendIgnorePersist
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Errorf(err.Error())
		}
		samples = append(samples, struct {
			a, b BackendIgnorePersist
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
			t.Errorf("Expected BackendIgnorePersist to be different, but it is not %s %s", a, b)
		}
	}
}

func TestBackendIgnorePersistDiff(t *testing.T) {
	samples := []struct {
		a, b BackendIgnorePersist
	}{}
	for i := 0; i < 2; i++ {
		var sample BackendIgnorePersist
		var result BackendIgnorePersist
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
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
			a, b BackendIgnorePersist
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
			t.Errorf("Expected BackendIgnorePersist to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestBackendIgnorePersistDiffFalse(t *testing.T) {
	samples := []struct {
		a, b BackendIgnorePersist
	}{}
	for i := 0; i < 2; i++ {
		var sample BackendIgnorePersist
		var result BackendIgnorePersist
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Errorf(err.Error())
		}
		samples = append(samples, struct {
			a, b BackendIgnorePersist
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
			t.Errorf("Expected BackendIgnorePersist to be different in 2 cases, but it is not (%d) %s %s", len(result), a, b)
		}
	}
}

func TestIgnorePersistEqual(t *testing.T) {
	samples := []struct {
		a, b IgnorePersist
	}{}
	for i := 0; i < 2; i++ {
		var sample IgnorePersist
		var result IgnorePersist
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
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
			a, b IgnorePersist
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
			t.Errorf("Expected IgnorePersist to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestIgnorePersistEqualFalse(t *testing.T) {
	samples := []struct {
		a, b IgnorePersist
	}{}
	for i := 0; i < 2; i++ {
		var sample IgnorePersist
		var result IgnorePersist
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Errorf(err.Error())
		}
		samples = append(samples, struct {
			a, b IgnorePersist
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
			t.Errorf("Expected IgnorePersist to be different, but it is not %s %s", a, b)
		}
	}
}

func TestIgnorePersistDiff(t *testing.T) {
	samples := []struct {
		a, b IgnorePersist
	}{}
	for i := 0; i < 2; i++ {
		var sample IgnorePersist
		var result IgnorePersist
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
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
			a, b IgnorePersist
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
			t.Errorf("Expected IgnorePersist to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestIgnorePersistDiffFalse(t *testing.T) {
	samples := []struct {
		a, b IgnorePersist
	}{}
	for i := 0; i < 2; i++ {
		var sample IgnorePersist
		var result IgnorePersist
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Errorf(err.Error())
		}
		samples = append(samples, struct {
			a, b IgnorePersist
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
			t.Errorf("Expected IgnorePersist to be different in 2 cases, but it is not (%d) %s %s", len(result), a, b)
		}
	}
}
