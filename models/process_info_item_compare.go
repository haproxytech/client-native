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

// Equal checks if two structs of type ProcessInfoItem are equal
//
//	var a, b ProcessInfoItem
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s ProcessInfoItem) Equal(t ProcessInfoItem, opts ...Options) bool {
	if !equalPointers(s.ActivePeers, t.ActivePeers) {
		return false
	}

	if !equalPointers(s.BusyPolling, t.BusyPolling) {
		return false
	}

	if !equalPointers(s.BytesOutRate, t.BytesOutRate) {
		return false
	}

	if !equalPointers(s.CompressBpsIn, t.CompressBpsIn) {
		return false
	}

	if !equalPointers(s.CompressBpsOut, t.CompressBpsOut) {
		return false
	}

	if !equalPointers(s.CompressBpsRateLim, t.CompressBpsRateLim) {
		return false
	}

	if !equalPointers(s.ConnRate, t.ConnRate) {
		return false
	}

	if !equalPointers(s.ConnRateLimit, t.ConnRateLimit) {
		return false
	}

	if !equalPointers(s.ConnectedPeers, t.ConnectedPeers) {
		return false
	}

	if !equalPointers(s.CumConns, t.CumConns) {
		return false
	}

	if !equalPointers(s.CumReq, t.CumReq) {
		return false
	}

	if !equalPointers(s.CumSslConns, t.CumSslConns) {
		return false
	}

	if !equalPointers(s.CurrConns, t.CurrConns) {
		return false
	}

	if !equalPointers(s.CurrSslConns, t.CurrSslConns) {
		return false
	}

	if !equalPointers(s.DroppedLogs, t.DroppedLogs) {
		return false
	}

	if !equalPointers(s.FailedResolutions, t.FailedResolutions) {
		return false
	}

	if !equalPointers(s.HardMaxConn, t.HardMaxConn) {
		return false
	}

	if !equalPointers(s.IdlePct, t.IdlePct) {
		return false
	}

	if !equalPointers(s.Jobs, t.Jobs) {
		return false
	}

	if !equalPointers(s.Listeners, t.Listeners) {
		return false
	}

	if !equalPointers(s.MaxConn, t.MaxConn) {
		return false
	}

	if !equalPointers(s.MaxConnRate, t.MaxConnRate) {
		return false
	}

	if !equalPointers(s.MaxPipes, t.MaxPipes) {
		return false
	}

	if !equalPointers(s.MaxSessRate, t.MaxSessRate) {
		return false
	}

	if !equalPointers(s.MaxSock, t.MaxSock) {
		return false
	}

	if !equalPointers(s.MaxSslConns, t.MaxSslConns) {
		return false
	}

	if !equalPointers(s.MaxSslRate, t.MaxSslRate) {
		return false
	}

	if !equalPointers(s.MaxZlibMemUsage, t.MaxZlibMemUsage) {
		return false
	}

	if !equalPointers(s.MemMaxMb, t.MemMaxMb) {
		return false
	}

	if !equalPointers(s.Nbthread, t.Nbthread) {
		return false
	}

	if s.Node != t.Node {
		return false
	}

	if !equalPointers(s.Pid, t.Pid) {
		return false
	}

	if !equalPointers(s.PipesFree, t.PipesFree) {
		return false
	}

	if !equalPointers(s.PipesUsed, t.PipesUsed) {
		return false
	}

	if !equalPointers(s.PoolAllocMb, t.PoolAllocMb) {
		return false
	}

	if !equalPointers(s.PoolFailed, t.PoolFailed) {
		return false
	}

	if !equalPointers(s.PoolUsedMb, t.PoolUsedMb) {
		return false
	}

	if !equalPointers(s.ProcessNum, t.ProcessNum) {
		return false
	}

	if !equalPointers(s.Processes, t.Processes) {
		return false
	}

	if !s.ReleaseDate.Equal(t.ReleaseDate) {
		return false
	}

	if !equalPointers(s.RunQueue, t.RunQueue) {
		return false
	}

	if !equalPointers(s.SessRate, t.SessRate) {
		return false
	}

	if !equalPointers(s.SessRateLimit, t.SessRateLimit) {
		return false
	}

	if !equalPointers(s.SslBackendKeyRate, t.SslBackendKeyRate) {
		return false
	}

	if !equalPointers(s.SslBackendMaxKeyRate, t.SslBackendMaxKeyRate) {
		return false
	}

	if !equalPointers(s.SslCacheLookups, t.SslCacheLookups) {
		return false
	}

	if !equalPointers(s.SslCacheMisses, t.SslCacheMisses) {
		return false
	}

	if !equalPointers(s.SslFrontendKeyRate, t.SslFrontendKeyRate) {
		return false
	}

	if !equalPointers(s.SslFrontendMaxKeyRate, t.SslFrontendMaxKeyRate) {
		return false
	}

	if !equalPointers(s.SslFrontendSessionReuse, t.SslFrontendSessionReuse) {
		return false
	}

	if !equalPointers(s.SslRate, t.SslRate) {
		return false
	}

	if !equalPointers(s.SslRateLimit, t.SslRateLimit) {
		return false
	}

	if !equalPointers(s.Stopping, t.Stopping) {
		return false
	}

	if !equalPointers(s.Tasks, t.Tasks) {
		return false
	}

	if !equalPointers(s.TotalBytesOut, t.TotalBytesOut) {
		return false
	}

	if !equalPointers(s.Ulimitn, t.Ulimitn) {
		return false
	}

	if !equalPointers(s.Unstoppable, t.Unstoppable) {
		return false
	}

	if !equalPointers(s.Uptime, t.Uptime) {
		return false
	}

	if s.Version != t.Version {
		return false
	}

	if !equalPointers(s.ZlibMemUsage, t.ZlibMemUsage) {
		return false
	}

	return true
}

