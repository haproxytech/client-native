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

func TestHTTPRequestRuleEqual(t *testing.T) {
	samples := []struct {
		a, b HTTPRequestRule
	}{}
	for i := 0; i < 2; i++ {
		var sample HTTPRequestRule
		var result HTTPRequestRule
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
			a, b HTTPRequestRule
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
			t.Errorf("Expected HTTPRequestRule to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestHTTPRequestRuleEqualFalse(t *testing.T) {
	samples := []struct {
		a, b HTTPRequestRule
	}{}
	for i := 0; i < 2; i++ {
		var sample HTTPRequestRule
		var result HTTPRequestRule
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Errorf(err.Error())
		}
		result.CaptureID = Ptr(*sample.CaptureID + 1)
		result.CaptureLen = sample.CaptureLen + 1
		result.DenyStatus = Ptr(*sample.DenyStatus + 1)
		result.NiceValue = sample.NiceValue + 1
		result.NormalizerFull = !sample.NormalizerFull
		result.NormalizerStrict = !sample.NormalizerStrict
		result.RedirCode = Ptr(*sample.RedirCode + 1)
		result.ReturnStatusCode = Ptr(*sample.ReturnStatusCode + 1)
		result.RstTTL = sample.RstTTL + 1
		result.ScID = sample.ScID + 1
		result.ScIdx = sample.ScIdx + 1
		result.ScInt = Ptr(*sample.ScInt + 1)
		result.TrackScStickCounter = Ptr(*sample.TrackScStickCounter + 1)
		result.WaitAtLeast = Ptr(*sample.WaitAtLeast + 1)
		result.WaitTime = Ptr(*sample.WaitTime + 1)
		samples = append(samples, struct {
			a, b HTTPRequestRule
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
			t.Errorf("Expected HTTPRequestRule to be different, but it is not %s %s", a, b)
		}
	}
}

func TestHTTPRequestRuleDiff(t *testing.T) {
	samples := []struct {
		a, b HTTPRequestRule
	}{}
	for i := 0; i < 2; i++ {
		var sample HTTPRequestRule
		var result HTTPRequestRule
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
			a, b HTTPRequestRule
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
			t.Errorf("Expected HTTPRequestRule to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestHTTPRequestRuleDiffFalse(t *testing.T) {
	samples := []struct {
		a, b HTTPRequestRule
	}{}
	for i := 0; i < 2; i++ {
		var sample HTTPRequestRule
		var result HTTPRequestRule
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Errorf(err.Error())
		}
		result.CaptureID = Ptr(*sample.CaptureID + 1)
		result.CaptureLen = sample.CaptureLen + 1
		result.DenyStatus = Ptr(*sample.DenyStatus + 1)
		result.NiceValue = sample.NiceValue + 1
		result.NormalizerFull = !sample.NormalizerFull
		result.NormalizerStrict = !sample.NormalizerStrict
		result.RedirCode = Ptr(*sample.RedirCode + 1)
		result.ReturnStatusCode = Ptr(*sample.ReturnStatusCode + 1)
		result.RstTTL = sample.RstTTL + 1
		result.ScID = sample.ScID + 1
		result.ScIdx = sample.ScIdx + 1
		result.ScInt = Ptr(*sample.ScInt + 1)
		result.TrackScStickCounter = Ptr(*sample.TrackScStickCounter + 1)
		result.WaitAtLeast = Ptr(*sample.WaitAtLeast + 1)
		result.WaitTime = Ptr(*sample.WaitTime + 1)
		samples = append(samples, struct {
			a, b HTTPRequestRule
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		if len(result) != 70 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected HTTPRequestRule to be different in 70 cases, but it is not (%d) %s %s", len(result), a, b)
		}
	}
}
