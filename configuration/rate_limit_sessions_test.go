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
	"github.com/haproxytech/client-native/v6/config-parser/types"
	"github.com/haproxytech/client-native/v6/models"
	"github.com/stretchr/testify/require"
)

func TestRateLimitSessions(t *testing.T) {
	t.Run("parse values", func(t *testing.T) {
		for _, tc := range []struct {
			name string
			line string
			want *int64
		}{
			{"value", "rate-limit sessions 10", rateLimitInt64Ptr(10)},
			{"explicit zero", "rate-limit sessions 0", rateLimitInt64Ptr(0)},
			{"absent", "maxconn 10", nil},
		} {
			t.Run(tc.name, func(t *testing.T) {
				defaults, frontend := parseRateLimitSections(t, tc.line)
				require.Equal(t, tc.want, defaults.RateLimitSessions)
				require.Equal(t, tc.want, frontend.RateLimitSessions)

				if tc.want == nil {
					require.NotContains(t, rateLimitSessionsJSON(t, defaults), "\"rate_limit_sessions\"")
					return
				}
				want, err := json.Marshal(tc.want)
				require.NoError(t, err)
				require.JSONEq(t, string(want), rateLimitSessionsJSON(t, defaults))
				require.JSONEq(t, string(want), rateLimitSessionsJSON(t, frontend))
			})
		}
	})

	t.Run("serialize and remove", func(t *testing.T) {
		p, err := parser.New(parser_options.String("defaults\nfrontend fe\n"))
		require.NoError(t, err)

		require.NoError(t, CreateEditSection(&models.DefaultsBase{Name: "defaults", RateLimitSessions: rateLimitInt64Ptr(10)}, parser.Defaults, parser.DefaultsSectionName(p), p, nil))
		require.NoError(t, CreateEditSection(&models.FrontendBase{Name: "fe", RateLimitSessions: rateLimitInt64Ptr(0)}, parser.Frontends, "fe", p, nil))
		require.Equal(t, 1, rateLimitLineCount(p.String(), "rate-limit sessions 10"))
		require.Equal(t, 1, rateLimitLineCount(p.String(), "rate-limit sessions 0"))

		reparsedDefaults, reparsedFrontend := parseRateLimitConfig(t, p.String())
		require.Equal(t, rateLimitInt64Ptr(10), reparsedDefaults.RateLimitSessions)
		require.Equal(t, rateLimitInt64Ptr(0), reparsedFrontend.RateLimitSessions)

		require.NoError(t, CreateEditSection(&models.DefaultsBase{Name: "defaults"}, parser.Defaults, parser.DefaultsSectionName(p), p, nil))
		require.NoError(t, CreateEditSection(&models.FrontendBase{Name: "fe"}, parser.Frontends, "fe", p, nil))
		require.NotContains(t, p.String(), "rate-limit sessions")
		reparsedDefaults, reparsedFrontend = parseRateLimitConfig(t, p.String())
		require.Nil(t, reparsedDefaults.RateLimitSessions)
		require.Nil(t, reparsedFrontend.RateLimitSessions)
	})

	t.Run("validation bounds", func(t *testing.T) {
		for _, value := range []int64{0, 10, 4294967295} {
			require.NoError(t, (&models.FrontendBase{Name: "fe", RateLimitSessions: &value}).Validate(strfmt.Default))
			require.NoError(t, (&models.DefaultsBase{Name: "defaults", RateLimitSessions: &value}).Validate(strfmt.Default))
		}
		for _, value := range []int64{-1, 4294967296} {
			require.Error(t, (&models.FrontendBase{Name: "fe", RateLimitSessions: &value}).Validate(strfmt.Default))
			require.Error(t, (&models.DefaultsBase{Name: "defaults", RateLimitSessions: &value}).Validate(strfmt.Default))
		}
	})

	t.Run("listen parser unchanged", func(t *testing.T) {
		p, err := parser.New(parser_options.String("listen li\n  rate-limit sessions 10\n"), parser_options.UseListenSectionParsers)
		require.NoError(t, err)
		require.True(t, p.HasParser(parser.Listen, "rate-limit sessions"))
		data, err := p.Get(parser.Listen, "li", "rate-limit sessions")
		require.NoError(t, err)
		require.Equal(t, int64(10), data.(*types.Int64C).Value)
	})
}

func parseRateLimitSections(t *testing.T, line string) (*models.DefaultsBase, *models.FrontendBase) {
	t.Helper()

	return parseRateLimitConfig(t, strings.Join([]string{
		"defaults",
		"  " + line,
		"frontend fe",
		"  " + line,
		"",
	}, "\n"))
}

func parseRateLimitConfig(t *testing.T, config string) (*models.DefaultsBase, *models.FrontendBase) {
	t.Helper()

	p, err := parser.New(parser_options.String(config))
	require.NoError(t, err)

	defaults := &models.DefaultsBase{}
	require.NoError(t, ParseSection(defaults, parser.Defaults, parser.DefaultsSectionName(p), p))
	frontend := &models.FrontendBase{}
	require.NoError(t, ParseSection(frontend, parser.Frontends, "fe", p))
	return defaults, frontend
}

func rateLimitInt64Ptr(value int64) *int64 {
	return &value
}

func rateLimitSessionsJSON(t *testing.T, section any) string {
	t.Helper()

	data, err := json.Marshal(section)
	require.NoError(t, err)
	var fields map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(data, &fields))
	return string(fields["rate_limit_sessions"])
}

func rateLimitLineCount(config, line string) int {
	count := 0
	for _, configLine := range strings.Split(config, "\n") {
		if strings.TrimSpace(configLine) == line {
			count++
		}
	}
	return count
}
