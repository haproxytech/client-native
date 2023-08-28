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

// Equal checks if two structs of type SpoeMessage are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b SpoeMessage
//	equal := a.Equal(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b SpoeMessage
//	equal := a.Equal(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s SpoeMessage) Equal(t SpoeMessage, opts ...Options) bool {
	opt := getOptions(opts...)

	if !s.ACL.Equal(t.ACL, opt) {
		return false
	}

	if s.Args != t.Args {
		return false
	}

	if !s.Event.Equal(*t.Event, opt) {
		return false
	}

	if !equalPointers(s.Name, t.Name) {
		return false
	}

	return true
}

// Diff checks if two structs of type SpoeMessage are equal
//
// By default empty arrays, maps and slices are equal to nil:
//
//	var a, b SpoeMessage
//	diff := a.Diff(b)
//
// For more advanced use case you can configure the options (default values are shown):
//
//	var a, b SpoeMessage
//	equal := a.Diff(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s SpoeMessage) Diff(t SpoeMessage, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})
	if !s.ACL.Equal(t.ACL, opt) {
		diff["ACL"] = []interface{}{s.ACL, t.ACL}
	}

	if s.Args != t.Args {
		diff["Args"] = []interface{}{s.Args, t.Args}
	}

	if !s.Event.Equal(*t.Event, opt) {
		diff["Event"] = []interface{}{ValueOrNil(s.Event), ValueOrNil(t.Event)}
	}

	if !equalPointers(s.Name, t.Name) {
		diff["Name"] = []interface{}{ValueOrNil(s.Name), ValueOrNil(t.Name)}
	}

	return diff
}

// Equal checks if two structs of type SpoeMessageEvent are equal
//
//	var a, b SpoeMessageEvent
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s SpoeMessageEvent) Equal(t SpoeMessageEvent, opts ...Options) bool {
	if s.Cond != t.Cond {
		return false
	}

	if s.CondTest != t.CondTest {
		return false
	}

	if !equalPointers(s.Name, t.Name) {
		return false
	}

	return true
}

// Diff checks if two structs of type SpoeMessageEvent are equal
//
//	var a, b SpoeMessageEvent
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s SpoeMessageEvent) Diff(t SpoeMessageEvent, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if s.Cond != t.Cond {
		diff["Cond"] = []interface{}{s.Cond, t.Cond}
	}

	if s.CondTest != t.CondTest {
		diff["CondTest"] = []interface{}{s.CondTest, t.CondTest}
	}

	if !equalPointers(s.Name, t.Name) {
		diff["Name"] = []interface{}{ValueOrNil(s.Name), ValueOrNil(t.Name)}
	}

	return diff
}
