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

	if !CheckSameNilAndLen(s.ErrorFiles, t.ErrorFiles, opt) {
		return false
	}
	for i := range s.ErrorFiles {
		if !s.ErrorFiles[i].Equal(*t.ErrorFiles[i], opt) {
			return false
		}
	}

	if !CheckSameNilAndLen(s.ErrorFilesFromHTTPErrors, t.ErrorFilesFromHTTPErrors, opt) {
		return false
	}
	for i := range s.ErrorFilesFromHTTPErrors {
		if !s.ErrorFilesFromHTTPErrors[i].Equal(*t.ErrorFilesFromHTTPErrors[i], opt) {
			return false
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

	if !s.Balance.Equal(*t.Balance, opt) {
		return false
	}

	if s.BindProcess != t.BindProcess {
		return false
	}

	if !equalPointers(s.CheckTimeout, t.CheckTimeout) {
		return false
	}

	if s.Checkcache != t.Checkcache {
		return false
	}

	if !s.Compression.Equal(*t.Compression, opt) {
		return false
	}

	if !equalPointers(s.ConnectTimeout, t.ConnectTimeout) {
		return false
	}

	if !s.Cookie.Equal(*t.Cookie, opt) {
		return false
	}

	if !s.DefaultServer.Equal(*t.DefaultServer, opt) {
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

	if !s.EmailAlert.Equal(*t.EmailAlert, opt) {
		return false
	}

	if s.Enabled != t.Enabled {
		return false
	}

	if !s.Errorloc302.Equal(*t.Errorloc302, opt) {
		return false
	}

	if !s.Errorloc303.Equal(*t.Errorloc303, opt) {
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

	if !s.ForcePersist.Equal(*t.ForcePersist, opt) {
		return false
	}

	if !s.Forwardfor.Equal(*t.Forwardfor, opt) {
		return false
	}

	if s.From != t.From {
		return false
	}

	if !equalPointers(s.Fullconn, t.Fullconn) {
		return false
	}

	if s.H1CaseAdjustBogusServer != t.H1CaseAdjustBogusServer {
		return false
	}

	if !s.HashType.Equal(*t.HashType, opt) {
		return false
	}

	if s.HTTPBufferRequest != t.HTTPBufferRequest {
		return false
	}

	if !s.HTTPCheck.Equal(*t.HTTPCheck, opt) {
		return false
	}

	if s.HTTPKeepAlive != t.HTTPKeepAlive {
		return false
	}

	if s.HTTPNoDelay != t.HTTPNoDelay {
		return false
	}

	if s.HTTPServerClose != t.HTTPServerClose {
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

	if !s.HttpchkParams.Equal(*t.HttpchkParams, opt) {
		return false
	}

	if s.Httpclose != t.Httpclose {
		return false
	}

	if !equalPointers(s.ID, t.ID) {
		return false
	}

	if !s.IgnorePersist.Equal(*t.IgnorePersist, opt) {
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

	if !s.MysqlCheckParams.Equal(*t.MysqlCheckParams, opt) {
		return false
	}

	if s.Name != t.Name {
		return false
	}

	if s.Nolinger != t.Nolinger {
		return false
	}

	if !s.Originalto.Equal(*t.Originalto, opt) {
		return false
	}

	if s.Persist != t.Persist {
		return false
	}

	if !s.PersistRule.Equal(*t.PersistRule, opt) {
		return false
	}

	if !s.PgsqlCheckParams.Equal(*t.PgsqlCheckParams, opt) {
		return false
	}

	if s.PreferLastServer != t.PreferLastServer {
		return false
	}

	if !equalPointers(s.QueueTimeout, t.QueueTimeout) {
		return false
	}

	if !s.Redispatch.Equal(*t.Redispatch, opt) {
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

	if !s.SmtpchkParams.Equal(*t.SmtpchkParams, opt) {
		return false
	}

	if !s.Source.Equal(*t.Source, opt) {
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

	if !s.StatsOptions.Equal(*t.StatsOptions, opt) {
		return false
	}

	if !s.StickTable.Equal(*t.StickTable, opt) {
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

// Diff checks if two structs of type Backend are equal
//
// By default empty arrays, maps and slices are equal to nil:
//
//	var a, b Backend
//	diff := a.Diff(b)
//
// For more advanced use case you can configure the options (default values are shown):
//
//	var a, b Backend
//	equal := a.Diff(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s Backend) Diff(t Backend, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})
	if len(s.ErrorFiles) != len(t.ErrorFiles) {
		diff["ErrorFiles"] = []interface{}{s.ErrorFiles, t.ErrorFiles}
	} else {
		diff2 := make(map[string][]interface{})
		for i := range s.ErrorFiles {
			diffSub := s.ErrorFiles[i].Diff(*t.ErrorFiles[i], opt)
			if len(diffSub) > 0 {
				diff2[strconv.Itoa(i)] = []interface{}{diffSub}
			}
		}
		if len(diff2) > 0 {
			diff["ErrorFiles"] = []interface{}{diff2}
		}
	}

	if len(s.ErrorFilesFromHTTPErrors) != len(t.ErrorFilesFromHTTPErrors) {
		diff["ErrorFilesFromHTTPErrors"] = []interface{}{s.ErrorFilesFromHTTPErrors, t.ErrorFilesFromHTTPErrors}
	} else {
		diff2 := make(map[string][]interface{})
		for i := range s.ErrorFilesFromHTTPErrors {
			diffSub := s.ErrorFilesFromHTTPErrors[i].Diff(*t.ErrorFilesFromHTTPErrors[i], opt)
			if len(diffSub) > 0 {
				diff2[strconv.Itoa(i)] = []interface{}{diffSub}
			}
		}
		if len(diff2) > 0 {
			diff["ErrorFilesFromHTTPErrors"] = []interface{}{diff2}
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

	if !s.Balance.Equal(*t.Balance, opt) {
		diff["Balance"] = []interface{}{s.Balance, t.Balance}
	}

	if s.BindProcess != t.BindProcess {
		diff["BindProcess"] = []interface{}{s.BindProcess, t.BindProcess}
	}

	if !equalPointers(s.CheckTimeout, t.CheckTimeout) {
		diff["CheckTimeout"] = []interface{}{s.CheckTimeout, t.CheckTimeout}
	}

	if s.Checkcache != t.Checkcache {
		diff["Checkcache"] = []interface{}{s.Checkcache, t.Checkcache}
	}

	if !s.Compression.Equal(*t.Compression, opt) {
		diff["Compression"] = []interface{}{s.Compression, t.Compression}
	}

	if !equalPointers(s.ConnectTimeout, t.ConnectTimeout) {
		diff["ConnectTimeout"] = []interface{}{s.ConnectTimeout, t.ConnectTimeout}
	}

	if !s.Cookie.Equal(*t.Cookie, opt) {
		diff["Cookie"] = []interface{}{s.Cookie, t.Cookie}
	}

	if !s.DefaultServer.Equal(*t.DefaultServer, opt) {
		diff["DefaultServer"] = []interface{}{s.DefaultServer, t.DefaultServer}
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

	if !s.EmailAlert.Equal(*t.EmailAlert, opt) {
		diff["EmailAlert"] = []interface{}{s.EmailAlert, t.EmailAlert}
	}

	if s.Enabled != t.Enabled {
		diff["Enabled"] = []interface{}{s.Enabled, t.Enabled}
	}

	if !s.Errorloc302.Equal(*t.Errorloc302, opt) {
		diff["Errorloc302"] = []interface{}{s.Errorloc302, t.Errorloc302}
	}

	if !s.Errorloc303.Equal(*t.Errorloc303, opt) {
		diff["Errorloc303"] = []interface{}{s.Errorloc303, t.Errorloc303}
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

	if !s.ForcePersist.Equal(*t.ForcePersist, opt) {
		diff["ForcePersist"] = []interface{}{s.ForcePersist, t.ForcePersist}
	}

	if !s.Forwardfor.Equal(*t.Forwardfor, opt) {
		diff["Forwardfor"] = []interface{}{s.Forwardfor, t.Forwardfor}
	}

	if s.From != t.From {
		diff["From"] = []interface{}{s.From, t.From}
	}

	if !equalPointers(s.Fullconn, t.Fullconn) {
		diff["Fullconn"] = []interface{}{s.Fullconn, t.Fullconn}
	}

	if s.H1CaseAdjustBogusServer != t.H1CaseAdjustBogusServer {
		diff["H1CaseAdjustBogusServer"] = []interface{}{s.H1CaseAdjustBogusServer, t.H1CaseAdjustBogusServer}
	}

	if !s.HashType.Equal(*t.HashType, opt) {
		diff["HashType"] = []interface{}{s.HashType, t.HashType}
	}

	if s.HTTPBufferRequest != t.HTTPBufferRequest {
		diff["HTTPBufferRequest"] = []interface{}{s.HTTPBufferRequest, t.HTTPBufferRequest}
	}

	if !s.HTTPCheck.Equal(*t.HTTPCheck, opt) {
		diff["HTTPCheck"] = []interface{}{s.HTTPCheck, t.HTTPCheck}
	}

	if s.HTTPKeepAlive != t.HTTPKeepAlive {
		diff["HTTPKeepAlive"] = []interface{}{s.HTTPKeepAlive, t.HTTPKeepAlive}
	}

	if s.HTTPNoDelay != t.HTTPNoDelay {
		diff["HTTPNoDelay"] = []interface{}{s.HTTPNoDelay, t.HTTPNoDelay}
	}

	if s.HTTPServerClose != t.HTTPServerClose {
		diff["HTTPServerClose"] = []interface{}{s.HTTPServerClose, t.HTTPServerClose}
	}

	if s.HTTPUseHtx != t.HTTPUseHtx {
		diff["HTTPUseHtx"] = []interface{}{s.HTTPUseHtx, t.HTTPUseHtx}
	}

	if s.HTTPConnectionMode != t.HTTPConnectionMode {
		diff["HTTPConnectionMode"] = []interface{}{s.HTTPConnectionMode, t.HTTPConnectionMode}
	}

	if !equalPointers(s.HTTPKeepAliveTimeout, t.HTTPKeepAliveTimeout) {
		diff["HTTPKeepAliveTimeout"] = []interface{}{s.HTTPKeepAliveTimeout, t.HTTPKeepAliveTimeout}
	}

	if s.HTTPPretendKeepalive != t.HTTPPretendKeepalive {
		diff["HTTPPretendKeepalive"] = []interface{}{s.HTTPPretendKeepalive, t.HTTPPretendKeepalive}
	}

	if s.HTTPProxy != t.HTTPProxy {
		diff["HTTPProxy"] = []interface{}{s.HTTPProxy, t.HTTPProxy}
	}

	if !equalPointers(s.HTTPRequestTimeout, t.HTTPRequestTimeout) {
		diff["HTTPRequestTimeout"] = []interface{}{s.HTTPRequestTimeout, t.HTTPRequestTimeout}
	}

	if s.HTTPRestrictReqHdrNames != t.HTTPRestrictReqHdrNames {
		diff["HTTPRestrictReqHdrNames"] = []interface{}{s.HTTPRestrictReqHdrNames, t.HTTPRestrictReqHdrNames}
	}

	if s.HTTPReuse != t.HTTPReuse {
		diff["HTTPReuse"] = []interface{}{s.HTTPReuse, t.HTTPReuse}
	}

	if !equalPointers(s.HTTPSendNameHeader, t.HTTPSendNameHeader) {
		diff["HTTPSendNameHeader"] = []interface{}{s.HTTPSendNameHeader, t.HTTPSendNameHeader}
	}

	if !s.HttpchkParams.Equal(*t.HttpchkParams, opt) {
		diff["HttpchkParams"] = []interface{}{s.HttpchkParams, t.HttpchkParams}
	}

	if s.Httpclose != t.Httpclose {
		diff["Httpclose"] = []interface{}{s.Httpclose, t.Httpclose}
	}

	if !equalPointers(s.ID, t.ID) {
		diff["ID"] = []interface{}{s.ID, t.ID}
	}

	if !s.IgnorePersist.Equal(*t.IgnorePersist, opt) {
		diff["IgnorePersist"] = []interface{}{s.IgnorePersist, t.IgnorePersist}
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
		diff["MaxKeepAliveQueue"] = []interface{}{s.MaxKeepAliveQueue, t.MaxKeepAliveQueue}
	}

	if s.Mode != t.Mode {
		diff["Mode"] = []interface{}{s.Mode, t.Mode}
	}

	if !s.MysqlCheckParams.Equal(*t.MysqlCheckParams, opt) {
		diff["MysqlCheckParams"] = []interface{}{s.MysqlCheckParams, t.MysqlCheckParams}
	}

	if s.Name != t.Name {
		diff["Name"] = []interface{}{s.Name, t.Name}
	}

	if s.Nolinger != t.Nolinger {
		diff["Nolinger"] = []interface{}{s.Nolinger, t.Nolinger}
	}

	if !s.Originalto.Equal(*t.Originalto, opt) {
		diff["Originalto"] = []interface{}{s.Originalto, t.Originalto}
	}

	if s.Persist != t.Persist {
		diff["Persist"] = []interface{}{s.Persist, t.Persist}
	}

	if !s.PersistRule.Equal(*t.PersistRule, opt) {
		diff["PersistRule"] = []interface{}{s.PersistRule, t.PersistRule}
	}

	if !s.PgsqlCheckParams.Equal(*t.PgsqlCheckParams, opt) {
		diff["PgsqlCheckParams"] = []interface{}{s.PgsqlCheckParams, t.PgsqlCheckParams}
	}

	if s.PreferLastServer != t.PreferLastServer {
		diff["PreferLastServer"] = []interface{}{s.PreferLastServer, t.PreferLastServer}
	}

	if !equalPointers(s.QueueTimeout, t.QueueTimeout) {
		diff["QueueTimeout"] = []interface{}{s.QueueTimeout, t.QueueTimeout}
	}

	if !s.Redispatch.Equal(*t.Redispatch, opt) {
		diff["Redispatch"] = []interface{}{s.Redispatch, t.Redispatch}
	}

	if !equalPointers(s.Retries, t.Retries) {
		diff["Retries"] = []interface{}{s.Retries, t.Retries}
	}

	if s.RetryOn != t.RetryOn {
		diff["RetryOn"] = []interface{}{s.RetryOn, t.RetryOn}
	}

	if !equalPointers(s.ServerFinTimeout, t.ServerFinTimeout) {
		diff["ServerFinTimeout"] = []interface{}{s.ServerFinTimeout, t.ServerFinTimeout}
	}

	if s.ServerStateFileName != t.ServerStateFileName {
		diff["ServerStateFileName"] = []interface{}{s.ServerStateFileName, t.ServerStateFileName}
	}

	if !equalPointers(s.ServerTimeout, t.ServerTimeout) {
		diff["ServerTimeout"] = []interface{}{s.ServerTimeout, t.ServerTimeout}
	}

	if !s.SmtpchkParams.Equal(*t.SmtpchkParams, opt) {
		diff["SmtpchkParams"] = []interface{}{s.SmtpchkParams, t.SmtpchkParams}
	}

	if !s.Source.Equal(*t.Source, opt) {
		diff["Source"] = []interface{}{s.Source, t.Source}
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
		diff["SrvtcpkaCnt"] = []interface{}{s.SrvtcpkaCnt, t.SrvtcpkaCnt}
	}

	if !equalPointers(s.SrvtcpkaIdle, t.SrvtcpkaIdle) {
		diff["SrvtcpkaIdle"] = []interface{}{s.SrvtcpkaIdle, t.SrvtcpkaIdle}
	}

	if !equalPointers(s.SrvtcpkaIntvl, t.SrvtcpkaIntvl) {
		diff["SrvtcpkaIntvl"] = []interface{}{s.SrvtcpkaIntvl, t.SrvtcpkaIntvl}
	}

	if !s.StatsOptions.Equal(*t.StatsOptions, opt) {
		diff["StatsOptions"] = []interface{}{s.StatsOptions, t.StatsOptions}
	}

	if !s.StickTable.Equal(*t.StickTable, opt) {
		diff["StickTable"] = []interface{}{s.StickTable, t.StickTable}
	}

	if !equalPointers(s.TarpitTimeout, t.TarpitTimeout) {
		diff["TarpitTimeout"] = []interface{}{s.TarpitTimeout, t.TarpitTimeout}
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
		diff["TunnelTimeout"] = []interface{}{s.TunnelTimeout, t.TunnelTimeout}
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
		diff["Cond"] = []interface{}{s.Cond, t.Cond}
	}

	if !equalPointers(s.CondTest, t.CondTest) {
		diff["CondTest"] = []interface{}{s.CondTest, t.CondTest}
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
		diff["Cond"] = []interface{}{s.Cond, t.Cond}
	}

	if !equalPointers(s.CondTest, t.CondTest) {
		diff["CondTest"] = []interface{}{s.CondTest, t.CondTest}
	}

	return diff
}
