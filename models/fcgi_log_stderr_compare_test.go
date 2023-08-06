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

func TestFCGILogStderrEqual(t *testing.T) {
	faker := gofakeit.NewCrypto()
	gofakeit.SetGlobalFaker(faker)
	samples := []struct {
		a, b FCGILogStderr
	}{}
	for i := 0; i < 100; i++ {
		var sample FCGILogStderr
		var result FCGILogStderr
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
			a, b FCGILogStderr
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
			t.Errorf("Expected FCGILogStderr to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestFCGILogStderrEqualFalse(t *testing.T) {
	faker := gofakeit.NewCrypto()
	gofakeit.SetGlobalFaker(faker)
	samples := []struct {
		a, b FCGILogStderr
	}{}
	for i := 0; i < 100; i++ {
		var sample FCGILogStderr
		var result FCGILogStderr
		err := gofakeit.Struct(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = gofakeit.Struct(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		result.Global = !sample.Global
		result.Len = sample.Len + 1
		samples = append(samples, struct {
			a, b FCGILogStderr
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
			t.Errorf("Expected FCGILogStderr to be different, but it is not %s %s", a, b)
		}
	}
}

func TestFCGILogStderrDiff(t *testing.T) {
	faker := gofakeit.NewCrypto()
	gofakeit.SetGlobalFaker(faker)
	samples := []struct {
		a, b FCGILogStderr
	}{}
	for i := 0; i < 100; i++ {
		var sample FCGILogStderr
		var result FCGILogStderr
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
			a, b FCGILogStderr
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
			t.Errorf("Expected FCGILogStderr to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestFCGILogStderrDiffFalse(t *testing.T) {
	faker := gofakeit.NewCrypto()
	gofakeit.SetGlobalFaker(faker)
	samples := []struct {
		a, b FCGILogStderr
	}{}
	for i := 0; i < 100; i++ {
		var sample FCGILogStderr
		var result FCGILogStderr
		err := gofakeit.Struct(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = gofakeit.Struct(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		result.Global = !sample.Global
		result.Len = sample.Len + 1
		samples = append(samples, struct {
			a, b FCGILogStderr
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
			t.Errorf("Expected FCGILogStderr to be different in 8 cases, but it is not (%d) %s %s", len(result), a, b)
		}
	}
}

func TestFCGILogStderrSampleEqual(t *testing.T) {
	faker := gofakeit.NewCrypto()
	gofakeit.SetGlobalFaker(faker)
	samples := []struct {
		a, b FCGILogStderrSample
	}{}
	for i := 0; i < 100; i++ {
		var sample FCGILogStderrSample
		var result FCGILogStderrSample
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
			a, b FCGILogStderrSample
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
			t.Errorf("Expected FCGILogStderrSample to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestFCGILogStderrSampleEqualFalse(t *testing.T) {
	faker := gofakeit.NewCrypto()
	gofakeit.SetGlobalFaker(faker)
	samples := []struct {
		a, b FCGILogStderrSample
	}{}
	for i := 0; i < 100; i++ {
		var sample FCGILogStderrSample
		var result FCGILogStderrSample
		err := gofakeit.Struct(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = gofakeit.Struct(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		result.Size = Ptr(*sample.Size + 1)
		samples = append(samples, struct {
			a, b FCGILogStderrSample
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
			t.Errorf("Expected FCGILogStderrSample to be different, but it is not %s %s", a, b)
		}
	}
}

func TestFCGILogStderrSampleDiff(t *testing.T) {
	faker := gofakeit.NewCrypto()
	gofakeit.SetGlobalFaker(faker)
	samples := []struct {
		a, b FCGILogStderrSample
	}{}
	for i := 0; i < 100; i++ {
		var sample FCGILogStderrSample
		var result FCGILogStderrSample
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
			a, b FCGILogStderrSample
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
			t.Errorf("Expected FCGILogStderrSample to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestFCGILogStderrSampleDiffFalse(t *testing.T) {
	faker := gofakeit.NewCrypto()
	gofakeit.SetGlobalFaker(faker)
	samples := []struct {
		a, b FCGILogStderrSample
	}{}
	for i := 0; i < 100; i++ {
		var sample FCGILogStderrSample
		var result FCGILogStderrSample
		err := gofakeit.Struct(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = gofakeit.Struct(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		result.Size = Ptr(*sample.Size + 1)
		samples = append(samples, struct {
			a, b FCGILogStderrSample
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
			t.Errorf("Expected FCGILogStderrSample to be different in 2 cases, but it is not (%d) %s %s", len(result), a, b)
		}
	}
}
