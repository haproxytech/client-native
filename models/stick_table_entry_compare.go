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

// Equal checks if two structs of type StickTableEntry are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b StickTableEntry
//	equal := a.Equal(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b StickTableEntry
//	equal := a.Equal(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s StickTableEntry) Equal(t StickTableEntry, opts ...Options) bool {
	opt := getOptions(opts...)

	if !equalPointers(s.BytesInCnt, t.BytesInCnt) {
		return false
	}

	if !equalPointers(s.BytesInRate, t.BytesInRate) {
		return false
	}

	if !equalPointers(s.BytesOutCnt, t.BytesOutCnt) {
		return false
	}

	if !equalPointers(s.BytesOutRate, t.BytesOutRate) {
		return false
	}

	if !equalPointers(s.ConnCnt, t.ConnCnt) {
		return false
	}

	if !equalPointers(s.ConnCur, t.ConnCur) {
		return false
	}

	if !equalPointers(s.ConnRate, t.ConnRate) {
		return false
	}

	if !equalPointers(s.Exp, t.Exp) {
		return false
	}

	if !equalPointers(s.GlitchCnt, t.GlitchCnt) {
		return false
	}

	if !equalPointers(s.GlitchRate, t.GlitchRate) {
		return false
	}

	if s.Gpc == nil || t.Gpc == nil {
		if s.Gpc != nil || t.Gpc != nil {
			if opt.NilSameAsEmpty {
				empty := &StickTableEntryGpc{}
				if s.Gpc == nil {
					if !(t.Gpc.Equal(*empty)) {
						return false
					}
				}
				if t.Gpc == nil {
					if !(s.Gpc.Equal(*empty)) {
						return false
					}
				}
			} else {
				return false
			}
		}
	} else if !s.Gpc.Equal(*t.Gpc, opt) {
		return false
	}

	if !equalPointers(s.Gpc0, t.Gpc0) {
		return false
	}

	if !equalPointers(s.Gpc0Rate, t.Gpc0Rate) {
		return false
	}

	if !equalPointers(s.Gpc1, t.Gpc1) {
		return false
	}

	if !equalPointers(s.Gpc1Rate, t.Gpc1Rate) {
		return false
	}

	if s.GpcRate == nil || t.GpcRate == nil {
		if s.GpcRate != nil || t.GpcRate != nil {
			if opt.NilSameAsEmpty {
				empty := &StickTableEntryGpcRate{}
				if s.GpcRate == nil {
					if !(t.GpcRate.Equal(*empty)) {
						return false
					}
				}
				if t.GpcRate == nil {
					if !(s.GpcRate.Equal(*empty)) {
						return false
					}
				}
			} else {
				return false
			}
		}
	} else if !s.GpcRate.Equal(*t.GpcRate, opt) {
		return false
	}

	if s.Gpt == nil || t.Gpt == nil {
		if s.Gpt != nil || t.Gpt != nil {
			if opt.NilSameAsEmpty {
				empty := &StickTableEntryGpt{}
				if s.Gpt == nil {
					if !(t.Gpt.Equal(*empty)) {
						return false
					}
				}
				if t.Gpt == nil {
					if !(s.Gpt.Equal(*empty)) {
						return false
					}
				}
			} else {
				return false
			}
		}
	} else if !s.Gpt.Equal(*t.Gpt, opt) {
		return false
	}

	if !equalPointers(s.Gpt0, t.Gpt0) {
		return false
	}

	if !equalPointers(s.HTTPErrCnt, t.HTTPErrCnt) {
		return false
	}

	if !equalPointers(s.HTTPErrRate, t.HTTPErrRate) {
		return false
	}

	if !equalPointers(s.HTTPFailCnt, t.HTTPFailCnt) {
		return false
	}

	if !equalPointers(s.HTTPFailRate, t.HTTPFailRate) {
		return false
	}

	if !equalPointers(s.HTTPReqCnt, t.HTTPReqCnt) {
		return false
	}

	if !equalPointers(s.HTTPReqRate, t.HTTPReqRate) {
		return false
	}

	if s.ID != t.ID {
		return false
	}

	if s.Key != t.Key {
		return false
	}

	if !equalPointers(s.ServerID, t.ServerID) {
		return false
	}

	if !equalPointers(s.SessCnt, t.SessCnt) {
		return false
	}

	if !equalPointers(s.SessRate, t.SessRate) {
		return false
	}

	if s.Use != t.Use {
		return false
	}

	return true
}

