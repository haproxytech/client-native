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
	"strconv"
)

// Equal checks if two structs of type Global are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b Global
//	equal := a.Equal(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b Global
//	equal := a.Equal(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s Global) Equal(t Global, opts ...Options) bool {
	opt := getOptions(opts...)

	if !CheckSameNilAndLen(s.CPUMaps, t.CPUMaps, opt) {
		return false
	} else {
		for i := range s.CPUMaps {
			if !s.CPUMaps[i].Equal(*t.CPUMaps[i], opt) {
				return false
			}
		}
	}

	if !CheckSameNilAndLen(s.H1CaseAdjusts, t.H1CaseAdjusts, opt) {
		return false
	} else {
		for i := range s.H1CaseAdjusts {
			if !s.H1CaseAdjusts[i].Equal(*t.H1CaseAdjusts[i], opt) {
				return false
			}
		}
	}

	if !CheckSameNilAndLen(s.PresetEnvs, t.PresetEnvs, opt) {
		return false
	} else {
		for i := range s.PresetEnvs {
			if !s.PresetEnvs[i].Equal(*t.PresetEnvs[i], opt) {
				return false
			}
		}
	}

	if !CheckSameNilAndLen(s.RuntimeAPIs, t.RuntimeAPIs, opt) {
		return false
	} else {
		for i := range s.RuntimeAPIs {
			if !s.RuntimeAPIs[i].Equal(*t.RuntimeAPIs[i], opt) {
				return false
			}
		}
	}

	if !CheckSameNilAndLen(s.SetEnvs, t.SetEnvs, opt) {
		return false
	} else {
		for i := range s.SetEnvs {
			if !s.SetEnvs[i].Equal(*t.SetEnvs[i], opt) {
				return false
			}
		}
	}

	if !CheckSameNilAndLen(s.SetVarFmts, t.SetVarFmts, opt) {
		return false
	} else {
		for i := range s.SetVarFmts {
			if !s.SetVarFmts[i].Equal(*t.SetVarFmts[i], opt) {
				return false
			}
		}
	}

	if !CheckSameNilAndLen(s.SetVars, t.SetVars, opt) {
		return false
	} else {
		for i := range s.SetVars {
			if !s.SetVars[i].Equal(*t.SetVars[i], opt) {
				return false
			}
		}
	}

	if !CheckSameNilAndLen(s.SslEngines, t.SslEngines, opt) {
		return false
	} else {
		for i := range s.SslEngines {
			if !s.SslEngines[i].Equal(*t.SslEngines[i], opt) {
				return false
			}
		}
	}

	if !CheckSameNilAndLen(s.ThreadGroupLines, t.ThreadGroupLines, opt) {
		return false
	} else {
		for i := range s.ThreadGroupLines {
			if !s.ThreadGroupLines[i].Equal(*t.ThreadGroupLines[i], opt) {
				return false
			}
		}
	}

	if !equalPointers(s.Anonkey, t.Anonkey) {
		return false
	}

	if s.BusyPolling != t.BusyPolling {
		return false
	}

	if s.CaBase != t.CaBase {
		return false
	}

	if s.Chroot != t.Chroot {
		return false
	}

	if !equalPointers(s.CloseSpreadTime, t.CloseSpreadTime) {
		return false
	}

	if s.ClusterSecret != t.ClusterSecret {
		return false
	}

	if s.CrtBase != t.CrtBase {
		return false
	}

	if s.Daemon != t.Daemon {
		return false
	}

	if s.DefaultPath == nil || t.DefaultPath == nil {
		if s.DefaultPath != nil || t.DefaultPath != nil {
			if opt.NilSameAsEmpty {
				empty := &GlobalDefaultPath{}
				if s.DefaultPath == nil {
					if !(t.DefaultPath.Equal(*empty)) {
						return false
					}
				}
				if t.DefaultPath == nil {
					if !(s.DefaultPath.Equal(*empty)) {
						return false
					}
				}
			} else {
				return false
			}
		}
	} else if !s.DefaultPath.Equal(*t.DefaultPath, opt) {
		return false
	}

	if s.Description != t.Description {
		return false
	}

	if s.DeviceAtlasOptions == nil || t.DeviceAtlasOptions == nil {
		if s.DeviceAtlasOptions != nil || t.DeviceAtlasOptions != nil {
			if opt.NilSameAsEmpty {
				empty := &GlobalDeviceAtlasOptions{}
				if s.DeviceAtlasOptions == nil {
					if !(t.DeviceAtlasOptions.Equal(*empty)) {
						return false
					}
				}
				if t.DeviceAtlasOptions == nil {
					if !(s.DeviceAtlasOptions.Equal(*empty)) {
						return false
					}
				}
			} else {
				return false
			}
		}
	} else if !s.DeviceAtlasOptions.Equal(*t.DeviceAtlasOptions, opt) {
		return false
	}

	if s.ExposeExperimentalDirectives != t.ExposeExperimentalDirectives {
		return false
	}

	if s.ExternalCheck != t.ExternalCheck {
		return false
	}

	if s.FiftyOneDegreesOptions == nil || t.FiftyOneDegreesOptions == nil {
		if s.FiftyOneDegreesOptions != nil || t.FiftyOneDegreesOptions != nil {
			if opt.NilSameAsEmpty {
				empty := &GlobalFiftyOneDegreesOptions{}
				if s.FiftyOneDegreesOptions == nil {
					if !(t.FiftyOneDegreesOptions.Equal(*empty)) {
						return false
					}
				}
				if t.FiftyOneDegreesOptions == nil {
					if !(s.FiftyOneDegreesOptions.Equal(*empty)) {
						return false
					}
				}
			} else {
				return false
			}
		}
	} else if !s.FiftyOneDegreesOptions.Equal(*t.FiftyOneDegreesOptions, opt) {
		return false
	}

	if s.Gid != t.Gid {
		return false
	}

	if !equalPointers(s.Grace, t.Grace) {
		return false
	}

	if s.Group != t.Group {
		return false
	}

	if s.H1CaseAdjustFile != t.H1CaseAdjustFile {
		return false
	}

	if s.H2WorkaroundBogusWebsocketClients != t.H2WorkaroundBogusWebsocketClients {
		return false
	}

	if !equalPointers(s.HardStopAfter, t.HardStopAfter) {
		return false
	}

	if s.HttpclientResolversDisabled != t.HttpclientResolversDisabled {
		return false
	}

	if s.HttpclientResolversID != t.HttpclientResolversID {
		return false
	}

	if s.HttpclientResolversPrefer != t.HttpclientResolversPrefer {
		return false
	}

	if s.HttpclientRetries != t.HttpclientRetries {
		return false
	}

	if s.HttpclientSslCaFile != t.HttpclientSslCaFile {
		return false
	}

	if !equalPointers(s.HttpclientSslVerify, t.HttpclientSslVerify) {
		return false
	}

	if !equalPointers(s.HttpclientTimeoutConnect, t.HttpclientTimeoutConnect) {
		return false
	}

	if s.InsecureForkWanted != t.InsecureForkWanted {
		return false
	}

	if s.InsecureSetuidWanted != t.InsecureSetuidWanted {
		return false
	}

	if s.IssuersChainPath != t.IssuersChainPath {
		return false
	}

	if s.LimitedQuic != t.LimitedQuic {
		return false
	}

	if s.LoadServerStateFromFile != t.LoadServerStateFromFile {
		return false
	}

	if s.Localpeer != t.Localpeer {
		return false
	}

	if s.LogSendHostname == nil || t.LogSendHostname == nil {
		if s.LogSendHostname != nil || t.LogSendHostname != nil {
			if opt.NilSameAsEmpty {
				empty := &GlobalLogSendHostname{}
				if s.LogSendHostname == nil {
					if !(t.LogSendHostname.Equal(*empty)) {
						return false
					}
				}
				if t.LogSendHostname == nil {
					if !(s.LogSendHostname.Equal(*empty)) {
						return false
					}
				}
			} else {
				return false
			}
		}
	} else if !s.LogSendHostname.Equal(*t.LogSendHostname, opt) {
		return false
	}

	if s.LuaLoadPerThread != t.LuaLoadPerThread {
		return false
	}

	if !CheckSameNilAndLen(s.LuaLoads, t.LuaLoads, opt) {
		return false
	} else {
		for i := range s.LuaLoads {
			if !s.LuaLoads[i].Equal(*t.LuaLoads[i], opt) {
				return false
			}
		}
	}

	if !CheckSameNilAndLen(s.LuaPrependPath, t.LuaPrependPath, opt) {
		return false
	} else {
		for i := range s.LuaPrependPath {
			if !s.LuaPrependPath[i].Equal(*t.LuaPrependPath[i], opt) {
				return false
			}
		}
	}

	if s.MasterWorker != t.MasterWorker {
		return false
	}

	if !equalPointers(s.MaxSpreadChecks, t.MaxSpreadChecks) {
		return false
	}

	if s.Maxcompcpuusage != t.Maxcompcpuusage {
		return false
	}

	if s.Maxcomprate != t.Maxcomprate {
		return false
	}

	if s.Maxconn != t.Maxconn {
		return false
	}

	if s.Maxconnrate != t.Maxconnrate {
		return false
	}

	if s.Maxpipes != t.Maxpipes {
		return false
	}

	if s.Maxsessrate != t.Maxsessrate {
		return false
	}

	if s.Maxsslconn != t.Maxsslconn {
		return false
	}

	if s.Maxsslrate != t.Maxsslrate {
		return false
	}

	if s.Maxzlibmem != t.Maxzlibmem {
		return false
	}

	if !equalPointers(s.MworkerMaxReloads, t.MworkerMaxReloads) {
		return false
	}

	if s.Nbproc != t.Nbproc {
		return false
	}

	if s.Nbthread != t.Nbthread {
		return false
	}

	if s.NoQuic != t.NoQuic {
		return false
	}

	if s.Node != t.Node {
		return false
	}

	if s.Noepoll != t.Noepoll {
		return false
	}

	if s.Noevports != t.Noevports {
		return false
	}

	if s.Nogetaddrinfo != t.Nogetaddrinfo {
		return false
	}

	if s.Nokqueue != t.Nokqueue {
		return false
	}

	if s.Nopoll != t.Nopoll {
		return false
	}

	if s.Noreuseport != t.Noreuseport {
		return false
	}

	if s.Nosplice != t.Nosplice {
		return false
	}

	if s.NumaCPUMapping != t.NumaCPUMapping {
		return false
	}

	if s.Pidfile != t.Pidfile {
		return false
	}

	if s.Pp2NeverSendLocal != t.Pp2NeverSendLocal {
		return false
	}

	if s.PreallocFd != t.PreallocFd {
		return false
	}

	if s.ProfilingTasks != t.ProfilingTasks {
		return false
	}

	if s.Quiet != t.Quiet {
		return false
	}

	if s.Resetenv != t.Resetenv {
		return false
	}

	if s.ServerStateBase != t.ServerStateBase {
		return false
	}

	if s.ServerStateFile != t.ServerStateFile {
		return false
	}

	if s.SetDumpable != t.SetDumpable {
		return false
	}

	if s.Setcap != t.Setcap {
		return false
	}

	if s.SpreadChecks != t.SpreadChecks {
		return false
	}

	if s.SslDefaultBindCiphers != t.SslDefaultBindCiphers {
		return false
	}

	if s.SslDefaultBindCiphersuites != t.SslDefaultBindCiphersuites {
		return false
	}

	if s.SslDefaultBindClientSigalgs != t.SslDefaultBindClientSigalgs {
		return false
	}

	if s.SslDefaultBindCurves != t.SslDefaultBindCurves {
		return false
	}

	if s.SslDefaultBindOptions != t.SslDefaultBindOptions {
		return false
	}

	if s.SslDefaultBindSigalgs != t.SslDefaultBindSigalgs {
		return false
	}

	if s.SslDefaultServerCiphers != t.SslDefaultServerCiphers {
		return false
	}

	if s.SslDefaultServerCiphersuites != t.SslDefaultServerCiphersuites {
		return false
	}

	if s.SslDefaultServerClientSigalgs != t.SslDefaultServerClientSigalgs {
		return false
	}

	if s.SslDefaultServerCurves != t.SslDefaultServerCurves {
		return false
	}

	if s.SslDefaultServerOptions != t.SslDefaultServerOptions {
		return false
	}

	if s.SslDefaultServerSigalgs != t.SslDefaultServerSigalgs {
		return false
	}

	if s.SslDhParamFile != t.SslDhParamFile {
		return false
	}

	if s.SslLoadExtraFiles != t.SslLoadExtraFiles {
		return false
	}

	if s.SslModeAsync != t.SslModeAsync {
		return false
	}

	if s.SslPropquery != t.SslPropquery {
		return false
	}

	if s.SslProvider != t.SslProvider {
		return false
	}

	if s.SslProviderPath != t.SslProviderPath {
		return false
	}

	if s.SslServerVerify != t.SslServerVerify {
		return false
	}

	if s.SslSkipSelfIssuedCa != t.SslSkipSelfIssuedCa {
		return false
	}

	if !equalPointers(s.StatsMaxconn, t.StatsMaxconn) {
		return false
	}

	if !equalPointers(s.StatsTimeout, t.StatsTimeout) {
		return false
	}

	if s.StrictLimits != t.StrictLimits {
		return false
	}

	if s.ThreadGroups != t.ThreadGroups {
		return false
	}

	if s.TuneOptions == nil || t.TuneOptions == nil {
		if s.TuneOptions != nil || t.TuneOptions != nil {
			if opt.NilSameAsEmpty {
				empty := &GlobalTuneOptions{}
				if s.TuneOptions == nil {
					if !(t.TuneOptions.Equal(*empty)) {
						return false
					}
				}
				if t.TuneOptions == nil {
					if !(s.TuneOptions.Equal(*empty)) {
						return false
					}
				}
			} else {
				return false
			}
		}
	} else if !s.TuneOptions.Equal(*t.TuneOptions, opt) {
		return false
	}

	if s.TuneSslDefaultDhParam != t.TuneSslDefaultDhParam {
		return false
	}

	if s.UID != t.UID {
		return false
	}

	if s.Ulimitn != t.Ulimitn {
		return false
	}

	if s.Unsetenv != t.Unsetenv {
		return false
	}

	if s.User != t.User {
		return false
	}

	if s.WurflOptions == nil || t.WurflOptions == nil {
		if s.WurflOptions != nil || t.WurflOptions != nil {
			if opt.NilSameAsEmpty {
				empty := &GlobalWurflOptions{}
				if s.WurflOptions == nil {
					if !(t.WurflOptions.Equal(*empty)) {
						return false
					}
				}
				if t.WurflOptions == nil {
					if !(s.WurflOptions.Equal(*empty)) {
						return false
					}
				}
			} else {
				return false
			}
		}
	} else if !s.WurflOptions.Equal(*t.WurflOptions, opt) {
		return false
	}

	if s.ZeroWarning != t.ZeroWarning {
		return false
	}

	return true
}

