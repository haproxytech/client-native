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
	config_options "github.com/haproxytech/client-native/v6/configuration/options"
	"github.com/haproxytech/client-native/v6/models"
	"github.com/stretchr/testify/require"
)

func TestClientHsTimeout(t *testing.T) {
	t.Run("parse", func(t *testing.T) {
		defaults, frontend := parseClientHsSections(t, "timeout client-hs 5s", "timeout client 3s", "timeout client-fin 2s")
		assertClientHsTimeout(t, "defaults", defaults.ClientHsTimeout, clientHsInt64Ptr(5000))
		assertClientHsTimeout(t, "frontend", frontend.ClientHsTimeout, clientHsInt64Ptr(5000))
		assertClientHsTimeout(t, "defaults client", defaults.ClientTimeout, clientHsInt64Ptr(3000))
		assertClientHsTimeout(t, "frontend client", frontend.ClientTimeout, clientHsInt64Ptr(3000))
		assertClientHsTimeout(t, "defaults client-fin", defaults.ClientFinTimeout, clientHsInt64Ptr(2000))
		assertClientHsTimeout(t, "frontend client-fin", frontend.ClientFinTimeout, clientHsInt64Ptr(2000))
	})

	t.Run("zero and missing", func(t *testing.T) {
		deflt, frontend := parseClientHsSections(t, "timeout client-hs 0")
		assertClientHsTimeout(t, "defaults zero", deflt.ClientHsTimeout, clientHsInt64Ptr(0))
		assertClientHsTimeout(t, "frontend zero", frontend.ClientHsTimeout, clientHsInt64Ptr(0))

		deflt, frontend = parseClientHsSections(t, "timeout client 3s")
		assertClientHsTimeout(t, "defaults absent", deflt.ClientHsTimeout, nil)
		assertClientHsTimeout(t, "frontend absent", frontend.ClientHsTimeout, nil)
	})

	t.Run("serialize with suffix and reparse", func(t *testing.T) {
		p, err := parser.New(parser_options.String("defaults\nfrontend fe\n"))
		require.NoError(t, err)
		opt := &config_options.ConfigurationOptions{PreferredTimeSuffix: "s"}

		require.NoError(t, CreateEditSection(&models.DefaultsBase{Name: "defaults", ClientHsTimeout: clientHsInt64Ptr(5000), ClientTimeout: clientHsInt64Ptr(3000)}, parser.Defaults, parser.DefaultsSectionName(p), p, opt))
		require.NoError(t, CreateEditSection(&models.FrontendBase{Name: "fe", ClientHsTimeout: clientHsInt64Ptr(5000), ClientTimeout: clientHsInt64Ptr(3000)}, parser.Frontends, "fe", p, opt))

		config := p.String()
		assertConfigLineCount(t, config, "timeout client-hs 5s", 2)
		assertConfigLineCount(t, config, "timeout client 3s", 2)

		reparsed, err := parser.New(parser_options.String(config))
		require.NoError(t, err)
		deflt := &models.DefaultsBase{}
		require.NoError(t, ParseSection(deflt, parser.Defaults, parser.DefaultsSectionName(reparsed), reparsed))
		frontend := &models.FrontendBase{}
		require.NoError(t, ParseSection(frontend, parser.Frontends, "fe", reparsed))
		assertClientHsTimeout(t, "defaults reparse", deflt.ClientHsTimeout, clientHsInt64Ptr(5000))
		assertClientHsTimeout(t, "frontend reparse", frontend.ClientHsTimeout, clientHsInt64Ptr(5000))
	})

	t.Run("removal keeps other timeouts", func(t *testing.T) {
		p, err := parser.New(parser_options.String("defaults\nfrontend fe\n"))
		require.NoError(t, err)
		opt := &config_options.ConfigurationOptions{PreferredTimeSuffix: "s"}

		require.NoError(t, CreateEditSection(&models.DefaultsBase{Name: "defaults", ClientHsTimeout: clientHsInt64Ptr(5000), ClientTimeout: clientHsInt64Ptr(3000)}, parser.Defaults, parser.DefaultsSectionName(p), p, opt))
		require.NoError(t, CreateEditSection(&models.FrontendBase{Name: "fe", ClientHsTimeout: clientHsInt64Ptr(5000), ClientTimeout: clientHsInt64Ptr(3000)}, parser.Frontends, "fe", p, opt))
		require.NoError(t, CreateEditSection(&models.DefaultsBase{Name: "defaults", ClientHsTimeout: nil, ClientTimeout: clientHsInt64Ptr(3000)}, parser.Defaults, parser.DefaultsSectionName(p), p, opt))
		require.NoError(t, CreateEditSection(&models.FrontendBase{Name: "fe", ClientHsTimeout: nil, ClientTimeout: clientHsInt64Ptr(3000)}, parser.Frontends, "fe", p, opt))

		config := p.String()
		assertConfigLineCount(t, config, "timeout client-hs 5s", 0)
		assertConfigLineCount(t, config, "timeout client 3s", 2)
	})

	t.Run("validation", func(t *testing.T) {
		for _, value := range []int64{0, 5000} {
			require.NoError(t, (&models.FrontendBase{Name: "fe", ClientHsTimeout: &value}).Validate(strfmt.Default))
			require.NoError(t, (&models.DefaultsBase{Name: "defaults", ClientHsTimeout: &value}).Validate(strfmt.Default))
		}
		invalid := int64(-1)
		require.Error(t, (&models.FrontendBase{Name: "fe", ClientHsTimeout: &invalid}).Validate(strfmt.Default))
		require.Error(t, (&models.DefaultsBase{Name: "defaults", ClientHsTimeout: &invalid}).Validate(strfmt.Default))
	})
}

func parseClientHsSections(t *testing.T, lines ...string) (*models.DefaultsBase, *models.FrontendBase) {
	t.Helper()

	var b strings.Builder
	b.WriteString("defaults\n")
	for _, line := range lines {
		b.WriteString("  " + line + "\n")
	}
	b.WriteString("frontend fe\n")
	for _, line := range lines {
		b.WriteString("  " + line + "\n")
	}
	p, err := parser.New(parser_options.String(b.String()))
	require.NoError(t, err)

	defaults := &models.DefaultsBase{}
	require.NoError(t, ParseSection(defaults, parser.Defaults, parser.DefaultsSectionName(p), p))
	frontend := &models.FrontendBase{}
	require.NoError(t, ParseSection(frontend, parser.Frontends, "fe", p))
	return defaults, frontend
}

func clientHsInt64Ptr(value int64) *int64 {
	return &value
}

func assertClientHsTimeout(t *testing.T, label string, got, want *int64) {
	t.Helper()

	if got == nil || want == nil {
		require.Equal(t, want, got, label)
		return
	}
	require.Equal(t, *want, *got, label)
}
