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

func TestSingleSpoe_GetMessages(t *testing.T) {
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

	eventName := "on-client-session"
	messageName := "check-client-ip"
	tests := []struct {
		name    string
		params  Params
		scope   string
		want    int64
		want1   models.SpoeMessages
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
			want1: models.SpoeMessages{
				&models.SpoeMessage{
					Args: "ip=src",
					Event: &models.SpoeMessageEvent{
						Cond:     "if",
						CondTest: "! { src -f /etc/haproxy/whitelist.lst }",
						Name:     &eventName,
					},
					Name: &messageName,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ss, err := newSingleSpoe(tt.params)
			if err != nil {
				t.Errorf("SingleSpoe.GetMessages() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got, got1, err := ss.GetMessages(tt.scope, "")
			if (err != nil) != tt.wantErr {
				t.Errorf("SingleSpoe.GetMessages() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SingleSpoe.GetMessages() got = %v, want %v", got, tt.want)
			}
			if !assert.EqualValues(t, got1, tt.want1) {
				t.Errorf("SingleSpoe.GetMessages() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestSingleSpoe_DeleteMessage(t *testing.T) { //nolint:dupl
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
		name        string
		params      Params
		scope       string
		messageName string
		version     int64
		wantErr     bool
	}{
		{
			name: "Should delete message",
			params: Params{
				SpoeDir:           dir,
				TransactionDir:    transactionDir,
				ConfigurationFile: filepath.Join(dir, configFile),
			},
			scope:       "[ip-reputation]",
			messageName: "check-client-ip",
			version:     1,
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ss, err := newSingleSpoe(tt.params)
			if err != nil {
				t.Errorf("SingleSpoe.DeleteMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err := ss.DeleteMessage(tt.scope, tt.messageName, "", tt.version); (err != nil) != tt.wantErr {
				t.Errorf("SingleSpoe.DeleteMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSingleSpoe_CreateMessage(t *testing.T) {
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
	messageName := "new-message"
	eventName := "on-server-session"
	tests := []struct {
		name    string
		params  Params
		data    *models.SpoeMessage
		scope   string
		version int64
		wantErr bool
	}{
		{
			name: "Shoud create a new message section",
			params: Params{
				SpoeDir:           dir,
				TransactionDir:    transactionDir,
				ConfigurationFile: filepath.Join(dir, configFile),
			},
			scope: "[ip-reputation]",
			data: &models.SpoeMessage{
				Args: "ip=dst",
				Event: &models.SpoeMessageEvent{
					Cond:     "unless",
					CondTest: "{ src -f /etc/haproxy/whitelist.lst }",
					Name:     &eventName,
				},
				Name: &messageName,
			},
			version: 1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ss, err := newSingleSpoe(tt.params)
			if err != nil {
				t.Errorf("SingleSpoe.CreateMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err := ss.CreateMessage(tt.scope, tt.data, "", tt.version); (err != nil) != tt.wantErr {
				t.Errorf("SingleSpoe.CreateMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSingleSpoe_EditMessage(t *testing.T) {
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
	messageName := "check-client-ip"
	eventName := "on-backend-tcp-request"
	editMessage := &models.SpoeMessage{
		Args: "ip=dst",
		Event: &models.SpoeMessageEvent{
			Cond:     "unless",
			CondTest: "{ src -f /etc/haproxy/whitelist.lst }",
			Name:     &eventName,
		},
		Name: &messageName,
	}
	tests := []struct {
		name        string
		params      Params
		scope       string
		data        *models.SpoeMessage
		messageName string
		version     int64
		wantErr     bool
	}{
		{
			name: "Should edit an existing message",
			params: Params{
				SpoeDir:           dir,
				TransactionDir:    transactionDir,
				ConfigurationFile: filepath.Join(dir, configFile),
			},
			scope:       "[ip-reputation]",
			data:        editMessage,
			version:     1,
			wantErr:     false,
			messageName: messageName,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ss, err := newSingleSpoe(tt.params)
			if err != nil {
				t.Errorf("SingleSpoe.EditMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err = ss.EditMessage(tt.scope, tt.data, tt.messageName, "", tt.version); (err != nil) != tt.wantErr {
				t.Errorf("SingleSpoe.EditMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
			_, got, err := ss.GetMessage(tt.scope, tt.messageName, "")
			if err != nil {
				t.Errorf("SingleSpoe.EditMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.EqualValues(t, got, editMessage) {
				t.Errorf("SingleSpoe.EditMessage() got = %v, want %v", got, editMessage)
			}
		})
	}
}
