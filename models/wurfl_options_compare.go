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

// Equal checks if two structs of type WurflOptions are equal
//
//	var a, b WurflOptions
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s WurflOptions) Equal(t WurflOptions, opts ...Options) bool {
	if s.CacheSize != t.CacheSize {
		return false
	}

	if s.DataFile != t.DataFile {
		return false
	}

	if s.InformationList != t.InformationList {
		return false
	}

	if s.InformationListSeparator != t.InformationListSeparator {
		return false
	}

	if s.PatchFile != t.PatchFile {
		return false
	}

	return true
}

// Diff checks if two structs of type WurflOptions are equal
//
//	var a, b WurflOptions
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s WurflOptions) Diff(t WurflOptions, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if s.CacheSize != t.CacheSize {
		diff["CacheSize"] = []interface{}{s.CacheSize, t.CacheSize}
	}

	if s.DataFile != t.DataFile {
		diff["DataFile"] = []interface{}{s.DataFile, t.DataFile}
	}

	if s.InformationList != t.InformationList {
		diff["InformationList"] = []interface{}{s.InformationList, t.InformationList}
	}

	if s.InformationListSeparator != t.InformationListSeparator {
		diff["InformationListSeparator"] = []interface{}{s.InformationListSeparator, t.InformationListSeparator}
	}

	if s.PatchFile != t.PatchFile {
		diff["PatchFile"] = []interface{}{s.PatchFile, t.PatchFile}
	}

	return diff
}
