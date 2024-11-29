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

// Equal checks if two structs of type BindParams are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b BindParams
//	equal := a.Equal(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b BindParams
//	equal := a.Equal(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s BindParams) Equal(t BindParams, opts ...Options) bool {
	opt := getOptions(opts...)

	if s.AcceptNetscalerCip != t.AcceptNetscalerCip {
		return false
	}

	if s.AcceptProxy != t.AcceptProxy {
		return false
	}

	if s.Allow0rtt != t.Allow0rtt {
		return false
	}

	if s.Alpn != t.Alpn {
		return false
	}

	if s.Backlog != t.Backlog {
		return false
	}

	if s.CaIgnoreErr != t.CaIgnoreErr {
		return false
	}

	if s.CaSignFile != t.CaSignFile {
		return false
	}

	if s.CaSignPass != t.CaSignPass {
		return false
	}

	if s.CaVerifyFile != t.CaVerifyFile {
		return false
	}

	if s.Ciphers != t.Ciphers {
		return false
	}

	if s.Ciphersuites != t.Ciphersuites {
		return false
	}

	if s.ClientSigalgs != t.ClientSigalgs {
		return false
	}

	if s.CrlFile != t.CrlFile {
		return false
	}

	if s.CrtIgnoreErr != t.CrtIgnoreErr {
		return false
	}

	if s.CrtList != t.CrtList {
		return false
	}

	if s.Curves != t.Curves {
		return false
	}

	if !equalComparableSlice(s.DefaultCrtList, t.DefaultCrtList, opt) {
		return false
	}

	if s.DeferAccept != t.DeferAccept {
		return false
	}

	if s.Ecdhe != t.Ecdhe {
		return false
	}

	if s.ExposeFdListeners != t.ExposeFdListeners {
		return false
	}

	if s.ForceSslv3 != t.ForceSslv3 {
		return false
	}

	if s.ForceTlsv10 != t.ForceTlsv10 {
		return false
	}

	if s.ForceTlsv11 != t.ForceTlsv11 {
		return false
	}

	if s.ForceTlsv12 != t.ForceTlsv12 {
		return false
	}

	if s.ForceTlsv13 != t.ForceTlsv13 {
		return false
	}

	if s.GenerateCertificates != t.GenerateCertificates {
		return false
	}

	if s.Gid != t.Gid {
		return false
	}

	if s.Group != t.Group {
		return false
	}

	if s.GUIDPrefix != t.GUIDPrefix {
		return false
	}

	if s.ID != t.ID {
		return false
	}

	if s.Interface != t.Interface {
		return false
	}

	if s.Level != t.Level {
		return false
	}

	if s.Maxconn != t.Maxconn {
		return false
	}

	if s.Mode != t.Mode {
		return false
	}

	if s.Mss != t.Mss {
		return false
	}

	if s.Name != t.Name {
		return false
	}

	if s.Namespace != t.Namespace {
		return false
	}

	if s.Nbconn != t.Nbconn {
		return false
	}

	if s.Nice != t.Nice {
		return false
	}

	if s.NoAlpn != t.NoAlpn {
		return false
	}

	if s.NoCaNames != t.NoCaNames {
		return false
	}

	if s.NoSslv3 != t.NoSslv3 {
		return false
	}

	if s.NoTLSTickets != t.NoTLSTickets {
		return false
	}

	if s.NoTlsv10 != t.NoTlsv10 {
		return false
	}

	if s.NoTlsv11 != t.NoTlsv11 {
		return false
	}

	if s.NoTlsv12 != t.NoTlsv12 {
		return false
	}

	if s.NoTlsv13 != t.NoTlsv13 {
		return false
	}

	if s.Npn != t.Npn {
		return false
	}

	if s.PreferClientCiphers != t.PreferClientCiphers {
		return false
	}

	if s.Proto != t.Proto {
		return false
	}

	if s.QuicCcAlgo != t.QuicCcAlgo {
		return false
	}

	if s.QuicForceRetry != t.QuicForceRetry {
		return false
	}

	if s.QuicSocket != t.QuicSocket {
		return false
	}

	if s.SeverityOutput != t.SeverityOutput {
		return false
	}

	if s.Sigalgs != t.Sigalgs {
		return false
	}

	if s.Ssl != t.Ssl {
		return false
	}

	if s.SslCafile != t.SslCafile {
		return false
	}

	if s.SslCertificate != t.SslCertificate {
		return false
	}

	if s.SslMaxVer != t.SslMaxVer {
		return false
	}

	if s.SslMinVer != t.SslMinVer {
		return false
	}

	if s.Sslv3 != t.Sslv3 {
		return false
	}

	if s.StrictSni != t.StrictSni {
		return false
	}

	if !equalPointers(s.TCPUserTimeout, t.TCPUserTimeout) {
		return false
	}

	if s.Tfo != t.Tfo {
		return false
	}

	if s.Thread != t.Thread {
		return false
	}

	if s.TLSTicketKeys != t.TLSTicketKeys {
		return false
	}

	if s.Tlsv10 != t.Tlsv10 {
		return false
	}

	if s.Tlsv11 != t.Tlsv11 {
		return false
	}

	if s.Tlsv12 != t.Tlsv12 {
		return false
	}

	if s.Tlsv13 != t.Tlsv13 {
		return false
	}

	if s.Transparent != t.Transparent {
		return false
	}

	if s.UID != t.UID {
		return false
	}

	if s.User != t.User {
		return false
	}

	if s.V4v6 != t.V4v6 {
		return false
	}

	if s.V6only != t.V6only {
		return false
	}

	if s.Verify != t.Verify {
		return false
	}

	return true
}

