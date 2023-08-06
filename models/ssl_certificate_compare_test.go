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
	"math/rand"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/go-openapi/strfmt"

	jsoniter "github.com/json-iterator/go"
)

func TestSslCertificateEqual(t *testing.T) {
	faker := gofakeit.NewCrypto()
	gofakeit.SetGlobalFaker(faker)
	samples := []struct {
		a, b SslCertificate
	}{}
	for i := 0; i < 100; i++ {
		var sample SslCertificate
		var result SslCertificate
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
			a, b SslCertificate
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
			t.Errorf("Expected SslCertificate to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestSslCertificateEqualFalse(t *testing.T) {
	faker := gofakeit.NewCrypto()
	gofakeit.SetGlobalFaker(faker)
	samples := []struct {
		a, b SslCertificate
	}{}
	for i := 0; i < 100; i++ {
		var sample SslCertificate
		var result SslCertificate
		err := gofakeit.Struct(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = gofakeit.Struct(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		result.NotAfter = strfmt.Date(time.Now().AddDate(rand.Intn(10), rand.Intn(12), rand.Intn(28)))
		result.NotBefore = strfmt.Date(time.Now().AddDate(rand.Intn(10), rand.Intn(12), rand.Intn(28)))
		result.Size = sample.Size + 1
		samples = append(samples, struct {
			a, b SslCertificate
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
			t.Errorf("Expected SslCertificate to be different, but it is not %s %s", a, b)
		}
	}
}

func TestSslCertificateDiff(t *testing.T) {
	faker := gofakeit.NewCrypto()
	gofakeit.SetGlobalFaker(faker)
	samples := []struct {
		a, b SslCertificate
	}{}
	for i := 0; i < 100; i++ {
		var sample SslCertificate
		var result SslCertificate
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
			a, b SslCertificate
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
			t.Errorf("Expected SslCertificate to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestSslCertificateDiffFalse(t *testing.T) {
	faker := gofakeit.NewCrypto()
	gofakeit.SetGlobalFaker(faker)
	samples := []struct {
		a, b SslCertificate
	}{}
	for i := 0; i < 100; i++ {
		var sample SslCertificate
		var result SslCertificate
		err := gofakeit.Struct(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = gofakeit.Struct(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		result.NotAfter = strfmt.Date(time.Now().AddDate(rand.Intn(10), rand.Intn(12), rand.Intn(28)))
		result.NotBefore = strfmt.Date(time.Now().AddDate(rand.Intn(10), rand.Intn(12), rand.Intn(28)))
		result.Size = sample.Size + 1
		samples = append(samples, struct {
			a, b SslCertificate
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		if len(result) != 9 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected SslCertificate to be different in 9 cases, but it is not (%d) %s %s", len(result), a, b)
		}
	}
}
