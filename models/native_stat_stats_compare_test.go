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

func TestNativeStatStatsEqual(t *testing.T) {
	samples := []struct {
		a, b NativeStatStats
	}{}
	for i := 0; i < 2; i++ {
		var sample NativeStatStats
		var result NativeStatStats
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
			a, b NativeStatStats
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
			t.Errorf("Expected NativeStatStats to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestNativeStatStatsEqualFalse(t *testing.T) {
	samples := []struct {
		a, b NativeStatStats
	}{}
	for i := 0; i < 2; i++ {
		var sample NativeStatStats
		var result NativeStatStats
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		result.Act = Ptr(*sample.Act + 1)
		result.AgentCode = Ptr(*sample.AgentCode + 1)
		result.AgentDuration = Ptr(*sample.AgentDuration + 1)
		result.AgentFall = Ptr(*sample.AgentFall + 1)
		result.AgentHealth = Ptr(*sample.AgentHealth + 1)
		result.AgentRise = Ptr(*sample.AgentRise + 1)
		result.Bck = Ptr(*sample.Bck + 1)
		result.Bin = Ptr(*sample.Bin + 1)
		result.Bout = Ptr(*sample.Bout + 1)
		result.CheckCode = Ptr(*sample.CheckCode + 1)
		result.CheckDuration = Ptr(*sample.CheckDuration + 1)
		result.CheckFall = Ptr(*sample.CheckFall + 1)
		result.CheckHealth = Ptr(*sample.CheckHealth + 1)
		result.CheckRise = Ptr(*sample.CheckRise + 1)
		result.Chkdown = Ptr(*sample.Chkdown + 1)
		result.Chkfail = Ptr(*sample.Chkfail + 1)
		result.CliAbrt = Ptr(*sample.CliAbrt + 1)
		result.CompByp = Ptr(*sample.CompByp + 1)
		result.CompIn = Ptr(*sample.CompIn + 1)
		result.CompOut = Ptr(*sample.CompOut + 1)
		result.CompRsp = Ptr(*sample.CompRsp + 1)
		result.ConnRate = Ptr(*sample.ConnRate + 1)
		result.ConnRateMax = Ptr(*sample.ConnRateMax + 1)
		result.ConnTot = Ptr(*sample.ConnTot + 1)
		result.Ctime = Ptr(*sample.Ctime + 1)
		result.Dcon = Ptr(*sample.Dcon + 1)
		result.Downtime = Ptr(*sample.Downtime + 1)
		result.Dreq = Ptr(*sample.Dreq + 1)
		result.Dresp = Ptr(*sample.Dresp + 1)
		result.Dses = Ptr(*sample.Dses + 1)
		result.Econ = Ptr(*sample.Econ + 1)
		result.Ereq = Ptr(*sample.Ereq + 1)
		result.Eresp = Ptr(*sample.Eresp + 1)
		result.Hrsp1xx = Ptr(*sample.Hrsp1xx + 1)
		result.Hrsp2xx = Ptr(*sample.Hrsp2xx + 1)
		result.Hrsp3xx = Ptr(*sample.Hrsp3xx + 1)
		result.Hrsp4xx = Ptr(*sample.Hrsp4xx + 1)
		result.Hrsp5xx = Ptr(*sample.Hrsp5xx + 1)
		result.HrspOther = Ptr(*sample.HrspOther + 1)
		result.Iid = Ptr(*sample.Iid + 1)
		result.Intercepted = Ptr(*sample.Intercepted + 1)
		result.Lastchg = Ptr(*sample.Lastchg + 1)
		result.Lastsess = Ptr(*sample.Lastsess + 1)
		result.Lbtot = Ptr(*sample.Lbtot + 1)
		result.Pid = Ptr(*sample.Pid + 1)
		result.Qcur = Ptr(*sample.Qcur + 1)
		result.Qlimit = Ptr(*sample.Qlimit + 1)
		result.Qmax = Ptr(*sample.Qmax + 1)
		result.Qtime = Ptr(*sample.Qtime + 1)
		result.Rate = Ptr(*sample.Rate + 1)
		result.RateLim = Ptr(*sample.RateLim + 1)
		result.RateMax = Ptr(*sample.RateMax + 1)
		result.ReqRate = Ptr(*sample.ReqRate + 1)
		result.ReqRateMax = Ptr(*sample.ReqRateMax + 1)
		result.ReqTot = Ptr(*sample.ReqTot + 1)
		result.Rtime = Ptr(*sample.Rtime + 1)
		result.Scur = Ptr(*sample.Scur + 1)
		result.Sid = Ptr(*sample.Sid + 1)
		result.Slim = Ptr(*sample.Slim + 1)
		result.Smax = Ptr(*sample.Smax + 1)
		result.SrvAbrt = Ptr(*sample.SrvAbrt + 1)
		result.Stot = Ptr(*sample.Stot + 1)
		result.Throttle = Ptr(*sample.Throttle + 1)
		result.Ttime = Ptr(*sample.Ttime + 1)
		result.Weight = Ptr(*sample.Weight + 1)
		result.Wredis = Ptr(*sample.Wredis + 1)
		result.Wretr = Ptr(*sample.Wretr + 1)
		samples = append(samples, struct {
			a, b NativeStatStats
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
			t.Errorf("Expected NativeStatStats to be different, but it is not %s %s", a, b)
		}
	}
}

func TestNativeStatStatsDiff(t *testing.T) {
	samples := []struct {
		a, b NativeStatStats
	}{}
	for i := 0; i < 2; i++ {
		var sample NativeStatStats
		var result NativeStatStats
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
			a, b NativeStatStats
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
			t.Errorf("Expected NativeStatStats to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestNativeStatStatsDiffFalse(t *testing.T) {
	samples := []struct {
		a, b NativeStatStats
	}{}
	for i := 0; i < 2; i++ {
		var sample NativeStatStats
		var result NativeStatStats
		err := faker.FakeData(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = faker.FakeData(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		result.Act = Ptr(*sample.Act + 1)
		result.AgentCode = Ptr(*sample.AgentCode + 1)
		result.AgentDuration = Ptr(*sample.AgentDuration + 1)
		result.AgentFall = Ptr(*sample.AgentFall + 1)
		result.AgentHealth = Ptr(*sample.AgentHealth + 1)
		result.AgentRise = Ptr(*sample.AgentRise + 1)
		result.Bck = Ptr(*sample.Bck + 1)
		result.Bin = Ptr(*sample.Bin + 1)
		result.Bout = Ptr(*sample.Bout + 1)
		result.CheckCode = Ptr(*sample.CheckCode + 1)
		result.CheckDuration = Ptr(*sample.CheckDuration + 1)
		result.CheckFall = Ptr(*sample.CheckFall + 1)
		result.CheckHealth = Ptr(*sample.CheckHealth + 1)
		result.CheckRise = Ptr(*sample.CheckRise + 1)
		result.Chkdown = Ptr(*sample.Chkdown + 1)
		result.Chkfail = Ptr(*sample.Chkfail + 1)
		result.CliAbrt = Ptr(*sample.CliAbrt + 1)
		result.CompByp = Ptr(*sample.CompByp + 1)
		result.CompIn = Ptr(*sample.CompIn + 1)
		result.CompOut = Ptr(*sample.CompOut + 1)
		result.CompRsp = Ptr(*sample.CompRsp + 1)
		result.ConnRate = Ptr(*sample.ConnRate + 1)
		result.ConnRateMax = Ptr(*sample.ConnRateMax + 1)
		result.ConnTot = Ptr(*sample.ConnTot + 1)
		result.Ctime = Ptr(*sample.Ctime + 1)
		result.Dcon = Ptr(*sample.Dcon + 1)
		result.Downtime = Ptr(*sample.Downtime + 1)
		result.Dreq = Ptr(*sample.Dreq + 1)
		result.Dresp = Ptr(*sample.Dresp + 1)
		result.Dses = Ptr(*sample.Dses + 1)
		result.Econ = Ptr(*sample.Econ + 1)
		result.Ereq = Ptr(*sample.Ereq + 1)
		result.Eresp = Ptr(*sample.Eresp + 1)
		result.Hrsp1xx = Ptr(*sample.Hrsp1xx + 1)
		result.Hrsp2xx = Ptr(*sample.Hrsp2xx + 1)
		result.Hrsp3xx = Ptr(*sample.Hrsp3xx + 1)
		result.Hrsp4xx = Ptr(*sample.Hrsp4xx + 1)
		result.Hrsp5xx = Ptr(*sample.Hrsp5xx + 1)
		result.HrspOther = Ptr(*sample.HrspOther + 1)
		result.Iid = Ptr(*sample.Iid + 1)
		result.Intercepted = Ptr(*sample.Intercepted + 1)
		result.Lastchg = Ptr(*sample.Lastchg + 1)
		result.Lastsess = Ptr(*sample.Lastsess + 1)
		result.Lbtot = Ptr(*sample.Lbtot + 1)
		result.Pid = Ptr(*sample.Pid + 1)
		result.Qcur = Ptr(*sample.Qcur + 1)
		result.Qlimit = Ptr(*sample.Qlimit + 1)
		result.Qmax = Ptr(*sample.Qmax + 1)
		result.Qtime = Ptr(*sample.Qtime + 1)
		result.Rate = Ptr(*sample.Rate + 1)
		result.RateLim = Ptr(*sample.RateLim + 1)
		result.RateMax = Ptr(*sample.RateMax + 1)
		result.ReqRate = Ptr(*sample.ReqRate + 1)
		result.ReqRateMax = Ptr(*sample.ReqRateMax + 1)
		result.ReqTot = Ptr(*sample.ReqTot + 1)
		result.Rtime = Ptr(*sample.Rtime + 1)
		result.Scur = Ptr(*sample.Scur + 1)
		result.Sid = Ptr(*sample.Sid + 1)
		result.Slim = Ptr(*sample.Slim + 1)
		result.Smax = Ptr(*sample.Smax + 1)
		result.SrvAbrt = Ptr(*sample.SrvAbrt + 1)
		result.Stot = Ptr(*sample.Stot + 1)
		result.Throttle = Ptr(*sample.Throttle + 1)
		result.Ttime = Ptr(*sample.Ttime + 1)
		result.Weight = Ptr(*sample.Weight + 1)
		result.Wredis = Ptr(*sample.Wredis + 1)
		result.Wretr = Ptr(*sample.Wretr + 1)
		samples = append(samples, struct {
			a, b NativeStatStats
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		if len(result) != 78 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected NativeStatStats to be different in 78 cases, but it is not (%d) %s %s", len(result), a, b)
		}
	}
}
