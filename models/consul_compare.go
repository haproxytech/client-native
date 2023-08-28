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

// Equal checks if two structs of type Consul are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b Consul
//	equal := a.Equal(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b Consul
//	equal := a.Equal(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s Consul) Equal(t Consul, opts ...Options) bool {
	opt := getOptions(opts...)

	if !equalPointers(s.Address, t.Address) {
		return false
	}

	if s.Defaults != t.Defaults {
		return false
	}

	if s.Description != t.Description {
		return false
	}

	if !equalPointers(s.Enabled, t.Enabled) {
		return false
	}

	if !equalPointers(s.HealthCheckPolicy, t.HealthCheckPolicy) {
		return false
	}

	if s.HealthCheckPolicyMin != t.HealthCheckPolicyMin {
		return false
	}

	if !equalPointers(s.ID, t.ID) {
		return false
	}

	if s.Name != t.Name {
		return false
	}

	if s.Namespace != t.Namespace {
		return false
	}

	if !equalPointers(s.Port, t.Port) {
		return false
	}

	if !equalPointers(s.RetryTimeout, t.RetryTimeout) {
		return false
	}

	if !equalPointers(s.ServerSlotsBase, t.ServerSlotsBase) {
		return false
	}

	if s.ServerSlotsGrowthIncrement != t.ServerSlotsGrowthIncrement {
		return false
	}

	if !equalPointers(s.ServerSlotsGrowthType, t.ServerSlotsGrowthType) {
		return false
	}

	if !equalComparableSlice(s.ServiceBlacklist, t.ServiceBlacklist, opt) {
		return false
	}

	if !equalComparableSlice(s.ServiceWhitelist, t.ServiceWhitelist, opt) {
		return false
	}

	if !equalComparableSlice(s.ServiceAllowlist, t.ServiceAllowlist, opt) {
		return false
	}

	if !equalComparableSlice(s.ServiceDenylist, t.ServiceDenylist, opt) {
		return false
	}

	if s.ServiceNameRegexp != t.ServiceNameRegexp {
		return false
	}

	if s.Token != t.Token {
		return false
	}

	return true
}

// Diff checks if two structs of type Consul are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b Consul
//	diff := a.Diff(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b Consul
//	diff := a.Diff(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s Consul) Diff(t Consul, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})
	if !equalPointers(s.Address, t.Address) {
		diff["Address"] = []interface{}{ValueOrNil(s.Address), ValueOrNil(t.Address)}
	}

	if s.Defaults != t.Defaults {
		diff["Defaults"] = []interface{}{s.Defaults, t.Defaults}
	}

	if s.Description != t.Description {
		diff["Description"] = []interface{}{s.Description, t.Description}
	}

	if !equalPointers(s.Enabled, t.Enabled) {
		diff["Enabled"] = []interface{}{ValueOrNil(s.Enabled), ValueOrNil(t.Enabled)}
	}

	if !equalPointers(s.HealthCheckPolicy, t.HealthCheckPolicy) {
		diff["HealthCheckPolicy"] = []interface{}{ValueOrNil(s.HealthCheckPolicy), ValueOrNil(t.HealthCheckPolicy)}
	}

	if s.HealthCheckPolicyMin != t.HealthCheckPolicyMin {
		diff["HealthCheckPolicyMin"] = []interface{}{s.HealthCheckPolicyMin, t.HealthCheckPolicyMin}
	}

	if !equalPointers(s.ID, t.ID) {
		diff["ID"] = []interface{}{ValueOrNil(s.ID), ValueOrNil(t.ID)}
	}

	if s.Name != t.Name {
		diff["Name"] = []interface{}{s.Name, t.Name}
	}

	if s.Namespace != t.Namespace {
		diff["Namespace"] = []interface{}{s.Namespace, t.Namespace}
	}

	if !equalPointers(s.Port, t.Port) {
		diff["Port"] = []interface{}{ValueOrNil(s.Port), ValueOrNil(t.Port)}
	}

	if !equalPointers(s.RetryTimeout, t.RetryTimeout) {
		diff["RetryTimeout"] = []interface{}{ValueOrNil(s.RetryTimeout), ValueOrNil(t.RetryTimeout)}
	}

	if !equalPointers(s.ServerSlotsBase, t.ServerSlotsBase) {
		diff["ServerSlotsBase"] = []interface{}{ValueOrNil(s.ServerSlotsBase), ValueOrNil(t.ServerSlotsBase)}
	}

	if s.ServerSlotsGrowthIncrement != t.ServerSlotsGrowthIncrement {
		diff["ServerSlotsGrowthIncrement"] = []interface{}{s.ServerSlotsGrowthIncrement, t.ServerSlotsGrowthIncrement}
	}

	if !equalPointers(s.ServerSlotsGrowthType, t.ServerSlotsGrowthType) {
		diff["ServerSlotsGrowthType"] = []interface{}{ValueOrNil(s.ServerSlotsGrowthType), ValueOrNil(t.ServerSlotsGrowthType)}
	}

	if !equalComparableSlice(s.ServiceBlacklist, t.ServiceBlacklist, opt) {
		diff["ServiceBlacklist"] = []interface{}{s.ServiceBlacklist, t.ServiceBlacklist}
	}

	if !equalComparableSlice(s.ServiceWhitelist, t.ServiceWhitelist, opt) {
		diff["ServiceWhitelist"] = []interface{}{s.ServiceWhitelist, t.ServiceWhitelist}
	}

	if !equalComparableSlice(s.ServiceAllowlist, t.ServiceAllowlist, opt) {
		diff["ServiceAllowlist"] = []interface{}{s.ServiceAllowlist, t.ServiceAllowlist}
	}

	if !equalComparableSlice(s.ServiceDenylist, t.ServiceDenylist, opt) {
		diff["ServiceDenylist"] = []interface{}{s.ServiceDenylist, t.ServiceDenylist}
	}

	if s.ServiceNameRegexp != t.ServiceNameRegexp {
		diff["ServiceNameRegexp"] = []interface{}{s.ServiceNameRegexp, t.ServiceNameRegexp}
	}

	if s.Token != t.Token {
		diff["Token"] = []interface{}{s.Token, t.Token}
	}

	return diff
}
