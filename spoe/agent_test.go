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

func TestSingleSpoe_GetAgents(t *testing.T) {
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
	agent := "iprep-agent"
	tests := []struct {
		name    string
		params  Params
		scope   string
		want    int64
		want1   models.SpoeAgents
		wantErr bool
	}{
		{
			name: "Should return all configured agents",
			params: Params{
				SpoeDir:           dir,
				TransactionDir:    transactionDir,
				ConfigurationFile: filepath.Join(dir, configFile),
			},
			scope: "[ip-reputation]",
			want:  1,
			want1: models.SpoeAgents{&models.SpoeAgent{
				Messages:          "check-client-ip",
				OptionVarPrefix:   "iprep",
				HelloTimeout:      *misc.ParseTimeout("2s"),
				IdleTimeout:       *misc.ParseTimeout("2m"),
				ProcessingTimeout: *misc.ParseTimeout("10ms"),
				UseBackend:        "agents",
				Log:               models.LogTargets{&models.LogTarget{Global: true, Index: misc.Int64P(0)}},
				Async:             "enabled",
				Name:              &agent,
				ContinueOnError:   "enabled",
				ForceSetVar:       "enabled",
				DontlogNormal:     "enabled",
				Pipelining:        "enabled",
				SendFragPayload:   "enabled",
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ss, err := newSingleSpoe(tt.params)
			if err != nil {
				t.Errorf("SingleSpoe.GetAgents() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got, got1, err := ss.GetAgents(tt.scope, "")
			if (err != nil) != tt.wantErr {
				t.Errorf("SingleSpoe.GetAgents() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SingleSpoe.GetAgents() got = %v, want %v", got, tt.want)
			}
			if !assert.EqualValues(t, got1, tt.want1) {
				t.Errorf("SingleSpoe.GetAgents() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestSingleSpoe_DeleteAgent(t *testing.T) { //nolint:dupl
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
		agentName string
		version   int64
		wantErr   bool
	}{
		{
			name: "Should delete an agent",
			params: Params{
				SpoeDir:           dir,
				TransactionDir:    transactionDir,
				ConfigurationFile: filepath.Join(dir, configFile),
			},
			scope:     "[ip-reputation]",
			agentName: "iprep-agent",
			version:   1,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ss, err := newSingleSpoe(tt.params)
			if err != nil {
				t.Errorf("SingleSpoe.DeleteAgent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err := ss.DeleteAgent(tt.scope, tt.agentName, "", tt.version); (err != nil) != tt.wantErr {
				t.Errorf("SingleSpoe.DeleteAgent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSingleSpoe_CreateAgent(t *testing.T) {
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
	agentName := "new-agent"
	tests := []struct {
		name    string
		params  Params
		data    *models.SpoeAgent
		scope   string
		version int64
		wantErr bool
	}{
		{
			name: "Should create a new agent",
			params: Params{
				SpoeDir:           dir,
				TransactionDir:    transactionDir,
				ConfigurationFile: filepath.Join(dir, configFile),
			},
			data: &models.SpoeAgent{
				Async:             "enabled",
				Name:              &agentName,
				UseBackend:        "mybackend",
				HelloTimeout:      2000,
				ContinueOnError:   "enabled",
				DontlogNormal:     "disabled",
				EngineName:        "myengine",
				ForceSetVar:       "enabled",
				IdleTimeout:       10000,
				Log:               models.LogTargets{&models.LogTarget{Nolog: true, Index: misc.Int64P(1)}},
				Maxconnrate:       3000,
				MaxFrameSize:      10,
				Maxerrrate:        11,
				Pipelining:        "disabled",
				ProcessingTimeout: 200,
				RegisterVarNames:  "names",
				SendFragPayload:   "disabled",
				OptionVarPrefix:   "pref",
			},
			scope:   "[ip-reputation]",
			version: 1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ss, err := newSingleSpoe(tt.params)
			if err != nil {
				t.Errorf("SingleSpoe.CreateAgent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err := ss.CreateAgent(tt.scope, tt.data, "", tt.version); (err != nil) != tt.wantErr {
				t.Errorf("SingleSpoe.CreateAgent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSingleSpoe_EditAgent(t *testing.T) {
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
	agentName := "iprep-agent"
	editAgent := &models.SpoeAgent{
		Async:                "enabled",
		Name:                 &agentName,
		UseBackend:           "newbackend",
		HelloTimeout:         1000,
		ContinueOnError:      "enabled", // single
		DontlogNormal:        "disabled",
		ForceSetVar:          "enabled", // single
		Groups:               "group1",
		IdleTimeout:          30000,
		Log:                  models.LogTargets{&models.LogTarget{Nolog: true, Index: misc.Int64P(0)}},
		Maxconnrate:          5000,
		MaxFrameSize:         15,
		Maxerrrate:           25,
		MaxWaitingFrames:     100,
		Messages:             "msg1",
		Pipelining:           "enabled",
		ProcessingTimeout:    500,
		RegisterVarNames:     "new-names",
		SendFragPayload:      "disabled",
		OptionVarPrefix:      "pref",
		OptionSetOnError:     "setonerror",
		OptionSetProcessTime: "setprocesstime",
		OptionSetTotalTime:   "settotaltime",
	}
	tests := []struct {
		name    string
		params  Params
		data    *models.SpoeAgent
		scope   string
		version int64
		wantErr bool
	}{
		{
			name: "Should edit an existing agent",
			params: Params{
				SpoeDir:           dir,
				TransactionDir:    transactionDir,
				ConfigurationFile: filepath.Join(dir, configFile),
			},
			data:    editAgent,
			scope:   "[ip-reputation]",
			version: 1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ss, err := newSingleSpoe(tt.params)
			if err != nil {
				t.Errorf("SingleSpoe.EditAgent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err = ss.EditAgent(tt.scope, tt.data, "", tt.version); (err != nil) != tt.wantErr {
				t.Errorf("SingleSpoe.EditAgent() error = %v, wantErr %v", err, tt.wantErr)
			}
			_, got, err := ss.GetAgent(tt.scope, agentName, "")
			if err != nil {
				t.Errorf("SingleSpoe.EditAgent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.EqualValues(t, got, editAgent) {
				t.Errorf("SingleSpoe.EditAgent() got = %v, want %v", got, editAgent)
			}
		})
	}
}
