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

// Equal checks if two structs of type BackendBase are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b BackendBase
//	equal := a.Equal(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b BackendBase
//	equal := a.Equal(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s BackendBase) Equal(t BackendBase, opts ...Options) bool {
	opt := getOptions(opts...)

	if !CheckSameNilAndLen(s.ErrorFiles, t.ErrorFiles, opt) {
		return false
	} else {
		for i := range s.ErrorFiles {
			if !s.ErrorFiles[i].Equal(*t.ErrorFiles[i], opt) {
				return false
			}
		}
	}

	if !CheckSameNilAndLen(s.ErrorFilesFromHTTPErrors, t.ErrorFilesFromHTTPErrors, opt) {
		return false
	} else {
		for i := range s.ErrorFilesFromHTTPErrors {
			if !s.ErrorFilesFromHTTPErrors[i].Equal(*t.ErrorFilesFromHTTPErrors[i], opt) {
				return false
			}
		}
	}

	if !CheckSameNilAndLen(s.ForcePersistList, t.ForcePersistList, opt) {
		return false
	} else {
		for i := range s.ForcePersistList {
			if !s.ForcePersistList[i].Equal(*t.ForcePersistList[i], opt) {
				return false
			}
		}
	}

	if !CheckSameNilAndLen(s.IgnorePersistList, t.IgnorePersistList, opt) {
		return false
	} else {
		for i := range s.IgnorePersistList {
			if !s.IgnorePersistList[i].Equal(*t.IgnorePersistList[i], opt) {
				return false
			}
		}
	}

	if s.Abortonclose != t.Abortonclose {
		return false
	}

	if s.AcceptInvalidHTTPResponse != t.AcceptInvalidHTTPResponse {
		return false
	}

	if s.AdvCheck != t.AdvCheck {
		return false
	}

	if s.Allbackups != t.Allbackups {
		return false
	}

	if s.Balance == nil || t.Balance == nil {
		if s.Balance != nil || t.Balance != nil {
			if opt.NilSameAsEmpty {
				empty := &Balance{}
				if s.Balance == nil {
					if !(t.Balance.Equal(*empty)) {
						return false
					}
				}
				if t.Balance == nil {
					if !(s.Balance.Equal(*empty)) {
						return false
					}
				}
			} else {
				return false
			}
		}
	} else if !s.Balance.Equal(*t.Balance, opt) {
		return false
	}

	if !equalPointers(s.CheckTimeout, t.CheckTimeout) {
		return false
	}

	if s.Checkcache != t.Checkcache {
		return false
	}

	if s.Compression == nil || t.Compression == nil {
		if s.Compression != nil || t.Compression != nil {
			if opt.NilSameAsEmpty {
				empty := &Compression{}
				if s.Compression == nil {
					if !(t.Compression.Equal(*empty)) {
						return false
					}
				}
				if t.Compression == nil {
					if !(s.Compression.Equal(*empty)) {
						return false
					}
				}
			} else {
				return false
			}
		}
	} else if !s.Compression.Equal(*t.Compression, opt) {
		return false
	}

	if !equalPointers(s.ConnectTimeout, t.ConnectTimeout) {
		return false
	}

	if s.Cookie == nil || t.Cookie == nil {
		if s.Cookie != nil || t.Cookie != nil {
			if opt.NilSameAsEmpty {
				empty := &Cookie{}
				if s.Cookie == nil {
					if !(t.Cookie.Equal(*empty)) {
						return false
					}
				}
				if t.Cookie == nil {
					if !(s.Cookie.Equal(*empty)) {
						return false
					}
				}
			} else {
				return false
			}
		}
	} else if !s.Cookie.Equal(*t.Cookie, opt) {
		return false
	}

	if s.DefaultServer == nil || t.DefaultServer == nil {
		if s.DefaultServer != nil || t.DefaultServer != nil {
			if opt.NilSameAsEmpty {
				empty := &DefaultServer{}
				if s.DefaultServer == nil {
					if !(t.DefaultServer.Equal(*empty)) {
						return false
					}
				}
				if t.DefaultServer == nil {
					if !(s.DefaultServer.Equal(*empty)) {
						return false
					}
				}
			} else {
				return false
			}
		}
	} else if !s.DefaultServer.Equal(*t.DefaultServer, opt) {
		return false
	}

	if s.Description != t.Description {
		return false
	}

	if s.Disabled != t.Disabled {
		return false
	}

	if s.DynamicCookieKey != t.DynamicCookieKey {
		return false
	}

	if s.EmailAlert == nil || t.EmailAlert == nil {
		if s.EmailAlert != nil || t.EmailAlert != nil {
			if opt.NilSameAsEmpty {
				empty := &EmailAlert{}
				if s.EmailAlert == nil {
					if !(t.EmailAlert.Equal(*empty)) {
						return false
					}
				}
				if t.EmailAlert == nil {
					if !(s.EmailAlert.Equal(*empty)) {
						return false
					}
				}
			} else {
				return false
			}
		}
	} else if !s.EmailAlert.Equal(*t.EmailAlert, opt) {
		return false
	}

	if s.Enabled != t.Enabled {
		return false
	}

	if s.Errorloc302 == nil || t.Errorloc302 == nil {
		if s.Errorloc302 != nil || t.Errorloc302 != nil {
			if opt.NilSameAsEmpty {
				empty := &Errorloc{}
				if s.Errorloc302 == nil {
					if !(t.Errorloc302.Equal(*empty)) {
						return false
					}
				}
				if t.Errorloc302 == nil {
					if !(s.Errorloc302.Equal(*empty)) {
						return false
					}
				}
			} else {
				return false
			}
		}
	} else if !s.Errorloc302.Equal(*t.Errorloc302, opt) {
		return false
	}

	if s.Errorloc303 == nil || t.Errorloc303 == nil {
		if s.Errorloc303 != nil || t.Errorloc303 != nil {
			if opt.NilSameAsEmpty {
				empty := &Errorloc{}
				if s.Errorloc303 == nil {
					if !(t.Errorloc303.Equal(*empty)) {
						return false
					}
				}
				if t.Errorloc303 == nil {
					if !(s.Errorloc303.Equal(*empty)) {
						return false
					}
				}
			} else {
				return false
			}
		}
	} else if !s.Errorloc303.Equal(*t.Errorloc303, opt) {
		return false
	}

	if s.ExternalCheck != t.ExternalCheck {
		return false
	}

	if s.ExternalCheckCommand != t.ExternalCheckCommand {
		return false
	}

	if s.ExternalCheckPath != t.ExternalCheckPath {
		return false
	}

	if s.ForcePersist == nil || t.ForcePersist == nil {
		if s.ForcePersist != nil || t.ForcePersist != nil {
			if opt.NilSameAsEmpty {
				empty := &BackendForcePersist{}
				if s.ForcePersist == nil {
					if !(t.ForcePersist.Equal(*empty)) {
						return false
					}
				}
				if t.ForcePersist == nil {
					if !(s.ForcePersist.Equal(*empty)) {
						return false
					}
				}
			} else {
				return false
			}
		}
	} else if !s.ForcePersist.Equal(*t.ForcePersist, opt) {
		return false
	}

	if s.Forwardfor == nil || t.Forwardfor == nil {
		if s.Forwardfor != nil || t.Forwardfor != nil {
			if opt.NilSameAsEmpty {
				empty := &Forwardfor{}
				if s.Forwardfor == nil {
					if !(t.Forwardfor.Equal(*empty)) {
						return false
					}
				}
				if t.Forwardfor == nil {
					if !(s.Forwardfor.Equal(*empty)) {
						return false
					}
				}
			} else {
				return false
			}
		}
	} else if !s.Forwardfor.Equal(*t.Forwardfor, opt) {
		return false
	}

	if s.From != t.From {
		return false
	}

	if !equalPointers(s.Fullconn, t.Fullconn) {
		return false
	}

	if s.GUID != t.GUID {
		return false
	}

	if s.H1CaseAdjustBogusServer != t.H1CaseAdjustBogusServer {
		return false
	}

	if !equalPointers(s.HashBalanceFactor, t.HashBalanceFactor) {
		return false
	}

	if s.HashType == nil || t.HashType == nil {
		if s.HashType != nil || t.HashType != nil {
			if opt.NilSameAsEmpty {
				empty := &HashType{}
				if s.HashType == nil {
					if !(t.HashType.Equal(*empty)) {
						return false
					}
				}
				if t.HashType == nil {
					if !(s.HashType.Equal(*empty)) {
						return false
					}
				}
			} else {
				return false
			}
		}
	} else if !s.HashType.Equal(*t.HashType, opt) {
		return false
	}

	if s.HTTPBufferRequest != t.HTTPBufferRequest {
		return false
	}

	if s.HTTPNoDelay != t.HTTPNoDelay {
		return false
	}

	if s.HTTPUseHtx != t.HTTPUseHtx {
		return false
	}

	if s.HTTPConnectionMode != t.HTTPConnectionMode {
		return false
	}

	if !equalPointers(s.HTTPKeepAliveTimeout, t.HTTPKeepAliveTimeout) {
		return false
	}

	if s.HTTPPretendKeepalive != t.HTTPPretendKeepalive {
		return false
	}

	if s.HTTPProxy != t.HTTPProxy {
		return false
	}

	if !equalPointers(s.HTTPRequestTimeout, t.HTTPRequestTimeout) {
		return false
	}

	if s.HTTPRestrictReqHdrNames != t.HTTPRestrictReqHdrNames {
		return false
	}

	if s.HTTPReuse != t.HTTPReuse {
		return false
	}

	if !equalPointers(s.HTTPSendNameHeader, t.HTTPSendNameHeader) {
		return false
	}

	if s.HttpchkParams == nil || t.HttpchkParams == nil {
		if s.HttpchkParams != nil || t.HttpchkParams != nil {
			if opt.NilSameAsEmpty {
				empty := &HttpchkParams{}
				if s.HttpchkParams == nil {
					if !(t.HttpchkParams.Equal(*empty)) {
						return false
					}
				}
				if t.HttpchkParams == nil {
					if !(s.HttpchkParams.Equal(*empty)) {
						return false
					}
				}
			} else {
				return false
			}
		}
	} else if !s.HttpchkParams.Equal(*t.HttpchkParams, opt) {
		return false
	}

	if !equalPointers(s.ID, t.ID) {
		return false
	}

	if s.IgnorePersist == nil || t.IgnorePersist == nil {
		if s.IgnorePersist != nil || t.IgnorePersist != nil {
			if opt.NilSameAsEmpty {
				empty := &BackendIgnorePersist{}
				if s.IgnorePersist == nil {
					if !(t.IgnorePersist.Equal(*empty)) {
						return false
					}
				}
				if t.IgnorePersist == nil {
					if !(s.IgnorePersist.Equal(*empty)) {
						return false
					}
				}
			} else {
				return false
			}
		}
	} else if !s.IgnorePersist.Equal(*t.IgnorePersist, opt) {
		return false
	}

	if s.IndependentStreams != t.IndependentStreams {
		return false
	}

	if s.LoadServerStateFromFile != t.LoadServerStateFromFile {
		return false
	}

	if s.LogHealthChecks != t.LogHealthChecks {
		return false
	}

	if s.LogTag != t.LogTag {
		return false
	}

	if !equalPointers(s.MaxKeepAliveQueue, t.MaxKeepAliveQueue) {
		return false
	}

	if s.Mode != t.Mode {
		return false
	}

	if s.MysqlCheckParams == nil || t.MysqlCheckParams == nil {
		if s.MysqlCheckParams != nil || t.MysqlCheckParams != nil {
			if opt.NilSameAsEmpty {
				empty := &MysqlCheckParams{}
				if s.MysqlCheckParams == nil {
					if !(t.MysqlCheckParams.Equal(*empty)) {
						return false
					}
				}
				if t.MysqlCheckParams == nil {
					if !(s.MysqlCheckParams.Equal(*empty)) {
						return false
					}
				}
			} else {
				return false
			}
		}
	} else if !s.MysqlCheckParams.Equal(*t.MysqlCheckParams, opt) {
		return false
	}

	if s.Name != t.Name {
		return false
	}

	if s.Nolinger != t.Nolinger {
		return false
	}

	if s.Originalto == nil || t.Originalto == nil {
		if s.Originalto != nil || t.Originalto != nil {
			if opt.NilSameAsEmpty {
				empty := &Originalto{}
				if s.Originalto == nil {
					if !(t.Originalto.Equal(*empty)) {
						return false
					}
				}
				if t.Originalto == nil {
					if !(s.Originalto.Equal(*empty)) {
						return false
					}
				}
			} else {
				return false
			}
		}
	} else if !s.Originalto.Equal(*t.Originalto, opt) {
		return false
	}

	if s.Persist != t.Persist {
		return false
	}

	if s.PersistRule == nil || t.PersistRule == nil {
		if s.PersistRule != nil || t.PersistRule != nil {
			if opt.NilSameAsEmpty {
				empty := &PersistRule{}
				if s.PersistRule == nil {
					if !(t.PersistRule.Equal(*empty)) {
						return false
					}
				}
				if t.PersistRule == nil {
					if !(s.PersistRule.Equal(*empty)) {
						return false
					}
				}
			} else {
				return false
			}
		}
	} else if !s.PersistRule.Equal(*t.PersistRule, opt) {
		return false
	}

	if s.PgsqlCheckParams == nil || t.PgsqlCheckParams == nil {
		if s.PgsqlCheckParams != nil || t.PgsqlCheckParams != nil {
			if opt.NilSameAsEmpty {
				empty := &PgsqlCheckParams{}
				if s.PgsqlCheckParams == nil {
					if !(t.PgsqlCheckParams.Equal(*empty)) {
						return false
					}
				}
				if t.PgsqlCheckParams == nil {
					if !(s.PgsqlCheckParams.Equal(*empty)) {
						return false
					}
				}
			} else {
				return false
			}
		}
	} else if !s.PgsqlCheckParams.Equal(*t.PgsqlCheckParams, opt) {
		return false
	}

	if s.PreferLastServer != t.PreferLastServer {
		return false
	}

	if !equalPointers(s.QueueTimeout, t.QueueTimeout) {
		return false
	}

	if s.Redispatch == nil || t.Redispatch == nil {
		if s.Redispatch != nil || t.Redispatch != nil {
			if opt.NilSameAsEmpty {
				empty := &Redispatch{}
				if s.Redispatch == nil {
					if !(t.Redispatch.Equal(*empty)) {
						return false
					}
				}
				if t.Redispatch == nil {
					if !(s.Redispatch.Equal(*empty)) {
						return false
					}
				}
			} else {
				return false
			}
		}
	} else if !s.Redispatch.Equal(*t.Redispatch, opt) {
		return false
	}

	if !equalPointers(s.Retries, t.Retries) {
		return false
	}

	if s.RetryOn != t.RetryOn {
		return false
	}

	if !equalPointers(s.ServerFinTimeout, t.ServerFinTimeout) {
		return false
	}

	if s.ServerStateFileName != t.ServerStateFileName {
		return false
	}

	if !equalPointers(s.ServerTimeout, t.ServerTimeout) {
		return false
	}

	if s.SmtpchkParams == nil || t.SmtpchkParams == nil {
		if s.SmtpchkParams != nil || t.SmtpchkParams != nil {
			if opt.NilSameAsEmpty {
				empty := &SmtpchkParams{}
				if s.SmtpchkParams == nil {
					if !(t.SmtpchkParams.Equal(*empty)) {
						return false
					}
				}
				if t.SmtpchkParams == nil {
					if !(s.SmtpchkParams.Equal(*empty)) {
						return false
					}
				}
			} else {
				return false
			}
		}
	} else if !s.SmtpchkParams.Equal(*t.SmtpchkParams, opt) {
		return false
	}

	if s.Source == nil || t.Source == nil {
		if s.Source != nil || t.Source != nil {
			if opt.NilSameAsEmpty {
				empty := &Source{}
				if s.Source == nil {
					if !(t.Source.Equal(*empty)) {
						return false
					}
				}
				if t.Source == nil {
					if !(s.Source.Equal(*empty)) {
						return false
					}
				}
			} else {
				return false
			}
		}
	} else if !s.Source.Equal(*t.Source, opt) {
		return false
	}

	if s.SpliceAuto != t.SpliceAuto {
		return false
	}

	if s.SpliceRequest != t.SpliceRequest {
		return false
	}

	if s.SpliceResponse != t.SpliceResponse {
		return false
	}

	if s.SpopCheck != t.SpopCheck {
		return false
	}

	if s.Srvtcpka != t.Srvtcpka {
		return false
	}

	if !equalPointers(s.SrvtcpkaCnt, t.SrvtcpkaCnt) {
		return false
	}

	if !equalPointers(s.SrvtcpkaIdle, t.SrvtcpkaIdle) {
		return false
	}

	if !equalPointers(s.SrvtcpkaIntvl, t.SrvtcpkaIntvl) {
		return false
	}

	if s.StatsOptions == nil || t.StatsOptions == nil {
		if s.StatsOptions != nil || t.StatsOptions != nil {
			if opt.NilSameAsEmpty {
				empty := &StatsOptions{}
				if s.StatsOptions == nil {
					if !(t.StatsOptions.Equal(*empty)) {
						return false
					}
				}
				if t.StatsOptions == nil {
					if !(s.StatsOptions.Equal(*empty)) {
						return false
					}
				}
			} else {
				return false
			}
		}
	} else if !s.StatsOptions.Equal(*t.StatsOptions, opt) {
		return false
	}

	if s.StickTable == nil || t.StickTable == nil {
		if s.StickTable != nil || t.StickTable != nil {
			if opt.NilSameAsEmpty {
				empty := &ConfigStickTable{}
				if s.StickTable == nil {
					if !(t.StickTable.Equal(*empty)) {
						return false
					}
				}
				if t.StickTable == nil {
					if !(s.StickTable.Equal(*empty)) {
						return false
					}
				}
			} else {
				return false
			}
		}
	} else if !s.StickTable.Equal(*t.StickTable, opt) {
		return false
	}

	if !equalPointers(s.TarpitTimeout, t.TarpitTimeout) {
		return false
	}

	if s.TCPSmartConnect != t.TCPSmartConnect {
		return false
	}

	if s.Tcpka != t.Tcpka {
		return false
	}

	if s.Transparent != t.Transparent {
		return false
	}

	if !equalPointers(s.TunnelTimeout, t.TunnelTimeout) {
		return false
	}

	if s.UseFCGIApp != t.UseFCGIApp {
		return false
	}

	return true
}

