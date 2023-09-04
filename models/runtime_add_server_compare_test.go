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

	jsoniter "github.com/json-iterator/go"
)

func TestRuntimeAddServerEqual(t *testing.T) {
	samples := []struct {
		a, b RuntimeAddServer
	}{}
	for i := 0; i < 2; i++ {
		var sample RuntimeAddServer
		var result RuntimeAddServer
		err := faker.FakeData(&sample)
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
			a, b RuntimeAddServer
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
			t.Errorf("Expected RuntimeAddServer to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestRuntimeAddServerEqualFalse(t *testing.T) {
	samples := []struct {
		a, b RuntimeAddServer
	}{}
	for i := 0; i < 2; i++ {
		var sample RuntimeAddServer
		var result RuntimeAddServer
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		result.AgentInter = Ptr(*sample.AgentInter + 1)
		result.AgentPort = Ptr(*sample.AgentPort + 1)
		result.Allow0rtt = !sample.Allow0rtt
		result.Downinter = Ptr(*sample.Downinter + 1)
		result.ErrorLimit = Ptr(*sample.ErrorLimit + 1)
		result.Fall = Ptr(*sample.Fall + 1)
		result.Fastinter = Ptr(*sample.Fastinter + 1)
		result.HealthCheckPort = Ptr(*sample.HealthCheckPort + 1)
		result.Inter = Ptr(*sample.Inter + 1)
		result.Maxconn = Ptr(*sample.Maxconn + 1)
		result.Maxqueue = Ptr(*sample.Maxqueue + 1)
		result.Minconn = Ptr(*sample.Minconn + 1)
		result.PoolLowConn = Ptr(*sample.PoolLowConn + 1)
		result.PoolMaxConn = Ptr(*sample.PoolMaxConn + 1)
		result.PoolPurgeDelay = Ptr(*sample.PoolPurgeDelay + 1)
		result.Port = Ptr(*sample.Port + 1)
		result.Rise = Ptr(*sample.Rise + 1)
		result.Slowstart = Ptr(*sample.Slowstart + 1)
		result.Weight = Ptr(*sample.Weight + 1)
		samples = append(samples, struct {
			a, b RuntimeAddServer
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
			t.Errorf("Expected RuntimeAddServer to be different, but it is not %s %s", a, b)
		}
	}
}

func TestRuntimeAddServerDiff(t *testing.T) {
	samples := []struct {
		a, b RuntimeAddServer
	}{}
	for i := 0; i < 2; i++ {
		var sample RuntimeAddServer
		var result RuntimeAddServer
		err := faker.FakeData(&sample)
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
			a, b RuntimeAddServer
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
			t.Errorf("Expected RuntimeAddServer to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestRuntimeAddServerDiffFalse(t *testing.T) {
	samples := []struct {
		a, b RuntimeAddServer
	}{}
	for i := 0; i < 2; i++ {
		var sample RuntimeAddServer
		var result RuntimeAddServer
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		result.AgentInter = Ptr(*sample.AgentInter + 1)
		result.AgentPort = Ptr(*sample.AgentPort + 1)
		result.Allow0rtt = !sample.Allow0rtt
		result.Downinter = Ptr(*sample.Downinter + 1)
		result.ErrorLimit = Ptr(*sample.ErrorLimit + 1)
		result.Fall = Ptr(*sample.Fall + 1)
		result.Fastinter = Ptr(*sample.Fastinter + 1)
		result.HealthCheckPort = Ptr(*sample.HealthCheckPort + 1)
		result.Inter = Ptr(*sample.Inter + 1)
		result.Maxconn = Ptr(*sample.Maxconn + 1)
		result.Maxqueue = Ptr(*sample.Maxqueue + 1)
		result.Minconn = Ptr(*sample.Minconn + 1)
		result.PoolLowConn = Ptr(*sample.PoolLowConn + 1)
		result.PoolMaxConn = Ptr(*sample.PoolMaxConn + 1)
		result.PoolPurgeDelay = Ptr(*sample.PoolPurgeDelay + 1)
		result.Port = Ptr(*sample.Port + 1)
		result.Rise = Ptr(*sample.Rise + 1)
		result.Slowstart = Ptr(*sample.Slowstart + 1)
		result.Weight = Ptr(*sample.Weight + 1)
		samples = append(samples, struct {
			a, b RuntimeAddServer
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		if len(result) != 74 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected RuntimeAddServer to be different in 74 cases, but it is not (%d) %s %s", len(result), a, b)
		}
	}
}
