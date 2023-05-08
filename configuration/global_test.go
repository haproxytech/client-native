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

package configuration

import (
	"fmt"
	"testing"

	"github.com/haproxytech/client-native/v4/misc"
	"github.com/haproxytech/client-native/v4/models"
)

func TestGetGlobal(t *testing.T) {
	v, global, err := clientTest.GetGlobalConfiguration("")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	if global.Daemon != "enabled" {
		t.Errorf("Daemon is %v, expected enabled", global.Daemon)
	}
	if *global.Anonkey != 25 {
		t.Errorf("Anonkey is %v, expected 25", *global.Anonkey)
	}
	if len(global.RuntimeAPIs) == 1 {
		if *global.RuntimeAPIs[0].Address != "/var/run/haproxy.sock" {
			t.Errorf("RuntimeAPI.Address is %v, expected /var/run/haproxy.sock", *global.RuntimeAPIs[0].Address)
		}
		if global.RuntimeAPIs[0].Level != "admin" {
			t.Errorf("RuntimeAPI.Level is %v, expected admin", global.RuntimeAPIs[0].Level)
		}
		if global.RuntimeAPIs[0].Mode != "0660" {
			t.Errorf("RuntimeAPI.Mode is %v, expected 0660", global.RuntimeAPIs[0].Mode)
		}
	} else {
		t.Errorf("RuntimeAPI is not set")
	}
	if global.ClusterSecret != "my_secret" {
		t.Errorf("ClusterSecret is %v, expected my_secret", global.ClusterSecret)
	}
	if global.CaBase != "/etc/ssl/certs" {
		t.Errorf("CaBase is %v, expected /etc/ssl/certs", global.CaBase)
	}
	if global.CrtBase != "/etc/ssl/private" {
		t.Errorf("CrtBase is %v, expected /etc/ssl/private", global.CrtBase)
	}
	if global.Nbproc != 4 {
		t.Errorf("Nbproc is %v, expected 4", global.Nbproc)
	}
	if global.Maxconn != 2000 {
		t.Errorf("Maxconn is %v, expected 2000", global.Maxconn)
	}
	if global.ExternalCheck != true {
		t.Errorf("ExternalCheck is false, expected true")
	}
	if len(global.LuaPrependPath) == 2 {
		if *global.LuaPrependPath[0].Path != "/usr/share/haproxy-lua/?/init.lua" {
			t.Errorf("LuaPrependPath.Path is %v, expected /usr/share/haproxy-lua/?/init.lua", *global.LuaPrependPath[0].Path)
		}
		if *global.LuaPrependPath[1].Path != "/usr/share/haproxy-lua/?.lua" {
			t.Errorf("LuaPrependPath.Path is %v, expected /usr/share/haproxy-lua/?.lua", global.LuaPrependPath[1].Path)
		}
		typ := "cpath"
		if global.LuaPrependPath[1].Type != typ {
			t.Errorf("LuaPrependPath.Type is %v, expected cpath", global.LuaPrependPath[1].Type)
		}
	} else {
		t.Errorf("%v LuaPrependPath returned, expected 2", len(global.LuaPrependPath))
	}
	if len(global.LuaLoads) == 2 {
		if *global.LuaLoads[0].File != "/etc/foo.lua" {
			t.Errorf("LuaLoad.File is %v, expected /etc/foo.lua", *global.LuaLoads[0].File)
		}
		if *global.LuaLoads[1].File != "/etc/bar.lua" {
			t.Errorf("LuaLoad.File is %v, expected /etc/bar.lua", global.LuaLoads[1].File)
		}
	} else {
		t.Errorf("%v LuaLoads returned, expected 2", len(global.LuaLoads))
	}
	if len(global.H1CaseAdjusts) == 2 {
		if *global.H1CaseAdjusts[0].From != "host" {
			t.Errorf("H1CaseAdjusts[0].From is %v, expected host", *global.H1CaseAdjusts[0].From)
		}
		if *global.H1CaseAdjusts[0].To != "Host" {
			t.Errorf("H1CaseAdjusts[0].To is %v, expected Host", *global.H1CaseAdjusts[0].To)
		}
		if *global.H1CaseAdjusts[1].From != "content-type" {
			t.Errorf("H1CaseAdjusts[1].From is %v, expected content-type", *global.H1CaseAdjusts[1].From)
		}
		if *global.H1CaseAdjusts[1].To != "Content-Type" {
			t.Errorf("H1CaseAdjusts[1].To is %v, expected Content-Type", *global.H1CaseAdjusts[1].To)
		}
	}
	if global.H1CaseAdjustFile != "/etc/headers.adjust" {
		t.Errorf("H1CaseAdjustFile is %v, expected /etc/headers.adjust", global.H1CaseAdjustFile)
	}
	if global.UID != 1 {
		t.Errorf("UID is %v, expected 1", global.UID)
	}
	if global.Gid != 1 {
		t.Errorf("Gid is %v, expected 1", global.Gid)
	}
	if global.BusyPolling != true {
		t.Errorf("BusyPolling is false, expected true")
	}
	if global.MaxSpreadChecks != 1 {
		t.Errorf("MaxSpreadChecks is %v, expected 1", global.MaxSpreadChecks)
	}
	if global.Maxconnrate != 2 {
		t.Errorf("Maxconnrate is %v, expected 2", global.Maxconnrate)
	}
	if global.Maxcomprate != 3 {
		t.Errorf("Maxcomprate is %v, expected 3", global.Maxcomprate)
	}
	if global.Maxcompcpuusage != 4 {
		t.Errorf("Maxcompcpuusage is %v, expected 4", global.Maxcompcpuusage)
	}
	if global.Maxpipes != 5 {
		t.Errorf("Maxpipes is %v, expected 5", global.Maxpipes)
	}
	if global.Maxsessrate != 6 {
		t.Errorf("Maxsessrate is %v, expected 6", global.Maxsessrate)
	}
	if global.Maxsslconn != 7 {
		t.Errorf("Maxsslconn is %v, expected 7", global.Maxsslconn)
	}
	if global.Maxsslrate != 8 {
		t.Errorf("Maxsslrate is %v, expected 8", global.Maxsslrate)
	}
	if global.Maxzlibmem != 9 {
		t.Errorf("Maxzlibmem is %v, expected 9", global.Maxzlibmem)
	}
	if global.NoQuic != true {
		t.Errorf("NoQuic is false, expected true")
	}
	if global.Noepoll != true {
		t.Errorf("Noepoll is false, expected true")
	}
	if global.Nosplice != true {
		t.Errorf("Nosplice is false, expected true")
	}
	if global.Nogetaddrinfo != true {
		t.Errorf("Nogetaddrinfo is false, expected true")
	}
	if global.Noreuseport != true {
		t.Errorf("Noreuseport is false, expected true")
	}
	if global.ProfilingTasks != "enabled" {
		t.Errorf("ProfilingTasks is %s, expected on", global.ProfilingTasks)
	}
	if global.SpreadChecks != 10 {
		t.Errorf("SpreadChecks is %v, expected 10", global.SpreadChecks)
	}
	if global.WurflOptions.DataFile != "path" {
		t.Errorf("WurflDataFile is %v, expected path", global.WurflOptions.DataFile)
	}
	if global.WurflOptions.InformationList != "wurfl_id,wurfl_root_id,wurfl_isdevroot,wurfl_useragent,wurfl_api_version,wurfl_info,wurfl_last_load_time,wurfl_normalized_useragent" {
		t.Errorf("WurflInformationList is %v, expected wurfl_id,wurfl_root_id,wurfl_isdevroot,wurfl_useragent,wurfl_api_version,wurfl_info,wurfl_last_load_time,wurfl_normalized_useragent", global.WurflOptions.InformationList)
	}
	if global.WurflOptions.InformationListSeparator != "," {
		t.Errorf("WurflInformationListSeparator is %v, expected ,", global.WurflOptions.InformationListSeparator)
	}
	if global.WurflOptions.PatchFile != "path1,path2" {
		t.Errorf("WurflPatchFile is %v, expected path1,path2", global.WurflOptions.PatchFile)
	}
	if global.WurflOptions.CacheSize != 64 {
		t.Errorf("WurflCacheSize is %v, expected 64", global.WurflOptions.CacheSize)
	}
	if global.SslDefaultBindCurves != "X25519:P-256" {
		t.Errorf("SslDefaultBindCurves is %v, expected X25519:P-256", global.SslDefaultBindCurves)
	}
	if global.SslSkipSelfIssuedCa != true {
		t.Errorf("SslSkipSelfIssuedCa is %v, expected enabled", global.SslSkipSelfIssuedCa)
	}
	if global.Node != "node" {
		t.Errorf("Node is %v, expected node", global.Node)
	}
	if global.Description != "description" {
		t.Errorf("Description is %v, expected description", global.Description)
	}
	if global.ExposeExperimentalDirectives != true {
		t.Errorf("ExposeExperimentalDirectives is %v, expected enabled", global.ExposeExperimentalDirectives)
	}
	if global.InsecureForkWanted != true {
		t.Errorf("InsecureForkWanted is %v, expected enabled", global.InsecureForkWanted)
	}
	if global.InsecureSetuidWanted != true {
		t.Errorf("InsecureSetuidWanted is %v, expected enabled", global.InsecureSetuidWanted)
	}
	if global.IssuersChainPath != "issuers-chain-path" {
		t.Errorf("IssuersChainPath is %v, expected issuers-chain-path", global.IssuersChainPath)
	}
	if global.H2WorkaroundBogusWebsocketClients != true {
		t.Errorf("H2WorkaroundBogusWebsocketClients is %v, expected enabled", global.H2WorkaroundBogusWebsocketClients)
	}
	if global.LuaLoadPerThread != "file.ext" {
		t.Errorf("LuaLoadPerThread is %v, expected file.ext", global.LuaLoadPerThread)
	}
	if *global.MworkerMaxReloads != 5 {
		t.Errorf("MworkerMaxReloads is %v, expected 5", global.MworkerMaxReloads)
	}
	if global.NumaCPUMapping != "enabled" {
		t.Errorf("NumaCPUMapping is %v, expected enabled", global.NumaCPUMapping)
	}
	if global.Pp2NeverSendLocal != true {
		t.Errorf("Pp2NeverSendLocal is %v, expected enabled", global.Pp2NeverSendLocal)
	}
	if global.Ulimitn != 10 {
		t.Errorf("Ulimitn is %v, expected 10", global.Ulimitn)
	}
	if global.SetDumpable != true {
		t.Errorf("SetDumpable is %v, expected enabled", global.SetDumpable)
	}
	if global.StrictLimits != true {
		t.Errorf("StrictLimits is %v, expected enabled", global.StrictLimits)
	}
	if *global.Grace != 10000 {
		t.Errorf("Grace is %v, expected 10000", global.Grace)
	}
	if global.SslDefaultServerCiphers != "ECDH+AESGCM:ECDH+CHACHA20:ECDH+AES256:ECDH+AES128:!aNULL:!SHA1:!AESCCM" {
		t.Errorf("SslDefaultServerCiphers is %v, expected %v", global.SslDefaultServerCiphers, "ECDH+AESGCM:ECDH+CHACHA20:ECDH+AES256:ECDH+AES128:!aNULL:!SHA1:!AESCCM")
	}
	if global.SslDefaultServerCiphersuites != "TLS_AES_128_GCM_SHA256:TLS_AES_256_GCM_SHA384:TLS_CHACHA20_POLY1305_SHA256" {
		t.Errorf("SslDefaultServerCiphersuites is %v, expected %v", global.SslDefaultServerCiphersuites, "TLS_AES_128_GCM_SHA256:TLS_AES_256_GCM_SHA384:TLS_CHACHA20_POLY1305_SHA256")
	}
	if global.Chroot != "/var/www" {
		t.Errorf("Chroot is %v, expected /var/www", global.Chroot)
	}
	if *global.HardStopAfter != 2000 {
		t.Errorf("HardStopAfter is %v, expected 20000", global.HardStopAfter)
	}
	if global.Localpeer != "test" {
		t.Errorf("Localpeer is %v, expected test", global.Localpeer)
	}
	if global.User != "thomas" {
		t.Errorf("User is %v, expected thomas", global.User)
	}
	if global.Group != "anderson" {
		t.Errorf("Group is %v, expected anderson", global.Group)
	}
	if global.Nbthread != 128 {
		t.Errorf("Nbthread is %v, expected 128", global.Nbthread)
	}
	if global.Pidfile != "pidfile.text" {
		t.Errorf("Pidfile is %v, expected pidfile.text", global.Pidfile)
	}
	if global.SslDefaultBindCiphers != "ECDH+AESGCM:ECDH+CHACHA20" {
		t.Errorf("SslDefaultBindCiphers is %v, expected %v", global.SslDefaultBindCiphers, "ECDH+AESGCM:ECDH+CHACHA20")
	}
	if global.SslDefaultBindCiphersuites != "TLS_AES_128_GCM_SHA256:TLS_AES_256_GCM_SHA384" {
		t.Errorf("SslDefaultBindCiphersuites is %v, expected %v", global.SslDefaultBindCiphersuites, "TLS_AES_128_GCM_SHA256:TLS_AES_256_GCM_SHA384")
	}
	if global.SslDefaultServerOptions != "ssl-min-ver TLSv1.1 no-tls-tickets" {
		t.Errorf("SslDefaultServerOptions is %v, expected ssl-min-ver TLSv1.1 no-tls-tickets", global.SslDefaultServerOptions)
	}
	if *global.TuneOptions.BuffersLimit != 11 {
		t.Errorf("BuffersLimit is %v, expected 11", global.TuneOptions.BuffersLimit)
	}
	if global.TuneOptions.BuffersReserve != 12 {
		t.Errorf("BuffersReserve is %v, expected 12", global.TuneOptions.BuffersReserve)
	}
	if global.TuneOptions.Bufsize != 13 {
		t.Errorf("Bufsize is %v, expected 13", global.TuneOptions.Bufsize)
	}
	if global.TuneOptions.CompMaxlevel != 14 {
		t.Errorf("CompMaxlevel is %v, expected 14", global.TuneOptions.CompMaxlevel)
	}
	if global.TuneOptions.H2HeaderTableSize != 15 {
		t.Errorf("H2HeaderTableSize is %v, expected 15", global.TuneOptions.H2HeaderTableSize)
	}
	if *global.TuneOptions.H2InitialWindowSize != 16 {
		t.Errorf("H2InitialWindowSize is %v, expected 16", global.TuneOptions.H2InitialWindowSize)
	}
	if global.TuneOptions.H2MaxConcurrentStreams != 17 {
		t.Errorf("H2MaxConcurrentStreams is %v, expected 17", global.TuneOptions.H2MaxConcurrentStreams)
	}
	if global.TuneOptions.H2MaxFrameSize != 18 {
		t.Errorf("H2MaxFrameSize is %v, expected 18", global.TuneOptions.H2MaxFrameSize)
	}
	if global.TuneOptions.HTTPCookielen != 19 {
		t.Errorf("HTTPCookielen is %v, expected 19", global.TuneOptions.HTTPCookielen)
	}
	if global.TuneOptions.HTTPLogurilen != 20 {
		t.Errorf("HTTPLogurilen is %v, expected 20", global.TuneOptions.HTTPLogurilen)
	}
	if global.TuneOptions.HTTPMaxhdr != 21 {
		t.Errorf("HTTPMaxhdr is %v, expected 21", global.TuneOptions.HTTPMaxhdr)
	}
	if *global.TuneOptions.Idletimer != 22 {
		t.Errorf("Idletimer is %v, expected 22", global.TuneOptions.Idletimer)
	}
	if global.TuneOptions.ListenerDefaultShards != "by-process" {
		t.Errorf("ListenerDefaultShards is %v, expected by-process", global.TuneOptions.ListenerDefaultShards)
	}
	if global.TuneOptions.LuaForcedYield != 23 {
		t.Errorf("LuaForcedYield is %v, expected 23", global.TuneOptions.LuaForcedYield)
	}
	if global.TuneOptions.LuaMaxmem != true {
		t.Errorf("Maxzlibmem is false, expected true")
	}
	if *global.TuneOptions.LuaSessionTimeout != 25 {
		t.Errorf("LuaSessionTimeout is %v, expected 25", global.TuneOptions.LuaSessionTimeout)
	}
	if *global.TuneOptions.LuaBurstTimeout != 205 {
		t.Errorf("LuaSessionTimeout is %v, expected 205", global.TuneOptions.LuaBurstTimeout)
	}
	if *global.TuneOptions.LuaTaskTimeout != 26 {
		t.Errorf("LuaTaskTimeout is %v, expected 26", global.TuneOptions.LuaTaskTimeout)
	}
	if *global.TuneOptions.LuaServiceTimeout != 27 {
		t.Errorf("LuaServiceTimeout is %v, expected 27", global.TuneOptions.LuaServiceTimeout)
	}
	if global.TuneOptions.Maxaccept != 28 {
		t.Errorf("Maxaccept is %v, expected 28", global.TuneOptions.Maxaccept)
	}
	if global.TuneOptions.Maxpollevents != 29 {
		t.Errorf("Maxpollevents is %v, expected 29", global.TuneOptions.Maxpollevents)
	}
	if global.TuneOptions.Maxrewrite != 30 {
		t.Errorf("Maxrewrite is %v, expected 30", global.TuneOptions.Maxrewrite)
	}
	if *global.TuneOptions.PatternCacheSize != 31 {
		t.Errorf("PatternCacheSize is %v, expected 31", global.TuneOptions.PatternCacheSize)
	}
	if global.TuneOptions.PeersMaxUpdatesAtOnce != 200 {
		t.Errorf("PeersMaxUpdatesAtOnce is %v, expected 200", global.TuneOptions.PeersMaxUpdatesAtOnce)
	}
	if global.TuneOptions.Pipesize != 32 {
		t.Errorf("Pipesize is %v, expected 32", global.TuneOptions.Pipesize)
	}
	if global.TuneOptions.PoolHighFdRatio != 33 {
		t.Errorf("PoolHighFdRatio is %v, expected 33", global.TuneOptions.PoolHighFdRatio)
	}
	if global.TuneOptions.PoolLowFdRatio != 34 {
		t.Errorf("PoolLowFdRatio is %v, expected 34", global.TuneOptions.PoolLowFdRatio)
	}
	if *global.TuneOptions.RcvbufClient != 35 {
		t.Errorf("RcvbufClient is %v, expected 35", global.TuneOptions.RcvbufClient)
	}
	if *global.TuneOptions.RcvbufServer != 36 {
		t.Errorf("RcvbufServer is %v, expected 36", global.TuneOptions.RcvbufServer)
	}
	if global.TuneOptions.RecvEnough != 37 {
		t.Errorf("RecvEnough is %v, expected 37", global.TuneOptions.RecvEnough)
	}
	if global.TuneOptions.RunqueueDepth != 38 {
		t.Errorf("RunqueueDepth is %v, expected 38", global.TuneOptions.RunqueueDepth)
	}
	if *global.TuneOptions.SndbufClient != 39 {
		t.Errorf("SndbufClient is %v, expected 39", global.TuneOptions.SndbufClient)
	}
	if *global.TuneOptions.SndbufServer != 40 {
		t.Errorf("SndbufServer is %v, expected 40", global.TuneOptions.SndbufServer)
	}
	if *global.TuneOptions.SslCachesize != 41 {
		t.Errorf("SslCachesize is %v, expected 41", global.TuneOptions.SslCachesize)
	}
	if global.TuneOptions.SslKeylog != "enabled" {
		t.Errorf("SslKeylog is %v, expected enabled", global.TuneOptions.SslKeylog)
	}
	if *global.TuneOptions.SslLifetime != 43 {
		t.Errorf("SslLifetime is %v, expected 43", global.TuneOptions.SslLifetime)
	}
	if *global.TuneOptions.SslMaxrecord != 44 {
		t.Errorf("SslMaxrecord is %v, expected 44", global.TuneOptions.SslMaxrecord)
	}
	if global.TuneOptions.SslDefaultDhParam != 45 {
		t.Errorf("SslDefaultDhParam is %v, expected 45", global.TuneOptions.SslDefaultDhParam)
	}
	if global.TuneOptions.SslCtxCacheSize != 46 {
		t.Errorf("SslCtxCacheSize is %v, expected 46", global.TuneOptions.SslCtxCacheSize)
	}
	if *global.TuneOptions.SslCaptureBufferSize != 47 {
		t.Errorf("SslCaptureBufferSize is %v, expected 47", global.TuneOptions.SslCaptureBufferSize)
	}
	if *global.TuneOptions.SslOcspUpdateMaxDelay != 48 {
		t.Errorf("SslCaptureBufferSize is %v, expected 48", global.TuneOptions.SslCaptureBufferSize)
	}
	if *global.TuneOptions.SslOcspUpdateMinDelay != 49 {
		t.Errorf("SslCaptureBufferSize is %v, expected 49", global.TuneOptions.SslCaptureBufferSize)
	}
	if *global.TuneOptions.StickCounters != 50 {
		t.Errorf("SslCaptureBufferSize is %v, expected 50", global.TuneOptions.SslCaptureBufferSize)
	}
	if *global.TuneOptions.VarsGlobalMaxSize != 51 {
		t.Errorf("VarsGlobalMaxSize is %v, expected 51", global.TuneOptions.VarsGlobalMaxSize)
	}
	if *global.TuneOptions.VarsProcMaxSize != 52 {
		t.Errorf("VarsProcMaxSize is %v, expected 52", global.TuneOptions.VarsProcMaxSize)
	}
	if *global.TuneOptions.VarsReqresMaxSize != 53 {
		t.Errorf("VarsReqresMaxSize is %v, expected 53", global.TuneOptions.VarsReqresMaxSize)
	}
	if *global.TuneOptions.VarsSessMaxSize != 54 {
		t.Errorf("VarsSessMaxSize is %v, expected 54", global.TuneOptions.VarsSessMaxSize)
	}
	if *global.TuneOptions.VarsTxnMaxSize != 55 {
		t.Errorf("VarsTxnMaxSize is %v, expected 55", global.TuneOptions.VarsTxnMaxSize)
	}
	if *global.TuneOptions.QuicFrontendConnTcBuffersLimit != 10 {
		t.Errorf("QuicFrontendConnTcBuffersLimit is %v, expected 10", global.TuneOptions.QuicFrontendConnTcBuffersLimit)
	}
	if *global.TuneOptions.QuicFrontendMaxIdleTimeout != 10000 {
		t.Errorf("QuicFrontendMaxIdleTimeout is %v, expected 10000", global.TuneOptions.QuicFrontendMaxIdleTimeout)
	}
	if *global.TuneOptions.QuicFrontendMaxStreamsBidi != 100 {
		t.Errorf("QuicFrontendMaxStreamsBidi is %v, expected 100", global.TuneOptions.QuicFrontendMaxStreamsBidi)
	}
	if *global.TuneOptions.QuicMaxFrameLoss != 5 {
		t.Errorf("QuicMaxFrameLoss is %v, expected 5", global.TuneOptions.QuicMaxFrameLoss)
	}
	if *global.TuneOptions.QuicRetryThreshold != 5 {
		t.Errorf("QuicRetryThreshold is %v, expected 5", global.TuneOptions.QuicRetryThreshold)
	}
	if global.TuneOptions.QuicSocketOwner != "connection" {
		t.Errorf("QuicSocketOwner is %v, expected connection", global.TuneOptions.QuicSocketOwner)
	}
	if global.TuneOptions.ZlibMemlevel != 54 {
		t.Errorf("ZlibMemlevel is %v, expected 54", global.TuneOptions.ZlibMemlevel)
	}
	if global.TuneOptions.FdEdgeTriggered != "enabled" {
		t.Errorf("FdEdgeTriggered is %v, expected enabled", global.TuneOptions.FdEdgeTriggered)
	}
	if global.TuneOptions.ZlibWindowsize != 55 {
		t.Errorf("ZlibWindowsize is %v, expected 55", global.TuneOptions.ZlibWindowsize)
	}
	if *global.TuneOptions.MemoryHotSize != 56 {
		t.Errorf("MemoryHotSize is %v, expected 56", global.TuneOptions.MemoryHotSize)
	}
	if global.TuneOptions.H2BeInitialWindowSize != 201 {
		t.Errorf("H2BeInitialWindowSize is %v, expected 201", global.TuneOptions.H2BeInitialWindowSize)
	}
	if global.TuneOptions.H2BeMaxConcurrentStreams != 202 {
		t.Errorf("H2BeMaxConcurrentStreams is %v, expected 202", global.TuneOptions.H2BeMaxConcurrentStreams)
	}
	if global.TuneOptions.H2FeInitialWindowSize != 203 {
		t.Errorf("H2FeInitialWindowSize is %v, expected 203", global.TuneOptions.H2FeInitialWindowSize)
	}
	if global.TuneOptions.H2FeMaxConcurrentStreams != 204 {
		t.Errorf("H2FeMaxConcurrentStreams is %v, expected 204", global.TuneOptions.H2FeMaxConcurrentStreams)
	}
	if global.ThreadGroups != 1 {
		t.Errorf("ThreadGroups is %v, expected 1", global.ThreadGroups)
	}
	if len(global.ThreadGroupLines) == 1 {
		if *global.ThreadGroupLines[0].Group != "first" {
			t.Errorf("ThreadGroup[0] Group is %v, expected first", *global.ThreadGroupLines[0].Group)
		}
		if *global.ThreadGroupLines[0].NumOrRange != "1-16" {
			t.Errorf("ThreadGroup[0] NumOrRange is %v, expected 1-16", *global.ThreadGroupLines[0].NumOrRange)
		}
	}
	if *global.StatsMaxconn != 20 {
		t.Errorf("StatsMaxconn is %v, expected 20", *global.StatsMaxconn)
	}
	if global.DeviceAtlasOptions.JSONFile != "atlas.json" {
		t.Errorf("DeviceAtlasOptions.JSONFile is %v, expected atlas.json", global.DeviceAtlasOptions.JSONFile)
	}
	if global.DeviceAtlasOptions.LogLevel != "1" {
		t.Errorf("DeviceAtlasOptions.LogLevel is %v, expected 1", global.DeviceAtlasOptions.LogLevel)
	}
	if global.DeviceAtlasOptions.Separator != "-" {
		t.Errorf("DeviceAtlasOptions.Separator is %v, expected -", global.DeviceAtlasOptions.Separator)
	}
	if global.DeviceAtlasOptions.PropertiesCookie != "chocolate" {
		t.Errorf("DeviceAtlasOptions.PropertiesCookie is %v, expected chocolate", global.DeviceAtlasOptions.PropertiesCookie)
	}
	if global.FiftyOneDegreesOptions.DataFile != "51.file" {
		t.Errorf("FiftyOneDegreesOptions.DataFile is %v, expected 51.file", global.FiftyOneDegreesOptions.DataFile)
	}
	if global.FiftyOneDegreesOptions.PropertyNameList != "first second third fourth fifth" {
		t.Errorf("FiftyOneDegreesOptions.PropertyNameList is %v, expected 'first second third fourth fifth'", global.FiftyOneDegreesOptions.PropertyNameList)
	}
	if global.FiftyOneDegreesOptions.PropertySeparator != "/" {
		t.Errorf("FiftyOneDegreesOptions.PropertySeparator is %v, expected /", global.FiftyOneDegreesOptions.PropertySeparator)
	}
	if global.FiftyOneDegreesOptions.CacheSize != 51 {
		t.Errorf("FiftyOneDegreesOptions.CacheSize is %v, expected 51", global.FiftyOneDegreesOptions.CacheSize)
	}
	if global.Quiet != true {
		t.Errorf("Quiet is false, expected true")
	}
	if global.ZeroWarning != true {
		t.Errorf("ZeroWarning is false, expected true")
	}
	if len(global.SslEngines) == 3 {
		if *global.SslEngines[0].Name != "first" {
			t.Errorf("SslEngines[0] Name is %v, expected first", *global.SslEngines[0].Name)
		}
		if *global.SslEngines[0].Algorithms != "" {
			t.Errorf("SslEngines[0] Algorithms is %v, expected empty", *global.SslEngines[0].Algorithms)
		}
		if *global.SslEngines[1].Name != "second" {
			t.Errorf("SslEngines[1] Name is %v, expected second", *global.SslEngines[1].Name)
		}
		if *global.SslEngines[1].Algorithms != "RSA,DSA,DH,EC,RAND" {
			t.Errorf("SslEngines[1] Algorithms is %v, expected RSA,DSA,DH,EC,RAND", *global.SslEngines[1].Algorithms)
		}
		if *global.SslEngines[2].Name != "third" {
			t.Errorf("SslEngines[2] Name is %v, expected third", *global.SslEngines[2].Name)
		}
		if *global.SslEngines[2].Algorithms != "CIPHERS,DIGESTS,PKEY,PKEY_CRYPTO,PKEY_ASN1" {
			t.Errorf("SslEngines[2] Algorithms is %v, expected CIPHERS,DIGESTS,PKEY,PKEY_CRYPTO,PKEY_ASN1", *global.SslEngines[2].Algorithms)
		}
	}
	if global.SslDhParamFile != "ssl-dh-param-file.txt" {
		t.Errorf("SslDhParamFile is %v, expected ssl-dh-param-file.txt", global.SslDhParamFile)
	}
	if global.SslServerVerify != "required" {
		t.Errorf("SslServerVerify is %v, expected required", global.SslServerVerify)
	}
	if len(global.SetVars) == 3 {
		if *global.SetVars[0].Name != "proc.current_state" {
			t.Errorf("SetVars[0] Name is %v, expected proc.current_state", *global.SetVars[0].Name)
		}
		if *global.SetVars[0].Expr != "str(primary)" {
			t.Errorf("SetVars[0] Expr is %v, expected str(primary)", *global.SetVars[0].Expr)
		}
		if *global.SetVars[1].Name != "proc.prio" {
			t.Errorf("SetVars[1] Name is %v, expected proc.prio", *global.SetVars[1].Name)
		}
		if *global.SetVars[1].Expr != "int(100)" {
			t.Errorf("SetVars[1] Expr is %v, expected int(100)", *global.SetVars[1].Expr)
		}
		if *global.SetVars[2].Name != "proc.threshold" {
			t.Errorf("SetVars[2] Name is %v, expected proc.threshold", *global.SetVars[2].Name)
		}
		if *global.SetVars[2].Expr != "int(200),sub(proc.prio)" {
			t.Errorf("SetVars[2] Expr is %v, expected int(200),sub(proc.prio)", *global.SetVars[2].Expr)
		}
	} else {
		t.Errorf("SetVars lenght is %v, expected 3", len(global.SetVars))
	}
	if len(global.SetVarFmts) == 2 {
		if *global.SetVarFmts[0].Name != "proc.bootid" {
			t.Errorf("SetVars[0] Name is %v, expected proc.bootid", *global.SetVarFmts[0].Name)
		}
		if *global.SetVarFmts[0].Format != `"%pid|%t"` {
			t.Errorf("SetVars[0] Format is %v, expected %%pid|%%t", *global.SetVarFmts[0].Format)
		}
		if *global.SetVarFmts[1].Name != "proc.current_state" {
			t.Errorf("SetVars[1] Name is %v, expected proc.current_state", *global.SetVarFmts[1].Name)
		}
		if *global.SetVarFmts[1].Format != `"primary"` {
			t.Errorf("SetVars[1] Format is %v, expected \"primary\"", *global.SetVarFmts[1].Format)
		}
	} else {
		t.Errorf("SetVarFmts lenght is %v, expected 2", len(global.SetVarFmts))
	}
	if len(global.PresetEnvs) == 1 {
		if *global.PresetEnvs[0].Name != "first" {
			t.Errorf("esetEnvs[0] Name is %v, expected first", *global.PresetEnvs[0].Name)
		}
		if *global.PresetEnvs[0].Value != "order" {
			t.Errorf("esetEnvs[0] Value is %v, expected order", *global.PresetEnvs[0].Name)
		}
	}
	if len(global.SetEnvs) == 1 {
		if *global.SetEnvs[0].Name != "third" {
			t.Errorf("setEnvs[0] Name is %v, expected third", *global.PresetEnvs[0].Name)
		}
		if *global.SetEnvs[0].Value != "sister" {
			t.Errorf("setEnvs[0] Value is %v, expected sister", *global.PresetEnvs[0].Name)
		}
	}
	if global.Resetenv != "first second" {
		t.Errorf("Resetenv is %v, expected first second", global.Resetenv)
	}
	if global.Unsetenv != "third fourth" {
		t.Errorf("Unsetenv is %v, expected third fourth", global.Resetenv)
	}
	if global.DefaultPath == nil {
		t.Error("DefaultPath is nil, expected not nil")
	} else {
		if global.DefaultPath.Type != "origin" {
			t.Errorf("DefaultPath Type is %v, expected origin", global.DefaultPath.Type)
		}
		if global.DefaultPath.Path != "/some/path" {
			t.Errorf("DefaultPath Path is %v, expected /some/path", global.DefaultPath.Path)
		}
	}
	if global.SslDefaultBindSigalgs != "RSA+SHA256" {
		t.Errorf("SslDefaultBindSigalgs is %v, expected RSA+SHA256", global.SslDefaultBindSigalgs)
	}
	if global.SslDefaultBindClientSigalgs != "ECDSA+SHA256:RSA+SHA256" {
		t.Errorf("SslDefaultBindClientSigalgs is %v, expected ECDSA+SHA256:RSA+SHA256", global.SslDefaultBindClientSigalgs)
	}
}

