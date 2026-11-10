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
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	parser "github.com/haproxytech/client-native/v6/config-parser"
	parser_options "github.com/haproxytech/client-native/v6/config-parser/options"
	config_options "github.com/haproxytech/client-native/v6/configuration/options"
	"github.com/haproxytech/client-native/v6/models"
)

func parseRingConfig(t *testing.T, config string) (parser.Parser, *models.Ring) {
	t.Helper()
	p, err := parser.New(parser_options.String(config))
	require.NoError(t, err)
	ring := &models.Ring{RingBase: models.RingBase{Name: "r"}}
	require.NoError(t, ParseRingSection(p, ring))
	return p, ring
}

func TestRingBackingFileRoundTrip(t *testing.T) {
	backing := filepath.Join(t.TempDir(), "ring.buf")
	p, ring := parseRingConfig(t, "ring r\n  backing-file "+backing+"\n")
	require.Equal(t, backing, ring.BackingFile)

	require.NoError(t, SerializeRingSection(p, ring, &config_options.ConfigurationOptions{}))
	out := p.String()
	require.Contains(t, out, "backing-file "+backing)

	_, got := parseRingConfig(t, out)
	require.Equal(t, backing, got.BackingFile)
}

func TestRingBackingFileUpdateAndRemoval(t *testing.T) {
	oldPath := "/tmp/old-ring"
	p, ring := parseRingConfig(t, "ring r\n  backing-file "+oldPath+"\n")

	newPath := "/tmp/new-ring"
	ring.BackingFile = newPath
	require.NoError(t, SerializeRingSection(p, ring, &config_options.ConfigurationOptions{}))
	out := p.String()
	require.Contains(t, out, "backing-file "+newPath)
	require.NotContains(t, out, oldPath)

	ring.BackingFile = ""
	require.NoError(t, SerializeRingSection(p, ring, &config_options.ConfigurationOptions{}))
	require.NotContains(t, p.String(), "backing-file")
}

func TestRingBackingFilePreservesOtherFields(t *testing.T) {
	backing := filepath.Join(t.TempDir(), "ring.buf")
	config := `ring r
  description test ring
  format rfc3164
  size 1000
  timeout connect 5s
  timeout server 5s
  server s1 127.0.0.1:514
  backing-file ` + backing + `
`
	p, ring := parseRingConfig(t, config)
	require.Equal(t, backing, ring.BackingFile)
	require.Equal(t, "test ring", ring.Description)
	require.Equal(t, "rfc3164", ring.Format)
	require.NotNil(t, ring.Size)
	require.NotNil(t, ring.TimeoutConnect)
	require.NotNil(t, ring.TimeoutServer)

	newBacking := filepath.Join(t.TempDir(), "ring-new.buf")
	ring.BackingFile = newBacking
	require.NoError(t, SerializeRingSection(p, ring, &config_options.ConfigurationOptions{}))
	out := p.String()
	require.Contains(t, out, "backing-file "+newBacking)
	require.NotContains(t, out, backing)
	require.Contains(t, out, "description test ring")
	require.Contains(t, out, "format rfc3164")
	require.Contains(t, out, "size 1000")
	require.Contains(t, out, "timeout connect")
	require.Contains(t, out, "timeout server")
	require.Contains(t, out, "server s1 127.0.0.1:514")

	_, got := parseRingConfig(t, out)
	require.Equal(t, newBacking, got.BackingFile)
	require.Equal(t, ring.Description, got.Description)
	require.Equal(t, ring.Format, got.Format)
	require.Equal(t, ring.Size, got.Size)
	require.Equal(t, ring.TimeoutConnect, got.TimeoutConnect)
	require.Equal(t, ring.TimeoutServer, got.TimeoutServer)
}
