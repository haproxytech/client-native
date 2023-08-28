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

// Equal checks if two structs of type Site are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b Site
//	equal := a.Equal(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b Site
//	equal := a.Equal(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s Site) Equal(t Site, opts ...Options) bool {
	opt := getOptions(opts...)

	if !CheckSameNilAndLen(s.Farms, t.Farms, opt) {
		return false
	} else {
		for i := range s.Farms {
			if !s.Farms[i].Equal(*t.Farms[i], opt) {
				return false
			}
		}
	}

	if s.Name != t.Name {
		return false
	}

	if !s.Service.Equal(*t.Service, opt) {
		return false
	}

	return true
}

// Diff checks if two structs of type Site are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b Site
//	diff := a.Diff(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b Site
//	diff := a.Diff(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s Site) Diff(t Site, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})
	if !CheckSameNilAndLen(s.Farms, t.Farms, opt) {
		diff["Farms"] = []interface{}{s.Farms, t.Farms}
	} else {
		diff2 := make(map[string][]interface{})
		for i := range s.Farms {
			if !s.Farms[i].Equal(*t.Farms[i], opt) {
				diffSub := s.Farms[i].Diff(*t.Farms[i], opt)
				if len(diffSub) > 0 {
					diff2[strconv.Itoa(i)] = []interface{}{diffSub}
				}
			}
		}
		if len(diff2) > 0 {
			diff["Farms"] = []interface{}{diff2}
		}
	}

	if s.Name != t.Name {
		diff["Name"] = []interface{}{s.Name, t.Name}
	}

	if !s.Service.Equal(*t.Service, opt) {
		diff["Service"] = []interface{}{ValueOrNil(s.Service), ValueOrNil(t.Service)}
	}

	return diff
}

// Equal checks if two structs of type SiteFarm are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b SiteFarm
//	equal := a.Equal(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b SiteFarm
//	equal := a.Equal(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s SiteFarm) Equal(t SiteFarm, opts ...Options) bool {
	opt := getOptions(opts...)

	if !s.Balance.Equal(*t.Balance, opt) {
		return false
	}

	if s.Cond != t.Cond {
		return false
	}

	if s.CondTest != t.CondTest {
		return false
	}

	if !s.Forwardfor.Equal(*t.Forwardfor, opt) {
		return false
	}

	if s.Mode != t.Mode {
		return false
	}

	if s.Name != t.Name {
		return false
	}

	if !CheckSameNilAndLen(s.Servers, t.Servers, opt) {
		return false
	} else {
		for i := range s.Servers {
			if !s.Servers[i].Equal(*t.Servers[i], opt) {
				return false
			}
		}
	}

	if s.UseAs != t.UseAs {
		return false
	}

	return true
}

// Diff checks if two structs of type SiteFarm are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b SiteFarm
//	diff := a.Diff(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b SiteFarm
//	diff := a.Diff(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s SiteFarm) Diff(t SiteFarm, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})
	if !s.Balance.Equal(*t.Balance, opt) {
		diff["Balance"] = []interface{}{ValueOrNil(s.Balance), ValueOrNil(t.Balance)}
	}

	if s.Cond != t.Cond {
		diff["Cond"] = []interface{}{s.Cond, t.Cond}
	}

	if s.CondTest != t.CondTest {
		diff["CondTest"] = []interface{}{s.CondTest, t.CondTest}
	}

	if !s.Forwardfor.Equal(*t.Forwardfor, opt) {
		diff["Forwardfor"] = []interface{}{ValueOrNil(s.Forwardfor), ValueOrNil(t.Forwardfor)}
	}

	if s.Mode != t.Mode {
		diff["Mode"] = []interface{}{s.Mode, t.Mode}
	}

	if s.Name != t.Name {
		diff["Name"] = []interface{}{s.Name, t.Name}
	}

	if !CheckSameNilAndLen(s.Servers, t.Servers, opt) {
		diff["Servers"] = []interface{}{s.Servers, t.Servers}
	} else {
		diff2 := make(map[string][]interface{})
		for i := range s.Servers {
			if !s.Servers[i].Equal(*t.Servers[i], opt) {
				diffSub := s.Servers[i].Diff(*t.Servers[i], opt)
				if len(diffSub) > 0 {
					diff2[strconv.Itoa(i)] = []interface{}{diffSub}
				}
			}
		}
		if len(diff2) > 0 {
			diff["Servers"] = []interface{}{diff2}
		}
	}

	if s.UseAs != t.UseAs {
		diff["UseAs"] = []interface{}{s.UseAs, t.UseAs}
	}

	return diff
}

// Equal checks if two structs of type SiteService are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b SiteService
//	equal := a.Equal(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b SiteService
//	equal := a.Equal(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s SiteService) Equal(t SiteService, opts ...Options) bool {
	opt := getOptions(opts...)

	if s.HTTPConnectionMode != t.HTTPConnectionMode {
		return false
	}

	if !CheckSameNilAndLen(s.Listeners, t.Listeners, opt) {
		return false
	} else {
		for i := range s.Listeners {
			if !s.Listeners[i].Equal(*t.Listeners[i], opt) {
				return false
			}
		}
	}

	if !equalPointers(s.Maxconn, t.Maxconn) {
		return false
	}

	if s.Mode != t.Mode {
		return false
	}

	return true
}

// Diff checks if two structs of type SiteService are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b SiteService
//	diff := a.Diff(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b SiteService
//	diff := a.Diff(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s SiteService) Diff(t SiteService, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})
	if s.HTTPConnectionMode != t.HTTPConnectionMode {
		diff["HTTPConnectionMode"] = []interface{}{s.HTTPConnectionMode, t.HTTPConnectionMode}
	}

	if !CheckSameNilAndLen(s.Listeners, t.Listeners, opt) {
		diff["Listeners"] = []interface{}{s.Listeners, t.Listeners}
	} else {
		diff2 := make(map[string][]interface{})
		for i := range s.Listeners {
			if !s.Listeners[i].Equal(*t.Listeners[i], opt) {
				diffSub := s.Listeners[i].Diff(*t.Listeners[i], opt)
				if len(diffSub) > 0 {
					diff2[strconv.Itoa(i)] = []interface{}{diffSub}
				}
			}
		}
		if len(diff2) > 0 {
			diff["Listeners"] = []interface{}{diff2}
		}
	}

	if !equalPointers(s.Maxconn, t.Maxconn) {
		diff["Maxconn"] = []interface{}{ValueOrNil(s.Maxconn), ValueOrNil(t.Maxconn)}
	}

	if s.Mode != t.Mode {
		diff["Mode"] = []interface{}{s.Mode, t.Mode}
	}

	return diff
}
