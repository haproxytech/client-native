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

import "reflect"

// Equal checks if two structs of type SSLFrontUse are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b SSLFrontUse
//	equal := a.Equal(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b SSLFrontUse
//	equal := a.Equal(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s SSLFrontUse) Equal(t SSLFrontUse, opts ...Options) bool {
	opt := getOptions(opts...)

	if s.Allow0rtt != t.Allow0rtt {
		return false
	}

	if s.Alpn != t.Alpn {
		return false
	}

	if s.CaFile != t.CaFile {
		return false
	}

	if s.Certificate != t.Certificate {
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

	if s.Curves != t.Curves {
		return false
	}

	if s.Ecdhe != t.Ecdhe {
		return false
	}

	if s.Issuer != t.Issuer {
		return false
	}

	if s.Key != t.Key {
		return false
	}

	if !CheckSameNilAndLenMap[string](s.Metadata, t.Metadata, opt) {
		return false
	}

	for k, v := range s.Metadata {
		if !reflect.DeepEqual(t.Metadata[k], v) {
			return false
		}
	}

	if s.NoAlpn != t.NoAlpn {
		return false
	}

	if s.NoCaNames != t.NoCaNames {
		return false
	}

	if s.Npn != t.Npn {
		return false
	}

	if s.Ocsp != t.Ocsp {
		return false
	}

	if s.OcspUpdate != t.OcspUpdate {
		return false
	}

	if s.Sctl != t.Sctl {
		return false
	}

	if s.Sigalgs != t.Sigalgs {
		return false
	}

	if s.SslMaxVer != t.SslMaxVer {
		return false
	}

	if s.SslMinVer != t.SslMinVer {
		return false
	}

	if s.Verify != t.Verify {
		return false
	}

	return true
}

// Diff checks if two structs of type SSLFrontUse are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b SSLFrontUse
//	diff := a.Diff(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b SSLFrontUse
//	diff := a.Diff(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s SSLFrontUse) Diff(t SSLFrontUse, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})
	if s.Allow0rtt != t.Allow0rtt {
		diff["Allow0rtt"] = []interface{}{s.Allow0rtt, t.Allow0rtt}
	}

	if s.Alpn != t.Alpn {
		diff["Alpn"] = []interface{}{s.Alpn, t.Alpn}
	}

	if s.CaFile != t.CaFile {
		diff["CaFile"] = []interface{}{s.CaFile, t.CaFile}
	}

	if s.Certificate != t.Certificate {
		diff["Certificate"] = []interface{}{s.Certificate, t.Certificate}
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

	if s.Curves != t.Curves {
		diff["Curves"] = []interface{}{s.Curves, t.Curves}
	}

	if s.Ecdhe != t.Ecdhe {
		diff["Ecdhe"] = []interface{}{s.Ecdhe, t.Ecdhe}
	}

	if s.Issuer != t.Issuer {
		diff["Issuer"] = []interface{}{s.Issuer, t.Issuer}
	}

	if s.Key != t.Key {
		diff["Key"] = []interface{}{s.Key, t.Key}
	}

	if !CheckSameNilAndLenMap[string](s.Metadata, t.Metadata, opt) {
		diff["Metadata"] = []interface{}{s.Metadata, t.Metadata}
	}

	for k, v := range s.Metadata {
		if !reflect.DeepEqual(t.Metadata[k], v) {
			diff["Metadata"] = []interface{}{s.Metadata, t.Metadata}
		}
	}

	if s.NoAlpn != t.NoAlpn {
		diff["NoAlpn"] = []interface{}{s.NoAlpn, t.NoAlpn}
	}

	if s.NoCaNames != t.NoCaNames {
		diff["NoCaNames"] = []interface{}{s.NoCaNames, t.NoCaNames}
	}

	if s.Npn != t.Npn {
		diff["Npn"] = []interface{}{s.Npn, t.Npn}
	}

	if s.Ocsp != t.Ocsp {
		diff["Ocsp"] = []interface{}{s.Ocsp, t.Ocsp}
	}

	if s.OcspUpdate != t.OcspUpdate {
		diff["OcspUpdate"] = []interface{}{s.OcspUpdate, t.OcspUpdate}
	}

	if s.Sctl != t.Sctl {
		diff["Sctl"] = []interface{}{s.Sctl, t.Sctl}
	}

	if s.Sigalgs != t.Sigalgs {
		diff["Sigalgs"] = []interface{}{s.Sigalgs, t.Sigalgs}
	}

	if s.SslMaxVer != t.SslMaxVer {
		diff["SslMaxVer"] = []interface{}{s.SslMaxVer, t.SslMaxVer}
	}

	if s.SslMinVer != t.SslMinVer {
		diff["SslMinVer"] = []interface{}{s.SslMinVer, t.SslMinVer}
	}

	if s.Verify != t.Verify {
		diff["Verify"] = []interface{}{s.Verify, t.Verify}
	}

	return diff
}