// Diff checks if two structs of type BackendBase are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b BackendBase
//	diff := a.Diff(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b BackendBase
//	diff := a.Diff(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s BackendBase) Diff(t BackendBase, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})
	if !CheckSameNilAndLen(s.ErrorFiles, t.ErrorFiles, opt) {
		diff["ErrorFiles"] = []interface{}{s.ErrorFiles, t.ErrorFiles}
	} else {
		diff2 := make(map[string][]interface{})
		for i := range s.ErrorFiles {
			if !s.ErrorFiles[i].Equal(*t.ErrorFiles[i], opt) {
				diffSub := s.ErrorFiles[i].Diff(*t.ErrorFiles[i], opt)
				if len(diffSub) > 0 {
					diff2[strconv.Itoa(i)] = []interface{}{diffSub}
				}
			}
		}
		if len(diff2) > 0 {
			diff["ErrorFiles"] = []interface{}{diff2}
		}
	}

	if !CheckSameNilAndLen(s.ErrorFilesFromHTTPErrors, t.ErrorFilesFromHTTPErrors, opt) {
		diff["ErrorFilesFromHTTPErrors"] = []interface{}{s.ErrorFilesFromHTTPErrors, t.ErrorFilesFromHTTPErrors}
	} else {
		diff2 := make(map[string][]interface{})
		for i := range s.ErrorFilesFromHTTPErrors {
			if !s.ErrorFilesFromHTTPErrors[i].Equal(*t.ErrorFilesFromHTTPErrors[i], opt) {
				diffSub := s.ErrorFilesFromHTTPErrors[i].Diff(*t.ErrorFilesFromHTTPErrors[i], opt)
				if len(diffSub) > 0 {
					diff2[strconv.Itoa(i)] = []interface{}{diffSub}
				}
			}
		}
		if len(diff2) > 0 {
			diff["ErrorFilesFromHTTPErrors"] = []interface{}{diff2}
		}
	}

	if !CheckSameNilAndLen(s.ForcePersistList, t.ForcePersistList, opt) {
		diff["ForcePersistList"] = []interface{}{s.ForcePersistList, t.ForcePersistList}
	} else {
		diff2 := make(map[string][]interface{})
		for i := range s.ForcePersistList {
			if !s.ForcePersistList[i].Equal(*t.ForcePersistList[i], opt) {
				diffSub := s.ForcePersistList[i].Diff(*t.ForcePersistList[i], opt)
				if len(diffSub) > 0 {
					diff2[strconv.Itoa(i)] = []interface{}{diffSub}
				}
			}
		}
		if len(diff2) > 0 {
			diff["ForcePersistList"] = []interface{}{diff2}
		}
	}

	if !CheckSameNilAndLen(s.IgnorePersistList, t.IgnorePersistList, opt) {
		diff["IgnorePersistList"] = []interface{}{s.IgnorePersistList, t.IgnorePersistList}
	} else {
		diff2 := make(map[string][]interface{})
		for i := range s.IgnorePersistList {
			if !s.IgnorePersistList[i].Equal(*t.IgnorePersistList[i], opt) {
				diffSub := s.IgnorePersistList[i].Diff(*t.IgnorePersistList[i], opt)
				if len(diffSub) > 0 {
					diff2[strconv.Itoa(i)] = []interface{}{diffSub}
				}
			}
		}
		if len(diff2) > 0 {
			diff["IgnorePersistList"] = []interface{}{diff2}
		}
	}

	if s.Abortonclose != t.Abortonclose {
		diff["Abortonclose"] = []interface{}{s.Abortonclose, t.Abortonclose}
	}

	if s.AcceptInvalidHTTPResponse != t.AcceptInvalidHTTPResponse {
		diff["AcceptInvalidHTTPResponse"] = []interface{}{s.AcceptInvalidHTTPResponse, t.AcceptInvalidHTTPResponse}
	}

	if s.AdvCheck != t.AdvCheck {
		diff["AdvCheck"] = []interface{}{s.AdvCheck, t.AdvCheck}
	}

	if s.Allbackups != t.Allbackups {
		diff["Allbackups"] = []interface{}{s.Allbackups, t.Allbackups}
	}

	if s.Balance == nil || t.Balance == nil {
		if s.Balance != nil || t.Balance != nil {
			if opt.NilSameAsEmpty {
				empty := &Balance{}
				if s.Balance == nil {
					if !(t.Balance.Equal(*empty)) {
						diff["Balance"] = []interface{}{ValueOrNil(s.Balance), ValueOrNil(t.Balance)}
					}
				}
				if t.Balance == nil {
					if !(s.Balance.Equal(*empty)) {
						diff["Balance"] = []interface{}{ValueOrNil(s.Balance), ValueOrNil(t.Balance)}
					}
				}
			} else {
				diff["Balance"] = []interface{}{ValueOrNil(s.Balance), ValueOrNil(t.Balance)}
			}
		}
	} else if !s.Balance.Equal(*t.Balance, opt) {
		diff["Balance"] = []interface{}{ValueOrNil(s.Balance), ValueOrNil(t.Balance)}
	}

	if !equalPointers(s.CheckTimeout, t.CheckTimeout) {
		diff["CheckTimeout"] = []interface{}{ValueOrNil(s.CheckTimeout), ValueOrNil(t.CheckTimeout)}
	}

	if s.Checkcache != t.Checkcache {
		diff["Checkcache"] = []interface{}{s.Checkcache, t.Checkcache}
	}

	if s.Compression == nil || t.Compression == nil {
		if s.Compression != nil || t.Compression != nil {
			if opt.NilSameAsEmpty {
				empty := &Compression{}
				if s.Compression == nil {
					if !(t.Compression.Equal(*empty)) {
						diff["Compression"] = []interface{}{ValueOrNil(s.Compression), ValueOrNil(t.Compression)}
					}
				}
				if t.Compression == nil {
					if !(s.Compression.Equal(*empty)) {
						diff["Compression"] = []interface{}{ValueOrNil(s.Compression), ValueOrNil(t.Compression)}
					}
				}
			} else {
				diff["Compression"] = []interface{}{ValueOrNil(s.Compression), ValueOrNil(t.Compression)}
			}
		}
	} else if !s.Compression.Equal(*t.Compression, opt) {
		diff["Compression"] = []interface{}{ValueOrNil(s.Compression), ValueOrNil(t.Compression)}
	}

	if !equalPointers(s.ConnectTimeout, t.ConnectTimeout) {
		diff["ConnectTimeout"] = []interface{}{ValueOrNil(s.ConnectTimeout), ValueOrNil(t.ConnectTimeout)}
	}

	if s.Cookie == nil || t.Cookie == nil {
		if s.Cookie != nil || t.Cookie != nil {
			if opt.NilSameAsEmpty {
				empty := &Cookie{}
				if s.Cookie == nil {
					if !(t.Cookie.Equal(*empty)) {
						diff["Cookie"] = []interface{}{ValueOrNil(s.Cookie), ValueOrNil(t.Cookie)}
					}
				}
				if t.Cookie == nil {
					if !(s.Cookie.Equal(*empty)) {
						diff["Cookie"] = []interface{}{ValueOrNil(s.Cookie), ValueOrNil(t.Cookie)}
					}
				}
			} else {
				diff["Cookie"] = []interface{}{ValueOrNil(s.Cookie), ValueOrNil(t.Cookie)}
			}
		}
	} else if !s.Cookie.Equal(*t.Cookie, opt) {
		diff["Cookie"] = []interface{}{ValueOrNil(s.Cookie), ValueOrNil(t.Cookie)}
	}

	if s.DefaultServer == nil || t.DefaultServer == nil {
		if s.DefaultServer != nil || t.DefaultServer != nil {
			if opt.NilSameAsEmpty {
				empty := &DefaultServer{}
				if s.DefaultServer == nil {
					if !(t.DefaultServer.Equal(*empty)) {
						diff["DefaultServer"] = []interface{}{ValueOrNil(s.DefaultServer), ValueOrNil(t.DefaultServer)}
					}
				}
				if t.DefaultServer == nil {
					if !(s.DefaultServer.Equal(*empty)) {
						diff["DefaultServer"] = []interface{}{ValueOrNil(s.DefaultServer), ValueOrNil(t.DefaultServer)}
					}
				}
			} else {
				diff["DefaultServer"] = []interface{}{ValueOrNil(s.DefaultServer), ValueOrNil(t.DefaultServer)}
			}
		}
	} else if !s.DefaultServer.Equal(*t.DefaultServer, opt) {
		diff["DefaultServer"] = []interface{}{ValueOrNil(s.DefaultServer), ValueOrNil(t.DefaultServer)}
	}

	if s.Description != t.Description {
		diff["Description"] = []interface{}{s.Description, t.Description}
	}

	if s.Disabled != t.Disabled {
		diff["Disabled"] = []interface{}{s.Disabled, t.Disabled}
	}

	if s.DynamicCookieKey != t.DynamicCookieKey {
		diff["DynamicCookieKey"] = []interface{}{s.DynamicCookieKey, t.DynamicCookieKey}
	}

	if s.EmailAlert == nil || t.EmailAlert == nil {
		if s.EmailAlert != nil || t.EmailAlert != nil {
			if opt.NilSameAsEmpty {
				empty := &EmailAlert{}
				if s.EmailAlert == nil {
					if !(t.EmailAlert.Equal(*empty)) {
						diff["EmailAlert"] = []interface{}{ValueOrNil(s.EmailAlert), ValueOrNil(t.EmailAlert)}
					}
				}
				if t.EmailAlert == nil {
					if !(s.EmailAlert.Equal(*empty)) {
						diff["EmailAlert"] = []interface{}{ValueOrNil(s.EmailAlert), ValueOrNil(t.EmailAlert)}
					}
				}
			} else {
				diff["EmailAlert"] = []interface{}{ValueOrNil(s.EmailAlert), ValueOrNil(t.EmailAlert)}
			}
		}
	} else if !s.EmailAlert.Equal(*t.EmailAlert, opt) {
		diff["EmailAlert"] = []interface{}{ValueOrNil(s.EmailAlert), ValueOrNil(t.EmailAlert)}
	}

	if s.Enabled != t.Enabled {
		diff["Enabled"] = []interface{}{s.Enabled, t.Enabled}
	}

	if s.Errorloc302 == nil || t.Errorloc302 == nil {
		if s.Errorloc302 != nil || t.Errorloc302 != nil {
			if opt.NilSameAsEmpty {
				empty := &Errorloc{}
				if s.Errorloc302 == nil {
					if !(t.Errorloc302.Equal(*empty)) {
						diff["Errorloc302"] = []interface{}{ValueOrNil(s.Errorloc302), ValueOrNil(t.Errorloc302)}
					}
				}
				if t.Errorloc302 == nil {
					if !(s.Errorloc302.Equal(*empty)) {
						diff["Errorloc302"] = []interface{}{ValueOrNil(s.Errorloc302), ValueOrNil(t.Errorloc302)}
					}
				}
			} else {
				diff["Errorloc302"] = []interface{}{ValueOrNil(s.Errorloc302), ValueOrNil(t.Errorloc302)}
			}
		}
	} else if !s.Errorloc302.Equal(*t.Errorloc302, opt) {
		diff["Errorloc302"] = []interface{}{ValueOrNil(s.Errorloc302), ValueOrNil(t.Errorloc302)}
	}

	if s.Errorloc303 == nil || t.Errorloc303 == nil {
		if s.Errorloc303 != nil || t.Errorloc303 != nil {
			if opt.NilSameAsEmpty {
				empty := &Errorloc{}
				if s.Errorloc303 == nil {
					if !(t.Errorloc303.Equal(*empty)) {
						diff["Errorloc303"] = []interface{}{ValueOrNil(s.Errorloc303), ValueOrNil(t.Errorloc303)}
					}
				}
				if t.Errorloc303 == nil {
					if !(s.Errorloc303.Equal(*empty)) {
						diff["Errorloc303"] = []interface{}{ValueOrNil(s.Errorloc303), ValueOrNil(t.Errorloc303)}
					}
				}
			} else {
				diff["Errorloc303"] = []interface{}{ValueOrNil(s.Errorloc303), ValueOrNil(t.Errorloc303)}
			}
		}
	} else if !s.Errorloc303.Equal(*t.Errorloc303, opt) {
		diff["Errorloc303"] = []interface{}{ValueOrNil(s.Errorloc303), ValueOrNil(t.Errorloc303)}
	}

	if s.ExternalCheck != t.ExternalCheck {
		diff["ExternalCheck"] = []interface{}{s.ExternalCheck, t.ExternalCheck}
	}

	if s.ExternalCheckCommand != t.ExternalCheckCommand {
		diff["ExternalCheckCommand"] = []interface{}{s.ExternalCheckCommand, t.ExternalCheckCommand}
	}

	if s.ExternalCheckPath != t.ExternalCheckPath {
		diff["ExternalCheckPath"] = []interface{}{s.ExternalCheckPath, t.ExternalCheckPath}
	}

	if s.ForcePersist == nil || t.ForcePersist == nil {
		if s.ForcePersist != nil || t.ForcePersist != nil {
			if opt.NilSameAsEmpty {
				empty := &BackendForcePersist{}
				if s.ForcePersist == nil {
					if !(t.ForcePersist.Equal(*empty)) {
						diff["ForcePersist"] = []interface{}{ValueOrNil(s.ForcePersist), ValueOrNil(t.ForcePersist)}
					}
				}
				if t.ForcePersist == nil {
					if !(s.ForcePersist.Equal(*empty)) {
						diff["ForcePersist"] = []interface{}{ValueOrNil(s.ForcePersist), ValueOrNil(t.ForcePersist)}
					}
				}
			} else {
				diff["ForcePersist"] = []interface{}{ValueOrNil(s.ForcePersist), ValueOrNil(t.ForcePersist)}
			}
		}
	} else if !s.ForcePersist.Equal(*t.ForcePersist, opt) {
		diff["ForcePersist"] = []interface{}{ValueOrNil(s.ForcePersist), ValueOrNil(t.ForcePersist)}
	}

	if s.Forwardfor == nil || t.Forwardfor == nil {
		if s.Forwardfor != nil || t.Forwardfor != nil {
			if opt.NilSameAsEmpty {
				empty := &Forwardfor{}
				if s.Forwardfor == nil {
					if !(t.Forwardfor.Equal(*empty)) {
						diff["Forwardfor"] = []interface{}{ValueOrNil(s.Forwardfor), ValueOrNil(t.Forwardfor)}
					}
				}
				if t.Forwardfor == nil {
					if !(s.Forwardfor.Equal(*empty)) {
						diff["Forwardfor"] = []interface{}{ValueOrNil(s.Forwardfor), ValueOrNil(t.Forwardfor)}
					}
				}
			} else {
				diff["Forwardfor"] = []interface{}{ValueOrNil(s.Forwardfor), ValueOrNil(t.Forwardfor)}
			}
		}
	} else if !s.Forwardfor.Equal(*t.Forwardfor, opt) {
		diff["Forwardfor"] = []interface{}{ValueOrNil(s.Forwardfor), ValueOrNil(t.Forwardfor)}
	}

	if s.From != t.From {
		diff["From"] = []interface{}{s.From, t.From}
	}

	if !equalPointers(s.Fullconn, t.Fullconn) {
		diff["Fullconn"] = []interface{}{ValueOrNil(s.Fullconn), ValueOrNil(t.Fullconn)}
	}

	if s.GUID != t.GUID {
		diff["GUID"] = []interface{}{s.GUID, t.GUID}
	}

	if s.H1CaseAdjustBogusServer != t.H1CaseAdjustBogusServer {
		diff["H1CaseAdjustBogusServer"] = []interface{}{s.H1CaseAdjustBogusServer, t.H1CaseAdjustBogusServer}
	}

	if !equalPointers(s.HashBalanceFactor, t.HashBalanceFactor) {
		diff["HashBalanceFactor"] = []interface{}{ValueOrNil(s.HashBalanceFactor), ValueOrNil(t.HashBalanceFactor)}
	}

	if s.HashType == nil || t.HashType == nil {
		if s.HashType != nil || t.HashType != nil {
			if opt.NilSameAsEmpty {
				empty := &HashType{}
				if s.HashType == nil {
					if !(t.HashType.Equal(*empty)) {
						diff["HashType"] = []interface{}{ValueOrNil(s.HashType), ValueOrNil(t.HashType)}
					}
				}
				if t.HashType == nil {
					if !(s.HashType.Equal(*empty)) {
						diff["HashType"] = []interface{}{ValueOrNil(s.HashType), ValueOrNil(t.HashType)}
					}
				}
			} else {
				diff["HashType"] = []interface{}{ValueOrNil(s.HashType), ValueOrNil(t.HashType)}
			}
		}
	} else if !s.HashType.Equal(*t.HashType, opt) {
		diff["HashType"] = []interface{}{ValueOrNil(s.HashType), ValueOrNil(t.HashType)}
	}

	if s.HTTPBufferRequest != t.HTTPBufferRequest {
		diff["HTTPBufferRequest"] = []interface{}{s.HTTPBufferRequest, t.HTTPBufferRequest}
	}

	if s.HTTPNoDelay != t.HTTPNoDelay {
		diff["HTTPNoDelay"] = []interface{}{s.HTTPNoDelay, t.HTTPNoDelay}
	}

	if s.HTTPUseHtx != t.HTTPUseHtx {
		diff["HTTPUseHtx"] = []interface{}{s.HTTPUseHtx, t.HTTPUseHtx}
	}

	if s.HTTPConnectionMode != t.HTTPConnectionMode {
		diff["HTTPConnectionMode"] = []interface{}{s.HTTPConnectionMode, t.HTTPConnectionMode}
	}

	if !equalPointers(s.HTTPKeepAliveTimeout, t.HTTPKeepAliveTimeout) {
		diff["HTTPKeepAliveTimeout"] = []interface{}{ValueOrNil(s.HTTPKeepAliveTimeout), ValueOrNil(t.HTTPKeepAliveTimeout)}
	}

	if s.HTTPPretendKeepalive != t.HTTPPretendKeepalive {
		diff["HTTPPretendKeepalive"] = []interface{}{s.HTTPPretendKeepalive, t.HTTPPretendKeepalive}
	}

	if s.HTTPProxy != t.HTTPProxy {
		diff["HTTPProxy"] = []interface{}{s.HTTPProxy, t.HTTPProxy}
	}

	if !equalPointers(s.HTTPRequestTimeout, t.HTTPRequestTimeout) {
		diff["HTTPRequestTimeout"] = []interface{}{ValueOrNil(s.HTTPRequestTimeout), ValueOrNil(t.HTTPRequestTimeout)}
	}

	if s.HTTPRestrictReqHdrNames != t.HTTPRestrictReqHdrNames {
		diff["HTTPRestrictReqHdrNames"] = []interface{}{s.HTTPRestrictReqHdrNames, t.HTTPRestrictReqHdrNames}
	}

	if s.HTTPReuse != t.HTTPReuse {
		diff["HTTPReuse"] = []interface{}{s.HTTPReuse, t.HTTPReuse}
	}

	if !equalPointers(s.HTTPSendNameHeader, t.HTTPSendNameHeader) {
		diff["HTTPSendNameHeader"] = []interface{}{ValueOrNil(s.HTTPSendNameHeader), ValueOrNil(t.HTTPSendNameHeader)}
	}

	if s.HttpchkParams == nil || t.HttpchkParams == nil {
		if s.HttpchkParams != nil || t.HttpchkParams != nil {
			if opt.NilSameAsEmpty {
				empty := &HttpchkParams{}
				if s.HttpchkParams == nil {
					if !(t.HttpchkParams.Equal(*empty)) {
						diff["HttpchkParams"] = []interface{}{ValueOrNil(s.HttpchkParams), ValueOrNil(t.HttpchkParams)}
					}
				}
				if t.HttpchkParams == nil {
					if !(s.HttpchkParams.Equal(*empty)) {
						diff["HttpchkParams"] = []interface{}{ValueOrNil(s.HttpchkParams), ValueOrNil(t.HttpchkParams)}
					}
				}
			} else {
				diff["HttpchkParams"] = []interface{}{ValueOrNil(s.HttpchkParams), ValueOrNil(t.HttpchkParams)}
			}
		}
	} else if !s.HttpchkParams.Equal(*t.HttpchkParams, opt) {
		diff["HttpchkParams"] = []interface{}{ValueOrNil(s.HttpchkParams), ValueOrNil(t.HttpchkParams)}
	}

	if !equalPointers(s.ID, t.ID) {
		diff["ID"] = []interface{}{ValueOrNil(s.ID), ValueOrNil(t.ID)}
	}

	if s.IgnorePersist == nil || t.IgnorePersist == nil {
		if s.IgnorePersist != nil || t.IgnorePersist != nil {
			if opt.NilSameAsEmpty {
				empty := &BackendIgnorePersist{}
				if s.IgnorePersist == nil {
					if !(t.IgnorePersist.Equal(*empty)) {
						diff["IgnorePersist"] = []interface{}{ValueOrNil(s.IgnorePersist), ValueOrNil(t.IgnorePersist)}
					}
				}
				if t.IgnorePersist == nil {
					if !(s.IgnorePersist.Equal(*empty)) {
						diff["IgnorePersist"] = []interface{}{ValueOrNil(s.IgnorePersist), ValueOrNil(t.IgnorePersist)}
					}
				}
			} else {
				diff["IgnorePersist"] = []interface{}{ValueOrNil(s.IgnorePersist), ValueOrNil(t.IgnorePersist)}
			}
		}
	} else if !s.IgnorePersist.Equal(*t.IgnorePersist, opt) {
		diff["IgnorePersist"] = []interface{}{ValueOrNil(s.IgnorePersist), ValueOrNil(t.IgnorePersist)}
	}

	if s.IndependentStreams != t.IndependentStreams {
		diff["IndependentStreams"] = []interface{}{s.IndependentStreams, t.IndependentStreams}
	}

	if s.LoadServerStateFromFile != t.LoadServerStateFromFile {
		diff["LoadServerStateFromFile"] = []interface{}{s.LoadServerStateFromFile, t.LoadServerStateFromFile}
	}

	if s.LogHealthChecks != t.LogHealthChecks {
		diff["LogHealthChecks"] = []interface{}{s.LogHealthChecks, t.LogHealthChecks}
	}

	if s.LogTag != t.LogTag {
		diff["LogTag"] = []interface{}{s.LogTag, t.LogTag}
	}

	if !equalPointers(s.MaxKeepAliveQueue, t.MaxKeepAliveQueue) {
		diff["MaxKeepAliveQueue"] = []interface{}{ValueOrNil(s.MaxKeepAliveQueue), ValueOrNil(t.MaxKeepAliveQueue)}
	}

	if s.Mode != t.Mode {
		diff["Mode"] = []interface{}{s.Mode, t.Mode}
	}

	if s.MysqlCheckParams == nil || t.MysqlCheckParams == nil {
		if s.MysqlCheckParams != nil || t.MysqlCheckParams != nil {
			if opt.NilSameAsEmpty {
				empty := &MysqlCheckParams{}
				if s.MysqlCheckParams == nil {
					if !(t.MysqlCheckParams.Equal(*empty)) {
						diff["MysqlCheckParams"] = []interface{}{ValueOrNil(s.MysqlCheckParams), ValueOrNil(t.MysqlCheckParams)}
					}
				}
				if t.MysqlCheckParams == nil {
					if !(s.MysqlCheckParams.Equal(*empty)) {
						diff["MysqlCheckParams"] = []interface{}{ValueOrNil(s.MysqlCheckParams), ValueOrNil(t.MysqlCheckParams)}
					}
				}
			} else {
				diff["MysqlCheckParams"] = []interface{}{ValueOrNil(s.MysqlCheckParams), ValueOrNil(t.MysqlCheckParams)}
			}
		}
	} else if !s.MysqlCheckParams.Equal(*t.MysqlCheckParams, opt) {
		diff["MysqlCheckParams"] = []interface{}{ValueOrNil(s.MysqlCheckParams), ValueOrNil(t.MysqlCheckParams)}
	}

	if s.Name != t.Name {
		diff["Name"] = []interface{}{s.Name, t.Name}
	}

	if s.Nolinger != t.Nolinger {
		diff["Nolinger"] = []interface{}{s.Nolinger, t.Nolinger}
	}

	if s.Originalto == nil || t.Originalto == nil {
		if s.Originalto != nil || t.Originalto != nil {
			if opt.NilSameAsEmpty {
				empty := &Originalto{}
				if s.Originalto == nil {
					if !(t.Originalto.Equal(*empty)) {
						diff["Originalto"] = []interface{}{ValueOrNil(s.Originalto), ValueOrNil(t.Originalto)}
					}
				}
				if t.Originalto == nil {
					if !(s.Originalto.Equal(*empty)) {
						diff["Originalto"] = []interface{}{ValueOrNil(s.Originalto), ValueOrNil(t.Originalto)}
					}
				}
			} else {
				diff["Originalto"] = []interface{}{ValueOrNil(s.Originalto), ValueOrNil(t.Originalto)}
			}
		}
	} else if !s.Originalto.Equal(*t.Originalto, opt) {
		diff["Originalto"] = []interface{}{ValueOrNil(s.Originalto), ValueOrNil(t.Originalto)}
	}

	if s.Persist != t.Persist {
		diff["Persist"] = []interface{}{s.Persist, t.Persist}
	}

	if s.PersistRule == nil || t.PersistRule == nil {
		if s.PersistRule != nil || t.PersistRule != nil {
			if opt.NilSameAsEmpty {
				empty := &PersistRule{}
				if s.PersistRule == nil {
					if !(t.PersistRule.Equal(*empty)) {
						diff["PersistRule"] = []interface{}{ValueOrNil(s.PersistRule), ValueOrNil(t.PersistRule)}
					}
				}
				if t.PersistRule == nil {
					if !(s.PersistRule.Equal(*empty)) {
						diff["PersistRule"] = []interface{}{ValueOrNil(s.PersistRule), ValueOrNil(t.PersistRule)}
					}
				}
			} else {
				diff["PersistRule"] = []interface{}{ValueOrNil(s.PersistRule), ValueOrNil(t.PersistRule)}
			}
		}
	} else if !s.PersistRule.Equal(*t.PersistRule, opt) {
		diff["PersistRule"] = []interface{}{ValueOrNil(s.PersistRule), ValueOrNil(t.PersistRule)}
	}

	if s.PgsqlCheckParams == nil || t.PgsqlCheckParams == nil {
		if s.PgsqlCheckParams != nil || t.PgsqlCheckParams != nil {
			if opt.NilSameAsEmpty {
				empty := &PgsqlCheckParams{}
				if s.PgsqlCheckParams == nil {
					if !(t.PgsqlCheckParams.Equal(*empty)) {
						diff["PgsqlCheckParams"] = []interface{}{ValueOrNil(s.PgsqlCheckParams), ValueOrNil(t.PgsqlCheckParams)}
					}
				}
				if t.PgsqlCheckParams == nil {
					if !(s.PgsqlCheckParams.Equal(*empty)) {
						diff["PgsqlCheckParams"] = []interface{}{ValueOrNil(s.PgsqlCheckParams), ValueOrNil(t.PgsqlCheckParams)}
					}
				}
			} else {
				diff["PgsqlCheckParams"] = []interface{}{ValueOrNil(s.PgsqlCheckParams), ValueOrNil(t.PgsqlCheckParams)}
			}
		}
	} else if !s.PgsqlCheckParams.Equal(*t.PgsqlCheckParams, opt) {
		diff["PgsqlCheckParams"] = []interface{}{ValueOrNil(s.PgsqlCheckParams), ValueOrNil(t.PgsqlCheckParams)}
	}

	if s.PreferLastServer != t.PreferLastServer {
		diff["PreferLastServer"] = []interface{}{s.PreferLastServer, t.PreferLastServer}
	}

	if !equalPointers(s.QueueTimeout, t.QueueTimeout) {
		diff["QueueTimeout"] = []interface{}{ValueOrNil(s.QueueTimeout), ValueOrNil(t.QueueTimeout)}
	}

	if s.Redispatch == nil || t.Redispatch == nil {
		if s.Redispatch != nil || t.Redispatch != nil {
			if opt.NilSameAsEmpty {
				empty := &Redispatch{}
				if s.Redispatch == nil {
					if !(t.Redispatch.Equal(*empty)) {
						diff["Redispatch"] = []interface{}{ValueOrNil(s.Redispatch), ValueOrNil(t.Redispatch)}
					}
				}
				if t.Redispatch == nil {
					if !(s.Redispatch.Equal(*empty)) {
						diff["Redispatch"] = []interface{}{ValueOrNil(s.Redispatch), ValueOrNil(t.Redispatch)}
					}
				}
			} else {
				diff["Redispatch"] = []interface{}{ValueOrNil(s.Redispatch), ValueOrNil(t.Redispatch)}
			}
		}
	} else if !s.Redispatch.Equal(*t.Redispatch, opt) {
		diff["Redispatch"] = []interface{}{ValueOrNil(s.Redispatch), ValueOrNil(t.Redispatch)}
	}

	if !equalPointers(s.Retries, t.Retries) {
		diff["Retries"] = []interface{}{ValueOrNil(s.Retries), ValueOrNil(t.Retries)}
	}

	if s.RetryOn != t.RetryOn {
		diff["RetryOn"] = []interface{}{s.RetryOn, t.RetryOn}
	}

	if !equalPointers(s.ServerFinTimeout, t.ServerFinTimeout) {
		diff["ServerFinTimeout"] = []interface{}{ValueOrNil(s.ServerFinTimeout), ValueOrNil(t.ServerFinTimeout)}
	}

	if s.ServerStateFileName != t.ServerStateFileName {
		diff["ServerStateFileName"] = []interface{}{s.ServerStateFileName, t.ServerStateFileName}
	}

	if !equalPointers(s.ServerTimeout, t.ServerTimeout) {
		diff["ServerTimeout"] = []interface{}{ValueOrNil(s.ServerTimeout), ValueOrNil(t.ServerTimeout)}
	}

	if s.SmtpchkParams == nil || t.SmtpchkParams == nil {
		if s.SmtpchkParams != nil || t.SmtpchkParams != nil {
			if opt.NilSameAsEmpty {
				empty := &SmtpchkParams{}
				if s.SmtpchkParams == nil {
					if !(t.SmtpchkParams.Equal(*empty)) {
						diff["SmtpchkParams"] = []interface{}{ValueOrNil(s.SmtpchkParams), ValueOrNil(t.SmtpchkParams)}
					}
				}
				if t.SmtpchkParams == nil {
					if !(s.SmtpchkParams.Equal(*empty)) {
						diff["SmtpchkParams"] = []interface{}{ValueOrNil(s.SmtpchkParams), ValueOrNil(t.SmtpchkParams)}
					}
				}
			} else {
				diff["SmtpchkParams"] = []interface{}{ValueOrNil(s.SmtpchkParams), ValueOrNil(t.SmtpchkParams)}
			}
		}
	} else if !s.SmtpchkParams.Equal(*t.SmtpchkParams, opt) {
		diff["SmtpchkParams"] = []interface{}{ValueOrNil(s.SmtpchkParams), ValueOrNil(t.SmtpchkParams)}
	}

	if s.Source == nil || t.Source == nil {
		if s.Source != nil || t.Source != nil {
			if opt.NilSameAsEmpty {
				empty := &Source{}
				if s.Source == nil {
					if !(t.Source.Equal(*empty)) {
						diff["Source"] = []interface{}{ValueOrNil(s.Source), ValueOrNil(t.Source)}
					}
				}
				if t.Source == nil {
					if !(s.Source.Equal(*empty)) {
						diff["Source"] = []interface{}{ValueOrNil(s.Source), ValueOrNil(t.Source)}
					}
				}
			} else {
				diff["Source"] = []interface{}{ValueOrNil(s.Source), ValueOrNil(t.Source)}
			}
		}
	} else if !s.Source.Equal(*t.Source, opt) {
		diff["Source"] = []interface{}{ValueOrNil(s.Source), ValueOrNil(t.Source)}
	}

	if s.SpliceAuto != t.SpliceAuto {
		diff["SpliceAuto"] = []interface{}{s.SpliceAuto, t.SpliceAuto}
	}

	if s.SpliceRequest != t.SpliceRequest {
		diff["SpliceRequest"] = []interface{}{s.SpliceRequest, t.SpliceRequest}
	}

	if s.SpliceResponse != t.SpliceResponse {
		diff["SpliceResponse"] = []interface{}{s.SpliceResponse, t.SpliceResponse}
	}

	if s.SpopCheck != t.SpopCheck {
		diff["SpopCheck"] = []interface{}{s.SpopCheck, t.SpopCheck}
	}

	if s.Srvtcpka != t.Srvtcpka {
		diff["Srvtcpka"] = []interface{}{s.Srvtcpka, t.Srvtcpka}
	}

	if !equalPointers(s.SrvtcpkaCnt, t.SrvtcpkaCnt) {
		diff["SrvtcpkaCnt"] = []interface{}{ValueOrNil(s.SrvtcpkaCnt), ValueOrNil(t.SrvtcpkaCnt)}
	}

	if !equalPointers(s.SrvtcpkaIdle, t.SrvtcpkaIdle) {
		diff["SrvtcpkaIdle"] = []interface{}{ValueOrNil(s.SrvtcpkaIdle), ValueOrNil(t.SrvtcpkaIdle)}
	}

	if !equalPointers(s.SrvtcpkaIntvl, t.SrvtcpkaIntvl) {
		diff["SrvtcpkaIntvl"] = []interface{}{ValueOrNil(s.SrvtcpkaIntvl), ValueOrNil(t.SrvtcpkaIntvl)}
	}

	if s.StatsOptions == nil || t.StatsOptions == nil {
		if s.StatsOptions != nil || t.StatsOptions != nil {
			if opt.NilSameAsEmpty {
				empty := &StatsOptions{}
				if s.StatsOptions == nil {
					if !(t.StatsOptions.Equal(*empty)) {
						diff["StatsOptions"] = []interface{}{ValueOrNil(s.StatsOptions), ValueOrNil(t.StatsOptions)}
					}
				}
				if t.StatsOptions == nil {
					if !(s.StatsOptions.Equal(*empty)) {
						diff["StatsOptions"] = []interface{}{ValueOrNil(s.StatsOptions), ValueOrNil(t.StatsOptions)}
					}
				}
			} else {
				diff["StatsOptions"] = []interface{}{ValueOrNil(s.StatsOptions), ValueOrNil(t.StatsOptions)}
			}
		}
	} else if !s.StatsOptions.Equal(*t.StatsOptions, opt) {
		diff["StatsOptions"] = []interface{}{ValueOrNil(s.StatsOptions), ValueOrNil(t.StatsOptions)}
	}

	if s.StickTable == nil || t.StickTable == nil {
		if s.StickTable != nil || t.StickTable != nil {
			if opt.NilSameAsEmpty {
				empty := &ConfigStickTable{}
				if s.StickTable == nil {
					if !(t.StickTable.Equal(*empty)) {
						diff["StickTable"] = []interface{}{ValueOrNil(s.StickTable), ValueOrNil(t.StickTable)}
					}
				}
				if t.StickTable == nil {
					if !(s.StickTable.Equal(*empty)) {
						diff["StickTable"] = []interface{}{ValueOrNil(s.StickTable), ValueOrNil(t.StickTable)}
					}
				}
			} else {
				diff["StickTable"] = []interface{}{ValueOrNil(s.StickTable), ValueOrNil(t.StickTable)}
			}
		}
	} else if !s.StickTable.Equal(*t.StickTable, opt) {
		diff["StickTable"] = []interface{}{ValueOrNil(s.StickTable), ValueOrNil(t.StickTable)}
	}

	if !equalPointers(s.TarpitTimeout, t.TarpitTimeout) {
		diff["TarpitTimeout"] = []interface{}{ValueOrNil(s.TarpitTimeout), ValueOrNil(t.TarpitTimeout)}
	}

	if s.TCPSmartConnect != t.TCPSmartConnect {
		diff["TCPSmartConnect"] = []interface{}{s.TCPSmartConnect, t.TCPSmartConnect}
	}

	if s.Tcpka != t.Tcpka {
		diff["Tcpka"] = []interface{}{s.Tcpka, t.Tcpka}
	}

	if s.Transparent != t.Transparent {
		diff["Transparent"] = []interface{}{s.Transparent, t.Transparent}
	}

	if !equalPointers(s.TunnelTimeout, t.TunnelTimeout) {
		diff["TunnelTimeout"] = []interface{}{ValueOrNil(s.TunnelTimeout), ValueOrNil(t.TunnelTimeout)}
	}

	if s.UseFCGIApp != t.UseFCGIApp {
		diff["UseFCGIApp"] = []interface{}{s.UseFCGIApp, t.UseFCGIApp}
	}

	return diff
}

