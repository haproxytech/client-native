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

func TestSslOcspResponseEqual(t *testing.T) {
	samples := []struct {
		a, b SslOcspResponse
	}{}
	for i := 0; i < 2; i++ {
		var sample SslOcspResponse
		var result SslOcspResponse
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
			a, b SslOcspResponse
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
			t.Errorf("Expected SslOcspResponse to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestSslOcspResponseEqualFalse(t *testing.T) {
	samples := []struct {
		a, b SslOcspResponse
	}{}
	for i := 0; i < 2; i++ {
		var sample SslOcspResponse
		var result SslOcspResponse
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		result.ProducedAt = strfmt.Date(time.Now().AddDate(rand.Intn(10), rand.Intn(12), rand.Intn(28)))
		samples = append(samples, struct {
			a, b SslOcspResponse
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
			t.Errorf("Expected SslOcspResponse to be different, but it is not %s %s", a, b)
		}
	}
}

func TestSslOcspResponseDiff(t *testing.T) {
	samples := []struct {
		a, b SslOcspResponse
	}{}
	for i := 0; i < 2; i++ {
		var sample SslOcspResponse
		var result SslOcspResponse
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
			a, b SslOcspResponse
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
			t.Errorf("Expected SslOcspResponse to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestSslOcspResponseDiffFalse(t *testing.T) {
	samples := []struct {
		a, b SslOcspResponse
	}{}
	for i := 0; i < 2; i++ {
		var sample SslOcspResponse
		var result SslOcspResponse
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		result.ProducedAt = strfmt.Date(time.Now().AddDate(rand.Intn(10), rand.Intn(12), rand.Intn(28)))
		samples = append(samples, struct {
			a, b SslOcspResponse
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		if len(result) != 7 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Error(err)
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Error(err)
			}
			t.Errorf("Expected SslOcspResponse to be different in 7 cases, but it is not (%d) %s %s", len(result), a, b)
		}
	}
}

func TestOCSPResponsesEqual(t *testing.T) {
	samples := []struct {
		a, b OCSPResponses
	}{}
	for i := 0; i < 2; i++ {
		var sample OCSPResponses
		var result OCSPResponses
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
			a, b OCSPResponses
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
			t.Errorf("Expected OCSPResponses to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestOCSPResponsesEqualFalse(t *testing.T) {
	samples := []struct {
		a, b OCSPResponses
	}{}
	for i := 0; i < 2; i++ {
		var sample OCSPResponses
		var result OCSPResponses
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		result.NextUpdate = strfmt.Date(time.Now().AddDate(rand.Intn(10), rand.Intn(12), rand.Intn(28)))
		result.ThisUpdate = strfmt.Date(time.Now().AddDate(rand.Intn(10), rand.Intn(12), rand.Intn(28)))
		samples = append(samples, struct {
			a, b OCSPResponses
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
			t.Errorf("Expected OCSPResponses to be different, but it is not %s %s", a, b)
		}
	}
}

func TestOCSPResponsesDiff(t *testing.T) {
	samples := []struct {
		a, b OCSPResponses
	}{}
	for i := 0; i < 2; i++ {
		var sample OCSPResponses
		var result OCSPResponses
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
			a, b OCSPResponses
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
			t.Errorf("Expected OCSPResponses to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestOCSPResponsesDiffFalse(t *testing.T) {
	samples := []struct {
		a, b OCSPResponses
	}{}
	for i := 0; i < 2; i++ {
		var sample OCSPResponses
		var result OCSPResponses
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		result.NextUpdate = strfmt.Date(time.Now().AddDate(rand.Intn(10), rand.Intn(12), rand.Intn(28)))
		result.ThisUpdate = strfmt.Date(time.Now().AddDate(rand.Intn(10), rand.Intn(12), rand.Intn(28)))
		samples = append(samples, struct {
			a, b OCSPResponses
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		if len(result) != 5 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Error(err)
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Error(err)
			}
			t.Errorf("Expected OCSPResponses to be different in 5 cases, but it is not (%d) %s %s", len(result), a, b)
		}
	}
}
