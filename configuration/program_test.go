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

package configuration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/haproxytech/client-native/v4/misc"
	"github.com/haproxytech/client-native/v4/models"
)

func TestGetPrograms(t *testing.T) {
	v, programs, err := clientTest.GetPrograms("")

	require.NoError(t, err)
	require.Len(t, programs, 2)
	require.Equal(t, version, v)
}

func TestGetProgram(t *testing.T) {
	v, program, err := clientTest.GetProgram("test", "")

	require.NoError(t, err)
	require.NotNil(t, program)
	require.Equal(t, version, v)

	assert.Equal(t, `echo "Hello, World!"`, *program.Command)
	assert.Equal(t, "hapee-lb", program.User)
	assert.Equal(t, "hapee", program.Group)
	assert.Equal(t, models.ProgramStartOnReloadEnabled, program.StartOnReload)
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
