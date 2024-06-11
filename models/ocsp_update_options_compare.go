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

// Equal checks if two structs of type OcspUpdateOptions are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b OcspUpdateOptions
//	equal := a.Equal(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b OcspUpdateOptions
//	equal := a.Equal(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s OcspUpdateOptions) Equal(t OcspUpdateOptions, opts ...Options) bool {
	opt := getOptions(opts...)

	if !equalPointers(s.Disable, t.Disable) {
		return false
	}

	if s.Httpproxy == nil || t.Httpproxy == nil {
		if s.Httpproxy != nil || t.Httpproxy != nil {
			if opt.NilSameAsEmpty {
				empty := &OcspUpdateOptionsHttpproxy{}
				if s.Httpproxy == nil {
					if !(t.Httpproxy.Equal(*empty)) {
						return false
					}
				}
				if t.Httpproxy == nil {
					if !(s.Httpproxy.Equal(*empty)) {
						return false
					}
				}
			} else {
				return false
			}
		}
	} else if !s.Httpproxy.Equal(*t.Httpproxy, opt) {
		return false
	}

	if !equalPointers(s.Maxdelay, t.Maxdelay) {
		return false
	}

	if !equalPointers(s.Mindelay, t.Mindelay) {
		return false
	}

	if s.Mode != t.Mode {
		return false
	}

	return true
}

// Diff checks if two structs of type OcspUpdateOptions are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b OcspUpdateOptions
//	diff := a.Diff(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b OcspUpdateOptions
//	diff := a.Diff(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s OcspUpdateOptions) Diff(t OcspUpdateOptions, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})
	if !equalPointers(s.Disable, t.Disable) {
		diff["Disable"] = []interface{}{ValueOrNil(s.Disable), ValueOrNil(t.Disable)}
	}

	if s.Httpproxy == nil || t.Httpproxy == nil {
		if s.Httpproxy != nil || t.Httpproxy != nil {
			if opt.NilSameAsEmpty {
				empty := &OcspUpdateOptionsHttpproxy{}
				if s.Httpproxy == nil {
					if !(t.Httpproxy.Equal(*empty)) {
						diff["Httpproxy"] = []interface{}{ValueOrNil(s.Httpproxy), ValueOrNil(t.Httpproxy)}
					}
				}
				if t.Httpproxy == nil {
					if !(s.Httpproxy.Equal(*empty)) {
						diff["Httpproxy"] = []interface{}{ValueOrNil(s.Httpproxy), ValueOrNil(t.Httpproxy)}
					}
				}
			} else {
				diff["Httpproxy"] = []interface{}{ValueOrNil(s.Httpproxy), ValueOrNil(t.Httpproxy)}
			}
		}
	} else if !s.Httpproxy.Equal(*t.Httpproxy, opt) {
		diff["Httpproxy"] = []interface{}{ValueOrNil(s.Httpproxy), ValueOrNil(t.Httpproxy)}
	}

	if !equalPointers(s.Maxdelay, t.Maxdelay) {
		diff["Maxdelay"] = []interface{}{ValueOrNil(s.Maxdelay), ValueOrNil(t.Maxdelay)}
	}

	if !equalPointers(s.Mindelay, t.Mindelay) {
		diff["Mindelay"] = []interface{}{ValueOrNil(s.Mindelay), ValueOrNil(t.Mindelay)}
	}

	if s.Mode != t.Mode {
		diff["Mode"] = []interface{}{s.Mode, t.Mode}
	}

	return diff
}

// Equal checks if two structs of type OcspUpdateOptionsHttpproxy are equal
//
//	var a, b OcspUpdateOptionsHttpproxy
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s OcspUpdateOptionsHttpproxy) Equal(t OcspUpdateOptionsHttpproxy, opts ...Options) bool {
	if s.Address != t.Address {
		return false
	}

	if !equalPointers(s.Port, t.Port) {
		return false
	}

	return true
}

// Diff checks if two structs of type OcspUpdateOptionsHttpproxy are equal
//
//	var a, b OcspUpdateOptionsHttpproxy
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s OcspUpdateOptionsHttpproxy) Diff(t OcspUpdateOptionsHttpproxy, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if s.Address != t.Address {
		diff["Address"] = []interface{}{s.Address, t.Address}
	}

	if !equalPointers(s.Port, t.Port) {
		diff["Port"] = []interface{}{ValueOrNil(s.Port), ValueOrNil(t.Port)}
	}

	return diff
}