// Equal checks if two structs of type BackendForcePersist are equal
//
//	var a, b BackendForcePersist
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s BackendForcePersist) Equal(t BackendForcePersist, opts ...Options) bool {
	if !equalPointers(s.Cond, t.Cond) {
		return false
	}

	if !equalPointers(s.CondTest, t.CondTest) {
		return false
	}

	return true
}

// Diff checks if two structs of type BackendForcePersist are equal
//
//	var a, b BackendForcePersist
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s BackendForcePersist) Diff(t BackendForcePersist, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if !equalPointers(s.Cond, t.Cond) {
		diff["Cond"] = []interface{}{ValueOrNil(s.Cond), ValueOrNil(t.Cond)}
	}

	if !equalPointers(s.CondTest, t.CondTest) {
		diff["CondTest"] = []interface{}{ValueOrNil(s.CondTest), ValueOrNil(t.CondTest)}
	}

	return diff
}

// Equal checks if two structs of type ForcePersist are equal
//
//	var a, b ForcePersist
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s ForcePersist) Equal(t ForcePersist, opts ...Options) bool {
	if !equalPointers(s.Cond, t.Cond) {
		return false
	}

	if !equalPointers(s.CondTest, t.CondTest) {
		return false
	}

	return true
}

// Diff checks if two structs of type ForcePersist are equal
//
//	var a, b ForcePersist
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s ForcePersist) Diff(t ForcePersist, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if !equalPointers(s.Cond, t.Cond) {
		diff["Cond"] = []interface{}{ValueOrNil(s.Cond), ValueOrNil(t.Cond)}
	}

	if !equalPointers(s.CondTest, t.CondTest) {
		diff["CondTest"] = []interface{}{ValueOrNil(s.CondTest), ValueOrNil(t.CondTest)}
	}

	return diff
}

