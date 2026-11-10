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

	"github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/client-native/v6/config-parser"
	parser_options "github.com/haproxytech/client-native/v6/config-parser/options"
	"github.com/haproxytech/client-native/v6/models"
	"github.com/stretchr/testify/require"
)

func TestDefaultsSpopCheck(t *testing.T) {
	t.Run("parse", func(t *testing.T) {
		for _, tc := range []struct {
			name string
			line string
			want string
		}{
			{"option", "option spop-check", "enabled"},
			{"no option", "no option spop-check", "disabled"},
			{"absent", "option tcplog", ""},
		} {
			t.Run(tc.name, func(t *testing.T) {
				deflt := parseDefaultsSpopCheck(t, tc.line)
				require.Equal(t, tc.want, deflt.SpopCheck)
				require.Empty(t, deflt.AdvCheck)
			})
		}
	})

	t.Run("serialize transitions", func(t *testing.T) {
		p, err := parser.New(parser_options.String("defaults\n"))
		require.NoError(t, err)

		require.NoError(t, CreateEditSection(&models.DefaultsBase{Name: "defaults", SpopCheck: "enabled"}, parser.Defaults, parser.DefaultsSectionName(p), p, nil))
		assertConfigLineCount(t, p.String(), "option spop-check", 1)

		require.NoError(t, CreateEditSection(&models.DefaultsBase{Name: "defaults", SpopCheck: "disabled"}, parser.Defaults, parser.DefaultsSectionName(p), p, nil))
		assertConfigLineCount(t, p.String(), "option spop-check", 0)
		assertConfigLineCount(t, p.String(), "no option spop-check", 1)

		require.NoError(t, CreateEditSection(&models.DefaultsBase{Name: "defaults"}, parser.Defaults, parser.DefaultsSectionName(p), p, nil))
		assertConfigLineCount(t, p.String(), "option spop-check", 0)
		assertConfigLineCount(t, p.String(), "no option spop-check", 0)
	})

	t.Run("frontend is not registered", func(t *testing.T) {
		config := "frontend fe\n  option spop-check\n"
		p, err := parser.New(parser_options.String(config))
		require.NoError(t, err)
		require.False(t, p.HasParser(parser.Frontends, "option spop-check"))
		require.Contains(t, p.String(), "option spop-check")
	})

	t.Run("validation", func(t *testing.T) {
		for _, value := range []string{"", "enabled", "disabled"} {
			require.NoError(t, (&models.DefaultsBase{Name: "defaults", SpopCheck: value}).Validate(strfmt.Default))
		}
		require.Error(t, (&models.DefaultsBase{Name: "defaults", SpopCheck: "bogus"}).Validate(strfmt.Default))
	})
}

func parseDefaultsSpopCheck(t *testing.T, line string) *models.DefaultsBase {
	t.Helper()

	config := strings.Join([]string{"defaults", "  " + line, ""}, "\n")
	p, err := parser.New(parser_options.String(config))
	require.NoError(t, err)
	deflt := &models.DefaultsBase{}
	require.NoError(t, ParseSection(deflt, parser.Defaults, parser.DefaultsSectionName(p), p))
	return deflt
}
