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

	"github.com/stretchr/testify/require"

	parser "github.com/haproxytech/client-native/v6/config-parser"
	"github.com/haproxytech/client-native/v6/config-parser/options"
	"github.com/haproxytech/client-native/v6/models"
)

func newHealthCheckParser(t *testing.T, config string) parser.Parser {
	t.Helper()
	p, err := parser.New(options.String(config))
	require.NoError(t, err)
	return p
}

// The check types whose parameters are optional in HAProxy must render the
// bare "type <check>" line when no parameters are given, and that line must
// be read back with the same type.
func TestSerializeHealthCheckSectionTypeWithoutParams(t *testing.T) {
	for _, typ := range []string{"httpchk", "smtpchk", "mysql-check"} {
		t.Run(typ, func(t *testing.T) {
			p := newHealthCheckParser(t, "healthcheck hc\n")
			hc := &models.HealthCheck{HealthCheckBase: models.HealthCheckBase{Name: "hc", Type: typ}}
			require.NoError(t, SerializeHealthCheckSection(p, hc))

			out := p.String()
			require.Contains(t, out, "\n  type "+typ+"\n", out)

			got := &models.HealthCheck{HealthCheckBase: models.HealthCheckBase{Name: "hc"}}
			require.NoError(t, ParseHealthcheckSection(newHealthCheckParser(t, out), got))
			require.Equal(t, typ, got.Type)
		})
	}
}

// An empty type only removes the existing type line, it is not an error.
func TestSerializeHealthCheckSectionEmptyType(t *testing.T) {
	p := newHealthCheckParser(t, "healthcheck hc\n  type httpchk GET /\n")
	hc := &models.HealthCheck{HealthCheckBase: models.HealthCheckBase{Name: "hc"}}
	require.NoError(t, SerializeHealthCheckSection(p, hc))
	require.NotContains(t, p.String(), "type")
}

// pgsql-check has a mandatory user in HAProxy: a bare "type pgsql-check"
// is rejected instead of being written to the configuration.
func TestSerializeHealthCheckSectionPgsqlCheckRequiresUser(t *testing.T) {
	for name, params := range map[string]*models.PgsqlCheckParams{
		"nil params":     nil,
		"empty username": {},
	} {
		t.Run(name, func(t *testing.T) {
			p := newHealthCheckParser(t, "healthcheck hc\n")
			hc := &models.HealthCheck{HealthCheckBase: models.HealthCheckBase{
				Name: "hc", Type: "pgsql-check", PgsqlCheckParams: params,
			}}
			require.ErrorIs(t, SerializeHealthCheckSection(p, hc), ErrValidationError)
			require.NotContains(t, p.String(), "type")
		})
	}
}

func TestSerializeHealthCheckSectionUnknownType(t *testing.T) {
	p := newHealthCheckParser(t, "healthcheck hc\n")
	hc := &models.HealthCheck{HealthCheckBase: models.HealthCheckBase{Name: "hc", Type: "foo-check"}}
	require.ErrorIs(t, SerializeHealthCheckSection(p, hc), ErrValidationError)
	require.NotContains(t, p.String(), "type")
}
