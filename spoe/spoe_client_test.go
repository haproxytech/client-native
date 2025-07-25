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
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func remove(file string) error {
	err := os.Remove(file)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}
	return nil
}

func Test_spoeclient_Create(t *testing.T) {
	fileName := "spoe.cfg"
	spoeDir := "/tmp"
	defer func() {
		_ = remove(filepath.Join(spoeDir, fileName))
	}()

	tests := []struct {
		name       string
		clients    map[string]*SingleSpoe
		initParams Params
		fileName   string
		readCloser io.ReadCloser
		want       string
		wantErr    bool
	}{
		{
			name:       "Should create a file with # _version prepended",
			clients:    make(map[string]*SingleSpoe),
			fileName:   "spoe.cfg",
			readCloser: io.NopCloser(bytes.NewReader([]byte("hello world"))),
			initParams: Params{
				SpoeDir:        "/tmp",
				TransactionDir: "/tmp",
			},
			want:    filepath.Join("/tmp", "spoe.cfg"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := spoeclient{
				clients:    tt.clients,
				initParams: tt.initParams,
			}
			got, err := c.Create(tt.fileName, tt.readCloser)
			if (err != nil) != tt.wantErr {
				t.Errorf("spoeclient.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("spoeclient.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}
