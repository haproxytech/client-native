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
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/haproxytech/client-native/v6/configuration"
	"github.com/haproxytech/client-native/v6/misc"
	"github.com/haproxytech/client-native/v6/models"
	"github.com/stretchr/testify/require"
)

func TestGetStructuredPeerSections(t *testing.T) {
	clientTest, filename, err := getTestClient()
	require.NoError(t, err)
	defer os.Remove(filename)
	version := int64(1)

	m := make(map[string]models.PeerSections)

	v, peerSections, err := clientTest.GetStructuredPeerSections("")
	require.NoError(t, err)

	if len(peerSections) != 1 {
		t.Errorf("%v peerSections returned, expected 1", len(peerSections))
	}

	require.Equal(t, version, v, "Version %v returned, expected %v", v, version)

	m[""] = peerSections

	checkStructuredPeerSections(t, m)
}

func TestGetStructuredPeerSection(t *testing.T) {
	clientTest, filename, err := getTestClient()
	require.NoError(t, err)
	defer os.Remove(filename)
	version := int64(1)

	m := make(map[string]models.PeerSections)
	v, l, err := clientTest.GetStructuredPeerSection("mycluster", "")
	require.NoError(t, err)
	require.Equal(t, version, v, "Version %v returned, expected %v", v, version)

	m["mycluster"] = models.PeerSections{l}
	checkStructuredPeerSections(t, m)

	_, err = l.MarshalBinary()
	require.NoError(t, err)

	_, _, err = clientTest.GetStructuredPeerSection("doesnotexist", "")
	require.Error(t, err, "Should throw error, non existent peer section")
}

func checkStructuredPeerSections(t *testing.T, got map[string]models.PeerSections) {
	exp := peerSectionExpectation()
	for k, v := range got {
		want, ok := exp[k]
		require.True(t, ok, "k=%s", k)
		require.Equal(t, len(want), len(v), "k=%s", k)
		for _, g := range v {
			for _, w := range want {
				if g.Name == w.Name {
					require.True(t, g.Equal(*w), "k=%s - diff %v", k, cmp.Diff(*g, *w))
					break
				}
			}
		}
	}
}

func TestCreateEditDeleteStructuredPeerSection(t *testing.T) {
	clientTest, filename, err := getTestClient()
	require.NoError(t, err)
	defer os.Remove(filename)
	version := int64(1)

	tOut := int64(5)
	ps := &models.PeerSection{
		PeerSectionBase: models.PeerSectionBase{
			Name:     "testcluster",
			Disabled: true,
			DefaultServer: &models.DefaultServer{
				ServerParams: models.ServerParams{
					Fall:  &tOut,
					Inter: &tOut,
				},
			},
			DefaultBind: &models.DefaultBind{
				BindParams: models.BindParams{
					Alpn:           "h2,http/1.1",
					Ssl:            true,
					SslCertificate: "/etc/haproxy/cluster.pem",
				},
			},
			Shards: 4,
		},
		Servers: map[string]models.Server{
			"webserv": {
				Address: "192.168.1.1",
				Name:    "webserv",
				ID:      misc.Ptr[int64](1234),
				Port:    misc.Ptr[int64](9200),
			},
		},
		Binds: map[string]models.Bind{
			"192.168.1.1:9200": {
				BindParams: models.BindParams{Name: "192.168.1.1:9200"},
				Address:    "192.168.1.1",
				Port:       misc.Ptr[int64](9200),
			},
		},
		PeerEntries: map[string]models.PeerEntry{
			"entry1": {
				Name:    "entry1",
				Address: misc.Ptr("192.168.1.1"),
				Port:    misc.Ptr[int64](9200),
				Shard:   4,
			},
		},
		LogTargetList: models.LogTargets{
			&models.LogTarget{
				Address:  "192.169.0.1",
				Facility: "mail",
				Global:   true,
			},
		},
	}
	err = clientTest.CreateStructuredPeerSection(ps, "", version)
	require.NoError(t, err)
	version++

	v, peerSection, err := clientTest.GetStructuredPeerSection("testcluster", "")
	require.NoError(t, err)
	require.True(t, peerSection.Equal(*ps), "peer_section=%s - diff %v", peerSection.Name, cmp.Diff(*peerSection, *ps))
	require.Equal(t, version, v, "Version %v returned, expected %v", v, version)

	err = clientTest.CreateStructuredPeerSection(ps, "", version)
	require.Error(t, err, "Should throw error peerSection already exists")

	err = clientTest.DeletePeerSection("testcluster", "", version)
	require.NoError(t, err)
	version++

	v, _ = clientTest.GetVersion("")
	require.Equal(t, version, v, "Version %v returned, expected %v", v, version)

	err = clientTest.DeletePeerSection("testcluster", "", 999999)
	require.Error(t, err, "Should throw error, non existent frontend")
	require.ErrorIs(t, err, configuration.ErrVersionMismatch, "Should throw configuration.ErrVersionMismatch error")

	_, _, err = clientTest.GetStructuredPeerSection("testcluster", "")
	require.Error(t, err, "DeletePeerSection failed, peerSection testcluster still exists")

	err = clientTest.DeletePeerSection("doesnotexist", "", version)
	require.Error(t, err, "Should throw error, non existent peerSection")
}
