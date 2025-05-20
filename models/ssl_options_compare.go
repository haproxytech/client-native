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

// Equal checks if two structs of type SslOptions are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b SslOptions
//	equal := a.Equal(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b SslOptions
//	equal := a.Equal(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s SslOptions) Equal(t SslOptions, opts ...Options) bool {
	opt := getOptions(opts...)

	if !CheckSameNilAndLen(s.SslEngines, t.SslEngines, opt) {
		return false
	} else {
		for i := range s.SslEngines {
			if !s.SslEngines[i].Equal(*t.SslEngines[i], opt) {
				return false
			}
		}
	}

	if s.AcmeScheduler != t.AcmeScheduler {
		return false
	}

	if s.CaBase != t.CaBase {
		return false
	}

	if s.CrtBase != t.CrtBase {
		return false
	}

	if s.DefaultBindCiphers != t.DefaultBindCiphers {
		return false
	}

	if s.DefaultBindCiphersuites != t.DefaultBindCiphersuites {
		return false
	}

	if s.DefaultBindClientSigalgs != t.DefaultBindClientSigalgs {
		return false
	}

	if s.DefaultBindCurves != t.DefaultBindCurves {
		return false
	}

	if s.DefaultBindOptions != t.DefaultBindOptions {
		return false
	}

	if s.DefaultBindSigalgs != t.DefaultBindSigalgs {
		return false
	}

	if s.DefaultServerCiphers != t.DefaultServerCiphers {
		return false
	}

	if s.DefaultServerCiphersuites != t.DefaultServerCiphersuites {
		return false
	}

	if s.DefaultServerClientSigalgs != t.DefaultServerClientSigalgs {
		return false
	}

	if s.DefaultServerCurves != t.DefaultServerCurves {
		return false
	}

	if s.DefaultServerOptions != t.DefaultServerOptions {
		return false
	}

	if s.DefaultServerSigalgs != t.DefaultServerSigalgs {
		return false
	}

	if s.DhParamFile != t.DhParamFile {
		return false
	}

	if s.IssuersChainPath != t.IssuersChainPath {
		return false
	}

	if s.LoadExtraFiles != t.LoadExtraFiles {
		return false
	}

	if s.Maxsslconn != t.Maxsslconn {
		return false
	}

	if s.Maxsslrate != t.Maxsslrate {
		return false
	}

	if s.ModeAsync != t.ModeAsync {
		return false
	}

	if s.Propquery != t.Propquery {
		return false
	}

	if s.Provider != t.Provider {
		return false
	}

	if s.ProviderPath != t.ProviderPath {
		return false
	}

	if !equalPointers(s.SecurityLevel, t.SecurityLevel) {
		return false
	}

	if s.ServerVerify != t.ServerVerify {
		return false
	}

	if s.SkipSelfIssuedCa != t.SkipSelfIssuedCa {
		return false
	}

	return true
}

