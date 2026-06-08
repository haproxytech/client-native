// Copyright 2019 HAProxy Technologies
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

package test

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/haproxytech/client-native/v6/models"
	"github.com/stretchr/testify/require"
)

func healthcheckExpectation() map[string]models.Healthchecks {
	initStructuredExpected()
	res := StructuredToHealthcheckMap()
	// Add individual entries
	for _, vs := range res {
		for _, v := range vs {
			key := v.Name
			res[key] = models.Healthchecks{v}
		}
	}
	return res
}

func TestGetHealthchecks(t *testing.T) {
	m := make(map[string]models.Healthchecks)
	v, healthchecks, err := clientTest.GetHealthchecks("")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("version %v returned, expected %v", v, version)
	}
	m[""] = healthchecks
	checkHealthchecks(t, m)
}

func checkHealthchecks(t *testing.T, got map[string]models.Healthchecks) {
	exp := healthcheckExpectation()

	for k, v := range got {
		want, ok := exp[k]
		// It's ok if there's no expectation (empty configuration)
		if !ok {
			if len(v) > 0 {
				t.Errorf("Got healthchecks but expected none for key=%s", k)
			}
			continue
		}
		require.Equal(t, len(want), len(v), "k=%s", k)
		for _, g := range v {
			for _, w := range want {
				if g.Name == w.Name {
					require.True(t, g.HealthCheckBase.Equal(w.HealthCheckBase), "k=%s - diff %v", k, cmp.Diff(*g, *w))
					break
				}
			}
		}
	}
}

