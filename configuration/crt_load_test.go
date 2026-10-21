// Copyright 2026 HAProxy Technologies
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

package configuration

import (
	"strings"
	"testing"

	"github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/client-native/v6/config-parser"
	"github.com/haproxytech/client-native/v6/config-parser/options"
	"github.com/haproxytech/client-native/v6/config-parser/types"
	"github.com/haproxytech/client-native/v6/models"
	"github.com/stretchr/testify/require"
)

func TestCrtLoadRoundTrip(t *testing.T) {
	tests := []struct {
		name string
		line string
		load models.CrtLoad
	}{
		{
			name: "omitted options",
			line: "load crt test.crt",
			load: models.CrtLoad{Certificate: "test.crt"},
		},
		{
			name: "RSA dummy",
			line: "load crt test.crt generate-dummy on keytype RSA bits 4096 curves P-384",
			load: models.CrtLoad{Certificate: "test.crt", GenerateDummy: "enabled", Keytype: "RSA", Bits: new(int64(4096)), Curves: "P-384"},
		},
		{
			name: "ECDSA dummy",
			line: "load crt test.crt generate-dummy on keytype ECDSA curves P-256",
			load: models.CrtLoad{Certificate: "test.crt", GenerateDummy: "enabled", Keytype: "ECDSA", Curves: "P-256"},
		},
		{
			name: "ACME IPs and JWT",
			line: "load crt test.crt acme LE domains example.com,example.org ips 192.0.2.1,2001:db8::1 jwt on # {\"comment\":\"certificate options\"}",
			load: models.CrtLoad{
				Certificate: "test.crt",
				Acme:        "LE",
				Domains:     []string{"example.com", "example.org"},
				IPs:         []string{"192.0.2.1", "2001:db8::1"},
				Jwt:         "enabled",
				Metadata:    map[string]interface{}{"comment": "certificate options"},
			},
		},
		{
			name: "explicitly disabled",
			line: "load crt test.crt ocsp-update off jwt off generate-dummy off",
			load: models.CrtLoad{Certificate: "test.crt", OcspUpdate: "disabled", Jwt: "disabled", GenerateDummy: "disabled"},
		},
		{
			name: "explicit zero bits",
			line: "load crt test.crt generate-dummy on bits 0",
			load: models.CrtLoad{Certificate: "test.crt", GenerateDummy: "enabled", Bits: new(int64(0))},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NoError(t, tt.load.Validate(strfmt.Default))
			p, err := parser.New(options.String("crt-store test\n  " + tt.line + "\n"))
			require.NoError(t, err)
			loads, err := ParseCrtLoads("test", p)
			require.NoError(t, err)
			require.Equal(t, models.CrtLoads{&tt.load}, loads)

			p, err = parser.New(options.String("crt-store test\n"))
			require.NoError(t, err)
			require.NoError(t, p.Insert(parser.CrtStore, "test", "load", SerializeCrtLoad(&tt.load), -1))
			out := p.String()
			require.Equal(t, "crt-store test\n  "+tt.line, strings.TrimSpace(out))

			p, err = parser.New(options.String(out))
			require.NoError(t, err)
			loads, err = ParseCrtLoads("test", p)
			require.NoError(t, err)
			require.Equal(t, models.CrtLoads{&tt.load}, loads)

			cleared := &models.CrtLoad{Certificate: "test.crt"}
			require.NoError(t, p.Set(parser.CrtStore, "test", "load", SerializeCrtLoad(cleared), 0))
			out = p.String()
			require.Equal(t, "crt-store test\n  load crt test.crt", strings.TrimSpace(out))
			p, err = parser.New(options.String(out))
			require.NoError(t, err)
			loads, err = ParseCrtLoads("test", p)
			require.NoError(t, err)
			require.Equal(t, models.CrtLoads{cleared}, loads)
		})
	}
}

func TestCrtLoadBitsDoNotAlias(t *testing.T) {
	t.Run("parse", func(t *testing.T) {
		load := types.LoadCert{Certificate: "test.crt", Bits: new(int64(4096))}
		model := ParseCrtLoad(load)
		*model.Bits = 2048
		require.Equal(t, int64(4096), *load.Bits)
	})
	t.Run("serialize", func(t *testing.T) {
		model := &models.CrtLoad{Certificate: "test.crt", Bits: new(int64(4096))}
		serialized := SerializeCrtLoad(model)
		*model.Bits = 2048
		require.Equal(t, int64(4096), *serialized.Bits)
	})
}

func TestCrtLoadValidation(t *testing.T) {
	for _, load := range []models.CrtLoad{
		{Certificate: "test.crt", Jwt: "on"},
		{Certificate: "test.crt", GenerateDummy: "on"},
		{Certificate: "test.crt", Keytype: "rsa"},
	} {
		require.Error(t, load.Validate(strfmt.Default))
	}
}