// Diff checks if two structs of type Global are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b Global
//	diff := a.Diff(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b Global
//	diff := a.Diff(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s Global) Diff(t Global, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})
	if !CheckSameNilAndLen(s.CPUMaps, t.CPUMaps, opt) {
		diff["CPUMaps"] = []interface{}{s.CPUMaps, t.CPUMaps}
	} else {
		diff2 := make(map[string][]interface{})
		for i := range s.CPUMaps {
			if !s.CPUMaps[i].Equal(*t.CPUMaps[i], opt) {
				diffSub := s.CPUMaps[i].Diff(*t.CPUMaps[i], opt)
				if len(diffSub) > 0 {
					diff2[strconv.Itoa(i)] = []interface{}{diffSub}
				}
			}
		}
		if len(diff2) > 0 {
			diff["CPUMaps"] = []interface{}{diff2}
		}
	}

	if !CheckSameNilAndLen(s.H1CaseAdjusts, t.H1CaseAdjusts, opt) {
		diff["H1CaseAdjusts"] = []interface{}{s.H1CaseAdjusts, t.H1CaseAdjusts}
	} else {
		diff2 := make(map[string][]interface{})
		for i := range s.H1CaseAdjusts {
			if !s.H1CaseAdjusts[i].Equal(*t.H1CaseAdjusts[i], opt) {
				diffSub := s.H1CaseAdjusts[i].Diff(*t.H1CaseAdjusts[i], opt)
				if len(diffSub) > 0 {
					diff2[strconv.Itoa(i)] = []interface{}{diffSub}
				}
			}
		}
		if len(diff2) > 0 {
			diff["H1CaseAdjusts"] = []interface{}{diff2}
		}
	}

	if !CheckSameNilAndLen(s.PresetEnvs, t.PresetEnvs, opt) {
		diff["PresetEnvs"] = []interface{}{s.PresetEnvs, t.PresetEnvs}
	} else {
		diff2 := make(map[string][]interface{})
		for i := range s.PresetEnvs {
			if !s.PresetEnvs[i].Equal(*t.PresetEnvs[i], opt) {
				diffSub := s.PresetEnvs[i].Diff(*t.PresetEnvs[i], opt)
				if len(diffSub) > 0 {
					diff2[strconv.Itoa(i)] = []interface{}{diffSub}
				}
			}
		}
		if len(diff2) > 0 {
			diff["PresetEnvs"] = []interface{}{diff2}
		}
	}

	if !CheckSameNilAndLen(s.RuntimeAPIs, t.RuntimeAPIs, opt) {
		diff["RuntimeAPIs"] = []interface{}{s.RuntimeAPIs, t.RuntimeAPIs}
	} else {
		diff2 := make(map[string][]interface{})
		for i := range s.RuntimeAPIs {
			if !s.RuntimeAPIs[i].Equal(*t.RuntimeAPIs[i], opt) {
				diffSub := s.RuntimeAPIs[i].Diff(*t.RuntimeAPIs[i], opt)
				if len(diffSub) > 0 {
					diff2[strconv.Itoa(i)] = []interface{}{diffSub}
				}
			}
		}
		if len(diff2) > 0 {
			diff["RuntimeAPIs"] = []interface{}{diff2}
		}
	}

	if !CheckSameNilAndLen(s.SetEnvs, t.SetEnvs, opt) {
		diff["SetEnvs"] = []interface{}{s.SetEnvs, t.SetEnvs}
	} else {
		diff2 := make(map[string][]interface{})
		for i := range s.SetEnvs {
			if !s.SetEnvs[i].Equal(*t.SetEnvs[i], opt) {
				diffSub := s.SetEnvs[i].Diff(*t.SetEnvs[i], opt)
				if len(diffSub) > 0 {
					diff2[strconv.Itoa(i)] = []interface{}{diffSub}
				}
			}
		}
		if len(diff2) > 0 {
			diff["SetEnvs"] = []interface{}{diff2}
		}
	}

	if !CheckSameNilAndLen(s.SetVarFmts, t.SetVarFmts, opt) {
		diff["SetVarFmts"] = []interface{}{s.SetVarFmts, t.SetVarFmts}
	} else {
		diff2 := make(map[string][]interface{})
		for i := range s.SetVarFmts {
			if !s.SetVarFmts[i].Equal(*t.SetVarFmts[i], opt) {
				diffSub := s.SetVarFmts[i].Diff(*t.SetVarFmts[i], opt)
				if len(diffSub) > 0 {
					diff2[strconv.Itoa(i)] = []interface{}{diffSub}
				}
			}
		}
		if len(diff2) > 0 {
			diff["SetVarFmts"] = []interface{}{diff2}
		}
	}

	if !CheckSameNilAndLen(s.SetVars, t.SetVars, opt) {
		diff["SetVars"] = []interface{}{s.SetVars, t.SetVars}
	} else {
		diff2 := make(map[string][]interface{})
		for i := range s.SetVars {
			if !s.SetVars[i].Equal(*t.SetVars[i], opt) {
				diffSub := s.SetVars[i].Diff(*t.SetVars[i], opt)
				if len(diffSub) > 0 {
					diff2[strconv.Itoa(i)] = []interface{}{diffSub}
				}
			}
		}
		if len(diff2) > 0 {
			diff["SetVars"] = []interface{}{diff2}
		}
	}

	if !CheckSameNilAndLen(s.SslEngines, t.SslEngines, opt) {
		diff["SslEngines"] = []interface{}{s.SslEngines, t.SslEngines}
	} else {
		diff2 := make(map[string][]interface{})
		for i := range s.SslEngines {
			if !s.SslEngines[i].Equal(*t.SslEngines[i], opt) {
				diffSub := s.SslEngines[i].Diff(*t.SslEngines[i], opt)
				if len(diffSub) > 0 {
					diff2[strconv.Itoa(i)] = []interface{}{diffSub}
				}
			}
		}
		if len(diff2) > 0 {
			diff["SslEngines"] = []interface{}{diff2}
		}
	}

	if !CheckSameNilAndLen(s.ThreadGroupLines, t.ThreadGroupLines, opt) {
		diff["ThreadGroupLines"] = []interface{}{s.ThreadGroupLines, t.ThreadGroupLines}
	} else {
		diff2 := make(map[string][]interface{})
		for i := range s.ThreadGroupLines {
			if !s.ThreadGroupLines[i].Equal(*t.ThreadGroupLines[i], opt) {
				diffSub := s.ThreadGroupLines[i].Diff(*t.ThreadGroupLines[i], opt)
				if len(diffSub) > 0 {
					diff2[strconv.Itoa(i)] = []interface{}{diffSub}
				}
			}
		}
		if len(diff2) > 0 {
			diff["ThreadGroupLines"] = []interface{}{diff2}
		}
	}

	if !equalPointers(s.Anonkey, t.Anonkey) {
		diff["Anonkey"] = []interface{}{ValueOrNil(s.Anonkey), ValueOrNil(t.Anonkey)}
	}

	if s.BusyPolling != t.BusyPolling {
		diff["BusyPolling"] = []interface{}{s.BusyPolling, t.BusyPolling}
	}

	if s.CaBase != t.CaBase {
		diff["CaBase"] = []interface{}{s.CaBase, t.CaBase}
	}

	if s.Chroot != t.Chroot {
		diff["Chroot"] = []interface{}{s.Chroot, t.Chroot}
	}

	if !equalPointers(s.CloseSpreadTime, t.CloseSpreadTime) {
		diff["CloseSpreadTime"] = []interface{}{ValueOrNil(s.CloseSpreadTime), ValueOrNil(t.CloseSpreadTime)}
	}

	if s.ClusterSecret != t.ClusterSecret {
		diff["ClusterSecret"] = []interface{}{s.ClusterSecret, t.ClusterSecret}
	}

	if s.CrtBase != t.CrtBase {
		diff["CrtBase"] = []interface{}{s.CrtBase, t.CrtBase}
	}

	if s.Daemon != t.Daemon {
		diff["Daemon"] = []interface{}{s.Daemon, t.Daemon}
	}

	if s.DefaultPath == nil || t.DefaultPath == nil {
		if s.DefaultPath != nil || t.DefaultPath != nil {
			if opt.NilSameAsEmpty {
				empty := &GlobalDefaultPath{}
				if s.DefaultPath == nil {
					if !(t.DefaultPath.Equal(*empty)) {
						diff["DefaultPath"] = []interface{}{ValueOrNil(s.DefaultPath), ValueOrNil(t.DefaultPath)}
					}
				}
				if t.DefaultPath == nil {
					if !(s.DefaultPath.Equal(*empty)) {
						diff["DefaultPath"] = []interface{}{ValueOrNil(s.DefaultPath), ValueOrNil(t.DefaultPath)}
					}
				}
			} else {
				diff["DefaultPath"] = []interface{}{ValueOrNil(s.DefaultPath), ValueOrNil(t.DefaultPath)}
			}
		}
	} else if !s.DefaultPath.Equal(*t.DefaultPath, opt) {
		diff["DefaultPath"] = []interface{}{ValueOrNil(s.DefaultPath), ValueOrNil(t.DefaultPath)}
	}

	if s.Description != t.Description {
		diff["Description"] = []interface{}{s.Description, t.Description}
	}

	if s.DeviceAtlasOptions == nil || t.DeviceAtlasOptions == nil {
		if s.DeviceAtlasOptions != nil || t.DeviceAtlasOptions != nil {
			if opt.NilSameAsEmpty {
				empty := &GlobalDeviceAtlasOptions{}
				if s.DeviceAtlasOptions == nil {
					if !(t.DeviceAtlasOptions.Equal(*empty)) {
						diff["DeviceAtlasOptions"] = []interface{}{ValueOrNil(s.DeviceAtlasOptions), ValueOrNil(t.DeviceAtlasOptions)}
					}
				}
				if t.DeviceAtlasOptions == nil {
					if !(s.DeviceAtlasOptions.Equal(*empty)) {
						diff["DeviceAtlasOptions"] = []interface{}{ValueOrNil(s.DeviceAtlasOptions), ValueOrNil(t.DeviceAtlasOptions)}
					}
				}
			} else {
				diff["DeviceAtlasOptions"] = []interface{}{ValueOrNil(s.DeviceAtlasOptions), ValueOrNil(t.DeviceAtlasOptions)}
			}
		}
	} else if !s.DeviceAtlasOptions.Equal(*t.DeviceAtlasOptions, opt) {
		diff["DeviceAtlasOptions"] = []interface{}{ValueOrNil(s.DeviceAtlasOptions), ValueOrNil(t.DeviceAtlasOptions)}
	}

	if s.ExposeExperimentalDirectives != t.ExposeExperimentalDirectives {
		diff["ExposeExperimentalDirectives"] = []interface{}{s.ExposeExperimentalDirectives, t.ExposeExperimentalDirectives}
	}

	if s.ExternalCheck != t.ExternalCheck {
		diff["ExternalCheck"] = []interface{}{s.ExternalCheck, t.ExternalCheck}
	}

	if s.FiftyOneDegreesOptions == nil || t.FiftyOneDegreesOptions == nil {
		if s.FiftyOneDegreesOptions != nil || t.FiftyOneDegreesOptions != nil {
			if opt.NilSameAsEmpty {
				empty := &GlobalFiftyOneDegreesOptions{}
				if s.FiftyOneDegreesOptions == nil {
					if !(t.FiftyOneDegreesOptions.Equal(*empty)) {
						diff["FiftyOneDegreesOptions"] = []interface{}{ValueOrNil(s.FiftyOneDegreesOptions), ValueOrNil(t.FiftyOneDegreesOptions)}
					}
				}
				if t.FiftyOneDegreesOptions == nil {
					if !(s.FiftyOneDegreesOptions.Equal(*empty)) {
						diff["FiftyOneDegreesOptions"] = []interface{}{ValueOrNil(s.FiftyOneDegreesOptions), ValueOrNil(t.FiftyOneDegreesOptions)}
					}
				}
			} else {
				diff["FiftyOneDegreesOptions"] = []interface{}{ValueOrNil(s.FiftyOneDegreesOptions), ValueOrNil(t.FiftyOneDegreesOptions)}
			}
		}
	} else if !s.FiftyOneDegreesOptions.Equal(*t.FiftyOneDegreesOptions, opt) {
		diff["FiftyOneDegreesOptions"] = []interface{}{ValueOrNil(s.FiftyOneDegreesOptions), ValueOrNil(t.FiftyOneDegreesOptions)}
	}

	if s.Gid != t.Gid {
		diff["Gid"] = []interface{}{s.Gid, t.Gid}
	}

	if !equalPointers(s.Grace, t.Grace) {
		diff["Grace"] = []interface{}{ValueOrNil(s.Grace), ValueOrNil(t.Grace)}
	}

	if s.Group != t.Group {
		diff["Group"] = []interface{}{s.Group, t.Group}
	}

	if s.H1CaseAdjustFile != t.H1CaseAdjustFile {
		diff["H1CaseAdjustFile"] = []interface{}{s.H1CaseAdjustFile, t.H1CaseAdjustFile}
	}

	if s.H2WorkaroundBogusWebsocketClients != t.H2WorkaroundBogusWebsocketClients {
		diff["H2WorkaroundBogusWebsocketClients"] = []interface{}{s.H2WorkaroundBogusWebsocketClients, t.H2WorkaroundBogusWebsocketClients}
	}

	if !equalPointers(s.HardStopAfter, t.HardStopAfter) {
		diff["HardStopAfter"] = []interface{}{ValueOrNil(s.HardStopAfter), ValueOrNil(t.HardStopAfter)}
	}

	if s.HttpclientResolversDisabled != t.HttpclientResolversDisabled {
		diff["HttpclientResolversDisabled"] = []interface{}{s.HttpclientResolversDisabled, t.HttpclientResolversDisabled}
	}

	if s.HttpclientResolversID != t.HttpclientResolversID {
		diff["HttpclientResolversID"] = []interface{}{s.HttpclientResolversID, t.HttpclientResolversID}
	}

	if s.HttpclientResolversPrefer != t.HttpclientResolversPrefer {
		diff["HttpclientResolversPrefer"] = []interface{}{s.HttpclientResolversPrefer, t.HttpclientResolversPrefer}
	}

	if s.HttpclientRetries != t.HttpclientRetries {
		diff["HttpclientRetries"] = []interface{}{s.HttpclientRetries, t.HttpclientRetries}
	}

	if s.HttpclientSslCaFile != t.HttpclientSslCaFile {
		diff["HttpclientSslCaFile"] = []interface{}{s.HttpclientSslCaFile, t.HttpclientSslCaFile}
	}

	if !equalPointers(s.HttpclientSslVerify, t.HttpclientSslVerify) {
		diff["HttpclientSslVerify"] = []interface{}{ValueOrNil(s.HttpclientSslVerify), ValueOrNil(t.HttpclientSslVerify)}
	}

	if !equalPointers(s.HttpclientTimeoutConnect, t.HttpclientTimeoutConnect) {
		diff["HttpclientTimeoutConnect"] = []interface{}{ValueOrNil(s.HttpclientTimeoutConnect), ValueOrNil(t.HttpclientTimeoutConnect)}
	}

	if s.InsecureForkWanted != t.InsecureForkWanted {
		diff["InsecureForkWanted"] = []interface{}{s.InsecureForkWanted, t.InsecureForkWanted}
	}

	if s.InsecureSetuidWanted != t.InsecureSetuidWanted {
		diff["InsecureSetuidWanted"] = []interface{}{s.InsecureSetuidWanted, t.InsecureSetuidWanted}
	}

	if s.IssuersChainPath != t.IssuersChainPath {
		diff["IssuersChainPath"] = []interface{}{s.IssuersChainPath, t.IssuersChainPath}
	}

	if s.LimitedQuic != t.LimitedQuic {
		diff["LimitedQuic"] = []interface{}{s.LimitedQuic, t.LimitedQuic}
	}

	if s.LoadServerStateFromFile != t.LoadServerStateFromFile {
		diff["LoadServerStateFromFile"] = []interface{}{s.LoadServerStateFromFile, t.LoadServerStateFromFile}
	}

	if s.Localpeer != t.Localpeer {
		diff["Localpeer"] = []interface{}{s.Localpeer, t.Localpeer}
	}

	if s.LogSendHostname == nil || t.LogSendHostname == nil {
		if s.LogSendHostname != nil || t.LogSendHostname != nil {
			if opt.NilSameAsEmpty {
				empty := &GlobalLogSendHostname{}
				if s.LogSendHostname == nil {
					if !(t.LogSendHostname.Equal(*empty)) {
						diff["LogSendHostname"] = []interface{}{ValueOrNil(s.LogSendHostname), ValueOrNil(t.LogSendHostname)}
					}
				}
				if t.LogSendHostname == nil {
					if !(s.LogSendHostname.Equal(*empty)) {
						diff["LogSendHostname"] = []interface{}{ValueOrNil(s.LogSendHostname), ValueOrNil(t.LogSendHostname)}
					}
				}
			} else {
				diff["LogSendHostname"] = []interface{}{ValueOrNil(s.LogSendHostname), ValueOrNil(t.LogSendHostname)}
			}
		}
	} else if !s.LogSendHostname.Equal(*t.LogSendHostname, opt) {
		diff["LogSendHostname"] = []interface{}{ValueOrNil(s.LogSendHostname), ValueOrNil(t.LogSendHostname)}
	}

	if s.LuaLoadPerThread != t.LuaLoadPerThread {
		diff["LuaLoadPerThread"] = []interface{}{s.LuaLoadPerThread, t.LuaLoadPerThread}
	}

	if !CheckSameNilAndLen(s.LuaLoads, t.LuaLoads, opt) {
		diff["LuaLoads"] = []interface{}{s.LuaLoads, t.LuaLoads}
	} else {
		diff2 := make(map[string][]interface{})
		for i := range s.LuaLoads {
			if !s.LuaLoads[i].Equal(*t.LuaLoads[i], opt) {
				diffSub := s.LuaLoads[i].Diff(*t.LuaLoads[i], opt)
				if len(diffSub) > 0 {
					diff2[strconv.Itoa(i)] = []interface{}{diffSub}
				}
			}
		}
		if len(diff2) > 0 {
			diff["LuaLoads"] = []interface{}{diff2}
		}
	}

	if !CheckSameNilAndLen(s.LuaPrependPath, t.LuaPrependPath, opt) {
		diff["LuaPrependPath"] = []interface{}{s.LuaPrependPath, t.LuaPrependPath}
	} else {
		diff2 := make(map[string][]interface{})
		for i := range s.LuaPrependPath {
			if !s.LuaPrependPath[i].Equal(*t.LuaPrependPath[i], opt) {
				diffSub := s.LuaPrependPath[i].Diff(*t.LuaPrependPath[i], opt)
				if len(diffSub) > 0 {
					diff2[strconv.Itoa(i)] = []interface{}{diffSub}
				}
			}
		}
		if len(diff2) > 0 {
			diff["LuaPrependPath"] = []interface{}{diff2}
		}
	}

	if s.MasterWorker != t.MasterWorker {
		diff["MasterWorker"] = []interface{}{s.MasterWorker, t.MasterWorker}
	}

	if !equalPointers(s.MaxSpreadChecks, t.MaxSpreadChecks) {
		diff["MaxSpreadChecks"] = []interface{}{ValueOrNil(s.MaxSpreadChecks), ValueOrNil(t.MaxSpreadChecks)}
	}

	if s.Maxcompcpuusage != t.Maxcompcpuusage {
		diff["Maxcompcpuusage"] = []interface{}{s.Maxcompcpuusage, t.Maxcompcpuusage}
	}

	if s.Maxcomprate != t.Maxcomprate {
		diff["Maxcomprate"] = []interface{}{s.Maxcomprate, t.Maxcomprate}
	}

	if s.Maxconn != t.Maxconn {
		diff["Maxconn"] = []interface{}{s.Maxconn, t.Maxconn}
	}

	if s.Maxconnrate != t.Maxconnrate {
		diff["Maxconnrate"] = []interface{}{s.Maxconnrate, t.Maxconnrate}
	}

	if s.Maxpipes != t.Maxpipes {
		diff["Maxpipes"] = []interface{}{s.Maxpipes, t.Maxpipes}
	}

	if s.Maxsessrate != t.Maxsessrate {
		diff["Maxsessrate"] = []interface{}{s.Maxsessrate, t.Maxsessrate}
	}

	if s.Maxsslconn != t.Maxsslconn {
		diff["Maxsslconn"] = []interface{}{s.Maxsslconn, t.Maxsslconn}
	}

	if s.Maxsslrate != t.Maxsslrate {
		diff["Maxsslrate"] = []interface{}{s.Maxsslrate, t.Maxsslrate}
	}

	if s.Maxzlibmem != t.Maxzlibmem {
		diff["Maxzlibmem"] = []interface{}{s.Maxzlibmem, t.Maxzlibmem}
	}

	if !equalPointers(s.MworkerMaxReloads, t.MworkerMaxReloads) {
		diff["MworkerMaxReloads"] = []interface{}{ValueOrNil(s.MworkerMaxReloads), ValueOrNil(t.MworkerMaxReloads)}
	}

	if s.Nbproc != t.Nbproc {
		diff["Nbproc"] = []interface{}{s.Nbproc, t.Nbproc}
	}

	if s.Nbthread != t.Nbthread {
		diff["Nbthread"] = []interface{}{s.Nbthread, t.Nbthread}
	}

	if s.NoQuic != t.NoQuic {
		diff["NoQuic"] = []interface{}{s.NoQuic, t.NoQuic}
	}

	if s.Node != t.Node {
		diff["Node"] = []interface{}{s.Node, t.Node}
	}

	if s.Noepoll != t.Noepoll {
		diff["Noepoll"] = []interface{}{s.Noepoll, t.Noepoll}
	}

	if s.Noevports != t.Noevports {
		diff["Noevports"] = []interface{}{s.Noevports, t.Noevports}
	}

	if s.Nogetaddrinfo != t.Nogetaddrinfo {
		diff["Nogetaddrinfo"] = []interface{}{s.Nogetaddrinfo, t.Nogetaddrinfo}
	}

	if s.Nokqueue != t.Nokqueue {
		diff["Nokqueue"] = []interface{}{s.Nokqueue, t.Nokqueue}
	}

	if s.Nopoll != t.Nopoll {
		diff["Nopoll"] = []interface{}{s.Nopoll, t.Nopoll}
	}

	if s.Noreuseport != t.Noreuseport {
		diff["Noreuseport"] = []interface{}{s.Noreuseport, t.Noreuseport}
	}

	if s.Nosplice != t.Nosplice {
		diff["Nosplice"] = []interface{}{s.Nosplice, t.Nosplice}
	}

	if s.NumaCPUMapping != t.NumaCPUMapping {
		diff["NumaCPUMapping"] = []interface{}{s.NumaCPUMapping, t.NumaCPUMapping}
	}

	if s.Pidfile != t.Pidfile {
		diff["Pidfile"] = []interface{}{s.Pidfile, t.Pidfile}
	}

	if s.Pp2NeverSendLocal != t.Pp2NeverSendLocal {
		diff["Pp2NeverSendLocal"] = []interface{}{s.Pp2NeverSendLocal, t.Pp2NeverSendLocal}
	}

	if s.PreallocFd != t.PreallocFd {
		diff["PreallocFd"] = []interface{}{s.PreallocFd, t.PreallocFd}
	}

	if s.ProfilingTasks != t.ProfilingTasks {
		diff["ProfilingTasks"] = []interface{}{s.ProfilingTasks, t.ProfilingTasks}
	}

	if s.Quiet != t.Quiet {
		diff["Quiet"] = []interface{}{s.Quiet, t.Quiet}
	}

	if s.Resetenv != t.Resetenv {
		diff["Resetenv"] = []interface{}{s.Resetenv, t.Resetenv}
	}

	if s.ServerStateBase != t.ServerStateBase {
		diff["ServerStateBase"] = []interface{}{s.ServerStateBase, t.ServerStateBase}
	}

	if s.ServerStateFile != t.ServerStateFile {
		diff["ServerStateFile"] = []interface{}{s.ServerStateFile, t.ServerStateFile}
	}

	if s.SetDumpable != t.SetDumpable {
		diff["SetDumpable"] = []interface{}{s.SetDumpable, t.SetDumpable}
	}

	if s.Setcap != t.Setcap {
		diff["Setcap"] = []interface{}{s.Setcap, t.Setcap}
	}

	if s.SpreadChecks != t.SpreadChecks {
		diff["SpreadChecks"] = []interface{}{s.SpreadChecks, t.SpreadChecks}
	}

	if s.SslDefaultBindCiphers != t.SslDefaultBindCiphers {
		diff["SslDefaultBindCiphers"] = []interface{}{s.SslDefaultBindCiphers, t.SslDefaultBindCiphers}
	}

	if s.SslDefaultBindCiphersuites != t.SslDefaultBindCiphersuites {
		diff["SslDefaultBindCiphersuites"] = []interface{}{s.SslDefaultBindCiphersuites, t.SslDefaultBindCiphersuites}
	}

	if s.SslDefaultBindClientSigalgs != t.SslDefaultBindClientSigalgs {
		diff["SslDefaultBindClientSigalgs"] = []interface{}{s.SslDefaultBindClientSigalgs, t.SslDefaultBindClientSigalgs}
	}

	if s.SslDefaultBindCurves != t.SslDefaultBindCurves {
		diff["SslDefaultBindCurves"] = []interface{}{s.SslDefaultBindCurves, t.SslDefaultBindCurves}
	}

	if s.SslDefaultBindOptions != t.SslDefaultBindOptions {
		diff["SslDefaultBindOptions"] = []interface{}{s.SslDefaultBindOptions, t.SslDefaultBindOptions}
	}

	if s.SslDefaultBindSigalgs != t.SslDefaultBindSigalgs {
		diff["SslDefaultBindSigalgs"] = []interface{}{s.SslDefaultBindSigalgs, t.SslDefaultBindSigalgs}
	}

	if s.SslDefaultServerCiphers != t.SslDefaultServerCiphers {
		diff["SslDefaultServerCiphers"] = []interface{}{s.SslDefaultServerCiphers, t.SslDefaultServerCiphers}
	}

	if s.SslDefaultServerCiphersuites != t.SslDefaultServerCiphersuites {
		diff["SslDefaultServerCiphersuites"] = []interface{}{s.SslDefaultServerCiphersuites, t.SslDefaultServerCiphersuites}
	}

	if s.SslDefaultServerClientSigalgs != t.SslDefaultServerClientSigalgs {
		diff["SslDefaultServerClientSigalgs"] = []interface{}{s.SslDefaultServerClientSigalgs, t.SslDefaultServerClientSigalgs}
	}

	if s.SslDefaultServerCurves != t.SslDefaultServerCurves {
		diff["SslDefaultServerCurves"] = []interface{}{s.SslDefaultServerCurves, t.SslDefaultServerCurves}
	}

	if s.SslDefaultServerOptions != t.SslDefaultServerOptions {
		diff["SslDefaultServerOptions"] = []interface{}{s.SslDefaultServerOptions, t.SslDefaultServerOptions}
	}

	if s.SslDefaultServerSigalgs != t.SslDefaultServerSigalgs {
		diff["SslDefaultServerSigalgs"] = []interface{}{s.SslDefaultServerSigalgs, t.SslDefaultServerSigalgs}
	}

	if s.SslDhParamFile != t.SslDhParamFile {
		diff["SslDhParamFile"] = []interface{}{s.SslDhParamFile, t.SslDhParamFile}
	}

	if s.SslLoadExtraFiles != t.SslLoadExtraFiles {
		diff["SslLoadExtraFiles"] = []interface{}{s.SslLoadExtraFiles, t.SslLoadExtraFiles}
	}

	if s.SslModeAsync != t.SslModeAsync {
		diff["SslModeAsync"] = []interface{}{s.SslModeAsync, t.SslModeAsync}
	}

	if s.SslPropquery != t.SslPropquery {
		diff["SslPropquery"] = []interface{}{s.SslPropquery, t.SslPropquery}
	}

	if s.SslProvider != t.SslProvider {
		diff["SslProvider"] = []interface{}{s.SslProvider, t.SslProvider}
	}

	if s.SslProviderPath != t.SslProviderPath {
		diff["SslProviderPath"] = []interface{}{s.SslProviderPath, t.SslProviderPath}
	}

	if s.SslServerVerify != t.SslServerVerify {
		diff["SslServerVerify"] = []interface{}{s.SslServerVerify, t.SslServerVerify}
	}

	if s.SslSkipSelfIssuedCa != t.SslSkipSelfIssuedCa {
		diff["SslSkipSelfIssuedCa"] = []interface{}{s.SslSkipSelfIssuedCa, t.SslSkipSelfIssuedCa}
	}

	if !equalPointers(s.StatsMaxconn, t.StatsMaxconn) {
		diff["StatsMaxconn"] = []interface{}{ValueOrNil(s.StatsMaxconn), ValueOrNil(t.StatsMaxconn)}
	}

	if !equalPointers(s.StatsTimeout, t.StatsTimeout) {
		diff["StatsTimeout"] = []interface{}{ValueOrNil(s.StatsTimeout), ValueOrNil(t.StatsTimeout)}
	}

	if s.StrictLimits != t.StrictLimits {
		diff["StrictLimits"] = []interface{}{s.StrictLimits, t.StrictLimits}
	}

	if s.ThreadGroups != t.ThreadGroups {
		diff["ThreadGroups"] = []interface{}{s.ThreadGroups, t.ThreadGroups}
	}

	if s.TuneOptions == nil || t.TuneOptions == nil {
		if s.TuneOptions != nil || t.TuneOptions != nil {
			if opt.NilSameAsEmpty {
				empty := &GlobalTuneOptions{}
				if s.TuneOptions == nil {
					if !(t.TuneOptions.Equal(*empty)) {
						diff["TuneOptions"] = []interface{}{ValueOrNil(s.TuneOptions), ValueOrNil(t.TuneOptions)}
					}
				}
				if t.TuneOptions == nil {
					if !(s.TuneOptions.Equal(*empty)) {
						diff["TuneOptions"] = []interface{}{ValueOrNil(s.TuneOptions), ValueOrNil(t.TuneOptions)}
					}
				}
			} else {
				diff["TuneOptions"] = []interface{}{ValueOrNil(s.TuneOptions), ValueOrNil(t.TuneOptions)}
			}
		}
	} else if !s.TuneOptions.Equal(*t.TuneOptions, opt) {
		diff["TuneOptions"] = []interface{}{ValueOrNil(s.TuneOptions), ValueOrNil(t.TuneOptions)}
	}

	if s.TuneSslDefaultDhParam != t.TuneSslDefaultDhParam {
		diff["TuneSslDefaultDhParam"] = []interface{}{s.TuneSslDefaultDhParam, t.TuneSslDefaultDhParam}
	}

	if s.UID != t.UID {
		diff["UID"] = []interface{}{s.UID, t.UID}
	}

	if s.Ulimitn != t.Ulimitn {
		diff["Ulimitn"] = []interface{}{s.Ulimitn, t.Ulimitn}
	}

	if s.Unsetenv != t.Unsetenv {
		diff["Unsetenv"] = []interface{}{s.Unsetenv, t.Unsetenv}
	}

	if s.User != t.User {
		diff["User"] = []interface{}{s.User, t.User}
	}

	if s.WurflOptions == nil || t.WurflOptions == nil {
		if s.WurflOptions != nil || t.WurflOptions != nil {
			if opt.NilSameAsEmpty {
				empty := &GlobalWurflOptions{}
				if s.WurflOptions == nil {
					if !(t.WurflOptions.Equal(*empty)) {
						diff["WurflOptions"] = []interface{}{ValueOrNil(s.WurflOptions), ValueOrNil(t.WurflOptions)}
					}
				}
				if t.WurflOptions == nil {
					if !(s.WurflOptions.Equal(*empty)) {
						diff["WurflOptions"] = []interface{}{ValueOrNil(s.WurflOptions), ValueOrNil(t.WurflOptions)}
					}
				}
			} else {
				diff["WurflOptions"] = []interface{}{ValueOrNil(s.WurflOptions), ValueOrNil(t.WurflOptions)}
			}
		}
	} else if !s.WurflOptions.Equal(*t.WurflOptions, opt) {
		diff["WurflOptions"] = []interface{}{ValueOrNil(s.WurflOptions), ValueOrNil(t.WurflOptions)}
	}

	if s.ZeroWarning != t.ZeroWarning {
		diff["ZeroWarning"] = []interface{}{s.ZeroWarning, t.ZeroWarning}
	}

	return diff
}

