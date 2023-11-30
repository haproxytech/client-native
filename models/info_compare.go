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

// Equal checks if two structs of type Info are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b Info
//	equal := a.Equal(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b Info
//	equal := a.Equal(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s Info) Equal(t Info, opts ...Options) bool {
	opt := getOptions(opts...)

	if s.API == nil || t.API == nil {
		if s.API != nil || t.API != nil {
			if opt.NilSameAsEmpty {
				empty := &InfoAPI{}
				if s.API == nil {
					if !(t.API.Equal(*empty)) {
						return false
					}
				}
				if t.API == nil {
					if !(s.API.Equal(*empty)) {
						return false
					}
				}
			} else {
				return false
			}
		}
	} else if !s.API.Equal(*t.API, opt) {
		return false
	}

	if s.System == nil || t.System == nil {
		if s.System != nil || t.System != nil {
			if opt.NilSameAsEmpty {
				empty := &InfoSystem{}
				if s.System == nil {
					if !(t.System.Equal(*empty)) {
						return false
					}
				}
				if t.System == nil {
					if !(s.System.Equal(*empty)) {
						return false
					}
				}
			} else {
				return false
			}
		}
	} else if !s.System.Equal(*t.System, opt) {
		return false
	}

	return true
}

// Diff checks if two structs of type Info are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b Info
//	diff := a.Diff(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b Info
//	diff := a.Diff(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s Info) Diff(t Info, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})

	if s.API == nil || t.API == nil {
		if s.API != nil || t.API != nil {
			if opt.NilSameAsEmpty {
				empty := &InfoAPI{}
				if s.API == nil {
					if !(t.API.Equal(*empty)) {
						diff["API"] = []interface{}{ValueOrNil(s.API), ValueOrNil(t.API)}
					}
				}
				if t.API == nil {
					if !(s.API.Equal(*empty)) {
						diff["API"] = []interface{}{ValueOrNil(s.API), ValueOrNil(t.API)}
					}
				}
			} else {
				diff["API"] = []interface{}{ValueOrNil(s.API), ValueOrNil(t.API)}
			}
		}
	} else if !s.API.Equal(*t.API, opt) {
		diff["API"] = []interface{}{ValueOrNil(s.API), ValueOrNil(t.API)}
	}

	if s.System == nil || t.System == nil {
		if s.System != nil || t.System != nil {
			if opt.NilSameAsEmpty {
				empty := &InfoSystem{}
				if s.System == nil {
					if !(t.System.Equal(*empty)) {
						diff["System"] = []interface{}{ValueOrNil(s.System), ValueOrNil(t.System)}
					}
				}
				if t.System == nil {
					if !(s.System.Equal(*empty)) {
						diff["System"] = []interface{}{ValueOrNil(s.System), ValueOrNil(t.System)}
					}
				}
			} else {
				diff["System"] = []interface{}{ValueOrNil(s.System), ValueOrNil(t.System)}
			}
		}
	} else if !s.System.Equal(*t.System, opt) {
		diff["System"] = []interface{}{ValueOrNil(s.System), ValueOrNil(t.System)}
	}

	return diff
}

// Equal checks if two structs of type InfoAPI are equal
//
//	var a, b InfoAPI
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s InfoAPI) Equal(t InfoAPI, opts ...Options) bool {

	if !s.BuildDate.Equal(t.BuildDate) {
		return false
	}

	if s.Version != t.Version {
		return false
	}

	return true
}

// Diff checks if two structs of type InfoAPI are equal
//
//	var a, b InfoAPI
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s InfoAPI) Diff(t InfoAPI, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})

	if !s.BuildDate.Equal(t.BuildDate) {
		diff["BuildDate"] = []interface{}{s.BuildDate, t.BuildDate}
	}

	if s.Version != t.Version {
		diff["Version"] = []interface{}{s.Version, t.Version}
	}

	return diff
}

// Equal checks if two structs of type InfoSystem are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b InfoSystem
//	equal := a.Equal(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b InfoSystem
//	equal := a.Equal(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s InfoSystem) Equal(t InfoSystem, opts ...Options) bool {
	opt := getOptions(opts...)

	if s.CPUInfo == nil || t.CPUInfo == nil {
		if s.CPUInfo != nil || t.CPUInfo != nil {
			if opt.NilSameAsEmpty {
				empty := &InfoSystemCPUInfo{}
				if s.CPUInfo == nil {
					if !(t.CPUInfo.Equal(*empty)) {
						return false
					}
				}
				if t.CPUInfo == nil {
					if !(s.CPUInfo.Equal(*empty)) {
						return false
					}
				}
			} else {
				return false
			}
		}
	} else if !s.CPUInfo.Equal(*t.CPUInfo, opt) {
		return false
	}

	if s.Hostname != t.Hostname {
		return false
	}

	if s.MemInfo == nil || t.MemInfo == nil {
		if s.MemInfo != nil || t.MemInfo != nil {
			if opt.NilSameAsEmpty {
				empty := &InfoSystemMemInfo{}
				if s.MemInfo == nil {
					if !(t.MemInfo.Equal(*empty)) {
						return false
					}
				}
				if t.MemInfo == nil {
					if !(s.MemInfo.Equal(*empty)) {
						return false
					}
				}
			} else {
				return false
			}
		}
	} else if !s.MemInfo.Equal(*t.MemInfo, opt) {
		return false
	}

	if s.OsString != t.OsString {
		return false
	}

	if s.Time != t.Time {
		return false
	}

	if !equalPointers(s.Uptime, t.Uptime) {
		return false
	}

	return true
}

