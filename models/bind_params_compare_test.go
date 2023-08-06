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
	"bytes"
	"encoding/gob"
	"testing"

	"github.com/brianvoe/gofakeit/v6"

	jsoniter "github.com/json-iterator/go"
)

func TestBindParamsEqual(t *testing.T) {
	faker := gofakeit.NewCrypto()
	gofakeit.SetGlobalFaker(faker)
	samples := []struct {
		a, b BindParams
	}{}
	for i := 0; i < 100; i++ {
		var sample BindParams
		var result BindParams
		err := gofakeit.Struct(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		buf := new(bytes.Buffer)
		enc := gob.NewEncoder(buf)
		err = enc.Encode(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		dec := gob.NewDecoder(buf)
		err = dec.Decode(&result)
		if err != nil {
			t.Errorf(err.Error())
		}

		samples = append(samples, struct {
			a, b BindParams
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Equal(sample.b)
		if !result {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected BindParams to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestBindParamsEqualFalse(t *testing.T) {
	faker := gofakeit.NewCrypto()
	gofakeit.SetGlobalFaker(faker)
	samples := []struct {
		a, b BindParams
	}{}
	for i := 0; i < 100; i++ {
		var sample BindParams
		var result BindParams
		err := gofakeit.Struct(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = gofakeit.Struct(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		result.AcceptNetscalerCip = sample.AcceptNetscalerCip + 1
		result.AcceptProxy = !sample.AcceptProxy
		result.Allow0rtt = !sample.Allow0rtt
		result.DeferAccept = !sample.DeferAccept
		result.ExposeFdListeners = !sample.ExposeFdListeners
		result.ForceSslv3 = !sample.ForceSslv3
		result.ForceTlsv10 = !sample.ForceTlsv10
		result.ForceTlsv11 = !sample.ForceTlsv11
		result.ForceTlsv12 = !sample.ForceTlsv12
		result.ForceTlsv13 = !sample.ForceTlsv13
		result.GenerateCertificates = !sample.GenerateCertificates
		result.Gid = sample.Gid + 1
		result.Maxconn = sample.Maxconn + 1
		result.Nice = sample.Nice + 1
		result.NoAlpn = !sample.NoAlpn
		result.NoCaNames = !sample.NoCaNames
		result.NoSslv3 = !sample.NoSslv3
		result.NoTLSTickets = !sample.NoTLSTickets
		result.NoTlsv10 = !sample.NoTlsv10
		result.NoTlsv11 = !sample.NoTlsv11
		result.NoTlsv12 = !sample.NoTlsv12
		result.NoTlsv13 = !sample.NoTlsv13
		result.PreferClientCiphers = !sample.PreferClientCiphers
		result.QuicForceRetry = !sample.QuicForceRetry
		result.Ssl = !sample.Ssl
		result.StrictSni = !sample.StrictSni
		result.TCPUserTimeout = Ptr(*sample.TCPUserTimeout + 1)
		result.Tfo = !sample.Tfo
		result.Transparent = !sample.Transparent
		result.V4v6 = !sample.V4v6
		result.V6only = !sample.V6only
		samples = append(samples, struct {
			a, b BindParams
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Equal(sample.b)
		if result {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected BindParams to be different, but it is not %s %s", a, b)
		}
	}
}

func TestBindParamsDiff(t *testing.T) {
	faker := gofakeit.NewCrypto()
	gofakeit.SetGlobalFaker(faker)
	samples := []struct {
		a, b BindParams
	}{}
	for i := 0; i < 100; i++ {
		var sample BindParams
		var result BindParams
		err := gofakeit.Struct(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		buf := new(bytes.Buffer)
		enc := gob.NewEncoder(buf)
		err = enc.Encode(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		dec := gob.NewDecoder(buf)
		err = dec.Decode(&result)
		if err != nil {
			t.Errorf(err.Error())
		}

		samples = append(samples, struct {
			a, b BindParams
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		if len(result) != 0 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected BindParams to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestBindParamsDiffFalse(t *testing.T) {
	faker := gofakeit.NewCrypto()
	gofakeit.SetGlobalFaker(faker)
	samples := []struct {
		a, b BindParams
	}{}
	for i := 0; i < 100; i++ {
		var sample BindParams
		var result BindParams
		err := gofakeit.Struct(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = gofakeit.Struct(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		result.AcceptNetscalerCip = sample.AcceptNetscalerCip + 1
		result.AcceptProxy = !sample.AcceptProxy
		result.Allow0rtt = !sample.Allow0rtt
		result.DeferAccept = !sample.DeferAccept
		result.ExposeFdListeners = !sample.ExposeFdListeners
		result.ForceSslv3 = !sample.ForceSslv3
		result.ForceTlsv10 = !sample.ForceTlsv10
		result.ForceTlsv11 = !sample.ForceTlsv11
		result.ForceTlsv12 = !sample.ForceTlsv12
		result.ForceTlsv13 = !sample.ForceTlsv13
		result.GenerateCertificates = !sample.GenerateCertificates
		result.Gid = sample.Gid + 1
		result.Maxconn = sample.Maxconn + 1
		result.Nice = sample.Nice + 1
		result.NoAlpn = !sample.NoAlpn
		result.NoCaNames = !sample.NoCaNames
		result.NoSslv3 = !sample.NoSslv3
		result.NoTLSTickets = !sample.NoTLSTickets
		result.NoTlsv10 = !sample.NoTlsv10
		result.NoTlsv11 = !sample.NoTlsv11
		result.NoTlsv12 = !sample.NoTlsv12
		result.NoTlsv13 = !sample.NoTlsv13
		result.PreferClientCiphers = !sample.PreferClientCiphers
		result.QuicForceRetry = !sample.QuicForceRetry
		result.Ssl = !sample.Ssl
		result.StrictSni = !sample.StrictSni
		result.TCPUserTimeout = Ptr(*sample.TCPUserTimeout + 1)
		result.Tfo = !sample.Tfo
		result.Transparent = !sample.Transparent
		result.V4v6 = !sample.V4v6
		result.V6only = !sample.V6only
		samples = append(samples, struct {
			a, b BindParams
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		if len(result) != 68 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected BindParams to be different in 68 cases, but it is not (%d) %s %s", len(result), a, b)
		}
	}
}
