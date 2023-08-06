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

	if !s.DefaultBind.Equal(*t.DefaultBind, opt) {
		return false
	}

	if !s.DefaultServer.Equal(*t.DefaultServer, opt) {
		return false
	}

	if s.Disabled != t.Disabled {
		return false
	}

	if s.Enabled != t.Enabled {
		return false
	}

	if s.Name != t.Name {
		return false
	}

	if s.Shards != t.Shards {
		return false
	}

	return true
}

// Diff checks if two structs of type PeerSection are equal
//
// By default empty arrays, maps and slices are equal to nil:
//
//	var a, b PeerSection
//	diff := a.Diff(b)
//
// For more advanced use case you can configure the options (default values are shown):
//
//	var a, b PeerSection
//	equal := a.Diff(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s PeerSection) Diff(t PeerSection, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})
	if !s.DefaultBind.Equal(*t.DefaultBind, opt) {
		diff["DefaultBind"] = []interface{}{s.DefaultBind, t.DefaultBind}
	}

	if !s.DefaultServer.Equal(*t.DefaultServer, opt) {
		diff["DefaultServer"] = []interface{}{s.DefaultServer, t.DefaultServer}
	}

	if s.Disabled != t.Disabled {
		diff["Disabled"] = []interface{}{s.Disabled, t.Disabled}
	}

	if s.Enabled != t.Enabled {
		diff["Enabled"] = []interface{}{s.Enabled, t.Enabled}
	}

	if s.Name != t.Name {
		diff["Name"] = []interface{}{s.Name, t.Name}
	}

	if s.Shards != t.Shards {
		diff["Shards"] = []interface{}{s.Shards, t.Shards}
	}

	return diff
}
