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

// Equal checks if two structs of type StatsOptions are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b StatsOptions
//	equal := a.Equal(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b StatsOptions
//	equal := a.Equal(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s StatsOptions) Equal(t StatsOptions, opts ...Options) bool {
	opt := getOptions(opts...)

	if s.StatsAdmin != t.StatsAdmin {
		return false
	}

	if s.StatsAdminCond != t.StatsAdminCond {
		return false
	}

	if s.StatsAdminCondTest != t.StatsAdminCondTest {
		return false
	}

	if !CheckSameNilAndLen(s.StatsAuths, t.StatsAuths, opt) {
		return false
	} else {
		for i := range s.StatsAuths {
			if !s.StatsAuths[i].Equal(*t.StatsAuths[i], opt) {
				return false
			}
		}
	}

	if s.StatsEnable != t.StatsEnable {
		return false
	}

	if s.StatsHideVersion != t.StatsHideVersion {
		return false
	}

	if !CheckSameNilAndLen(s.StatsHTTPRequests, t.StatsHTTPRequests, opt) {
		return false
	} else {
		for i := range s.StatsHTTPRequests {
			if !s.StatsHTTPRequests[i].Equal(*t.StatsHTTPRequests[i], opt) {
				return false
			}
		}
	}

	if s.StatsMaxconn != t.StatsMaxconn {
		return false
	}

	if s.StatsRealm != t.StatsRealm {
		return false
	}

	if !equalPointers(s.StatsRealmRealm, t.StatsRealmRealm) {
		return false
	}

	if !equalPointers(s.StatsRefreshDelay, t.StatsRefreshDelay) {
		return false
	}

	if !equalPointers(s.StatsShowDesc, t.StatsShowDesc) {
		return false
	}

	if s.StatsShowLegends != t.StatsShowLegends {
		return false
	}

	if s.StatsShowModules != t.StatsShowModules {
		return false
	}

	if !equalPointers(s.StatsShowNodeName, t.StatsShowNodeName) {
		return false
	}

	if s.StatsURIPrefix != t.StatsURIPrefix {
		return false
	}

	return true
}

// Diff checks if two structs of type StatsOptions are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b StatsOptions
//	diff := a.Diff(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b StatsOptions
//	diff := a.Diff(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s StatsOptions) Diff(t StatsOptions, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})
	if s.StatsAdmin != t.StatsAdmin {
		diff["StatsAdmin"] = []interface{}{s.StatsAdmin, t.StatsAdmin}
	}

	if s.StatsAdminCond != t.StatsAdminCond {
		diff["StatsAdminCond"] = []interface{}{s.StatsAdminCond, t.StatsAdminCond}
	}

	if s.StatsAdminCondTest != t.StatsAdminCondTest {
		diff["StatsAdminCondTest"] = []interface{}{s.StatsAdminCondTest, t.StatsAdminCondTest}
	}

	if !CheckSameNilAndLen(s.StatsAuths, t.StatsAuths, opt) {
		diff["StatsAuths"] = []interface{}{s.StatsAuths, t.StatsAuths}
	} else {
		diff2 := make(map[string][]interface{})
		for i := range s.StatsAuths {
			if !s.StatsAuths[i].Equal(*t.StatsAuths[i], opt) {
				diffSub := s.StatsAuths[i].Diff(*t.StatsAuths[i], opt)
				if len(diffSub) > 0 {
					diff2[strconv.Itoa(i)] = []interface{}{diffSub}
				}
			}
		}
		if len(diff2) > 0 {
			diff["StatsAuths"] = []interface{}{diff2}
		}
	}

	if s.StatsEnable != t.StatsEnable {
		diff["StatsEnable"] = []interface{}{s.StatsEnable, t.StatsEnable}
	}

	if s.StatsHideVersion != t.StatsHideVersion {
		diff["StatsHideVersion"] = []interface{}{s.StatsHideVersion, t.StatsHideVersion}
	}

	if !CheckSameNilAndLen(s.StatsHTTPRequests, t.StatsHTTPRequests, opt) {
		diff["StatsHTTPRequests"] = []interface{}{s.StatsHTTPRequests, t.StatsHTTPRequests}
	} else {
		diff2 := make(map[string][]interface{})
		for i := range s.StatsHTTPRequests {
			if !s.StatsHTTPRequests[i].Equal(*t.StatsHTTPRequests[i], opt) {
				diffSub := s.StatsHTTPRequests[i].Diff(*t.StatsHTTPRequests[i], opt)
				if len(diffSub) > 0 {
					diff2[strconv.Itoa(i)] = []interface{}{diffSub}
				}
			}
		}
		if len(diff2) > 0 {
			diff["StatsHTTPRequests"] = []interface{}{diff2}
		}
	}

	if s.StatsMaxconn != t.StatsMaxconn {
		diff["StatsMaxconn"] = []interface{}{s.StatsMaxconn, t.StatsMaxconn}
	}

	if s.StatsRealm != t.StatsRealm {
		diff["StatsRealm"] = []interface{}{s.StatsRealm, t.StatsRealm}
	}

	if !equalPointers(s.StatsRealmRealm, t.StatsRealmRealm) {
		diff["StatsRealmRealm"] = []interface{}{ValueOrNil(s.StatsRealmRealm), ValueOrNil(t.StatsRealmRealm)}
	}

	if !equalPointers(s.StatsRefreshDelay, t.StatsRefreshDelay) {
		diff["StatsRefreshDelay"] = []interface{}{ValueOrNil(s.StatsRefreshDelay), ValueOrNil(t.StatsRefreshDelay)}
	}

	if !equalPointers(s.StatsShowDesc, t.StatsShowDesc) {
		diff["StatsShowDesc"] = []interface{}{ValueOrNil(s.StatsShowDesc), ValueOrNil(t.StatsShowDesc)}
	}

	if s.StatsShowLegends != t.StatsShowLegends {
		diff["StatsShowLegends"] = []interface{}{s.StatsShowLegends, t.StatsShowLegends}
	}

	if s.StatsShowModules != t.StatsShowModules {
		diff["StatsShowModules"] = []interface{}{s.StatsShowModules, t.StatsShowModules}
	}

	if !equalPointers(s.StatsShowNodeName, t.StatsShowNodeName) {
		diff["StatsShowNodeName"] = []interface{}{ValueOrNil(s.StatsShowNodeName), ValueOrNil(t.StatsShowNodeName)}
	}

	if s.StatsURIPrefix != t.StatsURIPrefix {
		diff["StatsURIPrefix"] = []interface{}{s.StatsURIPrefix, t.StatsURIPrefix}
	}

	return diff
}
