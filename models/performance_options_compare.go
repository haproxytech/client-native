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

// Equal checks if two structs of type PerformanceOptions are equal
//
//	var a, b PerformanceOptions
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s PerformanceOptions) Equal(t PerformanceOptions, opts ...Options) bool {
	if s.BusyPolling != t.BusyPolling {
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

	if s.Maxzlibmem != t.Maxzlibmem {
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

	if s.ProfilingMemory != t.ProfilingMemory {
		return false
	}

	if s.ProfilingTasks != t.ProfilingTasks {
		return false
	}

	if s.ServerStateBase != t.ServerStateBase {
		return false
	}

	if s.ServerStateFile != t.ServerStateFile {
		return false
	}

	if s.SpreadChecks != t.SpreadChecks {
		return false
	}

	if !equalPointers(s.ThreadHardLimit, t.ThreadHardLimit) {
		return false
	}

	return true
}

// Diff checks if two structs of type PerformanceOptions are equal
//
//	var a, b PerformanceOptions
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s PerformanceOptions) Diff(t PerformanceOptions, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if s.BusyPolling != t.BusyPolling {
		diff["BusyPolling"] = []interface{}{s.BusyPolling, t.BusyPolling}
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

	if s.Maxzlibmem != t.Maxzlibmem {
		diff["Maxzlibmem"] = []interface{}{s.Maxzlibmem, t.Maxzlibmem}
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

	if s.ProfilingMemory != t.ProfilingMemory {
		diff["ProfilingMemory"] = []interface{}{s.ProfilingMemory, t.ProfilingMemory}
	}

	if s.ProfilingTasks != t.ProfilingTasks {
		diff["ProfilingTasks"] = []interface{}{s.ProfilingTasks, t.ProfilingTasks}
	}

	if s.ServerStateBase != t.ServerStateBase {
		diff["ServerStateBase"] = []interface{}{s.ServerStateBase, t.ServerStateBase}
	}

	if s.ServerStateFile != t.ServerStateFile {
		diff["ServerStateFile"] = []interface{}{s.ServerStateFile, t.ServerStateFile}
	}

	if s.SpreadChecks != t.SpreadChecks {
		diff["SpreadChecks"] = []interface{}{s.SpreadChecks, t.SpreadChecks}
	}

	if !equalPointers(s.ThreadHardLimit, t.ThreadHardLimit) {
		diff["ThreadHardLimit"] = []interface{}{ValueOrNil(s.ThreadHardLimit), ValueOrNil(t.ThreadHardLimit)}
	}

	return diff
}
