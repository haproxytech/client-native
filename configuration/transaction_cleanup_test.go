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
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/haproxytech/client-native/v6/configuration/options"
	"github.com/stretchr/testify/require"
)

type cleanupEnv struct {
	c        *client
	trDir    string
	failFlag string
}

// newCleanupEnv uses a fake haproxy binary whose check fails while failFlag exists.
func newCleanupEnv(t *testing.T, opts ...options.ConfigurationOption) cleanupEnv {
	t.Helper()
	root := t.TempDir()
	failFlag := filepath.Join(root, "fail")
	bin := filepath.Join(root, "haproxy")
	script := "#!/bin/sh\n" +
		"case \"$*\" in *-v*) echo 'HAProxy version 3.2.0 2025/01/01'; exit 0;; esac\n" +
		"[ -e " + failFlag + " ] && exit 1\n" +
		"exit 0\n"
	require.NoError(t, os.WriteFile(bin, []byte(script), 0o700))
	cfg := filepath.Join(root, "haproxy.cfg")
	require.NoError(t, os.WriteFile(cfg, []byte("# _version=1\nglobal\n\ndefaults unnamed_defaults_1\n"), 0o600))
	trDir := filepath.Join(root, "transactions")
	opts = append([]options.ConfigurationOption{
		options.ConfigurationFile(cfg),
		options.HAProxyBin(bin),
		options.TransactionsDir(trDir),
	}, opts...)
	c, err := New(context.Background(), opts...)
	require.NoError(t, err)
	cl, ok := c.(*client)
	require.True(t, ok)
	return cleanupEnv{c: cl, trDir: trDir, failFlag: failFlag}
}

func (e cleanupEnv) setFail(t *testing.T) {
	t.Helper()
	require.NoError(t, os.WriteFile(e.failFlag, nil, 0o600))
}

func (e cleanupEnv) start(t *testing.T) string {
	t.Helper()
	v, err := e.c.GetVersion("")
	require.NoError(t, err)
	tr, err := e.c.StartTransaction(v)
	require.NoError(t, err)
	return tr.ID
}

func TestCommitFailureReleasesParser(t *testing.T) {
	modes := map[string][]options.ConfigurationOption{
		"non-persistent": nil,
		"persistent":     {options.UsePersistentTransactions},
		"skip-failed":    {options.SkipFailedTransactions},
	}
	for name, opts := range modes {
		t.Run(name+"/version mismatch", func(t *testing.T) {
			e := newCleanupEnv(t, opts...)
			first, outdated := e.start(t), e.start(t)
			_, err := e.c.CommitTransaction(first)
			require.NoError(t, err)
			_, err = e.c.CommitTransaction(outdated)
			require.Error(t, err)
			require.False(t, e.c.HasParser(outdated), "outdated transaction parser kept in memory")
		})
		t.Run(name+"/check fails", func(t *testing.T) {
			e := newCleanupEnv(t, opts...)
			e.setFail(t)
			id := e.start(t)
			_, err := e.c.CommitTransaction(id)
			require.Error(t, err)
			require.False(t, e.c.HasParser(id), "failed transaction parser kept in memory")
		})
	}
}
