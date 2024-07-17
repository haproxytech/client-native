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

// Equal checks if two structs of type Frontend are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b Frontend
//	equal := a.Equal(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b Frontend
//	equal := a.Equal(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s Frontend) Equal(t Frontend, opts ...Options) bool {
	opt := getOptions(opts...)

	if !s.FrontendBase.Equal(t.FrontendBase, opt) {
		return false
	}

	if !s.ACLList.Equal(t.ACLList, opt) {
		return false
	}

	if !s.BackendSwitchingRuleList.Equal(t.BackendSwitchingRuleList, opt) {
		return false
	}

	if !s.CaptureList.Equal(t.CaptureList, opt) {
		return false
	}

	if !s.FilterList.Equal(t.FilterList, opt) {
		return false
	}

	if !s.HTTPAfterResponseRuleList.Equal(t.HTTPAfterResponseRuleList, opt) {
		return false
	}

	if !s.HTTPErrorRuleList.Equal(t.HTTPErrorRuleList, opt) {
		return false
	}

	if !s.HTTPRequestRuleList.Equal(t.HTTPRequestRuleList, opt) {
		return false
	}

	if !s.HTTPResponseRuleList.Equal(t.HTTPResponseRuleList, opt) {
		return false
	}

	if !s.LogTargetList.Equal(t.LogTargetList, opt) {
		return false
	}

	if !s.TCPRequestRuleList.Equal(t.TCPRequestRuleList, opt) {
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

	return true
}

// Diff checks if two structs of type Frontend are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b Frontend
//	diff := a.Diff(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b Frontend
//	diff := a.Diff(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s Frontend) Diff(t Frontend, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})

	if !s.FrontendBase.Equal(t.FrontendBase, opt) {
		diff["FrontendBase"] = []interface{}{s.FrontendBase, t.FrontendBase}
	}

	if !s.ACLList.Equal(t.ACLList, opt) {
		diff["ACLList"] = []interface{}{s.ACLList, t.ACLList}
	}

	if !s.BackendSwitchingRuleList.Equal(t.BackendSwitchingRuleList, opt) {
		diff["BackendSwitchingRuleList"] = []interface{}{s.BackendSwitchingRuleList, t.BackendSwitchingRuleList}
	}

	if !s.CaptureList.Equal(t.CaptureList, opt) {
		diff["CaptureList"] = []interface{}{s.CaptureList, t.CaptureList}
	}

	if !s.FilterList.Equal(t.FilterList, opt) {
		diff["FilterList"] = []interface{}{s.FilterList, t.FilterList}
	}

	if !s.HTTPAfterResponseRuleList.Equal(t.HTTPAfterResponseRuleList, opt) {
		diff["HTTPAfterResponseRuleList"] = []interface{}{s.HTTPAfterResponseRuleList, t.HTTPAfterResponseRuleList}
	}

	if !s.HTTPErrorRuleList.Equal(t.HTTPErrorRuleList, opt) {
		diff["HTTPErrorRuleList"] = []interface{}{s.HTTPErrorRuleList, t.HTTPErrorRuleList}
	}

	if !s.HTTPRequestRuleList.Equal(t.HTTPRequestRuleList, opt) {
		diff["HTTPRequestRuleList"] = []interface{}{s.HTTPRequestRuleList, t.HTTPRequestRuleList}
	}

	if !s.HTTPResponseRuleList.Equal(t.HTTPResponseRuleList, opt) {
		diff["HTTPResponseRuleList"] = []interface{}{s.HTTPResponseRuleList, t.HTTPResponseRuleList}
	}

	if !s.LogTargetList.Equal(t.LogTargetList, opt) {
		diff["LogTargetList"] = []interface{}{s.LogTargetList, t.LogTargetList}
	}

	if !s.TCPRequestRuleList.Equal(t.TCPRequestRuleList, opt) {
		diff["TCPRequestRuleList"] = []interface{}{s.TCPRequestRuleList, t.TCPRequestRuleList}
	}

	if !CheckSameNilAndLenMap[string, Bind](s.Binds, t.Binds, opt) {
		diff["Binds"] = []interface{}{s.Binds, t.Binds}
	}

	for k, v := range s.Binds {
		if !t.Binds[k].Equal(v, opt) {
			diff["Binds"] = []interface{}{s.Binds, t.Binds}
		}
	}

	return diff
}
