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

func TestGlobalKeyBase(t *testing.T) {
	p, err := parser.New(parser_options.String("global\n  key-base /etc/haproxy/keys\n"))
	require.NoError(t, err)
	global, err := ParseGlobalSection(p)
	require.NoError(t, err)
	data, err := json.Marshal(global)
	require.NoError(t, err)
	var fields map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(data, &fields))
	require.Contains(t, fields, "ssl_options")
	require.JSONEq(t, `{"key_base":"/etc/haproxy/keys"}`, string(fields["ssl_options"]))

	for _, config := range []string{"global\n  maxconn 100\n", "global\n"} {
		t.Run(config, func(t *testing.T) {
			p, err := parser.New(parser_options.String(config))
			require.NoError(t, err)
			global, err := ParseGlobalSection(p)
			require.NoError(t, err)
			data, err := json.Marshal(global)
			require.NoError(t, err)
			var fields map[string]json.RawMessage
			require.NoError(t, json.Unmarshal(data, &fields))
			require.NotContains(t, fields, "ssl_options")
		})
	}

	payload := `{"ssl_options":{"key_base":"/etc/haproxy/keys"}}`
	var model models.Global
	require.NoError(t, json.Unmarshal([]byte(payload), &model))
	p, err = parser.New(parser_options.String("global\n"))
	require.NoError(t, err)
	require.NoError(t, SerializeGlobalSection(p, &model, &config_options.ConfigurationOptions{}))
	rendered := p.String()
	require.Contains(t, rendered, "\n  key-base /etc/haproxy/keys\n")
	fresh, err := parser.New(parser_options.String(rendered))
	require.NoError(t, err)
	reparsed, err := ParseGlobalSection(fresh)
	require.NoError(t, err)
	data, err = json.Marshal(reparsed)
	require.NoError(t, err)
	require.JSONEq(t, payload, string(data))

	const config = "global\n  key-base /etc/haproxy/keys\n"
	for _, removal := range []string{`{}`, `{"ssl_options":{}}`} {
		t.Run(removal, func(t *testing.T) {
			p, err := parser.New(parser_options.String(config))
			require.NoError(t, err)
			var global models.Global
			require.NoError(t, json.Unmarshal([]byte(removal), &global))
			require.NoError(t, SerializeGlobalSection(p, &global, &config_options.ConfigurationOptions{}))
			require.NotContains(t, p.String(), "key-base")
		})
	}
}

func TestGlobalKeyBaseIndependentFromCrtStore(t *testing.T) {
	const config = "global\n  key-base /global/keys\n\ncrt-store cs\n  key-base /store/keys\n"
	p, err := parser.New(parser_options.String(config))
	require.NoError(t, err)
	global, err := ParseGlobalSection(p)
	require.NoError(t, err)
	require.NotNil(t, global.SslOptions)
	require.Equal(t, "/global/keys", global.SslOptions.KeyBase)

	store := &models.CrtStore{CrtStoreBase: models.CrtStoreBase{Name: "cs"}}
	require.NoError(t, ParseCrtStore(p, store))
	require.Equal(t, "/store/keys", store.KeyBase)

	require.NoError(t, SerializeGlobalSection(p, &models.Global{GlobalBase: models.GlobalBase{
		SslOptions: &models.SslOptions{KeyBase: "/other/keys"},
	}}, &config_options.ConfigurationOptions{}))
	require.Contains(t, p.String(), "\n  key-base /other/keys\n")
	require.Contains(t, p.String(), "crt-store cs\n  key-base /store/keys")
}