// Diff checks if two structs of type StickTableEntry are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b StickTableEntry
//	diff := a.Diff(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b StickTableEntry
//	diff := a.Diff(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s StickTableEntry) Diff(t StickTableEntry, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})
	if !equalPointers(s.BytesInCnt, t.BytesInCnt) {
		diff["BytesInCnt"] = []interface{}{ValueOrNil(s.BytesInCnt), ValueOrNil(t.BytesInCnt)}
	}

	if !equalPointers(s.BytesInRate, t.BytesInRate) {
		diff["BytesInRate"] = []interface{}{ValueOrNil(s.BytesInRate), ValueOrNil(t.BytesInRate)}
	}

	if !equalPointers(s.BytesOutCnt, t.BytesOutCnt) {
		diff["BytesOutCnt"] = []interface{}{ValueOrNil(s.BytesOutCnt), ValueOrNil(t.BytesOutCnt)}
	}

	if !equalPointers(s.BytesOutRate, t.BytesOutRate) {
		diff["BytesOutRate"] = []interface{}{ValueOrNil(s.BytesOutRate), ValueOrNil(t.BytesOutRate)}
	}

	if !equalPointers(s.ConnCnt, t.ConnCnt) {
		diff["ConnCnt"] = []interface{}{ValueOrNil(s.ConnCnt), ValueOrNil(t.ConnCnt)}
	}

	if !equalPointers(s.ConnCur, t.ConnCur) {
		diff["ConnCur"] = []interface{}{ValueOrNil(s.ConnCur), ValueOrNil(t.ConnCur)}
	}

	if !equalPointers(s.ConnRate, t.ConnRate) {
		diff["ConnRate"] = []interface{}{ValueOrNil(s.ConnRate), ValueOrNil(t.ConnRate)}
	}

	if !equalPointers(s.Exp, t.Exp) {
		diff["Exp"] = []interface{}{ValueOrNil(s.Exp), ValueOrNil(t.Exp)}
	}

	if !equalPointers(s.GlitchCnt, t.GlitchCnt) {
		diff["GlitchCnt"] = []interface{}{ValueOrNil(s.GlitchCnt), ValueOrNil(t.GlitchCnt)}
	}

	if !equalPointers(s.GlitchRate, t.GlitchRate) {
		diff["GlitchRate"] = []interface{}{ValueOrNil(s.GlitchRate), ValueOrNil(t.GlitchRate)}
	}

	if s.Gpc == nil || t.Gpc == nil {
		if s.Gpc != nil || t.Gpc != nil {
			if opt.NilSameAsEmpty {
				empty := &StickTableEntryGpc{}
				if s.Gpc == nil {
					if !(t.Gpc.Equal(*empty)) {
						diff["Gpc"] = []interface{}{ValueOrNil(s.Gpc), ValueOrNil(t.Gpc)}
					}
				}
				if t.Gpc == nil {
					if !(s.Gpc.Equal(*empty)) {
						diff["Gpc"] = []interface{}{ValueOrNil(s.Gpc), ValueOrNil(t.Gpc)}
					}
				}
			} else {
				diff["Gpc"] = []interface{}{ValueOrNil(s.Gpc), ValueOrNil(t.Gpc)}
			}
		}
	} else if !s.Gpc.Equal(*t.Gpc, opt) {
		diff["Gpc"] = []interface{}{ValueOrNil(s.Gpc), ValueOrNil(t.Gpc)}
	}

	if !equalPointers(s.Gpc0, t.Gpc0) {
		diff["Gpc0"] = []interface{}{ValueOrNil(s.Gpc0), ValueOrNil(t.Gpc0)}
	}

	if !equalPointers(s.Gpc0Rate, t.Gpc0Rate) {
		diff["Gpc0Rate"] = []interface{}{ValueOrNil(s.Gpc0Rate), ValueOrNil(t.Gpc0Rate)}
	}

	if !equalPointers(s.Gpc1, t.Gpc1) {
		diff["Gpc1"] = []interface{}{ValueOrNil(s.Gpc1), ValueOrNil(t.Gpc1)}
	}

	if !equalPointers(s.Gpc1Rate, t.Gpc1Rate) {
		diff["Gpc1Rate"] = []interface{}{ValueOrNil(s.Gpc1Rate), ValueOrNil(t.Gpc1Rate)}
	}

	if s.GpcRate == nil || t.GpcRate == nil {
		if s.GpcRate != nil || t.GpcRate != nil {
			if opt.NilSameAsEmpty {
				empty := &StickTableEntryGpcRate{}
				if s.GpcRate == nil {
					if !(t.GpcRate.Equal(*empty)) {
						diff["GpcRate"] = []interface{}{ValueOrNil(s.GpcRate), ValueOrNil(t.GpcRate)}
					}
				}
				if t.GpcRate == nil {
					if !(s.GpcRate.Equal(*empty)) {
						diff["GpcRate"] = []interface{}{ValueOrNil(s.GpcRate), ValueOrNil(t.GpcRate)}
					}
				}
			} else {
				diff["GpcRate"] = []interface{}{ValueOrNil(s.GpcRate), ValueOrNil(t.GpcRate)}
			}
		}
	} else if !s.GpcRate.Equal(*t.GpcRate, opt) {
		diff["GpcRate"] = []interface{}{ValueOrNil(s.GpcRate), ValueOrNil(t.GpcRate)}
	}

	if s.Gpt == nil || t.Gpt == nil {
		if s.Gpt != nil || t.Gpt != nil {
			if opt.NilSameAsEmpty {
				empty := &StickTableEntryGpt{}
				if s.Gpt == nil {
					if !(t.Gpt.Equal(*empty)) {
						diff["Gpt"] = []interface{}{ValueOrNil(s.Gpt), ValueOrNil(t.Gpt)}
					}
				}
				if t.Gpt == nil {
					if !(s.Gpt.Equal(*empty)) {
						diff["Gpt"] = []interface{}{ValueOrNil(s.Gpt), ValueOrNil(t.Gpt)}
					}
				}
			} else {
				diff["Gpt"] = []interface{}{ValueOrNil(s.Gpt), ValueOrNil(t.Gpt)}
			}
		}
	} else if !s.Gpt.Equal(*t.Gpt, opt) {
		diff["Gpt"] = []interface{}{ValueOrNil(s.Gpt), ValueOrNil(t.Gpt)}
	}

	if !equalPointers(s.Gpt0, t.Gpt0) {
		diff["Gpt0"] = []interface{}{ValueOrNil(s.Gpt0), ValueOrNil(t.Gpt0)}
	}

	if !equalPointers(s.HTTPErrCnt, t.HTTPErrCnt) {
		diff["HTTPErrCnt"] = []interface{}{ValueOrNil(s.HTTPErrCnt), ValueOrNil(t.HTTPErrCnt)}
	}

	if !equalPointers(s.HTTPErrRate, t.HTTPErrRate) {
		diff["HTTPErrRate"] = []interface{}{ValueOrNil(s.HTTPErrRate), ValueOrNil(t.HTTPErrRate)}
	}

	if !equalPointers(s.HTTPFailCnt, t.HTTPFailCnt) {
		diff["HTTPFailCnt"] = []interface{}{ValueOrNil(s.HTTPFailCnt), ValueOrNil(t.HTTPFailCnt)}
	}

	if !equalPointers(s.HTTPFailRate, t.HTTPFailRate) {
		diff["HTTPFailRate"] = []interface{}{ValueOrNil(s.HTTPFailRate), ValueOrNil(t.HTTPFailRate)}
	}

	if !equalPointers(s.HTTPReqCnt, t.HTTPReqCnt) {
		diff["HTTPReqCnt"] = []interface{}{ValueOrNil(s.HTTPReqCnt), ValueOrNil(t.HTTPReqCnt)}
	}

	if !equalPointers(s.HTTPReqRate, t.HTTPReqRate) {
		diff["HTTPReqRate"] = []interface{}{ValueOrNil(s.HTTPReqRate), ValueOrNil(t.HTTPReqRate)}
	}

	if s.ID != t.ID {
		diff["ID"] = []interface{}{s.ID, t.ID}
	}

	if s.Key != t.Key {
		diff["Key"] = []interface{}{s.Key, t.Key}
	}

	if !equalPointers(s.ServerID, t.ServerID) {
		diff["ServerID"] = []interface{}{ValueOrNil(s.ServerID), ValueOrNil(t.ServerID)}
	}

	if !equalPointers(s.SessCnt, t.SessCnt) {
		diff["SessCnt"] = []interface{}{ValueOrNil(s.SessCnt), ValueOrNil(t.SessCnt)}
	}

	if !equalPointers(s.SessRate, t.SessRate) {
		diff["SessRate"] = []interface{}{ValueOrNil(s.SessRate), ValueOrNil(t.SessRate)}
	}

	if s.Use != t.Use {
		diff["Use"] = []interface{}{s.Use, t.Use}
	}

	return diff
}

