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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/haproxytech/client-native/v2/misc"
	"github.com/haproxytech/client-native/v2/models"
)

func TestSingleSpoe_GetGroups(t *testing.T) {
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
	groupName := "mygroup"
	tests := []struct {
		name    string
		params  Params
		scope   string
		want    int64
		want1   models.SpoeGroups
		wantErr bool
	}{
		{
			name: "Should return all messages",
			params: Params{
				SpoeDir:           dir,
				TransactionDir:    transactionDir,
				ConfigurationFile: filepath.Join(dir, configFile),
			},
			scope: "[ip-reputation]",
			want:  1,
			want1: models.SpoeGroups{&models.SpoeGroup{
				Messages: "mymessage",
				Name:     &groupName,
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ss, err := newSingleSpoe(tt.params)
			if err != nil {
				t.Errorf("SingleSpoe.GetGroups() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got, got1, err := ss.GetGroups(tt.scope, "")
			if (err != nil) != tt.wantErr {
				t.Errorf("SingleSpoe.GetGroups() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SingleSpoe.GetGroups() got = %v, want %v", got, tt.want)
			}
			if !assert.EqualValues(t, got1, tt.want1) {
				t.Errorf("SingleSpoe.GetGroups() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestSingleSpoe_DeleteGroup(t *testing.T) { //nolint:dupl
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
		scope     string
		groupName string
		version   int64
		wantErr   bool
	}{
		{
			name: "Should delete a group",
			params: Params{
				SpoeDir:           dir,
				TransactionDir:    transactionDir,
				ConfigurationFile: filepath.Join(dir, configFile),
			},
			scope:     "[ip-reputation]",
			groupName: "mygroup",
			version:   1,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ss, err := newSingleSpoe(tt.params)
			if err != nil {
				t.Errorf("SingleSpoe.DeleteGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err := ss.DeleteGroup(tt.scope, tt.groupName, "", tt.version); (err != nil) != tt.wantErr {
				t.Errorf("SingleSpoe.DeleteGroup() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSingleSpoe_CreateGroup(t *testing.T) {
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
	groupName := "group2"
	tests := []struct {
		name    string
		params  Params
		scope   string
		data    *models.SpoeGroup
		version int64
		wantErr bool
	}{
		{
			name: "Should create a new group",
			params: Params{
				SpoeDir:           dir,
				TransactionDir:    transactionDir,
				ConfigurationFile: filepath.Join(dir, configFile),
			},
			scope: "[ip-reputation]",
			data: &models.SpoeGroup{
				Name:     &groupName,
				Messages: "msg2",
			},
			version: 1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ss, err := newSingleSpoe(tt.params)
			if err != nil {
				t.Errorf("SingleSpoe.CreateGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err := ss.CreateGroup(tt.scope, tt.data, "", tt.version); (err != nil) != tt.wantErr {
				t.Errorf("SingleSpoe.CreateGroup() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSingleSpoe_EditGroup(t *testing.T) {
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
	groupName := "mygroup"
	editGroup := &models.SpoeGroup{
		Name:     &groupName,
		Messages: "msg4",
	}
	tests := []struct {
		name      string
		params    Params
		scope     string
		data      *models.SpoeGroup
		groupName string
		version   int64
		wantErr   bool
	}{
		{
			name: "Should edit an existing group",
			params: Params{
				SpoeDir:           dir,
				TransactionDir:    transactionDir,
				ConfigurationFile: filepath.Join(dir, configFile),
			},
			scope:     "[ip-reputation]",
			groupName: groupName,
			data:      editGroup,
			version:   1,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ss, err := newSingleSpoe(tt.params)
			if err != nil {
				t.Errorf("SingleSpoe.EditGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err = ss.EditGroup(tt.scope, tt.data, tt.groupName, "", tt.version); (err != nil) != tt.wantErr {
				t.Errorf("SingleSpoe.EditGroup() error = %v, wantErr %v", err, tt.wantErr)
			}
			_, got, err := ss.GetGroup(tt.scope, tt.groupName, "")
			if err != nil {
				t.Errorf("SingleSpoe.EditGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.EqualValues(t, got, editGroup) {
				t.Errorf("SingleSpoe.EditGroup() got = %v, want %v", got, editGroup)
			}
		})
	}
}
