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

func TestGlobalQuickExit(t *testing.T) {
	p, err := parser.New(parser_options.String("global\n  quick-exit\n"))
	require.NoError(t, err)
	global, err := ParseGlobalSection(p)
	require.NoError(t, err)
	data, err := json.Marshal(global)
	require.NoError(t, err)
	var fields map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(data, &fields))
	require.Contains(t, fields, "debug_options")
	require.JSONEq(t, `{"quick_exit":true}`, string(fields["debug_options"]))

	p, err = parser.New(parser_options.String("global\n  maxconn 100\n"))
	require.NoError(t, err)
	global, err = ParseGlobalSection(p)
	require.NoError(t, err)
	data, err = json.Marshal(global)
	require.NoError(t, err)
	fields = nil
	require.NoError(t, json.Unmarshal(data, &fields))
	require.NotContains(t, fields, "debug_options")

	payload := `{"debug_options":{"quick_exit":true}}`
	var model models.Global
	require.NoError(t, json.Unmarshal([]byte(payload), &model))
	p, err = parser.New(parser_options.String("global\n"))
	require.NoError(t, err)
	require.NoError(t, SerializeGlobalSection(p, &model, &config_options.ConfigurationOptions{}))
	rendered := p.String()
	require.Contains(t, rendered, "\n  quick-exit\n")
	fresh, err := parser.New(parser_options.String(rendered))
	require.NoError(t, err)
	reparsed, err := ParseGlobalSection(fresh)
	require.NoError(t, err)
	data, err = json.Marshal(reparsed)
	require.NoError(t, err)
	require.JSONEq(t, payload, string(data))

	const config = "global\n  quick-exit\n"
	for _, removal := range []string{`{}`, `{"debug_options":{}}`} {
		t.Run(removal, func(t *testing.T) {
			p, err := parser.New(parser_options.String(config))
			require.NoError(t, err)
			var global models.Global
			require.NoError(t, json.Unmarshal([]byte(removal), &global))
			require.NoError(t, SerializeGlobalSection(p, &global, &config_options.ConfigurationOptions{}))
			require.NotContains(t, p.String(), "quick-exit")
		})
	}
}

func TestGlobalEmptyUpdateOmitsQuickExit(t *testing.T) {
	p, err := parser.New(parser_options.String("global\n"))
	require.NoError(t, err)
	require.NoError(t, SerializeGlobalSection(p, &models.Global{}, &config_options.ConfigurationOptions{}))
	require.NotContains(t, p.String(), "quick-exit")
}
