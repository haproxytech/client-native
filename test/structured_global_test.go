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
	_ "embed"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/haproxytech/client-native/v6/models"
	"github.com/stretchr/testify/require"
)

func TestStructuredGetGlobal(t *testing.T) {
	clientTest, filename, err := getTestClient()
	require.NoError(t, err)
	defer os.Remove(filename)
	version := int64(1)

	v, global, err := clientTest.GetStructuredGlobalConfiguration("")
	require.NoError(t, err)
	require.Equal(t, version, v, "Version %v returned, expected %v", v, version)

	checkStructuredGlobal(t, global)
}

func checkStructuredGlobal(t *testing.T, global *models.Global) {
	want := globalExpcectation()
	require.True(t, global.Equal(*want), "diff %v", cmp.Diff(*global, *want))
}

func TestPutStructuredGlobal(t *testing.T) {
	clientTest, filename, err := getTestClient()
	require.NoError(t, err)
	defer os.Remove(filename)
	version := int64(1)

	g := &models.Global{
		GlobalBase: getGlobalBase(),
		LogTargetList: models.LogTargets{
			&models.LogTarget{
				Address:  "192.169.0.1",
				Facility: "mail",
				Global:   true,
			},
		},
	}

	err = clientTest.PushStructuredGlobalConfiguration(g, "", version)

	require.NoError(t, err)
	version++

	ver, global, err := clientTest.GetStructuredGlobalConfiguration("")
	require.NoError(t, err)
	require.True(t, global.Equal(*g), "global - diff %v", cmp.Diff(*global, *g))
	require.Equal(t, version, ver, "Version %v returned, expected %v", ver, version)

	err = clientTest.PushStructuredGlobalConfiguration(g, "", 55)
	require.Error(t, err, "Should have returned version conflict.")

}
