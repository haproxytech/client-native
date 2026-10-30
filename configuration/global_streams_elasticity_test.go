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
	"strconv"
	"testing"

	"github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/client-native/v6/config-parser"
	parser_options "github.com/haproxytech/client-native/v6/config-parser/options"
	config_options "github.com/haproxytech/client-native/v6/configuration/options"
	"github.com/haproxytech/client-native/v6/models"
	"github.com/stretchr/testify/require"
)

func validateStreamsElasticity(value *int64) error {
	global := &models.Global{
		GlobalBase: models.GlobalBase{
			TuneOptions: &models.TuneOptions{
				StreamsElasticity: value,
			},
		},
	}
	if err := global.Validate(strfmt.Default); err != nil {
		return err
	}
	return ValidateGlobalSection(global)
}

func TestGlobalStreamElasticityValidation(t *testing.T) {
	t.Run("unset", func(t *testing.T) {
		require.NoError(t, validateStreamsElasticity(nil))
	})

	for _, valid := range []int64{0, 100, 120, 1000} {
		t.Run("valid "+strconv.FormatInt(valid, 10), func(t *testing.T) {
			require.NoError(t, validateStreamsElasticity(&valid))
		})
	}

	for _, invalid := range []int64{-10, 22, 99} {
		t.Run("invalid "+strconv.FormatInt(invalid, 10), func(t *testing.T) {
			require.Error(t, validateStreamsElasticity(&invalid))
		})
	}
}

func TestGlobalStreamElasticityOption(t *testing.T) {
	for _, tc := range []struct {
		name string
		line string
		want string
	}{
		{"zero", "tune.streams-elasticity 0", `{"streams_elasticity":0}`},
		{"minimum", "tune.streams-elasticity 100", `{"streams_elasticity":100}`},
		{"moderate", "tune.streams-elasticity 5000", `{"streams_elasticity":5000}`},
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
				require.NotContains(t, fields, "tune_options")
			} else {
				require.Contains(t, fields, "tune_options")
				require.JSONEq(t, tc.want, string(fields["tune_options"]))
			}
		})
	}
}

func TestGlobalStreamsElasticityRoundTrip(t *testing.T) {
	payload := `{"tune_options":{"streams_elasticity":120}}`
	var global models.Global
	require.NoError(t, json.Unmarshal([]byte(payload), &global))
	require.NoError(t, validateStreamsElasticity(global.TuneOptions.StreamsElasticity))
	p, err := parser.New(parser_options.String("global\n"))
	require.NoError(t, err)
	require.NoError(t, SerializeGlobalSection(p, &global, &config_options.ConfigurationOptions{}))
	rendered := p.String()
	require.Contains(t, rendered, "\n  tune.streams-elasticity 120\n")
	fresh, err := parser.New(parser_options.String(rendered))
	require.NoError(t, err)
	reparsed, err := ParseGlobalSection(fresh)
	require.NoError(t, err)
	data, err := json.Marshal(reparsed)
	require.NoError(t, err)
	require.JSONEq(t, payload, string(data))
}

func TestGlobalStreamsElasticityDeletion(t *testing.T) {
	const config = "global\n  tune.streams-elasticity 120\n"
	for _, payload := range []string{`{}`, `{"tune_options":{}}`} {
		t.Run(payload, func(t *testing.T) {
			p, err := parser.New(parser_options.String(config))
			require.NoError(t, err)
			var global models.Global
			require.NoError(t, json.Unmarshal([]byte(payload), &global))
			require.NoError(t, SerializeGlobalSection(p, &global, &config_options.ConfigurationOptions{}))
			require.NotContains(t, p.String(), "tune.streams-elasticity")
		})
	}
}
