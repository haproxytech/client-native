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

func TestFilterConfigEqual(t *testing.T) {
	samples := []struct {
		a, b FilterConfig
	}{}
	for i := 0; i < 2; i++ {
		var sample FilterConfig
		var result FilterConfig
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
			a, b FilterConfig
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
			t.Errorf("Expected FilterConfig to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestFilterConfigEqualFalse(t *testing.T) {
	samples := []struct {
		a, b FilterConfig
	}{}
	for i := 0; i < 2; i++ {
		var sample FilterConfig
		var result FilterConfig
		// Non-empty slices and maps: if both sides got an empty one, that
		// field would compare equal and the diff field count would be short.
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true), options.WithRandomMapAndSliceMinSize(1))
		if err != nil {
			t.Error(err)
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true), options.WithRandomMapAndSliceMinSize(1))
		if err != nil {
			t.Error(err)
		}
		result.DefaultLimit = sample.DefaultLimit + 1
		result.DefaultPeriod = sample.DefaultPeriod + 1
		result.Enabled = !sample.Enabled
		result.Limit = sample.Limit + 1
		result.MinSize = sample.MinSize + 1
		result.TraceHexdump = !sample.TraceHexdump
		result.TraceMaxFwd = Ptr(*sample.TraceMaxFwd + 1)
		result.TraceQuiet = !sample.TraceQuiet
		result.TraceRndForwarding = !sample.TraceRndForwarding
		result.TraceRndParsing = !sample.TraceRndParsing
		samples = append(samples, struct {
			a, b FilterConfig
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
			t.Errorf("Expected FilterConfig to be different, but it is not %s %s", a, b)
		}
	}
}

func TestFilterConfigDiff(t *testing.T) {
	samples := []struct {
		a, b FilterConfig
	}{}
	for i := 0; i < 2; i++ {
		var sample FilterConfig
		var result FilterConfig
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
			a, b FilterConfig
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
			t.Errorf("Expected FilterConfig to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestFilterConfigDiffFalse(t *testing.T) {
	samples := []struct {
		a, b FilterConfig
	}{}
	for i := 0; i < 2; i++ {
		var sample FilterConfig
		var result FilterConfig
		// Non-empty slices and maps: if both sides got an empty one, that
		// field would compare equal and the diff field count would be short.
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true), options.WithRandomMapAndSliceMinSize(1))
		if err != nil {
			t.Error(err)
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true), options.WithRandomMapAndSliceMinSize(1))
		if err != nil {
			t.Error(err)
		}
		result.DefaultLimit = sample.DefaultLimit + 1
		result.DefaultPeriod = sample.DefaultPeriod + 1
		result.Enabled = !sample.Enabled
		result.Limit = sample.Limit + 1
		result.MinSize = sample.MinSize + 1
		result.TraceHexdump = !sample.TraceHexdump
		result.TraceMaxFwd = Ptr(*sample.TraceMaxFwd + 1)
		result.TraceQuiet = !sample.TraceQuiet
		result.TraceRndForwarding = !sample.TraceRndForwarding
		result.TraceRndParsing = !sample.TraceRndParsing
		samples = append(samples, struct {
			a, b FilterConfig
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		listDiffFields := GetListOfDiffFields(result)
		if len(listDiffFields) != 17 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Error(err)
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Error(err)
			}
			t.Errorf("Expected FilterConfig to be different in 17 cases, but it is not (%d) %s %s", len(listDiffFields), a, b)
		}
	}
}