// Equal checks if two structs of type StickTableEntryGpc are equal
//
//	var a, b StickTableEntryGpc
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s StickTableEntryGpc) Equal(t StickTableEntryGpc, opts ...Options) bool {
	if s.Idx != t.Idx {
		return false
	}

	if !equalPointers(s.Value, t.Value) {
		return false
	}

	return true
}

// Diff checks if two structs of type StickTableEntryGpc are equal
//
//	var a, b StickTableEntryGpc
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s StickTableEntryGpc) Diff(t StickTableEntryGpc, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if s.Idx != t.Idx {
		diff["Idx"] = []interface{}{s.Idx, t.Idx}
	}

	if !equalPointers(s.Value, t.Value) {
		diff["Value"] = []interface{}{ValueOrNil(s.Value), ValueOrNil(t.Value)}
	}

	return diff
}

// Equal checks if two structs of type StickTableEntryGpcRate are equal
//
//	var a, b StickTableEntryGpcRate
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s StickTableEntryGpcRate) Equal(t StickTableEntryGpcRate, opts ...Options) bool {
	if s.Idx != t.Idx {
		return false
	}

	if !equalPointers(s.Value, t.Value) {
		return false
	}

	return true
}

// Diff checks if two structs of type StickTableEntryGpcRate are equal
//
//	var a, b StickTableEntryGpcRate
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s StickTableEntryGpcRate) Diff(t StickTableEntryGpcRate, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if s.Idx != t.Idx {
		diff["Idx"] = []interface{}{s.Idx, t.Idx}
	}

	if !equalPointers(s.Value, t.Value) {
		diff["Value"] = []interface{}{ValueOrNil(s.Value), ValueOrNil(t.Value)}
	}

	return diff
}

// Equal checks if two structs of type StickTableEntryGpt are equal
//
//	var a, b StickTableEntryGpt
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s StickTableEntryGpt) Equal(t StickTableEntryGpt, opts ...Options) bool {
	if s.Idx != t.Idx {
		return false
	}

	if !equalPointers(s.Value, t.Value) {
		return false
	}

	return true
}

// Diff checks if two structs of type StickTableEntryGpt are equal
//
//	var a, b StickTableEntryGpt
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s StickTableEntryGpt) Diff(t StickTableEntryGpt, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if s.Idx != t.Idx {
		diff["Idx"] = []interface{}{s.Idx, t.Idx}
	}

	if !equalPointers(s.Value, t.Value) {
		diff["Value"] = []interface{}{ValueOrNil(s.Value), ValueOrNil(t.Value)}
	}

	return diff
}
