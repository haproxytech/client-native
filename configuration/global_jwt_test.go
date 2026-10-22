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

package configuration

import (
	"encoding/json"
	"testing"

	parser "github.com/haproxytech/client-native/v6/config-parser"
	parser_options "github.com/haproxytech/client-native/v6/config-parser/options"
	config_options "github.com/haproxytech/client-native/v6/configuration/options"
	"github.com/haproxytech/client-native/v6/models"
	"github.com/stretchr/testify/require"
)

func TestGlobalJWTOptions(t *testing.T) {
	for _, tc := range []struct {
		name string
		line string
		want string
	}{
		{"allow none", "jwt.allow_none_alg", `{"allow_none_alg":true}`},
		{"algorithms", "jwt.decrypt_alg_list ALL:!ECDH-ES", `{"decrypt_alg_list":"ALL:!ECDH-ES"}`},
		{"encryption", "jwt.decrypt_enc_list A128GCM:A192GCM:A256GCM", `{"decrypt_enc_list":"A128GCM:A192GCM:A256GCM"}`},
		{"absent", "maxconn 100", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			p, err := parser.New(parser_options.String("global\n  " + tc.line + "\n"))
			require.NoError(t, err)
			global, err := ParseGlobalSection(p)
			require.NoError(t, err)
			data, err := json.Marshal(global)
			require.NoError(t, err)
			var fields map[string]json.RawMessage
			require.NoError(t, json.Unmarshal(data, &fields))
			if tc.want == "" {
				require.NotContains(t, fields, "jwt_options")
			} else {
				require.Contains(t, fields, "jwt_options")
				require.JSONEq(t, tc.want, string(fields["jwt_options"]))
			}
		})
	}
}

func TestGlobalJWTOptionsRoundTrip(t *testing.T) {
	const payload = `{"jwt_options":{"allow_none_alg":true,"decrypt_alg_list":"ALL:!ECDH-ES","decrypt_enc_list":"A128GCM:A192GCM:A256GCM"}}`
	var global models.Global
	require.NoError(t, json.Unmarshal([]byte(payload), &global))
	p, err := parser.New(parser_options.String("global\n"))
	require.NoError(t, err)
	require.NoError(t, SerializeGlobalSection(p, &global, &config_options.ConfigurationOptions{}))
	rendered := p.String()
	require.Contains(t, rendered, "\n  jwt.allow_none_alg\n")
	require.Contains(t, rendered, "\n  jwt.decrypt_alg_list ALL:!ECDH-ES\n")
	require.Contains(t, rendered, "\n  jwt.decrypt_enc_list A128GCM:A192GCM:A256GCM\n")
	fresh, err := parser.New(parser_options.String(rendered))
	require.NoError(t, err)
	reparsed, err := ParseGlobalSection(fresh)
	require.NoError(t, err)
	data, err := json.Marshal(reparsed)
	require.NoError(t, err)
	require.JSONEq(t, payload, string(data))
}

func TestGlobalJWTOptionsDeletion(t *testing.T) {
	const config = "global\n  jwt.allow_none_alg\n  jwt.decrypt_alg_list ALL:!ECDH-ES\n  jwt.decrypt_enc_list A128GCM:A192GCM:A256GCM\n"
	for _, payload := range []string{`{}`, `{"jwt_options":{}}`} {
		t.Run(payload, func(t *testing.T) {
			p, err := parser.New(parser_options.String(config))
			require.NoError(t, err)
			var global models.Global
			require.NoError(t, json.Unmarshal([]byte(payload), &global))
			require.NoError(t, SerializeGlobalSection(p, &global, &config_options.ConfigurationOptions{}))
			require.NotContains(t, p.String(), "jwt.")
		})
	}
}
