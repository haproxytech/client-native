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

package storage

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// Certificate bundles include the private key: owner only.
func TestSSLStorage_PrivateKeyNotWorldReadable(t *testing.T) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(1), Subject: pkix.Name{CommonName: "test.local"},
		NotBefore: time.Now(), NotAfter: time.Now().Add(time.Hour),
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	require.NoError(t, err)
	kb, err := x509.MarshalECPrivateKey(key)
	require.NoError(t, err)
	var bundle bytes.Buffer
	require.NoError(t, pem.Encode(&bundle, &pem.Block{Type: "CERTIFICATE", Bytes: der}))
	require.NoError(t, pem.Encode(&bundle, &pem.Block{Type: "EC PRIVATE KEY", Bytes: kb}))

	pemData := bundle.String()
	s, err := New(t.TempDir(), SSLType)
	require.NoError(t, err)
	f, _, err := s.Create("site.pem", io.NopCloser(strings.NewReader(pemData)))
	require.NoError(t, err)
	fi, err := os.Stat(f)
	require.NoError(t, err)
	require.Zero(t, fi.Mode().Perm()&0o077, "created: private key readable by group/others: %v", fi.Mode().Perm())

	f, err = s.Replace("site.pem", pemData)
	require.NoError(t, err)
	fi, err = os.Stat(f)
	require.NoError(t, err)
	require.Zero(t, fi.Mode().Perm()&0o077, "replaced: private key readable by group/others: %v", fi.Mode().Perm())
}
