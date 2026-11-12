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
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGlobalFilterAPI(t *testing.T) {
	c, path := newFilterTestClient(t)

	v, global, err := c.GetGlobalConfiguration("")
	require.NoError(t, err)
	require.Empty(t, global.FilterAPI)

	for _, mode := range []string{"default", "downgrade-to-legacy"} {
		v, global, err = c.GetGlobalConfiguration("")
		require.NoError(t, err)
		global.FilterAPI = mode
		require.NoError(t, c.PushGlobalConfiguration(global, "", v))

		_, global, err = c.GetGlobalConfiguration("")
		require.NoError(t, err)
		require.Equal(t, mode, global.FilterAPI)

		content, err := os.ReadFile(path)
		require.NoError(t, err)
		require.Contains(t, string(content), "filter-api "+mode)
	}

	// Structured global delegates through the same option.
	_, structured, err := c.GetStructuredGlobalConfiguration("")
	require.NoError(t, err)
	require.Equal(t, "downgrade-to-legacy", structured.FilterAPI)

	v, structured, err = c.GetStructuredGlobalConfiguration("")
	require.NoError(t, err)
	structured.FilterAPI = "default"
	require.NoError(t, c.PushStructuredGlobalConfiguration(structured, "", v))
	_, structured, err = c.GetStructuredGlobalConfiguration("")
	require.NoError(t, err)
	require.Equal(t, "default", structured.FilterAPI)

	// Clearing removes the line.
	v, global, err = c.GetGlobalConfiguration("")
	require.NoError(t, err)
	global.FilterAPI = ""
	require.NoError(t, c.PushGlobalConfiguration(global, "", v))
	_, global, err = c.GetGlobalConfiguration("")
	require.NoError(t, err)
	require.Empty(t, global.FilterAPI)

	content, err := os.ReadFile(path)
	require.NoError(t, err)
	require.False(t, strings.Contains(string(content), "filter-api"))
}

func TestGlobalFilterAPIRejectsUnknownMode(t *testing.T) {
	c, _ := newFilterTestClient(t)
	v, global, err := c.GetGlobalConfiguration("")
	require.NoError(t, err)
	global.FilterAPI = "bogus"
	require.Error(t, c.PushGlobalConfiguration(global, "", v))
}
