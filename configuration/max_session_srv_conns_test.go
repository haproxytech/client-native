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
	"strings"
	"testing"

	"github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/client-native/v6/config-parser"
	parser_options "github.com/haproxytech/client-native/v6/config-parser/options"
	"github.com/haproxytech/client-native/v6/models"
	"github.com/stretchr/testify/require"
)

func TestMaxSessionSrvConns(t *testing.T) {
	t.Run("parse values", func(t *testing.T) {
		for _, tc := range []struct {
			name string
			line string
			want *int64
		}{
			{"value", "max-session-srv-conns 5", maxSessionInt64Ptr(5)},
			{"explicit zero", "max-session-srv-conns 0", maxSessionInt64Ptr(0)},
			{"absent", "backlog 10", nil},
		} {
			t.Run(tc.name, func(t *testing.T) {
				defaults, frontend := parseMaxSessionSrvConnsSections(t, tc.line)
				require.Equal(t, tc.want, defaults.MaxSessionSrvConns)
				require.Equal(t, tc.want, frontend.MaxSessionSrvConns)

				if tc.want == nil {
					require.NotContains(t, maxSessionSrvConnsJSON(t, defaults), "max_session_srv_conns")
					require.NotContains(t, maxSessionSrvConnsJSON(t, frontend), "max_session_srv_conns")
					return
				}
				require.JSONEq(t, string(mustMaxSessionSrvConnsJSON(t, tc.want)), maxSessionSrvConnsJSON(t, defaults))
				require.JSONEq(t, string(mustMaxSessionSrvConnsJSON(t, tc.want)), maxSessionSrvConnsJSON(t, frontend))
			})
		}
	})

	t.Run("serialize and remove", func(t *testing.T) {
		p, err := parser.New(parser_options.String("defaults\nfrontend fe\n"))
		require.NoError(t, err)

		require.NoError(t, CreateEditSection(&models.DefaultsBase{MaxSessionSrvConns: maxSessionInt64Ptr(5)}, parser.Defaults, parser.DefaultsSectionName(p), p, nil))
		require.NoError(t, CreateEditSection(&models.FrontendBase{MaxSessionSrvConns: maxSessionInt64Ptr(0)}, parser.Frontends, "fe", p, nil))
		require.Equal(t, 1, maxSessionSrvConnsLineCount(p.String(), "max-session-srv-conns 5"))
		require.Equal(t, 1, maxSessionSrvConnsLineCount(p.String(), "max-session-srv-conns 0"))

		require.NoError(t, CreateEditSection(&models.DefaultsBase{}, parser.Defaults, parser.DefaultsSectionName(p), p, nil))
		require.NoError(t, CreateEditSection(&models.FrontendBase{}, parser.Frontends, "fe", p, nil))
		require.Equal(t, 0, maxSessionSrvConnsLineCount(p.String(), "max-session-srv-conns"))
	})

	t.Run("validation", func(t *testing.T) {
		for _, value := range []int64{0, 5} {
			require.NoError(t, (&models.FrontendBase{Name: "fe", MaxSessionSrvConns: &value}).Validate(strfmt.Default))
			require.NoError(t, (&models.DefaultsBase{Name: "defaults", MaxSessionSrvConns: &value}).Validate(strfmt.Default))
		}
		invalid := int64(-1)
		require.Error(t, (&models.FrontendBase{Name: "fe", MaxSessionSrvConns: &invalid}).Validate(strfmt.Default))
		require.Error(t, (&models.DefaultsBase{Name: "defaults", MaxSessionSrvConns: &invalid}).Validate(strfmt.Default))
	})

	t.Run("backend directive is not registered", func(t *testing.T) {
		p, err := parser.New(parser_options.String("backend be\n  max-session-srv-conns 5\n"))
		require.NoError(t, err)
		require.False(t, p.HasParser(parser.Backends, "max-session-srv-conns"))
	})
}

func parseMaxSessionSrvConnsSections(t *testing.T, line string) (*models.DefaultsBase, *models.FrontendBase) {
	t.Helper()

	config := strings.Join([]string{
		"defaults",
		"  " + line,
		"frontend fe",
		"  " + line,
		"",
	}, "\n")
	p, err := parser.New(parser_options.String(config))
	require.NoError(t, err)

	defaults := &models.DefaultsBase{}
	require.NoError(t, ParseSection(defaults, parser.Defaults, parser.DefaultsSectionName(p), p))
	frontend := &models.FrontendBase{}
	require.NoError(t, ParseSection(frontend, parser.Frontends, "fe", p))
	return defaults, frontend
}

func maxSessionInt64Ptr(value int64) *int64 {
	return &value
}

func mustMaxSessionSrvConnsJSON(t *testing.T, value *int64) []byte {
	t.Helper()
	data, err := json.Marshal(value)
	require.NoError(t, err)
	return data
}

func maxSessionSrvConnsJSON(t *testing.T, section any) string {
	t.Helper()

	data, err := json.Marshal(section)
	require.NoError(t, err)
	var fields map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(data, &fields))
	return string(fields["max_session_srv_conns"])
}

func maxSessionSrvConnsLineCount(config, line string) int {
	count := 0
	for _, configLine := range strings.Split(config, "\n") {
		trimmed := strings.TrimSpace(configLine)
		if trimmed == line || strings.HasPrefix(trimmed, line+" ") {
			count++
		}
	}
	return count
}
