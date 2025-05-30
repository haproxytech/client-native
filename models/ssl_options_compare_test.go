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

func TestSslOptionsEqual(t *testing.T) {
	samples := []struct {
		a, b SslOptions
	}{}
	for i := 0; i < 2; i++ {
		var sample SslOptions
		var result SslOptions
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
			a, b SslOptions
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
			t.Errorf("Expected SslOptions to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestSslOptionsEqualFalse(t *testing.T) {
	samples := []struct {
		a, b SslOptions
	}{}
	for i := 0; i < 2; i++ {
		var sample SslOptions
		var result SslOptions
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		result.Maxsslconn = sample.Maxsslconn + 1
		result.Maxsslrate = sample.Maxsslrate + 1
		result.SecurityLevel = Ptr(*sample.SecurityLevel + 1)
		result.SkipSelfIssuedCa = !sample.SkipSelfIssuedCa
		samples = append(samples, struct {
			a, b SslOptions
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
			t.Errorf("Expected SslOptions to be different, but it is not %s %s", a, b)
		}
	}
}

func TestSslOptionsDiff(t *testing.T) {
	samples := []struct {
		a, b SslOptions
	}{}
	for i := 0; i < 2; i++ {
		var sample SslOptions
		var result SslOptions
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
			a, b SslOptions
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
			t.Errorf("Expected SslOptions to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestSslOptionsDiffFalse(t *testing.T) {
	samples := []struct {
		a, b SslOptions
	}{}
	for i := 0; i < 2; i++ {
		var sample SslOptions
		var result SslOptions
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		result.Maxsslconn = sample.Maxsslconn + 1
		result.Maxsslrate = sample.Maxsslrate + 1
		result.SecurityLevel = Ptr(*sample.SecurityLevel + 1)
		result.SkipSelfIssuedCa = !sample.SkipSelfIssuedCa
		samples = append(samples, struct {
			a, b SslOptions
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		if len(result) != 28 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Error(err)
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Error(err)
			}
			t.Errorf("Expected SslOptions to be different in 28 cases, but it is not (%d) %s %s", len(result), a, b)
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
			a, b SslEngine
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
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
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
				t.Error(err)
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Error(err)
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
			a, b SslEngine
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
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
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
				t.Error(err)
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Error(err)
			}
			t.Errorf("Expected SslEngine to be different in 2 cases, but it is not (%d) %s %s", len(result), a, b)
		}
	}
}
