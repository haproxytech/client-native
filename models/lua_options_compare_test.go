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

func TestLuaOptionsEqual(t *testing.T) {
	samples := []struct {
		a, b LuaOptions
	}{}
	for i := 0; i < 2; i++ {
		var sample LuaOptions
		var result LuaOptions
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
			a, b LuaOptions
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
			t.Errorf("Expected LuaOptions to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestLuaOptionsEqualFalse(t *testing.T) {
	samples := []struct {
		a, b LuaOptions
	}{}
	for i := 0; i < 2; i++ {
		var sample LuaOptions
		var result LuaOptions
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		samples = append(samples, struct {
			a, b LuaOptions
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
			t.Errorf("Expected LuaOptions to be different, but it is not %s %s", a, b)
		}
	}
}

func TestLuaOptionsDiff(t *testing.T) {
	samples := []struct {
		a, b LuaOptions
	}{}
	for i := 0; i < 2; i++ {
		var sample LuaOptions
		var result LuaOptions
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
			a, b LuaOptions
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
			t.Errorf("Expected LuaOptions to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestLuaOptionsDiffFalse(t *testing.T) {
	samples := []struct {
		a, b LuaOptions
	}{}
	for i := 0; i < 2; i++ {
		var sample LuaOptions
		var result LuaOptions
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		samples = append(samples, struct {
			a, b LuaOptions
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		if len(result) != 3 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Error(err)
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Error(err)
			}
			t.Errorf("Expected LuaOptions to be different in 3 cases, but it is not (%d) %s %s", len(result), a, b)
		}
	}
}

func TestLuaLoadEqual(t *testing.T) {
	samples := []struct {
		a, b LuaLoad
	}{}
	for i := 0; i < 2; i++ {
		var sample LuaLoad
		var result LuaLoad
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
			a, b LuaLoad
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
			t.Errorf("Expected LuaLoad to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestLuaLoadEqualFalse(t *testing.T) {
	samples := []struct {
		a, b LuaLoad
	}{}
	for i := 0; i < 2; i++ {
		var sample LuaLoad
		var result LuaLoad
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		samples = append(samples, struct {
			a, b LuaLoad
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
			t.Errorf("Expected LuaLoad to be different, but it is not %s %s", a, b)
		}
	}
}

func TestLuaLoadDiff(t *testing.T) {
	samples := []struct {
		a, b LuaLoad
	}{}
	for i := 0; i < 2; i++ {
		var sample LuaLoad
		var result LuaLoad
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
			a, b LuaLoad
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
			t.Errorf("Expected LuaLoad to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestLuaLoadDiffFalse(t *testing.T) {
	samples := []struct {
		a, b LuaLoad
	}{}
	for i := 0; i < 2; i++ {
		var sample LuaLoad
		var result LuaLoad
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		samples = append(samples, struct {
			a, b LuaLoad
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		if len(result) != 1 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Error(err)
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Error(err)
			}
			t.Errorf("Expected LuaLoad to be different in 1 cases, but it is not (%d) %s %s", len(result), a, b)
		}
	}
}

func TestLuaPrependPathEqual(t *testing.T) {
	samples := []struct {
		a, b LuaPrependPath
	}{}
	for i := 0; i < 2; i++ {
		var sample LuaPrependPath
		var result LuaPrependPath
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
			a, b LuaPrependPath
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
			t.Errorf("Expected LuaPrependPath to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestLuaPrependPathEqualFalse(t *testing.T) {
	samples := []struct {
		a, b LuaPrependPath
	}{}
	for i := 0; i < 2; i++ {
		var sample LuaPrependPath
		var result LuaPrependPath
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		samples = append(samples, struct {
			a, b LuaPrependPath
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
			t.Errorf("Expected LuaPrependPath to be different, but it is not %s %s", a, b)
		}
	}
}

func TestLuaPrependPathDiff(t *testing.T) {
	samples := []struct {
		a, b LuaPrependPath
	}{}
	for i := 0; i < 2; i++ {
		var sample LuaPrependPath
		var result LuaPrependPath
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
			a, b LuaPrependPath
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
			t.Errorf("Expected LuaPrependPath to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestLuaPrependPathDiffFalse(t *testing.T) {
	samples := []struct {
		a, b LuaPrependPath
	}{}
	for i := 0; i < 2; i++ {
		var sample LuaPrependPath
		var result LuaPrependPath
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Error(err)
		}
		samples = append(samples, struct {
			a, b LuaPrependPath
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
			t.Errorf("Expected LuaPrependPath to be different in 2 cases, but it is not (%d) %s %s", len(result), a, b)
		}
	}
}