// Equal checks if two structs of type CPUMap are equal
//
//	var a, b CPUMap
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s CPUMap) Equal(t CPUMap, opts ...Options) bool {
	if !equalPointers(s.CPUSet, t.CPUSet) {
		return false
	}

	if !equalPointers(s.Process, t.Process) {
		return false
	}

	return true
}

// Diff checks if two structs of type CPUMap are equal
//
//	var a, b CPUMap
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s CPUMap) Diff(t CPUMap, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if !equalPointers(s.CPUSet, t.CPUSet) {
		diff["CPUSet"] = []interface{}{ValueOrNil(s.CPUSet), ValueOrNil(t.CPUSet)}
	}

	if !equalPointers(s.Process, t.Process) {
		diff["Process"] = []interface{}{ValueOrNil(s.Process), ValueOrNil(t.Process)}
	}

	return diff
}

// Equal checks if two structs of type GlobalDefaultPath are equal
//
//	var a, b GlobalDefaultPath
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s GlobalDefaultPath) Equal(t GlobalDefaultPath, opts ...Options) bool {
	if s.Path != t.Path {
		return false
	}

	if s.Type != t.Type {
		return false
	}

	return true
}

// Diff checks if two structs of type GlobalDefaultPath are equal
//
//	var a, b GlobalDefaultPath
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s GlobalDefaultPath) Diff(t GlobalDefaultPath, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if s.Path != t.Path {
		diff["Path"] = []interface{}{s.Path, t.Path}
	}

	if s.Type != t.Type {
		diff["Type"] = []interface{}{s.Type, t.Type}
	}

	return diff
}

