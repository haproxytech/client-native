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

// Equal checks if two structs of type Defaults are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b Defaults
//	equal := a.Equal(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b Defaults
//	equal := a.Equal(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s Defaults) Equal(t Defaults, opts ...Options) bool {
	opt := getOptions(opts...)

	if !s.DefaultsBase.Equal(t.DefaultsBase, opt) {
		return false
	}

	if !s.HTTPCheckList.Equal(t.HTTPCheckList, opt) {
		return false
	}

	if !s.HTTPErrorRuleList.Equal(t.HTTPErrorRuleList, opt) {
		return false
	}

	if !s.LogTargetList.Equal(t.LogTargetList, opt) {
		return false
	}

	if !s.QUICInitialRuleList.Equal(t.QUICInitialRuleList, opt) {
		return false
	}

	if !s.TCPCheckRuleList.Equal(t.TCPCheckRuleList, opt) {
		return false
	}

	return true
}

// Diff checks if two structs of type Defaults are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b Defaults
//	diff := a.Diff(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b Defaults
//	diff := a.Diff(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s Defaults) Diff(t Defaults, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})

	if !s.DefaultsBase.Equal(t.DefaultsBase, opt) {
		diff["DefaultsBase"] = []interface{}{s.DefaultsBase, t.DefaultsBase}
	}

	if !s.HTTPCheckList.Equal(t.HTTPCheckList, opt) {
		diff["HTTPCheckList"] = []interface{}{s.HTTPCheckList, t.HTTPCheckList}
	}

	if !s.HTTPErrorRuleList.Equal(t.HTTPErrorRuleList, opt) {
		diff["HTTPErrorRuleList"] = []interface{}{s.HTTPErrorRuleList, t.HTTPErrorRuleList}
	}

	if !s.LogTargetList.Equal(t.LogTargetList, opt) {
		diff["LogTargetList"] = []interface{}{s.LogTargetList, t.LogTargetList}
	}

	if !s.QUICInitialRuleList.Equal(t.QUICInitialRuleList, opt) {
		diff["QUICInitialRuleList"] = []interface{}{s.QUICInitialRuleList, t.QUICInitialRuleList}
	}

	if !s.TCPCheckRuleList.Equal(t.TCPCheckRuleList, opt) {
		diff["TCPCheckRuleList"] = []interface{}{s.TCPCheckRuleList, t.TCPCheckRuleList}
	}

	return diff
}
