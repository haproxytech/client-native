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

func TestSslCrlEntryEqual(t *testing.T) {
	samples := []struct {
		a, b SslCrlEntry
	}{}
	for i := 0; i < 2; i++ {
		var sample SslCrlEntry
		var result SslCrlEntry
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
			a, b SslCrlEntry
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
			t.Errorf("Expected SslCrlEntry to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestSslCrlEntryEqualFalse(t *testing.T) {
	samples := []struct {
		a, b SslCrlEntry
	}{}
	for i := 0; i < 2; i++ {
		var sample SslCrlEntry
		var result SslCrlEntry
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		result.LastUpdate = strfmt.Date(time.Now().AddDate(rand.Intn(10), rand.Intn(12), rand.Intn(28)))
		result.NextUpdate = strfmt.Date(time.Now().AddDate(rand.Intn(10), rand.Intn(12), rand.Intn(28)))
		samples = append(samples, struct {
			a, b SslCrlEntry
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
			t.Errorf("Expected SslCrlEntry to be different, but it is not %s %s", a, b)
		}
	}
}

func TestSslCrlEntryDiff(t *testing.T) {
	samples := []struct {
		a, b SslCrlEntry
	}{}
	for i := 0; i < 2; i++ {
		var sample SslCrlEntry
		var result SslCrlEntry
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
			a, b SslCrlEntry
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
			t.Errorf("Expected SslCrlEntry to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestSslCrlEntryDiffFalse(t *testing.T) {
	samples := []struct {
		a, b SslCrlEntry
	}{}
	for i := 0; i < 2; i++ {
		var sample SslCrlEntry
		var result SslCrlEntry
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		result.LastUpdate = strfmt.Date(time.Now().AddDate(rand.Intn(10), rand.Intn(12), rand.Intn(28)))
		result.NextUpdate = strfmt.Date(time.Now().AddDate(rand.Intn(10), rand.Intn(12), rand.Intn(28)))
		samples = append(samples, struct {
			a, b SslCrlEntry
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		listDiffFields := GetListOfDiffFields(result)
		if len(listDiffFields) != 8 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Error(err)
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Error(err)
			}
			t.Errorf("Expected SslCrlEntry to be different in 8 cases, but it is not (%d) %s %s", len(result), a, b)
		}
	}
}

func TestRevokedCertificatesEqual(t *testing.T) {
	samples := []struct {
		a, b RevokedCertificates
	}{}
	for i := 0; i < 2; i++ {
		var sample RevokedCertificates
		var result RevokedCertificates
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
			a, b RevokedCertificates
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
			t.Errorf("Expected RevokedCertificates to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestRevokedCertificatesEqualFalse(t *testing.T) {
	samples := []struct {
		a, b RevokedCertificates
	}{}
	for i := 0; i < 2; i++ {
		var sample RevokedCertificates
		var result RevokedCertificates
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		result.RevocationDate = strfmt.Date(time.Now().AddDate(rand.Intn(10), rand.Intn(12), rand.Intn(28)))
		samples = append(samples, struct {
			a, b RevokedCertificates
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
			t.Errorf("Expected RevokedCertificates to be different, but it is not %s %s", a, b)
		}
	}
}

func TestRevokedCertificatesDiff(t *testing.T) {
	samples := []struct {
		a, b RevokedCertificates
	}{}
	for i := 0; i < 2; i++ {
		var sample RevokedCertificates
		var result RevokedCertificates
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
			a, b RevokedCertificates
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
			t.Errorf("Expected RevokedCertificates to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestRevokedCertificatesDiffFalse(t *testing.T) {
	samples := []struct {
		a, b RevokedCertificates
	}{}
	for i := 0; i < 2; i++ {
		var sample RevokedCertificates
		var result RevokedCertificates
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		result.RevocationDate = strfmt.Date(time.Now().AddDate(rand.Intn(10), rand.Intn(12), rand.Intn(28)))
		samples = append(samples, struct {
			a, b RevokedCertificates
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
			t.Errorf("Expected RevokedCertificates to be different in 2 cases, but it is not (%d) %s %s", len(result), a, b)
		}
	}
}
