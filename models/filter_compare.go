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

// Equal checks if two structs of type Filter are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b Filter
//	equal := a.Equal(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b Filter
//	equal := a.Equal(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s Filter) Equal(t Filter, opts ...Options) bool {
	opt := getOptions(opts...)

	if s.AppName != t.AppName {
		return false
	}

	if s.BandwidthLimitName != t.BandwidthLimitName {
		return false
	}

	if s.CacheName != t.CacheName {
		return false
	}

	if s.DefaultLimit != t.DefaultLimit {
		return false
	}

	if s.DefaultPeriod != t.DefaultPeriod {
		return false
	}

	if s.Key != t.Key {
		return false
	}

	if s.Limit != t.Limit {
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

	if s.MinSize != t.MinSize {
		return false
	}

	if s.SpoeConfig != t.SpoeConfig {
		return false
	}

	if s.SpoeEngine != t.SpoeEngine {
		return false
	}

	if s.Table != t.Table {
		return false
	}

	if s.TraceHexdump != t.TraceHexdump {
		return false
	}

	if s.TraceName != t.TraceName {
		return false
	}

	if s.TraceRndForwarding != t.TraceRndForwarding {
		return false
	}

	if s.TraceRndParsing != t.TraceRndParsing {
		return false
	}

	if s.Type != t.Type {
		return false
	}

	return true
}

// Diff checks if two structs of type Filter are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b Filter
//	diff := a.Diff(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b Filter
//	diff := a.Diff(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s Filter) Diff(t Filter, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})
	if s.AppName != t.AppName {
		diff["AppName"] = []interface{}{s.AppName, t.AppName}
	}

	if s.BandwidthLimitName != t.BandwidthLimitName {
		diff["BandwidthLimitName"] = []interface{}{s.BandwidthLimitName, t.BandwidthLimitName}
	}

	if s.CacheName != t.CacheName {
		diff["CacheName"] = []interface{}{s.CacheName, t.CacheName}
	}

	if s.DefaultLimit != t.DefaultLimit {
		diff["DefaultLimit"] = []interface{}{s.DefaultLimit, t.DefaultLimit}
	}

	if s.DefaultPeriod != t.DefaultPeriod {
		diff["DefaultPeriod"] = []interface{}{s.DefaultPeriod, t.DefaultPeriod}
	}

	if s.Key != t.Key {
		diff["Key"] = []interface{}{s.Key, t.Key}
	}

	if s.Limit != t.Limit {
		diff["Limit"] = []interface{}{s.Limit, t.Limit}
	}

	if !CheckSameNilAndLenMap[string](s.Metadata, t.Metadata, opt) {
		diff["Metadata"] = []interface{}{s.Metadata, t.Metadata}
	}

	for k, v := range s.Metadata {
		if !reflect.DeepEqual(t.Metadata[k], v) {
			diff["Metadata"] = []interface{}{s.Metadata, t.Metadata}
		}
	}

	if s.MinSize != t.MinSize {
		diff["MinSize"] = []interface{}{s.MinSize, t.MinSize}
	}

	if s.SpoeConfig != t.SpoeConfig {
		diff["SpoeConfig"] = []interface{}{s.SpoeConfig, t.SpoeConfig}
	}

	if s.SpoeEngine != t.SpoeEngine {
		diff["SpoeEngine"] = []interface{}{s.SpoeEngine, t.SpoeEngine}
	}

	if s.Table != t.Table {
		diff["Table"] = []interface{}{s.Table, t.Table}
	}

	if s.TraceHexdump != t.TraceHexdump {
		diff["TraceHexdump"] = []interface{}{s.TraceHexdump, t.TraceHexdump}
	}

	if s.TraceName != t.TraceName {
		diff["TraceName"] = []interface{}{s.TraceName, t.TraceName}
	}

	if s.TraceRndForwarding != t.TraceRndForwarding {
		diff["TraceRndForwarding"] = []interface{}{s.TraceRndForwarding, t.TraceRndForwarding}
	}

	if s.TraceRndParsing != t.TraceRndParsing {
		diff["TraceRndParsing"] = []interface{}{s.TraceRndParsing, t.TraceRndParsing}
	}

	if s.Type != t.Type {
		diff["Type"] = []interface{}{s.Type, t.Type}
	}

	return diff
}
