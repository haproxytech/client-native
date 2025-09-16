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

// Equal checks if two structs of type RuntimeServer are equal
//
//	var a, b RuntimeServer
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s RuntimeServer) Equal(t RuntimeServer, opts ...Options) bool {
	if s.Address != t.Address {
		return false
	}

	if s.AdminState != t.AdminState {
		return false
	}

	if s.AgentAddr != t.AgentAddr {
		return false
	}

	if !equalPointers(s.AgentPort, t.AgentPort) {
		return false
	}

	if !equalPointers(s.AgentState, t.AgentState) {
		return false
	}

	if !equalPointers(s.BackendForcedID, t.BackendForcedID) {
		return false
	}

	if !equalPointers(s.BackendID, t.BackendID) {
		return false
	}

	if s.BackendName != t.BackendName {
		return false
	}

	if s.CheckAddr != t.CheckAddr {
		return false
	}

	if !equalPointers(s.CheckHealth, t.CheckHealth) {
		return false
	}

	if !equalPointers(s.CheckPort, t.CheckPort) {
		return false
	}

	if !equalPointers(s.CheckResult, t.CheckResult) {
		return false
	}

	if !equalPointers(s.CheckState, t.CheckState) {
		return false
	}

	if !equalPointers(s.CheckStatus, t.CheckStatus) {
		return false
	}

	if !equalPointers(s.ForecedID, t.ForecedID) {
		return false
	}

	if s.Fqdn != t.Fqdn {
		return false
	}

	if s.ID != t.ID {
		return false
	}

	if !equalPointers(s.Iweight, t.Iweight) {
		return false
	}

	if !equalPointers(s.LastTimeChange, t.LastTimeChange) {
		return false
	}

	if s.Name != t.Name {
		return false
	}

	if s.OperationalState != t.OperationalState {
		return false
	}

	if !equalPointers(s.Port, t.Port) {
		return false
	}

	if s.Srvrecord != t.Srvrecord {
		return false
	}

	if !equalPointers(s.UseSsl, t.UseSsl) {
		return false
	}

	if !equalPointers(s.Uweight, t.Uweight) {
		return false
	}

	return true
}

// Diff checks if two structs of type RuntimeServer are equal
//
//	var a, b RuntimeServer
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s RuntimeServer) Diff(t RuntimeServer, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if s.Address != t.Address {
		diff["Address"] = []interface{}{s.Address, t.Address}
	}

	if s.AdminState != t.AdminState {
		diff["AdminState"] = []interface{}{s.AdminState, t.AdminState}
	}

	if s.AgentAddr != t.AgentAddr {
		diff["AgentAddr"] = []interface{}{s.AgentAddr, t.AgentAddr}
	}

	if !equalPointers(s.AgentPort, t.AgentPort) {
		diff["AgentPort"] = []interface{}{ValueOrNil(s.AgentPort), ValueOrNil(t.AgentPort)}
	}

	if !equalPointers(s.AgentState, t.AgentState) {
		diff["AgentState"] = []interface{}{ValueOrNil(s.AgentState), ValueOrNil(t.AgentState)}
	}

	if !equalPointers(s.BackendForcedID, t.BackendForcedID) {
		diff["BackendForcedID"] = []interface{}{ValueOrNil(s.BackendForcedID), ValueOrNil(t.BackendForcedID)}
	}

	if !equalPointers(s.BackendID, t.BackendID) {
		diff["BackendID"] = []interface{}{ValueOrNil(s.BackendID), ValueOrNil(t.BackendID)}
	}

	if s.BackendName != t.BackendName {
		diff["BackendName"] = []interface{}{s.BackendName, t.BackendName}
	}

	if s.CheckAddr != t.CheckAddr {
		diff["CheckAddr"] = []interface{}{s.CheckAddr, t.CheckAddr}
	}

	if !equalPointers(s.CheckHealth, t.CheckHealth) {
		diff["CheckHealth"] = []interface{}{ValueOrNil(s.CheckHealth), ValueOrNil(t.CheckHealth)}
	}

	if !equalPointers(s.CheckPort, t.CheckPort) {
		diff["CheckPort"] = []interface{}{ValueOrNil(s.CheckPort), ValueOrNil(t.CheckPort)}
	}

	if !equalPointers(s.CheckResult, t.CheckResult) {
		diff["CheckResult"] = []interface{}{ValueOrNil(s.CheckResult), ValueOrNil(t.CheckResult)}
	}

	if !equalPointers(s.CheckState, t.CheckState) {
		diff["CheckState"] = []interface{}{ValueOrNil(s.CheckState), ValueOrNil(t.CheckState)}
	}

	if !equalPointers(s.CheckStatus, t.CheckStatus) {
		diff["CheckStatus"] = []interface{}{ValueOrNil(s.CheckStatus), ValueOrNil(t.CheckStatus)}
	}

	if !equalPointers(s.ForecedID, t.ForecedID) {
		diff["ForecedID"] = []interface{}{ValueOrNil(s.ForecedID), ValueOrNil(t.ForecedID)}
	}

	if s.Fqdn != t.Fqdn {
		diff["Fqdn"] = []interface{}{s.Fqdn, t.Fqdn}
	}

	if s.ID != t.ID {
		diff["ID"] = []interface{}{s.ID, t.ID}
	}

	if !equalPointers(s.Iweight, t.Iweight) {
		diff["Iweight"] = []interface{}{ValueOrNil(s.Iweight), ValueOrNil(t.Iweight)}
	}

	if !equalPointers(s.LastTimeChange, t.LastTimeChange) {
		diff["LastTimeChange"] = []interface{}{ValueOrNil(s.LastTimeChange), ValueOrNil(t.LastTimeChange)}
	}

	if s.Name != t.Name {
		diff["Name"] = []interface{}{s.Name, t.Name}
	}

	if s.OperationalState != t.OperationalState {
		diff["OperationalState"] = []interface{}{s.OperationalState, t.OperationalState}
	}

	if !equalPointers(s.Port, t.Port) {
		diff["Port"] = []interface{}{ValueOrNil(s.Port), ValueOrNil(t.Port)}
	}

	if s.Srvrecord != t.Srvrecord {
		diff["Srvrecord"] = []interface{}{s.Srvrecord, t.Srvrecord}
	}

	if !equalPointers(s.UseSsl, t.UseSsl) {
		diff["UseSsl"] = []interface{}{ValueOrNil(s.UseSsl), ValueOrNil(t.UseSsl)}
	}

	if !equalPointers(s.Uweight, t.Uweight) {
		diff["Uweight"] = []interface{}{ValueOrNil(s.Uweight), ValueOrNil(t.Uweight)}
	}

	return diff
}
