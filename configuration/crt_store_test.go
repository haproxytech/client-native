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
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	parser "github.com/haproxytech/client-native/v6/config-parser"
	"github.com/haproxytech/client-native/v6/config-parser/options"
	"github.com/haproxytech/client-native/v6/models"
)

func renderCrtStore(t *testing.T, config string, store *models.CrtStore) string {
	t.Helper()
	p, err := parser.New(options.String(config))
	require.NoError(t, err)
	require.NoError(t, SerializeCrtStore(p, store))
	return strings.TrimSpace(p.String())
}

// crt-base and key-base take a mandatory path: a store without them must
// render a bare section, not empty keywords that HAProxy rejects.
func TestSerializeCrtStoreWithoutBases(t *testing.T) {
	store := &models.CrtStore{CrtStoreBase: models.CrtStoreBase{Name: "cs"}}
	require.Equal(t, "crt-store cs", renderCrtStore(t, "crt-store cs\n", store))
}

// Clearing a base on edit removes its line instead of leaving the keyword.
func TestSerializeCrtStoreClearsBases(t *testing.T) {
	store := &models.CrtStore{CrtStoreBase: models.CrtStoreBase{Name: "cs", CrtBase: "/certs"}}
	out := renderCrtStore(t, "crt-store cs\n  crt-base /old\n  key-base /keys\n", store)
	require.Equal(t, "crt-store cs\n  crt-base /certs", out)
}

func TestSerializeCrtStoreRoundTrip(t *testing.T) {
	store := &models.CrtStore{CrtStoreBase: models.CrtStoreBase{Name: "cs", CrtBase: "/certs", KeyBase: "/keys"}}
	out := renderCrtStore(t, "crt-store cs\n", store)

	p, err := parser.New(options.String(out))
	require.NoError(t, err)
	got := &models.CrtStore{CrtStoreBase: models.CrtStoreBase{Name: "cs"}}
	require.NoError(t, ParseCrtStore(p, got))
	require.True(t, got.CrtStoreBase.Equal(store.CrtStoreBase))
}
