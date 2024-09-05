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

func TestCookieEqual(t *testing.T) {
	samples := []struct {
		a, b Cookie
	}{}
	for i := 0; i < 2; i++ {
		var sample Cookie
		var result Cookie
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
			a, b Cookie
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
			t.Errorf("Expected Cookie to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestCookieEqualFalse(t *testing.T) {
	samples := []struct {
		a, b Cookie
	}{}
	for i := 0; i < 2; i++ {
		var sample Cookie
		var result Cookie
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Errorf(err.Error())
		}
		result.Dynamic = !sample.Dynamic
		result.Httponly = !sample.Httponly
		result.Indirect = !sample.Indirect
		result.Maxidle = sample.Maxidle + 1
		result.Maxlife = sample.Maxlife + 1
		result.Nocache = !sample.Nocache
		result.Postonly = !sample.Postonly
		result.Preserve = !sample.Preserve
		result.Secure = !sample.Secure
		samples = append(samples, struct {
			a, b Cookie
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
			t.Errorf("Expected Cookie to be different, but it is not %s %s", a, b)
		}
	}
}

func TestCookieDiff(t *testing.T) {
	samples := []struct {
		a, b Cookie
	}{}
	for i := 0; i < 2; i++ {
		var sample Cookie
		var result Cookie
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
			a, b Cookie
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
			t.Errorf("Expected Cookie to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestCookieDiffFalse(t *testing.T) {
	samples := []struct {
		a, b Cookie
	}{}
	for i := 0; i < 2; i++ {
		var sample Cookie
		var result Cookie
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Errorf(err.Error())
		}
		result.Dynamic = !sample.Dynamic
		result.Httponly = !sample.Httponly
		result.Indirect = !sample.Indirect
		result.Maxidle = sample.Maxidle + 1
		result.Maxlife = sample.Maxlife + 1
		result.Nocache = !sample.Nocache
		result.Postonly = !sample.Postonly
		result.Preserve = !sample.Preserve
		result.Secure = !sample.Secure
		samples = append(samples, struct {
			a, b Cookie
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		if len(result) != 13 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected Cookie to be different in 13 cases, but it is not (%d) %s %s", len(result), a, b)
		}
	}
}

func TestAttrEqual(t *testing.T) {
	samples := []struct {
		a, b Attr
	}{}
	for i := 0; i < 2; i++ {
		var sample Attr
		var result Attr
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
			a, b Attr
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
			t.Errorf("Expected Attr to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestAttrEqualFalse(t *testing.T) {
	samples := []struct {
		a, b Attr
	}{}
	for i := 0; i < 2; i++ {
		var sample Attr
		var result Attr
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Errorf(err.Error())
		}
		samples = append(samples, struct {
			a, b Attr
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
			t.Errorf("Expected Attr to be different, but it is not %s %s", a, b)
		}
	}
}

func TestAttrDiff(t *testing.T) {
	samples := []struct {
		a, b Attr
	}{}
	for i := 0; i < 2; i++ {
		var sample Attr
		var result Attr
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
			a, b Attr
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
			t.Errorf("Expected Attr to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestAttrDiffFalse(t *testing.T) {
	samples := []struct {
		a, b Attr
	}{}
	for i := 0; i < 2; i++ {
		var sample Attr
		var result Attr
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Errorf(err.Error())
		}
		samples = append(samples, struct {
			a, b Attr
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		if len(result) != 1 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected Attr to be different in 1 cases, but it is not (%d) %s %s", len(result), a, b)
		}
	}
}

func TestDomainEqual(t *testing.T) {
	samples := []struct {
		a, b Domain
	}{}
	for i := 0; i < 2; i++ {
		var sample Domain
		var result Domain
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
			a, b Domain
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
			t.Errorf("Expected Domain to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestDomainEqualFalse(t *testing.T) {
	samples := []struct {
		a, b Domain
	}{}
	for i := 0; i < 2; i++ {
		var sample Domain
		var result Domain
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Errorf(err.Error())
		}
		samples = append(samples, struct {
			a, b Domain
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
			t.Errorf("Expected Domain to be different, but it is not %s %s", a, b)
		}
	}
}

func TestDomainDiff(t *testing.T) {
	samples := []struct {
		a, b Domain
	}{}
	for i := 0; i < 2; i++ {
		var sample Domain
		var result Domain
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
			a, b Domain
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
			t.Errorf("Expected Domain to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestDomainDiffFalse(t *testing.T) {
	samples := []struct {
		a, b Domain
	}{}
	for i := 0; i < 2; i++ {
		var sample Domain
		var result Domain
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Errorf(err.Error())
		}
		samples = append(samples, struct {
			a, b Domain
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		if len(result) != 1 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected Domain to be different in 1 cases, but it is not (%d) %s %s", len(result), a, b)
		}
	}
}