// Equal checks if two structs of type GlobalDeviceAtlasOptions are equal
//
//	var a, b GlobalDeviceAtlasOptions
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s GlobalDeviceAtlasOptions) Equal(t GlobalDeviceAtlasOptions, opts ...Options) bool {
	if s.JSONFile != t.JSONFile {
		return false
	}

	if s.LogLevel != t.LogLevel {
		return false
	}

	if s.PropertiesCookie != t.PropertiesCookie {
		return false
	}

	if s.Separator != t.Separator {
		return false
	}

	return true
}

// Diff checks if two structs of type GlobalDeviceAtlasOptions are equal
//
//	var a, b GlobalDeviceAtlasOptions
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s GlobalDeviceAtlasOptions) Diff(t GlobalDeviceAtlasOptions, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if s.JSONFile != t.JSONFile {
		diff["JSONFile"] = []interface{}{s.JSONFile, t.JSONFile}
	}

	if s.LogLevel != t.LogLevel {
		diff["LogLevel"] = []interface{}{s.LogLevel, t.LogLevel}
	}

	if s.PropertiesCookie != t.PropertiesCookie {
		diff["PropertiesCookie"] = []interface{}{s.PropertiesCookie, t.PropertiesCookie}
	}

	if s.Separator != t.Separator {
		diff["Separator"] = []interface{}{s.Separator, t.Separator}
	}

	return diff
}