// Diff checks if two structs of type SslOptions are equal
//
// By default empty maps and slices are equal to nil:
//
//	var a, b SslOptions
//	diff := a.Diff(b)
//
// For more advanced use case you can configure these options (default values are shown):
//
//	var a, b SslOptions
//	diff := a.Diff(b,Options{
//		NilSameAsEmpty: true,
//	})
func (s SslOptions) Diff(t SslOptions, opts ...Options) map[string][]interface{} {
	opt := getOptions(opts...)

	diff := make(map[string][]interface{})
	if !CheckSameNilAndLen(s.SslEngines, t.SslEngines, opt) {
		diff["SslEngines"] = []interface{}{s.SslEngines, t.SslEngines}
	} else {
		diff2 := make(map[string][]interface{})
		for i := range s.SslEngines {
			if !s.SslEngines[i].Equal(*t.SslEngines[i], opt) {
				diffSub := s.SslEngines[i].Diff(*t.SslEngines[i], opt)
				if len(diffSub) > 0 {
					diff2[strconv.Itoa(i)] = []interface{}{diffSub}
				}
			}
		}
		if len(diff2) > 0 {
			diff["SslEngines"] = []interface{}{diff2}
		}
	}

	if s.AcmeScheduler != t.AcmeScheduler {
		diff["AcmeScheduler"] = []interface{}{s.AcmeScheduler, t.AcmeScheduler}
	}

	if s.CaBase != t.CaBase {
		diff["CaBase"] = []interface{}{s.CaBase, t.CaBase}
	}

	if s.CrtBase != t.CrtBase {
		diff["CrtBase"] = []interface{}{s.CrtBase, t.CrtBase}
	}

	if s.DefaultBindCiphers != t.DefaultBindCiphers {
		diff["DefaultBindCiphers"] = []interface{}{s.DefaultBindCiphers, t.DefaultBindCiphers}
	}

	if s.DefaultBindCiphersuites != t.DefaultBindCiphersuites {
		diff["DefaultBindCiphersuites"] = []interface{}{s.DefaultBindCiphersuites, t.DefaultBindCiphersuites}
	}

	if s.DefaultBindClientSigalgs != t.DefaultBindClientSigalgs {
		diff["DefaultBindClientSigalgs"] = []interface{}{s.DefaultBindClientSigalgs, t.DefaultBindClientSigalgs}
	}

	if s.DefaultBindCurves != t.DefaultBindCurves {
		diff["DefaultBindCurves"] = []interface{}{s.DefaultBindCurves, t.DefaultBindCurves}
	}

	if s.DefaultBindOptions != t.DefaultBindOptions {
		diff["DefaultBindOptions"] = []interface{}{s.DefaultBindOptions, t.DefaultBindOptions}
	}

	if s.DefaultBindSigalgs != t.DefaultBindSigalgs {
		diff["DefaultBindSigalgs"] = []interface{}{s.DefaultBindSigalgs, t.DefaultBindSigalgs}
	}

	if s.DefaultServerCiphers != t.DefaultServerCiphers {
		diff["DefaultServerCiphers"] = []interface{}{s.DefaultServerCiphers, t.DefaultServerCiphers}
	}

	if s.DefaultServerCiphersuites != t.DefaultServerCiphersuites {
		diff["DefaultServerCiphersuites"] = []interface{}{s.DefaultServerCiphersuites, t.DefaultServerCiphersuites}
	}

	if s.DefaultServerClientSigalgs != t.DefaultServerClientSigalgs {
		diff["DefaultServerClientSigalgs"] = []interface{}{s.DefaultServerClientSigalgs, t.DefaultServerClientSigalgs}
	}

	if s.DefaultServerCurves != t.DefaultServerCurves {
		diff["DefaultServerCurves"] = []interface{}{s.DefaultServerCurves, t.DefaultServerCurves}
	}

	if s.DefaultServerOptions != t.DefaultServerOptions {
		diff["DefaultServerOptions"] = []interface{}{s.DefaultServerOptions, t.DefaultServerOptions}
	}

	if s.DefaultServerSigalgs != t.DefaultServerSigalgs {
		diff["DefaultServerSigalgs"] = []interface{}{s.DefaultServerSigalgs, t.DefaultServerSigalgs}
	}

	if s.DhParamFile != t.DhParamFile {
		diff["DhParamFile"] = []interface{}{s.DhParamFile, t.DhParamFile}
	}

	if s.IssuersChainPath != t.IssuersChainPath {
		diff["IssuersChainPath"] = []interface{}{s.IssuersChainPath, t.IssuersChainPath}
	}

	if s.LoadExtraFiles != t.LoadExtraFiles {
		diff["LoadExtraFiles"] = []interface{}{s.LoadExtraFiles, t.LoadExtraFiles}
	}

	if s.Maxsslconn != t.Maxsslconn {
		diff["Maxsslconn"] = []interface{}{s.Maxsslconn, t.Maxsslconn}
	}

	if s.Maxsslrate != t.Maxsslrate {
		diff["Maxsslrate"] = []interface{}{s.Maxsslrate, t.Maxsslrate}
	}

	if s.ModeAsync != t.ModeAsync {
		diff["ModeAsync"] = []interface{}{s.ModeAsync, t.ModeAsync}
	}

	if s.Propquery != t.Propquery {
		diff["Propquery"] = []interface{}{s.Propquery, t.Propquery}
	}

	if s.Provider != t.Provider {
		diff["Provider"] = []interface{}{s.Provider, t.Provider}
	}

	if s.ProviderPath != t.ProviderPath {
		diff["ProviderPath"] = []interface{}{s.ProviderPath, t.ProviderPath}
	}

	if !equalPointers(s.SecurityLevel, t.SecurityLevel) {
		diff["SecurityLevel"] = []interface{}{ValueOrNil(s.SecurityLevel), ValueOrNil(t.SecurityLevel)}
	}

	if s.ServerVerify != t.ServerVerify {
		diff["ServerVerify"] = []interface{}{s.ServerVerify, t.ServerVerify}
	}

	if s.SkipSelfIssuedCa != t.SkipSelfIssuedCa {
		diff["SkipSelfIssuedCa"] = []interface{}{s.SkipSelfIssuedCa, t.SkipSelfIssuedCa}
	}

	return diff
}

// Equal checks if two structs of type SslEngine are equal
//
//	var a, b SslEngine
//	equal := a.Equal(b)
//
// opts ...Options are ignored in this method
func (s SslEngine) Equal(t SslEngine, opts ...Options) bool {
	if !equalPointers(s.Algorithms, t.Algorithms) {
		return false
	}

	if !equalPointers(s.Name, t.Name) {
		return false
	}

	return true
}

// Diff checks if two structs of type SslEngine are equal
//
//	var a, b SslEngine
//	diff := a.Diff(b)
//
// opts ...Options are ignored in this method
func (s SslEngine) Diff(t SslEngine, opts ...Options) map[string][]interface{} {
	diff := make(map[string][]interface{})
	if !equalPointers(s.Algorithms, t.Algorithms) {
		diff["Algorithms"] = []interface{}{ValueOrNil(s.Algorithms), ValueOrNil(t.Algorithms)}
	}

	if !equalPointers(s.Name, t.Name) {
		diff["Name"] = []interface{}{ValueOrNil(s.Name), ValueOrNil(t.Name)}
	}

	return diff
}
