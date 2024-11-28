// Copyright 2024 HAProxy Technologies
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
	"github.com/haproxytech/client-native/v6/models"
	"github.com/stretchr/testify/require"
)

func TestStructuredGetTest(t *testing.T) {
	clientTest, filename, err := getTestClient()
	require.NoError(t, err)
	defer os.Remove(filename)
	version := int64(1)

	v, traces, err := clientTest.GetStructuredTraces("")
	require.NoError(t, err)
	require.Equal(t, version, v, "Version %v returned, expected %v", v, version)

	checkStructuredTraces(t, traces)
}

func checkStructuredTraces(t *testing.T, traces *models.Traces) {
	want := tracesExpcectation()
	require.True(t, traces.Equal(*want), "diff %v", cmp.Diff(*traces, *want))
}

func TestPutStructuredTraces(t *testing.T) {
	require := require.New(t)
	clientTest, filename, err := getTestClient()
	require.NoError(err)
	defer os.Remove(filename)
	version := int64(1)

	traces := &models.Traces{
		Entries: models.TraceEntries{
			&models.TraceEntry{Trace: "test trace 1"},
			&models.TraceEntry{Trace: "test trace 2"},
		},
	}

	err = clientTest.PushStructuredTraces(traces, "", version)

	require.NoError(err)
	version++

	ver, got, err := clientTest.GetStructuredTraces("")
	require.NoError(err)
	require.True(got.Equal(*traces), "global - diff %v", cmp.Diff(*got, *traces))
	require.Equal(version, ver, "Version %v returned, expected %v", ver, version)

	err = clientTest.PushStructuredTraces(traces, "", 55)
	require.Error(err, "Should have returned version conflict.")

}