// Equal checks if two structs of type GlobalFiftyOneDegreesOptions are equal
//
//	var a, b GlobalFiftyOneDegreesOptions
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s GlobalFiftyOneDegreesOptions) Equal(t GlobalFiftyOneDegreesOptions, opts ...Options) bool {
	if s.CacheSize != t.CacheSize {
		return false
	}

	if s.DataFile != t.DataFile {
		return false
	}

	if s.PropertyNameList != t.PropertyNameList {
		return false
	}

	if s.PropertySeparator != t.PropertySeparator {
		return false
	}

	return true
}

// Diff checks if two structs of type GlobalFiftyOneDegreesOptions are equal
//
//	var a, b GlobalFiftyOneDegreesOptions
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s GlobalFiftyOneDegreesOptions) Diff(t GlobalFiftyOneDegreesOptions, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if s.CacheSize != t.CacheSize {
		diff["CacheSize"] = []interface{}{s.CacheSize, t.CacheSize}
	}

	if s.DataFile != t.DataFile {
		diff["DataFile"] = []interface{}{s.DataFile, t.DataFile}
	}

	if s.PropertyNameList != t.PropertyNameList {
		diff["PropertyNameList"] = []interface{}{s.PropertyNameList, t.PropertyNameList}
	}

	if s.PropertySeparator != t.PropertySeparator {
		diff["PropertySeparator"] = []interface{}{s.PropertySeparator, t.PropertySeparator}
	}

	return diff
}

// Equal checks if two structs of type H1CaseAdjust are equal
//
//	var a, b H1CaseAdjust
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s H1CaseAdjust) Equal(t H1CaseAdjust, opts ...Options) bool {
	if !equalPointers(s.From, t.From) {
		return false
	}

	if !equalPointers(s.To, t.To) {
		return false
	}

	return true
}

// Diff checks if two structs of type H1CaseAdjust are equal
//
//	var a, b H1CaseAdjust
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s H1CaseAdjust) Diff(t H1CaseAdjust, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if !equalPointers(s.From, t.From) {
		diff["From"] = []interface{}{ValueOrNil(s.From), ValueOrNil(t.From)}
	}

	if !equalPointers(s.To, t.To) {
		diff["To"] = []interface{}{ValueOrNil(s.To), ValueOrNil(t.To)}
	}

	return diff
}

// Equal checks if two structs of type GlobalLogSendHostname are equal
//
//	var a, b GlobalLogSendHostname
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s GlobalLogSendHostname) Equal(t GlobalLogSendHostname, opts ...Options) bool {
	if !equalPointers(s.Enabled, t.Enabled) {
		return false
	}

	if s.Param != t.Param {
		return false
	}

	return true
}

// Diff checks if two structs of type GlobalLogSendHostname are equal
//
//	var a, b GlobalLogSendHostname
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s GlobalLogSendHostname) Diff(t GlobalLogSendHostname, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if !equalPointers(s.Enabled, t.Enabled) {
		diff["Enabled"] = []interface{}{ValueOrNil(s.Enabled), ValueOrNil(t.Enabled)}
	}

	if s.Param != t.Param {
		diff["Param"] = []interface{}{s.Param, t.Param}
	}

	return diff
}

// Equal checks if two structs of type LuaLoad are equal
//
//	var a, b LuaLoad
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s LuaLoad) Equal(t LuaLoad, opts ...Options) bool {
	if !equalPointers(s.File, t.File) {
		return false
	}

	return true
}

// Diff checks if two structs of type LuaLoad are equal
//
//	var a, b LuaLoad
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s LuaLoad) Diff(t LuaLoad, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if !equalPointers(s.File, t.File) {
		diff["File"] = []interface{}{ValueOrNil(s.File), ValueOrNil(t.File)}
	}

	return diff
}

// Equal checks if two structs of type LuaPrependPath are equal
//
//	var a, b LuaPrependPath
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s LuaPrependPath) Equal(t LuaPrependPath, opts ...Options) bool {
	if !equalPointers(s.Path, t.Path) {
		return false
	}

	if s.Type != t.Type {
		return false
	}

	return true
}

