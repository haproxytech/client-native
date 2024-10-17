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

func TestFrontendBaseEqual(t *testing.T) {
	samples := []struct {
		a, b FrontendBase
	}{}
	for i := 0; i < 2; i++ {
		var sample FrontendBase
		var result FrontendBase
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
			a, b FrontendBase
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
			t.Errorf("Expected FrontendBase to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestFrontendBaseEqualFalse(t *testing.T) {
	samples := []struct {
		a, b FrontendBase
	}{}
	for i := 0; i < 2; i++ {
		var sample FrontendBase
		var result FrontendBase
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Errorf(err.Error())
		}
		result.Backlog = Ptr(*sample.Backlog + 1)
		result.Clflog = !sample.Clflog
		result.ClientFinTimeout = Ptr(*sample.ClientFinTimeout + 1)
		result.ClientTimeout = Ptr(*sample.ClientTimeout + 1)
		result.ClitcpkaCnt = Ptr(*sample.ClitcpkaCnt + 1)
		result.ClitcpkaIdle = Ptr(*sample.ClitcpkaIdle + 1)
		result.ClitcpkaIntvl = Ptr(*sample.ClitcpkaIntvl + 1)
		result.Disabled = !sample.Disabled
		result.Enabled = !sample.Enabled
		result.HTTPKeepAliveTimeout = Ptr(*sample.HTTPKeepAliveTimeout + 1)
		result.HTTPRequestTimeout = Ptr(*sample.HTTPRequestTimeout + 1)
		result.Httplog = !sample.Httplog
		result.ID = Ptr(*sample.ID + 1)
		result.Maxconn = Ptr(*sample.Maxconn + 1)
		result.TarpitTimeout = Ptr(*sample.TarpitTimeout + 1)
		result.Tcplog = !sample.Tcplog
		samples = append(samples, struct {
			a, b FrontendBase
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
			t.Errorf("Expected FrontendBase to be different, but it is not %s %s", a, b)
		}
	}
}

func TestFrontendBaseDiff(t *testing.T) {
	samples := []struct {
		a, b FrontendBase
	}{}
	for i := 0; i < 2; i++ {
		var sample FrontendBase
		var result FrontendBase
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
			a, b FrontendBase
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
			t.Errorf("Expected FrontendBase to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestFrontendBaseDiffFalse(t *testing.T) {
	samples := []struct {
		a, b FrontendBase
	}{}
	for i := 0; i < 2; i++ {
		var sample FrontendBase
		var result FrontendBase
		err := faker.FakeData(&sample, options.WithIgnoreInterface(true))
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result, options.WithIgnoreInterface(true))
		if err != nil {
			t.Errorf(err.Error())
		}
		result.Backlog = Ptr(*sample.Backlog + 1)
		result.Clflog = !sample.Clflog
		result.ClientFinTimeout = Ptr(*sample.ClientFinTimeout + 1)
		result.ClientTimeout = Ptr(*sample.ClientTimeout + 1)
		result.ClitcpkaCnt = Ptr(*sample.ClitcpkaCnt + 1)
		result.ClitcpkaIdle = Ptr(*sample.ClitcpkaIdle + 1)
		result.ClitcpkaIntvl = Ptr(*sample.ClitcpkaIntvl + 1)
		result.Disabled = !sample.Disabled
		result.Enabled = !sample.Enabled
		result.HTTPKeepAliveTimeout = Ptr(*sample.HTTPKeepAliveTimeout + 1)
		result.HTTPRequestTimeout = Ptr(*sample.HTTPRequestTimeout + 1)
		result.Httplog = !sample.Httplog
		result.ID = Ptr(*sample.ID + 1)
		result.Maxconn = Ptr(*sample.Maxconn + 1)
		result.TarpitTimeout = Ptr(*sample.TarpitTimeout + 1)
		result.Tcplog = !sample.Tcplog
		samples = append(samples, struct {
			a, b FrontendBase
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		if len(result) != 66 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected FrontendBase to be different in 66 cases, but it is not (%d) %s %s", len(result), a, b)
		}
	}
}
