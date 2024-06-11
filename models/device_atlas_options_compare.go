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

// Equal checks if two structs of type DeviceAtlasOptions are equal
//
//	var a, b DeviceAtlasOptions
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s DeviceAtlasOptions) Equal(t DeviceAtlasOptions, opts ...Options) bool {
	if s.JSONFile != t.JSONFile {
		return false
	}

	if s.LogLevel != t.LogLevel {
		return false
	}

	if s.PropertiesCookie != t.PropertiesCookie {
		return false
	}

	if s.Separator != t.Separator {
		return false
	}

	return true
}

// Diff checks if two structs of type DeviceAtlasOptions are equal
//
//	var a, b DeviceAtlasOptions
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s DeviceAtlasOptions) Diff(t DeviceAtlasOptions, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if s.JSONFile != t.JSONFile {
		diff["JSONFile"] = []interface{}{s.JSONFile, t.JSONFile}
	}

	if s.LogLevel != t.LogLevel {
		diff["LogLevel"] = []interface{}{s.LogLevel, t.LogLevel}
	}

	if s.PropertiesCookie != t.PropertiesCookie {
		diff["PropertiesCookie"] = []interface{}{s.PropertiesCookie, t.PropertiesCookie}
	}

	if s.Separator != t.Separator {
		diff["Separator"] = []interface{}{s.Separator, t.Separator}
	}

	return diff
}
