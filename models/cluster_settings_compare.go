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

// Equal checks if two structs of type ClusterSettings are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b ClusterSettings
//	equal := a.Equal(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b ClusterSettings
//	equal := a.Equal(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s ClusterSettings) Equal(t ClusterSettings, opts ...Options) bool {
	opt := getOptions(opts...)

	if s.BootstrapKey != t.BootstrapKey {
		return false
	}

	if !s.Cluster.Equal(*t.Cluster, opt) {
		return false
	}

	if s.Mode != t.Mode {
		return false
	}

	if s.Status != t.Status {
		return false
	}

	return true
}

// Diff checks if two structs of type ClusterSettings are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b ClusterSettings
//	diff := a.Diff(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b ClusterSettings
//	diff := a.Diff(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s ClusterSettings) Diff(t ClusterSettings, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})
	if s.BootstrapKey != t.BootstrapKey {
		diff["BootstrapKey"] = []interface{}{s.BootstrapKey, t.BootstrapKey}
	}

	if !s.Cluster.Equal(*t.Cluster, opt) {
		diff["Cluster"] = []interface{}{ValueOrNil(s.Cluster), ValueOrNil(t.Cluster)}
	}

	if s.Mode != t.Mode {
		diff["Mode"] = []interface{}{s.Mode, t.Mode}
	}

	if s.Status != t.Status {
		diff["Status"] = []interface{}{s.Status, t.Status}
	}

	return diff
}

// Equal checks if two structs of type ClusterSettingsCluster are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b ClusterSettingsCluster
//	equal := a.Equal(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b ClusterSettingsCluster
//	equal := a.Equal(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s ClusterSettingsCluster) Equal(t ClusterSettingsCluster, opts ...Options) bool {
	opt := getOptions(opts...)

	if !CheckSameNilAndLen(s.ClusterLogTargets, t.ClusterLogTargets, opt) {
		return false
	} else {
		for i := range s.ClusterLogTargets {
			if !s.ClusterLogTargets[i].Equal(*t.ClusterLogTargets[i], opt) {
				return false
			}
		}
	}

	if s.Address != t.Address {
		return false
	}

	if s.APIBasePath != t.APIBasePath {
		return false
	}

	if s.ClusterID != t.ClusterID {
		return false
	}

	if s.Description != t.Description {
		return false
	}

	if s.Name != t.Name {
		return false
	}

	if !equalPointers(s.Port, t.Port) {
		return false
	}

	return true
}

// Diff checks if two structs of type ClusterSettingsCluster are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b ClusterSettingsCluster
//	diff := a.Diff(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b ClusterSettingsCluster
//	diff := a.Diff(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s ClusterSettingsCluster) Diff(t ClusterSettingsCluster, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})
	if !CheckSameNilAndLen(s.ClusterLogTargets, t.ClusterLogTargets, opt) {
		diff["ClusterLogTargets"] = []interface{}{s.ClusterLogTargets, t.ClusterLogTargets}
	} else {
		diff2 := make(map[string][]interface{})
		for i := range s.ClusterLogTargets {
			if !s.ClusterLogTargets[i].Equal(*t.ClusterLogTargets[i], opt) {
				diffSub := s.ClusterLogTargets[i].Diff(*t.ClusterLogTargets[i], opt)
				if len(diffSub) > 0 {
					diff2[strconv.Itoa(i)] = []interface{}{diffSub}
				}
			}
		}
		if len(diff2) > 0 {
			diff["ClusterLogTargets"] = []interface{}{diff2}
		}
	}

	if s.Address != t.Address {
		diff["Address"] = []interface{}{s.Address, t.Address}
	}

	if s.APIBasePath != t.APIBasePath {
		diff["APIBasePath"] = []interface{}{s.APIBasePath, t.APIBasePath}
	}

	if s.ClusterID != t.ClusterID {
		diff["ClusterID"] = []interface{}{s.ClusterID, t.ClusterID}
	}

	if s.Description != t.Description {
		diff["Description"] = []interface{}{s.Description, t.Description}
	}

	if s.Name != t.Name {
		diff["Name"] = []interface{}{s.Name, t.Name}
	}

	if !equalPointers(s.Port, t.Port) {
		diff["Port"] = []interface{}{ValueOrNil(s.Port), ValueOrNil(t.Port)}
	}

	return diff
}

// Equal checks if two structs of type ClusterLogTarget are equal
//
//	var a, b ClusterLogTarget
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s ClusterLogTarget) Equal(t ClusterLogTarget, opts ...Options) bool {
	if !equalPointers(s.Address, t.Address) {
		return false
	}

	if s.LogFormat != t.LogFormat {
		return false
	}

	if !equalPointers(s.Port, t.Port) {
		return false
	}

	if !equalPointers(s.Protocol, t.Protocol) {
		return false
	}

	return true
}

// Diff checks if two structs of type ClusterLogTarget are equal
//
//	var a, b ClusterLogTarget
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s ClusterLogTarget) Diff(t ClusterLogTarget, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if !equalPointers(s.Address, t.Address) {
		diff["Address"] = []interface{}{ValueOrNil(s.Address), ValueOrNil(t.Address)}
	}

	if s.LogFormat != t.LogFormat {
		diff["LogFormat"] = []interface{}{s.LogFormat, t.LogFormat}
	}

	if !equalPointers(s.Port, t.Port) {
		diff["Port"] = []interface{}{ValueOrNil(s.Port), ValueOrNil(t.Port)}
	}

	if !equalPointers(s.Protocol, t.Protocol) {
		diff["Protocol"] = []interface{}{ValueOrNil(s.Protocol), ValueOrNil(t.Protocol)}
	}

	return diff
}
