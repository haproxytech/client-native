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

func TestGlobalHardMaxrecord(t *testing.T) {
	for _, tc := range []struct {
		name string
		line string
		want string
	}{
		{"zero", "tune.ssl.hard-maxrecord 0", `{"hard_maxrecord":0}`},
		{"default record", "tune.ssl.hard-maxrecord 16384", `{"hard_maxrecord":16384}`},
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
			require.Contains(t, fields, "tune_ssl_options")
			require.JSONEq(t, tc.want, string(fields["tune_ssl_options"]))
		})
	}

	p, err := parser.New(parser_options.String("global\n  maxconn 100\n"))
	require.NoError(t, err)
	global, err := ParseGlobalSection(p)
	require.NoError(t, err)
	require.Nil(t, global.TuneSslOptions)

	for _, value := range []int64{0, 16384} {
		t.Run("round-trip "+strconv.FormatInt(value, 10), func(t *testing.T) {
			payload := `{"tune_ssl_options":{"hard_maxrecord":` + strconv.FormatInt(value, 10) + `}}`
			var model models.Global
			require.NoError(t, json.Unmarshal([]byte(payload), &model))
			p, err := parser.New(parser_options.String("global\n"))
			require.NoError(t, err)
			require.NoError(t, SerializeGlobalSection(p, &model, &config_options.ConfigurationOptions{}))
			rendered := p.String()
			require.Contains(t, rendered, "\n  tune.ssl.hard-maxrecord "+strconv.FormatInt(value, 10)+"\n")
			fresh, err := parser.New(parser_options.String(rendered))
			require.NoError(t, err)
			reparsed, err := ParseGlobalSection(fresh)
			require.NoError(t, err)
			data, err := json.Marshal(reparsed)
			require.NoError(t, err)
			require.JSONEq(t, payload, string(data))
		})
	}

	const config = "global\n  tune.ssl.hard-maxrecord 16384\n"
	for _, removal := range []string{`{}`, `{"tune_ssl_options":{}}`} {
		t.Run(removal, func(t *testing.T) {
			p, err := parser.New(parser_options.String(config))
			require.NoError(t, err)
			var global models.Global
			require.NoError(t, json.Unmarshal([]byte(removal), &global))
			require.NoError(t, SerializeGlobalSection(p, &global, &config_options.ConfigurationOptions{}))
			require.NotContains(t, p.String(), "hard-maxrecord")
		})
	}
}

func TestGlobalHardMaxrecordValidation(t *testing.T) {
	validate := func(value int64) error {
		return (&models.Global{GlobalBase: models.GlobalBase{
			TuneSslOptions: &models.TuneSslOptions{HardMaxrecord: &value},
		}}).Validate(strfmt.Default)
	}
	for _, valid := range []int64{0, 16384, 2147483647} {
		t.Run("valid "+strconv.FormatInt(valid, 10), func(t *testing.T) {
			require.NoError(t, validate(valid))
		})
	}
	for _, invalid := range []int64{-1, 2147483648} {
		t.Run("invalid "+strconv.FormatInt(invalid, 10), func(t *testing.T) {
			require.Error(t, validate(invalid))
		})
	}
}

func TestGlobalHardMaxrecordIndependentFromMaxrecord(t *testing.T) {
	const config = "global\n  tune.ssl.maxrecord 1460\n  tune.ssl.hard-maxrecord 16384\n"
	p, err := parser.New(parser_options.String(config))
	require.NoError(t, err)
	global, err := ParseGlobalSection(p)
	require.NoError(t, err)
	require.NotNil(t, global.TuneSslOptions)
	require.Equal(t, int64(1460), *global.TuneSslOptions.Maxrecord)
	require.Equal(t, int64(16384), *global.TuneSslOptions.HardMaxrecord)

	hard := int64(8192)
	require.NoError(t, SerializeGlobalSection(p, &models.Global{GlobalBase: models.GlobalBase{
		TuneSslOptions: &models.TuneSslOptions{HardMaxrecord: &hard},
	}}, &config_options.ConfigurationOptions{}))
	rendered := p.String()
	require.Contains(t, rendered, "\n  tune.ssl.hard-maxrecord 8192\n")
	require.NotContains(t, rendered, "tune.ssl.maxrecord")
}
