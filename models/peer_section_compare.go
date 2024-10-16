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

// Equal checks if two structs of type PeerSection are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b PeerSection
//	equal := a.Equal(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b PeerSection
//	equal := a.Equal(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s PeerSection) Equal(t PeerSection, opts ...Options) bool {
	opt := getOptions(opts...)

	if !s.PeerSectionBase.Equal(t.PeerSectionBase, opt) {
		return false
	}

	if !s.LogTargetList.Equal(t.LogTargetList, opt) {
		return false
	}

	if !CheckSameNilAndLenMap[string, Bind](s.Binds, t.Binds, opt) {
		return false
	}

	for k, v := range s.Binds {
		if !t.Binds[k].Equal(v, opt) {
			return false
		}
	}

	if !CheckSameNilAndLenMap[string, PeerEntry](s.PeerEntries, t.PeerEntries, opt) {
		return false
	}

	for k, v := range s.PeerEntries {
		if !t.PeerEntries[k].Equal(v, opt) {
			return false
		}
	}

	if !CheckSameNilAndLenMap[string, Server](s.Servers, t.Servers, opt) {
		return false
	}

	for k, v := range s.Servers {
		if !t.Servers[k].Equal(v, opt) {
			return false
		}
	}

	if !CheckSameNilAndLenMap[string, Table](s.Tables, t.Tables, opt) {
		return false
	}

	for k, v := range s.Tables {
		if !t.Tables[k].Equal(v, opt) {
			return false
		}
	}

	return true
}

// Diff checks if two structs of type PeerSection are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b PeerSection
//	diff := a.Diff(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b PeerSection
//	diff := a.Diff(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s PeerSection) Diff(t PeerSection, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})

	if !s.PeerSectionBase.Equal(t.PeerSectionBase, opt) {
		diff["PeerSectionBase"] = []interface{}{s.PeerSectionBase, t.PeerSectionBase}
	}

	if !s.LogTargetList.Equal(t.LogTargetList, opt) {
		diff["LogTargetList"] = []interface{}{s.LogTargetList, t.LogTargetList}
	}

	if !CheckSameNilAndLenMap[string, Bind](s.Binds, t.Binds, opt) {
		diff["Binds"] = []interface{}{s.Binds, t.Binds}
	}

	for k, v := range s.Binds {
		if !t.Binds[k].Equal(v, opt) {
			diff["Binds"] = []interface{}{s.Binds, t.Binds}
		}
	}

	if !CheckSameNilAndLenMap[string, PeerEntry](s.PeerEntries, t.PeerEntries, opt) {
		diff["PeerEntries"] = []interface{}{s.PeerEntries, t.PeerEntries}
	}

	for k, v := range s.PeerEntries {
		if !t.PeerEntries[k].Equal(v, opt) {
			diff["PeerEntries"] = []interface{}{s.PeerEntries, t.PeerEntries}
		}
	}

	if !CheckSameNilAndLenMap[string, Server](s.Servers, t.Servers, opt) {
		diff["Servers"] = []interface{}{s.Servers, t.Servers}
	}

	for k, v := range s.Servers {
		if !t.Servers[k].Equal(v, opt) {
			diff["Servers"] = []interface{}{s.Servers, t.Servers}
		}
	}

	if !CheckSameNilAndLenMap[string, Table](s.Tables, t.Tables, opt) {
		diff["Tables"] = []interface{}{s.Tables, t.Tables}
	}

	for k, v := range s.Tables {
		if !t.Tables[k].Equal(v, opt) {
			diff["Tables"] = []interface{}{s.Tables, t.Tables}
		}
	}

	return diff
}
