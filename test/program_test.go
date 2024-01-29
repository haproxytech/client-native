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

package test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/haproxytech/client-native/v6/misc"
	"github.com/haproxytech/client-native/v6/models"
)

func programExpectation() map[string]models.Programs {
	initStructuredExpected()
	res := StructuredToProgramMap()
	// Add individual entries
	for _, vs := range res {
		for _, v := range vs {
			key := v.Name
			res[key] = models.Programs{v}
		}
	}
	return res
}

func TestGetPrograms(t *testing.T) {
	m := make(map[string]models.Programs)
	v, programs, err := clientTest.GetPrograms("")

	require.NoError(t, err)
	require.Len(t, programs, 2)
	require.Equal(t, version, v)
	m[""] = programs
	checkPrograms(t, m)
}

func TestGetProgram(t *testing.T) {
	m := make(map[string]models.Programs)
	v, program, err := clientTest.GetProgram("test", "")

	require.NoError(t, err)
	require.NotNil(t, program)
	require.Equal(t, version, v)
	m["test"] = models.Programs{program}

	checkPrograms(t, m)
}

func checkPrograms(t *testing.T, got map[string]models.Programs) {
	exp := programExpectation()
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

func TestCreateEditDeleteProgram(t *testing.T) {
	// Creating a program
	program := &models.Program{
		Command:       misc.StringP("python script.py"),
		Group:         "xfs",
		Name:          "my-program",
		StartOnReload: "enabled",
		User:          "www-data",
	}
	require.NoError(t, clientTest.CreateProgram(program, "", version))
	version++
	// Ensuring data is correct
	configVersion, found, err := clientTest.GetProgram(program.Name, "")
	assert.Nil(t, err)
	assert.Equal(t, configVersion, version)
	assert.EqualValues(t, program, found)
	// Expecting failure due to duplicated entry
	assert.Error(t, clientTest.CreateProgram(program, "", version))
	// Updating and ensuring data is properly reflected
	program.Command = misc.StringP("sh script.sh")
	program.Group = ""
	program.User = ""
	program.StartOnReload = "disabled"

	require.NoError(t, clientTest.EditProgram(program.Name, program, "", version))
	version++

	configVersion, found, err = clientTest.GetProgram(program.Name, "")
	assert.Nil(t, err)
	assert.Equal(t, configVersion, version)
	assert.EqualValues(t, program, found)
	// Deleting and ensuring triggering not found error
	assert.NoError(t, clientTest.DeleteProgram(program.Name, "", version))
	version++

	_, _, err = clientTest.GetProgram(program.Name, "")
	assert.Error(t, err)
}