// Diff checks if two structs of type ProcessInfoItem are equal
//
//	var a, b ProcessInfoItem
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s ProcessInfoItem) Diff(t ProcessInfoItem, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if !equalPointers(s.ActivePeers, t.ActivePeers) {
		diff["ActivePeers"] = []interface{}{ValueOrNil(s.ActivePeers), ValueOrNil(t.ActivePeers)}
	}

	if !equalPointers(s.BusyPolling, t.BusyPolling) {
		diff["BusyPolling"] = []interface{}{ValueOrNil(s.BusyPolling), ValueOrNil(t.BusyPolling)}
	}

	if !equalPointers(s.BytesOutRate, t.BytesOutRate) {
		diff["BytesOutRate"] = []interface{}{ValueOrNil(s.BytesOutRate), ValueOrNil(t.BytesOutRate)}
	}

	if !equalPointers(s.CompressBpsIn, t.CompressBpsIn) {
		diff["CompressBpsIn"] = []interface{}{ValueOrNil(s.CompressBpsIn), ValueOrNil(t.CompressBpsIn)}
	}

	if !equalPointers(s.CompressBpsOut, t.CompressBpsOut) {
		diff["CompressBpsOut"] = []interface{}{ValueOrNil(s.CompressBpsOut), ValueOrNil(t.CompressBpsOut)}
	}

	if !equalPointers(s.CompressBpsRateLim, t.CompressBpsRateLim) {
		diff["CompressBpsRateLim"] = []interface{}{ValueOrNil(s.CompressBpsRateLim), ValueOrNil(t.CompressBpsRateLim)}
	}

	if !equalPointers(s.ConnRate, t.ConnRate) {
		diff["ConnRate"] = []interface{}{ValueOrNil(s.ConnRate), ValueOrNil(t.ConnRate)}
	}

	if !equalPointers(s.ConnRateLimit, t.ConnRateLimit) {
		diff["ConnRateLimit"] = []interface{}{ValueOrNil(s.ConnRateLimit), ValueOrNil(t.ConnRateLimit)}
	}

	if !equalPointers(s.ConnectedPeers, t.ConnectedPeers) {
		diff["ConnectedPeers"] = []interface{}{ValueOrNil(s.ConnectedPeers), ValueOrNil(t.ConnectedPeers)}
	}

	if !equalPointers(s.CumConns, t.CumConns) {
		diff["CumConns"] = []interface{}{ValueOrNil(s.CumConns), ValueOrNil(t.CumConns)}
	}

	if !equalPointers(s.CumReq, t.CumReq) {
		diff["CumReq"] = []interface{}{ValueOrNil(s.CumReq), ValueOrNil(t.CumReq)}
	}

	if !equalPointers(s.CumSslConns, t.CumSslConns) {
		diff["CumSslConns"] = []interface{}{ValueOrNil(s.CumSslConns), ValueOrNil(t.CumSslConns)}
	}

	if !equalPointers(s.CurrConns, t.CurrConns) {
		diff["CurrConns"] = []interface{}{ValueOrNil(s.CurrConns), ValueOrNil(t.CurrConns)}
	}

	if !equalPointers(s.CurrSslConns, t.CurrSslConns) {
		diff["CurrSslConns"] = []interface{}{ValueOrNil(s.CurrSslConns), ValueOrNil(t.CurrSslConns)}
	}

	if !equalPointers(s.DroppedLogs, t.DroppedLogs) {
		diff["DroppedLogs"] = []interface{}{ValueOrNil(s.DroppedLogs), ValueOrNil(t.DroppedLogs)}
	}

	if !equalPointers(s.FailedResolutions, t.FailedResolutions) {
		diff["FailedResolutions"] = []interface{}{ValueOrNil(s.FailedResolutions), ValueOrNil(t.FailedResolutions)}
	}

	if !equalPointers(s.HardMaxConn, t.HardMaxConn) {
		diff["HardMaxConn"] = []interface{}{ValueOrNil(s.HardMaxConn), ValueOrNil(t.HardMaxConn)}
	}

	if !equalPointers(s.IdlePct, t.IdlePct) {
		diff["IdlePct"] = []interface{}{ValueOrNil(s.IdlePct), ValueOrNil(t.IdlePct)}
	}

	if !equalPointers(s.Jobs, t.Jobs) {
		diff["Jobs"] = []interface{}{ValueOrNil(s.Jobs), ValueOrNil(t.Jobs)}
	}

	if !equalPointers(s.Listeners, t.Listeners) {
		diff["Listeners"] = []interface{}{ValueOrNil(s.Listeners), ValueOrNil(t.Listeners)}
	}

	if !equalPointers(s.MaxConn, t.MaxConn) {
		diff["MaxConn"] = []interface{}{ValueOrNil(s.MaxConn), ValueOrNil(t.MaxConn)}
	}

	if !equalPointers(s.MaxConnRate, t.MaxConnRate) {
		diff["MaxConnRate"] = []interface{}{ValueOrNil(s.MaxConnRate), ValueOrNil(t.MaxConnRate)}
	}

	if !equalPointers(s.MaxPipes, t.MaxPipes) {
		diff["MaxPipes"] = []interface{}{ValueOrNil(s.MaxPipes), ValueOrNil(t.MaxPipes)}
	}

	if !equalPointers(s.MaxSessRate, t.MaxSessRate) {
		diff["MaxSessRate"] = []interface{}{ValueOrNil(s.MaxSessRate), ValueOrNil(t.MaxSessRate)}
	}

	if !equalPointers(s.MaxSock, t.MaxSock) {
		diff["MaxSock"] = []interface{}{ValueOrNil(s.MaxSock), ValueOrNil(t.MaxSock)}
	}

	if !equalPointers(s.MaxSslConns, t.MaxSslConns) {
		diff["MaxSslConns"] = []interface{}{ValueOrNil(s.MaxSslConns), ValueOrNil(t.MaxSslConns)}
	}

	if !equalPointers(s.MaxSslRate, t.MaxSslRate) {
		diff["MaxSslRate"] = []interface{}{ValueOrNil(s.MaxSslRate), ValueOrNil(t.MaxSslRate)}
	}

	if !equalPointers(s.MaxZlibMemUsage, t.MaxZlibMemUsage) {
		diff["MaxZlibMemUsage"] = []interface{}{ValueOrNil(s.MaxZlibMemUsage), ValueOrNil(t.MaxZlibMemUsage)}
	}

	if !equalPointers(s.MemMaxMb, t.MemMaxMb) {
		diff["MemMaxMb"] = []interface{}{ValueOrNil(s.MemMaxMb), ValueOrNil(t.MemMaxMb)}
	}

	if !equalPointers(s.Nbthread, t.Nbthread) {
		diff["Nbthread"] = []interface{}{ValueOrNil(s.Nbthread), ValueOrNil(t.Nbthread)}
	}

	if s.Node != t.Node {
		diff["Node"] = []interface{}{s.Node, t.Node}
	}

	if !equalPointers(s.Pid, t.Pid) {
		diff["Pid"] = []interface{}{ValueOrNil(s.Pid), ValueOrNil(t.Pid)}
	}

	if !equalPointers(s.PipesFree, t.PipesFree) {
		diff["PipesFree"] = []interface{}{ValueOrNil(s.PipesFree), ValueOrNil(t.PipesFree)}
	}

	if !equalPointers(s.PipesUsed, t.PipesUsed) {
		diff["PipesUsed"] = []interface{}{ValueOrNil(s.PipesUsed), ValueOrNil(t.PipesUsed)}
	}

	if !equalPointers(s.PoolAllocMb, t.PoolAllocMb) {
		diff["PoolAllocMb"] = []interface{}{ValueOrNil(s.PoolAllocMb), ValueOrNil(t.PoolAllocMb)}
	}

	if !equalPointers(s.PoolFailed, t.PoolFailed) {
		diff["PoolFailed"] = []interface{}{ValueOrNil(s.PoolFailed), ValueOrNil(t.PoolFailed)}
	}

	if !equalPointers(s.PoolUsedMb, t.PoolUsedMb) {
		diff["PoolUsedMb"] = []interface{}{ValueOrNil(s.PoolUsedMb), ValueOrNil(t.PoolUsedMb)}
	}

	if !equalPointers(s.ProcessNum, t.ProcessNum) {
		diff["ProcessNum"] = []interface{}{ValueOrNil(s.ProcessNum), ValueOrNil(t.ProcessNum)}
	}

	if !equalPointers(s.Processes, t.Processes) {
		diff["Processes"] = []interface{}{ValueOrNil(s.Processes), ValueOrNil(t.Processes)}
	}

	if !s.ReleaseDate.Equal(t.ReleaseDate) {
		diff["ReleaseDate"] = []interface{}{s.ReleaseDate, t.ReleaseDate}
	}

	if !equalPointers(s.RunQueue, t.RunQueue) {
		diff["RunQueue"] = []interface{}{ValueOrNil(s.RunQueue), ValueOrNil(t.RunQueue)}
	}

	if !equalPointers(s.SessRate, t.SessRate) {
		diff["SessRate"] = []interface{}{ValueOrNil(s.SessRate), ValueOrNil(t.SessRate)}
	}

	if !equalPointers(s.SessRateLimit, t.SessRateLimit) {
		diff["SessRateLimit"] = []interface{}{ValueOrNil(s.SessRateLimit), ValueOrNil(t.SessRateLimit)}
	}

	if !equalPointers(s.SslBackendKeyRate, t.SslBackendKeyRate) {
		diff["SslBackendKeyRate"] = []interface{}{ValueOrNil(s.SslBackendKeyRate), ValueOrNil(t.SslBackendKeyRate)}
	}

	if !equalPointers(s.SslBackendMaxKeyRate, t.SslBackendMaxKeyRate) {
		diff["SslBackendMaxKeyRate"] = []interface{}{ValueOrNil(s.SslBackendMaxKeyRate), ValueOrNil(t.SslBackendMaxKeyRate)}
	}

	if !equalPointers(s.SslCacheLookups, t.SslCacheLookups) {
		diff["SslCacheLookups"] = []interface{}{ValueOrNil(s.SslCacheLookups), ValueOrNil(t.SslCacheLookups)}
	}

	if !equalPointers(s.SslCacheMisses, t.SslCacheMisses) {
		diff["SslCacheMisses"] = []interface{}{ValueOrNil(s.SslCacheMisses), ValueOrNil(t.SslCacheMisses)}
	}

	if !equalPointers(s.SslFrontendKeyRate, t.SslFrontendKeyRate) {
		diff["SslFrontendKeyRate"] = []interface{}{ValueOrNil(s.SslFrontendKeyRate), ValueOrNil(t.SslFrontendKeyRate)}
	}

	if !equalPointers(s.SslFrontendMaxKeyRate, t.SslFrontendMaxKeyRate) {
		diff["SslFrontendMaxKeyRate"] = []interface{}{ValueOrNil(s.SslFrontendMaxKeyRate), ValueOrNil(t.SslFrontendMaxKeyRate)}
	}

	if !equalPointers(s.SslFrontendSessionReuse, t.SslFrontendSessionReuse) {
		diff["SslFrontendSessionReuse"] = []interface{}{ValueOrNil(s.SslFrontendSessionReuse), ValueOrNil(t.SslFrontendSessionReuse)}
	}

	if !equalPointers(s.SslRate, t.SslRate) {
		diff["SslRate"] = []interface{}{ValueOrNil(s.SslRate), ValueOrNil(t.SslRate)}
	}

	if !equalPointers(s.SslRateLimit, t.SslRateLimit) {
		diff["SslRateLimit"] = []interface{}{ValueOrNil(s.SslRateLimit), ValueOrNil(t.SslRateLimit)}
	}

	if !equalPointers(s.Stopping, t.Stopping) {
		diff["Stopping"] = []interface{}{ValueOrNil(s.Stopping), ValueOrNil(t.Stopping)}
	}

	if !equalPointers(s.Tasks, t.Tasks) {
		diff["Tasks"] = []interface{}{ValueOrNil(s.Tasks), ValueOrNil(t.Tasks)}
	}

	if !equalPointers(s.TotalBytesOut, t.TotalBytesOut) {
		diff["TotalBytesOut"] = []interface{}{ValueOrNil(s.TotalBytesOut), ValueOrNil(t.TotalBytesOut)}
	}

	if !equalPointers(s.Ulimitn, t.Ulimitn) {
		diff["Ulimitn"] = []interface{}{ValueOrNil(s.Ulimitn), ValueOrNil(t.Ulimitn)}
	}

	if !equalPointers(s.Unstoppable, t.Unstoppable) {
		diff["Unstoppable"] = []interface{}{ValueOrNil(s.Unstoppable), ValueOrNil(t.Unstoppable)}
	}

	if !equalPointers(s.Uptime, t.Uptime) {
		diff["Uptime"] = []interface{}{ValueOrNil(s.Uptime), ValueOrNil(t.Uptime)}
	}

	if s.Version != t.Version {
		diff["Version"] = []interface{}{s.Version, t.Version}
	}

	if !equalPointers(s.ZlibMemUsage, t.ZlibMemUsage) {
		diff["ZlibMemUsage"] = []interface{}{ValueOrNil(s.ZlibMemUsage), ValueOrNil(t.ZlibMemUsage)}
	}

	return diff
}
