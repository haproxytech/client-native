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
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/haproxytech/client-native/v6/configuration"
	"github.com/haproxytech/client-native/v6/models"
	"github.com/stretchr/testify/require"
)

func peerSectionExpectation() map[string]models.PeerSections {
	initStructuredExpected()
	res := StructuredToPeerSectionMap()
	// Add individual entries
	for _, vs := range res {
		for _, v := range vs {
			key := v.Name
			res[key] = models.PeerSections{v}
		}
	}
	return res
}

func TestGetPeerSections(t *testing.T) {
	m := make(map[string]models.PeerSections)

	v, peerSections, err := clientTest.GetPeerSections("")
	if err != nil {
		t.Error(err.Error())
	}

	if len(peerSections) != 1 {
		t.Errorf("%v peerSections returned, expected 1", len(peerSections))
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	m[""] = peerSections

	checkPeerSections(t, m)
}

func TestGetPeerSection(t *testing.T) {
	m := make(map[string]models.PeerSections)
	v, l, err := clientTest.GetPeerSection("mycluster", "")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("%v: Version not %v: %v", l.Name, v, version)
	}
	m["mycluster"] = models.PeerSections{l}
	checkPeerSections(t, m)

	_, err = l.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	_, _, err = clientTest.GetPeerSection("doesnotexist", "")
	if err == nil {
		t.Error("Should throw error, non existent peer section")
	}
}

func checkPeerSections(t *testing.T, got map[string]models.PeerSections) {
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

func TestCreateEditDeletePeerSection(t *testing.T) {
	tOut := int64(5)
	f := &models.PeerSection{
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
	}
	err := clientTest.CreatePeerSection(f, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, peerSection, err := clientTest.GetPeerSection("testcluster", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(peerSection, f) {
		fmt.Printf("Created peerSection: %v\n", peerSection)
		fmt.Printf("Given peerSection: %v\n", f)
		t.Error("Created peerSection not equal to given peerSection")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	err = clientTest.CreatePeerSection(f, "", version)
	if err == nil {
		t.Error("Should throw error peerSection already exists")
		version++
	}

	err = clientTest.DeletePeerSection("testcluster", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ := clientTest.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	err = clientTest.DeletePeerSection("testcluster", "", 999999)
	if err != nil {
		if confErr, ok := err.(*configuration.ConfError); ok {
			if !confErr.Is(configuration.ErrVersionMismatch) {
				t.Error("Should throw configuration.ErrVersionMismatch error")
			}
		} else {
			t.Error("Should throw configuration.ErrVersionMismatch error")
		}
	}
	_, _, err = clientTest.GetPeerSection("testcluster", "")
	if err == nil {
		t.Error("DeletePeerSection failed, peerSection testcluster still exists")
	}

	err = clientTest.DeletePeerSection("doesnotexist", "", version)
	if err == nil {
		t.Error("Should throw error, non existent peerSection")
		version++
	}
}