// Diff checks if two structs of type LuaPrependPath are equal
//
//	var a, b LuaPrependPath
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s LuaPrependPath) Diff(t LuaPrependPath, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if !equalPointers(s.Path, t.Path) {
		diff["Path"] = []interface{}{ValueOrNil(s.Path), ValueOrNil(t.Path)}
	}

	if s.Type != t.Type {
		diff["Type"] = []interface{}{s.Type, t.Type}
	}

	return diff
}

// Equal checks if two structs of type PresetEnv are equal
//
//	var a, b PresetEnv
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s PresetEnv) Equal(t PresetEnv, opts ...Options) bool {
	if !equalPointers(s.Name, t.Name) {
		return false
	}

	if !equalPointers(s.Value, t.Value) {
		return false
	}

	return true
}

// Diff checks if two structs of type PresetEnv are equal
//
//	var a, b PresetEnv
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s PresetEnv) Diff(t PresetEnv, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if !equalPointers(s.Name, t.Name) {
		diff["Name"] = []interface{}{ValueOrNil(s.Name), ValueOrNil(t.Name)}
	}

	if !equalPointers(s.Value, t.Value) {
		diff["Value"] = []interface{}{ValueOrNil(s.Value), ValueOrNil(t.Value)}
	}

	return diff
}

// Equal checks if two structs of type RuntimeAPI are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b RuntimeAPI
//	equal := a.Equal(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b RuntimeAPI
//	equal := a.Equal(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s RuntimeAPI) Equal(t RuntimeAPI, opts ...Options) bool {
	opt := getOptions(opts...)

	if !s.BindParams.Equal(t.BindParams, opt) {
		return false
	}

	if !equalPointers(s.Address, t.Address) {
		return false
	}

	return true
}

// Diff checks if two structs of type RuntimeAPI are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b RuntimeAPI
//	diff := a.Diff(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b RuntimeAPI
//	diff := a.Diff(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s RuntimeAPI) Diff(t RuntimeAPI, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})

	if !s.BindParams.Equal(t.BindParams, opt) {
		diff["BindParams"] = []interface{}{s.BindParams, t.BindParams}
	}

	if !equalPointers(s.Address, t.Address) {
		diff["Address"] = []interface{}{ValueOrNil(s.Address), ValueOrNil(t.Address)}
	}

	return diff
}

// Equal checks if two structs of type SetVarFmt are equal
//
//	var a, b SetVarFmt
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s SetVarFmt) Equal(t SetVarFmt, opts ...Options) bool {
	if !equalPointers(s.Format, t.Format) {
		return false
	}

	if !equalPointers(s.Name, t.Name) {
		return false
	}

	return true
}

// Diff checks if two structs of type SetVarFmt are equal
//
//	var a, b SetVarFmt
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s SetVarFmt) Diff(t SetVarFmt, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if !equalPointers(s.Format, t.Format) {
		diff["Format"] = []interface{}{ValueOrNil(s.Format), ValueOrNil(t.Format)}
	}

	if !equalPointers(s.Name, t.Name) {
		diff["Name"] = []interface{}{ValueOrNil(s.Name), ValueOrNil(t.Name)}
	}

	return diff
}

// Equal checks if two structs of type SetVar are equal
//
//	var a, b SetVar
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s SetVar) Equal(t SetVar, opts ...Options) bool {
	if !equalPointers(s.Expr, t.Expr) {
		return false
	}

	if !equalPointers(s.Name, t.Name) {
		return false
	}

	return true
}

// Diff checks if two structs of type SetVar are equal
//
//	var a, b SetVar
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s SetVar) Diff(t SetVar, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if !equalPointers(s.Expr, t.Expr) {
		diff["Expr"] = []interface{}{ValueOrNil(s.Expr), ValueOrNil(t.Expr)}
	}

	if !equalPointers(s.Name, t.Name) {
		diff["Name"] = []interface{}{ValueOrNil(s.Name), ValueOrNil(t.Name)}
	}

	return diff
}

// Equal checks if two structs of type SetEnv are equal
//
//	var a, b SetEnv
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s SetEnv) Equal(t SetEnv, opts ...Options) bool {
	if !equalPointers(s.Name, t.Name) {
		return false
	}

	if !equalPointers(s.Value, t.Value) {
		return false
	}

	return true
}

// Diff checks if two structs of type SetEnv are equal
//
//	var a, b SetEnv
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s SetEnv) Diff(t SetEnv, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if !equalPointers(s.Name, t.Name) {
		diff["Name"] = []interface{}{ValueOrNil(s.Name), ValueOrNil(t.Name)}
	}

	if !equalPointers(s.Value, t.Value) {
		diff["Value"] = []interface{}{ValueOrNil(s.Value), ValueOrNil(t.Value)}
	}

	return diff
}

// Equal checks if two structs of type SslEngine are equal
//
//	var a, b SslEngine
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s SslEngine) Equal(t SslEngine, opts ...Options) bool {
	if !equalPointers(s.Algorithms, t.Algorithms) {
		return false
	}

	if !equalPointers(s.Name, t.Name) {
		return false
	}

	return true
}

// Diff checks if two structs of type SslEngine are equal
//
//	var a, b SslEngine
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s SslEngine) Diff(t SslEngine, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if !equalPointers(s.Algorithms, t.Algorithms) {
		diff["Algorithms"] = []interface{}{ValueOrNil(s.Algorithms), ValueOrNil(t.Algorithms)}
	}

	if !equalPointers(s.Name, t.Name) {
		diff["Name"] = []interface{}{ValueOrNil(s.Name), ValueOrNil(t.Name)}
	}

	return diff
}

// Equal checks if two structs of type ThreadGroup are equal
//
//	var a, b ThreadGroup
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s ThreadGroup) Equal(t ThreadGroup, opts ...Options) bool {
	if !equalPointers(s.Group, t.Group) {
		return false
	}

	if !equalPointers(s.NumOrRange, t.NumOrRange) {
		return false
	}

	return true
}

// Diff checks if two structs of type ThreadGroup are equal
//
//	var a, b ThreadGroup
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s ThreadGroup) Diff(t ThreadGroup, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if !equalPointers(s.Group, t.Group) {
		diff["Group"] = []interface{}{ValueOrNil(s.Group), ValueOrNil(t.Group)}
	}

	if !equalPointers(s.NumOrRange, t.NumOrRange) {
		diff["NumOrRange"] = []interface{}{ValueOrNil(s.NumOrRange), ValueOrNil(t.NumOrRange)}
	}

	return diff
}

// Equal checks if two structs of type GlobalTuneOptions are equal
//
//	var a, b GlobalTuneOptions
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s GlobalTuneOptions) Equal(t GlobalTuneOptions, opts ...Options) bool {
	if !equalPointers(s.BuffersLimit, t.BuffersLimit) {
		return false
	}

	if s.BuffersReserve != t.BuffersReserve {
		return false
	}

	if s.Bufsize != t.Bufsize {
		return false
	}

	if s.CompMaxlevel != t.CompMaxlevel {
		return false
	}

	if s.DisableZeroCopyForwarding != t.DisableZeroCopyForwarding {
		return false
	}

	if s.EventsMaxEventsAtOnce != t.EventsMaxEventsAtOnce {
		return false
	}

	if s.FailAlloc != t.FailAlloc {
		return false
	}

	if s.FdEdgeTriggered != t.FdEdgeTriggered {
		return false
	}

	if s.H1ZeroCopyFwdRecv != t.H1ZeroCopyFwdRecv {
		return false
	}

	if s.H2BeInitialWindowSize != t.H2BeInitialWindowSize {
		return false
	}

	if s.H2BeMaxConcurrentStreams != t.H2BeMaxConcurrentStreams {
		return false
	}

	if s.H2FeInitialWindowSize != t.H2FeInitialWindowSize {
		return false
	}

	if s.H2FeMaxConcurrentStreams != t.H2FeMaxConcurrentStreams {
		return false
	}

	if s.H2HeaderTableSize != t.H2HeaderTableSize {
		return false
	}

	if !equalPointers(s.H2InitialWindowSize, t.H2InitialWindowSize) {
		return false
	}

	if s.H2MaxConcurrentStreams != t.H2MaxConcurrentStreams {
		return false
	}

	if s.H2MaxFrameSize != t.H2MaxFrameSize {
		return false
	}

	if s.HTTPCookielen != t.HTTPCookielen {
		return false
	}

	if s.HTTPLogurilen != t.HTTPLogurilen {
		return false
	}

	if s.HTTPMaxhdr != t.HTTPMaxhdr {
		return false
	}

	if s.IdlePoolShared != t.IdlePoolShared {
		return false
	}

	if !equalPointers(s.Idletimer, t.Idletimer) {
		return false
	}

	if s.ListenerDefaultShards != t.ListenerDefaultShards {
		return false
	}

	if s.ListenerMultiQueue != t.ListenerMultiQueue {
		return false
	}

	if !equalPointers(s.LuaBurstTimeout, t.LuaBurstTimeout) {
		return false
	}

	if s.LuaForcedYield != t.LuaForcedYield {
		return false
	}

	if s.LuaLogLoggers != t.LuaLogLoggers {
		return false
	}

	if s.LuaLogStderr != t.LuaLogStderr {
		return false
	}

	if s.LuaMaxmem != t.LuaMaxmem {
		return false
	}

	if !equalPointers(s.LuaServiceTimeout, t.LuaServiceTimeout) {
		return false
	}

	if !equalPointers(s.LuaSessionTimeout, t.LuaSessionTimeout) {
		return false
	}

	if !equalPointers(s.LuaTaskTimeout, t.LuaTaskTimeout) {
		return false
	}

	if !equalPointers(s.MaxChecksPerThread, t.MaxChecksPerThread) {
		return false
	}

	if s.Maxaccept != t.Maxaccept {
		return false
	}

	if s.Maxpollevents != t.Maxpollevents {
		return false
	}

	if s.Maxrewrite != t.Maxrewrite {
		return false
	}

	if !equalPointers(s.MemoryHotSize, t.MemoryHotSize) {
		return false
	}

	if !equalPointers(s.PatternCacheSize, t.PatternCacheSize) {
		return false
	}

	if s.PeersMaxUpdatesAtOnce != t.PeersMaxUpdatesAtOnce {
		return false
	}

	if s.Pipesize != t.Pipesize {
		return false
	}

	if s.PoolHighFdRatio != t.PoolHighFdRatio {
		return false
	}

	if s.PoolLowFdRatio != t.PoolLowFdRatio {
		return false
	}

	if !equalPointers(s.QuicFrontendConnTxBuffersLimit, t.QuicFrontendConnTxBuffersLimit) {
		return false
	}

	if !equalPointers(s.QuicFrontendMaxIdleTimeout, t.QuicFrontendMaxIdleTimeout) {
		return false
	}

	if !equalPointers(s.QuicFrontendMaxStreamsBidi, t.QuicFrontendMaxStreamsBidi) {
		return false
	}

	if !equalPointers(s.QuicMaxFrameLoss, t.QuicMaxFrameLoss) {
		return false
	}

	if !equalPointers(s.QuicRetryThreshold, t.QuicRetryThreshold) {
		return false
	}

	if s.QuicSocketOwner != t.QuicSocketOwner {
		return false
	}

	if !equalPointers(s.RcvbufBackend, t.RcvbufBackend) {
		return false
	}

	if !equalPointers(s.RcvbufClient, t.RcvbufClient) {
		return false
	}

	if !equalPointers(s.RcvbufFrontend, t.RcvbufFrontend) {
		return false
	}

	if !equalPointers(s.RcvbufServer, t.RcvbufServer) {
		return false
	}

	if s.RecvEnough != t.RecvEnough {
		return false
	}

	if s.RunqueueDepth != t.RunqueueDepth {
		return false
	}

	if s.SchedLowLatency != t.SchedLowLatency {
		return false
	}

	if !equalPointers(s.SndbufBackend, t.SndbufBackend) {
		return false
	}

	if !equalPointers(s.SndbufClient, t.SndbufClient) {
		return false
	}

	if !equalPointers(s.SndbufFrontend, t.SndbufFrontend) {
		return false
	}

	if !equalPointers(s.SndbufServer, t.SndbufServer) {
		return false
	}

	if !equalPointers(s.SslCachesize, t.SslCachesize) {
		return false
	}

	if !equalPointers(s.SslCaptureBufferSize, t.SslCaptureBufferSize) {
		return false
	}

	if s.SslCtxCacheSize != t.SslCtxCacheSize {
		return false
	}

	if s.SslDefaultDhParam != t.SslDefaultDhParam {
		return false
	}

	if s.SslForcePrivateCache != t.SslForcePrivateCache {
		return false
	}

	if s.SslKeylog != t.SslKeylog {
		return false
	}

	if !equalPointers(s.SslLifetime, t.SslLifetime) {
		return false
	}

	if !equalPointers(s.SslMaxrecord, t.SslMaxrecord) {
		return false
	}

	if !equalPointers(s.SslOcspUpdateMaxDelay, t.SslOcspUpdateMaxDelay) {
		return false
	}

	if !equalPointers(s.SslOcspUpdateMinDelay, t.SslOcspUpdateMinDelay) {
		return false
	}

	if !equalPointers(s.StickCounters, t.StickCounters) {
		return false
	}

	if !equalPointers(s.VarsGlobalMaxSize, t.VarsGlobalMaxSize) {
		return false
	}

	if !equalPointers(s.VarsProcMaxSize, t.VarsProcMaxSize) {
		return false
	}

	if !equalPointers(s.VarsReqresMaxSize, t.VarsReqresMaxSize) {
		return false
	}

	if !equalPointers(s.VarsSessMaxSize, t.VarsSessMaxSize) {
		return false
	}

	if !equalPointers(s.VarsTxnMaxSize, t.VarsTxnMaxSize) {
		return false
	}

	if s.ZlibMemlevel != t.ZlibMemlevel {
		return false
	}

	if s.ZlibWindowsize != t.ZlibWindowsize {
		return false
	}

	return true
}