func TestCreateEditDeleteHealthcheck(t *testing.T) { //nolint:gocognit,gocyclo
	testCases := []struct {
		name     string
		typeName string
		build    func(name string) *models.HealthCheck
		edit     func(name string) *models.HealthCheck
	}{
		{
			name:     "httpchk",
			typeName: "httpchk",
			build: func(name string) *models.HealthCheck {
				return &models.HealthCheck{
					HealthCheckBase: models.HealthCheckBase{
						Name: name,
						Type: "httpchk",
						HttpchkParams: &models.HttpchkParams{
							Method:  "GET",
							URI:     "/health",
							Version: "HTTP/1.1",
						},
					},
				}
			},
			edit: func(name string) *models.HealthCheck {
				return &models.HealthCheck{
					HealthCheckBase: models.HealthCheckBase{
						Name: name,
						Type: "httpchk",
						HttpchkParams: &models.HttpchkParams{
							Method:  "POST",
							URI:     "/health",
							Version: "HTTP/1.1",
						},
					},
				}
			},
		},
		{
			name:     "mysql-check",
			typeName: "mysql-check",
			build: func(name string) *models.HealthCheck {
				return &models.HealthCheck{
					HealthCheckBase: models.HealthCheckBase{
						Name: name,
						Type: "mysql-check",
						MysqlCheckParams: &models.MysqlCheckParams{
							ClientVersion: "pre-41",
							Username:      "root",
						},
					},
				}
			},
			edit: func(name string) *models.HealthCheck {
				return &models.HealthCheck{
					HealthCheckBase: models.HealthCheckBase{
						Name: name,
						Type: "mysql-check",
						MysqlCheckParams: &models.MysqlCheckParams{
							ClientVersion: "post-41",
							Username:      "admin",
						},
					},
				}
			},
		},
		{
			name:     "pgsql-check",
			typeName: "pgsql-check",
			build: func(name string) *models.HealthCheck {
				return &models.HealthCheck{
					HealthCheckBase: models.HealthCheckBase{
						Name: name,
						Type: "pgsql-check",
						PgsqlCheckParams: &models.PgsqlCheckParams{
							Username: "pguser",
						},
					},
				}
			},
			edit: func(name string) *models.HealthCheck {
				return &models.HealthCheck{
					HealthCheckBase: models.HealthCheckBase{
						Name: name,
						Type: "pgsql-check",
						PgsqlCheckParams: &models.PgsqlCheckParams{
							Username: "pgadmin",
						},
					},
				}
			},
		},
		{
			name:     "smtpchk",
			typeName: "smtpchk",
			build: func(name string) *models.HealthCheck {
				return &models.HealthCheck{
					HealthCheckBase: models.HealthCheckBase{
						Name: name,
						Type: "smtpchk",
						SmtpchkParams: &models.SmtpchkParams{
							Hello:  "EHLO",
							Domain: "example.org",
						},
					},
				}
			},
			edit: func(name string) *models.HealthCheck {
				return &models.HealthCheck{
					HealthCheckBase: models.HealthCheckBase{
						Name: name,
						Type: "smtpchk",
						SmtpchkParams: &models.SmtpchkParams{
							Hello:  "HELO",
							Domain: "example.com",
						},
					},
				}
			},
		},
		{
			name:     "tcp-check",
			typeName: "tcp-check",
			build: func(name string) *models.HealthCheck {
				return &models.HealthCheck{
					HealthCheckBase: models.HealthCheckBase{
						Name: name,
						Type: "tcp-check",
					},
				}
			},
			edit: func(name string) *models.HealthCheck {
				return &models.HealthCheck{
					HealthCheckBase: models.HealthCheckBase{
						Name: name,
						Type: "tcp-check",
					},
				}
			},
		},
		{
			name:     "ssl-hello-chk",
			typeName: "ssl-hello-chk",
			build: func(name string) *models.HealthCheck {
				return &models.HealthCheck{
					HealthCheckBase: models.HealthCheckBase{
						Name: name,
						Type: "ssl-hello-chk",
					},
				}
			},
			edit: func(name string) *models.HealthCheck {
				return &models.HealthCheck{
					HealthCheckBase: models.HealthCheckBase{
						Name: name,
						Type: "ssl-hello-chk",
					},
				}
			},
		},
		{
			name:     "redis-check",
			typeName: "redis-check",
			build: func(name string) *models.HealthCheck {
				return &models.HealthCheck{
					HealthCheckBase: models.HealthCheckBase{
						Name: name,
						Type: "redis-check",
					},
				}
			},
			edit: func(name string) *models.HealthCheck {
				return &models.HealthCheck{
					HealthCheckBase: models.HealthCheckBase{
						Name: name,
						Type: "redis-check",
					},
				}
			},
		},
		{
			name:     "ldap-check",
			typeName: "ldap-check",
			build: func(name string) *models.HealthCheck {
				return &models.HealthCheck{
					HealthCheckBase: models.HealthCheckBase{
						Name: name,
						Type: "ldap-check",
					},
				}
			},
			edit: func(name string) *models.HealthCheck {
				return &models.HealthCheck{
					HealthCheckBase: models.HealthCheckBase{
						Name: name,
						Type: "ldap-check",
					},
				}
			},
		},
		{
			name:     "spop-check",
			typeName: "spop-check",
			build: func(name string) *models.HealthCheck {
				return &models.HealthCheck{
					HealthCheckBase: models.HealthCheckBase{
						Name: name,
						Type: "spop-check",
					},
				}
			},
			edit: func(name string) *models.HealthCheck {
				return &models.HealthCheck{
					HealthCheckBase: models.HealthCheckBase{
						Name: name,
						Type: "spop-check",
					},
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			name := fmt.Sprintf("created_%s_%d", tc.typeName, version)
			hc := tc.build(name)

			err := clientTest.CreateHealthcheck(hc, "", version)
			require.NoError(t, err)

			version++

			v, healthcheck, err := clientTest.GetHealthcheck(name, "")
			require.NoError(t, err)
			require.NotNil(t, healthcheck)
			require.True(t, hc.HealthCheckBase.Equal(healthcheck.HealthCheckBase), "diff %v", cmp.Diff(*hc, *healthcheck))
			require.Equal(t, version, v)

			// Edit
			hcEdit := tc.edit(name)
			err = clientTest.EditHealthcheck(name, hcEdit, "", version)
			require.NoError(t, err)

			version++

			v, healthcheck, err = clientTest.GetHealthcheck(name, "")
			require.NoError(t, err)
			require.NotNil(t, healthcheck)
			require.True(t, hcEdit.HealthCheckBase.Equal(healthcheck.HealthCheckBase), "diff %v", cmp.Diff(*hcEdit, *healthcheck))
			require.Equal(t, version, v)

			err = clientTest.DeleteHealthcheck(name, "", version)
			require.NoError(t, err)

			version++

			_, _, err = clientTest.GetHealthcheck(name, "")
			require.Error(t, err)
		})
	}
}
