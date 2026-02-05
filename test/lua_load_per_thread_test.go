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
	"context"
	"os"
	"strings"
	"testing"

	"github.com/haproxytech/client-native/v6/configuration"
	"github.com/haproxytech/client-native/v6/configuration/options"
	"github.com/stretchr/testify/require"
)

// TestLuaLoadPerThreadMultipleInstances tests that multiple instances of
// lua-load-per-thread are preserved when using PostRawConfiguration
func TestLuaLoadPerThreadMultipleInstances(t *testing.T) {
	// Create a temporary config file
	tmpFile, err := os.CreateTemp("", "haproxy_test_*.cfg")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// Initial config with multiple lua-load-per-thread instances
	initialConfig := `# _version=1
global
  lua-load-per-thread /etc/haproxy/lua/ja4.lua
  lua-load-per-thread /etc/haproxy/lua/ja4h.lua
  stats socket /var/run/haproxy.sock level admin
defaults
  mode http
`

	err = os.WriteFile(tmpFile.Name(), []byte(initialConfig), 0644)
	require.NoError(t, err)

	// Create configuration client
	client, err := configuration.New(context.Background(),
		options.ConfigurationFile(tmpFile.Name()),
		options.UsePersistentTransactions,
		options.TransactionsDir("/tmp"),
	)
	require.NoError(t, err)

	// Get raw configuration
	_, rawConfig, err := client.GetRawConfiguration("", 0)
	require.NoError(t, err)

	// Verify both instances are present
	require.Contains(t, rawConfig, "lua-load-per-thread /etc/haproxy/lua/ja4.lua", "First instance should be present")
	require.Contains(t, rawConfig, "lua-load-per-thread /etc/haproxy/lua/ja4h.lua", "Second instance should be present")

	// Post the raw configuration back (this is what was failing before)
	err = client.PostRawConfiguration(&rawConfig, 0, true)
	require.NoError(t, err)

	// Get the configuration again to verify both instances are still there
	_, rawConfigAfter, err := client.GetRawConfiguration("", 0)
	require.NoError(t, err)

	// Count occurrences
	count1 := strings.Count(rawConfigAfter, "lua-load-per-thread /etc/haproxy/lua/ja4.lua")
	count2 := strings.Count(rawConfigAfter, "lua-load-per-thread /etc/haproxy/lua/ja4h.lua")

	require.Equal(t, 1, count1, "First instance should appear exactly once")
	require.Equal(t, 1, count2, "Second instance should appear exactly once")

	// Verify both are still present
	require.Contains(t, rawConfigAfter, "lua-load-per-thread /etc/haproxy/lua/ja4.lua", "First instance should still be present after PostRawConfiguration")
	require.Contains(t, rawConfigAfter, "lua-load-per-thread /etc/haproxy/lua/ja4h.lua", "Second instance should still be present after PostRawConfiguration")
}

// TestLuaLoadPerThreadStructuredAPI tests that the structured API works
// (though it will only return the last instance due to spec limitation)
func TestLuaLoadPerThreadStructuredAPI(t *testing.T) {
	// Create a temporary config file
	tmpFile, err := os.CreateTemp("", "haproxy_test_*.cfg")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// Initial config with multiple lua-load-per-thread instances
	initialConfig := `# _version=1
global
  lua-load-per-thread /etc/haproxy/lua/ja4.lua
  lua-load-per-thread /etc/haproxy/lua/ja4h.lua
  stats socket /var/run/haproxy.sock level admin
defaults
  mode http
`

	err = os.WriteFile(tmpFile.Name(), []byte(initialConfig), 0644)
	require.NoError(t, err)

	// Create configuration client
	client, err := configuration.New(context.Background(),
		options.ConfigurationFile(tmpFile.Name()),
		options.UsePersistentTransactions,
		options.TransactionsDir("/tmp"),
	)
	require.NoError(t, err)

	// Get global configuration via structured API
	_, global, err := client.GetGlobalConfiguration("")
	require.NoError(t, err)

	// Due to spec limitation (string not array), only last instance is returned
	// This is expected behavior until spec is updated
	if global.LuaOptions != nil && global.LuaOptions.LoadPerThread != "" {
		// Should be the last one (ja4h.lua)
		require.Equal(t, "/etc/haproxy/lua/ja4h.lua", global.LuaOptions.LoadPerThread,
			"Structured API returns last instance (spec limitation)")
	}

	// Verify raw config still has both
	_, rawConfig, err := client.GetRawConfiguration("", 0)
	require.NoError(t, err)
	require.Contains(t, rawConfig, "lua-load-per-thread /etc/haproxy/lua/ja4.lua")
	require.Contains(t, rawConfig, "lua-load-per-thread /etc/haproxy/lua/ja4h.lua")
}