// Diff checks if two structs of type GlobalTuneOptions are equal
//
//	var a, b GlobalTuneOptions
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s GlobalTuneOptions) Diff(t GlobalTuneOptions, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if !equalPointers(s.BuffersLimit, t.BuffersLimit) {
		diff["BuffersLimit"] = []interface{}{ValueOrNil(s.BuffersLimit), ValueOrNil(t.BuffersLimit)}
	}

	if s.BuffersReserve != t.BuffersReserve {
		diff["BuffersReserve"] = []interface{}{s.BuffersReserve, t.BuffersReserve}
	}

	if s.Bufsize != t.Bufsize {
		diff["Bufsize"] = []interface{}{s.Bufsize, t.Bufsize}
	}

	if s.CompMaxlevel != t.CompMaxlevel {
		diff["CompMaxlevel"] = []interface{}{s.CompMaxlevel, t.CompMaxlevel}
	}

	if s.DisableZeroCopyForwarding != t.DisableZeroCopyForwarding {
		diff["DisableZeroCopyForwarding"] = []interface{}{s.DisableZeroCopyForwarding, t.DisableZeroCopyForwarding}
	}

	if s.EventsMaxEventsAtOnce != t.EventsMaxEventsAtOnce {
		diff["EventsMaxEventsAtOnce"] = []interface{}{s.EventsMaxEventsAtOnce, t.EventsMaxEventsAtOnce}
	}

	if s.FailAlloc != t.FailAlloc {
		diff["FailAlloc"] = []interface{}{s.FailAlloc, t.FailAlloc}
	}

	if s.FdEdgeTriggered != t.FdEdgeTriggered {
		diff["FdEdgeTriggered"] = []interface{}{s.FdEdgeTriggered, t.FdEdgeTriggered}
	}

	if s.H1ZeroCopyFwdRecv != t.H1ZeroCopyFwdRecv {
		diff["H1ZeroCopyFwdRecv"] = []interface{}{s.H1ZeroCopyFwdRecv, t.H1ZeroCopyFwdRecv}
	}

	if s.H2BeInitialWindowSize != t.H2BeInitialWindowSize {
		diff["H2BeInitialWindowSize"] = []interface{}{s.H2BeInitialWindowSize, t.H2BeInitialWindowSize}
	}

	if s.H2BeMaxConcurrentStreams != t.H2BeMaxConcurrentStreams {
		diff["H2BeMaxConcurrentStreams"] = []interface{}{s.H2BeMaxConcurrentStreams, t.H2BeMaxConcurrentStreams}
	}

	if s.H2FeInitialWindowSize != t.H2FeInitialWindowSize {
		diff["H2FeInitialWindowSize"] = []interface{}{s.H2FeInitialWindowSize, t.H2FeInitialWindowSize}
	}

	if s.H2FeMaxConcurrentStreams != t.H2FeMaxConcurrentStreams {
		diff["H2FeMaxConcurrentStreams"] = []interface{}{s.H2FeMaxConcurrentStreams, t.H2FeMaxConcurrentStreams}
	}

	if s.H2HeaderTableSize != t.H2HeaderTableSize {
		diff["H2HeaderTableSize"] = []interface{}{s.H2HeaderTableSize, t.H2HeaderTableSize}
	}

	if !equalPointers(s.H2InitialWindowSize, t.H2InitialWindowSize) {
		diff["H2InitialWindowSize"] = []interface{}{ValueOrNil(s.H2InitialWindowSize), ValueOrNil(t.H2InitialWindowSize)}
	}

	if s.H2MaxConcurrentStreams != t.H2MaxConcurrentStreams {
		diff["H2MaxConcurrentStreams"] = []interface{}{s.H2MaxConcurrentStreams, t.H2MaxConcurrentStreams}
	}

	if s.H2MaxFrameSize != t.H2MaxFrameSize {
		diff["H2MaxFrameSize"] = []interface{}{s.H2MaxFrameSize, t.H2MaxFrameSize}
	}

	if s.HTTPCookielen != t.HTTPCookielen {
		diff["HTTPCookielen"] = []interface{}{s.HTTPCookielen, t.HTTPCookielen}
	}

	if s.HTTPLogurilen != t.HTTPLogurilen {
		diff["HTTPLogurilen"] = []interface{}{s.HTTPLogurilen, t.HTTPLogurilen}
	}

	if s.HTTPMaxhdr != t.HTTPMaxhdr {
		diff["HTTPMaxhdr"] = []interface{}{s.HTTPMaxhdr, t.HTTPMaxhdr}
	}

	if s.IdlePoolShared != t.IdlePoolShared {
		diff["IdlePoolShared"] = []interface{}{s.IdlePoolShared, t.IdlePoolShared}
	}

	if !equalPointers(s.Idletimer, t.Idletimer) {
		diff["Idletimer"] = []interface{}{ValueOrNil(s.Idletimer), ValueOrNil(t.Idletimer)}
	}

	if s.ListenerDefaultShards != t.ListenerDefaultShards {
		diff["ListenerDefaultShards"] = []interface{}{s.ListenerDefaultShards, t.ListenerDefaultShards}
	}

	if s.ListenerMultiQueue != t.ListenerMultiQueue {
		diff["ListenerMultiQueue"] = []interface{}{s.ListenerMultiQueue, t.ListenerMultiQueue}
	}

	if !equalPointers(s.LuaBurstTimeout, t.LuaBurstTimeout) {
		diff["LuaBurstTimeout"] = []interface{}{ValueOrNil(s.LuaBurstTimeout), ValueOrNil(t.LuaBurstTimeout)}
	}

	if s.LuaForcedYield != t.LuaForcedYield {
		diff["LuaForcedYield"] = []interface{}{s.LuaForcedYield, t.LuaForcedYield}
	}

	if s.LuaLogLoggers != t.LuaLogLoggers {
		diff["LuaLogLoggers"] = []interface{}{s.LuaLogLoggers, t.LuaLogLoggers}
	}

	if s.LuaLogStderr != t.LuaLogStderr {
		diff["LuaLogStderr"] = []interface{}{s.LuaLogStderr, t.LuaLogStderr}
	}

	if s.LuaMaxmem != t.LuaMaxmem {
		diff["LuaMaxmem"] = []interface{}{s.LuaMaxmem, t.LuaMaxmem}
	}

	if !equalPointers(s.LuaServiceTimeout, t.LuaServiceTimeout) {
		diff["LuaServiceTimeout"] = []interface{}{ValueOrNil(s.LuaServiceTimeout), ValueOrNil(t.LuaServiceTimeout)}
	}

	if !equalPointers(s.LuaSessionTimeout, t.LuaSessionTimeout) {
		diff["LuaSessionTimeout"] = []interface{}{ValueOrNil(s.LuaSessionTimeout), ValueOrNil(t.LuaSessionTimeout)}
	}

	if !equalPointers(s.LuaTaskTimeout, t.LuaTaskTimeout) {
		diff["LuaTaskTimeout"] = []interface{}{ValueOrNil(s.LuaTaskTimeout), ValueOrNil(t.LuaTaskTimeout)}
	}

	if !equalPointers(s.MaxChecksPerThread, t.MaxChecksPerThread) {
		diff["MaxChecksPerThread"] = []interface{}{ValueOrNil(s.MaxChecksPerThread), ValueOrNil(t.MaxChecksPerThread)}
	}

	if s.Maxaccept != t.Maxaccept {
		diff["Maxaccept"] = []interface{}{s.Maxaccept, t.Maxaccept}
	}

	if s.Maxpollevents != t.Maxpollevents {
		diff["Maxpollevents"] = []interface{}{s.Maxpollevents, t.Maxpollevents}
	}

	if s.Maxrewrite != t.Maxrewrite {
		diff["Maxrewrite"] = []interface{}{s.Maxrewrite, t.Maxrewrite}
	}

	if !equalPointers(s.MemoryHotSize, t.MemoryHotSize) {
		diff["MemoryHotSize"] = []interface{}{ValueOrNil(s.MemoryHotSize), ValueOrNil(t.MemoryHotSize)}
	}

	if !equalPointers(s.PatternCacheSize, t.PatternCacheSize) {
		diff["PatternCacheSize"] = []interface{}{ValueOrNil(s.PatternCacheSize), ValueOrNil(t.PatternCacheSize)}
	}

	if s.PeersMaxUpdatesAtOnce != t.PeersMaxUpdatesAtOnce {
		diff["PeersMaxUpdatesAtOnce"] = []interface{}{s.PeersMaxUpdatesAtOnce, t.PeersMaxUpdatesAtOnce}
	}

	if s.Pipesize != t.Pipesize {
		diff["Pipesize"] = []interface{}{s.Pipesize, t.Pipesize}
	}

	if s.PoolHighFdRatio != t.PoolHighFdRatio {
		diff["PoolHighFdRatio"] = []interface{}{s.PoolHighFdRatio, t.PoolHighFdRatio}
	}

	if s.PoolLowFdRatio != t.PoolLowFdRatio {
		diff["PoolLowFdRatio"] = []interface{}{s.PoolLowFdRatio, t.PoolLowFdRatio}
	}

	if !equalPointers(s.QuicFrontendConnTxBuffersLimit, t.QuicFrontendConnTxBuffersLimit) {
		diff["QuicFrontendConnTxBuffersLimit"] = []interface{}{ValueOrNil(s.QuicFrontendConnTxBuffersLimit), ValueOrNil(t.QuicFrontendConnTxBuffersLimit)}
	}

	if !equalPointers(s.QuicFrontendMaxIdleTimeout, t.QuicFrontendMaxIdleTimeout) {
		diff["QuicFrontendMaxIdleTimeout"] = []interface{}{ValueOrNil(s.QuicFrontendMaxIdleTimeout), ValueOrNil(t.QuicFrontendMaxIdleTimeout)}
	}

	if !equalPointers(s.QuicFrontendMaxStreamsBidi, t.QuicFrontendMaxStreamsBidi) {
		diff["QuicFrontendMaxStreamsBidi"] = []interface{}{ValueOrNil(s.QuicFrontendMaxStreamsBidi), ValueOrNil(t.QuicFrontendMaxStreamsBidi)}
	}

	if !equalPointers(s.QuicMaxFrameLoss, t.QuicMaxFrameLoss) {
		diff["QuicMaxFrameLoss"] = []interface{}{ValueOrNil(s.QuicMaxFrameLoss), ValueOrNil(t.QuicMaxFrameLoss)}
	}

	if !equalPointers(s.QuicRetryThreshold, t.QuicRetryThreshold) {
		diff["QuicRetryThreshold"] = []interface{}{ValueOrNil(s.QuicRetryThreshold), ValueOrNil(t.QuicRetryThreshold)}
	}

	if s.QuicSocketOwner != t.QuicSocketOwner {
		diff["QuicSocketOwner"] = []interface{}{s.QuicSocketOwner, t.QuicSocketOwner}
	}

	if !equalPointers(s.RcvbufBackend, t.RcvbufBackend) {
		diff["RcvbufBackend"] = []interface{}{ValueOrNil(s.RcvbufBackend), ValueOrNil(t.RcvbufBackend)}
	}

	if !equalPointers(s.RcvbufClient, t.RcvbufClient) {
		diff["RcvbufClient"] = []interface{}{ValueOrNil(s.RcvbufClient), ValueOrNil(t.RcvbufClient)}
	}

	if !equalPointers(s.RcvbufFrontend, t.RcvbufFrontend) {
		diff["RcvbufFrontend"] = []interface{}{ValueOrNil(s.RcvbufFrontend), ValueOrNil(t.RcvbufFrontend)}
	}

	if !equalPointers(s.RcvbufServer, t.RcvbufServer) {
		diff["RcvbufServer"] = []interface{}{ValueOrNil(s.RcvbufServer), ValueOrNil(t.RcvbufServer)}
	}

	if s.RecvEnough != t.RecvEnough {
		diff["RecvEnough"] = []interface{}{s.RecvEnough, t.RecvEnough}
	}

	if s.RunqueueDepth != t.RunqueueDepth {
		diff["RunqueueDepth"] = []interface{}{s.RunqueueDepth, t.RunqueueDepth}
	}

	if s.SchedLowLatency != t.SchedLowLatency {
		diff["SchedLowLatency"] = []interface{}{s.SchedLowLatency, t.SchedLowLatency}
	}

	if !equalPointers(s.SndbufBackend, t.SndbufBackend) {
		diff["SndbufBackend"] = []interface{}{ValueOrNil(s.SndbufBackend), ValueOrNil(t.SndbufBackend)}
	}

	if !equalPointers(s.SndbufClient, t.SndbufClient) {
		diff["SndbufClient"] = []interface{}{ValueOrNil(s.SndbufClient), ValueOrNil(t.SndbufClient)}
	}

	if !equalPointers(s.SndbufFrontend, t.SndbufFrontend) {
		diff["SndbufFrontend"] = []interface{}{ValueOrNil(s.SndbufFrontend), ValueOrNil(t.SndbufFrontend)}
	}

	if !equalPointers(s.SndbufServer, t.SndbufServer) {
		diff["SndbufServer"] = []interface{}{ValueOrNil(s.SndbufServer), ValueOrNil(t.SndbufServer)}
	}

	if !equalPointers(s.SslCachesize, t.SslCachesize) {
		diff["SslCachesize"] = []interface{}{ValueOrNil(s.SslCachesize), ValueOrNil(t.SslCachesize)}
	}

	if !equalPointers(s.SslCaptureBufferSize, t.SslCaptureBufferSize) {
		diff["SslCaptureBufferSize"] = []interface{}{ValueOrNil(s.SslCaptureBufferSize), ValueOrNil(t.SslCaptureBufferSize)}
	}

	if s.SslCtxCacheSize != t.SslCtxCacheSize {
		diff["SslCtxCacheSize"] = []interface{}{s.SslCtxCacheSize, t.SslCtxCacheSize}
	}

	if s.SslDefaultDhParam != t.SslDefaultDhParam {
		diff["SslDefaultDhParam"] = []interface{}{s.SslDefaultDhParam, t.SslDefaultDhParam}
	}

	if s.SslForcePrivateCache != t.SslForcePrivateCache {
		diff["SslForcePrivateCache"] = []interface{}{s.SslForcePrivateCache, t.SslForcePrivateCache}
	}

	if s.SslKeylog != t.SslKeylog {
		diff["SslKeylog"] = []interface{}{s.SslKeylog, t.SslKeylog}
	}

	if !equalPointers(s.SslLifetime, t.SslLifetime) {
		diff["SslLifetime"] = []interface{}{ValueOrNil(s.SslLifetime), ValueOrNil(t.SslLifetime)}
	}

	if !equalPointers(s.SslMaxrecord, t.SslMaxrecord) {
		diff["SslMaxrecord"] = []interface{}{ValueOrNil(s.SslMaxrecord), ValueOrNil(t.SslMaxrecord)}
	}

	if !equalPointers(s.SslOcspUpdateMaxDelay, t.SslOcspUpdateMaxDelay) {
		diff["SslOcspUpdateMaxDelay"] = []interface{}{ValueOrNil(s.SslOcspUpdateMaxDelay), ValueOrNil(t.SslOcspUpdateMaxDelay)}
	}

	if !equalPointers(s.SslOcspUpdateMinDelay, t.SslOcspUpdateMinDelay) {
		diff["SslOcspUpdateMinDelay"] = []interface{}{ValueOrNil(s.SslOcspUpdateMinDelay), ValueOrNil(t.SslOcspUpdateMinDelay)}
	}

	if !equalPointers(s.StickCounters, t.StickCounters) {
		diff["StickCounters"] = []interface{}{ValueOrNil(s.StickCounters), ValueOrNil(t.StickCounters)}
	}

	if !equalPointers(s.VarsGlobalMaxSize, t.VarsGlobalMaxSize) {
		diff["VarsGlobalMaxSize"] = []interface{}{ValueOrNil(s.VarsGlobalMaxSize), ValueOrNil(t.VarsGlobalMaxSize)}
	}

	if !equalPointers(s.VarsProcMaxSize, t.VarsProcMaxSize) {
		diff["VarsProcMaxSize"] = []interface{}{ValueOrNil(s.VarsProcMaxSize), ValueOrNil(t.VarsProcMaxSize)}
	}

	if !equalPointers(s.VarsReqresMaxSize, t.VarsReqresMaxSize) {
		diff["VarsReqresMaxSize"] = []interface{}{ValueOrNil(s.VarsReqresMaxSize), ValueOrNil(t.VarsReqresMaxSize)}
	}

	if !equalPointers(s.VarsSessMaxSize, t.VarsSessMaxSize) {
		diff["VarsSessMaxSize"] = []interface{}{ValueOrNil(s.VarsSessMaxSize), ValueOrNil(t.VarsSessMaxSize)}
	}

	if !equalPointers(s.VarsTxnMaxSize, t.VarsTxnMaxSize) {
		diff["VarsTxnMaxSize"] = []interface{}{ValueOrNil(s.VarsTxnMaxSize), ValueOrNil(t.VarsTxnMaxSize)}
	}

	if s.ZlibMemlevel != t.ZlibMemlevel {
		diff["ZlibMemlevel"] = []interface{}{s.ZlibMemlevel, t.ZlibMemlevel}
	}

	if s.ZlibWindowsize != t.ZlibWindowsize {
		diff["ZlibWindowsize"] = []interface{}{s.ZlibWindowsize, t.ZlibWindowsize}
	}

	return diff
}

