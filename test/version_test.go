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
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func generateConfig(config string) (string, error) {
	f, err := ioutil.TempFile("/tmp", "version")
	if err != nil {
		return "", err
	}
	err = prepareTestFile(config, f.Name())
	if err != nil {
		return "", err
	}
	return f.Name(), nil
}

func TestClient_PostRawConfiguration(t *testing.T) {
	fVersion, err := generateConfig("")
	require.NoError(t, err, "generateConfig")

	t.Cleanup(func() {
		assert.NoError(t, deleteTestFile(fVersion), "clean-up")
	})

	c, err := prepareClient(fVersion)
	require.NoError(t, err, "prepareClient")

	v, vErr := c.GetVersion("")
	assert.Equal(t, int64(1), v, "initialized configuration must be 1")
	// The user is providing a raw configuration with a wrong version such as metadata:
	// this must be ignored and removed by Client Native
	configWithVersion := `# _version=123
global
	daemon`

	err = c.PostRawConfiguration(&configWithVersion, 1, false)
	require.NoError(t, err, "PostRawConfiguration")

	v, vErr = c.GetVersion("")
	require.NoError(t, vErr, "GetVersion")

	assert.Equal(t, int64(2), v, "123 should be dropped, and version bumped")
}

func TestClient_GetConfigurationVersion(t *testing.T) {
	configWithVersion := `# _version=10
global
	daemon
`
	fVersion, err := generateConfig(configWithVersion)
	if err != nil {
		t.Error(err.Error())
	}
	defer func() {
		_ = deleteTestFile(fVersion)
	}()

	configWithoutVersion := `
global
	daemon
`
	fNoVersion, err := generateConfig(configWithoutVersion)
	if err != nil {
		t.Error(err.Error())
	}
	defer func() {
		_ = deleteTestFile(fNoVersion)
	}()

	tests := []struct {
		name              string
		configurationFile string
		want              int64
		wantErr           bool
	}{
		{
			name:              "Pass with version",
			configurationFile: fVersion,
			want:              10,
			wantErr:           false,
		},
		{
			name:              "Pass without version",
			configurationFile: fNoVersion,
			want:              1, // config without version should add `# _version=1`
			wantErr:           false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := prepareClient(tt.configurationFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("prepareClient error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got, err := c.GetConfigurationVersion("")
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetConfigurationVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Client.GetConfigurationVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