// Equal checks if two structs of type BackendIgnorePersist are equal
//
//	var a, b BackendIgnorePersist
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s BackendIgnorePersist) Equal(t BackendIgnorePersist, opts ...Options) bool {
	if !equalPointers(s.Cond, t.Cond) {
		return false
	}

	if !equalPointers(s.CondTest, t.CondTest) {
		return false
	}

	return true
}

// Diff checks if two structs of type BackendIgnorePersist are equal
//
//	var a, b BackendIgnorePersist
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s BackendIgnorePersist) Diff(t BackendIgnorePersist, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if !equalPointers(s.Cond, t.Cond) {
		diff["Cond"] = []interface{}{ValueOrNil(s.Cond), ValueOrNil(t.Cond)}
	}

	if !equalPointers(s.CondTest, t.CondTest) {
		diff["CondTest"] = []interface{}{ValueOrNil(s.CondTest), ValueOrNil(t.CondTest)}
	}

	return diff
}

// Equal checks if two structs of type IgnorePersist are equal
//
//	var a, b IgnorePersist
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s IgnorePersist) Equal(t IgnorePersist, opts ...Options) bool {
	if !equalPointers(s.Cond, t.Cond) {
		return false
	}

	if !equalPointers(s.CondTest, t.CondTest) {
		return false
	}

	return true
}

// Diff checks if two structs of type IgnorePersist are equal
//
//	var a, b IgnorePersist
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s IgnorePersist) Diff(t IgnorePersist, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if !equalPointers(s.Cond, t.Cond) {
		diff["Cond"] = []interface{}{ValueOrNil(s.Cond), ValueOrNil(t.Cond)}
	}

	if !equalPointers(s.CondTest, t.CondTest) {
		diff["CondTest"] = []interface{}{ValueOrNil(s.CondTest), ValueOrNil(t.CondTest)}
	}

	return diff
}
