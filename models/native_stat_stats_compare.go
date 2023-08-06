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

// Equal checks if two structs of type NativeStatStats are equal
//
//	var a, b NativeStatStats
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s NativeStatStats) Equal(t NativeStatStats, opts ...Options) bool {
	if !equalPointers(s.Act, t.Act) {
		return false
	}

	if s.Addr != t.Addr {
		return false
	}

	if !equalPointers(s.AgentCode, t.AgentCode) {
		return false
	}

	if s.AgentDesc != t.AgentDesc {
		return false
	}

	if !equalPointers(s.AgentDuration, t.AgentDuration) {
		return false
	}

	if !equalPointers(s.AgentFall, t.AgentFall) {
		return false
	}

	if !equalPointers(s.AgentHealth, t.AgentHealth) {
		return false
	}

	if !equalPointers(s.AgentRise, t.AgentRise) {
		return false
	}

	if s.AgentStatus != t.AgentStatus {
		return false
	}

	if s.Algo != t.Algo {
		return false
	}

	if !equalPointers(s.Bck, t.Bck) {
		return false
	}

	if !equalPointers(s.Bin, t.Bin) {
		return false
	}

	if !equalPointers(s.Bout, t.Bout) {
		return false
	}

	if !equalPointers(s.CheckCode, t.CheckCode) {
		return false
	}

	if s.CheckDesc != t.CheckDesc {
		return false
	}

	if !equalPointers(s.CheckDuration, t.CheckDuration) {
		return false
	}

	if !equalPointers(s.CheckFall, t.CheckFall) {
		return false
	}

	if !equalPointers(s.CheckHealth, t.CheckHealth) {
		return false
	}

	if !equalPointers(s.CheckRise, t.CheckRise) {
		return false
	}

	if s.CheckStatus != t.CheckStatus {
		return false
	}

	if !equalPointers(s.Chkdown, t.Chkdown) {
		return false
	}

	if !equalPointers(s.Chkfail, t.Chkfail) {
		return false
	}

	if !equalPointers(s.CliAbrt, t.CliAbrt) {
		return false
	}

	if !equalPointers(s.CompByp, t.CompByp) {
		return false
	}

	if !equalPointers(s.CompIn, t.CompIn) {
		return false
	}

	if !equalPointers(s.CompOut, t.CompOut) {
		return false
	}

	if !equalPointers(s.CompRsp, t.CompRsp) {
		return false
	}

	if !equalPointers(s.ConnRate, t.ConnRate) {
		return false
	}

	if !equalPointers(s.ConnRateMax, t.ConnRateMax) {
		return false
	}

	if !equalPointers(s.ConnTot, t.ConnTot) {
		return false
	}

	if s.Cookie != t.Cookie {
		return false
	}

	if !equalPointers(s.Ctime, t.Ctime) {
		return false
	}

	if !equalPointers(s.Dcon, t.Dcon) {
		return false
	}

	if !equalPointers(s.Downtime, t.Downtime) {
		return false
	}

	if !equalPointers(s.Dreq, t.Dreq) {
		return false
	}

	if !equalPointers(s.Dresp, t.Dresp) {
		return false
	}

	if !equalPointers(s.Dses, t.Dses) {
		return false
	}

	if !equalPointers(s.Econ, t.Econ) {
		return false
	}

	if !equalPointers(s.Ereq, t.Ereq) {
		return false
	}

	if !equalPointers(s.Eresp, t.Eresp) {
		return false
	}

	if s.Hanafail != t.Hanafail {
		return false
	}

	if !equalPointers(s.Hrsp1xx, t.Hrsp1xx) {
		return false
	}

	if !equalPointers(s.Hrsp2xx, t.Hrsp2xx) {
		return false
	}

	if !equalPointers(s.Hrsp3xx, t.Hrsp3xx) {
		return false
	}

	if !equalPointers(s.Hrsp4xx, t.Hrsp4xx) {
		return false
	}

	if !equalPointers(s.Hrsp5xx, t.Hrsp5xx) {
		return false
	}

	if !equalPointers(s.HrspOther, t.HrspOther) {
		return false
	}

	if !equalPointers(s.Iid, t.Iid) {
		return false
	}

	if !equalPointers(s.Intercepted, t.Intercepted) {
		return false
	}

	if !equalPointers(s.Lastchg, t.Lastchg) {
		return false
	}

	if !equalPointers(s.Lastsess, t.Lastsess) {
		return false
	}

	if !equalPointers(s.Lbtot, t.Lbtot) {
		return false
	}

	if s.Mode != t.Mode {
		return false
	}

	if !equalPointers(s.Pid, t.Pid) {
		return false
	}

	if !equalPointers(s.Qcur, t.Qcur) {
		return false
	}

	if !equalPointers(s.Qlimit, t.Qlimit) {
		return false
	}

	if !equalPointers(s.Qmax, t.Qmax) {
		return false
	}

	if !equalPointers(s.Qtime, t.Qtime) {
		return false
	}

	if !equalPointers(s.Rate, t.Rate) {
		return false
	}

	if !equalPointers(s.RateLim, t.RateLim) {
		return false
	}

	if !equalPointers(s.RateMax, t.RateMax) {
		return false
	}

	if !equalPointers(s.ReqRate, t.ReqRate) {
		return false
	}

	if !equalPointers(s.ReqRateMax, t.ReqRateMax) {
		return false
	}

	if !equalPointers(s.ReqTot, t.ReqTot) {
		return false
	}

	if !equalPointers(s.Rtime, t.Rtime) {
		return false
	}

	if !equalPointers(s.Scur, t.Scur) {
		return false
	}

	if !equalPointers(s.Sid, t.Sid) {
		return false
	}

	if !equalPointers(s.Slim, t.Slim) {
		return false
	}

	if !equalPointers(s.Smax, t.Smax) {
		return false
	}

	if !equalPointers(s.SrvAbrt, t.SrvAbrt) {
		return false
	}

	if s.Status != t.Status {
		return false
	}

	if !equalPointers(s.Stot, t.Stot) {
		return false
	}

	if !equalPointers(s.Throttle, t.Throttle) {
		return false
	}

	if s.Tracked != t.Tracked {
		return false
	}

	if !equalPointers(s.Ttime, t.Ttime) {
		return false
	}

	if !equalPointers(s.Weight, t.Weight) {
		return false
	}

	if !equalPointers(s.Wredis, t.Wredis) {
		return false
	}

	if !equalPointers(s.Wretr, t.Wretr) {
		return false
	}

	return true
}

