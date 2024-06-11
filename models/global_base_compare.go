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

// Equal checks if two structs of type GlobalBase are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b GlobalBase
//	equal := a.Equal(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b GlobalBase
//	equal := a.Equal(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s GlobalBase) Equal(t GlobalBase, opts ...Options) bool {
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

	if !CheckSameNilAndLen(s.RuntimeAPIs, t.RuntimeAPIs, opt) {
		return false
	} else {
		for i := range s.RuntimeAPIs {
			if !s.RuntimeAPIs[i].Equal(*t.RuntimeAPIs[i], opt) {
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

	if !CheckSameNilAndLen(s.ThreadGroupLines, t.ThreadGroupLines, opt) {
		return false
	} else {
		for i := range s.ThreadGroupLines {
			if !s.ThreadGroupLines[i].Equal(*t.ThreadGroupLines[i], opt) {
				return false
			}
		}
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

	if s.Daemon != t.Daemon {
		return false
	}

	if s.DebugOptions == nil || t.DebugOptions == nil {
		if s.DebugOptions != nil || t.DebugOptions != nil {
			if opt.NilSameAsEmpty {
				empty := &DebugOptions{}
				if s.DebugOptions == nil {
					if !(t.DebugOptions.Equal(*empty)) {
						return false
					}
				}
				if t.DebugOptions == nil {
					if !(s.DebugOptions.Equal(*empty)) {
						return false
					}
				}
			} else {
				return false
			}
		}
	} else if !s.DebugOptions.Equal(*t.DebugOptions, opt) {
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
				empty := &DeviceAtlasOptions{}
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

	if s.EnvironmentOptions == nil || t.EnvironmentOptions == nil {
		if s.EnvironmentOptions != nil || t.EnvironmentOptions != nil {
			if opt.NilSameAsEmpty {
				empty := &EnvironmentOptions{}
				if s.EnvironmentOptions == nil {
					if !(t.EnvironmentOptions.Equal(*empty)) {
						return false
					}
				}
				if t.EnvironmentOptions == nil {
					if !(s.EnvironmentOptions.Equal(*empty)) {
						return false
					}
				}
			} else {
				return false
			}
		}
	} else if !s.EnvironmentOptions.Equal(*t.EnvironmentOptions, opt) {
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
				empty := &FiftyOneDegreesOptions{}
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

	if s.Harden == nil || t.Harden == nil {
		if s.Harden != nil || t.Harden != nil {
			if opt.NilSameAsEmpty {
				empty := &GlobalHarden{}
				if s.Harden == nil {
					if !(t.Harden.Equal(*empty)) {
						return false
					}
				}
				if t.Harden == nil {
					if !(s.Harden.Equal(*empty)) {
						return false
					}
				}
			} else {
				return false
			}
		}
	} else if !s.Harden.Equal(*t.Harden, opt) {
		return false
	}

	if s.HTTPClientOptions == nil || t.HTTPClientOptions == nil {
		if s.HTTPClientOptions != nil || t.HTTPClientOptions != nil {
			if opt.NilSameAsEmpty {
				empty := &HTTPClientOptions{}
				if s.HTTPClientOptions == nil {
					if !(t.HTTPClientOptions.Equal(*empty)) {
						return false
					}
				}
				if t.HTTPClientOptions == nil {
					if !(s.HTTPClientOptions.Equal(*empty)) {
						return false
					}
				}
			} else {
				return false
			}
		}
	} else if !s.HTTPClientOptions.Equal(*t.HTTPClientOptions, opt) {
		return false
	}

	if !CheckSameNilAndLen(s.HTTPErrCodes, t.HTTPErrCodes, opt) {
		return false
	} else {
		for i := range s.HTTPErrCodes {
			if !s.HTTPErrCodes[i].Equal(*t.HTTPErrCodes[i], opt) {
				return false
			}
		}
	}

	if !CheckSameNilAndLen(s.HTTPFailCodes, t.HTTPFailCodes, opt) {
		return false
	} else {
		for i := range s.HTTPFailCodes {
			if !s.HTTPFailCodes[i].Equal(*t.HTTPFailCodes[i], opt) {
				return false
			}
		}
	}

	if s.InsecureForkWanted != t.InsecureForkWanted {
		return false
	}

	if s.InsecureSetuidWanted != t.InsecureSetuidWanted {
		return false
	}

	if s.LimitedQuic != t.LimitedQuic {
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

	if s.LuaOptions == nil || t.LuaOptions == nil {
		if s.LuaOptions != nil || t.LuaOptions != nil {
			if opt.NilSameAsEmpty {
				empty := &LuaOptions{}
				if s.LuaOptions == nil {
					if !(t.LuaOptions.Equal(*empty)) {
						return false
					}
				}
				if t.LuaOptions == nil {
					if !(s.LuaOptions.Equal(*empty)) {
						return false
					}
				}
			} else {
				return false
			}
		}
	} else if !s.LuaOptions.Equal(*t.LuaOptions, opt) {
		return false
	}

	if s.MasterWorker != t.MasterWorker {
		return false
	}

	if !equalPointers(s.MworkerMaxReloads, t.MworkerMaxReloads) {
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

	if s.NumaCPUMapping != t.NumaCPUMapping {
		return false
	}

	if s.OcspUpdateOptions == nil || t.OcspUpdateOptions == nil {
		if s.OcspUpdateOptions != nil || t.OcspUpdateOptions != nil {
			if opt.NilSameAsEmpty {
				empty := &OcspUpdateOptions{}
				if s.OcspUpdateOptions == nil {
					if !(t.OcspUpdateOptions.Equal(*empty)) {
						return false
					}
				}
				if t.OcspUpdateOptions == nil {
					if !(s.OcspUpdateOptions.Equal(*empty)) {
						return false
					}
				}
			} else {
				return false
			}
		}
	} else if !s.OcspUpdateOptions.Equal(*t.OcspUpdateOptions, opt) {
		return false
	}

	if s.PerformanceOptions == nil || t.PerformanceOptions == nil {
		if s.PerformanceOptions != nil || t.PerformanceOptions != nil {
			if opt.NilSameAsEmpty {
				empty := &PerformanceOptions{}
				if s.PerformanceOptions == nil {
					if !(t.PerformanceOptions.Equal(*empty)) {
						return false
					}
				}
				if t.PerformanceOptions == nil {
					if !(s.PerformanceOptions.Equal(*empty)) {
						return false
					}
				}
			} else {
				return false
			}
		}
	} else if !s.PerformanceOptions.Equal(*t.PerformanceOptions, opt) {
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

	if s.SetDumpable != t.SetDumpable {
		return false
	}

	if s.Setcap != t.Setcap {
		return false
	}

	if s.SslOptions == nil || t.SslOptions == nil {
		if s.SslOptions != nil || t.SslOptions != nil {
			if opt.NilSameAsEmpty {
				empty := &SslOptions{}
				if s.SslOptions == nil {
					if !(t.SslOptions.Equal(*empty)) {
						return false
					}
				}
				if t.SslOptions == nil {
					if !(s.SslOptions.Equal(*empty)) {
						return false
					}
				}
			} else {
				return false
			}
		}
	} else if !s.SslOptions.Equal(*t.SslOptions, opt) {
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

	if s.TuneBufferOptions == nil || t.TuneBufferOptions == nil {
		if s.TuneBufferOptions != nil || t.TuneBufferOptions != nil {
			if opt.NilSameAsEmpty {
				empty := &TuneBufferOptions{}
				if s.TuneBufferOptions == nil {
					if !(t.TuneBufferOptions.Equal(*empty)) {
						return false
					}
				}
				if t.TuneBufferOptions == nil {
					if !(s.TuneBufferOptions.Equal(*empty)) {
						return false
					}
				}
			} else {
				return false
			}
		}
	} else if !s.TuneBufferOptions.Equal(*t.TuneBufferOptions, opt) {
		return false
	}

	if s.TuneLuaOptions == nil || t.TuneLuaOptions == nil {
		if s.TuneLuaOptions != nil || t.TuneLuaOptions != nil {
			if opt.NilSameAsEmpty {
				empty := &TuneLuaOptions{}
				if s.TuneLuaOptions == nil {
					if !(t.TuneLuaOptions.Equal(*empty)) {
						return false
					}
				}
				if t.TuneLuaOptions == nil {
					if !(s.TuneLuaOptions.Equal(*empty)) {
						return false
					}
				}
			} else {
				return false
			}
		}
	} else if !s.TuneLuaOptions.Equal(*t.TuneLuaOptions, opt) {
		return false
	}

	if s.TuneOptions == nil || t.TuneOptions == nil {
		if s.TuneOptions != nil || t.TuneOptions != nil {
			if opt.NilSameAsEmpty {
				empty := &TuneOptions{}
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

	if s.TuneQuicOptions == nil || t.TuneQuicOptions == nil {
		if s.TuneQuicOptions != nil || t.TuneQuicOptions != nil {
			if opt.NilSameAsEmpty {
				empty := &TuneQuicOptions{}
				if s.TuneQuicOptions == nil {
					if !(t.TuneQuicOptions.Equal(*empty)) {
						return false
					}
				}
				if t.TuneQuicOptions == nil {
					if !(s.TuneQuicOptions.Equal(*empty)) {
						return false
					}
				}
			} else {
				return false
			}
		}
	} else if !s.TuneQuicOptions.Equal(*t.TuneQuicOptions, opt) {
		return false
	}

	if s.TuneSslOptions == nil || t.TuneSslOptions == nil {
		if s.TuneSslOptions != nil || t.TuneSslOptions != nil {
			if opt.NilSameAsEmpty {
				empty := &TuneSslOptions{}
				if s.TuneSslOptions == nil {
					if !(t.TuneSslOptions.Equal(*empty)) {
						return false
					}
				}
				if t.TuneSslOptions == nil {
					if !(s.TuneSslOptions.Equal(*empty)) {
						return false
					}
				}
			} else {
				return false
			}
		}
	} else if !s.TuneSslOptions.Equal(*t.TuneSslOptions, opt) {
		return false
	}

	if s.TuneVarsOptions == nil || t.TuneVarsOptions == nil {
		if s.TuneVarsOptions != nil || t.TuneVarsOptions != nil {
			if opt.NilSameAsEmpty {
				empty := &TuneVarsOptions{}
				if s.TuneVarsOptions == nil {
					if !(t.TuneVarsOptions.Equal(*empty)) {
						return false
					}
				}
				if t.TuneVarsOptions == nil {
					if !(s.TuneVarsOptions.Equal(*empty)) {
						return false
					}
				}
			} else {
				return false
			}
		}
	} else if !s.TuneVarsOptions.Equal(*t.TuneVarsOptions, opt) {
		return false
	}

	if s.TuneZlibOptions == nil || t.TuneZlibOptions == nil {
		if s.TuneZlibOptions != nil || t.TuneZlibOptions != nil {
			if opt.NilSameAsEmpty {
				empty := &TuneZlibOptions{}
				if s.TuneZlibOptions == nil {
					if !(t.TuneZlibOptions.Equal(*empty)) {
						return false
					}
				}
				if t.TuneZlibOptions == nil {
					if !(s.TuneZlibOptions.Equal(*empty)) {
						return false
					}
				}
			} else {
				return false
			}
		}
	} else if !s.TuneZlibOptions.Equal(*t.TuneZlibOptions, opt) {
		return false
	}

	if s.UID != t.UID {
		return false
	}

	if s.Ulimitn != t.Ulimitn {
		return false
	}

	if s.User != t.User {
		return false
	}

	if s.WurflOptions == nil || t.WurflOptions == nil {
		if s.WurflOptions != nil || t.WurflOptions != nil {
			if opt.NilSameAsEmpty {
				empty := &WurflOptions{}
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

	return true
}

// Diff checks if two structs of type GlobalBase are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b GlobalBase
//	diff := a.Diff(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b GlobalBase
//	diff := a.Diff(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s GlobalBase) Diff(t GlobalBase, opts ...Options) map[string][]interface{} {
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

	if s.Chroot != t.Chroot {
		diff["Chroot"] = []interface{}{s.Chroot, t.Chroot}
	}

	if !equalPointers(s.CloseSpreadTime, t.CloseSpreadTime) {
		diff["CloseSpreadTime"] = []interface{}{ValueOrNil(s.CloseSpreadTime), ValueOrNil(t.CloseSpreadTime)}
	}

	if s.ClusterSecret != t.ClusterSecret {
		diff["ClusterSecret"] = []interface{}{s.ClusterSecret, t.ClusterSecret}
	}

	if s.Daemon != t.Daemon {
		diff["Daemon"] = []interface{}{s.Daemon, t.Daemon}
	}

	if s.DebugOptions == nil || t.DebugOptions == nil {
		if s.DebugOptions != nil || t.DebugOptions != nil {
			if opt.NilSameAsEmpty {
				empty := &DebugOptions{}
				if s.DebugOptions == nil {
					if !(t.DebugOptions.Equal(*empty)) {
						diff["DebugOptions"] = []interface{}{ValueOrNil(s.DebugOptions), ValueOrNil(t.DebugOptions)}
					}
				}
				if t.DebugOptions == nil {
					if !(s.DebugOptions.Equal(*empty)) {
						diff["DebugOptions"] = []interface{}{ValueOrNil(s.DebugOptions), ValueOrNil(t.DebugOptions)}
					}
				}
			} else {
				diff["DebugOptions"] = []interface{}{ValueOrNil(s.DebugOptions), ValueOrNil(t.DebugOptions)}
			}
		}
	} else if !s.DebugOptions.Equal(*t.DebugOptions, opt) {
		diff["DebugOptions"] = []interface{}{ValueOrNil(s.DebugOptions), ValueOrNil(t.DebugOptions)}
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
				empty := &DeviceAtlasOptions{}
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

	if s.EnvironmentOptions == nil || t.EnvironmentOptions == nil {
		if s.EnvironmentOptions != nil || t.EnvironmentOptions != nil {
			if opt.NilSameAsEmpty {
				empty := &EnvironmentOptions{}
				if s.EnvironmentOptions == nil {
					if !(t.EnvironmentOptions.Equal(*empty)) {
						diff["EnvironmentOptions"] = []interface{}{ValueOrNil(s.EnvironmentOptions), ValueOrNil(t.EnvironmentOptions)}
					}
				}
				if t.EnvironmentOptions == nil {
					if !(s.EnvironmentOptions.Equal(*empty)) {
						diff["EnvironmentOptions"] = []interface{}{ValueOrNil(s.EnvironmentOptions), ValueOrNil(t.EnvironmentOptions)}
					}
				}
			} else {
				diff["EnvironmentOptions"] = []interface{}{ValueOrNil(s.EnvironmentOptions), ValueOrNil(t.EnvironmentOptions)}
			}
		}
	} else if !s.EnvironmentOptions.Equal(*t.EnvironmentOptions, opt) {
		diff["EnvironmentOptions"] = []interface{}{ValueOrNil(s.EnvironmentOptions), ValueOrNil(t.EnvironmentOptions)}
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
				empty := &FiftyOneDegreesOptions{}
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

	if s.Harden == nil || t.Harden == nil {
		if s.Harden != nil || t.Harden != nil {
			if opt.NilSameAsEmpty {
				empty := &GlobalHarden{}
				if s.Harden == nil {
					if !(t.Harden.Equal(*empty)) {
						diff["Harden"] = []interface{}{ValueOrNil(s.Harden), ValueOrNil(t.Harden)}
					}
				}
				if t.Harden == nil {
					if !(s.Harden.Equal(*empty)) {
						diff["Harden"] = []interface{}{ValueOrNil(s.Harden), ValueOrNil(t.Harden)}
					}
				}
			} else {
				diff["Harden"] = []interface{}{ValueOrNil(s.Harden), ValueOrNil(t.Harden)}
			}
		}
	} else if !s.Harden.Equal(*t.Harden, opt) {
		diff["Harden"] = []interface{}{ValueOrNil(s.Harden), ValueOrNil(t.Harden)}
	}

	if s.HTTPClientOptions == nil || t.HTTPClientOptions == nil {
		if s.HTTPClientOptions != nil || t.HTTPClientOptions != nil {
			if opt.NilSameAsEmpty {
				empty := &HTTPClientOptions{}
				if s.HTTPClientOptions == nil {
					if !(t.HTTPClientOptions.Equal(*empty)) {
						diff["HTTPClientOptions"] = []interface{}{ValueOrNil(s.HTTPClientOptions), ValueOrNil(t.HTTPClientOptions)}
					}
				}
				if t.HTTPClientOptions == nil {
					if !(s.HTTPClientOptions.Equal(*empty)) {
						diff["HTTPClientOptions"] = []interface{}{ValueOrNil(s.HTTPClientOptions), ValueOrNil(t.HTTPClientOptions)}
					}
				}
			} else {
				diff["HTTPClientOptions"] = []interface{}{ValueOrNil(s.HTTPClientOptions), ValueOrNil(t.HTTPClientOptions)}
			}
		}
	} else if !s.HTTPClientOptions.Equal(*t.HTTPClientOptions, opt) {
		diff["HTTPClientOptions"] = []interface{}{ValueOrNil(s.HTTPClientOptions), ValueOrNil(t.HTTPClientOptions)}
	}

	if !CheckSameNilAndLen(s.HTTPErrCodes, t.HTTPErrCodes, opt) {
		diff["HTTPErrCodes"] = []interface{}{s.HTTPErrCodes, t.HTTPErrCodes}
	} else {
		diff2 := make(map[string][]interface{})
		for i := range s.HTTPErrCodes {
			if !s.HTTPErrCodes[i].Equal(*t.HTTPErrCodes[i], opt) {
				diffSub := s.HTTPErrCodes[i].Diff(*t.HTTPErrCodes[i], opt)
				if len(diffSub) > 0 {
					diff2[strconv.Itoa(i)] = []interface{}{diffSub}
				}
			}
		}
		if len(diff2) > 0 {
			diff["HTTPErrCodes"] = []interface{}{diff2}
		}
	}

	if !CheckSameNilAndLen(s.HTTPFailCodes, t.HTTPFailCodes, opt) {
		diff["HTTPFailCodes"] = []interface{}{s.HTTPFailCodes, t.HTTPFailCodes}
	} else {
		diff2 := make(map[string][]interface{})
		for i := range s.HTTPFailCodes {
			if !s.HTTPFailCodes[i].Equal(*t.HTTPFailCodes[i], opt) {
				diffSub := s.HTTPFailCodes[i].Diff(*t.HTTPFailCodes[i], opt)
				if len(diffSub) > 0 {
					diff2[strconv.Itoa(i)] = []interface{}{diffSub}
				}
			}
		}
		if len(diff2) > 0 {
			diff["HTTPFailCodes"] = []interface{}{diff2}
		}
	}

	if s.InsecureForkWanted != t.InsecureForkWanted {
		diff["InsecureForkWanted"] = []interface{}{s.InsecureForkWanted, t.InsecureForkWanted}
	}

	if s.InsecureSetuidWanted != t.InsecureSetuidWanted {
		diff["InsecureSetuidWanted"] = []interface{}{s.InsecureSetuidWanted, t.InsecureSetuidWanted}
	}

	if s.LimitedQuic != t.LimitedQuic {
		diff["LimitedQuic"] = []interface{}{s.LimitedQuic, t.LimitedQuic}
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

	if s.LuaOptions == nil || t.LuaOptions == nil {
		if s.LuaOptions != nil || t.LuaOptions != nil {
			if opt.NilSameAsEmpty {
				empty := &LuaOptions{}
				if s.LuaOptions == nil {
					if !(t.LuaOptions.Equal(*empty)) {
						diff["LuaOptions"] = []interface{}{ValueOrNil(s.LuaOptions), ValueOrNil(t.LuaOptions)}
					}
				}
				if t.LuaOptions == nil {
					if !(s.LuaOptions.Equal(*empty)) {
						diff["LuaOptions"] = []interface{}{ValueOrNil(s.LuaOptions), ValueOrNil(t.LuaOptions)}
					}
				}
			} else {
				diff["LuaOptions"] = []interface{}{ValueOrNil(s.LuaOptions), ValueOrNil(t.LuaOptions)}
			}
		}
	} else if !s.LuaOptions.Equal(*t.LuaOptions, opt) {
		diff["LuaOptions"] = []interface{}{ValueOrNil(s.LuaOptions), ValueOrNil(t.LuaOptions)}
	}

	if s.MasterWorker != t.MasterWorker {
		diff["MasterWorker"] = []interface{}{s.MasterWorker, t.MasterWorker}
	}

	if !equalPointers(s.MworkerMaxReloads, t.MworkerMaxReloads) {
		diff["MworkerMaxReloads"] = []interface{}{ValueOrNil(s.MworkerMaxReloads), ValueOrNil(t.MworkerMaxReloads)}
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

	if s.NumaCPUMapping != t.NumaCPUMapping {
		diff["NumaCPUMapping"] = []interface{}{s.NumaCPUMapping, t.NumaCPUMapping}
	}

	if s.OcspUpdateOptions == nil || t.OcspUpdateOptions == nil {
		if s.OcspUpdateOptions != nil || t.OcspUpdateOptions != nil {
			if opt.NilSameAsEmpty {
				empty := &OcspUpdateOptions{}
				if s.OcspUpdateOptions == nil {
					if !(t.OcspUpdateOptions.Equal(*empty)) {
						diff["OcspUpdateOptions"] = []interface{}{ValueOrNil(s.OcspUpdateOptions), ValueOrNil(t.OcspUpdateOptions)}
					}
				}
				if t.OcspUpdateOptions == nil {
					if !(s.OcspUpdateOptions.Equal(*empty)) {
						diff["OcspUpdateOptions"] = []interface{}{ValueOrNil(s.OcspUpdateOptions), ValueOrNil(t.OcspUpdateOptions)}
					}
				}
			} else {
				diff["OcspUpdateOptions"] = []interface{}{ValueOrNil(s.OcspUpdateOptions), ValueOrNil(t.OcspUpdateOptions)}
			}
		}
	} else if !s.OcspUpdateOptions.Equal(*t.OcspUpdateOptions, opt) {
		diff["OcspUpdateOptions"] = []interface{}{ValueOrNil(s.OcspUpdateOptions), ValueOrNil(t.OcspUpdateOptions)}
	}

	if s.PerformanceOptions == nil || t.PerformanceOptions == nil {
		if s.PerformanceOptions != nil || t.PerformanceOptions != nil {
			if opt.NilSameAsEmpty {
				empty := &PerformanceOptions{}
				if s.PerformanceOptions == nil {
					if !(t.PerformanceOptions.Equal(*empty)) {
						diff["PerformanceOptions"] = []interface{}{ValueOrNil(s.PerformanceOptions), ValueOrNil(t.PerformanceOptions)}
					}
				}
				if t.PerformanceOptions == nil {
					if !(s.PerformanceOptions.Equal(*empty)) {
						diff["PerformanceOptions"] = []interface{}{ValueOrNil(s.PerformanceOptions), ValueOrNil(t.PerformanceOptions)}
					}
				}
			} else {
				diff["PerformanceOptions"] = []interface{}{ValueOrNil(s.PerformanceOptions), ValueOrNil(t.PerformanceOptions)}
			}
		}
	} else if !s.PerformanceOptions.Equal(*t.PerformanceOptions, opt) {
		diff["PerformanceOptions"] = []interface{}{ValueOrNil(s.PerformanceOptions), ValueOrNil(t.PerformanceOptions)}
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

	if s.SetDumpable != t.SetDumpable {
		diff["SetDumpable"] = []interface{}{s.SetDumpable, t.SetDumpable}
	}

	if s.Setcap != t.Setcap {
		diff["Setcap"] = []interface{}{s.Setcap, t.Setcap}
	}

	if s.SslOptions == nil || t.SslOptions == nil {
		if s.SslOptions != nil || t.SslOptions != nil {
			if opt.NilSameAsEmpty {
				empty := &SslOptions{}
				if s.SslOptions == nil {
					if !(t.SslOptions.Equal(*empty)) {
						diff["SslOptions"] = []interface{}{ValueOrNil(s.SslOptions), ValueOrNil(t.SslOptions)}
					}
				}
				if t.SslOptions == nil {
					if !(s.SslOptions.Equal(*empty)) {
						diff["SslOptions"] = []interface{}{ValueOrNil(s.SslOptions), ValueOrNil(t.SslOptions)}
					}
				}
			} else {
				diff["SslOptions"] = []interface{}{ValueOrNil(s.SslOptions), ValueOrNil(t.SslOptions)}
			}
		}
	} else if !s.SslOptions.Equal(*t.SslOptions, opt) {
		diff["SslOptions"] = []interface{}{ValueOrNil(s.SslOptions), ValueOrNil(t.SslOptions)}
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

	if s.TuneBufferOptions == nil || t.TuneBufferOptions == nil {
		if s.TuneBufferOptions != nil || t.TuneBufferOptions != nil {
			if opt.NilSameAsEmpty {
				empty := &TuneBufferOptions{}
				if s.TuneBufferOptions == nil {
					if !(t.TuneBufferOptions.Equal(*empty)) {
						diff["TuneBufferOptions"] = []interface{}{ValueOrNil(s.TuneBufferOptions), ValueOrNil(t.TuneBufferOptions)}
					}
				}
				if t.TuneBufferOptions == nil {
					if !(s.TuneBufferOptions.Equal(*empty)) {
						diff["TuneBufferOptions"] = []interface{}{ValueOrNil(s.TuneBufferOptions), ValueOrNil(t.TuneBufferOptions)}
					}
				}
			} else {
				diff["TuneBufferOptions"] = []interface{}{ValueOrNil(s.TuneBufferOptions), ValueOrNil(t.TuneBufferOptions)}
			}
		}
	} else if !s.TuneBufferOptions.Equal(*t.TuneBufferOptions, opt) {
		diff["TuneBufferOptions"] = []interface{}{ValueOrNil(s.TuneBufferOptions), ValueOrNil(t.TuneBufferOptions)}
	}

	if s.TuneLuaOptions == nil || t.TuneLuaOptions == nil {
		if s.TuneLuaOptions != nil || t.TuneLuaOptions != nil {
			if opt.NilSameAsEmpty {
				empty := &TuneLuaOptions{}
				if s.TuneLuaOptions == nil {
					if !(t.TuneLuaOptions.Equal(*empty)) {
						diff["TuneLuaOptions"] = []interface{}{ValueOrNil(s.TuneLuaOptions), ValueOrNil(t.TuneLuaOptions)}
					}
				}
				if t.TuneLuaOptions == nil {
					if !(s.TuneLuaOptions.Equal(*empty)) {
						diff["TuneLuaOptions"] = []interface{}{ValueOrNil(s.TuneLuaOptions), ValueOrNil(t.TuneLuaOptions)}
					}
				}
			} else {
				diff["TuneLuaOptions"] = []interface{}{ValueOrNil(s.TuneLuaOptions), ValueOrNil(t.TuneLuaOptions)}
			}
		}
	} else if !s.TuneLuaOptions.Equal(*t.TuneLuaOptions, opt) {
		diff["TuneLuaOptions"] = []interface{}{ValueOrNil(s.TuneLuaOptions), ValueOrNil(t.TuneLuaOptions)}
	}

	if s.TuneOptions == nil || t.TuneOptions == nil {
		if s.TuneOptions != nil || t.TuneOptions != nil {
			if opt.NilSameAsEmpty {
				empty := &TuneOptions{}
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

	if s.TuneQuicOptions == nil || t.TuneQuicOptions == nil {
		if s.TuneQuicOptions != nil || t.TuneQuicOptions != nil {
			if opt.NilSameAsEmpty {
				empty := &TuneQuicOptions{}
				if s.TuneQuicOptions == nil {
					if !(t.TuneQuicOptions.Equal(*empty)) {
						diff["TuneQuicOptions"] = []interface{}{ValueOrNil(s.TuneQuicOptions), ValueOrNil(t.TuneQuicOptions)}
					}
				}
				if t.TuneQuicOptions == nil {
					if !(s.TuneQuicOptions.Equal(*empty)) {
						diff["TuneQuicOptions"] = []interface{}{ValueOrNil(s.TuneQuicOptions), ValueOrNil(t.TuneQuicOptions)}
					}
				}
			} else {
				diff["TuneQuicOptions"] = []interface{}{ValueOrNil(s.TuneQuicOptions), ValueOrNil(t.TuneQuicOptions)}
			}
		}
	} else if !s.TuneQuicOptions.Equal(*t.TuneQuicOptions, opt) {
		diff["TuneQuicOptions"] = []interface{}{ValueOrNil(s.TuneQuicOptions), ValueOrNil(t.TuneQuicOptions)}
	}

	if s.TuneSslOptions == nil || t.TuneSslOptions == nil {
		if s.TuneSslOptions != nil || t.TuneSslOptions != nil {
			if opt.NilSameAsEmpty {
				empty := &TuneSslOptions{}
				if s.TuneSslOptions == nil {
					if !(t.TuneSslOptions.Equal(*empty)) {
						diff["TuneSslOptions"] = []interface{}{ValueOrNil(s.TuneSslOptions), ValueOrNil(t.TuneSslOptions)}
					}
				}
				if t.TuneSslOptions == nil {
					if !(s.TuneSslOptions.Equal(*empty)) {
						diff["TuneSslOptions"] = []interface{}{ValueOrNil(s.TuneSslOptions), ValueOrNil(t.TuneSslOptions)}
					}
				}
			} else {
				diff["TuneSslOptions"] = []interface{}{ValueOrNil(s.TuneSslOptions), ValueOrNil(t.TuneSslOptions)}
			}
		}
	} else if !s.TuneSslOptions.Equal(*t.TuneSslOptions, opt) {
		diff["TuneSslOptions"] = []interface{}{ValueOrNil(s.TuneSslOptions), ValueOrNil(t.TuneSslOptions)}
	}

	if s.TuneVarsOptions == nil || t.TuneVarsOptions == nil {
		if s.TuneVarsOptions != nil || t.TuneVarsOptions != nil {
			if opt.NilSameAsEmpty {
				empty := &TuneVarsOptions{}
				if s.TuneVarsOptions == nil {
					if !(t.TuneVarsOptions.Equal(*empty)) {
						diff["TuneVarsOptions"] = []interface{}{ValueOrNil(s.TuneVarsOptions), ValueOrNil(t.TuneVarsOptions)}
					}
				}
				if t.TuneVarsOptions == nil {
					if !(s.TuneVarsOptions.Equal(*empty)) {
						diff["TuneVarsOptions"] = []interface{}{ValueOrNil(s.TuneVarsOptions), ValueOrNil(t.TuneVarsOptions)}
					}
				}
			} else {
				diff["TuneVarsOptions"] = []interface{}{ValueOrNil(s.TuneVarsOptions), ValueOrNil(t.TuneVarsOptions)}
			}
		}
	} else if !s.TuneVarsOptions.Equal(*t.TuneVarsOptions, opt) {
		diff["TuneVarsOptions"] = []interface{}{ValueOrNil(s.TuneVarsOptions), ValueOrNil(t.TuneVarsOptions)}
	}

	if s.TuneZlibOptions == nil || t.TuneZlibOptions == nil {
		if s.TuneZlibOptions != nil || t.TuneZlibOptions != nil {
			if opt.NilSameAsEmpty {
				empty := &TuneZlibOptions{}
				if s.TuneZlibOptions == nil {
					if !(t.TuneZlibOptions.Equal(*empty)) {
						diff["TuneZlibOptions"] = []interface{}{ValueOrNil(s.TuneZlibOptions), ValueOrNil(t.TuneZlibOptions)}
					}
				}
				if t.TuneZlibOptions == nil {
					if !(s.TuneZlibOptions.Equal(*empty)) {
						diff["TuneZlibOptions"] = []interface{}{ValueOrNil(s.TuneZlibOptions), ValueOrNil(t.TuneZlibOptions)}
					}
				}
			} else {
				diff["TuneZlibOptions"] = []interface{}{ValueOrNil(s.TuneZlibOptions), ValueOrNil(t.TuneZlibOptions)}
			}
		}
	} else if !s.TuneZlibOptions.Equal(*t.TuneZlibOptions, opt) {
		diff["TuneZlibOptions"] = []interface{}{ValueOrNil(s.TuneZlibOptions), ValueOrNil(t.TuneZlibOptions)}
	}

	if s.UID != t.UID {
		diff["UID"] = []interface{}{s.UID, t.UID}
	}

	if s.Ulimitn != t.Ulimitn {
		diff["Ulimitn"] = []interface{}{s.Ulimitn, t.Ulimitn}
	}

	if s.User != t.User {
		diff["User"] = []interface{}{s.User, t.User}
	}

	if s.WurflOptions == nil || t.WurflOptions == nil {
		if s.WurflOptions != nil || t.WurflOptions != nil {
			if opt.NilSameAsEmpty {
				empty := &WurflOptions{}
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

// Equal checks if two structs of type GlobalHarden are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b GlobalHarden
//	equal := a.Equal(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b GlobalHarden
//	equal := a.Equal(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s GlobalHarden) Equal(t GlobalHarden, opts ...Options) bool {
	opt := getOptions(opts...)

	if s.RejectPrivilegedPorts == nil || t.RejectPrivilegedPorts == nil {
		if s.RejectPrivilegedPorts != nil || t.RejectPrivilegedPorts != nil {
			if opt.NilSameAsEmpty {
				empty := &GlobalHardenRejectPrivilegedPorts{}
				if s.RejectPrivilegedPorts == nil {
					if !(t.RejectPrivilegedPorts.Equal(*empty)) {
						return false
					}
				}
				if t.RejectPrivilegedPorts == nil {
					if !(s.RejectPrivilegedPorts.Equal(*empty)) {
						return false
					}
				}
			} else {
				return false
			}
		}
	} else if !s.RejectPrivilegedPorts.Equal(*t.RejectPrivilegedPorts, opt) {
		return false
	}

	return true
}

// Diff checks if two structs of type GlobalHarden are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b GlobalHarden
//	diff := a.Diff(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b GlobalHarden
//	diff := a.Diff(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s GlobalHarden) Diff(t GlobalHarden, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})

	if s.RejectPrivilegedPorts == nil || t.RejectPrivilegedPorts == nil {
		if s.RejectPrivilegedPorts != nil || t.RejectPrivilegedPorts != nil {
			if opt.NilSameAsEmpty {
				empty := &GlobalHardenRejectPrivilegedPorts{}
				if s.RejectPrivilegedPorts == nil {
					if !(t.RejectPrivilegedPorts.Equal(*empty)) {
						diff["RejectPrivilegedPorts"] = []interface{}{ValueOrNil(s.RejectPrivilegedPorts), ValueOrNil(t.RejectPrivilegedPorts)}
					}
				}
				if t.RejectPrivilegedPorts == nil {
					if !(s.RejectPrivilegedPorts.Equal(*empty)) {
						diff["RejectPrivilegedPorts"] = []interface{}{ValueOrNil(s.RejectPrivilegedPorts), ValueOrNil(t.RejectPrivilegedPorts)}
					}
				}
			} else {
				diff["RejectPrivilegedPorts"] = []interface{}{ValueOrNil(s.RejectPrivilegedPorts), ValueOrNil(t.RejectPrivilegedPorts)}
			}
		}
	} else if !s.RejectPrivilegedPorts.Equal(*t.RejectPrivilegedPorts, opt) {
		diff["RejectPrivilegedPorts"] = []interface{}{ValueOrNil(s.RejectPrivilegedPorts), ValueOrNil(t.RejectPrivilegedPorts)}
	}

	return diff
}

// Equal checks if two structs of type GlobalHardenRejectPrivilegedPorts are equal
//
//	var a, b GlobalHardenRejectPrivilegedPorts
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s GlobalHardenRejectPrivilegedPorts) Equal(t GlobalHardenRejectPrivilegedPorts, opts ...Options) bool {
	if s.Quic != t.Quic {
		return false
	}

	if s.TCP != t.TCP {
		return false
	}

	return true
}

// Diff checks if two structs of type GlobalHardenRejectPrivilegedPorts are equal
//
//	var a, b GlobalHardenRejectPrivilegedPorts
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s GlobalHardenRejectPrivilegedPorts) Diff(t GlobalHardenRejectPrivilegedPorts, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if s.Quic != t.Quic {
		diff["Quic"] = []interface{}{s.Quic, t.Quic}
	}

	if s.TCP != t.TCP {
		diff["TCP"] = []interface{}{s.TCP, t.TCP}
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
