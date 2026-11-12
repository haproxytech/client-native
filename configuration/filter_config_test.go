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
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/haproxytech/client-native/v6/configuration/options"
	"github.com/haproxytech/client-native/v6/models"
)

const filterTestConfig = `global
  maxconn 10

defaults named
  mode http

frontend fe
  mode http

backend be
  mode http
`

func newFilterTestClient(t *testing.T) (Configuration, string) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "haproxy.cfg")
	if err := os.WriteFile(path, []byte(filterTestConfig), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	c, err := New(
		context.Background(),
		options.ConfigurationFile(path),
		options.HAProxyBin("echo"),
		options.UseModelsValidation,
		options.TransactionsDir(filepath.Join(dir, "transactions")),
	)
	if err != nil {
		t.Fatalf("new client: %v", err)
	}
	return c, path
}

func mustGetFilterConfigs(t *testing.T, c Configuration, parentType, parentName string) (int64, models.FilterConfigs) {
	t.Helper()
	v, configs, err := c.GetFilterConfigs(parentType, parentName, "")
	require.NoError(t, err)
	return v, configs
}

func filterConfigClassNames(configs models.FilterConfigs) []string {
	classes := make([]string, 0, len(configs))
	for _, config := range configs {
		classes = append(classes, config.Class)
	}
	return classes
}

func TestFilterConfigCRUD(t *testing.T) {
	parents := []struct {
		parentType string
		parentName string
	}{
		{FrontendParentName, "fe"},
		{BackendParentName, "be"},
		{DefaultsParentName, "named"},
	}
	for _, parent := range parents {
		t.Run(parent.parentType+"/"+parent.parentName, func(t *testing.T) {
			c, _ := newFilterTestClient(t)
			parentType, parentName := parent.parentType, parent.parentName

			maxFwd := int64(0)
			trace := &models.FilterConfig{
				Class:        "trace",
				ID:           "A",
				Enabled:      true,
				TraceQuiet:   true,
				TraceMaxFwd:  &maxFwd,
				Metadata:     map[string]interface{}{"k": "v"},
				TraceHexdump: true,
			}
			v, _ := mustGetFilterConfigs(t, c, parentType, parentName)
			require.NoError(t, c.CreateFilterConfig(0, parentType, parentName, trace, "", v))

			cache := &models.FilterConfig{Class: "cache-store", ID: "C"}
			v, _ = mustGetFilterConfigs(t, c, parentType, parentName)
			require.NoError(t, c.CreateFilterConfig(1, parentType, parentName, cache, "", v))

			compression := &models.FilterConfig{Class: "http-comp-req", Enabled: true}
			v, _ = mustGetFilterConfigs(t, c, parentType, parentName)
			require.NoError(t, c.CreateFilterConfig(1, parentType, parentName, compression, "", v))

			v, configs := mustGetFilterConfigs(t, c, parentType, parentName)
			require.Equal(t, []string{"trace", "http-comp-req", "cache-store"}, filterConfigClassNames(configs))

			_, got, err := c.GetFilterConfig(1, parentType, parentName, "")
			require.NoError(t, err)
			require.Equal(t, "http-comp-req", got.Class)

			fcgi := &models.FilterConfig{Class: "fcgi", ID: "F"}
			require.NoError(t, c.EditFilterConfig(1, parentType, parentName, fcgi, "", v))
			_, got, err = c.GetFilterConfig(1, parentType, parentName, "")
			require.NoError(t, err)
			require.Equal(t, "fcgi", got.Class)

			v, _ = mustGetFilterConfigs(t, c, parentType, parentName)
			require.NoError(t, c.DeleteFilterConfig(0, parentType, parentName, "", v))
			v, configs = mustGetFilterConfigs(t, c, parentType, parentName)
			require.Equal(t, []string{"fcgi", "cache-store"}, filterConfigClassNames(configs))

			reversed := models.FilterConfigs{cache, fcgi}
			require.NoError(t, c.ReplaceFilterConfigs(parentType, parentName, reversed, "", v))
			v, configs = mustGetFilterConfigs(t, c, parentType, parentName)
			require.Equal(t, []string{"cache-store", "fcgi"}, filterConfigClassNames(configs))

			require.NoError(t, c.ReplaceFilterConfigs(parentType, parentName, models.FilterConfigs{}, "", v))
			v, configs = mustGetFilterConfigs(t, c, parentType, parentName)
			require.Empty(t, configs)

			require.NoError(t, c.CreateFilterConfig(0, parentType, parentName, trace, "", v))
			_, got, err = c.GetFilterConfig(0, parentType, parentName, "")
			require.NoError(t, err)
			require.True(t, got.Enabled)
			require.True(t, got.TraceQuiet)
			require.True(t, got.TraceHexdump)
			require.NotNil(t, got.TraceMaxFwd)
			require.Equal(t, int64(0), *got.TraceMaxFwd)
			require.Equal(t, map[string]interface{}{"k": "v"}, got.Metadata)

			v, _ = mustGetFilterConfigs(t, c, parentType, parentName)
			_, _, err = c.GetFilterConfig(99, parentType, parentName, "")
			require.Error(t, err)
			require.Error(t, c.EditFilterConfig(99, parentType, parentName, fcgi, "", v))
			require.Error(t, c.DeleteFilterConfig(99, parentType, parentName, "", v))
			require.Error(t, c.CreateFilterConfig(0, parentType, parentName, trace, "", v+100))
		})
	}

	t.Run("missing parent", func(t *testing.T) {
		c, _ := newFilterTestClient(t)
		_, _, err := c.GetFilterConfigs(FrontendParentName, "doesnotexist", "")
		require.Error(t, err)
	})

	t.Run("explicit transaction", func(t *testing.T) {
		c, _ := newFilterTestClient(t)
		v, _ := mustGetFilterConfigs(t, c, FrontendParentName, "fe")
		transaction, err := c.StartTransaction(v)
		require.NoError(t, err)

		trace := &models.FilterConfig{Class: "trace", ID: "T"}
		require.NoError(t, c.CreateFilterConfig(0, FrontendParentName, "fe", trace, transaction.ID, 0))
		_, configs, err := c.GetFilterConfigs(FrontendParentName, "fe", transaction.ID)
		require.NoError(t, err)
		require.Len(t, configs, 1)

		_, err = c.CommitTransaction(transaction.ID)
		require.NoError(t, err)
		_, configs, err = c.GetFilterConfigs(FrontendParentName, "fe", "")
		require.NoError(t, err)
		require.Len(t, configs, 1)
		require.Equal(t, "trace", configs[0].Class)
	})
}
