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

package runtime

import (
	"testing"

	"github.com/haproxytech/client-native/v6/models"
	"github.com/stretchr/testify/assert"
)

func Test_parseMapPayload(t *testing.T) {
	type args struct {
		entries    models.MapEntries
		maxBufSize int
	}
	tests := []struct {
		name             string
		args             args
		wantExceededSize bool
		wantPayload      []string
	}{
		{
			name: "Return false and non empty slice when size < maxBufSize",
			args: args{entries: models.MapEntries{
				&models.MapEntry{Key: "k", Value: "v"},
			}, maxBufSize: 10},
			wantExceededSize: false,
			wantPayload:      []string{"k v\n"},
		},
		{
			name: "Return false and non empty slice when size == maxBufSize",
			args: args{entries: models.MapEntries{
				&models.MapEntry{Key: "key1", Value: "val1"},
				&models.MapEntry{Key: "key2", Value: "val2"},
			}, maxBufSize: 10},
			wantExceededSize: false,
			wantPayload: []string{
				"key1 val1\n",
				"key2 val2\n",
			},
		},
		{
			name: "Return true and non empty slice when size > maxBufSize",
			args: args{entries: models.MapEntries{
				&models.MapEntry{Key: "key1", Value: "val1"},
				&models.MapEntry{Key: "key2", Value: "val2"},
			}, maxBufSize: 5},
			wantExceededSize: true,
			wantPayload: []string{
				"key1 val1\n",
				"key2 val2\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotExceededSize, gotPayload := parseMapPayload(tt.args.entries, tt.args.maxBufSize)
			if gotExceededSize != tt.wantExceededSize {
				t.Errorf("parseMapPayload() gotExceededSize = %v, want %v", gotExceededSize, tt.wantExceededSize)
			}
			if !assert.EqualValues(t, gotPayload, tt.wantPayload) {
				t.Errorf("parseMapPayload() gotPayload gotPayload = %v, want %v", gotPayload, tt.wantPayload)
			}
		})
	}
}