// Diff checks if two structs of type NativeStatStats are equal
//
//	var a, b NativeStatStats
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s NativeStatStats) Diff(t NativeStatStats, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if !equalPointers(s.Act, t.Act) {
		diff["Act"] = []interface{}{s.Act, t.Act}
	}

	if s.Addr != t.Addr {
		diff["Addr"] = []interface{}{s.Addr, t.Addr}
	}

	if !equalPointers(s.AgentCode, t.AgentCode) {
		diff["AgentCode"] = []interface{}{s.AgentCode, t.AgentCode}
	}

	if s.AgentDesc != t.AgentDesc {
		diff["AgentDesc"] = []interface{}{s.AgentDesc, t.AgentDesc}
	}

	if !equalPointers(s.AgentDuration, t.AgentDuration) {
		diff["AgentDuration"] = []interface{}{s.AgentDuration, t.AgentDuration}
	}

	if !equalPointers(s.AgentFall, t.AgentFall) {
		diff["AgentFall"] = []interface{}{s.AgentFall, t.AgentFall}
	}

	if !equalPointers(s.AgentHealth, t.AgentHealth) {
		diff["AgentHealth"] = []interface{}{s.AgentHealth, t.AgentHealth}
	}

	if !equalPointers(s.AgentRise, t.AgentRise) {
		diff["AgentRise"] = []interface{}{s.AgentRise, t.AgentRise}
	}

	if s.AgentStatus != t.AgentStatus {
		diff["AgentStatus"] = []interface{}{s.AgentStatus, t.AgentStatus}
	}

	if s.Algo != t.Algo {
		diff["Algo"] = []interface{}{s.Algo, t.Algo}
	}

	if !equalPointers(s.Bck, t.Bck) {
		diff["Bck"] = []interface{}{s.Bck, t.Bck}
	}

	if !equalPointers(s.Bin, t.Bin) {
		diff["Bin"] = []interface{}{s.Bin, t.Bin}
	}

	if !equalPointers(s.Bout, t.Bout) {
		diff["Bout"] = []interface{}{s.Bout, t.Bout}
	}

	if !equalPointers(s.CheckCode, t.CheckCode) {
		diff["CheckCode"] = []interface{}{s.CheckCode, t.CheckCode}
	}

	if s.CheckDesc != t.CheckDesc {
		diff["CheckDesc"] = []interface{}{s.CheckDesc, t.CheckDesc}
	}

	if !equalPointers(s.CheckDuration, t.CheckDuration) {
		diff["CheckDuration"] = []interface{}{s.CheckDuration, t.CheckDuration}
	}

	if !equalPointers(s.CheckFall, t.CheckFall) {
		diff["CheckFall"] = []interface{}{s.CheckFall, t.CheckFall}
	}

	if !equalPointers(s.CheckHealth, t.CheckHealth) {
		diff["CheckHealth"] = []interface{}{s.CheckHealth, t.CheckHealth}
	}

	if !equalPointers(s.CheckRise, t.CheckRise) {
		diff["CheckRise"] = []interface{}{s.CheckRise, t.CheckRise}
	}

	if s.CheckStatus != t.CheckStatus {
		diff["CheckStatus"] = []interface{}{s.CheckStatus, t.CheckStatus}
	}

	if !equalPointers(s.Chkdown, t.Chkdown) {
		diff["Chkdown"] = []interface{}{s.Chkdown, t.Chkdown}
	}

	if !equalPointers(s.Chkfail, t.Chkfail) {
		diff["Chkfail"] = []interface{}{s.Chkfail, t.Chkfail}
	}

	if !equalPointers(s.CliAbrt, t.CliAbrt) {
		diff["CliAbrt"] = []interface{}{s.CliAbrt, t.CliAbrt}
	}

	if !equalPointers(s.CompByp, t.CompByp) {
		diff["CompByp"] = []interface{}{s.CompByp, t.CompByp}
	}

	if !equalPointers(s.CompIn, t.CompIn) {
		diff["CompIn"] = []interface{}{s.CompIn, t.CompIn}
	}

	if !equalPointers(s.CompOut, t.CompOut) {
		diff["CompOut"] = []interface{}{s.CompOut, t.CompOut}
	}

	if !equalPointers(s.CompRsp, t.CompRsp) {
		diff["CompRsp"] = []interface{}{s.CompRsp, t.CompRsp}
	}

	if !equalPointers(s.ConnRate, t.ConnRate) {
		diff["ConnRate"] = []interface{}{s.ConnRate, t.ConnRate}
	}

	if !equalPointers(s.ConnRateMax, t.ConnRateMax) {
		diff["ConnRateMax"] = []interface{}{s.ConnRateMax, t.ConnRateMax}
	}

	if !equalPointers(s.ConnTot, t.ConnTot) {
		diff["ConnTot"] = []interface{}{s.ConnTot, t.ConnTot}
	}

	if s.Cookie != t.Cookie {
		diff["Cookie"] = []interface{}{s.Cookie, t.Cookie}
	}

	if !equalPointers(s.Ctime, t.Ctime) {
		diff["Ctime"] = []interface{}{s.Ctime, t.Ctime}
	}

	if !equalPointers(s.Dcon, t.Dcon) {
		diff["Dcon"] = []interface{}{s.Dcon, t.Dcon}
	}

	if !equalPointers(s.Downtime, t.Downtime) {
		diff["Downtime"] = []interface{}{s.Downtime, t.Downtime}
	}

	if !equalPointers(s.Dreq, t.Dreq) {
		diff["Dreq"] = []interface{}{s.Dreq, t.Dreq}
	}

	if !equalPointers(s.Dresp, t.Dresp) {
		diff["Dresp"] = []interface{}{s.Dresp, t.Dresp}
	}

	if !equalPointers(s.Dses, t.Dses) {
		diff["Dses"] = []interface{}{s.Dses, t.Dses}
	}

	if !equalPointers(s.Econ, t.Econ) {
		diff["Econ"] = []interface{}{s.Econ, t.Econ}
	}

	if !equalPointers(s.Ereq, t.Ereq) {
		diff["Ereq"] = []interface{}{s.Ereq, t.Ereq}
	}

	if !equalPointers(s.Eresp, t.Eresp) {
		diff["Eresp"] = []interface{}{s.Eresp, t.Eresp}
	}

	if s.Hanafail != t.Hanafail {
		diff["Hanafail"] = []interface{}{s.Hanafail, t.Hanafail}
	}

	if !equalPointers(s.Hrsp1xx, t.Hrsp1xx) {
		diff["Hrsp1xx"] = []interface{}{s.Hrsp1xx, t.Hrsp1xx}
	}

	if !equalPointers(s.Hrsp2xx, t.Hrsp2xx) {
		diff["Hrsp2xx"] = []interface{}{s.Hrsp2xx, t.Hrsp2xx}
	}

	if !equalPointers(s.Hrsp3xx, t.Hrsp3xx) {
		diff["Hrsp3xx"] = []interface{}{s.Hrsp3xx, t.Hrsp3xx}
	}

	if !equalPointers(s.Hrsp4xx, t.Hrsp4xx) {
		diff["Hrsp4xx"] = []interface{}{s.Hrsp4xx, t.Hrsp4xx}
	}

	if !equalPointers(s.Hrsp5xx, t.Hrsp5xx) {
		diff["Hrsp5xx"] = []interface{}{s.Hrsp5xx, t.Hrsp5xx}
	}

	if !equalPointers(s.HrspOther, t.HrspOther) {
		diff["HrspOther"] = []interface{}{s.HrspOther, t.HrspOther}
	}

	if !equalPointers(s.Iid, t.Iid) {
		diff["Iid"] = []interface{}{s.Iid, t.Iid}
	}

	if !equalPointers(s.Intercepted, t.Intercepted) {
		diff["Intercepted"] = []interface{}{s.Intercepted, t.Intercepted}
	}

	if !equalPointers(s.Lastchg, t.Lastchg) {
		diff["Lastchg"] = []interface{}{s.Lastchg, t.Lastchg}
	}

	if !equalPointers(s.Lastsess, t.Lastsess) {
		diff["Lastsess"] = []interface{}{s.Lastsess, t.Lastsess}
	}

	if !equalPointers(s.Lbtot, t.Lbtot) {
		diff["Lbtot"] = []interface{}{s.Lbtot, t.Lbtot}
	}

	if s.Mode != t.Mode {
		diff["Mode"] = []interface{}{s.Mode, t.Mode}
	}

	if !equalPointers(s.Pid, t.Pid) {
		diff["Pid"] = []interface{}{s.Pid, t.Pid}
	}

	if !equalPointers(s.Qcur, t.Qcur) {
		diff["Qcur"] = []interface{}{s.Qcur, t.Qcur}
	}

	if !equalPointers(s.Qlimit, t.Qlimit) {
		diff["Qlimit"] = []interface{}{s.Qlimit, t.Qlimit}
	}

	if !equalPointers(s.Qmax, t.Qmax) {
		diff["Qmax"] = []interface{}{s.Qmax, t.Qmax}
	}

	if !equalPointers(s.Qtime, t.Qtime) {
		diff["Qtime"] = []interface{}{s.Qtime, t.Qtime}
	}

	if !equalPointers(s.Rate, t.Rate) {
		diff["Rate"] = []interface{}{s.Rate, t.Rate}
	}

	if !equalPointers(s.RateLim, t.RateLim) {
		diff["RateLim"] = []interface{}{s.RateLim, t.RateLim}
	}

	if !equalPointers(s.RateMax, t.RateMax) {
		diff["RateMax"] = []interface{}{s.RateMax, t.RateMax}
	}

	if !equalPointers(s.ReqRate, t.ReqRate) {
		diff["ReqRate"] = []interface{}{s.ReqRate, t.ReqRate}
	}

	if !equalPointers(s.ReqRateMax, t.ReqRateMax) {
		diff["ReqRateMax"] = []interface{}{s.ReqRateMax, t.ReqRateMax}
	}

	if !equalPointers(s.ReqTot, t.ReqTot) {
		diff["ReqTot"] = []interface{}{s.ReqTot, t.ReqTot}
	}

	if !equalPointers(s.Rtime, t.Rtime) {
		diff["Rtime"] = []interface{}{s.Rtime, t.Rtime}
	}

	if !equalPointers(s.Scur, t.Scur) {
		diff["Scur"] = []interface{}{s.Scur, t.Scur}
	}

	if !equalPointers(s.Sid, t.Sid) {
		diff["Sid"] = []interface{}{s.Sid, t.Sid}
	}

	if !equalPointers(s.Slim, t.Slim) {
		diff["Slim"] = []interface{}{s.Slim, t.Slim}
	}

	if !equalPointers(s.Smax, t.Smax) {
		diff["Smax"] = []interface{}{s.Smax, t.Smax}
	}

	if !equalPointers(s.SrvAbrt, t.SrvAbrt) {
		diff["SrvAbrt"] = []interface{}{s.SrvAbrt, t.SrvAbrt}
	}

	if s.Status != t.Status {
		diff["Status"] = []interface{}{s.Status, t.Status}
	}

	if !equalPointers(s.Stot, t.Stot) {
		diff["Stot"] = []interface{}{s.Stot, t.Stot}
	}

	if !equalPointers(s.Throttle, t.Throttle) {
		diff["Throttle"] = []interface{}{s.Throttle, t.Throttle}
	}

	if s.Tracked != t.Tracked {
		diff["Tracked"] = []interface{}{s.Tracked, t.Tracked}
	}

	if !equalPointers(s.Ttime, t.Ttime) {
		diff["Ttime"] = []interface{}{s.Ttime, t.Ttime}
	}

	if !equalPointers(s.Weight, t.Weight) {
		diff["Weight"] = []interface{}{s.Weight, t.Weight}
	}

	if !equalPointers(s.Wredis, t.Wredis) {
		diff["Wredis"] = []interface{}{s.Wredis, t.Wredis}
	}

	if !equalPointers(s.Wretr, t.Wretr) {
		diff["Wretr"] = []interface{}{s.Wretr, t.Wretr}
	}

	return diff
}