// Diff checks if two structs of type InfoSystem are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b InfoSystem
//	diff := a.Diff(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b InfoSystem
//	diff := a.Diff(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s InfoSystem) Diff(t InfoSystem, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})

	if s.CPUInfo == nil || t.CPUInfo == nil {
		if s.CPUInfo != nil || t.CPUInfo != nil {
			if opt.NilSameAsEmpty {
				empty := &InfoSystemCPUInfo{}
				if s.CPUInfo == nil {
					if !(t.CPUInfo.Equal(*empty)) {
						diff["CPUInfo"] = []interface{}{ValueOrNil(s.CPUInfo), ValueOrNil(t.CPUInfo)}
					}
				}
				if t.CPUInfo == nil {
					if !(s.CPUInfo.Equal(*empty)) {
						diff["CPUInfo"] = []interface{}{ValueOrNil(s.CPUInfo), ValueOrNil(t.CPUInfo)}
					}
				}
			} else {
				diff["CPUInfo"] = []interface{}{ValueOrNil(s.CPUInfo), ValueOrNil(t.CPUInfo)}
			}
		}
	} else if !s.CPUInfo.Equal(*t.CPUInfo, opt) {
		diff["CPUInfo"] = []interface{}{ValueOrNil(s.CPUInfo), ValueOrNil(t.CPUInfo)}
	}

	if s.Hostname != t.Hostname {
		diff["Hostname"] = []interface{}{s.Hostname, t.Hostname}
	}

	if s.MemInfo == nil || t.MemInfo == nil {
		if s.MemInfo != nil || t.MemInfo != nil {
			if opt.NilSameAsEmpty {
				empty := &InfoSystemMemInfo{}
				if s.MemInfo == nil {
					if !(t.MemInfo.Equal(*empty)) {
						diff["MemInfo"] = []interface{}{ValueOrNil(s.MemInfo), ValueOrNil(t.MemInfo)}
					}
				}
				if t.MemInfo == nil {
					if !(s.MemInfo.Equal(*empty)) {
						diff["MemInfo"] = []interface{}{ValueOrNil(s.MemInfo), ValueOrNil(t.MemInfo)}
					}
				}
			} else {
				diff["MemInfo"] = []interface{}{ValueOrNil(s.MemInfo), ValueOrNil(t.MemInfo)}
			}
		}
	} else if !s.MemInfo.Equal(*t.MemInfo, opt) {
		diff["MemInfo"] = []interface{}{ValueOrNil(s.MemInfo), ValueOrNil(t.MemInfo)}
	}

	if s.OsString != t.OsString {
		diff["OsString"] = []interface{}{s.OsString, t.OsString}
	}

	if s.Time != t.Time {
		diff["Time"] = []interface{}{s.Time, t.Time}
	}

	if !equalPointers(s.Uptime, t.Uptime) {
		diff["Uptime"] = []interface{}{ValueOrNil(s.Uptime), ValueOrNil(t.Uptime)}
	}

	return diff
}

// Equal checks if two structs of type InfoSystemCPUInfo are equal
//
//	var a, b InfoSystemCPUInfo
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s InfoSystemCPUInfo) Equal(t InfoSystemCPUInfo, opts ...Options) bool {
	if s.Model != t.Model {
		return false
	}

	if s.NumCpus != t.NumCpus {
		return false
	}

	return true
}

// Diff checks if two structs of type InfoSystemCPUInfo are equal
//
//	var a, b InfoSystemCPUInfo
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s InfoSystemCPUInfo) Diff(t InfoSystemCPUInfo, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if s.Model != t.Model {
		diff["Model"] = []interface{}{s.Model, t.Model}
	}

	if s.NumCpus != t.NumCpus {
		diff["NumCpus"] = []interface{}{s.NumCpus, t.NumCpus}
	}

	return diff
}

// Equal checks if two structs of type InfoSystemMemInfo are equal
//
//	var a, b InfoSystemMemInfo
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s InfoSystemMemInfo) Equal(t InfoSystemMemInfo, opts ...Options) bool {
	if s.DataplaneapiMemory != t.DataplaneapiMemory {
		return false
	}

	if s.FreeMemory != t.FreeMemory {
		return false
	}

	if s.TotalMemory != t.TotalMemory {
		return false
	}

	return true
}

// Diff checks if two structs of type InfoSystemMemInfo are equal
//
//	var a, b InfoSystemMemInfo
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s InfoSystemMemInfo) Diff(t InfoSystemMemInfo, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if s.DataplaneapiMemory != t.DataplaneapiMemory {
		diff["DataplaneapiMemory"] = []interface{}{s.DataplaneapiMemory, t.DataplaneapiMemory}
	}

	if s.FreeMemory != t.FreeMemory {
		diff["FreeMemory"] = []interface{}{s.FreeMemory, t.FreeMemory}
	}

	if s.TotalMemory != t.TotalMemory {
		diff["TotalMemory"] = []interface{}{s.TotalMemory, t.TotalMemory}
	}

	return diff
}
