// Copyright 2026 HAProxy Technologies
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package spoe

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// SPOE file names must not resolve outside SpoeDir.
func TestSpoeDelete_NameConfinedToSpoeDir(t *testing.T) {
	root := t.TempDir()
	spoeDir := filepath.Join(root, "spoe")
	victim := filepath.Join(root, "victim")
	require.NoError(t, os.WriteFile(victim, []byte("x"), 0o600))

	c, err := NewSpoe(Params{SpoeDir: spoeDir, TransactionDir: filepath.Join(root, "tr")})
	require.NoError(t, err)
	for _, name := range []string{"../victim", "..", ".", ""} {
		require.Error(t, c.Delete(name), "name=%q", name)
	}
	require.FileExists(t, victim, "file outside SpoeDir deleted via spoe file name")
	require.DirExists(t, spoeDir)
}
