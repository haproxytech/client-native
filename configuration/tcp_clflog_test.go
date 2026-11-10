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
	"strings"
	"testing"

	parser "github.com/haproxytech/client-native/v6/config-parser"
	parser_options "github.com/haproxytech/client-native/v6/config-parser/options"
	"github.com/haproxytech/client-native/v6/config-parser/parsers"
	"github.com/haproxytech/client-native/v6/config-parser/types"
	"github.com/haproxytech/client-native/v6/models"
	"github.com/stretchr/testify/require"
)

func TestTCPClflog(t *testing.T) {
	t.Run("parse states", func(t *testing.T) {
		for _, tc := range []struct {
			name      string
			line      string
			tcplog    bool
			tcpClflog bool
		}{
			{"clf", "option tcplog clf", true, true},
			{"plain", "option tcplog", true, false},
			{"no option", "no option tcplog", false, false},
			{"no option clf", "no option tcplog clf", false, false},
			{"absent", "option httplog", false, false},
		} {
			t.Run(tc.name, func(t *testing.T) {
				defaults, frontend := parseTCPClflogSections(t, tc.line)
				require.Equal(t, tc.tcplog, defaults.Tcplog)
				require.Equal(t, tc.tcpClflog, defaults.TCPClflog)
				require.Equal(t, tc.tcplog, frontend.Tcplog)
				require.Equal(t, tc.tcpClflog, frontend.TCPClflog)
			})
		}
	})

	t.Run("comments and no-option are preserved", func(t *testing.T) {
		config := strings.Join([]string{
			"defaults",
			"  option tcplog clf # clf comment",
			"frontend fe",
			"  # pre comment",
			"  no option tcplog # no option comment",
			"",
		}, "\n")
		p, err := parser.New(parser_options.String(config))
		require.NoError(t, err)

		data, err := p.Get(parser.Defaults, parser.DefaultsSectionName(p), "option tcplog")
		require.NoError(t, err)
		require.Equal(t, "clf comment", data.(*types.OptionTCPLog).Comment)
		require.True(t, data.(*types.OptionTCPLog).Clf)
		require.Contains(t, p.String(), "option tcplog clf # clf comment")

		data, err = p.Get(parser.Frontends, "fe", "option tcplog")
		require.NoError(t, err)
		require.Equal(t, "no option comment", data.(*types.OptionTCPLog).Comment)
		require.True(t, data.(*types.OptionTCPLog).NoOption)
		preComments, err := p.GetPreComments(parser.Frontends, "fe", "option tcplog")
		require.NoError(t, err)
		require.Equal(t, []string{"pre comment"}, preComments)
		require.Contains(t, p.String(), "no option tcplog # no option comment")
	})

	t.Run("typed round-trip states", func(t *testing.T) {
		for _, tc := range []struct {
			name      string
			line      string
			tcplog    bool
			tcpClflog bool
			wantLine  string
		}{
			{"plain", "option tcplog", true, false, "option tcplog"},
			{"clf", "option tcplog clf", true, true, "option tcplog clf"},
			{"disabled", "no option tcplog", false, false, ""},
			{"disabled clf", "no option tcplog clf", false, false, ""},
		} {
			t.Run(tc.name, func(t *testing.T) {
				parsed, reparsed, config := roundTripTCPClflogFrontend(t, tc.line)
				require.Equal(t, tc.tcplog, parsed.Tcplog)
				require.Equal(t, tc.tcpClflog, parsed.TCPClflog)
				require.Equal(t, tc.tcplog, reparsed.Tcplog)
				require.Equal(t, tc.tcpClflog, reparsed.TCPClflog)

				if tc.wantLine == "" {
					require.NotContains(t, config, "option tcplog")
					return
				}
				require.Contains(t, config, tc.wantLine)
			})
		}
	})

	t.Run("typed http clf unchanged", func(t *testing.T) {
		parsed, reparsed, config := roundTripTCPClflogFrontend(t, "option httplog clf")
		require.True(t, parsed.Clflog)
		require.False(t, parsed.Httplog)
		require.False(t, parsed.TCPClflog)
		require.True(t, reparsed.Clflog)
		require.False(t, reparsed.Httplog)
		require.False(t, reparsed.TCPClflog)
		require.Contains(t, config, "option httplog clf")
	})

	t.Run("independent from http clf", func(t *testing.T) {
		defaults, frontend := parseTCPClflogSections(t, "option httplog clf")
		require.False(t, defaults.Tcplog)
		require.False(t, defaults.TCPClflog)
		require.True(t, defaults.Clflog)
		require.False(t, frontend.Tcplog)
		require.False(t, frontend.TCPClflog)
		require.True(t, frontend.Clflog)

		defaults, frontend = parseTCPClflogSections(t, "option tcplog clf")
		require.True(t, defaults.Tcplog)
		require.True(t, defaults.TCPClflog)
		require.False(t, defaults.Clflog)
		require.True(t, frontend.Tcplog)
		require.True(t, frontend.TCPClflog)
		require.False(t, frontend.Clflog)
	})

	t.Run("transitions clf to plain to removed", func(t *testing.T) {
		p, err := parser.New(parser_options.String("frontend fe\n"))
		require.NoError(t, err)

		require.NoError(t, CreateEditSection(&models.FrontendBase{Name: "fe", Tcplog: true, TCPClflog: true}, parser.Frontends, "fe", p, nil))
		assertConfigLineCount(t, p.String(), "option tcplog clf", 1)
		assertConfigLineCount(t, p.String(), "option tcplog", 0)

		require.NoError(t, CreateEditSection(&models.FrontendBase{Name: "fe", Tcplog: true}, parser.Frontends, "fe", p, nil))
		assertConfigLineCount(t, p.String(), "option tcplog clf", 0)
		assertConfigLineCount(t, p.String(), "option tcplog", 1)

		require.NoError(t, CreateEditSection(&models.FrontendBase{Name: "fe"}, parser.Frontends, "fe", p, nil))
		assertConfigLineCount(t, p.String(), "option tcplog", 0)
	})

	t.Run("parser rejects invalid suffixes", func(t *testing.T) {
		for _, parts := range [][]string{
			{"option"},
			{"option", "tcplog", "bogus"},
			{"no", "option", "tcplog", "bogus"},
			{"no", "option"},
		} {
			o := &parsers.OptionTCPLog{}
			o.Init()
			_, err := o.Parse(strings.Join(parts, " "), parts, "")
			require.Error(t, err)
		}
	})
}

func parseTCPClflogSections(t *testing.T, line string) (*models.DefaultsBase, *models.FrontendBase) {
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

func roundTripTCPClflogFrontend(t *testing.T, line string) (*models.FrontendBase, *models.FrontendBase, string) {
	t.Helper()

	source, err := parser.New(parser_options.String("frontend fe\n  " + line + "\n"))
	require.NoError(t, err)
	parsed := &models.FrontendBase{}
	require.NoError(t, ParseSection(parsed, parser.Frontends, "fe", source))

	p, err := parser.New(parser_options.String("frontend fe\n"))
	require.NoError(t, err)
	require.NoError(t, CreateEditSection(parsed, parser.Frontends, "fe", p, nil))

	fresh, err := parser.New(parser_options.String(p.String()))
	require.NoError(t, err)
	reparsed := &models.FrontendBase{}
	require.NoError(t, ParseSection(reparsed, parser.Frontends, "fe", fresh))
	return parsed, reparsed, p.String()
}
