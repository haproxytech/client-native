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

func TestGetStructuredResolvers(t *testing.T) {
	clientTest, filename, err := getTestClient()
	require.NoError(t, err)
	defer os.Remove(filename)
	version := int64(1)

	m := make(map[string]models.Resolvers)

	v, resolvers, err := clientTest.GetStructuredResolvers("")
	require.NoError(t, err)

	require.Equal(t, version, v, "Version %v returned, expected %v", v, version)
	m[""] = resolvers
	checkStructuredResolvers(t, m)
}

func TestGetStructuredResolver(t *testing.T) {
	clientTest, filename, err := getTestClient()
	require.NoError(t, err)
	defer os.Remove(filename)
	version := int64(1)

	m := make(map[string]models.Resolvers)
	v, l, err := clientTest.GetStructuredResolver("test", "")
	require.NoError(t, err)
	require.Equal(t, version, v, "Version %v returned, expected %v", v, version)

	m["test"] = models.Resolvers{l}
	checkStructuredResolvers(t, m)

	_, _, err = clientTest.GetStructuredResolver("doesnotexist", "")
	require.Error(t, err, "Should throw error, non existent resolvers section")
}

func checkStructuredResolvers(t *testing.T, got map[string]models.Resolvers) {
	exp := resolverExpectation()
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

func TestCreateEditDeleteStructuredResolver(t *testing.T) {
	clientTest, filename, err := getTestClient()
	require.NoError(t, err)
	defer os.Remove(filename)
	version := int64(1)

	r := &models.Resolver{
		ResolverBase: models.ResolverBase{
			Name:                "created_resolver",
			AcceptedPayloadSize: 4096,
			HoldNx:              misc.Int64P(10),
			HoldObsolete:        misc.Int64P(10),
			HoldOther:           misc.Int64P(10),
			HoldRefused:         misc.Int64P(10),
			HoldTimeout:         misc.Int64P(10),
			HoldValid:           misc.Int64P(100),
			ResolveRetries:      10,
			ParseResolvConf:     true,
			TimeoutResolve:      10,
			TimeoutRetry:        10,
		},
		Nameservers: map[string]models.Nameserver{
			"my-nameserver": {
				Address: misc.Ptr("192.168.1.1"),
				Name:    "my-nameserver",
				Port:    misc.Ptr[int64](9200),
			},
		},
	}
	err = clientTest.CreateStructuredResolver(r, "", version)
	require.NoError(t, err)
	version++

	v, resolver, err := clientTest.GetStructuredResolver("created_resolver", "")
	require.NoError(t, err)
	require.True(t, resolver.Equal(*r), "resolver=%s - diff %v", resolver.Name, cmp.Diff(*resolver, *r))
	require.Equal(t, version, v, "Version %v returned, expected %v", v, version)

	err = clientTest.CreateStructuredResolver(r, "", version)
	require.Error(t, err, "Should throw error resolver already exists")

	err = clientTest.DeleteResolver("created_resolver", "", version)
	require.NoError(t, err)
	version++

	v, _ = clientTest.GetVersion("")
	require.Equal(t, version, v, "Version %v returned, expected %v", v, version)

	err = clientTest.DeleteResolver("created_resolver", "", 999999)
	require.Error(t, err, "Should throw error, non existent frontend")
	require.ErrorIs(t, err, configuration.ErrVersionMismatch, "Should throw configuration.ErrVersionMismatch error")

	_, _, err = clientTest.GetStructuredResolver("created_resolver", "")
	require.Error(t, err, "DeleteResolver failed, resolver created_resolver still exists")

	err = clientTest.DeleteResolver("doesnotexist", "", version)
	require.Error(t, err, "Should throw error, non existent resolver")
}
