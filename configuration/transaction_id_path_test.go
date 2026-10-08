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
package configuration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/haproxytech/client-native/v6/configuration/options"
	"github.com/stretchr/testify/require"
)

// Transaction IDs must not resolve to files outside TransactionDir.
func TestDeleteTransaction_IDConfinedToTransactionDir(t *testing.T) {
	root := t.TempDir()
	trDir := filepath.Join(root, "transactions")
	require.NoError(t, os.MkdirAll(trDir, 0o750))
	victim := filepath.Join(root, "victim")
	require.NoError(t, os.WriteFile(victim, []byte("x"), 0o600))

	tr := &Transaction{ConfigurationOptions: options.ConfigurationOptions{
		ConfigurationFile:      filepath.Join(root, "haproxy.cfg"),
		TransactionDir:         trDir,
		PersistentTransactions: true,
	}}
	require.Error(t, tr.DeleteTransaction("/../../victim"))
	require.FileExists(t, victim, "file outside TransactionDir deleted via transaction id")
}
