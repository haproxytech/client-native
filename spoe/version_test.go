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

package spoe

import (
	"path"
	"path/filepath"
	"testing"

	"github.com/haproxytech/client-native/v4/misc"
)

func TestSingleSpoe_GetConfigurationVersion(t *testing.T) {
	dir, configFile, err := misc.CreateTempDir(basicConfig, true)
	if err != nil {
		t.Error(err.Error())
	}
	transactionDir, _, err := misc.CreateTempDir("", false)
	if err != nil {
		t.Error(err.Error())
	}
	defer func() {
		_ = remove(path.Join(dir, configFile))
		_ = remove(dir)
		_ = remove(transactionDir)
	}()

	tests := []struct {
		params        Params
		name          string
		transactionID string
		want          int64
		wantErr       bool
	}{
		{
			name: "Should return correct version",
			params: Params{
				SpoeDir:           dir,
				TransactionDir:    transactionDir,
				ConfigurationFile: filepath.Join(dir, configFile),
			},
			transactionID: "",
			want:          1,
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ss, err := newSingleSpoe(tt.params)
			if err != nil {
				t.Errorf("SingleSpoe.GetConfigurationVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got, err := ss.GetConfigurationVersion(tt.transactionID)
			if (err != nil) != tt.wantErr {
				t.Errorf("SingleSpoe.GetConfigurationVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SingleSpoe.GetConfigurationVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
