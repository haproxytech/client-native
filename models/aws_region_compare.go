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

// Equal checks if two structs of type AwsRegion are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b AwsRegion
//	equal := a.Equal(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b AwsRegion
//	equal := a.Equal(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s AwsRegion) Equal(t AwsRegion, opts ...Options) bool {
	opt := getOptions(opts...)

	if s.AccessKeyID != t.AccessKeyID {
		return false
	}

	if !CheckSameNilAndLen(s.Allowlist, t.Allowlist, opt) {
		return false
	}
	for i := range s.Allowlist {
		if !s.Allowlist[i].Equal(*t.Allowlist[i], opt) {
			return false
		}
	}

	if !CheckSameNilAndLen(s.Denylist, t.Denylist, opt) {
		return false
	}
	for i := range s.Denylist {
		if !s.Denylist[i].Equal(*t.Denylist[i], opt) {
			return false
		}
	}

	if s.Description != t.Description {
		return false
	}

	if !equalPointers(s.Enabled, t.Enabled) {
		return false
	}

	if !equalPointers(s.ID, t.ID) {
		return false
	}

	if !equalPointers(s.IPV4Address, t.IPV4Address) {
		return false
	}

	if !equalPointers(s.Name, t.Name) {
		return false
	}

	if !equalPointers(s.Region, t.Region) {
		return false
	}

	if !equalPointers(s.RetryTimeout, t.RetryTimeout) {
		return false
	}

	if s.SecretAccessKey != t.SecretAccessKey {
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

	return true
}

// Diff checks if two structs of type AwsRegion are equal
//
// By default empty arrays, maps and slices are equal to nil:
//
//	var a, b AwsRegion
//	diff := a.Diff(b)
//
// For more advanced use case you can configure the options (default values are shown):
//
//	var a, b AwsRegion
//	equal := a.Diff(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s AwsRegion) Diff(t AwsRegion, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})
	if s.AccessKeyID != t.AccessKeyID {
		diff["AccessKeyID"] = []interface{}{s.AccessKeyID, t.AccessKeyID}
	}

	if !CheckSameNilAndLen(s.Allowlist, t.Allowlist, opt) {
		diff["Allowlist"] = []interface{}{s.Allowlist, t.Allowlist}
	} else {
		diff2 := make(map[string][]interface{})
		for i := range s.Allowlist {
			diffSub := s.Allowlist[i].Diff(*t.Allowlist[i], opt)
			if len(diffSub) > 0 {
				diff2[strconv.Itoa(i)] = []interface{}{diffSub}
			}
		}
		if len(diff2) > 0 {
			diff["Allowlist"] = []interface{}{diff2}
		}
	}

	if !CheckSameNilAndLen(s.Denylist, t.Denylist, opt) {
		diff["Denylist"] = []interface{}{s.Denylist, t.Denylist}
	} else {
		diff2 := make(map[string][]interface{})
		for i := range s.Denylist {
			diffSub := s.Denylist[i].Diff(*t.Denylist[i], opt)
			if len(diffSub) > 0 {
				diff2[strconv.Itoa(i)] = []interface{}{diffSub}
			}
		}
		if len(diff2) > 0 {
			diff["Denylist"] = []interface{}{diff2}
		}
	}

	if s.Description != t.Description {
		diff["Description"] = []interface{}{s.Description, t.Description}
	}

	if !equalPointers(s.Enabled, t.Enabled) {
		diff["Enabled"] = []interface{}{ValueOrNil(s.Enabled), ValueOrNil(t.Enabled)}
	}

	if !equalPointers(s.ID, t.ID) {
		diff["ID"] = []interface{}{ValueOrNil(s.ID), ValueOrNil(t.ID)}
	}

	if !equalPointers(s.IPV4Address, t.IPV4Address) {
		diff["IPV4Address"] = []interface{}{ValueOrNil(s.IPV4Address), ValueOrNil(t.IPV4Address)}
	}

	if !equalPointers(s.Name, t.Name) {
		diff["Name"] = []interface{}{ValueOrNil(s.Name), ValueOrNil(t.Name)}
	}

	if !equalPointers(s.Region, t.Region) {
		diff["Region"] = []interface{}{ValueOrNil(s.Region), ValueOrNil(t.Region)}
	}

	if !equalPointers(s.RetryTimeout, t.RetryTimeout) {
		diff["RetryTimeout"] = []interface{}{ValueOrNil(s.RetryTimeout), ValueOrNil(t.RetryTimeout)}
	}

	if s.SecretAccessKey != t.SecretAccessKey {
		diff["SecretAccessKey"] = []interface{}{s.SecretAccessKey, t.SecretAccessKey}
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

	return diff
}
