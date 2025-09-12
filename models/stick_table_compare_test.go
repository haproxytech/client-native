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

func TestStickTableEqual(t *testing.T) {
	samples := []struct {
		a, b StickTable
	}{}
	for i := 0; i < 2; i++ {
		var sample StickTable
		var result StickTable
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
			a, b StickTable
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
			t.Errorf("Expected StickTable to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestStickTableEqualFalse(t *testing.T) {
	samples := []struct {
		a, b StickTable
	}{}
	for i := 0; i < 2; i++ {
		var sample StickTable
		var result StickTable
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		result.Size = Ptr(*sample.Size + 1)
		result.Used = Ptr(*sample.Used + 1)
		samples = append(samples, struct {
			a, b StickTable
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
			t.Errorf("Expected StickTable to be different, but it is not %s %s", a, b)
		}
	}
}

func TestStickTableDiff(t *testing.T) {
	samples := []struct {
		a, b StickTable
	}{}
	for i := 0; i < 2; i++ {
		var sample StickTable
		var result StickTable
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
			a, b StickTable
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
			t.Errorf("Expected StickTable to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestStickTableDiffFalse(t *testing.T) {
	samples := []struct {
		a, b StickTable
	}{}
	for i := 0; i < 2; i++ {
		var sample StickTable
		var result StickTable
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		result.Size = Ptr(*sample.Size + 1)
		result.Used = Ptr(*sample.Used + 1)
		samples = append(samples, struct {
			a, b StickTable
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
			t.Errorf("Expected StickTable to be different in 5 cases, but it is not (%d) %s %s", len(result), a, b)
		}
	}
}

func TestStickTableFieldEqual(t *testing.T) {
	samples := []struct {
		a, b StickTableField
	}{}
	for i := 0; i < 2; i++ {
		var sample StickTableField
		var result StickTableField
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
			a, b StickTableField
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
			t.Errorf("Expected StickTableField to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestStickTableFieldEqualFalse(t *testing.T) {
	samples := []struct {
		a, b StickTableField
	}{}
	for i := 0; i < 2; i++ {
		var sample StickTableField
		var result StickTableField
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		result.Idx = sample.Idx + 1
		result.Period = sample.Period + 1
		samples = append(samples, struct {
			a, b StickTableField
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
			t.Errorf("Expected StickTableField to be different, but it is not %s %s", a, b)
		}
	}
}

func TestStickTableFieldDiff(t *testing.T) {
	samples := []struct {
		a, b StickTableField
	}{}
	for i := 0; i < 2; i++ {
		var sample StickTableField
		var result StickTableField
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
			a, b StickTableField
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
			t.Errorf("Expected StickTableField to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestStickTableFieldDiffFalse(t *testing.T) {
	samples := []struct {
		a, b StickTableField
	}{}
	for i := 0; i < 2; i++ {
		var sample StickTableField
		var result StickTableField
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		result.Idx = sample.Idx + 1
		result.Period = sample.Period + 1
		samples = append(samples, struct {
			a, b StickTableField
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		if len(result) != 4 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Error(err)
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Error(err)
			}
			t.Errorf("Expected StickTableField to be different in 4 cases, but it is not (%d) %s %s", len(result), a, b)
		}
	}
}
