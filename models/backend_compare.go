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

// Equal checks if two structs of type Backend are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b Backend
//	equal := a.Equal(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b Backend
//	equal := a.Equal(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s Backend) Equal(t Backend, opts ...Options) bool {
	opt := getOptions(opts...)

	if !s.BackendBase.Equal(t.BackendBase, opt) {
		return false
	}

	if !s.ACLList.Equal(t.ACLList, opt) {
		return false
	}

	if !s.FilterList.Equal(t.FilterList, opt) {
		return false
	}

	if !s.HTTPAfterResponseRuleList.Equal(t.HTTPAfterResponseRuleList, opt) {
		return false
	}

	if !s.HTTPCheckList.Equal(t.HTTPCheckList, opt) {
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

	if !s.ServerSwitchingRuleList.Equal(t.ServerSwitchingRuleList, opt) {
		return false
	}

	if !s.StickRuleList.Equal(t.StickRuleList, opt) {
		return false
	}

	if !s.TCPCheckRuleList.Equal(t.TCPCheckRuleList, opt) {
		return false
	}

	if !s.TCPRequestRuleList.Equal(t.TCPRequestRuleList, opt) {
		return false
	}

	if !s.TCPResponseRuleList.Equal(t.TCPResponseRuleList, opt) {
		return false
	}

	if !CheckSameNilAndLenMap[string, ServerTemplate](s.ServerTemplates, t.ServerTemplates, opt) {
		return false
	}

	for k, v := range s.ServerTemplates {
		if !t.ServerTemplates[k].Equal(v, opt) {
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

	return true
}

// Diff checks if two structs of type Backend are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b Backend
//	diff := a.Diff(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b Backend
//	diff := a.Diff(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s Backend) Diff(t Backend, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})

	if !s.BackendBase.Equal(t.BackendBase, opt) {
		diff["BackendBase"] = []interface{}{s.BackendBase, t.BackendBase}
	}

	if !s.ACLList.Equal(t.ACLList, opt) {
		diff["ACLList"] = []interface{}{s.ACLList, t.ACLList}
	}

	if !s.FilterList.Equal(t.FilterList, opt) {
		diff["FilterList"] = []interface{}{s.FilterList, t.FilterList}
	}

	if !s.HTTPAfterResponseRuleList.Equal(t.HTTPAfterResponseRuleList, opt) {
		diff["HTTPAfterResponseRuleList"] = []interface{}{s.HTTPAfterResponseRuleList, t.HTTPAfterResponseRuleList}
	}

	if !s.HTTPCheckList.Equal(t.HTTPCheckList, opt) {
		diff["HTTPCheckList"] = []interface{}{s.HTTPCheckList, t.HTTPCheckList}
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

	if !s.ServerSwitchingRuleList.Equal(t.ServerSwitchingRuleList, opt) {
		diff["ServerSwitchingRuleList"] = []interface{}{s.ServerSwitchingRuleList, t.ServerSwitchingRuleList}
	}

	if !s.StickRuleList.Equal(t.StickRuleList, opt) {
		diff["StickRuleList"] = []interface{}{s.StickRuleList, t.StickRuleList}
	}

	if !s.TCPCheckRuleList.Equal(t.TCPCheckRuleList, opt) {
		diff["TCPCheckRuleList"] = []interface{}{s.TCPCheckRuleList, t.TCPCheckRuleList}
	}

	if !s.TCPRequestRuleList.Equal(t.TCPRequestRuleList, opt) {
		diff["TCPRequestRuleList"] = []interface{}{s.TCPRequestRuleList, t.TCPRequestRuleList}
	}

	if !s.TCPResponseRuleList.Equal(t.TCPResponseRuleList, opt) {
		diff["TCPResponseRuleList"] = []interface{}{s.TCPResponseRuleList, t.TCPResponseRuleList}
	}

	if !CheckSameNilAndLenMap[string, ServerTemplate](s.ServerTemplates, t.ServerTemplates, opt) {
		diff["ServerTemplates"] = []interface{}{s.ServerTemplates, t.ServerTemplates}
	}

	for k, v := range s.ServerTemplates {
		if !t.ServerTemplates[k].Equal(v, opt) {
			diff["ServerTemplates"] = []interface{}{s.ServerTemplates, t.ServerTemplates}
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

	return diff
}
