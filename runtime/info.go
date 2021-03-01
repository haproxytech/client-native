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

package runtime

import (
	"strconv"
	"strings"

	"github.com/go-openapi/strfmt"

	"github.com/haproxytech/client-native/v2/models"
)

// GetInfo fetches HAProxy info from runtime API
func (s *SingleRuntime) GetInfo() models.ProcessInfo {
	dataStr, err := s.ExecuteRaw("show info typed")
	data := models.ProcessInfo{RuntimeAPI: s.socketPath}
	if err != nil {
		data.Error = err.Error()
		return data
	}

	data.Info = parseInfo(dataStr)
	return data
}

func parseInfo(info string) *models.ProcessInfoItem { //nolint:gocognit,gocyclo
	data := &models.ProcessInfoItem{}

	for _, line := range strings.Split(info, "\n") {
		fields := strings.Split(line, ":")
		fID := strings.TrimSpace(strings.Split(fields[0], ".")[0])
		switch fID {
		case "1":
			data.Version = fields[3]
		case "2":
			d := strfmt.Date{}
			err := d.Scan(strings.ReplaceAll(fields[3], "/", "-"))
			if err == nil {
				data.ReleaseDate = d
			}
		case "3":
			nbthread, err := strconv.ParseInt(fields[3], 10, 64)
			if err == nil {
				data.Nbthread = &nbthread
			}
		case "4":
			nbproc, err := strconv.ParseInt(fields[3], 10, 64)
			if err == nil {
				data.Processes = &nbproc
			}
		case "5":
			procNum, err := strconv.ParseInt(fields[3], 10, 64)
			if err == nil {
				data.ProcessNum = &procNum
			}
		case "6":
			pid, err := strconv.ParseInt(fields[3], 10, 64)
			if err == nil {
				data.Pid = &pid
			}
		case "8":
			uptime, err := strconv.ParseInt(fields[3], 10, 64)
			if err == nil {
				data.Uptime = &uptime
			}
		case "9":
			mmMB, err := strconv.ParseInt(fields[3], 10, 64)
			if err == nil {
				data.MemMaxMb = &mmMB
			}
		case "10":
			pAllocMB, err := strconv.ParseInt(fields[3], 10, 64)
			if err == nil {
				data.PoolAllocMb = &pAllocMB
			}
		case "11":
			pUsedMB, err := strconv.ParseInt(fields[3], 10, 64)
			if err == nil {
				data.PoolUsedMb = &pUsedMB
			}
		case "12":
			pFailed, err := strconv.ParseInt(fields[3], 10, 64)
			if err == nil {
				data.PoolFailed = &pFailed
			}
		case "13":
			uLimitN, err := strconv.ParseInt(fields[3], 10, 64)
			if err == nil {
				data.Ulimitn = &uLimitN
			}
		case "14":
			maxSock, err := strconv.ParseInt(fields[3], 10, 64)
			if err == nil {
				data.MaxSock = &maxSock
			}
		case "15":
			maxConn, err := strconv.ParseInt(fields[3], 10, 64)
			if err == nil {
				data.MaxConn = &maxConn
			}
		case "16":
			hMaxConn, err := strconv.ParseInt(fields[3], 10, 64)
			if err == nil {
				data.HardMaxConn = &hMaxConn
			}
		case "17":
			currConn, err := strconv.ParseInt(fields[3], 10, 64)
			if err == nil {
				data.CurrConns = &currConn
			}
		case "18":
			cumConn, err := strconv.ParseInt(fields[3], 10, 64)
			if err == nil {
				data.CumConns = &cumConn
			}
		case "19":
			cumReq, err := strconv.ParseInt(fields[3], 10, 64)
			if err == nil {
				data.CumReq = &cumReq
			}
		case "20":
			maxSSLConn, err := strconv.ParseInt(fields[3], 10, 64)
			if err == nil {
				data.MaxSslConns = &maxSSLConn
			}
		case "21":
			curSSLConn, err := strconv.ParseInt(fields[3], 10, 64)
			if err == nil {
				data.CurrSslConns = &curSSLConn
			}
		case "22":
			cumSSLCons, err := strconv.ParseInt(fields[3], 10, 64)
			if err == nil {
				data.CumSslConns = &cumSSLCons
			}
		case "23":
			maxPipes, err := strconv.ParseInt(fields[3], 10, 64)
			if err == nil {
				data.MaxPipes = &maxPipes
			}
		case "24":
			pipesUsed, err := strconv.ParseInt(fields[3], 10, 64)
			if err == nil {
				data.PipesUsed = &pipesUsed
			}
		case "25":
			pipesFree, err := strconv.ParseInt(fields[3], 10, 64)
			if err == nil {
				data.PipesFree = &pipesFree
			}
		case "26":
			connRate, err := strconv.ParseInt(fields[3], 10, 64)
			if err == nil {
				data.ConnRate = &connRate
			}
		case "27":
			connRateLimit, err := strconv.ParseInt(fields[3], 10, 64)
			if err == nil {
				data.ConnRateLimit = &connRateLimit
			}
		case "28":
			maxConnRate, err := strconv.ParseInt(fields[3], 10, 64)
			if err == nil {
				data.MaxConnRate = &maxConnRate
			}
		case "29":
			sessRate, err := strconv.ParseInt(fields[3], 10, 64)
			if err == nil {
				data.SessRate = &sessRate
			}
		case "30":
			sessRateLimit, err := strconv.ParseInt(fields[3], 10, 64)
			if err == nil {
				data.SessRateLimit = &sessRateLimit
			}
		case "31":
			maxSessRate, err := strconv.ParseInt(fields[3], 10, 64)
			if err == nil {
				data.MaxSessRate = &maxSessRate
			}
		case "32":
			sslRate, err := strconv.ParseInt(fields[3], 10, 64)
			if err == nil {
				data.SslRate = &sslRate
			}
		case "33":
			sslRateLimit, err := strconv.ParseInt(fields[3], 10, 64)
			if err == nil {
				data.SslRateLimit = &sslRateLimit
			}
		case "34":
			maxSSLRate, err := strconv.ParseInt(fields[3], 10, 64)
			if err == nil {
				data.MaxSslRate = &maxSSLRate
			}
		case "35":
			sslFrKeyRate, err := strconv.ParseInt(fields[3], 10, 64)
			if err == nil {
				data.SslFrontendKeyRate = &sslFrKeyRate
			}
		case "36":
			sslFrMaxKeyRate, err := strconv.ParseInt(fields[3], 10, 64)
			if err == nil {
				data.SslFrontendMaxKeyRate = &sslFrMaxKeyRate
			}
		case "37":
			sslFrSessionReusePct, err := strconv.ParseInt(fields[3], 10, 64)
			if err == nil {
				data.SslFrontendSessionReuse = &sslFrSessionReusePct
			}
		case "38":
			sslBckKeyRate, err := strconv.ParseInt(fields[3], 10, 64)
			if err == nil {
				data.SslBackendKeyRate = &sslBckKeyRate
			}
		case "39":
			sslBckMaxKeyRate, err := strconv.ParseInt(fields[3], 10, 64)
			if err == nil {
				data.SslBackendMaxKeyRate = &sslBckMaxKeyRate
			}
		case "40":
			sslCacheLookups, err := strconv.ParseInt(fields[3], 10, 64)
			if err == nil {
				data.SslCacheLookups = &sslCacheLookups
			}
		case "41":
			sslCacheMisses, err := strconv.ParseInt(fields[3], 10, 64)
			if err == nil {
				data.SslCacheMisses = &sslCacheMisses
			}
		case "42":
			compressBpsIn, err := strconv.ParseInt(fields[3], 10, 64)
			if err == nil {
				data.CompressBpsIn = &compressBpsIn
			}
		case "43":
			compressBpsOut, err := strconv.ParseInt(fields[3], 10, 64)
			if err == nil {
				data.CompressBpsOut = &compressBpsOut
			}
		case "44":
			compressBpsRateLim, err := strconv.ParseInt(fields[3], 10, 64)
			if err == nil {
				data.CompressBpsRateLim = &compressBpsRateLim
			}
		case "45":
			zlibMemUsage, err := strconv.ParseInt(fields[3], 10, 64)
			if err == nil {
				data.ZlibMemUsage = &zlibMemUsage
			}
		case "46":
			maxZlibMemUsage, err := strconv.ParseInt(fields[3], 10, 64)
			if err == nil {
				data.MaxZlibMemUsage = &maxZlibMemUsage
			}
		case "47":
			tasks, err := strconv.ParseInt(fields[3], 10, 64)
			if err == nil {
				data.Tasks = &tasks
			}
		case "48":
			runQ, err := strconv.ParseInt(fields[3], 10, 64)
			if err == nil {
				data.RunQueue = &runQ
			}
		case "49":
			idle, err := strconv.ParseInt(fields[3], 10, 64)
			if err == nil {
				data.IdlePct = &idle
			}
		case "50":
			data.Node = fields[3]
		case "52":
			stopping, err := strconv.ParseInt(fields[3], 10, 64)
			if err == nil {
				data.Stopping = &stopping
			}
		case "53":
			jobs, err := strconv.ParseInt(fields[3], 10, 64)
			if err == nil {
				data.Jobs = &jobs
			}
		case "54":
			unstoppableJ, err := strconv.ParseInt(fields[3], 10, 64)
			if err == nil {
				data.Unstoppable = &unstoppableJ
			}
		case "55":
			listeners, err := strconv.ParseInt(fields[3], 10, 64)
			if err == nil {
				data.Listeners = &listeners
			}
		case "56":
			activePeers, err := strconv.ParseInt(fields[3], 10, 64)
			if err == nil {
				data.ActivePeers = &activePeers
			}
		case "57":
			connPeers, err := strconv.ParseInt(fields[3], 10, 64)
			if err == nil {
				data.ConnectedPeers = &connPeers
			}
		case "58":
			droppedLogs, err := strconv.ParseInt(fields[3], 10, 64)
			if err == nil {
				data.DroppedLogs = &droppedLogs
			}
		case "59":
			busyPolling, err := strconv.ParseInt(fields[3], 10, 64)
			if err == nil {
				data.BusyPolling = &busyPolling
			}
		case "60":
			failedRes, err := strconv.ParseInt(fields[3], 10, 64)
			if err == nil {
				data.FailedResolutions = &failedRes
			}
		case "61":
			totalBOut, err := strconv.ParseInt(fields[3], 10, 64)
			if err == nil {
				data.TotalBytesOut = &totalBOut
			}
		case "62":
			bOutRate, err := strconv.ParseInt(fields[3], 10, 64)
			if err == nil {
				data.BytesOutRate = &bOutRate
			}
		}
	}

	return data
}
