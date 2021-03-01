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
	"path/filepath"
	"reflect"
	"testing"

	"github.com/haproxytech/client-native/v2/misc"
	"github.com/haproxytech/client-native/v2/models"
)

func TestSingleSpoe_GetScopes(t *testing.T) {
	dir, configFile, err := misc.CreateTempDir(basicConfig, true)
	if err != nil {
		t.Error(err.Error())
	}
	transactionDir, _, err := misc.CreateTempDir("", false)
	if err != nil {
		t.Error(err.Error())
	}
	defer func() {
		_ = remove(configFile)
		_ = remove(dir)
		_ = remove(transactionDir)
	}()
	tests := []struct {
		name          string
		params        Params
		transactionID string
		want          int64
		want1         models.SpoeScopes
		wantErr       bool
	}{
		{
			name: "Should return scopes",
			params: Params{
				SpoeDir:           dir,
				TransactionDir:    transactionDir,
				ConfigurationFile: filepath.Join(dir, configFile),
			},
			want:    1,
			want1:   models.SpoeScopes{models.SpoeScope("[ip-reputation]")},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ss, err := newSingleSpoe(tt.params)
			if err != nil {
				t.Errorf("SingleSpoe.GetScopes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			got, got1, err := ss.GetScopes(tt.transactionID)
			if (err != nil) != tt.wantErr {
				t.Errorf("SingleSpoe.GetScopes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SingleSpoe.GetScopes() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("SingleSpoe.GetScopes() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestSingleSpoe_DeleteScope(t *testing.T) {
	dir, configFile, err := misc.CreateTempDir(basicConfig, true)
	if err != nil {
		t.Error(err.Error())
	}
	transactionDir, _, err := misc.CreateTempDir("", false)
	if err != nil {
		t.Error(err.Error())
	}
	defer func() {
		_ = remove(configFile)
		_ = remove(dir)
		_ = remove(transactionDir)
	}()
	tests := []struct {
		name      string
		params    Params
		scopeName string
		version   int64
		wantErr   bool
	}{
		{
			name: "Should delete scope by its name and provided version",
			params: Params{
				SpoeDir:           dir,
				TransactionDir:    transactionDir,
				ConfigurationFile: filepath.Join(dir, configFile),
			},
			scopeName: "[ip-reputation]",
			version:   1,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ss, err := newSingleSpoe(tt.params)
			if err != nil {
				t.Errorf("SingleSpoe.DeleteScope() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err := ss.DeleteScope(tt.scopeName, "", tt.version); (err != nil) != tt.wantErr {
				t.Errorf("SingleSpoe.DeleteScope() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSingleSpoe_CreateScope(t *testing.T) {
	dir, configFile, err := misc.CreateTempDir(basicConfig, true)
	if err != nil {
		t.Error(err.Error())
	}
	transactionDir, _, err := misc.CreateTempDir("", false)
	if err != nil {
		t.Error(err.Error())
	}

	defer func() {
		_ = remove(configFile)
		_ = remove(dir)
		_ = remove(transactionDir)
	}()
	tests := []struct {
		name    string
		params  Params
		data    models.SpoeScope
		version int64
		wantErr bool
	}{
		{
			name: "Should create a new scope",
			params: Params{
				SpoeDir:           dir,
				TransactionDir:    transactionDir,
				ConfigurationFile: filepath.Join(dir, configFile),
			},
			data:    models.SpoeScope("[new-scope]"),
			version: 1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ss, err := newSingleSpoe(tt.params)
			if err != nil {
				t.Errorf("SingleSpoe.CreateScope() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err := ss.CreateScope(&tt.data, "", tt.version); (err != nil) != tt.wantErr {
				t.Errorf("SingleSpoe.CreateScope() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSingleSpoe_GetScope(t *testing.T) {
	dir, configFile, err := misc.CreateTempDir(basicConfig, true)
	if err != nil {
		t.Error(err.Error())
	}
	transactionDir, _, err := misc.CreateTempDir("", false)
	if err != nil {
		t.Error(err.Error())
	}
	defer func() {
		_ = remove(configFile)
		_ = remove(dir)
		_ = remove(transactionDir)
	}()

	tests := []struct {
		name      string
		params    Params
		scopeName string
		version   int64
		want1     models.SpoeScope
		wantErr   bool
	}{
		{
			name: "Should return scope if exists",
			params: Params{
				SpoeDir:           dir,
				TransactionDir:    transactionDir,
				ConfigurationFile: filepath.Join(dir, configFile),
			},
			scopeName: "[ip-reputation]",
			version:   1,
			want1:     models.SpoeScope("[ip-reputation]"),
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ss, err := newSingleSpoe(tt.params)
			if err != nil {
				t.Errorf("SingleSpoe.GetScope() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got, got1, err := ss.GetScope(tt.scopeName, "")
			if (err != nil) != tt.wantErr {
				t.Errorf("SingleSpoe.GetScope() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.version {
				t.Errorf("SingleSpoe.GetScope() got = %v, want %v", got, tt.version)
			}
			if !reflect.DeepEqual(got1, &tt.want1) {
				t.Errorf("SingleSpoe.GetScope() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