// Diff checks if two structs of type BindParams are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b BindParams
//	diff := a.Diff(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b BindParams
//	diff := a.Diff(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s BindParams) Diff(t BindParams, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})
	if s.AcceptNetscalerCip != t.AcceptNetscalerCip {
		diff["AcceptNetscalerCip"] = []interface{}{s.AcceptNetscalerCip, t.AcceptNetscalerCip}
	}

	if s.AcceptProxy != t.AcceptProxy {
		diff["AcceptProxy"] = []interface{}{s.AcceptProxy, t.AcceptProxy}
	}

	if s.Allow0rtt != t.Allow0rtt {
		diff["Allow0rtt"] = []interface{}{s.Allow0rtt, t.Allow0rtt}
	}

	if s.Alpn != t.Alpn {
		diff["Alpn"] = []interface{}{s.Alpn, t.Alpn}
	}

	if s.Backlog != t.Backlog {
		diff["Backlog"] = []interface{}{s.Backlog, t.Backlog}
	}

	if s.CaIgnoreErr != t.CaIgnoreErr {
		diff["CaIgnoreErr"] = []interface{}{s.CaIgnoreErr, t.CaIgnoreErr}
	}

	if s.CaSignFile != t.CaSignFile {
		diff["CaSignFile"] = []interface{}{s.CaSignFile, t.CaSignFile}
	}

	if s.CaSignPass != t.CaSignPass {
		diff["CaSignPass"] = []interface{}{s.CaSignPass, t.CaSignPass}
	}

	if s.CaVerifyFile != t.CaVerifyFile {
		diff["CaVerifyFile"] = []interface{}{s.CaVerifyFile, t.CaVerifyFile}
	}

	if s.Ciphers != t.Ciphers {
		diff["Ciphers"] = []interface{}{s.Ciphers, t.Ciphers}
	}

	if s.Ciphersuites != t.Ciphersuites {
		diff["Ciphersuites"] = []interface{}{s.Ciphersuites, t.Ciphersuites}
	}

	if s.ClientSigalgs != t.ClientSigalgs {
		diff["ClientSigalgs"] = []interface{}{s.ClientSigalgs, t.ClientSigalgs}
	}

	if s.CrlFile != t.CrlFile {
		diff["CrlFile"] = []interface{}{s.CrlFile, t.CrlFile}
	}

	if s.CrtIgnoreErr != t.CrtIgnoreErr {
		diff["CrtIgnoreErr"] = []interface{}{s.CrtIgnoreErr, t.CrtIgnoreErr}
	}

	if s.CrtList != t.CrtList {
		diff["CrtList"] = []interface{}{s.CrtList, t.CrtList}
	}

	if s.Curves != t.Curves {
		diff["Curves"] = []interface{}{s.Curves, t.Curves}
	}

	if !equalComparableSlice(s.DefaultCrtList, t.DefaultCrtList, opt) {
		diff["DefaultCrtList"] = []interface{}{s.DefaultCrtList, t.DefaultCrtList}
	}

	if s.DeferAccept != t.DeferAccept {
		diff["DeferAccept"] = []interface{}{s.DeferAccept, t.DeferAccept}
	}

	if s.Ecdhe != t.Ecdhe {
		diff["Ecdhe"] = []interface{}{s.Ecdhe, t.Ecdhe}
	}

	if s.ExposeFdListeners != t.ExposeFdListeners {
		diff["ExposeFdListeners"] = []interface{}{s.ExposeFdListeners, t.ExposeFdListeners}
	}

	if s.ForceSslv3 != t.ForceSslv3 {
		diff["ForceSslv3"] = []interface{}{s.ForceSslv3, t.ForceSslv3}
	}

	if s.ForceTlsv10 != t.ForceTlsv10 {
		diff["ForceTlsv10"] = []interface{}{s.ForceTlsv10, t.ForceTlsv10}
	}

	if s.ForceTlsv11 != t.ForceTlsv11 {
		diff["ForceTlsv11"] = []interface{}{s.ForceTlsv11, t.ForceTlsv11}
	}

	if s.ForceTlsv12 != t.ForceTlsv12 {
		diff["ForceTlsv12"] = []interface{}{s.ForceTlsv12, t.ForceTlsv12}
	}

	if s.ForceTlsv13 != t.ForceTlsv13 {
		diff["ForceTlsv13"] = []interface{}{s.ForceTlsv13, t.ForceTlsv13}
	}

	if s.GenerateCertificates != t.GenerateCertificates {
		diff["GenerateCertificates"] = []interface{}{s.GenerateCertificates, t.GenerateCertificates}
	}

	if s.Gid != t.Gid {
		diff["Gid"] = []interface{}{s.Gid, t.Gid}
	}

	if s.Group != t.Group {
		diff["Group"] = []interface{}{s.Group, t.Group}
	}

	if s.GUIDPrefix != t.GUIDPrefix {
		diff["GUIDPrefix"] = []interface{}{s.GUIDPrefix, t.GUIDPrefix}
	}

	if s.ID != t.ID {
		diff["ID"] = []interface{}{s.ID, t.ID}
	}

	if s.Interface != t.Interface {
		diff["Interface"] = []interface{}{s.Interface, t.Interface}
	}

	if s.Level != t.Level {
		diff["Level"] = []interface{}{s.Level, t.Level}
	}

	if s.Maxconn != t.Maxconn {
		diff["Maxconn"] = []interface{}{s.Maxconn, t.Maxconn}
	}

	if s.Mode != t.Mode {
		diff["Mode"] = []interface{}{s.Mode, t.Mode}
	}

	if s.Mss != t.Mss {
		diff["Mss"] = []interface{}{s.Mss, t.Mss}
	}

	if s.Name != t.Name {
		diff["Name"] = []interface{}{s.Name, t.Name}
	}

	if s.Namespace != t.Namespace {
		diff["Namespace"] = []interface{}{s.Namespace, t.Namespace}
	}

	if s.Nbconn != t.Nbconn {
		diff["Nbconn"] = []interface{}{s.Nbconn, t.Nbconn}
	}

	if s.Nice != t.Nice {
		diff["Nice"] = []interface{}{s.Nice, t.Nice}
	}

	if s.NoAlpn != t.NoAlpn {
		diff["NoAlpn"] = []interface{}{s.NoAlpn, t.NoAlpn}
	}

	if s.NoCaNames != t.NoCaNames {
		diff["NoCaNames"] = []interface{}{s.NoCaNames, t.NoCaNames}
	}

	if s.NoSslv3 != t.NoSslv3 {
		diff["NoSslv3"] = []interface{}{s.NoSslv3, t.NoSslv3}
	}

	if s.NoTLSTickets != t.NoTLSTickets {
		diff["NoTLSTickets"] = []interface{}{s.NoTLSTickets, t.NoTLSTickets}
	}

	if s.NoTlsv10 != t.NoTlsv10 {
		diff["NoTlsv10"] = []interface{}{s.NoTlsv10, t.NoTlsv10}
	}

	if s.NoTlsv11 != t.NoTlsv11 {
		diff["NoTlsv11"] = []interface{}{s.NoTlsv11, t.NoTlsv11}
	}

	if s.NoTlsv12 != t.NoTlsv12 {
		diff["NoTlsv12"] = []interface{}{s.NoTlsv12, t.NoTlsv12}
	}

	if s.NoTlsv13 != t.NoTlsv13 {
		diff["NoTlsv13"] = []interface{}{s.NoTlsv13, t.NoTlsv13}
	}

	if s.Npn != t.Npn {
		diff["Npn"] = []interface{}{s.Npn, t.Npn}
	}

	if s.PreferClientCiphers != t.PreferClientCiphers {
		diff["PreferClientCiphers"] = []interface{}{s.PreferClientCiphers, t.PreferClientCiphers}
	}

	if s.Proto != t.Proto {
		diff["Proto"] = []interface{}{s.Proto, t.Proto}
	}

	if s.QuicCcAlgo != t.QuicCcAlgo {
		diff["QuicCcAlgo"] = []interface{}{s.QuicCcAlgo, t.QuicCcAlgo}
	}

	if s.QuicForceRetry != t.QuicForceRetry {
		diff["QuicForceRetry"] = []interface{}{s.QuicForceRetry, t.QuicForceRetry}
	}

	if s.QuicSocket != t.QuicSocket {
		diff["QuicSocket"] = []interface{}{s.QuicSocket, t.QuicSocket}
	}

	if s.SeverityOutput != t.SeverityOutput {
		diff["SeverityOutput"] = []interface{}{s.SeverityOutput, t.SeverityOutput}
	}

	if s.Sigalgs != t.Sigalgs {
		diff["Sigalgs"] = []interface{}{s.Sigalgs, t.Sigalgs}
	}

	if s.Ssl != t.Ssl {
		diff["Ssl"] = []interface{}{s.Ssl, t.Ssl}
	}

	if s.SslCafile != t.SslCafile {
		diff["SslCafile"] = []interface{}{s.SslCafile, t.SslCafile}
	}

	if s.SslCertificate != t.SslCertificate {
		diff["SslCertificate"] = []interface{}{s.SslCertificate, t.SslCertificate}
	}

	if s.SslMaxVer != t.SslMaxVer {
		diff["SslMaxVer"] = []interface{}{s.SslMaxVer, t.SslMaxVer}
	}

	if s.SslMinVer != t.SslMinVer {
		diff["SslMinVer"] = []interface{}{s.SslMinVer, t.SslMinVer}
	}

	if s.Sslv3 != t.Sslv3 {
		diff["Sslv3"] = []interface{}{s.Sslv3, t.Sslv3}
	}

	if s.StrictSni != t.StrictSni {
		diff["StrictSni"] = []interface{}{s.StrictSni, t.StrictSni}
	}

	if !equalPointers(s.TCPUserTimeout, t.TCPUserTimeout) {
		diff["TCPUserTimeout"] = []interface{}{ValueOrNil(s.TCPUserTimeout), ValueOrNil(t.TCPUserTimeout)}
	}

	if s.Tfo != t.Tfo {
		diff["Tfo"] = []interface{}{s.Tfo, t.Tfo}
	}

	if s.Thread != t.Thread {
		diff["Thread"] = []interface{}{s.Thread, t.Thread}
	}

	if s.TLSTicketKeys != t.TLSTicketKeys {
		diff["TLSTicketKeys"] = []interface{}{s.TLSTicketKeys, t.TLSTicketKeys}
	}

	if s.Tlsv10 != t.Tlsv10 {
		diff["Tlsv10"] = []interface{}{s.Tlsv10, t.Tlsv10}
	}

	if s.Tlsv11 != t.Tlsv11 {
		diff["Tlsv11"] = []interface{}{s.Tlsv11, t.Tlsv11}
	}

	if s.Tlsv12 != t.Tlsv12 {
		diff["Tlsv12"] = []interface{}{s.Tlsv12, t.Tlsv12}
	}

	if s.Tlsv13 != t.Tlsv13 {
		diff["Tlsv13"] = []interface{}{s.Tlsv13, t.Tlsv13}
	}

	if s.Transparent != t.Transparent {
		diff["Transparent"] = []interface{}{s.Transparent, t.Transparent}
	}

	if s.UID != t.UID {
		diff["UID"] = []interface{}{s.UID, t.UID}
	}

	if s.User != t.User {
		diff["User"] = []interface{}{s.User, t.User}
	}

	if s.V4v6 != t.V4v6 {
		diff["V4v6"] = []interface{}{s.V4v6, t.V4v6}
	}

	if s.V6only != t.V6only {
		diff["V6only"] = []interface{}{s.V6only, t.V6only}
	}

	if s.Verify != t.Verify {
		diff["Verify"] = []interface{}{s.Verify, t.Verify}
	}

	return diff
}
