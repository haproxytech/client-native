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

import "reflect"

// Equal checks if two structs of type PeerSectionBase are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b PeerSectionBase
//	equal := a.Equal(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b PeerSectionBase
//	equal := a.Equal(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s PeerSectionBase) Equal(t PeerSectionBase, opts ...Options) bool {
	opt := getOptions(opts...)

	if s.DefaultBind == nil || t.DefaultBind == nil {
		if s.DefaultBind != nil || t.DefaultBind != nil {
			if opt.NilSameAsEmpty {
				empty := &DefaultBind{}
				if s.DefaultBind == nil {
					if !(t.DefaultBind.Equal(*empty)) {
						return false
					}
				}
				if t.DefaultBind == nil {
					if !(s.DefaultBind.Equal(*empty)) {
						return false
					}
				}
			} else {
				return false
			}
		}
	} else if !s.DefaultBind.Equal(*t.DefaultBind, opt) {
		return false
	}

	if s.DefaultServer == nil || t.DefaultServer == nil {
		if s.DefaultServer != nil || t.DefaultServer != nil {
			if opt.NilSameAsEmpty {
				empty := &DefaultServer{}
				if s.DefaultServer == nil {
					if !(t.DefaultServer.Equal(*empty)) {
						return false
					}
				}
				if t.DefaultServer == nil {
					if !(s.DefaultServer.Equal(*empty)) {
						return false
					}
				}
			} else {
				return false
			}
		}
	} else if !s.DefaultServer.Equal(*t.DefaultServer, opt) {
		return false
	}

	if s.Disabled != t.Disabled {
		return false
	}

	if s.Enabled != t.Enabled {
		return false
	}

	if !CheckSameNilAndLenMap[string](s.Metadata, t.Metadata, opt) {
		return false
	}

	for k, v := range s.Metadata {
		if !reflect.DeepEqual(t.Metadata[k], v) {
			return false
		}
	}

	if s.Name != t.Name {
		return false
	}

	if s.Shards != t.Shards {
		return false
	}

	return true
}

// Diff checks if two structs of type PeerSectionBase are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b PeerSectionBase
//	diff := a.Diff(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b PeerSectionBase
//	diff := a.Diff(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s PeerSectionBase) Diff(t PeerSectionBase, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})

	if s.DefaultBind == nil || t.DefaultBind == nil {
		if s.DefaultBind != nil || t.DefaultBind != nil {
			if opt.NilSameAsEmpty {
				empty := &DefaultBind{}
				if s.DefaultBind == nil {
					if !(t.DefaultBind.Equal(*empty)) {
						diff["DefaultBind"] = []interface{}{ValueOrNil(s.DefaultBind), ValueOrNil(t.DefaultBind)}
					}
				}
				if t.DefaultBind == nil {
					if !(s.DefaultBind.Equal(*empty)) {
						diff["DefaultBind"] = []interface{}{ValueOrNil(s.DefaultBind), ValueOrNil(t.DefaultBind)}
					}
				}
			} else {
				diff["DefaultBind"] = []interface{}{ValueOrNil(s.DefaultBind), ValueOrNil(t.DefaultBind)}
			}
		}
	} else if !s.DefaultBind.Equal(*t.DefaultBind, opt) {
		diff["DefaultBind"] = []interface{}{ValueOrNil(s.DefaultBind), ValueOrNil(t.DefaultBind)}
	}

	if s.DefaultServer == nil || t.DefaultServer == nil {
		if s.DefaultServer != nil || t.DefaultServer != nil {
			if opt.NilSameAsEmpty {
				empty := &DefaultServer{}
				if s.DefaultServer == nil {
					if !(t.DefaultServer.Equal(*empty)) {
						diff["DefaultServer"] = []interface{}{ValueOrNil(s.DefaultServer), ValueOrNil(t.DefaultServer)}
					}
				}
				if t.DefaultServer == nil {
					if !(s.DefaultServer.Equal(*empty)) {
						diff["DefaultServer"] = []interface{}{ValueOrNil(s.DefaultServer), ValueOrNil(t.DefaultServer)}
					}
				}
			} else {
				diff["DefaultServer"] = []interface{}{ValueOrNil(s.DefaultServer), ValueOrNil(t.DefaultServer)}
			}
		}
	} else if !s.DefaultServer.Equal(*t.DefaultServer, opt) {
		diff["DefaultServer"] = []interface{}{ValueOrNil(s.DefaultServer), ValueOrNil(t.DefaultServer)}
	}

	if s.Disabled != t.Disabled {
		diff["Disabled"] = []interface{}{s.Disabled, t.Disabled}
	}

	if s.Enabled != t.Enabled {
		diff["Enabled"] = []interface{}{s.Enabled, t.Enabled}
	}

	if !CheckSameNilAndLenMap[string](s.Metadata, t.Metadata, opt) {
		diff["Metadata"] = []interface{}{s.Metadata, t.Metadata}
	}

	for k, v := range s.Metadata {
		if !reflect.DeepEqual(t.Metadata[k], v) {
			diff["Metadata"] = []interface{}{s.Metadata, t.Metadata}
		}
	}

	if s.Name != t.Name {
		diff["Name"] = []interface{}{s.Name, t.Name}
	}

	if s.Shards != t.Shards {
		diff["Shards"] = []interface{}{s.Shards, t.Shards}
	}

	return diff
}
