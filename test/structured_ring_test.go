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

func TestGetStructuredRings(t *testing.T) {
	clientTest, filename, err := getTestClient()
	require.NoError(t, err)
	defer os.Remove(filename)
	version := int64(1)

	m := make(map[string]models.Rings)
	v, rings, err := clientTest.GetStructuredRings("")
	require.NoError(t, err)
	require.Equal(t, version, v, "Version %v returned, expected %v", v, version)

	m[""] = rings
	checkStructuredRings(t, m)
}

func TestGetStructuredRing(t *testing.T) {
	clientTest, filename, err := getTestClient()
	require.NoError(t, err)
	defer os.Remove(filename)
	version := int64(1)

	m := make(map[string]models.Rings)

	v, r, err := clientTest.GetStructuredRing("myring", "")
	require.NoError(t, err)
	require.Equal(t, version, v, "Version %v returned, expected %v", v, version)

	m["myring"] = models.Rings{r}
	checkStructuredRings(t, m)

	_, _, err = clientTest.GetRing("doesnotexist", "")
	require.Error(t, err, "should throw error, non existent rings section")
}

func checkStructuredRings(t *testing.T, got map[string]models.Rings) {
	exp := ringExpectation()

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

func TestCreateEditDeleteStructuredRing(t *testing.T) {
	clientTest, filename, err := getTestClient()
	require.NoError(t, err)
	defer os.Remove(filename)
	version := int64(1)

	maxlen := int64(1300)
	size := int64(32765)
	timeoutConnect := int64(5)
	timeoutServer := int64(10)

	r := &models.Ring{
		RingBase: models.RingBase{
			Name:           "created_ring",
			Description:    "My local buffer 2",
			Format:         "rfc3164",
			Maxlen:         &maxlen,
			Size:           &size,
			TimeoutConnect: &timeoutConnect,
			TimeoutServer:  &timeoutServer,
		},
		Servers: map[string]models.Server{
			"webserv": {
				Address: "192.168.1.1",
				Name:    "webserv",
				ID:      misc.Ptr[int64](1234),
				Port:    misc.Ptr[int64](9200),
			},
		},
	}
	err = clientTest.CreateStructuredRing(r, "", version)
	require.NoError(t, err)
	version++

	v, ring, err := clientTest.GetStructuredRing("created_ring", "")
	require.NoError(t, err)
	require.True(t, ring.Equal(*r), "ring=%s - diff %v", ring.Name, cmp.Diff(*ring, *r))
	require.Equal(t, version, v, "Version %v returned, expected %v", v, version)

	err = clientTest.CreateStructuredRing(r, "", version)
	require.Error(t, err, "should throw error ring already exists")

	err = clientTest.DeleteRing("created_ring", "", version)
	require.NoError(t, err)
	version++

	v, _ = clientTest.GetVersion("")
	require.Equal(t, version, v, "Version %v returned, expected %v", v, version)

	err = clientTest.DeleteRing("created_ring", "", 999999)
	require.Error(t, err, "Should throw error, non existent frontend")
	require.ErrorIs(t, err, configuration.ErrVersionMismatch, "Should throw configuration.ErrVersionMismatch error")

	_, _, err = clientTest.GetStructuredRing("created_ring", "")
	require.Error(t, err, "deleteRing failed, ring created_ring still exists")

	require.Error(t, err, "should throw error, non existent ring")
}
