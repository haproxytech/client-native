// Copyright 2025 HAProxy Technologies
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
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	parser "github.com/haproxytech/client-native/v6/config-parser"
	parser_options "github.com/haproxytech/client-native/v6/config-parser/options"
	config_options "github.com/haproxytech/client-native/v6/configuration/options"
	"github.com/haproxytech/client-native/v6/models"
	"github.com/stretchr/testify/require"
)

func Test_serializeAcmeVars(t *testing.T) {
	tests := []struct {
		vars    map[string]string
		want    string
		wantErr bool
	}{
		{
			vars:    map[string]string{"foo": "bar", "ApiKey": "FEFF,==\""},
			want:    `"ApiKey=FEFF\,==\",foo=bar"`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got, err := serializeAcmeVars(tt.vars)
			if tt.wantErr != (err != nil) {
				t.Errorf("serializeAcmeVars() got error '%v', wantErr=%v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("serializeAcmeVars() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ParseAcmeVars(t *testing.T) {
	tests := []struct {
		vars string
		want map[string]string
	}{
		{
			vars: `"ApiKey=FEFF\,==\",foo=bar"`,
			want: map[string]string{"foo": "bar", "ApiKey": "FEFF,==\""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.vars, func(t *testing.T) {
			got := ParseAcmeVars(tt.vars)
			if !cmp.Equal(got, tt.want) {
				t.Errorf("parseAcmeVars() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func Test_AcmeEabOptionsRoundTrip(t *testing.T) {
	const payload = `{
		"name": "test",
		"directory": "https://acme.example.com/directory",
		"eab_key_id": "/etc/haproxy/eab.keyid",
		"eab_mac_key": "/etc/haproxy/eab.mackey",
		"eab_mac_alg": "HS512"
	}`
	var acmeProvider models.AcmeProvider
	require.NoError(t, json.Unmarshal([]byte(payload), &acmeProvider))
	p, err := parser.New(parser_options.String("acme test\n"))
	require.NoError(t, err)
	require.NoError(t, SerializeAcmeProvider(p, &acmeProvider, &config_options.ConfigurationOptions{}))
	rendered := p.String()
	require.Contains(t, rendered, "\n  eab-key-id /etc/haproxy/eab.keyid\n")
	require.Contains(t, rendered, "\n  eab-mac-key /etc/haproxy/eab.mackey\n")
	require.Contains(t, rendered, "\n  eab-mac-alg HS512\n")
	fresh, err := parser.New(parser_options.String(rendered))
	require.NoError(t, err)
	reparsed, err := ParseAcmeProvider(fresh, "test")
	require.NoError(t, err)
	data, err := json.Marshal(reparsed)
	require.NoError(t, err)
	require.JSONEq(t, payload, string(data))
}

func Test_AcmeEabOptionsDeletion(t *testing.T) {
	const config = `acme test
  eab-key-id /etc/haproxy/eab.keyid
  eab-mac-key /etc/haproxy/eab.mackey
  eab-mac-alg HS512
`
	for _, payload := range []string{`{"name":"test"}`} {
		t.Run(payload, func(t *testing.T) {
			p, err := parser.New(parser_options.String(config))
			require.NoError(t, err)
			var acmeProvider models.AcmeProvider
			require.NoError(t, json.Unmarshal([]byte(payload), &acmeProvider))
			require.NoError(t, SerializeAcmeProvider(p, &acmeProvider, &config_options.ConfigurationOptions{}))
			require.NotContains(t, p.String(), "eab-")
		})
	}
}
