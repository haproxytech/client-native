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

// Equal checks if two structs of type SpoeAgent are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b SpoeAgent
//	equal := a.Equal(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b SpoeAgent
//	equal := a.Equal(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s SpoeAgent) Equal(t SpoeAgent, opts ...Options) bool {
	opt := getOptions(opts...)

	if s.Async != t.Async {
		return false
	}

	if s.ContinueOnError != t.ContinueOnError {
		return false
	}

	if s.DontlogNormal != t.DontlogNormal {
		return false
	}

	if s.EngineName != t.EngineName {
		return false
	}

	if s.ForceSetVar != t.ForceSetVar {
		return false
	}

	if s.Groups != t.Groups {
		return false
	}

	if s.HelloTimeout != t.HelloTimeout {
		return false
	}

	if s.IdleTimeout != t.IdleTimeout {
		return false
	}

	if !s.Log.Equal(t.Log, opt) {
		return false
	}

	if s.MaxFrameSize != t.MaxFrameSize {
		return false
	}

	if s.MaxWaitingFrames != t.MaxWaitingFrames {
		return false
	}

	if s.Maxconnrate != t.Maxconnrate {
		return false
	}

	if s.Maxerrrate != t.Maxerrrate {
		return false
	}

	if s.Messages != t.Messages {
		return false
	}

	if !equalPointers(s.Name, t.Name) {
		return false
	}

	if s.OptionSetOnError != t.OptionSetOnError {
		return false
	}

	if s.OptionSetProcessTime != t.OptionSetProcessTime {
		return false
	}

	if s.OptionSetTotalTime != t.OptionSetTotalTime {
		return false
	}

	if s.OptionVarPrefix != t.OptionVarPrefix {
		return false
	}

	if s.Pipelining != t.Pipelining {
		return false
	}

	if s.ProcessingTimeout != t.ProcessingTimeout {
		return false
	}

	if s.RegisterVarNames != t.RegisterVarNames {
		return false
	}

	if s.SendFragPayload != t.SendFragPayload {
		return false
	}

	if s.UseBackend != t.UseBackend {
		return false
	}

	return true
}

// Diff checks if two structs of type SpoeAgent are equal
//
// By default empty arrays, maps and slices are equal to nil:
//
//	var a, b SpoeAgent
//	diff := a.Diff(b)
//
// For more advanced use case you can configure the options (default values are shown):
//
//	var a, b SpoeAgent
//	equal := a.Diff(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s SpoeAgent) Diff(t SpoeAgent, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})
	if s.Async != t.Async {
		diff["Async"] = []interface{}{s.Async, t.Async}
	}

	if s.ContinueOnError != t.ContinueOnError {
		diff["ContinueOnError"] = []interface{}{s.ContinueOnError, t.ContinueOnError}
	}

	if s.DontlogNormal != t.DontlogNormal {
		diff["DontlogNormal"] = []interface{}{s.DontlogNormal, t.DontlogNormal}
	}

	if s.EngineName != t.EngineName {
		diff["EngineName"] = []interface{}{s.EngineName, t.EngineName}
	}

	if s.ForceSetVar != t.ForceSetVar {
		diff["ForceSetVar"] = []interface{}{s.ForceSetVar, t.ForceSetVar}
	}

	if s.Groups != t.Groups {
		diff["Groups"] = []interface{}{s.Groups, t.Groups}
	}

	if s.HelloTimeout != t.HelloTimeout {
		diff["HelloTimeout"] = []interface{}{s.HelloTimeout, t.HelloTimeout}
	}

	if s.IdleTimeout != t.IdleTimeout {
		diff["IdleTimeout"] = []interface{}{s.IdleTimeout, t.IdleTimeout}
	}

	if !s.Log.Equal(t.Log, opt) {
		diff["Log"] = []interface{}{s.Log, t.Log}
	}

	if s.MaxFrameSize != t.MaxFrameSize {
		diff["MaxFrameSize"] = []interface{}{s.MaxFrameSize, t.MaxFrameSize}
	}

	if s.MaxWaitingFrames != t.MaxWaitingFrames {
		diff["MaxWaitingFrames"] = []interface{}{s.MaxWaitingFrames, t.MaxWaitingFrames}
	}

	if s.Maxconnrate != t.Maxconnrate {
		diff["Maxconnrate"] = []interface{}{s.Maxconnrate, t.Maxconnrate}
	}

	if s.Maxerrrate != t.Maxerrrate {
		diff["Maxerrrate"] = []interface{}{s.Maxerrrate, t.Maxerrrate}
	}

	if s.Messages != t.Messages {
		diff["Messages"] = []interface{}{s.Messages, t.Messages}
	}

	if !equalPointers(s.Name, t.Name) {
		diff["Name"] = []interface{}{s.Name, t.Name}
	}

	if s.OptionSetOnError != t.OptionSetOnError {
		diff["OptionSetOnError"] = []interface{}{s.OptionSetOnError, t.OptionSetOnError}
	}

	if s.OptionSetProcessTime != t.OptionSetProcessTime {
		diff["OptionSetProcessTime"] = []interface{}{s.OptionSetProcessTime, t.OptionSetProcessTime}
	}

	if s.OptionSetTotalTime != t.OptionSetTotalTime {
		diff["OptionSetTotalTime"] = []interface{}{s.OptionSetTotalTime, t.OptionSetTotalTime}
	}

	if s.OptionVarPrefix != t.OptionVarPrefix {
		diff["OptionVarPrefix"] = []interface{}{s.OptionVarPrefix, t.OptionVarPrefix}
	}

	if s.Pipelining != t.Pipelining {
		diff["Pipelining"] = []interface{}{s.Pipelining, t.Pipelining}
	}

	if s.ProcessingTimeout != t.ProcessingTimeout {
		diff["ProcessingTimeout"] = []interface{}{s.ProcessingTimeout, t.ProcessingTimeout}
	}

	if s.RegisterVarNames != t.RegisterVarNames {
		diff["RegisterVarNames"] = []interface{}{s.RegisterVarNames, t.RegisterVarNames}
	}

	if s.SendFragPayload != t.SendFragPayload {
		diff["SendFragPayload"] = []interface{}{s.SendFragPayload, t.SendFragPayload}
	}

	if s.UseBackend != t.UseBackend {
		diff["UseBackend"] = []interface{}{s.UseBackend, t.UseBackend}
	}

	return diff
}
