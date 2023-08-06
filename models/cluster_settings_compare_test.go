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
	"bytes"
	"encoding/gob"
	"testing"

	"github.com/brianvoe/gofakeit/v6"

	jsoniter "github.com/json-iterator/go"
)

func TestClusterSettingsEqual(t *testing.T) {
	faker := gofakeit.NewCrypto()
	gofakeit.SetGlobalFaker(faker)
	samples := []struct {
		a, b ClusterSettings
	}{}
	for i := 0; i < 100; i++ {
		var sample ClusterSettings
		var result ClusterSettings
		err := gofakeit.Struct(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		buf := new(bytes.Buffer)
		enc := gob.NewEncoder(buf)
		err = enc.Encode(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		dec := gob.NewDecoder(buf)
		err = dec.Decode(&result)
		if err != nil {
			t.Errorf(err.Error())
		}

		samples = append(samples, struct {
			a, b ClusterSettings
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
			t.Errorf("Expected ClusterSettings to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestClusterSettingsEqualFalse(t *testing.T) {
	faker := gofakeit.NewCrypto()
	gofakeit.SetGlobalFaker(faker)
	samples := []struct {
		a, b ClusterSettings
	}{}
	for i := 0; i < 100; i++ {
		var sample ClusterSettings
		var result ClusterSettings
		err := gofakeit.Struct(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = gofakeit.Struct(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		samples = append(samples, struct {
			a, b ClusterSettings
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
			t.Errorf("Expected ClusterSettings to be different, but it is not %s %s", a, b)
		}
	}
}

func TestClusterSettingsDiff(t *testing.T) {
	faker := gofakeit.NewCrypto()
	gofakeit.SetGlobalFaker(faker)
	samples := []struct {
		a, b ClusterSettings
	}{}
	for i := 0; i < 100; i++ {
		var sample ClusterSettings
		var result ClusterSettings
		err := gofakeit.Struct(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		buf := new(bytes.Buffer)
		enc := gob.NewEncoder(buf)
		err = enc.Encode(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		dec := gob.NewDecoder(buf)
		err = dec.Decode(&result)
		if err != nil {
			t.Errorf(err.Error())
		}

		samples = append(samples, struct {
			a, b ClusterSettings
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
			t.Errorf("Expected ClusterSettings to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestClusterSettingsDiffFalse(t *testing.T) {
	faker := gofakeit.NewCrypto()
	gofakeit.SetGlobalFaker(faker)
	samples := []struct {
		a, b ClusterSettings
	}{}
	for i := 0; i < 100; i++ {
		var sample ClusterSettings
		var result ClusterSettings
		err := gofakeit.Struct(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = gofakeit.Struct(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		samples = append(samples, struct {
			a, b ClusterSettings
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
			t.Errorf("Expected ClusterSettings to be different in 4 cases, but it is not (%d) %s %s", len(result), a, b)
		}
	}
}

func TestClusterSettingsClusterEqual(t *testing.T) {
	faker := gofakeit.NewCrypto()
	gofakeit.SetGlobalFaker(faker)
	samples := []struct {
		a, b ClusterSettingsCluster
	}{}
	for i := 0; i < 100; i++ {
		var sample ClusterSettingsCluster
		var result ClusterSettingsCluster
		err := gofakeit.Struct(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		buf := new(bytes.Buffer)
		enc := gob.NewEncoder(buf)
		err = enc.Encode(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		dec := gob.NewDecoder(buf)
		err = dec.Decode(&result)
		if err != nil {
			t.Errorf(err.Error())
		}

		samples = append(samples, struct {
			a, b ClusterSettingsCluster
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
			t.Errorf("Expected ClusterSettingsCluster to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestClusterSettingsClusterEqualFalse(t *testing.T) {
	faker := gofakeit.NewCrypto()
	gofakeit.SetGlobalFaker(faker)
	samples := []struct {
		a, b ClusterSettingsCluster
	}{}
	for i := 0; i < 100; i++ {
		var sample ClusterSettingsCluster
		var result ClusterSettingsCluster
		err := gofakeit.Struct(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = gofakeit.Struct(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		result.Port = Ptr(*sample.Port + 1)
		samples = append(samples, struct {
			a, b ClusterSettingsCluster
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
			t.Errorf("Expected ClusterSettingsCluster to be different, but it is not %s %s", a, b)
		}
	}
}

func TestClusterSettingsClusterDiff(t *testing.T) {
	faker := gofakeit.NewCrypto()
	gofakeit.SetGlobalFaker(faker)
	samples := []struct {
		a, b ClusterSettingsCluster
	}{}
	for i := 0; i < 100; i++ {
		var sample ClusterSettingsCluster
		var result ClusterSettingsCluster
		err := gofakeit.Struct(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		buf := new(bytes.Buffer)
		enc := gob.NewEncoder(buf)
		err = enc.Encode(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		dec := gob.NewDecoder(buf)
		err = dec.Decode(&result)
		if err != nil {
			t.Errorf(err.Error())
		}

		samples = append(samples, struct {
			a, b ClusterSettingsCluster
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
			t.Errorf("Expected ClusterSettingsCluster to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestClusterSettingsClusterDiffFalse(t *testing.T) {
	faker := gofakeit.NewCrypto()
	gofakeit.SetGlobalFaker(faker)
	samples := []struct {
		a, b ClusterSettingsCluster
	}{}
	for i := 0; i < 100; i++ {
		var sample ClusterSettingsCluster
		var result ClusterSettingsCluster
		err := gofakeit.Struct(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = gofakeit.Struct(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		result.Port = Ptr(*sample.Port + 1)
		samples = append(samples, struct {
			a, b ClusterSettingsCluster
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		if len(result) != 7 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected ClusterSettingsCluster to be different in 7 cases, but it is not (%d) %s %s", len(result), a, b)
		}
	}
}

func TestClusterLogTargetEqual(t *testing.T) {
	faker := gofakeit.NewCrypto()
	gofakeit.SetGlobalFaker(faker)
	samples := []struct {
		a, b ClusterLogTarget
	}{}
	for i := 0; i < 100; i++ {
		var sample ClusterLogTarget
		var result ClusterLogTarget
		err := gofakeit.Struct(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		buf := new(bytes.Buffer)
		enc := gob.NewEncoder(buf)
		err = enc.Encode(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		dec := gob.NewDecoder(buf)
		err = dec.Decode(&result)
		if err != nil {
			t.Errorf(err.Error())
		}

		samples = append(samples, struct {
			a, b ClusterLogTarget
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
			t.Errorf("Expected ClusterLogTarget to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestClusterLogTargetEqualFalse(t *testing.T) {
	faker := gofakeit.NewCrypto()
	gofakeit.SetGlobalFaker(faker)
	samples := []struct {
		a, b ClusterLogTarget
	}{}
	for i := 0; i < 100; i++ {
		var sample ClusterLogTarget
		var result ClusterLogTarget
		err := gofakeit.Struct(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = gofakeit.Struct(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		result.Port = Ptr(*sample.Port + 1)
		samples = append(samples, struct {
			a, b ClusterLogTarget
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
			t.Errorf("Expected ClusterLogTarget to be different, but it is not %s %s", a, b)
		}
	}
}

func TestClusterLogTargetDiff(t *testing.T) {
	faker := gofakeit.NewCrypto()
	gofakeit.SetGlobalFaker(faker)
	samples := []struct {
		a, b ClusterLogTarget
	}{}
	for i := 0; i < 100; i++ {
		var sample ClusterLogTarget
		var result ClusterLogTarget
		err := gofakeit.Struct(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		buf := new(bytes.Buffer)
		enc := gob.NewEncoder(buf)
		err = enc.Encode(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		dec := gob.NewDecoder(buf)
		err = dec.Decode(&result)
		if err != nil {
			t.Errorf(err.Error())
		}

		samples = append(samples, struct {
			a, b ClusterLogTarget
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
			t.Errorf("Expected ClusterLogTarget to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestClusterLogTargetDiffFalse(t *testing.T) {
	faker := gofakeit.NewCrypto()
	gofakeit.SetGlobalFaker(faker)
	samples := []struct {
		a, b ClusterLogTarget
	}{}
	for i := 0; i < 100; i++ {
		var sample ClusterLogTarget
		var result ClusterLogTarget
		err := gofakeit.Struct(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = gofakeit.Struct(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		result.Port = Ptr(*sample.Port + 1)
		samples = append(samples, struct {
			a, b ClusterLogTarget
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
			t.Errorf("Expected ClusterLogTarget to be different in 4 cases, but it is not (%d) %s %s", len(result), a, b)
		}
	}
}