func TestPutGlobal(t *testing.T) {
	tOut := int64(3600)
	n := "1/1"
	v := "0"
	a := "/var/run/haproxy.sock"
	f := "/etc/foo.lua"
	luaPrependPath := "/usr/share/haproxy-lua/?/init.lua"
	enabled := "enabled"
	g := &models.Global{
		Daemon: "enabled",
		CPUMaps: []*models.CPUMap{
			{
				Process: &n,
				CPUSet:  &v,
			},
		},
		RuntimeAPIs: []*models.RuntimeAPI{
			{
				Address: &a,
				BindParams: models.BindParams{
					Level: "admin",
				},
			},
		},
		Maxconn:               1000,
		SslDefaultBindCiphers: "test",
		SslDefaultBindOptions: "ssl-min-ver TLSv1.0 no-tls-tickets",
		StatsTimeout:          &tOut,
		TuneSslDefaultDhParam: 1024,
		ExternalCheck:         false,
		LuaPrependPath: []*models.LuaPrependPath{
			{
				Path: &luaPrependPath,
				Type: "cpath",
			},
		},
		LuaLoads: []*models.LuaLoad{
			{
				File: &f,
			},
		},
		LogSendHostname: &models.GlobalLogSendHostname{
			Enabled: &enabled,
			Param:   "something",
		},
		TuneOptions: &models.GlobalTuneOptions{
			PeersMaxUpdatesAtOnce:          100,
			QuicFrontendConnTcBuffersLimit: nil,
			QuicFrontendMaxIdleTimeout:     misc.Int64P(5000),
			QuicSocketOwner:                "listener",
			SslOcspUpdateMaxDelay:          misc.Int64P(48),
			SslOcspUpdateMinDelay:          misc.Int64P(49),
			StickCounters:                  misc.Int64P(50),
		},
		UID:                    1234,
		WurflOptions:           &models.GlobalWurflOptions{},
		DeviceAtlasOptions:     &models.GlobalDeviceAtlasOptions{},
		FiftyOneDegreesOptions: &models.GlobalFiftyOneDegreesOptions{},
		StatsMaxconn:           misc.Int64P(30),
		Anonkey:                misc.Int64P(40),
		NumaCPUMapping:         "disabled",
		DefaultPath: &models.GlobalDefaultPath{
			Type: "origin",
			Path: "/some/other/path",
		},
		NoQuic:        false,
		ClusterSecret: "",
	}

	err := clientTest.PushGlobalConfiguration(g, "", version)

	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	ver, global, err := clientTest.GetGlobalConfiguration("")
	if err != nil {
		t.Error(err.Error())
	}

	var givenJSON []byte
	givenJSON, err = g.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	var onDiskJSON []byte
	onDiskJSON, err = global.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	if string(givenJSON) != string(onDiskJSON) {
		fmt.Printf("Created global: %v\n", string(onDiskJSON))
		fmt.Printf("Given global: %v\n", string(givenJSON))
		t.Error("Created global not equal to given global")
	}

	if ver != version {
		t.Error("Version not incremented!")
	}

	err = clientTest.PushGlobalConfiguration(g, "", 55)

	if err == nil {
		t.Error("Should have returned version conflict.")
	}
}