// Equal checks if two structs of type GlobalWurflOptions are equal
//
//	var a, b GlobalWurflOptions
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s GlobalWurflOptions) Equal(t GlobalWurflOptions, opts ...Options) bool {
	if s.CacheSize != t.CacheSize {
		return false
	}

	if s.DataFile != t.DataFile {
		return false
	}

	if s.InformationList != t.InformationList {
		return false
	}

	if s.InformationListSeparator != t.InformationListSeparator {
		return false
	}

	if s.PatchFile != t.PatchFile {
		return false
	}

	return true
}

// Diff checks if two structs of type GlobalWurflOptions are equal
//
//	var a, b GlobalWurflOptions
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s GlobalWurflOptions) Diff(t GlobalWurflOptions, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if s.CacheSize != t.CacheSize {
		diff["CacheSize"] = []interface{}{s.CacheSize, t.CacheSize}
	}

	if s.DataFile != t.DataFile {
		diff["DataFile"] = []interface{}{s.DataFile, t.DataFile}
	}

	if s.InformationList != t.InformationList {
		diff["InformationList"] = []interface{}{s.InformationList, t.InformationList}
	}

	if s.InformationListSeparator != t.InformationListSeparator {
		diff["InformationListSeparator"] = []interface{}{s.InformationListSeparator, t.InformationListSeparator}
	}

	if s.PatchFile != t.PatchFile {
		diff["PatchFile"] = []interface{}{s.PatchFile, t.PatchFile}
	}

	return diff
}
