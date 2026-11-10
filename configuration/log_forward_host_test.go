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
//

package configuration

import (
	"testing"

	strfmt "github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/require"

	parser "github.com/haproxytech/client-native/v6/config-parser"
	parser_options "github.com/haproxytech/client-native/v6/config-parser/options"
	config_options "github.com/haproxytech/client-native/v6/configuration/options"
	"github.com/haproxytech/client-native/v6/models"
)

func parseLogForwardConfig(t *testing.T, config string) (parser.Parser, *models.LogForward) {
	t.Helper()
	p, err := parser.New(parser_options.String(config))
	require.NoError(t, err)
	lf := &models.LogForward{LogForwardBase: models.LogForwardBase{Name: "lf"}}
	require.NoError(t, ParseLogForward(p, lf))
	return p, lf
}

func TestLogForwardHostRoundTrip(t *testing.T) {
	for _, strategy := range []string{"append", "fill", "keep", "replace"} {
		t.Run(strategy, func(t *testing.T) {
			p, lf := parseLogForwardConfig(t, "log-forward lf\n  option host "+strategy+"\n")
			require.Equal(t, strategy, lf.Host)

			require.NoError(t, SerializeLogForwardSection(p, lf, &config_options.ConfigurationOptions{}))
			out := p.String()
			require.Contains(t, out, "option host "+strategy)

			fresh, err := parser.New(parser_options.String(out))
			require.NoError(t, err)
			got := &models.LogForward{LogForwardBase: models.LogForwardBase{Name: "lf"}}
			require.NoError(t, ParseLogForward(fresh, got))
			require.Equal(t, strategy, got.Host)
		})
	}
}

func TestLogForwardHostUpdateAndRemoval(t *testing.T) {
	p, lf := parseLogForwardConfig(t, "log-forward lf\n  option host append\n")
	lf.Host = "replace"
	require.NoError(t, SerializeLogForwardSection(p, lf, &config_options.ConfigurationOptions{}))
	out := p.String()
	require.Contains(t, out, "option host replace")
	require.NotContains(t, out, "option host append")

	lf.Host = ""
	require.NoError(t, SerializeLogForwardSection(p, lf, &config_options.ConfigurationOptions{}))
	require.NotContains(t, p.String(), "option host")
}

func TestLogForwardHostOmittedByDefault(t *testing.T) {
	p, lf := parseLogForwardConfig(t, "log-forward lf\n")
	require.Empty(t, lf.Host)
	require.NoError(t, SerializeLogForwardSection(p, lf, &config_options.ConfigurationOptions{}))
	require.NotContains(t, p.String(), "option host")
}

func TestLogForwardHostRejectsUnknownStrategy(t *testing.T) {
	lf := &models.LogForward{LogForwardBase: models.LogForwardBase{Name: "lf", Host: "prepend"}}
	err := lf.Validate(strfmt.Default)
	require.Error(t, err)
}

func TestLogForwardHostPreservesOtherFields(t *testing.T) {
	config := `log-forward lf
  option host keep
  option assume-rfc6587-ntf
  option dont-parse-log
  timeout client 10s
`
	p, lf := parseLogForwardConfig(t, config)
	require.Equal(t, "keep", lf.Host)
	require.True(t, lf.AssumeRfc6587Ntf)
	require.True(t, lf.DontParseLog)
	require.NotNil(t, lf.TimeoutClient)

	require.NoError(t, SerializeLogForwardSection(p, lf, &config_options.ConfigurationOptions{}))
	out := p.String()
	for _, line := range []string{"option host keep", "option assume-rfc6587-ntf", "option dont-parse-log", "timeout client"} {
		require.Contains(t, out, line)
	}
	require.Contains(t, out, "log-forward lf")

	_, got := parseLogForwardConfig(t, out)
	require.Equal(t, "keep", got.Host)
	require.True(t, got.AssumeRfc6587Ntf)
	require.True(t, got.DontParseLog)
	require.Equal(t, lf.TimeoutClient, got.TimeoutClient)
}
